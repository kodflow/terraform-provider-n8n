#!/bin/bash
# ============================================================================
# initialize.sh - Runs on HOST machine BEFORE container build
# ============================================================================
# This script runs on the host machine before Docker Compose starts.
# Use it for: .env setup, project name configuration, feature validation.
# ============================================================================

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
HOOKS_DIR="$(dirname "$SCRIPT_DIR")"
DEVCONTAINER_DIR="$(dirname "$HOOKS_DIR")"
ENV_FILE="$DEVCONTAINER_DIR/.env"

# Extract project name from git remote URL (with fallback if no remote)
REPO_NAME=$(basename "$(git config --get remote.origin.url 2>/dev/null || echo "devcontainer")" .git)

# Sanitize project name for Docker Compose requirements:
# - Must start with a letter or number
# - Only lowercase alphanumeric, hyphens, and underscores allowed
REPO_NAME=$(echo "$REPO_NAME" | sed 's/^[^a-zA-Z0-9]*//' | tr '[:upper:]' '[:lower:]' | sed 's/[^a-z0-9_-]/-/g')

# If name is empty after sanitization, use a default
if [ -z "$REPO_NAME" ]; then
  REPO_NAME="devcontainer"
fi

echo "Initializing devcontainer environment..."
echo "Project name: $REPO_NAME"

# If .env doesn't exist, create it from .env.example
if [ ! -f "$ENV_FILE" ]; then
  echo "Creating .env from .env.example..."
  cp "$HOOKS_DIR/shared/.env.example" "$ENV_FILE"
fi

# Update or add COMPOSE_PROJECT_NAME in .env
if grep -q "^COMPOSE_PROJECT_NAME=" "$ENV_FILE"; then
  # Update existing line
  if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS
    sed -i '' "s|^COMPOSE_PROJECT_NAME=.*|COMPOSE_PROJECT_NAME=$REPO_NAME|" "$ENV_FILE"
  else
    # Linux
    sed -i "s|^COMPOSE_PROJECT_NAME=.*|COMPOSE_PROJECT_NAME=$REPO_NAME|" "$ENV_FILE"
  fi
  echo "Updated COMPOSE_PROJECT_NAME=$REPO_NAME in .env"
else
  # Add at the beginning of the file
  echo "COMPOSE_PROJECT_NAME=$REPO_NAME" | cat - "$ENV_FILE" >"$ENV_FILE.tmp" && mv "$ENV_FILE.tmp" "$ENV_FILE"
  echo "Added COMPOSE_PROJECT_NAME=$REPO_NAME to .env"
fi

# ============================================================================
# Validate devcontainer features
# ============================================================================
echo ""
echo "Validating devcontainer features..."

FEATURES_DIR="$DEVCONTAINER_DIR/features"
ERRORS=0
FIXED=0

