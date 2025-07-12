// Ukrainian Voice Transcriber
// Copyright (c) 2025 Ihor Dvoretskyi
//
// Licensed under MIT License

// Package storage provides Google Cloud Storage functionality.
package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	storageapi "cloud.google.com/go/storage"

	"github.com/idvoretskyi/ukrainian-voice-transcriber/pkg/config"
)

// Service handles Google Cloud Storage operations.
type Service struct {
	client    *storageapi.Client
	config    *config.Config
	projectID string
}

// NewService creates a new storage service.
func NewService(client *storageapi.Client, cfg *config.Config, projectID string) *Service {
	return &Service{
		client:    client,
		config:    cfg,
		projectID: projectID,
	}
}

// EnsureBucket creates bucket if it doesn't exist.
func (s *Service) EnsureBucket(ctx context.Context) error {
	bucket := s.client.Bucket(s.config.BucketName)

	// Check if bucket exists
	_, err := bucket.Attrs(ctx)
	if err == nil {
		return nil // Bucket exists
	}

	// Create bucket
	if err := bucket.Create(ctx, s.projectID, &storageapi.BucketAttrs{
		Location: "US",
		Lifecycle: storageapi.Lifecycle{
			Rules: []storageapi.LifecycleRule{
				{
					Action: storageapi.LifecycleAction{Type: "Delete"},
					Condition: storageapi.LifecycleCondition{
						AgeInDays: 1, // Auto-delete after 1 day
					},
				},
			},
		},
	}); err != nil {
		return fmt.Errorf("failed to create bucket: %v", err)
	}

	return nil
}

// UploadFile uploads file to Google Cloud Storage.
func (s *Service) UploadFile(ctx context.Context, filePath string) (string, error) {
	fileName := fmt.Sprintf("audio_%d_%s", time.Now().Unix(), filepath.Base(filePath))

	if s.config.Verbose && !s.config.Quiet {
		fmt.Printf("🔍 Uploading to GCS: %s\n", fileName)
	}

	bucket := s.client.Bucket(s.config.BucketName)
	obj := bucket.Object(fileName)

	// Open local file
	file, err := os.Open(filePath) //nolint:gosec // File path is validated by caller
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}

	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			if s.config.Verbose && !s.config.Quiet {
				fmt.Printf("🔍 Warning: Failed to close file: %v\n", closeErr)
			}
		}
	}()

	// Upload to GCS
	writer := obj.NewWriter(ctx)
	if _, err := io.Copy(writer, file); err != nil {
		return "", fmt.Errorf("failed to upload file: %v", err)
	}

	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("failed to close writer: %v", err)
	}

	gcsURI := fmt.Sprintf("gs://%s/%s", s.config.BucketName, fileName)

	if s.config.Verbose && !s.config.Quiet {
		fmt.Printf("🔍 File uploaded to: %s\n", gcsURI)
	}

	return gcsURI, nil
}

// CleanupFile removes file from Google Cloud Storage.
func (s *Service) CleanupFile(ctx context.Context, gcsURI string) {
	// Extract object name from URI
	parts := strings.Split(gcsURI, "/")
	if len(parts) < 4 {
		return
	}

	objectName := parts[len(parts)-1]

	bucket := s.client.Bucket(s.config.BucketName)
	obj := bucket.Object(objectName)

	if err := obj.Delete(ctx); err != nil {
		if s.config.Verbose && !s.config.Quiet {
			fmt.Printf("🔍 Failed to cleanup %s: %v\n", objectName, err)
		}
	} else {
		if s.config.Verbose && !s.config.Quiet {
			fmt.Printf("🔍 Cleaned up: %s\n", objectName)
		}
	}
}
