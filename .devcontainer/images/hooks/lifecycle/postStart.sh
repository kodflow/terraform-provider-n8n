#!/bin/bash
# shellcheck disable=SC1090,SC1091
# ============================================================================
# postStart.sh - Runs EVERY TIME the container starts
# ============================================================================
# This script runs after postCreateCommand and before postAttachCommand.
# Runs each time the container is successfully started.
# Use it for: MCP setup, services startup, recurring initialization.
#
# Uses run_step pattern: each step runs in an isolated subshell so that
# failures never block the DevContainer lifecycle.
# ============================================================================

set -u

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/../shared/utils.sh"

log_info "postStart: Container starting..."

init_steps

# ============================================================================
# Helpers
# ============================================================================

# Resolve 1Password vault ID
# Priority: OP_VAULT_ID env var > dynamic lookup of "CI" vault
resolve_vault_id() {
  # Priority 1: OP_VAULT_ID env var (validated)
  if [ -n "${OP_VAULT_ID:-}" ]; then
    local sanitized
    sanitized=$(echo "$OP_VAULT_ID" | tr -d '[:space:]')
    # Reject empty after trim or values starting with '-' (option injection)
    if [ -z "$sanitized" ] || [[ "$sanitized" == -* ]]; then
      log_warning "OP_VAULT_ID contains invalid value, ignoring" >&2
    else
      echo "$sanitized"
      return 0
    fi
  fi
  # Priority 2: dynamic resolution of vault named "CI"
  if command -v op &>/dev/null && [ -n "${OP_SERVICE_ACCOUNT_TOKEN:-}" ]; then
    if ! command -v jq &>/dev/null; then
      log_info "jq not available; set OP_VAULT_ID to skip auto-lookup" >&2
      echo ""
      return 0
    fi
    local vault_id
    vault_id=$(op vault get "CI" --format=json 2>/dev/null | jq -r '.id // empty' 2>/dev/null || echo "")
    if [ -n "$vault_id" ]; then
      echo "$vault_id"
      return 0
    fi
  fi
  echo ""
}

# ============================================================================
# Step functions
# ============================================================================

# Restore Claude commands/scripts from image defaults OR host installation
# Priority:
#   1. Host installation ($HOME/.claude/ with .template-version)
#   2. Image defaults (/etc/claude-defaults/)
step_restore_claude_config() {
  local CLAUDE_DEFAULTS="/etc/claude-defaults"

  # Check if user has a host-installed Claude configuration
  # (installed via install.sh - indicated by .template-version file)
  if [ -f "$HOME/.claude/.template-version" ]; then
    log_info "Detected host-installed Claude configuration (via install.sh)"
    log_info "Skipping image defaults restore - using host installation"
    log_success "Claude configuration from host (priority)"
    return 0
  fi

  if [ ! -d "$CLAUDE_DEFAULTS" ]; then
    log_warning "No Claude configuration found (neither host nor image defaults)"
    return 0
  fi

  log_info "Restoring Claude configuration from image defaults..."

  # Ensure base directory exists
  mkdir -p "$HOME/.claude"

  # CLEAN commands, scripts, agents and docs to avoid legacy pollution
  # Only these directories are managed by the image - sessions/plans are user data
  rm -rf "$HOME/.claude/commands" "$HOME/.claude/scripts" "$HOME/.claude/agents" "$HOME/.claude/docs"

  # Restore commands (fresh copy from image)
  if [ -d "$CLAUDE_DEFAULTS/commands" ]; then
    mkdir -p "$HOME/.claude/commands"
    cp -r "$CLAUDE_DEFAULTS/commands/"* "$HOME/.claude/commands/" 2>/dev/null || true
  fi

  # Restore scripts (fresh copy from image)
  if [ -d "$CLAUDE_DEFAULTS/scripts" ]; then
    mkdir -p "$HOME/.claude/scripts"
    cp -r "$CLAUDE_DEFAULTS/scripts/"* "$HOME/.claude/scripts/" 2>/dev/null || true
    chmod -R 755 "$HOME/.claude/scripts/"
  fi

  # Restore agents (fresh copy from image)
  if [ -d "$CLAUDE_DEFAULTS/agents" ]; then
    mkdir -p "$HOME/.claude/agents"
    cp -r "$CLAUDE_DEFAULTS/agents/"* "$HOME/.claude/agents/" 2>/dev/null || true
    chmod -R 755 "$HOME/.claude/agents/"
  fi

  # Restore docs (Design Patterns Knowledge Base - fresh copy from image)
  if [ -d "$CLAUDE_DEFAULTS/docs" ]; then
    mkdir -p "$HOME/.claude/docs"
    cp -r "$CLAUDE_DEFAULTS/docs/"* "$HOME/.claude/docs/" 2>/dev/null || true
  fi

  # Restore templates (Documentation and C4 templates - fresh copy from image)
  if [ -d "$CLAUDE_DEFAULTS/templates" ]; then
    mkdir -p "$HOME/.claude/templates"
    cp -r "$CLAUDE_DEFAULTS/templates/"* "$HOME/.claude/templates/" 2>/dev/null || true
  fi

  # Restore settings.json only if it does not exist (user customizations preserved)
  if [ -f "$CLAUDE_DEFAULTS/settings.json" ] && [ ! -f "$HOME/.claude/settings.json" ]; then
    cp "$CLAUDE_DEFAULTS/settings.json" "$HOME/.claude/settings.json"
  fi

  log_success "Claude configuration restored from image defaults"
}

# Ensure Claude directories exist (volume mount point)
step_init_claude_dirs() {
  mkdir -p "$HOME/.claude/sessions" "$HOME/.claude/plans" "$HOME/.claude/contexts"
  log_success "Claude directories initialized"
}

# Ensure shell environment is properly configured (repair mechanism)
# postCreate.sh creates ~/.devcontainer-env.sh with super-claude() and other
# shell functions. If the source line is missing from RC files (stale image,
# volume issue, or .bashrc never configured), inject it here on every start.
step_shell_env_repair() {
  local DC_SOURCE_LINE='[[ -f ~/.devcontainer-env.sh ]] && source ~/.devcontainer-env.sh'

  for rc_file in "$HOME/.zshrc" "$HOME/.bashrc"; do
    if [ -f "$rc_file" ] && ! grep -q "devcontainer-env.sh" "$rc_file" 2>/dev/null; then
      printf '\n%s\n' "$DC_SOURCE_LINE" >>"$rc_file"
      log_info "Added devcontainer-env source to $(basename "$rc_file")"
    fi
  done

  # Upgrade devcontainer-env.sh: v1→v2 (two-phase) or v2→v3 (lazy wrappers)
  # v1: heavy init unconditionally (caused VS Code ptyHost timeout)
  # v2: two-phase (Phase 1 fast, Phase 2 terminal-only) — still eager init
  # v3: lazy wrappers + cached completions (no eager init, no source <(...))
  local DC_ENV="$HOME/.devcontainer-env.sh"
  local need_regen=false

  # v1→v3: no two-phase marker at all
  if [ -f "$DC_ENV" ] && ! grep -q "two-phase\|lazy wrappers" "$DC_ENV" 2>/dev/null; then
    log_info "Upgrading devcontainer-env.sh from v1 to v3..."
    need_regen=true
  fi

  # v2→v3: has two-phase but still uses source <(cmd completion)
  if [ -f "$DC_ENV" ] && grep -q "two-phase" "$DC_ENV" 2>/dev/null && ! grep -q "lazy wrappers" "$DC_ENV" 2>/dev/null; then
    log_info "Upgrading devcontainer-env.sh from v2 to v3 (lazy wrappers)..."
    need_regen=true
  fi

  if [ "$need_regen" = true ]; then
    if [ -f "$HOME/.devcontainer-initialized" ]; then
      rm -f "$HOME/.devcontainer-initialized"
      if [ -f "${WORKSPACE_FOLDER:-/workspace}/.devcontainer/hooks/lifecycle/postCreate.sh" ]; then
        bash "${WORKSPACE_FOLDER:-/workspace}/.devcontainer/hooks/lifecycle/postCreate.sh" 2>/dev/null || true
      fi
      touch "$HOME/.devcontainer-initialized"
    fi
    log_success "devcontainer-env.sh upgraded to v3 (lazy wrappers + cached completions)"
  fi

  # Remove duplicate NVM from .zshrc (added by nodejs feature, already in env.sh)
  if [ -f "$HOME/.zshrc" ] && grep -c "NVM_DIR" "$HOME/.zshrc" 2>/dev/null | grep -q "^[2-9]"; then
    sed -i '/^# NVM (Node Version Manager)$/,/^\[ -s "\$NVM_DIR\/nvm\.sh" \]/d' "$HOME/.zshrc" 2>/dev/null || true
    sed -i '/^export NVM_SYMLINK_CURRENT=true$/d' "$HOME/.zshrc" 2>/dev/null || true
    log_info "Removed duplicate NVM from .zshrc (handled by devcontainer-env.sh)"
  fi

  # Remove duplicate pyenv init from .zshrc (added by python feature, handled by lazy wrapper)
  if [ -f "$HOME/.zshrc" ] && grep -q "pyenv init" "$HOME/.zshrc" 2>/dev/null; then
    sed -i '/^# Pyenv initialization$/,/^eval "\$(pyenv init -)"/d' "$HOME/.zshrc" 2>/dev/null || true
    sed -i '/^eval "\$(pyenv virtualenv-init -)"/d' "$HOME/.zshrc" 2>/dev/null || true
    sed -i '/^export PYENV_ROOT=.*$/d' "$HOME/.zshrc" 2>/dev/null || true
    log_info "Removed duplicate pyenv init from .zshrc (handled by devcontainer-env.sh)"
  fi

  # Remove duplicate rbenv init from .zshrc (added by ruby feature, handled by lazy wrapper)
  if [ -f "$HOME/.zshrc" ] && grep -q "rbenv init" "$HOME/.zshrc" 2>/dev/null; then
    sed -i '/^# Rbenv initialization$/,/^eval "\$(rbenv init -)"/d' "$HOME/.zshrc" 2>/dev/null || true
    sed -i '/^export RBENV_ROOT=.*$/d' "$HOME/.zshrc" 2>/dev/null || true
    log_info "Removed duplicate rbenv init from .zshrc (handled by devcontainer-env.sh)"
  fi

  # Remove duplicate SDKMAN init from .zshrc (injected by get.sdkman.io installer)
  if [ -f "$HOME/.zshrc" ] && grep -q "sdkman-init.sh" "$HOME/.zshrc" 2>/dev/null; then
    sed -i '/^#THIS MUST BE AT THE END OF THE FILE/,/sdkman-init\.sh/d' "$HOME/.zshrc" 2>/dev/null || true
    log_info "Removed duplicate SDKMAN init from .zshrc (handled by devcontainer-env.sh)"
  fi

  # Ensure zsh is the default login shell
  if [ "$(getent passwd "$(whoami)" | cut -d: -f7)" != "/bin/zsh" ]; then
    sudo chsh -s /bin/zsh "$(whoami)" 2>/dev/null || true
    log_success "Default shell set to zsh"
  fi
}

