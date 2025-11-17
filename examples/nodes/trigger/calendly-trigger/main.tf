# Test workflow for Calendly Trigger
# Category: Trigger
# Type: n8n-nodes-base.calendlyTrigger

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

# TESTED NODE: Calendly Trigger
resource "n8n_workflow_node" "test_node" {
  name     = "Calendly Trigger"
  type     = "n8n-nodes-base.calendlyTrigger"
  position = [250, 300]

  parameters = jsonencode(
    {
        "note": "Configure Calendly Trigger parameters here"
    }
  )
}

# OUTPUT: Display result
resource "n8n_workflow_node" "display_result" {
  name     = "Display Result"
  type     = "n8n-nodes-base.set"
  position = [450, 300]

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
resource "n8n_workflow_connection" "test_to_output" {
  source_node         = n8n_workflow_node.test_node.name
  source_output       = "main"
  source_output_index = 0
  target_node         = n8n_workflow_node.display_result.name
  target_input        = "main"
  target_input_index  = 0
}

# WORKFLOW
resource "n8n_workflow" "test_calendly-trigger" {
  name   = "Test: Calendly Trigger"
  active = false

  nodes_json = jsonencode([
    jsondecode(n8n_workflow_node.test_node.node_json),
    jsondecode(n8n_workflow_node.display_result.node_json)
  ])

  connections_json = jsonencode({
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
  value       = n8n_workflow.test_calendly-trigger.id
  description = "ID of the test workflow"
}

output "workflow_name" {
  value       = n8n_workflow.test_calendly-trigger.name
  description = "Name of the test workflow"
}
