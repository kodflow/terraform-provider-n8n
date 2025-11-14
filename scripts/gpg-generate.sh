#!/bin/bash
# Copyright (c) 2024 Florent (Kodflow). All rights reserved.
# Licensed under the Sustainable Use License 1.0
# See LICENSE.md in the project root for license information.

set -e

NAME="${1:-$(git config --global user.name || echo "Kodflow")}"
EMAIL="${2:-$(git config --global user.email || echo "133899878+kodflow@users.noreply.github.com")}"

echo "🔐 Generating GPG key..."
echo ""
echo "Name:  $NAME"
echo "Email: $EMAIL"
echo ""

# Check if key already exists
EXISTING_KEY=$(gpg --list-secret-keys --with-colons "$EMAIL" 2>/dev/null | awk -F: '/sec:/ {print $5; exit}')
if [ -n "$EXISTING_KEY" ]; then
    echo "⚠️  GPG key already exists: $EXISTING_KEY"
    echo "   Run 'make gpg/delete' to remove it first"
    exit 1
fi

# Create batch file for key generation
cat > /tmp/gpg-batch <<EOF
%echo Generating OpenPGP key for $EMAIL
Key-Type: eddsa
Key-Curve: Ed25519
Key-Usage: sign
Name-Real: $NAME
Name-Email: $EMAIL
Expire-Date: 0
%no-protection
%commit
%echo done
EOF

# Ensure GPG home exists
mkdir -p ~/.gnupg
chmod 700 ~/.gnupg

# Generate the key
echo "Generating key (this may take a moment)..."
gpg --batch --gen-key /tmp/gpg-batch 2>&1 || {
    echo ""
    echo "❌ GPG key generation failed!"
    echo ""
    echo "Please run this command manually:"
    echo ""
    echo "  gpg --full-gen-key"
    echo ""
    echo "Then select:"
    echo "  - (9) ECC (sign only)"
    echo "  - (1) Curve 25519"
    echo "  - Key does not expire"
    echo "  - Real name: $NAME"
    echo "  - Email: $EMAIL"
    echo "  - No passphrase (press ENTER twice)"
    echo ""
    rm -f /tmp/gpg-batch
    exit 1
}

# Clean up
rm -f /tmp/gpg-batch

echo ""
echo "✅ GPG key generated successfully"
echo ""

# Show key info
NEW_KEY=$(gpg --list-secret-keys --with-colons "$EMAIL" 2>/dev/null | awk -F: '/sec:/ {print $5; exit}')
echo "Key ID: $NEW_KEY"
echo ""
gpg --list-secret-keys "$EMAIL"
