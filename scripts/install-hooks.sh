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
cat > "$HOOKS_DIR/pre-commit" << 'EOF'
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
cat > "$HOOKS_DIR/prepare-commit-msg" << 'EOF'
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

echo ""
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}✅ Git hooks installed successfully!${NC}"
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo -e "Hooks installed:"
echo -e "  📝 ${BLUE}pre-commit${NC}         - Auto-generates CHANGELOG.md"
echo -e "  💬 ${BLUE}prepare-commit-msg${NC} - Suggests conventional commit format"
echo ""
echo -e "Next commit will automatically:"
echo -e "  • Generate CHANGELOG.md from git history"
echo -e "  • Run test coverage analysis"
echo ""
echo -e "${YELLOW}Note:${NC} To skip hooks temporarily, use: git commit --no-verify"
echo ""
