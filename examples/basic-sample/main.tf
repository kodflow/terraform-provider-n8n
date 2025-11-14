# Basic Sample - Complete n8n Terraform Example
# Demonstrates tags, workflows with nodes, and data sources

terraform {
  required_providers {
    n8n = {
      source  = "kodflow/n8n"
      version = "0.0.1-dev"
    }
  }
}

provider "n8n" {
  api_key  = var.n8n_api_key
  base_url = var.n8n_base_url
}

# ============================================================================
# Tags for Organization
# ============================================================================

resource "n8n_tag" "basic_sample" {
  name = "tf:basic-sample"
}

resource "n8n_tag" "environment_dev" {
  name = "env:dev"
}

resource "n8n_tag" "automated" {
  name = "automated"
}

# ============================================================================
# Simple Workflow (no nodes)
# ============================================================================

resource "n8n_workflow" "simple" {
  name   = "TF Basic Sample - Simple Workflow"
  active = false
  tags   = [n8n_tag.basic_sample.id, n8n_tag.environment_dev.id]
}

# ============================================================================
# API Test Workflow (with nodes and connections)
# ============================================================================

resource "n8n_workflow" "api_test" {
  name   = "TF Basic Sample - API Test"
  active = false
  tags   = [n8n_tag.basic_sample.id, n8n_tag.automated.id]

  nodes_json = jsonencode([
    {
      id          = "start"
      name        = "Start"
      type        = "n8n-nodes-base.manualTrigger"
      position    = [240, 300]
      typeVersion = 1
      parameters  = {}
    },
    {
      id       = "http-request"
      name     = "HTTP Request"
      type     = "n8n-nodes-base.httpRequest"
      position = [460, 300]
      typeVersion = 4.2
      parameters = {
        url    = "https://api.github.com/repos/n8n-io/n8n"
        method = "GET"
        options = {}
      }
    },
    {
      id       = "process-data"
      name     = "Process Data"
      type     = "n8n-nodes-base.set"
      position = [680, 300]
      typeVersion = 3.4
      parameters = {
        mode = "manual"
        duplicateItem = false
        assignments = {
          assignments = [
            {
              id = "field1"
              name = "repo_name"
              type = "string"
              value = "={{ $json.name }}"
            },
            {
              id = "field2"
              name = "stars"
              type = "number"
              value = "={{ $json.stargazers_count }}"
            }
          ]
        }
      }
    }
  ])

  connections_json = jsonencode({
    Start = {
      main = [[{
        node  = "HTTP Request"
        type  = "main"
        index = 0
      }]]
    }
    "HTTP Request" = {
      main = [[{
        node  = "Process Data"
        type  = "main"
        index = 0
      }]]
    }
  })

  settings_json = jsonencode({
    executionOrder = "v1"
    saveDataErrorExecution = "all"
    saveDataSuccessExecution = "all"
    saveManualExecutions = true
  })
}

# ============================================================================
# Data Processing Workflow
# ============================================================================

resource "n8n_workflow" "data_processor" {
  name   = "TF Basic Sample - Data Processor"
  active = false
  tags   = [n8n_tag.basic_sample.id, n8n_tag.automated.id]

  nodes_json = jsonencode([
    {
      id          = "manual-trigger"
      name        = "When clicking Test workflow"
      type        = "n8n-nodes-base.manualTrigger"
      position    = [240, 300]
      typeVersion = 1
      parameters  = {}
    },
    {
      id       = "create-data"
      name     = "Create Sample Data"
      type     = "n8n-nodes-base.set"
      position = [460, 300]
      typeVersion = 3.4
      parameters = {
        mode = "manual"
        duplicateItem = false
        assignments = {
          assignments = [
            {
              id = "field1"
              name = "users"
              type = "array"
              value = jsonencode([
                {name = "Alice", age = 30},
                {name = "Bob", age = 25}
              ])
            }
          ]
        }
      }
    }
  ])

  connections_json = jsonencode({
    "When clicking Test workflow" = {
      main = [[{
        node  = "Create Sample Data"
        type  = "main"
        index = 0
      }]]
    }
  })

  settings_json = jsonencode({
    executionOrder = "v1"
  })
}

# ============================================================================
# Data Sources
# ============================================================================

# Query all workflows
data "n8n_workflows" "all" {
  depends_on = [
    n8n_workflow.simple,
    n8n_workflow.api_test,
    n8n_workflow.data_processor
  ]
}

# Query all tags
data "n8n_tags" "all" {
  depends_on = [
    n8n_tag.basic_sample,
    n8n_tag.environment_dev,
    n8n_tag.automated
  ]
}
