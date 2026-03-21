#!/bin/bash
set -euo pipefail

echo "========================================="
echo "Installing Infrastructure Development Tools"
echo "========================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

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
  if [[ "$value" =~ ^v?[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9.]+)?$ ]]; then
    return 0
  fi
  echo -e "${RED}✗ Invalid $name version format: $value${NC}"
  echo -e "${RED}  Expected: 'latest' or semver (e.g., 1.2.3 or v1.2.3)${NC}"
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
TERRAGRUNT_VERSION="${TERRAGRUNTVERSION:-latest}"
TFLINT_VERSION="${TFLINTVERSION:-latest}"
INFRACOST_VERSION="${INFRACOSTVERSION:-latest}"
ENABLE_CFSSL="${ENABLECFSSL:-true}"
ENABLE_ANSIBLE_TOOLS="${ENABLEANSIBLETOOLS:-true}"

# Validate all inputs
validate_version "terragruntVersion" "$TERRAGRUNT_VERSION"
validate_version "tflintVersion" "$TFLINT_VERSION"
validate_version "infracostVersion" "$INFRACOST_VERSION"
validate_boolean "enableCfssl" "$ENABLE_CFSSL"
validate_boolean "enableAnsibleTools" "$ENABLE_ANSIBLE_TOOLS"

# ─────────────────────────────────────────────────────────────────────────────
# Detect architecture
# ─────────────────────────────────────────────────────────────────────────────
ARCH=$(uname -m)
case "$ARCH" in
  x86_64)
    ARCH_DL="amd64"
    ARCH_CFSSL="amd64"
    ;;
  aarch64 | arm64)
    ARCH_DL="arm64"
    ARCH_CFSSL="arm64"
    ;;
  *)
    echo -e "${RED}Unsupported architecture: $ARCH${NC}"
    exit 1
    ;;
esac

# ─────────────────────────────────────────────────────────────────────────────
# Helper functions
# ─────────────────────────────────────────────────────────────────────────────

# Get latest version from GitHub (no hardcoded fallbacks - fail loudly)
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

# Ensure unzip is available (needed for tflint)
ensure_unzip() {
  if command -v unzip &>/dev/null; then
    return 0
  fi
  echo -e "${YELLOW}Installing unzip (required for tflint)...${NC}"
  if command -v apt-get &>/dev/null; then
    apt-get update -qq && apt-get install -y -qq unzip >/dev/null 2>&1
  elif command -v apk &>/dev/null; then
    apk add --no-cache unzip >/dev/null 2>&1
  else
    echo -e "${RED}✗ Cannot install unzip - unknown package manager${NC}"
    return 1
  fi
}

# ─────────────────────────────────────────────────────────────────────────────
# Resolve versions
# ─────────────────────────────────────────────────────────────────────────────
echo -e "${YELLOW}Resolving tool versions...${NC}"

if [[ "$TERRAGRUNT_VERSION" == "latest" ]]; then
  TERRAGRUNT_VERSION=$(get_github_version "gruntwork-io/terragrunt")
fi
[[ "$TERRAGRUNT_VERSION" != v* ]] && TERRAGRUNT_VERSION="v${TERRAGRUNT_VERSION}"

if [[ "$TFLINT_VERSION" == "latest" ]]; then
  TFLINT_VERSION=$(get_github_version "terraform-linters/tflint")
fi
[[ "$TFLINT_VERSION" != v* ]] && TFLINT_VERSION="v${TFLINT_VERSION}"

if [[ "$INFRACOST_VERSION" == "latest" ]]; then
  INFRACOST_VERSION=$(get_github_version "infracost/infracost")
fi
[[ "$INFRACOST_VERSION" != v* ]] && INFRACOST_VERSION="v${INFRACOST_VERSION}"

CFSSL_VERSION=""
if [[ "$ENABLE_CFSSL" == "true" ]]; then
  CFSSL_VERSION=$(get_github_version "cloudflare/cfssl")
  [[ "$CFSSL_VERSION" != v* ]] && CFSSL_VERSION="v${CFSSL_VERSION}"
fi

