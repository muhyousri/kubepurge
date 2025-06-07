# Dependency Management

This document outlines how dependencies are managed in the kubepurge project.

## Overview

We use automated dependency management to keep the project secure and up-to-date while maintaining stability.

## Tools Used

### Dependabot
- **Purpose**: Automated dependency updates
- **Configuration**: `.github/dependabot.yml`
- **Schedule**: Weekly updates on different days
- **Scope**: Go modules, Docker images, GitHub Actions

### Security Scanning
- **Trivy**: Container vulnerability scanning
- **Gosec**: Go static security analysis  
- **govulncheck**: Go vulnerability database checking
- **Nancy**: OSS Index vulnerability scanning
- **TruffleHog**: Secrets detection

## Update Strategy

### Automatic Updates

#### Patch Updates (Auto-merge)
- **Go modules**: Patch versions (x.y.Z)
- **GitHub Actions**: Patch and minor versions for core actions
- **Testing dependencies**: Ginkgo/Gomega minor updates

#### Manual Review Required
- **Major versions**: Breaking changes possible
- **Kubernetes dependencies**: Require careful testing
- **Security tools**: Need validation
- **Docker base images**: Verify compatibility

### Grouping Strategy

Dependencies are grouped to reduce PR noise:

1. **kubernetes-deps**: All k8s.io and sigs.k8s.io packages
2. **testing-deps**: Ginkgo and Gomega test frameworks  
3. **actions-core**: Core GitHub Actions
4. **docker-actions**: Docker-related actions
5. **security-actions**: Security scanning actions

## Review Process

### Patch Updates
1. **Automated CI**: Must pass all tests
2. **Security scan**: No new vulnerabilities
3. **Auto-merge**: If all checks pass

### Minor Updates  
1. **Automated CI**: Must pass all tests
2. **Manual review**: Check for behavioral changes
3. **Testing**: Run E2E tests
4. **Approval**: Maintainer approval required

### Major Updates
1. **Breaking changes**: Review changelog thoroughly
2. **Compatibility**: Check Kubernetes version support
3. **Testing**: Extended testing required
4. **Documentation**: Update if APIs change
5. **Staged rollout**: Consider gradual deployment

## Security Considerations

### Vulnerability Response

1. **Critical/High severity**:
   - Immediate update required
   - Emergency release if needed
   - Security advisory issued

2. **Medium severity**:
   - Update within 7 days
   - Include in next regular release

3. **Low severity**:
   - Update in next maintenance cycle
   - Monitor for exploitation

### Supply Chain Security

- **Verification**: Only trusted sources
- **Checksums**: Verify package integrity  
- **Signatures**: Validate signed packages (when available)
- **SBOM**: Generate Software Bill of Materials
- **Provenance**: Track build origins

## Monitoring

### Automated Checks
- **Daily**: Vulnerability scanning
- **Weekly**: Dependency updates
- **On PR**: Security analysis
- **On release**: Full security audit

### Alerts
- **High/Critical vulnerabilities**: Immediate notification
- **Failed security scans**: Investigation required
- **License violations**: Legal review needed
- **Supply chain attacks**: Emergency response

## Manual Dependency Updates

When manual updates are needed:

```bash
# Update specific dependency
go get -u github.com/example/package@v1.2.3

# Update all dependencies
go get -u ./...

# Tidy modules
go mod tidy

# Verify no vulnerabilities
govulncheck ./...

# Run tests
make test

# Update go.sum
go mod verify
```

## Kubernetes Dependencies

Special considerations for Kubernetes-related packages:

### Version Compatibility
- **Kubernetes**: Support N and N-1 versions
- **controller-runtime**: Match Kubernetes versions
- **client-go**: Compatible with target K8s versions
- **API machinery**: Consistent version across packages

### Update Strategy
- **Batch updates**: Update K8s deps together
- **Testing**: Extensive compatibility testing
- **Documentation**: Update supported versions
- **Migration**: Provide upgrade guides if needed

## CI/CD Integration

### Pre-merge Checks
- Security scanning
- License compliance
- Vulnerability assessment
- Test execution
- Build verification

### Post-merge Actions
- Update security baseline
- Regenerate SBOM
- Update dependency graphs
- Monitor for issues

## Troubleshooting

### Common Issues

1. **Conflicting versions**:
   ```bash
   go mod why -m github.com/example/package
   go mod graph | grep package
   ```

2. **Indirect dependency issues**:
   ```bash
   go mod edit -replace github.com/old/package=github.com/new/package@version
   ```

3. **Version constraints**:
   ```bash
   go mod edit -require github.com/example/package@version
   ```

### Recovery Procedures

1. **Rollback problematic update**:
   ```bash
   git revert <commit-hash>
   go mod tidy
   ```

2. **Pin to working version**:
   ```bash
   go mod edit -require github.com/example/package@v1.2.3
   ```

3. **Clean module cache**:
   ```bash
   go clean -modcache
   go mod download
   ```

## Best Practices

### For Maintainers
- Review dependency changes carefully
- Test updates in development environment
- Monitor security advisories
- Keep documentation updated
- Communicate breaking changes

### For Contributors  
- Don't update dependencies unless necessary
- Include justification for updates
- Test thoroughly before submitting
- Follow semantic versioning
- Document any breaking changes

## Resources

- [Go Modules Reference](https://go.dev/ref/mod)
- [Dependabot Documentation](https://docs.github.com/en/code-security/dependabot)
- [Kubernetes API Compatibility](https://kubernetes.io/docs/reference/using-api/api-concepts/)
- [Go Security Policy](https://go.dev/security)
- [NIST Vulnerability Database](https://nvd.nist.gov/)