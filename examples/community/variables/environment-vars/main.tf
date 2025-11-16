# Variables example - n8n Community Edition
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
# Note: This requires appropriate API permissions
# Commented out by default due to 403 Forbidden on some instances
# data "n8n_variables" "all" {
#   depends_on = [
#     n8n_variable.api_endpoint,
#     n8n_variable.api_timeout,
#     n8n_variable.environment
#   ]
# }

# Query specific variable
data "n8n_variable" "api_endpoint" {
  id = n8n_variable.api_endpoint.id

  depends_on = [n8n_variable.api_endpoint]
}

# Output showing all created variables
output "created_variables" {
  value = {
    api_endpoint = {
      id  = n8n_variable.api_endpoint.id
      key = n8n_variable.api_endpoint.key
    }
    api_timeout = {
      id  = n8n_variable.api_timeout.id
      key = n8n_variable.api_timeout.key
    }
    environment = {
      id  = n8n_variable.environment.id
      key = n8n_variable.environment.key
    }
  }
  description = "All created variables (values hidden for security)"
}

output "api_endpoint_id" {
  value       = n8n_variable.api_endpoint.id
  description = "The ID of the API endpoint variable"
}

output "queried_api_endpoint_key" {
  value       = data.n8n_variable.api_endpoint.key
  description = "API endpoint key queried by ID"
}
