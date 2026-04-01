// Ukrainian Voice Transcriber
// Copyright (c) 2025 Ihor Dvoretskyi
//
// Licensed under MIT License

// Package speech exports internal symbols for testing.
package speech

import "cloud.google.com/go/speech/apiv1/speechpb"

// ExtractTranscript exposes extractTranscript for black-box tests.
func (s *Service) ExtractTranscript(results []*speechpb.SpeechRecognitionResult) string {
	return s.extractTranscript(results)
}
