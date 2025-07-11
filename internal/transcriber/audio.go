// Ukrainian Voice Transcriber

// Package transcriber provides audio transcription functionality.
package transcriber

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/idvoretskyi/ukrainian-voice-transcriber/pkg/config"
)

// extractAudio extracts audio from video file using FFmpeg.
func extractAudio(videoPath string, cfg *config.Config) (string, error) {
	// Validate and sanitize the input path
	if filepath.IsAbs(videoPath) {
		return "", fmt.Errorf("absolute paths not allowed")
	}

	// Use filepath.Clean to normalize the path
	videoPath = filepath.Clean(videoPath)

	// Ensure the path doesn't contain directory traversal
	if strings.Contains(videoPath, "..") {
		return "", fmt.Errorf("directory traversal not allowed")
	}

	// Validate file exists
	if _, err := os.Stat(videoPath); err != nil {
		return "", fmt.Errorf("input file does not exist: %s", videoPath)
	}

	if cfg.Verbose && !cfg.Quiet {
		fmt.Printf("🔍 Extracting audio from: %s\n", videoPath)
	}

	// Check if FFmpeg is available
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return "", fmt.Errorf("FFmpeg not found. Please install FFmpeg first")
	}

	// Create temporary audio file with sanitized path
	audioPath := strings.TrimSuffix(videoPath, filepath.Ext(videoPath)) + "_audio.wav"

	// Run FFmpeg command
	cmd := exec.Command("ffmpeg",
		"-i", videoPath,
		"-acodec", "pcm_s16le",
		"-ar", "16000",
		"-ac", "1",
		"-y", // Overwrite output file
		audioPath,
	)

	if cfg.Verbose {
		cmd.Stderr = os.Stderr
	}

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("FFmpeg failed: %v", err)
	}

	if cfg.Verbose && !cfg.Quiet {
		fmt.Printf("🔍 Audio extracted to: %s\n", audioPath)
	}

	return audioPath, nil
}
