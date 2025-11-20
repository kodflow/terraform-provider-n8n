#!/bin/bash
# Copyright (c) 2024 Florent (Kodflow). All rights reserved.
# Licensed under the Sustainable Use License 1.0
# See LICENSE in the project root for license information.
# Copyright 2024 Kodflow
# SPDX-License-Identifier: MIT
#
# Update Go version across the entire project
set -e

# Colors
CYAN='\033[36m'
GREEN='\033[32m'
YELLOW='\033[33m'
RESET='\033[0m'
BOLD='\033[1m'

echo ""
echo -e "${BOLD}${CYAN}🔄 Updating Go version across project...${RESET}"
echo ""

# Get latest stable Go version from official API
echo -e "  ${CYAN}→${RESET} Fetching latest Go version..."
LATEST_VERSION=$(curl -s https://go.dev/VERSION?m=text | head -n 1 | sed 's/go//')

if [ -z "$LATEST_VERSION" ]; then
  echo -e "  ${YELLOW}⚠${RESET}  Failed to fetch Go version from go.dev"
  echo -e "  ${YELLOW}ℹ${RESET}  Trying alternative source..."
  LATEST_VERSION=$(curl -s https://golang.org/VERSION?m=text | head -n 1 | sed 's/go//')
fi

if [ -z "$LATEST_VERSION" ]; then
  echo -e "  ${YELLOW}✗${RESET} Failed to fetch latest Go version"
  exit 1
fi

echo -e "  ${GREEN}✓${RESET} Latest Go version: ${BOLD}${LATEST_VERSION}${RESET}"
echo ""

# Update root go.mod
echo -e "  ${CYAN}→${RESET} Updating go.mod..."
if [ -f "go.mod" ]; then
  CURRENT_VERSION=$(grep -oP 'go \K[0-9]+\.[0-9]+(\.[0-9]+)?' go.mod || echo "unknown")
  if [ "$CURRENT_VERSION" != "$LATEST_VERSION" ]; then
    sed -i "s/^go [0-9]\+\.[0-9]\+\(\.[0-9]\+\)\?$/go $LATEST_VERSION/" go.mod
    echo -e "  ${GREEN}✓${RESET} Updated: ${CURRENT_VERSION} → ${LATEST_VERSION}"
  else
    echo -e "  ${GREEN}✓${RESET} Already up to date: ${CURRENT_VERSION}"
  fi
else
  echo -e "  ${YELLOW}⚠${RESET}  go.mod not found"
fi

# Update SDK go.mod
echo -e "  ${CYAN}→${RESET} Updating sdk/n8nsdk/go.mod..."
if [ -f "sdk/n8nsdk/go.mod" ]; then
  CURRENT_VERSION=$(grep -oP 'go \K[0-9]+\.[0-9]+(\.[0-9]+)?' sdk/n8nsdk/go.mod || echo "unknown")
  if [ "$CURRENT_VERSION" != "$LATEST_VERSION" ]; then
    sed -i "s/^go [0-9]\+\.[0-9]\+\(\.[0-9]\+\)\?$/go $LATEST_VERSION/" sdk/n8nsdk/go.mod
    echo -e "  ${GREEN}✓${RESET} Updated: ${CURRENT_VERSION} → ${LATEST_VERSION}"
  else
    echo -e "  ${GREEN}✓${RESET} Already up to date: ${CURRENT_VERSION}"
  fi
else
  echo -e "  ${YELLOW}⚠${RESET}  sdk/n8nsdk/go.mod not found"
fi

# Update devcontainer.json if it specifies a Go version
echo -e "  ${CYAN}→${RESET} Checking .devcontainer/devcontainer.json..."
if [ -f ".devcontainer/devcontainer.json" ]; then
  # The devcontainer uses "latest" by default, which is fine
  # We only update if there's a specific version pinned
  if grep -q '"version".*"[0-9]' .devcontainer/devcontainer.json; then
    echo -e "  ${CYAN}ℹ${RESET}  Specific Go version found in devcontainer.json"
    echo -e "  ${YELLOW}⚠${RESET}  Note: devcontainer.json uses 'latest' by default (recommended)"
  else
    echo -e "  ${GREEN}✓${RESET} Using 'latest' (recommended)"
  fi
else
  echo -e "  ${YELLOW}⚠${RESET}  .devcontainer/devcontainer.json not found"
fi

# Run go mod tidy on all modules
echo ""
echo -e "  ${CYAN}→${RESET} Running go mod tidy on root module..."
if go mod tidy 2>&1; then
  echo -e "  ${GREEN}✓${RESET} Root module dependencies cleaned"
else
  echo -e "  ${YELLOW}⚠${RESET}  Failed to tidy root module"
fi

echo -e "  ${CYAN}→${RESET} Running go mod tidy on SDK module..."
if (cd sdk/n8nsdk && go mod tidy 2>&1); then
  echo -e "  ${GREEN}✓${RESET} SDK module dependencies cleaned"
else
  echo -e "  ${YELLOW}⚠${RESET}  Failed to tidy SDK module"
fi

echo ""
echo -e "${GREEN}✓${RESET} Go version update completed"
echo ""
