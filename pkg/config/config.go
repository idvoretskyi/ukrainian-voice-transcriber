// Ukrainian Voice Transcriber
// Copyright (c) {{ YEAR }} Ihor Dvoretskyi
//
// Licensed under MIT License

// Package config provides configuration structures and utilities.
package config

import (
	"os"
)

// Config holds application configuration.
type Config struct {
	ServiceAccountPath string
	DriveCredentials   string
	BucketName         string
	Verbose            bool
	Quiet              bool
}

// FindServiceAccount looks for Google Cloud service account key.
func FindServiceAccount() string {
	candidates := []string{
		"service-account.json",
		"service_account.json",
		"gcloud-key.json",
		"key.json",
	}

	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}

	return ""
}

// FindDriveCredentials looks for Google Drive OAuth credentials.
func FindDriveCredentials() string {
	candidates := []string{
		"credentials.json",
		"drive_credentials.json",
		"client_secret.json",
	}

	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}

	return ""
}
