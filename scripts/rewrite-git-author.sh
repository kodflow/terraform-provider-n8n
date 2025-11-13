#!/bin/bash
# Copyright (c) 2024 Florent (Kodflow). All rights reserved.
# Licensed under the Sustainable Use License 1.0
# See LICENSE.md in the project root for license information.

# Script to rewrite git history to update author email

set -e

# Colors
CYAN='\033[36m'
GREEN='\033[32m'
YELLOW='\033[33m'
RED='\033[31m'
RESET='\033[0m'
BOLD='\033[1m'

# New author information
NEW_NAME="Kodflow"
NEW_EMAIL="133899878+kodflow@users.noreply.github.com"

echo ""
echo -e "${BOLD}${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${RESET}"
echo -e "${BOLD}  Git History Rewrite - Author Update${RESET}"
echo -e "${BOLD}${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${RESET}"
echo ""

# Check if we're in a git repository
if ! git rev-parse --git-dir >/dev/null 2>&1; then
  echo -e "${RED}✗${RESET} Not in a git repository"
  exit 1
fi

# Get current branch
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
echo -e "${CYAN}Current branch:${RESET} $CURRENT_BRANCH"

# Get base branch (usually main or master)
BASE_BRANCH="main"
if ! git rev-parse --verify $BASE_BRANCH >/dev/null 2>&1; then
  BASE_BRANCH="master"
  if ! git rev-parse --verify $BASE_BRANCH >/dev/null 2>&1; then
    echo -e "${RED}✗${RESET} Could not find base branch (main or master)"
    exit 1
  fi
fi
echo -e "${CYAN}Base branch:${RESET} $BASE_BRANCH"
echo ""

# Check if there are uncommitted changes
if ! git diff-index --quiet HEAD --; then
  echo -e "${RED}✗${RESET} You have uncommitted changes. Please commit or stash them first."
  exit 1
fi

# Count commits to rewrite
COMMIT_COUNT=$(git rev-list --count $BASE_BRANCH..$CURRENT_BRANCH 2>/dev/null || echo "0")
if [ "$COMMIT_COUNT" = "0" ]; then
  echo -e "${YELLOW}⚠${RESET}  No commits to rewrite on this branch"
  exit 0
fi

echo -e "${BOLD}Commits to rewrite:${RESET} $COMMIT_COUNT"
echo ""

# Show warning
echo -e "${YELLOW}⚠${RESET}  ${BOLD}WARNING:${RESET} This will rewrite git history!"
echo -e "   All commit SHAs will change on branch: ${BOLD}$CURRENT_BRANCH${RESET}"
echo -e "   New author: ${BOLD}$NEW_NAME <$NEW_EMAIL>${RESET}"
echo ""
echo -e "${YELLOW}⚠${RESET}  If you have already pushed this branch, you'll need to force push!"
echo ""

# Ask for confirmation
read -p "Do you want to continue? (yes/no): " CONFIRM
if [ "$CONFIRM" != "yes" ]; then
  echo -e "${YELLOW}✗${RESET} Operation cancelled"
  exit 0
fi

echo ""
echo -e "${CYAN}→${RESET} Rewriting commit history..."

# Rewrite history using git filter-branch
git filter-branch --force --env-filter '
export GIT_AUTHOR_NAME="'"$NEW_NAME"'"
export GIT_AUTHOR_EMAIL="'"$NEW_EMAIL"'"
export GIT_COMMITTER_NAME="'"$NEW_NAME"'"
export GIT_COMMITTER_EMAIL="'"$NEW_EMAIL"'"
' -- $BASE_BRANCH..$CURRENT_BRANCH

echo ""
echo -e "${GREEN}✓${RESET} Git history rewritten successfully!"
echo ""
echo -e "${BOLD}Next steps:${RESET}"
echo -e "  1. Review the changes: ${CYAN}git log --oneline -10${RESET}"
echo -e "  2. Force push to remote: ${CYAN}git push --force-with-lease${RESET}"
echo ""
echo -e "${YELLOW}⚠${RESET}  Backup refs saved in: ${CYAN}.git/refs/original/${RESET}"
echo -e "   To remove backups: ${CYAN}git update-ref -d refs/original/refs/heads/$CURRENT_BRANCH${RESET}"
echo ""
