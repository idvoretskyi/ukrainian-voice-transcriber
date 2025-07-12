#!/bin/bash

# Function to add MIT header to Go files
add_header() {
    local file="$1"
    local package_comment="$2"
    
    # Create temp file with header
    cat > /tmp/header.go << 'EOF'
// Ukrainian Voice Transcriber
//
// Copyright (c) 2025 Ihor Dvoretskyi
//
// Licensed under MIT License

EOF
    
    # Add package comment
    echo "// $package_comment" >> /tmp/header.go
    
    # Find the package line and preserve everything from there
    grep -n "^package " "$file" | head -1 | cut -d: -f1 | read package_line
    tail -n +$package_line "$file" >> /tmp/header.go
    
    # Replace original file
    mv /tmp/header.go "$file"
}

# Fix all Go files that need headers
add_header "internal/auth/oauth.go" "Package auth provides OAuth authentication functionality."
add_header "internal/cli/auth.go" "Package cli provides command-line interface functionality."
add_header "internal/cli/root.go" "Package cli provides command-line interface functionality."
add_header "internal/cli/setup.go" "Package cli provides command-line interface functionality."
add_header "internal/cli/transcribe.go" "Package cli provides command-line interface functionality."
add_header "internal/cli/version.go" "Package cli provides command-line interface functionality."
add_header "internal/speech/service.go" "Package speech provides Google Cloud Speech-to-Text functionality."
add_header "internal/storage/service.go" "Package storage provides Google Cloud Storage functionality."
add_header "internal/transcriber/audio.go" "Package transcriber provides audio transcription functionality."
add_header "internal/transcriber/transcriber.go" "Package transcriber provides video transcription functionality."
add_header "pkg/config/config.go" "Package config provides configuration structures and utilities."

echo "Headers added to all files"