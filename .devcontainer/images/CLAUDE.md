<!-- updated: 2026-03-14T12:00:00Z -->

# DevContainer Images

## Purpose

Docker base image with all development tools pre-installed. Claude Code and MCP servers are included; languages added via features.

## Structure

```text
.devcontainer/images/
├── Dockerfile          # Base image (~15 layers, consolidated)
├── .dockerignore       # Build context exclusions
├── mcp.json.tpl        # MCP server template
├── grepai.config.yaml  # GrepAI search config (12 languages)
├── .p10k.zsh           # Powerlevel10k config
├── scripts/vpn/        # VPN helper scripts
├── hooks/              # Image-embedded lifecycle hooks
│   ├── shared/utils.sh # Shared utilities
│   └── lifecycle/      # onCreate, postCreate, postStart, etc.
└── .claude/            # Claude Code configuration
    ├── commands/       # Slash commands (17 skills)
    ├── scripts/        # Hook scripts (31 scripts)
    ├── agents/         # Agent definitions (79 agents)
    ├── docs/           # Design Patterns Knowledge Base (170+ patterns)
    ├── templates/      # Project/docs/terraform templates
    └── settings.json   # Claude settings
```

## Container Paths (Runtime)

| Source (Build)      | Container Path             | Backup Location                  |
| ------------------- | -------------------------- | -------------------------------- |
| `.claude/`          | `/home/vscode/.claude/`    | `/etc/claude-defaults/`          |
| `.claude/commands/` | `~/.claude/commands/`      | `/etc/claude-defaults/commands/` |
| `.claude/scripts/`  | `~/.claude/scripts/`       | `/etc/claude-defaults/scripts/`  |
| `.claude/agents/`   | `~/.claude/agents/`        | `/etc/claude-defaults/agents/`   |
| `.claude/docs/`     | `~/.claude/docs/`          | `/etc/claude-defaults/docs/`     |
| `mcp.json.tpl`      | `/etc/mcp/mcp.json.tpl`    | -                                |
| `hooks/`            | `/etc/devcontainer-hooks/` | -                                |

**Note:** Claude files restored from `/etc/claude-defaults/` at each start via `postStart.sh`. Lifecycle hooks delegate from workspace stubs to
`/etc/devcontainer-hooks/` (image-embedded).

## Design Patterns Knowledge Base

**Container Location:** `~/.claude/docs/` (restored at startup)

170+ pattern files across 19 categories, consulted by `/plan` and `/review`.

| Category                     | Files                              | Examples                      |
| ---------------------------- | ---------------------------------- | ----------------------------- |
| GoF (25 files)               | creational, structural, behavioral | Factory, Observer, Strategy   |
| Architectural (9)            | architectural/                     | MVC, Hexagonal, CQRS          |
| Cloud + Resilience (27)      | cloud/, resilience/                | Circuit Breaker, Saga, Retry  |
| Concurrency (8)              | concurrency/                       | Thread Pool, Actor, Mutex     |
| DDD (8)                      | ddd/                               | Aggregate, Repository, Entity |
| DevOps (14)                  | devops/                            | Feature Toggles, Blue-Green   |
| Enterprise (12)              | enterprise/                        | PoEAA (Martin Fowler)         |
| Functional (5)               | functional/                        | Monad, Either, Lens           |
| Integration + Messaging (15) | integration/, messaging/           | API Gateway, EIP              |
| Performance (8)              | performance/                       | Cache, Lazy Load, Pool        |
| Principles (7)               | principles/, conventions/          | SOLID, DRY, KISS              |
| Security (8)                 | security/                          | OAuth, JWT, RBAC              |
| Testing (8)                  | testing/                           | Mock, Stub, Fixture           |

**Agent usage:** See `.claude/docs/CLAUDE.md`

## Installed Tools

| Category     | Tools                                                                       |
| ------------ | --------------------------------------------------------------------------- |
| Cloud CLIs   | AWS, GCP, Azure, 1Password                                                  |
| IaC          | Terraform, Vault, Consul, Nomad, Packer, Ansible                            |
| Container    | Docker (via feature), kubectl, Helm                                         |
| Network      | ping, dig, nmap, traceroute, mtr, tcpdump, netcat, whois, iperf3, net-tools |
| VPN          | OpenVPN, WireGuard, StrongSwan (IPsec), PPTP                                |
| Code Quality | ShellCheck, ktn-linter, grepai, CodeRabbit, Qodo, RTK                       |
| Shell        | Zsh (default `$SHELL`) + Oh My Zsh + Powerlevel10k                          |

## Shell Startup Optimization (v3)

`~/.devcontainer-env.sh` uses a three-pillar architecture for fast shell startup:

| Pillar                | Mechanism                                                    | Effect       |
| --------------------- | ------------------------------------------------------------ | ------------ |
| Lazy wrappers         | `nvm`, `pyenv`, `rbenv`, `sdk` load on first use             | ~1.4s saved  |
| Cached completions    | `~/.zsh_completions/` via fpath (pre-generated by postStart) | ~750ms saved |
| Dynamic p10k segments | `~/.p10k-segments.zsh` based on installed tools              | ~70ms saved  |

**Phase 1** (always): PATH, env vars, fpath — no subprocesses. **Phase 2** (terminal only): lazy wrappers, aliases, fast `complete -C` for HashiCorp/AWS.

Tool binaries (`node`, `python`, `ruby`, `java`) work immediately via Phase 1 PATH/shims. Management commands (`nvm use`, `pyenv install`) trigger lazy-load on
first call.

## MCP Servers (Runtime)

Configured in `mcp.json.tpl`:

