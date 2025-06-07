# Release v0.0.1

## 🚀 First Release of kubepurge!

This is the initial release of kubepurge, a Kubernetes controller for automated resource cleanup and cost optimization in development environments.

### ✨ Key Features

- **📅 Scheduled Cleanup**: Cron-based automatic resource deletion
- **🎯 Targeted Purging**: Namespace and resource type selection
- **🛡️ Safe Guards**: Label-based exclusion for critical resources
- **📊 Status Tracking**: Monitor operations with PurgeStatus resources
- **🔒 Security First**: RBAC-enabled with minimal permissions

### 📦 What's Included

- **Custom Resources**: PurgePolicy and PurgeStatus CRDs
- **Controller**: Complete Kubernetes controller implementation
- **Container Images**: Multi-platform builds (linux/amd64, linux/arm64)
- **Installation**: Ready-to-use Kubernetes manifests
- **Documentation**: Comprehensive README and examples

### 🛠️ Installation

```bash
kubectl apply -f https://github.com/muhyousri/kubepurge/releases/download/v0.0.1/kubepurge-install.yaml
```

### 📋 Supported Resources

- Pods
- Deployments
- Services  
- ConfigMaps
- Secrets

### 🔧 Quick Start

Create a PurgePolicy to clean development resources daily:

```yaml
apiVersion: kubepurge.xyz/v1
kind: PurgePolicy
metadata:
  name: dev-cleanup
spec:
  targetNamespace: "dev-environment"
  schedule: "0 18 * * *"  # Daily at 6 PM
  resources:
    - "pods"
    - "deployments"
    - "services"
```

### 🛡️ Safety Features

- Non-root container execution
- Read-only root filesystem
- Label-based resource exclusion
- Comprehensive error handling
- Detailed logging and status reporting

### 📊 Container Images

- `ghcr.io/muhyousri/kubepurge:v0.0.1`
- `ghcr.io/muhyousri/kubepurge:latest`

### 🔍 What's Next

Check out our [roadmap](https://github.com/muhyousri/kubepurge#-roadmap) for upcoming features including:
- Helm chart support
- Web UI for policy management
- Advanced scheduling options
- Integration with cost management tools

### 🤝 Contributing

We welcome contributions! See our [contributing guide](https://github.com/muhyousri/kubepurge/blob/main/CONTRIBUTING.md) for details.

---

**⚠️ Important**: kubepurge permanently deletes Kubernetes resources. Always test in non-production environments first!