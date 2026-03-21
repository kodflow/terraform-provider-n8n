#!/bin/bash
# ============================================================================
# Project-specific postCreate hook - terraform-provider-n8n
# ============================================================================
# Runs after the base postCreate lifecycle hook.
# Install project-specific Go tools and configure git hooks.
# ============================================================================

set -e

echo "========================================="
echo "Installing terraform-provider-n8n tools"
echo "========================================="

# Locate Go binary
GO_BIN="/usr/local/go/bin/go"
if [ ! -x "$GO_BIN" ]; then
  GO_BIN="$(command -v go 2>/dev/null || echo "")"
fi

if [ ! -x "$GO_BIN" ]; then
  echo "⚠ Go not found, skipping Go tool installation"
  exit 0
fi

echo "→ Go found at: $GO_BIN"
echo ""

# Install goreleaser (for releasing the Terraform provider)
if [ ! -f "$HOME/.cache/go/bin/goreleaser" ]; then
  echo "→ Installing goreleaser v2.6.1..."
  "$GO_BIN" install github.com/goreleaser/goreleaser/v2@v2.6.1
  echo "  ✓ goreleaser installed"
else
  echo "  ✓ goreleaser already installed"
fi

# Install buildifier (Bazel BUILD file formatter)
if [ ! -f "$HOME/.cache/go/bin/buildifier" ]; then
  echo "→ Installing buildifier..."
  "$GO_BIN" install github.com/bazelbuild/buildtools/buildifier@latest
  echo "  ✓ buildifier installed"
else
  echo "  ✓ buildifier already installed"
fi

# Install shfmt (shell script formatter)
if [ ! -f "$HOME/.cache/go/bin/shfmt" ]; then
  echo "→ Installing shfmt..."
  "$GO_BIN" install mvdan.cc/sh/v3/cmd/shfmt@latest
  echo "  ✓ shfmt installed"
else
  echo "  ✓ shfmt already installed"
fi

# Configure git hooks
echo ""
echo "→ Configuring git hooks..."
if [ -f "/workspace/scripts/install-hooks.sh" ]; then
  chmod +x /workspace/scripts/install-hooks.sh
  /workspace/scripts/install-hooks.sh
  echo "  ✓ Git hooks configured"
else
  echo "  ⚠ scripts/install-hooks.sh not found, skipping"
fi

echo ""
echo "========================================="
echo "✓ Project tools installed successfully"
echo "========================================="
echo ""
echo "Available make commands:"
echo "  make help   - Show all commands"
echo "  make test   - Run test suite"
echo "  make lint   - Run linters"
echo "  make build  - Build provider"
echo ""
