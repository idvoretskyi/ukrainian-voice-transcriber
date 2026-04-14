// Voice Transcriber
// Copyright (c) 2025 Ihor Dvoretskyi
//
// Licensed under MIT License

package gemini_test

import (
	"strings"
	"testing"

	"github.com/idvoretskyi/voice-transcriber/internal/gemini"
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

func TestBuildPrompt(t *testing.T) {
	t.Parallel()

	t.Run("auto returns original spoken language prompt", func(t *testing.T) {
		t.Parallel()

		p := gemini.BuildPrompt("auto")
		if !strings.Contains(p, "original spoken language") {
			t.Errorf("BuildPrompt(%q) = %q; want it to mention 'original spoken language'", "auto", p)
		}
	})

	t.Run("empty string returns original spoken language prompt", func(t *testing.T) {
		t.Parallel()

		p := gemini.BuildPrompt("")
		if !strings.Contains(p, "original spoken language") {
			t.Errorf("BuildPrompt(%q) = %q; want it to mention 'original spoken language'", "", p)
		}
	})

	t.Run("ISO code produces language-specific prompt", func(t *testing.T) {
		t.Parallel()

		p := gemini.BuildPrompt("uk")
		if strings.Contains(p, "original spoken language") {
			t.Errorf("BuildPrompt(%q) should not mention 'original spoken language'", "uk")
		}

		if !strings.Contains(p, "uk") {
			t.Errorf("BuildPrompt(%q) = %q; want it to contain the language code", "uk", p)
		}
	})

	t.Run("prompt always contains transcription instructions", func(t *testing.T) {
		t.Parallel()

		for _, lang := range []string{"auto", "", "uk", "en", "de"} {
			p := gemini.BuildPrompt(lang)
			if !strings.Contains(p, "Output only the transcription text") {
				t.Errorf("BuildPrompt(%q) missing standard instructions", lang)
			}
		}
	})
}
