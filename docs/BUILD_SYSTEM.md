# Build System Documentation

## Overview

The Ukrainian Voice Transcriber uses a comprehensive build system with automated CI/CD, multi-platform builds, and professional release management.

## Requirements

- **Go 1.24+** - Programming language
- **FFmpeg** - Audio/video processing
- **Docker** - Containerization (optional)
- **Make** - Build automation

## Quick Start

### Local Development
```bash
# Setup development environment
./scripts/dev.sh setup

# Build for current platform
make build

# Build for all platforms
make build-all

# Run full development cycle
./scripts/dev.sh full
```

### Testing
```bash
# Run tests
make test

# Run tests with coverage
make test-coverage

# Run linting
make lint

# Run security scan
make security
```

## Build System Components

### 1. Makefile
Comprehensive build automation with 20+ targets:

```bash
make build           # Build for current platform
make build-all       # Build all platforms (6 total)
make test           # Run tests with race detection
make lint           # Code linting with golangci-lint
make security       # Security scanning with gosec
make docker-build   # Build Docker image
make release-prepare # Prepare release
```

**Supported Platforms:**
- Linux (AMD64, ARM64)
- macOS (Intel, Apple Silicon)
- Windows (AMD64, ARM64)

### 2. GitHub Actions

#### Continuous Integration (`.github/workflows/ci.yml`)
- **Triggers**: Push/PR to main branch
- **Jobs**: test, lint, build, security, dependency-check, security-scan
- **Features**: 
  - Go 1.24 testing
  - Race condition detection
  - Code coverage reporting
  - Multi-platform build verification
  - Security scanning (gosec + Trivy)
  - Vulnerability checking (govulncheck)
  - SARIF upload to GitHub Security tab

#### Security Analysis (`.github/workflows/security.yml`)
- **Triggers**: Push/PR to main branch, weekly schedule
- **Jobs**: codeql, dependency-review, govulncheck, trivy, gosec
- **Features**:
  - CodeQL analysis with GitHub Advanced Security
  - Dependency review on pull requests
  - Comprehensive vulnerability scanning
  - Container and filesystem security scanning
  - SARIF results uploaded to Security tab

#### Release Workflow (`.github/workflows/goreleaser.yml`)
- **Triggers**: Version tags (`v*`)
- **Jobs**: test, release (using GoReleaser)
- **Features**:
  - Multi-platform binary builds (6 platforms)
  - Automated GitHub releases with professional formatting
  - Docker image publishing to GitHub Container Registry
  - Homebrew tap publishing
  - Linux package publishing (DEB, RPM, APK)
  - Automated changelog generation
  - Checksum generation and verification
  - Release asset signing

### 3. Docker Support

**Dockerfile Features:**
- Multi-stage build (Go 1.24 + Alpine)
- FFmpeg included
- Non-root user
- Health checks
- Optimized for size

**Docker Commands:**
```bash
make docker-build    # Build image
make docker-push     # Push to registry
make docker-run      # Run container
```

**Registry**: `ghcr.io/idvoretskyi/ukrainian-voice-transcriber`

### 4. Release Management

#### Version Management
- **Source**: Git tags (`v1.0.0`)
- **Build flags**: Version, build date, commit hash
- **Semantic versioning**: Automatic from tags

#### Release Process
```bash
# 1. Prepare release
make release-check
git tag v1.0.0

# 2. Push tag (triggers automation)
git push origin v1.0.0

# 3. Automation handles:
# - Building all platforms
# - Creating GitHub release
# - Publishing Docker images
# - Publishing packages
```

#### Distribution Channels
1. **GitHub Releases** - Binary downloads
2. **GitHub Container Registry** - Docker images
3. **Homebrew Tap** - macOS package manager
4. **Linux Packages** - DEB, RPM, APK formats

### 5. Development Tools

#### Development Script (`scripts/dev.sh`)
```bash
./scripts/dev.sh setup      # Setup environment
./scripts/dev.sh check      # Check dependencies
./scripts/dev.sh test       # Run tests
./scripts/dev.sh lint       # Run linting
./scripts/dev.sh security   # Run security scan
./scripts/dev.sh build      # Build project
./scripts/dev.sh full       # Full development cycle
```

