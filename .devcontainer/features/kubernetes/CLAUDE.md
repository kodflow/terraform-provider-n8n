<!-- updated: 2026-02-27T13:00:00Z -->

# Kubernetes Feature (kind)

## Purpose

Local Kubernetes development using kind (Kubernetes in Docker).

## Components

| Tool    | Version | Description                   |
| ------- | ------- | ----------------------------- |
| kind    | latest  | Kubernetes clusters in Docker |
| kubectl | latest  | Kubernetes CLI                |
| Helm    | latest  | Package manager for K8s       |

## Quick Start

```bash
# Create a cluster
kind create cluster --name dev

# Check cluster
kubectl cluster-info
kubectl get nodes

# Delete cluster
kind delete cluster --name dev
```

## Configuration

Enable in `devcontainer.json`:

```json
"features": {
  "./features/kubernetes": {
    "kindVersion": "latest",
    "enableHelm": true,
    "enableRegistry": true
  }
}
```

## Options

| Option           | Default | Description             |
| ---------------- | ------- | ----------------------- |
| `kindVersion`    | latest  | kind version            |
| `kubectlVersion` | latest  | kubectl version         |
| `helmVersion`    | latest  | Helm version            |
| `enableHelm`     | true    | Install Helm            |
| `enableRegistry` | true    | Local registry on :5001 |
| `clusterName`    | dev     | Default cluster name    |

## Advanced Usage

### Custom Cluster Config

```bash
# Use provided config
kind create cluster --config /workspace/.devcontainer/features/kubernetes/kind-config.yaml

# Multi-node cluster
kind create cluster --config - <<EOF
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
- role: worker
- role: worker
EOF
```

### Local Registry

```bash
# Push image to local registry
docker tag myapp:latest localhost:5001/myapp:latest
docker push localhost:5001/myapp:latest

# Use in K8s
kubectl create deployment myapp --image=localhost:5001/myapp:latest
```

### Helm Charts

```bash
# Add repo
helm repo add bitnami https://charts.bitnami.com/bitnami

# Install chart
helm install myrelease bitnami/nginx

# List releases
helm list
```

## CI Integration

For GitHub Actions, kind works out of the box:

```yaml
- name: Create kind cluster
  uses: helm/kind-action@v1
  with:
    cluster_name: test
```
