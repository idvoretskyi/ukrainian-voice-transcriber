// Voice Transcriber
// Copyright (c) 2025 Ihor Dvoretskyi
//
// Licensed under MIT License

// Package config provides configuration structures, loading, and validation.
package config

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// iso639Re matches exactly two lowercase ASCII letters (ISO 639-1 code).
var iso639Re = regexp.MustCompile(`^[a-z]{2}$`)

// Config holds application configuration.
type Config struct {
	Verbose bool
	Quiet   bool

	// Language for transcription. "auto" or "" means automatic detection.
	// Otherwise use an ISO 639-1 code (e.g. "uk", "en", "de").
	Language string

	// Gemini model selection
	GeminiModel string // e.g., "gemini-3.1-flash-lite-preview", "gemini-3-flash-preview"
	GCPLocation string // Vertex AI location, e.g., "global", "us-central1"

	// GCPProject is the Google Cloud project ID. Populated by FromEnv or
	// resolved at runtime via gcloud when empty.
	GCPProject string
}

// FromEnv returns a Config pre-populated from well-known environment variables.
// It does not validate — call Validate() on the result if needed.
//
//   - GOOGLE_CLOUD_PROJECT  → GCPProject
func FromEnv() *Config {
	return &Config{
		GCPProject: os.Getenv("GOOGLE_CLOUD_PROJECT"),
	}
}

// Validate returns an error if the Config contains contradictory or clearly
// invalid field values.
func (c *Config) Validate() error {
	if c.Verbose && c.Quiet {
		return fmt.Errorf("--verbose and --quiet are mutually exclusive")
	}

	lang := strings.ToLower(strings.TrimSpace(c.Language))
	if lang != "" && lang != "auto" && !iso639Re.MatchString(lang) {
		return fmt.Errorf(
			"invalid --language %q: must be 'auto', empty, or a two-letter ISO 639-1 code (e.g. 'uk', 'en')",
			c.Language,
		)
	}

	if c.GeminiModel != "" && strings.TrimSpace(c.GeminiModel) == "" {
		return fmt.Errorf("--model must not be blank")
	}

	return nil
}
