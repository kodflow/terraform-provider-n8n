#!/bin/bash
# Exit on error, but continue if non-critical commands fail
set -e

# Function to safely run commands that shouldn't stop setup
safe_run() {
  if ! "$@"; then
    echo "⚠️  Warning: Command failed but continuing setup: $*"
    return 0
  fi
  return 0
}

echo "🔧 Installing development tools..."

# Configure git identity
echo "👤 Configuring git identity..."
git config --global user.name "Kodflow"
git config --global user.email "133899878+kodflow@users.noreply.github.com"

# Configure GPG signing if GPG key is available
if [ -f "/host-gpg/gpg-config.env" ]; then
  echo "🔐 Configuring GPG signing..."
  source /host-gpg/gpg-config.env

  # Verify GPG key is imported
  echo "   Checking GPG key import..."
  if gpg --list-secret-keys "$KEYID" >/dev/null 2>&1; then
    echo "   ✅ GPG key $KEYID is imported"

    # Display key info
    KEY_INFO=$(gpg --list-keys "$KEYID" 2>/dev/null | grep -A 1 "^pub" | tail -n 1 | xargs)
    echo "   📋 Key: $KEY_INFO"
  else
    echo "   ❌ GPG key $KEYID not found in keyring"
    exit 1
  fi

  # Configure Git
  git config --global user.signingkey "$KEYID"
  git config --global gpg.program gpg
  git config --global commit.gpgsign true
  git config --global tag.gpgsign true

  # Verify Git configuration
  echo "   Verifying Git GPG configuration..."
  SIGNING_KEY=$(git config --global user.signingkey)
  COMMIT_SIGN=$(git config --global commit.gpgsign)
  TAG_SIGN=$(git config --global tag.gpgsign)

  if [ "$SIGNING_KEY" = "$KEYID" ] && [ "$COMMIT_SIGN" = "true" ] && [ "$TAG_SIGN" = "true" ]; then
    echo "   ✅ Git configured to sign commits and tags with key $KEYID"
  else
    echo "   ❌ Git GPG configuration verification failed"
    exit 1
  fi
else
  echo "ℹ️  No GPG key found, skipping GPG configuration"
fi

# Configure npm for local global packages (no sudo needed)
# Using environment variable instead of npm config to avoid conflicts with nvm
echo "⚙️  Configuring npm..."

# Create directories with proper ownership
mkdir -p "$HOME/.local/share/npm-global" 2>/dev/null || true
chmod -R 755 "$HOME/.local" 2>/dev/null || true

# Set npm prefix for this session and future sessions
export NPM_CONFIG_PREFIX="$HOME/.local/share/npm-global"

# Install npm global packages with explicit prefix
echo "📦 Installing npm packages..."
if npm install -g --prefix "$HOME/.local/share/npm-global" @anthropic-ai/claude-code@latest 2>&1 | tee /tmp/claude-install.log; then
  echo "✅ claude-code installed successfully"
else
  echo "⚠️  Failed to install claude-code, check /tmp/claude-install.log for details"
fi

if npm install -g --prefix "$HOME/.local/share/npm-global" @commitlint/cli@latest @commitlint/config-conventional@latest 2>&1 | tee /tmp/commitlint-install.log; then
  echo "✅ commitlint installed successfully"
else
  echo "⚠️  Failed to install commitlint, check /tmp/commitlint-install.log for details"
fi

# Install ktn-linter from GitHub releases (doesn't require Go)
echo "🧹 Installing ktn-linter..."
ARCH="$(dpkg --print-architecture)"
case "$ARCH" in
  amd64) KTN_ARCH="amd64" ;;
  arm64) KTN_ARCH="arm64" ;;
  *) echo "⚠️  Unsupported architecture: $ARCH, skipping ktn-linter" && KTN_ARCH="" ;;
esac

if [ -n "$KTN_ARCH" ]; then
  KTN_VERSION=$(curl -s https://api.github.com/repos/kodflow/ktn-linter/releases/latest | grep '"tag_name"' | sed -E 's/.*"v([^"]+)".*/\1/')
  if [ -n "$KTN_VERSION" ]; then
    safe_run curl -fsSL "https://github.com/kodflow/ktn-linter/releases/download/v${KTN_VERSION}/ktn-linter-linux-${KTN_ARCH}" -o "$HOME/.local/bin/ktn-linter" || echo "⚠️  Failed to download ktn-linter"
    chmod +x "$HOME/.local/bin/ktn-linter" 2>/dev/null || true
    echo "✅ ktn-linter installed successfully"
  else
    echo "⚠️  Failed to get ktn-linter version, skipping"
  fi
fi

# Install golangci-lint v2 (supports Go 1.25+)
echo "🧹 Installing golangci-lint v2..."
safe_run bash -c 'curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b "$HOME/.local/bin" v2.6.1' || echo "⚠️  Failed to install golangci-lint"

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

# Display GPG configuration summary
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📋 Configuration Summary"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Git User:  $(git config --global user.name) <$(git config --global user.email)>"

if [ -f "/host-gpg/gpg-config.env" ]; then
  source /host-gpg/gpg-config.env
  echo "GPG Key:   $KEYID"
  echo "Signing:   Commits ✅ | Tags ✅"
  echo ""
  echo "🔐 All commits and tags will be automatically signed!"
else
  echo "GPG Key:   Not configured"
  echo "Signing:   Disabled"
fi
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "💡 Quick Tips:"
echo "  • Use 'super-claude' alias to run Claude with MCP servers"
echo "  • Run 'make help' to see all available commands"
echo "  • Run 'make test' to execute the test suite"
echo ""
