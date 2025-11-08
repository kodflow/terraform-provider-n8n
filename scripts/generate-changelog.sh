#!/bin/bash
# Generate CHANGELOG.md from git commits
# Follows Conventional Commits and Keep a Changelog format

set -e

CHANGELOG_FILE="CHANGELOG.md"
TEMP_FILE=$(mktemp)
BRANCH=${1:-feat/codegen-pipeline}
BASE_BRANCH=${2:-main}

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}🔄 Generating CHANGELOG from git history...${NC}"

# Function to categorize commit by type
categorize_commit() {
  local commit_msg="$1"
  local commit_hash="$2"

  # Extract type from conventional commit format
  if [[ $commit_msg =~ ^(feat|feature)(\(.+\))?:\ (.+) ]]; then
    echo "### 🚀 Features|${BASH_REMATCH[3]} (\`${commit_hash:0:7}\`)"
  elif [[ $commit_msg =~ ^fix(\(.+\))?:\ (.+) ]]; then
    echo "### 🐛 Bug Fixes|${BASH_REMATCH[2]} (\`${commit_hash:0:7}\`)"
  elif [[ $commit_msg =~ ^test(\(.+\))?:\ (.+) ]]; then
    echo "### ✅ Tests|${BASH_REMATCH[2]} (\`${commit_hash:0:7}\`)"
  elif [[ $commit_msg =~ ^docs(\(.+\))?:\ (.+) ]]; then
    echo "### 📚 Documentation|${BASH_REMATCH[2]} (\`${commit_hash:0:7}\`)"
  elif [[ $commit_msg =~ ^refactor(\(.+\))?:\ (.+) ]]; then
    echo "### ♻️ Refactoring|${BASH_REMATCH[2]} (\`${commit_hash:0:7}\`)"
  elif [[ $commit_msg =~ ^perf(\(.+\))?:\ (.+) ]]; then
    echo "### ⚡ Performance|${BASH_REMATCH[2]} (\`${commit_hash:0:7}\`)"
  elif [[ $commit_msg =~ ^build(\(.+\))?:\ (.+) ]]; then
    echo "### 🔧 Build|${BASH_REMATCH[2]} (\`${commit_hash:0:7}\`)"
  elif [[ $commit_msg =~ ^ci(\(.+\))?:\ (.+) ]]; then
    echo "### 🤖 CI/CD|${BASH_REMATCH[2]} (\`${commit_hash:0:7}\`)"
  elif [[ $commit_msg =~ ^chore(\(.+\))?:\ (.+) ]]; then
    echo "### 🔨 Chore|${BASH_REMATCH[2]} (\`${commit_hash:0:7}\`)"
  elif [[ $commit_msg =~ ^style(\(.+\))?:\ (.+) ]]; then
    echo "### 💄 Style|${BASH_REMATCH[2]} (\`${commit_hash:0:7}\`)"
  else
    echo "### 📝 Other|${commit_msg} (\`${commit_hash:0:7}\`)"
  fi
}

# Start building changelog
cat >"$TEMP_FILE" <<'EOF'
# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/).

---

## [Unreleased]

EOF

# Get commits from current branch that are not in base branch
echo -e "${YELLOW}📊 Analyzing commits on ${BRANCH} (not in ${BASE_BRANCH})...${NC}"

# Collect commits by category
declare -A categories
declare -A category_order=(
  ["### 🚀 Features"]=1
  ["### 🐛 Bug Fixes"]=2
  ["### ✅ Tests"]=3
  ["### ♻️ Refactoring"]=4
  ["### ⚡ Performance"]=5
  ["### 📚 Documentation"]=6
  ["### 🔧 Build"]=7
  ["### 🤖 CI/CD"]=8
  ["### 💄 Style"]=9
  ["### 🔨 Chore"]=10
  ["### 📝 Other"]=11
)

while IFS= read -r line; do
  commit_hash=$(echo "$line" | awk '{print $1}')
  commit_msg=$(echo "$line" | cut -d' ' -f2-)

  result=$(categorize_commit "$commit_msg" "$commit_hash")
  category=$(echo "$result" | cut -d'|' -f1)
  item=$(echo "$result" | cut -d'|' -f2)

  if [[ -z "${categories[$category]}" ]]; then
    categories["$category"]="$item"
  else
    categories["$category"]="${categories[$category]}"$'\n'"$item"
  fi
done < <(git log --oneline --no-merges "$BRANCH" ^"$BASE_BRANCH" 2>/dev/null || echo "")

# Write categories in order
for category in "${!category_order[@]}"; do
  if [[ -n "${categories[$category]}" ]]; then
    echo "" >>"$TEMP_FILE"
    echo "$category" >>"$TEMP_FILE"
    echo "" >>"$TEMP_FILE"
    echo "${categories[$category]}" | while IFS= read -r item; do
      echo "- $item" >>"$TEMP_FILE"
    done
  fi
done

# Add statistics section
echo "" >>"$TEMP_FILE"
echo "---" >>"$TEMP_FILE"
echo "" >>"$TEMP_FILE"
echo "### 📊 Statistics" >>"$TEMP_FILE"
echo "" >>"$TEMP_FILE"

# Count commits by type
total_commits=$(git log --oneline --no-merges "$BRANCH" ^"$BASE_BRANCH" 2>/dev/null | wc -l || echo "0")
feat_count=$(git log --oneline --no-merges --grep="^feat" "$BRANCH" ^"$BASE_BRANCH" 2>/dev/null | wc -l || echo "0")
fix_count=$(git log --oneline --no-merges --grep="^fix" "$BRANCH" ^"$BASE_BRANCH" 2>/dev/null | wc -l || echo "0")
test_count=$(git log --oneline --no-merges --grep="^test" "$BRANCH" ^"$BASE_BRANCH" 2>/dev/null | wc -l || echo "0")
refactor_count=$(git log --oneline --no-merges --grep="^refactor" "$BRANCH" ^"$BASE_BRANCH" 2>/dev/null | wc -l || echo "0")

echo "- **Total commits:** $total_commits" >>"$TEMP_FILE"
echo "- **Features:** $feat_count" >>"$TEMP_FILE"
echo "- **Bug fixes:** $fix_count" >>"$TEMP_FILE"
echo "- **Tests:** $test_count" >>"$TEMP_FILE"
echo "- **Refactoring:** $refactor_count" >>"$TEMP_FILE"

# Add coverage info if exists
if [ -f "COVERAGE.MD" ]; then
  coverage=$(grep "TOTAL (Global)" COVERAGE.MD | grep -oP '\d+\.\d+%' || echo "N/A")
  echo "- **Test coverage:** $coverage" >>"$TEMP_FILE"
fi

# Add contributors
echo "" >>"$TEMP_FILE"
echo "### 👥 Contributors" >>"$TEMP_FILE"
echo "" >>"$TEMP_FILE"
git log --format='%aN <%aE>' "$BRANCH" ^"$BASE_BRANCH" 2>/dev/null | sort -u | while read -r contributor; do
  echo "- $contributor" >>"$TEMP_FILE"
done || true

# Add generation timestamp
echo "" >>"$TEMP_FILE"
echo "---" >>"$TEMP_FILE"
echo "" >>"$TEMP_FILE"
echo "*Changelog generated automatically on $(date '+%Y-%m-%d %H:%M:%S')*" >>"$TEMP_FILE"

# Move temp file to final location
mv "$TEMP_FILE" "$CHANGELOG_FILE"

echo -e "${GREEN}✅ CHANGELOG.md generated successfully!${NC}"
echo -e "${BLUE}📝 Total commits processed: $total_commits${NC}"
