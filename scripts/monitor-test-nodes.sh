#!/usr/bin/env bash

# Monitor test-nodes.sh progress

LOG_FILE="${1:-/tmp/test-nodes-full.log}"

if [ ! -f "$LOG_FILE" ]; then
  echo "❌ Log file not found: $LOG_FILE"
  exit 1
fi

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'
BOLD='\033[1m'

echo -e "${BOLD}${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BOLD}  Node Testing Monitor${NC}"
echo -e "${BOLD}${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

# Get current progress
CURRENT=$(grep -oP '\[\d+/\d+\]' "$LOG_FILE" | tail -1 | grep -oP '\d+' | head -1)
TOTAL=$(grep -oP '\[\d+/\d+\]' "$LOG_FILE" | tail -1 | grep -oP '\d+' | tail -1)
SUCCESS=$(grep -c "✓ Success" "$LOG_FILE" || echo 0)
FAILED=$(grep -c "✗.*failed" "$LOG_FILE" || echo 0)

if [ -n "$CURRENT" ] && [ -n "$TOTAL" ]; then
  PERCENT=$(awk "BEGIN {printf \"%.1f\", ($CURRENT/$TOTAL)*100}")
  echo -e "${BLUE}Progress:${NC} ${CURRENT}/${TOTAL} nodes (${PERCENT}%)"
else
  echo -e "${YELLOW}Initializing...${NC}"
fi

echo -e "${GREEN}Success:${NC} ${SUCCESS}"
echo -e "${RED}Failed:${NC} ${FAILED}"
echo ""

# Last 10 tested nodes
echo -e "${BOLD}Last 10 nodes:${NC}"
grep -E "Testing:" "$LOG_FILE" | tail -10 | while read -r line; do
  NODE=$(echo "$line" | grep -oP 'Testing: \K.*')
  if [ -n "$NODE" ]; then
    echo "  • $NODE"
  fi
done
echo ""

# Any failures
if [ "$FAILED" -gt 0 ]; then
  echo -e "${RED}${BOLD}Failed nodes:${NC}"
  grep "✗" "$LOG_FILE" | tail -10
  echo ""
fi

# Check if still running
if ps aux | grep -q "[t]est-nodes.sh"; then
  echo -e "${BLUE}Status:${NC} ${GREEN}Running${NC}"

  # Estimate time remaining
  if [ -n "$CURRENT" ] && [ -n "$TOTAL" ] && [ "$CURRENT" -gt 0 ]; then
    # Get start time from log
    START_TIME=$(stat -c %Y "$LOG_FILE" 2>/dev/null || stat -f %m "$LOG_FILE" 2>/dev/null)
    CURRENT_TIME=$(date +%s)
    ELAPSED=$((CURRENT_TIME - START_TIME))

    if [ "$ELAPSED" -gt 0 ]; then
      AVG_TIME_PER_NODE=$(awk "BEGIN {printf \"%.1f\", $ELAPSED/$CURRENT}")
      REMAINING_NODES=$((TOTAL - CURRENT))
      ESTIMATED_SECONDS=$(awk "BEGIN {printf \"%.0f\", $REMAINING_NODES*$AVG_TIME_PER_NODE}")

      HOURS=$((ESTIMATED_SECONDS / 3600))
      MINUTES=$(((ESTIMATED_SECONDS % 3600) / 60))

      echo -e "${BLUE}Elapsed:${NC} $((ELAPSED / 60))m $((ELAPSED % 60))s"
      echo -e "${BLUE}Estimated remaining:${NC} ${HOURS}h ${MINUTES}m"
    fi
  fi
else
  echo -e "${BLUE}Status:${NC} ${YELLOW}Completed or not running${NC}"

  # Show final summary if available
  if grep -q "Test Summary" "$LOG_FILE"; then
    echo ""
    echo -e "${BOLD}Final Summary:${NC}"
    grep -A 10 "Test Summary" "$LOG_FILE" | tail -10
  fi
fi

echo ""
echo -e "${BOLD}${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
