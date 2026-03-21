# Enterprise Credentials Example
# Demonstrates reading credentials using datasources
#
# Note: Credentials must exist in n8n before being read by these datasources.
# The datasources provide read-only access to credential metadata (not secret values).

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
# Credential Resource (for testing datasources)
# ============================================================================

resource "n8n_credential" "test" {
  name = "ci_${var.run_id}_http_header"
  type = "httpHeaderAuth"
  data = jsonencode({
    name  = "Authorization"
    value = "Bearer test-token"
  })
}

# ============================================================================
# Credential Datasource (single by ID)
# ============================================================================

data "n8n_credential" "by_id" {
  id         = n8n_credential.test.id
  depends_on = [n8n_credential.test]
}

# ============================================================================
# Credentials Datasource (list)
# ============================================================================

data "n8n_credentials" "all" {
  depends_on = [n8n_credential.test]
}

# ============================================================================
# Outputs
# ============================================================================

output "credential_id" {
  description = "The ID of the created credential"
  value       = n8n_credential.test.id
}

output "credential_name" {
  description = "The name of the credential (from datasource)"
  value       = data.n8n_credential.by_id.name
}

output "credentials_count" {
  description = "Total number of credentials"
  value       = length(data.n8n_credentials.all.credentials)
}
