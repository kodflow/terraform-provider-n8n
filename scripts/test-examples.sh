#!/usr/bin/env bash
# Test all Terraform examples against n8n server
# Usage: ./scripts/test-examples.sh [community|enterprise|all]

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
N8N_API_URL="${N8N_API_URL:-http://localhost:5678}"
N8N_API_KEY="${N8N_API_KEY:-}"

# Check prerequisites
check_prerequisites() {
  echo -e "${BLUE}Checking prerequisites...${NC}"

  # Check if terraform is installed
  if ! command -v terraform &>/dev/null; then
    echo -e "${RED}❌ Terraform not found. Please install terraform or opentofu.${NC}"
    exit 1
  fi

  # Check if N8N_API_KEY is set
  if [ -z "$N8N_API_KEY" ]; then
    echo -e "${RED}❌ N8N_API_KEY environment variable is not set.${NC}"
    echo "   Set it with: export N8N_API_KEY=your-api-key"
    exit 1
  fi

  # Check if n8n server is accessible
  if ! curl -s -f -H "X-N8N-API-KEY: $N8N_API_KEY" "$N8N_API_URL/api/v1/workflows" >/dev/null 2>&1; then
    echo -e "${YELLOW}⚠️  Warning: Cannot connect to n8n server at $N8N_API_URL${NC}"
    echo "   Make sure your n8n instance is running."
    read -p "Continue anyway? (y/n) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
      exit 1
    fi
  fi

  echo -e "${GREEN}✓ Prerequisites check passed${NC}\n"
}

# Test a single example
test_example() {
  local example_path=$1
  local example_name=$(basename "$example_path")
  local example_category=$(basename "$(dirname "$example_path")")

  echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
  echo -e "${BLUE}Testing: $example_category/$example_name${NC}"
  echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"

  cd "$example_path"

  # Export variables
  export TF_VAR_n8n_api_url="$N8N_API_URL"
  export TF_VAR_n8n_api_key="$N8N_API_KEY"

  # Clean any existing state
  rm -rf .terraform .terraform.lock.hcl terraform.tfstate terraform.tfstate.backup 2>/dev/null || true

  # Initialize
  echo -e "\n${YELLOW}→ terraform init${NC}"
  if ! terraform init -no-color >/dev/null 2>&1; then
    echo -e "${RED}❌ Init failed${NC}"
    return 1
  fi
  echo -e "${GREEN}✓ Init successful${NC}"

  # Plan
  echo -e "\n${YELLOW}→ terraform plan${NC}"
  if ! terraform plan -no-color -out=tfplan >/dev/null 2>&1; then
    echo -e "${RED}❌ Plan failed${NC}"
    return 1
  fi
  echo -e "${GREEN}✓ Plan successful${NC}"

  # Apply
  echo -e "\n${YELLOW}→ terraform apply${NC}"
  if ! terraform apply -no-color -auto-approve tfplan >/dev/null 2>&1; then
    echo -e "${RED}❌ Apply failed${NC}"
    # Try to destroy partial state
    terraform destroy -auto-approve >/dev/null 2>&1 || true
    return 1
  fi
  echo -e "${GREEN}✓ Apply successful${NC}"

  # Show outputs
  echo -e "\n${YELLOW}→ terraform output${NC}"
  terraform output -no-color

  # Destroy
  echo -e "\n${YELLOW}→ terraform destroy${NC}"
  if ! terraform destroy -no-color -auto-approve >/dev/null 2>&1; then
    echo -e "${RED}❌ Destroy failed${NC}"
    return 1
  fi
  echo -e "${GREEN}✓ Destroy successful${NC}"

  # Clean up
  rm -rf .terraform .terraform.lock.hcl terraform.tfstate terraform.tfstate.backup tfplan 2>/dev/null || true

  echo -e "\n${GREEN}✅ Example test passed: $example_category/$example_name${NC}\n"
  return 0
}

# Test all examples in a directory
test_examples_in_dir() {
  local base_dir=$1
  local category=$2

  echo -e "${BLUE}═══════════════════════════════════════════════════${NC}"
  echo -e "${BLUE}Testing $category Edition Examples${NC}"
  echo -e "${BLUE}═══════════════════════════════════════════════════${NC}\n"

  local passed=0
  local failed=0
  local failed_examples=()

  # Find all example directories (those containing main.tf)
  while IFS= read -r -d '' example_path; do
    if test_example "$(dirname "$example_path")"; then
      ((passed++))
    else
      ((failed++))
      failed_examples+=("$(basename "$(dirname "$(dirname "$example_path")")")/$(basename "$(dirname "$example_path")")")
    fi
  done < <(find "$base_dir" -name "main.tf" -print0 | sort -z)

  echo -e "${BLUE}═══════════════════════════════════════════════════${NC}"
  echo -e "${BLUE}$category Edition Results${NC}"
  echo -e "${GREEN}✓ Passed: $passed${NC}"
  if [ $failed -gt 0 ]; then
    echo -e "${RED}✗ Failed: $failed${NC}"
    echo -e "${RED}Failed examples:${NC}"
    for ex in "${failed_examples[@]}"; do
      echo -e "${RED}  - $ex${NC}"
    done
  fi
  echo -e "${BLUE}═══════════════════════════════════════════════════${NC}\n"

  return $failed
}

# Main
main() {
  local mode="${1:-all}"
  local workspace_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
  local examples_dir="$workspace_root/examples"

  check_prerequisites

  local total_failed=0

  case $mode in
    community)
      test_examples_in_dir "$examples_dir/community" "Community"
      total_failed=$?
      ;;
    enterprise)
      test_examples_in_dir "$examples_dir/enterprise" "Enterprise"
      total_failed=$?
      ;;
    all)
      test_examples_in_dir "$examples_dir/community" "Community"
      community_failed=$?

      test_examples_in_dir "$examples_dir/enterprise" "Enterprise"
      enterprise_failed=$?

      total_failed=$((community_failed + enterprise_failed))
      ;;
    *)
      echo -e "${RED}Invalid mode: $mode${NC}"
      echo "Usage: $0 [community|enterprise|all]"
      exit 1
      ;;
  esac

  echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
  if [ $total_failed -eq 0 ]; then
    echo -e "${GREEN}🎉 All tests passed!${NC}"
  else
    echo -e "${RED}❌ Some tests failed (total: $total_failed)${NC}"
  fi
  echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"

  exit $total_failed
}

main "$@"
