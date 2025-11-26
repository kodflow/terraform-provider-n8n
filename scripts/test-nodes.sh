#!/usr/bin/env bash
# Copyright (c) 2024 Florent (Kodflow). All rights reserved.
# Licensed under the Sustainable Use License 1.0
# See LICENSE in the project root for license information.

# Test all node examples with Terraform plan/apply/destroy
# This script finds all node examples and runs Terraform lifecycle on each
#
# Environment variables:
#   NODES_CATEGORY    - Category within examples/nodes/ (core, trigger, integration)
#   NODES_GROUP       - Specific group from CATEGORIES.json (crm-marketing, dev-devops, etc.)
#   NODES_LIST        - Space-separated list of specific node names to test
#   TEST_NODES_LIMIT  - Limit number of nodes to test (for debugging)

set -euo pipefail

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color
BOLD='\033[1m'

# Counters
TOTAL=0
SUCCESS=0
FAILED=0
declare -a FAILED_NODES=()

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
NODES_DIR="${PROJECT_ROOT}/examples/nodes"

# Load environment variables from .env if it exists
if [ -f "${PROJECT_ROOT}/.env" ]; then
  echo -e "${BLUE}📁 Loading environment from .env${NC}"
  # shellcheck disable=SC1091
  set -a
  source "${PROJECT_ROOT}/.env"
  set +a
fi

# Check required environment variables
if [ -z "${N8N_API_URL:-}" ]; then
  echo -e "${RED}❌ N8N_API_URL is not set${NC}"
  echo "Please set N8N_API_URL in .env or environment"
  exit 1
fi

if [ -z "${N8N_API_KEY:-}" ]; then
  echo -e "${RED}❌ N8N_API_KEY is not set${NC}"
  echo "Please set N8N_API_KEY in .env or environment"
  exit 1
fi

# Check if provider is built
if [ ! -f "${HOME}/.terraform.d/plugins/registry.terraform.io/kodflow/n8n/1.1.0/linux_arm64/terraform-provider-n8n_v1.1.0" ] \
  && [ ! -f "${HOME}/.terraform.d/plugins/registry.terraform.io/kodflow/n8n/1.1.0/linux_amd64/terraform-provider-n8n_v1.1.0" ] \
  && [ ! -f "${HOME}/.terraform.d/plugins/registry.terraform.io/kodflow/n8n/1.1.0/darwin_arm64/terraform-provider-n8n_v1.1.0" ] \
  && [ ! -f "${HOME}/.terraform.d/plugins/registry.terraform.io/kodflow/n8n/1.1.0/darwin_amd64/terraform-provider-n8n_v1.1.0" ]; then
  echo -e "${YELLOW}⚠️  Provider not found, building it first...${NC}"
  cd "${PROJECT_ROOT}"
  make build
fi

echo -e "${BOLD}${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BOLD}  Testing Node Examples${NC}"
echo -e "${BOLD}${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo -e "${BLUE}🌐 N8N Instance:${NC} ${N8N_API_URL}"
echo -e "${BLUE}📦 Testing from:${NC} ${NODES_DIR}"
echo ""

# Determine which nodes to test
declare -a NODE_DIRS=()

# Option 1: Specific list of nodes provided via NODES_LIST
if [ -n "${NODES_LIST:-}" ]; then
  echo -e "${BLUE}📂 Testing specific nodes from NODES_LIST${NC}"
  CATEGORY="${NODES_CATEGORY:-integration}"
  for node in ${NODES_LIST}; do
    NODE_PATH="${NODES_DIR}/${CATEGORY}/${node}"
    if [ -d "${NODE_PATH}" ] && [ -f "${NODE_PATH}/main.tf" ]; then
      NODE_DIRS+=("${NODE_PATH}")
    else
      echo -e "${YELLOW}⚠️  Node not found: ${node}${NC}"
    fi
  done

# Option 2: Group from CATEGORIES.json
elif [ -n "${NODES_GROUP:-}" ]; then
  CATEGORIES_FILE="${NODES_DIR}/integration/CATEGORIES.json"
  if [ ! -f "${CATEGORIES_FILE}" ]; then
    echo -e "${RED}❌ CATEGORIES.json not found at ${CATEGORIES_FILE}${NC}"
    exit 1
  fi

  echo -e "${BLUE}📂 Testing group: ${NODES_GROUP}${NC}"

  # Extract nodes for this group using jq
  if ! command -v jq &>/dev/null; then
    echo -e "${RED}❌ jq is required for NODES_GROUP feature${NC}"
    exit 1
  fi

  GROUP_NODES=$(jq -r ".categories[\"${NODES_GROUP}\"].nodes[]" "${CATEGORIES_FILE}" 2>/dev/null || echo "")
  if [ -z "${GROUP_NODES}" ]; then
    echo -e "${RED}❌ Group '${NODES_GROUP}' not found in CATEGORIES.json${NC}"
    echo "Available groups:"
    jq -r '.categories | keys[]' "${CATEGORIES_FILE}"
    exit 1
  fi

  GROUP_DESC=$(jq -r ".categories[\"${NODES_GROUP}\"].description // \"\"" "${CATEGORIES_FILE}")
  echo -e "${BLUE}📝 Description: ${GROUP_DESC}${NC}"

  for node in ${GROUP_NODES}; do
    NODE_PATH="${NODES_DIR}/integration/${node}"
    if [ -d "${NODE_PATH}" ] && [ -f "${NODE_PATH}/main.tf" ]; then
      NODE_DIRS+=("${NODE_PATH}")
    else
      echo -e "${YELLOW}⚠️  Node not found: ${node}${NC}"
    fi
  done

