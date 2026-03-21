<!-- updated: 2026-02-27T13:00:00Z -->

# DevContainer Features

## Purpose

Modular features for languages, tools, and architectures.

## Structure

```text
features/
├── languages/      # 25 languages + shared utility library
│   ├── shared/     # feature-utils.sh (colors, logging, arch, GitHub API)
│   └── <lang>/     # install.sh + devcontainer-feature.json
├── architectures/  # Architecture patterns (14 patterns)
├── claude/         # Claude Code standalone integration
└── kubernetes/     # Local K8s via kind
```

## Key Components

- Each language has `install.sh` + `devcontainer-feature.json`
- All install.sh source `shared/feature-utils.sh` (with inline fallback)
- Downloads parallelized with `&` + `wait` for faster builds
- Conventions enforced by specialist agents (e.g., `developer-specialist-go`)

## Adding a Language

1. Create `languages/<name>/`
2. Add `devcontainer-feature.json` for metadata
3. Add `install.sh` sourcing `shared/feature-utils.sh`
