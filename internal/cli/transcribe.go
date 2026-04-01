// Ukrainian Voice Transcriber
// Copyright (c) 2025 Ihor Dvoretskyi
//
// Licensed under MIT License

// Package cli provides command-line interface functionality.
package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"

	"github.com/idvoretskyi/ukrainian-voice-transcriber/internal/transcriber"
)

var outputFile string

// sanitizeRe matches characters not allowed in a sanitized filename.
// \p{L} matches any Unicode letter (including Cyrillic), \p{N} any Unicode digit.
var sanitizeRe = regexp.MustCompile(`[^\p{L}\p{N}_\-.]`)

// transcribeCmd represents the transcribe command.
var transcribeCmd = &cobra.Command{
	Use:   "transcribe [video-file]",
	Short: "Transcribe a video file to text",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		videoFile := args[0]

		// Initialize transcriber
		t, err := transcriber.New(&globalConfig)
		if err != nil {
			return fmt.Errorf("initialization failed: %w", err)
		}

		defer func() {
			if closeErr := t.Close(); closeErr != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Warning: Failed to close transcriber: %v\n", closeErr)
			}
		}()

		// Transcribe file (input validation is performed inside TranscribeLocalFile)
		result := t.TranscribeLocalFile(context.Background(), videoFile)

		if !result.Success {
			return fmt.Errorf("transcription failed: %s", result.Error)
		}

		// Display results
		if !globalConfig.Quiet {
			fmt.Printf("\n📝 Transcription completed:\n")
			fmt.Printf("   Words: %d\n", result.WordCount)
			fmt.Printf("   Characters: %d\n", len(result.Text))
			fmt.Printf("   Processing time: %v\n", result.ProcessingTime)
			fmt.Println(strings.Repeat("-", 50))
		}

		// Determine output path
		var transcriptPath string
		if outputFile != "" {
			// User specified output file
			transcriptPath = outputFile
		} else {
			// Create path in output/ directory based on video filename
			videoBaseName := filepath.Base(videoFile)
			videoNameWithoutExt := strings.TrimSuffix(videoBaseName, filepath.Ext(videoBaseName))

			// Sanitize filename: replace spaces with underscores and remove special characters
			sanitizedName := sanitizeFilename(videoNameWithoutExt)

			// Create output subdirectory: output/<sanitized-name>/
			outputSubDir := filepath.Join("output", sanitizedName)

			// Create the full output path
			transcriptPath = filepath.Join(outputSubDir, sanitizedName+".txt")
		}

		// Ensure the directory for the output file exists
		outputDir := filepath.Dir(transcriptPath)
		if err := os.MkdirAll(outputDir, 0o750); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}

		// Save transcript with secure file permissions (0600 = rw-------)
		if err := os.WriteFile(transcriptPath, []byte(result.Text), 0o600); err != nil {
			return fmt.Errorf("failed to save transcript: %w", err)
		}

		if !globalConfig.Quiet {
			fmt.Printf("✅ Transcript saved to: %s\n", transcriptPath)
		}

		return nil
	},
}

// sanitizeFilename removes special characters and replaces spaces with underscores
// to create a safe filename for use in the filesystem.
// Preserves Unicode letters (including Cyrillic) for Ukrainian filenames.
func sanitizeFilename(filename string) string {
	// Replace spaces with underscores
	filename = strings.ReplaceAll(filename, " ", "_")

	// Remove any character that's not Unicode letter, digit, underscore, hyphen, or period
	filename = sanitizeRe.ReplaceAllString(filename, "")

	// Ensure the filename is not empty
	if filename == "" {
		return "transcript"
	}

	return filename
}

func init() {
	transcribeCmd.Flags().StringVarP(&outputFile, "output", "o", "",
		"Output file path (default: creates directory based on video filename)")
}
