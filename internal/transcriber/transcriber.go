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

// TranscriptionResult represents the result of a transcription.
type TranscriptionResult struct {
	Text           string        `json:"text"`
	Error          string        `json:"error,omitempty"`
	ProcessingTime time.Duration `json:"processing_time,omitempty"`
	WordCount      int           `json:"word_count"`
	Success        bool          `json:"success"`
}

// Transcriber handles the main transcription logic.
type Transcriber struct {
	config        *config.Config
	geminiService *gemini.Service
}

// getProjectIDFromGcloud gets the current project ID from gcloud.
func getProjectIDFromGcloud() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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
func New(cfg *config.Config) (*Transcriber, error) {
	ctx := context.Background()

	// Resolve the GCP project ID
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		var err error

		projectID, err = getProjectIDFromGcloud()
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

// logInfo logs info messages if not in quiet mode.
func (t *Transcriber) logInfo(msg string) {
	if !t.config.Quiet {
		fmt.Printf("ℹ️  %s\n", msg)
	}
}

// logVerbose logs messages only when verbose mode is enabled and not suppressed by quiet.
func logVerbose(cfg *config.Config, format string, args ...any) {
	if cfg.Verbose && !cfg.Quiet {
		fmt.Printf("ℹ️  "+format+"\n", args...)
	}
}

// TranscribeLocalFile transcribes a local video or audio file.
// ctx controls the lifetime of the entire operation.
func (t *Transcriber) TranscribeLocalFile(ctx context.Context, inputPath string) *TranscriptionResult {
	startTime := time.Now()

	t.logInfo(fmt.Sprintf("Processing: %s", inputPath))

	// Prepare audio (extract if video, read directly if audio)
	prepared, err := prepareAudio(inputPath, t.config)
	if err != nil {
		return &TranscriptionResult{
			Success: false,
			Error:   err.Error(),
		}
	}

	defer prepared.Cleanup()

	// Send to Gemini for transcription
	transcript, err := t.geminiService.TranscribeAudio(ctx, prepared.Data, prepared.MIMEType)
	if err != nil {
		return &TranscriptionResult{
			Success: false,
			Error:   err.Error(),
		}
	}

	processingTime := time.Since(startTime)
	wordCount := len(strings.Fields(transcript))

	return &TranscriptionResult{
		Text:           transcript,
		Success:        true,
		WordCount:      wordCount,
		ProcessingTime: processingTime,
	}
}
