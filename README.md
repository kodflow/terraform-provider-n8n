# Terraform Provider for n8n

Terraform provider to manage n8n resources (workflows, credentials, etc.).

[![Bazel](https://img.shields.io/badge/Build-Bazel%209.0-43A047?logo=bazel)](https://bazel.build/)
[![Go](https://img.shields.io/badge/Go-1.24-00ADD8?logo=go)](https://go.dev/)
[![Terraform](https://img.shields.io/badge/Terraform-Plugin%20Framework-7B42BC?logo=terraform)](https://developer.hashicorp.com/terraform/plugin/framework)

## Table of Contents

- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Development](#development)
- [Build and Tests](#build-and-tests)
- [Project Structure](#project-structure)
- [Release](#release)
- [Versioning and Releases](#versioning-and-releases)
- [Contributing](#contributing)
- [License](#license)

## Prerequisites

### Required Versions

- **Go 1.24.0+** (required by terraform-plugin-framework v1.16+)
- **Bazel 9.0+** (build system)
- **Terraform 1.0+** or **OpenTofu 1.0+**

### DevContainer (Recommended)

The project is configured with a DevContainer including all necessary tools:

- **Go 1.25.3** (compatible 1.24+)
- **Bazel 9.0.0rc1** (via Bazelisk)
- **Terraform & OpenTofu** (pre-installed)
- VS Code Extensions:
  - `golang.go` - Official Go support
  - `hashicorp.terraform` - Terraform support
  - `BazelBuild.vscode-bazel` - Bazel support

**To use the DevContainer:**
1. Open the project in VS Code
2. Accept the prompt to open in container
3. Wait for container build (first time only)

### Manual Installation

If you're not using the DevContainer:

```bash
# Install Go 1.24+
# See: https://go.dev/doc/install

# Install Bazelisk (recommended for managing Bazel versions)
go install github.com/bazelbuild/bazelisk@latest

# Verify versions
go version        # should display go1.24 or higher
bazel version     # should display Bazel 9.0+
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

### Local Installation for Development

```bash
# Build and install locally
make build

# The provider will be installed in:
# ~/.terraform.d/plugins/registry.terraform.io/kodflow/n8n/0.0.1/<OS>_<ARCH>/
```

## Development

### Available Make Commands

The `Makefile` provides essential commands:

```bash
make help    # Display help with all available commands
make test    # Run tests with Bazel
make build   # Build and install provider locally
make clean   # Clean Bazel artifacts
```

### Bazel Configuration

The project uses **Bazel 9** with **bzlmod** (the new dependency management system):

- **`.bazelversion`**: Bazel version (9.0.0rc1)
- **`MODULE.bazel`**: Dependencies and bzlmod configuration
- **`BUILD.bazel`**: Root build configuration
- **`.bazelrc`**: Bazel build options

**Bazel Dependencies:**
- `rules_go v0.58.3` - Go rules for Bazel (with Bazel 9 support)
- `gazelle v0.46.0` - Automatic BUILD files generator
- `rules_proto v7.1.0` - Protocol Buffers support
- `bazel_features v1.33.0` - Bazel features detection

### Project Architecture

```
.
├── .bazelrc              # Bazel configuration
├── .bazelversion         # Bazel version (9.0.0rc1)
├── MODULE.bazel          # bzlmod dependencies
├── BUILD.bazel           # Root build configuration
├── go.mod                # Go dependencies
├── Makefile              # Build commands
├── .devcontainer/        # DevContainer configuration
│   ├── Dockerfile        # Development image
│   └── devcontainer.json # VS Code configuration
├── src/                  # Provider source code
│   ├── main.go           # Entry point
│   ├── BUILD.bazel       # Source build configuration
│   └── internal/
│       └── provider/     # Provider implementation
│           ├── provider.go
│           ├── provider_test.go
│           └── BUILD.bazel
└── .github/
    └── workflows/        # CI/CD GitHub Actions
        ├── semver.yml    # Automatic semantic versioning
        └── release.yml   # Automatic release workflow
```

## Build and Tests

### Running Tests

```bash
# Via Make (recommended)
make test

# Directly with Bazel
bazel test //src/...

# Tests with verbose timeout warnings
bazel test --test_verbose_timeout_warnings //src/...
```

### Building the Provider

```bash
# Build with Bazel
bazel build //src:terraform-provider-n8n

# Binary will be available at:
# bazel-bin/src/terraform-provider-n8n
```

### Cleanup

```bash
# Clean Bazel artifacts
bazel clean

# Complete cleanup (including cache)
bazel clean --expunge
```

## Project Structure

### Source Code

The provider follows Terraform Plugin Framework best practices:

```
src/
├── main.go                    # Provider entry point
└── internal/
    └── provider/
        ├── provider.go        # Main provider implementation
        ├── provider_test.go   # Provider tests
        └── BUILD.bazel        # Build configuration
```

### Terraform Configuration

To use the provider in local development:

```hcl
terraform {
  required_providers {
    n8n = {
      source  = "registry.terraform.io/kodflow/n8n"
      version = "0.0.1"
    }
  }
}

provider "n8n" {
  # Provider configuration
}
```

## Release

### Release Workflow

The project uses **GoReleaser** via GitHub Actions to automate releases:

1. **Create a tag**:
   ```bash
   git tag -a v0.1.0 -m "Release v0.1.0"
   git push origin v0.1.0
   ```

2. **GitHub Actions** automatically triggers:
   - Cross-platform compilation (Linux, macOS, Windows, FreeBSD)
   - SHA256 checksums generation
   - GPG signature of checksums
   - GitHub release creation with artifacts

### GPG Configuration

To sign releases, configure a GPG key:

```bash
# Generate a GPG key
gpg --full-generate-key
# Choose: RSA and RSA, 4096 bits, no expiration

# Export private key
gpg --armor --export-secret-keys YOUR_EMAIL > private-key.asc

# Export public key
gpg --armor --export YOUR_EMAIL
```

Add GitHub secrets (Settings > Secrets and variables > Actions):
- `GPG_PRIVATE_KEY`: Content of `private-key.asc`
- `GPG_PASSPHRASE`: GPG key passphrase

### Release Artifacts

GoReleaser automatically generates:

```
terraform-provider-n8n_0.1.0_darwin_amd64.zip
terraform-provider-n8n_0.1.0_darwin_arm64.zip
terraform-provider-n8n_0.1.0_linux_amd64.zip
terraform-provider-n8n_0.1.0_linux_arm64.zip
terraform-provider-n8n_0.1.0_windows_amd64.zip
terraform-provider-n8n_0.1.0_SHA256SUMS
terraform-provider-n8n_0.1.0_SHA256SUMS.sig
```

### Registry Registration

#### Terraform Registry (Official)

1. Login to [registry.terraform.io](https://registry.terraform.io)
2. Go to "Publish" > "Provider"
3. Connect GitHub repository
4. Add GPG public key
5. Registry will automatically detect releases

#### OpenTofu Registry

OpenTofu uses the same format. Follow documentation at [github.com/opentofu/registry](https://github.com/opentofu/registry).

## Versioning and Releases

### Automatic Semantic Versioning

This project uses **semantic-release** to automate versioning according to [Semantic Versioning 2.0.0](https://semver.org/).

#### How It Works

Each merge into `main` automatically triggers:

1. **Commit analysis** since last release
2. **Version determination** based on commit types:
   - `fix:`, `perf:`, `refactor:`, `build:` → **Patch** (0.1.0 → 0.1.1)
   - `feat:` → **Minor** (0.1.0 → 0.2.0)
   - `BREAKING CHANGE:` or `!` → **Major** (0.1.0 → 1.0.0)
3. **Git tag creation** (e.g., `v0.2.0`)
4. **Automatic update** of [CHANGELOG.md](CHANGELOG.md)
5. **GitHub Release creation** with detailed notes
6. **Multi-platform binary compilation** and publication via GoReleaser

#### Commit Conventions

We use [Conventional Commits](https://www.conventionalcommits.org/):

| Type | Description | Version Impact |
|------|-------------|----------------|
| `feat:` | New feature | Minor (0.1.0 → 0.2.0) |
| `fix:` | Bug fix | Patch (0.1.0 → 0.1.1) |
| `perf:` | Performance improvement | Patch |
| `refactor:` | Refactoring | Patch |
| `build:` | Build changes | Patch |
| `docs:` | Documentation | No release |
| `test:` | Tests | No release |
| `chore:` | Maintenance | No release |
| `ci:` | CI/CD | No release |

**Breaking Change** (Major):
```bash
git commit -m "feat!: change workflows API"
# or
git commit -m "feat: change API

BREAKING CHANGE: API changed, see docs"
```

#### Release Workflow

```bash
# 1. Develop on a branch
git checkout -b feat/my-feature
git commit -m "feat: add my feature"
git push

# 2. Create PR with title: "feat: add my feature"

# 3. After review, merge into main

# 4. Automatically:
#    ✅ Version calculated (e.g., 0.2.0)
#    ✅ Tag created (v0.2.0)
#    ✅ CHANGELOG updated
#    ✅ GitHub Release published
#    ✅ Binaries compiled and signed (GPG)
```

See [CHANGELOG.md](CHANGELOG.md) for complete version history.

## Contributing

**Contributions welcome!**

### Quick Start

1. **Fork** the repository
2. **Clone** your fork
3. **Create a branch**: `git checkout -b feat/my-feature`
4. **Develop** and test: `make test && make build`
5. **Commit** with convention: `git commit -m "feat: my feature"`
6. **Push**: `git push origin feat/my-feature`
7. **Create a Pull Request** with conventional title

### Standards

- ✅ Go code with `gofmt` and `golint`
- ✅ Tests for any new feature
- ✅ [Conventional Commits](https://www.conventionalcommits.org/)
- ✅ PR title follows convention: `<type>: <description>`
- ✅ Updated documentation if necessary

### Valid PR Example

**Title**: `feat(workflows): add tags support`

**Description**:
- Implementation of tag management on workflows
- Added unit tests
- Updated documentation

**Impact**: Minor version bump (0.1.0 → 0.2.0)

See [PR template](.github/pull_request_template.md) for more details.

## Dependencies

### Main

- `github.com/hashicorp/terraform-plugin-framework` v1.16.1 - Terraform provider framework
- `github.com/hashicorp/terraform-plugin-docs` v0.24.0 - Documentation generation

### Build

- Bazel 9.0.0rc1 - Build system
- Go 1.24.0 - Programming language

See `go.mod` for complete dependencies list.

## CI/CD

### GitHub Actions

The project uses GitHub Actions for complete automation:

- **`.github/workflows/semver.yml`**: Automatic semantic versioning
  - Triggers on push to `main` (after PR merge)
  - Analyzes conventional commits
  - Determines new version (major, minor, patch)
  - Creates Git tag and GitHub Release
  - Automatically updates CHANGELOG.md

- **`.github/workflows/release.yml`**: Binary publication
  - Triggers on `v*` tags created by semantic-release
  - Compiles for all platforms (Linux, macOS, Windows, FreeBSD)
  - Generates SHA256 checksums
  - Signs with GPG
  - Publishes artifacts on GitHub Releases

### Bazel

Bazel build ensures:
- ✅ Reproducible builds
- ✅ Distributed cache
- ✅ Incremental compilation
- ✅ Parallel tests
- ✅ Multi-platform support

## Troubleshooting

### Bazel Won't Compile

```bash
# Clean cache
bazel clean --expunge

# Check version
bazel version  # Should display 9.0.0rc1 or higher

# Check .bazelversion
cat .bazelversion
```

### Tests Failing

```bash
# Run tests with more details
bazel test --test_output=all //src/...

# Check logs
bazel test --test_verbose_timeout_warnings //src/...
```

### DevContainer Won't Start

```bash
# Rebuild container
CMD/CTRL + Shift + P > "Dev Containers: Rebuild Container"

# Check logs
CMD/CTRL + Shift + P > "Dev Containers: Show Log"
```

## License

MPL-2.0

---

**Developed with ❤️ by [KodFlow](https://github.com/kodflow)**
