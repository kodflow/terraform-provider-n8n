#!/bin/bash
# Copyright (c) 2024 Florent (Kodflow). All rights reserved.
# Licensed under the Sustainable Use License 1.0
# See LICENSE in the project root for license information.

set -e

# Colors
CYAN='\033[36m'
GREEN='\033[32m'
YELLOW='\033[33m'
RESET='\033[0m'

echo -e "${CYAN}🔄 Updating n8n version badge in README.md...${RESET}"

# Extract n8n version from codegen/download-only.py
N8N_VERSION=$(grep -oP '# Frozen commit for API stability \(n8n@\K[\d.]+' codegen/download-only.py | head -1)

if [ -z "$N8N_VERSION" ]; then
  echo -e "${YELLOW}⚠${RESET}  Could not extract n8n version from codegen/download-only.py"
  exit 1
fi

echo -e "  ${CYAN}→${RESET} Detected n8n version: ${N8N_VERSION}"

# Update README.md badge
if grep -q "badge/n8n-.*-EA4B71" README.md; then
  sed -i "s|badge/n8n-[^-]*-EA4B71|badge/n8n-${N8N_VERSION}-EA4B71|g" README.md
  echo -e "${GREEN}✓${RESET} README.md badge updated to n8n ${N8N_VERSION}"
else
  echo -e "${YELLOW}⚠${RESET}  n8n badge not found in README.md"
  exit 1
fi
