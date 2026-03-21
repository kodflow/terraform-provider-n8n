# Enterprise Workflow Versions Example
# Demonstrates reading workflow version history
#
# Note: Workflow history is available with n8n Enterprise license.
# This example reads a specific version of an existing workflow.

terraform {
  required_providers {
    n8n = {
      source  = "kodflow/n8n"
      version = ">= 2.0"
    }
  }
}

provider "n8n" {
  api_key  = var.n8n_api_key
  base_url = var.n8n_base_url
}

# ============================================================================
# Workflow Resource (to have a workflow to read versions from)
# ============================================================================

resource "n8n_workflow" "versioned" {
  name        = "ci_${var.run_id}_versioned_workflow"
  active      = false
  nodes       = "[]"
  connections = jsonencode({})
  settings    = jsonencode({})
}

# ============================================================================
# Workflow Version Datasource
# ============================================================================

data "n8n_workflow_version" "latest" {
  workflow_id = n8n_workflow.versioned.id
  version_id  = var.version_id

  depends_on = [n8n_workflow.versioned]
}

# ============================================================================
# Outputs
# ============================================================================

output "workflow_id" {
  description = "The ID of the workflow"
  value       = n8n_workflow.versioned.id
}

output "version_name" {
  description = "The workflow name at this version"
  value       = data.n8n_workflow_version.latest.name
}

output "version_created_at" {
  description = "When this workflow version was created"
  value       = data.n8n_workflow_version.latest.created_at
}
