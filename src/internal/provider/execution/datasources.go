// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE.md in the project root for license information.

// Package execution provides data sources for querying n8n workflow executions.
package execution

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/execution/models"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/client"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/constants"
)

// Ensure ExecutionsDataSource implements required interfaces.
var (
	_ datasource.DataSource              = &ExecutionsDataSource{}
	_ ExecutionsDataSourceInterface      = &ExecutionsDataSource{}
	_ datasource.DataSourceWithConfigure = &ExecutionsDataSource{}
)

// ExecutionsDataSourceInterface defines the interface for ExecutionsDataSource.
type ExecutionsDataSourceInterface interface {
	datasource.DataSource
	Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse)
	Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse)
	Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse)
	Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse)
}

// ExecutionsDataSource is a Terraform datasource that provides read-only access to n8n executions.
// It enables querying and filtering workflow executions from the n8n API.
type ExecutionsDataSource struct {
	client *client.N8nClient
}

// NewExecutionsDataSource creates a new ExecutionsDataSource instance.
//
// Returns:
//   - datasource.DataSource: new ExecutionsDataSource instance
func NewExecutionsDataSource() *ExecutionsDataSource {
	// Return result.
	return &ExecutionsDataSource{}
}

// NewExecutionsDataSourceWrapper creates a new ExecutionsDataSource instance for Terraform.
// This wrapper function is used by the provider to maintain compatibility with the framework.
//
// Returns:
//   - datasource.DataSource: the wrapped ExecutionsDataSource instance
func NewExecutionsDataSourceWrapper() datasource.DataSource {
	// Return the wrapped datasource instance.
	return NewExecutionsDataSource()
}

// Metadata returns the data source type name.
//
// Params:
//   - ctx: context.Context for cancellation and timeout control
//   - req: datasource.MetadataRequest containing provider type name
//   - resp: datasource.MetadataResponse to populate with metadata
func (d *ExecutionsDataSource) Metadata(_ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_executions"
}

// Schema defines the schema for the data source.
//
// Params:
//   - ctx: context.Context for cancellation and timeout control
//   - req: datasource.SchemaRequest for schema definition
//   - resp: datasource.SchemaResponse to populate with schema
func (d *ExecutionsDataSource) Schema(_ctx context.Context, _req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches a list of n8n workflow executions with optional filtering",
		Attributes:          d.schemaAttributes(),
	}
}

// schemaAttributes returns the attribute definitions for the executions datasource schema.
//
// Returns:
//   - map[string]schema.Attribute: the datasource attribute definitions
func (d *ExecutionsDataSource) schemaAttributes() map[string]schema.Attribute {
	// Return executions datasource schema attributes.
	return map[string]schema.Attribute{
		"workflow_id": schema.StringAttribute{
			MarkdownDescription: "Filter executions by workflow ID",
			Optional:            true,
		},
		"project_id": schema.StringAttribute{
			MarkdownDescription: "Filter executions by project ID",
			Optional:            true,
		},
		"status": schema.StringAttribute{
			MarkdownDescription: "Filter executions by status (e.g., 'success', 'error', 'running')",
			Optional:            true,
		},
		"include_data": schema.BoolAttribute{
			MarkdownDescription: "Whether to include execution data in the response",
			Optional:            true,
		},
		"executions": schema.ListNestedAttribute{
			MarkdownDescription: "List of executions",
			Computed:            true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: d.executionItemAttributes(),
			},
		},
	}
}

// executionItemAttributes returns the nested attribute definitions for execution items.
//
// Returns:
//   - map[string]schema.Attribute: the execution item attribute definitions
func (d *ExecutionsDataSource) executionItemAttributes() map[string]schema.Attribute {
	// Return execution item schema attributes.
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "Execution identifier",
			Computed:            true,
		},
		"workflow_id": schema.StringAttribute{
			MarkdownDescription: "ID of the workflow that was executed",
			Computed:            true,
		},
		"finished": schema.BoolAttribute{
			MarkdownDescription: "Whether the execution finished",
			Computed:            true,
		},
		"mode": schema.StringAttribute{
			MarkdownDescription: "Execution mode",
			Computed:            true,
		},
		"started_at": schema.StringAttribute{
			MarkdownDescription: "Timestamp when the execution started",
			Computed:            true,
		},
		"stopped_at": schema.StringAttribute{
			MarkdownDescription: "Timestamp when the execution stopped",
			Computed:            true,
		},
		"status": schema.StringAttribute{
			MarkdownDescription: "Execution status",
			Computed:            true,
		},
	}
}

// Configure adds the provider configured client to the data source.
//
// Params:
//   - ctx: context.Context for cancellation and timeout control
//   - req: datasource.ConfigureRequest containing provider data
//   - resp: datasource.ConfigureResponse to populate with diagnostics
func (d *ExecutionsDataSource) Configure(_ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Check for nil value.
	if req.ProviderData == nil {
		// Return with error.
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
//   - ctx: context.Context for cancellation and timeout control
//   - req: datasource.ReadRequest containing configuration data
//   - resp: datasource.ReadResponse to populate with state and diagnostics
func (d *ExecutionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Initialize data model
	data := &models.DataSources{}

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	// Check for errors in diagnostics
	if resp.Diagnostics.HasError() {
		// Return with error.
		return
	}

	apiReq := d.buildExecutionsAPIRequest(ctx, data)

	executionList, httpResp, err := apiReq.Execute()
	// Close HTTP response body if present
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}
	// Check for API execution errors
	if err != nil {
		resp.Diagnostics.AddError(
			"Error listing executions",
			fmt.Sprintf("Could not list executions: %s\nHTTP Response: %v", err.Error(), httpResp),
		)
		// Return with error.
		return
	}

	d.populateExecutionsList(data, executionList)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// buildExecutionsAPIRequest builds the API request with optional filters.
//
// Params:
//   - ctx: context.Context for cancellation
//   - data: models.DataSources containing filter parameters
//
// Returns:
//   - apiReq: Configured API request
func (d *ExecutionsDataSource) buildExecutionsAPIRequest(
	ctx context.Context,
	data *models.DataSources,
) n8nsdk.ExecutionAPIExecutionsGetRequest {
	apiReq := d.client.APIClient.ExecutionAPI.ExecutionsGet(ctx)
	// Filter by workflow ID if provided
	if !data.WorkflowID.IsNull() {
		apiReq = apiReq.WorkflowId(data.WorkflowID.ValueString())
	}
	// Filter by project ID if provided
	if !data.ProjectID.IsNull() {
		apiReq = apiReq.ProjectId(data.ProjectID.ValueString())
	}
	// Filter by status if provided
	if !data.Status.IsNull() {
		apiReq = apiReq.Status(data.Status.ValueString())
	}
	// Set include data flag if provided
	if !data.IncludeData.IsNull() {
		apiReq = apiReq.IncludeData(data.IncludeData.ValueBool())
	}
	// Return configured API request.
	return apiReq
}

// populateExecutionsList populates the executions list from API response.
//
// Params:
//   - data: models.DataSources to populate
//   - executionList: API response containing execution data
func (d *ExecutionsDataSource) populateExecutionsList(
	data *models.DataSources,
	executionList *n8nsdk.ExecutionList,
) {
	data.Executions = make([]models.Item, 0, constants.DEFAULT_LIST_CAPACITY)
	// Check if execution data is available and populate executions
	if executionList.Data != nil {
		// Process each execution and add to the executions list
		for _, execution := range executionList.Data {
			item := mapExecutionToItem(&execution)
			data.Executions = append(data.Executions, item)
		}
	}
}
