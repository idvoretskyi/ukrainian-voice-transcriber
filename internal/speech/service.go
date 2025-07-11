// Ukrainian Voice Transcriber

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

// TranscribeAudio transcribes audio using Google Cloud Speech-to-Text.
func (s *Service) TranscribeAudio(ctx context.Context, gcsURI string) (string, error) {
	if !s.config.Quiet {
		fmt.Println("ℹ️  Starting transcription...")
	}

	// Create recognition request
	req := &speechpb.RecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:                   speechpb.RecognitionConfig_LINEAR16,
			SampleRateHertz:            16000,
			LanguageCode:               "uk-UA",   // Ukrainian
			Model:                      "default", // Cost-efficient standard model
			EnableAutomaticPunctuation: true,
			UseEnhanced:                false, // Keep costs down
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Uri{
				Uri: gcsURI,
			},
		},
	}

	// Try synchronous recognition first
	resp, err := s.client.Recognize(ctx, req)
	if err != nil {
		// Try long-running recognition for longer audio
		return s.transcribeLongRunning(ctx, gcsURI)
	}

	// Extract transcript
	transcript := s.extractTranscript(resp.Results)

	if !s.config.Quiet {
		fmt.Printf("ℹ️  Transcription completed: %d characters\n", len(transcript))
	}

	return transcript, nil
}

// transcribeLongRunning handles long-running transcription.
func (s *Service) transcribeLongRunning(ctx context.Context, gcsURI string) (string, error) {
	if !s.config.Quiet {
		fmt.Println("ℹ️  Using long-running recognition...")
	}

	req := &speechpb.LongRunningRecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:                   speechpb.RecognitionConfig_LINEAR16,
			SampleRateHertz:            16000,
			LanguageCode:               "uk-UA",
			Model:                      "default",
			EnableAutomaticPunctuation: true,
			UseEnhanced:                false,
		},
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
