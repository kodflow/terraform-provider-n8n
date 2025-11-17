<!-- markdownlint-disable MD043 -->

# Terraform Provider for n8n

[![Bazel](https://img.shields.io/badge/Build-Bazel%209.0-43A047?logo=bazel)](https://bazel.build/)
[![Go](https://img.shields.io/badge/Go-1.25.4-00ADD8?logo=go)](https://pkg.go.dev/github.com/kodflow/terraform-provider-n8n)
[![Terraform Registry](https://img.shields.io/badge/dynamic/json?url=https://registry.terraform.io/v1/providers/kodflow/n8n&query=$.version&label=terraform&logo=terraform&color=7B42BC)](https://registry.terraform.io/providers/kodflow/n8n/latest)
[![CI](https://github.com/kodflow/terraform-provider-n8n/actions/workflows/ci.yml/badge.svg)](https://github.com/kodflow/terraform-provider-n8n/actions/workflows/ci.yml)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/6ad65f0b28b64849ad2799943e8ad338)](https://app.codacy.com/gh/kodflow/terraform-provider-n8n/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade)
[![Codacy Badge](https://app.codacy.com/project/badge/Coverage/6ad65f0b28b64849ad2799943e8ad338)](https://app.codacy.com/gh/kodflow/terraform-provider-n8n/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_coverage)

Manage your n8n workflows, credentials, and resources as code with Terraform.

## Features

### Complete n8n Node Support

✅ **296 n8n nodes** fully supported with comprehensive workflow composition:

| Category        | Count | Description                                                         |
| --------------- | ----- | ------------------------------------------------------------------- |
| **Core**        | 5     | Essential workflow building blocks (Code, If, Merge, Set, Switch)   |
| **Trigger**     | 25    | Event-based workflow initiators (Webhook, Manual, Cron, etc.)       |
| **Integration** | 266   | Third-party service integrations (Slack, GitHub, AWS, Google, etc.) |

Each node includes:

- ✅ Complete Terraform workflow example
- ✅ Full lifecycle testing (init/plan/apply/destroy)
- ✅ 100% test pass rate with real n8n validation
- ✅ Ready-to-use templates

**📊 [View test coverage →](COVERAGE.MD)** | **📚 [Browse all nodes →](examples/nodes/)**

### Community Edition Resources

Full support for **n8n Community Edition** (free/self-hosted):

| Resource                  | Status | Description                             |
| ------------------------- | ------ | --------------------------------------- |
| `n8n_workflow`            | ✅     | Create and manage workflows             |
| `n8n_workflow_node`       | ✅     | Modular node composition                |
| `n8n_workflow_connection` | ✅     | Connect nodes in workflows              |
| `n8n_credential`          | ✅     | Store API credentials securely          |
| `n8n_tag`                 | ✅     | Organize resources with tags            |
| `n8n_variable`            | ✅     | Manage environment variables            |
| `n8n_execution`           | ✅     | Query workflow executions (data source) |

### Enterprise Edition Resources

Enterprise features require n8n Enterprise license:

| Resource             | Status            | License Required |
| -------------------- | ----------------- | ---------------- |
| `n8n_project`        | 🚧 In Development | Enterprise       |
| `n8n_user`           | 🚧 In Development | Enterprise       |
| `n8n_source_control` | 🚧 In Development | Enterprise       |

## Quick Start

### Installation

#### Via Terraform Registry

```hcl
terraform {
  required_providers {
    n8n = {
      source  = "kodflow/n8n"
      version = "~> 1.0"
    }
  }
}

provider "n8n" {
  base_url = "https://your-n8n-instance.com"
  api_key  = var.n8n_api_key
}
```

#### Local Development

```bash
make build
# Provider installed at: ~/.terraform.d/plugins/registry.terraform.io/kodflow/n8n/
```

### Get Your n8n API Key

1. Open your n8n instance
2. Go to **Settings** > **API**
3. Click **Create API Key**
4. Set as `N8N_API_KEY` environment variable

### Run Your First Example

```bash
export N8N_API_URL="http://localhost:5678"
export N8N_API_KEY="your-api-key"

cd examples/community/workflows/basic-workflow
terraform init
terraform apply
```

## Examples

### Node Workflow Examples

Browse **296 complete workflow examples** in [`examples/nodes/`](examples/nodes/):

- **[Core Nodes](examples/nodes/#core-nodes-5)** - Essential workflow building blocks
- **[Trigger Nodes](examples/nodes/#trigger-nodes-25)** - Event-based workflow initiators
- **[Integration Nodes](examples/nodes/#integration-nodes-266)** - Third-party service integrations
- **[MEGA Workflow](examples/mega-workflow/)** - All 296 nodes in a single workflow

Each node example includes `main.tf`, `variables.tf`, and documentation.

**Testing status:** ✅ All 296 workflows tested and passing

### Community Edition Examples

Browse complete examples in [`examples/community/`](examples/community/):

- **[Workflows](examples/community/workflows/)** - Basic webhook and scheduled workflows
- **[Credentials](examples/community/credentials/)** - HTTP Basic Auth and API credentials
- **[Tags](examples/community/tags/)** - Workflow organization with tags
- **[Variables](examples/community/variables/)** - Environment variable management
- **[Executions](examples/community/executions/)** - Query and filter workflow executions

### Comprehensive Examples

Production-ready examples at [`examples/comprehensive/`](examples/comprehensive/):

- **[Complete Modular Workflow](examples/comprehensive/complete-modular-workflow/)** - Advanced multi-node workflow with error handling

## Development

### Prerequisites

**Recommended:** Use the included DevContainer with all tools pre-installed:

- Go 1.25.4
- Bazel 9.0
- Terraform & OpenTofu
- All development tools

Just open the project in VS Code and rebuild the container.

### Essential Commands

```bash
make help          # Display all available commands
make build         # Build and install provider locally
make test          # Run full test suite
make fmt           # Format all source files
make lint          # Run code linters (zero tolerance)
make docs          # Generate CHANGELOG.md and COVERAGE.MD
```

### Node Management

```bash
make nodes                   # Synchronize n8n nodes from official repository
make nodes/fetch             # Fetch latest n8n repository
make nodes/parse             # Parse nodes and generate registry
make nodes/workflows         # Generate 296 per-node workflow examples
make nodes/mega-workflow     # Generate MEGA workflow with all 296 nodes
make nodes/validate-coverage # Validate test coverage completeness
make nodes/docs              # Generate node documentation
make nodes/stats             # Display node statistics
```

### Quality Standards

**Critical requirements:**

- ✅ All tests must pass: `make test`
- ✅ Code must be formatted: `make fmt`
- ✅ Zero linting errors: `make lint`
- ✅ Maximum test coverage (no `t.Skip()` allowed)

### SDK Generation

Auto-generate Go SDK from n8n OpenAPI specification:

```bash
make openapi       # Download and prepare n8n OpenAPI spec
make sdk           # Generate Go SDK from OpenAPI spec
```

### Git Workflow

Git hooks enforce quality:

- **Pre-commit**: Formats code, generates docs, validates changes
- **Commit-msg**: Validates commit message format
- **Pre-push**: Runs tests before pushing

Hooks are automatically installed in DevContainer.

## Project Architecture

```text
.
├── src/                          # Provider source code
│   ├── main.go                   # Entry point
│   └── internal/provider/        # Provider implementation
│       ├── credential/           # Credential resource
│       ├── execution/            # Execution data source
│       ├── tag/                  # Tag resource
│       ├── variable/             # Variable resource
│       ├── workflow/             # Workflow resource
│       └── shared/               # Shared utilities
├── sdk/n8nsdk/                   # Auto-generated n8n SDK
├── codegen/                      # SDK generation scripts
├── examples/                     # Terraform examples
│   ├── nodes/                    # 296 node examples
│   ├── community/                # Community edition examples
│   ├── comprehensive/            # Production-ready examples
│   └── mega-workflow/            # All nodes in one workflow
├── scripts/                      # Build and automation scripts
├── Makefile                      # Main development commands
└── .devcontainer/                # DevContainer configuration
```

## Release Process

Releases are fully automated via GitHub Actions using semantic versioning:

- `feat:` → Minor version bump (v0.1.0 → v0.2.0)
- `fix:` → Patch version bump (v0.1.0 → v0.1.1)
- `BREAKING CHANGE:` → Major version bump (v0.1.0 → v1.0.0)

The CI/CD pipeline automatically:

1. Analyzes commit messages
2. Determines next version
3. Updates CHANGELOG.md
4. Creates signed tags
5. Compiles multi-platform binaries
6. Generates checksums and signatures
7. Creates GitHub Release

View all releases at [GitHub Releases](../../releases).

## Contributing

Contributions are welcome! Follow these steps:

1. Fork the repository
2. Create a feature branch: `git checkout -b feat/my-feature`
3. Develop and test: `make test && make lint && make build`
4. Commit your changes (follow [Conventional Commits](https://www.conventionalcommits.org/))
5. Push and create a Pull Request

### Contribution Standards

- ✅ All tests must pass
- ✅ Code must be formatted and linted
- ✅ Tests required for new features
- ✅ Follow Conventional Commits format

See [CLAUDE.md](CLAUDE.md) for detailed development guidelines.

## Support This Project

If you find this project useful, consider sponsoring its development:

- ❤️ [GitHub Sponsors](https://github.com/sponsors/kodflow)
- ☕ [Ko-fi](https://ko-fi.com/kodflow)

Your support helps:

- ⏰ Dedicate more time to development and maintenance
- 🐛 Fix bugs faster and implement new features
- 📚 Improve documentation and examples
- 🆘 Provide better community support

Every contribution makes a difference! Thank you! 🙏

## License

Sustainable Use License 1.0

See [LICENSE](LICENSE) for details.

---

**Developed with ❤️ by [KodFlow](https://github.com/kodflow)**
