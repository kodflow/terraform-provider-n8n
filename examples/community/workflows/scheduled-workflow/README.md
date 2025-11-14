# Scheduled Workflow Example

Creates a workflow that runs automatically every hour using a schedule trigger.

## What This Does

1. Creates a workflow with two nodes:
   - **Schedule Trigger node**: Executes every hour
   - **Code node**: Generates a timestamped message

2. Activates the workflow automatically

3. Workflow will execute every hour in the background

## Usage

```bash
# Set your credentials
export TF_VAR_n8n_api_key="your-api-key-here"

# Initialize and apply
terraform init
terraform apply

# Check executions in n8n UI after an hour
```

## Expected Behavior

The workflow will execute automatically every hour, creating an execution with:

```json
{
  "timestamp": "2025-11-13T10:00:00.000Z",
  "message": "Scheduled execution"
}
```

## Cleanup

```bash
terraform destroy
```
