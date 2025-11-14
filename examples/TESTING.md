# Testing N8N Terraform Provider Examples

This guide explains how to test the examples against a real n8n server.

## Prerequisites

### 1. Running n8n Instance

You need a running n8n instance. Choose one option:

#### Option A: Docker (Recommended for Testing)

```bash
# Start n8n with Docker
docker run -d \
  --name n8n \
  -p 5678:5678 \
  -e N8N_BASIC_AUTH_ACTIVE=true \
  -e N8N_BASIC_AUTH_USER=admin \
  -e N8N_BASIC_AUTH_PASSWORD=password \
  n8nio/n8n

# Wait for n8n to be ready (check logs)
docker logs -f n8n
```

#### Option B: NPX

```bash
npx n8n
```

#### Option C: Existing Instance

Use your existing n8n instance (cloud or self-hosted).

### 2. Create API Key

1. Open n8n: http://localhost:5678
2. Log in (if required)
3. Go to **Settings** > **API**
4. Click **Create API Key**
5. Copy the generated key

### 3. Set Environment Variables

```bash
# Set n8n connection details
export N8N_API_URL="http://localhost:5678"
export N8N_API_KEY="your-api-key-here"

# Verify connection
curl -H "X-N8N-API-KEY: $N8N_API_KEY" \
  $N8N_API_URL/api/v1/workflows
```

Expected output: `{"data": [], "nextCursor": null}`

## Testing Methods

### Method 1: Automated Testing (Recommended)

Use the provided test script to validate all examples:

```bash
# Test all examples (community + enterprise)
./scripts/test-examples.sh all

# Test only community examples
./scripts/test-examples.sh community

# Test only enterprise examples (requires enterprise license)
./scripts/test-examples.sh enterprise
```

The script will:

- ✅ Verify prerequisites (Terraform, API key, connectivity)
- ✅ Run `terraform init` for each example
- ✅ Run `terraform plan` to validate configuration
- ✅ Run `terraform apply` to create resources
- ✅ Display outputs
- ✅ Run `terraform destroy` to clean up
- ✅ Report pass/fail status for each example

### Method 2: Manual Testing

Test individual examples manually:

```bash
# Navigate to example directory
cd examples/community/workflows/basic-workflow

# Export credentials
export TF_VAR_n8n_api_url="http://localhost:5678"
export TF_VAR_n8n_api_key="your-api-key-here"

# Initialize Terraform
terraform init

# Review the plan
terraform plan

# Apply the configuration
terraform apply

# Test the created resource (for webhook example)
curl -X POST http://localhost:5678/webhook/example-webhook \
  -H "Content-Type: application/json" \
  -d '{"test": "data"}'

# Expected output:
# {"message": "Hello from Terraform!"}

# Clean up
terraform destroy
```

## Community Edition Examples

### Workflows

#### Basic Webhook

```bash
cd examples/community/workflows/basic-workflow
terraform init && terraform plan && terraform apply
curl -X POST http://localhost:5678/webhook/example-webhook \
  -H "Content-Type: application/json" \
  -d '{"test": "data"}'
terraform destroy
```

#### Scheduled Workflow

```bash
cd examples/community/workflows/scheduled-workflow
terraform init && terraform plan && terraform apply
# Wait for hourly execution or check n8n UI
terraform destroy
```

### Credentials

#### HTTP Basic Auth

```bash
cd examples/community/credentials/basic-auth
terraform init && terraform plan && terraform apply
terraform output credential_id
terraform destroy
```

### Tags

#### Workflow Tags

```bash
cd examples/community/tags/workflow-tags
terraform init && terraform plan && terraform apply
terraform output all_tags
terraform destroy
```

### Variables

#### Environment Variables

```bash
cd examples/community/variables/environment-vars
terraform init && terraform plan && terraform apply
terraform output all_variables
terraform destroy
```

### Executions

#### Query Executions

```bash
cd examples/community/executions/query-executions
terraform init && terraform plan && terraform apply
# Execute the workflow manually in n8n UI
terraform refresh
terraform output recent_executions
terraform destroy
```

## Enterprise Edition Examples

**Note**: Enterprise examples are currently not available. They will be added once enterprise license access is obtained for testing.

Planned enterprise examples include:

