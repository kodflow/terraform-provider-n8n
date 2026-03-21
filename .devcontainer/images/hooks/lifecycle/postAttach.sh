#!/bin/bash
# shellcheck disable=SC1091
# ============================================================================
# postAttach.sh - Runs when IDE connects to the container
# ============================================================================
# This script runs when a tool (like VS Code) connects to the dev container.
# It's the only command that consistently allows user interaction.
# Use it for: Welcome messages, status display, interactive prompts.
# ============================================================================

set -u

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/../shared/utils.sh"

echo ""
echo -e "${CYAN}=========================================${NC}"
echo -e "${CYAN}   DevContainer Ready${NC}"
echo -e "${CYAN}=========================================${NC}"
echo ""

# Display useful information
log_success "IDE connected to DevContainer"
echo ""

# Verify super-claude is available before advertising it
if [ -f "$HOME/.devcontainer-env.sh" ] && grep -q "super-claude" "$HOME/.devcontainer-env.sh" 2>/dev/null; then
  echo -e "Tip: Use ${GREEN}super-claude${NC} for Claude CLI with MCP config"
else
  log_warning "super-claude not available - environment setup may have failed"
  log_info "Try running: source ~/.devcontainer-env.sh"
fi
echo ""
