// Voice Transcriber
// Copyright (c) 2025 Ihor Dvoretskyi
//
// Licensed under MIT License

package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// VersionInfo carries the build-time metadata stamped in by -ldflags.
// Pass it to Execute or NewRootCmd; empty fields fall back to "dev" / "unknown".
type VersionInfo struct {
	Version string
	Date    string
	Commit  string
}

// withDefaults returns a copy of v where empty strings are replaced with
// sensible defaults so that unlinked dev builds still print useful output.
func (v VersionInfo) withDefaults() VersionInfo {
	if v.Version == "" {
		v.Version = "dev"
	}

	if v.Date == "" {
		v.Date = "unknown"
	}

	if v.Commit == "" {
		v.Commit = "unknown"
	}

	return v
}

// newVersionCmd constructs the version subcommand, closing over info.
func newVersionCmd(info VersionInfo) *cobra.Command {
	d := info.withDefaults()

	return &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Printf("%s %s\n", appName, d.Version)
			fmt.Printf("Build Date: %s\n", d.Date)
			fmt.Printf("Git Commit: %s\n", d.Commit)
		},
	}
}
