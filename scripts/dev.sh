#!/bin/bash
set -e

# Development script for Ukrainian Voice Transcriber
# This script helps with common development tasks

PROJECT_NAME="ukrainian-voice-transcriber"
PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper functions
log_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

log_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

log_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

log_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Check if required tools are installed
check_dependencies() {
    log_info "Checking dependencies..."
    
    # Check Go
    if ! command -v go &> /dev/null; then
        log_error "Go is not installed. Please install Go 1.24 or later."
        exit 1
    fi
    
    # Check Go version
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    log_info "Go version: $GO_VERSION"
    
    # Check FFmpeg
    if ! command -v ffmpeg &> /dev/null; then
        log_warning "FFmpeg is not installed. Some features may not work."
        log_info "Install with: brew install ffmpeg (macOS) or apt install ffmpeg (Ubuntu)"
    else
        log_success "FFmpeg is installed"
    fi
    
    # Check Docker
    if ! command -v docker &> /dev/null; then
        log_warning "Docker is not installed. Docker-related commands will not work."
    else
        log_success "Docker is installed"
    fi
    
    # Check gcloud
    if ! command -v gcloud &> /dev/null; then
        log_warning "gcloud CLI is not installed. Authentication may require manual setup."
    else
        log_success "gcloud CLI is installed"
    fi
}

# Setup development environment
setup_dev() {
    log_info "Setting up development environment..."
    
    cd "$PROJECT_DIR"
    
    # Download dependencies
    log_info "Downloading Go dependencies..."
    go mod download
    go mod tidy
    
    # Install development tools
    log_info "Installing development tools..."
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
    go install golang.org/x/vuln/cmd/govulncheck@latest
    
    log_success "Development environment setup complete!"
}

# Run all tests
run_tests() {
    log_info "Running tests..."
    cd "$PROJECT_DIR"
    
    go test -v -race -coverprofile=coverage.out ./...
    
    if [ $? -eq 0 ]; then
        log_success "All tests passed!"
        
        # Generate coverage report
        go tool cover -html=coverage.out -o coverage.html
        log_info "Coverage report generated: coverage.html"
    else
        log_error "Some tests failed!"
        exit 1
    fi
}

# Run linting
run_lint() {
    log_info "Running linter..."
    cd "$PROJECT_DIR"
    
    if command -v golangci-lint &> /dev/null; then
        golangci-lint run
        log_success "Linting completed!"
    else
        log_error "golangci-lint not found. Run 'dev.sh setup' to install it."
        exit 1
    fi
}

# Run security scan
run_security() {
    log_info "Running security scan..."
    cd "$PROJECT_DIR"
    
    if command -v gosec &> /dev/null; then
        gosec ./...
        log_success "Security scan completed!"
    else
        log_error "gosec not found. Run 'dev.sh setup' to install it."
        exit 1
    fi
}

# Build the project
build_project() {
    log_info "Building project..."
    cd "$PROJECT_DIR"
    
    make build
    
    if [ $? -eq 0 ]; then
        log_success "Build completed!"
        ./ukrainian-voice-transcriber version
    else
        log_error "Build failed!"
        exit 1
    fi
}

# Run full development cycle
run_full_check() {
    log_info "Running full development check..."
    
    check_dependencies
    run_tests
    run_lint
    run_security
    build_project
    
    log_success "Full development check completed! 🎉"
}

# Show help
show_help() {
    echo "Ukrainian Voice Transcriber - Development Script"
    echo "=============================================="
    echo ""
    echo "Usage: ./scripts/dev.sh [command]"
    echo ""
    echo "Commands:"
    echo "  setup      - Set up development environment"
    echo "  check      - Check dependencies"
    echo "  test       - Run tests"
    echo "  lint       - Run linter"
    echo "  security   - Run security scan"
    echo "  build      - Build project"
    echo "  full       - Run full development cycle"
    echo "  help       - Show this help"
    echo ""
    echo "Examples:"
    echo "  ./scripts/dev.sh setup"
    echo "  ./scripts/dev.sh full"
    echo "  ./scripts/dev.sh test"
}

# Main script logic
case "${1:-help}" in
    setup)
        check_dependencies
        setup_dev
        ;;
    check)
        check_dependencies
        ;;
    test)
        run_tests
        ;;
    lint)
        run_lint
        ;;
    security)
        run_security
        ;;
    build)
        build_project
        ;;
    full)
        run_full_check
        ;;
    help|*)
        show_help
        ;;
esac