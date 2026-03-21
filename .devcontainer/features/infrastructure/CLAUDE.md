<!-- updated: 2026-02-19T12:00:00Z -->

# Infrastructure Feature

## Purpose

Infrastructure development tools for Terraform/Terragrunt and Ansible workflows.

## Components

| Tool         | Source                   | Description                                             |
| ------------ | ------------------------ | ------------------------------------------------------- |
| terragrunt   | gruntwork-io/terragrunt  | Terraform wrapper for DRY configurations                |
| tflint       | terraform-linters/tflint | Terraform linter with cloud provider rules              |
| infracost    | infracost/infracost      | Cloud cost estimation from Terraform code               |
| cfssl        | cloudflare/cfssl         | PKI/TLS certificate management (optional)               |
| ansible-lint | pip                      | Ansible playbook linter (optional, needs Python)        |
| molecule     | pip                      | Ansible role testing framework (optional, needs Python) |

## Quick Start

```bash
# Terraform/Terragrunt workflow
terragrunt run-all plan          # Plan all modules
terragrunt run-all apply         # Apply all modules
tflint --recursive               # Lint all modules

# Cost estimation
infracost breakdown --path .     # Estimate costs

# PKI (if cfssl enabled)
cfssl gencert -initca ca-csr.json | cfssljson -bare ca

# Ansible (if Python feature enabled)
ansible-lint                     # Lint playbooks
molecule test                    # Test roles
```

## Configuration

Enable in `devcontainer.json`:

```json
"features": {
  "./features/infrastructure": {
    "terragruntVersion": "latest",
    "tflintVersion": "latest",
    "infracostVersion": "latest",
    "enableCfssl": true,
    "enableAnsibleTools": true
  }
}
```

## Options

| Option               | Default | Description                             |
| -------------------- | ------- | --------------------------------------- |
| `terragruntVersion`  | latest  | Terragrunt version (semver or latest)   |
| `tflintVersion`      | latest  | TFLint version (semver or latest)       |
| `infracostVersion`   | latest  | Infracost version (semver or latest)    |
| `enableCfssl`        | true    | Install cfssl + cfssljson               |
| `enableAnsibleTools` | true    | Install ansible-lint + molecule via pip |

## Recommended Companions

| Feature | Why                                         |
| ------- | ------------------------------------------- |
| Python  | Required for ansible-lint and molecule      |
| Go      | Required for Terratest infrastructure tests |

## Pre-installed (Base Image)

These tools are already in the base image (no feature needed):

- **Terraform** - IaC provisioning
- **Vault** - Secrets management (needs IPC_LOCK capability)
- **Consul** - Service discovery
- **Nomad** - Workload orchestration
- **Packer** - Image building
- **Ansible** - Configuration management
