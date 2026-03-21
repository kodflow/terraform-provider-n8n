// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// Package execution implements execution management resources.
package execution

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/execution/models"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/client"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/constants"
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

// ExecutionDataSource is a Terraform datasource that provides read-only access to a single n8n execution.
// It fetches execution metadata from the n8n API by execution ID.
type ExecutionDataSource struct {
	// client is the N8n API client used for operations.
	client *client.N8nClient
}

// NewExecutionDataSource creates a new ExecutionDataSource instance.
//
// Returns:
//   - *ExecutionDataSource: A new ExecutionDataSource instance
func NewExecutionDataSource() (executionDataSource *ExecutionDataSource) {
	//: Return a new empty datasource instance.
	return &ExecutionDataSource{}
}

// NewExecutionDataSourceWrapper creates a new ExecutionDataSource instance for Terraform.
// This wrapper function is used by the provider to maintain compatibility with the framework.
//
// Returns:
//   - datasource.DataSource: the wrapped ExecutionDataSource instance
func NewExecutionDataSourceWrapper() (dataSource datasource.DataSource) {
	//: Return the wrapped datasource instance.
	return NewExecutionDataSource()
}

// Metadata returns the data source type name.
//
// Params:
//   - ctx: The request context
//   - req: The metadata request containing provider type information
//   - resp: The metadata response to populate with type name
//
// Returns:
//   - None
func (d *ExecutionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	//: Return early if context is cancelled.
	if ctx.Err() != nil {
		//: Return result.
		return
	}
	resp.TypeName = req.ProviderTypeName + "_execution"
}

// Schema defines the schema for the data source.
//
// Params:
//   - ctx: The request context
//   - req: The schema request from Terraform
//   - resp: The schema response to populate
//
// Returns:
//   - None
func (d *ExecutionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	//: Return early if context is cancelled.
	if ctx.Err() != nil {
		//: Return result.
		return
	}
	//: Apply the schema definition when the request is the default schema request.
	if req == (datasource.SchemaRequest{}) {
		resp.Schema = schema.Schema{
			MarkdownDescription: "Fetches a single n8n execution by ID.",
			Attributes: map[string]schema.Attribute{
				"id": schema.StringAttribute{
					MarkdownDescription: "Execution identifier",
					Required:            true,
				},
				"mode": schema.StringAttribute{
					MarkdownDescription: "Execution mode (e.g. manual, trigger, webhook)",
					Computed:            true,
				},
				"status": schema.StringAttribute{
					MarkdownDescription: "Execution status (e.g. success, error, running, waiting)",
					Computed:            true,
				},
				"workflow_id": schema.StringAttribute{
					MarkdownDescription: "Workflow identifier associated with this execution",
					Computed:            true,
				},
				"finished": schema.BoolAttribute{
					MarkdownDescription: "Whether the execution has finished",
					Computed:            true,
				},
				"created_at": schema.StringAttribute{
					MarkdownDescription: "Timestamp when the execution was created",
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
			},
		}
	}
}

// Configure adds the provider configured client to the data source.
//
// Params:
//   - ctx: The request context
//   - req: The configure request containing provider data
//   - resp: The configure response to handle errors
//
// Returns:
//   - None
func (d *ExecutionDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	//: Return early if context is cancelled.
	if ctx.Err() != nil {
		//: Return result.
		return
	}
	//: Skip configuration when provider data is not yet available.
	if req.ProviderData == nil {
		//: Return early without error when provider data is nil.
		return
	}

	clientData, ok := req.ProviderData.(*client.N8nClient)
	//: Validate that the provider data is the expected N8nClient type.
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.N8nClient, got: %T", req.ProviderData),
		)
		//: Return after adding type mismatch error to diagnostics.
		return
	}

	d.client = clientData
}

// Read refreshes the Terraform state with the latest data.
//
// Params:
//   - ctx: The request context
//   - req: The read request containing configuration
//   - resp: The read response to populate with state
//
// Returns:
//   - None
func (d *ExecutionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	data := &models.DataSource{}

	resp.Diagnostics.Append(req.Config.Get(ctx, data)...)
	//: Return early if config parsing produced errors.
	if resp.Diagnostics.HasError() {
		//: Return with error after diagnostics are populated.
		return
	}

	//: Return early if execution retrieval failed.
	if !d.executeReadLogic(ctx, data, resp) {
		//: Return after read failure.
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

// executeReadLogic fetches an execution by ID and populates the data model.
//
// Params:
//   - ctx: The request context
//   - data: The data source model to populate
//   - resp: The read response
//
// Returns:
//   - bool: True if retrieval succeeded, false otherwise
func (d *ExecutionDataSource) executeReadLogic(ctx context.Context, data *models.DataSource, resp *datasource.ReadResponse) (ok bool) {
	//: Parse the execution ID from string to float32 for the API.
	val, err := strconv.ParseFloat(data.ID.ValueString(), constants.Float32BitSize)
	//: Return failure when the ID cannot be parsed.
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Execution ID",
			fmt.Sprintf("Could not parse execution ID %s: %s", data.ID.ValueString(), err.Error()),
		)
		//: Return false to signal parse failure.
		return false
	}

	exec, httpResp, err := d.client.APIClient.ExecutionAPI.ExecutionsIdGet(ctx, float32(val)).Execute()
	//: Close the HTTP response body to avoid resource leaks.
	if httpResp != nil && httpResp.Body != nil {
		defer func() {
			//: Silently discard close error on response body.
			if closeErr := httpResp.Body.Close(); closeErr != nil {
				resp.Diagnostics.AddWarning("Failed to close response body", closeErr.Error())
			}
		}()
	}
	//: Return failure when the API call produced an error.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error retrieving execution",
			fmt.Sprintf("Could not retrieve execution with ID %s: %s\nHTTP Response: %v",
				data.ID.ValueString(), err.Error(), httpResp),
		)
		//: Return false to signal retrieval failure.
		return false
	}

	//: Map the API response to the datasource model.
	mapExecutionToDataSource(exec, data)

	//: Return true to signal successful retrieval.
	return true
}
