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

# Update root module dependencies
echo -e "  ${CYAN}→${RESET} Updating root module dependencies..."
if go get -u ./... 2>&1; then
  echo -e "  ${GREEN}✓${RESET} Root dependencies updated"
else
  echo -e "  ${YELLOW}⚠${RESET}  Failed to update root dependencies"
fi

echo -e "  ${CYAN}→${RESET} Running go mod tidy on root module..."
if go mod tidy 2>&1; then
  echo -e "  ${GREEN}✓${RESET} Root module cleaned"
else
  echo -e "  ${YELLOW}⚠${RESET}  Failed to tidy root module"
fi

# Update SDK module dependencies
echo ""
echo -e "  ${CYAN}→${RESET} Updating SDK module dependencies..."
if (cd sdk/n8nsdk && go get -u ./... 2>&1); then
  echo -e "  ${GREEN}✓${RESET} SDK dependencies updated"
else
  echo -e "  ${YELLOW}⚠${RESET}  Failed to update SDK dependencies"
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
