#!/usr/bin/env bash
# Copyright (c) 2024 Florent (Kodflow). All rights reserved.
# Licensed under the Sustainable Use License 1.0
# See LICENSE in the project root for license information.

# Local n8n Setup Script
# Creates the owner account and an API key on the local n8n instance,
# then writes N8N_BASE_URL and N8N_API_KEY to /workspace/.env.
#
# Usage:
#   make n8n/setup          (recommended)
#   ./scripts/setup-local-n8n.sh
#
# Environment variables (optional overrides):
#   N8N_URL       - n8n base URL (default: http://n8n:5678)
#   N8N_EMAIL     - owner email  (default: admin@local.dev)
#   N8N_PASSWORD  - owner pass   (default: Admin1234!)

set -euo pipefail

# ── Configuration ────────────────────────────────────────────────────────────

N8N_URL="${N8N_URL:-http://n8n:5678}"
N8N_EMAIL="${N8N_EMAIL:-admin@local.dev}"
N8N_PASSWORD="${N8N_PASSWORD:-Admin1234!}"
N8N_FIRST_NAME="Admin"
N8N_LAST_NAME="Local"
COOKIE_JAR="/tmp/n8n-setup-cookies.txt"
ENV_FILE="/workspace/.env"
MAX_WAIT=60

# ── Colors ───────────────────────────────────────────────────────────────────

RED='\033[31m'
GREEN='\033[32m'
CYAN='\033[36m'
YELLOW='\033[33m'
BOLD='\033[1m'
RESET='\033[0m'

log() { echo -e "${CYAN}→${RESET} $*"; }
ok() { echo -e "${GREEN}✓${RESET} $*"; }
warn() { echo -e "${YELLOW}⚠${RESET}  $*"; }
fail() {
  echo -e "${RED}✗${RESET} $*" >&2
  exit 1
}

# ── Wait for n8n health ───────────────────────────────────────────────────────

wait_for_n8n() {
  log "Waiting for n8n to be ready at ${N8N_URL}..."
  local attempt=0
  until curl -sf "${N8N_URL}/healthz" >/dev/null 2>&1; do
    attempt=$((attempt + 1))
    if [ "${attempt}" -ge "${MAX_WAIT}" ]; then
      fail "n8n did not become healthy after ${MAX_WAIT} attempts. Is the service running?"
    fi
    sleep 3
  done
  ok "n8n is ready"
}

# ── Create owner account (idempotent) ─────────────────────────────────────────

setup_owner() {
  log "Creating owner account (${N8N_EMAIL})..."
  local response http_code
  response=$(
    curl -s -o /tmp/n8n-owner-response.json -w "%{http_code}" \
      -X POST "${N8N_URL}/rest/owner-setup" \
      -H "Content-Type: application/json" \
      -d "{
        \"email\": \"${N8N_EMAIL}\",
        \"password\": \"${N8N_PASSWORD}\",
        \"firstName\": \"${N8N_FIRST_NAME}\",
        \"lastName\": \"${N8N_LAST_NAME}\"
      }"
  )
  http_code="${response}"

  if [ "${http_code}" = "200" ]; then
    ok "Owner account created"
    return 0
  fi

  # n8n returns 400/409 when owner already exists — that is fine
  local msg
  msg=$(jq -r '.message // empty' /tmp/n8n-owner-response.json 2>/dev/null || echo "")
  if echo "${msg}" | grep -qi "already\|exists\|owner"; then
    warn "Owner account already exists — skipping creation"
    return 0
  fi

  fail "Owner setup failed (HTTP ${http_code}): ${msg}"
}

# ── Login and capture session cookie ─────────────────────────────────────────

login() {
  log "Logging in as ${N8N_EMAIL}..."
  rm -f "${COOKIE_JAR}"
  local http_code
  http_code=$(
    curl -s -o /dev/null -w "%{http_code}" \
      -c "${COOKIE_JAR}" \
      -X POST "${N8N_URL}/rest/login" \
      -H "Content-Type: application/json" \
      -d "{\"email\": \"${N8N_EMAIL}\", \"password\": \"${N8N_PASSWORD}\"}"
  )

  if [ "${http_code}" != "200" ]; then
    fail "Login failed (HTTP ${http_code}). Check your N8N_EMAIL and N8N_PASSWORD."
  fi
  ok "Logged in successfully"
}

# ── Create API key ────────────────────────────────────────────────────────────

create_api_key() {
  log "Creating API key..."
  local response http_code api_key
  http_code=$(
    curl -s -o /tmp/n8n-apikey-response.json -w "%{http_code}" \
      -b "${COOKIE_JAR}" \
      -X POST "${N8N_URL}/rest/api-key" \
      -H "Content-Type: application/json" \
      -d '{"label": "devcontainer"}'
  )

  if [ "${http_code}" != "200" ] && [ "${http_code}" != "201" ]; then
    fail "API key creation failed (HTTP ${http_code}). Check the n8n logs."
  fi

  api_key=$(jq -r '.data.apiKey // .apiKey // empty' /tmp/n8n-apikey-response.json 2>/dev/null || echo "")
  if [ -z "${api_key}" ]; then
    fail "Could not extract API key from response. Raw: $(cat /tmp/n8n-apikey-response.json)"
  fi

  ok "API key created"
  echo "${api_key}"
}

# ── Write credentials to .env ─────────────────────────────────────────────────

update_env() {
  local api_key="$1"

  log "Writing credentials to ${ENV_FILE}..."

  # Preserve all non-n8n lines from existing .env
  local existing=""
  if [ -f "${ENV_FILE}" ]; then
    existing=$(grep -v '^N8N_BASE_URL=' "${ENV_FILE}" | grep -v '^N8N_API_KEY=' || true)
  fi

  {
    if [ -n "${existing}" ]; then
      echo "${existing}"
    fi
    echo "N8N_BASE_URL=${N8N_URL}"
    echo "N8N_API_KEY=${api_key}"
  } >"${ENV_FILE}"

  ok "Credentials written to ${ENV_FILE}"
}

# ── Cleanup temp files ────────────────────────────────────────────────────────

cleanup() {
  rm -f "${COOKIE_JAR}" /tmp/n8n-owner-response.json /tmp/n8n-apikey-response.json
}

# ── Main ──────────────────────────────────────────────────────────────────────

main() {
  echo ""
  echo -e "${BOLD}${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${RESET}"
  echo -e "${BOLD}  Local n8n Setup${RESET}"
  echo -e "${BOLD}${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${RESET}"
  echo ""

  trap cleanup EXIT

  wait_for_n8n
  setup_owner
  login
  local api_key
  api_key=$(create_api_key)
  update_env "${api_key}"

  echo ""
  echo -e "${BOLD}${GREEN}✅ n8n is ready for acceptance tests${RESET}"
  echo ""
  echo -e "  URL:     ${CYAN}${N8N_URL}${RESET}"
  echo -e "  Browser: ${CYAN}http://localhost:5678${RESET}"
  echo -e "  Email:   ${CYAN}${N8N_EMAIL}${RESET}"
  echo ""
  echo "Run acceptance tests with: make test/acceptance"
  echo ""
}

main
