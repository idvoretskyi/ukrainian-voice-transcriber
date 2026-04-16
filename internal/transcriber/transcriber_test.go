// Voice Transcriber
// Copyright (c) 2025 Ihor Dvoretskyi
//
// Licensed under MIT License

package transcriber_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/idvoretskyi/voice-transcriber/internal/transcriber"
)

func TestGenerateAudioPath(t *testing.T) {
	t.Parallel()

	inputPath := "/some/dir/my video file.mp4"
	audioPath, err := transcriber.GenerateAudioPath(inputPath)
	if err != nil {
		t.Fatalf("GenerateAudioPath() unexpected error: %v", err)
	}
	defer os.Remove(audioPath)

	// Must be in the system temp dir
	if !strings.HasPrefix(audioPath, os.TempDir()) {
		t.Errorf("GenerateAudioPath() = %q; want prefix %q", audioPath, os.TempDir())
	}

	// Must end with .wav
	if filepath.Ext(audioPath) != ".wav" {
		t.Errorf("GenerateAudioPath() = %q; want .wav extension", audioPath)
	}

	// Must contain the base name (without extension) of the input
	base := strings.TrimSuffix(filepath.Base(inputPath), filepath.Ext(inputPath))
	if !strings.Contains(filepath.Base(audioPath), base) {
		t.Errorf("GenerateAudioPath() = %q; want base name %q in output", audioPath, base)
	}

	// Two calls must produce different paths (os.CreateTemp uniqueness)
	audioPath2, err := transcriber.GenerateAudioPath(inputPath)
	if err != nil {
		t.Fatalf("GenerateAudioPath() second call unexpected error: %v", err)
	}
	defer os.Remove(audioPath2)

	if audioPath == audioPath2 {
		t.Errorf("GenerateAudioPath() returned identical paths on two calls: %q", audioPath)
	}
}

func TestValidateInputPath(t *testing.T) {
	t.Parallel()

	t.Run("non-existent file returns error", func(t *testing.T) {
		t.Parallel()

		_, err := transcriber.ValidateInputPath("/non/existent/file.mp4")
		if err == nil {
			t.Error("expected error for non-existent file, got nil")
		}
	})

	t.Run("regular file is accepted", func(t *testing.T) {
		t.Parallel()

		f, err := os.CreateTemp(t.TempDir(), "test-media-*.mp4")
		if err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}

		if err := f.Close(); err != nil {
			t.Fatalf("failed to close temp file: %v", err)
		}

		got, err := transcriber.ValidateInputPath(f.Name())
		if err != nil {
			t.Errorf("unexpected error for valid file: %v", err)
		}

		if got == "" {
			t.Error("expected non-empty cleaned path")
		}
	})

	t.Run("directory is rejected", func(t *testing.T) {
		t.Parallel()

		_, err := transcriber.ValidateInputPath(t.TempDir())
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

		dirty := filepath.Join(dir, ".", filepath.Base(f.Name()))
		clean := filepath.Clean(dirty)

		got, err := transcriber.ValidateInputPath(dirty)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got != clean {
			t.Errorf("got %q; want cleaned path %q", got, clean)
		}
	})
}

func TestClassifyInputFile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		path             string
		wantType         transcriber.InputType
		wantMIMENotEmpty bool
	}{
		{path: "recording.wav", wantType: transcriber.InputTypeAudio, wantMIMENotEmpty: true},
		{path: "recording.mp3", wantType: transcriber.InputTypeAudio, wantMIMENotEmpty: true},
		{path: "recording.flac", wantType: transcriber.InputTypeAudio, wantMIMENotEmpty: true},
		{path: "recording.ogg", wantType: transcriber.InputTypeAudio, wantMIMENotEmpty: true},
		{path: "recording.m4a", wantType: transcriber.InputTypeAudio, wantMIMENotEmpty: true},
		{path: "recording.aac", wantType: transcriber.InputTypeAudio, wantMIMENotEmpty: true},
		{path: "video.mp4", wantType: transcriber.InputTypeVideo, wantMIMENotEmpty: false},
		{path: "video.mkv", wantType: transcriber.InputTypeVideo, wantMIMENotEmpty: false},
		{path: "video.mov", wantType: transcriber.InputTypeVideo, wantMIMENotEmpty: false},
		{path: "video.avi", wantType: transcriber.InputTypeVideo, wantMIMENotEmpty: false},
		// Mixed-case extension
		{path: "recording.WAV", wantType: transcriber.InputTypeAudio, wantMIMENotEmpty: true},
		{path: "video.MP4", wantType: transcriber.InputTypeVideo, wantMIMENotEmpty: false},
	}

	for _, tc := range tests {
		t.Run(tc.path, func(t *testing.T) {
			t.Parallel()

			gotType, gotMIME := transcriber.ClassifyInputFile(tc.path)

			if gotType != tc.wantType {
				t.Errorf("ClassifyInputFile(%q) type = %v; want %v", tc.path, gotType, tc.wantType)
			}

			if tc.wantMIMENotEmpty && gotMIME == "" {
				t.Errorf("ClassifyInputFile(%q) MIME = empty; want non-empty", tc.path)
			}

			if !tc.wantMIMENotEmpty && gotMIME != "" {
				t.Errorf("ClassifyInputFile(%q) MIME = %q; want empty", tc.path, gotMIME)
			}
		})
	}
}
