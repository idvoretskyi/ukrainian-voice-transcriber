package auth

import (
	"context"
	"fmt"
)

// OAuthManager handles Google OAuth authentication
type OAuthManager struct {
	projectID string
}

// NewOAuthManager creates a new OAuth manager with a simple approach
func NewOAuthManager() *OAuthManager {
	// Use a simple approach: just recommend gcloud auth
	// Most users already have gcloud installed
	return &OAuthManager{}
}

// StartAuthFlow uses gcloud for simple authentication
func (om *OAuthManager) StartAuthFlow(ctx context.Context) error {
	return fmt.Errorf(`For the simplest setup, please use gcloud authentication:

1. Install gcloud CLI: https://cloud.google.com/sdk/docs/install
2. Run: gcloud auth login
3. Run: gcloud auth application-default login

Then this app will automatically use your gcloud credentials.

Alternatively, place a service-account.json file in the current directory.`)
}

