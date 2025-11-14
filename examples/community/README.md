# N8N Terraform Provider - Community Edition Examples

These examples work with n8n Community Edition (self-hosted) and don't require an Enterprise license.

## Available Resources

### ✅ Workflows

- Create and manage n8n workflows
- Examples: [workflows/](workflows/)

### ✅ Credentials

- Manage API credentials and authentication
- Examples: [credentials/](credentials/)

### ✅ Tags

- Organize workflows with tags
- Examples: [tags/](tags/)

### ✅ Variables

- Manage environment variables
- Examples: [variables/](variables/)

### ✅ Executions (Read-Only)

- Query workflow executions
- Examples: [executions/](executions/)

## Quick Start

1. **Set up provider configuration:**

```bash
export N8N_API_URL="http://localhost:5678"
export N8N_API_KEY="your-api-key"
```

2. **Initialize Terraform:**

```bash
cd workflows/basic-workflow
terraform init
terraform plan
terraform apply
```

## Prerequisites

- n8n instance (self-hosted or cloud)
- API key (Settings > API > Create API Key)
- Terraform or OpenTofu installed

## Testing

All examples can be tested with:

```bash
make test/n8n  # Requires running n8n instance
```
