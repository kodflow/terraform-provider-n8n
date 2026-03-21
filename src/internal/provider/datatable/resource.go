// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// Package datatable implements data table management resources and data sources.
package datatable

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/datatable/models"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/client"
)

// Ensure DataTableResource implements required interfaces.
var (
	_ resource.Resource                = &DataTableResource{}
	_ DataTableResourceInterface       = &DataTableResource{}
	_ resource.ResourceWithConfigure   = &DataTableResource{}
	_ resource.ResourceWithImportState = &DataTableResource{}
)

// DataTableResourceInterface defines the interface for DataTableResource.
type DataTableResourceInterface interface {
	resource.Resource
	Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse)
	Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse)
	Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse)
	Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse)
	Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse)
	Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse)
	Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse)
	ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse)
}

// DataTableResource defines the resource implementation for n8n data tables.
// Terraform resource that manages CRUD operations for n8n data tables via the n8n API.
type DataTableResource struct {
	// client is the N8n API client used for operations.
	client *client.N8nClient
}

// NewDataTableResource creates a new DataTableResource instance.
//
// Returns:
//   - *DataTableResource: new DataTableResource instance
func NewDataTableResource() (dataTableResource *DataTableResource) {
	//: Return a new empty resource instance.
	return &DataTableResource{}
}

// NewDataTableResourceWrapper creates a new DataTableResource instance for Terraform.
// This wrapper function is used by the provider to maintain compatibility with the framework.
//
// Returns:
//   - resource.Resource: the wrapped DataTableResource instance
func NewDataTableResourceWrapper() (res resource.Resource) {
	//: Return the wrapped resource instance.
	return NewDataTableResource()
}

// Metadata returns the resource type name.
//
// Params:
//   - ctx: context
//   - req: metadata request
//   - resp: metadata response
//
// Returns:
//   - none
func (r *DataTableResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	//: Return early if context is cancelled.
	if ctx.Err() != nil {
		//: Return result.
		return
	}
	resp.TypeName = req.ProviderTypeName + "_data_table"
}

// dataTableResourceSchema returns the schema definition for the data table resource.
//
// Returns:
//   - schema.Schema: The schema definition for the data table resource
func dataTableResourceSchema() (s schema.Schema) {
	//: Return the schema definition for the data table resource.
	return schema.Schema{
		MarkdownDescription: "n8n data table resource. Note: columns cannot be changed after creation; to change columns, delete and recreate the resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Data table identifier",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Data table name",
				Required:            true,
			},
			"project_id": schema.StringAttribute{
				MarkdownDescription: "Project ID the data table belongs to",
				Optional:            true,
				Computed:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "Timestamp when the data table was created",
				Computed:            true,
			},
			"updated_at": schema.StringAttribute{
				MarkdownDescription: "Timestamp when the data table was last updated",
				Computed:            true,
			},
		},
		Blocks: map[string]schema.Block{
			"columns": schema.ListNestedBlock{
				MarkdownDescription: "Column definitions for the data table. Cannot be changed after creation.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "Column name",
							Required:            true,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "Column data type",
							Required:            true,
						},
					},
				},
			},
		},
	}
}

// Schema defines the schema for the resource.
//
// Params:
//   - ctx: context
//   - req: schema request
//   - resp: schema response
//
// Returns:
//   - none
func (r *DataTableResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	//: Return early if context is cancelled.
	if ctx.Err() != nil {
		//: Return result.
		return
	}
	//: Apply the schema definition when the request is the default schema request.
	if req == (resource.SchemaRequest{}) {
		//: Populate the response schema with the resource schema definition.
		resp.Schema = dataTableResourceSchema()
	}
}

// Configure adds the provider configured client to the resource.
//
// Params:
//   - ctx: context
//   - req: configure request
//   - resp: configure response
//
// Returns:
//   - none
func (r *DataTableResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.N8nClient, got: %T", req.ProviderData),
		)
		//: Return after adding type mismatch error to diagnostics.
		return
	}

	r.client = clientData
}

