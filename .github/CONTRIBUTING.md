<!-- markdownlint-disable MD043 -->
<!-- Copyright (c) 2024 Florent (Kodflow). All rights reserved. -->
<!-- Licensed under the Sustainable Use License 1.0 -->
<!-- See LICENSE.md in the project root for license information. -->

# Contributing to n8n Terraform Provider

Thank you for your interest in contributing to the n8n Terraform Provider! This document provides guidelines and instructions for contributing to the project.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Development Workflow](#development-workflow)
- [Code Standards](#code-standards)
- [Testing Requirements](#testing-requirements)
- [Submitting Changes](#submitting-changes)
- [Review Process](#review-process)

## Code of Conduct

This project follows a professional code of conduct. We expect all contributors to be respectful, constructive, and collaborative.

## Getting Started

Before contributing, please:

1. **Search existing issues** to see if your bug/feature has already been reported
2. **Read the documentation** in the README and CLAUDE.md
3. **Understand the project structure** by reviewing the codebase
4. **Check the n8n API documentation** to understand the underlying API

## Development Setup

### Prerequisites

- **Go** 1.21+ (check `go.mod` for exact version)
- **Terraform** 1.0+
- **Bazel** (for build system)
- **Make** (for development commands)
- **Git** with GPG signing configured
- **n8n instance** for testing (API URL and API key required)

### Initial Setup

```bash
# Clone the repository
git clone https://github.com/kodflow/terraform-provider-n8n.git
cd terraform-provider-n8n

# Install development tools
make tools

# Build the provider
make build

# Run tests
make test
```

### Environment Configuration

Create a `.env` file in the project root with your n8n credentials:

```bash
N8N_API_URL=https://your-n8n-instance.com
N8N_API_KEY=your-api-key-here
```

**Important:** Never commit the `.env` file! It's already in `.gitignore`.

## Development Workflow

### 1. Create a Branch

Always create a new branch for your changes:

```bash
git checkout -b feat/your-feature-name
# or
git checkout -b fix/bug-description
```

Branch naming convention:

- `feat/` - New features
- `fix/` - Bug fixes
- `docs/` - Documentation updates
- `refactor/` - Code refactoring
- `test/` - Test improvements
- `chore/` - Maintenance tasks

### 2. Make Your Changes

Follow the development standards outlined in [CLAUDE.md](../CLAUDE.md):

- Write clean, readable code
- Add comprehensive tests
- Follow Go best practices
- Comment complex logic

### 3. Run the Quality Pipeline

**CRITICAL:** Before committing, you MUST run the following commands and fix ALL issues:

```bash
# Format code (required)
make fmt

# Run linters (zero errors, zero warnings)
make lint

# Run tests (all tests must pass)
make test

# Build the provider (must succeed)
make build
```

**Important:** The linters (`golangci-lint` and `ktn-linter`) must show ZERO errors and ZERO warnings. Do not skip this step!

### 4. Commit Your Changes

Follow [Conventional Commits](https://www.conventionalcommits.org/) format:

```text
<type>(<scope>): <description>

[optional body]

[optional footer]
```

**Types:**

- `feat`: New feature (minor version bump)
- `fix`: Bug fix (patch version bump)
- `docs`: Documentation only
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks
- `ci`: CI/CD changes
- `perf`: Performance improvements

**Examples:**

```bash
git commit -m "feat(workflow): add support for workflow tags"
git commit -m "fix(credential): handle empty credential data"
git commit -m "docs(readme): update installation instructions"
```

**Pre-commit hooks** will automatically:

- Add copyright headers
- Generate coverage documentation
- Validate commit message format
- Check for GPG signature
- Prevent AI mentions in commit messages

### 5. Push and Create a Pull Request

```bash
git push -u origin your-branch-name
```

Then create a Pull Request on GitHub with:

- Clear title following Conventional Commits
- Detailed description using the PR template
- All checklist items completed
- Links to related issues

## Code Standards

### Go Code Style

- Follow official [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` and `goimports` (included in `make fmt`)
- Keep functions small and focused
- Use descriptive variable names
- Add godoc comments for exported functions

### Project-Specific Standards

See [CLAUDE.md](../CLAUDE.md) for detailed standards including:

- **No `t.Skip()`** - Write proper unit tests with mocks
- **Maximum test coverage** - Test all code paths
- **Pure unit tests** - Mock all external dependencies
- **Fix ALL linter errors** - Zero tolerance for warnings
- **Complete the workflow** - fmt → test → lint → build

### File Organization

```
src/internal/provider/
├── <resource-name>/
│   ├── resource.go          # Resource implementation
│   ├── resource_test.go     # Unit tests
│   ├── data_source.go       # Data source (if applicable)
│   └── BUILD.bazel          # Bazel build file
```

## Testing Requirements

### Unit Tests

**ALL code must have unit tests.** Requirements:

- **Never use `t.Skip()`** - If it's hard to test, mock everything needed
- **Test all code paths** - Every line, every branch, every error case
- **Mock external dependencies** - HTTP clients, APIs, file systems, etc.
- **Aim for 100% coverage** - Keep adding tests until coverage is maximized

Example test structure:

```go
func TestResourceCreate(t *testing.T) {
    // Mock the API client
    mockClient := &MockN8NClient{
        CreateFunc: func(ctx context.Context, req *Request) (*Response, error) {
            return &Response{ID: "123"}, nil
        },
    }

    // Test the resource
    // ... test implementation
}
```

### Acceptance Tests

For integration testing with a real n8n instance:

```bash
# Set environment variables
export N8N_API_URL=https://your-n8n-instance.com
export N8N_API_KEY=your-api-key
export TF_ACC=1

# Run acceptance tests
make test/n8n
```

### Running Tests

```bash
# Run all unit tests
make test

# Run tests for specific package
go test ./src/internal/provider/workflow/...

# Run with coverage
go test -cover ./...

# Run acceptance tests (requires n8n instance)
make test/n8n
```

## Submitting Changes

### Pull Request Guidelines

1. **Use the PR template** - Fill out all sections
2. **Link related issues** - Use "Fixes #123" or "Closes #456"
3. **Keep PRs focused** - One feature/fix per PR
4. **Update documentation** - If you change behavior
5. **Add/update tests** - For all code changes
6. **Ensure CI passes** - All checks must be green

### PR Checklist

Before submitting, verify:

- [ ] Code follows project conventions (CLAUDE.md)
- [ ] All tests pass: `make test`
- [ ] Code is formatted: `make fmt`
- [ ] All linters pass: `make lint` (zero errors, zero warnings)
- [ ] Build succeeds: `make build`
- [ ] Documentation is updated
- [ ] Commit messages follow Conventional Commits
- [ ] PR title follows Conventional Commits
- [ ] Test coverage is maintained or improved

## Review Process

### What to Expect

1. **Automated checks** - CI/CD pipeline runs automatically
2. **Code review** - Maintainers will review your code
3. **Feedback** - You may be asked to make changes
4. **Approval** - Once approved, your PR will be merged

### Review Timeline

- Initial response: Within 3-5 business days
- Full review: Depends on PR complexity
- Merge: After approval and passing all checks

### After Merge

- Your changes will be included in the next release
- You'll be credited in the release notes
- Thank you for your contribution! 🎉

## Getting Help

If you need help:

1. **Check the documentation** - README.md and CLAUDE.md
2. **Search existing issues** - Someone may have asked before
3. **Ask in discussions** - Use GitHub Discussions for questions
4. **Create an issue** - For bugs or feature requests

## Additional Resources

- [CLAUDE.md](../CLAUDE.md) - Detailed development standards
- [README.md](../README.md) - Project overview and usage
- [n8n API Documentation](https://docs.n8n.io/api/) - API reference
- [Terraform Plugin Development](https://developer.hashicorp.com/terraform/plugin) - Terraform guides

---

Thank you for contributing to the n8n Terraform Provider! Your contributions help make infrastructure-as-code better for everyone. 🚀
