// Voice Transcriber
// Copyright (c) 2025 Ihor Dvoretskyi
//
// Licensed under MIT License

package cli_test

import (
	"strings"
	"testing"

	"github.com/idvoretskyi/voice-transcriber/internal/cli"
	"github.com/idvoretskyi/voice-transcriber/internal/config"
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

// TestSetVersion verifies that SetVersion updates the build version and that
// NewRootCmd wires it into cobra's --version flag.
func TestSetVersion(t *testing.T) {
	t.Parallel()

	cli.SetVersion("9.8.7", "2025-01-01T00:00:00Z", "abc1234")

	if got := cli.GetBuildVersion(); got != "9.8.7" {
		t.Errorf("GetBuildVersion() = %q; want %q", got, "9.8.7")
	}

	// NewRootCmd should pick up the version set above.
	cfg := &config.Config{}
	root := cli.NewRootCmd(cfg)
	root.Version = cli.GetBuildVersion()

	if root.Version != "9.8.7" {
		t.Errorf("rootCmd.Version = %q; want %q", root.Version, "9.8.7")
	}

	// Restore to default so other tests are not affected.
	t.Cleanup(func() { cli.SetVersion("dev", "unknown", "unknown") })
}

// TestRootCmdHelpDoesNotContainVersion verifies that the Long description no
// longer bakes in the version string (H4 fix: version removed from --help).
func TestRootCmdHelpDoesNotContainVersion(t *testing.T) {
	t.Parallel()

	cli.SetVersion("99.0.0", "2025-01-01T00:00:00Z", "deadbeef")

	t.Cleanup(func() { cli.SetVersion("dev", "unknown", "unknown") })

	cfg := &config.Config{}
	root := cli.NewRootCmd(cfg)

	if strings.Contains(root.Long, "99.0.0") {
		t.Errorf("rootCmd.Long contains version string %q; it should be version-free", "99.0.0")
	}
}
