// Voice Transcriber
// Copyright (c) 2025 Ihor Dvoretskyi
//
// Licensed under MIT License

// Package transcriber provides media transcription functionality.
package transcriber

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/idvoretskyi/voice-transcriber/internal/gemini"
	"github.com/idvoretskyi/voice-transcriber/pkg/config"
)

const gcloudTimeout = 10 * time.Second

// TranscriptionResult represents the result of a transcription.
type TranscriptionResult struct {
	Text           string
	Error          string
	ProcessingTime time.Duration
	WordCount      int
	Success        bool
}

// Transcriber handles the main transcription logic.
type Transcriber struct {
	config        *config.Config
	geminiService *gemini.Service
}

// getProjectIDFromGcloud gets the current project ID from gcloud.
func getProjectIDFromGcloud(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, gcloudTimeout)
	defer cancel()

	gcloudPath, err := exec.LookPath("gcloud")
	if err != nil {
		return "", fmt.Errorf("gcloud command not found: %w", err)
	}

	cmd := exec.CommandContext(ctx, gcloudPath, "config", "get-value", "project") // #nosec G204 gcloudPath is validated
	cmd.Env = os.Environ()

	var stdout, stderr strings.Builder

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to get project ID from gcloud: %v, stderr: %s", err, stderr.String())
	}

	projectID := strings.TrimSpace(stdout.String())
	if projectID == "" {
		return "", fmt.Errorf("no project ID configured in gcloud. Run: gcloud config set project PROJECT_ID")
	}

	return projectID, nil
}

// New creates a new Transcriber instance, initializing the Gemini service.
func New(ctx context.Context, cfg *config.Config) (*Transcriber, error) {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		var err error

		projectID, err = getProjectIDFromGcloud(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve GCP project ID: %w\n\n"+
				"Set it with one of:\n"+
				"  export GOOGLE_CLOUD_PROJECT=your-project-id\n"+
				"  gcloud config set project your-project-id", err)
		}
	}

	if !cfg.Quiet {
		fmt.Printf("ℹ️  Project: %s\n", projectID)
	}

	geminiService, err := gemini.NewService(ctx, cfg, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Gemini service: %w", err)
	}

	return &Transcriber{
		config:        cfg,
		geminiService: geminiService,
	}, nil
}

// logVerbose logs a message only when verbose mode is enabled and quiet is not set.
func logVerbose(cfg *config.Config, format string, args ...any) {
	if cfg.Verbose && !cfg.Quiet {
		fmt.Printf("ℹ️  "+format+"\n", args...)
	}
}

// TranscribeLocalFile transcribes a local video or audio file.
// ctx controls the lifetime of the entire operation.
func (t *Transcriber) TranscribeLocalFile(ctx context.Context, inputPath string) *TranscriptionResult {
	startTime := time.Now()

	if !t.config.Quiet {
		fmt.Printf("ℹ️  Processing: %s\n", inputPath)
	}

	prepared, err := prepareAudio(ctx, inputPath, t.config)
	if err != nil {
		return &TranscriptionResult{
			Success: false,
			Error:   err.Error(),
		}
	}

	defer prepared.Cleanup()

	transcript, err := t.geminiService.TranscribeAudio(ctx, prepared.Data, prepared.MIMEType)
	if err != nil {
		return &TranscriptionResult{
			Success: false,
			Error:   err.Error(),
		}
	}

	return &TranscriptionResult{
		Text:           transcript,
		Success:        true,
		WordCount:      len(strings.Fields(transcript)),
		ProcessingTime: time.Since(startTime),
	}
}
