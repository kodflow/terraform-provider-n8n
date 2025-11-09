{
  "mcpServers": {
    "codacy": {
      "command": "npx",
      "args": [
        "-y",
        "@codacy/codacy-mcp@latest"
      ],
      "env": {
        "CODACY_ACCOUNT_TOKEN": "{{ with secret "secret/mcp/codacy" }}{{ .Data.data.token }}{{ end }}"
      }
    },
    "github": {
      "command": "npx",
      "args": [
        "-y",
        "@modelcontextprotocol/server-github"
      ],
      "env": {
        "GITHUB_PERSONAL_ACCESS_TOKEN": "{{ with secret "secret/mcp/github" }}{{ .Data.data.token }}{{ end }}"
      }
    },
    "repomix": {
      "command": "npx",
      "args": [
        "-y",
        "repomix",
        "--mcp"
      ]
    }
  }
}
