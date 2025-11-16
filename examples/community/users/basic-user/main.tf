# User Management Example
# Demonstrates creating and managing n8n users with different roles

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
# Users with Different Roles
# ============================================================================

# Create an admin user
resource "n8n_user" "admin" {
  email = "admin-ci-${var.timestamp}@example.com"
  role  = "global:admin"
}

# Create a member user
resource "n8n_user" "member" {
  email = "member-ci-${var.timestamp}@example.com"
  role  = "global:member"
}

# Create a user without specifying role (will use instance default)
resource "n8n_user" "default_role" {
  email = "user-ci-${var.timestamp}@example.com"
}

# ============================================================================
# Data Sources - Query Created Users
# ============================================================================

# Query user by ID
data "n8n_user" "admin_by_id" {
  id = n8n_user.admin.id

  depends_on = [n8n_user.admin]
}

# Query user by email
data "n8n_user" "member_by_email" {
  email = n8n_user.member.email

  depends_on = [n8n_user.member]
}

# Query all users
data "n8n_users" "all" {
  depends_on = [
    n8n_user.admin,
    n8n_user.member,
    n8n_user.default_role
  ]
}
