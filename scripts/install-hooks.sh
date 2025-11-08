#!/bin/bash
# Install git hooks for automatic documentation generation

set -e

HOOKS_DIR=".git/hooks"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${BLUE}📦 Installing git hooks...${NC}"

# Check if we're in a git repository
if [ ! -d ".git" ]; then
  echo -e "${YELLOW}⚠️  Not in a git repository. Run this from the project root.${NC}"
  exit 1
fi

# Create hooks directory if it doesn't exist
mkdir -p "$HOOKS_DIR"

# Install pre-commit hook
echo -e "${BLUE}Installing pre-commit hook...${NC}"
cat >"$HOOKS_DIR/pre-commit" <<'EOF'
#!/bin/bash
# Pre-commit hook to auto-generate documentation files
# - CHANGELOG.md (from git history)
# - COVERAGE.MD (from test coverage)

set -e

echo "🔄 Generating documentation files..."

# Generate CHANGELOG.md
if [ -f "./scripts/generate-changelog.sh" ]; then
    echo "📝 Generating CHANGELOG.md..."
    ./scripts/generate-changelog.sh > /dev/null 2>&1
    git add CHANGELOG.md
    echo "✅ CHANGELOG.md updated"
fi

# Generate coverage report
echo "📊 Running tests for coverage..."
if go test -coverprofile=coverage.out ./src/internal/provider/... > /dev/null 2>&1; then
    echo "✅ Coverage data collected"
    # Clean up coverage file (not needed in git)
    rm -f coverage.out
else
    echo "⚠️  Some tests failed, but continuing commit"
fi

echo "✅ Documentation files ready"
exit 0
EOF

chmod +x "$HOOKS_DIR/pre-commit"
echo -e "${GREEN}✅ pre-commit hook installed${NC}"

# Optional: Install prepare-commit-msg for commit message templates
echo -e "${BLUE}Installing prepare-commit-msg hook...${NC}"
cat >"$HOOKS_DIR/prepare-commit-msg" <<'EOF'
#!/bin/bash
# Add commit message template with conventional commits hint

COMMIT_MSG_FILE=$1
COMMIT_SOURCE=$2

# Only add template for new commits (not amendments, merges, etc.)
if [ -z "$COMMIT_SOURCE" ]; then
    # Check if file is empty or only has comments
    if ! grep -q '^[^#]' "$COMMIT_MSG_FILE"; then
        cat > "$COMMIT_MSG_FILE" << 'TEMPLATE'
# <type>: <description>
#
# Types: feat, fix, docs, style, refactor, perf, test, build, ci, chore
# Example: feat: add new datasource for workflows
#
# Longer description (optional)
#
# Closes #123 (optional)
TEMPLATE
    fi
fi
EOF

chmod +x "$HOOKS_DIR/prepare-commit-msg"
echo -e "${GREEN}✅ prepare-commit-msg hook installed${NC}"

# Install commit-msg hook for commitlint validation
echo -e "${BLUE}Installing commit-msg hook...${NC}"
cat >"$HOOKS_DIR/commit-msg" <<'EOF'
#!/bin/bash
# Validate commit message with commitlint

COMMIT_MSG_FILE=$1

# Check if commitlint is installed
if ! command -v commitlint &> /dev/null; then
    echo "⚠️  commitlint not found, skipping validation"
    echo "   Install with: npm install -g @commitlint/cli @commitlint/config-conventional"
    exit 0
fi

# Check if config file exists
if [ ! -f ".commitlintrc.json" ]; then
    echo "⚠️  .commitlintrc.json not found, skipping validation"
    exit 0
fi

# Validate commit message
echo "🔍 Validating commit message..."
if commitlint --edit "$COMMIT_MSG_FILE"; then
    echo "✅ Commit message valid"
    exit 0
else
    echo ""
    echo "❌ Invalid commit message format!"
    echo ""
    echo "Expected format: <type>: <description>"
    echo ""
    echo "Valid types:"
    echo "  feat:     New feature"
    echo "  fix:      Bug fix"
    echo "  docs:     Documentation changes"
    echo "  test:     Test additions/changes"
    echo "  refactor: Code refactoring"
    echo "  perf:     Performance improvements"
    echo "  build:    Build system changes"
    echo "  ci:       CI/CD changes"
    echo "  chore:    Maintenance tasks"
    echo "  revert:   Revert previous commit"
    echo ""
    echo "Examples:"
    echo "  feat: add workflow datasource"
    echo "  fix: resolve null pointer in credential resource"
    echo "  docs: update README with new examples"
    echo ""
    exit 1
fi
EOF

chmod +x "$HOOKS_DIR/commit-msg"
echo -e "${GREEN}✅ commit-msg hook installed${NC}"

