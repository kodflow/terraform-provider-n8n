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

echo "✨ Environment initialization complete!"
