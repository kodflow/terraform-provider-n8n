#!/bin/bash
set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo "========================================="
echo "Installing Kubernetes Development Tools"
echo "========================================="

# ─────────────────────────────────────────────────────────────────────────────
# Pre-installed tool detection (base image may already include kubectl/Helm)
# ─────────────────────────────────────────────────────────────────────────────
KUBECTL_PREINSTALLED=false
HELM_PREINSTALLED=false

if command -v kubectl &>/dev/null; then
  echo -e "${GREEN}✓ kubectl already installed: $(kubectl version --client --short 2>/dev/null || kubectl version --client 2>/dev/null | head -1)${NC}"
  KUBECTL_PREINSTALLED=true
fi

if command -v helm &>/dev/null; then
  echo -e "${GREEN}✓ Helm already installed: $(helm version --short 2>/dev/null || echo 'unknown')${NC}"
  HELM_PREINSTALLED=true
fi

# ─────────────────────────────────────────────────────────────────────────────
# Input validation helpers
# ─────────────────────────────────────────────────────────────────────────────

# Validate version format: "latest" or semver (v1.2.3 or 1.2.3)
validate_version() {
  local name=$1
  local value=$2
  if [[ "$value" == "latest" ]]; then
    return 0
  fi
  # Allow v1.2.3, 1.2.3, v1.2.3-beta.1, etc.
  if [[ "$value" =~ ^v?[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9.]+)?$ ]]; then
    return 0
  fi
  echo -e "${RED}✗ Invalid $name version format: $value${NC}"
  echo -e "${RED}  Expected: 'latest' or semver (e.g., 1.2.3 or v1.2.3)${NC}"
  exit 1
}

# Validate cluster/registry name: alphanumeric, dash, underscore only
validate_name() {
  local name=$1
  local value=$2
  if [[ "$value" =~ ^[a-zA-Z0-9][a-zA-Z0-9_-]*$ ]] && [[ ${#value} -le 63 ]]; then
    return 0
  fi
  echo -e "${RED}✗ Invalid $name: $value${NC}"
  echo -e "${RED}  Must be alphanumeric (with - or _), max 63 chars${NC}"
  exit 1
}

# Validate boolean: true or false only
validate_boolean() {
  local name=$1
  local value=$2
  if [[ "$value" == "true" ]] || [[ "$value" == "false" ]]; then
    return 0
  fi
  echo -e "${RED}✗ Invalid $name: $value${NC}"
  echo -e "${RED}  Expected: 'true' or 'false'${NC}"
  exit 1
}

# ─────────────────────────────────────────────────────────────────────────────
# Feature options (from devcontainer-feature.json)
# ─────────────────────────────────────────────────────────────────────────────
KIND_VERSION="${KINDVERSION:-latest}"
KUBECTL_VERSION="${KUBECTLVERSION:-latest}"
HELM_VERSION="${HELMVERSION:-latest}"
ENABLE_HELM="${ENABLEHELM:-true}"
ENABLE_REGISTRY="${ENABLEREGISTRY:-true}"
CLUSTER_NAME="${CLUSTERNAME:-dev}"

# Validate all inputs
validate_version "kindVersion" "$KIND_VERSION"
validate_version "kubectlVersion" "$KUBECTL_VERSION"
validate_version "helmVersion" "$HELM_VERSION"
validate_boolean "enableHelm" "$ENABLE_HELM"
validate_boolean "enableRegistry" "$ENABLE_REGISTRY"
validate_name "clusterName" "$CLUSTER_NAME"

# Target user for kubeconfig (devcontainer provides _REMOTE_USER)
TARGET_USER="${_REMOTE_USER:-vscode}"
TARGET_HOME="$(getent passwd "$TARGET_USER" 2>/dev/null | cut -d: -f6 || echo "/home/${TARGET_USER}")"

# ─────────────────────────────────────────────────────────────────────────────
# Detect architecture
# ─────────────────────────────────────────────────────────────────────────────
ARCH=$(uname -m)
case "$ARCH" in
  x86_64)
    ARCH_KIND="amd64"
    ARCH_KUBECTL="amd64"
    ARCH_HELM="amd64"
    ;;
  aarch64 | arm64)
    ARCH_KIND="arm64"
    ARCH_KUBECTL="arm64"
    ARCH_HELM="arm64"
    ;;
  *)
    echo -e "${RED}Unsupported architecture: $ARCH${NC}"
    exit 1
    ;;
esac

