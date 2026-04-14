// Ukrainian Voice Transcriber
// Copyright (c) 2025 Ihor Dvoretskyi
//
// Licensed under MIT License

// Package transcriber provides audio transcription functionality.
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

	"github.com/idvoretskyi/ukrainian-voice-transcriber/pkg/config"
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
// video file that needs FFmpeg extraction.  It returns the InputType and, for
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
func prepareAudio(inputPath string, cfg *config.Config) (*PreparedAudio, error) {
	// Validate path first
	cleanPath, err := validateAndSanitizeVideoPath(inputPath)
	if err != nil {
		return nil, err
	}

	inputType, mimeType := classifyInputFile(cleanPath)

	switch inputType {
	case InputTypeAudio:
		// Read audio bytes directly — no FFmpeg needed.
		if cfg.Verbose && !cfg.Quiet {
			fmt.Printf("ℹ️  Audio file detected (%s), skipping FFmpeg extraction\n", mimeType)
		}

		data, err := os.ReadFile(cleanPath) // #nosec G304 -- cleanPath validated above
		if err != nil {
			return nil, fmt.Errorf("failed to read audio file: %w", err)
		}

		return &PreparedAudio{
			Data:     data,
			MIMEType: mimeType,
			Cleanup:  func() {}, // nothing to clean up
		}, nil

	case InputTypeVideo:
		// Extract audio via FFmpeg to a temporary WAV file, then read it.
		audioPath, err := extractAudio(cleanPath, cfg)
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
					if cfg.Verbose && !cfg.Quiet {
						fmt.Printf("ℹ️  Warning: failed to remove temp audio file: %v\n", removeErr)
					}
				}
			},
		}, nil

	default:
		return nil, fmt.Errorf("unsupported input type for %q", cleanPath)
	}
}

// extractAudio extracts audio from a video file using FFmpeg.
func extractAudio(videoPath string, cfg *config.Config) (string, error) {
	// Create a context with timeout for FFmpeg operation (30 minutes max)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	if cfg.Verbose && !cfg.Quiet {
		fmt.Printf("ℹ️  Extracting audio from video: %s\n", videoPath)
	}

	// Find full path to FFmpeg executable for security
	ffmpegPath, err := exec.LookPath("ffmpeg")
	if err != nil {
		return "", fmt.Errorf("FFmpeg not found. Please install FFmpeg first")
	}

	// Create temporary audio file
	audioPath := generateAudioPath(videoPath)

	// Run FFmpeg and verify output
	if err := runFFmpegCommand(ctx, ffmpegPath, videoPath, audioPath, cfg); err != nil {
		return "", err
	}

	if cfg.Verbose && !cfg.Quiet {
		fmt.Printf("ℹ️  Audio extracted to: %s\n", audioPath)
	}

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

	// 10 GB limit — enough for multi-hour videos
	if fileInfo.Size() > 10*1024*1024*1024 {
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
		"-y", // overwrite output file
		audioPath,
	)

	var stderr strings.Builder
	if cfg.Verbose {
		cmd.Stderr = io.MultiWriter(os.Stderr, &stderr)
	} else {
		cmd.Stderr = &stderr
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("FFmpeg failed: %v, stderr: %s", err, stderr.String())
	}

	if _, err := os.Stat(audioPath); err != nil {
		return fmt.Errorf("FFmpeg did not create output file: %w", err)
	}

	return nil
}
