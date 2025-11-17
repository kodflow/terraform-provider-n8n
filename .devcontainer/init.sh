#!/bin/bash
# Initialize .env file before devcontainer build
# This script runs on the host machine before Docker Compose

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ENV_FILE="$SCRIPT_DIR/.env"

# Extract project name from git remote URL
REPO_NAME=$(basename "$(git config --get remote.origin.url)" .git)

# Sanitize project name for Docker Compose requirements:
# - Must start with a letter or number
# - Only lowercase alphanumeric, hyphens, and underscores allowed
REPO_NAME=$(echo "$REPO_NAME" | sed 's/^[^a-zA-Z0-9]*//' | tr '[:upper:]' '[:lower:]' | sed 's/[^a-z0-9_-]/-/g')

# If name is empty after sanitization, use a default
if [ -z "$REPO_NAME" ]; then
  REPO_NAME="devcontainer"
fi

echo "🔧 Initializing devcontainer environment..."
echo "📦 Project name: $REPO_NAME"

# If .env doesn't exist, create it from .env.example
if [ ! -f "$ENV_FILE" ]; then
  echo "📝 Creating .env from .env.example..."
  cp "$SCRIPT_DIR/.env.example" "$ENV_FILE"
fi

# Update or add COMPOSE_PROJECT_NAME in .env
if grep -q "^COMPOSE_PROJECT_NAME=" "$ENV_FILE"; then
  # Update existing line
  if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS
    sed -i '' "s|^COMPOSE_PROJECT_NAME=.*|COMPOSE_PROJECT_NAME=$REPO_NAME|" "$ENV_FILE"
  else
    # Linux
    sed -i "s|^COMPOSE_PROJECT_NAME=.*|COMPOSE_PROJECT_NAME=$REPO_NAME|" "$ENV_FILE"
  fi
  echo "✅ Updated COMPOSE_PROJECT_NAME=$REPO_NAME in .env"
else
  # Add at the beginning of the file
  echo "COMPOSE_PROJECT_NAME=$REPO_NAME" | cat - "$ENV_FILE" >"$ENV_FILE.tmp" && mv "$ENV_FILE.tmp" "$ENV_FILE"
  echo "✅ Added COMPOSE_PROJECT_NAME=$REPO_NAME to .env"
fi

# Create and configure GPG directory
GPG_DIR="$HOME/.config/devcontainer-gpg"
if [ ! -d "$GPG_DIR" ]; then
  echo "📁 Creating GPG directory: $GPG_DIR"
  mkdir -p "$GPG_DIR"
  chmod 700 "$GPG_DIR"
  echo "✅ GPG directory created"
else
  # Check if directory is writable
  if [ ! -w "$GPG_DIR" ]; then
    echo "⚠️  $GPG_DIR exists but is not writable (wrong ownership)"

    # Detect current user
    CURRENT_USER=$(whoami)

    # Backup strategy: rename old directory and create new one
    BACKUP_DIR="${GPG_DIR}.backup.$(date +%s)"

    echo "🔧 Attempting to fix automatically..."

    # Try method 1: sudo without password (if configured)
    if sudo -n chown -R "$CURRENT_USER:$CURRENT_USER" "$GPG_DIR" 2>/dev/null; then
      echo "✅ Ownership fixed with sudo (no password)"
      chmod 700 "$GPG_DIR"
    # Try method 2: rename old dir and create new one (works if parent dir is writable)
    elif mv "$GPG_DIR" "$BACKUP_DIR" 2>/dev/null; then
      echo "✅ Moved old directory to $BACKUP_DIR"
      mkdir -p "$GPG_DIR"
      chmod 700 "$GPG_DIR"
      echo "✅ Created new GPG directory"
      echo "ℹ️  Old directory backed up, you can delete it with: sudo rm -rf $BACKUP_DIR"
    else
      # All automatic methods failed
      echo "❌ Automatic fix failed - manual intervention required"
      echo ""
      echo "Please run ONE of these commands:"
      echo "  1. sudo chown -R $CURRENT_USER:$CURRENT_USER $GPG_DIR"
      echo "  2. sudo rm -rf $GPG_DIR  (will delete existing GPG keys)"
      echo ""
      exit 1
    fi
  fi
fi

# Auto-create or validate GPG key (on host machine)
echo "🔑 Checking GPG configuration..."

# Check if gpg is available on host
if ! command -v gpg &> /dev/null; then
  echo "⚠️  GPG not found on host machine"
  echo "ℹ️  Install GPG: apt install gnupg (Linux) or brew install gnupg (macOS)"
  echo "ℹ️  Skipping GPG key creation - commits will not be signed"
