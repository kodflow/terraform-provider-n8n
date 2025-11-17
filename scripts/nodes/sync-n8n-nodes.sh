#!/usr/bin/env bash
# Copyright (c) 2024 Florent (Kodflow). All rights reserved.
# Licensed under the Sustainable Use License 1.0
# See LICENSE in the project root for license information.

set -euo pipefail

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;36m'
NC='\033[0m' # No Color

# Directories
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/../.." && pwd)"
CACHE_DIR="${ROOT_DIR}/.n8n-repo-cache"
DATA_DIR="${ROOT_DIR}/data"
N8N_REPO="https://github.com/n8n-io/n8n.git"

# Data files
REGISTRY_FILE="${DATA_DIR}/n8n-nodes-registry.json"
METADATA_FILE="${DATA_DIR}/n8n-nodes-metadata.json"
VERSION_FILE="${DATA_DIR}/n8n-nodes-version.txt"
CHANGELOG_FILE="${DATA_DIR}/n8n-nodes-changelog.md"

# Logging functions
log_info() {
  echo -e "${BLUE}→${NC} $1"
}

log_success() {
  echo -e "${GREEN}✓${NC} $1"
}

log_warn() {
  echo -e "${YELLOW}⚠${NC} $1"
}

log_error() {
  echo -e "${RED}✗${NC} $1"
}

# Create necessary directories
mkdir -p "${DATA_DIR}" "${CACHE_DIR}"

# Fetch n8n repository
fetch_n8n_repo() {
  log_info "Fetching n8n repository..."

  if [ -d "${CACHE_DIR}/.git" ]; then
    log_info "Updating existing repository..."
    cd "${CACHE_DIR}"
    git fetch --depth 1 origin master
    git reset --hard origin/master
  else
    log_info "Cloning n8n repository (shallow clone)..."
    git clone --depth 1 --single-branch --branch master "${N8N_REPO}" "${CACHE_DIR}"
  fi

  # Get version
  cd "${CACHE_DIR}"
  local version
  version=$(git describe --tags --abbrev=0 2>/dev/null || echo "unknown")
  echo "${version}" >"${VERSION_FILE}"

  log_success "Repository fetched successfully (version: ${version})"
}

# Parse nodes using Node.js script
parse_nodes() {
  log_info "Parsing nodes from repository..."

  if [ ! -f "${SCRIPT_DIR}/parse-nodes.js" ]; then
    log_error "parse-nodes.js not found!"
    exit 1
  fi

  # Run Node.js parser
  cd "${SCRIPT_DIR}"
  node parse-nodes.js "${CACHE_DIR}" "${DATA_DIR}"

  log_success "Nodes parsed successfully"
}

# Generate diff between versions
generate_diff() {
  log_info "Generating changelog..."

  if [ ! -f "${REGISTRY_FILE}" ]; then
    log_warn "No previous registry found, skipping diff"
    return
  fi

  # Run diff script
  if [ -f "${SCRIPT_DIR}/generate-diff.js" ]; then
    node "${SCRIPT_DIR}/generate-diff.js" "${DATA_DIR}"
    log_success "Changelog generated"
  else
    log_warn "Diff generator not found, skipping"
  fi
}

# Generate Go code
generate_code() {
  log_info "Generating Go code..."

  if [ ! -f "${ROOT_DIR}/codegen/nodes/generator.go" ]; then
    log_warn "Go generator not found yet, skipping code generation"
    return
  fi

  cd "${ROOT_DIR}"
  go run ./codegen/nodes/generator.go

  log_success "Go code generated"
}

# Generate examples
generate_examples() {
  log_info "Generating Terraform examples..."

  if [ ! -f "${SCRIPT_DIR}/generate-examples.js" ]; then
    log_warn "Example generator not found yet, skipping"
    return
  fi

  node "${SCRIPT_DIR}/generate-examples.js" "${DATA_DIR}" "${ROOT_DIR}/examples/nodes"

  log_success "Examples generated"
}

# Display statistics
show_stats() {
  if [ ! -f "${METADATA_FILE}" ]; then
    log_warn "No metadata file found"
    return
  fi

  log_info "Node Statistics:"

  local total
  total=$(jq -r '.total_nodes' "${METADATA_FILE}")
  echo "  Total Nodes: ${total}"

  local categories
  categories=$(jq -r '.categories | to_entries[] | "  - \(.key): \(.value)"' "${METADATA_FILE}")
  echo "${categories}"
}

# Main command handler
case "${1:-help}" in
  fetch)
    fetch_n8n_repo
    ;;
  parse)
    parse_nodes
    ;;
  diff)
    generate_diff
    ;;
  generate)
    generate_code
    generate_examples
    ;;
  stats)
    show_stats
    ;;
  all)
    echo -e "${GREEN}[1m╔════════════════════════════════════════╗${NC}"
    echo -e "${GREEN}[1m║  N8N Nodes Synchronization             ║${NC}"
    echo -e "${GREEN}[1m╚════════════════════════════════════════╝${NC}"
    echo ""
    fetch_n8n_repo
    echo ""
    parse_nodes
    echo ""
    generate_diff
    echo ""
    show_stats
    echo ""
    generate_code
    echo ""
    generate_examples
    echo ""
    log_success "Synchronization completed!"
    ;;
  clean)
    log_info "Cleaning cache..."
    rm -rf "${CACHE_DIR}"
    log_success "Cache cleaned"
    ;;
  help | *)
    echo "N8N Nodes Synchronization Script"
    echo ""
    echo "Usage: $0 [command]"
    echo ""
    echo "Commands:"
    echo "  fetch      - Fetch n8n repository"
    echo "  parse      - Parse nodes and generate registry"
    echo "  diff       - Generate changelog from differences"
    echo "  generate   - Generate Go code and examples"
    echo "  stats      - Display node statistics"
    echo "  all        - Run all steps (default)"
    echo "  clean      - Clean cache directory"
    echo "  help       - Show this help"
    ;;
esac
