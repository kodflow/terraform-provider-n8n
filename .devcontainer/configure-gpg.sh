#!/bin/bash
# Helper script to configure GPG for signing commits
# Run this script inside the devcontainer to generate/export your GPG key

set -e

echo "🔐 GPG Configuration Helper"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Check if we're inside the container
if [ ! -f "/.dockerenv" ] && [ ! -f "/run/.containerenv" ]; then
  echo "⚠️  This script should be run inside the devcontainer"
  echo "💡 Open the project in VS Code and reopen in container first"
  exit 1
fi

# Check if GPG is installed
if ! command -v gpg &>/dev/null; then
  echo "❌ GPG is not installed"
  exit 1
fi

# Git user info - read from container's git config
echo "📋 Reading Git configuration..."
GIT_NAME=$(git config --global user.name 2>/dev/null || echo "")
GIT_EMAIL=$(git config --global user.email 2>/dev/null || echo "")

# Validate and prompt if needed
if [ -z "$GIT_NAME" ]; then
  read -p "Enter your name: " GIT_NAME
  if [ -z "$GIT_NAME" ]; then
    echo "❌ Name is required"
    exit 1
  fi
fi

if [ -z "$GIT_EMAIL" ]; then
  read -p "Enter your email: " GIT_EMAIL
  if [ -z "$GIT_EMAIL" ]; then
    echo "❌ Email is required"
    exit 1
  fi
fi

echo "📋 Git Identity:"
echo "   Name:  $GIT_NAME"
echo "   Email: $GIT_EMAIL"
echo ""

# Check if a GPG key already exists for this email
EXISTING_KEY=$(gpg --list-secret-keys --keyid-format=long "$GIT_EMAIL" 2>/dev/null | grep "sec" | head -1 | awk '{print $2}' | cut -d'/' -f2)

if [ -n "$EXISTING_KEY" ]; then
  echo "✅ Found existing GPG key: $EXISTING_KEY"
  echo ""

  # Ask if user wants to use existing key or create new one
  read -p "Do you want to use this existing key? (y/n): " -n 1 -r
  echo ""

  if [[ $REPLY =~ ^[Yy]$ ]]; then
    KEYID="$EXISTING_KEY"
  else
    echo "Creating a new key..."
    KEYID=""
  fi
fi

# Generate new key if needed
if [ -z "$KEYID" ]; then
  echo "🔑 Generating new GPG key..."
  echo ""
  echo "Please enter a passphrase when prompted (can be empty for no passphrase)"
  echo ""

  # Generate key with batch mode
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
  KEYID=$(gpg --list-secret-keys --keyid-format=long "$GIT_EMAIL" 2>/dev/null | grep "sec" | head -1 | awk '{print $2}' | cut -d'/' -f2)

  if [ -z "$KEYID" ]; then
    echo "❌ Failed to generate GPG key"
    exit 1
  fi

  echo ""
  echo "✅ GPG key generated successfully: $KEYID"
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📤 Exporting GPG Public Key for GitHub"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Export public key
echo "Copy this public key and add it to GitHub:"
echo "https://github.com/settings/gpg/new"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
gpg --armor --export "$KEYID"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Save to file
GPG_PUBLIC_KEY_FILE="/workspace/.devcontainer/gpg-public-key.asc"
gpg --armor --export "$KEYID" >"$GPG_PUBLIC_KEY_FILE"
echo "✅ Public key also saved to: $GPG_PUBLIC_KEY_FILE"
echo ""

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "💾 Configuring Host GPG Directory"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Instructions for host configuration
if [ -d "/host-gpg" ]; then
  echo "✅ /host-gpg directory is mounted"

  # Try to export key to host directory
  if [ -w "/host-gpg" ]; then
    echo "📝 Creating configuration files in /host-gpg..."

    # Export private key
    gpg --export-secret-keys --armor "$KEYID" >/host-gpg/private.key
    chmod 600 /host-gpg/private.key
    echo "✅ Private key exported to /host-gpg/private.key"

    # Create config file with Git identity (quote values for safety)
    cat >/host-gpg/gpg-config.env <<EOF
KEYID="$KEYID"
GIT_NAME="$GIT_NAME"
GIT_EMAIL="$GIT_EMAIL"
EOF
    chmod 644 /host-gpg/gpg-config.env
    echo "✅ Configuration saved to /host-gpg/gpg-config.env"

    echo ""
    echo "✅ GPG is now configured for this devcontainer!"
  else
    echo "⚠️  /host-gpg directory is read-only"
    echo ""
    echo "Please run these commands ON YOUR HOST MACHINE (Windows/WSL):"
    echo ""
    echo "  mkdir -p ~/.config/devcontainer-gpg"
    echo "  gpg --export-secret-keys --armor $KEYID > ~/.config/devcontainer-gpg/private.key"
    echo "  chmod 600 ~/.config/devcontainer-gpg/private.key"
    echo "  cat > ~/.config/devcontainer-gpg/gpg-config.env <<'EOF'"
    echo "KEYID=\"$KEYID\""
    echo "GIT_NAME=\"$GIT_NAME\""
    echo "GIT_EMAIL=\"$GIT_EMAIL\""
    echo "EOF"
    echo ""
  fi
else
  echo "ℹ️  /host-gpg directory is not mounted"
  echo ""
  echo "To enable GPG signing on container restart:"
  echo ""
  echo "1. Create the directory on your host:"
  echo "   mkdir -p ~/.config/devcontainer-gpg"
  echo ""
  echo "2. Export your private key:"
  echo "   (Run this command INSIDE the container now)"
  echo "   gpg --export-secret-keys --armor $KEYID > /tmp/private.key"
  echo ""
  echo "   Then copy it to host:"
  echo "   (From your host machine)"
  echo "   docker cp \$(docker ps -q -f name=devcontainer):/tmp/private.key ~/.config/devcontainer-gpg/"
  echo "   chmod 600 ~/.config/devcontainer-gpg/private.key"
  echo ""
  echo "3. Create configuration file on host:"
  echo "   cat > ~/.config/devcontainer-gpg/gpg-config.env <<'EOF'"
  echo "KEYID=\"$KEYID\""
  echo "GIT_NAME=\"$GIT_NAME\""
  echo "GIT_EMAIL=\"$GIT_EMAIL\""
  echo "EOF"
  echo ""
  echo "4. Update .devcontainer/devcontainer.json to include the mount:"
  echo "   (Already configured, just rebuild the container)"
  echo ""
  echo "5. Rebuild the devcontainer"
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ GPG Configuration Complete!"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Next steps:"
echo "1. Add the public key to GitHub: https://github.com/settings/gpg/new"
echo "2. Restart the devcontainer if you configured the host directory"
echo "3. Your commits will now be signed automatically!"
echo ""
echo "Key ID: $KEYID"
echo ""
