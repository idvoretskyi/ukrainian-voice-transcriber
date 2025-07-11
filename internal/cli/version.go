// Ukrainian Voice Transcriber
// Copyright (c) {{ YEAR }} Ihor Dvoretskyi
//
// Licensed under MIT License

// Package cli provides command-line interface functionality.
package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Build information.
var (
	buildVersion = "dev"
	buildDate    = "unknown"
	buildCommit  = "unknown"
)

// SetVersion sets the version information from build flags.
func SetVersion(version, date, commit string) {
	buildVersion = version
	buildDate = date
	buildCommit = commit
}

// versionCmd represents the version command.
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Run: func(_ *cobra.Command, args []string) {
		fmt.Printf("%s %s\n", appName, buildVersion)
		fmt.Printf("Build Date: %s\n", buildDate)
		fmt.Printf("Git Commit: %s\n", buildCommit)
	},
}
