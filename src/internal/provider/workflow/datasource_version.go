// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// Package workflow implements workflow management resources and data sources.
package workflow

import (
	"context"
	"encoding/json"
	"fmt"

	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/client"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/workflow/models"
)

// Ensure WorkflowVersionDataSource implements required interfaces.
var (
	_ datasource.DataSource              = &WorkflowVersionDataSource{}
	_ WorkflowVersionDataSourceInterface = &WorkflowVersionDataSource{}
	_ datasource.DataSourceWithConfigure = &WorkflowVersionDataSource{}
)

// WorkflowVersionDataSourceInterface defines the interface for WorkflowVersionDataSource.
type WorkflowVersionDataSourceInterface interface {
	datasource.DataSource
	Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse)
	Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse)
	Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse)
	Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse)
}

// WorkflowVersionDataSource is a Terraform datasource that provides read-only access to a specific workflow version.
// It retrieves a specific version snapshot of a workflow from the n8n API.
type WorkflowVersionDataSource struct {
	// client is the N8n API client used for operations.
	client *client.N8nClient
}

// NewWorkflowVersionDataSource creates a new WorkflowVersionDataSource instance.
//
// Returns:
//   - *WorkflowVersionDataSource: A new WorkflowVersionDataSource instance
func NewWorkflowVersionDataSource() (workflowVersionDataSource *WorkflowVersionDataSource) {
	//: Return result.
	return &WorkflowVersionDataSource{}
}

// NewWorkflowVersionDataSourceWrapper creates a new WorkflowVersionDataSource instance for Terraform.
// This wrapper function is used by the provider to maintain compatibility with the framework.
//
// Returns:
//   - datasource.DataSource: the wrapped WorkflowVersionDataSource instance
func NewWorkflowVersionDataSourceWrapper() (dataSource datasource.DataSource) {
	//: Return the wrapped datasource instance.
	return NewWorkflowVersionDataSource()
}

// Metadata returns the data source type name.
//
// Params:
//   - ctx: The request context
//   - req: The metadata request containing provider type information
//   - resp: The metadata response to populate with type name
func (d *WorkflowVersionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	//: Guard against cancelled context.
	if ctx.Err() != nil {
		//: Return early when context is cancelled.
		return
	}
	resp.TypeName = req.ProviderTypeName + "_workflow_version"
}

