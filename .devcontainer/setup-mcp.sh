#!/bin/bash
set -e

VAULT_ID="ypahjj334ixtiyjkytu5hij2im"
MCP_TPL="/workspace/.devcontainer/mcp.json.tpl"
MCP_OUTPUT="/workspace/.devcontainer/mcp.json"

# Initialiser les tokens
CODACY_TOKEN=""
GITHUB_TOKEN=""

# Essayer 1Password si OP_SERVICE_ACCOUNT_TOKEN est défini
if [ -n "$OP_SERVICE_ACCOUNT_TOKEN" ] && command -v op &>/dev/null; then
  echo "🔐 Récupération des secrets depuis 1Password..."

  echo "  → Récupération du token Codacy..."
  CODACY_TOKEN=$(op item get "mcp-codacy" --vault "$VAULT_ID" --fields credential --reveal 2>/dev/null || echo "")

  echo "  → Récupération du token GitHub..."
  GITHUB_TOKEN=$(op item get "mcp-github" --vault "$VAULT_ID" --fields credential --reveal 2>/dev/null || echo "")
fi

# Utiliser les variables d'environnement en fallback
if [ -z "$CODACY_TOKEN" ] && [ -n "$CODACY_API_TOKEN" ]; then
  echo "📌 Utilisation du token Codacy depuis CODACY_API_TOKEN"
  CODACY_TOKEN="$CODACY_API_TOKEN"
fi

if [ -z "$GITHUB_TOKEN" ] && [ -n "$GITHUB_API_TOKEN" ]; then
  echo "📌 Utilisation du token GitHub depuis GITHUB_API_TOKEN"
  GITHUB_TOKEN="$GITHUB_API_TOKEN"
fi

# Afficher les avertissements seulement si aucun token n'a été trouvé
if [ -z "$CODACY_TOKEN" ]; then
  echo "⚠️  Token Codacy non disponible"
fi

if [ -z "$GITHUB_TOKEN" ]; then
  echo "⚠️  Token GitHub non disponible"
fi

# Générer le fichier mcp.json à partir du template
echo "📝 Génération du fichier mcp.json..."
mkdir -p "$(dirname "$MCP_OUTPUT")"
sed "s|{{ with secret \"secret/mcp/codacy\" }}{{ .Data.data.token }}{{ end }}|${CODACY_TOKEN}|g" "$MCP_TPL" \
  | sed "s|{{ with secret \"secret/mcp/github\" }}{{ .Data.data.token }}{{ end }}|${GITHUB_TOKEN}|g" \
    >"$MCP_OUTPUT"

echo "✅ Fichier mcp.json généré avec succès!"

# Configurer les paramètres Claude CLI
echo "⚙️  Configuration de Claude CLI..."
cat >/home/vscode/.claude/settings.json <<'EOF'
{
  "enableAllProjectMcpServers": true,
  "alwaysThinkingEnabled": true
}
EOF
echo "✅ Paramètres Claude CLI configurés!"
