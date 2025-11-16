output "development_project_id" {
  description = "ID of the development project"
  value       = n8n_project.development.id
}

output "development_project_name" {
  description = "Name of the development project"
  value       = n8n_project.development.name
}

output "development_project_type" {
  description = "Type of the development project"
  value       = n8n_project.development.type
}

output "production_project_id" {
  description = "ID of the production project"
  value       = n8n_project.production.id
}

output "production_project_name" {
  description = "Name of the production project"
  value       = n8n_project.production.name
}

output "testing_project_id" {
  description = "ID of the testing project"
  value       = n8n_project.testing.id
}

output "all_projects_count" {
  description = "Total number of projects in the instance"
  value       = length(data.n8n_projects.all.projects)
}

output "queried_dev_project_name" {
  description = "Development project name queried by ID"
  value       = data.n8n_project.dev_project.name
}

output "queried_dev_project_type" {
  description = "Development project type queried by ID"
  value       = data.n8n_project.dev_project.type
}
