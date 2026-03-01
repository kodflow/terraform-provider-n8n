// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// Package datatable implements data table management resources.
package datatable

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
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
type DataTableResource struct {
	client *client.N8nClient
}

// NewDataTableResource creates a new DataTableResource instance.
func NewDataTableResource() *DataTableResource {
	return &DataTableResource{}
}

// NewDataTableResourceWrapper creates a new DataTableResource instance for Terraform.
func NewDataTableResourceWrapper() resource.Resource {
	return NewDataTableResource()
}

// Metadata returns the resource type name.
func (r *DataTableResource) Metadata(_ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_data_table"
}

// Schema defines the schema for the resource.
func (r *DataTableResource) Schema(_ctx context.Context, _req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "n8n data table resource using generated SDK",

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
				MarkdownDescription: "Project ID to associate this variable with",
				Computed:            true,
			},
		},
		Blocks: map[string]schema.Block{
			"columns": schema.ListNestedBlock{
				MarkdownDescription: "Column definitions. Changing columns requires replacing the resource (API does not support in-place column updates).",
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
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *DataTableResource) Configure(_ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	clientData, ok := req.ProviderData.(*client.N8nClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.N8nClient, got: %T", req.ProviderData),
		)
		return
	}

	r.client = clientData
}

// Create creates the resource and sets the initial Terraform state.
func (r *DataTableResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan *models.Resource

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !r.executeCreateLogic(ctx, plan, resp) {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// executeCreateLogic contains the main logic for creating a data table.
func (r *DataTableResource) executeCreateLogic(ctx context.Context, plan *models.Resource, resp *resource.CreateResponse) bool {
	columns := make([]n8nsdk.CreateDataTableRequestColumnsInner, 0, len(plan.Columns))
	for _, c := range plan.Columns {
		columns = append(columns, n8nsdk.CreateDataTableRequestColumnsInner{
			Name: c.Name.ValueString(),
			Type: c.Type.ValueString(),
		})
	}

	createReq := n8nsdk.CreateDataTableRequest{
		Name:    plan.Name.ValueString(),
		Columns: columns,
	}

	table, httpResp, err := r.client.APIClient.DataTableAPI.CreateDataTable(ctx).
		CreateDataTableRequest(createReq).
		Execute()
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating data table",
			fmt.Sprintf("Could not create data table: %s\nHTTP Response: %v", err.Error(), httpResp),
		)
		return false
	}

	plan.ID = types.StringValue(table.Id)
	plan.Name = types.StringValue(table.Name)
	plan.ProjectID = types.StringValue(table.ProjectId)
	plan.Columns = mapColumnsToModel(table.Columns)

	return true
}

// Read refreshes the Terraform state with the latest data.
func (r *DataTableResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *models.Resource

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !r.executeReadLogic(ctx, state, resp) {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// executeReadLogic contains the main logic for reading a data table.
func (r *DataTableResource) executeReadLogic(ctx context.Context, state *models.Resource, resp *resource.ReadResponse) bool {
	table, httpResp, err := r.client.APIClient.DataTableAPI.GetDataTable(ctx, state.ID.ValueString()).Execute()
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading data table",
			fmt.Sprintf("Could not read data table ID %s: %s\nHTTP Response: %v", state.ID.ValueString(), err.Error(), httpResp),
		)
		return false
	}

	state.Name = types.StringValue(table.Name)
	state.ProjectID = types.StringValue(table.ProjectId)
	state.Columns = mapColumnsToModel(table.Columns)

	return true
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *DataTableResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state *models.Resource

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.Name.Equal(state.Name) && !r.executeUpdateLogic(ctx, plan, state, resp) {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// executeUpdateLogic contains the main logic for updating a data table.
// Note: n8n API UpdateDataTable only supports updating the name; column changes require replace.
func (r *DataTableResource) executeUpdateLogic(ctx context.Context, plan, state *models.Resource, resp *resource.UpdateResponse) bool {
	tableID := state.ID.ValueString()
	plan.ID = state.ID

	updateReq := n8nsdk.UpdateDataTableRequest{
		Name: plan.Name.ValueString(),
	}

	table, httpResp, err := r.client.APIClient.DataTableAPI.UpdateDataTable(ctx, tableID).
		UpdateDataTableRequest(updateReq).
		Execute()
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating data table",
			fmt.Sprintf("Could not update data table ID %s: %s\nHTTP Response: %v", tableID, err.Error(), httpResp),
		)
		return false
	}

	plan.Name = types.StringValue(table.Name)
	plan.ProjectID = types.StringValue(table.ProjectId)
	plan.Columns = mapColumnsToModel(table.Columns)

	return true
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *DataTableResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *models.Resource

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.executeDeleteLogic(ctx, state, resp)
}

// executeDeleteLogic contains the main logic for deleting a data table.
func (r *DataTableResource) executeDeleteLogic(ctx context.Context, state *models.Resource, resp *resource.DeleteResponse) bool {
	httpResp, err := r.client.APIClient.DataTableAPI.DeleteDataTable(ctx, state.ID.ValueString()).Execute()
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting data table",
			fmt.Sprintf("Could not delete data table ID %s: %s\nHTTP Response: %v", state.ID.ValueString(), err.Error(), httpResp),
		)
		return false
	}

	return true
}

// ImportState imports the resource into Terraform state.
func (r *DataTableResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
