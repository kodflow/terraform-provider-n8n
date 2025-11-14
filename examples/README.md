# N8N Terraform Provider - Examples

This directory contains comprehensive examples for the n8n Terraform provider, organized by license type.

## Directory Structure

```
examples/
├── community/          # Community Edition examples (free, self-hosted)
│   ├── workflows/      # Workflow examples (webhook, scheduled, etc.)
│   ├── credentials/    # Credential management
│   ├── tags/           # Tag organization
│   ├── variables/      # Environment variables
│   └── executions/     # Execution queries (read-only)
└── enterprise/         # Enterprise Edition examples (requires license)
    ├── projects/       # Project management
    ├── users/          # User administration
    └── source-control/ # Git integration
```

## Community Edition vs Enterprise Edition

### Community Edition

All resources in [`community/`](community/) work with:

- n8n Community Edition (self-hosted, free)
- n8n Cloud (free tier and paid plans)

Available resources:

- ✅ Workflows
- ✅ Credentials
- ✅ Tags
- ✅ Variables
- ✅ Executions (read-only)

### Enterprise Edition

Resources in [`enterprise/`](enterprise/) are planned for future release:

- n8n Enterprise Edition required
- Valid enterprise license needed

Planned resources:

- 🚧 Projects (organize workflows and credentials)
- 🚧 Users (user management and access control)
- 🚧 Source Control (Git integration)

**Status**: Coming soon - awaiting enterprise license access for testing

## Quick Start

### Prerequisites

1. **Running n8n instance**:

   ```bash
   # Using Docker
   docker run -it --rm \
     --name n8n \
     -p 5678:5678 \
     n8nio/n8n

   # Or using npx
   npx n8n
   ```

2. **API Key**:
   - Open n8n: http://localhost:5678
   - Go to Settings > API
   - Click "Create API Key"
   - Copy the key

3. **Set environment variables**:
   ```bash
   export N8N_API_URL="http://localhost:5678"
   export N8N_API_KEY="your-api-key-here"
   ```

### Try an Example

```bash
# Choose an example
cd community/workflows/basic-workflow

# Initialize Terraform
terraform init

# Review the plan
terraform plan

# Apply the configuration
terraform apply

# Test the webhook (for webhook examples)
curl -X POST http://localhost:5678/webhook/example-webhook \
  -H "Content-Type: application/json" \
  -d '{"test": "data"}'

# Clean up
terraform destroy
```

## Testing All Examples

Use the provided test script to validate all examples:

```bash
# Test all examples
./scripts/test-examples.sh all

# Test only community examples
./scripts/test-examples.sh community

# Test only enterprise examples
./scripts/test-examples.sh enterprise
```

The script will:

- ✅ Check prerequisites (Terraform, n8n connectivity)
- ✅ Run `terraform init/plan/apply/destroy` for each example
- ✅ Report results with pass/fail status
- ✅ Clean up resources automatically

## Example Categories

### 1. Workflows

Create and manage n8n workflows with various triggers:

- **Webhook**: HTTP endpoints for external integrations
- **Schedule**: Time-based automated workflows
- **Manual**: Workflows triggered manually

### 2. Credentials

Manage authentication credentials for API integrations:

- **HTTP Basic Auth**: Username/password authentication
- **API Key**: Token-based authentication
- **OAuth2**: OAuth 2.0 authentication

### 3. Tags

Organize and categorize workflows:

- Create tags for workflow categorization
- Assign multiple tags to workflows
- Query workflows by tags

### 4. Variables

Manage environment variables accessible in workflows:

- Store configuration values
- Use in workflow expressions
- Environment-specific settings

### 5. Executions (Read-Only)

Query workflow execution history:

- View execution status
- Filter by workflow, status, date
- Monitor workflow performance

### 6. Projects (Enterprise)

Organize workflows and credentials into projects:

- Team-based organization
- Project-level permissions
- Workflow isolation

### 7. Users (Enterprise)

Manage user accounts and access:

- Create and manage users
- Assign roles and permissions
- User lifecycle management

### 8. Source Control (Enterprise)

Git integration for workflow versioning:

- Pull workflows from Git
- Synchronize multiple instances
- Version control for workflows

## Best Practices

1. **Use Variables**: Store sensitive data in environment variables, not in Terraform files
2. **Tag Resources**: Use tags to organize workflows by team, environment, or purpose
3. **State Management**: Use remote state (S3, Terraform Cloud) for team collaboration
4. **Modularize**: Create reusable modules for common workflow patterns
5. **Version Control**: Commit Terraform files to Git, exclude `.terraform/` and `*.tfstate`

## Common Patterns

### Multi-Environment Setup

```hcl
# environments/dev/terraform.tfvars
n8n_api_url = "http://dev.n8n.example.com"
environment = "development"

# environments/prod/terraform.tfvars
n8n_api_url = "https://prod.n8n.example.com"
environment = "production"
```

### Workflow with Dependencies

```hcl
# Create credentials first
resource "n8n_credential" "api_key" {
  name = "External API Key"
  type = "httpCustomAuth"
  data = jsonencode({
    headerAuth = {
      name  = "X-API-Key"
      value = var.api_key
    }
  })
}

# Use credentials in workflow
resource "n8n_workflow" "api_workflow" {
  name = "API Integration"
  nodes = jsonencode([{
    type        = "n8n-nodes-base.httpRequest"
    credentials = {
      httpCustomAuth = {
        id   = n8n_credential.api_key.id
        name = n8n_credential.api_key.name
      }
    }
  }])
}
```

### Tagged Workflows

```hcl
# Create tags
resource "n8n_tag" "production" {
  name = "production"
}

resource "n8n_tag" "critical" {
  name = "critical"
}

# Apply to workflow
resource "n8n_workflow" "important_workflow" {
  name = "Critical Production Workflow"
  tags = [
    n8n_tag.production.id,
    n8n_tag.critical.id
  ]
}
```

## Troubleshooting

### Connection Issues

```bash
# Test n8n API connectivity
curl -H "X-N8N-API-KEY: $N8N_API_KEY" \
  $N8N_API_URL/api/v1/workflows

# Should return: {"data": [...], "nextCursor": null}
```

### Authentication Errors

```bash
# Verify API key is set
echo $N8N_API_KEY

# Check API key in n8n
# Settings > API > Your API keys
```

### Provider Not Found

```bash
# Install provider locally
cd /workspace
make build

# Verify installation
ls ~/.terraform.d/plugins/registry.terraform.io/kodflow/n8n/
```

## Contributing Examples

Have a useful pattern or example? Contributions welcome!

1. Create example in appropriate category
2. Include `main.tf`, `variables.tf`, and `README.md`
3. Test with `./scripts/test-examples.sh`
4. Submit pull request

## Resources

- **Main README**: [/workspace/README.md](../README.md)
- **Provider Documentation**: https://registry.terraform.io/providers/kodflow/n8n
- **n8n Documentation**: https://docs.n8n.io/
- **Terraform Documentation**: https://developer.hashicorp.com/terraform