for category in "$FEATURES_DIR"/*; do
  [ ! -d "$category" ] && continue

  for feature in "$category"/*; do
    [ ! -d "$feature" ] && continue

    # Skip utility directories (not actual features)
    [ "$(basename "$feature")" = "shared" ] && continue

    feature_name="$(basename "$category")/$(basename "$feature")"

    # Check devcontainer-feature.json
    if [ ! -f "$feature/devcontainer-feature.json" ]; then
      echo "ERROR: $feature_name: Missing devcontainer-feature.json"
      ERRORS=$((ERRORS + 1))
      continue
    fi

    # Check install.sh
    if [ ! -f "$feature/install.sh" ]; then
      echo "ERROR: $feature_name: Missing install.sh"
      ERRORS=$((ERRORS + 1))
      continue
    fi

    # Fix permissions if needed
    if [ ! -x "$feature/install.sh" ]; then
      chmod +x "$feature/install.sh"
      FIXED=$((FIXED + 1))
    fi
  done
done

if [ $ERRORS -gt 0 ]; then
  echo ""
  echo "ERROR: Found $ERRORS critical error(s) in features!"
  echo "Please fix missing files before building the devcontainer."
  exit 1
fi

if [ $FIXED -gt 0 ]; then
  echo "Fixed permissions on $FIXED install.sh file(s)"
fi

echo "All features validated successfully"

# ============================================================================
# Ollama Installation (Host GPU Acceleration for grepai)
# ============================================================================
# Ollama runs on the HOST to leverage GPU (Metal on Mac, CUDA on Linux)
# The DevContainer connects via host.docker.internal:11434
# ============================================================================
# Extract Ollama model from grepai config (single source of truth)
GREPAI_CONFIG="$DEVCONTAINER_DIR/images/grepai.config.yaml"
OLLAMA_MODEL=$(grep -E '^\s+model:' "$GREPAI_CONFIG" 2>/dev/null | awk '{print $2}' | head -1)
OLLAMA_MODEL="${OLLAMA_MODEL:-bge-m3}"

echo ""
echo "Setting up Ollama for GPU-accelerated semantic search..."

# Detect OS
detect_os() {
  case "$OSTYPE" in
    darwin*) echo "macos" ;;
    linux*) echo "linux" ;;
    msys* | cygwin* | mingw*) echo "windows" ;;
    *) echo "unknown" ;;
  esac
}

# Check if Ollama is installed
check_ollama_installed() {
  command -v ollama &>/dev/null
}

# Check if Ollama is running
check_ollama_running() {
  curl -sf --connect-timeout 2 http://localhost:11434/api/tags &>/dev/null
}

# Install Ollama based on OS
install_ollama() {
  local os="$1"
  echo "Installing Ollama..."

  case "$os" in
    macos)
      if command -v brew &>/dev/null; then
        brew install ollama
      else
        curl -fsSL https://ollama.com/install.sh | sh
      fi
      ;;
    linux)
      curl -fsSL https://ollama.com/install.sh | sh
      ;;
    windows)
      echo "Windows detected. Please install Ollama manually:"
      echo "  Download from: https://ollama.com/download/windows"
      echo "  Or via winget: winget install Ollama.Ollama"
      return 1
      ;;
    *)
      echo "Unknown OS. Please install Ollama manually from https://ollama.com"
      return 1
      ;;
  esac
}

# Start Ollama daemon
start_ollama() {
  local os="$1"
  echo "Starting Ollama daemon..."

  case "$os" in
    macos)
      # On macOS, ollama serve runs as a background service
      # Check if launchd service exists
      if launchctl list 2>/dev/null | grep -q "com.ollama"; then
        launchctl start com.ollama.ollama 2>/dev/null || true
      else
        # Start manually in background
        nohup ollama serve >/dev/null 2>&1 &
      fi
      ;;
    linux)
      # Check if systemd service exists
      if systemctl list-unit-files 2>/dev/null | grep -q "ollama"; then
        sudo systemctl start ollama 2>/dev/null || nohup ollama serve >/dev/null 2>&1 &
      else
        nohup ollama serve >/dev/null 2>&1 &
      fi
      ;;
    windows)
      # On Windows, Ollama typically runs as a service after installation
      echo "Please ensure Ollama is running (check system tray)"
      ;;
  esac

  # Wait for Ollama to be ready (max 30 seconds)
  local retries=15
  while [ $retries -gt 0 ]; do
    if check_ollama_running; then
      echo "Ollama is ready"
      return 0
    fi
    retries=$((retries - 1))
    sleep 2
  done

  echo "Warning: Ollama did not start in time"
  return 1
}

# Pull embedding model if not present
pull_model() {
  local model="$1"
  echo "Checking for embedding model: $model..."

  if curl -sf http://localhost:11434/api/tags 2>/dev/null | grep -qw "$model"; then
    echo "Model $model already available"
  else
    echo "Pulling model $model (this may take a few minutes)..."
    ollama pull "$model"
  fi
}

# Main Ollama setup flow
OS=$(detect_os)
echo "Detected OS: $OS"

if check_ollama_installed; then
  echo "Ollama is installed"
else
  echo "Ollama not found, installing..."
  if ! install_ollama "$OS"; then
    echo "Warning: Could not install Ollama automatically"
    echo "Semantic search will use CPU-only sidecar (slower)"
  fi
fi

if check_ollama_installed; then
  if check_ollama_running; then
    echo "Ollama is running"
  else
    start_ollama "$OS"
  fi

  # Pull model if Ollama is running
  if check_ollama_running; then
    pull_model "$OLLAMA_MODEL"
    echo "Ollama setup complete - GPU acceleration enabled"
  fi
else
  echo "Warning: Ollama not available - will use CPU-only sidecar"
  echo "To enable GPU acceleration, install Ollama manually:"
  echo "  macOS: brew install ollama"
  echo "  Linux: curl -fsSL https://ollama.com/install.sh | sh"
  echo "  Windows: https://ollama.com/download/windows"
fi

# ============================================================================
# Pull latest Docker image (bypass Docker cache on rebuild)
# ============================================================================
echo ""
echo "Pulling latest devcontainer image..."
docker pull ghcr.io/kodflow/devcontainer-template:latest 2>/dev/null || echo "Warning: Could not pull latest image, using cached version"

# ============================================================================
# Clean up existing containers to prevent race conditions during rebuild
# ============================================================================
echo ""
echo "Cleaning up existing devcontainer instances..."
docker compose -f "$DEVCONTAINER_DIR/docker-compose.yml" --project-name "$REPO_NAME" down --remove-orphans --timeout 5 2>/dev/null || true
echo "Cleanup complete"

echo ""
echo "Environment initialization complete!"
