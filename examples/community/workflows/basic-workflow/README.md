# Basic Workflow Example

Creates a simple webhook workflow that responds with a message.

## What This Does

1. Creates a workflow with two nodes:
   - **Webhook node**: Listens for POST requests
   - **Set node**: Returns a JSON response

2. Activates the workflow automatically

3. Outputs the webhook URL for testing

## Usage

```bash
# Set your credentials
export TF_VAR_n8n_api_key="your-api-key-here"

# Initialize and apply
terraform init
terraform apply

# Test the webhook
curl -X POST http://localhost:5678/webhook/example-webhook \
  -H "Content-Type: application/json" \
  -d '{"test": "data"}'
```

## Expected Output

```json
{
  "message": "Hello from Terraform!"
}
```

## Cleanup

```bash
terraform destroy
```
