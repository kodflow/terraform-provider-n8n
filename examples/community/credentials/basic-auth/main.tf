# HTTP Basic Auth credential example - n8n Community Edition
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

# Create HTTP Basic Auth credential
resource "n8n_credential" "http_basic_auth" {
  name = "ci-${var.run_id}-Example Basic Auth"
  type = "httpBasicAuth"

  data = {
    user     = var.basic_auth_user
    password = var.basic_auth_password
  }
}

output "credential_id" {
  value       = n8n_credential.http_basic_auth.id
  description = "The ID of the created credential"
}

output "credential_name" {
  value       = n8n_credential.http_basic_auth.name
  description = "The name of the credential"
}

output "credential_type" {
  value       = n8n_credential.http_basic_auth.type
  description = "The type of the credential"
}
