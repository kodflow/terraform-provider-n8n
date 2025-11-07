package project

// ProjectsDataSourceModel maps the Terraform schema to the datasource response.
// It represents a list of projects retrieved from the n8n API with all project details.
type ProjectsDataSourceModel struct {
	Projects []ProjectItemModel `tfsdk:"projects"`
}
