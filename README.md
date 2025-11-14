# Terraform Provider for n8n

Terraform provider to manage n8n resources (workflows, credentials, projects, users, and more).

[![Bazel](https://img.shields.io/badge/Build-Bazel%209.0-43A047?logo=bazel)](https://bazel.build/)
[![Go](https://img.shields.io/badge/Go-1.24-00ADD8?logo=go)](https://go.dev/)
[![Terraform](https://img.shields.io/badge/Terraform-Plugin%20Framework-7B42BC?logo=terraform)](https://developer.hashicorp.com/terraform/plugin/framework)

## Features

### Community Edition Support

The provider fully supports **n8n Community Edition** (free/self-hosted):

| Resource/Data Source | Status | Description |
|---------------------|---------|-------------|
| `n8n_workflow` | ✅ Available | Create and manage workflows |
| `n8n_credential` | ✅ Available | Store API credentials securely |
| `n8n_tag` | ✅ Available | Organize resources with tags |
| `n8n_variable` | ✅ Available | Manage environment variables |
| `n8n_execution` | ✅ Available | Query workflow executions (read-only) |

### Enterprise Edition Support

**Enterprise features** require an n8n Enterprise license:

| Resource/Data Source | Status | License Required |
|---------------------|---------|-----------------|
| `n8n_project` | 🚧 In Development | Enterprise |
| `n8n_user` | 🚧 In Development | Enterprise |
| `n8n_source_control` | 🚧 In Development | Enterprise |

> **Note:** Enterprise features are in development and will be available once enterprise license access is obtained for testing.

## Prerequisites

- **Go 1.24.0+**
- **Bazel 9.0+**
- **Terraform 1.0+** or **OpenTofu 1.0+**

### DevContainer (Recommended)

The project includes a preconfigured DevContainer with all tools (Go 1.25.3, Bazel 9.0, Terraform, OpenTofu).

**Quick start:**
1. Open in VS Code
2. Accept "Reopen in Container" prompt
3. Start developing

### Manual Installation

```bash
# Install Bazelisk (Bazel version manager)
go install github.com/bazelbuild/bazelisk@latest

# Verify versions
go version    # go1.24 or higher
bazel version # Bazel 9.0+
```

## Installation

### Via Terraform Registry (Coming Soon)

```hcl
terraform {
  required_providers {
    n8n = {
      source  = "kodflow/n8n"
      version = "~> 0.1.0"
    }
  }
}

provider "n8n" {
  api_url = "https://your-n8n-instance.com"
  api_key = var.n8n_api_key
}
```

### Local Development

```bash
make build
# Provider installed at: ~/.terraform.d/plugins/registry.terraform.io/kodflow/n8n/
```

## Quick Start

### Get Your API Key

1. Open your n8n instance
2. Go to **Settings** > **API**
3. Click **Create API Key**
4. Set as `N8N_API_KEY` environment variable

### Run Examples

```bash
export N8N_API_URL="http://localhost:5678"
export N8N_API_KEY="your-api-key"

cd examples/community/workflows/basic-workflow
terraform init
terraform apply
```

See [examples/](examples/) directory for more examples.

## Examples

The provider includes comprehensive examples for different use cases:

### Community Edition Examples

Browse complete examples in [`examples/community/`](examples/community/):

- **[Workflows](examples/community/workflows/)** - Basic webhook and scheduled workflows
- **[Credentials](examples/community/credentials/)** - HTTP Basic Auth and API credentials
- **[Tags](examples/community/tags/)** - Workflow organization with tags
- **[Variables](examples/community/variables/)** - Environment variable management
- **[Executions](examples/community/executions/)** - Query and filter workflow executions

### Enterprise Edition Examples

Enterprise examples are currently in development at [`examples/enterprise/`](examples/enterprise/).

## Development

### Essential Commands

```bash
make help          # Display all available commands
make build         # Build and install provider locally
make test          # Run full test suite
make fmt           # Format all source files
make lint          # Run code linters (zero tolerance)
make docs          # Generate CHANGELOG.md and COVERAGE.MD
make openapi       # Regenerate SDK from n8n OpenAPI spec
```

### Quality Standards

**Critical requirements:**
- ✅ All tests must pass: `make test`
- ✅ Code must be formatted: `make fmt`
- ✅ Zero linting errors: `make lint`
- ✅ Maximum test coverage (no `t.Skip()` allowed)

### SDK Generation

The provider uses auto-generated Go SDK from n8n OpenAPI specification:

```bash
make openapi       # Download and prepare n8n OpenAPI spec
make sdk           # Generate Go SDK from OpenAPI spec
```

**Auto-generated files:**
- `sdk/n8nsdk/*.go` - Go client for n8n API
- `sdk/n8nsdk/api/openapi.yaml` - Patched OpenAPI spec (not committed)

See [`codegen/`](codegen/) for generation scripts and patches.

### Git Workflow

The project uses git hooks for quality enforcement:

- **Pre-commit**: Formats code, generates documentation, validates changes
- **Commit-msg**: Validates commit message format
- **Pre-push**: Runs tests before pushing

Hooks are automatically installed in DevContainer.

## Project Architecture

```
.
├── src/                          # Provider source code
│   ├── main.go                   # Entry point
│   └── internal/provider/        # Provider implementation
│       ├── credential/           # Credential resource
│       ├── execution/            # Execution data source
│       ├── project/              # Project resource (Enterprise)
│       ├── sourcecontrol/        # Source control (Enterprise)
│       ├── tag/                  # Tag resource
│       ├── user/                 # User resource (Enterprise)
│       ├── variable/             # Variable resource
│       ├── workflow/             # Workflow resource
│       └── shared/               # Shared utilities
├── sdk/n8nsdk/                   # Auto-generated n8n SDK
├── codegen/                      # SDK generation scripts
├── examples/                     # Terraform examples
│   ├── community/                # Community edition examples
│   └── enterprise/               # Enterprise edition examples
├── scripts/                      # Build and automation scripts
├── Makefile                      # Main development commands
└── .devcontainer/                # DevContainer configuration
```

## Release Process

Releases are automated via GitHub Actions:

1. Create and push a git tag (e.g., `v0.1.0`)
2. GitHub Actions automatically:
   - Compiles multi-platform binaries
   - Generates checksums and signatures
   - Creates GitHub Release with artifacts

View all releases and binaries at [GitHub Releases](../../releases).

## Contributing

Contributions are welcome! Follow these steps:

1. Fork the repository
2. Create a feature branch: `git checkout -b feat/my-feature`
3. Develop and test: `make test && make lint && make build`
4. Commit your changes (follow conventional commits)
5. Push and create a Pull Request

### Contribution Standards

- ✅ All tests must pass
- ✅ Code must be formatted and linted
- ✅ Tests required for new features
- ✅ Follow [Conventional Commits](https://www.conventionalcommits.org/) format

See [CLAUDE.md](CLAUDE.md) for detailed development guidelines.

## Dependencies

### Main Dependencies

- `github.com/hashicorp/terraform-plugin-framework` v1.16.1 - Terraform provider framework
- `github.com/hashicorp/terraform-plugin-docs` v0.24.0 - Documentation generation
- Local module: `github.com/kodflow/n8n/sdk/n8nsdk` - Auto-generated n8n API client

### Build Tools

- **Bazel 9.0.0rc1** - Build system with bzlmod
- **Go 1.24.0** - Programming language
- **OpenAPI Generator 7.11.0** - SDK code generation
- **golangci-lint** + **ktn-linter** - Code quality tools

See `go.mod` and `MODULE.bazel` for complete dependencies.

## Troubleshooting

### Bazel Issues

```bash
bazel clean --expunge        # Clean cache
bazel version                # Verify version 9.0+
cat .bazelversion            # Check configured version
```

### Test Failures

```bash
bazel test --test_output=all //src/...              # Detailed output
bazel test --test_verbose_timeout_warnings //src/... # Verbose logs
```

### DevContainer Issues

Use VS Code Command Palette:
- "Dev Containers: Rebuild Container"
- "Dev Containers: Show Log"

## License

MPL-2.0

---

**Developed with ❤️ by [KodFlow](https://github.com/kodflow)**
