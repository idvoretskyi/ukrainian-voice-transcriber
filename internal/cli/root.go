// Ukrainian Voice Transcriber
// Copyright (c) 2025 Ihor Dvoretskyi
//
// Licensed under MIT License

// Package cli provides command-line interface functionality.
package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/idvoretskyi/ukrainian-voice-transcriber/pkg/config"
)

const (
	appName = "Ukrainian Voice Transcriber"
)

var globalConfig config.Config

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "ukrainian-voice-transcriber",
	Short: "AI-powered Ukrainian video-to-text transcription",
	Long: fmt.Sprintf(`%s v%s

Professional Ukrainian video-to-text transcription using Google Cloud Speech-to-Text API.

Features:
• Ukrainian language optimized (uk-UA locale)
• Cost-efficient with automatic cleanup
• Single binary - no dependencies to install
• FFmpeg integration for audio extraction
• Google Cloud Storage for temporary files

Prerequisites:
• FFmpeg installed (brew install ffmpeg / apt install ffmpeg)
• Google Cloud authentication (gcloud auth application-default login)
• Enabled APIs: Speech-to-Text, Cloud Storage

Examples:
  ukrainian-voice-transcriber transcribe input/video.mp4
  ukrainian-voice-transcriber transcribe input/video.mp4 -o output.txt
  ukrainian-voice-transcriber transcribe input/video.mp4 --verbose
  ukrainian-voice-transcriber version`, appName, buildVersion),
	SilenceUsage: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		return fmt.Errorf("command execution failed: %w", err)
	}

	return nil
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().BoolVarP(&globalConfig.Verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().BoolVarP(&globalConfig.Quiet, "quiet", "q", false, "Suppress all output except results")
	rootCmd.PersistentFlags().StringVar(&globalConfig.BucketName, "bucket", "",
		"Google Cloud Storage bucket name (default: ukr-voice-transcriber-temp)")

	// Speech-to-Text model selection
	rootCmd.PersistentFlags().StringVar(&globalConfig.STTModel, "model", "default",
		"Speech-to-Text model: 'default' (supports Ukrainian), 'latest_long', or 'latest_short'")

	// Add subcommands
	rootCmd.AddCommand(transcribeCmd)
	rootCmd.AddCommand(versionCmd)
}