# Pre-generate zsh completions to disk cache (~/.zsh_completions/)
# Replaces expensive 'source <(cmd completion zsh)' calls at shell startup.
# Each tool's completion is generated once per container start; stale files
# regenerated when the tool binary is newer than the cache (mtime-based).
step_cache_completions() {
  local COMP_DIR="$HOME/.zsh_completions"
  mkdir -p "$COMP_DIR"

  # Helper: regenerate completion file if tool binary is newer than cache
  cache_completion() {
    local tool_cmd="$1"
    local out_file="$2"
    shift 2
    local gen_cmd=("$@")

    command -v "$tool_cmd" &>/dev/null || return 0

    local tool_bin
    tool_bin=$(command -v "$tool_cmd" 2>/dev/null)

    # Regenerate if: file missing, tool newer than cache, or file empty
    if [ ! -f "$out_file" ] || [ "$tool_bin" -nt "$out_file" ] || [ ! -s "$out_file" ]; then
      if "${gen_cmd[@]}" >"$out_file" 2>/dev/null && [ -s "$out_file" ]; then
        log_info "Cached completion: $(basename "$out_file")"
      else
        rm -f "$out_file"
      fi
    fi
  }

  cache_completion kubectl "$COMP_DIR/_kubectl" kubectl completion zsh
  cache_completion helm "$COMP_DIR/_helm" helm completion zsh
  cache_completion docker "$COMP_DIR/_docker" docker completion zsh
  cache_completion gh "$COMP_DIR/_gh" gh completion -s zsh
  cache_completion rustup "$COMP_DIR/_rustup" rustup completions zsh
  cache_completion npm "$COMP_DIR/_npm" npm completion
  cache_completion pnpm "$COMP_DIR/_pnpm" pnpm completion zsh

  # Cargo completion via rustup (separate command)
  if command -v rustup &>/dev/null; then
    local rustup_bin
    rustup_bin=$(command -v rustup 2>/dev/null)
    if [ ! -f "$COMP_DIR/_cargo" ] || [ "$rustup_bin" -nt "$COMP_DIR/_cargo" ] || [ ! -s "$COMP_DIR/_cargo" ]; then
      if rustup completions zsh cargo >"$COMP_DIR/_cargo" 2>/dev/null && [ -s "$COMP_DIR/_cargo" ]; then
        log_info "Cached completion: _cargo"
      else
        rm -f "$COMP_DIR/_cargo"
      fi
    fi
  fi

  # Go completion: copy static file (no subprocess needed)
  if command -v go &>/dev/null; then
    local goroot
    goroot=$(go env GOROOT 2>/dev/null)
    if [ -f "$goroot/misc/zsh/go" ]; then
      if [ ! -f "$COMP_DIR/_go" ] || [ "$goroot/misc/zsh/go" -nt "$COMP_DIR/_go" ]; then
        cp "$goroot/misc/zsh/go" "$COMP_DIR/_go" 2>/dev/null && log_info "Cached completion: _go"
      fi
    fi
  fi

  local count
  count=$(find "$COMP_DIR" -name '_*' -type f 2>/dev/null | wc -l)
  log_success "ZSH completions cached ($count files in $COMP_DIR)"
}

