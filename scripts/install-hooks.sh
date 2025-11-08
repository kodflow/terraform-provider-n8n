#!/bin/bash
# Configure Git to use hooks from $HOME/.git-hooks/ directory

set -e

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${BLUE}📦 Configuring git hooks...${NC}"

# Check if we're in a git repository
if [ ! -d ".git" ]; then
  echo -e "${YELLOW}⚠️  Not in a git repository. Run this from the project root.${NC}"
  exit 1
fi

# Define hooks directory
HOOKS_DIR="$HOME/.git-hooks"

# Check if hooks directory exists
if [ ! -d "$HOOKS_DIR" ]; then
  echo -e "${YELLOW}⚠️  $HOOKS_DIR directory not found${NC}"
  exit 1
fi

# Configure Git to use $HOME/.git-hooks directory
echo -e "${BLUE}Setting core.hooksPath to $HOOKS_DIR${NC}"
git config core.hooksPath "$HOOKS_DIR"

# Permissions are already set in the Docker image, no need to chmod

echo ""
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}✅ Git hooks configured successfully!${NC}"
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo -e "Hooks enabled from ${BLUE}$HOOKS_DIR${NC} directory:"
echo -e "  📝 ${BLUE}pre-commit${NC}         - Auto-generates CHANGELOG.md and coverage"
echo -e "  💬 ${BLUE}prepare-commit-msg${NC} - Suggests conventional commit format"
echo -e "  ✅ ${BLUE}commit-msg${NC}         - Validates commit message with commitlint"
echo -e "  🚫 ${BLUE}pre-push${NC}           - Blocks push if AI mentions or Co-Authored-By found"
echo ""
echo -e "Next commit will automatically:"
echo -e "  • Generate CHANGELOG.md from git history"
echo -e "  • Run test coverage analysis"
echo -e "  • Validate commit message format (requires commitlint)"
echo ""
echo -e "Before push will check:"
echo -e "  • No AI mentions (Claude, ChatGPT, Copilot, etc.)"
echo -e "  • No Co-Authored-By with bots or AI"
echo ""
echo -e "${YELLOW}Note:${NC} To skip hooks temporarily, use: git commit --no-verify / git push --no-verify"
echo -e "${YELLOW}Note:${NC} commitlint requires: npm install -g @commitlint/cli @commitlint/config-conventional"
echo -e "${YELLOW}Note:${NC} Hooks are embedded in the Docker image with proper permissions"
echo ""
