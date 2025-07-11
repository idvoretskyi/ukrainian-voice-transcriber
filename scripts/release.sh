#!/bin/bash

# Ukrainian Voice Transcriber - Release Script
# Usage: ./scripts/release.sh [patch|minor|major] [version]

set -e

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

# Parse arguments
RELEASE_TYPE=${1:-"minor"}
VERSION=${2:-""}

log_info "Ukrainian Voice Transcriber Release Script"
log_info "=========================================="

# Validate release type
if [[ ! "$RELEASE_TYPE" =~ ^(patch|minor|major)$ ]]; then
    log_error "Invalid release type. Use: patch, minor, or major"
    exit 1
fi

# Get current version from git tags
CURRENT_VERSION=$(git describe --tags --abbrev=0 2>/dev/null | sed 's/^v//' || echo "0.0.0")
log_info "Current version: $CURRENT_VERSION"

# Calculate next version if not provided
if [[ -z "$VERSION" ]]; then
    IFS='.' read -r -a VERSION_PARTS <<< "$CURRENT_VERSION"
    MAJOR=${VERSION_PARTS[0]:-0}
    MINOR=${VERSION_PARTS[1]:-0}
    PATCH=${VERSION_PARTS[2]:-0}
    
    case $RELEASE_TYPE in
        major)
            MAJOR=$((MAJOR + 1))
            MINOR=0
            PATCH=0
            ;;
        minor)
            MINOR=$((MINOR + 1))
            PATCH=0
            ;;
        patch)
            PATCH=$((PATCH + 1))
            ;;
    esac
    
    VERSION="$MAJOR.$MINOR.$PATCH"
fi

log_info "Next version: $VERSION"

# Validate version format
if [[ ! "$VERSION" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    log_error "Invalid version format. Use semantic versioning (e.g., 1.2.0)"
    exit 1
fi

# Check if tag already exists
if git rev-parse "v$VERSION" >/dev/null 2>&1; then
    log_error "Tag v$VERSION already exists"
    exit 1
fi

# Confirm release
log_warning "About to create release v$VERSION"
read -p "Continue? (y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    log_info "Release cancelled"
    exit 0
fi

# Check working directory is clean
if [[ -n $(git status --porcelain) ]]; then
    log_error "Working directory is not clean. Commit or stash changes first."
    exit 1
fi

# Create and push tag
log_info "Creating tag v$VERSION..."
git tag -a "v$VERSION" -m "Release v$VERSION"

log_info "Pushing tag to origin..."
git push origin "v$VERSION"

log_success "Release v$VERSION created successfully!"
log_info "🚀 GoReleaser will now build and publish the release"
log_info "📦 Check progress at: https://github.com/idvoretskyi/ukrainian-voice-transcriber/actions"
log_info "🔗 Release will be available at: https://github.com/idvoretskyi/ukrainian-voice-transcriber/releases/tag/v$VERSION"