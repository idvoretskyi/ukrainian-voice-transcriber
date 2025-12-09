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

// extractAudio extracts audio from video file using FFmpeg.
func extractAudio(videoPath string, cfg *config.Config) (string, error) {
	// Create a context with timeout for FFmpeg operation (5 minutes)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Validate and sanitize the input path
	cleanPath, err := validateAndSanitizeVideoPath(videoPath)
	if err != nil {
		return "", err
	}

	if cfg.Verbose && !cfg.Quiet {
		fmt.Printf("🔍 Extracting audio from: %s\n", cleanPath)
	}

	// Find full path to FFmpeg executable for security
	ffmpegPath, err := exec.LookPath("ffmpeg")
	if err != nil {
		return "", fmt.Errorf("FFmpeg not found. Please install FFmpeg first")
	}

	// Create temporary audio file with sanitized path
	audioPath := generateAudioPath(cleanPath)

	// Run FFmpeg command and verify output
	err = runFFmpegCommand(ctx, ffmpegPath, cleanPath, audioPath, cfg)
	if err != nil {
		return "", err
	}

	if cfg.Verbose && !cfg.Quiet {
		fmt.Printf("🔍 Audio extracted to: %s\n", audioPath)
	}

	return audioPath, nil
}

// validateAndSanitizeVideoPath validates and sanitizes the input video path.
func validateAndSanitizeVideoPath(videoPath string) (string, error) {
	// Use filepath.Clean to normalize the path
	videoPath = filepath.Clean(videoPath)

	// Validate file exists and is readable
	fileInfo, err := os.Stat(videoPath)
	if err != nil {
		return "", fmt.Errorf("input file error: %v", err)
	}

	// Check if it's a regular file
	if !fileInfo.Mode().IsRegular() {
		return "", fmt.Errorf("not a regular file: %s", videoPath)
	}

	// Check file size (prevent processing extremely large files)
	// 5GB limit should be sufficient for most video files
	if fileInfo.Size() > 5*1024*1024*1024 {
		return "", fmt.Errorf("file too large (>5GB): %s", videoPath)
	}

	return videoPath, nil
}

// generateAudioPath creates a unique audio file path in the system temp directory.
func generateAudioPath(videoPath string) string {
	timestamp := time.Now().UnixNano()
	baseFileName := strings.TrimSuffix(filepath.Base(videoPath), filepath.Ext(videoPath))
	tempFileName := fmt.Sprintf("%s_%d_audio.wav", baseFileName, timestamp)

	// Use system temp directory instead of current directory
	return filepath.Join(os.TempDir(), tempFileName)
}

// runFFmpegCommand executes FFmpeg command and verifies the output.
func runFFmpegCommand(ctx context.Context, ffmpegPath, videoPath, audioPath string, cfg *config.Config) error {
	// Run FFmpeg command with context for timeout
	cmd := exec.CommandContext(ctx, ffmpegPath,
		"-i", videoPath,
		"-acodec", "pcm_s16le",
		"-ar", "16000",
		"-ac", "1",
		"-y", // Overwrite output file
		audioPath,
	)

	// Capture both stdout and stderr
	var stderr strings.Builder
	if cfg.Verbose {
		cmd.Stderr = io.MultiWriter(os.Stderr, &stderr)
	} else {
		cmd.Stderr = &stderr
	}

	// Run the command
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("FFmpeg failed: %v, stderr: %s", err, stderr.String())
	}

	// Verify the output file was created
	if _, err := os.Stat(audioPath); err != nil {
		return fmt.Errorf("FFmpeg did not create output file: %v", err)
	}

	return nil
}
