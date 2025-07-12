// Ukrainian Voice Transcriber
// Copyright (c) 2025 Ihor Dvoretskyi
//
// Licensed under MIT License

// Package cli provides command-line interface functionality.
package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"

	"github.com/idvoretskyi/ukrainian-voice-transcriber/internal/transcriber"
)

var outputFile string

// transcribeCmd represents the transcribe command.
var transcribeCmd = &cobra.Command{
	Use:   "transcribe [video-file]",
	Short: "Transcribe a video file to text",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		videoFile := args[0]

		// Validate input file exists and is accessible
		fileInfo, err := os.Stat(videoFile)
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("video file not found: %s", videoFile)
			}

			return fmt.Errorf("cannot access video file: %v", err)
		}

		// Check if it's a regular file (not a directory or device)
		if !fileInfo.Mode().IsRegular() {
			return fmt.Errorf("not a regular file: %s", videoFile)
		}

		// Initialize transcriber
		t, err := transcriber.New(&globalConfig)
		if err != nil {
			return fmt.Errorf("initialization failed: %v", err)
		}
		defer func() {
			if closeErr := t.Close(); closeErr != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Warning: Failed to close transcriber: %v\n", closeErr)
			}
		}()

		// Transcribe file
		result := t.TranscribeLocalFile(videoFile)

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

		// Create directory based on video filename and save transcript
		if outputFile != "" {
			// Ensure the directory for the output file exists
			outputDir := filepath.Dir(outputFile)
			if outputDir != "." && outputDir != "" {
				if err := os.MkdirAll(outputDir, 0o750); err != nil {
					return fmt.Errorf("failed to create directory for output file: %v", err)
				}
			}

			// Use secure file permissions (0600 = rw-------)
			if err := os.WriteFile(outputFile, []byte(result.Text), 0o600); err != nil {
				return fmt.Errorf("failed to save transcript: %v", err)
			}
			if !globalConfig.Quiet {
				fmt.Printf("✅ Transcript saved to: %s\n", outputFile)
			}
		} else {
			// Create directory based on video filename
			videoBaseName := filepath.Base(videoFile)
			videoNameWithoutExt := strings.TrimSuffix(videoBaseName, filepath.Ext(videoBaseName))

			// Sanitize directory name: replace spaces with underscores and remove special characters
			dirName := sanitizeFilename(videoNameWithoutExt)

			// Create directory with secure permissions (0750 = rwxr-x---)
			if err := os.MkdirAll(dirName, 0o750); err != nil {
				return fmt.Errorf("failed to create directory %s: %v", dirName, err)
			}

			// Save transcript to file in the new directory with sanitized filename
			transcriptPath := filepath.Join(dirName, dirName+".txt")
			// Use secure file permissions (0600 = rw-------)
			if err := os.WriteFile(transcriptPath, []byte(result.Text), 0o600); err != nil {
				return fmt.Errorf("failed to save transcript: %v", err)
			}

			if !globalConfig.Quiet {
				fmt.Printf("✅ Transcript saved to: %s\n", transcriptPath)
			}
		}

		return nil
	},
}

// sanitizeFilename removes special characters and replaces spaces with underscores
// to create a safe filename for use in the filesystem.
func sanitizeFilename(filename string) string {
	// Replace spaces with underscores
	filename = strings.ReplaceAll(filename, " ", "_")

	// Remove any character that's not alphanumeric, underscore, hyphen, or period
	reg := regexp.MustCompile(`[^a-zA-Z0-9_\-.]`)
	filename = reg.ReplaceAllString(filename, "")

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
