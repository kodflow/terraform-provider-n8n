#!/bin/bash
# Use set -e for safety; retry logic handles expected failures internally
set -e

FEATURE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=../shared/feature-utils.sh
source "${FEATURE_DIR}/../shared/feature-utils.sh" 2>/dev/null || {
  RED='\033[0;31m'
  GREEN='\033[0;32m'
  YELLOW='\033[1;33m'
  BLUE='\033[0;34m'
  NC='\033[0m'
  log_info() { echo -e "${BLUE}[INFO]${NC} $*"; }
  log_success() { echo -e "${GREEN}[SUCCESS]${NC} $*"; }
  log_warning() { echo -e "${YELLOW}[WARNING]${NC} $*"; }
  log_error() { echo -e "${RED}[ERROR]${NC} $*" >&2; }
  ok() { echo -e "${GREEN}✓${NC} $*"; }
  warn() { echo -e "${YELLOW}⚠${NC} $*"; }
}

# Retry function
retry() {
  local max_attempts=$1
  local delay=$2
  shift 2
  local attempt=1
  local exit_code=0

  while [ $attempt -le $max_attempts ]; do
    if "$@"; then
      if [ $attempt -gt 1 ]; then
        log_success "Command succeeded on attempt $attempt"
      fi
      return 0
    fi

    exit_code=$?

    if [ $attempt -lt $max_attempts ]; then
      log_warning "Command failed (exit code: $exit_code), retrying in ${delay}s... (attempt $attempt/$max_attempts)"
      sleep "$delay"
    else
      log_error "Command failed after $max_attempts attempts"
    fi

    ((attempt++))
  done

  return $exit_code
}

# apt-get with retry and lock handling
apt_get_retry() {
  local max_attempts=5
  local attempt=1
  local delay=10

  while [ $attempt -le $max_attempts ]; do
    # Wait for apt locks to be released
    local lock_wait=0
    while sudo fuser /var/lib/dpkg/lock-frontend >/dev/null 2>&1 \
      || sudo fuser /var/lib/apt/lists/lock >/dev/null 2>&1 \
      || sudo fuser /var/cache/apt/archives/lock >/dev/null 2>&1; do
      if [ $lock_wait -eq 0 ]; then
        log_warning "Waiting for apt locks to be released..."
      fi
      sleep 2
      lock_wait=$((lock_wait + 2))

      if [ $lock_wait -ge 60 ]; then
        log_warning "Forcing apt lock release after 60s wait"
        sudo rm -f /var/lib/dpkg/lock-frontend
        sudo rm -f /var/lib/apt/lists/lock
        sudo rm -f /var/cache/apt/archives/lock
        sudo dpkg --configure -a || true
        break
      fi
    done

    # Try apt-get command
    if sudo apt-get "$@"; then
      if [ $attempt -gt 1 ]; then
        log_success "apt-get succeeded on attempt $attempt"
      fi
      return 0
    fi

    exit_code=$?

    if [ $attempt -lt $max_attempts ]; then
      log_warning "apt-get failed, running update and retrying in ${delay}s... (attempt $attempt/$max_attempts)"
      sudo apt-get update --fix-missing || true
      sudo dpkg --configure -a || true
      sleep "$delay"
    else
      log_error "apt-get failed after $max_attempts attempts"
    fi

    ((attempt++))
  done

  return $exit_code
}

# Download and pipe to shell with retry
download_and_pipe() {
  local url=$1
  shift
  local shell_cmd=("$@")

  log_info "Downloading and executing: $url"

  local temp_file
  temp_file=$(mktemp)

  if retry 3 5 curl -fsSL --connect-timeout 30 --max-time 300 -o "$temp_file" "$url"; then
    "${shell_cmd[@]}" <"$temp_file"
    local exit_code=$?
    rm -f "$temp_file"
    return $exit_code
  else
    rm -f "$temp_file"
    return 1
  fi
}

