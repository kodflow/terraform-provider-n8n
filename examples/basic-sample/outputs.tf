# Tags
output "tag_basic_id" {
  description = "ID of the basic-sample tag"
  value       = n8n_tag.basic_sample.id
}

output "tag_env_id" {
  description = "ID of the environment tag"
  value       = n8n_tag.environment_dev.id
}

output "tag_test_id" {
  description = "ID of the test tag"
  value       = n8n_tag.test.id
}

# Workflows
output "workflow_api_test_id" {
  description = "ID of the API test workflow"
  value       = n8n_workflow.basic_example.id
}

output "workflow_data_processor_id" {
  description = "ID of the data processor workflow"
  value       = n8n_workflow.data_processor.id
}

# URLs
output "workflow_api_test_url" {
  description = "Direct URL to the API test workflow"
  value       = "${var.n8n_base_url}/workflow/${n8n_workflow.basic_example.id}"
}

output "workflow_data_processor_url" {
  description = "Direct URL to the data processor workflow"
  value       = "${var.n8n_base_url}/workflow/${n8n_workflow.data_processor.id}"
}

# Summary
output "summary" {
  description = "Summary of all created resources"
  value = {
    tags = {
      basic_sample = {
        id   = n8n_tag.basic_sample.id
        name = n8n_tag.basic_sample.name
      }
      environment_dev = {
        id   = n8n_tag.environment_dev.id
        name = n8n_tag.environment_dev.name
      }
      test = {
        id   = n8n_tag.test.id
        name = n8n_tag.test.name
      }
    }
    workflows = {
      api_test = {
        id     = n8n_workflow.basic_example.id
        name   = n8n_workflow.basic_example.name
        active = n8n_workflow.basic_example.active
        url    = "${var.n8n_base_url}/workflow/${n8n_workflow.basic_example.id}"
      }
      data_processor = {
        id     = n8n_workflow.data_processor.id
        name   = n8n_workflow.data_processor.name
        active = n8n_workflow.data_processor.active
        url    = "${var.n8n_base_url}/workflow/${n8n_workflow.data_processor.id}"
      }
    }
    total_resources = 5
  }
}
