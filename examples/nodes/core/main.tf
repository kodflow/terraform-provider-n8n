# Core Nodes Showcase
# Auto-generated from n8n repository

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

# Code
# Run custom JavaScript or Python code
resource "n8n_workflow_node" "code" {
  name     = "Code"
  type     = "code"
  position = [250, 300]

  parameters = jsonencode({
    # Add node-specific parameters here
  })
}

# If
# Route items to different branches (true/false)
resource "n8n_workflow_node" "if" {
  name     = "If"
  type     = "if"
  position = [450, 300]

  parameters = jsonencode({
    # Add node-specific parameters here
  })
}

# Merge
# Merges data of multiple streams once data from both is available
resource "n8n_workflow_node" "merge" {
  name     = "Merge"
  type     = "merge"
  position = [650, 300]

  parameters = jsonencode({
    # Add node-specific parameters here
  })
}

# Set
# Add or edit fields on an input item and optionally remove other fields
resource "n8n_workflow_node" "set" {
  name     = "Set"
  type     = "set"
  position = [850, 300]

  parameters = jsonencode({
    # Add node-specific parameters here
  })
}

# Switch
# Route items depending on defined expression or rules
resource "n8n_workflow_node" "switch" {
  name     = "Switch"
  type     = "switch"
  position = [1050, 300]

  parameters = jsonencode({
    # Add node-specific parameters here
  })
}

# Example workflow combining the nodes above
resource "n8n_workflow" "core_showcase" {
  name   = "ci-${var.run_id}-Core Showcase"
  active = false

  nodes_json = jsonencode([
    jsondecode(n8n_workflow_node.code.node_json),
    jsondecode(n8n_workflow_node.if.node_json),
    jsondecode(n8n_workflow_node.merge.node_json),
    jsondecode(n8n_workflow_node.set.node_json),
    jsondecode(n8n_workflow_node.switch.node_json)
  ])

  connections_json = jsonencode({})
}