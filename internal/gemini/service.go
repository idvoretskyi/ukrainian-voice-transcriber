// Voice Transcriber
// Copyright (c) 2025 Ihor Dvoretskyi
//
// Licensed under MIT License

// Package gemini provides Google Gemini transcription via Vertex AI.
package gemini

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"strings"

	"google.golang.org/genai"

	"github.com/idvoretskyi/voice-transcriber/internal/config"
)

const (
	// DefaultModel is the default Gemini model for transcription.
	// gemini-3.1-flash-lite-preview is optimized for ASR tasks and is the
	// most cost-effective model with explicit audio/ASR quality improvements.
	DefaultModel = "gemini-3.1-flash-lite-preview"

	// DefaultLocation is the default Vertex AI location.
	// Gemini 3.x models are only available on the global endpoint.
	DefaultLocation = "global"

	// autoLang is the sentinel value meaning automatic language detection.
	autoLang = "auto"

	// roleUser is the Gemini content role for user turns.
	roleUser = "user"
)

// iso639Re matches exactly two lowercase ASCII letters (ISO 639-1 code).
var iso639Re = regexp.MustCompile(`^[a-z]{2}$`)

// AudioTranscriber is the interface for sending audio to a transcription backend.
// It is satisfied by *Service and can be replaced in tests by a stub.
type AudioTranscriber interface {
	TranscribeAudio(ctx context.Context, audioData []byte, mimeType string) (string, error)
}

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
	client *genai.Client
	config *config.Config
	logger *slog.Logger
}

// NewService creates a new Gemini service and initializes the Vertex AI client.
// The client uses Application Default Credentials automatically.
// If logger is nil, slog.Default() is used.
func NewService(ctx context.Context, cfg *config.Config, projectID string, logger *slog.Logger) (*Service, error) {
	if logger == nil {
		logger = slog.Default()
	}

	location := cfg.GCPLocation
	if location == "" {
		location = DefaultLocation
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		Project:  projectID,
		Location: location,
		Backend:  genai.BackendVertexAI,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Vertex AI client: %w", err)
	}

	return &Service{
		client: client,
		config: cfg,
		logger: logger,
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

	s.logger.InfoContext(ctx, "sending audio to Gemini",
		slog.String("model", model),
		slog.String("size", formatBytes(len(audioData))),
	)

	parts := []*genai.Part{
		{Text: buildPrompt(s.config.Language)},
		{InlineData: &genai.Blob{MIMEType: mimeType, Data: audioData}},
	}
	contents := []*genai.Content{{Role: roleUser, Parts: parts}}

	resp, err := s.client.Models.GenerateContent(ctx, model, contents, nil)
	if err != nil {
		return "", fmt.Errorf("gemini generation failed: %w", err)
	}

	transcript := strings.TrimSpace(resp.Text())
	if transcript == "" {
		return "", fmt.Errorf("gemini returned empty transcript")
	}

	s.logger.InfoContext(ctx, "transcription received",
		slog.Int("characters", len(transcript)),
	)

	return transcript, nil
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
