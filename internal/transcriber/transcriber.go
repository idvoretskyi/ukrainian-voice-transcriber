// Ukrainian Voice Transcriber
// Copyright (c) {{ YEAR }} Ihor Dvoretskyi
//
// Licensed under MIT License

// Package transcriber provides video transcription functionality.
package transcriber

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	speechapi "cloud.google.com/go/speech/apiv1"
	storageapi "cloud.google.com/go/storage"
	"google.golang.org/api/option"

	"github.com/idvoretskyi/ukrainian-voice-transcriber/internal/speech"
	"github.com/idvoretskyi/ukrainian-voice-transcriber/internal/storage"
	"github.com/idvoretskyi/ukrainian-voice-transcriber/pkg/config"
)

// TranscriptionResult represents the result of a transcription.
type TranscriptionResult struct {
	Text           string        `json:"text"`
	Success        bool          `json:"success"`
	Error          string        `json:"error,omitempty"`
	Duration       time.Duration `json:"duration,omitempty"`
	WordCount      int           `json:"word_count"`
	ProcessingTime time.Duration `json:"processing_time,omitempty"`
}

// Transcriber handles the main transcription logic.
type Transcriber struct {
	config         *config.Config
	speechClient   *speechapi.Client
	storageClient  *storageapi.Client
	speechService  *speech.Service
	storageService *storage.Service
	ctx            context.Context
}

// getProjectIDFromGcloud gets the current project ID from gcloud.
func getProjectIDFromGcloud() (string, error) {
	// Set timeout to prevent hanging
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "gcloud", "config", "get-value", "project")

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get project ID from gcloud: %v", err)
	}

	projectID := strings.TrimSpace(string(output))
	if projectID == "" {
		return "", fmt.Errorf("no project ID configured in gcloud. Run: gcloud config set project PROJECT_ID")
	}

	return projectID, nil
}

// New creates a new transcriber instance.
func New(cfg *config.Config) (*Transcriber, error) {
	ctx := context.Background()

	var speechClient *speechapi.Client

	var storageClient *storageapi.Client

	var err error

	var projectID string

	// Try Application Default Credentials (works with gcloud auth)
	speechClient, err = speechapi.NewClient(ctx)
	if err != nil {
		// Fall back to service account
		serviceAccountPath := config.FindServiceAccount()
		if serviceAccountPath == "" {
			return nil, fmt.Errorf(`authentication required. Choose one option:

1. Use gcloud (Recommended):
   gcloud auth login
   gcloud auth application-default login

2. Service Account:
   Place service-account.json in current directory

3. OAuth setup:
   ukrainian-voice-transcriber auth`)
		}

		cfg.ServiceAccountPath = serviceAccountPath

		// Initialize Google Cloud clients with service account
		speechClient, err = speechapi.NewClient(ctx, option.WithCredentialsFile(serviceAccountPath))
		if err != nil {
			return nil, fmt.Errorf("failed to create speech client: %v", err)
		}

		storageClient, err = storageapi.NewClient(ctx, option.WithCredentialsFile(serviceAccountPath))
		if err != nil {
			return nil, fmt.Errorf("failed to create storage client: %v", err)
		}

		if !cfg.Quiet {
			fmt.Printf("🔑 Using service account authentication\n")
		}

		// For service account, we need to extract project ID from the service account file
		// or use environment variable
		projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")
		if projectID == "" {
			projectID = "ukrainian-voice-transcriber" // Default project name
		}
	} else {
		// Use Application Default Credentials (gcloud)
		storageClient, err = storageapi.NewClient(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to create storage client: %v", err)
		}

		// Get project ID from gcloud
		projectID, err = getProjectIDFromGcloud()
		if err != nil {
			return nil, fmt.Errorf("failed to get project ID: %v", err)
		}

		if !cfg.Quiet {
			fmt.Printf("🔐 Using Application Default Credentials (gcloud)\n")
			fmt.Printf("📊 Project: %s\n", projectID)
		}
	}

	// Set default bucket name if not provided
	if cfg.BucketName == "" {
		if projectID != "" {
			cfg.BucketName = fmt.Sprintf("%s-ukr-voice-transcriber", projectID)
		} else {
			cfg.BucketName = "ukr-voice-transcriber-temp"
		}
	}

	// Initialize services
	speechService := speech.NewService(speechClient, cfg)
	storageService := storage.NewService(storageClient, cfg, projectID)

	// Ensure bucket exists
	if err := storageService.EnsureBucket(ctx); err != nil {
		return nil, fmt.Errorf("failed to ensure bucket exists: %v", err)
	}

	transcriber := &Transcriber{
		config:         cfg,
		speechClient:   speechClient,
		storageClient:  storageClient,
		speechService:  speechService,
		storageService: storageService,
		ctx:            ctx,
	}

	// Check for Drive credentials (for future use)
	driveCredPath := config.FindDriveCredentials()
	if driveCredPath != "" {
		cfg.DriveCredentials = driveCredPath
		transcriber.logInfo("Google Drive credentials found (ready for future Drive support)")
	}

	return transcriber, nil
}

// logInfo logs info messages if not in quiet mode.
func (t *Transcriber) logInfo(msg string) {
	if !t.config.Quiet {
		fmt.Printf("ℹ️  %s\n", msg)
	}
}

// TranscribeLocalFile transcribes a local video file.
func (t *Transcriber) TranscribeLocalFile(videoPath string) *TranscriptionResult {
	startTime := time.Now()

	// Check if file exists
	if _, err := os.Stat(videoPath); os.IsNotExist(err) {
		return &TranscriptionResult{
			Success: false,
			Error:   fmt.Sprintf("File not found: %s", videoPath),
		}
	}

	t.logInfo(fmt.Sprintf("Processing: %s", videoPath))

	// Extract audio
	audioPath, err := extractAudio(videoPath, t.config)
	if err != nil {
		return &TranscriptionResult{
			Success: false,
			Error:   fmt.Sprintf("Audio extraction failed: %v", err),
		}
	}
	defer os.Remove(audioPath) // Cleanup local audio file

	// Upload to storage
	gcsURI, err := t.storageService.UploadFile(t.ctx, audioPath)
	if err != nil {
		return &TranscriptionResult{
			Success: false,
			Error:   fmt.Sprintf("Storage upload failed: %v", err),
		}
	}
	defer t.storageService.CleanupFile(t.ctx, gcsURI) // Cleanup cloud storage

	// Transcribe
	transcript, err := t.speechService.TranscribeAudio(t.ctx, gcsURI)
	if err != nil {
		return &TranscriptionResult{
			Success: false,
			Error:   fmt.Sprintf("Transcription failed: %v", err),
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

// Close closes all clients.
func (t *Transcriber) Close() error {
	if t.speechClient != nil {
		t.speechClient.Close()
	}
	if t.storageClient != nil {
		t.storageClient.Close()
	}

	return nil
}
