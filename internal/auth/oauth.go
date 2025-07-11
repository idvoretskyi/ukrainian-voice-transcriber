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
	projectID string //nolint:unused
}

// NewOAuthManager creates a new OAuth manager with a simple approach.
func NewOAuthManager() *OAuthManager {
	// Use a simple approach: just recommend gcloud auth
	// Most users already have gcloud installed
	return &OAuthManager{}
}

// StartAuthFlow uses gcloud for simple authentication.
func (om *OAuthManager) StartAuthFlow(_ context.Context) error {
	return fmt.Errorf(`for the simplest setup, please use gcloud authentication:

1. Install gcloud CLI: https://cloud.google.com/sdk/docs/install
2. Run: gcloud auth login
3. Run: gcloud auth application-default login

Then this app will automatically use your gcloud credentials.

Alternatively, place a service-account.json file in the current directory.`)
}
