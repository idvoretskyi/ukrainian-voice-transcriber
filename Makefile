# Ukrainian Voice Transcriber - Makefile

BINARY_NAME=ukrainian-voice-transcriber
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
DOCKER_IMAGE=ghcr.io/idvoretskyi/ukrainian-voice-transcriber

.PHONY: all build clean test help install lint fmt vet security docker docker-build docker-push

# Default target
all: build

# Build the binary
build:
	@echo "🔨 Building $(BINARY_NAME)..."
	go build $(BUILD_FLAGS) -o $(BINARY_NAME) ./cmd/transcriber
	@echo "✅ Built: $(BINARY_NAME)"

# Build for multiple platforms
build-all: clean
	@echo "🔨 Building for multiple platforms..."
	@mkdir -p $(DIST_DIR)
	@echo "Building Linux AMD64..."
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build $(BUILD_FLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/transcriber
	@echo "Building Linux ARM64..."
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build $(BUILD_FLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-arm64 ./cmd/transcriber
	@echo "Building macOS AMD64..."
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build $(BUILD_FLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/transcriber
	@echo "Building macOS ARM64..."
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build $(BUILD_FLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/transcriber
	@echo "Building Windows AMD64..."
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build $(BUILD_FLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/transcriber
	@echo "Building Windows ARM64..."
	GOOS=windows GOARCH=arm64 CGO_ENABLED=0 go build $(BUILD_FLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-windows-arm64.exe ./cmd/transcriber
	@echo "✅ Built all platforms"
	@ls -la $(DIST_DIR)/

# Build static binary (for Docker/Alpine)
build-static:
	@echo "🔨 Building static binary..."
	CGO_ENABLED=0 GOOS=linux go build $(BUILD_FLAGS) -a -installsuffix cgo -o $(BINARY_NAME)-static ./cmd/transcriber
	@echo "✅ Built: $(BINARY_NAME)-static"

# Install to $GOPATH/bin
install:
	@echo "📦 Installing to GOPATH/bin..."
	go install $(BUILD_FLAGS) ./cmd/transcriber
	@echo "✅ Installed: $(BINARY_NAME)"

# Clean build artifacts
clean:
	@echo "🧹 Cleaning..."
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME)-*
	rm -rf $(DIST_DIR)
	go clean
	@echo "✅ Cleaned"

# Download dependencies
deps:
	@echo "📥 Downloading dependencies..."
	go mod download
	go mod tidy
	@echo "✅ Dependencies ready"

# Run tests
test:
	@echo "🧪 Running tests..."
	go test -v -race -coverprofile=coverage.out ./...
	@echo "✅ Tests completed"

# Run tests with coverage
test-coverage: test
	@echo "📊 Generating coverage report..."
	go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report: coverage.html"

# Format code
fmt:
	@echo "🎨 Formatting code..."
	go fmt ./...
	@echo "✅ Code formatted"

# Vet code
vet:
	@echo "🔍 Vetting code..."
	go vet ./...
	@echo "✅ Code vetted"

# Lint code
lint:
	@echo "🔍 Linting code..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "⚠️  golangci-lint not found. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi
	@echo "✅ Linting completed"

# Security scan
security:
	@echo "🔒 Running security scan..."
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "⚠️  gosec not found. Install with: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"; \
	fi
	@echo "✅ Security scan completed"

# Run setup check
setup: build
	@echo "🔧 Running setup check..."
	./$(BINARY_NAME) setup

# Run with example (requires video file)
demo: build
	@echo "🎬 Running demo..."
	@if [ -f "example.mp4" ]; then \
		./$(BINARY_NAME) transcribe example.mp4 -v; \
	else \
		echo "❌ example.mp4 not found. Add a video file to test."; \
	fi

# Docker targets
docker-build:
	@echo "🐳 Building Docker image..."
	docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg BUILD_DATE=$(BUILD_DATE) \
		--build-arg GIT_COMMIT=$(GIT_COMMIT) \
		-t $(DOCKER_IMAGE):$(VERSION) \
		-t $(DOCKER_IMAGE):latest \
		.
	@echo "✅ Docker image built: $(DOCKER_IMAGE):$(VERSION)"

docker-push: docker-build
	@echo "🐳 Pushing Docker image..."
	docker push $(DOCKER_IMAGE):$(VERSION)
	docker push $(DOCKER_IMAGE):latest
	@echo "✅ Docker image pushed"

docker-run: docker-build
	@echo "🐳 Running Docker container..."
	docker run --rm -it $(DOCKER_IMAGE):$(VERSION)

# Release targets
release-check:
	@echo "🚀 Checking release readiness..."
	@if [ -z "$(VERSION)" ]; then echo "❌ VERSION not set"; exit 1; fi
	@if ! git diff --quiet; then echo "❌ Git working directory not clean"; exit 1; fi
	@echo "✅ Release check passed"

release-prepare: release-check clean fmt vet lint test
	@echo "🚀 Preparing release $(VERSION)..."
	@git tag -a $(VERSION) -m "Release $(VERSION)"
	@echo "✅ Release $(VERSION) prepared. Push with: git push origin $(VERSION)"

checksums:
	@echo "🔐 Generating checksums..."
	@if [ ! -d "$(DIST_DIR)" ]; then echo "❌ No dist directory. Run 'make build-all' first"; exit 1; fi
	@cd $(DIST_DIR) && sha256sum * > checksums.txt
	@echo "✅ Checksums generated: $(DIST_DIR)/checksums.txt"

# Show help
help:
	@echo "Ukrainian Voice Transcriber - Build Commands"
	@echo "==========================================="
	@echo ""
	@echo "Building:"
	@echo "  make build           - Build binary for current platform"
	@echo "  make build-all       - Build for all platforms"
	@echo "  make build-static    - Build static binary (for Docker/Alpine)"
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
	@echo "Docker:"
	@echo "  make docker-build    - Build Docker image"
	@echo "  make docker-push     - Push Docker image"
	@echo "  make docker-run      - Run Docker container"
	@echo ""
	@echo "Release:"
	@echo "  make release-check   - Check release readiness"
	@echo "  make release-prepare - Prepare release (tag)"
	@echo "  make checksums       - Generate checksums"
	@echo ""
	@echo "Utilities:"
	@echo "  make setup           - Run setup check"
	@echo "  make demo            - Run demo (needs example.mp4)"
	@echo "  make clean           - Clean build artifacts"
	@echo "  make help            - Show this help"

# Development shortcuts
dev: deps build

# Release preparation
release: clean deps test build-all
	@echo "🚀 Release ready!"
	@ls -la $(BINARY_NAME)-*