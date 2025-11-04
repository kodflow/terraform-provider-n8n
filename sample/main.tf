terraform {
  required_version = ">= 1.0"

  required_providers {
    n8n = {
      source  = "registry.terraform.io/kodflow/n8n"
      version = "0.0.1"
    }
  }
}

provider "n8n" {
  # Configuration du provider n8n
  # api_url = var.n8n_api_url
  # api_key = var.n8n_api_key
}

# Exemple de ressources (à implémenter dans le provider)
# resource "n8n_workflow" "example" {
#   name        = "Example Workflow"
#   active      = true
#   description = "Un workflow d'exemple créé via Terraform"
# }

# data "n8n_workflow" "existing" {
#   id = "workflow-id"
# }

# output "workflow_id" {
#   value = n8n_workflow.example.id
# }
