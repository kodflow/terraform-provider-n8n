# Basic workflow example - n8n Community Edition
terraform {
  required_providers {
    n8n = {
      source  = "kodflow/terraform-provider-n8n"
      version = "~> 0.1.0"
    }
  }
}

provider "n8n" {
  api_url = var.n8n_api_url
  api_key = var.n8n_api_key
}

# Create a simple webhook workflow
resource "n8n_workflow" "webhook_example" {
  name   = "Simple Webhook Workflow"
  active = true

  nodes = jsonencode([
    {
      id       = "webhook-node"
      name     = "Webhook"
      type     = "n8n-nodes-base.webhook"
      position = [250, 300]
      parameters = {
        path         = "example-webhook"
        httpMethod   = "POST"
        responseMode = "onReceived"
        responseData = "firstEntryJson"
      }
      webhookId = "example-webhook-id"
    },
    {
      id       = "set-node"
      name     = "Set Values"
      type     = "n8n-nodes-base.set"
      position = [450, 300]
      parameters = {
        values = {
          string = [
            {
              name  = "message"
              value = "Hello from Terraform!"
            }
          ]
        }
      }
    }
  ])

  connections = jsonencode({
    "Webhook" = {
      main = [[{
        node  = "Set Values"
        type  = "main"
        index = 0
      }]]
    }
  })

  settings = jsonencode({
    saveExecutionProgress    = true
    saveManualExecutions     = true
    saveDataErrorExecution   = "all"
    saveDataSuccessExecution = "all"
    executionTimeout         = 3600
    timezone                 = "America/New_York"
  })
}

# Data source to read the workflow back
data "n8n_workflow" "webhook_example" {
  id = n8n_workflow.webhook_example.id
}

output "workflow_id" {
  value       = n8n_workflow.webhook_example.id
  description = "The ID of the created workflow"
}

output "workflow_webhook_url" {
  value       = "http://localhost:5678/webhook/example-webhook"
  description = "The webhook URL to trigger the workflow"
}

output "workflow_active" {
  value       = data.n8n_workflow.webhook_example.active
  description = "Workflow activation status"
}
