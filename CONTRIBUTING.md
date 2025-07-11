# Contributing to Ukrainian Voice Transcriber

🇺🇦 Thank you for your interest in contributing to the Ukrainian Voice Transcriber project!

## Quick Start

1. **Fork the repository**
2. **Clone your fork**: `git clone https://github.com/YOUR_USERNAME/ukrainian-voice-transcriber.git`
3. **Create a branch**: `git checkout -b feature/your-feature-name`
4. **Make your changes**
5. **Test your changes**: `make build && ./ukrainian-voice-transcriber setup`
6. **Commit and push**: `git commit -m "Add feature" && git push origin feature/your-feature-name`
7. **Open a Pull Request**

## Development Setup

### Prerequisites
- Go 1.24+
- FFmpeg
- Make

### Build and Test
```bash
# Install dependencies
go mod tidy

# Build the project
make build

# Run tests (when available)
go test ./...

# Test the binary
./ukrainian-voice-transcriber setup
```

## Contributing Guidelines

### Code Style
- Follow standard Go conventions
- Use `gofmt` to format your code
- Add comments for exported functions
- Keep functions focused and small

### Commit Messages
- Use clear, descriptive commit messages
- Start with a verb (Add, Fix, Update, Remove)
- Reference issues when applicable: `Fix #123: Resolve OAuth timeout`

### Pull Requests
- Provide a clear description of the changes
- Include steps to test the changes
- Link to related issues
- Ensure all checks pass

## Types of Contributions

### 🐛 Bug Reports
- Use the GitHub Issues template
- Include reproduction steps
- Provide system information (OS, Go version)
- Include relevant logs

### ✨ Feature Requests
- Describe the problem you're solving
- Explain the proposed solution
- Consider backward compatibility
- Discuss performance implications

### 📖 Documentation
- Fix typos and improve clarity
- Add missing documentation
- Update outdated information
- Translate to other languages

### 🔧 Code Contributions
- Bug fixes
- Performance improvements
- New features
- Refactoring

## Development Areas

### Priority Areas
1. **Error handling improvements**
2. **Performance optimizations**
3. **Additional audio formats support**
4. **Better OAuth user experience**
5. **Cross-platform compatibility**

### Technical Areas
- **Audio processing**: FFmpeg integration improvements
- **Cloud integration**: Additional cloud providers
- **Authentication**: OAuth flow enhancements
- **CLI**: User experience improvements
- **Testing**: Unit and integration tests

## Security Considerations

### Credentials
- Never commit real credentials
- Use placeholder values in code
- Document credential setup clearly
- Follow OAuth best practices

### Dependencies
- Keep dependencies up to date
- Review security advisories
- Minimize dependency footprint

## Testing

### Manual Testing
```bash
# Test basic functionality
./ukrainian-voice-transcriber setup
./ukrainian-voice-transcriber auth --status

# Test with sample video (if available)
./ukrainian-voice-transcriber transcribe sample.mp4
```

### Automated Testing
- Unit tests for core logic
- Integration tests for Google Cloud APIs
- CLI interface tests

## Documentation

### README Updates
- Keep usage examples current
- Update feature lists
- Maintain installation instructions

### Code Documentation
- Document public APIs
- Add inline comments for complex logic
- Update godoc comments

## Community Guidelines

### Be Respectful
- Use welcoming and inclusive language
- Respect different viewpoints and experiences
- Accept constructive criticism gracefully

### Be Helpful
- Help others learn and contribute
- Share knowledge and resources
- Provide constructive feedback

### Ukrainian Context
- This project supports Ukrainian content creators
- Contributions that improve Ukrainian language support are especially welcome
- Consider accessibility for non-technical users

## Getting Help

### Questions
- Check existing documentation first
- Search closed issues for similar questions
- Open a discussion for general questions
- Use issues for specific bug reports

### Contact
- GitHub Issues: Bug reports and feature requests
- GitHub Discussions: General questions and ideas
- Pull Requests: Code contributions

## Release Process

### Versioning
- Follow semantic versioning (SemVer)
- Tag releases with `git tag v1.2.3`
- Update version in code

### Release Notes
- Summarize new features
- List bug fixes
- Note breaking changes
- Include upgrade instructions

## Recognition

Contributors will be recognized in:
- README acknowledgments
- Release notes
- Git commit history

Thank you for helping make Ukrainian Voice Transcriber better! 🎉

---

**Слава Україні!** 🇺🇦