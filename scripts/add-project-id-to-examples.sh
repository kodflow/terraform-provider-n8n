#!/usr/bin/env bash
# Copyright (c) 2024 Florent (Kodflow). All rights reserved.
# Licensed under the Sustainable Use License 1.0
# See LICENSE in the project root for license information.

# Script to add project_id variable and attribute to all E2E examples
# This enables E2E test isolation by running all tests in a single project

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
EXAMPLES_DIR="${PROJECT_ROOT}/examples"

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'
BOLD='\033[1m'

echo -e "${BLUE}Adding project_id to all E2E examples...${NC}"

# Project ID variable block to add to variables.tf files
PROJECT_ID_VAR='
variable "project_id" {
  description = "Project ID for E2E test isolation"
  type        = string
  default     = ""
}'

# Count of files modified
modified_vars=0
modified_main=0
skipped=0

# Function to check if file already has project_id variable
has_project_id_var() {
  grep -q 'variable "project_id"' "$1" 2>/dev/null
}

# Function to add project_id variable to variables.tf
add_project_id_var() {
  local file="$1"

  # Skip if already has project_id
  if has_project_id_var "$file"; then
    return 1
  fi

  # Append project_id variable
  echo "$PROJECT_ID_VAR" >>"$file"
  return 0
}

# Function to add project_id to n8n resources using awk
add_project_id_to_resources() {
  local file="$1"

  # Check if file has any n8n resources that need project_id
  if ! grep -qE 'resource "n8n_(workflow|credential|variable)"' "$file"; then
    return 1
  fi

  # Check if already has project_id in any resource block (at resource level)
  if grep -q '^  project_id\s*=' "$file"; then
    return 1
  fi

  # Use awk to add project_id after the name/type line at resource level (2 spaces indent)
  awk '
    BEGIN { in_resource = 0; resource_type = ""; added = 0 }

    # Match start of n8n_workflow, n8n_credential, or n8n_variable resource
    /^resource "n8n_workflow"/ { in_resource = 1; resource_type = "workflow"; added = 0 }
    /^resource "n8n_credential"/ { in_resource = 1; resource_type = "credential"; added = 0 }
    /^resource "n8n_variable"/ { in_resource = 1; resource_type = "variable"; added = 0 }

    # Print current line
    { print }

    # After name line at resource level (2 spaces indent), add project_id once for workflow
    in_resource == 1 && added == 0 && resource_type == "workflow" && /^  name[[:space:]]*=/ {
      print "  project_id = var.project_id != \"\" ? var.project_id : null"
      added = 1
    }

    # After value line at resource level (2 spaces indent), add project_id once for variable
    in_resource == 1 && added == 0 && resource_type == "variable" && /^  value[[:space:]]*=/ {
      print "  project_id = var.project_id != \"\" ? var.project_id : null"
      added = 1
    }

    # After type line at resource level (2 spaces indent), add project_id once
    in_resource == 1 && added == 0 && resource_type == "credential" && /^  type[[:space:]]*=/ {
      print "  project_id = var.project_id != \"\" ? var.project_id : null"
      added = 1
    }

    # End of resource block (line starts with })
    /^}$/ { in_resource = 0; resource_type = "" }
  ' "$file" >"${file}.tmp"

  mv "${file}.tmp" "$file"
  return 0
}

# Process all variables.tf files
echo -e "${YELLOW}Processing variables.tf files...${NC}"
while IFS= read -r file; do
  # Skip enterprise/projects (exception - projects can't be in projects)
  if [[ "$file" == *"enterprise/projects"* ]]; then
    echo -e "  ${YELLOW}Skipping (project exception):${NC} $file"
    skipped=$((skipped + 1))
    continue
  fi

  if add_project_id_var "$file"; then
    echo -e "  ${GREEN}Added variable:${NC} $file"
    modified_vars=$((modified_vars + 1))
  fi
done < <(find "$EXAMPLES_DIR" -name "variables.tf" -type f)

# Process all main.tf files
echo -e "${YELLOW}Processing main.tf files...${NC}"
while IFS= read -r file; do
  # Skip enterprise/projects (exception - projects can't be in projects)
  if [[ "$file" == *"enterprise/projects"* ]]; then
    echo -e "  ${YELLOW}Skipping (project exception):${NC} $file"
    continue
  fi

  if add_project_id_to_resources "$file"; then
    echo -e "  ${GREEN}Modified:${NC} $file"
    modified_main=$((modified_main + 1))
  fi
done < <(find "$EXAMPLES_DIR" -name "main.tf" -type f)

echo ""
echo -e "${BOLD}${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}Summary:${NC}"
echo -e "  Variables.tf files modified: ${modified_vars}"
echo -e "  Main.tf files modified: ${modified_main}"
echo -e "  Files skipped: ${skipped}"
echo -e "${BOLD}${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
