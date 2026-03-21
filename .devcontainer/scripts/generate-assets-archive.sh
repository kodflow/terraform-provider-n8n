#!/bin/bash
# ============================================================================
# generate-assets-archive.sh - Creates tar.gz of Claude Code assets
# ============================================================================
# Bundles all assets into a single downloadable archive.
# Used by CI (release.yml) to attach assets to GitHub Releases.
#
# Usage:
#   ./generate-assets-archive.sh                    # Default output path
#   ./generate-assets-archive.sh --output out.tar.gz # Custom output path
# ============================================================================

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
CLAUDE_DIR="$REPO_ROOT/.devcontainer/images/.claude"
OUTPUT_FILE="$REPO_ROOT/.devcontainer/claude-assets.tar.gz"

# Parse arguments
while [[ $# -gt 0 ]]; do
  case "$1" in
    --output)
      OUTPUT_FILE="$2"
      # Resolve relative paths from caller's working directory
      if [[ "$OUTPUT_FILE" != /* ]]; then
        OUTPUT_FILE="$(pwd)/$OUTPUT_FILE"
      fi
      shift 2
      ;;
    *)
      echo "Usage: $0 [--output <path>]"
      exit 1
      ;;
  esac
done

# Check if source directory exists
if [ ! -d "$CLAUDE_DIR" ]; then
  echo "Error: Claude assets directory not found: $CLAUDE_DIR"
  exit 1
fi

echo "→ Generating Claude Code assets archive..."

# Create tar.gz with relative paths
# Contents: agents/, commands/, scripts/, docs/, settings.json
cd "$CLAUDE_DIR"

tar -czf "$OUTPUT_FILE" \
  --exclude='*.pyc' \
  --exclude='__pycache__' \
  --exclude='.DS_Store' \
  agents/ \
  commands/ \
  scripts/ \
  docs/ \
  settings.json \
  2>/dev/null || {
  echo "Error: Failed to create archive"
  exit 1
}

# Show archive info
ARCHIVE_SIZE=$(du -h "$OUTPUT_FILE" | cut -f1)
FILE_COUNT=$(tar -tzf "$OUTPUT_FILE" | wc -l)

echo "  ✓ Archive created: $OUTPUT_FILE"
echo "  ✓ Size: $ARCHIVE_SIZE"
echo "  ✓ Files: $FILE_COUNT"

echo "→ Done"
