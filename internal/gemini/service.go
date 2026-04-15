// Voice Transcriber
// Copyright (c) 2025 Ihor Dvoretskyi
//
// Licensed under MIT License

// Package gemini provides Google Gemini transcription via Vertex AI.
package gemini

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"cloud.google.com/go/vertexai/genai"

	"github.com/idvoretskyi/voice-transcriber/pkg/config"
)

const (
	// DefaultModel is the default Gemini model for transcription.
	// gemini-3.1-flash-lite-preview is optimized for ASR tasks and is the
	// most cost-effective model with explicit audio/ASR quality improvements.
	DefaultModel = "gemini-3.1-flash-lite-preview"

	// DefaultLocation is the default Vertex AI region.
	DefaultLocation = "us-central1"

	// autoLang is the sentinel value meaning automatic language detection.
	autoLang = "auto"
)

// iso639Re matches exactly two lowercase ASCII letters (ISO 639-1 code).
var iso639Re = regexp.MustCompile(`^[a-z]{2}$`)

// buildPrompt returns the transcription prompt for the given language.
// When language is "auto" or empty, Gemini detects the language automatically.
// Otherwise language must be a two-letter ISO 639-1 code (e.g. "uk", "en", "de").
// Inputs are normalized (trimmed, lowercased) and validated; invalid values fall
// back to automatic detection.
func buildPrompt(language string) string {
	const suffix = `
Output only the transcription text with no commentary, labels, or metadata.
Preserve natural sentence structure and add punctuation where appropriate.
Do not translate, summarize, or modify the content in any way.`

	lang := strings.ToLower(strings.TrimSpace(language))

	// Validate: accept "auto", empty string, or a two-letter ISO 639-1 code.
	if lang != "" && lang != autoLang {
		if !iso639Re.MatchString(lang) {
			// Invalid input — fall back to automatic detection.
			lang = autoLang
		}
	}

	if lang == "" || lang == autoLang {
		return "Transcribe the following audio recording verbatim in its original spoken language." + suffix
	}

	return "Transcribe the following audio recording verbatim in " + lang + "." + suffix
}

// Service handles Gemini transcription via Vertex AI.
type Service struct {
	client    *genai.Client
	config    *config.Config
	projectID string
}

// NewService creates a new Gemini service and initializes the Vertex AI client.
func NewService(ctx context.Context, cfg *config.Config, projectID string) (*Service, error) {
	location := cfg.GCPLocation
	if location == "" {
		location = DefaultLocation
	}

	client, err := genai.NewClient(ctx, projectID, location)
	if err != nil {
		return nil, fmt.Errorf("failed to create Vertex AI client: %w", err)
	}

	return &Service{
		client:    client,
		config:    cfg,
		projectID: projectID,
	}, nil
}

// TranscribeAudio sends audio bytes to Gemini and returns the transcript.
// mimeType must be one of: audio/wav, audio/mp3, audio/flac, audio/ogg,
// audio/m4a, audio/aac, audio/webm, audio/pcm.
func (s *Service) TranscribeAudio(ctx context.Context, audioData []byte, mimeType string) (string, error) {
	model := s.config.GeminiModel
	if model == "" {
		model = DefaultModel
	}

	if !s.config.Quiet {
		fmt.Printf("ℹ️  Sending audio to Gemini (model: %s, size: %s)\n",
			model, formatBytes(len(audioData)))
	}

	gm := s.client.GenerativeModel(model)

	// Build the request: text prompt + audio blob
	prompt := genai.Text(buildPrompt(s.config.Language))
	audioBlob := genai.Blob{
		MIMEType: mimeType,
		Data:     audioData,
	}

	resp, err := gm.GenerateContent(ctx, prompt, audioBlob)
	if err != nil {
		return "", fmt.Errorf("gemini generation failed: %w", err)
	}

	transcript := extractText(resp)
	if transcript == "" {
		return "", fmt.Errorf("gemini returned empty transcript")
	}

	if !s.config.Quiet {
		fmt.Printf("ℹ️  Transcription received: %d characters\n", len(transcript))
	}

	return transcript, nil
}

// Close closes the Vertex AI client.
func (s *Service) Close() error {
	if s.client != nil {
		if err := s.client.Close(); err != nil {
			return fmt.Errorf("closing vertex AI client: %w", err)
		}
	}

	return nil
}

// extractText pulls all text parts from a GenerateContentResponse.
func extractText(resp *genai.GenerateContentResponse) string {
	if resp == nil {
		return ""
	}

	var sb strings.Builder

	for _, cand := range resp.Candidates {
		if cand.Content == nil {
			continue
		}

		for _, part := range cand.Content.Parts {
			if t, ok := part.(genai.Text); ok {
				sb.WriteString(string(t))
			}
		}
	}

	return strings.TrimSpace(sb.String())
}

// formatBytes returns a human-readable byte size string.
func formatBytes(n int) string {
	const unit = 1024
	if n < unit {
		return fmt.Sprintf("%d B", n)
	}

	div, exp := unit, 0
	for v := n / unit; v >= unit; v /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %cB", float64(n)/float64(div), "KMGTPE"[exp])
}
