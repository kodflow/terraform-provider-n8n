// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// Package project implements n8n project management resources and data sources.
package project

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/project/models"
)

// stringValue defines the minimal interface for string value types.
// This interface allows functions to accept any string value type rather
// than being bound to a concrete basetypes.StringValue implementation.
type stringValue interface {
	// IsNull returns true if the value is null.
	IsNull() bool
	// ValueString returns the underlying string value.
	ValueString() string
}

// findProjectByIDOrName searches for a project by ID or name in a project list.
//
// Params:
//   - projects: List of projects to search through
//   - id: Project ID to match
//   - name: Project name to match
//
// Returns:
//   - project: Pointer to the found project
//   - ok: True if project was found, false otherwise
func findProjectByIDOrName(projects []n8nsdk.Project, id, name stringValue) (project *n8nsdk.Project, ok bool) {
	//: Iterate over all projects to find a match by ID or name.
	for _, p := range projects {
		matchByID := !id.IsNull() && p.Id != nil && *p.Id == id.ValueString()
		matchByName := !name.IsNull() && p.Name == name.ValueString()

		//: Return immediately when a match is found.
		if matchByID || matchByName {
			//: Return pointer to matched project and success indicator.
			return &p, true
		}
	}
	//: Return nil and false when no project matches the search criteria.
	return nil, false
}

// nullableProjectIcon defines the minimal interface for nullable project icon values.
type nullableProjectIcon interface {
	// Get returns the project icon or nil.
	Get() *n8nsdk.ProjectIcon
	// IsSet returns true if the value is set.
	IsSet() bool
}

// nullableString defines the minimal interface for nullable string values.
type nullableString interface {
	// Get returns the string pointer or nil.
	Get() *string
	// IsSet returns true if the value is set.
	IsSet() bool
}

// extractIconValue extracts the icon string from a nullable ProjectIcon.
//
// Params:
//   - icon: nullable project icon value
//
// Returns:
//   - value: the icon as a Terraform string value
func extractIconValue(icon nullableProjectIcon) (value types.String) {
	//: Return null when icon is not set or nil.
	if !icon.IsSet() || icon.Get() == nil {
		//: Return null string for unset icon.
		return types.StringNull()
	}
	projectIcon := icon.Get()
	//: Return the icon string value if not nil, otherwise null.
	if projectIcon.Value != nil {
		//: Return the icon string value.
		return types.StringValue(*projectIcon.Value)
	}
	//: Return null string for nil icon value.
	return types.StringNull()
}

// extractDescriptionValue extracts the description string from a nullable string.
//
// Params:
//   - desc: nullable string description value
//
// Returns:
//   - value: the description as a Terraform string value
func extractDescriptionValue(desc nullableString) (value types.String) {
	//: Return description value when set and non-nil.
	if desc.IsSet() && desc.Get() != nil {
		//: Return the description string value.
		return types.StringPointerValue(desc.Get())
	}
	//: Return null string when description is not set.
	return types.StringNull()
}

// mapProjectToDataSourceModel maps an SDK project to the datasource model.
//
// Params:
//   - project: SDK project object to map
//   - data: Target datasource model to populate
func mapProjectToDataSourceModel(project *n8nsdk.Project, data *models.DataSource) {
	//: Set ID if the project has one.
	if project.Id != nil {
		data.ID = types.StringValue(*project.Id)
	}
	data.Name = types.StringValue(project.Name)
	//: Set Type if the project has one.
	if project.Type != nil {
		data.Type = types.StringPointerValue(project.Type)
	}
	//: Set CreatedAt if the project has a creation timestamp.
	if project.CreatedAt != nil {
		data.CreatedAt = types.StringValue(project.CreatedAt.String())
	}
	//: Set UpdatedAt if the project has an update timestamp.
	if project.UpdatedAt != nil {
		data.UpdatedAt = types.StringValue(project.UpdatedAt.String())
	}
	//: Extract icon and description using helper functions.
	data.Icon = extractIconValue(project.Icon)
	data.Description = extractDescriptionValue(project.Description)
}

// mapProjectToItem maps an SDK project to the project item model for datasources.
//
// Params:
//   - project: SDK project object to map
//
// Returns:
//   - item: Mapped project item model
func mapProjectToItem(project *n8nsdk.Project) (item models.Item) {
	result := &models.Item{
		Name: types.StringValue(project.Name),
	}

	//: Set ID if the project has one.
	if project.Id != nil {
		result.ID = types.StringValue(*project.Id)
	}
	//: Set Type if the project has one.
	if project.Type != nil {
		result.Type = types.StringPointerValue(project.Type)
	}
	//: Set CreatedAt if the project has a creation timestamp.
	if project.CreatedAt != nil {
		result.CreatedAt = types.StringValue(project.CreatedAt.String())
	}
	//: Set UpdatedAt if the project has an update timestamp.
	if project.UpdatedAt != nil {
		result.UpdatedAt = types.StringValue(project.UpdatedAt.String())
	}
	//: Extract icon and description using helper functions.
	result.Icon = extractIconValue(project.Icon)
	result.Description = extractDescriptionValue(project.Description)

	//: Return the populated item.
	return *result
}
