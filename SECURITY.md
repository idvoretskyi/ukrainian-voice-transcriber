# Security Policy

## Supported Versions

No stable releases have been cut yet. Security fixes are applied to the
`main` branch and will be included in the first and all subsequent releases.

| Version | Supported |
|---------|-----------|
| `main` (latest) | Yes |
| older commits | No — please update to `main` |

## Reporting a Vulnerability

**Please do not open a public GitHub issue for security vulnerabilities.**

Use one of the following channels:

### GitHub Private Vulnerability Reporting (preferred)

Open a [private security advisory](https://github.com/idvoretskyi/ukrainian-voice-transcriber/security/advisories/new)
directly in this repository. GitHub keeps the report private until a fix is
published.

### Email

Send details to **ihor@dvoretskyi.com** with the subject line:
`[SECURITY] ukrainian-voice-transcriber: <brief description>`

Please include:

- A description of the vulnerability and its potential impact
- Steps to reproduce or a proof-of-concept
- Any suggested fix or mitigation, if you have one

## Response Timeline

| Milestone | Target |
|-----------|--------|
| Initial acknowledgement | Within **48 hours** |
| Severity assessment | Within **5 business days** |
| Fix for Critical / High | Within **7 days** of confirmation |
| Fix for Medium / Low | Within **30 days** of confirmation |
| Public disclosure | After fix is released, coordinated with reporter |

## Automated Security Scanning

The following tools run automatically on every push and pull request:

| Tool | What it checks |
|------|---------------|
| [Trivy](https://github.com/aquasecurity/trivy) v0.35 | Dependency vulnerabilities (CRITICAL, HIGH, MEDIUM) |
| [gosec](https://github.com/securego/gosec) | Go source code security issues |
| [govulncheck](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck) | Known vulnerabilities in Go module dependencies |
| [CodeQL](https://codeql.github.com) | Static analysis for security and quality issues |
| [Dependency Review](https://github.com/actions/dependency-review-action) | New vulnerable dependencies introduced in PRs |

SARIF results from Trivy, gosec, and CodeQL are uploaded to
[GitHub Code Scanning](https://github.com/idvoretskyi/ukrainian-voice-transcriber/security/code-scanning).

## Scope

The following are **in scope** for security reports:

- Vulnerabilities in the transcriber binary itself
- Unsafe handling of user-supplied file paths or filenames
- Credential or token exposure in logs or output
- Dependency vulnerabilities with a realistic attack vector in this tool's use case

The following are **out of scope**:

- Vulnerabilities in Google Cloud / Vertex AI infrastructure
- Issues requiring physical access to the machine running the binary
- Social engineering attacks
