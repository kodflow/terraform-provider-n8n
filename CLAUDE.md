# N8N Terraform Provider - Claude Configuration

This document provides development guidelines for the n8n terraform provider project, a Terraform provider for managing n8n workflows and resources.

## Key Absolute Rules

The configuration establishes three critical prohibitions:
- No markdown files outside `/workspace/README.md`, `/workspace/CHANGELOG.md`, `/workspace/COVERAGE.MD`, and `/workspace/CLAUDE.md`
- No generated reports or documentation in folders (except CHANGELOG.md and COVERAGE.MD at root)
- Documentation updates must follow Conventional Commits format

## Mandatory Workflow

Each development iteration requires:
1. Write/modify code
2. Run `make fmt` to format all files (goimports, go fmt, gazelle, buildifier, etc.)
3. Run `make test` to execute all tests
4. Run `make lint` and fix all warnings, errors, and info messages
5. Verify test coverage remains high (target: >70%)
6. Update CHANGELOG.md using `make docs` if adding features
7. Remove temporary files (bazel-*, *.out, *.html)
8. Repeat until zero critical errors

## Self-Verification Checklist

Before completing tasks, Claude must verify:
- All tests pass: `make test`
- Code is properly formatted: `make fmt`
- Linters pass: `make lint` (golangci-lint + ktn-linter)
- BUILD.bazel files are up to date (gazelle runs in `make fmt`)
- Provider builds successfully: `make build`
- No uncommitted changes remain

## Code Quality Standards

All Go code must follow:
- Maximum 35 lines per function (KTN-FUNC-001)
- Maximum 5 parameters (KTN-FUNC-002)
- Complete documentation with params/returns (KTN-FUNC-007)
- Comments on all control flow blocks (KTN-FUNC-011)
- No else after return (KTN-FUNC-012)
- Proper error handling
- Test coverage for all new features

## Project Structure

```
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
- `make docs` - Generate CHANGELOG.md and COVERAGE.MD
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

## Current Quality Metrics

- 22 test packages passing (0 failures)
- 70.9% overall coverage
- All linters passing with exceptions for test helpers
- Provider successfully builds with Bazel (240 actions)

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
2. Run `make docs` to regenerate CHANGELOG.md and COVERAGE.MD
3. Commit documentation updates