# Generate dynamic p10k right-prompt segment list based on installed tools.
# Writes ~/.p10k-segments.zsh sourced by .p10k.zsh to override the static
# 47-segment RIGHT_PROMPT_ELEMENTS with only segments for present tools.
step_generate_p10k_segments() {
  local SEGMENTS_FILE="$HOME/.p10k-segments.zsh"
  local segments=()

  # Always-on utility segments
  segments+=(status command_execution_time background_jobs)

  # Python
  if [ -d "$HOME/.cache/pyenv" ] || command -v python3 &>/dev/null; then
    segments+=(virtualenv pyenv)
  fi

  # Node.js
  if [ -d "/usr/local/share/nvm" ] || command -v node &>/dev/null; then
    segments+=(nvm node_version)
  fi

  # Go
  command -v go &>/dev/null && segments+=(go_version)

  # Rust
  command -v rustc &>/dev/null && segments+=(rust_version)

  # .NET
  command -v dotnet &>/dev/null && segments+=(dotnet_version)

  # PHP
  command -v php &>/dev/null && segments+=(php_version)

  # Java/JVM
  if [ -d "$HOME/.cache/sdkman/candidates" ] || command -v java &>/dev/null; then
    segments+=(java_version)
  fi

  # Ruby
  if [ -d "$HOME/.cache/rbenv" ] || command -v ruby &>/dev/null; then
    segments+=(rbenv)
  fi

  # Flutter/Dart
  if [ -d "$HOME/.cache/flutter" ] || command -v flutter &>/dev/null; then
    segments+=(fvm)
  fi

  # Lua
  command -v luaenv &>/dev/null && segments+=(luaenv)

  # Scala
  command -v scalaenv &>/dev/null && segments+=(scalaenv)

  # Perl
  if command -v plenv &>/dev/null || command -v perlbrew &>/dev/null; then
    segments+=(plenv perlbrew)
  fi

  # Kubernetes (SHOW_ON_COMMAND guarded, lightweight)
  command -v kubectl &>/dev/null && segments+=(kubecontext)

  # Terraform (SHOW_ON_COMMAND guarded)
  command -v terraform &>/dev/null && segments+=(terraform)

  # Cloud CLIs (SHOW_ON_COMMAND guarded)
  command -v aws &>/dev/null && segments+=(aws)
  command -v az &>/dev/null && segments+=(azure)
  command -v gcloud &>/dev/null && segments+=(gcloud)

  # Always-on context and time
  segments+=(context time)

  # Line 2 separator
  segments+=(newline)

  # Write the override file (no console output — instant prompt constraint)
  {
    printf '# Auto-generated by postStart.sh — do not edit manually\n'
    printf '# Regenerated on each container start based on installed tools\n'
    printf 'typeset -g POWERLEVEL9K_RIGHT_PROMPT_ELEMENTS=(\n'
    local seg
    for seg in "${segments[@]}"; do
      printf '    %s\n' "$seg"
    done
    printf ')\n'
  } >"$SEGMENTS_FILE"

  local count=${#segments[@]}
  log_success "p10k segments: $count active → ~/.p10k-segments.zsh"
}

# Reload .env file to get updated tokens
step_reload_env() {
  local ENV_FILE="${WORKSPACE_FOLDER:-/workspace}/.devcontainer/.env"
  if [ -f "$ENV_FILE" ]; then
    log_info "Reloading environment from .env..."
    set -a
    source "$ENV_FILE"
    set +a
  fi
}

# Fix 1Password CLI config directory permissions
# Docker named volumes create directories with root ownership.
# 1Password CLI requires: ownership by current user + permissions 700.
step_1password_permissions() {
  local OP_CONFIG_DIRS=("$HOME/.config/op" "$HOME/.op")

  for OP_DIR in "${OP_CONFIG_DIRS[@]}"; do
    if [ -d "$OP_DIR" ]; then
      # Fix ownership if not current user
      if [ "$(stat -c '%U' "$OP_DIR" 2>/dev/null)" != "$(whoami)" ]; then
        log_info "Fixing ownership of $OP_DIR..."
        sudo chown -R "$(whoami):$(whoami)" "$OP_DIR"
      fi
      # Ensure correct permissions (700 = owner only)
      chmod 700 "$OP_DIR"
    fi
  done
  log_success "1Password config directories configured"
}

# Fix npm cache permissions
# Docker named volumes create directories with root ownership.
# npm requires write access to its cache for npx/MCP servers to work.
step_npm_cache_permissions() {
  local NPM_CACHE_DIR="$HOME/.cache/npm"

  if [ -d "$NPM_CACHE_DIR" ]; then
    # Fix ownership if not current user
    if [ "$(stat -c '%U' "$NPM_CACHE_DIR" 2>/dev/null)" != "$(whoami)" ]; then
      log_info "Fixing ownership of npm cache..."
      sudo chown -R "$(whoami):$(whoami)" "$NPM_CACHE_DIR"
    fi
  fi
  log_success "npm cache configured"
}

# MCP configuration setup (inject secrets into template)
step_mcp_configuration() {
  # 1Password vault ID (env var or dynamic "CI" vault lookup)
  local VAULT_ID
  VAULT_ID=$(resolve_vault_id)
  local MCP_TPL="/etc/mcp/mcp.json.tpl"
  local MCP_OUTPUT="${WORKSPACE_FOLDER:-/workspace}/mcp.json"

  # Skip regeneration if mcp.json exists and template hasn't changed
  if [ -f "$MCP_OUTPUT" ] && [ -f "$MCP_TPL" ]; then
    local tpl_hash output_marker="/tmp/.mcp-tpl-hash"
    tpl_hash=$(md5sum "$MCP_TPL" 2>/dev/null | cut -d" " -f1)
    if [ -f "$output_marker" ] && [ "$(cat "$output_marker" 2>/dev/null)" = "$tpl_hash" ]; then
      log_info "mcp.json already up-to-date (template unchanged), skipping"
      return 0
    fi
  fi

  # Helper function to get 1Password field (tries multiple field names)
  get_1password_field() {
    local item="$1"
    local vault="$2"
    local fields=("credential" "password" "identifiant" "mot de passe")

    for field in "${fields[@]}"; do
      local value
      value=$(op item get "$item" --vault "$vault" --fields "$field" --reveal 2>/dev/null || echo "")
      if [ -n "$value" ]; then
        echo "$value"
        return 0
      fi
    done
    echo ""
  }

  # Initialize tokens from environment variables (fallback)
  local CODACY_TOKEN="${CODACY_API_TOKEN:-}"
  local GITHUB_TOKEN="${GITHUB_API_TOKEN:-}"

  # Try 1Password if OP_SERVICE_ACCOUNT_TOKEN is defined
  if [ -n "${OP_SERVICE_ACCOUNT_TOKEN:-}" ] && command -v op &>/dev/null; then
    log_info "Retrieving secrets from 1Password..."
    if [ -z "$VAULT_ID" ]; then
      log_warning "Vault ID could not be resolved (OP_VAULT_ID unset, 'CI' vault not found, or jq missing), skipping 1Password secrets"
    else

      local OP_CODACY
      OP_CODACY=$(get_1password_field "mcp-codacy" "$VAULT_ID")
      local OP_GITHUB
      OP_GITHUB=$(get_1password_field "mcp-github" "$VAULT_ID")

      [ -n "$OP_CODACY" ] && CODACY_TOKEN="$OP_CODACY"
      [ -n "$OP_GITHUB" ] && GITHUB_TOKEN="$OP_GITHUB"
    fi
  fi

  # Show status of tokens (INFO for optional, WARNING for essential)
  [ -z "$CODACY_TOKEN" ] && log_info "Codacy token not configured (optional)"
  [ -z "$GITHUB_TOKEN" ] && log_warning "GitHub token not available"

  # Helper: escape special chars for sed replacement
  escape_for_sed() {
    LC_ALL=C printf '%s' "$1" | tr -d '\n\r' | sed -e 's/[&/|\\]/\\&/g'
  }

  # Security: refuse to write secrets through symlinks or unsafe directories
  local MCP_DIR
  MCP_DIR=$(dirname -- "$MCP_OUTPUT")
  if [ ! -d "$MCP_DIR" ] || [ -L "$MCP_DIR" ]; then
    log_error "Refusing to write mcp.json: unsafe parent directory ($MCP_DIR)"
    return 0
  elif [ -e "$MCP_OUTPUT" ] && { [ -L "$MCP_OUTPUT" ] || [ ! -f "$MCP_OUTPUT" ]; }; then
    log_error "Refusing to write mcp.json: not a regular file ($MCP_OUTPUT)"
    return 0
  fi

  # Migrate legacy .mcp.json to mcp.json (renamed in v2)
  if [ -f "${WORKSPACE_FOLDER:-/workspace}/.mcp.json" ] && [ ! -e "$MCP_OUTPUT" ]; then
    log_info "Migrating legacy .mcp.json to mcp.json..."

    if ! command -v jq >/dev/null 2>&1; then
      log_warning "jq not found; migrating without JSON validation"
      if cp "${WORKSPACE_FOLDER:-/workspace}/.mcp.json" "$MCP_OUTPUT"; then
        chown "$(id -u):$(id -g)" "$MCP_OUTPUT" 2>/dev/null || true
        chmod 600 "$MCP_OUTPUT"
        rm -f "${WORKSPACE_FOLDER:-/workspace}/.mcp.json" || log_warning "Could not remove legacy .mcp.json (permissions?)"
        log_success "Migration complete: .mcp.json → mcp.json"
      else
        log_error "Migration failed: unable to copy legacy file"
      fi
    else
      local MCP_MIG_TMP
      MCP_MIG_TMP=$(mktemp "${MCP_OUTPUT}.migrate.XXXXXX") || {
        log_error "Migration failed: unable to create temp file"
        MCP_MIG_TMP=""
      }
      if [ -n "$MCP_MIG_TMP" ] && cp "${WORKSPACE_FOLDER:-/workspace}/.mcp.json" "$MCP_MIG_TMP"; then
        if jq empty "$MCP_MIG_TMP" 2>/dev/null; then
          mv "$MCP_MIG_TMP" "$MCP_OUTPUT"
          chown "$(id -u):$(id -g)" "$MCP_OUTPUT" 2>/dev/null || true
          chmod 600 "$MCP_OUTPUT"
          rm -f "${WORKSPACE_FOLDER:-/workspace}/.mcp.json" || log_warning "Could not remove legacy .mcp.json (permissions?)"
          log_success "Migration complete: .mcp.json → mcp.json"
        else
          log_error "Legacy .mcp.json is invalid JSON; keeping legacy file"
          rm -f "$MCP_MIG_TMP"
        fi
      elif [ -n "$MCP_MIG_TMP" ]; then
        log_error "Migration failed"
        rm -f "$MCP_MIG_TMP"
      fi
    fi
  fi

  # Generate mcp.json from template (baked in Docker image)
  # ALWAYS regenerate from template to ensure latest MCP config is applied
  if [ -f "$MCP_TPL" ]; then
    if [ -z "$CODACY_TOKEN" ] && [ -z "$GITHUB_TOKEN" ]; then
      log_warning "No tokens available, creating minimal mcp.json"
      printf '%s\n' '{"mcpServers":{}}' >"$MCP_OUTPUT"
      chown "$(id -u):$(id -g)" "$MCP_OUTPUT" 2>/dev/null || true
      chmod 600 "$MCP_OUTPUT"
      log_info "Created minimal mcp.json (optional MCPs will be added below)"
    else
      generate_mcp_from_template() {
        local escaped_codacy escaped_github mcp_tmp
        escaped_codacy=$(escape_for_sed "${CODACY_TOKEN}")
        escaped_github=$(escape_for_sed "${GITHUB_TOKEN}")

        mcp_tmp=$(mktemp "${MCP_OUTPUT}.tmp.XXXXXX") || {
          log_error "Failed to create temp file for mcp.json generation"
          return 0
        }

        trap 'rm -f "${mcp_tmp:-}" 2>/dev/null || true' RETURN

        if ! sed -e "s|{{CODACY_TOKEN}}|${escaped_codacy}|g" \
          -e "s|{{GITHUB_TOKEN}}|${escaped_github}|g" \
          "$MCP_TPL" >"$mcp_tmp"; then
          log_error "Failed to render mcp.json template"
          return 0
        fi

        if jq empty "$mcp_tmp" 2>/dev/null; then
          mv "$mcp_tmp" "$MCP_OUTPUT"
          chown "$(id -u):$(id -g)" "$MCP_OUTPUT" 2>/dev/null || true
          chmod 600 "$MCP_OUTPUT"
          log_success "mcp.json generated successfully"
          # Save template hash to skip regeneration on next start
          md5sum "$MCP_TPL" 2>/dev/null | cut -d" " -f1 >/tmp/.mcp-tpl-hash 2>/dev/null || true
        else
          log_error "Generated mcp.json is invalid JSON, keeping original"
        fi
      }
      log_info "Regenerating mcp.json from template (forced)..."
      generate_mcp_from_template
    fi

    # Add optional MCPs based on installed features
    add_optional_mcp() {
      local name="$1"
      local binary="$2"
      local output="$3"

      [ -f "$output" ] || return 0

      if ! command -v jq >/dev/null 2>&1; then
        log_warning "Skipping $name MCP injection (jq not found)"
        return 0
      fi

      if [ -x "$binary" ]; then
        log_info "Adding $name MCP (binary found at $binary)"
        local tmp_file
        tmp_file=$(mktemp "${output}.tmp.XXXXXX") || {
          log_warning "Failed to add $name MCP (unable to create temp file)"
          return 0
        }
        if jq --arg name "$name" --arg bin "$binary" \
          '.mcpServers = (.mcpServers // {}) | .mcpServers[$name] = {"command": $bin, "args": [], "env": {}}' \
          "$output" >"$tmp_file" && jq empty "$tmp_file" 2>/dev/null; then
          mv "$tmp_file" "$output"
          chown "$(id -u):$(id -g)" "$output" 2>/dev/null || true
          chmod 600 "$output" 2>/dev/null || true
        else
          log_warning "Failed to add $name MCP, keeping original"
          rm -f "$tmp_file"
        fi
      else
        log_info "Skipping $name MCP (binary not found)"
      fi
    }

    add_optional_mcp "rust-analyzer" "$HOME/.cache/cargo/bin/rust-analyzer-mcp" "$MCP_OUTPUT"
  else
    log_warning "MCP template not found at $MCP_TPL"
  fi
}

# Taskmaster AI configuration
# Generates .taskmaster/config.json with MCP provider (uses Claude Code session)
# Optional Ollama fallback for streaming if available on host
step_taskmaster_config() {
  local TASKMASTER_DIR="${WORKSPACE_FOLDER:-/workspace}/.taskmaster"
  local TASKMASTER_CONFIG="$TASKMASTER_DIR/config.json"
  local tm_tmp=""

  if ! command -v jq >/dev/null 2>&1; then
    log_warning "Taskmaster config skipped: jq not available"
    return 0
  fi

  # Idempotent: skip only when existing config is valid JSON
  if [ -f "$TASKMASTER_CONFIG" ] && jq empty "$TASKMASTER_CONFIG" >/dev/null 2>&1; then
    log_info "Taskmaster config exists, skipping"
    return 0
  fi

  # Safety: refuse symlinked or non-directory target
  if [ -e "$TASKMASTER_DIR" ] && { [ -L "$TASKMASTER_DIR" ] || [ ! -d "$TASKMASTER_DIR" ]; }; then
    log_warning "Taskmaster config skipped: unsafe directory target ($TASKMASTER_DIR)"
    return 0
  fi
  mkdir -p "$TASKMASTER_DIR" 2>/dev/null || return 0
  if [ -e "$TASKMASTER_CONFIG" ] && { [ -L "$TASKMASTER_CONFIG" ] || [ ! -f "$TASKMASTER_CONFIG" ]; }; then
    log_warning "Taskmaster config skipped: unsafe file target ($TASKMASTER_CONFIG)"
    return 0
  fi
  tm_tmp=$(mktemp "${TASKMASTER_CONFIG}.tmp.XXXXXX") || return 0

  # Detect Ollama for optional fallback streaming
  local has_ollama=false
  local ollama_url="${OLLAMA_HOST:-host.docker.internal:11434}"
  local model="${OLLAMA_CHAT_MODEL:-llama3.2}"
  if curl -sf --connect-timeout 2 "http://$ollama_url/api/tags" >/dev/null 2>&1; then
    has_ollama=true
  fi

  if [ "$has_ollama" = true ]; then
    jq -n --arg url "http://$ollama_url/api" --arg model "$model" --arg root "${WORKSPACE_FOLDER:-/workspace}" '{
            models: {
                main:     { provider: "mcp", modelId: "claude-code" },
                research: { provider: "mcp", modelId: "claude-code" },
                fallback: { provider: "ollama", modelId: $model, baseURL: $url }
            },
            global: { projectRoot: $root, logLevel: "info" }
        }' >"$tm_tmp"
  else
    jq -n --arg root "${WORKSPACE_FOLDER:-/workspace}" '{
            models: {
                main:     { provider: "mcp", modelId: "claude-code" },
                research: { provider: "mcp", modelId: "claude-code" }
            },
            global: { projectRoot: $root, logLevel: "info" }
        }' >"$tm_tmp"
  fi

  if jq empty "$tm_tmp" >/dev/null 2>&1; then
    if mv "$tm_tmp" "$TASKMASTER_CONFIG"; then
      chmod 600 "$TASKMASTER_CONFIG" 2>/dev/null || true
      log_success "Taskmaster config generated"
    else
      rm -f "$tm_tmp"
      log_warning "Taskmaster config generation failed (install step failed)"
    fi
  else
    rm -f "$tm_tmp"
    log_warning "Taskmaster config generation failed (invalid JSON)"
  fi
}

