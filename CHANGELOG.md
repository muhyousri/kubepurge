# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.0.1] - 2025-01-07

### Added
- Initial release of kubepurge Kubernetes controller
- **Core Controller Features:**
  - PurgePolicy custom resource for defining cleanup schedules
  - PurgeStatus custom resource for tracking operations
  - Cron-based scheduling for automatic resource deletion
  - Support for pods, deployments, services, configmaps, and secrets
  - Label-based exclusion mechanism (`kubepurge.xyz/exclude: "true"`)
  
- **Resource Management:**
  - Targeted namespace cleanup
  - Flexible resource type selection
  - Safe deletion with error handling
  - Resource count tracking and reporting
  
- **Security & RBAC:**
  - Minimal required permissions
  - Non-root container execution
  - Read-only root filesystem
  - Distroless base image for security
  
- **Observability:**
  - Structured logging with context
  - Health check endpoints (liveness/readiness)
  - Prometheus metrics support
  - Detailed status reporting
  
- **CI/CD Pipeline:**
  - Automated testing with GitHub Actions
  - Multi-platform container builds (linux/amd64, linux/arm64)
  - Security scanning with Trivy and Gosec
  - Automated releases with GitHub releases
  - Cross-platform binary builds
  - SBOM generation
  
- **Documentation:**
  - Comprehensive README with usage examples
  - Release process documentation
  - API reference for custom resources
  - Security best practices guide
  - Development and contribution guidelines
  
- **Developer Tools:**
  - Release script for easy version management
  - Makefile with build, test, and lint targets
  - E2E testing with kind cluster
  - Kustomize configuration for deployment

### Technical Details
- Built with Go 1.22+ and Kubernetes controller-runtime
- Supports Kubernetes 1.25+
- Uses robfig/cron for scheduling
- Implements Kubernetes controller patterns
- Container images published to GitHub Container Registry

### Installation
- Quick install via kubectl apply
- Kustomize support for customization
- Complete installation manifests included
- Sample configurations provided

[Unreleased]: https://github.com/muhyousri/kubepurge/compare/v0.0.1...HEAD
[0.0.1]: https://github.com/muhyousri/kubepurge/releases/tag/v0.0.1