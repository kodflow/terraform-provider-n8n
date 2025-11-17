# Trigger Nodes Showcase
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

# Acuity Scheduling Trigger
# Handle Acuity Scheduling events via webhooks
resource "n8n_workflow_node" "acuity_scheduling_trigger" {
  name     = "Acuity Scheduling Trigger"
  type     = "acuitySchedulingTrigger"
  position = [250, 300]

  parameters = jsonencode({
    # Add node-specific parameters here
  })
}

# Bitbucket Trigger
# Handle Bitbucket events via webhooks
resource "n8n_workflow_node" "bitbucket_trigger" {
  name     = "Bitbucket Trigger"
  type     = "bitbucketTrigger"
  position = [450, 300]

  parameters = jsonencode({
    # Add node-specific parameters here
  })
}

# Cal.com Trigger
# Handle Cal.com events via webhooks
resource "n8n_workflow_node" "cal_com_trigger" {
  name     = "Cal.com Trigger"
  type     = "calTrigger"
  position = [650, 300]

  parameters = jsonencode({
    # Add node-specific parameters here
  })
}

# Calendly Trigger
# Starts the workflow when Calendly events occur
resource "n8n_workflow_node" "calendly_trigger" {
  name     = "Calendly Trigger"
  type     = "calendlyTrigger"
  position = [850, 300]

  parameters = jsonencode({
    # Add node-specific parameters here
  })
}

# Email Trigger (IMAP)
# Triggers the workflow when a new email is received
resource "n8n_workflow_node" "email_trigger__imap_" {
  name     = "Email Trigger (IMAP)"
  type     = "emailReadImap"
  position = [1050, 300]

  parameters = jsonencode({
    # Add node-specific parameters here
  })
}

# Example workflow combining the nodes above
resource "n8n_workflow" "trigger_showcase" {
  name   = "ci-${var.run_id}-Trigger Showcase"
  active = false

  nodes_json = jsonencode([
    jsondecode(n8n_workflow_node.acuity_scheduling_trigger.node_json),
    jsondecode(n8n_workflow_node.bitbucket_trigger.node_json),
    jsondecode(n8n_workflow_node.cal_com_trigger.node_json),
    jsondecode(n8n_workflow_node.calendly_trigger.node_json),
    jsondecode(n8n_workflow_node.email_trigger__imap_.node_json)
  ])

  connections_json = jsonencode({})
}