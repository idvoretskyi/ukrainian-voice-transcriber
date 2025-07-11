// Ukrainian Voice Transcriber
//
// Copyright (c) 2025 Ihor Dvoretskyi
//
// Licensed under MIT License

// Package cli provides command-line interface functionality.
package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/idvoretskyi/ukrainian-voice-transcriber/internal/transcriber"
	"github.com/spf13/cobra"
)

var outputFile string

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
			return fmt.Errorf("initialization failed: %v", err)
		}
		defer t.Close()

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
			if err := os.WriteFile(outputFile, []byte(result.Text), 0600); err != nil {
				return fmt.Errorf("failed to save transcript: %v", err)
			}
			if !globalConfig.Quiet {
				fmt.Printf("✅ Transcript saved to: %s\n", outputFile)
			}
		} else {
			// Create directory based on video filename
			videoBaseName := filepath.Base(videoFile)
			videoNameWithoutExt := strings.TrimSuffix(videoBaseName, filepath.Ext(videoBaseName))

			// Replace spaces with underscores
			dirName := strings.ReplaceAll(videoNameWithoutExt, " ", "_")

			// Create directory
			if err := os.MkdirAll(dirName, 0o755); err != nil {
				return fmt.Errorf("failed to create directory %s: %v", dirName, err)
			}

			// Save transcript to file in the new directory with original filename (spaces replaced with underscores)
			transcriptPath := filepath.Join(dirName, dirName+".txt")
			if err := os.WriteFile(transcriptPath, []byte(result.Text), 0600); err != nil {
				return fmt.Errorf("failed to save transcript: %v", err)
			}

			if !globalConfig.Quiet {
				fmt.Printf("✅ Transcript saved to: %s\n", transcriptPath)
			}
		}

		return nil
	},
}

func init() {
	transcribeCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file path (default: stdout)")
}