| Server         | Package                                     | Usage                                               | Auth                |
| -------------- | ------------------------------------------- | --------------------------------------------------- | ------------------- |
| **grepai**     | `grepai` (binary)                           | Semantic code search, Call graph                    | None (local)        |
| **context7**   | `@upstash/context7-mcp`                     | Up-to-date documentation for prompts                | None (rate-limited) |
| **GitHub**     | `ghcr.io/github/github-mcp-server` (Docker) | PR, Issues, Repos                                   | `GITHUB_TOKEN`      |
| **GitLab**     | `@zereight/mcp-gitlab`                      | MR, Issues, Pipelines, Wiki                         | `GITLAB_TOKEN`      |
| **Codacy**     | `@codacy/codacy-mcp`                        | Code quality, Security                              | `CODACY_TOKEN`      |
| **Playwright** | `@playwright/mcp`                           | Browser automation, E2E tests                       | None                |
| **Taskmaster** | `task-master-ai`                            | Persistent task planning, dependencies, PRD parsing | None (MCP provider) |

**grepai tools (MANDATORY - use instead of Grep):**

| Tool                   | Description           | Use Case                   |
| ---------------------- | --------------------- | -------------------------- |
| `grepai_search`        | Semantic code search  | Natural language queries   |
| `grepai_trace_callers` | Find function callers | Impact analysis            |
| `grepai_trace_callees` | Find called functions | Dependency analysis        |
| `grepai_trace_graph`   | Build call graph      | Architecture understanding |
| `grepai_index_status`  | Check index health    | Debugging                  |

**GitLab tools (when GITLAB_TOKEN configured):**

| Tool                         | Description              | Use Case          |
| ---------------------------- | ------------------------ | ----------------- |
| `gitlab_list_projects`       | List accessible projects | Project discovery |
| `gitlab_get_project`         | Get project details      | Project info      |
| `gitlab_list_merge_requests` | List MRs                 | Code review       |
| `gitlab_get_merge_request`   | Get MR details           | Review analysis   |
| `gitlab_list_issues`         | List project issues      | Issue tracking    |
| `gitlab_list_pipelines`      | List CI pipelines        | CI/CD status      |

**GitLab env vars:** `GITLAB_TOKEN`, `GITLAB_API_URL` (default: gitlab.com)

**Context7 usage:** Add "use context7" in prompts to fetch up-to-date documentation.

**Taskmaster tools (standard tier - 15 tools):** `initialize_project`, `parse_prd`, `get_tasks`, `get_task`, `add_task`, `set_task_status`, `next_task`,
`expand_task`, `expand_all`, `add_subtask`, `update_subtask`, `remove_task`, `analyze_project_complexity`, `complexity_report`, `generate`

**Playwright capabilities:** `core`, `pdf`, `testing`, `tracing` (headless mode)

## Skills (Slash Commands)

| Skill      | Description                                                         |
| ---------- | ------------------------------------------------------------------- |
| `/init`    | Project initialization check                                        |
| `/plan`    | Planning mode for implementation strategy                           |
| `/do`      | Iterative task execution loop (RLM)                                 |
| `/review`  | AI-powered code review (3-tier: agents + Qodo + CodeRabbit)         |
| `/git`     | Workflow Git automation (commit, push, PR, merge)                   |
| `/search`  | Documentation research with official sources                        |
| `/docs`    | Deep project documentation generation (multi-agent)                 |
| `/test`    | E2E testing with Playwright MCP                                     |
| `/lint`    | Intelligent linting with ktn-linter (148 rules)                     |
| `/infra`   | Infrastructure automation (Terraform/Terragrunt)                    |
| `/secret`  | Secure secret management (1Password + Vault-like paths)             |
| `/vpn`     | Multi-protocol VPN management (OpenVPN, WireGuard, IPsec, PPTP)     |
| `/warmup`  | Context pre-loading and CLAUDE.md update                            |
| `/update`  | DevContainer update from template                                   |
| `/improve` | Continuous docs enhancement & anti-pattern detection                |
| `/feature` | Feature tracking (RTM) with --add, --edit, --del, --list, --checkup |
| `/prompt`  | Generate ideal prompt structure for /plan                           |

## Hooks (31 scripts, 17 event types)

Core hooks (always active):

| Hook                 | Trigger                  | Action                                            |
| -------------------- | ------------------------ | ------------------------------------------------- |
| `git-guard.sh`       | PreToolUse (Bash)        | Block AI commits + secret scan + force-with-lease |
| `rtk-rewrite.sh`     | PreToolUse (Bash)        | Rewrite commands via RTK for token savings        |
| `pre-validate.sh`    | PreToolUse (Write/Edit)  | Protect sensitive files                           |
| `post-edit.sh`       | PostToolUse (Write/Edit) | **Format only** (fast, ~100-500ms)                |
| `on-stop-quality.sh` | Stop (\*)                | **Batch lint + typecheck + test** (deduplicated)  |
| `on-stop.sh`         | Stop (\*)                | Terminal bell + session summary                   |
| `session-init.sh`    | SessionStart (all)       | Cache git metadata as env vars                    |
| `post-compact.sh`    | SessionStart (compact)   | Restore RLM context rules                         |
| `notification.sh`    | Notification (\*)        | Terminal bell + notification log                  |

Full inventory: See `.devcontainer/hooks/CLAUDE.md` and `CLAUDE.md` (root).

**Makefile-first pattern:** Scripts check `make fmt/lint/typecheck/test FILE=<path>` first, then fall back to direct tool invocation.

## Build

```bash
docker buildx build --platform linux/amd64,linux/arm64 \
  -t ghcr.io/kodflow/devcontainer-template:latest .
```
