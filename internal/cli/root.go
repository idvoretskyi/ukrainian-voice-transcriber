// Voice Transcriber
// Copyright (c) 2025 Ihor Dvoretskyi
//
// Licensed under MIT License

// Package cli provides command-line interface functionality.
package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/idvoretskyi/voice-transcriber/internal/gemini"
	"github.com/idvoretskyi/voice-transcriber/pkg/config"
)

const (
	appName = "Voice Transcriber"
)

var globalConfig config.Config

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "voice-transcriber",
	Short: "AI-powered media-to-text transcription with automatic language detection",
	Long: fmt.Sprintf(`%s v%s

Multilingual media-to-text transcription using Google Gemini via Vertex AI.
Language is detected automatically from the audio by default.

Features:
• Automatic language detection (default) or specify with --language
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
  voice-transcriber transcribe input/video.mp4
  voice-transcriber transcribe input/recording.wav
  voice-transcriber transcribe input/video.mp4 -o output.txt
  voice-transcriber transcribe input/video.mp4 --verbose
  voice-transcriber transcribe input/video.mp4 --model gemini-2.5-flash
  voice-transcriber transcribe input/video.mp4 --language uk
  voice-transcriber version`, appName, buildVersion),
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

	// Language selection
	rootCmd.PersistentFlags().StringVar(&globalConfig.Language, "language", "auto",
		"Language for transcription: 'auto' for automatic detection, or ISO 639-1 code (e.g. uk, en, de)")

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