// Schema defines the schema for the data source.
//
// Params:
//   - ctx: The request context
//   - req: The schema request from Terraform
//   - resp: The schema response to populate
func (d *WorkflowVersionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	//: Guard against cancelled context.
	if ctx.Err() != nil {
		//: Return early when context is cancelled.
		return
	}
	//: Acknowledge empty schema request carried by framework convention.
	schemaReq := req
	_ = schemaReq
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches a specific version of an n8n workflow from workflow history.",
		Attributes: map[string]schema.Attribute{
			"workflow_id": schema.StringAttribute{
				MarkdownDescription: "Workflow identifier",
				Required:            true,
			},
			"version_id": schema.StringAttribute{
				MarkdownDescription: "Version identifier for the specific workflow snapshot",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Workflow name at this version",
				Computed:            true,
			},
			"authors": schema.StringAttribute{
				MarkdownDescription: "Authors who created this version",
				Computed:            true,
			},
			"nodes_json": schema.StringAttribute{
				MarkdownDescription: "Workflow nodes as a JSON string at this version",
				Computed:            true,
			},
			"connections_json": schema.StringAttribute{
				MarkdownDescription: "Workflow connections as a JSON string at this version",
				Computed:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "Timestamp when this version was created",
				Computed:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
//
// Params:
//   - ctx: The request context
//   - req: The configure request containing provider data
//   - resp: The configure response to handle errors
func (d *WorkflowVersionDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	//: Guard against cancelled context.
	if ctx.Err() != nil {
		resp.Diagnostics.AddError("context cancelled", ctx.Err().Error())
		//: Return early when context is cancelled.
		return
	}
	//: Return early when provider data is nil.
	if req.ProviderData == nil {
		//: Return early without error when provider data is nil.
		return
	}

	clientData, ok := req.ProviderData.(*client.N8nClient)
	//: Check if provider data is correct type.
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.N8nClient, got: %T", req.ProviderData),
		)
		//: Return early on type mismatch.
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
func (d *WorkflowVersionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	data := &models.DataSourceVersion{}

	resp.Diagnostics.Append(req.Config.Get(ctx, data)...)
	//: If there are errors from config parsing, return early.
	if resp.Diagnostics.HasError() {
		//: Return with error.
		return
	}

	//: Execute read logic.
	if !d.executeReadLogic(ctx, data, resp) {
		//: Return with error.
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

// executeReadLogic retrieves a specific workflow version from the API.
//
// Params:
//   - ctx: The request context
//   - data: The data source model
//   - resp: The read response
//
// Returns:
//   - bool: True if read succeeded, false otherwise
func (d *WorkflowVersionDataSource) executeReadLogic(ctx context.Context, data *models.DataSourceVersion, resp *datasource.ReadResponse) (ok bool) {
	version, httpResp, err := d.client.APIClient.WorkflowAPI.WorkflowsIdVersionIdGet(
		ctx,
		data.WorkflowID.ValueString(),
		data.VersionID.ValueString(),
	).Execute()
	//: Close HTTP response body if present.
	if httpResp != nil && httpResp.Body != nil {
		defer func() {
			//: Handle body close error.
			if closeErr := httpResp.Body.Close(); closeErr != nil {
				resp.Diagnostics.AddWarning("Failed to close response body", closeErr.Error())
			}
		}()
	}
	//: Check if API call returned an error.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading workflow version",
			fmt.Sprintf("Could not read workflow %s version %s: %s\nHTTP Response: %v",
				data.WorkflowID.ValueString(), data.VersionID.ValueString(), err.Error(), httpResp),
		)
		//: Return failure.
		return false
	}

	//: Map version metadata fields.
	data.WorkflowID = types.StringValue(version.WorkflowID)
	data.VersionID = types.StringValue(version.VersionID)
	mapVersionOptionalFields(version.Name, version.Authors, version.CreatedAt, data)

	//: Serialize nodes and then connections, return result.
	return serializeVersionFields(version, data, resp)
}

// serializeVersionFields serializes both nodes and connections to JSON strings.
//
// Params:
//   - version: the workflow version from SDK
//   - data: the data source model to update
//   - resp: the read response for error reporting
//
// Returns:
//   - bool: True if serialization succeeded, false otherwise
func serializeVersionFields(version *n8nsdk.WorkflowVersion, data *models.DataSourceVersion, resp *datasource.ReadResponse) (ok bool) {
	//: Serialize nodes to JSON.
	if !serializeVersionNodes(version.Nodes, data, resp) {
		//: Return failure.
		return false
	}

	//: Serialize connections inline to avoid passing interface{} as parameter.
	connectionsJSON, jsonErr := json.Marshal(version.Connections)
	//: Check for connections serialization error.
	if jsonErr != nil {
		resp.Diagnostics.AddError(
			"Error serializing workflow connections",
			fmt.Sprintf("Could not serialize workflow connections to JSON: %s", jsonErr.Error()),
		)
		//: Return failure.
		return false
	}
	data.ConnectionsJSON = types.StringValue(string(connectionsJSON))

	//: Return success.
	return true
}

// serializeVersionNodes serializes workflow nodes to a JSON string.
//
// Params:
//   - nodes: the workflow nodes from SDK
//   - data: the data source model to update
//   - resp: the read response for error reporting
//
// Returns:
//   - bool: True if serialization succeeded, false otherwise
func serializeVersionNodes(nodes []n8nsdk.Node, data *models.DataSourceVersion, resp *datasource.ReadResponse) (ok bool) {
	nodesJSON, err := json.Marshal(nodes)
	//: Check for nodes serialization error.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error serializing workflow nodes",
			fmt.Sprintf("Could not serialize workflow nodes to JSON: %s", err.Error()),
		)
		//: Return failure.
		return false
	}
	data.NodesJSON = types.StringValue(string(nodesJSON))
	//: Return success.
	return true
}

// mapVersionOptionalFields maps optional version fields to the data model.
//
// Params:
//   - name: optional workflow name pointer
//   - authors: optional authors pointer
//   - createdAt: optional created-at time pointer
//   - data: the data source model to update
func mapVersionOptionalFields(name, authors *string, createdAt *time.Time, data *models.DataSourceVersion) {
	//: Map name field.
	if name != nil {
		data.Name = types.StringValue(*name)
	} else {
		//: Set null when name is not available.
		data.Name = types.StringNull()
	}

	//: Map authors field.
	if authors != nil {
		data.Authors = types.StringValue(*authors)
	} else {
		//: Set null when authors is not available.
		data.Authors = types.StringNull()
	}

	//: Map created at field.
	if createdAt != nil {
		data.CreatedAt = types.StringValue(createdAt.Format("2006-01-02T15:04:05Z07:00"))
	} else {
		//: Set null when created at is not available.
		data.CreatedAt = types.StringNull()
	}
}
