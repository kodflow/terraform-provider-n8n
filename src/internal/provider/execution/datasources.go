// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// Package execution implements execution management resources.
package execution

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/execution/models"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/client"
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

// ExecutionsDataSource is a Terraform datasource that provides read-only access to a list of n8n executions.
// It supports optional filtering by workflow ID and status.
type ExecutionsDataSource struct {
	// client is the N8n API client used for operations.
	client *client.N8nClient
}

// NewExecutionsDataSource creates a new ExecutionsDataSource instance.
//
// Returns:
//   - *ExecutionsDataSource: A new ExecutionsDataSource instance
func NewExecutionsDataSource() (executionsDataSource *ExecutionsDataSource) {
	//: Return a new empty datasource instance.
	return &ExecutionsDataSource{}
}

// NewExecutionsDataSourceWrapper creates a new ExecutionsDataSource instance for Terraform.
// This wrapper function is used by the provider to maintain compatibility with the framework.
//
// Returns:
//   - datasource.DataSource: the wrapped ExecutionsDataSource instance
func NewExecutionsDataSourceWrapper() (dataSource datasource.DataSource) {
	//: Return the wrapped datasource instance.
	return NewExecutionsDataSource()
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
func (d *ExecutionsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	//: Return early if context is cancelled.
	if ctx.Err() != nil {
		//: Return result.
		return
	}
	resp.TypeName = req.ProviderTypeName + "_executions"
}

// executionItemSchema returns the nested attribute schema for a single execution in the list.
//
// Returns:
//   - schema.NestedAttributeObject: The nested attribute object definition for an execution item
func executionItemSchema() (nestedAttr schema.NestedAttributeObject) {
	//: Return the schema definition for execution item nested attributes.
	return schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Execution identifier",
				Computed:            true,
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

// Schema defines the schema for the data source.
//
// Params:
//   - ctx: The request context
//   - req: The schema request from Terraform
//   - resp: The schema response to populate
//
// Returns:
//   - None
func (d *ExecutionsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	//: Return early if context is cancelled.
	if ctx.Err() != nil {
		//: Return result.
		return
	}
	//: Apply the schema definition when the request is the default schema request.
	if req == (datasource.SchemaRequest{}) {
		resp.Schema = schema.Schema{
			MarkdownDescription: "Fetches a list of n8n executions with optional filtering.",
			Attributes: map[string]schema.Attribute{
				"workflow_id": schema.StringAttribute{
					MarkdownDescription: "Filter executions by workflow ID",
					Optional:            true,
				},
				"status": schema.StringAttribute{
					MarkdownDescription: "Filter executions by status (e.g. success, error, running, waiting)",
					Optional:            true,
				},
				"executions": schema.ListNestedAttribute{
					MarkdownDescription: "List of executions",
					Computed:            true,
					NestedObject:        executionItemSchema(),
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
func (d *ExecutionsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
//   - req: The read request from Terraform
//   - resp: The read response to populate with data
//
// Returns:
//   - None
func (d *ExecutionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.DataSources

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	//: Return early if config parsing produced errors.
	if resp.Diagnostics.HasError() {
		//: Return with error after diagnostics are populated.
		return
	}

	//: Honour deferred reads when the Terraform client supports deferral.
	if req.ClientCapabilities.DeferralAllowed {
		//: Return early when the client requests deferred evaluation.
		return
	}

	//: Return early if listing executions failed.
	if !d.executeListLogic(ctx, &data, resp) {
		//: Return after list failure.
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// executeListLogic fetches executions from the API with optional filters.
//
// Params:
//   - ctx: The request context
//   - data: The data sources model to populate
//   - resp: The read response
//
// Returns:
//   - bool: True if listing succeeded, false otherwise
func (d *ExecutionsDataSource) executeListLogic(ctx context.Context, data *models.DataSources, resp *datasource.ReadResponse) (ok bool) {
	req := d.client.APIClient.ExecutionAPI.ExecutionsGet(ctx)

	//: Apply workflow ID filter when provided.
	if isValueSet(data.WorkflowID) {
		req = req.WorkflowId(data.WorkflowID.ValueString())
	}

	//: Apply status filter when provided.
	if isValueSet(data.Status) {
		req = req.Status(data.Status.ValueString())
	}

	execList, httpResp, err := req.Execute()
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
			"Error listing executions",
			fmt.Sprintf("Could not list executions: %s\nHTTP Response: %v", err.Error(), httpResp),
		)
		//: Return false to signal listing failure.
		return false
	}

	//: Build the execution items slice from the API response.
	data.Executions = buildExecutionItems(execList)
	//: Return true to signal successful listing.
	return true
}
