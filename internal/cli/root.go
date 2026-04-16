// Voice Transcriber
// Copyright (c) 2025 Ihor Dvoretskyi
//
// Licensed under MIT License

// Package cli provides command-line interface functionality.
package cli

import (
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/spf13/cobra"

	"github.com/idvoretskyi/voice-transcriber/internal/config"
	"github.com/idvoretskyi/voice-transcriber/internal/gemini"
)

const appName = "Voice Transcriber"

// NewRootCmd builds and returns the root Cobra command with all subcommands
// wired in. cfg is the shared configuration that persistent flags write into.
func NewRootCmd(cfg *config.Config) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "voice-transcriber",
		Short: "AI-powered media-to-text transcription with automatic language detection",
		Long: `Multilingual media-to-text transcription using Google Gemini via Vertex AI.
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
  voice-transcriber transcribe input/video.mp4 --model gemini-3-flash-preview
  voice-transcriber transcribe input/video.mp4 --language uk
  voice-transcriber version`,
		SilenceUsage: true,
		// Validate config flags before any subcommand runs.
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			return cfg.Validate()
		},
	}

	// Persistent flags — write directly into the shared cfg.
	rootCmd.PersistentFlags().BoolVarP(&cfg.Verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().BoolVarP(&cfg.Quiet, "quiet", "q", false, "Suppress all output except results")
	rootCmd.PersistentFlags().StringVar(&cfg.Language, "language", "auto",
		"Language for transcription: 'auto' for automatic detection, or ISO 639-1 code (e.g. uk, en, de)")
	rootCmd.PersistentFlags().StringVar(&cfg.GeminiModel, "model", gemini.DefaultModel,
		"Gemini model to use for transcription (e.g. gemini-3.1-flash-lite-preview, gemini-3-flash-preview)")
	rootCmd.PersistentFlags().StringVar(&cfg.GCPLocation, "location", gemini.DefaultLocation,
		"Vertex AI location (e.g. global, us-central1, europe-west4); Gemini 3.x models require global")

	rootCmd.AddCommand(newTranscribeCmd(cfg))
	rootCmd.AddCommand(newVersionCmd())

	return rootCmd
}

// Execute builds the command tree and runs it.
// This is called by main.main().
func Execute() error {
	cfg := config.FromEnv()
	root := NewRootCmd(cfg)
	// Wire the runtime version into Cobra's built-in --version flag.
	root.Version = buildVersion

	if err := root.Execute(); err != nil {
		return fmt.Errorf("executing root command: %w", err)
	}

	return nil
}

// newLogger returns a slog.Logger appropriate for the current config:
//   - Quiet: all output discarded
//   - Verbose: Debug level to stderr
//   - Default: Info level to stderr
func newLogger(cfg *config.Config) *slog.Logger {
	if cfg.Quiet {
		return slog.New(slog.NewTextHandler(io.Discard, nil))
	}

	level := slog.LevelInfo
	if cfg.Verbose {
		level = slog.LevelDebug
	}

	return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level}))
}