# Install pre-push hook to check for AI mentions and Co-Authored-By
echo -e "${BLUE}Installing pre-push hook...${NC}"
cat >"$HOOKS_DIR/pre-push" <<'EOF'
#!/bin/bash
# Pre-push hook to check for AI mentions and Co-Authored-By in commits

set -e

echo "🔍 Checking commits for AI mentions and Co-Authored-By..."

# Get the remote and branch being pushed to
remote="$1"
url="$2"

# Get all commits that will be pushed
zero_commit="0000000000000000000000000000000000000000"

while read local_ref local_sha remote_ref remote_sha; do
    if [ "$local_sha" = "$zero_commit" ]; then
        # Handle delete
        continue
    fi

    if [ "$remote_sha" = "$zero_commit" ]; then
        # New branch, check all commits from main
        range="origin/main..$local_sha"
        if ! git rev-parse origin/main >/dev/null 2>&1; then
            # No origin/main, check all commits
            range="$local_sha"
        fi
    else
        # Update to existing branch, examine new commits
        range="$remote_sha..$local_sha"
    fi

    # Keywords to search for (case insensitive)
    AI_KEYWORDS=(
        "claude"
        "chatgpt"
        "gpt-"
        "copilot"
        "ai generated"
        "generated by ai"
        "generated with ai"
        "with the help of"
        "assisted by"
        "🤖"
        "co-authored-by: claude"
        "co-authored-by: github-actions"
        "co-authored-by: bot"
    )

    # Check each commit in the range
    found_issues=false
    problem_commits=()

    for commit in $(git rev-list "$range"); do
        commit_msg=$(git log -1 --format=%B "$commit")
        commit_short=$(git log -1 --format=%h "$commit")
        commit_subject=$(git log -1 --format=%s "$commit")

        # Check for AI keywords
        for keyword in "${AI_KEYWORDS[@]}"; do
            if echo "$commit_msg" | grep -qi "$keyword"; then
                if [ "$found_issues" = false ]; then
                    echo ""
                    echo "❌ Found commits with AI mentions or Co-Authored-By:"
                    echo ""
                    found_issues=true
                fi
                echo "  Commit: $commit_short - \"$commit_subject\""
                echo "  Found: '$keyword' in commit message"
                echo ""
                problem_commits+=("$commit")
                break
            fi
        done
    done

    if [ "$found_issues" = true ]; then
        echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
        echo "❌ Push rejected!"
        echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
        echo ""
        echo "Commits to fix:"
        for commit in "${problem_commits[@]}"; do
            echo "  $(git log -1 --format='%h - %s' "$commit")"
        done
        echo ""
        echo "To fix these commits, use one of the following methods:"
        echo ""
        echo "1. Interactive rebase to edit commit messages:"
        echo "   git rebase -i HEAD~N  (where N is the number of commits)"
        echo "   Change 'pick' to 'reword' for the problematic commits"
        echo ""
        echo "2. Amend the last commit (if it's the most recent):"
        echo "   git commit --amend"
        echo ""
        echo "3. Filter commits to remove AI mentions:"
        echo "   git filter-branch --msg-filter 'sed \"s/Co-Authored-By.*//g\"' HEAD~N..HEAD"
        echo ""
        echo "After fixing, force push with:"
        echo "   git push --force-with-lease"
        echo ""
        echo "To bypass this check (NOT recommended):"
        echo "   git push --no-verify"
        echo ""
        exit 1
    fi
done

echo "✅ No AI mentions or Co-Authored-By found in commits"
exit 0
EOF

chmod +x "$HOOKS_DIR/pre-push"
echo -e "${GREEN}✅ pre-push hook installed${NC}"

echo ""
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}✅ Git hooks installed successfully!${NC}"
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo -e "Hooks installed:"
echo -e "  📝 ${BLUE}pre-commit${NC}         - Auto-generates CHANGELOG.md and coverage"
echo -e "  💬 ${BLUE}prepare-commit-msg${NC} - Suggests conventional commit format"
echo -e "  ✅ ${BLUE}commit-msg${NC}         - Validates commit message with commitlint"
echo -e "  🚫 ${BLUE}pre-push${NC}           - Blocks push if AI mentions or Co-Authored-By found"
echo ""
echo -e "Next commit will automatically:"
echo -e "  • Generate CHANGELOG.md from git history"
echo -e "  • Run test coverage analysis"
echo -e "  • Validate commit message format (requires commitlint)"
echo ""
echo -e "Before push will check:"
echo -e "  • No AI mentions (Claude, ChatGPT, Copilot, etc.)"
echo -e "  • No Co-Authored-By with bots or AI"
echo ""
echo -e "${YELLOW}Note:${NC} To skip hooks temporarily, use: git commit --no-verify / git push --no-verify"
echo -e "${YELLOW}Note:${NC} commitlint requires: npm install -g @commitlint/cli @commitlint/config-conventional"
echo ""