# Option 3: Category directory (original behavior)
elif [ -n "${NODES_CATEGORY:-}" ]; then
  echo -e "${BLUE}📂 Filtering nodes by category: ${NODES_CATEGORY}${NC}"
  mapfile -t NODE_DIRS < <(find "${NODES_DIR}/${NODES_CATEGORY}" -type f -name "main.tf" -exec dirname {} \; | sort)

# Option 4: All nodes
else
  mapfile -t NODE_DIRS < <(find "${NODES_DIR}" -type f -name "main.tf" -exec dirname {} \; | sort)
fi

echo -e "${BLUE}📊 Found ${#NODE_DIRS[@]} node examples to test${NC}"

# Limit number of nodes to test if TEST_NODES_LIMIT is set
if [ -n "${TEST_NODES_LIMIT:-}" ]; then
  echo -e "${YELLOW}⚠️  TEST_NODES_LIMIT set to ${TEST_NODES_LIMIT} - testing only first ${TEST_NODES_LIMIT} nodes${NC}"
  NODE_DIRS=("${NODE_DIRS[@]:0:$TEST_NODES_LIMIT}")
  echo -e "${BLUE}📊 Will test ${#NODE_DIRS[@]} nodes${NC}"
fi

echo ""

# Function to test a single node
test_node() {
  local node_dir="$1"
  local node_name
  node_name="$(basename "${node_dir}")"
  local category
  category="$(basename "$(dirname "${node_dir}")")"
  local full_name="${category}/${node_name}"

  TOTAL=$((TOTAL + 1))

  echo -e "${BOLD}${BLUE}[${TOTAL}/${#NODE_DIRS[@]}]${NC} Testing: ${full_name}"

  # Create a temporary directory for this test
  local tmp_dir
  tmp_dir=$(mktemp -d)

  # Copy files to temp directory
  cp -r "${node_dir}"/* "${tmp_dir}/"
  cd "${tmp_dir}"

  # Initialize Terraform (suppress output unless error)
  if ! terraform init -no-color >/dev/null 2>&1; then
    echo -e "  ${RED}✗ Init failed${NC}"
    FAILED=$((FAILED + 1))
    FAILED_NODES+=("${full_name} (init)")
    cd "${PROJECT_ROOT}"
    rm -rf "${tmp_dir}"
    return 1
  fi
  echo -e "  ${GREEN}✓${NC} Init"

  # Plan
  if ! terraform plan -no-color \
    -var="n8n_base_url=${N8N_API_URL}" \
    -var="n8n_api_key=${N8N_API_KEY}" \
    -out=tfplan >/dev/null 2>&1; then
    echo -e "  ${RED}✗ Plan failed${NC}"
    FAILED=$((FAILED + 1))
    FAILED_NODES+=("${full_name} (plan)")
    cd "${PROJECT_ROOT}"
    rm -rf "${tmp_dir}"
    return 1
  fi
  echo -e "  ${GREEN}✓${NC} Plan"

  # Apply
  if ! terraform apply -no-color -auto-approve tfplan >/dev/null 2>&1; then
    echo -e "  ${RED}✗ Apply failed${NC}"
    FAILED=$((FAILED + 1))
    FAILED_NODES+=("${full_name} (apply)")
    # Try to cleanup anyway
    terraform destroy -no-color -auto-approve \
      -var="n8n_base_url=${N8N_API_URL}" \
      -var="n8n_api_key=${N8N_API_KEY}" >/dev/null 2>&1 || true
    cd "${PROJECT_ROOT}"
    rm -rf "${tmp_dir}"
    return 1
  fi
  echo -e "  ${GREEN}✓${NC} Apply"

  # Destroy
  if ! terraform destroy -no-color -auto-approve \
    -var="n8n_base_url=${N8N_API_URL}" \
    -var="n8n_api_key=${N8N_API_KEY}" >/dev/null 2>&1; then
    echo -e "  ${RED}✗ Destroy failed${NC}"
    FAILED=$((FAILED + 1))
    FAILED_NODES+=("${full_name} (destroy)")
    cd "${PROJECT_ROOT}"
    rm -rf "${tmp_dir}"
    return 1
  fi
  echo -e "  ${GREEN}✓${NC} Destroy"

  SUCCESS=$((SUCCESS + 1))
  echo -e "  ${GREEN}✓ Success${NC}"

  # Cleanup
  cd "${PROJECT_ROOT}"
  rm -rf "${tmp_dir}"

  return 0
}

# Test all nodes
for node_dir in "${NODE_DIRS[@]}"; do
  test_node "${node_dir}" || true
  echo ""
done

# Summary
echo -e "${BOLD}${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BOLD}  Test Summary${NC}"
echo -e "${BOLD}${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo -e "${BLUE}Total nodes tested:${NC} ${TOTAL}"
echo -e "${GREEN}Successful:${NC} ${SUCCESS}"
echo -e "${RED}Failed:${NC} ${FAILED}"
echo ""

if [ ${FAILED} -gt 0 ]; then
  echo -e "${RED}Failed nodes:${NC}"
  for failed_node in "${FAILED_NODES[@]}"; do
    echo -e "  ${RED}✗${NC} ${failed_node}"
  done
  echo ""
  exit 1
else
  echo -e "${GREEN}✓ All nodes tested successfully!${NC}"
  exit 0
fi
