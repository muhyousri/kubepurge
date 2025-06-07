# kubepurge

**Kubernetes Resource Purging Controller**

[![CI](https://github.com/muhyousri/kubepurge/actions/workflows/ci.yml/badge.svg)](https://github.com/muhyousri/kubepurge/actions/workflows/ci.yml)
[![Release](https://github.com/muhyousri/kubepurge/actions/workflows/release.yml/badge.svg)](https://github.com/muhyousri/kubepurge/actions/workflows/release.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/muhyousri/kubepurge)](https://goreportcard.com/report/github.com/muhyousri/kubepurge)

Optimize costs for development environments by automatically cleaning Kubernetes resources on a schedule. kubepurge helps prevent resource waste and reduces cloud costs by purging unused resources in non-production environments.

## 🚀 Features

- **📅 Scheduled Cleanup**: Define cron-based schedules for automatic resource deletion
- **🎯 Targeted Purging**: Specify which namespaces and resource types to clean
- **🛡️ Safe Guards**: Label-based exclusion to protect critical resources
- **📊 Status Tracking**: Monitor purge operations with detailed status reporting
- **🔧 Flexible Configuration**: Support for pods, deployments, services, configmaps, and secrets
- **🔒 Security First**: RBAC-enabled with minimal required permissions

## 📦 Installation

### Quick Install

```bash
kubectl apply -f https://github.com/muhyousri/kubepurge/releases/latest/download/kubepurge-install.yaml
```

### Using Kustomize

```bash
kustomize build github.com/muhyousri/kubepurge/config/default | kubectl apply -f -
```

### Using Helm (coming soon)

```bash
helm repo add kubepurge https://muhyousri.github.io/kubepurge
helm install kubepurge kubepurge/kubepurge
```

## 🛠️ Usage

### Basic Example

Create a `PurgePolicy` to clean development resources daily:

```yaml
apiVersion: kubepurge.xyz/v1
kind: PurgePolicy
metadata:
  name: dev-cleanup
  namespace: kubepurge-system
spec:
  targetNamespace: "dev-environment"
  schedule: "0 18 * * *"  # Daily at 6 PM
  resources:
    - "pods"
    - "deployments" 
    - "services"
```

### Advanced Example

Clean multiple resource types with specific timing:

```yaml
apiVersion: kubepurge.xyz/v1
kind: PurgePolicy
metadata:
  name: staging-weekend-cleanup
  namespace: kubepurge-system
spec:
  targetNamespace: "staging"
  schedule: "0 0 * * 0"  # Every Sunday at midnight
  resources:
    - "pods"
    - "deployments"
    - "services"
    - "configmaps"
    - "secrets"
```

### Protecting Resources

Add the exclusion label to protect important resources:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: critical-service
  namespace: dev-environment
  labels:
    kubepurge.xyz/exclude: "true"  # This deployment will be skipped
spec:
  # ... deployment spec
```

## 📋 Custom Resources

### PurgePolicy

Defines what resources to purge and when:

| Field | Type | Description | Required |
|-------|------|-------------|----------|
| `targetNamespace` | string | Namespace to target for cleanup | ✅ |
| `schedule` | string | Cron format schedule | ✅ |
| `resources` | []string | Resource types to purge | ✅ |

**Supported Resource Types:**
- `pods`
- `deployments` 
- `services`
- `configmaps`
- `secrets`

### PurgeStatus

Tracks purge operation results (automatically created):

| Field | Type | Description |
|-------|------|-------------|
| `cleanedNamespace` | string | Namespace that was cleaned |
| `lastPurgeTime` | Time | Timestamp of last operation |
| `purgedResources` | map[string]string | Count of resources purged |

## 🔧 Configuration

### Cron Schedule Format

The `schedule` field uses standard cron format:

```
# ┌───────────── minute (0 - 59)
# │ ┌───────────── hour (0 - 23)
# │ │ ┌───────────── day of the month (1 - 31)
# │ │ │ ┌───────────── month (1 - 12)
# │ │ │ │ ┌───────────── day of the week (0 - 6) (Sunday to Saturday)
# │ │ │ │ │
# │ │ │ │ │
# * * * * *
```

**Examples:**
- `"0 18 * * *"` - Daily at 6 PM
- `"0 0 * * 0"` - Every Sunday at midnight  
- `"0 */6 * * *"` - Every 6 hours
- `"*/30 * * * *"` - Every 30 minutes

### Environment Variables

Configure the controller behavior:

| Variable | Default | Description |
|----------|---------|-------------|
| `METRICS_ADDR` | `:8080` | Metrics server address |
| `PROBE_ADDR` | `:8081` | Health probe address |
| `LEADER_ELECT` | `false` | Enable leader election |

## 🚦 RBAC Permissions

kubepurge requires the following permissions:

```yaml
# Core resources
- apiGroups: [""]
  resources: ["pods", "services", "configmaps", "secrets"]
  verbs: ["get", "list", "delete"]

# Apps resources  
- apiGroups: ["apps"]
  resources: ["deployments", "replicasets"]
  verbs: ["get", "list", "delete"]

# Custom resources
- apiGroups: ["kubepurge.xyz"]
  resources: ["purgepolicies", "purgestatuses"]
  verbs: ["get", "list", "watch", "create", "update", "patch", "delete"]
```

## 📊 Monitoring

### Metrics

kubepurge exposes Prometheus metrics at `:8080/metrics`:

- `kubepurge_purges_total` - Total number of purge operations
- `kubepurge_resources_deleted_total` - Total resources deleted by type
- `kubepurge_purge_duration_seconds` - Time taken for purge operations

### Health Checks

- **Liveness**: `:8081/healthz`
- **Readiness**: `:8081/readyz`

### Observing Operations

Check purge status:

```bash
# List all purge policies
kubectl get purgepolicies -A

# Check purge status
kubectl get purgestatuses -A

# View controller logs
kubectl logs -n kubepurge-system deployment/kubepurge-controller-manager
```

## 🛡️ Security

### Best Practices

1. **Namespace Isolation**: Deploy kubepurge in its own namespace
2. **RBAC**: Use minimal required permissions
3. **Resource Protection**: Label critical resources with `kubepurge.xyz/exclude: "true"`
4. **Testing**: Test purge policies in development before production use
5. **Monitoring**: Set up alerts for unexpected purge operations

### Container Security

- Runs as non-root user (UID 65532)
- Read-only root filesystem
- No privilege escalation
- Minimal distroless base image
- Regular security scanning

## 🧪 Development

### Prerequisites

- Go 1.22+
- Kubernetes 1.25+
- Docker or Podman
- kubectl
- kind (for testing)

### Building

```bash
# Build binary
make build

# Build container image
make docker-build

# Run tests
make test

# Run linter
make lint
```

### Testing

```bash
# Unit tests
make test

# E2E tests
make test-e2e

# Test with kind cluster
kind create cluster
make install run
```

## 🚀 Releasing

Create a new release:

```bash
./scripts/release.sh v1.0.0
```

This will:
1. Validate the version format
2. Run tests
3. Create and push a git tag
4. Trigger automated release pipeline

See [RELEASE.md](RELEASE.md) for detailed release process.

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### Development Workflow

1. Fork the repository
2. Create a feature branch
3. Make changes with tests
4. Run `make test lint`
5. Submit a pull request

## 📝 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## 🙋‍♂️ Support

- **Issues**: [GitHub Issues](https://github.com/muhyousri/kubepurge/issues)
- **Discussions**: [GitHub Discussions](https://github.com/muhyousri/kubepurge/discussions)
- **Documentation**: [Wiki](https://github.com/muhyousri/kubepurge/wiki)

## 🗺️ Roadmap

- [ ] Helm chart support
- [ ] Web UI for policy management  
- [ ] Advanced scheduling (maintenance windows)
- [ ] Resource usage-based purging
- [ ] Integration with cost management tools
- [ ] Support for custom resources
- [ ] Backup before purge option

## ⚠️ Disclaimer

kubepurge permanently deletes Kubernetes resources. Always test in non-production environments first and ensure proper backup procedures are in place. Use the exclusion label (`kubepurge.xyz/exclude: "true"`) to protect critical resources.

---

**Made with ❤️ for the Kubernetes community**