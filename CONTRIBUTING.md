# Contributing to n8n Terraform Provider

Thank you for your interest in contributing! This guide will help you get started with development.

## Prerequisites

### Recommended: DevContainer

The easiest way to start is using the included DevContainer with all tools pre-installed:

- Go 1.25.4
- Bazel 9.0
- Terraform & OpenTofu
- All development tools (goimports, buildifier, prettier, etc.)
- Pre-configured Git hooks

Just open the project in VS Code and rebuild the container.

### Manual Setup

If you prefer manual setup, install:

- [Go 1.25.4+](https://go.dev/dl/)
- [Bazel 9.0+](https://bazel.build/)
- [Terraform](https://www.terraform.io/downloads) or [OpenTofu](https://opentofu.org/)
- Development tools: `make tools`

## Essential Commands

### Build & Test

```bash
make help          # Display all available commands
make build         # Build and install provider locally
make test          # Run full test suite (unit + acceptance)
make test/unit     # Run unit tests only
make fmt           # Format all source files
make lint          # Run code linters (zero tolerance)
make docs          # Generate CHANGELOG.md and COVERAGE.MD
```

### Update Dependencies

```bash
make update        # Update ALL (n8n SDK + ktn-linter + README badge)
```

This command:

1. Updates n8n commit to latest version
2. Updates ktn-linter to latest version
3. Updates n8n version badge in README
4. Regenerates SDK and documentation

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

### SDK Generation

Auto-generate Go SDK from n8n OpenAPI specification:

```bash
make openapi       # Download and prepare n8n OpenAPI spec
make sdk           # Generate Go SDK from OpenAPI spec
```

## Quality Standards

**Critical requirements** - ALL must pass:

- ✅ **All tests must pass**: `make test`
- ✅ **Code must be formatted**: `make fmt`
- ✅ **Zero linting errors**: `make lint`
- ✅ **Maximum test coverage**: No `t.Skip()` allowed

### Testing Standards

- **NEVER use `t.Skip()`** - Mock all dependencies instead
- **Write pure unit tests** - No external dependencies (HTTP, databases, APIs)
- **Aim for maximum coverage** - Every line, every branch, every error path
- **Mock everything** - APIClient, HTTP responses, file systems, time, random
- **Test all error paths** - Edge cases, nil checks, boundary conditions

## Git Workflow

### Hooks

Git hooks are automatically installed in the DevContainer and enforce quality:

- **Pre-commit**: Formats code, generates docs, validates changes
- **Commit-msg**: Validates commit message format
- **Pre-push**: Runs tests before pushing

### Commit Message Format

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

**Types:**

- `feat:` - New feature (minor version bump)
- `fix:` - Bug fix (patch version bump)
- `docs:` - Documentation only
- `style:` - Code style (formatting, semicolons, etc.)
- `refactor:` - Code refactoring
- `test:` - Adding or updating tests
- `chore:` - Maintenance tasks
- `ci:` - CI/CD changes
- `BREAKING CHANGE:` - Breaking compatibility (major version bump)

**Examples:**

```bash
feat(workflow): add support for error handling nodes
fix(credential): correct basic auth encoding
docs: update README with new examples
```

### Branch Naming

Use descriptive branch names:

- `feat/description` - New features
- `fix/description` - Bug fixes
- `docs/description` - Documentation updates
- `refactor/description` - Code refactoring

## Development Workflow

### 1. Fork & Clone

```bash
git clone https://github.com/YOUR_USERNAME/terraform-provider-n8n.git
cd terraform-provider-n8n
```

### 2. Create Branch

```bash
git checkout -b feat/my-feature
```

### 3. Make Changes

```bash
# Edit code
vim src/internal/provider/workflow/resource.go

# Format
make fmt

# Test
make test

# Lint
make lint
```

### 4. Commit & Push

```bash
git add .
git commit -m "feat(workflow): add new feature"
git push origin feat/my-feature
```

### 5. Create Pull Request

Create a PR on GitHub with:

- Clear description of changes
- Link to related issues
- Screenshots/examples if applicable

## Project Structure

```text
.
├── src/                          # Provider source code
│   ├── main.go                   # Entry point
│   └── internal/provider/        # Provider implementation
│       ├── credential/           # Credential resource
│       ├── execution/            # Execution data source
│       ├── project/              # Project resource (Enterprise)
│       ├── sourcecontrol/        # Source control resource (Enterprise)
│       ├── tag/                  # Tag resource
│       ├── user/                 # User resource (Enterprise)
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
├── CLAUDE.md                     # AI assistant guidelines
└── .devcontainer/                # DevContainer configuration
```

## Adding a New Resource

### 1. Create Resource File

```bash
# Create resource structure
mkdir -p src/internal/provider/my_resource
touch src/internal/provider/my_resource/resource.go
touch src/internal/provider/my_resource/resource_test.go
```

### 2. Implement CRUD Operations

See existing resources for examples:

- `src/internal/provider/workflow/resource.go`
- `src/internal/provider/credential/resource.go`

### 3. Add Comprehensive Tests

```bash
# Unit tests with 100% coverage
# Mock all external dependencies
# Test all error paths
```

### 4. Update BUILD Files

```bash
make fmt  # Gazelle will update BUILD.bazel files automatically
```

### 5. Test & Lint

```bash
make test
make lint
make build
```

## Getting Help

- **Issues**: [GitHub Issues](https://github.com/kodflow/terraform-provider-n8n/issues)
- **Discussions**: [GitHub Discussions](https://github.com/kodflow/terraform-provider-n8n/discussions)
- **Documentation**: [CLAUDE.md](CLAUDE.md) - Detailed development guidelines

## Code of Conduct

Be respectful, collaborative, and professional. We're all here to build something great together.

---

Thank you for contributing! 🎉
