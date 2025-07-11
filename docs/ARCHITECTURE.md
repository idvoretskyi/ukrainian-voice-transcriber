# Architecture Overview

## Project Structure

```
ukrainian-voice-transcriber/
├── cmd/
│   └── transcriber/
│       └── main.go              # Application entry point
├── internal/
│   ├── cli/                     # Command-line interface
│   │   ├── root.go             # Root command and global flags
│   │   ├── transcribe.go       # Transcribe command
│   │   ├── setup.go            # Setup command
│   │   └── version.go          # Version command
│   ├── transcriber/            # Core transcription logic
│   │   ├── transcriber.go      # Main transcriber service
│   │   └── audio.go            # Audio extraction utilities
│   ├── speech/                 # Google Cloud Speech-to-Text wrapper
│   │   └── service.go          # Speech service implementation
│   └── storage/                # Google Cloud Storage wrapper
│       └── service.go          # Storage service implementation
├── pkg/
│   └── config/                 # Shared configuration
│       └── config.go           # Config types and utilities
├── examples/                   # Usage examples
├── docs/                       # Documentation
├── Makefile                    # Build automation
├── go.mod                      # Go module definition
└── README.md                   # Main documentation
```

## Design Principles

### 1. **Single Responsibility**
Each package has a clear, focused responsibility:
- `cmd/transcriber`: Application entry point
- `internal/cli`: Command-line interface logic
- `internal/transcriber`: Core business logic
- `internal/speech`: Google Cloud Speech-to-Text abstraction
- `internal/storage`: Google Cloud Storage abstraction
- `pkg/config`: Shared configuration types

### 2. **Dependency Injection**
Services are initialized with their dependencies explicitly:
```go
speechService := speech.NewService(speechClient, config)
storageService := storage.NewService(storageClient, config)
transcriber := transcriber.New(config, speechService, storageService)
```

### 3. **Interface Segregation**
Each service exposes only the methods it needs:
- Speech service: `TranscribeAudio()`
- Storage service: `UploadFile()`, `CleanupFile()`, `EnsureBucket()`
- Transcriber: `TranscribeLocalFile()`, `Close()`

### 4. **Error Handling**
Consistent error handling throughout:
- Errors are wrapped with context
- User-friendly error messages
- Structured error responses

## Data Flow

```
1. CLI Command → 2. Transcriber → 3. Audio Extraction → 4. Storage Upload
                                                             ↓
8. Result ← 7. Cleanup ← 6. Storage Cleanup ← 5. Speech-to-Text
```

### Detailed Flow:

1. **CLI Command**: User runs transcribe command
2. **Transcriber**: Validates input and orchestrates the process
3. **Audio Extraction**: FFmpeg extracts audio from video
4. **Storage Upload**: Audio uploaded to Google Cloud Storage
5. **Speech-to-Text**: Google Cloud Speech API transcribes audio
6. **Storage Cleanup**: Temporary audio file removed from GCS
7. **Local Cleanup**: Local audio file removed
8. **Result**: Transcript returned to user

## Component Interactions

### CLI Layer
- **Root Command**: Global flags and command routing
- **Subcommands**: Specific functionality (transcribe, setup, version)
- **Flag Parsing**: Configuration from command-line arguments

### Service Layer
- **Transcriber**: Orchestrates the entire transcription process
- **Speech Service**: Abstracts Google Cloud Speech-to-Text API
- **Storage Service**: Abstracts Google Cloud Storage operations

### Infrastructure Layer
- **Google Cloud APIs**: Speech-to-Text and Storage
- **FFmpeg**: Audio extraction from video files
- **File System**: Local file operations

## Configuration Management

### Configuration Sources (Priority Order):
1. Command-line flags (`--verbose`, `--bucket`, etc.)
2. Environment variables (`GOOGLE_APPLICATION_CREDENTIALS`)
3. Auto-detected files (`service-account.json`)
4. Default values

### Configuration Structure:
```go
type Config struct {
    ServiceAccountPath string  // Path to service account JSON
    DriveCredentials   string  // Path to Drive credentials (future)
    BucketName         string  // GCS bucket name
    Verbose            bool    // Enable verbose output
    Quiet              bool    // Suppress output
}
```

## Build Process

### Multi-Platform Builds:
```bash
make build-all  # Builds for Linux, macOS, Windows (amd64/arm64)
```

### Build Targets:
- `ukrainian-voice-transcriber` (current platform)
- `ukrainian-voice-transcriber-linux-amd64`
- `ukrainian-voice-transcriber-linux-arm64`
- `ukrainian-voice-transcriber-darwin-amd64`
- `ukrainian-voice-transcriber-darwin-arm64`
- `ukrainian-voice-transcriber-windows-amd64.exe`

## Security Considerations

### Credential Management:
- Service account keys auto-detected from local files
- No credentials stored in code or version control
- Credentials passed securely to Google Cloud APIs

### Temporary File Handling:
- Local audio files cleaned up immediately after upload
- GCS files have 1-day lifecycle policy for automatic cleanup
- No sensitive data persisted locally

### Network Security:
- All API calls use HTTPS/TLS
- Google Cloud client libraries handle authentication
- No custom network protocols

## Performance Optimizations

### Cost Efficiency:
- Uses standard Speech-to-Text model (not enhanced)
- Automatic cleanup prevents storage costs
- Efficient audio encoding (16kHz mono)

### Processing Efficiency:
- Single-pass audio extraction
- Streaming upload to GCS
- Async operations where possible

### Memory Management:
- Streams large files instead of loading into memory
- Immediate cleanup of temporary resources
- Efficient string building for transcripts

## Extensibility

### Adding New Commands:
1. Create new file in `internal/cli/`
2. Implement cobra command
3. Add to root command in `init()`

### Adding New Services:
1. Create new package in `internal/`
2. Define service interface
3. Implement service with dependency injection
4. Add to transcriber initialization

### Adding New Configuration:
1. Update `Config` struct in `pkg/config/`
2. Add command-line flags in `internal/cli/root.go`
3. Update validation and defaults

This architecture provides a solid foundation for a production-ready transcription service while maintaining simplicity and extensibility.