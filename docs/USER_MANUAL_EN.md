# Ukrainian Voice Transcriber - User Manual

## Table of Contents
1. [Introduction](#introduction)
2. [Installation](#installation)
3. [Authentication Setup](#authentication-setup)
4. [Basic Usage](#basic-usage)
5. [Advanced Features](#advanced-features)
6. [File Organization](#file-organization)
7. [Troubleshooting](#troubleshooting)
8. [Cost Management](#cost-management)
9. [FAQ](#faq)

## Introduction

Ukrainian Voice Transcriber is a command-line tool that converts Ukrainian audio and video files into text using Google Cloud Speech-to-Text API. It's designed for content creators, journalists, researchers, and anyone who needs to transcribe Ukrainian language content.

### Key Features
- **Ukrainian Language Optimized**: Specifically configured for Ukrainian (`uk-UA`) language recognition
- **Multiple Input Formats**: Supports all video formats that FFmpeg can process (MP4, AVI, MOV, MKV, etc.)
- **Smart File Organization**: Automatically creates organized directories for your transcripts
- **Cost-Effective**: Uses standard (not premium) Speech-to-Text models with automatic cleanup
- **Simple Setup**: Works with gcloud CLI - no complex OAuth configuration needed

### System Requirements
- **Operating System**: macOS, Linux, or Windows
- **FFmpeg**: Required for audio extraction from video files
- **Google Cloud Project**: With Speech-to-Text and Cloud Storage APIs enabled
- **Internet Connection**: Required for API calls

## Installation

### Step 1: Install FFmpeg

**macOS:**
```bash
brew install ffmpeg
```

**Ubuntu/Debian:**
```bash
sudo apt update
sudo apt install ffmpeg
```

**Windows:**
1. Download FFmpeg from https://ffmpeg.org/download.html
2. Extract to a folder (e.g., `C:\ffmpeg`)
3. Add `C:\ffmpeg\bin` to your PATH environment variable

**Verify Installation:**
```bash
ffmpeg -version
```

### Step 2: Install Ukrainian Voice Transcriber

**Option A: Download Pre-built Binary**
1. Go to the [Releases page](https://github.com/idvoretskyi/ukrainian-voice-transcriber/releases)
2. Download the binary for your operating system
3. Make it executable (Linux/macOS):
   ```bash
   chmod +x ukrainian-voice-transcriber
   ```

**Option B: Build from Source**
```bash
# Install Go (1.24 or later)
# Download from https://golang.org/dl/

# Clone and build
git clone https://github.com/idvoretskyi/ukrainian-voice-transcriber.git
cd ukrainian-voice-transcriber
go build -o ukrainian-voice-transcriber cmd/transcriber/main.go
```

### Step 3: Verify Installation
```bash
./ukrainian-voice-transcriber --help
```

## Authentication Setup

The application supports two authentication methods. We recommend using gcloud CLI for simplicity.

### Method 1: gcloud CLI (Recommended)

**Step 1: Install gcloud CLI**
```bash
# macOS
brew install google-cloud-sdk

# Linux
curl https://sdk.cloud.google.com | bash
exec -l $SHELL

# Windows - Download installer from:
# https://cloud.google.com/sdk/docs/install
```

**Step 2: Authenticate and Configure**
```bash
# Login to your Google account
gcloud auth login

# Set up application default credentials
gcloud auth application-default login

# Set your project ID (replace with your actual project ID)
gcloud config set project YOUR_PROJECT_ID

# Enable required APIs
gcloud services enable speech.googleapis.com
gcloud services enable storage.googleapis.com
```

**Step 3: Verify Setup**
```bash
# Check authentication status
./ukrainian-voice-transcriber auth --status

# Run setup verification
./ukrainian-voice-transcriber setup
```

### Method 2: Service Account (Advanced)

**Step 1: Create Service Account**
1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Select your project or create a new one
3. Navigate to "IAM & Admin" → "Service Accounts"
4. Click "Create Service Account"
5. Enter a name (e.g., "ukrainian-transcriber")
6. Click "Create and Continue"

**Step 2: Assign Roles**
Add these roles:
- `Cloud Speech Client`
- `Storage Admin`

**Step 3: Create Key**
1. Click on the created service account
2. Go to "Keys" tab
3. Click "Add Key" → "Create New Key"
4. Select "JSON" format
5. Download the key file

**Step 4: Configure**
```bash
# Place the key file in the application directory
cp ~/Downloads/service-account-key.json ./service-account.json

# Or set environment variable
export GOOGLE_APPLICATION_CREDENTIALS="/path/to/service-account.json"
```

## Basic Usage

### First-Time Setup
```bash
# Verify everything is configured correctly
./ukrainian-voice-transcriber setup
```

### Basic Transcription
```bash
# Transcribe a video file
./ukrainian-voice-transcriber transcribe video.mp4
```

This will:
1. Create a directory named `video/`
2. Save the transcript as `video/video.txt`
3. Display progress and results

### Video with Spaces in Name
```bash
# Handles spaces automatically
./ukrainian-voice-transcriber transcribe "My Interview 2024.mp4"
```

This will:
1. Create a directory named `My_Interview_2024/`
2. Save the transcript as `My_Interview_2024/My_Interview_2024.txt`

### Save to Specific File
```bash
# Override default file organization
./ukrainian-voice-transcriber transcribe video.mp4 -o custom_transcript.txt
```

### Verbose Output
```bash
# See detailed processing information
./ukrainian-voice-transcriber transcribe video.mp4 --verbose
```

### Quiet Mode
```bash
# Only show the final transcript
./ukrainian-voice-transcriber transcribe video.mp4 --quiet
```

## Advanced Features

### Batch Processing
```bash
# Process all MP4 files in current directory
for video in *.mp4; do
    echo "Processing: $video"
    ./ukrainian-voice-transcriber transcribe "$video"
done
```

### Custom Storage Bucket
```bash
# Use your own Cloud Storage bucket
./ukrainian-voice-transcriber transcribe video.mp4 --bucket my-custom-bucket
```

### Integration with Scripts
```bash
#!/bin/bash
# Example: Process and count words

VIDEO_FILE="$1"
if [ -z "$VIDEO_FILE" ]; then
    echo "Usage: $0 <video-file>"
    exit 1
fi

echo "Transcribing $VIDEO_FILE..."
./ukrainian-voice-transcriber transcribe "$VIDEO_FILE" --quiet > transcript.txt

WORD_COUNT=$(wc -w < transcript.txt)
echo "Transcription complete. Word count: $WORD_COUNT"
```

## File Organization

### Default Organization
When you run:
```bash
./ukrainian-voice-transcriber transcribe "My Video.mp4"
```

The application creates:
```
My_Video/
└── My_Video.txt
```

### Multiple Files
```bash
./ukrainian-voice-transcriber transcribe interview1.mp4
./ukrainian-voice-transcriber transcribe interview2.mp4
./ukrainian-voice-transcriber transcribe "Final Discussion.mp4"
```

Results in:
```
interview1/
└── interview1.txt
interview2/
└── interview2.txt
Final_Discussion/
└── Final_Discussion.txt
```

### Custom Output
```bash
# Save to specific location
./ukrainian-voice-transcriber transcribe video.mp4 -o transcripts/my-transcript.txt
```

## Troubleshooting

### Common Issues

**1. "FFmpeg not found"**
```bash
# Check if FFmpeg is installed
ffmpeg -version

# Install if missing
brew install ffmpeg  # macOS
sudo apt install ffmpeg  # Ubuntu
```

**2. "Authentication required"**
```bash
# Check authentication status
./ukrainian-voice-transcriber auth --status

# If using gcloud
gcloud auth login
gcloud auth application-default login

# If using service account
ls -la service-account.json
```

**3. "Permission denied"**
```bash
# Make binary executable
chmod +x ukrainian-voice-transcriber

# Or run with full path
./ukrainian-voice-transcriber transcribe video.mp4
```

**4. "Project not found" or "Unknown project id"**
```bash
# Set project ID
gcloud config set project YOUR_PROJECT_ID

# Verify project
gcloud config get-value project
```

**5. "API not enabled"**
```bash
# Enable required APIs
gcloud services enable speech.googleapis.com
gcloud services enable storage.googleapis.com
```

### Debug Mode
```bash
# Get detailed error information
./ukrainian-voice-transcriber transcribe video.mp4 --verbose
```

### Log Files
The application doesn't create log files by default. For debugging:
```bash
# Redirect output to file
./ukrainian-voice-transcriber transcribe video.mp4 --verbose > debug.log 2>&1
```

## Cost Management

### Understanding Costs
- **Speech-to-Text**: ~$0.006 per 15-second chunk
- **Cloud Storage**: ~$0.020 per GB/month (temporary files)
- **Typical 1-hour video**: ~$1.44

### Cost Optimization Tips

**1. Use Standard Model**
The application uses the standard (not enhanced) Speech-to-Text model by default for cost efficiency.

**2. Automatic Cleanup**
Temporary files are automatically deleted after processing and have a 1-day lifecycle policy.

**3. Monitor Usage**
```bash
# Check your Google Cloud billing
gcloud billing accounts list
gcloud billing budgets list
```

**4. Set Budget Alerts**
1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Navigate to "Billing" → "Budgets & alerts"
3. Create a budget with email notifications

### Cost Estimation
```bash
# For planning purposes
echo "Duration: 60 minutes"
echo "Cost estimate: $1.44 (60 min × $0.006/15-sec × 4 chunks/min)"
```

## FAQ

### General Questions

**Q: What video formats are supported?**
A: All formats supported by FFmpeg (MP4, AVI, MOV, MKV, WebM, etc.)

**Q: What audio formats are supported?**
A: MP3, WAV, FLAC, M4A, and others supported by FFmpeg

**Q: Can I transcribe audio-only files?**
A: Yes, the application works with both video and audio files

**Q: How accurate is the transcription?**
A: Accuracy depends on audio quality, speaker clarity, and background noise. Ukrainian recognition is optimized for the `uk-UA` locale.

**Q: Can I transcribe other languages?**
A: The application is specifically optimized for Ukrainian. For other languages, you'd need to modify the language settings in the code.

### Technical Questions

**Q: Where are temporary files stored?**
A: In Google Cloud Storage in a bucket named `{project-id}-ukr-voice-transcriber`

**Q: How long are temporary files kept?**
A: Temporary files are deleted immediately after processing and have a 1-day lifecycle policy as backup

**Q: Can I use my own storage bucket?**
A: Yes, use the `--bucket` flag to specify a custom bucket name

**Q: What happens if transcription fails?**
A: Temporary files are cleaned up automatically, and detailed error messages are provided

### Security Questions

**Q: Is my data secure?**
A: Yes, files are processed through Google Cloud's secure infrastructure and temporary files are automatically deleted

**Q: Are credentials stored locally?**
A: Only application default credentials (via gcloud) or service account keys are used. No custom credentials are stored

**Q: Can I use this in production?**
A: Yes, but ensure you have proper authentication, monitoring, and error handling in place

### Performance Questions

**Q: How long does transcription take?**
A: Processing time varies but is typically 1-2x the audio duration (e.g., 10 minutes for a 5-minute video)

**Q: Can I process multiple files simultaneously?**
A: The application processes one file at a time, but you can run multiple instances in parallel

**Q: What are the file size limits?**
A: Google Cloud Speech-to-Text has a 10MB limit for synchronous requests, but the application handles longer files by uploading to Cloud Storage

---

**Need more help?**
- Check the [troubleshooting section](#troubleshooting)
- Visit the [GitHub repository](https://github.com/idvoretskyi/ukrainian-voice-transcriber)
- Review the [README.md](../README.md) for technical details

**🇺🇦 Created with ❤️ for Ukrainian content creators**