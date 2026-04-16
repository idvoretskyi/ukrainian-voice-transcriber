// Voice Transcriber
// Copyright (c) 2025 Ihor Dvoretskyi
//
// Licensed under MIT License

package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"

	"github.com/idvoretskyi/voice-transcriber/internal/config"
	"github.com/idvoretskyi/voice-transcriber/internal/transcriber"
)

// outputSeparatorWidth is the width of the separator line printed after transcription stats.
const outputSeparatorWidth = 50

// sanitizeRe matches characters not allowed in a sanitized filename.
// \p{L} matches any Unicode letter (including Cyrillic), \p{N} any Unicode digit.
var sanitizeRe = regexp.MustCompile(`[^\p{L}\p{N}_\-.]`)

// newTranscribeCmd constructs the transcribe subcommand.
// cfg is the shared config populated by persistent flags on the root command.
func newTranscribeCmd(cfg *config.Config) *cobra.Command {
	var outputFile string

	cmd := &cobra.Command{
		Use:   "transcribe [media-file]",
		Short: "Transcribe a video or audio file to text",
		Long: `Transcribe a video or audio file to text using Google Gemini.

Language is detected automatically from the audio by default.
Use --language to specify an ISO 639-1 code (e.g. uk, en, de).

Supported input formats:
  Video: mp4, mkv, mov, avi, wmv, flv, ts, mpeg, 3gp (audio extracted via FFmpeg)
  Audio: wav, mp3, flac, ogg, m4a, aac, pcm, webm (sent directly to Gemini, no FFmpeg needed)`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return runTranscribe(cfg, args[0], outputFile)
		},
	}

	cmd.Flags().StringVarP(&outputFile, "output", "o", "",
		"Output file path (default: creates directory based on media filename)")

	return cmd
}

// runTranscribe is the extracted body of the transcribe RunE, making it testable.
func runTranscribe(cfg *config.Config, mediaFile, outputFile string) error {
	ctx := context.Background()
	logger := newLogger(cfg)

	t, err := transcriber.New(ctx, cfg, logger)
	if err != nil {
		return fmt.Errorf("initialization failed: %w", err)
	}

	result, err := t.TranscribeLocalFile(ctx, mediaFile)
	if err != nil {
		return fmt.Errorf("transcription failed: %w", err)
	}

	if !cfg.Quiet {
		fmt.Printf("\n📝 Transcription completed:\n")
		fmt.Printf("   Words: %d\n", result.WordCount)
		fmt.Printf("   Characters: %d\n", len(result.Text))
		fmt.Printf("   Processing time: %v\n", result.ProcessingTime)
		fmt.Println(strings.Repeat("-", outputSeparatorWidth))
	}

	// Determine output path.
	transcriptPath := outputFile
	if transcriptPath == "" {
		mediaBaseName := filepath.Base(mediaFile)
		mediaNameWithoutExt := strings.TrimSuffix(mediaBaseName, filepath.Ext(mediaBaseName))
		sanitizedName := sanitizeFilename(mediaNameWithoutExt)
		outputSubDir := filepath.Join("output", sanitizedName)
		transcriptPath = filepath.Join(outputSubDir, sanitizedName+".txt")
	}

	// Ensure the directory for the output file exists.
	outputDir := filepath.Dir(transcriptPath)
	if err := os.MkdirAll(outputDir, 0o750); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Save transcript with secure file permissions (0600 = rw-------)
	if err := os.WriteFile(transcriptPath, []byte(result.Text), 0o600); err != nil {
		return fmt.Errorf("failed to save transcript: %w", err)
	}

	if !cfg.Quiet {
		fmt.Printf("✅ Transcript saved to: %s\n", transcriptPath)
	}

	return nil
}

// sanitizeFilename removes special characters and replaces spaces with underscores
// to create a safe filename for use in the filesystem.
// Preserves Unicode letters (including multilingual scripts) for internationalized filenames.
func sanitizeFilename(filename string) string {
	filename = strings.ReplaceAll(filename, " ", "_")
	filename = sanitizeRe.ReplaceAllString(filename, "")

	if filename == "" {
		return "transcript"
	}

	return filename
}