# CodeRabbit CLI authentication (inject token from 1Password)
# Retrieves the API token from 1Password and writes ~/.coderabbit/auth.json
# This enables `coderabbit review` and `cr review` without manual login
# Graceful degradation: skips silently if 1Password/token/vault unavailable
step_coderabbit_auth() {
  local VAULT_ID
  VAULT_ID=$(resolve_vault_id)
  local CR_AUTH_DIR="$HOME/.coderabbit"
  local CR_AUTH_FILE="$CR_AUTH_DIR/auth.json"

  # Skip if CodeRabbit CLI not installed
  if ! command -v coderabbit &>/dev/null; then
    log_info "CodeRabbit CLI not installed, skipping auth"
    return 0
  fi

  # Skip if already authenticated (check both OAuth and API key formats)
  if [ -f "$CR_AUTH_FILE" ]; then
    if command -v jq >/dev/null 2>&1; then
      local existing_token
      existing_token=$(jq -r '.accessToken // .apiKey // ""' "$CR_AUTH_FILE" 2>/dev/null || echo "")
      if [ -n "$existing_token" ]; then
        log_success "CodeRabbit already authenticated"
        return 0
      fi
    else
      # No jq but auth file exists - assume valid (don't overwrite)
      log_info "CodeRabbit auth file exists (cannot validate without jq)"
      return 0
    fi
  fi

  # Guard: 1Password CLI must be available
  if ! command -v op &>/dev/null; then
    log_info "CodeRabbit: op CLI not installed, skipping auto-auth"
    return 0
  fi

  # Guard: OP_SERVICE_ACCOUNT_TOKEN must be set
  if [ -z "${OP_SERVICE_ACCOUNT_TOKEN:-}" ]; then
    log_info "CodeRabbit: OP_SERVICE_ACCOUNT_TOKEN not set, skipping auto-auth"
    return 0
  fi

  # Guard: VAULT_ID must be resolved
  if [ -z "$VAULT_ID" ]; then
    log_info "CodeRabbit: OP_VAULT_ID not set and vault 'CI' not found, skipping auto-auth"
    return 0
  fi

  # Guard: vault must be accessible (quick connectivity check)
  if ! op vault get "$VAULT_ID" --format=json >/dev/null 2>&1; then
    log_warning "CodeRabbit: 1Password vault inaccessible (vault: $VAULT_ID), skipping"
    return 0
  fi

  # Retrieve token from 1Password (try multiple field names)
  local CR_TOKEN=""
  local fields=("credential" "password" "identifiant" "mot de passe")
  for field in "${fields[@]}"; do
    CR_TOKEN=$(op item get "mcp-coderabbit" --vault "$VAULT_ID" --fields "$field" --reveal 2>/dev/null || echo "")
    [ -n "$CR_TOKEN" ] && break
  done

  if [ -z "$CR_TOKEN" ]; then
    log_info "CodeRabbit: item 'mcp-coderabbit' not found in 1Password, skipping"
    return 0
  fi

  # Guard: jq required for safe JSON generation
  if ! command -v jq >/dev/null 2>&1; then
    log_warning "CodeRabbit: jq not available, cannot generate auth.json safely"
    return 0
  fi

  # Create auth directory (with error handling)
  if ! mkdir -p "$CR_AUTH_DIR" 2>/dev/null; then
    log_warning "CodeRabbit: cannot create $CR_AUTH_DIR, skipping"
    return 0
  fi
  chmod 700 "$CR_AUTH_DIR" 2>/dev/null || true

  # Write auth.json using jq for safe JSON escaping (handles special chars in token)
  local cr_tmp
  cr_tmp=$(mktemp "${CR_AUTH_FILE}.tmp.XXXXXX" 2>/dev/null) || {
    log_warning "CodeRabbit: cannot create temp file, skipping"
    return 0
  }
  trap 'rm -f "${cr_tmp:-}" 2>/dev/null || true' RETURN

  # Write OAuth auth format (accessToken + provider)
  # Note: both OAuth and API key tokens start with "cr-" so prefix detection
  # cannot distinguish them. We always use OAuth format because:
  #   - OAuth: works without paid addon
  #   - API key: requires "Hoppy CLI" subscription
  local cr_provider="github"
  local remote_url
  remote_url=$(git config --get remote.origin.url 2>/dev/null || echo "")
  if echo "$remote_url" | grep -qi "gitlab"; then
    cr_provider="gitlab"
  fi

  if ! jq -n --arg token "$CR_TOKEN" --arg provider "$cr_provider" \
    '{"accessToken": $token, "provider": $provider}' >"$cr_tmp" 2>/dev/null; then
    log_warning "CodeRabbit: failed to generate auth.json"
    rm -f "$cr_tmp" 2>/dev/null || true
    return 0
  fi

  # Validate and install
  if jq empty "$cr_tmp" 2>/dev/null; then
    mv "$cr_tmp" "$CR_AUTH_FILE"
    chmod 600 "$CR_AUTH_FILE" 2>/dev/null || true
    log_success "CodeRabbit authenticated via 1Password"
  else
    log_warning "CodeRabbit: generated auth.json is invalid, skipping"
    rm -f "$cr_tmp" 2>/dev/null || true
  fi
}

# Authenticate Qodo Command CLI via 1Password
# Item: mcp-qodo in 1Password vault
# Env var: QODO_API_KEY (written to ~/.devcontainer-env.sh)
step_qodo_auth() {
  local VAULT_ID
  VAULT_ID=$(resolve_vault_id)

  # Skip if Qodo CLI not installed
  if ! command -v qodo &>/dev/null; then
    log_info "Qodo CLI not installed, skipping auth"
    return 0
  fi

  # Skip if already configured
  if grep -q "^export QODO_API_KEY=" "$HOME/.devcontainer-env.sh" 2>/dev/null; then
    local existing_key
    existing_key=$(grep "^export QODO_API_KEY=" "$HOME/.devcontainer-env.sh" | head -1 | cut -d= -f2- | tr -d '"')
    if [ -n "$existing_key" ]; then
      log_success "Qodo already authenticated"
      return 0
    fi
  fi

  # Guard: 1Password CLI must be available
  if ! command -v op &>/dev/null; then
    log_info "Qodo: op CLI not installed, skipping auto-auth"
    return 0
  fi

  # Guard: OP_SERVICE_ACCOUNT_TOKEN must be set
  if [ -z "${OP_SERVICE_ACCOUNT_TOKEN:-}" ]; then
    log_info "Qodo: OP_SERVICE_ACCOUNT_TOKEN not set, skipping auto-auth"
    return 0
  fi

  # Guard: VAULT_ID must be resolved
  if [ -z "$VAULT_ID" ]; then
    log_info "Qodo: OP_VAULT_ID not set and vault 'CI' not found, skipping auto-auth"
    return 0
  fi

  # Guard: vault must be accessible (quick connectivity check)
  if ! op vault get "$VAULT_ID" --format=json >/dev/null 2>&1; then
    log_warning "Qodo: 1Password vault inaccessible (vault: $VAULT_ID), skipping"
    return 0
  fi

  # Retrieve API key from 1Password (try multiple field names)
  local QODO_KEY=""
  local fields=("credential" "password" "identifiant" "mot de passe")
  for field in "${fields[@]}"; do
    QODO_KEY=$(op item get "mcp-qodo" --vault "$VAULT_ID" --fields "$field" --reveal 2>/dev/null || echo "")
    [ -n "$QODO_KEY" ] && break
  done

  if [ -z "$QODO_KEY" ]; then
    log_info "Qodo: item 'mcp-qodo' not found in 1Password, skipping"
    return 0
  fi

  # Write to environment file (safe append with dedup)
  sed -i '/^export QODO_API_KEY=/d' "$HOME/.devcontainer-env.sh" 2>/dev/null || true
  printf "export QODO_API_KEY='%s'\n" "$QODO_KEY" >>"$HOME/.devcontainer-env.sh"
  chmod 600 "$HOME/.devcontainer-env.sh" 2>/dev/null || true

  # Export for current session
  export QODO_API_KEY="$QODO_KEY"

  log_success "Qodo authenticated via 1Password"
}

