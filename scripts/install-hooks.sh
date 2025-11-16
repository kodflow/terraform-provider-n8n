#!/bin/bash
# Copyright (c) 2024 Florent (Kodflow). All rights reserved.
# Licensed under the Sustainable Use License 1.0
# See LICENSE.md in the project root for license information.
# Configure Git to use hooks from .github/hooks/ directory

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

# Define hooks directory (relative to workspace root)
HOOKS_DIR=".github/hooks"

# Check if hooks directory exists
if [ ! -d "$HOOKS_DIR" ]; then
  echo -e "${YELLOW}⚠️  $HOOKS_DIR directory not found${NC}"
  exit 1
fi

# Make all hooks executable
echo -e "${BLUE}Setting executable permissions on hooks...${NC}"
chmod +x "$HOOKS_DIR"/*

# Configure Git to use .github/hooks directory
echo -e "${BLUE}Setting core.hooksPath to $HOOKS_DIR${NC}"
git config core.hooksPath "$HOOKS_DIR"

echo ""
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}✅ Git hooks configured successfully!${NC}"
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo -e "Hooks enabled from ${BLUE}$HOOKS_DIR${NC} directory:"
echo -e "  📝 ${BLUE}pre-commit${NC}         - ${YELLOW}Verifies GPG signatures${NC}, adds headers, generates docs"
echo -e "  💬 ${BLUE}prepare-commit-msg${NC} - Suggests conventional commit format"
echo -e "  ✅ ${BLUE}commit-msg${NC}         - Validates commit message + enforces GPG signature"
echo -e "  🔐 ${BLUE}post-commit${NC}        - Verifies GPG signature on commit"
echo -e "  🚫 ${BLUE}pre-push${NC}           - Blocks push if AI mentions found"
echo ""
echo -e "Next commit will automatically:"
echo -e "  • ${YELLOW}Verify ALL commits in branch are GPG signed${NC}"
echo -e "  • Add copyright headers to modified files"
echo -e "  • Generate CHANGELOG.md from git history"
echo -e "  • Generate COVERAGE.MD from test coverage"
echo -e "  • Unstage sdk/n8nsdk/api/openapi.yaml (auto-generated)"
echo -e "  • Validate commit message (Conventional Commits)"
echo -e "  • ${YELLOW}Enforce GPG signature${NC} (commit will be REJECTED if not signed)"
echo -e "  • Verify GPG signature after commit creation"
echo ""
echo -e "Before push will check:"
echo -e "  • No AI mentions (Claude, ChatGPT, Copilot, etc.)"
echo -e "  • No Co-Authored-By with bots or AI"
echo ""
echo -e "${YELLOW}Note:${NC} To skip hooks temporarily, use: git commit --no-verify / git push --no-verify"
echo -e "${YELLOW}Note:${NC} commitlint requires: npm install -g @commitlint/cli @commitlint/config-conventional"
echo -e "${YELLOW}Note:${NC} Hooks are version-controlled in .github/hooks/ and shared across the team"
echo ""
