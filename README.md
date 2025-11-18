<!-- markdownlint-disable MD043 -->

# Terraform Provider for n8n

[![Bazel](https://img.shields.io/badge/Build-Bazel%209.0-43A047?logo=bazel)](https://bazel.build/)
[![Go](https://img.shields.io/badge/Go-1.25.4-00ADD8?logo=go)](https://pkg.go.dev/github.com/kodflow/terraform-provider-n8n)
[![n8n](https://img.shields.io/badge/n8n-1.119.2-EA4B71?logo=n8n)](https://n8n.io/)
[![Terraform Registry](https://img.shields.io/badge/dynamic/json?url=https://registry.terraform.io/v1/providers/kodflow/n8n&query=$.version&label=terraform&logo=terraform&color=7B42BC)](https://registry.terraform.io/providers/kodflow/n8n/latest)
[![CI](https://github.com/kodflow/terraform-provider-n8n/actions/workflows/ci.yml/badge.svg)](https://github.com/kodflow/terraform-provider-n8n/actions/workflows/ci.yml)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/6ad65f0b28b64849ad2799943e8ad338)](https://app.codacy.com/gh/kodflow/terraform-provider-n8n/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade)
[![Codacy Badge](https://app.codacy.com/project/badge/Coverage/6ad65f0b28b64849ad2799943e8ad338)](https://app.codacy.com/gh/kodflow/terraform-provider-n8n/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_coverage)

Manage your n8n workflows, credentials, and resources as code with Terraform.

## Why This Provider?

### 🎯 Standard Terraform Provider

Full support for n8n resources with standard Terraform workflows:

| Resource                  | Description                                              |
| ------------------------- | -------------------------------------------------------- |
| `n8n_workflow`            | Create and manage workflows                              |
| `n8n_workflow_node`       | Modular node composition                                 |
| `n8n_workflow_connection` | Connect nodes in workflows                               |
| `n8n_credential`          | Store API credentials securely                           |
| `n8n_tag`                 | Organize resources with tags                             |
| `n8n_variable`            | Manage environment variables                             |
| `n8n_project` 🚧          | Project management (Enterprise - not tested, no license) |
| `n8n_user` 🚧             | User management (Enterprise - not tested, no license)    |
| `n8n_source_control` 🚧   | Git integration (Enterprise - not tested, no license)    |

### 🚀 Advanced Node Composition

**296 n8n nodes** fully supported with comprehensive workflow examples:

| Category        | Count | Description                                                         |
| --------------- | ----- | ------------------------------------------------------------------- |
| **Core**        | 5     | Essential workflow building blocks (Code, If, Merge, Set, Switch)   |
| **Trigger**     | 25    | Event-based workflow initiators (Webhook, Manual, Cron, etc.)       |
| **Integration** | 266   | Third-party service integrations (Slack, GitHub, AWS, Google, etc.) |

**Every node includes:**

- ✅ Complete Terraform workflow example
- ✅ Full lifecycle testing (init/plan/apply/destroy)
- ✅ 100% test pass rate with real n8n validation
- ✅ Ready-to-use templates

**📊 [View complete test coverage →](COVERAGE.MD)** | **📚 [Browse all 296 nodes →](examples/nodes/)**

## Quick Start

### Installation

#### Via Terraform Registry (Recommended)

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

#### Via OpenTofu Registry

```hcl
terraform {
  required_providers {
    n8n = {
      source  = "kodflow/n8n"
      version = "~> 1.0"
    }
  }
}
```

### Get Your n8n API Key

1. Open your n8n instance
2. Go to **Settings** > **API**
3. Click **Create API Key**
4. Set as `N8N_API_KEY` environment variable

### Run Your First Workflow

```bash
export N8N_API_URL="http://localhost:5678"
export N8N_API_KEY="your-api-key"

cd examples/community/workflows/basic-workflow
terraform init
terraform apply
```

## Examples

### 📦 Node Examples (296 workflows)

Explore **all 296 n8n nodes** with complete Terraform examples:

- **[Core Nodes](examples/nodes/#core-nodes-5)** - If, Code, Set, Merge, Switch
- **[Trigger Nodes](examples/nodes/#trigger-nodes-25)** - Webhook, Cron, Manual, Email, etc.
- **[Integration Nodes](examples/nodes/#integration-nodes-266)** - Slack, GitHub, AWS, Google Cloud, etc.
- **[MEGA Workflow](examples/mega-workflow/)** - All 296 nodes in a single workflow (testing)

📊 **Testing status:** All 296 workflows tested and passing

### 🎓 Community Examples

Ready-to-use examples for common use cases:

- **[Workflows](examples/community/workflows/)** - Basic webhook and scheduled workflows
- **[Credentials](examples/community/credentials/)** - HTTP Basic Auth and API credentials
- **[Tags](examples/community/tags/)** - Workflow organization with tags
- **[Variables](examples/community/variables/)** - Environment variable management

### 🏗️ Production Examples

Advanced production-ready examples:

- **[Complete Modular Workflow](examples/comprehensive/complete-modular-workflow/)** - Multi-node workflow with error handling, validation, and external API
  integration

## Documentation

- **[Terraform Registry Docs](https://registry.terraform.io/providers/kodflow/n8n/latest/docs)** - Complete provider documentation
- **[Test Coverage](COVERAGE.MD)** - Detailed test coverage report (96.4%)
- **[All Nodes](examples/nodes/)** - Complete catalog of 296 supported nodes
- **[Contributing Guide](CONTRIBUTING.md)** - Development setup and guidelines

## Contributing

Contributions are welcome! See [CONTRIBUTING.md](CONTRIBUTING.md) for:

- Development environment setup
- Essential commands and workflow
- Quality standards and testing requirements
- Git hooks and commit conventions

Quick start for contributors:

```bash
make build    # Build provider locally
make test     # Run full test suite
make help     # Display all available commands
```

## Support This Project

If you find this project useful, consider sponsoring its development:

- ❤️ [GitHub Sponsors](https://github.com/sponsors/kodflow)
- ☕ [Ko-fi](https://ko-fi.com/kodflow)

Your support helps maintain and improve this provider. Thank you! 🙏

## License

Sustainable Use License 1.0 - See [LICENSE](LICENSE) for details.

---

**Developed with ❤️ by [KodFlow](https://github.com/kodflow)**
