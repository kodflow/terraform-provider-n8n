# Examples Architecture

This document explains the organization of Terraform examples in this provider.

## Directory Structure

```
examples/
├── nodes/                           # Per-node complete workflow examples (296 nodes)
│   ├── core/
│   │   ├── code/                    # 1 folder per node
│   │   │   ├── main.tf             # Complete workflow testing this node
│   │   │   ├── variables.tf
│   │   │   └── README.md           # Node-specific documentation
│   │   ├── if/
│   │   ├── merge/
│   │   ├── set/
│   │   └── switch/
│   ├── trigger/
│   │   ├── webhook/
│   │   ├── schedule/
│   │   ├── manual/
│   │   └── ... (25 trigger nodes)
│   └── integration/
│       ├── http-request/
│       ├── postgres/
│       ├── slack/
│       └── ... (266 integration nodes)
│
├── community/                       # Community-contributed examples
│   ├── workflows/
│   │   └── getting-started/        # Simple 3-node tutorial workflow
│   ├── credentials/
│   └── tags/
│
├── comprehensive/                   # Production-ready complete examples
│   ├── data-pipeline/              # Complete ETL pipeline
│   ├── api-integration/            # Full API workflow with auth
│   └── notification-system/        # Multi-channel notifications
│
├── enterprise/                      # Enterprise features examples
│   ├── projects/
│   ├── source-control/
│   └── users/
│
└── basic-sample/                    # Quick start example
```

## Example Types

### 1. Per-Node Examples (`examples/nodes/`)

**Purpose**: Test and document EACH individual node with a complete workflow.

**Requirements for each node:**
- ✅ Complete, working workflow
- ✅ Demonstrates node's primary function
- ✅ Shows all important parameters
- ✅ Includes input nodes if required
- ✅ Includes output nodes if applicable
- ✅ Full connection graph
- ✅ Ready to `terraform apply`

**Example structure for `examples/nodes/core/code/main.tf`:**
```hcl
# Workflow: Test Code node
# Demonstrates: JavaScript execution, data transformation

# INPUT: Manual trigger to start workflow
resource "n8n_workflow_node" "manual" {
  name     = "Manual Trigger"
  type     = "n8n-nodes-base.manualTrigger"
  position = [250, 300]
}

# TESTED NODE: Code node with JavaScript
resource "n8n_workflow_node" "code" {
  name     = "Process Data"
  type     = "n8n-nodes-base.code"
  position = [450, 300]

  parameters = jsonencode({
    mode   = "runOnceForAllItems"
    jsCode = <<-EOT
      // Example: Transform input data
      const items = $input.all();
      return items.map(item => ({
        json: {
          original: item.json,
          processed: true,
          timestamp: new Date().toISOString()
        }
      }));
    EOT
  })
}

# OUTPUT: Display result
resource "n8n_workflow_node" "display" {
  name     = "Display Result"
  type     = "n8n-nodes-base.set"
  position = [650, 300]

  parameters = jsonencode({
    mode = "manual"
    fields = {
      values = [{
        name  = "result"
        value = "={{ $json }}"
      }]
    }
  })
}

# Connections
resource "n8n_workflow_connection" "manual_to_code" {
  source_node = n8n_workflow_node.manual.name
  target_node = n8n_workflow_node.code.name
}

resource "n8n_workflow_connection" "code_to_display" {
  source_node = n8n_workflow_node.code.name
  target_node = n8n_workflow_node.display.name
}

# Assemble workflow
resource "n8n_workflow" "test_code_node" {
  name   = "Test: Code Node"
  active = false

  nodes_json = jsonencode([
    jsondecode(n8n_workflow_node.manual.node_json),
    jsondecode(n8n_workflow_node.code.node_json),
    jsondecode(n8n_workflow_node.display.node_json),
  ])

  connections_json = jsonencode({
    (n8n_workflow_node.manual.name) = {
      main = [[{
        node  = n8n_workflow_node.code.name
        type  = "main"
        index = 0
      }]]
    }
    (n8n_workflow_node.code.name) = {
      main = [[{
        node  = n8n_workflow_node.display.name
        type  = "main"
        index = 0
      }]]
    }
  })
}
```

### 2. Community Examples (`examples/community/`)

**Purpose**: Pedagogical examples for learning.

- Simple, easy to understand
- Step-by-step tutorials
- Focus on concepts, not production use

### 3. Comprehensive Examples (`examples/comprehensive/`)

**Purpose**: Production-ready, real-world scenarios.

- Complete business logic
- Error handling
- Multiple integrations
- Best practices

## Generation Strategy

### Automated Generation

For the 296 per-node examples, we use:

1. **`scripts/nodes/generate-node-examples.js`**
   - Reads `data/n8n-nodes-registry.json`
   - For each node:
     - Determines required inputs/outputs
     - Generates appropriate trigger node
     - Creates workflow with connections
     - Adds example parameters
     - Writes main.tf, variables.tf, README.md

2. **Smart defaults based on node category:**
   - **Trigger nodes**: No input needed, add output node
   - **Core nodes**: Add manual trigger input + display output
   - **Integration nodes**: Add manual trigger + credentials example

### Manual Curation

Some nodes require manual customization:
- Nodes with complex credentials
- Nodes requiring external services
- Nodes with special connection patterns

## Testing

Each per-node example should:
1. ✅ Pass `terraform validate`
2. ✅ Pass `terraform plan`
3. ⚠️ `terraform apply` may require credentials/services

## Maintenance

- Auto-regenerate when n8n adds/updates nodes
- Manual review for node-specific optimizations
- Community contributions welcome
