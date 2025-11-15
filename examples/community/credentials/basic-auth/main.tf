# HTTP Basic Auth credential example - n8n Community Edition
terraform {
  required_providers {
    n8n = {
      source  = "kodflow/n8n"
      version = "~> 0.1.0"
    }
  }
}

provider "n8n" {
  api_url = var.n8n_api_url
  api_key = var.n8n_api_key
}

# Create HTTP Basic Auth credential
resource "n8n_credential" "http_basic_auth" {
  name = "Example Basic Auth"
  type = "httpBasicAuth"

  data = jsonencode({
    user     = var.basic_auth_user
    password = var.basic_auth_password
  })
}

# Data source to read the credential back
data "n8n_credential" "http_basic_auth" {
  id = n8n_credential.http_basic_auth.id
}

output "credential_id" {
  value       = n8n_credential.http_basic_auth.id
  description = "The ID of the created credential"
}

output "credential_name" {
  value       = data.n8n_credential.http_basic_auth.name
  description = "The name of the credential"
}

output "credential_type" {
  value       = data.n8n_credential.http_basic_auth.type
  description = "The type of the credential"
}
