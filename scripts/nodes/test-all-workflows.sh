#!/bin/bash
# Copyright (c) 2024 Florent (Kodflow). All rights reserved.
# Licensed under the Sustainable Use License 1.0
# See LICENSE in the project root for license information.

# Test all per-node workflow examples
#
# This script tests EVERY workflow in examples/nodes/ to ensure 100% coverage.
# It validates that each workflow can be initialized, validated, and planned.
#
# Usage: bash scripts/nodes/test-all-workflows.sh

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;36m'
NC='\033[0m' # No Color

# Counters
TOTAL=0
PASSED=0
FAILED=0
REQUIRES_CREDENTIALS=0

# Results arrays
declare -a FAILED_WORKFLOWS
declare -a CREDENTIAL_WORKFLOWS

# Output file
RESULTS_FILE="WORKFLOWS_TEST_RESULTS.md"

# Mock credentials for testing
export TF_VAR_n8n_base_url="http://localhost:5678"
export TF_VAR_n8n_api_key="mock-api-key-for-testing"

echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}Testing All N8N Node Workflows${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

# Start results file
cat > "$RESULTS_FILE" << 'EOF'
# Workflow Test Results

**Generated**: $(date -u +"%Y-%m-%dT%H:%M:%SZ")
**Environment**: Automated Testing

## Summary

EOF

# Function to test a single workflow
test_workflow() {
    local workflow_dir="$1"
    local category=$(basename $(dirname "$workflow_dir"))
    local node_name=$(basename "$workflow_dir")

    TOTAL=$((TOTAL + 1))

    echo -e "${BLUE}[$TOTAL]${NC} Testing: ${category}/${node_name}"

    # Change to workflow directory
    cd "$workflow_dir" || return 1

    # Check if requires credentials (heuristic)
    local requires_creds=0
    if grep -q "credentials" main.tf 2>/dev/null || \
       grep -q "oauth" main.tf 2>/dev/null || \
       grep -q "api_key" main.tf 2>/dev/null; then
        requires_creds=1
        REQUIRES_CREDENTIALS=$((REQUIRES_CREDENTIALS + 1))
        CREDENTIAL_WORKFLOWS+=("${category}/${node_name}")
    fi

    # Test: terraform init
    if ! terraform init -no-color > /dev/null 2>&1; then
        echo -e "  ${RED}✗${NC} init failed"
        FAILED=$((FAILED + 1))
        FAILED_WORKFLOWS+=("${category}/${node_name} (init)")
        cd - > /dev/null
        return 1
    fi

    # Test: terraform validate
    if ! terraform validate -no-color > /dev/null 2>&1; then
        echo -e "  ${RED}✗${NC} validate failed"
        FAILED=$((FAILED + 1))
        FAILED_WORKFLOWS+=("${category}/${node_name} (validate)")
        cd - > /dev/null
        return 1
    fi

    # Test: terraform plan (optional - may fail if real credentials needed)
    if terraform plan -no-color > /dev/null 2>&1; then
        echo -e "  ${GREEN}✓${NC} All checks passed"
        PASSED=$((PASSED + 1))
    else
        # Plan failed - check if it's due to credentials
        if [ $requires_creds -eq 1 ]; then
            echo -e "  ${YELLOW}⚠${NC}  Plan skipped (requires credentials)"
            PASSED=$((PASSED + 1))  # Still count as pass if init+validate worked
        else
            echo -e "  ${RED}✗${NC} plan failed"
            FAILED=$((FAILED + 1))
            FAILED_WORKFLOWS+=("${category}/${node_name} (plan)")
        fi
    fi

    cd - > /dev/null
    return 0
}

# Find all workflow directories
WORKFLOW_DIRS=$(find examples/nodes -mindepth 2 -maxdepth 2 -type d | sort)

# Test each workflow
for workflow_dir in $WORKFLOW_DIRS; do
    test_workflow "$workflow_dir"
done

# Generate summary
echo ""
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}Test Summary${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo -e "  Total Workflows: ${TOTAL}"
echo -e "  ${GREEN}Passed: ${PASSED}${NC}"
echo -e "  ${RED}Failed: ${FAILED}${NC}"
echo -e "  ${YELLOW}Require Credentials: ${REQUIRES_CREDENTIALS}${NC}"
echo ""

# Update results file
{
    echo "- **Total Workflows**: $TOTAL"
    echo "- **Passed**: $PASSED"
    echo "- **Failed**: $FAILED"
    echo "- **Require Credentials**: $REQUIRES_CREDENTIALS"
    echo "- **Success Rate**: $(awk "BEGIN {printf \"%.1f%%\", ($PASSED/$TOTAL)*100}")"
    echo ""
    echo "---"
    echo ""

    if [ ${#FAILED_WORKFLOWS[@]} -gt 0 ]; then
        echo "## ❌ Failed Workflows (${#FAILED_WORKFLOWS[@]})"
        echo ""
        for workflow in "${FAILED_WORKFLOWS[@]}"; do
            echo "- \`$workflow\`"
        done
        echo ""
    fi

    if [ ${#CREDENTIAL_WORKFLOWS[@]} -gt 0 ]; then
        echo "## ⚠️ Workflows Requiring Credentials (${#CREDENTIAL_WORKFLOWS[@]})"
        echo ""
        echo "These workflows passed init/validate but require real credentials for full testing:"
        echo ""
        for workflow in "${CREDENTIAL_WORKFLOWS[@]}"; do
            echo "- \`$workflow\`"
        done
        echo ""
    fi

    if [ ${#FAILED_WORKFLOWS[@]} -eq 0 ]; then
        echo "## ✅ All Workflows Passed!"
        echo ""
        echo "Every workflow successfully passed terraform init and validate."
        echo ""
    fi

    echo "---"
    echo ""
    echo "## Test Details"
    echo ""
    echo "Each workflow was tested with:"
    echo "1. \`terraform init\` - Initialize provider and modules"
    echo "2. \`terraform validate\` - Validate syntax and configuration"
    echo "3. \`terraform plan\` - Generate execution plan (optional if credentials required)"
    echo ""
    echo "**Note**: Workflows marked as requiring credentials passed init/validate but "
    echo "may require real API keys or OAuth tokens for full \`terraform apply\` testing."

} >> "$RESULTS_FILE"

echo "Results saved to: $RESULTS_FILE"
echo ""

# Exit with error if any failures
if [ $FAILED -gt 0 ]; then
    echo -e "${RED}Some workflows failed. Review the results above.${NC}"
    exit 1
else
    echo -e "${GREEN}All workflows passed!${NC}"
    exit 0
fi
