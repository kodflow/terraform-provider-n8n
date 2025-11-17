# Modular workflow example - using n8n_workflow_node and n8n_workflow_connection
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

# Define individual nodes
resource "n8n_workflow_node" "webhook_trigger" {
  name     = "Webhook"
  type     = "n8n-nodes-base.webhook"
  position = [250, 300]

  parameters = jsonencode({
    path         = "modular-webhook"
    httpMethod   = "POST"
    responseMode = "onReceived"
    responseData = "firstEntryJson"
  })

  webhook_id = "modular-webhook-id"
}

resource "n8n_workflow_node" "process_data" {
  name     = "Process Data"
  type     = "n8n-nodes-base.code"
  position = [450, 300]

  parameters = jsonencode({
    mode   = "runOnceForAllItems"
    jsCode = <<-EOT
      // Transform incoming webhook data
      const items = $input.all();
      return items.map(item => ({
        json: {
          ...item.json,
          processed: true,
          timestamp: new Date().toISOString()
        }
      }));
    EOT
  })
}

resource "n8n_workflow_node" "send_response" {
  name     = "Send Response"
  type     = "n8n-nodes-base.respondToWebhook"
  position = [650, 300]

  parameters = jsonencode({
    respondWith = "json"
    responseBody = jsonencode({
      success = true
      message = "Data processed successfully"
    })
  })
}

# Define connections between nodes
resource "n8n_workflow_connection" "webhook_to_process" {
  source_node         = n8n_workflow_node.webhook_trigger.name
  source_output       = "main"
  source_output_index = 0

  target_node        = n8n_workflow_node.process_data.name
  target_input       = "main"
  target_input_index = 0
}

resource "n8n_workflow_connection" "process_to_response" {
  source_node         = n8n_workflow_node.process_data.name
  source_output       = "main"
  source_output_index = 0

  target_node        = n8n_workflow_node.send_response.name
  target_input       = "main"
  target_input_index = 0
}

# Assemble the workflow using the modular components
resource "n8n_workflow" "modular_example" {
  name   = "ci-${var.run_id}-Modular Workflow Example"
  active = true

  # Use the node JSON from our workflow_node resources
  nodes_json = jsonencode([
    jsondecode(n8n_workflow_node.webhook_trigger.node_json),
    jsondecode(n8n_workflow_node.process_data.node_json),
    jsondecode(n8n_workflow_node.send_response.node_json),
  ])

  # Build connections JSON from our connection resources
  connections_json = jsonencode({
    (n8n_workflow_node.webhook_trigger.name) = {
      main = [[{
        node  = n8n_workflow_node.process_data.name
        type  = "main"
        index = 0
      }]]
    }
    (n8n_workflow_node.process_data.name) = {
      main = [[{
        node  = n8n_workflow_node.send_response.name
        type  = "main"
        index = 0
      }]]
    }
  })
}

output "workflow_id" {
  value       = n8n_workflow.modular_example.id
  description = "The ID of the created workflow"
}

output "workflow_webhook_url" {
  value       = "http://localhost:5678/webhook/modular-webhook"
  description = "The webhook URL to trigger the workflow"
}

output "nodes" {
  value = {
    webhook = n8n_workflow_node.webhook_trigger.node_json
    process = n8n_workflow_node.process_data.node_json
    respond = n8n_workflow_node.send_response.node_json
  }
  description = "JSON representations of all nodes"
}
