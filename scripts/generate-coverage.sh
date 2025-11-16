#!/bin/bash
set -e

# Colors for output
CYAN='\033[36m'
GREEN='\033[32m'
YELLOW='\033[33m'
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
TOTAL_VALUE=${TOTAL_COVERAGE//%/}

echo -e "${CYAN}→${RESET} Parsing package coverage..."

# Get date
REPORT_DATE=$(date +%Y-%m-%d)

# Start building the markdown file
cat >COVERAGE.MD <<EOF
# Coverage Report

Automatically generated coverage report.

**Legend:**
- 🟢 ≥90% - Excellent coverage
- 🟡 70-89% - Good coverage
- 🔴 <70% - Insufficient coverage

---

## Global Coverage

| Metric | Value |
|--------|-------|
| **Total Coverage** | **${TOTAL_COVERAGE}** |
| **Threshold** | **70.0%** |
| **Status** | $(if [ "$(awk "BEGIN {print ($TOTAL_VALUE >= 70.0)}")" -eq 1 ]; then echo "✅ PASSED"; else echo "❌ FAILED"; fi) |

---

## Acceptance Tests (E2E)

Acceptance tests validate the real behavior of the provider against an n8n instance.

| Resource | Status | Tests |
|----------|:------:|-------|
EOF

# Find all acceptance tests
ACCEPTANCE_TESTS=$(find src/internal/provider -name "*_acceptance_test.go" -type f | sort)

for test_file in $ACCEPTANCE_TESTS; do
  # Extract package name (e.g., credential, workflow, tag)
  temp=${test_file#src/internal/provider/}
  pkg_name=${temp%%/*}

  # Count TestAcc* functions in the file
  test_count=$(grep -c "^func TestAcc" "$test_file" 2>/dev/null || echo "0")

  # Get test function names
  test_names=$(grep "^func TestAcc" "$test_file" 2>/dev/null | awk '{gsub(/func /, ""); gsub(/\(.*$/, ""); print}' || echo "")

  if [ "$test_count" -gt 0 ]; then
    # Format test names as a comma-separated list with backticks
    test_list=$(echo "$test_names" | sed 's/^/`/' | sed 's/$/`/' | tr '\n' ',' | sed 's/,$//' | sed 's/,/, /g')
    echo "| \`$pkg_name\` | ✅ | $test_list |" >>COVERAGE.MD
  fi
done

cat >>COVERAGE.MD <<'EOF'

**Legend:**
- ✅ Acceptance tests present
- Acceptance tests verify real operations via the n8n API

---

## Coverage by Package

| Icon | Package | Coverage |
|:----:|---------|----------|
EOF

# Parse coverage by package
# Use the actual coverage output from go test, which gives accurate per-package percentages
echo "$COVERAGE_BY_PKG" | while IFS= read -r line; do
  # Extract package name and coverage percentage
  # Format: "ok  	github.com/kodflow/terraform-provider-n8n/src/internal/provider/variable	0.123s	coverage: 98.4% of statements"
  pkg=$(echo "$line" | awk '{print $2}')
  coverage=$(echo "$line" | grep -oP 'coverage: \K[0-9.]+%')

  # Skip if we couldn't extract coverage
  if [ -z "$coverage" ]; then
    continue
  fi

  PKG_VALUE=${coverage//%/}

  # Determine icon
  if [ "$(awk "BEGIN {print ($PKG_VALUE >= 90.0)}")" -eq 1 ]; then
    ICON="🟢"
  elif [ "$(awk "BEGIN {print ($PKG_VALUE >= 70.0)}")" -eq 1 ]; then
    ICON="🟡"
  else
    ICON="🔴"
  fi

  # Add to table
  echo "| $ICON | \`$pkg\` | $coverage |" >>COVERAGE.MD
done

# Add detailed coverage by category with public functions only
cat >>COVERAGE.MD <<EOF

---

## Detailed Coverage by Category

This section lists only **public functions** (exported) to quickly identify untested functions.
Tables are organized by resource category to facilitate understanding of the provider architecture.

EOF

# Parse coverage data to extract public functions grouped by package/file
# Format: file.go:line:	FunctionName	coverage%
# We only want public functions (starting with uppercase after package.)

# Create a temporary directory for processing
TMP_DIR=$(mktemp -d)
trap 'rm -rf "$TMP_DIR"' EXIT

# Get unique packages
PACKAGES=$(echo "$COVERAGE_DATA" | grep -E "^github.com/kodflow/terraform-provider-n8n/src/internal/provider/" | grep -v "total:" | awk -F: '{print $1}' | sed 's|/[^/]*\.go$||' | sort -u)

# First pass: collect all data organized by filename
for pkg in $PACKAGES; do
  temp=${pkg#github.com/kodflow/terraform-provider-n8n/src/internal/provider}
  PKG_SHORT=${temp#/}
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

# Define all provider packages for complete overview
ALL_PROVIDER_PACKAGES="credential execution project sourcecontrol tag user variable workflow"

# Helper function to generate coverage table for a file across packages
generate_coverage_table() {
  local FILE_SHORT="$1"
  shift
  local PACKAGES="$@"

  local FILE_DATA="$TMP_DIR/$FILE_SHORT"

  if [ ! -d "$FILE_DATA" ]; then
    return
  fi

  # Collect all unique function names for this file
  local ALL_FUNCS
  ALL_FUNCS=$(cat "$FILE_DATA"/* 2>/dev/null | awk -F'\t' '{print $1}' | sort -u)

  if [ -z "$ALL_FUNCS" ]; then
    return
  fi

  # Build table header
  local HEADER="| Function |"
  local SEPARATOR="|----------|"
  for pkg in $PACKAGES; do
    HEADER="$HEADER $pkg |"
    SEPARATOR="$SEPARATOR:--------:|"
  done

  echo "$HEADER" >>COVERAGE.MD
  echo "$SEPARATOR" >>COVERAGE.MD

  # Build table rows for each function
  echo "$ALL_FUNCS" | while read -r func; do
    local ROW="| \`$func\` |"

    for pkg in $PACKAGES; do
      local PKG_SAFE
      PKG_SAFE=$(echo "$pkg" | tr '/' '_')
      local COV
      COV=$(grep "^$func"$'\t' "$FILE_DATA/$PKG_SAFE" 2>/dev/null | awk -F'\t' '{print $2}')

      if [ -z "$COV" ]; then
        local PKG_PATH="src/internal/provider/$pkg/$FILE_SHORT"
        if [ -f "$PKG_PATH" ] && grep -q "func.*[[:space:]]$func(" "$PKG_PATH"; then
          local FUNC_BODY
          FUNC_BODY=$(awk "/func.*[[:space:]]$func\(/,/^}/" "$PKG_PATH" | grep -v "^//" | grep -v "^[[:space:]]*//")
          local EXECUTABLE_LINES
          EXECUTABLE_LINES=$(echo "$FUNC_BODY" | tail -n +2 | head -n -1 | grep -v "^[[:space:]]*$" | wc -l)

          if [ "$EXECUTABLE_LINES" -eq 0 ]; then
            ROW="$ROW 🔵 N/A |"
          else
            ROW="$ROW 🔴 0.0% |"
          fi
        else
          ROW="$ROW 🔵 N/A |"
        fi
      else
        local COV_VALUE
        COV_VALUE=${COV//%/}
        if [ "$(awk "BEGIN {print ($COV_VALUE >= 90.0)}")" -eq 1 ]; then
          local ICON="🟢"
        elif [ "$(awk "BEGIN {print ($COV_VALUE >= 70.0)}")" -eq 1 ]; then
          local ICON="🟡"
        elif [ "$(awk "BEGIN {print ($COV_VALUE > 0.0)}")" -eq 1 ]; then
          local ICON="🟠"
        else
          local PKG_PATH="src/internal/provider/$pkg/$FILE_SHORT"
          if [ -f "$PKG_PATH" ]; then
            local FUNC_BODY=$(awk "/func.*[[:space:]]$func\(/,/^}/" "$PKG_PATH" | grep -v "^//" | grep -v "^[[:space:]]*//")
            local EXECUTABLE_LINES=$(echo "$FUNC_BODY" | tail -n +2 | head -n -1 | grep -v "^[[:space:]]*$" | wc -l)

            if [ "$EXECUTABLE_LINES" -eq 0 ]; then
              ICON="🔵"
              COV="N/A"
            else
              ICON="🔴"
            fi
          else
            ICON="🔴"
          fi
        fi
        ROW="$ROW $ICON $COV |"
      fi
    done

    echo "$ROW" >>COVERAGE.MD
  done

  echo "" >>COVERAGE.MD
}

# Second pass: generate tables organized by resource category

# === PRIMARY RESOURCES (CRUD Entities) ===
echo "## 📦 Primary Resources (CRUD Entities)" >>COVERAGE.MD
echo "" >>COVERAGE.MD
echo "Complete lifecycle management of n8n resources (Create, Read, Update, Delete)." >>COVERAGE.MD
echo "" >>COVERAGE.MD

# Find all packages with resource.go
PRIMARY_PACKAGES=""
for pkg in $ALL_PROVIDER_PACKAGES; do
  if [ -f "src/internal/provider/$pkg/resource.go" ]; then
    PRIMARY_PACKAGES="$PRIMARY_PACKAGES $pkg"
  fi
done

if [ -n "$PRIMARY_PACKAGES" ]; then
  generate_coverage_table "resource.go" "$PRIMARY_PACKAGES"
fi

# === SECONDARY RESOURCES (Operations/Relations) ===
echo "---" >>COVERAGE.MD
echo "" >>COVERAGE.MD
echo "## 🔧 Secondary Resources (Operations/Relations)" >>COVERAGE.MD
echo "" >>COVERAGE.MD
echo "Special operations and resource relationship management." >>COVERAGE.MD
echo "" >>COVERAGE.MD

# Find all *_resource.go files (excluding resource.go)
SECONDARY_FILES=$(find "$TMP_DIR" -maxdepth 1 -name "*_resource.go" ! -name "resource.go" -printf "%f\n" 2>/dev/null | sort)

for FILE_SHORT in $SECONDARY_FILES; do
  # Extract resource type (e.g., retry_resource.go -> retry)
  RES_TYPE=${FILE_SHORT%_resource.go}

  # Find packages that have this file
  SEC_PACKAGES=""
  for pkg in $ALL_PROVIDER_PACKAGES; do
    if [ -f "src/internal/provider/$pkg/$FILE_SHORT" ]; then
      SEC_PACKAGES="$SEC_PACKAGES $pkg"
    fi
  done

  if [ -n "$SEC_PACKAGES" ]; then
    # Generate a nice title based on resource type and package
    case "$RES_TYPE" in
      pull)
        TITLE="Source Control Pull"
        ;;
      retry)
        TITLE="Execution Retry"
        ;;
      transfer)
        TITLE="Transfer"
        ;;
      user)
        TITLE="Project User"
        ;;
      *)
        # Default: capitalize first letter
        TITLE=$(echo "$RES_TYPE" | sed 's/\b\(.\)/\u\1/')
        ;;
    esac

    echo "### $TITLE" >>COVERAGE.MD
    echo "" >>COVERAGE.MD
    generate_coverage_table "$FILE_SHORT" "$SEC_PACKAGES"
  fi
