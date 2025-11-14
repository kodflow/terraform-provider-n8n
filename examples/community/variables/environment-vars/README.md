# Environment Variables Example

Demonstrates how to manage n8n environment variables with Terraform.

## What This Does

1. Creates three environment variables:
   - `API_ENDPOINT`: API base URL
   - `API_TIMEOUT`: Request timeout in milliseconds
   - `ENVIRONMENT`: Current environment name

2. Queries all variables in the system
3. Queries a specific variable by ID

## Usage

```bash
# Set your credentials
export TF_VAR_n8n_api_key="your-api-key-here"
export TF_VAR_environment="production"

# Initialize and apply
terraform init
terraform apply
```

## Using Variables in Workflows

Environment variables can be referenced in workflows using the expression syntax:

```javascript
// In a Code node or expression
const apiEndpoint = $vars.API_ENDPOINT;
const timeout = parseInt($vars.API_TIMEOUT);
const env = $vars.ENVIRONMENT;

console.log(`Making request to ${apiEndpoint} in ${env} environment`);
```

## Benefits

- **Configuration Management**: Centralize environment-specific settings
- **Security**: Store configuration separately from workflows
- **Flexibility**: Change values without modifying workflows
- **Multi-Environment**: Support dev, staging, production configurations

## Expected Output

```
all_variables = [
  {
    key = "API_ENDPOINT"
    id  = "1"
  },
  {
    key = "API_TIMEOUT"
    id  = "2"
  },
  {
    key = "ENVIRONMENT"
    id  = "3"
  }
]
```

## Cleanup

```bash
terraform destroy
```