// Create creates the resource and sets the initial Terraform state.
//
// Params:
//   - ctx: context
//   - req: create request
//   - resp: create response
//
// Returns:
//   - none
func (r *DataTableResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan *models.Resource

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	//: Return early if plan parsing produced errors.
	if resp.Diagnostics.HasError() {
		//: Return with error after diagnostics are populated.
		return
	}

	//: Return early if the create API call failed.
	if !r.executeCreateLogic(ctx, plan, resp) {
		//: Return after create failure.
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// executeCreateLogic contains the main logic for creating a data table.
// This helper function is separated for testability.
//
// Params:
//   - ctx: Context for the request
//   - plan: The planned resource data
//   - resp: Create response
//
// Returns:
//   - bool: True if creation succeeded, false otherwise
func (r *DataTableResource) executeCreateLogic(ctx context.Context, plan *models.Resource, resp *resource.CreateResponse) (ok bool) {
	//: Build column request slice from the plan column definitions.
	columns := make([]n8nsdk.CreateDataTableColumnRequest, 0, len(plan.Columns))
	//: Append each column definition to the request slice.
	for _, col := range plan.Columns {
		columns = append(columns, n8nsdk.CreateDataTableColumnRequest{
			Name: col.Name.ValueString(),
			Type: col.Type.ValueString(),
		})
	}

	createReq := n8nsdk.CreateDataTableRequest{
		Name:    plan.Name.ValueString(),
		Columns: columns,
	}

	dt, httpResp, err := r.client.APIClient.DataTableAPI.DataTablesPost(ctx).
		CreateDataTableRequest(createReq).
		Execute()
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
			"Error creating data table",
			fmt.Sprintf("Could not create data table: %s\nHTTP Response: %v", err.Error(), httpResp),
		)
		//: Return false to signal creation failure.
		return false
	}

	//: Map the API response to the plan model.
	mapDataTableToResourceModel(dt, plan)

	//: Return true to signal successful creation.
	return true
}

// Read refreshes the Terraform state with the latest data.
//
// Params:
//   - ctx: context
//   - req: read request
//   - resp: read response
//
// Returns:
//   - none
func (r *DataTableResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *models.Resource

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	//: Return early if state parsing produced errors.
	if resp.Diagnostics.HasError() {
		//: Return with error after diagnostics are populated.
		return
	}

	//: Handle read failure, including removal when resource no longer exists.
	if !r.executeReadLogic(ctx, state, resp) {
		//: Remove the resource from state when it has been deleted externally.
		if state.ID.IsNull() {
			resp.State.RemoveResource(ctx)
		}
		//: Return after read failure or external deletion.
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// executeReadLogic contains the main logic for reading a data table.
// This helper function is separated for testability.
//
// Params:
//   - ctx: Context for the request
//   - state: The current resource state
//   - resp: Read response
//
// Returns:
//   - bool: True if read succeeded, false otherwise
func (r *DataTableResource) executeReadLogic(ctx context.Context, state *models.Resource, resp *resource.ReadResponse) (ok bool) {
	dt, httpResp, err := r.client.APIClient.DataTableAPI.DataTablesIDGet(ctx, state.ID.ValueString()).Execute()
	//: Close the HTTP response body to avoid resource leaks.
	if httpResp != nil && httpResp.Body != nil {
		defer func() {
			//: Silently discard close error on response body.
			if closeErr := httpResp.Body.Close(); closeErr != nil {
				resp.Diagnostics.AddWarning("Failed to close response body", closeErr.Error())
			}
		}()
	}

	//: Mark state as removed when the resource no longer exists on the API.
	if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
		//: Set ID to null so Read() removes the resource from state.
		state.ID = types.StringNull()
		//: Return false to indicate the resource is gone without an error.
		return false
	}

	//: Return failure when the API call produced an error.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading data table",
			fmt.Sprintf("Could not read data table ID %s: %s\nHTTP Response: %v", state.ID.ValueString(), err.Error(), httpResp),
		)
		//: Return false to signal read failure.
		return false
	}

	//: Map the API response to the state model.
	mapDataTableToResourceModel(dt, state)

	//: Return true to signal successful read.
	return true
}

