#!/bin/bash
set -e

VAULT_ID="ypahjj334ixtiyjkytu5hij2im"
MCP_TPL="/workspace/.devcontainer/mcp.json.tpl"
MCP_OUTPUT="/workspace/.devcontainer/mcp.json"

# Initialize tokens
GITHUB_TOKEN=""

# Try 1Password if OP_SERVICE_ACCOUNT_TOKEN is defined
if [ -n "$OP_SERVICE_ACCOUNT_TOKEN" ] && command -v op &>/dev/null; then
  echo "🔐 Retrieving secrets from 1Password..."

  echo "  → Retrieving GitHub token..."
  GITHUB_TOKEN=$(op item get "mcp-github" --vault "$VAULT_ID" --fields credential --reveal 2>/dev/null || echo "")
fi

# Use environment variables as fallback
if [ -z "$GITHUB_TOKEN" ] && [ -n "$GITHUB_API_TOKEN" ]; then
  echo "📌 Using GitHub token from GITHUB_API_TOKEN"
  GITHUB_TOKEN="$GITHUB_API_TOKEN"
fi

# Display warnings only if no token was found
if [ -z "$GITHUB_TOKEN" ]; then
  echo "⚠️  GitHub token not available"
fi

# Generate mcp.json file from template
echo "📝 Generating mcp.json file..."
mkdir -p "$(dirname "$MCP_OUTPUT")"
sed "s|{{ with secret \"secret/mcp/github\" }}{{ .Data.data.token }}{{ end }}|${GITHUB_TOKEN}|g" "$MCP_TPL" \
    >"$MCP_OUTPUT"

echo "✅ mcp.json file generated successfully!"

# Configure Claude CLI settings
echo "⚙️  Configuring Claude CLI..."
cat >/home/vscode/.claude/settings.json <<'EOF'
{
  "enableAllProjectMcpServers": true,
  "alwaysThinkingEnabled": true
}
EOF
echo "✅ Claude CLI settings configured!"
