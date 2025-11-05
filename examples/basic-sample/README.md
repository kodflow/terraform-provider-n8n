# Basic Sample - n8n Terraform Provider

This example demonstrates a minimal working setup using the n8n Terraform provider.

## What's Included

- **Tag**: Organizes resources with `tf:basic-sample`
- **Workflow**: Basic workflow that can be configured in n8n UI

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                     n8n Instance                         │
├─────────────────────────────────────────────────────────┤
│                                                           │
│  Tag: tf:basic-sample                                    │
│  │                                                        │
│  └── Workflow: TF Basic Sample - Mocky.io Test         │
│      └── Configure nodes manually in n8n UI             │
│                                                           │
└─────────────────────────────────────────────────────────┘
```

## Current Limitations

This example reflects the current state of the provider:

- ❌ **Credentials**: Data field not yet supported (create manually in UI)
- ❌ **Variables**: Requires n8n enterprise license
- ❌ **Workflow Nodes**: Not yet supported (configure in UI)
- ❌ **Workflow Connections**: Not yet supported (configure in UI)
- ❌ **Workflow Tags**: Not yet supported (assign in UI)

## Prerequisites

1. **n8n instance** running and accessible
2. **n8n API key** with appropriate permissions
3. **Terraform** installed (>= 1.0)
4. **Provider built** and installed locally

## Setup

### 1. Load Environment

Create or use existing `.env` file:

```bash
export N8N_URL="https://your-n8n-instance.com"
export N8N_API_TOKEN="your-api-token-here"
```

Load it:

```bash
source ../../.env
```

### 2. Initialize (Optional with Dev Override)

```bash
cd examples/basic-sample
terraform init  # May show warnings with dev override - that's normal
```

## Usage

### Deploy Resources

```bash
# Plan
terraform plan \
  -var="n8n_api_key=$N8N_API_TOKEN" \
  -var="n8n_base_url=$N8N_URL"

# Apply
terraform apply \
  -var="n8n_api_key=$N8N_API_TOKEN" \
  -var="n8n_base_url=$N8N_URL"
```

### View Created Resources

```
Outputs:

summary = {
  tag = {
    id   = "UFqb7E5zX6pqLyzx"
    name = "tf:basic-sample"
  }
  workflow = {
    active = false
    id     = "lD7NX9WciXeTFyBx"
    name   = "TF Basic Sample - Mocky.io Test"
    url    = "https://your-n8n.com/workflow/lD7NX9WciXeTFyBx"
  }
}
```

### Configure the Workflow

After Terraform creates the workflow:

1. Open the workflow URL (shown in outputs)
2. Add nodes manually in n8n UI:
   - **Manual Trigger** node
   - **HTTP Request** node pointing to any API
   - **Set** node to format response
3. Connect the nodes
4. (Optional) Manually tag with `tf:basic-sample`

### Clean Up

```bash
terraform destroy \
  -var="n8n_api_key=$N8N_API_TOKEN" \
  -var="n8n_base_url=$N8N_URL"
```

## What You'll Learn

This minimal example demonstrates:

1. ✅ **Basic Resource Creation**: Tag and Workflow
2. ✅ **Resource Dependencies**: Using `depends_on`
3. ✅ **Outputs**: Exposing resource information
4. ✅ **Provider Configuration**: API authentication

## Future Enhancements

As the provider matures, this example will be enhanced with:

- Credential data support
- Workflow node definitions
- Workflow connections
- Workflow tags
- Variables (with enterprise license)

## Files

- `main.tf` - Main Terraform configuration
- `variables.tf` - Input variable definitions
- `outputs.tf` - Output definitions
- `.gitignore` - Git ignore rules
- `terraform.tfvars.example` - Example variables file
- `README.md` - This file

## Resources Created

| Type | Name | Purpose |
|------|------|---------|
| `n8n_tag` | `tf:basic-sample` | Organizes resources |
| `n8n_workflow` | `TF Basic Sample - Mocky.io Test` | Empty workflow to configure |

## Troubleshooting

### Variables Error (403 Forbidden)

```
Error: feat:variables requires enterprise license
```

**Solution**: Variables are commented out by default. They require an n8n enterprise license.

### Provider Development Override Warning

```
Warning: Provider development overrides are in effect
```

**Info**: This is normal when testing with `make build`. The provider is installed locally for development.

## Next Steps

1. Configure the workflow nodes in n8n UI
2. Add credentials manually in n8n UI
3. Test the workflow execution
4. Try modifying the tag name and re-applying
5. Explore other examples in `/examples/`

## See Also

- [n8n API Documentation](https://docs.n8n.io/api/)
- [Terraform Documentation](https://www.terraform.io/docs)
- [Provider Source Code](/src/)
