#!/bin/bash
# Setup GPG signing for git commits and tags
# GPG keys are created by init.sh on the host machine before container build

set -e

echo "🔐 Setting up GPG signing..."

# Check if host GPG directory is mounted
if [ ! -d "/host-gpg" ]; then
  echo "ℹ️  No /host-gpg directory found"
  echo "ℹ️  Skipping GPG setup - commits will not be signed"
  exit 0
fi

# Check if GPG configuration file exists
if [ ! -f "/host-gpg/gpg-config.env" ]; then
  echo "ℹ️  No GPG configuration found in /host-gpg/gpg-config.env"
  echo "ℹ️  Skipping GPG setup - commits will not be signed"
  echo ""
  echo "💡 GPG keys are auto-created by init.sh on the host machine"
  echo "   If you see this message, GPG might not be installed on your host"
  exit 0
fi

# Load GPG configuration
# shellcheck disable=SC1091
source /host-gpg/gpg-config.env

if [ -z "$KEYID" ]; then
  echo "⚠️  KEYID not set in /host-gpg/gpg-config.env"
  echo "ℹ️  Skipping GPG setup"
  exit 0
fi

echo "📋 GPG Key ID: $KEYID"

# Clean up any stale GPG processes and sockets that might interfere
# This prevents "gpg-agent is older than us" and "Forbidden" errors
echo "🧹 Cleaning up stale GPG processes..."
pkill -9 -f gpg-agent 2>/dev/null || true
pkill -9 -f dirmngr 2>/dev/null || true
rm -f ~/.gnupg/*.lock ~/.gnupg/S.* ~/.gnupg/d.* 2>/dev/null || true
sleep 1

# Check if private key file exists
if [ -f "/host-gpg/private.key" ]; then
  echo "🔑 Importing GPG private key..."

  # Copy key to writable location to avoid "Forbidden" errors with readonly mounts
  cp /host-gpg/private.key /tmp/private.key.import
  chmod 600 /tmp/private.key.import

  # Import the key (suppress output for security)
  if gpg --batch --import /tmp/private.key.import >/dev/null 2>&1; then
    echo "✅ GPG key imported successfully"
  else
    echo "⚠️  Failed to import GPG key, but continuing..."
    rm -f /tmp/private.key.import
    exit 0
  fi

  # Clean up temporary key file
  rm -f /tmp/private.key.import
else
  echo "ℹ️  No private.key file found in /host-gpg/"
  echo "ℹ️  Checking if key is already in keyring..."
fi

# Verify the key is available
if ! gpg --list-secret-keys "$KEYID" >/dev/null 2>&1; then
  echo "⚠️  GPG key $KEYID not found in keyring"
  echo "ℹ️  Skipping GPG configuration"
  exit 0
fi

echo "✅ GPG key $KEYID is available"

# Set ultimate trust on the key (required for signing)
echo "🔒 Setting trust level..."
echo "$KEYID:6:" | gpg --import-ownertrust >/dev/null 2>&1 || true

echo "✅ GPG setup completed successfully"
