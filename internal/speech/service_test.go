// Ukrainian Voice Transcriber
// Copyright (c) 2025 Ihor Dvoretskyi
//
// Licensed under MIT License

package speech_test

import (
	"testing"

	"cloud.google.com/go/speech/apiv1/speechpb"

	"github.com/idvoretskyi/ukrainian-voice-transcriber/internal/speech"
)

func TestExtractTranscript(t *testing.T) {
	t.Parallel()

	svc := &speech.Service{} // ExtractTranscript has no side effects; no config needed

	tests := []struct {
		name     string
		results  []*speechpb.SpeechRecognitionResult
		expected string
	}{
		{
			name:     "nil results returns empty string",
			results:  nil,
			expected: "",
		},
		{
			name:     "empty results returns empty string",
			results:  []*speechpb.SpeechRecognitionResult{},
			expected: "",
		},
		{
			name: "single result",
			results: []*speechpb.SpeechRecognitionResult{
				{
					Alternatives: []*speechpb.SpeechRecognitionAlternative{
						{Transcript: "hello world"},
					},
				},
			},
			expected: "hello world",
		},
		{
			name: "multiple results are joined",
			results: []*speechpb.SpeechRecognitionResult{
				{
					Alternatives: []*speechpb.SpeechRecognitionAlternative{
						{Transcript: "first segment"},
					},
				},
				{
					Alternatives: []*speechpb.SpeechRecognitionAlternative{
						{Transcript: "second segment"},
					},
				},
			},
			expected: "first segment second segment",
		},
		{
			name: "result with no alternatives is skipped",
			results: []*speechpb.SpeechRecognitionResult{
				{Alternatives: nil},
				{
					Alternatives: []*speechpb.SpeechRecognitionAlternative{
						{Transcript: "only this"},
					},
				},
			},
			expected: "only this",
		},
		{
			name: "Ukrainian text preserved",
			results: []*speechpb.SpeechRecognitionResult{
				{
					Alternatives: []*speechpb.SpeechRecognitionAlternative{
						{Transcript: "Привіт світ"},
					},
				},
			},
			expected: "Привіт світ",
		},
		{
			name: "only first alternative per result is used",
			results: []*speechpb.SpeechRecognitionResult{
				{
					Alternatives: []*speechpb.SpeechRecognitionAlternative{
						{Transcript: "best alternative"},
						{Transcript: "second choice"},
					},
				},
			},
			expected: "best alternative",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := svc.ExtractTranscript(tc.results)
			if got != tc.expected {
				t.Errorf("extractTranscript() = %q; want %q", got, tc.expected)
			}
		})
	}
}
