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
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/datatable/models"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/client"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/constants"
)

// Ensure DataTablesDataSource implements required interfaces.
var (
	_ datasource.DataSource              = &DataTablesDataSource{}
	_ DataTablesDataSourceInterface      = &DataTablesDataSource{}
	_ datasource.DataSourceWithConfigure = &DataTablesDataSource{}
)

// DataTablesDataSourceInterface defines the interface for DataTablesDataSource.
type DataTablesDataSourceInterface interface {
	datasource.DataSource
	Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse)
	Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse)
	Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse)
	Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse)
}

// DataTablesDataSource is a Terraform datasource implementation for listing data tables.
type DataTablesDataSource struct {
	client *client.N8nClient
}

// NewDataTablesDataSource creates a new DataTablesDataSource instance.
func NewDataTablesDataSource() *DataTablesDataSource {
	return &DataTablesDataSource{}
}

// NewDataTablesDataSourceWrapper creates a new DataTablesDataSource instance for Terraform.
func NewDataTablesDataSourceWrapper() datasource.DataSource {
	return NewDataTablesDataSource()
}

// Metadata returns the data source type name.
func (d *DataTablesDataSource) Metadata(_ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_data_tables"
}

// Schema defines the schema for the data source.
func (d *DataTablesDataSource) Schema(_ctx context.Context, _req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches a list of all n8n data tables",

		Attributes: map[string]schema.Attribute{
			"data_tables": schema.ListNestedAttribute{
				MarkdownDescription: "List of data tables",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Data table identifier",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Data table name",
							Computed:            true,
						},
						"project_id": schema.StringAttribute{
							MarkdownDescription: "ID of the project this table belongs to",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *DataTablesDataSource) Configure(_ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *DataTablesDataSource) Read(ctx context.Context, _req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.DataSources

	tableList, httpResp, err := d.client.APIClient.DataTableAPI.ListDataTables(ctx).Execute()
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Error listing data tables",
			fmt.Sprintf("Could not list data tables: %s\nHTTP Response: %v", err.Error(), httpResp),
		)
		return
	}

	data.DataTables = make([]models.Item, 0, constants.DEFAULT_LIST_CAPACITY)
	if tableList.Data != nil {
		for _, table := range tableList.Data {
			item := models.Item{
				ID:        types.StringValue(table.Id),
				Name:      types.StringValue(table.Name),
				ProjectID: types.StringValue(table.ProjectId),
			}
			data.DataTables = append(data.DataTables, item)
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
