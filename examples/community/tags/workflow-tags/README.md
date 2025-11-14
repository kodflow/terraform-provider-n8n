# Workflow Tags Example

Demonstrates how to create and use tags to organize workflows in n8n.

## What This Does

1. Creates three tags: "production", "automated", and "api"
2. Creates a workflow with all three tags assigned
3. Queries all tags in the system
4. Queries workflows filtered by the "production" tag

## Usage

```bash
# Set your credentials
export TF_VAR_n8n_api_key="your-api-key-here"

# Initialize and apply
terraform init
terraform apply
```

## Benefits of Tags

- **Organization**: Group related workflows together
- **Filtering**: Quickly find workflows by category
- **Documentation**: Self-documenting workflow purpose
- **Access Control**: Can be used for permission management (Enterprise)

## Expected Output

```
all_tags = [
  {
    id   = "1"
    name = "production"
  },
  {
    id   = "2"
    name = "automated"
  },
  {
    id   = "3"
    name = "api"
  }
]

production_workflows = [
  {
    id   = "..."
    name = "Tagged Workflow Example"
    tags = ["1", "2", "3"]
  }
]
```

## Cleanup

```bash
terraform destroy
```