// Update updates the resource and sets the updated Terraform state on success.
//
// Params:
//   - ctx: context
//   - req: update request
//   - resp: update response
//
// Returns:
//   - none
func (r *DataTableResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state *models.Resource

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	//: Return early if plan or state parsing produced errors.
	if resp.Diagnostics.HasError() {
		//: Return with error after diagnostics are populated.
		return
	}

	//: Return early if the update API call failed.
	if !r.executeUpdateLogic(ctx, plan, state, resp) {
		//: Return after update failure.
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// executeUpdateLogic contains the main logic for updating a data table.
// This helper function is separated for testability.
//
// Params:
//   - ctx: Context for the request
//   - plan: The planned resource data
//   - state: The current resource state
//   - resp: Update response
//
// Returns:
//   - bool: True if update succeeded, false otherwise
func (r *DataTableResource) executeUpdateLogic(ctx context.Context, plan, state *models.Resource, resp *resource.UpdateResponse) (ok bool) {
	//: Use the state ID since the plan ID may be Unknown for Computed attributes.
	tableID := state.ID.ValueString()

	//: Copy immutable fields from state since they cannot change on update.
	plan.ID = state.ID
	plan.Columns = state.Columns

	updateReq := n8nsdk.UpdateDataTableRequest{
		Name: plan.Name.ValueString(),
	}

	dt, httpResp, err := r.client.APIClient.DataTableAPI.DataTablesIDPatch(ctx, tableID).
		UpdateDataTableRequest(updateReq).
		Execute()
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
			"Error updating data table",
			fmt.Sprintf("Could not update data table ID %s: %s\nHTTP Response: %v", tableID, err.Error(), httpResp),
		)
		//: Return false to signal update failure.
		return false
	}

	//: Map the API response to the plan model.
	mapDataTableToResourceModel(dt, plan)

	//: Restore columns from state since the update response may not include them.
	plan.Columns = state.Columns

	//: Return true to signal successful update.
	return true
}

// Delete deletes the resource and removes the Terraform state on success.
//
// Params:
//   - ctx: context
//   - req: delete request
//   - resp: delete response
//
// Returns:
//   - none
func (r *DataTableResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *models.Resource

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	//: Return early if state parsing produced errors.
	if resp.Diagnostics.HasError() {
		//: Return with error after diagnostics are populated.
		return
	}

	//: Execute the delete API call.
	r.executeDeleteLogic(ctx, state, resp)
}

// executeDeleteLogic contains the main logic for deleting a data table.
// This helper function is separated for testability.
//
// Params:
//   - ctx: Context for the request
//   - state: The current resource state
//   - resp: Delete response
//
// Returns:
//   - bool: True if delete succeeded, false otherwise
func (r *DataTableResource) executeDeleteLogic(ctx context.Context, state *models.Resource, resp *resource.DeleteResponse) (ok bool) {
	httpResp, err := r.client.APIClient.DataTableAPI.DataTablesIDDelete(ctx, state.ID.ValueString()).Execute()
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
			"Error deleting data table",
			fmt.Sprintf("Could not delete data table ID %s: %s\nHTTP Response: %v", state.ID.ValueString(), err.Error(), httpResp),
		)
		//: Return false to signal deletion failure.
		return false
	}

	//: Return true to signal successful deletion.
	return true
}

// ImportState imports the resource into Terraform state.
//
// Params:
//   - ctx: context
//   - req: import state request
//   - resp: import state response
//
// Returns:
//   - none
func (r *DataTableResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
