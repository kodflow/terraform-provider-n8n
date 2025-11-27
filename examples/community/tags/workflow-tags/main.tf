# Tags example - n8n Community Edition
terraform {
  required_providers {
    n8n = {
      source  = "kodflow/n8n"
      version = ">= 1.0"
    }
  }
}

provider "n8n" {
  base_url = var.n8n_base_url
  api_key  = var.n8n_api_key
}

# Create tags for organization with unique CI-prefixed names
# Using run_id (GitHub run number) + timestamp for guaranteed uniqueness
resource "n8n_tag" "production" {
  name = "ci-${var.run_id}-${var.timestamp}-prod"
}

resource "n8n_tag" "automated" {
  name = "ci-${var.run_id}-${var.timestamp}-auto"
}

resource "n8n_tag" "api" {
  name = "ci-${var.run_id}-${var.timestamp}-api"
}

# Create a workflow with tags
resource "n8n_workflow" "tagged_workflow" {
  name       = "ci-${var.run_id}-Tagged Workflow Example"
  project_id = var.project_id != "" ? var.project_id : null
  active     = false
  tags       = [n8n_tag.production.id, n8n_tag.automated.id, n8n_tag.api.id]

  nodes_json = jsonencode([
    {
      id         = "manual-trigger"
      name       = "Manual Trigger"
      type       = "n8n-nodes-base.manualTrigger"
      position   = [250, 300]
      parameters = {}
    },
    {
      id       = "set-node"
      name     = "Set Values"
      type     = "n8n-nodes-base.set"
      position = [450, 300]
      parameters = {
        values = {
          string = [{
            name  = "status"
            value = "tagged workflow executed"
          }]
        }
      }
    }
  ])

  connections_json = jsonencode({
    "Manual Trigger" = {
      main = [[{
        node  = "Set Values"
        type  = "main"
        index = 0
      }]]
    }
  })
}

# Query all tags
data "n8n_tags" "all" {}

# Query all workflows
data "n8n_workflows" "all_workflows" {}

output "all_tags" {
  value       = data.n8n_tags.all.tags
  description = "All tags in the system"
}

output "tagged_workflow_id" {
  value       = n8n_workflow.tagged_workflow.id
  description = "The ID of the tagged workflow"
}

output "all_workflows" {
  value       = data.n8n_workflows.all_workflows.workflows
  description = "All workflows in the system"
}
