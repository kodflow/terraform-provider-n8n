#!/bin/bash
# DevContainer Utility Functions
# Provides retry mechanisms, error handling, and logging

# Colors for output
export RED='\033[0;31m'
export GREEN='\033[0;32m'
export YELLOW='\033[1;33m'
export BLUE='\033[0;34m'
export CYAN='\033[0;36m'
export NC='\033[0m' # No Color

# Logging functions
log_info() {
  echo -e "${BLUE}[INFO]${NC} $*"
}

log_success() {
  echo -e "${GREEN}[SUCCESS]${NC} $*"
}

log_warning() {
  echo -e "${YELLOW}[WARNING]${NC} $*"
}

log_error() {
  echo -e "${RED}[ERROR]${NC} $*"
}

log_debug() {
  if [ "${DEBUG:-0}" = "1" ]; then
    echo -e "${CYAN}[DEBUG]${NC} $*"
  fi
}

# Retry function for any command
# Usage: retry <max_attempts> <delay_seconds> <command...>
# Example: retry 3 5 curl -o file.tar.gz https://example.com/file.tar.gz
retry() {
  local max_attempts=$1
  local delay=$2
  shift 2
  local attempt=1
  local exit_code=0

  while [ $attempt -le $max_attempts ]; do
    log_debug "Attempt $attempt/$max_attempts: $*"

    # Execute command
    if "$@"; then
      if [ $attempt -gt 1 ]; then
        log_success "Command succeeded on attempt $attempt"
      fi
      return 0
    fi

    exit_code=$?

    if [ $attempt -lt $max_attempts ]; then
      log_warning "Command failed (exit code: $exit_code), retrying in ${delay}s... (attempt $attempt/$max_attempts)"
      sleep "$delay"
    else
      log_error "Command failed after $max_attempts attempts"
    fi

    ((attempt++))
  done

  return $exit_code
}

# Retry with exponential backoff
# Usage: retry_exponential <max_attempts> <initial_delay> <command...>
# Example: retry_exponential 5 2 curl -o file.tar.gz https://example.com/file.tar.gz
retry_exponential() {
  local max_attempts=$1
  local initial_delay=$2
  shift 2
  local attempt=1
  local delay=$initial_delay
  local exit_code=0

  while [ $attempt -le $max_attempts ]; do
    log_debug "Attempt $attempt/$max_attempts (delay: ${delay}s): $*"

    # Execute command
    if "$@"; then
      if [ $attempt -gt 1 ]; then
        log_success "Command succeeded on attempt $attempt"
      fi
      return 0
    fi

    exit_code=$?

    if [ $attempt -lt $max_attempts ]; then
      log_warning "Command failed, retrying in ${delay}s... (attempt $attempt/$max_attempts)"
      sleep "$delay"
      delay=$((delay * 2)) # Double the delay for next attempt
    else
      log_error "Command failed after $max_attempts attempts"
    fi

    ((attempt++))
  done

  return $exit_code
}

# apt-get with retry and lock handling
# Usage: apt_get_retry <apt-get arguments...>
# Example: apt_get_retry install -y curl wget
apt_get_retry() {
  local max_attempts=5
  local attempt=1
  local delay=10

  while [ $attempt -le $max_attempts ]; do
    log_debug "apt-get attempt $attempt/$max_attempts: $*"

    # Wait for apt locks to be released
    local lock_wait=0
    while sudo fuser /var/lib/dpkg/lock-frontend >/dev/null 2>&1 \
      || sudo fuser /var/lib/apt/lists/lock >/dev/null 2>&1 \
      || sudo fuser /var/cache/apt/archives/lock >/dev/null 2>&1; do
      if [ $lock_wait -eq 0 ]; then
        log_warning "Waiting for apt locks to be released..."
      fi
      sleep 2
      lock_wait=$((lock_wait + 2))

      # If waiting too long, try to force unlock
      if [ $lock_wait -ge 60 ]; then
        log_warning "Forcing apt lock release after 60s wait"
        sudo rm -f /var/lib/dpkg/lock-frontend
        sudo rm -f /var/lib/apt/lists/lock
        sudo rm -f /var/cache/apt/archives/lock
        sudo dpkg --configure -a || true
        break
      fi
    done

    # Try apt-get command
    if sudo apt-get "$@"; then
      if [ $attempt -gt 1 ]; then
        log_success "apt-get succeeded on attempt $attempt"
      fi
      return 0
    fi

    exit_code=$?

    if [ $attempt -lt $max_attempts ]; then
      log_warning "apt-get failed, running update and retrying in ${delay}s... (attempt $attempt/$max_attempts)"
      sudo apt-get update --fix-missing || true
      sudo dpkg --configure -a || true
      sleep "$delay"
    else
      log_error "apt-get failed after $max_attempts attempts"
    fi

    ((attempt++))
  done

  return $exit_code
}

# Download with retry and resume support
# Usage: download_retry <url> <output_file> [curl_options...]
# Example: download_retry "https://example.com/file.tar.gz" "/tmp/file.tar.gz"
download_retry() {
  local url=$1
  local output=$2
  shift 2
  local extra_args=("$@")

  log_info "Downloading: $url"

  # Use curl with retry, continue from partial downloads, and follow redirects
  retry_exponential 5 3 curl -fsSL \
    --connect-timeout 30 \
    --max-time 300 \
    --retry 3 \
    --retry-delay 5 \
    --retry-max-time 60 \
    -C - \
    "${extra_args[@]}" \
    -o "$output" \
    "$url"
}

