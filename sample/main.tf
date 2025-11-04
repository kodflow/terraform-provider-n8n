terraform {
  required_version = ">= 1.0"

  required_providers {
    n8n = {
      source  = "registry.terraform.io/kodflow/n8n"
      version = "0.0.1-dev"
    }
  }
}

provider "n8n" {
  # n8n provider configuration
  # api_url = var.n8n_api_url
  # api_key = var.n8n_api_key
}

# Example resources (to be implemented in the provider)
# resource "n8n_workflow" "example" {
#   name        = "Example Workflow"
#   active      = true
#   description = "An example workflow created via Terraform"
# }

# data "n8n_workflow" "existing" {
#   id = "workflow-id"
# }

# output "workflow_id" {
#   value = n8n_workflow.example.id
# }
