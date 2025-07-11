# Repository Overview

## 🏗️ GitHub Configuration Structure

This repository uses GitHub's advanced features for automated management, security, and collaboration.

### 📁 File Structure

```
.github/
├── workflows/                 # GitHub Actions workflows
│   ├── ci.yml                # Continuous integration
│   ├── security.yml          # Security analysis
│   └── goreleaser.yml        # Automated releases
├── ISSUE_TEMPLATE/           # Issue templates
│   ├── bug_report.yml        # Bug report template
│   ├── feature_request.yml   # Feature request template
│   ├── question.yml          # Question template
│   └── config.yml            # Issue template configuration
├── codeql/                   # CodeQL configuration
│   └── codeql-config.yml     # Security scanning config
├── settings.yml              # Repository settings as code
├── dependabot.yml            # Dependency automation
├── CODEOWNERS                # Code review assignments
├── PULL_REQUEST_TEMPLATE.md  # PR template
└── REPOSITORY_OVERVIEW.md    # This file
```

### 🔧 Repository Settings (settings.yml)

**Key Features:**
- **Branch Protection**: Main branch protected with required reviews
- **Labels**: Comprehensive labeling system for issues and PRs
- **Milestones**: Release planning and tracking
- **Security**: Automated vulnerability alerts and security fixes
- **Environments**: Production and staging deployment environments

**Branch Protection Rules:**
- Require 1 approving review
- Require status checks (tests, linting, security)
- Dismiss stale reviews
- Require conversation resolution

### 🤖 Automation (dependabot.yml)

**Automated Updates:**
- **Go modules**: Weekly updates on Mondays
- **GitHub Actions**: Weekly updates on Mondays  
- **Docker**: Weekly base image updates
- **Security**: Automatic security patches

**Configuration:**
- Max 10 open PRs for Go modules
- Max 5 open PRs for GitHub Actions/Docker
- Auto-assign to maintainers
- Semantic commit messages

### 👥 Code Reviews (CODEOWNERS)

**Automatic Assignment:**
- All code changes require review from `@idvoretskyi`
- Component-specific ownership
- Documentation changes tracked separately
- Security files require special attention

### 🎯 Issue Templates

**Template Types:**
1. **Bug Report** - Structured bug reporting with environment details
2. **Feature Request** - Feature proposals with priority and complexity
3. **Question** - General questions with categorization
4. **Config** - Links to documentation and community resources

### 🔒 Security Configuration

**CodeQL Analysis:**
- Extended security queries
- Focus on application code
- Exclude tests and documentation
- Upload results to Security tab

**Multi-layered Security:**
- CodeQL (semantic analysis)
- Trivy (vulnerability scanning)
- gosec (Go security analysis)
- govulncheck (Go vulnerability database)
- Dependabot (dependency security)

### 🚀 Workflows

**CI/CD Pipeline:**
1. **Continuous Integration** - Tests, linting, building
2. **Security Analysis** - Multiple security scanning tools
3. **Automated Releases** - Professional releases with GoReleaser

**Workflow Triggers:**
- Push to main branch
- Pull requests to main branch
- Git tags (for releases)
- Weekly scheduled scans

### 🏷️ Labels System

**Categories:**
- **Priority**: high, medium, low
- **Type**: bug, enhancement, documentation
- **Component**: cli, auth, transcription, build, docker
- **Status**: help wanted, good first issue, wontfix
- **Development**: dependencies, ci/cd, testing, performance

### 📊 Repository Features

**Enabled Features:**
- Issues and Projects
- Security advisories
- Dependabot alerts
- Automated security fixes
- Vulnerability alerts
- Downloads

**Disabled Features:**
- Wiki (using docs/ directory instead)
- GitHub Pages (not needed for CLI tool)

### 🔄 Deployment Environments

**Production Environment:**
- Protected branch deployment
- No wait time (immediate deployment)
- Environment variables configured

**Staging Environment:**
- Multi-branch deployment (main, develop)
- Testing environment for pre-release

### 💡 Best Practices

**This repository follows:**
- Semantic versioning
- Conventional commits
- Automated testing
- Security-first approach
- Comprehensive documentation
- Community-friendly issue tracking

**Maintenance:**
- Weekly dependency updates
- Regular security scans
- Automated code quality checks
- Professional release management

---

**For more information:**
- [Build System Documentation](../docs/BUILD_SYSTEM.md)
- [Contributing Guidelines](../CONTRIBUTING.md)
- [User Manual](../docs/USER_MANUAL_EN.md)
- [Security Policy](../SECURITY.md)