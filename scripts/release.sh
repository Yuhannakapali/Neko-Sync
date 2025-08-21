#!/bin/bash

# Neko-Sync Release Script
# Usage: ./scripts/release.sh [version]
# Example: ./scripts/release.sh v1.0.0

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if version is provided
if [ -z "$1" ]; then
    log_error "Version is required!"
    echo "Usage: $0 <version>"
    echo "Example: $0 v1.0.0"
    exit 1
fi

VERSION="$1"

# Validate version format (should start with 'v' followed by semantic version)
if [[ ! $VERSION =~ ^v[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9]+)?$ ]]; then
    log_error "Invalid version format. Use semantic versioning: v1.0.0"
    exit 1
fi

log_info "Starting release process for version: $VERSION"

# Check if we're on the main branch
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
if [ "$CURRENT_BRANCH" != "main" ]; then
    log_warning "Not on main branch (current: $CURRENT_BRANCH)"
    read -p "Continue anyway? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        log_info "Release cancelled"
        exit 0
    fi
fi

# Check if working directory is clean
if [ -n "$(git status --porcelain)" ]; then
    log_error "Working directory is not clean. Please commit or stash changes."
    git status --short
    exit 1
fi

# Check if tag already exists
if git tag -l | grep -q "^$VERSION$"; then
    log_error "Tag $VERSION already exists!"
    exit 1
fi

# Pull latest changes
log_info "Pulling latest changes..."
git pull origin main

# Run tests
log_info "Running tests..."
make test || {
    log_error "Tests failed! Please fix before releasing."
    exit 1
}

# Run release build
log_info "Building release artifacts..."
make release || {
    log_error "Release build failed!"
    exit 1
}

# Generate release notes
log_info "Generating release notes..."
make release-notes > RELEASE_NOTES.tmp

# Show release notes
echo ""
log_info "Release notes:"
echo "----------------------------------------"
cat RELEASE_NOTES.tmp
echo "----------------------------------------"
echo ""

# Confirm release
read -p "Proceed with release? (y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    log_info "Release cancelled"
    rm -f RELEASE_NOTES.tmp
    exit 0
fi

# Create and push tag
log_info "Creating and pushing tag..."
git tag -a "$VERSION" -m "Release $VERSION"
git push origin "$VERSION"

# Clean up
rm -f RELEASE_NOTES.tmp

log_success "Release $VERSION initiated!"
echo ""
log_info "Next steps:"
echo "1. GitHub Actions will automatically build and create the release"
echo "2. Monitor the release workflow: https://github.com/$(git remote get-url origin | sed 's/.*github.com[:\/]//' | sed 's/\.git$//')/actions"
echo "3. The release will be available at: https://github.com/$(git remote get-url origin | sed 's/.*github.com[:\/]//' | sed 's/\.git$//')/releases/tag/$VERSION"
echo ""
log_success "Release process completed!"