#### Code Quality Tools
- **golangci-lint** - Comprehensive linting
- **gosec** - Security scanning
- **govulncheck** - Vulnerability checking
- **Race detector** - Concurrency bug detection

## Configuration Files

### `.goreleaser.yml`
Professional release configuration:
- Multi-platform builds
- Archive generation
- Checksum creation
- Package publishing
- Docker image building

### `.golangci.yml`
Linting configuration:
- 30+ enabled linters
- Custom rules and exclusions
- Performance optimizations
- Security checks

### `.github/settings.yml`
Repository settings as code:
- Branch protection rules
- Labels and milestones
- Security settings
- Collaborator permissions
- Issue templates

### `.github/dependabot.yml`
Automated dependency updates:
- Go modules weekly updates
- GitHub Actions updates
- Docker base image updates
- Security vulnerability patches

### `.github/CODEOWNERS`
Code review assignments:
- Automatic reviewer assignment
- Component-specific owners
- Documentation maintainers

### Build Flags
Version information embedded at build time:
```go
var (
    version   = "dev"        // Set by: -X 'main.version=v1.0.0'
    buildDate = "unknown"    // Set by: -X 'main.buildDate=2024-01-01T00:00:00Z'
    gitCommit = "unknown"    // Set by: -X 'main.gitCommit=abc123'
)
```

## Continuous Integration Matrix

### Test Matrix
- **Go Version**: 1.24
- **Platforms**: ubuntu-latest
- **Features**: race detection, coverage, benchmarks

### Build Matrix
- **Platforms**: Linux, macOS, Windows
- **Architectures**: AMD64, ARM64
- **Output**: Static binaries with embedded metadata

### Security Matrix
- **CodeQL Analysis**: GitHub Advanced Security scanning
- **Static Analysis**: gosec with SARIF output
- **Vulnerability Scanning**: govulncheck + Trivy
- **Dependency Review**: Automated dependency security checks
- **Container Scanning**: Trivy filesystem scanning
- **Dependency Checking**: go mod verify

## Troubleshooting

### Common Issues

**Build Failures:**
```bash
# Check Go version
go version  # Should be 1.24+

# Update dependencies
go mod tidy

# Clean and rebuild
make clean && make build
```

**Test Failures:**
```bash
# Run with verbose output
go test -v ./...

# Run specific test
go test -v ./internal/cli -run TestSpecificFunction
```

**Docker Issues:**
```bash
# Check Docker daemon
docker version

# Clean Docker cache
docker system prune -a
```

### Development Environment Setup

**macOS:**
```bash
# Install dependencies
brew install go ffmpeg docker make

# Setup project
./scripts/dev.sh setup
```

**Ubuntu:**
```bash
# Install dependencies
sudo apt update
sudo apt install golang-1.24 ffmpeg docker.io make

# Setup project
./scripts/dev.sh setup
```

## Performance Optimization

### Build Optimization
- Static linking (`CGO_ENABLED=0`)
- Binary stripping (`-ldflags="-w -s"`)
- Parallel builds
- Build caching

### CI/CD Optimization
- Go module caching
- Artifact caching
- Parallel job execution
- Conditional builds

## Monitoring and Metrics

### Build Metrics
- Build time tracking
- Binary size monitoring
- Test coverage reporting
- Security scan results

### Release Metrics
- Download statistics
- Platform usage
- Error reporting
- Performance metrics

---

**For more information:**
- [Makefile](../Makefile) - Build targets
- [CI Workflow](../.github/workflows/ci.yml) - Continuous integration
- [Security Workflow](../.github/workflows/security.yml) - Security analysis
- [Release Workflow](../.github/workflows/goreleaser.yml) - Automated releases
- [GoReleaser Config](../.goreleaser.yml) - Release configuration
- [CodeQL Config](../.github/codeql/codeql-config.yml) - Security scanning configuration
- [Repository Settings](../.github/settings.yml) - Repository configuration
- [Dependabot Config](../.github/dependabot.yml) - Dependency automation
- [Code Owners](../.github/CODEOWNERS) - Review assignments
- [Development Script](../scripts/dev.sh) - Development automation