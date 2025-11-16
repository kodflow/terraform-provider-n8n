# Project Management Example
# Demonstrates creating and managing n8n projects

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
# Projects
# ============================================================================

# Create a project for development workflows
resource "n8n_project" "development" {
  name = "ci-${var.timestamp}-dev-project"
}

# Create a project for production workflows
resource "n8n_project" "production" {
  name = "ci-${var.timestamp}-prod-project"
}

# Create a project for testing
resource "n8n_project" "testing" {
  name = "ci-${var.timestamp}-test-project"
}

# ============================================================================
# Data Sources - Query Created Projects
# ============================================================================

# Query single project by ID
data "n8n_project" "dev_project" {
  id = n8n_project.development.id

  depends_on = [n8n_project.development]
}

# Query all projects
data "n8n_projects" "all" {
  depends_on = [
    n8n_project.development,
    n8n_project.production,
    n8n_project.testing
  ]
}
