# Build stage
FROM golang:1.24-alpine AS builder

# Set build arguments
ARG VERSION=dev
ARG BUILD_DATE=unknown
ARG GIT_COMMIT=unknown

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build \
    -a -installsuffix cgo \
    -ldflags="-w -s -X 'main.version=${VERSION}' -X 'main.buildDate=${BUILD_DATE}' -X 'main.gitCommit=${GIT_COMMIT}'" \
    -o ukrainian-voice-transcriber \
    ./cmd/transcriber

# Runtime stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    ffmpeg \
    tzdata \
    && rm -rf /var/cache/apk/*

# Create non-root user
RUN addgroup -g 1001 -S appuser && \
    adduser -u 1001 -S appuser -G appuser

# Set working directory
WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/ukrainian-voice-transcriber .

# Copy documentation
COPY --from=builder /app/README.md .
COPY --from=builder /app/docs/ ./docs/

# Change ownership
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Set environment variables
ENV PATH="/app:${PATH}"

# Expose any necessary ports (if needed in future)
# EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD ./ukrainian-voice-transcriber version || exit 1

# Default command
ENTRYPOINT ["./ukrainian-voice-transcriber"]
CMD ["--help"]