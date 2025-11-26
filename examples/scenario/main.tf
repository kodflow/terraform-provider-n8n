# Scenario Test - E2E Rename/Update Operations
# This example tests the complete lifecycle including rename operations:
# 1. Plan/Apply with initial names (version = "v1")
# 2. Plan/Apply with renamed resources (version = "v2")
# 3. Destroy all resources
#
# This validates that terraform update operations work correctly,
# specifically the fix for plan.ID vs state.ID in update operations.

terraform {
  required_providers {
    n8n = {
      source  = "kodflow/n8n"
      version = "~> 1.0"
    }
  }
}

provider "n8n" {
  api_key  = var.n8n_api_key
  base_url = var.n8n_base_url
}

# ============================================================================
# Tags - Testing rename operations
# ============================================================================

resource "n8n_tag" "scenario_tag" {
  name = "ci-${var.run_id}-scenario-tag-${var.name_suffix}"
}

resource "n8n_tag" "environment_tag" {
  name = "ci-${var.run_id}-env-${var.name_suffix}"
}

# ============================================================================
# Workflow - Testing rename operations
# ============================================================================

resource "n8n_workflow" "scenario_workflow" {
  name   = "ci-${var.run_id}-Scenario Workflow ${var.name_suffix}"
  active = false
  tags   = [n8n_tag.scenario_tag.id, n8n_tag.environment_tag.id]

  nodes_json = jsonencode([
    {
      id          = "manual-trigger"
      name        = "Start"
      type        = "n8n-nodes-base.manualTrigger"
      position    = [240, 300]
      typeVersion = 1
      parameters  = {}
    },
    {
      id          = "set-node"
      name        = "Set Suffix ${var.name_suffix}"
      type        = "n8n-nodes-base.set"
      position    = [460, 300]
      typeVersion = 3.4
      parameters = {
        mode          = "manual"
        duplicateItem = false
        assignments = {
          assignments = [
            {
              id    = "suffix"
              name  = "name_suffix"
              type  = "string"
              value = var.name_suffix
            },
            {
              id    = "run_id"
              name  = "run_id"
              type  = "string"
              value = var.run_id
            }
          ]
        }
      }
    }
  ])

  connections_json = jsonencode({
    Start = {
      main = [[{
        node  = "Set Suffix ${var.name_suffix}"
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
# Data Sources - Verify resources exist
# ============================================================================

data "n8n_tags" "all" {
  depends_on = [
    n8n_tag.scenario_tag,
    n8n_tag.environment_tag
  ]
}

data "n8n_workflows" "all" {
  depends_on = [n8n_workflow.scenario_workflow]
}
