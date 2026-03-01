// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// Package datatable implements data table management resources and data sources.
package datatable

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/datatable/models"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/client"
)

// Ensure DataTableDataSource implements required interfaces.
var (
	_ datasource.DataSource              = &DataTableDataSource{}
	_ DataTableDataSourceInterface       = &DataTableDataSource{}
	_ datasource.DataSourceWithConfigure = &DataTableDataSource{}
)

// DataTableDataSourceInterface defines the interface for DataTableDataSource.
type DataTableDataSourceInterface interface {
	datasource.DataSource
	Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse)
	Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse)
	Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse)
	Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse)
}

// DataTableDataSource is a Terraform datasource that provides read-only access to a single n8n data table.
type DataTableDataSource struct {
	client *client.N8nClient
}

// NewDataTableDataSource creates a new DataTableDataSource instance.
func NewDataTableDataSource() *DataTableDataSource {
	return &DataTableDataSource{}
}

// NewDataTableDataSourceWrapper creates a new DataTableDataSource instance for Terraform.
func NewDataTableDataSourceWrapper() datasource.DataSource {
	return NewDataTableDataSource()
}

// Metadata returns the data source type name.
func (d *DataTableDataSource) Metadata(_ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_data_table"
}

// Schema defines the schema for the data source.
func (d *DataTableDataSource) Schema(_ctx context.Context, _req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches a single n8n data table by ID or name. When using ID, the API's GET /data-tables/{id} endpoint is used. When using name, the LIST endpoint is used with client-side filtering.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Data table identifier. Either `id` or `name` must be specified.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Data table name. Either `id` or `name` must be specified.",
				Optional:            true,
				Computed:            true,
			},
			"project_id": schema.StringAttribute{
				MarkdownDescription: "ID of the project this table belongs to",
				Computed:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *DataTableDataSource) Configure(_ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	clientData, ok := req.ProviderData.(*client.N8nClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.N8nClient, got: %T", req.ProviderData),
		)
		return
	}

	d.client = clientData
}

// Read refreshes the Terraform state with the latest data.
func (d *DataTableDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	data := &models.DataSource{}

	resp.Diagnostics.Append(req.Config.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !d.validateIdentifier(data, resp) {
		return
	}

	var table *n8nsdk.DataTable
	if !data.ID.IsNull() {
		table = d.fetchDataTableByID(ctx, data, resp)
	} else {
		table = d.fetchDataTableByName(ctx, data, resp)
	}

	if table == nil {
		return
	}

	mapDataTableToDataSourceModel(table, data)
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

// validateIdentifier ensures at least one identifier is provided.
func (d *DataTableDataSource) validateIdentifier(data *models.DataSource, resp *datasource.ReadResponse) bool {
	if data.ID.IsNull() && data.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Missing Required Attribute",
			"Either 'id' or 'name' must be specified",
		)
		return false
	}
	return true
}

// fetchDataTableByID retrieves a data table using the GET endpoint.
func (d *DataTableDataSource) fetchDataTableByID(ctx context.Context, data *models.DataSource, resp *datasource.ReadResponse) *n8nsdk.DataTable {
	table, httpResp, err := d.client.APIClient.DataTableAPI.GetDataTable(ctx, data.ID.ValueString()).Execute()
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Error retrieving data table",
			fmt.Sprintf("Could not retrieve data table with ID %s: %s\nHTTP Response: %v", data.ID.ValueString(), err.Error(), httpResp),
		)
		return nil
	}
	return table
}

// fetchDataTableByName retrieves a data table by listing and filtering by name.
func (d *DataTableDataSource) fetchDataTableByName(ctx context.Context, data *models.DataSource, resp *datasource.ReadResponse) *n8nsdk.DataTable {
	tableList, httpResp, err := d.client.APIClient.DataTableAPI.ListDataTables(ctx).Execute()
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Error listing data tables",
			fmt.Sprintf("Could not list data tables: %s\nHTTP Response: %v", err.Error(), httpResp),
		)
		return nil
	}

	var table *n8nsdk.DataTable
	var found bool
	if tableList.Data != nil {
		table, found = findDataTableByName(tableList.Data, data.Name.ValueString())
	}

	if !found {
		resp.Diagnostics.AddError(
			"Data Table Not Found",
			fmt.Sprintf("Could not find data table with name: %s", data.Name.ValueString()),
		)
		return nil
	}

	return table
}
