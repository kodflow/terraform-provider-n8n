#!/usr/bin/env bash
# Copyright (c) 2024 Florent (Kodflow). All rights reserved.
# Licensed under the Sustainable Use License 1.0
# See LICENSE in the project root for license information.

# E2E Project Setup/Cleanup Script
# This script manages the "Provider Terraform" project for E2E tests isolation
#
# Actions:
#   setup   - Cleanup existing project (if any) and create fresh one
#   cleanup - Delete the project and all its contents
#
# Environment variables:
#   N8N_BASE_URL - n8n instance URL
#   N8N_API_KEY  - n8n API key
#
# Outputs (setup only):
#   PROJECT_ID - The ID of the created project (for GitHub Actions)

set -euo pipefail

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color
BOLD='\033[1m'

# Project name
PROJECT_NAME="Provider Terraform"

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

# Load environment variables from .env if it exists
if [ -f "${PROJECT_ROOT}/.env" ]; then
  set -a
  # shellcheck disable=SC1091
  source "${PROJECT_ROOT}/.env"
  set +a
fi

# Check required environment variables
check_env() {
  if [ -z "${N8N_BASE_URL:-}" ]; then
    echo -e "${RED}N8N_BASE_URL is not set${NC}"
    exit 1
  fi

  if [ -z "${N8N_API_KEY:-}" ]; then
    echo -e "${RED}N8N_API_KEY is not set${NC}"
    exit 1
  fi
}

# API helper function
api_call() {
  local method="$1"
  local endpoint="$2"
  local data="${3:-}"

  local url="${N8N_BASE_URL}/api/v1${endpoint}"

  if [ -n "$data" ]; then
    curl -s -X "$method" \
      -H "X-N8N-API-KEY: ${N8N_API_KEY}" \
      -H "Content-Type: application/json" \
      -d "$data" \
      "$url"
  else
    curl -s -X "$method" \
      -H "X-N8N-API-KEY: ${N8N_API_KEY}" \
      "$url"
  fi
}

# Find project by name
find_project() {
  local name="$1"
  api_call GET "/projects" | jq -r ".data[] | select(.name == \"$name\") | .id" 2>/dev/null || echo ""
}

# Delete all workflows in a project
delete_project_workflows() {
  local project_id="$1"
  echo -e "${BLUE}  Deleting workflows in project...${NC}"

  # Get all workflows
  local workflows
  workflows=$(api_call GET "/workflows?limit=250" | jq -r ".data[] | select(.shared[]?.projectId == \"$project_id\") | .id" 2>/dev/null || echo "")

  if [ -z "$workflows" ]; then
    echo -e "${GREEN}    No workflows to delete${NC}"
    return
  fi

  local count=0
  for wf_id in $workflows; do
    api_call DELETE "/workflows/${wf_id}" >/dev/null 2>&1 || true
    count=$((count + 1))
  done
  echo -e "${GREEN}    Deleted $count workflow(s)${NC}"
}

# Delete all credentials in a project
delete_project_credentials() {
  local project_id="$1"
  echo -e "${BLUE}  Deleting credentials in project...${NC}"

  # Get all credentials
  local credentials
  credentials=$(api_call GET "/credentials" | jq -r ".data[] | select(.sharedWithProjects[]?.id == \"$project_id\") | .id" 2>/dev/null || echo "")

  if [ -z "$credentials" ]; then
    echo -e "${GREEN}    No credentials to delete${NC}"
    return
  fi

  local count=0
  for cred_id in $credentials; do
    api_call DELETE "/credentials/${cred_id}" >/dev/null 2>&1 || true
    count=$((count + 1))
  done
  echo -e "${GREEN}    Deleted $count credential(s)${NC}"
}

