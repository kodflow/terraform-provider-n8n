# Integration Nodes Showcase
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

# Action Network
# Consume the Action Network API
resource "n8n_workflow_node" "action_network" {
  name     = "Action Network"
  type     = "actionNetwork"
  position = [250, 300]

  parameters = jsonencode({
    # Add node-specific parameters here
  })
}

# ActiveCampaign
# Create and edit data in ActiveCampaign
resource "n8n_workflow_node" "activecampaign" {
  name     = "ActiveCampaign"
  type     = "activeCampaign"
  position = [450, 300]

  parameters = jsonencode({
    # Add node-specific parameters here
  })
}

# Adalo
# Consume Adalo API
resource "n8n_workflow_node" "adalo" {
  name     = "Adalo"
  type     = "adalo"
  position = [650, 300]

  parameters = jsonencode({
    # Add node-specific parameters here
  })
}

# Affinity
# Consume Affinity API
resource "n8n_workflow_node" "affinity" {
  name     = "Affinity"
  type     = "affinity"
  position = [850, 300]

  parameters = jsonencode({
    # Add node-specific parameters here
  })
}

# Agile CRM
# Consume Agile CRM API
resource "n8n_workflow_node" "agile_crm" {
  name     = "Agile CRM"
  type     = "agileCrm"
  position = [1050, 300]

  parameters = jsonencode({
    # Add node-specific parameters here
  })
}

# Example workflow combining the nodes above
resource "n8n_workflow" "integration_showcase" {
  name   = "ci-${var.run_id}-Integration Showcase"
  active = false

  nodes_json = jsonencode([
    jsondecode(n8n_workflow_node.action_network.node_json),
    jsondecode(n8n_workflow_node.activecampaign.node_json),
    jsondecode(n8n_workflow_node.adalo.node_json),
    jsondecode(n8n_workflow_node.affinity.node_json),
    jsondecode(n8n_workflow_node.agile_crm.node_json)
  ])

  connections_json = jsonencode({})
}