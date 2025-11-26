# Enterprise Project Users Example
# Requires n8n Enterprise license AND instance owner credentials
#
# Note: This example demonstrates assigning users to projects with roles.
# Projects require Enterprise license, and user management requires instance owner.
# This test may fail with Community Edition or non-owner credentials.

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
# Project Resource
# ============================================================================

resource "n8n_project" "team" {
  name = "ci-${var.run_id}-Team Project"
}

# ============================================================================
# User Resources
# ============================================================================

resource "n8n_user" "project_admin" {
  email = "project-admin-ci-${var.run_id}@example.com"
  role  = "global:member"
}

resource "n8n_user" "project_member" {
  email = "project-member-ci-${var.run_id}@example.com"
  role  = "global:member"
}

# ============================================================================
# Project User Assignments
# ============================================================================

resource "n8n_project_user" "admin_assignment" {
  project_id = n8n_project.team.id
  user_id    = n8n_user.project_admin.id
  role       = "project:admin"
}

resource "n8n_project_user" "member_assignment" {
  project_id = n8n_project.team.id
  user_id    = n8n_user.project_member.id
  role       = "project:editor"
}

# ============================================================================
# Data Sources
# ============================================================================

data "n8n_project" "team" {
  id         = n8n_project.team.id
  depends_on = [n8n_project.team]
}

# ============================================================================
# Outputs
# ============================================================================

output "project" {
  value = {
    id   = n8n_project.team.id
    name = n8n_project.team.name
    type = n8n_project.team.type
  }
  description = "Created project"
}

output "users" {
  value = {
    admin = {
      id         = n8n_user.project_admin.id
      email      = n8n_user.project_admin.email
      role       = n8n_user.project_admin.role
      is_pending = n8n_user.project_admin.is_pending
    }
    member = {
      id         = n8n_user.project_member.id
      email      = n8n_user.project_member.email
      role       = n8n_user.project_member.role
      is_pending = n8n_user.project_member.is_pending
    }
  }
  description = "Created users"
  sensitive   = true
}

output "assignments" {
  value = {
    admin_assignment = {
      id         = n8n_project_user.admin_assignment.id
      project_id = n8n_project_user.admin_assignment.project_id
      user_id    = n8n_project_user.admin_assignment.user_id
      role       = n8n_project_user.admin_assignment.role
    }
    member_assignment = {
      id         = n8n_project_user.member_assignment.id
      project_id = n8n_project_user.member_assignment.project_id
      user_id    = n8n_project_user.member_assignment.user_id
      role       = n8n_project_user.member_assignment.role
    }
  }
  description = "Project user assignments"
}

output "summary" {
  value = {
    resources_created = 5
    data_sources      = ["n8n_project"]
  }
  description = "Test summary"
}