# ─────────────────────────────────────────────────────────────────────────────
# Parallel binary downloads
# ─────────────────────────────────────────────────────────────────────────────
echo -e "${YELLOW}Downloading tools in parallel...${NC}"

DOWNLOAD_DIR=$(mktemp -d)
PIDS=()
FAIL=0

# Terragrunt - direct binary
TERRAGRUNT_URL="https://github.com/gruntwork-io/terragrunt/releases/download/${TERRAGRUNT_VERSION}/terragrunt_linux_${ARCH_DL}"
curl -fsSL --connect-timeout 10 --max-time 120 -o "${DOWNLOAD_DIR}/terragrunt" "$TERRAGRUNT_URL" &
PIDS+=($!)

# TFLint - zip archive
ensure_unzip
TFLINT_URL="https://github.com/terraform-linters/tflint/releases/download/${TFLINT_VERSION}/tflint_linux_${ARCH_DL}.zip"
curl -fsSL --connect-timeout 10 --max-time 120 -o "${DOWNLOAD_DIR}/tflint.zip" "$TFLINT_URL" &
PIDS+=($!)

# Infracost - tar.gz archive
INFRACOST_URL="https://github.com/infracost/infracost/releases/download/${INFRACOST_VERSION}/infracost-linux-${ARCH_DL}.tar.gz"
curl -fsSL --connect-timeout 10 --max-time 120 -o "${DOWNLOAD_DIR}/infracost.tar.gz" "$INFRACOST_URL" &
PIDS+=($!)

# cfssl - direct binaries (if enabled)
if [[ "$ENABLE_CFSSL" == "true" ]]; then
  CFSSL_VERSION_NUM="${CFSSL_VERSION#v}"
  CFSSL_URL="https://github.com/cloudflare/cfssl/releases/download/${CFSSL_VERSION}/cfssl_${CFSSL_VERSION_NUM}_linux_${ARCH_CFSSL}"
  CFSSLJSON_URL="https://github.com/cloudflare/cfssl/releases/download/${CFSSL_VERSION}/cfssljson_${CFSSL_VERSION_NUM}_linux_${ARCH_CFSSL}"
  curl -fsSL --connect-timeout 10 --max-time 120 -o "${DOWNLOAD_DIR}/cfssl" "$CFSSL_URL" &
  PIDS+=($!)
  curl -fsSL --connect-timeout 10 --max-time 120 -o "${DOWNLOAD_DIR}/cfssljson" "$CFSSLJSON_URL" &
  PIDS+=($!)
fi

# Wait for all downloads
for pid in "${PIDS[@]}"; do
  if ! wait "$pid"; then
    FAIL=1
  fi
done

if [[ "$FAIL" -ne 0 ]]; then
  echo -e "${RED}✗ One or more downloads failed${NC}"
  rm -rf "$DOWNLOAD_DIR"
  exit 1
fi

echo -e "${GREEN}✓ All downloads completed${NC}"

# ─────────────────────────────────────────────────────────────────────────────
# Install terragrunt
# ─────────────────────────────────────────────────────────────────────────────
echo -e "${YELLOW}Installing terragrunt...${NC}"
mv "${DOWNLOAD_DIR}/terragrunt" /usr/local/bin/terragrunt
chmod +x /usr/local/bin/terragrunt
echo -e "${GREEN}✓ terragrunt ${TERRAGRUNT_VERSION} installed${NC}"

# ─────────────────────────────────────────────────────────────────────────────
# Install tflint
# ─────────────────────────────────────────────────────────────────────────────
echo -e "${YELLOW}Installing tflint...${NC}"
unzip -o -q "${DOWNLOAD_DIR}/tflint.zip" -d "${DOWNLOAD_DIR}/tflint_extracted"
mv "${DOWNLOAD_DIR}/tflint_extracted/tflint" /usr/local/bin/tflint
chmod +x /usr/local/bin/tflint
echo -e "${GREEN}✓ tflint ${TFLINT_VERSION} installed${NC}"

