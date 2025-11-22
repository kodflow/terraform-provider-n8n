#!/usr/bin/env bash
# Copyright (c) 2024 Florent (Kodflow). All rights reserved.
# Licensed under the Sustainable Use License 1.0
# See LICENSE in the project root for license information.
# Test n8n API connection

set -e

# Colors
RED='\033[31m'
GREEN='\033[32m'
YELLOW='\033[33m'
CYAN='\033[36m'
RESET='\033[0m'
BOLD='\033[1m'

# Timeout for curl commands (seconds)
CURL_TIMEOUT=30

echo ""
echo -e "${BOLD}${CYAN}🔍 Testing n8n API Connection${RESET}"
echo ""

# Check environment variables
if [ -z "$N8N_API_URL" ]; then
  echo -e "${RED}✗${RESET} N8N_API_URL not set"
  exit 1
fi

if [ -z "$N8N_API_KEY" ]; then
  echo -e "${RED}✗${RESET} N8N_API_KEY not set"
  exit 1
fi

# Normalize URL (remove trailing slash)
N8N_API_URL="${N8N_API_URL%/}"

echo -e "${CYAN}→${RESET} Testing connection to: ${BOLD}$N8N_API_URL${RESET}"
echo ""

# Test 1: Basic connectivity
echo -e "${CYAN}Test 1:${RESET} Basic connectivity..."
HTTP_CODE=$(curl -s --max-time "$CURL_TIMEOUT" -o /dev/null -w "%{http_code}" "$N8N_API_URL" || echo "000")
if [ "$HTTP_CODE" = "000" ]; then
  echo -e "${RED}✗${RESET} Cannot reach $N8N_API_URL (network error or timeout)"
  exit 1
elif [ "$HTTP_CODE" = "200" ] || [ "$HTTP_CODE" = "302" ]; then
  echo -e "${GREEN}✓${RESET} Instance is reachable (HTTP $HTTP_CODE)"
else
  echo -e "${YELLOW}⚠${RESET}  Instance returned HTTP $HTTP_CODE"
fi
echo ""

# Test 2: API endpoint accessibility
echo -e "${CYAN}Test 2:${RESET} API endpoint accessibility..."
API_URL="$N8N_API_URL/api/v1/workflows"
HTTP_CODE=$(curl -s --max-time "$CURL_TIMEOUT" -o /dev/null -w "%{http_code}" "$API_URL" || echo "000")
if [ "$HTTP_CODE" = "401" ]; then
  echo -e "${GREEN}✓${RESET} API endpoint exists (requires authentication)"
elif [ "$HTTP_CODE" = "000" ]; then
  echo -e "${RED}✗${RESET} Cannot reach API endpoint (network error or timeout)"
  exit 1
else
  echo -e "${YELLOW}⚠${RESET}  API endpoint returned HTTP $HTTP_CODE"
fi
echo ""

# Test 3: Authentication test
echo -e "${CYAN}Test 3:${RESET} Authentication test..."
echo -e "${CYAN}→${RESET} Using API key: [${#N8N_API_KEY} characters]"

# Use distinctive separator to avoid conflicts with API response content
RESPONSE=$(curl -s --max-time "$CURL_TIMEOUT" -w "\n###HTTP_CODE###:%{http_code}" \
  -H "X-N8N-API-KEY: $N8N_API_KEY" "$API_URL" || echo "###HTTP_CODE###:000")
HTTP_CODE=$(echo "$RESPONSE" | grep "###HTTP_CODE###" | cut -d: -f2)
BODY=$(echo "$RESPONSE" | sed '/###HTTP_CODE###/d')

if [ "$HTTP_CODE" = "000" ]; then
  echo -e "${RED}✗${RESET} Request failed (network error or timeout)"
  exit 1
elif [ "$HTTP_CODE" = "200" ]; then
  echo -e "${GREEN}✓${RESET} Authentication successful!"
  echo -e "${GREEN}→${RESET} Response preview:"
  echo "$BODY" | head -5
elif [ "$HTTP_CODE" = "401" ]; then
  echo -e "${RED}✗${RESET} Authentication failed (HTTP 401 Unauthorized)"
  echo -e "${RED}→${RESET} Response:"
  echo "$BODY"
  echo ""
  echo -e "${YELLOW}Possible causes:${RESET}"
  echo "  1. API key is invalid or expired"
  echo "  2. API key format is incorrect"
  echo "  3. n8n instance requires different authentication"
  exit 1
else
  echo -e "${RED}✗${RESET} Unexpected response (HTTP $HTTP_CODE)"
  echo "$BODY"
  exit 1
fi

echo ""
echo -e "${BOLD}${GREEN}✅ All tests passed!${RESET}"
echo ""
