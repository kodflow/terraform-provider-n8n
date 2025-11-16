# Tags
output "tag_basic_id" {
  description = "ID of the basic-sample tag"
  value       = n8n_tag.basic_sample.id
}

output "tag_env_id" {
  description = "ID of the environment tag"
  value       = n8n_tag.environment_dev.id
}

output "tag_automated_id" {
  description = "ID of the automated tag"
  value       = n8n_tag.automated.id
}

# Workflows
output "workflow_simple_id" {
  description = "ID of the simple workflow"
  value       = n8n_workflow.simple.id
}

output "workflow_api_test_id" {
  description = "ID of the API test workflow"
  value       = n8n_workflow.api_test.id
}

output "workflow_data_processor_id" {
  description = "ID of the data processor workflow"
  value       = n8n_workflow.data_processor.id
}

# URLs
output "workflow_simple_url" {
  description = "Direct URL to the simple workflow"
  value       = "${var.n8n_base_url}/workflow/${n8n_workflow.simple.id}"
}

output "workflow_api_test_url" {
  description = "Direct URL to the API test workflow"
  value       = "${var.n8n_base_url}/workflow/${n8n_workflow.api_test.id}"
}

output "workflow_data_processor_url" {
  description = "Direct URL to the data processor workflow"
  value       = "${var.n8n_base_url}/workflow/${n8n_workflow.data_processor.id}"
}

# Projects - Commented out (requires Enterprise license)
# output "project_id" {
#   description = "ID of the sample project"
#   value       = n8n_project.sample_project.id
# }

# output "project_name" {
#   description = "Name of the sample project"
#   value       = n8n_project.sample_project.name
# }

# output "project_type" {
#   description = "Type of the sample project"
#   value       = n8n_project.sample_project.type
# }

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
      automated = {
        id   = n8n_tag.automated.id
        name = n8n_tag.automated.name
      }
    }
    workflows = {
      simple = {
        id     = n8n_workflow.simple.id
        name   = n8n_workflow.simple.name
        active = n8n_workflow.simple.active
        url    = "${var.n8n_base_url}/workflow/${n8n_workflow.simple.id}"
      }
      api_test = {
        id     = n8n_workflow.api_test.id
        name   = n8n_workflow.api_test.name
        active = n8n_workflow.api_test.active
        url    = "${var.n8n_base_url}/workflow/${n8n_workflow.api_test.id}"
      }
      data_processor = {
        id     = n8n_workflow.data_processor.id
        name   = n8n_workflow.data_processor.name
        active = n8n_workflow.data_processor.active
        url    = "${var.n8n_base_url}/workflow/${n8n_workflow.data_processor.id}"
      }
    }
    # projects = {  # Commented out - requires Enterprise license
    #   sample_project = {
    #     id   = n8n_project.sample_project.id
    #     name = n8n_project.sample_project.name
    #     type = n8n_project.sample_project.type
    #   }
    # }
    total_tags      = length(data.n8n_tags.all.tags)
    total_workflows = length(data.n8n_workflows.all.workflows)
    # total_projects  = length(data.n8n_projects.all.projects)  # Commented out - requires Enterprise license
  }
}
