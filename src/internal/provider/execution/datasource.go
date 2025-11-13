// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE.md in the project root for license information.

// Package execution implements data sources and resources for n8n workflow executions.
package execution

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/execution/models"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
	"github.com/kodflow/n8n/src/internal/provider/shared/constants"
)

// Ensure ExecutionDataSource implements required interfaces.
var (
	_ datasource.DataSource              = &ExecutionDataSource{}
	_ ExecutionDataSourceInterface       = &ExecutionDataSource{}
	_ datasource.DataSourceWithConfigure = &ExecutionDataSource{}
)

// ExecutionDataSourceInterface defines the interface for ExecutionDataSource.
type ExecutionDataSourceInterface interface {
	datasource.DataSource
	Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse)
	Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse)
	Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse)
	Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse)
}

// ExecutionDataSource defines the data source implementation for a single execution.
// It provides read-only access to workflow execution details in n8n, including
// status, timing information, and execution mode via the n8n API.
type ExecutionDataSource struct {
	// client is the N8n API client used for execution operations.
	client *client.N8nClient
}

// NewExecutionDataSource creates a new ExecutionDataSource instance.
//
// Returns:
//   - datasource.DataSource: new ExecutionDataSource instance
func NewExecutionDataSource() *ExecutionDataSource {
	// Return result.
	return &ExecutionDataSource{}
}

// NewExecutionDataSourceWrapper creates a new ExecutionDataSource instance for Terraform.
// This wrapper function is used by the provider to maintain compatibility with the framework.
//
// Returns:
//   - datasource.DataSource: the wrapped ExecutionDataSource instance
func NewExecutionDataSourceWrapper() datasource.DataSource {
	// Return the wrapped datasource instance.
	return NewExecutionDataSource()
}

// Metadata returns the data source type name.
//
// Params:
//   - ctx: context.Context for cancellation and timeout control
//   - req: datasource.MetadataRequest containing provider type name
//   - resp: datasource.MetadataResponse to populate with metadata
func (d *ExecutionDataSource) Metadata(_ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_execution"
}

// Schema defines the schema for the data source.
//
// Params:
//   - ctx: context.Context for cancellation and timeout control
//   - req: datasource.SchemaRequest for schema definition
//   - resp: datasource.SchemaResponse to populate with schema
func (d *ExecutionDataSource) Schema(_ctx context.Context, _req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches a single n8n workflow execution by ID.",
		Attributes:          d.schemaAttributes(),
	}
}

// schemaAttributes returns the attribute definitions for the execution datasource schema.
//
// Returns:
//   - map[string]schema.Attribute: the datasource attribute definitions
func (d *ExecutionDataSource) schemaAttributes() map[string]schema.Attribute {
	// Return execution datasource schema attributes.
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "Execution identifier",
			Required:            true,
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
			MarkdownDescription: "Execution mode (e.g., 'manual', 'trigger', 'webhook')",
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
			MarkdownDescription: "Execution status (e.g., 'success', 'error', 'running')",
			Computed:            true,
		},
		"include_data": schema.BoolAttribute{
			MarkdownDescription: "Whether to include execution data in the response",
			Optional:            true,
		},
	}
}

// Configure adds the provider configured client to the data source.
//
// Params:
//   - ctx: context.Context for cancellation and timeout control
//   - req: datasource.ConfigureRequest containing provider data
//   - resp: datasource.ConfigureResponse to populate with diagnostics
func (d *ExecutionDataSource) Configure(_ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *ExecutionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Initialize data model
	data := &models.DataSource{}

	resp.Diagnostics.Append(req.Config.Get(ctx, data)...)
	// Check condition.
	if resp.Diagnostics.HasError() {
		// Return with error.
		return
	}

	// Convert ID string to float32 as required by the API
	executionID, err := strconv.ParseFloat(data.ID.ValueString(), constants.FLOAT32_BIT_SIZE)
	// Check for error.
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Execution ID",
			fmt.Sprintf("Could not parse execution ID %s as number: %s", data.ID.ValueString(), err.Error()),
		)
		// Return result.
		return
	}

	execution, httpResp, err := d.client.APIClient.ExecutionAPI.ExecutionsIdGet(ctx, float32(executionID)).Execute()
	// Check for non-nil value.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}
	// Check for error.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error retrieving execution",
			fmt.Sprintf("Could not retrieve execution with ID %s: %s\nHTTP Response: %v", data.ID.ValueString(), err.Error(), httpResp),
		)
		// Return result.
		return
	}

	// Populate data model
	populateExecutionData(execution, data)

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

// populateExecutionData populates the execution data model from the API response.
//
// Params:
//   - execution: the execution response from the API
//   - data: the data model to populate
func populateExecutionData(execution *n8nsdk.Execution, data *models.DataSource) {
	// Check if execution ID is available.
	if execution.Id != nil {
		data.ID = types.StringValue(fmt.Sprintf("%v", *execution.Id))
	}
	// Check if workflow ID is available.
	if execution.WorkflowId != nil {
		data.WorkflowID = types.StringValue(fmt.Sprintf("%v", *execution.WorkflowId))
	}
	// Check if finished status is available.
	if execution.Finished != nil {
		data.Finished = types.BoolPointerValue(execution.Finished)
	}
	// Check if execution mode is available.
	if execution.Mode != nil {
		data.Mode = types.StringPointerValue(execution.Mode)
	}
	// Check if started timestamp is available.
	if execution.StartedAt != nil {
		data.StartedAt = types.StringValue(execution.StartedAt.String())
	}
	// Check if stopped timestamp is available.
	if execution.StoppedAt.IsSet() {
		stoppedAt := execution.StoppedAt.Get()
		// Check if stopped timestamp value is not nil.
		if stoppedAt != nil {
			data.StoppedAt = types.StringValue(stoppedAt.String())
		}
	}
	// Check if execution status is available.
	if execution.Status != nil {
		data.Status = types.StringPointerValue(execution.Status)
	}
}