# Safe directory creation
mkdir_safe() {
  local dir=$1
  if [ ! -d "$dir" ]; then
    mkdir -p "$dir" 2>/dev/null || sudo mkdir -p "$dir"

    if [ "$(whoami)" = "vscode" ]; then
      sudo chown -R vscode:vscode "$dir" 2>/dev/null || true
    fi
  fi
}

print_banner "Node.js Development Environment" 2>/dev/null || {
  echo "========================================="
  echo "Installing Node.js Development Environment"
  echo "========================================="
}

# Environment variables
# NVM installed in system location (not volume) - Microsoft best practice
# See: https://github.com/microsoft/vscode-dev-containers/blob/main/script-library/docs/node.md
export NVM_DIR="/usr/local/share/nvm"
export NVM_SYMLINK_CURRENT=true
export NODE_VERSION="${NODE_VERSION:-lts/*}"
export npm_config_cache="${npm_config_cache:-/home/vscode/.cache/npm}"

# Install dependencies with retry
log_info "Installing dependencies..."
apt_get_retry update
apt_get_retry install -y curl git build-essential libssl-dev || {
  log_error "Failed to install dependencies"
  exit 1
}

# Install NVM (Node Version Manager)
log_info "Installing NVM..."
mkdir_safe "$NVM_DIR"
# Fetch latest NVM version from GitHub API
NVM_LATEST=$(curl -fsSL "https://api.github.com/repos/nvm-sh/nvm/releases/latest" 2>/dev/null | grep -o '"tag_name": *"[^"]*"' | head -1 | cut -d'"' -f4) || true
if [ -z "$NVM_LATEST" ]; then
  log_error "Failed to resolve latest NVM version"
  exit 1
fi
log_info "Using NVM ${NVM_LATEST}"
download_and_pipe "https://raw.githubusercontent.com/nvm-sh/nvm/${NVM_LATEST}/install.sh" bash || {
  log_error "Failed to install NVM"
  exit 1
}

