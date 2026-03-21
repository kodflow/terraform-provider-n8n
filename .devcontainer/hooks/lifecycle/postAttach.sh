#!/bin/bash
# ============================================================================
# Lifecycle delegation stub (DO NOT add logic here)
# ============================================================================
# Real logic: /etc/devcontainer-hooks/lifecycle/${HOOK}.sh (image-embedded)
# This stub auto-updates behavior when the Docker image is rebuilt.
# ============================================================================

HOOK="$(basename "${BASH_SOURCE[0]}" .sh)"

# Priority 1: Template dev source (only exists in template repo)
DEV="/workspace/.devcontainer/images/hooks/lifecycle/${HOOK}.sh"
# Priority 2: Image-embedded (exists in all downstream containers)
IMG="/etc/devcontainer-hooks/lifecycle/${HOOK}.sh"

if [ -x "$DEV" ]; then
  "$DEV" "$@"
elif [ -x "$IMG" ]; then
  "$IMG" "$@"
else
  echo "[WARNING] No ${HOOK} hook implementation found"
fi

# Project-specific extension (optional)
EXT="/workspace/.devcontainer/hooks/project/${HOOK}.sh"
if [ -x "$EXT" ]; then "$EXT" "$@"; fi