# Clean git credential helpers (remove macOS-specific helpers)
step_git_credential_cleanup() {
  log_info "Cleaning git credential helpers..."
  git config --global --unset-all credential.https://github.com.helper 2>/dev/null || true
  git config --global --unset-all credential.https://gist.github.com.helper 2>/dev/null || true
  log_success "Git credential helpers cleaned"
}

# Auto-update Claude Code to latest version
# Runs in background to avoid blocking container startup
# Fail-open: never blocks the DevContainer lifecycle
step_update_claude_code() {
  local PID_FILE="/tmp/.claude-update.pid"
  local INIT_MARKER="$HOME/.devcontainer-init-done"

  if ! command -v claude &>/dev/null; then
    log_info "Claude Code not installed, skipping update"
    return 0
  fi

  if [ -n "${CI:-}" ]; then
    log_info "CI environment detected, skipping Claude Code update"
    return 0
  fi

  # Skip if initial project init is still pending (avoid CLI contention)
  if [ ! -f "$INIT_MARKER" ]; then
    log_info "Initial project init still pending, deferring Claude Code update"
    return 0
  fi

  # Skip if a previous update is still running (prevent concurrent npm writes)
  if [ -f "$PID_FILE" ] && kill -0 "$(cat "$PID_FILE" 2>/dev/null)" 2>/dev/null; then
    log_info "Claude Code update already in progress, skipping"
    return 0
  fi

  local CURRENT_VERSION
  CURRENT_VERSION=$(claude --version 2>/dev/null | head -1 | grep -oP '[\d.]+' || echo "unknown")
  log_info "Current Claude Code version: $CURRENT_VERSION"

  # Run update in background (fail-open, non-blocking)
  nohup bash -c '
        UPDATE_LOG="$HOME/.claude-update.log"
        echo "[$(date -Iseconds)] Starting Claude Code update..." >> "$UPDATE_LOG"
        OLD_VER=$(claude --version 2>/dev/null | head -1 || echo "unknown")
        if npm update -g @anthropic-ai/claude-code >> "$UPDATE_LOG" 2>&1; then
            NEW_VER=$(claude --version 2>/dev/null | head -1 || echo "unknown")
            echo "[$(date -Iseconds)] Update complete: $OLD_VER -> $NEW_VER" >> "$UPDATE_LOG"
        else
            echo "[$(date -Iseconds)] Update failed (non-blocking)" >> "$UPDATE_LOG"
        fi
        rm -f /tmp/.claude-update.pid
    ' >>/tmp/claude-update.log 2>&1 &
  echo $! >"$PID_FILE"

  log_success "Claude Code update scheduled (logs: ~/.claude-update.log)"
}

# Auto-run /init for project initialization check
# Runs at every container start to verify project is properly initialized
# Skipped in CI environment
step_auto_init_check() {
  local INIT_LOG="$HOME/.devcontainer-init.log"
  local INIT_MARKER="$HOME/.devcontainer-init-done"

  # Only run /init on FIRST container start, not every restart
  if [ -f "$INIT_MARKER" ]; then
    log_info "Init already completed (marker: ~/.devcontainer-init-done), skipping"
    return 0
  fi

  if command -v claude &>/dev/null && [ -z "${CI:-}" ]; then
    log_info "Running project initialization check (first start)..."
    nohup bash -c "sleep 2 && claude \"/init\" && touch \"$INIT_MARKER\" || echo \"[\$(date -Iseconds)] Init check failed with exit code \$?\" >> \"$INIT_LOG\"" >>"$INIT_LOG" 2>&1 &
    log_success "Init check scheduled (logs: ~/.devcontainer-init.log)"
  elif [ -n "${CI:-}" ]; then
    log_info "CI environment detected, skipping init"
  fi
}

# ============================================================================
# Ollama + grepai Initialization (for semantic code search MCP)
# ============================================================================
GREPAI_BIN="/usr/local/bin/grepai"
GREPAI_CONFIG_TPL="/etc/grepai/config.yaml"
OLLAMA_HOST_ENDPOINT="${OLLAMA_HOST:-host.docker.internal:11434}"

detect_ollama_endpoint() {
  local endpoint=""
  local source=""

  if [ -n "${OLLAMA_HOST:-}" ]; then
    endpoint="$OLLAMA_HOST"
    source="OLLAMA_HOST env var"
    if curl -sf --connect-timeout 3 "http://${endpoint}/api/tags" >/dev/null 2>&1; then
      echo "$endpoint|$source"
      return 0
    else
      log_warning "OLLAMA_HOST=$endpoint not responding"
    fi
  fi

  endpoint="$OLLAMA_HOST_ENDPOINT"
  if curl -sf --connect-timeout 3 "http://${endpoint}/api/tags" >/dev/null 2>&1; then
    source="host (GPU-accelerated)"
    echo "$endpoint|$source"
    return 0
  fi

  echo ""
  return 1
}

check_model_available() {
  local endpoint="$1"
  local model="$2"
  curl -sf "http://${endpoint}/api/tags" 2>/dev/null | grep -q "$model"
}

show_ollama_instructions() {
  log_warning "==============================================================================="
  log_warning "  Ollama not running - Semantic search (grepai) will be disabled"
  log_warning "==============================================================================="
  log_info ""
  log_info "Ollama should be installed automatically via initialize.sh"
  log_info "If not running, start it manually on your host machine:"
  log_info ""
  log_info "  macOS:"
  log_info "    brew services start ollama"
  log_info "    # or: ollama serve"
  log_info ""
  log_info "  Linux:"
  log_info "    sudo systemctl start ollama"
  log_info "    # or: ollama serve"
  log_info ""
  log_info "  Then restart the DevContainer or run: /init"
  log_info ""
  log_warning "==============================================================================="
}

# --- Health stamp helpers ---
# The health stamp tracks 3 invalidation factors: model, binary version, config hash.
# If any factor changes between container starts, the index is purged and rebuilt.

compute_config_hash() {
  # Hash the config EXCLUDING the endpoint line (which changes per environment)
  grep -v '^\s*endpoint:' "$1" 2>/dev/null | md5sum | awk '{print $1}'
}

get_grepai_version() {
  "$GREPAI_BIN" version 2>/dev/null | grep -oE '[0-9]+\.[0-9]+\.[0-9]+' | head -1 || echo "unknown"
}

read_health_stamp() {
  # Reads .health-stamp into STAMP_* variables in the caller's scope
  # Returns 1 if stamp doesn't exist
  local stamp_file="$1"
  STAMP_MODEL=""
  STAMP_GREPAI_VERSION=""
  STAMP_CONFIG_HASH=""
  STAMP_DAEMON_PID=""
  STAMP_LAST_HEALTHY=""

  [ -f "$stamp_file" ] || return 1

  STAMP_MODEL=$(grep '^MODEL=' "$stamp_file" 2>/dev/null | cut -d= -f2-)
  STAMP_GREPAI_VERSION=$(grep '^GREPAI_VERSION=' "$stamp_file" 2>/dev/null | cut -d= -f2-)
  STAMP_CONFIG_HASH=$(grep '^CONFIG_HASH=' "$stamp_file" 2>/dev/null | cut -d= -f2-)
  # shellcheck disable=SC2034  # STAMP_DAEMON_PID used by grepai_watchdog via read_health_stamp
  STAMP_DAEMON_PID=$(grep '^DAEMON_PID=' "$stamp_file" 2>/dev/null | cut -d= -f2-)
  # shellcheck disable=SC2034  # STAMP_LAST_HEALTHY used by write_health_stamp
  STAMP_LAST_HEALTHY=$(grep '^LAST_HEALTHY=' "$stamp_file" 2>/dev/null | cut -d= -f2-)
  return 0
}

write_health_stamp() {
  local stamp_file="$1"
  local model="$2"
  local version="$3"
  local config_hash="$4"
  local daemon_pid="$5"

  cat >"$stamp_file" <<STAMP_EOF
MODEL=$model
GREPAI_VERSION=$version
CONFIG_HASH=$config_hash
DAEMON_PID=$daemon_pid
LAST_HEALTHY=$(date +%s)
STAMP_EOF
  chmod 644 "$stamp_file"
}

