// Voice Transcriber
// Copyright (c) 2025 Ihor Dvoretskyi
//
// Licensed under MIT License

// Package config provides configuration structures and utilities.
package config

// Config holds application configuration.
type Config struct {
	Verbose bool
	Quiet   bool

	// Language for transcription. "auto" or "" means automatic detection.
	// Otherwise use an ISO 639-1 code (e.g. "uk", "en", "de").
	Language string

	// Gemini model selection
	GeminiModel string // e.g., "gemini-3.1-flash-lite-preview", "gemini-2.5-flash"
	GCPLocation string // Vertex AI region, e.g., "us-central1"
}
