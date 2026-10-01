{
  "mcpServers": {
    "github": {
      "command": "npx",
      "args": [
        "-y",
        "@modelcontextprotocol/server-github"
      ],
      "env": {
        "GITHUB_PERSONAL_ACCESS_TOKEN": "{{ with secret "secret/mcp/github" }}{{ .Data.data.token }}{{ end }}"
      }
    }
  }
}