stop_grepai_daemon() {
  local pid
  pid=$(pgrep -f "$GREPAI_BIN watch" 2>/dev/null || true)
  if [ -n "$pid" ]; then
    log_info "Stopping grepai daemon (PID: $pid)..."
    kill "$pid" 2>/dev/null || true
    sleep 2
    # Force kill if still alive
    if kill -0 "$pid" 2>/dev/null; then
      kill -9 "$pid" 2>/dev/null || true
      sleep 1
    fi
  fi
  # Clean stale lock (daemon may have held it when killed)
  rm -f "${WORKSPACE_FOLDER:-/workspace}/.grepai/index.gob.lock"
}

_grepai_init_core() {
  local quiet="${1:-false}"
  local grepai_dir="${WORKSPACE_FOLDER:-/workspace}/.grepai"
  local grepai_config="${grepai_dir}/config.yaml"
  local health_stamp="${grepai_dir}/.health-stamp"
  local grepai_log="/tmp/grepai.log"
  local detected_result=""
  local ollama_endpoint=""
  local ollama_source=""

  # --- Step 1: Detect Ollama endpoint ---
  [ "$quiet" = "false" ] && log_info "Checking Ollama on host ($OLLAMA_HOST_ENDPOINT)..."
  detected_result=$(detect_ollama_endpoint)

  if [ -n "$detected_result" ]; then
    ollama_endpoint=$(echo "$detected_result" | cut -d'|' -f1)
    ollama_source=$(echo "$detected_result" | cut -d'|' -f2)
    log_success "Ollama connected: $ollama_endpoint ($ollama_source)"
  else
    [ "$quiet" = "false" ] && show_ollama_instructions
    # Pre-initialize config for when Ollama becomes available
    if [ -x "$GREPAI_BIN" ] && [ -f "$GREPAI_CONFIG_TPL" ]; then
      mkdir -p "$grepai_dir"
      cp "$GREPAI_CONFIG_TPL" "$grepai_config" 2>/dev/null || true
      sed -i -E "s|(endpoint: http://)[^[:space:]]+|\1${OLLAMA_HOST_ENDPOINT}|" "$grepai_config" 2>/dev/null || true
      [ "$quiet" = "false" ] && log_info "grepai config initialized (waiting for Ollama)"
    fi
    return 1
  fi

  # --- Step 2: Verify grepai binary ---
  if [ ! -x "$GREPAI_BIN" ]; then
    [ "$quiet" = "false" ] && log_warning "grepai binary not found at $GREPAI_BIN"
    return 2
  fi

  # --- Step 3: Sync config from template (always, to prevent drift) ---
  mkdir -p "$grepai_dir"

  if [ -f "$GREPAI_CONFIG_TPL" ]; then
    cp "$GREPAI_CONFIG_TPL" "$grepai_config"
    sed -i -E "s|(endpoint: http://)[^[:space:]]+|\1${ollama_endpoint}|" "$grepai_config"
    log_success "grepai config synced from template (endpoint: http://$ollama_endpoint)"
  else
    log_warning "Config template not found at $GREPAI_CONFIG_TPL, using grepai init..."
    (cd "${WORKSPACE_FOLDER:-/workspace}" && "$GREPAI_BIN" init --provider ollama --backend gob --yes 2>/dev/null) || true
    if [ -f "$grepai_config" ]; then
      sed -i -E "s|(endpoint: http://)[^[:space:]]+|\1${ollama_endpoint}|" "$grepai_config"
    fi
  fi

  # --- Step 4: Compute current state ---
  local current_model current_version current_config_hash
  current_model=$(grep -E '^\s+model:' "$grepai_config" 2>/dev/null | awk '{print $2}' | head -1)
  current_version=$(get_grepai_version)
  current_config_hash=$(compute_config_hash "$grepai_config")

  [ "$quiet" = "false" ] && log_info "grepai state: model=$current_model version=$current_version config=$current_config_hash"

  # --- Step 5-6: Multi-factor invalidation detection ---
  local need_rebuild=false
  local rebuild_reasons=""

  if read_health_stamp "$health_stamp"; then
    # Factor 1: model change (embeddings become incompatible)
    if [ -n "$current_model" ] && [ "$STAMP_MODEL" != "$current_model" ]; then
      log_warning "Model changed: ${STAMP_MODEL:-unknown} -> $current_model"
      need_rebuild=true
      rebuild_reasons="${rebuild_reasons}model_change "
    fi

    # Factor 2: grepai binary version change (index format may change)
    if [ "$STAMP_GREPAI_VERSION" != "$current_version" ]; then
      log_warning "grepai version changed: ${STAMP_GREPAI_VERSION:-unknown} -> $current_version"
      need_rebuild=true
      rebuild_reasons="${rebuild_reasons}version_change "
    fi

    # Factor 3: config change (chunk size, ignore patterns, etc.)
    if [ "$STAMP_CONFIG_HASH" != "$current_config_hash" ]; then
      log_warning "Config changed: ${STAMP_CONFIG_HASH:-unknown} -> $current_config_hash"
      need_rebuild=true
      rebuild_reasons="${rebuild_reasons}config_change "
    fi
  else
    # No health stamp = fresh install or first run after migration
    [ "$quiet" = "false" ] && log_info "No health stamp found (fresh install or first run)"

    # Migrate from legacy .model-stamp if present
    local legacy_stamp="${grepai_dir}/.model-stamp"
    if [ -f "$legacy_stamp" ]; then
      local legacy_model
      legacy_model=$(cat "$legacy_stamp" 2>/dev/null || echo "")
      if [ -n "$legacy_model" ] && [ "$legacy_model" != "$current_model" ]; then
        log_warning "Legacy stamp model mismatch: $legacy_model -> $current_model"
        need_rebuild=true
        rebuild_reasons="${rebuild_reasons}legacy_model_change "
      fi
      rm -f "$legacy_stamp"
      log_info "Migrated from legacy .model-stamp"
    fi

    # If index exists but no stamp, we can't trust it — rebuild
    if [ -f "${grepai_dir}/index.gob" ] && [ "$need_rebuild" = "false" ]; then
      log_warning "Index exists but no health stamp — cannot verify integrity"
      need_rebuild=true
      rebuild_reasons="${rebuild_reasons}missing_stamp "
    fi
  fi

  # --- Step 7: Handle invalidation ---
  if [ "$need_rebuild" = "true" ]; then
    log_warning "Index rebuild required: ${rebuild_reasons}"
    stop_grepai_daemon
    rm -f "${grepai_dir}/index.gob" "${grepai_dir}/symbols.gob" "${grepai_dir}/index.gob.lock"
    log_success "Index cleared — will rebuild from scratch"
  fi

  # --- Step 8: Verify model is available before starting daemon ---
  if [ -n "$current_model" ]; then
    if check_model_available "$ollama_endpoint" "$current_model"; then
      log_success "Model $current_model available on Ollama"
    else
      log_warning "Model $current_model not found on Ollama"
      [ "$quiet" = "false" ] && log_info "Pull the model on your host: ollama pull $current_model"
      [ "$quiet" = "false" ] && log_warning "Skipping daemon start until model is available"
      return 2
    fi
  fi

  # --- Step 9: Start or verify daemon ---
  local grepai_pid
  grepai_pid=$(pgrep -f "$GREPAI_BIN watch" 2>/dev/null || true)

  if [ -z "$grepai_pid" ]; then
    # Clean stale lock from previous crashed daemon (no process = lock is stale)
    rm -f "${grepai_dir}/index.gob.lock"
    log_info "Starting grepai watch daemon..."
    : >"$grepai_log"
    (cd "${WORKSPACE_FOLDER:-/workspace}" && nohup "$GREPAI_BIN" watch >"$grepai_log" 2>&1 &)

    # Retry loop: wait up to 5 seconds for daemon to start
    local retries=10
    while [ $retries -gt 0 ]; do
      grepai_pid=$(pgrep -f "$GREPAI_BIN watch" 2>/dev/null || true)
      [ -n "$grepai_pid" ] && break
      retries=$((retries - 1))
      sleep 0.5
    done

    if [ -n "$grepai_pid" ]; then
      log_success "grepai watch daemon started (PID: $grepai_pid)"
    else
      log_warning "grepai daemon failed to start (check $grepai_log)"
      return 2
    fi
  else
    log_info "grepai watch daemon already running (PID: $grepai_pid)"
  fi

  # --- Step 10: Save health stamp ---
  write_health_stamp "$health_stamp" \
    "$current_model" "$current_version" "$current_config_hash" "$grepai_pid"
  log_success "Health stamp saved (model=$current_model ver=$current_version)"
}

init_semantic_search() {
  _grepai_init_core "false"
  # Always launch watchdog (handles both daemon monitoring AND deferred init)
  grepai_watchdog &
}

