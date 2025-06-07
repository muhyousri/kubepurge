# Release Process

This document describes the release process for kubepurge.

## Automated Release Pipeline

The project uses GitHub Actions to automate the release process. When a new version tag is pushed, the following happens automatically:

### 1. Container Image Build & Push
- Multi-platform builds (linux/amd64, linux/arm64)
- Images pushed to GitHub Container Registry
- Images tagged with version and `latest`

### 2. Security Scanning
- Trivy vulnerability scanning
- Gosec static analysis
- SARIF results uploaded to GitHub Security tab

### 3. Release Artifacts
- Cross-platform binaries (Linux, macOS, Windows)
- SHA256 checksums for all binaries
- Complete Kubernetes installation manifest
- Sample configuration files

### 4. GitHub Release
- Automatic release notes generation
- Asset uploads
- Draft/pre-release handling for RC/beta versions

## Creating a Release

### Prerequisites
- Push access to the main branch
- Clean working directory on main branch
- All tests passing

### Process

1. **Use the release script (recommended):**
   ```bash
   ./scripts/release.sh v1.0.0
   ```

2. **Manual process:**
   ```bash
   # Ensure you're on main and up to date
   git checkout main
   git pull origin main
   
   # Run tests
   make test
   
   # Create and push tag
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```

### Version Format
- Use semantic versioning: `vMAJOR.MINOR.PATCH`
- Pre-releases: `v1.0.0-rc1`, `v1.0.0-beta1`, `v1.0.0-alpha1`

## Container Images

Released images are available at:
- `ghcr.io/muhyousri/kubepurge:latest` (latest stable)
- `ghcr.io/muhyousri/kubepurge:v1.0.0` (specific version)

## Installation

### Using the installation manifest:
```bash
kubectl apply -f https://github.com/muhyousri/kubepurge/releases/latest/download/kubepurge-install.yaml
```

### Using Kustomize:
```bash
kustomize build github.com/muhyousri/kubepurge/config/default | kubectl apply -f -
```

## Release Checklist

- [ ] All tests pass
- [ ] Documentation is up to date
- [ ] CHANGELOG.md is updated
- [ ] Version follows semantic versioning
- [ ] Tag is properly formatted
- [ ] GitHub release is created successfully
- [ ] Container images are built and pushed
- [ ] Security scans pass
- [ ] Installation manifest works correctly

## Troubleshooting

### Failed Image Build
- Check Dockerfile syntax
- Verify all dependencies are available
- Review build logs in GitHub Actions

### Failed Security Scan
- Review Trivy/Gosec reports
- Fix identified vulnerabilities
- Re-run the release process

### Missing Artifacts
- Check GitHub Actions workflow logs
- Verify permissions and secrets
- Ensure all jobs completed successfully

## Rollback Process

If a release needs to be rolled back:

1. **Delete the problematic tag:**
   ```bash
   git tag -d v1.0.0
   git push origin :refs/tags/v1.0.0
   ```

2. **Delete the GitHub release**
3. **Create a new patch release with fixes**