// Ukrainian Voice Transcriber
// Copyright (c) 2025 Ihor Dvoretskyi
//
// Licensed under MIT License

// Package transcriber provides video transcription functionality.
package transcriber

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
	Error          string        `json:"error,omitempty"`
	Duration       time.Duration `json:"duration,omitempty"`
	ProcessingTime time.Duration `json:"processing_time,omitempty"`
	WordCount      int           `json:"word_count"`
	Success        bool          `json:"success"`
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

	// Use a fixed path for gcloud to avoid path traversal
	gcloudPath, err := exec.LookPath("gcloud")
	if err != nil {
		return "", fmt.Errorf("gcloud command not found: %v", err)
	}

	// Use only fixed arguments to avoid command injection
	cmd := exec.CommandContext(ctx, gcloudPath, "config", "get-value", "project") // #nosec G204 gcloudPath is validated

	// Restrict command execution environment
	cmd.Env = []string{"PATH=" + filepath.Dir(gcloudPath)}

	// Capture both stdout and stderr
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

// New creates a new transcriber instance.
func New(cfg *config.Config) (*Transcriber, error) {
	ctx := context.Background()

	// Initialize Google Cloud clients
	speechClient, storageClient, projectID, err := initializeClients(ctx, cfg)
	if err != nil {
		return nil, err
	}

	// Set default bucket name if not provided
	setBucketName(cfg, projectID)

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
	configureDriveCredentials(cfg, transcriber)

	return transcriber, nil
}

// initializeClients initializes Google Cloud clients and returns them with project ID.
func initializeClients(ctx context.Context, cfg *config.Config) (
	speechClient *speechapi.Client, storageClient *storageapi.Client, projectID string, err error,
) {
	// Try Application Default Credentials (works with gcloud auth)
	speechClient, err = speechapi.NewClient(ctx)
	if err != nil {
		return initializeWithServiceAccount(ctx, cfg)
	}

	return initializeWithDefaultCredentials(ctx, cfg, speechClient)
}

// initializeWithServiceAccount initializes clients using service account authentication.
func initializeWithServiceAccount(ctx context.Context, cfg *config.Config) (
	speechClient *speechapi.Client, storageClient *storageapi.Client, projectID string, err error,
) {
	serviceAccountPath := config.FindServiceAccount()
	if serviceAccountPath == "" {
		return nil, nil, "", fmt.Errorf("authentication required. Choose one option:\n\n" +
			"1. Use gcloud (Recommended):\n   gcloud auth login\n   gcloud auth application-default login\n\n" +
			"2. Service Account:\n   Place service-account.json in current directory\n\n" +
			"3. OAuth setup:\n   ukrainian-voice-transcriber auth")
	}

	cfg.ServiceAccountPath = serviceAccountPath

	// Initialize Google Cloud clients with service account
	speechClient, err = speechapi.NewClient(ctx, option.WithCredentialsFile(serviceAccountPath))
	if err != nil {
		return nil, nil, "", fmt.Errorf("failed to create speech client: %v", err)
	}

	storageClient, err = storageapi.NewClient(ctx, option.WithCredentialsFile(serviceAccountPath))
	if err != nil {
		return nil, nil, "", fmt.Errorf("failed to create storage client: %v", err)
	}

	if !cfg.Quiet {
		fmt.Printf("🔑 Using service account authentication\n")
	}

	// For service account, we need to extract project ID from the service account file
	// or use environment variable
	projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		// Try to get project ID from service account file (would require parsing JSON)
		// For now, return an error if project ID is not set
		return nil, nil, "", fmt.Errorf("project ID not set. Please set GOOGLE_CLOUD_PROJECT environment variable")
	}

	return speechClient, storageClient, projectID, nil
}

// initializeWithDefaultCredentials initializes clients using Application Default Credentials.
func initializeWithDefaultCredentials(
	ctx context.Context, cfg *config.Config, speechClient *speechapi.Client,
) (retSpeechClient *speechapi.Client, storageClient *storageapi.Client, projectID string, err error) {
	// Use Application Default Credentials (gcloud)
	storageClient, err = storageapi.NewClient(ctx)
	if err != nil {
		return nil, nil, "", fmt.Errorf("failed to create storage client: %v", err)
	}

	// Get project ID from gcloud
	projectID, err = getProjectIDFromGcloud()
	if err != nil {
		return nil, nil, "", fmt.Errorf("failed to get project ID: %v", err)
	}

	if !cfg.Quiet {
		fmt.Printf("🔐 Using Application Default Credentials (gcloud)\n")
		fmt.Printf("📊 Project: %s\n", projectID)
	}

	return speechClient, storageClient, projectID, nil
}