# Watchdog: monitors grepai daemon and handles deferred initialization.
# When health stamp exists: restarts crashed daemon.
# When no health stamp: retries init when Ollama becomes available.
# Runs in background for the lifetime of the container.
grepai_watchdog() {
  local grepai_dir="${WORKSPACE_FOLDER:-/workspace}/.grepai"
  local health_stamp="${grepai_dir}/.health-stamp"
  local grepai_log="/tmp/grepai.log"
  local deferred_attempts=0

  # Write PID file for discoverability (shell functions are invisible to pgrep)
  echo $$ >/tmp/grepai-watchdog.pid

  # Let the container finish starting before first check
  sleep 30

  while true; do
    sleep 60

    if [ ! -f "$health_stamp" ]; then
      # Deferred init: Ollama may have become available since startup
      if detect_ollama_endpoint >/dev/null 2>&1; then
        deferred_attempts=$((deferred_attempts + 1))
        if [ "$deferred_attempts" -eq 1 ]; then
          log_info "[WATCHDOG] Ollama now available — initializing semantic search"
        elif [ $((deferred_attempts % 5)) -eq 0 ]; then
          log_warning "[WATCHDOG] Deferred init retry #${deferred_attempts}"
        fi

        if _grepai_init_core "true"; then
          log_success "[WATCHDOG] Deferred initialization complete"
          deferred_attempts=0
        fi
      fi
      continue
    fi

    # Reset counter once healthy
    deferred_attempts=0

    local current_pid
    current_pid=$(pgrep -f "$GREPAI_BIN watch" 2>/dev/null || true)

    if [ -z "$current_pid" ]; then
      log_warning "[WATCHDOG] grepai daemon not running — restarting..."

      # Verify Ollama is still reachable before restarting
      if ! detect_ollama_endpoint >/dev/null 2>&1; then
        log_warning "[WATCHDOG] Ollama not reachable, skipping restart"
        continue
      fi

      # Clean stale lock from crashed daemon before restart
      rm -f "${grepai_dir}/index.gob.lock"
      (cd "${WORKSPACE_FOLDER:-/workspace}" && nohup "$GREPAI_BIN" watch >>"$grepai_log" 2>&1 &)
      sleep 3

      current_pid=$(pgrep -f "$GREPAI_BIN watch" 2>/dev/null || true)
      if [ -n "$current_pid" ]; then
        log_success "[WATCHDOG] Daemon restarted (PID: $current_pid)"
        # Update stamp with new PID
        if read_health_stamp "$health_stamp"; then
          write_health_stamp "$health_stamp" \
            "$STAMP_MODEL" "$STAMP_GREPAI_VERSION" \
            "$STAMP_CONFIG_HASH" "$current_pid"
        fi
      else
        log_warning "[WATCHDOG] Failed to restart daemon (check $grepai_log)"
      fi
    fi
  done
}

# ============================================================================
# VPN Auto-Connect (optional - skipped if no config found)
# ============================================================================
# Multi-protocol VPN support: OpenVPN, WireGuard, IPsec/IKEv2, PPTP
# Config source: VPN_CONFIG_REF=op://VAULT/PROFILE (or legacy OPENVPN_CONFIG_REF)
#   - DOCUMENT "PROFILE" in vault → config file (protocol determined by tags)
#   - LOGIN "PROFILE" in vault → username/password (not needed for WireGuard)
#   - Tags: "openvpn" (default), "wireguard", "ipsec", "pptp"
#   - Fallback: file on disk at $OPENVPN_CONFIG
#   - Nothing found → skip silently

# --- OpenVPN connect (extracted for multi-protocol support) ---
connect_openvpn() {
  local vault="$1" profile="$2" doc_uuid="$3" login_uuid="$4"
  local ovpn_config="${OPENVPN_CONFIG:-$HOME/.config/openvpn/client.ovpn}"
  local ovpn_auth="${OPENVPN_AUTH:-/tmp/vpn-auth.txt}"
  local ovpn_dir
  ovpn_dir=$(dirname "$ovpn_config")

  mkdir -p "$ovpn_dir"

  # Download .ovpn config
  if [ -n "$doc_uuid" ] && op document get "$doc_uuid" --vault "$vault" >"$ovpn_config" 2>/dev/null; then
    chmod 600 "$ovpn_config"
    log_success "Resolved .ovpn config ($vault/$profile)"
  else
    log_warning "No DOCUMENT '$profile' in vault '$vault', skipping VPN"
    return 0
  fi

  # Resolve credentials
  if [ -n "$login_uuid" ]; then
    local vpn_user vpn_pass
    vpn_user=$(op read "op://$vault/$login_uuid/username" 2>/dev/null || echo "")
    vpn_pass=$(op read "op://$vault/$login_uuid/password" 2>/dev/null || echo "")
    if [ -n "$vpn_user" ] && [ -n "$vpn_pass" ]; then
      printf '%s\n%s\n' "$vpn_user" "$vpn_pass" >"$ovpn_auth"
      chmod 600 "$ovpn_auth"
      log_success "Resolved VPN credentials ($vault/$profile)"
    fi
  fi

  # Validate config
  if [ ! -s "$ovpn_config" ]; then
    log_warning "OpenVPN config is empty: $ovpn_config"
    return 0
  fi

  local -a vpn_args=(
    --config "$ovpn_config"
    --daemon ovpn-client
    --log /tmp/openvpn.log
    --script-security 2
    --up /etc/openvpn/update-dns
    --down /etc/openvpn/update-dns
    --keepalive 10 60
    --connect-retry 5
    --connect-retry-max 0
    --persist-tun
    --persist-key
    --resolv-retry infinite
  )
  if [ -f "$ovpn_auth" ] && [ -s "$ovpn_auth" ]; then
    vpn_args+=(--auth-user-pass "$ovpn_auth")
  fi

  log_info "Starting OpenVPN..."
  if sudo openvpn "${vpn_args[@]}"; then
    local attempt=0
    while [ $attempt -lt 15 ]; do
      if ip link show tun0 &>/dev/null; then
        local vpn_ip
        vpn_ip=$(ip -4 addr show tun0 2>/dev/null | grep -oP 'inet \K[\d.]+' || echo "unknown")
        log_success "VPN connected via OpenVPN (tun0: $vpn_ip)"
        return 0
      fi
      sleep 1
      ((attempt++))
    done
    log_warning "OpenVPN started but tun0 not detected after 15s (check /tmp/openvpn.log)"
  else
    log_warning "OpenVPN failed to start (check /tmp/openvpn.log)"
  fi
}

# --- WireGuard connect ---
connect_wireguard() {
  local vault="$1" profile="$2" doc_uuid="$3"
  local wg_config="$HOME/.config/wireguard/wg0.conf"

  mkdir -p "$(dirname "$wg_config")"

  if [ -n "$doc_uuid" ] && op document get "$doc_uuid" --vault "$vault" >"$wg_config" 2>/dev/null; then
    chmod 600 "$wg_config"
    log_success "Resolved WireGuard config ($vault/$profile)"
  else
    log_warning "No DOCUMENT '$profile' in vault '$vault', skipping WireGuard"
    return 0
  fi

  log_info "Starting WireGuard..."
  if sudo wg-quick up "$wg_config" 2>/dev/null; then
    local attempt=0
    while [ $attempt -lt 10 ]; do
      if ip link show wg0 &>/dev/null; then
        local vpn_ip
        vpn_ip=$(ip -4 addr show wg0 2>/dev/null | grep -oP 'inet \K[\d.]+' || echo "unknown")
        log_success "VPN connected via WireGuard (wg0: $vpn_ip)"
        return 0
      fi
      sleep 1
      ((attempt++))
    done
    log_warning "WireGuard started but wg0 not detected after 10s"
  else
    log_warning "WireGuard failed to start"
  fi
}

# --- IPsec/IKEv2 connect ---
connect_ipsec() {
  local vault="$1" profile="$2" doc_uuid="$3" login_uuid="$4"
  local ipsec_config="$HOME/.config/strongswan/ipsec.conf"

  mkdir -p "$(dirname "$ipsec_config")"

  if [ -n "$doc_uuid" ] && op document get "$doc_uuid" --vault "$vault" >"$ipsec_config" 2>/dev/null; then
    chmod 600 "$ipsec_config"
    log_success "Resolved IPsec config ($vault/$profile)"
  else
    log_warning "No DOCUMENT '$profile' in vault '$vault', skipping IPsec"
    return 0
  fi

  # Copy config and secrets to strongswan dir
  sudo cp "$ipsec_config" /etc/ipsec.d/profile.conf
  if [ -n "$login_uuid" ]; then
    local vpn_user vpn_pass
    vpn_user=$(op read "op://$vault/$login_uuid/username" 2>/dev/null || echo "")
    vpn_pass=$(op read "op://$vault/$login_uuid/password" 2>/dev/null || echo "")
    if [ -n "$vpn_user" ] && [ -n "$vpn_pass" ]; then
      printf '%s : EAP "%s"\n' "$vpn_user" "$vpn_pass" | sudo tee /etc/ipsec.d/profile.secrets >/dev/null
      sudo chmod 600 /etc/ipsec.d/profile.secrets
    fi
  fi

  local conn_name
  conn_name=$(grep -oP '(?<=^conn )\S+' "$ipsec_config" | head -1)
  if [ -z "$conn_name" ]; then
    log_warning "No connection name found in IPsec config"
    return 0
  fi

  log_info "Starting IPsec ($conn_name)..."
  sudo ipsec restart 2>/dev/null
  sleep 2
  if sudo ipsec up "$conn_name" 2>/dev/null; then
    log_success "VPN connected via IPsec ($conn_name)"
  else
    log_warning "IPsec connection '$conn_name' failed"
  fi
}

# --- PPTP connect ---
connect_pptp() {
  local vault="$1" profile="$2" doc_uuid="$3" login_uuid="$4"
  local pptp_config="$HOME/.config/pptp/tunnel.conf"

  mkdir -p "$(dirname "$pptp_config")"

  if [ -n "$doc_uuid" ] && op document get "$doc_uuid" --vault "$vault" >"$pptp_config" 2>/dev/null; then
    chmod 600 "$pptp_config"
    log_success "Resolved PPTP config ($vault/$profile)"
  else
    log_warning "No DOCUMENT '$profile' in vault '$vault', skipping PPTP"
    return 0
  fi

  if [ -n "$login_uuid" ]; then
    local vpn_user vpn_pass
    vpn_user=$(op read "op://$vault/$login_uuid/username" 2>/dev/null || echo "")
    vpn_pass=$(op read "op://$vault/$login_uuid/password" 2>/dev/null || echo "")
    if [ -n "$vpn_user" ] && [ -n "$vpn_pass" ]; then
      printf '%s\n%s\n' "$vpn_user" "$vpn_pass" >/tmp/vpn-auth.txt
      chmod 600 /tmp/vpn-auth.txt
    fi
  fi

  log_info "Starting PPTP..."
  # shellcheck disable=SC2024
  sudo pppd call tunnel nodetach </dev/null >/tmp/pptp.log 2>&1 &

  local attempt=0
  while [ $attempt -lt 15 ]; do
    if ip link show ppp0 &>/dev/null; then
      local vpn_ip
      vpn_ip=$(ip -4 addr show ppp0 2>/dev/null | grep -oP 'inet \K[\d.]+' || echo "unknown")
      log_success "VPN connected via PPTP (ppp0: $vpn_ip)"
      return 0
    fi
    sleep 1
    ((attempt++))
  done
  log_warning "PPTP started but ppp0 not detected after 15s"
}