else
  # Git user info - always read from host git config
  echo "📋 Reading Git configuration from host..."
  GIT_NAME=$(git config --global user.name 2>/dev/null || echo "")
  GIT_EMAIL=$(git config --global user.email 2>/dev/null || echo "")

  # Validate and provide defaults if needed
  if [ -z "$GIT_NAME" ]; then
    echo "⚠️  Git user.name not configured on host"
    GIT_NAME="Developer"
    echo "ℹ️  Using default name: $GIT_NAME"
    echo "💡 Configure with: git config --global user.name \"Your Name\""
  fi

  if [ -z "$GIT_EMAIL" ]; then
    echo "⚠️  Git user.email not configured on host"
    GIT_EMAIL="developer@example.com"
    echo "ℹ️  Using default email: $GIT_EMAIL"
    echo "💡 Configure with: git config --global user.email \"your.email@example.com\""
  fi

  echo "✅ Git identity: $GIT_NAME <$GIT_EMAIL>"

  # Check if gpg-config.env exists and compare with current git config
  NEEDS_REGENERATION=false
  if [ -f "$GPG_DIR/gpg-config.env" ]; then
    # Save current git config before sourcing
    CURRENT_GIT_NAME="$GIT_NAME"
    CURRENT_GIT_EMAIL="$GIT_EMAIL"

    # shellcheck disable=SC1090
    source "$GPG_DIR/gpg-config.env"

    # Remove quotes from stored values for comparison
    STORED_NAME=$(echo "$GIT_NAME" | tr -d '"')
    STORED_EMAIL=$(echo "$GIT_EMAIL" | tr -d '"')
    STORED_KEYID=$(echo "$KEYID" | tr -d '"')

    # Restore current git config
    GIT_NAME="$CURRENT_GIT_NAME"
    GIT_EMAIL="$CURRENT_GIT_EMAIL"

    # Compare with current git config
    if [ "$STORED_NAME" != "$GIT_NAME" ] || [ "$STORED_EMAIL" != "$GIT_EMAIL" ]; then
      echo "⚠️  Git identity has changed!"
      echo "   Old: $STORED_NAME <$STORED_EMAIL>"
      echo "   New: $GIT_NAME <$GIT_EMAIL>"
      NEEDS_REGENERATION=true
    else
      # Identity matches, but verify GPG key exists in keyring
      if ! gpg --list-secret-keys --keyid-format=long "$STORED_KEYID" >/dev/null 2>&1; then
        echo "⚠️  GPG key $STORED_KEYID not found in keyring"
        NEEDS_REGENERATION=true
      else
        echo "✅ GPG configuration is up to date"
        # No need to regenerate, we're done
        echo "✨ Environment initialization complete!"
        exit 0
      fi
    fi
  else
    echo "ℹ️  No existing GPG configuration found"
    NEEDS_REGENERATION=true
  fi

  # Generate or regenerate GPG key if needed
  if [ "$NEEDS_REGENERATION" = true ]; then
    echo "🔄 Setting up GPG key for: $GIT_NAME <$GIT_EMAIL>"

    # Check if key already exists in host GPG keyring for this email
    EXISTING_KEY=$(gpg --list-secret-keys --keyid-format=long "$GIT_EMAIL" 2>/dev/null | grep "sec" | head -1 | awk '{print $2}' | cut -d'/' -f2 || true)

    if [ -n "$EXISTING_KEY" ]; then
      echo "✅ Found existing GPG key for $GIT_EMAIL: $EXISTING_KEY"
      KEYID="$EXISTING_KEY"
    else
      echo "📝 Generating new GPG key (no passphrase)..."

      # Generate key with batch mode (no passphrase)
      gpg --batch --gen-key <<EOF
Key-Type: RSA
Key-Length: 4096
Subkey-Type: RSA
Subkey-Length: 4096
Name-Real: $GIT_NAME
Name-Email: $GIT_EMAIL
Expire-Date: 0
%no-protection
%commit
EOF

      # Get the new key ID
      KEYID=$(gpg --list-secret-keys --keyid-format=long "$GIT_EMAIL" 2>/dev/null | grep "sec" | head -1 | awk '{print $2}' | cut -d'/' -f2 || true)

      if [ -z "$KEYID" ]; then
        echo "❌ Failed to generate GPG key"
        exit 1
      else
        echo "✅ GPG key generated: $KEYID"
      fi
    fi

    # Export keys and configuration
    if [ -n "$KEYID" ]; then
      # Export private key
      gpg --export-secret-keys --armor "$KEYID" > "$GPG_DIR/private.key"
      chmod 600 "$GPG_DIR/private.key"
      echo "✅ Private key exported to $GPG_DIR/private.key"

      # Export public key
      gpg --armor --export "$KEYID" > "$GPG_DIR/public.key"
      chmod 644 "$GPG_DIR/public.key"
      echo "✅ Public key exported to $GPG_DIR/public.key"

      # Create config file with Git identity (quote values for safety)
      cat > "$GPG_DIR/gpg-config.env" <<EOF
KEYID="$KEYID"
GIT_NAME="$GIT_NAME"
GIT_EMAIL="$GIT_EMAIL"
EOF
      chmod 644 "$GPG_DIR/gpg-config.env"
      echo "✅ Configuration saved to $GPG_DIR/gpg-config.env"

      echo ""
      echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
      echo "📤 Add this GPG public key to GitHub:"
      echo "   https://github.com/settings/gpg/new"
      echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
      echo ""
      cat "$GPG_DIR/public.key"
      echo ""
      echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
      echo "💡 Public key saved to: $GPG_DIR/public.key"
      echo ""
    fi
  fi
fi

echo "✨ Environment initialization complete!"
