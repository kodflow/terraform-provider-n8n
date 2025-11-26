# Enterprise Projects Example
# Requires n8n Enterprise license
#
# Note: Projects are only available with an Enterprise license.
# This test will fail with Community Edition with error:
# "Your license does not allow for feat:projectRole:admin"

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
# Project Resources
# ============================================================================

resource "n8n_project" "main" {
  name = "ci-${var.run_id}-Enterprise Main Project"
}

resource "n8n_project" "dev" {
  name = "ci-${var.run_id}-Enterprise Dev Project"
  # Sequential creation to avoid n8n database concurrency issues
  depends_on = [n8n_project.main]
}

# ============================================================================
# Data Sources
# ============================================================================

data "n8n_projects" "all" {
  depends_on = [n8n_project.main, n8n_project.dev]
}

data "n8n_project" "main" {
  id         = n8n_project.main.id
  depends_on = [n8n_project.main]
}

# ============================================================================
# Outputs
# ============================================================================

output "projects" {
  value = {
    main = {
      id   = n8n_project.main.id
      name = n8n_project.main.name
      type = n8n_project.main.type
    }
    dev = {
      id   = n8n_project.dev.id
      name = n8n_project.dev.name
      type = n8n_project.dev.type
    }
  }
  description = "Created projects"
}

output "total_projects" {
  value       = length(data.n8n_projects.all.projects)
  description = "Total number of projects"
}

output "summary" {
  value = {
    resources_created = 2
    data_sources      = ["n8n_projects", "n8n_project"]
  }
  description = "Test summary"
}
