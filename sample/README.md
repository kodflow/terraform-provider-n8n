# n8n Provider Usage Example

This directory contains an example of using the Terraform provider for n8n in local development.

## Prerequisites

1. Build and install the provider locally:
   ```bash
   cd ..
   make build
   ```

2. Verify installation:
   ```bash
   ls -la ~/.terraform.d/plugins/registry.terraform.io/kodflow/n8n/0.0.1/
   ```

## Usage

1. Initialize Terraform:
   ```bash
   terraform init
   ```

2. View execution plan:
   ```bash
   terraform plan
   ```

3. Apply configuration:
   ```bash
   terraform apply
   ```

## Configuration

The n8n provider requires the following configuration:

- `api_url`: URL of your n8n instance (e.g., `https://your-n8n.com`)
- `api_key`: n8n API key for authentication

### Environment Variables

You can also use environment variables:

```bash
export N8N_API_URL="https://your-n8n.com"
export N8N_API_KEY="your-api-key"
```

### Variables File

Create a `terraform.tfvars` file (ignored by Git):

```hcl
n8n_api_url = "https://your-n8n.com"
n8n_api_key = "your-api-key"
```

## Cleanup

To destroy created resources:

```bash
terraform destroy
```
