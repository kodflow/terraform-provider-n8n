# Modular Workflow Example

This example demonstrates the new **modular workflow composition** feature in the n8n Terraform provider.

## Overview

Instead of defining workflows as large JSON blobs, you can now define individual nodes and connections as separate Terraform resources. This makes your workflow
configurations:

- **More readable**: Each node is clearly defined with its own resource block
- **More maintainable**: Easy to modify individual nodes without touching a giant JSON blob
- **More reusable**: Nodes and connections can be referenced and composed
- **Type-safe**: Terraform validates your node and connection configurations

## Resources

### `n8n_workflow_node`

Defines an individual workflow node. This resource is **local-only** and doesn't make API calls. It exists in Terraform state to generate node JSON.

**Attributes:**

- `name` (required): Display name of the node (used in connections)
- `type` (required): n8n node type (e.g., `n8n-nodes-base.webhook`)
- `position` (required): `[x, y]` coordinates for UI display
- `parameters` (optional): Node parameters as JSON string
- `type_version` (optional): Node type version (default: 1)
- `webhook_id` (optional): Webhook identifier for webhook nodes
- `disabled` (optional): Whether the node is disabled
- `notes` (optional): User notes about the node
- `node_json` (computed): Generated JSON representation

### `n8n_workflow_connection`

Defines a connection between two workflow nodes. This resource is also **local-only**.

**Attributes:**

- `source_node` (required): Name of the source node
- `target_node` (required): Name of the destination node
- `source_output` (optional): Output type (default: `"main"`)
- `source_output_index` (optional): Output index for multi-output nodes like Switch (default: 0)
- `target_input` (optional): Input type (default: `"main"`)
- `target_input_index` (optional): Input index (default: 0)
- `connection_json` (computed): Generated JSON representation

## Examples

### Basic Modular Workflow

See `main.tf` for a simple webhook → process → respond workflow using modular nodes.

```bash
terraform init
terraform plan
terraform apply
```

### Advanced: Switch with Multiple Branches

See `advanced-switch.tf` for a complex workflow with a Switch node routing to multiple branches.

This demonstrates how to handle nodes with multiple outputs:

```hcl
resource "n8n_workflow_connection" "switch_to_branch_1" {
  source_node         = n8n_workflow_node.switch.name
  source_output_index = 0  # First branch
  target_node         = n8n_workflow_node.branch_1.name
}

resource "n8n_workflow_connection" "switch_to_branch_2" {
  source_node         = n8n_workflow_node.switch.name
  source_output_index = 1  # Second branch
  target_node         = n8n_workflow_node.branch_2.name
}
```

## Migration from JSON Blobs

### Before (JSON Blob approach):

```hcl
resource "n8n_workflow" "example" {
  name = "My Workflow"

  nodes_json = jsonencode([
    {
      id = "webhook-1"
      name = "Webhook"
      type = "n8n-nodes-base.webhook"
      position = [250, 300]
      parameters = { ... }
    },
    {
      id = "code-1"
      name = "Process"
      type = "n8n-nodes-base.code"
      position = [450, 300]
      parameters = { ... }
    }
  ])

  connections_json = jsonencode({
    "Webhook" = { main = [[{ node = "Process", type = "main", index = 0 }]] }
  })
}
```

### After (Modular approach):

```hcl
resource "n8n_workflow_node" "webhook" {
  name     = "Webhook"
  type     = "n8n-nodes-base.webhook"
  position = [250, 300]
  parameters = jsonencode({ ... })
}

resource "n8n_workflow_node" "process" {
  name     = "Process"
  type     = "n8n-nodes-base.code"
  position = [450, 300]
  parameters = jsonencode({ ... })
}

resource "n8n_workflow_connection" "webhook_to_process" {
  source_node = n8n_workflow_node.webhook.name
  target_node = n8n_workflow_node.process.name
}

resource "n8n_workflow" "example" {
  name = "My Workflow"

  nodes_json = jsonencode([
    jsondecode(n8n_workflow_node.webhook.node_json),
    jsondecode(n8n_workflow_node.process.node_json),
  ])

  connections_json = jsonencode({
    (n8n_workflow_node.webhook.name) = {
      main = [[{
        node  = n8n_workflow_node.process.name
        type  = "main"
        index = 0
      }]]
    }
  })
}
```

## Benefits

1. **Clarity**: Each node is self-contained and easy to understand
2. **Reusability**: Define common nodes once, reference them multiple times
3. **Validation**: Terraform validates node and connection configurations
4. **Refactoring**: Easy to move nodes around or change connections
5. **Collaboration**: Easier to review changes in version control (git diff shows exactly what changed)

## Testing

```bash
# Format
terraform fmt

# Validate
terraform validate

# Plan
terraform plan

# Apply
terraform apply
```

## Notes

- Both approaches (JSON blob and modular) are supported
- You can mix both approaches if needed during migration
- The `n8n_workflow_node` and `n8n_workflow_connection` resources are local-only and don't make API calls
- Only the `n8n_workflow` resource makes API calls to n8n