- **Projects**: Project creation and user assignment
- **Users**: User account management
- **Source Control**: Git integration

To test community examples only:

```bash
./scripts/test-examples.sh community
```

## Troubleshooting

### Connection Refused

```bash
# Check if n8n is running
curl http://localhost:5678

# Check Docker container
docker ps | grep n8n
docker logs n8n
```

### Authentication Failed

```bash
# Verify API key
echo $N8N_API_KEY

# Test API key manually
curl -H "X-N8N-API-KEY: $N8N_API_KEY" \
  http://localhost:5678/api/v1/workflows
```

### Provider Not Found

```bash
# Build and install provider locally
cd /workspace
make build

# Verify installation
ls ~/.terraform.d/plugins/registry.terraform.io/kodflow/n8n/
```

### Example Apply Failed

```bash
# Check Terraform state
terraform show

# View detailed error
terraform apply -no-color

# Force destroy if stuck
terraform destroy -auto-approve

# Clean up
rm -rf .terraform .terraform.lock.hcl terraform.tfstate*
```

### Enterprise Features Not Available

```bash
# Verify enterprise license
curl -H "X-N8N-API-KEY: $N8N_API_KEY" \
  $N8N_API_URL/api/v1/license

# If no license:
# - Use only community examples
# - Or obtain enterprise license
```

## Expected Results

When all tests pass, you should see output similar to:

```
═══════════════════════════════════════════════════
Testing Community Edition Examples
═══════════════════════════════════════════════════

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Testing: workflows/basic-workflow
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

→ terraform init
✓ Init successful

→ terraform plan
✓ Plan successful

→ terraform apply
✓ Apply successful

→ terraform output
workflow_id = "1"
workflow_webhook_url = "http://localhost:5678/webhook/example-webhook"
workflow_active = true

→ terraform destroy
✓ Destroy successful

✅ Example test passed: workflows/basic-workflow

[... more examples ...]

═══════════════════════════════════════════════════
Community Edition Results
✓ Passed: 6
═══════════════════════════════════════════════════

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🎉 All tests passed!
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

## CI/CD Integration

### GitHub Actions Example

```yaml
name: Test Examples

on:
  pull_request:
    paths:
      - "examples/**"

jobs:
  test-examples:
    runs-on: ubuntu-latest
    services:
      n8n:
        image: n8nio/n8n
        ports:
          - 5678:5678
        env:
          N8N_BASIC_AUTH_ACTIVE: true
          N8N_BASIC_AUTH_USER: admin
          N8N_BASIC_AUTH_PASSWORD: password

    steps:
      - uses: actions/checkout@v4

      - name: Setup Terraform
        uses: hashicorp/setup-terraform@v3

      - name: Build Provider
        run: make build

      - name: Wait for n8n
        run: |
          for i in {1..30}; do
            if curl -f http://localhost:5678; then break; fi
            sleep 2
          done

      - name: Create API Key
        id: api-key
        run: |
          # Create API key via n8n API
          # Store in $GITHUB_OUTPUT

      - name: Test Examples
        env:
          N8N_API_URL: http://localhost:5678
          N8N_API_KEY: ${{ steps.api-key.outputs.key }}
        run: ./scripts/test-examples.sh community
```

## Performance Notes

- **Basic workflow**: ~5 seconds (init + plan + apply + destroy)
- **Full community suite**: ~30-60 seconds (6 examples)
- **Full test suite**: ~1-2 minutes (all examples)

## Cleanup

After testing, clean up your n8n instance:

```bash
# Remove all workflows
curl -X DELETE -H "X-N8N-API-KEY: $N8N_API_KEY" \
  $N8N_API_URL/api/v1/workflows

# Stop Docker container
docker stop n8n
docker rm n8n
```

## Next Steps

- ✅ Automated testing in CI/CD
- ✅ Add more complex examples
- ✅ Create reusable modules
- ✅ Document common patterns
- ✅ Share your own examples!

## Support

Having issues? Check:

1. [Main README](../README.md) - Provider documentation
2. [Examples README](README.md) - Example overview
3. [n8n Documentation](https://docs.n8n.io/) - n8n API reference
4. [GitHub Issues](https://github.com/kodflow/n8n/issues) - Report bugs
