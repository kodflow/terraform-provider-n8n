# Test workflow for Switch
# Category: Core
# Type: n8n-nodes-base.switch

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

# TESTED NODE: Switch
resource "n8n_workflow_node" "test_node" {
  name     = "Switch"
  type     = "n8n-nodes-base.switch"
  position = [450, 300]

  parameters = jsonencode(
    {
      "mode" : "rules",
      "rules" : {
        "values" : [
          {
            "value" : "={{ $json.type === \"A\" }}"
          },
          {
            "value" : "={{ $json.type === \"B\" }}"
          }
        ]
      }
    }
  )
}

# OUTPUT 1: Output 1 (First matching rule)
resource "n8n_workflow_node" "output_0" {
  name     = "Output: Output 1"
  type     = "n8n-nodes-base.set"
  position = [650, 75]

  parameters = jsonencode({
    mode = "manual"
    fields = {
      values = [{
        name  = "output_type"
        type  = "string"
        value = "Output 1"
        }, {
        name  = "result"
        type  = "string"
        value = "={{ $json }}"
      }]
    }
  })
}

# OUTPUT 2: Output 2 (Second matching rule)
resource "n8n_workflow_node" "output_1" {
  name     = "Output: Output 2"
  type     = "n8n-nodes-base.set"
  position = [650, 225]

  parameters = jsonencode({
    mode = "manual"
    fields = {
      values = [{
        name  = "output_type"
        type  = "string"
        value = "Output 2"
        }, {
        name  = "result"
        type  = "string"
        value = "={{ $json }}"
      }]
    }
  })
}

# OUTPUT 3: Output 3 (Third matching rule)
resource "n8n_workflow_node" "output_2" {
  name     = "Output: Output 3"
  type     = "n8n-nodes-base.set"
  position = [650, 375]

  parameters = jsonencode({
    mode = "manual"
    fields = {
      values = [{
        name  = "output_type"
        type  = "string"
        value = "Output 3"
        }, {
        name  = "result"
        type  = "string"
        value = "={{ $json }}"
      }]
    }
  })
}

# OUTPUT 4: Fallback (No rules matched)
resource "n8n_workflow_node" "output_3" {
  name     = "Output: Fallback"
  type     = "n8n-nodes-base.set"
  position = [650, 525]

  parameters = jsonencode({
    mode = "manual"
    fields = {
      values = [{
        name  = "output_type"
        type  = "string"
        value = "Fallback"
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

# Connection from test node output[0] (Output 1) to output node
resource "n8n_workflow_connection" "test_to_output_0" {
  source_node         = n8n_workflow_node.test_node.name
  source_output       = "main"
  source_output_index = 0
  target_node         = n8n_workflow_node.output_0.name
  target_input        = "main"
  target_input_index  = 0
}

# Connection from test node output[1] (Output 2) to output node
resource "n8n_workflow_connection" "test_to_output_1" {
  source_node         = n8n_workflow_node.test_node.name
  source_output       = "main"
  source_output_index = 1
  target_node         = n8n_workflow_node.output_1.name
  target_input        = "main"
  target_input_index  = 0
}

# Connection from test node output[2] (Output 3) to output node
resource "n8n_workflow_connection" "test_to_output_2" {
  source_node         = n8n_workflow_node.test_node.name
  source_output       = "main"
  source_output_index = 2
  target_node         = n8n_workflow_node.output_2.name
  target_input        = "main"
  target_input_index  = 0
}

# Connection from test node output[3] (Fallback) to output node
resource "n8n_workflow_connection" "test_to_output_3" {
  source_node         = n8n_workflow_node.test_node.name
  source_output       = "main"
  source_output_index = 3
  target_node         = n8n_workflow_node.output_3.name
  target_input        = "main"
  target_input_index  = 0
}

# WORKFLOW
resource "n8n_workflow" "test_switch" {
  name   = "Test: Switch"
  active = false

  nodes_json = jsonencode([
    jsondecode(n8n_workflow_node.manual_trigger.node_json),
    jsondecode(n8n_workflow_node.test_node.node_json),
    jsondecode(n8n_workflow_node.output_0.node_json),
    jsondecode(n8n_workflow_node.output_1.node_json),
    jsondecode(n8n_workflow_node.output_2.node_json),
    jsondecode(n8n_workflow_node.output_3.node_json)
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
        }],
        [{
          node  = n8n_workflow_node.output_2.name
          type  = "main"
          index = 0
        }],
        [{
          node  = n8n_workflow_node.output_3.name
          type  = "main"
          index = 0
        }]
      ]
    }
  })
}

# OUTPUTS
output "workflow_id" {
  value       = n8n_workflow.test_switch.id
  description = "ID of the test workflow"
}

output "workflow_name" {
  value       = n8n_workflow.test_switch.name
  description = "Name of the test workflow"
}
