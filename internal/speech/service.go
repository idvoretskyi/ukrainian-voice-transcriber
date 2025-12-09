// Ukrainian Voice Transcriber
// Copyright (c) 2025 Ihor Dvoretskyi
//
// Licensed under MIT License

// Package speech provides Google Cloud Speech-to-Text functionality.
package speech

import (
	"context"
	"fmt"
	"strings"

	speechapi "cloud.google.com/go/speech/apiv1"
	"cloud.google.com/go/speech/apiv1/speechpb"

	"github.com/idvoretskyi/ukrainian-voice-transcriber/pkg/config"
)

// Service handles Google Cloud Speech-to-Text operations.
type Service struct {
	client *speechapi.Client
	config *config.Config
}

// NewService creates a new speech service.
func NewService(client *speechapi.Client, cfg *config.Config) *Service {
	return &Service{
		client: client,
		config: cfg,
	}
}

// TranscribeFromGCS transcribes audio from Google Cloud Storage URI.
func (s *Service) TranscribeFromGCS(ctx context.Context, gcsURI string) (string, error) {
	if !s.config.Quiet {
		fmt.Println("ℹ️  Starting transcription...")
	}

	// Determine model to use
	// Note: 'video' and 'phone_call' models don't support Ukrainian (uk-UA)
	// Using 'default' model which supports all languages
	model := s.config.STTModel
	if model == "" || model == "video" {
		model = "default"
	}

	// Create recognition config
	config := &speechpb.RecognitionConfig{
		Encoding:                   speechpb.RecognitionConfig_LINEAR16,
		SampleRateHertz:            16000,
		LanguageCode:               "uk-UA",
		Model:                      model,
		EnableAutomaticPunctuation: true,
		UseEnhanced:                false, // Enhanced models don't support Ukrainian
	}

	// Use long-running recognition for all audio (handles both short and long files)
	return s.transcribeLongRunning(ctx, gcsURI, config)
}

// Close closes the speech client (no-op since client is managed externally).
func (s *Service) Close() error {
	// Client is closed by the transcriber
	return nil
}

// transcribeLongRunning handles long-running transcription.
func (s *Service) transcribeLongRunning(
	ctx context.Context,
	gcsURI string,
	config *speechpb.RecognitionConfig,
) (string, error) {
	if !s.config.Quiet {
		fmt.Println("ℹ️  Using long-running recognition...")
	}

	req := &speechpb.LongRunningRecognizeRequest{
		Config: config,
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Uri{
				Uri: gcsURI,
			},
		},
	}

	op, err := s.client.LongRunningRecognize(ctx, req)
	if err != nil {
		return "", fmt.Errorf("long-running recognition failed: %v", err)
	}

	if !s.config.Quiet {
		fmt.Println("ℹ️  Waiting for transcription to complete...")
	}

	resp, err := op.Wait(ctx)
	if err != nil {
		return "", fmt.Errorf("transcription operation failed: %v", err)
	}

	transcript := s.extractTranscript(resp.Results)

	if !s.config.Quiet {
		fmt.Printf("ℹ️  Long transcription completed: %d characters\n", len(transcript))
	}

	return transcript, nil
}

// extractTranscript extracts text from speech recognition results.
func (s *Service) extractTranscript(results []*speechpb.SpeechRecognitionResult) string {
	var transcript strings.Builder

	for _, result := range results {
		if len(result.Alternatives) > 0 {
			transcript.WriteString(result.Alternatives[0].Transcript)
			transcript.WriteString(" ")
		}
	}

	return strings.TrimSpace(transcript.String())
}
