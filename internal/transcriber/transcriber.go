// Voice Transcriber
// Copyright (c) 2025 Ihor Dvoretskyi
//
// Licensed under MIT License

// Package transcriber provides media transcription functionality.
package transcriber

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/idvoretskyi/voice-transcriber/internal/config"
	"github.com/idvoretskyi/voice-transcriber/internal/gemini"
)

const gcloudTimeout = 10 * time.Second

// TranscriptionResult holds the output of a successful transcription.
// On failure, TranscribeLocalFile returns a non-nil error instead.
type TranscriptionResult struct {
	Text           string
	ProcessingTime time.Duration
	WordCount      int
}

// projectIDResolver is the function type used to obtain a GCP project ID
// at runtime. The default implementation calls gcloud; tests can inject a stub.
type projectIDResolver func(ctx context.Context) (string, error)

// Transcriber handles the main transcription logic.
type Transcriber struct {
	config    *config.Config
	backend   gemini.AudioTranscriber
	logger    *slog.Logger
	resolveID projectIDResolver
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
		return "", fmt.Errorf("failed to get project ID from gcloud: %w (stderr: %s)", err, stderr.String())
	}

	projectID := strings.TrimSpace(stdout.String())
	if projectID == "" {
		return "", fmt.Errorf("no project ID configured in gcloud; run: gcloud config set project PROJECT_ID")
	}

	return projectID, nil
}

// New creates a new Transcriber instance, resolving the GCP project ID and
// initializing the Gemini service.
// If logger is nil, slog.Default() is used.
func New(ctx context.Context, cfg *config.Config, logger *slog.Logger) (*Transcriber, error) {
	if logger == nil {
		logger = slog.Default()
	}

	t := &Transcriber{
		config:    cfg,
		logger:    logger,
		resolveID: getProjectIDFromGcloud,
	}

	// Prefer GCPProject already on the config (e.g. from FromEnv), then env
	// var, then gcloud CLI (via the injected resolver).
	projectID := cfg.GCPProject
	if projectID == "" {
		projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")
	}

	if projectID == "" {
		var err error

		projectID, err = t.resolveID(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve GCP project ID: %w\n\n"+
				"Set it with one of:\n"+
				"  export GOOGLE_CLOUD_PROJECT=your-project-id\n"+
				"  gcloud config set project your-project-id", err)
		}
	}

	logger.DebugContext(ctx, "resolved GCP project", slog.String("project", projectID))

	backend, err := gemini.NewService(ctx, cfg, projectID, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Gemini service: %w", err)
	}

	t.backend = backend

	return t, nil
}

// TranscribeLocalFile transcribes a local video or audio file.
// ctx controls the lifetime of the entire operation.
// It returns a *TranscriptionResult on success, or a non-nil error on failure.
func (t *Transcriber) TranscribeLocalFile(ctx context.Context, inputPath string) (*TranscriptionResult, error) {
	startTime := time.Now()

	t.logger.InfoContext(ctx, "processing file", slog.String("path", inputPath))

	prepared, err := prepareAudio(ctx, inputPath, t.logger)
	if err != nil {
		return nil, fmt.Errorf("preparing audio: %w", err)
	}

	defer func() {
		if closeErr := prepared.Close(); closeErr != nil {
			t.logger.WarnContext(ctx, "failed to remove temp audio file", slog.Any("error", closeErr))
		}
	}()

	transcript, err := t.backend.TranscribeAudio(ctx, prepared.Data, prepared.MIMEType)
	if err != nil {
		return nil, fmt.Errorf("transcribing audio: %w", err)
	}

	return &TranscriptionResult{
		Text:           transcript,
		WordCount:      len(strings.Fields(transcript)),
		ProcessingTime: time.Since(startTime),
	}, nil
}
