# Test workflow for Merge
# Category: Core
# Type: n8n-nodes-base.merge

terraform {
  required_providers {
    n8n = {
      source  = "kodflow/n8n"
      version = "~> 1.0"
    }
  }
}

provider "n8n" {
  base_url = var.n8n_base_url
  api_key  = var.n8n_api_key
}

# INPUT 1: Input 1
resource "n8n_workflow_node" "input_0" {
  name     = "Input 1"
  type     = "n8n-nodes-base.manualTrigger"
  position = [250, 300]
}

# INPUT 2: Input 2
resource "n8n_workflow_node" "input_1" {
  name     = "Input 2"
  type     = "n8n-nodes-base.manualTrigger"
  position = [250, 450]
}

# TESTED NODE: Merge
resource "n8n_workflow_node" "test_node" {
  name     = "Merge"
  type     = "n8n-nodes-base.merge"
  position = [450, 375]

  parameters = jsonencode(
    {
        "mode": "combine",
        "mergeByFields": {
            "values": [
                {
                    "field1": "id",
                    "field2": "id"
                }
            ]
        }
    }
  )
}

# OUTPUT: Display result
resource "n8n_workflow_node" "display_result" {
  name     = "Display Result"
  type     = "n8n-nodes-base.set"
  position = [650, 300]

  parameters = jsonencode({
    mode = "manual"
    fields = {
      values = [{
        name  = "result"
        type  = "string"
        value = "={{ $json }}"
      }]
    }
  })
}

# CONNECTIONS
# Connection from Input 1 to test node
resource "n8n_workflow_connection" "input_0_to_test" {
  source_node         = n8n_workflow_node.input_0.name
  source_output       = "main"
  source_output_index = 0
  target_node         = n8n_workflow_node.test_node.name
  target_input        = "main"
  target_input_index  = 0
}

# Connection from Input 2 to test node
resource "n8n_workflow_connection" "input_1_to_test" {
  source_node         = n8n_workflow_node.input_1.name
  source_output       = "main"
  source_output_index = 0
  target_node         = n8n_workflow_node.test_node.name
  target_input        = "main"
  target_input_index  = 1
}

resource "n8n_workflow_connection" "test_to_output" {
  source_node         = n8n_workflow_node.test_node.name
  source_output       = "main"
  source_output_index = 0
  target_node         = n8n_workflow_node.display_result.name
  target_input        = "main"
  target_input_index  = 0
}

# WORKFLOW
resource "n8n_workflow" "test_merge" {
  name   = "Test: Merge"
  active = false

  nodes_json = jsonencode([
    jsondecode(n8n_workflow_node.input_0.node_json),
    jsondecode(n8n_workflow_node.input_1.node_json),
    jsondecode(n8n_workflow_node.test_node.node_json),
    jsondecode(n8n_workflow_node.display_result.node_json)
  ])

  connections_json = jsonencode({
    (n8n_workflow_node.input_0.name) = {
      main = [[{
        node  = n8n_workflow_node.test_node.name
        type  = "main"
        index = 0
      }]]
    }
    (n8n_workflow_node.input_1.name) = {
      main = [[{
        node  = n8n_workflow_node.test_node.name
        type  = "main"
        index = 1
      }]]
    }
    (n8n_workflow_node.test_node.name) = {
      main = [[{
        node  = n8n_workflow_node.display_result.name
        type  = "main"
        index = 0
      }]]
    }
  })
}

# OUTPUTS
output "workflow_id" {
  value       = n8n_workflow.test_merge.id
  description = "ID of the test workflow"
}

output "workflow_name" {
  value       = n8n_workflow.test_merge.name
  description = "Name of the test workflow"
}
