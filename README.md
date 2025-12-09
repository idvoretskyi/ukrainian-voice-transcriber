# Ukrainian Voice Transcriber

Single-binary Ukrainian video-to-text transcription using Google Cloud Speech-to-Text API.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Features

- Ukrainian language optimized (`uk-UA` locale)
- Single binary - no dependencies
- FFmpeg integration for video processing
- Clean directory structure: `input/` → `output/`
- Simple gcloud authentication
- Auto-cleanup of temporary files

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
# Basic transcription (output goes to output/<video-name>/<video-name>.txt)
transcriber transcribe input/video.mp4

# Use different Speech-to-Text model
transcriber transcribe input/video.mp4 --model=latest_long

# Save to specific file
transcriber transcribe input/video.mp4 -o custom-output.txt

# Verbose mode (show detailed progress)
transcriber transcribe input/video.mp4 --verbose
```

## CLI Commands

```bash
# Main commands
transcriber transcribe <video-file> [-o output.txt] [--verbose|--quiet]
transcriber transcribe <video-file> [--model default|latest_long|latest_short]
transcriber version
```

## Directory Structure

```
ukrainian-voice-transcriber/
├── input/              # Place your video files here
│   └── video.mp4
├── output/             # Transcripts are automatically saved here
│   └── video/
│       └── video.txt
└── ukrainian-voice-transcriber  # The binary
```

## Examples

```bash
# Basic usage - transcribe a video file
transcriber transcribe input/meeting.mp4
# Output: output/meeting/meeting.txt

# Use latest_long model for very long audio files
transcriber transcribe input/interview.mp4 --model=latest_long

# Batch processing - transcribe all videos in input/
for video in input/*.mp4; do
    transcriber transcribe "$video"
done

# Save to custom location
transcriber transcribe input/presentation.mp4 -o transcripts/my-transcript.txt
```

## Troubleshooting

**FFmpeg not found**: `brew install ffmpeg` or `sudo apt install ffmpeg`

**Authentication error**: Run `gcloud auth application-default login`

**Permission denied**: Make sure `$GOPATH/bin` is in your `$PATH`

## Building

```bash
git clone https://github.com/idvoretskyi/ukrainian-voice-transcriber.git
cd ukrainian-voice-transcriber
make build
```

## License

MIT License - see LICENSE file for details.

---

Made for Ukrainian content creators
