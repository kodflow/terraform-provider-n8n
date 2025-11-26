# Enterprise Variables Example
# Requires n8n Enterprise license with feat:variables enabled
#
# Note: Environment variables are only available with an Enterprise license.
# This test will fail with Community Edition with error:
# "Your license does not allow for feat:variables"

terraform {
  required_providers {
    n8n = {
      source  = "kodflow/n8n"
      version = "~> 1.0"
    }
  }
}

provider "n8n" {
  api_key  = var.n8n_api_key
  base_url = var.n8n_base_url
}

# ============================================================================
# Variable Resources
# ============================================================================

resource "n8n_variable" "api_url" {
  key   = "ci-${var.run_id}-API_URL"
  value = "https://api.example.com"
  # Note: type is computed by n8n, not settable via API
}

resource "n8n_variable" "api_timeout" {
  key   = "ci-${var.run_id}-API_TIMEOUT"
  value = "30"
  # Sequential creation to avoid n8n database concurrency issues
  depends_on = [n8n_variable.api_url]
}

resource "n8n_variable" "debug_enabled" {
  key   = "ci-${var.run_id}-DEBUG_ENABLED"
  value = "true"
  # Sequential creation to avoid n8n database concurrency issues
  depends_on = [n8n_variable.api_timeout]
}

# ============================================================================
# Data Sources
# ============================================================================

data "n8n_variables" "all" {
  depends_on = [
    n8n_variable.api_url,
    n8n_variable.api_timeout,
    n8n_variable.debug_enabled
  ]
}

# ============================================================================
# Outputs
# ============================================================================

output "variables" {
  value = {
    api_url = {
      id    = n8n_variable.api_url.id
      key   = n8n_variable.api_url.key
      value = n8n_variable.api_url.value
      type  = n8n_variable.api_url.type
    }
    api_timeout = {
      id    = n8n_variable.api_timeout.id
      key   = n8n_variable.api_timeout.key
      value = n8n_variable.api_timeout.value
      type  = n8n_variable.api_timeout.type
    }
    debug_enabled = {
      id    = n8n_variable.debug_enabled.id
      key   = n8n_variable.debug_enabled.key
      value = n8n_variable.debug_enabled.value
      type  = n8n_variable.debug_enabled.type
    }
  }
  description = "Created variables"
  sensitive   = true
}

output "total_variables" {
  value       = length(data.n8n_variables.all.variables)
  description = "Total number of variables"
}

output "summary" {
  value = {
    resources_created = 3
    data_sources      = ["n8n_variables"]
  }
  description = "Test summary"
}
