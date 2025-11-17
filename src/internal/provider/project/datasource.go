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
	"github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/project/models"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/client"
)

// Ensure ProjectDataSource implements required interfaces.
var (
	_ datasource.DataSource              = &ProjectDataSource{}
	_ ProjectDataSourceInterface         = &ProjectDataSource{}
	_ datasource.DataSourceWithConfigure = &ProjectDataSource{}
)

// ProjectDataSourceInterface defines the interface for ProjectDataSource.
type ProjectDataSourceInterface interface {
	datasource.DataSource
	Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse)
	Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse)
	Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse)
	Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse)
}

// ProjectDataSource is a Terraform datasource that provides read-only access to a single n8n project.
// It fetches project details from the n8n API using ID or name-based filtering.
type ProjectDataSource struct {
	// client is the N8n API client used for operations.
	client *client.N8nClient
}

// NewProjectDataSource creates a new ProjectDataSource instance.
//
// Returns:
//   - datasource.DataSource: New ProjectDataSource instance
func NewProjectDataSource() *ProjectDataSource {
	// Return result.
	return &ProjectDataSource{}
}

// NewProjectDataSourceWrapper creates a new ProjectDataSource instance for Terraform.
// This wrapper function is used by the provider to maintain compatibility with the framework.
//
// Returns:
//   - datasource.DataSource: the wrapped ProjectDataSource instance
func NewProjectDataSourceWrapper() datasource.DataSource {
	// Return the wrapped datasource instance.
	return NewProjectDataSource()
}

// Metadata returns the data source type name.
//
// Params:
//   - ctx: Context for the request
//   - req: Metadata request
//   - resp: Metadata response to populate
func (d *ProjectDataSource) Metadata(_ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

// Schema defines the schema for the data source.
//
// Params:
//   - ctx: Context for the request
//   - req: Schema request
//   - resp: Schema response to populate
func (d *ProjectDataSource) Schema(_ctx context.Context, _req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches a single n8n project by ID or name. Since the n8n API doesn't provide a GET /projects/{id} endpoint, this datasource uses the LIST endpoint with client-side filtering.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Project identifier. Either `id` or `name` must be specified.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Project name. Either `id` or `name` must be specified.",
				Optional:            true,
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
		},
	}
}

// Configure adds the provider configured client to the data source.
//
// Params:
//   - ctx: Context for the request
//   - req: Configure request
//   - resp: Configure response to populate
func (d *ProjectDataSource) Configure(_ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Check for nil value.
	if req.ProviderData == nil {
		// Return result.
		return
	}

	clientData, ok := req.ProviderData.(*client.N8nClient)
	// Check condition.
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.N8nClient, got: %T", req.ProviderData),
		)
		// Return result.
		return
	}

	d.client = clientData
}

// Read refreshes the Terraform state with the latest data.
//
// Params:
//   - ctx: Context for the request
//   - req: Read request
//   - resp: Read response to populate
func (d *ProjectDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	data := &models.DataSource{}

	resp.Diagnostics.Append(req.Config.Get(ctx, data)...)
	// Check for diagnostics errors.
	if resp.Diagnostics.HasError() {
		// Return with error.
		return
	}

	// Validate that at least one identifier is provided.
	if data.ID.IsNull() && data.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Missing Required Attribute",
			"Either 'id' or 'name' must be specified",
		)
		// Return with error.
		return
	}

	// List all projects and filter client-side (API limitation).
	projectList, httpResp, err := d.client.APIClient.ProjectsAPI.ProjectsGet(ctx).Execute()
	// Close response body if available.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}
	// Check for API errors.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error listing projects",
			fmt.Sprintf("Could not list projects: %s\nHTTP Response: %v", err.Error(), httpResp),
		)
		// Return with error.
		return
	}

	// Find project by ID or name.
	var project *n8nsdk.Project
	var found bool
	// Search for project in the list.
	if projectList.Data != nil {
		project, found = findProjectByIDOrName(projectList.Data, data.ID, data.Name)
	}

	// Check if project was found.
	if !found {
		identifier := data.ID.ValueString()
		// Use name if ID is empty.
		if identifier == "" {
			identifier = data.Name.ValueString()
		}
		resp.Diagnostics.AddError(
			"Project Not Found",
			fmt.Sprintf("Could not find project with identifier: %s", identifier),
		)
		// Return with error.
		return
	}

	// Map project to model.
	mapProjectToDataSourceModel(project, data)

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}
