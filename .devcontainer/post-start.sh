#!/bin/bash
set -e

echo "🚀 Starting post-start configuration..."

# Ensure npm global directory has correct permissions
if [ -d "$HOME/.local/share/npm-global" ]; then
  echo "🔧 Checking npm global directory permissions..."
  chmod -R 755 "$HOME/.local/share/npm-global" 2>/dev/null || true
fi

# Check if Claude CLI is installed, if not install it
if ! command -v claude &>/dev/null; then
  echo "📦 Installing Claude CLI..."
  export NPM_CONFIG_PREFIX="$HOME/.local/share/npm-global"
  if npm install -g --prefix "$HOME/.local/share/npm-global" @anthropic-ai/claude-code@latest; then
    echo "✅ Claude CLI installed successfully"
  else
    echo "⚠️  Failed to install Claude CLI"
  fi
fi

# Check if commitlint is installed, if not install it
if ! command -v commitlint &>/dev/null; then
  echo "📦 Installing commitlint..."
  export NPM_CONFIG_PREFIX="$HOME/.local/share/npm-global"
  if npm install -g --prefix "$HOME/.local/share/npm-global" @commitlint/cli@latest @commitlint/config-conventional@latest; then
    echo "✅ commitlint installed successfully"
  else
    echo "⚠️  Failed to install commitlint"
  fi
fi

# Clean up old or duplicate binaries from previous builds
echo "🧹 Cleaning up old binaries..."
rm -f "$HOME/.cache/go/bin/golangci-lint-real" # Old golangci-lint wrapper
rm -f "$HOME/.cache/go/bin/ktn-linter"         # Duplicate (should be in .local/bin)

# Install Go tools (using absolute path to ensure Go is found)
GO_BIN="/usr/local/go/bin/go"
if [ -x "$GO_BIN" ]; then
  # Check if tools are already installed to avoid reinstalling every time
  if [ ! -f "$HOME/.cache/go/bin/golangci-lint" ] || ! "$HOME/.cache/go/bin/golangci-lint" version &>/dev/null; then
    echo "🔨 Installing Go tools..."

    # Clean up any corrupted golangci-lint binary
    rm -f "$HOME/.cache/go/bin/golangci-lint"

    # Install golangci-lint v2.6.1 from GitHub releases
    echo "📦 Installing golangci-lint v2.6.1..."
    ARCH="$(uname -m)"
    case "$ARCH" in
      x86_64) ARCH="amd64" ;;
      aarch64) ARCH="arm64" ;;
      *) echo "❌ Unsupported architecture: $ARCH" && exit 1 ;;
    esac

    GOLANGCI_VERSION="2.6.1"
    curl -fsSL "https://github.com/golangci/golangci-lint/releases/download/v${GOLANGCI_VERSION}/golangci-lint-${GOLANGCI_VERSION}-linux-${ARCH}.tar.gz" -o /tmp/golangci-lint.tar.gz
    tar -xzf /tmp/golangci-lint.tar.gz -C /tmp
    cp "/tmp/golangci-lint-${GOLANGCI_VERSION}-linux-${ARCH}/golangci-lint" "$HOME/.cache/go/bin/golangci-lint"
    chmod +x "$HOME/.cache/go/bin/golangci-lint"
    rm -rf /tmp/golangci-lint.tar.gz "/tmp/golangci-lint-${GOLANGCI_VERSION}-linux-${ARCH}"

    # Install other Go tools using absolute path
    "$GO_BIN" install github.com/bazelbuild/buildtools/buildifier@latest
    "$GO_BIN" install mvdan.cc/sh/v3/cmd/shfmt@latest
    echo "✅ Go tools installed successfully!"
  else
    echo "✅ Go tools already installed"
  fi
else
  echo "⚠️  Go not found at $GO_BIN, skipping Go tools installation"
fi

# Ensure git hooks are configured
if [ -d "/workspace/.git" ] && [ -f "/workspace/scripts/install-hooks.sh" ]; then
  echo "🪝 Configuring git hooks..."
  chmod +x /workspace/scripts/install-hooks.sh
  /workspace/scripts/install-hooks.sh
elif [ ! -d "/workspace/.git" ]; then
  echo "⚠️  Git repository not found, skipping git hooks configuration"
elif [ ! -f "/workspace/scripts/install-hooks.sh" ]; then
  echo "⚠️  Install hooks script not found, skipping git hooks configuration"
fi

# Setup MCP configuration
if [ -f "/workspace/.devcontainer/setup-mcp.sh" ]; then
  echo "⚙️  Setting up MCP..."
  /workspace/.devcontainer/setup-mcp.sh
fi

echo "✅ DevContainer ready!"
