// Ukrainian Voice Transcriber
// Copyright (c) 2025 Ihor Dvoretskyi
//
// Licensed under MIT License

// Package config provides configuration structures and utilities.
package config

import (
	"os"
)

// Config holds application configuration.
type Config struct {
	ServiceAccountPath string
	Verbose            bool
	Quiet              bool

	// Gemini model selection
	GeminiModel string // e.g., "gemini-3.1-flash-lite-preview", "gemini-2.5-flash"
	GCPLocation string // Vertex AI region, e.g., "us-central1"
}

// FindServiceAccount looks for Google Cloud service account key.
func FindServiceAccount() string {
	candidates := []string{
		"service-account.json",
		"service_account.json",
		"gcloud-key.json",
		"key.json",
	}

	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}

	return ""
}
