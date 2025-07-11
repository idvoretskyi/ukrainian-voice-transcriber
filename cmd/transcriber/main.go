// Ukrainian Voice Transcriber
//
// Copyright (c) 2025 Ihor Dvoretskyi
//
// Licensed under MIT License

// Package main provides the application entry point.
package main

import (
	"os"

	"github.com/idvoretskyi/ukrainian-voice-transcriber/internal/cli"
)

// Version information - set by build flags.
var (
	version   = "dev"
	buildDate = "unknown"
	gitCommit = "unknown"
)

func main() {
	// Set version information for CLI
	cli.SetVersion(version, buildDate, gitCommit)

	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