// setBucketName sets the default bucket name if not provided.
// Note: If both BucketName and projectID are empty, this will be caught later
// when trying to ensure the bucket exists.
func setBucketName(cfg *config.Config, projectID string) {
	if cfg.BucketName == "" && projectID != "" {
		// Create a bucket name based on project ID with a consistent suffix
		// This avoids hardcoding a default bucket name
		cfg.BucketName = fmt.Sprintf("%s-voice-transcriber-data", projectID)
	}
}

// configureDriveCredentials checks for Drive credentials and configures them.
func configureDriveCredentials(cfg *config.Config, transcriber *Transcriber) {
	driveCredPath := config.FindDriveCredentials()
	if driveCredPath != "" {
		cfg.DriveCredentials = driveCredPath

		transcriber.logInfo("Google Drive credentials found (ready for future Drive support)")
	}
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

	// Validate input file
	if err := t.validateInputFile(videoPath); err != nil {
		return &TranscriptionResult{
			Success: false,
			Error:   err.Error(),
		}
	}

	t.logInfo(fmt.Sprintf("Processing: %s", videoPath))

	// Create a context with timeout for the entire operation (30 minutes)
	ctx, cancel := context.WithTimeout(t.ctx, 30*time.Minute)
	defer cancel()

	// Process audio (extract and upload)
	gcsURI, audioCleanup, storageCleanup, err := t.processAudio(ctx, videoPath)
	if err != nil {
		return &TranscriptionResult{
			Success: false,
			Error:   err.Error(),
		}
	}
	defer audioCleanup()
	defer storageCleanup()

	// Perform transcription
	transcript, err := t.performTranscription(ctx, gcsURI)
	if err != nil {
		return &TranscriptionResult{
			Success: false,
			Error:   err.Error(),
		}
	}

	// Calculate results
	processingTime := time.Since(startTime)
	wordCount := len(strings.Fields(transcript))

	return &TranscriptionResult{
		Text:           transcript,
		Success:        true,
		WordCount:      wordCount,
		ProcessingTime: processingTime,
	}
}

// validateInputFile checks if the input file exists and is accessible.
func (t *Transcriber) validateInputFile(videoPath string) error {
	if _, err := os.Stat(videoPath); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", videoPath)
	}

	return nil
}

// processAudio extracts audio and uploads it to storage, returning cleanup functions.
func (t *Transcriber) processAudio(ctx context.Context, videoPath string) (string, func(), func(), error) {
	// Extract audio
	audioPath, err := extractAudio(videoPath, t.config)
	if err != nil {
		return "", nil, nil, fmt.Errorf("audio extraction failed: %v", err)
	}

	// Setup audio cleanup
	audioCleanup := func() {
		if removeErr := os.Remove(audioPath); removeErr != nil && !os.IsNotExist(removeErr) {
			t.logInfo(fmt.Sprintf("Warning: Failed to remove temporary audio file: %v", removeErr))
		}
	}

	// Create a context with timeout for upload (5 minutes)
	uploadCtx, uploadCancel := context.WithTimeout(ctx, 5*time.Minute)
	defer uploadCancel()

	// Upload to storage
	gcsURI, err := t.storageService.UploadFile(uploadCtx, audioPath)
	if err != nil {
		audioCleanup() // Clean up audio file on upload failure

		return "", nil, nil, fmt.Errorf("storage upload failed: %v", err)
	}

	// Setup storage cleanup
	storageCleanup := func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cleanupCancel()

		t.storageService.CleanupFile(cleanupCtx, gcsURI)
	}

	return gcsURI, audioCleanup, storageCleanup, nil
}

// performTranscription executes the transcription process.
func (t *Transcriber) performTranscription(ctx context.Context, gcsURI string) (string, error) {
	// Create a context with timeout for transcription (20 minutes)
	transcribeCtx, transcribeCancel := context.WithTimeout(ctx, 20*time.Minute)
	defer transcribeCancel()

	// Transcribe
	transcript, err := t.speechService.TranscribeAudio(transcribeCtx, gcsURI)
	if err != nil {
		return "", fmt.Errorf("transcription failed: %v", err)
	}

	return transcript, nil
}

// Close closes all clients.
func (t *Transcriber) Close() error {
	var errs []string

	if t.speechClient != nil {
		if err := t.speechClient.Close(); err != nil {
			errs = append(errs, fmt.Sprintf("speech client close error: %v", err))
		}
	}

	if t.storageClient != nil {
		if err := t.storageClient.Close(); err != nil {
			errs = append(errs, fmt.Sprintf("storage client close error: %v", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("close errors: %s", strings.Join(errs, "; "))
	}

	return nil
}
