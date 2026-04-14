# Ukrainian Voice Transcriber

Single-binary Ukrainian media-to-text transcription powered by Google Gemini via Vertex AI.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Features

- Ukrainian language optimized (`uk-UA`)
- Accepts **both video and audio files** as input
- **No Google Cloud Storage required** — audio is sent inline to Gemini
- FFmpeg used only for video-to-audio extraction (skipped for audio files)
- Cost-efficient default model: `gemini-3.1-flash-lite-preview` (~$0.03/hr of audio)
- Handles files up to ~8.4 hours in a single request (no chunking needed)
- Single binary — no extra runtime dependencies
- All Gemini models selectable via `--model` flag

## Supported Input Formats

| Type | Extensions |
|------|-----------|
| **Audio** (sent directly to Gemini) | `.wav`, `.mp3`, `.flac`, `.ogg`, `.m4a`, `.aac`, `.webm`, `.pcm` |
| **Video** (audio extracted via FFmpeg first) | `.mp4`, `.mkv`, `.mov`, `.avi`, `.wmv`, `.flv`, `.ts`, `.mpeg`, `.3gp` |

## Quick Start

### Prerequisites

```bash
# Install FFmpeg (only needed for video files)
brew install ffmpeg        # macOS
sudo apt install ffmpeg    # Ubuntu/Debian

# Install Go 1.26+
brew install go            # macOS
```

### Installation

```bash
go install github.com/idvoretskyi/ukrainian-voice-transcriber/cmd/transcriber@latest
```

> **Note**: Ensure `$(go env GOPATH)/bin` is in your `$PATH`:
> ```bash
> export PATH="$(go env GOPATH)/bin:$PATH"
> ```

### Authentication & Setup

```bash
# Authenticate with Google Cloud
gcloud auth login
gcloud auth application-default login
gcloud config set project YOUR_PROJECT_ID

# Enable the Vertex AI API
gcloud services enable aiplatform.googleapis.com
```

### Usage

```bash
# Transcribe a video file (FFmpeg extracts audio automatically)
transcriber transcribe input/video.mp4

# Transcribe an audio file directly (no FFmpeg needed)
transcriber transcribe input/recording.wav
transcriber transcribe input/interview.mp3

# Save to a specific output file
transcriber transcribe input/video.mp4 -o output.txt

# Use a different Gemini model
transcriber transcribe input/video.mp4 --model gemini-2.5-flash

# Use a different Vertex AI region
transcriber transcribe input/video.mp4 --location europe-west4

# Verbose mode (show detailed progress)
transcriber transcribe input/video.mp4 --verbose

# Quiet mode (only output the transcript path)
transcriber transcribe input/video.mp4 --quiet
```

## CLI Reference

```
Flags:
  --model string      Gemini model to use (default: gemini-3.1-flash-lite-preview)
  --location string   Vertex AI region (default: us-central1)
  -o, --output string Output file path (default: output/<name>/<name>.txt)
  -v, --verbose       Enable verbose output
  -q, --quiet         Suppress all output except results
```

## Model Selection Guide

| Model | ~Cost/hr audio | Quality | Status | Best for |
|-------|---------------|---------|--------|----------|
| `gemini-3.1-flash-lite-preview` | ~$0.03 | Excellent (ASR-optimized) | Preview | **Default** — best quality/cost |
| `gemini-2.5-flash-lite` | ~$0.01 | Good | GA | Maximum cost savings |
| `gemini-2.5-flash` | ~$0.04 | Very good | GA | Production stability |

> Pricing estimates based on 25 tokens/second of audio at published Gemini API rates.

## Output Directory Structure

```
ukrainian-voice-transcriber/
├── input/               # Place your media files here
│   └── video.mp4
│   └── recording.wav
├── output/              # Transcripts are automatically saved here
│   └── video/
│       └── video.txt
└── ukrainian-voice-transcriber  # The binary
```

## Examples

```bash
# Transcribe a long meeting recording (3-4 hours — handled in a single API call)
transcriber transcribe input/board-meeting.mp4

# Transcribe an interview audio file
transcriber transcribe input/interview.m4a

# Batch process all video files
for f in input/*.mp4; do
    transcriber transcribe "$f"
done

# Batch process mixed audio and video
for f in input/*; do
    transcriber transcribe "$f"
done
```

## Troubleshooting

**FFmpeg not found**: Only needed for video files. `brew install ffmpeg` or `sudo apt install ffmpeg`

**Authentication error**: Run `gcloud auth application-default login`

**Project not set**: Run `gcloud config set project YOUR_PROJECT_ID` or `export GOOGLE_CLOUD_PROJECT=your-project-id`

**API not enabled**: Run `gcloud services enable aiplatform.googleapis.com`

**Permission denied on binary**: Make sure `$GOPATH/bin` is in your `$PATH`

## Building from Source

```bash
git clone https://github.com/idvoretskyi/ukrainian-voice-transcriber.git
cd ukrainian-voice-transcriber
make build
```

## License

MIT License — see LICENSE file for details.

---

Made for Ukrainian content creators
