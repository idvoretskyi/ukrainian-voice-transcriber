// Ukrainian Voice Transcriber
// Copyright (c) 2025 Ihor Dvoretskyi
//
// Licensed under MIT License

// Package cli provides command-line interface functionality.
package cli

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	speechapi "cloud.google.com/go/speech/apiv1"
	"github.com/spf13/cobra"

	"github.com/idvoretskyi/ukrainian-voice-transcriber/internal/transcriber"
	"github.com/idvoretskyi/ukrainian-voice-transcriber/pkg/config"
)

// setupCmd represents the setup command.
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Check setup and configuration",
	RunE: func(_ *cobra.Command, _ []string) error {
		fmt.Printf("🚀 %s v%s - Setup Check\n", appName, buildVersion)
		fmt.Println(strings.Repeat("=", 50))

		// Check FFmpeg
		if _, err := exec.LookPath("ffmpeg"); err != nil {
			fmt.Println("❌ FFmpeg not found")
			fmt.Println("   Install: brew install ffmpeg (macOS) or apt install ffmpeg (Ubuntu)")

			return fmt.Errorf("FFmpeg required")
		}
		fmt.Println("✅ FFmpeg found")

		// Check authentication
		fmt.Printf("Checking authentication...\n")

		// Try Application Default Credentials first
		ctx := context.Background()
		_, err := speechapi.NewClient(ctx)
		if err == nil {
			fmt.Printf("✅ Application Default Credentials working\n")
		} else {
			// Check service account as fallback
			serviceAccount := config.FindServiceAccount()
			if serviceAccount != "" {
				fmt.Printf("✅ Service account found: %s\n", serviceAccount)
			} else {
				fmt.Println("❌ No authentication found")
				fmt.Println("   Option 1 (Recommended): gcloud auth application-default login")
				fmt.Println("   Option 2: Place service-account.json in current directory")
				fmt.Println("   Option 3: ukrainian-voice-transcriber auth")

				return fmt.Errorf("authentication required")
			}
		}

		// Check drive credentials
		driveCredentials := config.FindDriveCredentials()
		if driveCredentials != "" {
			fmt.Printf("✅ Google Drive credentials found: %s\n", driveCredentials)
		} else {
			fmt.Println("ℹ️  Google Drive credentials not found (optional)")
		}

		// Test initialization
		t, err := transcriber.New(&globalConfig)
		if err != nil {
			return fmt.Errorf("initialization test failed: %v", err)
		}
		defer func() {
			if closeErr := t.Close(); closeErr != nil {
				fmt.Printf("Warning: Failed to close transcriber: %v\n", closeErr)
			}
		}()

		fmt.Println("✅ Google Cloud clients initialized successfully")
		fmt.Printf("✅ Storage bucket ready: %s\n", globalConfig.BucketName)

		fmt.Println("\n🎉 Setup completed successfully!")
		fmt.Println("\nUsage:")
		fmt.Println("  ukrainian-voice-transcriber transcribe video.mp4")
		fmt.Println("  ukrainian-voice-transcriber transcribe video.mp4 -o transcript.txt")

		return nil
	},
}
