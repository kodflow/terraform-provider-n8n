<!-- markdownlint-disable MD043 -->

# N8N Terraform Provider - Claude Configuration

This document provides development guidelines for the n8n terraform provider project, a Terraform provider for managing n8n workflows and resources.

## Key Absolute Rules

The configuration establishes three critical prohibitions:

- No markdown files outside `/workspace/README.md`, `/workspace/COVERAGE.MD`, and `/workspace/CLAUDE.md`
- No generated reports or documentation in folders (except COVERAGE.MD at root)
- Generated documentation files: `docs/` (Terraform docs), `examples/nodes/README.md` (nodes catalog)
- Documentation updates must follow Conventional Commits format

## Mandatory Workflow

Each development iteration requires:

1. Write/modify code
2. Run `make fmt` to format all files (goimports, go fmt, gazelle, buildifier, etc.)
3. Run `make test` to execute all tests - ALL tests must pass
4. Run `make lint` and fix ALL errors, warnings, and info messages - NO exceptions
5. Verify test coverage is maintained or improved
6. Run `make docs` to regenerate all documentation (COVERAGE.MD, Terraform docs, nodes README)
7. Remove temporary files (bazel-\*, \*.out, \*.html)
8. **Repeat until EVERYTHING is perfect** - zero errors, zero warnings

## Self-Verification Checklist

Before completing tasks, Claude must verify:

- All tests pass: `make test`
- Code is properly formatted: `make fmt`
- Linters pass: `make lint` (golangci-lint + ktn-linter)
- BUILD.bazel files are up to date (gazelle runs in `make fmt`)
- Provider builds successfully: `make build`
- No uncommitted changes remain

## Code Quality Standards

**CRITICAL: Execute `make lint` and fix ALL errors and warnings reported by ktn-linter and golangci-lint.**

- **NEVER** consider any linting error as a false positive
- **NEVER** consider any linting warning as optional
- If ktn-linter or golangci-lint reports it, it MUST be fixed
- No exceptions, no debates - fix everything until `make lint` is clean

## Testing Standards

### Pure Unit Tests with Maximum Coverage

- **NEVER EVER use t.Skip()** - if a test is hard, mock everything needed
- **ALWAYS write pure unit tests** - mock all external dependencies (HTTP, databases, APIs, etc.)
- **ALWAYS aim for maximum coverage** - every line, every branch, every error path
- **NEVER stop until coverage is maximized** - if coverage is not 100%, keep adding tests
- Mock APIClient, HTTP responses, file systems, time, random - EVERYTHING
- Test all error paths, edge cases, nil checks, boundary conditions
- If something seems impossible to test, you're not mocking enough

## Project Structure

```text
/workspace/
├── src/
│   ├── cmd/                          # CLI entry point
│   ├── internal/
│   │   └── provider/                 # Terraform provider implementation
│   │       ├── credential/           # Credential resource
│   │       ├── execution/            # Execution resources
│   │       ├── project/              # Project resource
│   │       ├── sourcecontrol/        # Source control resource
│   │       ├── tag/                  # Tag resource
│   │       ├── user/                 # User resource
│   │       ├── variable/             # Variable resource
│   │       ├── workflow/             # Workflow resource
│   │       └── shared/               # Shared utilities
│   └── BUILD.bazel                   # Bazel build configuration
├── sdk/n8nsdk/                       # Generated N8N SDK
├── codegen/                          # SDK generation scripts
├── examples/                         # Terraform examples
├── scripts/                          # Build and automation scripts
├── Makefile                          # Main development commands
├── WORKSPACE                         # Bazel workspace
└── BUILD.bazel                       # Root build file
```

## Make Commands

- `make help` - Display all available commands
- `make build` - Build and install the provider
- `make test` - Run the full test suite
- `make fmt` - Format all source files (Go, Bazel, Shell, Markdown, etc.)
- `make lint` - Run code linters (golangci-lint + ktn-linter)
- `make docs` - Generate ALL documentation (Terraform docs + COVERAGE.MD + examples/nodes/README.md)
- `make openapi` - Regenerate SDK from n8n OpenAPI spec (includes fmt)
- `make update` - Update ktn-linter to latest version

## SDK Generation

The project uses OpenAPI Generator to create the N8N SDK:

1. Downloads OpenAPI spec from GitHub
2. Bundles YAML files
3. Applies patches
4. Fixes schema aliases
5. Generates Go SDK
6. Updates BUILD files with gazelle
7. Formats code automatically

Run `make openapi` to regenerate the SDK after n8n API changes.

## Testing Strategy

- Unit tests for all resources and data sources
- Mock HTTP clients for testing provider logic
- Integration tests with real n8n instance (optional): `make test/n8n`
- Test coverage tracked in COVERAGE.MD

## Git Workflow

- Use Conventional Commits format: `type(scope): description`
- Types: feat, fix, docs, style, refactor, test, chore
- Hooks validate commit messages and run tests
- CHANGELOG.md auto-generated from commit history
- No AI mentions or Co-Authored-By in commit messages

## Autonomy and Decision Making

**Be autonomous and make the best decisions:**

- If you have doubts about implementation details, make the best technical decision
- If multiple approaches exist, choose the most maintainable and performant one
- **NEVER stop until the objective is fully satisfied**
- **NEVER give up because of complexity or uncertainty**
- Research, analyze, and solve problems independently

**CRITICAL GIT RULES:**

- **NEVER perform a git reset without explicit user confirmation**
- **NEVER discard work without asking the user first**
- **ALWAYS commit work progressively** to avoid data loss
- If something goes wrong, ask the user before reverting changes

## Common Tasks

### Adding a new resource

1. Create resource file in `src/internal/provider/<name>/`
2. Implement CRUD operations
3. Add comprehensive tests
4. Update BUILD.bazel with `make fmt` (gazelle)
5. Run tests: `make test`
6. Test with real provider: `make test/n8n`

### Fixing linting errors

1. Run `make lint` to see issues
2. Fix code quality issues (godot, thelper, unused, etc.)
3. Re-run `make lint` until clean
4. Run `make test` to ensure no breakage

### Updating documentation

1. Make code changes with proper commit messages
2. Run `make docs` to regenerate ALL documentation:
   - `COVERAGE.MD` - Test coverage report
   - `docs/**/*.md` - Terraform provider documentation (20 files)
   - `examples/nodes/README.md` - Complete nodes catalog (296 nodes)
3. Commit documentation updates
