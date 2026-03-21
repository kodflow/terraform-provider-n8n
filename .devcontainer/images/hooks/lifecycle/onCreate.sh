#!/bin/bash
# shellcheck disable=SC1091
# ============================================================================
# onCreate.sh - Runs INSIDE container immediately after first start
# ============================================================================
# This script runs inside the container after it starts for the first time.
# Use it for: Initial container setup that doesn't need user-specific config.
# ============================================================================

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/../shared/utils.sh"

log_info "onCreate: Container created, performing initial setup..."

# Ensure cache directories exist with proper permissions
CACHE_DIRS=(
  "/home/vscode/.cache"
  "/home/vscode/.config"
  "/home/vscode/.local/bin"
  "/home/vscode/.local/share"
)

for dir in "${CACHE_DIRS[@]}"; do
  mkdir_safe "$dir"
done

# Note: .claude/ is now baked into the Docker image at /home/vscode/.claude/
# No longer needs injection from devcontainer feature

# Inject CLAUDE.md from devcontainer if not present in project
CLAUDE_FEATURE_DIR="/workspace/.devcontainer/features/claude"

if [ -f "$CLAUDE_FEATURE_DIR/CLAUDE.md" ]; then
  if [ ! -f "/workspace/CLAUDE.md" ]; then
    log_info "Injecting CLAUDE.md from devcontainer..."
    cp "$CLAUDE_FEATURE_DIR/CLAUDE.md" /workspace/CLAUDE.md
    log_success "CLAUDE.md injected to /workspace/"
  else
    log_info "Project has its own CLAUDE.md, skipping injection"
  fi
fi

log_success "onCreate: Initial container setup complete"
