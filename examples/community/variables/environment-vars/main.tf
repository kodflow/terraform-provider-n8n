# Variables example - n8n Community Edition
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

# Create environment variables
resource "n8n_variable" "api_endpoint" {
  key   = "API_ENDPOINT"
  value = "https://api.example.com/v1"
}

resource "n8n_variable" "api_timeout" {
  key   = "API_TIMEOUT"
  value = "30000"
}

resource "n8n_variable" "environment" {
  key   = "ENVIRONMENT"
  value = var.environment
}

# Query all variables
data "n8n_variables" "all" {}

# Query specific variable
data "n8n_variable" "api_endpoint" {
  id = n8n_variable.api_endpoint.id
}

output "all_variables" {
  value = [
    for v in data.n8n_variables.all.variables : {
      key = v.key
      id  = v.id
    }
  ]
  description = "All environment variables (values hidden)"
}

output "api_endpoint_id" {
  value       = n8n_variable.api_endpoint.id
  description = "The ID of the API endpoint variable"
}
