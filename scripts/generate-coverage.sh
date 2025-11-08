#!/bin/bash
set -e

# Colors for output
CYAN='\033[36m'
GREEN='\033[32m'
YELLOW='\033[33m'
RED='\033[31m'
RESET='\033[0m'
BOLD='\033[1m'

echo ""
echo -e "${BOLD}${CYAN}📊 Generating Coverage Report...${RESET}"
echo ""

# Run tests with coverage
echo -e "${CYAN}→${RESET} Running tests with coverage..."
go test -coverprofile=coverage.out -covermode=atomic ./src/internal/provider/... >/dev/null 2>&1

# Generate coverage report
echo -e "${CYAN}→${RESET} Generating coverage data..."
COVERAGE_DATA=$(go tool cover -func=coverage.out)

# Extract total coverage
TOTAL_COVERAGE=$(echo "$COVERAGE_DATA" | tail -1 | awk '{print $3}')
TOTAL_VALUE=$(echo "$TOTAL_COVERAGE" | sed 's/%//')

echo -e "${CYAN}→${RESET} Parsing package coverage..."

# Get date
REPORT_DATE=$(date +%Y-%m-%d)

# Start building the markdown file
cat >COVERAGE.MD <<EOF
# Coverage Report

Rapport de couverture généré automatiquement.

**Légende:**
- 🟢 ≥90% - Excellente couverture
- 🟡 70-89% - Bonne couverture
- 🔴 <70% - Couverture insuffisante

---

## Coverage Global

| Metric | Value |
|--------|-------|
| **Total Coverage** | **${TOTAL_COVERAGE}** |
| **Threshold** | **70.0%** |
| **Status** | $(if [ $(awk "BEGIN {print ($TOTAL_VALUE >= 70.0)}") -eq 1 ]; then echo "✅ PASSED"; else echo "❌ FAILED"; fi) |

---

## Coverage par Package

| Icon | Package | Coverage |
|:----:|---------|----------|
EOF

# Parse coverage by package
PACKAGES=$(echo "$COVERAGE_DATA" | grep -v "^total:" | awk '{print $1}' | sed 's/\/[^\/]*$//' | sort -u)

for pkg in $PACKAGES; do
  # Get coverage for this package
  PKG_COV=$(echo "$COVERAGE_DATA" | grep "^$pkg" | awk '{sum += $3; count++} END {if (count > 0) printf "%.1f%%", sum/count; else print "0.0%"}')
  PKG_VALUE=$(echo "$PKG_COV" | sed 's/%//')

  # Determine icon
  if [ $(awk "BEGIN {print ($PKG_VALUE >= 90.0)}") -eq 1 ]; then
    ICON="🟢"
  elif [ $(awk "BEGIN {print ($PKG_VALUE >= 70.0)}") -eq 1 ]; then
    ICON="🟡"
  else
    ICON="🔴"
  fi

  # Add to table
  echo "| $ICON | \`$pkg\` | $PKG_COV |" >>COVERAGE.MD
done

# Add footer
cat >>COVERAGE.MD <<EOF

---

## How to Improve Coverage

To improve coverage, focus on:

1. **Add unit tests** for untested functions
2. **Test error paths** and edge cases
3. **Add integration tests** for complex scenarios
4. **Test validation logic** and input handling

Run \`make test\` to see detailed coverage per file.

---

*Rapport généré le: ${REPORT_DATE}*
*Threshold: 70.0%*
EOF

echo -e "${GREEN}✓${RESET} COVERAGE.MD generated successfully"
echo -e "  ${CYAN}Total Coverage:${RESET} ${TOTAL_COVERAGE}"
echo ""

# Check if coverage meets threshold
if [ $(awk "BEGIN {print ($TOTAL_VALUE < 70.0)}") -eq 1 ]; then
  echo -e "${RED}⚠️  Warning: Coverage is below 70% threshold${RESET}"
  exit 1
fi