done

# === DATA SOURCES ===
echo "---" >>COVERAGE.MD
echo "" >>COVERAGE.MD
echo "## 📊 Data Sources" >>COVERAGE.MD
echo "" >>COVERAGE.MD
echo "Reading n8n resources without managing their lifecycle." >>COVERAGE.MD
echo "" >>COVERAGE.MD

# datasource.go (singular)
DS_PACKAGES=""
for pkg in $ALL_PROVIDER_PACKAGES; do
  if [ -f "src/internal/provider/$pkg/datasource.go" ]; then
    DS_PACKAGES="$DS_PACKAGES $pkg"
  fi
done

if [ -n "$DS_PACKAGES" ]; then
  echo "### datasource (singular)" >>COVERAGE.MD
  echo "" >>COVERAGE.MD
  generate_coverage_table "datasource.go" "$DS_PACKAGES"
fi

# datasources.go (plural)
DSS_PACKAGES=""
for pkg in $ALL_PROVIDER_PACKAGES; do
  if [ -f "src/internal/provider/$pkg/datasources.go" ]; then
    DSS_PACKAGES="$DSS_PACKAGES $pkg"
  fi
done

if [ -n "$DSS_PACKAGES" ]; then
  echo "### datasources (plural)" >>COVERAGE.MD
  echo "" >>COVERAGE.MD
  generate_coverage_table "datasources.go" "$DSS_PACKAGES"
fi

# Add footer
cat >>COVERAGE.MD <<EOF

---

## Icon Legend

- 🟢 **≥90%** - Excellent coverage
- 🟡 **70-89%** - Good coverage
- 🟠 **1-69%** - Partial coverage (needs improvement)
- 🔴 **0%** - Untested function (implemented but no tests)
- 🔵 **N/A** - Not applicable (not in API or intentionally empty)

**Note:** Only public functions (exported) are listed. Private functions and constructors (\`New*\`) are excluded.

---

*Report generated on: ${REPORT_DATE}*
*Threshold: 70.0%*
EOF

echo -e "${GREEN}✓${RESET} COVERAGE.MD generated successfully"
echo -e "  ${CYAN}Total Coverage:${RESET} ${TOTAL_COVERAGE}"
echo ""

# Check if coverage meets threshold (warning only, no exit)
if [ "$(awk "BEGIN {print ($TOTAL_VALUE < 70.0)}")" -eq 1 ]; then
  echo -e "${YELLOW}⚠️  Info: Coverage is below 70% threshold${RESET}"
fi
