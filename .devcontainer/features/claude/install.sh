#!/bin/bash
# ============================================================================
# Claude Code Marketplace - One-liner Install
# ============================================================================
# curl -sL https://raw.githubusercontent.com/kodflow/devcontainer-template/main/.devcontainer/features/claude/install.sh | bash
# ============================================================================

set -e

REPO="kodflow/devcontainer-template"
BRANCH="main"
BASE="https://raw.githubusercontent.com/${REPO}/${BRANCH}/.devcontainer/images"

# DC_TARGET: Override installation directory (defaults to current working directory)
# Usage: DC_TARGET=/path/to/project ./install.sh
TARGET="${DC_TARGET:-$(pwd)}"

echo "═══════════════════════════════════════════"
echo "  Claude Code Marketplace"
echo "═══════════════════════════════════════════"
echo ""

# ─────────────────────────────────────────────────────────────────────────────
# 0. Read feature options from devcontainer-feature.json (passed as env vars)
# ─────────────────────────────────────────────────────────────────────────────
INSTALL_CLAUDE="${INSTALLCLAUDE:-true}"
BUNDLE_LEVEL="${BUNDLE:-full}"
INSTALL_STATUSLINE="${STATUSLINE:-true}"

echo "  Options: claude=${INSTALL_CLAUDE}, bundle=${BUNDLE_LEVEL}, statusline=${INSTALL_STATUSLINE}"
echo ""

# ─────────────────────────────────────────────────────────────────────────────
# 1. Install Claude CLI (if not already installed)
# ─────────────────────────────────────────────────────────────────────────────
if [ "$INSTALL_CLAUDE" = "true" ] && ! command -v claude &>/dev/null; then
  echo "→ Installing Claude CLI..."
  npm install -g @anthropic-ai/claude-code 2>/dev/null \
    || curl -fsSL https://claude.ai/install.sh | sh 2>/dev/null || true
elif [ "$INSTALL_CLAUDE" != "true" ]; then
  echo "→ Skipping Claude CLI installation (installClaude=false)"
fi

# ─────────────────────────────────────────────────────────────────────────────
# 2. Create directories
# ─────────────────────────────────────────────────────────────────────────────
echo "→ Setting up $TARGET/.claude/..."
mkdir -p "$TARGET/.claude/commands"
mkdir -p "$TARGET/.claude/scripts"

# ─────────────────────────────────────────────────────────────────────────────
# 3. Download commands
# ─────────────────────────────────────────────────────────────────────────────
echo "→ Downloading commands..."
for cmd in git search prompt; do
  curl -sL "$BASE/.claude/commands/$cmd.md" -o "$TARGET/.claude/commands/$cmd.md" 2>/dev/null && echo "  ✓ /$cmd"
done

# ─────────────────────────────────────────────────────────────────────────────
# 4. Download scripts (hooks) - requires standard or full bundle
# ─────────────────────────────────────────────────────────────────────────────
if [ "$BUNDLE_LEVEL" = "minimal" ]; then
  echo "→ Skipping scripts download (bundle=minimal)"
else
  echo "→ Downloading scripts..."
  download_failed=0
  for script in format imports lint post-edit pre-validate security test bash-validate commit-validate post-compact on-stop notification session-init; do
    script_tmp="$(mktemp)"
    if curl -fsL --retry 2 "$BASE/.claude/scripts/$script.sh" -o "$script_tmp" 2>/dev/null; then
      install -m 0755 "$script_tmp" "$TARGET/.claude/scripts/$script.sh"
    else
      echo "  ⚠ Failed to download: $script.sh" >&2
      download_failed=1
    fi
    rm -f "$script_tmp"
  done
  [ "$download_failed" -eq 0 ] && echo "  ✓ hooks (format, lint, security...)" || echo "  ⚠ Some hooks failed to download"
fi

# ─────────────────────────────────────────────────────────────────────────────
# 5. Download settings.json
# ─────────────────────────────────────────────────────────────────────────────
echo "→ Downloading settings..."
curl -sL "$BASE/.claude/settings.json" -o "$TARGET/.claude/settings.json" 2>/dev/null
echo "  ✓ settings.json"

# ─────────────────────────────────────────────────────────────────────────────
# 6. Download CLAUDE.md (if not existing)
# ─────────────────────────────────────────────────────────────────────────────
if [ ! -f "$TARGET/CLAUDE.md" ]; then
  curl -sL "https://raw.githubusercontent.com/${REPO}/${BRANCH}/CLAUDE.md" -o "$TARGET/CLAUDE.md" 2>/dev/null
  echo "  ✓ CLAUDE.md"
fi

# ─────────────────────────────────────────────────────────────────────────────
# 7. Install grepai (semantic code search MCP) — skip if pre-installed
# ─────────────────────────────────────────────────────────────────────────────
if command -v grepai &>/dev/null; then
  echo "  ✓ grepai already installed: $(grepai --version 2>/dev/null || echo 'pre-installed')"
