// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// Package project implements n8n project management resources and data sources.
package project

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/project/models"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/client"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/constants"
)

// Ensure ProjectsDataSource implements required interfaces.
var (
	_ datasource.DataSource              = &ProjectsDataSource{}
	_ ProjectsDataSourceInterface        = &ProjectsDataSource{}
	_ datasource.DataSourceWithConfigure = &ProjectsDataSource{}
)

// ProjectsDataSourceInterface defines the interface for ProjectsDataSource.
type ProjectsDataSourceInterface interface {
	datasource.DataSource
	Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse)
	Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse)
	Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse)
	Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse)
}

// ProjectsDataSource is a Terraform datasource that provides read-only access to all n8n projects.
// It enables querying and iterating through all available projects from the n8n API.
type ProjectsDataSource struct {
	// client is the N8n API client used for operations.
	client *client.N8nClient
}

// NewProjectsDataSource creates a new ProjectsDataSource instance.
//
// Returns:
//   - datasource.DataSource: a new ProjectsDataSource instance configured for accessing n8n projects
func NewProjectsDataSource() (projectsDataSource *ProjectsDataSource) {
	//: Return new empty ProjectsDataSource instance.
	return &ProjectsDataSource{}
}

// NewProjectsDataSourceWrapper creates a new ProjectsDataSource instance for Terraform.
// This wrapper function is used by the provider to maintain compatibility with the framework.
//
// Returns:
//   - datasource.DataSource: the wrapped ProjectsDataSource instance
func NewProjectsDataSourceWrapper() (dataSource datasource.DataSource) {
	//: Return the wrapped datasource instance.
	return NewProjectsDataSource()
}

// Metadata returns the data source type name.
//
// Params:
//   - ctx: context for the request
//   - req: metadata request containing provider type name
//   - resp: metadata response to populate
func (d *ProjectsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	//: Return early if context is cancelled.
	if ctx.Err() != nil {
		//: Return when context is cancelled.
		return
	}
	resp.TypeName = req.ProviderTypeName + "_projects"
}

// Schema defines the schema for the data source.
//
// Params:
//   - ctx: context for the request
//   - req: schema request from the framework
//   - resp: schema response to populate with schema definition
func (d *ProjectsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	//: Return early if context is cancelled.
	if ctx.Err() != nil {
		//: Return when context is cancelled.
		return
	}
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches a list of all n8n projects",
		//: Return schema attributes for the projects datasource.
		Attributes: d.schemaAttributes(),
	}
}

// schemaAttributes returns the attribute definitions for this datasource.
//
// Returns:
//   - m: the datasource attribute definitions
func (d *ProjectsDataSource) schemaAttributes() (m map[string]schema.Attribute) {
	//: Return schema attributes for the projects list.
	return map[string]schema.Attribute{
		"projects": schema.ListNestedAttribute{
			MarkdownDescription: "List of projects",
			Computed:            true,
			NestedObject: schema.NestedAttributeObject{
				//: Return nested project item attributes.
				Attributes: d.projectAttributes(),
			},
		},
	}
}

// projectAttributes returns the attribute definitions for a project item.
//
// Returns:
//   - m: the project item attribute definitions
func (d *ProjectsDataSource) projectAttributes() (m map[string]schema.Attribute) {
	//: Return schema attributes for a single project item.
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "Project identifier",
			Computed:            true,
		},
		"name": schema.StringAttribute{
			MarkdownDescription: "Project name",
			Computed:            true,
		},
		"type": schema.StringAttribute{
			MarkdownDescription: "Project type (e.g., 'team', 'personal')",
			Computed:            true,
		},
		"created_at": schema.StringAttribute{
			MarkdownDescription: "Timestamp when the project was created",
			Computed:            true,
		},
		"updated_at": schema.StringAttribute{
			MarkdownDescription: "Timestamp when the project was last updated",
			Computed:            true,
		},
		"icon": schema.StringAttribute{
			MarkdownDescription: "Project icon",
			Computed:            true,
		},
		"description": schema.StringAttribute{
			MarkdownDescription: "Project description",
			Computed:            true,
		},
	}
}

// Configure adds the provider configured client to the data source.
//
// Params:
//   - ctx: context for the request
//   - req: configure request containing provider data
//   - resp: configure response to report errors
func (d *ProjectsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	//: Return early if context is cancelled.
	if ctx.Err() != nil {
		//: Return when context is cancelled.
		return
	}
	//: Skip configuration when provider data is not yet available.
	if req.ProviderData == nil {
		//: Return early when no provider data is set.
		return
	}

	clientData, ok := req.ProviderData.(*client.N8nClient)
	//: Validate the provider data type is the expected N8nClient.
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.N8nClient, got: %T", req.ProviderData),
		)
		//: Return after reporting the type mismatch error.
		return
	}

	d.client = clientData
}

// Read refreshes the Terraform state with the latest data.
//
// Params:
//   - ctx: context for the request
//   - req: read request from Terraform
//   - resp: read response to populate with data
func (d *ProjectsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.DataSources

	projectList, httpResp, err := d.client.APIClient.ProjectsAPI.ProjectsGet(ctx).Execute()
	//: Close the HTTP response body if it is not nil to prevent resource leaks.
	if httpResp != nil && httpResp.Body != nil {
		defer func() {
			//: Silently discard close error on response body.
			if closeErr := httpResp.Body.Close(); closeErr != nil {
				resp.Diagnostics.AddWarning("Failed to close response body", closeErr.Error())
			}
		}()
	}
	//: Check for error from the list API call.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error listing projects",
			fmt.Sprintf("Could not list projects: %s\nHTTP Response: %v", err.Error(), httpResp),
		)
		//: Return after reporting the API error.
		return
	}

	data.Projects = make([]models.Item, 0, constants.DefaultListCapacity)
	//: Iterate through all projects and convert them to the model format.
	if projectList.Data != nil {
		//: Convert each project from the API response to the Item format.
		for _, project := range projectList.Data {
			item := mapProjectToItem(&project)
			data.Projects = append(data.Projects, item)
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
