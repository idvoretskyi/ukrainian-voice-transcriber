# OAuth Setup Guide

This guide explains how to set up OAuth credentials for the Ukrainian Voice Transcriber to enable the zero-configuration authentication flow.

## Overview

The transcriber can use OAuth for authentication, which provides these benefits:
- **No manual credential files** - Users authenticate through browser
- **Automatic token refresh** - Handles expired tokens seamlessly  
- **Secure storage** - Tokens stored locally with proper permissions
- **Easy revocation** - Users can revoke access anytime

## Setting Up OAuth Credentials

### 1. Create OAuth Credentials

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Select or create a project
3. Navigate to **APIs & Services > Credentials**
4. Click **+ CREATE CREDENTIALS > OAuth client ID**
5. Choose **Desktop application**
6. Set name: "Ukrainian Voice Transcriber"
7. Download the credentials JSON file

### 2. Configure OAuth in Code

**IMPORTANT**: Replace the placeholder credentials in `internal/auth/oauth.go`:

```go
config := &oauth2.Config{
    ClientID:     "your-actual-client-id.apps.googleusercontent.com",
    ClientSecret: "your-actual-client-secret", 
    RedirectURL:  "http://localhost:8080/callback",
    Scopes: []string{
        scopeCloudPlatform,
        scopeDriveReadonly,
    },
    Endpoint: google.Endpoint,
}
```

⚠️ **Security Note**: The current code contains placeholder values that will not work. You must replace them with your actual OAuth credentials from Google Cloud Console.

### 3. Enable Required APIs

Ensure these APIs are enabled in your Google Cloud project:
- **Cloud Speech-to-Text API**
- **Cloud Storage API** 
- **Cloud Resource Manager API** (for project detection)

### 4. Configure Consent Screen

1. Go to **APIs & Services > OAuth consent screen**
2. Choose **External** (unless using Google Workspace)
3. Fill required fields:
   - **Application name**: Ukrainian Voice Transcriber
   - **User support email**: Your email
   - **Developer contact information**: Your email
4. Add scopes:
   - `https://www.googleapis.com/auth/cloud-platform`
   - `https://www.googleapis.com/auth/drive.readonly`
5. Add test users (during development)

## OAuth Flow Details

### Authentication Process

1. **User runs**: `ukrainian-voice-transcriber auth`
2. **Browser opens** to Google authentication page
3. **User grants permissions** for Speech-to-Text and Cloud Storage
4. **Callback received** at `http://localhost:8080/callback`
5. **Token saved** to `.transcriber-token.json`
6. **Future requests** use saved token automatically

### Token Management

- **Automatic refresh**: Expired tokens are refreshed transparently
- **Secure storage**: Tokens stored with 0600 permissions
- **Project detection**: Automatically detects available Google Cloud projects
- **API verification**: Checks that required APIs are enabled

### Fallback to Service Account

If OAuth fails, the application falls back to service account authentication:

```
1. Try OAuth authentication
2. If no token found, try service account
3. If neither available, show helpful error message
```

## Security Considerations

### OAuth Client Type

- Uses **Desktop Application** OAuth client type
- Safe to embed Client ID in binary (industry standard)
- Client Secret provides additional security layer
- Redirect URI restricted to localhost

### Token Storage

- Tokens stored in `.transcriber-token.json`
- File permissions set to 0600 (owner read/write only)
- Contains access token, refresh token, and expiry
- Automatically excluded from git via `.gitignore`

### Scope Permissions

- **Cloud Platform**: Minimal scope for Speech-to-Text and Storage
- **Drive Readonly**: Optional, for future Google Drive integration
- No sensitive scopes (Gmail, Calendar, etc.)

## Troubleshooting

### Common Issues

**"Invalid client ID"**
- Verify Client ID in `oauth.go` matches Google Cloud Console
- Ensure OAuth client type is "Desktop Application"

**"Redirect URI mismatch"**
- Verify redirect URI is exactly `http://localhost:8080/callback`
- Check for typos in Google Cloud Console configuration

**"Access blocked"**
- Add your email as test user in OAuth consent screen
- Publish app for production use

**"API not enabled"**
- Enable Speech-to-Text and Cloud Storage APIs
- Wait a few minutes for APIs to propagate

### Debug Mode

Enable verbose logging to debug OAuth issues:

```bash
./ukrainian-voice-transcriber auth --verbose
```

### Manual Token Inspection

View token details:

```bash
./ukrainian-voice-transcriber auth --status
```

Revoke and re-authenticate:

```bash
./ukrainian-voice-transcriber auth --revoke
./ukrainian-voice-transcriber auth
```

## Production Deployment

### Publishing OAuth App

For production deployment:

1. **Complete OAuth consent screen** with all required information
2. **Add privacy policy** and terms of service URLs
3. **Submit for verification** if using sensitive scopes
4. **Update app status** from "Testing" to "In production"

### Distribution Considerations

- **Client ID can be public** (standard for desktop apps)
- **Client Secret should be protected** but can be embedded
- **Users authenticate with their own Google accounts**
- **No shared credentials** between users

This OAuth implementation provides a seamless authentication experience while maintaining security best practices for a desktop application.