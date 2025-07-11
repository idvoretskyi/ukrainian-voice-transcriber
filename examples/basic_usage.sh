#!/bin/bash
# Basic usage examples for Ukrainian Voice Transcriber

echo "🎬 Ukrainian Voice Transcriber - Usage Examples"
echo "=============================================="

# Check if binary exists
if [ ! -f "../ukrainian-voice-transcriber" ]; then
    echo "❌ Binary not found. Please run 'make build' first."
    exit 1
fi

TRANSCRIBER="../ukrainian-voice-transcriber"

echo "1. Check setup and configuration:"
echo "   $TRANSCRIBER setup"
echo ""

echo "2. Basic transcription (output to stdout):"
echo "   $TRANSCRIBER transcribe video.mp4"
echo ""

echo "3. Save transcript to file:"
echo "   $TRANSCRIBER transcribe video.mp4 -o transcript.txt"
echo ""

echo "4. Verbose mode (show detailed progress):"
echo "   $TRANSCRIBER transcribe video.mp4 --verbose"
echo ""

echo "5. Quiet mode (only show results):"
echo "   $TRANSCRIBER transcribe video.mp4 --quiet"
echo ""

echo "6. Custom storage bucket:"
echo "   $TRANSCRIBER transcribe video.mp4 --bucket my-custom-bucket"
echo ""

echo "7. Show version:"
echo "   $TRANSCRIBER version"
echo ""

echo "8. Get help for specific command:"
echo "   $TRANSCRIBER transcribe --help"
echo ""

echo "📝 Batch processing example:"
cat << 'EOF'
# Process multiple video files
for video in *.mp4; do
    echo "Processing $video..."
    ./ukrainian-voice-transcriber transcribe "$video" -o "${video%.mp4}_transcript.txt"
done
EOF

echo ""
echo "🔧 Prerequisites:"
echo "• FFmpeg installed (brew install ffmpeg)"
echo "• Google Cloud service account JSON in current directory"
echo "• Speech-to-Text and Cloud Storage APIs enabled"