# --- RTK CLI proxy initialization ---
init_rtk() {
  if ! command -v rtk &>/dev/null; then
    log_info "RTK not found, installing latest from GitHub..."
    if ! command -v jq >/dev/null 2>&1; then
      log_warning "RTK: jq not available, skipping installation"
      return 0
    fi
    local rtk_arch
    case "$(uname -m)" in
      x86_64) rtk_arch="x86_64-unknown-linux-musl" ;;
      aarch64) rtk_arch="aarch64-unknown-linux-musl" ;;
      *)
        log_warning "RTK: unsupported architecture $(uname -m)"
        return 0
        ;;
    esac
    local curl_auth_args=()
    if [ -n "${GITHUB_API_TOKEN:-}" ]; then
      curl_auth_args=(-H "Authorization: token ${GITHUB_API_TOKEN}")
    fi
    local rtk_tag
    rtk_tag=$(curl -fsSL --connect-timeout 5 --max-time 15 "${curl_auth_args[@]}" \
      "https://api.github.com/repos/rtk-ai/rtk/releases/latest" 2>/dev/null | jq -r '.tag_name // empty')
    if [ -z "$rtk_tag" ] || ! [[ "$rtk_tag" =~ ^v?[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
      log_warning "RTK: failed to fetch valid release tag"
      return 0
    fi
    if curl -fsSL --connect-timeout 5 --max-time 60 "https://github.com/rtk-ai/rtk/releases/download/${rtk_tag}/rtk-${rtk_arch}.tar.gz" \
      | sudo tar xz -C /usr/local/bin rtk 2>/dev/null; then
      log_success "RTK ${rtk_tag} installed to /usr/local/bin/rtk"
    else
      log_warning "RTK: installation failed (non-blocking)"
      return 0
    fi
  fi

  # Sync config from template (with hash tracking for updates)
  local RTK_CONFIG_DIR="$HOME/.config/rtk"
  local RTK_CONFIG_SRC="/etc/rtk/config.toml"
  if [ -f "$RTK_CONFIG_SRC" ]; then
    mkdir -p "$RTK_CONFIG_DIR"
    local template_hash user_hash stored_hash
    template_hash=$(md5sum "$RTK_CONFIG_SRC" 2>/dev/null | cut -d' ' -f1)
    if [ ! -f "$RTK_CONFIG_DIR/config.toml" ]; then
      cp "$RTK_CONFIG_SRC" "$RTK_CONFIG_DIR/config.toml"
      echo "$template_hash" >"$RTK_CONFIG_DIR/.template_hash"
      log_info "RTK config initialized from template"
    else
      stored_hash=$(cat "$RTK_CONFIG_DIR/.template_hash" 2>/dev/null || echo "")
      user_hash=$(md5sum "$RTK_CONFIG_DIR/config.toml" 2>/dev/null | cut -d' ' -f1)
      if [ "$stored_hash" != "$template_hash" ] && [ "$user_hash" = "$stored_hash" ]; then
        cp "$RTK_CONFIG_SRC" "$RTK_CONFIG_DIR/config.toml"
        echo "$template_hash" >"$RTK_CONFIG_DIR/.template_hash"
        log_info "RTK config updated from new template"
      elif [ "$stored_hash" != "$template_hash" ]; then
        echo "$template_hash" >"$RTK_CONFIG_DIR/.template_hash"
        log_info "RTK config preserved (user customized)"
      fi
    fi
  fi

  local RTK_VER
  RTK_VER=$(rtk --version 2>/dev/null | grep -oE '[0-9]+\.[0-9]+\.[0-9]+' | head -1)
  log_success "RTK ${RTK_VER:-unknown} ready (token savings active)"
}

# --- Main VPN auto-connect orchestrator ---
init_vpn() {
  # Skip if no VPN tools installed at all
  local has_vpn=false
  command -v openvpn &>/dev/null && has_vpn=true
  command -v wg &>/dev/null && has_vpn=true
  command -v ipsec &>/dev/null && has_vpn=true
  command -v pptp &>/dev/null && has_vpn=true
  if [ "$has_vpn" = "false" ]; then
    log_debug "No VPN clients installed, skipping"
    return 0
  fi

  log_info "VPN clients detected, checking configuration..."

  # Skip if already connected (any protocol)
  if pgrep -x openvpn &>/dev/null || ip link show wg0 &>/dev/null 2>&1 \
    || pgrep -x charon &>/dev/null || pgrep -x pppd &>/dev/null; then
    log_info "VPN already connected, skipping"
    return 0
  fi

  # Backward compatible: support both VPN_CONFIG_REF and OPENVPN_CONFIG_REF
  local vpn_ref="${VPN_CONFIG_REF:-${OPENVPN_CONFIG_REF:-}}"

  # Source 1: Resolve from 1Password vault
  if [ -n "$vpn_ref" ] && [ -n "${OP_SERVICE_ACCOUNT_TOKEN:-}" ] && command -v op &>/dev/null; then
    local ref="${vpn_ref#op://}"
    local vault profile
    vault=$(echo "$ref" | cut -d'/' -f1)
    profile=$(echo "$ref" | cut -d'/' -f2)

    log_info "Resolving VPN profile '$profile' from 1Password ($vault)..."

    # Resolve DOCUMENT UUID and detect protocol from tags
    local doc_item doc_uuid protocol tags
    doc_item=$(op item list --vault "$vault" --categories DOCUMENT --format json 2>/dev/null \
      | jq -r --arg t "$profile" '.[] | select(.title==$t)' 2>/dev/null || echo "")
    doc_uuid=$(echo "$doc_item" | jq -r '.id // empty' 2>/dev/null || echo "")
    tags=$(echo "$doc_item" | jq -r '.tags // [] | .[]' 2>/dev/null || echo "")

    # Detect protocol (default: openvpn)
    protocol="openvpn"
    for tag in $tags; do
      case "$tag" in
        wireguard | ipsec | pptp)
          protocol="$tag"
          break
          ;;
      esac
    done

    # Resolve LOGIN UUID
    local login_uuid
    login_uuid=$(op item list --vault "$vault" --categories LOGIN --format json 2>/dev/null \
      | jq -r --arg t "$profile" '.[] | select(.title==$t) | .id' 2>/dev/null || echo "")

    log_info "VPN profile '$profile' → protocol: $protocol"

    # Protocol-specific connect
    case "$protocol" in
      openvpn) connect_openvpn "$vault" "$profile" "$doc_uuid" "$login_uuid" ;;
      wireguard) connect_wireguard "$vault" "$profile" "$doc_uuid" ;;
      ipsec) connect_ipsec "$vault" "$profile" "$doc_uuid" "$login_uuid" ;;
      pptp) connect_pptp "$vault" "$profile" "$doc_uuid" "$login_uuid" ;;
    esac
    return 0
  fi

  # Source 2: File on disk (OpenVPN fallback for backward compat)
  local ovpn_config="${OPENVPN_CONFIG:-$HOME/.config/openvpn/client.ovpn}"
  if [ -f "$ovpn_config" ] && [ -s "$ovpn_config" ] && command -v openvpn &>/dev/null; then
    log_info "Found OpenVPN config on disk, connecting..."
    connect_openvpn "" "" "" ""
  else
    log_info "No VPN config found, skipping"
  fi

  return 0
}

# ============================================================================
# Execution
# ============================================================================

run_step "Restore Claude config" step_restore_claude_config
run_step "Init Claude dirs" step_init_claude_dirs
run_step "Shell env repair" step_shell_env_repair
run_step "Cache completions" step_cache_completions
run_step "p10k segments" step_generate_p10k_segments
# NOTE: Environment reload MUST NOT run inside run_step (which uses a subshell).
# Variables sourced in a subshell are lost when it exits — tokens would never
# reach step_mcp_configuration. Source .env directly in the main shell.
_ENV_FILE="${WORKSPACE_FOLDER:-/workspace}/.devcontainer/.env"
if [ -f "$_ENV_FILE" ]; then
  log_info "Reloading environment from .env..."
  set -a
  # shellcheck source=/dev/null
  source "$_ENV_FILE"
  set +a
  log_success "Environment reloaded from .env"
else
  log_info "No .env file found, skipping environment reload"
fi
run_step "1Password permissions" step_1password_permissions
run_step "npm cache permissions" step_npm_cache_permissions
run_step "MCP configuration" step_mcp_configuration
run_step "Taskmaster config" step_taskmaster_config
run_step "CodeRabbit auth" step_coderabbit_auth
run_step "Qodo auth" step_qodo_auth
run_step "Git credential cleanup" step_git_credential_cleanup

# Background tasks (tracked via PID files for diagnostics)
init_semantic_search >>/tmp/grepai-init.log 2>&1 &
echo $! >/tmp/.grepai-init.pid
init_vpn >>/tmp/vpn-init.log 2>&1 &
echo $! >/tmp/.vpn-init.pid
run_step "RTK init" init_rtk

# Export dynamic environment variables (appended to ~/.devcontainer-env.sh)
# Note: ~/.devcontainer-env.sh is created by postCreate.sh with static content

run_step "Claude Code update" step_update_claude_code
run_step "Project init check" step_auto_init_check

print_step_summary "postStart"

log_success "postStart: Container ready!"