# Download and pipe to shell with retry
# Usage: download_and_pipe <url> <shell_command...>
# Example: download_and_pipe "https://get.pnpm.io/install.sh" "sh" "-"
download_and_pipe() {
  local url=$1
  shift
  local shell_cmd=("$@")

  log_info "Downloading and executing: $url"

  # Download to temp file first, then execute
  local temp_file
  temp_file=$(mktemp)

  if download_retry "$url" "$temp_file"; then
    "${shell_cmd[@]}" <"$temp_file"
    local exit_code=$?
    rm -f "$temp_file"
    return $exit_code
  else
    rm -f "$temp_file"
    return 1
  fi
}

# Check if command exists
# Usage: command_exists <command_name>
# Example: if command_exists node; then echo "Node is installed"; fi
command_exists() {
  command -v "$1" &>/dev/null
}

# Wait for network connectivity
# Usage: wait_for_network [max_wait_seconds]
# Example: wait_for_network 60
wait_for_network() {
  local max_wait=${1:-60}
  local waited=0

  log_info "Checking network connectivity..."

  while [ $waited -lt $max_wait ]; do
    if curl -sSf --connect-timeout 5 https://www.google.com >/dev/null 2>&1; then
      log_success "Network is available"
      return 0
    fi

    log_warning "Waiting for network... (${waited}s/${max_wait}s)"
    sleep 5
    waited=$((waited + 5))
  done

  log_error "Network not available after ${max_wait}s"
  return 1
}

# Safe directory creation
# Usage: mkdir_safe <directory_path>
# Example: mkdir_safe "/home/vscode/.cache/npm"
mkdir_safe() {
  local dir=$1
  if [ ! -d "$dir" ]; then
    mkdir -p "$dir" 2>/dev/null || sudo mkdir -p "$dir"

    # Fix permissions if we're the vscode user
    if [ "$(whoami)" = "vscode" ]; then
      sudo chown -R vscode:vscode "$dir" 2>/dev/null || true
    fi
  fi
}

# ============================================================================
# Step tracking functions (non-blocking lifecycle pattern)
# ============================================================================
# Usage:
#   init_steps
#   run_step "Step name" step_function
#   print_step_summary "scriptName"
#   if step_ok "Step name"; then ... fi

# Internal arrays for step tracking
_STEP_NAMES=()
_STEP_RESULTS=()

# Reset step tracking (call at start of script)
init_steps() {
  _STEP_NAMES=()
  _STEP_RESULTS=()
}

# Run a function in an isolated subshell, track PASS/FAIL
# Usage: run_step "Display name" function_name [args...]
run_step() {
  local name="$1"
  shift
  local func="$1"
  shift

  local rc=0
  # Run in subshell with errexit/pipefail disabled for isolation
  (
    set +e +o pipefail
    "$func" "$@"
  ) || rc=$?

  _STEP_NAMES+=("$name")
  if [ "$rc" -eq 0 ]; then
    _STEP_RESULTS+=("PASS")
  else
    _STEP_RESULTS+=("FAIL")
    log_warning "Step failed: $name (exit code $rc)"
  fi
  return 0
}

# Print summary table of all steps
# Usage: print_step_summary "postCreate"
print_step_summary() {
  local label="${1:-Summary}"
  local total=${#_STEP_NAMES[@]}
  local fail_count=0

  echo ""
  echo -e "${CYAN}=========================================${NC}"
  echo -e "${CYAN}   ${label} Summary${NC}"
  echo -e "${CYAN}=========================================${NC}"

  local i
  for i in $(seq 0 $((total - 1))); do
    if [ "${_STEP_RESULTS[$i]}" = "PASS" ]; then
      echo -e "  ${GREEN}[PASS]${NC} ${_STEP_NAMES[$i]}"
    else
      echo -e "  ${RED}[FAIL]${NC} ${_STEP_NAMES[$i]}"
      ((fail_count++))
    fi
  done

  echo ""
  if [ "$fail_count" -gt 0 ]; then
    log_warning "$fail_count of $total steps failed (non-blocking)"
  else
    log_success "All $total steps passed"
  fi
  echo ""
}

# Check if a named step succeeded
# Usage: if step_ok "Git safe directory"; then ...
step_ok() {
  local name="$1"
  local i
  for i in $(seq 0 $((${#_STEP_NAMES[@]} - 1))); do
    if [ "${_STEP_NAMES[$i]}" = "$name" ]; then
      [ "${_STEP_RESULTS[$i]}" = "PASS" ]
      return $?
    fi
  done
  return 1
}

# Export all functions for use in subscripts
export -f log_info
export -f log_success
export -f log_warning
export -f log_error
export -f log_debug
export -f retry
export -f retry_exponential
export -f apt_get_retry
export -f download_retry
export -f download_and_pipe
export -f command_exists
export -f wait_for_network
export -f mkdir_safe
export -f init_steps
export -f run_step
export -f print_step_summary
export -f step_ok
