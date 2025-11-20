#!/bin/bash
# Copyright (c) 2024 Florent (Kodflow). All rights reserved.
# Licensed under the Sustainable Use License 1.0
# See LICENSE in the project root for license information.
# Copyright 2024 Kodflow
# SPDX-License-Identifier: MIT
#
# Update Go dependencies across the entire project
set -e

# Colors
CYAN='\033[36m'
GREEN='\033[32m'
YELLOW='\033[33m'
RESET='\033[0m'
BOLD='\033[1m'

echo ""
echo -e "${BOLD}${CYAN}📦 Updating Go dependencies...${RESET}"
echo ""

# Save current Go version before updating dependencies
ROOT_GO_VERSION=$(grep -oP 'go \K[0-9]+\.[0-9]+(\.[0-9]+)?' go.mod 2>/dev/null || echo "")

# Update root module dependencies
echo -e "  ${CYAN}→${RESET} Updating root module dependencies..."
if go get -u ./... 2>&1; then
  echo -e "  ${GREEN}✓${RESET} Root dependencies updated"
else
  echo -e "  ${YELLOW}⚠${RESET}  Failed to update root dependencies"
fi

# Restore Go version if it was changed by go get
if [ -n "$ROOT_GO_VERSION" ]; then
  CURRENT_GO_VERSION=$(grep -oP 'go \K[0-9]+\.[0-9]+(\.[0-9]+)?' go.mod 2>/dev/null || echo "")
  if [ "$CURRENT_GO_VERSION" != "$ROOT_GO_VERSION" ]; then
    sed -i "s/^go [0-9]\+\.[0-9]\+\(\.[0-9]\+\)\?$/go $ROOT_GO_VERSION/" go.mod
    echo -e "  ${CYAN}ℹ${RESET}  Restored Go version to $ROOT_GO_VERSION"
  fi
fi

echo -e "  ${CYAN}→${RESET} Running go mod tidy on root module..."
if go mod tidy 2>&1; then
  echo -e "  ${GREEN}✓${RESET} Root module cleaned"
else
  echo -e "  ${YELLOW}⚠${RESET}  Failed to tidy root module"
fi

# Update SDK module dependencies
echo ""
# Save current Go version before updating SDK dependencies
SDK_GO_VERSION=$(grep -oP 'go \K[0-9]+\.[0-9]+(\.[0-9]+)?' sdk/n8nsdk/go.mod 2>/dev/null || echo "")

echo -e "  ${CYAN}→${RESET} Updating SDK module dependencies..."
if (cd sdk/n8nsdk && go get -u ./... 2>&1); then
  echo -e "  ${GREEN}✓${RESET} SDK dependencies updated"
else
  echo -e "  ${YELLOW}⚠${RESET}  Failed to update SDK dependencies"
fi

# Restore Go version if it was changed by go get
if [ -n "$SDK_GO_VERSION" ]; then
  CURRENT_SDK_GO_VERSION=$(grep -oP 'go \K[0-9]+\.[0-9]+(\.[0-9]+)?' sdk/n8nsdk/go.mod 2>/dev/null || echo "")
  if [ "$CURRENT_SDK_GO_VERSION" != "$SDK_GO_VERSION" ]; then
    sed -i "s/^go [0-9]\+\.[0-9]\+\(\.[0-9]\+\)\?$/go $SDK_GO_VERSION/" sdk/n8nsdk/go.mod
    echo -e "  ${CYAN}ℹ${RESET}  Restored SDK Go version to $SDK_GO_VERSION"
  fi
fi

echo -e "  ${CYAN}→${RESET} Running go mod tidy on SDK module..."
if (cd sdk/n8nsdk && go mod tidy 2>&1); then
  echo -e "  ${GREEN}✓${RESET} SDK module cleaned"
else
  echo -e "  ${YELLOW}⚠${RESET}  Failed to tidy SDK module"
fi

echo ""
echo -e "${GREEN}✓${RESET} Go dependencies update completed"
echo ""
