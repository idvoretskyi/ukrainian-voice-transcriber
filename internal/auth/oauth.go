// Ukrainian Voice Transcriber
//
// Copyright (c) 2025 Ihor Dvoretskyi
//
// Licensed under MIT License

// Package auth provides OAuth authentication functionality.
package auth

import (
	"context"
	"fmt"
)

// OAuthManager handles Google OAuth authentication.
type OAuthManager struct {
	// projectID is not currently used but kept for future expansion
	projectID string //nolint:unused // This field is reserved for future use
}

// NewOAuthManager creates a new OAuth manager with a simple approach.
func NewOAuthManager() *OAuthManager {
	// Use a simple approach: just recommend gcloud auth
	// Most users already have gcloud installed
	return &OAuthManager{}
}

// StartAuthFlow uses gcloud for simple authentication.
func (om *OAuthManager) StartAuthFlow(_ context.Context) error {
	return fmt.Errorf("for the simplest setup, please use gcloud authentication:\n\n1. Install gcloud CLI: https://cloud.google.com/sdk/docs/install\n2. Run: gcloud auth login\n3. Run: gcloud auth application-default login\n\nThen this app will automatically use your gcloud credentials.\n\nAlternatively, place a service-account.json file in the current directory")
}
