# Tags example - n8n Community Edition
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

# Create tags for organization
resource "n8n_tag" "production" {
  name = "production"
}

resource "n8n_tag" "automated" {
  name = "automated"
}

resource "n8n_tag" "api" {
  name = "api"
}

# Create a workflow with tags
resource "n8n_workflow" "tagged_workflow" {
  name   = "Tagged Workflow Example"
  active = false
  tags   = [n8n_tag.production.id, n8n_tag.automated.id, n8n_tag.api.id]

  nodes = jsonencode([
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

  connections = jsonencode({
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

# Query workflows by tag
data "n8n_workflows" "production_workflows" {
  tags = [n8n_tag.production.id]
}

output "all_tags" {
  value       = data.n8n_tags.all.tags
  description = "All tags in the system"
}

output "tagged_workflow_id" {
  value       = n8n_workflow.tagged_workflow.id
  description = "The ID of the tagged workflow"
}

output "production_workflows" {
  value       = data.n8n_workflows.production_workflows.workflows
  description = "All workflows with production tag"
}
