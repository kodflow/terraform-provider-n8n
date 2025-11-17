# Test workflow for Split In Batches
# Category: Integration
# Type: n8n-nodes-base.splitInBatches

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

# INPUT: Manual trigger to start the workflow
resource "n8n_workflow_node" "manual_trigger" {
  name     = "Manual Trigger"
  type     = "n8n-nodes-base.manualTrigger"
  position = [250, 300]
}

# TESTED NODE: Split In Batches
resource "n8n_workflow_node" "test_node" {
  name     = "Split In Batches"
  type     = "n8n-nodes-base.splitInBatches"
  position = [450, 300]

  parameters = jsonencode(
    {
      "note" : "Configure Split In Batches parameters here"
    }
  )
}

# OUTPUT 1: Batch (Current batch items)
resource "n8n_workflow_node" "output_0" {
  name     = "Output: Batch"
  type     = "n8n-nodes-base.set"
  position = [650, 225]

  parameters = jsonencode({
    mode = "manual"
    fields = {
      values = [{
        name  = "output_type"
        type  = "string"
        value = "Batch"
        }, {
        name  = "result"
        type  = "string"
        value = "={{ $json }}"
      }]
    }
  })
}

# OUTPUT 2: Done (All batches processed)
resource "n8n_workflow_node" "output_1" {
  name     = "Output: Done"
  type     = "n8n-nodes-base.set"
  position = [650, 375]

  parameters = jsonencode({
    mode = "manual"
    fields = {
      values = [{
        name  = "output_type"
        type  = "string"
        value = "Done"
        }, {
        name  = "result"
        type  = "string"
        value = "={{ $json }}"
      }]
    }
  })
}

# CONNECTIONS
resource "n8n_workflow_connection" "input_to_test" {
  source_node         = n8n_workflow_node.manual_trigger.name
  source_output       = "main"
  source_output_index = 0
  target_node         = n8n_workflow_node.test_node.name
  target_input        = "main"
  target_input_index  = 0
}

# Connection from test node output[0] (Batch) to output node
resource "n8n_workflow_connection" "test_to_output_0" {
  source_node         = n8n_workflow_node.test_node.name
  source_output       = "main"
  source_output_index = 0
  target_node         = n8n_workflow_node.output_0.name
  target_input        = "main"
  target_input_index  = 0
}

# Connection from test node output[1] (Done) to output node
resource "n8n_workflow_connection" "test_to_output_1" {
  source_node         = n8n_workflow_node.test_node.name
  source_output       = "main"
  source_output_index = 1
  target_node         = n8n_workflow_node.output_1.name
  target_input        = "main"
  target_input_index  = 0
}

# WORKFLOW
resource "n8n_workflow" "test_split-in-batches" {
  name   = "Test: Split In Batches"
  active = false

  nodes_json = jsonencode([
    jsondecode(n8n_workflow_node.manual_trigger.node_json),
    jsondecode(n8n_workflow_node.test_node.node_json),
    jsondecode(n8n_workflow_node.output_0.node_json),
    jsondecode(n8n_workflow_node.output_1.node_json)
  ])

  connections_json = jsonencode({
    (n8n_workflow_node.manual_trigger.name) = {
      main = [[{
        node  = n8n_workflow_node.test_node.name
        type  = "main"
        index = 0
      }]]
    }
    (n8n_workflow_node.test_node.name) = {
      main = [
        [{
          node  = n8n_workflow_node.output_0.name
          type  = "main"
          index = 0
        }],
        [{
          node  = n8n_workflow_node.output_1.name
          type  = "main"
          index = 0
        }]
      ]
    }
  })
}

# OUTPUTS
output "workflow_id" {
  value       = n8n_workflow.test_split-in-batches.id
  description = "ID of the test workflow"
}

output "workflow_name" {
  value       = n8n_workflow.test_split-in-batches.name
  description = "Name of the test workflow"
}
