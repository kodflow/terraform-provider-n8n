# Community Edition Integration Test
# Tests all community resources: tags, workflows, credentials (read-only)

terraform {
  required_providers {
    n8n = {
      source  = "kodflow/n8n"
      version = "~> 1.0"
    }
  }
}

provider "n8n" {
  base_url = var.n8n_api_url
  api_key  = var.n8n_api_key
}

# ============================================================================
# Tags
# ============================================================================

resource "n8n_tag" "test" {
  name = "community-test"
}

resource "n8n_tag" "environment" {
  name = "test-env"
}

# ============================================================================
# Workflows
# ============================================================================

# Basic workflow without nodes (inactive)
resource "n8n_workflow" "basic" {
  name   = "Community Test - Basic"
  active = false
  tags   = [n8n_tag.test.id]
}

# Workflow with nodes
resource "n8n_workflow" "with_nodes" {
  name   = "Community Test - With Nodes"
  active = false
  tags   = [n8n_tag.test.id, n8n_tag.environment.id]

  nodes_json = jsonencode([
    {
      id          = "manual-trigger"
      name        = "Manual Trigger"
      type        = "n8n-nodes-base.manualTrigger"
      position    = [250, 300]
      typeVersion = 1
      parameters  = {}
    },
    {
      id          = "set-node"
      name        = "Set Values"
      type        = "n8n-nodes-base.set"
      position    = [450, 300]
      typeVersion = 1
      parameters = {
        values = {
          string = [{
            name  = "status"
            value = "workflow executed"
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

  settings_json = jsonencode({
    executionOrder = "v1"
  })
}

# ============================================================================
# Data Sources
# ============================================================================

# Query all tags
data "n8n_tags" "all" {
  depends_on = [n8n_tag.test, n8n_tag.environment]
}

# Query workflows
data "n8n_workflows" "all" {
  depends_on = [n8n_workflow.basic, n8n_workflow.with_nodes]
}

# ============================================================================
# Outputs
# ============================================================================

output "tags" {
  value = {
    test_id        = n8n_tag.test.id
    environment_id = n8n_tag.environment.id
    all_tags_count = length(data.n8n_tags.all.tags)
  }
  description = "Created tags"
}

output "workflows" {
  value = {
    basic_id      = n8n_workflow.basic.id
    with_nodes_id = n8n_workflow.with_nodes.id
    total_count   = length(data.n8n_workflows.all.workflows)
  }
  description = "Created workflows"
}

output "summary" {
  value = {
    resources_created = {
      tags      = 2
      workflows = 2
    }
    data_sources_queried = {
      all_tags      = true
      all_workflows = true
    }
  }
  description = "Test summary"
}
