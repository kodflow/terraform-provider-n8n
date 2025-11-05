#!/bin/bash
set -e

echo "🚀 Starting post-start configuration..."

# Install Go tools (now that Go is available in PATH)
if command -v go &>/dev/null; then
  # Check if tools are already installed to avoid reinstalling every time
  if [ ! -f "$HOME/.cache/go/bin/golangci-lint" ]; then
    echo "🔨 Installing Go tools..."

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

    # Install other Go tools
    go install github.com/bazelbuild/buildtools/buildifier@latest
    go install mvdan.cc/sh/v3/cmd/shfmt@latest
    echo "✅ Go tools installed successfully!"
  else
    echo "✅ Go tools already installed"
  fi
else
  echo "⚠️  Go not found in PATH, skipping Go tools installation"
fi

# Setup MCP configuration
if [ -f "/workspace/.devcontainer/setup-mcp.sh" ]; then
  echo "⚙️  Setting up MCP..."
  /workspace/.devcontainer/setup-mcp.sh
fi

echo "✅ DevContainer ready!"
