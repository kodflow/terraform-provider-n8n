package project

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/project/models"
)

// findProjectByIDOrName searches for a project by ID or name in a project list.
//
// Params:
//   - projects: List of projects to search through
//   - id: Project ID to match
//   - name: Project name to match
//
// Returns:
//   - *n8nsdk.Project: Pointer to the found project
//   - bool: True if project was found, false otherwise
func findProjectByIDOrName(projects []n8nsdk.Project, id, name types.String) (*n8nsdk.Project, bool) {
	// Iterate over items.
	for _, project := range projects {
		matchByID := !id.IsNull() && project.Id != nil && *project.Id == id.ValueString()
		matchByName := !name.IsNull() && project.Name == name.ValueString()

		// Check condition.
		if matchByID || matchByName {
			// Return result.
			return &project, true
		}
	}
	// Return failure status.
	return nil, false
}

// mapProjectToDataSourceModel maps an SDK project to the datasource model.
//
// Params:
//   - project: SDK project object to map
//   - data: Target datasource model to populate
func mapProjectToDataSourceModel(project *n8nsdk.Project, data *models.DataSource) {
	// Check for non-nil value.
	if project.Id != nil {
		data.ID = types.StringValue(*project.Id)
	}
	data.Name = types.StringValue(project.Name)
	// Check for non-nil value.
	if project.Type != nil {
		data.Type = types.StringPointerValue(project.Type)
	}
	// Check for non-nil value.
	if project.CreatedAt != nil {
		data.CreatedAt = types.StringValue(project.CreatedAt.String())
	}
	// Check for non-nil value.
	if project.UpdatedAt != nil {
		data.UpdatedAt = types.StringValue(project.UpdatedAt.String())
	}
	// Check for non-nil value.
	if project.Icon != nil {
		data.Icon = types.StringPointerValue(project.Icon)
	} else {
		// Icon is nil, explicitly set to null for Terraform state.
		data.Icon = types.StringNull()
	}
	// Check for non-nil value.
	if project.Description != nil {
		data.Description = types.StringPointerValue(project.Description)
	} else {
		// Description is nil, explicitly set to null for Terraform state.
		data.Description = types.StringNull()
	}
}

// mapProjectToItem maps an SDK project to the project item model for datasources.
//
// Params:
//   - project: SDK project object to map
//
// Returns:
//   - models.Item: Mapped project item model
func mapProjectToItem(project *n8nsdk.Project) models.Item {
	item := &models.Item{
		Name: types.StringValue(project.Name),
	}

	// Check for non-nil value.
	if project.Id != nil {
		item.ID = types.StringValue(*project.Id)
	}
	// Check for non-nil value.
	if project.Type != nil {
		item.Type = types.StringPointerValue(project.Type)
	}
	// Check for non-nil value.
	if project.CreatedAt != nil {
		item.CreatedAt = types.StringValue(project.CreatedAt.String())
	}
	// Check for non-nil value.
	if project.UpdatedAt != nil {
		item.UpdatedAt = types.StringValue(project.UpdatedAt.String())
	}
	// Check for non-nil value.
	if project.Icon != nil {
		item.Icon = types.StringPointerValue(project.Icon)
	}
	// Check for non-nil value.
	if project.Description != nil {
		item.Description = types.StringPointerValue(project.Description)
	}

	// Return result.
	return *item
}
