// Voice Transcriber
// Copyright (c) 2025 Ihor Dvoretskyi
//
// Licensed under MIT License

// Package main provides the application entry point.
package main

import (
	"os"

	"github.com/idvoretskyi/voice-transcriber/internal/cli"
)

// Version information - set by build flags.
var (
	version   = "dev"
	buildDate = "unknown"
	gitCommit = "unknown"
)

func main() {
	if err := cli.Execute(cli.VersionInfo{
		Version: version,
		Date:    buildDate,
		Commit:  gitCommit,
	}); err != nil {
		os.Exit(1)
	}
}
