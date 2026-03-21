#!/bin/bash
# shellcheck disable=SC1091
# ============================================================================
# updateContent.sh - Runs INSIDE container for cache/content refresh
# ============================================================================
# This script runs inside the container to refresh cached or prebuilt content.
# Cloud services may run this periodically to keep prebuilds fresh.
# Use it for: Updating package caches, refreshing prebuilt dependencies.
# ============================================================================

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/../shared/utils.sh"

log_info "updateContent: Refreshing cached content..."

# Refresh package manager caches if needed
if command_exists npm; then
  log_info "Checking npm cache..."
  npm cache verify 2>/dev/null || true
fi

if command_exists pip; then
  log_info "Checking pip cache..."
  pip cache info 2>/dev/null || true
fi

log_success "updateContent: Content refresh complete"
