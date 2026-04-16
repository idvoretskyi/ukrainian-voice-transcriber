// Voice Transcriber
// Copyright (c) 2025 Ihor Dvoretskyi
//
// Licensed under MIT License

package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Build information — set at binary link time via -ldflags from main.go.
var (
	buildVersion = "dev"
	buildDate    = "unknown"
	buildCommit  = "unknown"
)

// SetVersion stores build-time metadata and wires the version into the root
// command so that `voice-transcriber --version` prints the correct string.
func SetVersion(version, date, commit string) {
	if version != "" {
		buildVersion = version
	}

	if date != "" {
		buildDate = date
	}

	if commit != "" {
		buildCommit = commit
	}
}

// newVersionCmd constructs the version subcommand.
func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Printf("%s %s\n", appName, buildVersion)
			fmt.Printf("Build Date: %s\n", buildDate)
			fmt.Printf("Git Commit: %s\n", buildCommit)
		},
	}
}
