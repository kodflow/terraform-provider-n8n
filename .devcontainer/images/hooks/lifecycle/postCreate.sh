#!/bin/bash
# shellcheck disable=SC1090,SC1091
# ============================================================================
# postCreate.sh - Runs ONCE after container is assigned to user
# ============================================================================
# This script runs once after the dev container is assigned to a user.
# Use it for: User-specific setup, environment variables, shell config.
# Has access to user-specific secrets and permissions.
#
# Uses run_step pattern: each step runs in an isolated subshell so that
# failures (e.g. unconfigured git email, missing GPG keys) never kill
# the entire script. The container always starts successfully.
# ============================================================================

set -u

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/../shared/utils.sh"

echo ""
echo -e "${CYAN}=========================================${NC}"
echo -e "${CYAN}   DevContainer Setup${NC}"
echo -e "${CYAN}=========================================${NC}"
echo ""

init_steps

# ============================================================================
# Step functions
# ============================================================================

# Prevents "dubious ownership" errors when container user differs from
# directory owner (common in Docker where /workspace may be owned by root)
step_git_safe_directory() {
  if ! git config --global --get-all safe.directory 2>/dev/null | grep -q "^/workspace$"; then
    git config --global --add safe.directory /workspace
    log_success "Git safe.directory configured for /workspace"
  else
    log_info "Git safe.directory already configured"
  fi
}

# Conditionally disable SSL verification (for corporate proxies/self-signed certs)
# Only applies when GIT_SSL_NO_VERIFY=1 is set in .env or environment
step_git_ssl_config() {
  if [ "${GIT_SSL_NO_VERIFY:-0}" = "1" ]; then
    git config --global http.sslVerify false
    log_success "Git SSL verification disabled (GIT_SSL_NO_VERIFY=1)"
  else
    log_info "Git SSL verification kept enabled (set GIT_SSL_NO_VERIFY=1 to disable)"
  fi
}

# GPG commit signing configuration
step_gpg_signing() {
  if [ ! -d "/home/vscode/.gnupg" ] || ! gpg --list-secret-keys --keyid-format LONG 2>/dev/null | grep -q "^sec"; then
    log_info "No GPG keys available - commit signing disabled"
    return 0
  fi

  # Get GIT_EMAIL from .env or git config
  local git_email=""
  if [ -f "/workspace/.env" ]; then
    git_email=$(grep -E "^GIT_EMAIL=" /workspace/.env 2>/dev/null | cut -d'=' -f2 | tr -d '"' || true)
  fi
  if [ -z "$git_email" ]; then
    git_email=$(git config --global user.email 2>/dev/null || true)
  fi

  local gpg_key=""
  if [ -n "$git_email" ]; then
    # Priority: Find GPG key matching the configured GIT_EMAIL
    gpg_key=$(gpg --list-secret-keys --keyid-format LONG 2>/dev/null \
      | grep -B1 "$git_email" 2>/dev/null \
      | grep -E "^sec" 2>/dev/null | head -1 | awk '{print $2}' | cut -d'/' -f2 || true)
  fi

  if [ -n "$gpg_key" ]; then
    git config --global user.signingkey "$gpg_key"
    git config --global commit.gpgsign true
    git config --global tag.forceSignAnnotated true
    git config --global gpg.program gpg
    log_success "Git GPG signing configured with key: $gpg_key (matching $git_email)"
  else
    # No matching key found - GPG signing will be configured via /git skill
    log_warning "No GPG key found for email '$git_email' - configure via /git skill"
  fi
}

