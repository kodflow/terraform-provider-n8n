<!-- updated: 2026-02-23T12:00:00Z -->

# DevContainer Configuration

## Purpose

Development container setup for consistent dev environments across languages.

## Structure

```text
.devcontainer/
├── devcontainer.json    # Main config
├── docker-compose.yml   # Multi-service setup
├── Dockerfile           # Extends images/ base
├── install.sh           # Standalone Claude installer
├── scripts/             # Build utilities
│   └── generate-assets-archive.sh
├── features/            # Language & tool features
│   ├── languages/       # 25 languages + shared/
│   ├── architectures/   # 14 architecture patterns
│   ├── claude/          # Standalone Claude feature
│   └── kubernetes/      # Local K8s via kind
├── hooks/               # Lifecycle scripts (delegation stubs)
└── images/              # Docker base image + Claude config
```

## Key Files

- `devcontainer.json`: VS Code devcontainer config
- `docker-compose.yml`: Services (app, MCP servers)
- `.env`: Environment variables (git-ignored)
- `scripts/generate-assets-archive.sh`: Generates `claude-assets.tar.gz` (published via GitHub Releases)

## Usage

Features are enabled in `devcontainer.json` under `features`. Language conventions are enforced by specialist agents (e.g., `developer-specialist-go`).
