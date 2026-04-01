// Ukrainian Voice Transcriber
// Copyright (c) 2025 Ihor Dvoretskyi
//
// Licensed under MIT License

package transcriber_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/idvoretskyi/ukrainian-voice-transcriber/internal/transcriber"
	"github.com/idvoretskyi/ukrainian-voice-transcriber/pkg/config"
)

func TestSetBucketName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		initialBucket  string
		projectID      string
		expectedBucket string
	}{
		{
			name:           "sets bucket from project ID when empty",
			initialBucket:  "",
			projectID:      "my-project",
			expectedBucket: "my-project-voice-transcriber-data",
		},
		{
			name:           "does not override existing bucket name",
			initialBucket:  "my-existing-bucket",
			projectID:      "my-project",
			expectedBucket: "my-existing-bucket",
		},
		{
			name:           "no-op when project ID is empty",
			initialBucket:  "",
			projectID:      "",
			expectedBucket: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			cfg := &config.Config{BucketName: tc.initialBucket}
			transcriber.SetBucketName(cfg, tc.projectID)

			if cfg.BucketName != tc.expectedBucket {
				t.Errorf("transcriber.SetBucketName() BucketName = %q; want %q", cfg.BucketName, tc.expectedBucket)
			}
		})
	}
}

func TestGenerateAudioPath(t *testing.T) {
	t.Parallel()

	videoPath := "/some/dir/my video file.mp4"
	audioPath := transcriber.GenerateAudioPath(videoPath)

	// Must be in the system temp dir
	if !strings.HasPrefix(audioPath, os.TempDir()) {
		t.Errorf("transcriber.GenerateAudioPath() = %q; want prefix %q", audioPath, os.TempDir())
	}

	// Must end with .wav
	if filepath.Ext(audioPath) != ".wav" {
		t.Errorf("transcriber.GenerateAudioPath() = %q; want .wav extension", audioPath)
	}

	// Must contain the base name (without extension) of the input
	base := strings.TrimSuffix(filepath.Base(videoPath), filepath.Ext(videoPath))
	if !strings.Contains(filepath.Base(audioPath), base) {
		t.Errorf("transcriber.GenerateAudioPath() = %q; want base name %q in output", audioPath, base)
	}

	// Two calls must produce different paths (timestamp-based uniqueness)
	audioPath2 := transcriber.GenerateAudioPath(videoPath)
	if audioPath == audioPath2 {
		t.Errorf("transcriber.GenerateAudioPath() returned identical paths on two calls: %q", audioPath)
	}
}

func TestValidateAndSanitizeVideoPath(t *testing.T) {
	t.Parallel()

	t.Run("non-existent file returns error", func(t *testing.T) {
		t.Parallel()

		_, err := transcriber.ValidateAndSanitizeVideoPath("/non/existent/file.mp4")
		if err == nil {
			t.Error("expected error for non-existent file, got nil")
		}
	})

	t.Run("regular file is accepted", func(t *testing.T) {
		t.Parallel()

		f, err := os.CreateTemp(t.TempDir(), "test-video-*.mp4")
		if err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}

		if err := f.Close(); err != nil {
			t.Fatalf("failed to close temp file: %v", err)
		}

		got, err := transcriber.ValidateAndSanitizeVideoPath(f.Name())
		if err != nil {
			t.Errorf("unexpected error for valid file: %v", err)
		}

		if got == "" {
			t.Error("expected non-empty cleaned path")
		}
	})

	t.Run("directory is rejected", func(t *testing.T) {
		t.Parallel()

		_, err := transcriber.ValidateAndSanitizeVideoPath(t.TempDir())
		if err == nil {
			t.Error("expected error for directory input, got nil")
		}
	})

	t.Run("path is cleaned", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()

		f, err := os.CreateTemp(dir, "test-*.mp4")
		if err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}

		if err := f.Close(); err != nil {
			t.Fatalf("failed to close temp file: %v", err)
		}

		// Provide a path with redundant slashes / dots
		dirty := filepath.Join(dir, ".", filepath.Base(f.Name()))
		clean := filepath.Clean(dirty)

		got, err := transcriber.ValidateAndSanitizeVideoPath(dirty)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got != clean {
			t.Errorf("got %q; want cleaned path %q", got, clean)
		}
	})
}
