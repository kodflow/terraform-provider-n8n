# Scheduled workflow example - n8n Community Edition
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

# Create a scheduled workflow that runs every hour
resource "n8n_workflow" "scheduled_example" {
  name   = "ci-${var.run_id}-Hourly Scheduled Workflow"
  active = true

  nodes_json = jsonencode([
    {
      id       = "schedule-node"
      name     = "Schedule Trigger"
      type     = "n8n-nodes-base.scheduleTrigger"
      position = [250, 300]
      parameters = {
        rule = {
          interval = [{
            field         = "hours"
            hoursInterval = 1
          }]
        }
      }
    },
    {
      id       = "code-node"
      name     = "Process Data"
      type     = "n8n-nodes-base.code"
      position = [450, 300]
      parameters = {
        mode   = "runOnceForAllItems"
        jsCode = "const now = new Date();\nreturn [{ json: { timestamp: now.toISOString(), message: 'Scheduled execution' } }];"
      }
    }
  ])

  connections_json = jsonencode({
    "Schedule Trigger" = {
      main = [[{
        node  = "Process Data"
        type  = "main"
        index = 0
      }]]
    }
  })
}

output "workflow_id" {
  value       = n8n_workflow.scheduled_example.id
  description = "The ID of the created workflow"
}

output "workflow_name" {
  value       = n8n_workflow.scheduled_example.name
  description = "The name of the workflow"
}

output "workflow_active" {
  value       = n8n_workflow.scheduled_example.active
  description = "Workflow activation status"
}
