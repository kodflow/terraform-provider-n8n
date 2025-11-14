# Executions Query Example

Demonstrates how to query workflow executions using Terraform data sources (read-only).

## What This Does

1. Creates a test workflow
2. Queries all executions in the system
3. Queries executions for a specific workflow
4. Queries only successful executions
5. Outputs execution statistics and recent executions

## Usage

```bash
# Set your credentials
export TF_VAR_n8n_api_key="your-api-key-here"

# Initialize and apply
terraform init
terraform apply

# Execute the workflow manually in n8n UI, then refresh
terraform refresh
terraform output
```

## Note

This example demonstrates **read-only** access to executions. The n8n Terraform provider does not create or manage executions directly - they are created by:

- Manual workflow execution in the UI
- Active workflows triggered by webhooks, schedules, etc.
- API calls to execute workflows

## Query Filters

The executions data source supports various filters:

```hcl
data "n8n_executions" "filtered" {
  workflow_id = "workflow-id"  # Filter by workflow
  status      = "success"      # Filter by status (success, error, waiting)
  limit       = 10             # Limit number of results
}
```

## Expected Output

After executing the workflow a few times:

```
total_executions = 5
workflow_executions_count = 3
successful_executions_count = 3

recent_executions = [
  {
    id          = "exec-123"
    workflow_id = "workflow-456"
    status      = "success"
    mode        = "manual"
  },
  ...
]
```

## Cleanup

```bash
terraform destroy
```
