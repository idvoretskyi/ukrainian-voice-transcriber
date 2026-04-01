// Ukrainian Voice Transcriber
// Copyright (c) 2025 Ihor Dvoretskyi
//
// Licensed under MIT License

// Package transcriber exports internal symbols for testing.
package transcriber

import "github.com/idvoretskyi/ukrainian-voice-transcriber/pkg/config"

// SetBucketName exposes setBucketName for black-box tests.
func SetBucketName(cfg *config.Config, projectID string) { setBucketName(cfg, projectID) }

// GenerateAudioPath exposes generateAudioPath for black-box tests.
func GenerateAudioPath(videoPath string) string { return generateAudioPath(videoPath) }

// ValidateAndSanitizeVideoPath exposes validateAndSanitizeVideoPath for black-box tests.
func ValidateAndSanitizeVideoPath(videoPath string) (string, error) {
	return validateAndSanitizeVideoPath(videoPath)
}
