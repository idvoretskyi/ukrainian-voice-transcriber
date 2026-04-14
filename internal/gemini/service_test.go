// Ukrainian Voice Transcriber
// Copyright (c) 2025 Ihor Dvoretskyi
//
// Licensed under MIT License

package gemini_test

import (
	"strings"
	"testing"

	"github.com/idvoretskyi/ukrainian-voice-transcriber/internal/gemini"
)

func TestDefaultModel(t *testing.T) {
	t.Parallel()

	if gemini.DefaultModel == "" {
		t.Error("DefaultModel must not be empty")
	}

	if !strings.Contains(gemini.DefaultModel, "gemini") {
		t.Errorf("DefaultModel %q does not look like a Gemini model ID", gemini.DefaultModel)
	}
}

func TestDefaultLocation(t *testing.T) {
	t.Parallel()

	if gemini.DefaultLocation == "" {
		t.Error("DefaultLocation must not be empty")
	}
}
