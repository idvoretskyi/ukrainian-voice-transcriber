// Voice Transcriber
// Copyright (c) 2025 Ihor Dvoretskyi
//
// Licensed under MIT License

package config_test

import (
	"testing"

	"github.com/idvoretskyi/voice-transcriber/internal/config"
)

func TestValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		cfg     config.Config
		wantErr bool
	}{
		{
			name:    "zero value is valid",
			cfg:     config.Config{},
			wantErr: false,
		},
		{
			name:    "verbose only is valid",
			cfg:     config.Config{Verbose: true},
			wantErr: false,
		},
		{
			name:    "quiet only is valid",
			cfg:     config.Config{Quiet: true},
			wantErr: false,
		},
		{
			name:    "verbose and quiet together is invalid",
			cfg:     config.Config{Verbose: true, Quiet: true},
			wantErr: true,
		},
		{
			name:    "language auto is valid",
			cfg:     config.Config{Language: "auto"},
			wantErr: false,
		},
		{
			name:    "language empty is valid",
			cfg:     config.Config{Language: ""},
			wantErr: false,
		},
		{
			name:    "two-letter ISO code is valid",
			cfg:     config.Config{Language: "uk"},
			wantErr: false,
		},
		{
			name:    "uppercase ISO code is valid (normalized internally)",
			cfg:     config.Config{Language: "EN"},
			wantErr: false,
		},
		{
			name:    "full language name is invalid",
			cfg:     config.Config{Language: "english"},
			wantErr: true,
		},
		{
			name:    "numeric language code is invalid",
			cfg:     config.Config{Language: "42"},
			wantErr: true,
		},
		{
			name:    "three-letter code is invalid",
			cfg:     config.Config{Language: "ukr"},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := tc.cfg.Validate()

			if tc.wantErr && err == nil {
				t.Errorf("Validate() = nil; want error")
			}

			if !tc.wantErr && err != nil {
				t.Errorf("Validate() = %v; want nil", err)
			}
		})
	}
}

func TestFromEnv(t *testing.T) {
	// t.Setenv is incompatible with t.Parallel on subtests; run sequentially.
	t.Run("GOOGLE_CLOUD_PROJECT is read", func(t *testing.T) {
		t.Setenv("GOOGLE_CLOUD_PROJECT", "my-test-project")

		cfg := config.FromEnv()

		if cfg.GCPProject != "my-test-project" {
			t.Errorf("FromEnv().GCPProject = %q; want %q", cfg.GCPProject, "my-test-project")
		}
	})

	t.Run("unset GOOGLE_CLOUD_PROJECT returns empty string", func(t *testing.T) {
		t.Setenv("GOOGLE_CLOUD_PROJECT", "")

		cfg := config.FromEnv()

		if cfg.GCPProject != "" {
			t.Errorf("FromEnv().GCPProject = %q; want empty", cfg.GCPProject)
		}
	})
}
