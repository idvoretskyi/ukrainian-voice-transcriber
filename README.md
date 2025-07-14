# Ukrainian Voice Transcriber

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![CI](https://github.com/idvoretskyi/ukrainian-voice-transcriber/actions/workflows/ci.yml/badge.svg)](https://github.com/idvoretskyi/ukrainian-voice-transcriber/actions)

Single-binary Ukrainian video-to-text transcription using Google Cloud Speech-to-Text API.

## Features

- Ukrainian language optimized (`uk-UA` locale)
- Single binary deployment (no dependencies)
- Cost-efficient with auto-cleanup
- FFmpeg integration for video processing
- Simple authentication via gcloud CLI

## Quick Start

### Prerequisites

```bash
# Install FFmpeg
brew install ffmpeg  # macOS
sudo apt install ffmpeg  # Ubuntu

# Install Go 1.24+
brew install go  # macOS
```

### Installation

```bash
go install github.com/idvoretskyi/ukrainian-voice-transcriber/cmd/transcriber@latest
```

> **Note**: Ensure `$(go env GOPATH)/bin` is in your `$PATH`. If not, add this to your shell profile:
> ```bash
> export PATH="$(go env GOPATH)/bin:$PATH"
> ```

### Authentication

```bash
# Install and authenticate with gcloud
curl https://sdk.cloud.google.com | bash
gcloud auth login
gcloud auth application-default login
gcloud config set project YOUR_PROJECT_ID
gcloud services enable speech.googleapis.com storage.googleapis.com
```

### Usage

```bash
# Basic transcription
transcriber transcribe video.mp4

# Save to specific file
transcriber transcribe video.mp4 -o transcript.txt

# Check setup
transcriber setup
```

## Documentation

- 📖 [English User Manual](docs/USER_MANUAL_EN.md)
- 📖 [Ukrainian User Manual](docs/USER_MANUAL_UK.md)

## CLI Commands

```bash
# Main commands
transcriber transcribe video.mp4 [-o output.txt] [--verbose|--quiet]
transcriber auth [--status]
transcriber setup
transcriber version
```

## Examples

```bash
# Basic usage
transcriber transcribe meeting.mp4

# Batch processing
for video in *.mp4; do
    transcriber transcribe "$video"
done
```

## Troubleshooting

**FFmpeg not found**: `brew install ffmpeg` or `sudo apt install ffmpeg`

**Authentication required**: Run `transcriber auth --status` and follow setup instructions

**Permission denied**: Make sure `$GOPATH/bin` is in your `$PATH`
.

## Building

```bash
git clone https://github.com/idvoretskyi/ukrainian-voice-transcriber.git
cd ukrainian-voice-transcriber
make build
```

## License

MIT License - see LICENSE file for details.

---

🇺🇦 Made with ❤️ for Ukrainian content creators
