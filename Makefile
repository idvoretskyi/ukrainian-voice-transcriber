# Voice Transcriber - Makefile

BINARY_NAME=voice-transcriber
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_DATE ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
GIT_COMMIT ?= $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")

# Build flags with version information
BUILD_FLAGS=-ldflags="-w -s -X 'main.version=$(VERSION)' -X 'main.buildDate=$(BUILD_DATE)' -X 'main.gitCommit=$(GIT_COMMIT)'"

# Go settings
GO_VERSION ?= 1.26
GOPATH ?= $(shell go env GOPATH)
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

# Directories
DIST_DIR=dist

.PHONY: all build clean test help install lint fmt vet security

# Default target
all: build

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	go build $(BUILD_FLAGS) -o $(BINARY_NAME) ./cmd/transcriber
	@echo "Built: $(BINARY_NAME)"

# Build for multiple platforms (linux + darwin, amd64 + arm64)
build-all: clean
	@echo "Building for multiple platforms..."
	@mkdir -p $(DIST_DIR)
	@echo "Building Linux AMD64..."
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build $(BUILD_FLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/transcriber
	@echo "Building Linux ARM64..."
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build $(BUILD_FLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-arm64 ./cmd/transcriber
	@echo "Building macOS AMD64..."
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build $(BUILD_FLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/transcriber
	@echo "Building macOS ARM64..."
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build $(BUILD_FLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/transcriber
	@echo "Built all platforms"
	@ls -la $(DIST_DIR)/

# Install to $GOPATH/bin
install:
	@echo "Installing to GOPATH/bin..."
	go install $(BUILD_FLAGS) ./cmd/transcriber
	@echo "Installed: $(BINARY_NAME)"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME)-*
	rm -rf $(DIST_DIR)
	go clean
	@echo "Cleaned"

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy
	@echo "Dependencies ready"

# Run tests
test:
	@echo "Running tests..."
	go test -v -race -coverprofile=coverage.out ./...
	@echo "Tests completed"

# Run tests with coverage
test-coverage: test
	@echo "Generating coverage report..."
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...
	@echo "Code formatted"

# Vet code
vet:
	@echo "Vetting code..."
	go vet ./...
	@echo "Code vetted"

# Lint code
lint:
	@echo "Linting code..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi
	@echo "Linting completed"

# Security scan
security:
	@echo "Running security scan..."
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "gosec not found. Install with: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"; \
	fi
	@echo "Security scan completed"

# Release targets
release-check:
	@echo "Checking release readiness..."
	@if [ -z "$(VERSION)" ]; then echo "VERSION not set"; exit 1; fi
	@if ! git diff --quiet; then echo "Git working directory not clean"; exit 1; fi
	@echo "Release check passed"

release-prepare: release-check clean fmt vet lint test
	@echo "Preparing release $(VERSION)..."
	@git tag -a $(VERSION) -m "Release $(VERSION)"
	@echo "Release $(VERSION) prepared. Push with: git push origin $(VERSION)"

checksums:
	@echo "Generating checksums..."
	@if [ ! -d "$(DIST_DIR)" ]; then echo "No dist directory. Run 'make build-all' first"; exit 1; fi
	@cd $(DIST_DIR) && sha256sum * > checksums.txt
	@echo "Checksums generated: $(DIST_DIR)/checksums.txt"

# Show help
help:
	@echo "Voice Transcriber - Build Commands"
	@echo "==================================="
	@echo ""
	@echo "Building:"
	@echo "  make build           - Build binary for current platform"
	@echo "  make build-all       - Build for linux+darwin x amd64+arm64"
	@echo "  make install         - Install to GOPATH/bin"
	@echo ""
	@echo "Development:"
	@echo "  make deps            - Download dependencies"
	@echo "  make fmt             - Format code"
	@echo "  make vet             - Vet code"
	@echo "  make lint            - Lint code"
	@echo "  make test            - Run tests"
	@echo "  make test-coverage   - Run tests with coverage"
	@echo "  make security        - Run security scan"
	@echo ""
	@echo "Release:"
	@echo "  make release-check   - Check release readiness"
	@echo "  make release-prepare - Prepare release (tag)"
	@echo "  make checksums       - Generate checksums"
	@echo ""
	@echo "Utilities:"
	@echo "  make clean           - Clean build artifacts"
	@echo "  make help            - Show this help"

# Development shortcuts
dev: deps build

# Release preparation
release: clean deps test build-all
	@echo "Release ready!"
	@ls -la $(DIST_DIR)/
