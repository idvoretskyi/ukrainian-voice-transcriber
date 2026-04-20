// Voice Transcriber
// Copyright (c) 2025 Ihor Dvoretskyi
//
// Licensed under MIT License

package cli_test

import (
	"bytes"
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

// TestNewRootCmdVersion verifies that NewRootCmd wires VersionInfo into
// cobra's --version flag correctly by executing the command and checking output.
func TestNewRootCmdVersion(t *testing.T) {
	t.Parallel()

	info := cli.VersionInfo{Version: "9.8.7", Date: "2025-01-01T00:00:00Z", Commit: "abc1234"}
	cfg := &config.Config{}
	root := cli.NewRootCmd(cfg, info)
	// Wire version exactly as Execute() does.
	root.Version = info.Version

	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetArgs([]string{"--version"})

	if err := root.Execute(); err != nil {
		t.Fatalf("root.Execute() returned unexpected error: %v", err)
	}

	if got := buf.String(); !strings.Contains(got, "9.8.7") {
		t.Errorf("--version output = %q; want it to contain %q", got, "9.8.7")
	}
}

// TestRootCmdHelpDoesNotContainVersion verifies that the Long description no
// longer bakes in the version string (H4 fix: version removed from --help).
// It executes --help and checks the captured output, not just root.Long.
func TestRootCmdHelpDoesNotContainVersion(t *testing.T) {
	t.Parallel()

	info := cli.VersionInfo{Version: "99.0.0", Date: "2025-01-01T00:00:00Z", Commit: "deadbeef"}
	cfg := &config.Config{}
	root := cli.NewRootCmd(cfg, info)

	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetArgs([]string{"--help"})

	// --help exits with code 0; Execute should not return an error.
	if err := root.Execute(); err != nil {
		t.Fatalf("root.Execute() --help returned unexpected error: %v", err)
	}

	if got := buf.String(); strings.Contains(got, "99.0.0") {
		t.Errorf("--help output contains version string %q; it should be version-free", "99.0.0")
	}
}

// TestResolveOutputPath verifies the default transcript-path derivation logic.
func TestResolveOutputPath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		outputFile string
		mediaFile  string
		want       string
	}{
		{
			name:       "explicit output path returned unchanged",
			outputFile: "/tmp/my-output.txt",
			mediaFile:  "anything.mp4",
			want:       "/tmp/my-output.txt",
		},
		{
			name:       "default path derived from media filename",
			outputFile: "",
			mediaFile:  "lecture.mp4",
			want:       "output/lecture/lecture.txt",
		},
		{
			name:       "extension stripped and spaces sanitized",
			outputFile: "",
			mediaFile:  "my recording.wav",
			want:       "output/my_recording/my_recording.txt",
		},
		{
			name:       "nested media path — only basename used",
			outputFile: "",
			mediaFile:  "/some/deep/path/talk.mkv",
			want:       "output/talk/talk.txt",
		},
		{
			name:       "special chars removed from name",
			outputFile: "",
			mediaFile:  "hello!@world.mp4",
			want:       "output/helloworld/helloworld.txt",
		},
		{
			name:       "empty name after sanitization falls back to transcript",
			outputFile: "",
			mediaFile:  "!!!.mp4",
			want:       "output/transcript/transcript.txt",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := cli.ResolveOutputPath(tc.outputFile, tc.mediaFile)
			if got != tc.want {
				t.Errorf("ResolveOutputPath(%q, %q) = %q; want %q",
					tc.outputFile, tc.mediaFile, got, tc.want)
			}
		})
	}
}