# ─────────────────────────────────────────────────────────────────────────────
# Install infracost
# ─────────────────────────────────────────────────────────────────────────────
echo -e "${YELLOW}Installing infracost...${NC}"
tar -xzf "${DOWNLOAD_DIR}/infracost.tar.gz" -C "${DOWNLOAD_DIR}"
mv "${DOWNLOAD_DIR}/infracost-linux-${ARCH_DL}" /usr/local/bin/infracost
chmod +x /usr/local/bin/infracost
echo -e "${GREEN}✓ infracost ${INFRACOST_VERSION} installed${NC}"

# ─────────────────────────────────────────────────────────────────────────────
# Install cfssl (if enabled)
# ─────────────────────────────────────────────────────────────────────────────
if [[ "$ENABLE_CFSSL" == "true" ]]; then
  echo -e "${YELLOW}Installing cfssl...${NC}"
  mv "${DOWNLOAD_DIR}/cfssl" /usr/local/bin/cfssl
  mv "${DOWNLOAD_DIR}/cfssljson" /usr/local/bin/cfssljson
  chmod +x /usr/local/bin/cfssl /usr/local/bin/cfssljson
  echo -e "${GREEN}✓ cfssl ${CFSSL_VERSION} installed${NC}"
fi

# ─────────────────────────────────────────────────────────────────────────────
# Cleanup download directory
# ─────────────────────────────────────────────────────────────────────────────
rm -rf "$DOWNLOAD_DIR"

# ─────────────────────────────────────────────────────────────────────────────
# Install Python infra tools (if enabled)
# ─────────────────────────────────────────────────────────────────────────────
if [[ "$ENABLE_ANSIBLE_TOOLS" == "true" ]]; then
  echo -e "${YELLOW}Installing Python infra tools...${NC}"
  if command -v pip3 &>/dev/null; then
    pip3 install --no-cache-dir ansible-lint molecule molecule-docker \
      && echo -e "${GREEN}✓ ansible-lint, molecule installed${NC}" \
      || echo -e "${YELLOW}⚠ pip install failed - Python infra tools skipped${NC}"
  elif command -v pip &>/dev/null; then
    pip install --no-cache-dir ansible-lint molecule molecule-docker \
      && echo -e "${GREEN}✓ ansible-lint, molecule installed${NC}" \
      || echo -e "${YELLOW}⚠ pip install failed - Python infra tools skipped${NC}"
  else
    echo -e "${YELLOW}⚠ Python/pip not found - ansible-lint and molecule skipped${NC}"
    echo -e "${YELLOW}  Enable the Python feature to install these tools${NC}"
  fi
fi

# ─────────────────────────────────────────────────────────────────────────────
# Summary
# ─────────────────────────────────────────────────────────────────────────────
echo ""
echo -e "${GREEN}=========================================${NC}"
echo -e "${GREEN}Infrastructure tools installed successfully!${NC}"
echo -e "${GREEN}=========================================${NC}"
echo ""
echo "Installed components:"
echo "  - terragrunt $(terragrunt --version 2>/dev/null | head -1 || echo "$TERRAGRUNT_VERSION")"
echo "  - tflint $(tflint --version 2>/dev/null | head -1 || echo "$TFLINT_VERSION")"
echo "  - infracost $(infracost --version 2>/dev/null | head -1 || echo "$INFRACOST_VERSION")"
if [[ "$ENABLE_CFSSL" == "true" ]]; then
  echo "  - cfssl $(cfssl version 2>/dev/null | head -1 || echo "$CFSSL_VERSION")"
fi
if [[ "$ENABLE_ANSIBLE_TOOLS" == "true" ]]; then
  if command -v ansible-lint &>/dev/null; then
    echo "  - ansible-lint $(ansible-lint --version 2>/dev/null | head -1 || echo 'installed')"
  fi
  if command -v molecule &>/dev/null; then
    echo "  - molecule $(molecule --version 2>/dev/null | head -1 || echo 'installed')"
  fi
fi
echo ""
echo "Quick start:"
echo "  terragrunt run-all plan     # Plan all modules"
echo "  tflint --recursive          # Lint Terraform files"
echo "  infracost breakdown --path . # Estimate costs"
if [[ "$ENABLE_CFSSL" == "true" ]]; then
  echo "  cfssl gencert -initca ca-csr.json | cfssljson -bare ca  # Generate CA"
fi
echo ""
