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

# Ensure Oh My Zsh is installed (handles empty volume on first run)
if [ ! -f "$HOME/.oh-my-zsh/oh-my-zsh.sh" ]; then
  echo "📦 Installing Oh My Zsh..."
  rm -rf "$HOME/.oh-my-zsh"
  sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)" "" --unattended --keep-zshrc
  # Install themes
  git clone --depth=1 https://github.com/romkatv/powerlevel10k.git \
    "${ZSH_CUSTOM:-$HOME/.oh-my-zsh/custom}/themes/powerlevel10k" 2>/dev/null || true
  git clone --depth=1 https://github.com/dracula/zsh.git \
    "${ZSH_CUSTOM:-$HOME/.oh-my-zsh/custom}/themes/dracula" 2>/dev/null || true
  echo "✅ Oh My Zsh installed"
fi

# Configure git identity
echo "👤 Configuring git identity..."

# Read Git identity from GPG config if available
if [ -f "/host-gpg/gpg-config.env" ]; then
  # shellcheck disable=SC1091
  source /host-gpg/gpg-config.env

  if [ -n "$GIT_NAME" ] && [ -n "$GIT_EMAIL" ]; then
    git config --global user.name "$GIT_NAME"
    git config --global user.email "$GIT_EMAIL"
    echo "   ✅ Using Git identity from host: $GIT_NAME <$GIT_EMAIL>"
  else
    # Fallback to defaults if not set in config
    git config --global user.name "Developer"
    git config --global user.email "developer@example.com"
    echo "   ⚠️  Git identity not found in config, using defaults"
  fi
else
  # No GPG config, use defaults
  git config --global user.name "Developer"
  git config --global user.email "developer@example.com"
  echo "   ⚠️  No GPG config found, using default Git identity"
  echo "   💡 Run init.sh on host to configure Git identity and GPG"
fi

# Configure GPG signing if GPG key is available (imported by .devcontainer/setup-gpg.sh)
if [ -f "/host-gpg/gpg-config.env" ]; then
  echo "🔐 Configuring Git GPG signing..."
  # shellcheck disable=SC1091
  source /host-gpg/gpg-config.env

  # Verify GPG key is available (should be imported by scripts/setup-gpg.sh)
  if gpg --list-secret-keys "$KEYID" >/dev/null 2>&1; then
    # Configure Git to use GPG signing
    git config --global user.signingkey "$KEYID"
    git config --global gpg.program gpg
    git config --global commit.gpgsign true
    git config --global tag.gpgsign true

    echo "   ✅ Git configured to sign commits and tags with key $KEYID"
  else
    echo "   ⚠️  GPG key $KEYID not found - commits will not be signed"
    echo "   ℹ️  Run scripts/setup-gpg.sh manually if you want to enable signing"
  fi
else
  echo "ℹ️  No GPG configuration found, commits will not be signed"
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
# shellcheck disable=SC2016
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
  # shellcheck disable=SC1091
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
