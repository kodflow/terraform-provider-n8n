# Executions query example - n8n Community Edition
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

# Create a simple workflow to generate executions
resource "n8n_workflow" "test_workflow" {
  name   = "Test Workflow for Executions"
  active = false

  nodes = jsonencode([
    {
      id         = "manual-trigger"
      name       = "When clicking 'Test workflow'"
      type       = "n8n-nodes-base.manualTrigger"
      position   = [250, 300]
      parameters = {}
    },
    {
      id       = "set-node"
      name     = "Set"
      type     = "n8n-nodes-base.set"
      position = [450, 300]
      parameters = {
        values = {
          string = [{
            name  = "result"
            value = "execution completed"
          }]
        }
      }
    }
  ])

  connections = jsonencode({
    "When clicking 'Test workflow'" = {
      main = [[{
        node  = "Set"
        type  = "main"
        index = 0
      }]]
    }
  })
}

# Query all executions
data "n8n_executions" "all" {
  depends_on = [n8n_workflow.test_workflow]
}

# Query executions for specific workflow
data "n8n_executions" "workflow_executions" {
  workflow_id = n8n_workflow.test_workflow.id
  depends_on  = [n8n_workflow.test_workflow]
}

# Query only successful executions
data "n8n_executions" "successful" {
  status     = "success"
  depends_on = [n8n_workflow.test_workflow]
}

output "workflow_id" {
  value       = n8n_workflow.test_workflow.id
  description = "The ID of the test workflow"
}

output "total_executions" {
  value       = length(data.n8n_executions.all.executions)
  description = "Total number of executions"
}

output "workflow_executions_count" {
  value       = length(data.n8n_executions.workflow_executions.executions)
  description = "Number of executions for the test workflow"
}

output "successful_executions_count" {
  value       = length(data.n8n_executions.successful.executions)
  description = "Number of successful executions"
}

output "recent_executions" {
  value = [
    for exec in slice(data.n8n_executions.all.executions, 0, min(3, length(data.n8n_executions.all.executions))) : {
      id          = exec.id
      workflow_id = exec.workflow_id
      status      = exec.status
      mode        = exec.mode
    }
  ]
  description = "Most recent 3 executions (limited output)"
}