# Create environment initialization script (~/.devcontainer-env.sh)
step_create_env_script() {
  log_info "Setting up environment variables and aliases..."

  cat >/home/vscode/.devcontainer-env.sh <<'ENVEOF'
# DevContainer Environment Initialization (v3 - lazy wrappers + cached completions)
# This file is sourced by ~/.zshrc and ~/.bashrc
#
# Architecture: Two-phase loading for fast shell startup
#   Phase 1 (always): PATH exports, env vars, fpath — fast, no subprocesses
#   Phase 2 (real terminal only): lazy wrappers, aliases, fast completions
#
# Why: VS Code's ptyHost spawns a shell to resolve env vars with a 10s timeout.
# Heavy init (eval, source <(...), nvm.sh) easily exceeds this on ARM64.
# Phase 1 gives VS Code the PATH/env it needs; Phase 2 only runs in terminals.
#
# v3 changes (from v2):
#   - Version managers (NVM, pyenv, rbenv, SDKMAN) use lazy wrappers instead of
#     eager init. Management commands load on first use; tool binaries (node,
#     python, ruby, java) work immediately via Phase 1 PATH/shims.
#   - Completions (kubectl, helm, docker, etc.) pre-cached to ~/.zsh_completions/
#     by postStart.sh and loaded via fpath (no more source <(...) subprocesses).

# ============================================================================
# Phase 1: Fast PATH and Environment Variables (no subprocesses)
# ============================================================================

# NVM (Node.js Version Manager)
export NVM_DIR="/usr/local/share/nvm"
export NVM_SYMLINK_CURRENT=true
# Add NVM current bin to PATH directly (no need to source heavy nvm.sh)
[ -d "$NVM_DIR/current/bin" ] && export PATH="$NVM_DIR/current/bin:$PATH"

# pyenv (Python Version Manager)
export PYENV_ROOT="/home/vscode/.cache/pyenv"
if [ -d "$PYENV_ROOT" ]; then
    export PATH="$PYENV_ROOT/shims:$PYENV_ROOT/bin:$PATH"
fi

# rbenv (Ruby Version Manager)
export RBENV_ROOT="/home/vscode/.cache/rbenv"
if [ -d "$RBENV_ROOT" ]; then
    export PATH="$RBENV_ROOT/shims:$RBENV_ROOT/bin:$PATH"
fi

# SDKMAN (Java/JVM SDK Manager)
export SDKMAN_DIR="/home/vscode/.cache/sdkman"
if [ -d "$SDKMAN_DIR/candidates" ]; then
    for _sdk_bin in "$SDKMAN_DIR"/candidates/*/current/bin; do
        [ -d "$_sdk_bin" ] && PATH="$_sdk_bin:$PATH"
    done
    unset _sdk_bin
fi

# Rust/Cargo
export CARGO_HOME="/home/vscode/.cache/cargo"
export RUSTUP_HOME="/home/vscode/.cache/rustup"
[ -d "$CARGO_HOME/bin" ] && export PATH="$CARGO_HOME/bin:$PATH"

# Go
export GOPATH="/home/vscode/.cache/go"
if [ -d "/usr/local/go" ]; then
    export GOROOT="/usr/local/go"
    export PATH="$GOROOT/bin:$GOPATH/bin:$PATH"
fi

# Flutter/Dart
export FLUTTER_ROOT="/home/vscode/.cache/flutter"
export PUB_CACHE="/home/vscode/.cache/pub-cache"
if [ -d "$FLUTTER_ROOT" ]; then
    export PATH="$FLUTTER_ROOT/bin:$PUB_CACHE/bin:$PATH"
fi

# Composer (PHP)
export COMPOSER_HOME="/home/vscode/.cache/composer"
export PATH="$COMPOSER_HOME/vendor/bin:$PATH"

# Mix (Elixir)
export MIX_HOME="/home/vscode/.cache/mix"
export PATH="$MIX_HOME/escripts:$PATH"

# npm global packages
export PATH="/home/vscode/.local/share/npm-global/bin:$PATH"

# pnpm
export PNPM_HOME="/home/vscode/.cache/pnpm"
export PATH="$PNPM_HOME:$PATH"

# Local bin
export PATH="/home/vscode/.local/bin:$PATH"

# vcpkg
export VCPKG_ROOT="/home/vscode/.cache/vcpkg"
export PATH="$VCPKG_ROOT:$PATH"

# Scala (SBT)
export SBT_HOME="/home/vscode/.cache/sbt"
[ -d "$SBT_HOME/bin" ] && export PATH="$SBT_HOME/bin:$PATH"

# .NET (C#, VB.NET)
export DOTNET_ROOT="/usr/share/dotnet"
[ -d "$DOTNET_ROOT" ] && export PATH="$DOTNET_ROOT:$HOME/.dotnet/tools:$PATH"

# R
export R_HOME="/usr/lib/R"

# Cached completions: pre-generated by postStart.sh, loaded via fpath
# Must be set before compinit (which runs inside Oh My Zsh)
if [ -d "$HOME/.zsh_completions" ]; then
    fpath=("$HOME/.zsh_completions" $fpath)
fi

# ============================================================================
# Phase 2: Interactive Terminal Features (lazy wrappers, aliases, fast completions)
# ============================================================================
# Skip when stdout is not a real terminal (e.g., VS Code env resolution).
# This is the key optimization: VS Code only needs PATH/env from Phase 1.
if [ ! -t 1 ]; then
    return 0 2>/dev/null || true
fi

# ----------------------------------------------------------------------------
# Lazy-load wrappers for version managers
# Phase 1 PATH already covers tool binaries (node, python, ruby, java) via
# symlinks and shims. These wrappers only load the full manager when the
# management command itself is first used (nvm, pyenv, rbenv, sdk).
# ----------------------------------------------------------------------------

# NVM: lazy-load on first 'nvm' call (~500ms saved per shell)
nvm() {
    unfunction nvm 2>/dev/null || unset -f nvm 2>/dev/null
    [ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"
    [ -s "$NVM_DIR/bash_completion" ] && \. "$NVM_DIR/bash_completion"
    nvm "$@"
}

# pyenv: lazy-load on first 'pyenv' call (~300ms saved per shell)
pyenv() {
    unfunction pyenv 2>/dev/null || unset -f pyenv 2>/dev/null
    if [ -d "$PYENV_ROOT" ]; then
        eval "$(command pyenv init -)" 2>/dev/null || true
        eval "$(command pyenv virtualenv-init -)" 2>/dev/null || true
    fi
    command pyenv "$@"
}

# rbenv: lazy-load on first 'rbenv' call (~150ms saved per shell)
rbenv() {
    unfunction rbenv 2>/dev/null || unset -f rbenv 2>/dev/null
    if [ -d "$RBENV_ROOT" ]; then
        eval "$(command rbenv init -)" 2>/dev/null || true
    fi
    command rbenv "$@"
}

# SDKMAN: lazy-load on first 'sdk' call (~400ms saved per shell)
sdk() {
    unfunction sdk 2>/dev/null || unset -f sdk 2>/dev/null
    [[ -s "$SDKMAN_DIR/bin/sdkman-init.sh" ]] && source "$SDKMAN_DIR/bin/sdkman-init.sh"
    sdk "$@"
}

# ----------------------------------------------------------------------------
# Aliases
# ----------------------------------------------------------------------------

# super-claude: runs claude with MCP config if available, otherwise without
super-claude() {
    local mcp_config="/workspace/mcp.json"

    # Check if jq is available for JSON validation
    if ! command -v jq &>/dev/null; then
        echo "Warning: jq not found, skipping MCP config validation" >&2
        # Still use mcp config if it looks like JSON (skip leading whitespace/newlines)
        if [ -s "$mcp_config" ] && LC_ALL=C tr -d ' \t\r\n' < "$mcp_config" 2>/dev/null | head -c 1 | grep -q '{'; then
            claude --dangerously-skip-permissions --mcp-config "$mcp_config" "$@"
        else
            claude --dangerously-skip-permissions "$@"
        fi
        return
    fi

    if [ -f "$mcp_config" ] && jq empty "$mcp_config" 2>/dev/null; then
        claude --dangerously-skip-permissions --mcp-config "$mcp_config" "$@"
    else
        claude --dangerously-skip-permissions "$@"
    fi
}

# ----------------------------------------------------------------------------
# Fast completions (native complete -C, ~1ms each — no subprocess overhead)
# Heavier completions (kubectl, helm, docker, etc.) are pre-cached to
# ~/.zsh_completions/ by postStart.sh and loaded via fpath above.
# ----------------------------------------------------------------------------

# HashiCorp tools (native binary completion)
if command -v terraform &> /dev/null; then
    complete -o nospace -C "$(which terraform)" terraform 2>/dev/null || true
fi
if command -v vault &> /dev/null; then
    complete -o nospace -C "$(which vault)" vault 2>/dev/null || true
fi
if command -v consul &> /dev/null; then
    complete -o nospace -C "$(which consul)" consul 2>/dev/null || true
fi
if command -v nomad &> /dev/null; then
    complete -o nospace -C "$(which nomad)" nomad 2>/dev/null || true
fi
if command -v packer &> /dev/null; then
    complete -o nospace -C "$(which packer)" packer 2>/dev/null || true
fi

# AWS CLI (native binary completion)
if command -v aws_completer &> /dev/null; then
    complete -C aws_completer aws 2>/dev/null || true
fi

# Google Cloud SDK (static file, fast)
if [ -f "/usr/share/google-cloud-sdk/completion.zsh.inc" ]; then
    source "/usr/share/google-cloud-sdk/completion.zsh.inc" 2>/dev/null || true
fi
ENVEOF

  log_success "Environment script created at ~/.devcontainer-env.sh"
}

# Mark container as initialized
step_mark_initialized() {
  touch /home/vscode/.devcontainer-initialized
  log_success "DevContainer marked as initialized"
}

# ============================================================================
# Execution (always runs git steps, skips env if already initialized)
# ============================================================================

# Git steps run every time (safe directory, SSL, GPG)
run_step "Git safe directory" step_git_safe_directory
run_step "Git SSL configuration" step_git_ssl_config
run_step "GPG signing" step_gpg_signing

# Note: Tools (status-line, ktn-linter) are now baked into the Docker image
# No longer need to update on each rebuild

# Check if already initialized (but only if env file also exists)
# If ~/.devcontainer-env.sh is missing, we must recreate it even if marker exists
if [ -f /home/vscode/.devcontainer-initialized ] && [ -f /home/vscode/.devcontainer-env.sh ]; then
  log_success "DevContainer already initialized"
  echo ""
  exit 0
fi

run_step "Environment script" step_create_env_script
run_step "Mark initialized" step_mark_initialized

print_step_summary "postCreate"

exit 0
