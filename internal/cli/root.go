// Ukrainian Voice Transcriber

// Package cli provides command-line interface functionality.
package cli

import (
	"fmt"

	"github.com/idvoretskyi/ukrainian-voice-transcriber/pkg/config"
	"github.com/spf13/cobra"
)

const (
	version = "1.0.0"
	appName = "Ukrainian Voice Transcriber"
)

var globalConfig config.Config

// rootCmd represents the base command when called without any subcommands
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
• Google Cloud service account JSON file in current directory
• Enabled APIs: Speech-to-Text, Cloud Storage

Examples:
  ukrainian-voice-transcriber transcribe video.mp4
  ukrainian-voice-transcriber transcribe video.mp4 -o transcript.txt
  ukrainian-voice-transcriber transcribe video.mp4 --verbose
  ukrainian-voice-transcriber version`, appName, version),
	SilenceUsage: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().BoolVarP(&globalConfig.Verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().BoolVarP(&globalConfig.Quiet, "quiet", "q", false, "Suppress all output except results")
	rootCmd.PersistentFlags().StringVar(&globalConfig.BucketName, "bucket", "", "Google Cloud Storage bucket name (default: ukr-voice-transcriber-temp)")

	// Add subcommands
	rootCmd.AddCommand(authCmd)
	rootCmd.AddCommand(transcribeCmd)
	rootCmd.AddCommand(setupCmd)
	rootCmd.AddCommand(versionCmd)
}
