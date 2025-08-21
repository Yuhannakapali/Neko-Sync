# Neko-Sync Release Guide

This guide explains how to create and manage releases for the Neko-Sync backend application.

## Quick Release

### Option 1: Using the Release Script (Recommended)
```bash
# Create a new release
./scripts/release.sh v1.0.0
```

### Option 2: Manual Release
```bash
# Check if ready for release
make check-release-ready

# Create and push tag
make tag VERSION=v1.0.0
```

## Release Process Overview

When you create a new tag (starting with `v`), GitHub Actions automatically:

1. **Creates a GitHub Release** with auto-generated release notes
2. **Builds cross-platform binaries** for:
   - Linux (AMD64, ARM64)
   - macOS (AMD64, ARM64) 
   - Windows (AMD64)
3. **Publishes Docker images** to GitHub Container Registry
4. **Uploads release assets** (compressed binaries)

## Release Types

### Stable Releases
- Format: `v1.0.0`, `v2.1.3`, etc.
- Automatically marked as latest release
- Published to all distribution channels

### Pre-releases
- Format: `v1.0.0-beta`, `v2.0.0-rc1`, etc.
- Marked as pre-release on GitHub
- Useful for testing before stable release

## Available Make Targets

```bash
# Check what would be released
make release-dry-run

# Generate release notes
make release-notes

# Check if repository is ready for release
make check-release-ready

# Build release artifacts locally
make release

# Create and push a tag
make tag VERSION=v1.0.0
```

## Installation for Users

### Binary Download
Users can download pre-built binaries from the [Releases page](https://github.com/Yuhannakapali/Neko-Sync/releases).

#### Quick Install (Linux/macOS)
```bash
# Auto-detect platform and install
curl -L -o nekosync https://github.com/Yuhannakapali/Neko-Sync/releases/latest/download/nekosync-$(uname -s | tr '[:upper:]' '[:lower:]')-$(uname -m | sed 's/x86_64/amd64/')
chmod +x nekosync
sudo mv nekosync /usr/local/bin/
```

#### Manual Download
1. Go to [Releases](https://github.com/Yuhannakapali/Neko-Sync/releases)
2. Download the appropriate binary for your platform:
   - `nekosync-linux-amd64.tar.gz` - Linux x64
   - `nekosync-linux-arm64.tar.gz` - Linux ARM64
   - `nekosync-darwin-amd64.tar.gz` - macOS Intel
   - `nekosync-darwin-arm64.tar.gz` - macOS Apple Silicon
   - `nekosync-windows-amd64.zip` - Windows x64

### Docker
```bash
# Latest stable release
docker pull ghcr.io/yuhannakapali/neko-sync:latest

# Specific version
docker pull ghcr.io/yuhannakapali/neko-sync:v1.0.0

# Run the container
docker run -p 8080:8080 ghcr.io/yuhannakapali/neko-sync:latest
```

## Release Checklist

Before creating a release:

- [ ] All tests pass (`make test`)
- [ ] Code is on `main` branch
- [ ] Working directory is clean
- [ ] Version follows semantic versioning
- [ ] CHANGELOG is updated (optional)
- [ ] Documentation is up to date

## Semantic Versioning

We follow [Semantic Versioning](https://semver.org/):

- **MAJOR** version (`v2.0.0`): Incompatible API changes
- **MINOR** version (`v1.1.0`): New functionality, backward compatible
- **PATCH** version (`v1.0.1`): Bug fixes, backward compatible

## Troubleshooting

### Release Workflow Failed
1. Check the [Actions tab](https://github.com/Yuhannakapali/Neko-Sync/actions)
2. Look for the failed workflow run
3. Check the logs for specific errors
4. Common issues:
   - Build failures: Check Go compilation errors
   - Docker build failures: Check Dockerfile syntax
   - Permission issues: Ensure repository has proper secrets

### Tag Already Exists
```bash
# Delete local tag
git tag -d v1.0.0

# Delete remote tag
git push origin --delete v1.0.0

# Recreate the tag
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

## Advanced Configuration

### Custom Release Notes
Create a `.github/release.yml` file to customize release note generation:

```yaml
changelog:
  exclude:
    labels:
      - ignore-for-release
  categories:
    - title: Breaking Changes 🛠
      labels:
        - breaking-change
    - title: New Features 🎉
      labels:
        - enhancement
    - title: Bug Fixes 🐛
      labels:
        - bug
    - title: Other Changes
      labels:
        - "*"
```

### Repository Secrets
The release workflow uses these GitHub secrets:
- `GITHUB_TOKEN` - Automatically provided by GitHub
- Additional secrets can be added for external services

## Support

For release-related issues:
1. Check existing [Issues](https://github.com/Yuhannakapali/Neko-Sync/issues)
2. Create a new issue with the `release` label
3. Include relevant workflow logs and error messages
