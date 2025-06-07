# Security Policy

## Supported Versions

We actively support the following versions of kubepurge:

| Version | Supported          |
| ------- | ------------------ |
| 0.0.x   | :white_check_mark: |
| < 0.0.1 | :x:                |

## Reporting a Vulnerability

We take security vulnerabilities seriously. If you discover a security issue, please report it responsibly.

### How to Report

1. **Do NOT create a public GitHub issue**
2. **Email**: Send details to [security contact email]
3. **Include**: 
   - Description of the vulnerability
   - Steps to reproduce
   - Potential impact
   - Suggested fix (if any)

### What to Expect

- **Acknowledgment**: Within 24 hours
- **Initial Assessment**: Within 72 hours  
- **Regular Updates**: Every 7 days until resolved
- **Resolution Timeline**: Varies by severity
  - Critical: 1-7 days
  - High: 7-14 days
  - Medium: 14-30 days
  - Low: 30-90 days

## Security Measures

### Container Security

- **Non-root execution**: Controller runs as UID 65532
- **Read-only filesystem**: Root filesystem is read-only
- **Minimal base**: Distroless image reduces attack surface
- **No privilege escalation**: Security contexts prevent escalation
- **Capability dropping**: All capabilities dropped

### RBAC & Permissions

- **Principle of least privilege**: Minimal required permissions
- **Namespace isolation**: Recommended deployment in dedicated namespace
- **Resource-specific permissions**: Only required resource types
- **No cluster-admin**: Never requires cluster-admin privileges

### Supply Chain Security

- **Dependency scanning**: Automated with Dependabot
- **Container scanning**: Trivy vulnerability scanning
- **SBOM generation**: Software Bill of Materials included
- **Signed releases**: Container images and binaries (planned)
- **Provenance**: Build provenance attestation (planned)

### Runtime Security

- **Input validation**: All CRD fields validated
- **Error handling**: Graceful handling of edge cases
- **Resource limits**: Memory and CPU limits enforced
- **Network policies**: Recommended network restrictions
- **Audit logging**: Kubernetes audit logging compatible

## Security Best Practices

### Deployment

1. **Namespace Isolation**:
   ```bash
   kubectl create namespace kubepurge-system
   ```

2. **Network Policies**:
   ```yaml
   apiVersion: networking.k8s.io/v1
   kind: NetworkPolicy
   metadata:
     name: kubepurge-deny-all
     namespace: kubepurge-system
   spec:
     podSelector: {}
     policyTypes:
     - Ingress
     - Egress
   ```

3. **Resource Quotas**:
   ```yaml
   apiVersion: v1
   kind: ResourceQuota
   metadata:
     name: kubepurge-quota
     namespace: kubepurge-system
   spec:
     hard:
       requests.cpu: 100m
       requests.memory: 128Mi
       limits.cpu: 500m
       limits.memory: 256Mi
   ```

### Operation

1. **Resource Protection**: Use exclusion labels
   ```yaml
   metadata:
     labels:
       kubepurge.xyz/exclude: "true"
   ```

2. **Monitoring**: Set up alerts for:
   - Unexpected resource deletions
   - Controller errors or crashes
   - Failed purge operations
   - Security events

3. **Backup Strategy**: 
   - Regular cluster backups
   - Test restore procedures
   - Document critical resources

4. **Testing**:
   - Test policies in development first
   - Use dry-run mode (when available)
   - Monitor logs during initial deployment

## Vulnerability Disclosure

### Public Disclosure

- Security vulnerabilities will be disclosed publicly after:
  - Fix is available and tested
  - Reasonable time for users to update
  - Coordination with security researchers

### CVE Process

- For qualifying vulnerabilities, we will:
  - Request CVE assignment
  - Provide CVSS scoring
  - Include in security advisories
  - Update documentation

## Security Updates

### Notification Channels

- **GitHub Security Advisories**: Primary channel
- **Release Notes**: Include security fixes
- **Documentation**: Security best practices updates

### Update Recommendations

- **Critical/High**: Update immediately
- **Medium**: Update within 30 days
- **Low**: Update during next maintenance window

## Contact

For security-related questions or concerns:
- **Issues**: Use GitHub Issues for non-sensitive topics
- **Discussions**: GitHub Discussions for general security questions
- **Vulnerabilities**: Follow responsible disclosure process above

---

**Last Updated**: January 2025