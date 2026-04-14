// Ukrainian Voice Transcriber
// Copyright (c) 2025 Ihor Dvoretskyi
//
// Licensed under MIT License

// Package cli provides command-line interface functionality.
package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/idvoretskyi/ukrainian-voice-transcriber/internal/gemini"
	"github.com/idvoretskyi/ukrainian-voice-transcriber/pkg/config"
)

const (
	appName = "Ukrainian Voice Transcriber"
)

var globalConfig config.Config

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "ukrainian-voice-transcriber",
	Short: "AI-powered Ukrainian media-to-text transcription",
	Long: fmt.Sprintf(`%s v%s

Professional Ukrainian media-to-text transcription using Google Gemini via Vertex AI.

Features:
• Ukrainian language optimized (uk-UA)
• Accepts both video files (mp4, mkv, mov, ...) and audio files (wav, mp3, flac, ...)
• No Google Cloud Storage required — audio sent inline to Gemini
• FFmpeg used only for video-to-audio extraction
• Cost-efficient: default model ~$0.03/hr of audio
• Single binary - no extra runtime dependencies

Prerequisites:
• FFmpeg installed (brew install ffmpeg / apt install ffmpeg) — for video files only
• Google Cloud authentication (gcloud auth application-default login)
• Vertex AI API enabled (gcloud services enable aiplatform.googleapis.com)
• GCP project configured (gcloud config set project YOUR_PROJECT_ID)

Examples:
  ukrainian-voice-transcriber transcribe input/video.mp4
  ukrainian-voice-transcriber transcribe input/recording.wav
  ukrainian-voice-transcriber transcribe input/video.mp4 -o output.txt
  ukrainian-voice-transcriber transcribe input/video.mp4 --verbose
  ukrainian-voice-transcriber transcribe input/video.mp4 --model gemini-2.5-flash
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

	// Gemini model selection
	rootCmd.PersistentFlags().StringVar(&globalConfig.GeminiModel, "model", gemini.DefaultModel,
		"Gemini model to use for transcription (e.g. gemini-3.1-flash-lite-preview, gemini-2.5-flash, gemini-2.5-flash-lite)")

	// Vertex AI region
	rootCmd.PersistentFlags().StringVar(&globalConfig.GCPLocation, "location", gemini.DefaultLocation,
		"Vertex AI region (e.g. us-central1, europe-west4)")

	// Add subcommands
	rootCmd.AddCommand(transcribeCmd)
	rootCmd.AddCommand(versionCmd)
}