# Load NVM
export NVM_DIR="$NVM_DIR"
[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"

# Install Node.js (latest LTS by default)
log_info "Installing Node.js ${NODE_VERSION}..."
retry 3 5 nvm install "$NODE_VERSION" || {
  log_error "Failed to install Node.js"
  exit 1
}
nvm use "$NODE_VERSION"
nvm alias default "$NODE_VERSION"

# Get installed Node and npm versions
NODE_INSTALLED=$(node --version)
NPM_INSTALLED=$(npm --version)

log_success "Node.js ${NODE_INSTALLED} installed"
log_success "npm ${NPM_INSTALLED} installed"

# Create cache directory
mkdir_safe "$npm_config_cache"

# ─────────────────────────────────────────────────────────────────────────────
# Install Node.js Development Tools — batched for speed
# ─────────────────────────────────────────────────────────────────────────────
log_info "Installing Node.js development tools..."

npm install -g pnpm@latest typescript@latest eslint@latest prettier@latest tsx@latest
log_success "Core tools installed (pnpm, typescript, eslint, prettier, tsx)"

# Install additional global packages from devcontainer-feature.json option
GLOBAL_PACKAGES="${GLOBALPACKAGES:-}"
if [ -n "$GLOBAL_PACKAGES" ]; then
  log_info "Installing additional global packages: ${GLOBAL_PACKAGES}"
  IFS=',' read -ra PKGS <<<"$GLOBAL_PACKAGES"
  for pkg in "${PKGS[@]}"; do
    pkg=$(echo "$pkg" | xargs) # trim whitespace
    if [ -n "$pkg" ] && ! command -v "$pkg" &>/dev/null; then
      npm install -g "$pkg" 2>/dev/null && log_success "Installed $pkg" || log_warning "Failed to install $pkg"
    fi
  done
fi

# ─────────────────────────────────────────────────────────────────────────────
# Install Desktop & WASM Tools — parallel
# ─────────────────────────────────────────────────────────────────────────────
log_info "Installing Desktop & WASM tools..."

(
  if npm install -g electron@latest && command -v electron &>/dev/null; then
    ELECTRON_VERSION=$(electron --version 2>/dev/null || echo "installed")
    log_success "Electron ${ELECTRON_VERSION} installed"
  else
    log_warning "Electron installation failed or not in PATH"
  fi
) &
ELECTRON_PID=$!

(
  if npm install -g assemblyscript@latest && command -v asc &>/dev/null; then
    ASC_VERSION=$(asc --version 2>/dev/null | head -n 1 || echo "installed")
    log_success "AssemblyScript ${ASC_VERSION} installed"
  else
    log_warning "AssemblyScript installation failed or not in PATH"
  fi
) &
ASC_PID=$!

wait "$ELECTRON_PID" 2>/dev/null || true
wait "$ASC_PID" 2>/dev/null || true
log_success "Desktop & WASM tools installed"

# Create global symlinks for node, npm, and npx
# This ensures they're available for subsequent devcontainer features
log_info "Creating global symlinks..."
NVM_NODE_DIR=$(nvm which current | xargs dirname)
if [ -n "$NVM_NODE_DIR" ] && [ -d "$NVM_NODE_DIR" ]; then
  # Create symlinks in /usr/local/bin (which is in default PATH)
  sudo ln -sf "$NVM_NODE_DIR/node" /usr/local/bin/node
  sudo ln -sf "$NVM_NODE_DIR/npm" /usr/local/bin/npm
  sudo ln -sf "$NVM_NODE_DIR/npx" /usr/local/bin/npx

  # Verify symlinks were created
  if [ -L /usr/local/bin/node ] && [ -L /usr/local/bin/npm ]; then
    log_success "Global symlinks created in /usr/local/bin"
  else
    log_warning "Failed to create global symlinks, but NVM is still available"
  fi
else
  log_warning "Could not determine NVM node directory, skipping symlink creation"
fi

# Ensure vscode user can update NVM files (especially the 'current' symlink)
# This is required because NVM_SYMLINK_CURRENT=true needs write access
log_info "Setting NVM directory ownership for vscode user..."
if [ -d "$NVM_DIR" ]; then
  sudo chown -R vscode:vscode "$NVM_DIR" 2>/dev/null || true
  log_success "NVM directory ownership set to vscode"
fi

# NVM is loaded via ~/.devcontainer-env.sh (sourced by both .zshrc and .bashrc)
# No need to inject NVM into .zshrc directly — avoids double-loading which
# causes VS Code ptyHost shell environment resolution timeout
log_info "NVM will be loaded via ~/.devcontainer-env.sh (two-phase loading)"

print_success_banner "Node.js environment" 2>/dev/null || {
  echo ""
  echo -e "${GREEN}=========================================${NC}"
  echo -e "${GREEN}Node.js environment installed successfully!${NC}"
  echo -e "${GREEN}=========================================${NC}"
  echo ""
}
log_success "Installation complete!"
echo ""
echo "Installed components:"
echo "  - NVM (Node Version Manager)"
echo "  - Node.js ${NODE_INSTALLED}"
echo "  - npm ${NPM_INSTALLED}"
echo ""
echo "Development tools:"
echo "  - pnpm (package manager)"
echo "  - TypeScript (type checker)"
echo "  - ESLint (linter)"
echo "  - Prettier (formatter)"
echo "  - tsx (TypeScript runner)"
echo ""
echo "Desktop & WASM tools:"
echo "  - electron (desktop GUI framework)"
echo "  - assemblyscript (TypeScript to WASM compiler)"
echo ""
echo "Global availability:"
echo "  - node, npm, npx, pnpm, tsc, eslint, prettier available globally"
echo "  - NVM loaded in interactive shells"
echo ""
echo "Cache directory:"
echo "  - npm: $npm_config_cache"
echo ""

# Exit successfully
exit 0
