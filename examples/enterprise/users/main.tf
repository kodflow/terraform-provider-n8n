# Enterprise Users Example
# Requires instance owner credentials
#
# Note: User management is only available for the instance owner.
# Users can be invited with specific roles.
# The API only supports updating the user's role, not other fields.

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
# User Resources
# ============================================================================

# Create an admin user
resource "n8n_user" "admin" {
  email = "admin-ci-${var.run_id}@example.com"
  role  = "global:admin"
}

# Create a member user
resource "n8n_user" "member" {
  email = "member-ci-${var.run_id}@example.com"
  role  = "global:member"
}

# ============================================================================
# Data Sources
# ============================================================================

data "n8n_users" "all" {
  depends_on = [n8n_user.admin, n8n_user.member]
}

data "n8n_user" "admin" {
  id         = n8n_user.admin.id
  depends_on = [n8n_user.admin]
}

# ============================================================================
# Outputs
# ============================================================================

output "users" {
  value = {
    admin = {
      id         = n8n_user.admin.id
      email      = n8n_user.admin.email
      role       = n8n_user.admin.role
      is_pending = n8n_user.admin.is_pending
    }
    member = {
      id         = n8n_user.member.id
      email      = n8n_user.member.email
      role       = n8n_user.member.role
      is_pending = n8n_user.member.is_pending
    }
  }
  description = "Created users"
  sensitive   = true
}

output "total_users" {
  value       = length(data.n8n_users.all.users)
  description = "Total number of users"
}

output "summary" {
  value = {
    resources_created = 2
    data_sources      = ["n8n_users", "n8n_user"]
  }
  description = "Test summary"
}
