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

# Get coverage per package using go test directly
# Filter only lines starting with "ok" to avoid workflow aggregate line
COVERAGE_BY_PKG=$(go test -cover ./src/internal/provider/... 2>&1 | grep "^ok" | grep "coverage:" | grep -v "\[no statements\]")

# Also get total coverage (allow tests to fail)
go test -coverprofile=coverage.out -covermode=atomic ./src/internal/provider/... >/dev/null 2>&1 || true
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
# Use the actual coverage output from go test, which gives accurate per-package percentages
echo "$COVERAGE_BY_PKG" | while IFS= read -r line; do
  # Extract package name and coverage percentage
  # Format: "ok  	github.com/kodflow/n8n/src/internal/provider/variable	0.123s	coverage: 98.4% of statements"
  pkg=$(echo "$line" | awk '{print $2}')
  coverage=$(echo "$line" | grep -oP 'coverage: \K[0-9.]+%')

  # Skip if we couldn't extract coverage
  if [ -z "$coverage" ]; then
    continue
  fi

  PKG_VALUE=$(echo "$coverage" | sed 's/%//')

  # Determine icon
  if [ $(awk "BEGIN {print ($PKG_VALUE >= 90.0)}") -eq 1 ]; then
    ICON="🟢"
  elif [ $(awk "BEGIN {print ($PKG_VALUE >= 70.0)}") -eq 1 ]; then
    ICON="🟡"
  else
    ICON="🔴"
  fi

  # Add to table
  echo "| $ICON | \`$pkg\` | $coverage |" >>COVERAGE.MD
done

# Add detailed coverage by file with public functions only
cat >>COVERAGE.MD <<EOF

---

## Coverage Détaillée par Fichier

Cette section liste uniquement les **fonctions publiques** (exportées) pour identifier rapidement les fonctions non testées.
Les tableaux sont organisés par type de fichier pour faciliter la comparaison entre packages.

EOF

# Parse coverage data to extract public functions grouped by package/file
# Format: file.go:line:	FunctionName	coverage%
# We only want public functions (starting with uppercase after package.)

# Create a temporary directory for processing
TMP_DIR=$(mktemp -d)
trap "rm -rf $TMP_DIR" EXIT

# Get unique packages
PACKAGES=$(echo "$COVERAGE_DATA" | grep -E "^github.com/kodflow/n8n/src/internal/provider/" | grep -v "total:" | awk -F: '{print $1}' | sed 's|/[^/]*\.go$||' | sort -u)

# First pass: collect all data organized by filename
for pkg in $PACKAGES; do
  PKG_SHORT=$(echo "$pkg" | sed 's|github.com/kodflow/n8n/src/internal/provider||' | sed 's|^/||')
  if [ -z "$PKG_SHORT" ]; then
    PKG_SHORT="provider"
  fi

  # Get files in this package
  FILES=$(echo "$COVERAGE_DATA" | grep "^$pkg/" | awk -F: '{print $1}' | sort -u)

  for file in $FILES; do
    FILE_SHORT=$(basename "$file")

    # Extract public functions from this file
    PUBLIC_FUNCS=$(echo "$COVERAGE_DATA" | grep "^$file:" | awk '{
      gsub(/^[ \t]+|[ \t]+$/, "", $2);
      if ($2 ~ /^[A-Z]/ && $2 !~ /^New/) {
        print $2 "\t" $3
      }
    }')

    # Skip if no public functions
    if [ -z "$PUBLIC_FUNCS" ]; then
      continue
    fi

    # Save data to temp file organized by filename
    # Replace / with _ in package name to avoid directory issues
    PKG_SAFE=$(echo "$PKG_SHORT" | tr '/' '_')
    FILE_DATA="$TMP_DIR/$FILE_SHORT"
    mkdir -p "$FILE_DATA"
    echo "$PUBLIC_FUNCS" >"$FILE_DATA/$PKG_SAFE"
  done
done

# Second pass: generate tables organized by filename
for FILE_SHORT in $(ls "$TMP_DIR" 2>/dev/null | sort); do
  FILE_DATA="$TMP_DIR/$FILE_SHORT"

  echo "### 📄 $FILE_SHORT" >>COVERAGE.MD
  echo "" >>COVERAGE.MD

  # Collect all unique function names across all packages for this file
  ALL_FUNCS=$(cat "$FILE_DATA"/* 2>/dev/null | awk -F'\t' '{print $1}' | sort -u)

  # Skip if no functions
  if [ -z "$ALL_FUNCS" ]; then
    continue
  fi

  # Get list of packages that have this file (restore original names)
  PKG_LIST=$(ls "$FILE_DATA" | tr '_' '/' | sort)

  # Build table header
  HEADER="| Function |"
  SEPARATOR="|----------|"
  for pkg in $PKG_LIST; do
    HEADER="$HEADER $pkg |"
    SEPARATOR="$SEPARATOR:--------:|"
  done

  echo "$HEADER" >>COVERAGE.MD
  echo "$SEPARATOR" >>COVERAGE.MD

  # Build table rows for each function
  echo "$ALL_FUNCS" | while read -r func; do
    ROW="| \`$func\` |"

    for pkg in $PKG_LIST; do
      # Convert package name back to safe version for file lookup
      PKG_SAFE=$(echo "$pkg" | tr '/' '_')
      # Get coverage for this function in this package
      COV=$(grep "^$func"$'\t' "$FILE_DATA/$PKG_SAFE" 2>/dev/null | awk -F'\t' '{print $2}')

      if [ -z "$COV" ]; then
        # Function doesn't exist in this package
        ROW="$ROW - |"
      else
        # Function exists, add icon based on coverage
        COV_VALUE=$(echo "$COV" | sed 's/%//')
        if [ $(awk "BEGIN {print ($COV_VALUE >= 90.0)}") -eq 1 ]; then
          ICON="🟢"
        elif [ $(awk "BEGIN {print ($COV_VALUE >= 70.0)}") -eq 1 ]; then
          ICON="🟡"
        elif [ $(awk "BEGIN {print ($COV_VALUE > 0.0)}") -eq 1 ]; then
          ICON="🟠"
        else
          ICON="🔴"
        fi
        ROW="$ROW $ICON $COV |"
      fi
    done

    echo "$ROW" >>COVERAGE.MD
  done

  echo "" >>COVERAGE.MD
done

# Add footer
cat >>COVERAGE.MD <<EOF

---

## Légende des Icônes

- 🟢 **≥90%** - Excellente couverture
- 🟡 **70-89%** - Bonne couverture
- 🟠 **1-69%** - Couverture partielle (à améliorer)
- 🔴 **0%** - Fonction non testée

**Note:** Seules les fonctions publiques (exportées) sont listées. Les fonctions privées et constructeurs (\`New*\`) sont exclus.

---

*Rapport généré le: ${REPORT_DATE}*
*Threshold: 70.0%*
EOF

echo -e "${GREEN}✓${RESET} COVERAGE.MD generated successfully"
echo -e "  ${CYAN}Total Coverage:${RESET} ${TOTAL_COVERAGE}"
echo ""

# Check if coverage meets threshold (warning only, no exit)
if [ $(awk "BEGIN {print ($TOTAL_VALUE < 70.0)}") -eq 1 ]; then
  echo -e "${YELLOW}⚠️  Info: Coverage is below 70% threshold${RESET}"
fi
