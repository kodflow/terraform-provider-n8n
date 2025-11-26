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

# Log function - always goes to stderr
log() {
  echo "[e2e-setup] $*" >&2
}

# Check required environment variables
check_env() {
  if [ -z "${N8N_BASE_URL:-}" ]; then
    log "ERROR: N8N_BASE_URL is not set"
    exit 1
  fi

  if [ -z "${N8N_API_KEY:-}" ]; then
    log "ERROR: N8N_API_KEY is not set"
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
  log "Deleting workflows in project..."

  local workflows
  workflows=$(api_call GET "/workflows?limit=250" | jq -r ".data[] | select(.shared[]?.projectId == \"$project_id\") | .id" 2>/dev/null || echo "")

  if [ -z "$workflows" ]; then
    log "No workflows to delete"
    return
  fi

  local count=0
  for wf_id in $workflows; do
    api_call DELETE "/workflows/${wf_id}" >/dev/null 2>&1 || true
    count=$((count + 1))
  done
  log "Deleted $count workflow(s)"
}

# Delete all credentials in a project
delete_project_credentials() {
  local project_id="$1"
  log "Deleting credentials in project..."

  local credentials
  credentials=$(api_call GET "/credentials" | jq -r ".data[] | select(.sharedWithProjects[]?.id == \"$project_id\") | .id" 2>/dev/null || echo "")

  if [ -z "$credentials" ]; then
    log "No credentials to delete"
    return
  fi

  local count=0
  for cred_id in $credentials; do
    api_call DELETE "/credentials/${cred_id}" >/dev/null 2>&1 || true
    count=$((count + 1))
  done
  log "Deleted $count credential(s)"
}

# Delete all variables in a project (Enterprise only)
delete_project_variables() {
  local project_id="$1"
  log "Deleting variables in project..."

  local variables
  variables=$(api_call GET "/variables" | jq -r ".data[] | select(.project?.id == \"$project_id\") | .id" 2>/dev/null || echo "")

  if [ -z "$variables" ]; then
    log "No variables to delete"
    return
  fi

  local count=0
  for var_id in $variables; do
    api_call DELETE "/variables/${var_id}" >/dev/null 2>&1 || true
    count=$((count + 1))
  done
  log "Deleted $count variable(s)"
}

# Delete a project and all its contents
delete_project() {
  local project_id="$1"
  local project_name="$2"

  log "Cleaning up project: $project_name ($project_id)"

  delete_project_workflows "$project_id"
  delete_project_credentials "$project_id"
  delete_project_variables "$project_id"

  log "Deleting project..."
  api_call DELETE "/projects/${project_id}" >/dev/null 2>&1 || true
  log "Project deleted"
}

# Create a new project - returns project_id on stdout
create_project() {
  local name="$1"

  log "Creating project: $name"

  local response
  response=$(api_call POST "/projects" "{\"name\": \"$name\"}")

  local project_id
  project_id=$(echo "$response" | jq -r '.id' 2>/dev/null || echo "")

  if [ -z "$project_id" ] || [ "$project_id" = "null" ]; then
    log "ERROR: Failed to create project: $response"
    exit 1
  fi

  log "Project created: $project_id"
  # Return project_id on stdout
  echo "$project_id"
}

# Setup action: cleanup existing + create new
setup() {
  log "=========================================="
  log "E2E Project Setup"
  log "=========================================="

  check_env

  log "Checking for existing '$PROJECT_NAME' project..."
  local existing_id
  existing_id=$(find_project "$PROJECT_NAME")

  if [ -n "$existing_id" ]; then
    log "Found existing project: $existing_id"
    delete_project "$existing_id" "$PROJECT_NAME"
  else
    log "No existing project found"
  fi

  local project_id
  project_id=$(create_project "$PROJECT_NAME")

  log "=========================================="
  log "Setup Complete! Project ID: $project_id"
  log "=========================================="

  # Output for GitHub Actions
  if [ -n "${GITHUB_OUTPUT:-}" ]; then
    echo "project_id=$project_id" >>"$GITHUB_OUTPUT"
  fi

  # Plain text output for logging (goes to stdout)
  echo "E2E_PROJECT_ID=$project_id"
}

# Cleanup action: delete project and all contents
cleanup() {
  log "=========================================="
  log "E2E Project Cleanup"
  log "=========================================="

  check_env

  log "Looking for '$PROJECT_NAME' project..."
  local project_id
  project_id=$(find_project "$PROJECT_NAME")

  if [ -z "$project_id" ]; then
    log "No project found to cleanup"
    return 0
  fi

  delete_project "$project_id" "$PROJECT_NAME"

  log "=========================================="
  log "Cleanup Complete!"
  log "=========================================="
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
      echo "Usage: $0 {setup|cleanup}" >&2
      exit 1
      ;;
  esac
}

main "$@"
