#!/bin/bash
# Copyright (c) 2024 Florent (Kodflow). All rights reserved.
# Licensed under the Sustainable Use License 1.0
# See LICENSE.md in the project root for license information.
set -e

# Colors for output
CYAN='\033[36m'
GREEN='\033[32m'
YELLOW='\033[33m'
RED='\033[31m'
RESET='\033[0m'

# Copyright header templates
GO_HEADER='// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE.md in the project root for license information.'

SHELL_HEADER='#!/bin/bash
# Copyright (c) 2024 Florent (Kodflow). All rights reserved.
# Licensed under the Sustainable Use License 1.0
# See LICENSE.md in the project root for license information.'

YAML_HEADER='# Copyright (c) 2024 Florent (Kodflow). All rights reserved.
# Licensed under the Sustainable Use License 1.0
# See LICENSE.md in the project root for license information.'

# Files/directories to ignore (regex patterns)
IGNORE_PATTERNS=(
  "^\.git/"
  "^bazel-"
  "^vendor/"
  "^node_modules/"
  "/\..*"                    # Hidden files/dirs
  "\.md$"                    # Markdown files
  "\.txt$"                   # Text files
  "\.json$"                  # JSON files
  "\.lock$"                  # Lock files
  "\.sum$"                   # Checksum files
  "\.mod$"                   # Go module files
  "^LICENSE"                 # License files
  "^CHANGELOG"               # Changelog files
  "^COVERAGE"                # Coverage files
  "\.gitignore$"            # Gitignore files
  "\.prettierignore$"       # Prettier ignore
  "^sdk/n8nsdk/"            # Generated SDK code
  "^codegen/"               # Code generation scripts
  "_test\.go$"              # Test files (already have package declaration)
  "BUILD\.bazel$"           # Bazel BUILD files
  "WORKSPACE$"              # Bazel WORKSPACE
  "\.bzl$"                  # Bazel .bzl files
  "^examples/"              # Example files
  "\.tfvars$"               # Terraform variables
  "\.hcl$"                  # HCL files
  "^\.devcontainer/"        # Devcontainer config
  "^\.github/"              # GitHub workflows (YAML with specific syntax)
  "^scripts/generate-"      # Auto-generated script output
  "\.patch$"                # Patch files
  "\.diff$"                 # Diff files
  "openapi.*\.ya?ml$"       # OpenAPI spec files
)

# Check if file should be ignored
should_ignore() {
  local file="$1"
  for pattern in "${IGNORE_PATTERNS[@]}"; do
    if [[ "$file" =~ $pattern ]]; then
      return 0
    fi
  done
  return 1
}

# Check if file already has copyright header
has_copyright() {
  local file="$1"
  head -n 5 "$file" 2>/dev/null | grep -q "Copyright (c) 2024 Florent (Kodflow)"
}

# Add copyright header to Go file
add_go_header() {
  local file="$1"

  # Check if file starts with package declaration
  if head -n 1 "$file" | grep -q "^package "; then
    # Insert header before package declaration
    echo "$GO_HEADER" > "$file.tmp"
    echo "" >> "$file.tmp"
    cat "$file" >> "$file.tmp"
    mv "$file.tmp" "$file"
  elif head -n 1 "$file" | grep -q "^//go:build "; then
    # Insert header after build constraint
    head -n 1 "$file" > "$file.tmp"
    echo "" >> "$file.tmp"
    echo "$GO_HEADER" >> "$file.tmp"
    echo "" >> "$file.tmp"
    tail -n +2 "$file" >> "$file.tmp"
    mv "$file.tmp" "$file"
  else
    # Insert at beginning
    echo "$GO_HEADER" > "$file.tmp"
    echo "" >> "$file.tmp"
    cat "$file" >> "$file.tmp"
    mv "$file.tmp" "$file"
  fi
}

# Add copyright header to shell script
add_shell_header() {
  local file="$1"

  # Check if file starts with shebang
  if head -n 1 "$file" | grep -q "^#!/"; then
    # Insert header after shebang
    head -n 1 "$file" > "$file.tmp"
    echo "$SHELL_HEADER" | tail -n +2 >> "$file.tmp"
    tail -n +2 "$file" >> "$file.tmp"
    mv "$file.tmp" "$file"
  else
    # Insert at beginning
    echo "$SHELL_HEADER" > "$file.tmp"
    cat "$file" >> "$file.tmp"
    mv "$file.tmp" "$file"
  fi
}

# Add copyright header to YAML/other files
add_yaml_header() {
  local file="$1"
  echo "$YAML_HEADER" > "$file.tmp"
  echo "" >> "$file.tmp"
  cat "$file" >> "$file.tmp"
  mv "$file.tmp" "$file"
}

# Process a single file
process_file() {
  local file="$1"

  # Skip if file doesn't exist or is not a regular file
  if [ ! -f "$file" ]; then
    return 0
  fi

  # Skip if should be ignored
  if should_ignore "$file"; then
    return 0
  fi

  # Skip if already has copyright
  if has_copyright "$file"; then
    return 0
  fi

  # Determine file type and add appropriate header
  case "$file" in
    *.go)
      echo -e "  ${CYAN}→${RESET} Adding header to: $file"
      add_go_header "$file"
      return 1  # File was modified
      ;;
    *.sh)
      echo -e "  ${CYAN}→${RESET} Adding header to: $file"
      add_shell_header "$file"
      return 1  # File was modified
      ;;
    *.yml|*.yaml)
      # Skip GitHub workflows and OpenAPI specs
      if [[ "$file" =~ ^\.github/ ]] || [[ "$file" =~ openapi ]]; then
        return 0
      fi
      echo -e "  ${CYAN}→${RESET} Adding header to: $file"
      add_yaml_header "$file"
      return 1  # File was modified
      ;;
    *)
      return 0  # Unknown file type, skip
      ;;
  esac
}

# Main logic
main() {
  local files_modified=0

  # Get list of staged files or all files if not in git hook mode
  if git rev-parse --git-dir > /dev/null 2>&1 && [ -n "$(git diff --cached --name-only)" ]; then
    # Git hook mode: process staged files
    files=$(git diff --cached --name-only --diff-filter=ACM)
  else
    # Manual mode: process all tracked files
    if git rev-parse --git-dir > /dev/null 2>&1; then
      files=$(git ls-files)
    else
      # Not in a git repo, process all files in current directory
      files=$(find . -type f -not -path '*/\.*' -not -path '*/bazel-*' -not -path '*/vendor/*' -not -path '*/node_modules/*')
    fi
  fi

  echo -e "${CYAN}Checking copyright headers...${RESET}"

  for file in $files; do
    if process_file "$file"; then
      continue
    else
      files_modified=$((files_modified + 1))
      # Re-stage the file if in git hook mode
      if git rev-parse --git-dir > /dev/null 2>&1; then
        git add "$file"
      fi
    fi
  done

  if [ $files_modified -gt 0 ]; then
    echo -e "${GREEN}✓${RESET} Added copyright headers to $files_modified file(s)"
    echo -e "${YELLOW}ℹ${RESET}  Modified files have been re-staged"
  else
    echo -e "${GREEN}✓${RESET} All files already have copyright headers"
  fi
}

# Run main function
main "$@"
