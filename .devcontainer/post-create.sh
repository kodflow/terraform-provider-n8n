#!/bin/bash
set -e

echo "🔧 Installing development tools..."

# Configure git identity
echo "👤 Configuring git identity..."
git config --global user.name "Kodflow"
git config --global user.email "133899878+kodflow@users.noreply.github.com"

# Configure GPG signing if GPG key is available
if [ -f "/host-gpg/gpg-config.env" ]; then
  echo "🔐 Configuring GPG signing..."
  source /host-gpg/gpg-config.env
  git config --global user.signingkey "$KEYID"
  git config --global gpg.program gpg
  git config --global commit.gpgsign true
  git config --global tag.gpgsign true
  echo "✅ GPG signing configured with key $KEYID"
else
  echo "ℹ️  No GPG key found, skipping GPG configuration"
fi

# Configure npm for local global packages (no sudo needed)
# Using environment variable instead of npm config to avoid conflicts with nvm
echo "⚙️  Configuring npm..."
mkdir -p "$HOME/.local/share/npm-global"
export NPM_CONFIG_PREFIX="$HOME/.local/share/npm-global"

# Install npm global packages
echo "📦 Installing npm packages..."
npm install -g @anthropic-ai/claude-code@latest
npm install -g @commitlint/cli@latest @commitlint/config-conventional@latest

# Install ktn-linter from GitHub releases (doesn't require Go)
echo "🧹 Installing ktn-linter..."
ARCH="$(dpkg --print-architecture)"
case "$ARCH" in
  amd64) KTN_ARCH="amd64" ;;
  arm64) KTN_ARCH="arm64" ;;
  *) echo "❌ Unsupported architecture: $ARCH" && exit 1 ;;
esac

KTN_VERSION=$(curl -s https://api.github.com/repos/kodflow/ktn-linter/releases/latest | grep '"tag_name"' | sed -E 's/.*"v([^"]+)".*/\1/')
curl -fsSL "https://github.com/kodflow/ktn-linter/releases/download/v${KTN_VERSION}/ktn-linter-linux-${KTN_ARCH}" -o "$HOME/.local/bin/ktn-linter"
chmod +x "$HOME/.local/bin/ktn-linter"

# Install golangci-lint v2 (supports Go 1.25+)
echo "🧹 Installing golangci-lint v2..."
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b "$HOME/.local/bin" v2.6.1

# Install git hooks
echo "🪝 Installing git hooks..."
if [ -f "./scripts/install-hooks.sh" ]; then
  chmod +x ./scripts/install-hooks.sh
  ./scripts/install-hooks.sh
else
  echo "⚠️  scripts/install-hooks.sh not found, skipping hooks installation"
fi

echo "✅ Development tools installed successfully!"
echo "ℹ️  Go tools will be installed on container start..."
