# Ukrainian Voice Transcriber

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![CI](https://github.com/idvoretskyi/ukrainian-voice-transcriber/workflows/CI/badge.svg)](https://github.com/idvoretskyi/ukrainian-voice-transcriber/actions)

🎯 **Single-binary Ukrainian video-to-text transcription using Google Cloud Speech-to-Text API**

Built with Go for **zero-dependency deployment** and **maximum simplicity**.

## ✨ Features

- 🎯 **Ukrainian Language Optimized** - Specifically tuned for `uk-UA` locale
- 🚀 **Single Binary** - No Python, no virtual environments, no dependency hell
- 💰 **Cost-Efficient** - Uses standard Speech-to-Text model with auto-cleanup
- 🎵 **FFmpeg Integration** - Automatic audio extraction from video files
- 📊 **Detailed Results** - Word counts, processing time, success metrics
- 🔧 **Built-in Help** - Complete CLI help with `-h` flag
- 🧹 **Auto-Cleanup** - Temporary files removed automatically
- 📁 **Smart File Organization** - Creates directories based on video filenames
- 🔐 **Simplified Authentication** - Uses gcloud CLI (no complex OAuth setup required)

## 🚀 Quick Start

### 1. Prerequisites

**Install FFmpeg:**

```bash
# macOS
brew install ffmpeg

# Ubuntu/Debian
sudo apt install ffmpeg

# Windows
# Download from https://ffmpeg.org/download.html
```

**Install Go 1.24+:**

```bash
# macOS
brew install go

# Ubuntu/Debian
sudo apt install golang-1.24

# Or download from https://golang.org/dl/
```

### 2. Installation

**Option A: Install from source (Recommended)**

```bash
# Clone and build
git clone https://github.com/idvoretskyi/ukrainian-voice-transcriber.git
cd ukrainian-voice-transcriber
go build -o ukrainian-voice-transcriber

# Or build and install to $GOPATH/bin
go install
```

**Option B: Direct go install**

```bash
go install github.com/idvoretskyi/ukrainian-voice-transcriber/cmd/transcriber@latest
```

### 3. Authentication Setup

**Option A: gcloud CLI (Recommended - Simplest)**

```bash
# Install gcloud CLI
curl https://sdk.cloud.google.com | bash

# Authenticate with Google Cloud
gcloud auth login

# Set up application default credentials
gcloud auth application-default login

# Set your project ID
gcloud config set project YOUR_PROJECT_ID

# Enable required APIs
gcloud services enable speech.googleapis.com storage.googleapis.com
```

**Option B: Service Account (Advanced)**

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create new project or select existing
3. Enable APIs: **Speech-to-Text** and **Cloud Storage**
4. Create Service Account with roles:
   - **Speech-to-Text Client**
   - **Storage Admin**
5. Download JSON key file
6. **Place in directory as `service-account.json`**

### 4. Usage

**First-time setup:**

```bash
# Check setup and authentication
./ukrainian-voice-transcriber setup

# Check authentication status
./ukrainian-voice-transcriber auth --status
```

**Transcribe video:**

```bash
# Basic transcription - creates directory with video name
./ukrainian-voice-transcriber transcribe video.mp4
# → Creates: video/ directory with video.txt inside

# Video with spaces in name
./ukrainian-voice-transcriber transcribe "My Video File.mp4"
# → Creates: My_Video_File/ directory with My_Video_File.txt inside

# Save to specific file instead
./ukrainian-voice-transcriber transcribe video.mp4 -o transcript.txt

# Verbose mode
./ukrainian-voice-transcriber transcribe video.mp4 --verbose

# Quiet mode (only results)
./ukrainian-voice-transcriber transcribe video.mp4 --quiet
```

**Authentication management:**

```bash
# Check auth status
./ukrainian-voice-transcriber auth --status

# Show setup instructions
./ukrainian-voice-transcriber auth
```

**Get help:**

```bash
./ukrainian-voice-transcriber -h
./ukrainian-voice-transcriber auth -h
./ukrainian-voice-transcriber transcribe -h
```

## 📚 Documentation

**Comprehensive User Manuals:**

- 📖 [**English User Manual**](docs/USER_MANUAL_EN.md) - Complete guide with examples
- 📖 [**Ukrainian User Manual**](docs/USER_MANUAL_UK.md) - Повний посібник з прикладами

## 📖 CLI Reference

### Global Options

```
-v, --verbose    Enable verbose output
-q, --quiet      Suppress all output except results
    --bucket     Custom GCS bucket name (optional)
-h, --help       Show help information
```

### Commands

**`auth`** - Authentication setup and status

```bash
./ukrainian-voice-transcriber auth [OPTIONS]

Options:
      --status    Show current authentication status
```

**`transcribe [video-file]`** - Transcribe video to text

```bash
./ukrainian-voice-transcriber transcribe video.mp4 [OPTIONS]

Options:
  -o, --output string   Output file path (default: stdout)
```

**`setup`** - Check configuration and dependencies

```bash
./ukrainian-voice-transcriber setup
```

**`version`** - Show version information

```bash
./ukrainian-voice-transcriber version
```

## 🔧 Configuration

### Authentication Files

- **`service-account.json`** - Google Cloud service account key (optional)
- **Application Default Credentials** - Automatically used when available (via gcloud)

### Environment Variables

- **`GOOGLE_APPLICATION_CREDENTIALS`** - Path to service account JSON (optional)
- **`GCS_BUCKET_NAME`** - Custom bucket name (optional, default: auto-generated)

### File Organization

- **Output directories** are created based on video filenames
- **Spaces in filenames** are replaced with underscores
- **Transcripts** are saved as `filename.txt` in the created directory

## 💡 Examples

### Basic Usage

```bash
# Simple transcription - creates meeting/ directory with meeting.txt
./ukrainian-voice-transcriber transcribe meeting.mp4

# Video with spaces - creates My_Meeting_2024/ directory with My_Meeting_2024.txt
./ukrainian-voice-transcriber transcribe "My Meeting 2024.mp4"

# Save to specific file instead of directory
./ukrainian-voice-transcriber transcribe lecture.mp4 -o lecture_transcript.txt
```

### Batch Processing

```bash
# Process multiple files - each creates its own directory
for video in *.mp4; do
    echo "Processing $video..."
    ./ukrainian-voice-transcriber transcribe "$video"
done

# Or save to specific files
for video in *.mp4; do
    echo "Processing $video..."
    ./ukrainian-voice-transcriber transcribe "$video" -o "${video%.mp4}_transcript.txt"
done
```

### Integration with Scripts

```bash
#!/bin/bash
# Transcribe and process
TRANSCRIPT=$(./ukrainian-voice-transcriber transcribe video.mp4 --quiet)
echo "Word count: $(echo "$TRANSCRIPT" | wc -w)"
echo "Content: $TRANSCRIPT"
```

## 💰 Cost Optimization

**Current Settings:**

- Uses **standard Speech-to-Text model** (not enhanced)
- **Auto-cleanup** of temporary Cloud Storage files (1-day lifecycle)
- **Efficient audio encoding** (16kHz mono WAV)
- **Single API call** per file

**Estimated Costs:**

- Speech-to-Text: ~$0.006 per 15 seconds
- Cloud Storage: ~$0.020 per GB/month (temporary files)
- **1 hour video**: ~$1.44

## 🐛 Troubleshooting

### Common Issues

**"FFmpeg not found"**

```bash
# Check FFmpeg installation
ffmpeg -version

# Install if missing
brew install ffmpeg  # macOS
sudo apt install ffmpeg  # Ubuntu
```

**"Authentication required"**

```bash
# Check authentication status
./ukrainian-voice-transcriber auth --status

# Set up gcloud authentication
gcloud auth login
gcloud auth application-default login

# Or ensure service account file exists
ls -la service-account.json

# Or set environment variable
export GOOGLE_APPLICATION_CREDENTIALS="/path/to/service-account.json"
```

**"Permission denied" errors**

```bash
# Make binary executable
chmod +x ukrainian-voice-transcriber

# Or run with go
go run main.go transcribe video.mp4
```

**"Bucket creation failed"**

- Ensure billing is enabled in Google Cloud Console
- Check IAM permissions for Storage Admin role
- Verify project ID and region settings

### Debug Mode

```bash
# Enable verbose logging
./ukrainian-voice-transcriber transcribe video.mp4 --verbose

# Check Google Cloud setup
./ukrainian-voice-transcriber setup
```

## 🏗️ Project Structure

```
ukrainian-voice-transcriber/
├── cmd/
│   └── transcriber/            # Application entry point
│       └── main.go
├── internal/
│   ├── cli/                    # Command-line interface
│   │   ├── root.go            # Root command and global flags
│   │   ├── transcribe.go      # Transcribe command
│   │   ├── setup.go           # Setup command
│   │   └── version.go         # Version command
│   ├── transcriber/           # Core transcription logic
│   │   ├── transcriber.go     # Main transcriber service
│   │   └── audio.go           # Audio extraction utilities
│   ├── speech/                # Google Cloud Speech-to-Text wrapper
│   │   └── service.go
│   └── storage/               # Google Cloud Storage wrapper
│       └── service.go
├── pkg/
│   └── config/                # Shared configuration
│       └── config.go
├── examples/                  # Usage examples and scripts
├── docs/                      # Documentation
├── Makefile                   # Build automation
├── go.mod                     # Go module definition
└── README.md                  # This file
```

## 🏗️ Building from Source

### Quick Build

```bash
# Clone repository
git clone https://github.com/idvoretskyi/ukrainian-voice-transcriber.git
cd ukrainian-voice-transcriber

# Build with Make
make build

# Or build directly
go build -o ukrainian-voice-transcriber ./cmd/transcriber
```

### Development Setup

```bash
# Download dependencies
go mod tidy

# Run without building
go run ./cmd/transcriber setup

# Build and test
make build
./ukrainian-voice-transcriber setup
```

### Multi-Platform Builds

```bash
# Build for all platforms
make build-all

# Individual platform builds
make build                    # Current platform
GOOS=linux GOARCH=amd64 go build -o ukrainian-voice-transcriber-linux ./cmd/transcriber
GOOS=windows GOARCH=amd64 go build -o ukrainian-voice-transcriber.exe ./cmd/transcriber
GOOS=darwin GOARCH=arm64 go build -o ukrainian-voice-transcriber-macos ./cmd/transcriber
```

### Dependencies

- **Go 1.24+**
- **FFmpeg** (system dependency)
- Google Cloud Go SDK (auto-downloaded)
- Cobra CLI framework (auto-downloaded)

## 🚀 Deployment

### Single Binary Deployment

```bash
# Build static binary
CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o ukrainian-voice-transcriber

# Deploy anywhere
scp ukrainian-voice-transcriber user@server:/usr/local/bin/
scp service-account.json user@server:/app/
```

### Docker Deployment

```bash
# Build container
docker build -t ukrainian-voice-transcriber .

# Run container
docker run -v $(pwd)/service-account.json:/app/service-account.json \
           -v $(pwd)/videos:/app/videos \
           ukrainian-voice-transcriber transcribe /app/videos/video.mp4
```

## 📄 License

MIT License - see LICENSE file for details.

## 🔒 Security

This project implements comprehensive security scanning:

- **🔍 CodeQL Analysis** - GitHub Advanced Security scanning
- **🛡️ Dependency Review** - Automated security checks on dependencies
- **🔒 Vulnerability Scanning** - Multiple tools (gosec, Trivy, govulncheck)
- **📊 SARIF Reporting** - Security results uploaded to GitHub Security tab
- **⏰ Scheduled Scans** - Weekly automated security analysis

Security findings are automatically reported in the **Security** tab of the GitHub repository.

## 🤝 Contributing

1. Fork the repository: https://github.com/idvoretskyi/ukrainian-voice-transcriber
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Make changes and test (`go test ./...`)
4. Commit changes (`git commit -m 'Add amazing feature'`)
5. Push to branch (`git push origin feature/amazing-feature`)
6. Open Pull Request

**Security**: All contributions are automatically scanned for security vulnerabilities.

## 🆘 Support

- **Issues**: Report bugs on [GitHub Issues](https://github.com/idvoretskyi/ukrainian-voice-transcriber/issues)
- **Documentation**: This README and built-in help (`-h` flag)
- **Examples**: See usage examples above
- **Discussions**: Join conversations on [GitHub Discussions](https://github.com/idvoretskyi/ukrainian-voice-transcriber/discussions)

---

## 🌟 Open Source

This project is **open source** and welcomes contributions!

- **License**: MIT License (see [LICENSE](LICENSE))
- **Contributing**: See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines
- **Issues**: Report bugs and request features via [GitHub Issues](https://github.com/idvoretskyi/ukrainian-voice-transcriber/issues)
- **Security**: Report security issues responsibly

---

**🇺🇦 Made with ❤️ for Ukrainian content creators and teams.**

**Simple. Fast. Reliable. Open Source.**
