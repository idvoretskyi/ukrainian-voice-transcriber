// Voice Transcriber
// Copyright (c) 2025 Ihor Dvoretskyi
//
// Licensed under MIT License

package cli_test

import (
	"testing"

	"github.com/idvoretskyi/voice-transcriber/internal/cli"
)

func TestSanitizeFilename(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "plain ASCII",
			input:    "hello",
			expected: "hello",
		},
		{
			name:     "spaces replaced by underscores",
			input:    "hello world",
			expected: "hello_world",
		},
		{
			name:     "Ukrainian Cyrillic preserved",
			input:    "привіт світ",
			expected: "привіт_світ",
		},
		{
			name:     "mixed ASCII and Cyrillic",
			input:    "video урок 1",
			expected: "video_урок_1",
		},
		{
			name:     "special characters removed",
			input:    "hello!@#world",
			expected: "helloworld",
		},
		{
			name:     "hyphens and dots preserved",
			input:    "file-name.ext",
			expected: "file-name.ext",
		},
		{
			name:     "underscores preserved",
			input:    "file_name",
			expected: "file_name",
		},
		{
			name:     "empty string returns transcript",
			input:    "",
			expected: "transcript",
		},
		{
			name:     "only special chars returns transcript",
			input:    "!!!@@@",
			expected: "transcript",
		},
		{
			name:     "numbers preserved",
			input:    "file123",
			expected: "file123",
		},
		{
			name:     "multiple spaces collapsed to underscores",
			input:    "a b c",
			expected: "a_b_c",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := cli.SanitizeFilename(tc.input)
			if got != tc.expected {
				t.Errorf("sanitizeFilename(%q) = %q; want %q", tc.input, got, tc.expected)
			}
		})
	}
}
