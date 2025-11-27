<!-- markdownlint-disable MD043 -->

# Terraform Provider for n8n

[![Bazel](https://img.shields.io/badge/Build-Bazel%209.0-43A047?logo=bazel)](https://bazel.build/)
[![Go](https://img.shields.io/badge/Go-1.25.4-00ADD8?logo=go)](https://pkg.go.dev/github.com/kodflow/terraform-provider-n8n)
[![n8n](https://img.shields.io/badge/n8n-1.121.2-EA4B71?logo=n8n)](https://n8n.io/)
[![Terraform Registry](https://img.shields.io/badge/dynamic/json?url=https://registry.terraform.io/v1/providers/kodflow/n8n&query=$.version&label=terraform&logo=terraform&color=7B42BC)](https://registry.terraform.io/providers/kodflow/n8n/latest)
[![CI](https://github.com/kodflow/terraform-provider-n8n/actions/workflows/ci.yml/badge.svg)](https://github.com/kodflow/terraform-provider-n8n/actions/workflows/ci.yml)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/6ad65f0b28b64849ad2799943e8ad338)](https://app.codacy.com/gh/kodflow/terraform-provider-n8n/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade)
[![Codacy Badge](https://app.codacy.com/project/badge/Coverage/6ad65f0b28b64849ad2799943e8ad338)](https://app.codacy.com/gh/kodflow/terraform-provider-n8n/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_coverage)

**Manage your n8n workflows as code.** Version control, PR reviews, GitOps workflows, and reproducible environments for your automation platform.

---

## Why This Provider?

Low-code is great for building fast. It's terrible for:

- **Versioning workflows cleanly** — giant JSON exports are unreadable in diffs
- **Reviewing changes in PRs** — no way to track what actually changed
- **Reproducing environments** — manual clicks don't scale across teams
- **Collaborating without clicking through a UI** — engineers need code

This provider doesn't replace the visual builder — it gives you a **proper, maintainable backbone** behind it. Define workflows as structured code while keeping
the visual magic for rapid prototyping.

---

## ⚠️ Work in Progress

> **This provider is currently under active development.**
>
> I use it for my personal needs and continuously improve it. **Every issue you create helps me fix bugs and improve the quality of the product** — please don't
> hesitate to report problems!

### 🚨 Important Limitation: Credentials

> **The official n8n API does not support GET/PUT/PATCH operations on credentials.**
>
> As a result, the `n8n_credential` resource has **degraded behavior**:
>
> - ✅ Create (`POST`): works normally
> - ❌ Read (`GET`): not possible — state cannot be retrieved from n8n
> - ❌ Update (`PUT/PATCH`): not possible — changes require resource recreation
> - ✅ Delete (`DELETE`): works normally
>
> **Impact**: Terraform cannot detect drift on existing credentials. Any modification forces a resource recreation.

---

## Resources

| Resource                  | Description                                    |
| ------------------------- | ---------------------------------------------- |
| `n8n_workflow`            | Create and manage workflows                    |
| `n8n_workflow_node`       | Modular node composition                       |
| `n8n_workflow_connection` | Connect nodes in workflows                     |
| `n8n_credential` ⚠️       | Store API credentials securely (limited API)   |
| `n8n_tag`                 | Organize resources with tags                   |
| `n8n_variable`            | Manage environment variables (Enterprise)      |
| `n8n_project`             | Project management (Enterprise)                |
| `n8n_user`                | User management (Enterprise)                   |
| `n8n_source_control` 🚧   | Git integration (Enterprise - not implemented) |

---

## Quick Start

### 1. Configure the Provider

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

> Works with both **Terraform** and **OpenTofu** — same configuration, same registry.

### 2. Get Your n8n API Key

1. Open your n8n instance
2. Go to **Settings** → **API**
3. Click **Create API Key**
4. Export it:

```bash
export TF_VAR_n8n_api_key="your-api-key"
export TF_VAR_n8n_base_url="http://localhost:5678"
```

### 3. Run Your First Workflow

```bash
cd examples/community/workflows/basic-workflow
terraform init
terraform apply
```

---

## Examples

### 🎓 Community Examples

Basic examples for n8n Community Edition:

| Example                                                                | Description                           |
| ---------------------------------------------------------------------- | ------------------------------------- |
| [Basic Workflow](examples/community/workflows/basic-workflow/)         | Simple webhook workflow with response |
| [Scheduled Workflow](examples/community/workflows/scheduled-workflow/) | Cron-triggered automation             |
| [Credentials](examples/community/credentials/)                         | HTTP Basic Auth setup                 |
| [Tags](examples/community/tags/)                                       | Workflow organization                 |

### 🏢 Enterprise Examples

Examples requiring n8n Enterprise license:

| Example                                     | Description                |
| ------------------------------------------- | -------------------------- |
| [Projects](examples/enterprise/projects/)   | Multi-project organization |
| [Users](examples/enterprise/users/)         | User management and roles  |
| [Variables](examples/enterprise/variables/) | Environment variables      |

### 📦 Node Examples (296 workflows)

Complete Terraform examples for **all 296 n8n nodes**:

| Category                                       | Count | Examples                               |
| ---------------------------------------------- | ----- | -------------------------------------- |
| **[Core](examples/nodes/core/)**               | 5     | Code, If, Merge, Set, Switch           |
| **[Trigger](examples/nodes/trigger/)**         | 25    | Webhook, Cron, Manual, Email, etc.     |
| **[Integration](examples/nodes/integration/)** | 266   | Slack, GitHub, AWS, Google Cloud, etc. |

Every node includes a complete workflow example with full lifecycle testing.

**📚 [Browse all 296 nodes →](examples/nodes/)**

---

## Documentation

- **[Terraform Registry Docs](https://registry.terraform.io/providers/kodflow/n8n/latest/docs)** — Complete provider documentation
- **[Test Coverage](COVERAGE.MD)** — Detailed coverage report (97.9%)
- **[All Nodes](examples/nodes/)** — Complete catalog of 296 supported nodes
- **[Contributing Guide](CONTRIBUTING.md)** — Development setup and guidelines

---

## Support This Project

If you find this project useful, consider sponsoring its development:

- ❤️ [GitHub Sponsors](https://github.com/sponsors/kodflow)
- ☕ [Ko-fi](https://ko-fi.com/kodflow)

---

## License

Sustainable Use License 1.0 — See [LICENSE](LICENSE) for details.

---

**Developed with ❤️ by [KodFlow](https://github.com/kodflow)**
