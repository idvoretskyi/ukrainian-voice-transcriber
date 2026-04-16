// Voice Transcriber
// Copyright (c) 2025 Ihor Dvoretskyi
//
// Licensed under MIT License

package transcriber

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	// ffmpegTimeout is the maximum time allowed for a single FFmpeg extraction.
	ffmpegTimeout = 30 * time.Minute

	// maxFileSize is the maximum accepted input file size (10 GB).
	maxFileSize = 10 * 1024 * 1024 * 1024

	// ffmpegSampleRate is the PCM sample rate used for audio extraction.
	ffmpegSampleRate = "16000"

	// ffmpegChannels is the number of audio channels (mono) used for extraction.
	ffmpegChannels = "1"
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

// PreparedAudio holds the audio bytes and MIME type ready to send to Gemini.
// Call Close() to remove any temporary file that was created during preparation.
type PreparedAudio struct {
	Data     []byte
	MIMEType string
	// tempPath is the path of a temporary file to remove on Close, or empty
	// when no temporary file was created (e.g. native audio input).
	tempPath string
}

// Close removes the temporary audio file if one was created during preparation.
// It is safe to call Close on a PreparedAudio that has no temporary file.
// Implements io.Closer.
func (p *PreparedAudio) Close() error {
	if p.tempPath == "" {
		return nil
	}

	if err := os.Remove(p.tempPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove temp audio file %q: %w", p.tempPath, err)
	}

	return nil
}

// prepareAudio reads the input file, extracting audio via FFmpeg when the
// input is a video, and returns a PreparedAudio ready to pass to the Gemini
// service.
func prepareAudio(ctx context.Context, inputPath string, logger *slog.Logger) (*PreparedAudio, error) {
	cleanPath, err := validateInputPath(inputPath)
	if err != nil {
		return nil, err
	}

	inputType, mimeType := classifyInputFile(cleanPath)

	switch inputType {
	case InputTypeAudio:
		logger.InfoContext(ctx, "audio file detected, skipping FFmpeg extraction",
			slog.String("mime", mimeType))

		data, err := os.ReadFile(cleanPath) // #nosec G304 -- cleanPath validated above
		if err != nil {
			return nil, fmt.Errorf("failed to read audio file: %w", err)
		}

		return &PreparedAudio{
			Data:     data,
			MIMEType: mimeType,
		}, nil

	case InputTypeVideo:
		audioPath, err := extractAudio(ctx, cleanPath, logger)
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
			tempPath: audioPath,
		}, nil

	default:
		// unreachable unless InputType is extended
		return nil, fmt.Errorf("unsupported input type for %q", cleanPath)
	}
}

// extractAudio extracts audio from a video file using FFmpeg.
func extractAudio(ctx context.Context, videoPath string, logger *slog.Logger) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, ffmpegTimeout)
	defer cancel()

	logger.InfoContext(ctx, "extracting audio from video", slog.String("path", videoPath))

	ffmpegPath, err := exec.LookPath("ffmpeg")
	if err != nil {
		return "", fmt.Errorf("ffmpeg not found; install ffmpeg first: %w", err)
	}

	audioPath, err := generateAudioPath(videoPath)
	if err != nil {
		return "", fmt.Errorf("failed to create temp audio file: %w", err)
	}

	if err := runFFmpegCommand(ctx, ffmpegPath, videoPath, audioPath, logger); err != nil {
		return "", err
	}

	logger.InfoContext(ctx, "audio extracted", slog.String("output", audioPath))

	return audioPath, nil
}

// validateInputPath validates and sanitizes any input media path (audio or video).
func validateInputPath(inputPath string) (string, error) {
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

// generateAudioPath creates a unique WAV file path in the system temp directory
// using os.CreateTemp to avoid any TOCTOU race between path generation and file creation.
// The caller is responsible for removing the file when done.
func generateAudioPath(inputPath string) (string, error) {
	baseFileName := strings.TrimSuffix(filepath.Base(inputPath), filepath.Ext(inputPath))
	pattern := baseFileName + "_*_audio.wav"

	f, err := os.CreateTemp("", pattern)
	if err != nil {
		return "", fmt.Errorf("creating temp audio file: %w", err)
	}

	// Close immediately; the file is only needed as a reserved path for FFmpeg.
	if err := f.Close(); err != nil {
		_ = os.Remove(f.Name())

		return "", fmt.Errorf("closing temp audio file: %w", err)
	}

	return f.Name(), nil
}

// runFFmpegCommand executes FFmpeg and verifies the output was created.
func runFFmpegCommand(ctx context.Context, ffmpegPath, videoPath, audioPath string, logger *slog.Logger) error {
	cmd := exec.CommandContext(ctx, ffmpegPath, // #nosec G204 -- ffmpegPath resolved via exec.LookPath
		"-i", videoPath,
		"-acodec", "pcm_s16le",
		"-ar", ffmpegSampleRate,
		"-ac", ffmpegChannels,
		"-y",
		audioPath,
	)

	var stderr strings.Builder

	// Always capture stderr; tee to os.Stderr at debug level via slog when verbose.
	if logger.Enabled(ctx, slog.LevelDebug) {
		cmd.Stderr = io.MultiWriter(os.Stderr, &stderr)
	} else {
		cmd.Stderr = &stderr
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ffmpeg failed: %w (stderr: %s)", err, stderr.String())
	}

	if _, err := os.Stat(audioPath); err != nil {
		return fmt.Errorf("ffmpeg did not create output file: %w", err)
	}

	return nil
}
