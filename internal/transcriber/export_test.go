// Voice Transcriber
// Copyright (c) 2025 Ihor Dvoretskyi
//
// Licensed under MIT License

// Package transcriber exports internal symbols for testing.
package transcriber

import (
	"context"
	"log/slog"

	"github.com/idvoretskyi/voice-transcriber/internal/config"
	"github.com/idvoretskyi/voice-transcriber/internal/gemini"
)

// GenerateAudioPath exposes generateAudioPath for black-box tests.
func GenerateAudioPath(inputPath string) (string, error) { return generateAudioPath(inputPath) }

// ValidateInputPath exposes validateInputPath for black-box tests.
func ValidateInputPath(inputPath string) (string, error) {
	return validateInputPath(inputPath)
}

// ClassifyInputFile exposes classifyInputFile for black-box tests.
func ClassifyInputFile(inputPath string) (InputType, string) {
	return classifyInputFile(inputPath)
}

// NewForTesting constructs a Transcriber with an injected AudioTranscriber
// backend, bypassing real Gemini and gcloud resolution.
// For use in unit tests only.
func NewForTesting(cfg *config.Config, backend gemini.AudioTranscriber, logger *slog.Logger) *Transcriber {
	if logger == nil {
		logger = slog.Default()
	}

	return &Transcriber{
		config:    cfg,
		backend:   backend,
		logger:    logger,
		resolveID: func(_ context.Context) (string, error) { return "test-project", nil },
	}
}
