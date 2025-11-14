#!/bin/bash
set -e

echo "🔧 Installing development tools..."

# Configure git identity
echo "👤 Configuring git identity..."
git config --global user.name "Kodflow"
git config --global user.email "133899878+kodflow@users.noreply.github.com"

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

# Configure GPG for unattended key generation
echo "🔐 Configuring GPG..."
mkdir -p "$HOME/.gnupg"
chmod 700 "$HOME/.gnupg"

cat > "$HOME/.gnupg/gpg.conf" <<GPGCONF
# Allow generation of keys without passphrase
allow-freeform-uid
GPGCONF

cat > "$HOME/.gnupg/gpg-agent.conf" <<AGENTCONF
# Allow unattended passphrase entry
allow-preset-passphrase
allow-loopback-pinentry
max-cache-ttl 34560000
default-cache-ttl 34560000
pinentry-mode loopback
AGENTCONF

chmod 600 "$HOME/.gnupg/gpg.conf" "$HOME/.gnupg/gpg-agent.conf"

# Kill any existing gpg-agent
gpgconf --kill gpg-agent 2>/dev/null || true

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
