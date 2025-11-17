#!/usr/bin/env bash
# Copyright (c) 2024 Florent (Kodflow). All rights reserved.
# Licensed under the Sustainable Use License 1.0
# See LICENSE.md in the project root for license information.
# Copyright 2025 Kodflow
# SPDX-License-Identifier: MIT

# Terraform Examples Validation Script
# Validates all Terraform examples syntax using locally built provider

set -euo pipefail

# Colors
CYAN='\033[36m'
GREEN='\033[32m'
RED='\033[31m'
YELLOW='\033[33m'
RESET='\033[0m'

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
WORKSPACE_DIR="$(dirname "$SCRIPT_DIR")"

# Detect OS and ARCH
GOOS=$(go env GOOS)
GOARCH=$(go env GOARCH)

# Get version from git tags
LAST_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "v1.0.0")
VERSION="${LAST_TAG#v}"

# Provider paths
PROVIDER_BIN="$WORKSPACE_DIR/bazel-bin/src/terraform-provider-n8n_/terraform-provider-n8n"
PROVIDER_DIR="$HOME/.terraform.d/plugins/registry.terraform.io/kodflow/n8n/$VERSION/${GOOS}_${GOARCH}"

# Validate provider binary exists
if [ ! -f "$PROVIDER_BIN" ]; then
    printf "  ${RED}✗${RESET} Provider binary not found at: $PROVIDER_BIN\n"
    printf "  ${CYAN}ℹ${RESET}  Run 'make build' first\n"
    exit 1
fi

# Install provider to plugin directory
printf "  ${CYAN}→${RESET} Installing provider to plugin directory\n"
mkdir -p "$PROVIDER_DIR"
cp "$PROVIDER_BIN" "$PROVIDER_DIR/"
chmod +x "$PROVIDER_DIR/terraform-provider-n8n"

# Setup dev overrides config
TF_CLI_CONFIG_FILE=$(mktemp)
trap "rm -f $TF_CLI_CONFIG_FILE" EXIT

cat > "$TF_CLI_CONFIG_FILE" <<EOF
provider_installation {
  dev_overrides {
    "kodflow/n8n" = "$PROVIDER_DIR"
  }
  direct {}
}
EOF

export TF_CLI_CONFIG_FILE

# Find all example directories
EXAMPLES=$(find "$WORKSPACE_DIR/examples" -name "*.tf" -exec dirname {} \; | sort -u)

FAILED_EXAMPLES=()
PASSED_COUNT=0
TOTAL_COUNT=0

# Validate each example
for example_dir in $EXAMPLES; do
    TOTAL_COUNT=$((TOTAL_COUNT + 1))
    EXAMPLE_NAME=$(echo "$example_dir" | sed "s|$WORKSPACE_DIR/examples/||")

    printf "  ${CYAN}→${RESET} Validating: $EXAMPLE_NAME\n"

    # Clean terraform artifacts
    rm -rf "$example_dir/.terraform" "$example_dir/.terraform.lock.hcl" 2>/dev/null || true

    # Try to validate
    if (cd "$example_dir" && terraform init -no-color -upgrade > /dev/null 2>&1 && terraform validate -no-color > /dev/null 2>&1); then
        printf "    ${GREEN}✓${RESET} Valid\n"
        PASSED_COUNT=$((PASSED_COUNT + 1))
    else
        printf "    ${RED}✗${RESET} Invalid\n"
        # Show error details
        (cd "$example_dir" && terraform init -no-color -upgrade > /dev/null 2>&1 && terraform validate -no-color 2>&1) | sed 's/^/      /'
        FAILED_EXAMPLES+=("$EXAMPLE_NAME")
    fi

    # Cleanup
    rm -rf "$example_dir/.terraform" "$example_dir/.terraform.lock.hcl" 2>/dev/null || true
done

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

if [ ${#FAILED_EXAMPLES[@]} -eq 0 ]; then
    printf "${GREEN}✓${RESET} All $TOTAL_COUNT examples are valid!\n"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    exit 0
else
    printf "${RED}✗${RESET} $PASSED_COUNT/$TOTAL_COUNT examples passed validation\n"
    echo ""
    printf "Failed examples:\n"
    for failed in "${FAILED_EXAMPLES[@]}"; do
        printf "  - $failed\n"
    done
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    exit 1
fi
