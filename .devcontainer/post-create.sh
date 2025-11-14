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
