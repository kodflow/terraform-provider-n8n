terraform {
  required_providers {
    n8n = {
      source  = "kodflow/n8n"
      version = "~> 0.0.1-dev"
    }
  }
}

provider "n8n" {
  api_key  = var.n8n_api_key
  base_url = var.n8n_base_url
}

# Tags to organize resources
resource "n8n_tag" "basic_sample" {
  name = "tf:basic-sample"
}

resource "n8n_tag" "environment_dev" {
  name = "env:dev"
}

resource "n8n_tag" "test" {
  name = "test"
}

# Note: Credentials require the 'data' field which is not yet supported
# Credentials must be created manually in the n8n UI
# resource "n8n_credential" "test_api" {
#   name = "TF Test API Credential"
#   type = "httpHeaderAuth"
# }

# Note: Variables require an n8n enterprise license - commented out
# resource "n8n_variable" "api_endpoint" {
#   key   = "MOCKY_API_ENDPOINT"
#   value = "https://run.mocky.io/v3/5b065803-be37-47c5-bdfd-45ee2a3cf34e"
# }

# Basic workflow for testing
# Note: Workflow nodes and connections are not yet supported by the provider
# These must be configured manually in the n8n UI after creation
resource "n8n_workflow" "basic_example" {
  name   = "TF Basic Sample - API Test"
  active = false
  tags   = [n8n_tag.basic_sample.id, n8n_tag.environment_dev.id]
}

# Second workflow to test multiple resources
resource "n8n_workflow" "data_processor" {
  name   = "TF Basic Sample - Data Processor"
  active = false
  tags   = [n8n_tag.basic_sample.id, n8n_tag.test.id]
}

# Data source to list all workflows
data "n8n_workflows" "all_workflows" {
  depends_on = [
    n8n_workflow.basic_example,
    n8n_workflow.data_processor
  ]
}
