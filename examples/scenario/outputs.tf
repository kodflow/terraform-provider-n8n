# Tags
output "tag_scenario_id" {
  description = "ID of the scenario tag"
  value       = n8n_tag.scenario_tag.id
}

output "tag_scenario_name" {
  description = "Name of the scenario tag"
  value       = n8n_tag.scenario_tag.name
}

output "tag_environment_id" {
  description = "ID of the environment tag"
  value       = n8n_tag.environment_tag.id
}

output "tag_environment_name" {
  description = "Name of the environment tag"
  value       = n8n_tag.environment_tag.name
}

# Workflow
output "workflow_id" {
  description = "ID of the scenario workflow"
  value       = n8n_workflow.scenario_workflow.id
}

output "workflow_name" {
  description = "Name of the scenario workflow"
  value       = n8n_workflow.scenario_workflow.name
}

output "workflow_url" {
  description = "Direct URL to the scenario workflow"
  value       = "${var.n8n_base_url}/workflow/${n8n_workflow.scenario_workflow.id}"
}

# Summary
output "summary" {
  description = "Summary of test resources"
  value = {
    name_suffix = var.name_suffix
    run_id      = var.run_id
    tags = {
      scenario = {
        id   = n8n_tag.scenario_tag.id
        name = n8n_tag.scenario_tag.name
      }
      environment = {
        id   = n8n_tag.environment_tag.id
        name = n8n_tag.environment_tag.name
      }
    }
    workflow = {
      id   = n8n_workflow.scenario_workflow.id
      name = n8n_workflow.scenario_workflow.name
    }
    data_sources = {
      total_tags      = length(data.n8n_tags.all.tags)
      total_workflows = length(data.n8n_workflows.all.workflows)
    }
  }
}
