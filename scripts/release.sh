#!/bin/bash

# Release script for kubepurge
# Usage: ./scripts/release.sh v1.0.0

set -e

VERSION="$1"

if [ -z "$VERSION" ]; then
    echo "Usage: $0 <version>"
    echo "Example: $0 v1.0.0"
    exit 1
fi

# Validate version format
if [[ ! "$VERSION" =~ ^v[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9-]+)?$ ]]; then
    echo "Error: Version must be in format vX.Y.Z or vX.Y.Z-suffix (e.g., v1.0.0, v1.0.0-rc1)"
    exit 1
fi

echo "Preparing release $VERSION..."

# Check if we're on main branch
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
if [ "$CURRENT_BRANCH" != "main" ]; then
    echo "Error: Must be on main branch to create a release (currently on $CURRENT_BRANCH)"
    exit 1
fi

# Check if working directory is clean
if [ -n "$(git status --porcelain)" ]; then
    echo "Error: Working directory is not clean. Please commit or stash changes."
    exit 1
fi

# Check if tag already exists
if git tag -l | grep -q "^$VERSION$"; then
    echo "Error: Tag $VERSION already exists"
    exit 1
fi

# Pull latest changes
echo "Pulling latest changes..."
git pull origin main

# Run tests
echo "Running tests..."
make test

# Generate and validate manifests
echo "Generating manifests..."
make manifests

# Check if manifests changed
if [ -n "$(git status --porcelain config/)" ]; then
    echo "Manifests have been updated. Please review and commit the changes:"
    git status config/
    exit 1
fi

# Create and push tag
echo "Creating and pushing tag $VERSION..."
git tag -a "$VERSION" -m "Release $VERSION"
git push origin "$VERSION"

echo "✅ Release $VERSION has been created!"
echo ""
echo "The GitHub Actions workflow will now:"
echo "1. Create a GitHub release"
echo "2. Build and push container images"
echo "3. Generate release artifacts"
echo "4. Run security scans"
echo ""
echo "Monitor the progress at: https://github.com/$(git remote get-url origin | sed 's/.*github.com[:/]//' | sed 's/\.git$//')/actions"