# ─────────────────────────────────────────────────────────────────────────────
# Helper functions
# ─────────────────────────────────────────────────────────────────────────────

# Validate port number (1-65535)
validate_port() {
  local name=$1
  local value=$2
  if [[ "$value" =~ ^[0-9]+$ ]] && ((value >= 1 && value <= 65535)); then
    return 0
  fi
  echo "Error: Invalid $name '$value'. Must be an integer 1-65535." >&2
  return 1
}

# Get latest version from GitHub
get_github_version() {
  local repo=$1
  local version response

  response="$(curl -fsSL --connect-timeout 5 --max-time 10 \
    "https://api.github.com/repos/${repo}/releases/latest" 2>/dev/null)" || response=""

  if command -v jq &>/dev/null; then
    version="$(echo "$response" | jq -r '.tag_name // empty' 2>/dev/null)" || version=""
  else
    version="$(echo "$response" | sed -n 's/.*"tag_name": *"\([^"]*\)".*/\1/p' | head -n 1)" || version=""
  fi

  if [[ -z "$version" ]]; then
    echo -e "${RED}✗ Failed to resolve latest version for ${repo}${NC}" >&2
    echo -e "${RED}  GitHub API may be rate-limited. Try setting an explicit version.${NC}" >&2
    exit 1
  fi

  echo "$version"
}

# Verify SHA256 checksum
verify_sha256() {
  local file=$1
  local checksum_url=$2
  local filename expected actual checksums

  # Check sha256sum availability
  if ! command -v sha256sum &>/dev/null; then
    echo -e "${YELLOW}⚠ sha256sum not available, skipping verification${NC}"
    return 0
  fi

  filename="$(basename "$file")"
  checksums="$(curl -fsSL --connect-timeout 5 --max-time 10 "$checksum_url" 2>/dev/null)" || true

  if [[ -z "$checksums" ]]; then
    echo -e "${YELLOW}⚠ Checksum not available, skipping verification${NC}"
    return 0
  fi

  # Try to find checksum for specific filename, fallback to first entry
  expected=$(echo "$checksums" | awk -v f="$filename" '$2==f || $2=="*"f {print $1; exit}')
  if [[ -z "$expected" ]]; then
    # Single-line checksum file (just the hash)
    expected=$(echo "$checksums" | awk 'NF>=1 {print $1; exit}')
  fi

  if [[ -z "$expected" ]]; then
    echo -e "${YELLOW}⚠ Could not parse checksum, skipping verification${NC}"
    return 0
  fi

  actual=$(sha256sum "$file" | awk '{print $1}')
  if [[ "$expected" == "$actual" ]]; then
    echo -e "${GREEN}  ✓ SHA256 verified${NC}"
    return 0
  fi

  echo -e "${RED}  ✗ SHA256 mismatch!${NC}"
  echo -e "${RED}    Expected: $expected${NC}"
  echo -e "${RED}    Actual:   $actual${NC}"
  return 1
}

# ─────────────────────────────────────────────────────────────────────────────
# Install kind
# ─────────────────────────────────────────────────────────────────────────────
echo -e "${YELLOW}Installing kind...${NC}"

if [[ "$KIND_VERSION" == "latest" ]]; then
  KIND_VERSION=$(get_github_version "kubernetes-sigs/kind")
fi
# Ensure version starts with 'v'
[[ "$KIND_VERSION" != v* ]] && KIND_VERSION="v${KIND_VERSION}"

KIND_URL="https://kind.sigs.k8s.io/dl/${KIND_VERSION}/kind-linux-${ARCH_KIND}"
KIND_SHA_URL="${KIND_URL}.sha256sum"
if curl -fsSL --connect-timeout 10 --max-time 120 -o /tmp/kind "$KIND_URL"; then
  if verify_sha256 /tmp/kind "$KIND_SHA_URL"; then
    mv /tmp/kind /usr/local/bin/kind
    chmod +x /usr/local/bin/kind
    echo -e "${GREEN}✓ kind ${KIND_VERSION} installed${NC}"
  else
    rm -f /tmp/kind
    exit 1
  fi
else
  echo -e "${RED}✗ kind installation failed${NC}"
  exit 1
fi

# ─────────────────────────────────────────────────────────────────────────────
# Install kubectl
# ─────────────────────────────────────────────────────────────────────────────
echo -e "${YELLOW}Installing kubectl...${NC}"

if [[ "$KUBECTL_VERSION" == "latest" ]]; then
  KUBECTL_VERSION="$(curl -fsSL --connect-timeout 5 --max-time 10 \
    https://dl.k8s.io/release/stable.txt 2>/dev/null)" || true
  if [[ -z "$KUBECTL_VERSION" ]]; then
    echo -e "${RED}✗ Failed to resolve latest kubectl version${NC}"
    echo -e "${RED}  Try setting an explicit kubectlVersion.${NC}"
    exit 1
  fi
fi
[[ "$KUBECTL_VERSION" != v* ]] && KUBECTL_VERSION="v${KUBECTL_VERSION}"

KUBECTL_URL="https://dl.k8s.io/release/${KUBECTL_VERSION}/bin/linux/${ARCH_KUBECTL}/kubectl"
KUBECTL_SHA_URL="${KUBECTL_URL}.sha256"
if curl -fsSL --connect-timeout 10 --max-time 120 -o /tmp/kubectl "$KUBECTL_URL"; then
  if verify_sha256 /tmp/kubectl "$KUBECTL_SHA_URL"; then
    mv /tmp/kubectl /usr/local/bin/kubectl
    chmod +x /usr/local/bin/kubectl
    echo -e "${GREEN}✓ kubectl ${KUBECTL_VERSION} installed${NC}"
  else
    rm -f /tmp/kubectl
    exit 1
  fi
else
  echo -e "${RED}✗ kubectl installation failed${NC}"
  exit 1
fi

# ─────────────────────────────────────────────────────────────────────────────
# Install Helm (if enabled)
# ─────────────────────────────────────────────────────────────────────────────
if [[ "$ENABLE_HELM" == "true" ]]; then
  echo -e "${YELLOW}Installing Helm...${NC}"

  if [[ "$HELM_VERSION" == "latest" ]]; then
    HELM_VERSION=$(get_github_version "helm/helm")
  fi
  [[ "$HELM_VERSION" != v* ]] && HELM_VERSION="v${HELM_VERSION}"

  HELM_URL="https://get.helm.sh/helm-${HELM_VERSION}-linux-${ARCH_HELM}.tar.gz"
  HELM_SHA_URL="${HELM_URL}.sha256sum"
  if curl -fsSL --connect-timeout 10 --max-time 120 -o /tmp/helm.tar.gz "$HELM_URL"; then
    if verify_sha256 /tmp/helm.tar.gz "$HELM_SHA_URL"; then
      tar -xzf /tmp/helm.tar.gz -C /tmp
      mv "/tmp/linux-${ARCH_HELM}/helm" /usr/local/bin/helm
      chmod +x /usr/local/bin/helm
      rm -rf /tmp/helm.tar.gz "/tmp/linux-${ARCH_HELM}"
      echo -e "${GREEN}✓ Helm ${HELM_VERSION} installed${NC}"
    else
      rm -f /tmp/helm.tar.gz
      echo -e "${RED}✗ Helm checksum verification failed${NC}"
      exit 1
    fi
  else
    echo -e "${RED}✗ Helm installation failed${NC}"
    exit 1
  fi
fi

# ─────────────────────────────────────────────────────────────────────────────
# Setup kubeconfig directory
# ─────────────────────────────────────────────────────────────────────────────
echo -e "${YELLOW}Setting up kubeconfig...${NC}"

mkdir -p "${TARGET_HOME}/.kube"
touch "${TARGET_HOME}/.kube/config"
chown -R "${TARGET_USER}:${TARGET_USER}" "${TARGET_HOME}/.kube" 2>/dev/null || true
chmod 600 "${TARGET_HOME}/.kube/config" 2>/dev/null || true
echo -e "${GREEN}✓ kubeconfig directory created${NC}"

# ─────────────────────────────────────────────────────────────────────────────
# Create local registry setup script (if enabled)
# ─────────────────────────────────────────────────────────────────────────────
if [[ "$ENABLE_REGISTRY" == "true" ]]; then
  echo -e "${YELLOW}Creating local registry script...${NC}"

  cat >/usr/local/bin/kind-with-registry <<'REGISTRY_SCRIPT'
#!/bin/bash
set -euo pipefail

# Validate cluster name (alphanumeric, dash, underscore, max 63 chars)
validate_cluster_name() {
    local name=$1
    if [[ "$name" =~ ^[a-zA-Z0-9][a-zA-Z0-9_-]*$ ]] && [[ ${#name} -le 63 ]]; then
        return 0
    fi
    echo "Error: Invalid cluster name '$name'. Must be alphanumeric (with - or _), max 63 chars." >&2
    exit 1
}

# Validate port number (1-65535)
validate_port() {
    local name=$1
    local value=$2
    if [[ "$value" =~ ^[0-9]+$ ]] && (( value >= 1 && value <= 65535 )); then
        return 0
    fi
    echo "Error: Invalid $name '$value'. Must be an integer 1-65535." >&2
    exit 1
}

# Check Docker dependency
if ! command -v docker &>/dev/null; then
    echo "Error: docker is required for kind-with-registry (docker-outside-of-docker)." >&2
    exit 1
fi

CLUSTER_NAME="${1:-dev}"
validate_cluster_name "$CLUSTER_NAME"

REGISTRY_NAME="kind-registry"
REGISTRY_PORT="${REGISTRY_PORT:-5001}"
validate_port "REGISTRY_PORT" "$REGISTRY_PORT"

# Create registry container if not exists
if [[ "$(docker inspect -f '{{.State.Running}}' "${REGISTRY_NAME}" 2>/dev/null || echo 'false')" != 'true' ]]; then
    docker run -d --restart=always -p "127.0.0.1:${REGISTRY_PORT}:5000" --network bridge --name "${REGISTRY_NAME}" registry:2
fi

# Create kind cluster with containerd registry config
cat <<EOF | kind create cluster --name "${CLUSTER_NAME}" --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
containerdConfigPatches:
- |-
  [plugins."io.containerd.grpc.v1.cri".registry]
    config_path = "/etc/containerd/certs.d"
nodes:
- role: control-plane
  extraPortMappings:
  - containerPort: 30000
    hostPort: 30000
    protocol: TCP
  - containerPort: 30001
    hostPort: 30001
    protocol: TCP
  - containerPort: 80
    hostPort: 8080
    protocol: TCP
  - containerPort: 443
    hostPort: 8443
    protocol: TCP
EOF

# Configure registry in cluster (avoid word-splitting with while read)
REGISTRY_DIR="/etc/containerd/certs.d/localhost:${REGISTRY_PORT}"
kind get nodes --name "${CLUSTER_NAME}" | while read -r node; do
    docker exec "${node}" mkdir -p "${REGISTRY_DIR}"
    cat <<EOF | docker exec -i "${node}" cp /dev/stdin "${REGISTRY_DIR}/hosts.toml"
[host."http://${REGISTRY_NAME}:5000"]
EOF
done

# Connect registry to kind network
if [[ "$(docker inspect -f='{{json .NetworkSettings.Networks.kind}}' "${REGISTRY_NAME}" 2>/dev/null || echo 'null')" == 'null' ]]; then
    docker network connect "kind" "${REGISTRY_NAME}"
fi

# Document registry in cluster
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: ConfigMap
metadata:
  name: local-registry-hosting
  namespace: kube-public
data:
  localRegistryHosting.v1: |
    host: "localhost:${REGISTRY_PORT}"
    help: "https://kind.sigs.k8s.io/docs/user/local-registry/"
EOF

echo "Cluster '${CLUSTER_NAME}' created with local registry at localhost:${REGISTRY_PORT}"
REGISTRY_SCRIPT

  chmod +x /usr/local/bin/kind-with-registry
  echo -e "${GREEN}✓ kind-with-registry script created${NC}"
fi

# ─────────────────────────────────────────────────────────────────────────────
# Summary
# ─────────────────────────────────────────────────────────────────────────────
echo ""
echo -e "${GREEN}=========================================${NC}"
echo -e "${GREEN}Kubernetes tools installed successfully!${NC}"
echo -e "${GREEN}=========================================${NC}"
echo ""
echo "Installed components:"
echo "  - kind $(kind version 2>/dev/null || echo "$KIND_VERSION")"
echo "  - kubectl $(kubectl version --client -o yaml 2>/dev/null | grep gitVersion | awk '{print $2}' || echo "$KUBECTL_VERSION")"
if [[ "$ENABLE_HELM" == "true" ]]; then
  echo "  - helm $(helm version --short 2>/dev/null || echo "$HELM_VERSION")"
fi
echo ""
echo "Quick start:"
echo "  kind create cluster --name ${CLUSTER_NAME}"
if [[ "$ENABLE_REGISTRY" == "true" ]]; then
  echo "  # Or with local registry:"
  echo "  kind-with-registry ${CLUSTER_NAME}"
fi
echo ""