# Delete all variables in a project (Enterprise only)
delete_project_variables() {
  local project_id="$1"
  echo -e "${BLUE}  Deleting variables in project...${NC}"

  # Get all variables
  local variables
  variables=$(api_call GET "/variables" | jq -r ".data[] | select(.project?.id == \"$project_id\") | .id" 2>/dev/null || echo "")

  if [ -z "$variables" ]; then
    echo -e "${GREEN}    No variables to delete${NC}"
    return
  fi

  local count=0
  for var_id in $variables; do
    api_call DELETE "/variables/${var_id}" >/dev/null 2>&1 || true
    count=$((count + 1))
  done
  echo -e "${GREEN}    Deleted $count variable(s)${NC}"
}

# Delete a project and all its contents
delete_project() {
  local project_id="$1"
  local project_name="$2"

  echo -e "${YELLOW}Cleaning up project: $project_name ($project_id)${NC}"

  # Delete all resources in the project
  delete_project_workflows "$project_id"
  delete_project_credentials "$project_id"
  delete_project_variables "$project_id"

  # Delete the project itself
  echo -e "${BLUE}  Deleting project...${NC}"
  api_call DELETE "/projects/${project_id}" >/dev/null 2>&1 || true
  echo -e "${GREEN}    Project deleted${NC}"
}

# Create a new project
create_project() {
  local name="$1"

  echo -e "${BLUE}Creating project: $name${NC}"

  local response
  response=$(api_call POST "/projects" "{\"name\": \"$name\"}")

  local project_id
  project_id=$(echo "$response" | jq -r '.id' 2>/dev/null || echo "")

  if [ -z "$project_id" ] || [ "$project_id" = "null" ]; then
    echo -e "${RED}Failed to create project: $response${NC}"
    exit 1
  fi

  echo -e "${GREEN}  Project created: $project_id${NC}"
  echo "$project_id"
}

# Setup action: cleanup existing + create new
setup() {
  echo -e "${BOLD}${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
  echo -e "${BOLD}  E2E Project Setup${NC}"
  echo -e "${BOLD}${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
  echo ""

  check_env

  # Check for existing project
  echo -e "${BLUE}Checking for existing '$PROJECT_NAME' project...${NC}"
  local existing_id
  existing_id=$(find_project "$PROJECT_NAME")

  if [ -n "$existing_id" ]; then
    echo -e "${YELLOW}Found existing project: $existing_id${NC}"
    delete_project "$existing_id" "$PROJECT_NAME"
  else
    echo -e "${GREEN}No existing project found${NC}"
  fi

  echo ""

  # Create new project
  local project_id
  project_id=$(create_project "$PROJECT_NAME")

  echo ""
  echo -e "${BOLD}${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
  echo -e "${BOLD}${GREEN}  Setup Complete!${NC}"
  echo -e "${BOLD}${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
  echo ""
  echo -e "${BLUE}Project ID:${NC} $project_id"
  echo ""

  # Output for GitHub Actions
  if [ -n "${GITHUB_OUTPUT:-}" ]; then
    echo "project_id=$project_id" >>"$GITHUB_OUTPUT"
  fi

  # Also output to stdout for capture
  echo "E2E_PROJECT_ID=$project_id"
}

# Cleanup action: delete project and all contents
cleanup() {
  echo -e "${BOLD}${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
  echo -e "${BOLD}  E2E Project Cleanup${NC}"
  echo -e "${BOLD}${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
  echo ""

  check_env

  # Find project
  echo -e "${BLUE}Looking for '$PROJECT_NAME' project...${NC}"
  local project_id
  project_id=$(find_project "$PROJECT_NAME")

  if [ -z "$project_id" ]; then
    echo -e "${GREEN}No project found to cleanup${NC}"
    return 0
  fi

  delete_project "$project_id" "$PROJECT_NAME"

  echo ""
  echo -e "${BOLD}${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
  echo -e "${BOLD}${GREEN}  Cleanup Complete!${NC}"
  echo -e "${BOLD}${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
}

# Main entry point
main() {
  local action="${1:-setup}"

  case "$action" in
    setup)
      setup
      ;;
    cleanup)
      cleanup
      ;;
    *)
      echo "Usage: $0 {setup|cleanup}"
      exit 1
      ;;
  esac
}

main "$@"
