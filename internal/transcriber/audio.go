// Voice Transcriber
// Copyright (c) 2025 Ihor Dvoretskyi
//
// Licensed under MIT License

package transcriber

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/idvoretskyi/voice-transcriber/pkg/config"
)

const (
	// ffmpegTimeout is the maximum time allowed for a single FFmpeg extraction.
	ffmpegTimeout = 30 * time.Minute

	// maxFileSize is the maximum accepted input file size (10 GB).
	maxFileSize = 10 * 1024 * 1024 * 1024
)

// InputType represents the kind of media file provided by the user.
type InputType int

const (
	// InputTypeAudio is a native audio file that Gemini can consume directly.
	InputTypeAudio InputType = iota
	// InputTypeVideo is a video file that requires FFmpeg audio extraction.
	InputTypeVideo
)

// audioExtensions lists file extensions that Gemini accepts as native audio.
// See: https://cloud.google.com/vertex-ai/generative-ai/docs/multimodal/audio-understanding
var audioExtensions = map[string]string{
	".wav":  "audio/wav",
	".mp3":  "audio/mp3",
	".flac": "audio/flac",
	".ogg":  "audio/ogg",
	".m4a":  "audio/m4a",
	".aac":  "audio/aac",
	".pcm":  "audio/pcm",
	".webm": "audio/webm",
}

// classifyInputFile determines whether the path is a native audio file or a
// video file that needs FFmpeg extraction. It returns the InputType and, for
// audio files, the MIME type string required by the Gemini API.
func classifyInputFile(inputPath string) (InputType, string) {
	ext := strings.ToLower(filepath.Ext(inputPath))

	if mimeType, ok := audioExtensions[ext]; ok {
		return InputTypeAudio, mimeType
	}

	return InputTypeVideo, ""
}

// PreparedAudio holds the audio bytes and MIME type ready to send to Gemini,
// plus a cleanup function to remove any temporary file that was created.
type PreparedAudio struct {
	Data     []byte
	MIMEType string
	Cleanup  func()
}

// prepareAudio reads the input file, extracting audio via FFmpeg when the
// input is a video, and returns a PreparedAudio ready to pass to the Gemini
// service.
func prepareAudio(ctx context.Context, inputPath string, cfg *config.Config) (*PreparedAudio, error) {
	cleanPath, err := validateAndSanitizeVideoPath(inputPath)
	if err != nil {
		return nil, err
	}

	inputType, mimeType := classifyInputFile(cleanPath)

	switch inputType {
	case InputTypeAudio:
		logVerbose(cfg, "Audio file detected (%s), skipping FFmpeg extraction", mimeType)

		data, err := os.ReadFile(cleanPath) // #nosec G304 -- cleanPath validated above
		if err != nil {
			return nil, fmt.Errorf("failed to read audio file: %w", err)
		}

		return &PreparedAudio{
			Data:     data,
			MIMEType: mimeType,
			Cleanup:  func() {},
		}, nil

	case InputTypeVideo:
		audioPath, err := extractAudio(ctx, cleanPath, cfg)
		if err != nil {
			return nil, err
		}

		data, err := os.ReadFile(audioPath) // #nosec G304 -- audioPath from extractAudio
		if err != nil {
			_ = os.Remove(audioPath)

			return nil, fmt.Errorf("failed to read extracted audio: %w", err)
		}

		return &PreparedAudio{
			Data:     data,
			MIMEType: "audio/wav",
			Cleanup: func() {
				if removeErr := os.Remove(audioPath); removeErr != nil && !os.IsNotExist(removeErr) {
					logVerbose(cfg, "Warning: failed to remove temp audio file: %v", removeErr)
				}
			},
		}, nil

	default:
		// unreachable unless InputType is extended
		return nil, fmt.Errorf("unsupported input type for %q", cleanPath)
	}
}

// extractAudio extracts audio from a video file using FFmpeg.
func extractAudio(ctx context.Context, videoPath string, cfg *config.Config) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, ffmpegTimeout)
	defer cancel()

	logVerbose(cfg, "Extracting audio from video: %s", videoPath)

	ffmpegPath, err := exec.LookPath("ffmpeg")
	if err != nil {
		return "", fmt.Errorf("ffmpeg not found; install ffmpeg first: %w", err)
	}

	audioPath := generateAudioPath(videoPath)

	if err := runFFmpegCommand(ctx, ffmpegPath, videoPath, audioPath, cfg); err != nil {
		return "", err
	}

	logVerbose(cfg, "Audio extracted to: %s", audioPath)

	return audioPath, nil
}

// validateAndSanitizeVideoPath validates and sanitizes any input media path.
func validateAndSanitizeVideoPath(inputPath string) (string, error) {
	inputPath = filepath.Clean(inputPath)

	fileInfo, err := os.Stat(inputPath)
	if err != nil {
		return "", fmt.Errorf("input file error: %w", err)
	}

	if !fileInfo.Mode().IsRegular() {
		return "", fmt.Errorf("not a regular file: %s", inputPath)
	}

	if fileInfo.Size() > maxFileSize {
		return "", fmt.Errorf("file too large (>10GB): %s", inputPath)
	}

	return inputPath, nil
}

// generateAudioPath creates a unique WAV file path in the system temp directory.
func generateAudioPath(inputPath string) string {
	timestamp := time.Now().UnixNano()
	baseFileName := strings.TrimSuffix(filepath.Base(inputPath), filepath.Ext(inputPath))

	return filepath.Join(os.TempDir(), fmt.Sprintf("%s_%d_audio.wav", baseFileName, timestamp))
}

// runFFmpegCommand executes FFmpeg and verifies the output was created.
func runFFmpegCommand(ctx context.Context, ffmpegPath, videoPath, audioPath string, cfg *config.Config) error {
	cmd := exec.CommandContext(ctx, ffmpegPath, // #nosec G204 -- ffmpegPath resolved via exec.LookPath
		"-i", videoPath,
		"-acodec", "pcm_s16le",
		"-ar", "16000",
		"-ac", "1",
		"-y",
		audioPath,
	)

	var stderr strings.Builder
	if cfg.Verbose {
		cmd.Stderr = io.MultiWriter(os.Stderr, &stderr)
	} else {
		cmd.Stderr = &stderr
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ffmpeg failed: %v, stderr: %s", err, stderr.String())
	}

	if _, err := os.Stat(audioPath); err != nil {
		return fmt.Errorf("ffmpeg did not create output file: %w", err)
	}

	return nil
}
