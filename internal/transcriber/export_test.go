// Voice Transcriber
// Copyright (c) 2025 Ihor Dvoretskyi
//
// Licensed under MIT License

// Package transcriber exports internal symbols for testing.
package transcriber

// GenerateAudioPath exposes generateAudioPath for black-box tests.
func GenerateAudioPath(inputPath string) (string, error) { return generateAudioPath(inputPath) }

// ValidateInputPath exposes validateInputPath for black-box tests.
func ValidateInputPath(inputPath string) (string, error) {
	return validateInputPath(inputPath)
}

// ClassifyInputFile exposes classifyInputFile for black-box tests.
func ClassifyInputFile(inputPath string) (InputType, string) {
	return classifyInputFile(inputPath)
}