else
  echo "→ Installing grepai..."
  mkdir -p "$HOME/.local/bin"

  # Detect OS
  case "$(uname -s)" in
    Linux*) GREPAI_OS="linux" ;;
    Darwin*) GREPAI_OS="darwin" ;;
    MINGW* | MSYS* | CYGWIN*) GREPAI_OS="windows" ;;
    *) GREPAI_OS="linux" ;;
  esac

  # Detect architecture
  case "$(uname -m)" in
    x86_64 | amd64) GREPAI_ARCH="amd64" ;;
    aarch64 | arm64) GREPAI_ARCH="arm64" ;;
    *) GREPAI_ARCH="amd64" ;;
  esac

  # Extension for Windows
  GREPAI_EXT=""
  [ "$GREPAI_OS" = "windows" ] && GREPAI_EXT=".exe"

  # Fetch latest version and download correct tar.gz asset
  GREPAI_LATEST=$(curl -fsSL "https://api.github.com/repos/yoanbernabeu/grepai/releases/latest" 2>/dev/null | grep -o '"tag_name": *"[^"]*"' | head -1 | cut -d'"' -f4)
  GREPAI_URL=""
  if [ -n "$GREPAI_LATEST" ]; then
    GREPAI_URL="https://github.com/yoanbernabeu/grepai/releases/download/${GREPAI_LATEST}/grepai_${GREPAI_LATEST#v}_${GREPAI_OS}_${GREPAI_ARCH}.tar.gz"
  fi
  grepai_tmp="$(mktemp)"
  grepai_extract="$(mktemp -d)"
  if curl -fsL --retry 3 --retry-delay 1 --proto '=https' --tlsv1.2 "$GREPAI_URL" -o "$grepai_tmp" 2>/dev/null \
    && tar -xzf "$grepai_tmp" -C "$grepai_extract" grepai 2>/dev/null; then
    install -m 0755 "$grepai_extract/grepai" "$HOME/.local/bin/grepai${GREPAI_EXT}"
    echo "  ✓ grepai ${GREPAI_LATEST} (${GREPAI_OS}/${GREPAI_ARCH})"
  else
    # Fallback: try go install with binary discovery
    if command -v go &>/dev/null; then
      if go install github.com/yoanbernabeu/grepai/cmd/grepai@latest 2>/dev/null; then
        # Ensure binary is in expected location
        GREPAI_BIN_PATH="$(go env GOBIN 2>/dev/null || true)"
        [ -z "${GREPAI_BIN_PATH}" ] && GREPAI_BIN_PATH="$(go env GOPATH 2>/dev/null)/bin"
        if [ -x "${GREPAI_BIN_PATH}/grepai${GREPAI_EXT}" ]; then
          cp -f "${GREPAI_BIN_PATH}/grepai${GREPAI_EXT}" "$HOME/.local/bin/grepai${GREPAI_EXT}"
          chmod +x "$HOME/.local/bin/grepai${GREPAI_EXT}"
        fi
        echo "  ✓ grepai (go install)"
      else
        echo "  ⚠ grepai install failed (optional)"
      fi
    else
      echo "  ⚠ grepai download failed (optional)"
    fi
  fi
  rm -f "$grepai_tmp"
  rm -rf "$grepai_extract"
fi

# ─────────────────────────────────────────────────────────────────────────────
# 8. Install status-line (official binary) - controlled by statusline option
# ─────────────────────────────────────────────────────────────────────────────
if [ "$INSTALL_STATUSLINE" != "true" ]; then
  echo "→ Skipping status-line installation (statusline=false)"
else
  if command -v status-line &>/dev/null; then
    echo "  ✓ status-line already installed"
  else
    echo "→ Installing status-line..."
    mkdir -p "$HOME/.local/bin"

    # Detect OS
    case "$(uname -s)" in
      Linux*) STATUS_OS="linux" ;;
      Darwin*) STATUS_OS="darwin" ;;
      MINGW* | MSYS* | CYGWIN*) STATUS_OS="windows" ;;
      *) STATUS_OS="linux" ;;
    esac

    # Detect architecture
    case "$(uname -m)" in
      x86_64 | amd64) STATUS_ARCH="amd64" ;;
      aarch64 | arm64) STATUS_ARCH="arm64" ;;
      *) STATUS_ARCH="amd64" ;;
    esac

    # Extension for Windows
    STATUS_EXT=""
    [ "$STATUS_OS" = "windows" ] && STATUS_EXT=".exe"

    # Download from official releases (with secure download)
    STATUS_URL="https://github.com/kodflow/status-line/releases/latest/download/status-line-${STATUS_OS}-${STATUS_ARCH}${STATUS_EXT}"
    status_tmp="$(mktemp)"
    if curl -fsL --retry 3 --retry-delay 1 --proto '=https' --tlsv1.2 "$STATUS_URL" -o "$status_tmp" 2>/dev/null; then
      install -m 0755 "$status_tmp" "$HOME/.local/bin/status-line${STATUS_EXT}"
      echo "  ✓ status-line (${STATUS_OS}/${STATUS_ARCH})"
    else
      echo "  ⚠ status-line download failed (optional)"
    fi
    rm -f "$status_tmp"

  fi
fi

# Add to PATH if needed
if [[ ":$PATH:" != *":$HOME/.local/bin:"* ]]; then
  # shellcheck disable=SC2016 # $HOME must be resolved at shell runtime, not now
  echo 'export PATH="$HOME/.local/bin:$PATH"' >>"$HOME/.bashrc" 2>/dev/null || true
  # shellcheck disable=SC2016 # $HOME must be resolved at shell runtime, not now
  echo 'export PATH="$HOME/.local/bin:$PATH"' >>"$HOME/.zshrc" 2>/dev/null || true
fi

# ─────────────────────────────────────────────────────────────────────────────
# Done
# ─────────────────────────────────────────────────────────────────────────────
echo ""
echo "═══════════════════════════════════════════"
echo "  ✓ Installation complete!"
echo ""
echo "  Available commands:"
echo "    /git     - Git workflow (commit, branch, PR)"
echo "    /search  - Documentation search"
echo "    /prompt  - Write better /plan descriptions"
echo ""
echo "  Native Claude 2.x features:"
echo "    EnterPlanMode - Built-in planning"
echo "    TaskCreate    - Task tracking with progress"
echo "    Task agents   - Parallelization"
echo ""
echo "  → Restart 'claude' to load commands"
echo "═══════════════════════════════════════════"
