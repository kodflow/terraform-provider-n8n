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
// It provides read-only access to all n8n data tables through the n8n API.
type DataTablesDataSource struct {
	// client is the N8n API client used for operations.
	client *client.N8nClient
}

// NewDataTablesDataSource creates a new DataTablesDataSource instance.
//
// Returns:
//   - *DataTablesDataSource: a new DataTablesDataSource instance
func NewDataTablesDataSource() (dataTablesDataSource *DataTablesDataSource) {
	//: Return a new empty datasource instance.
	return &DataTablesDataSource{}
}

// NewDataTablesDataSourceWrapper creates a new DataTablesDataSource instance for Terraform.
// This wrapper function is used by the provider to maintain compatibility with the framework.
//
// Returns:
//   - datasource.DataSource: the wrapped DataTablesDataSource instance
func NewDataTablesDataSourceWrapper() (dataSource datasource.DataSource) {
	//: Return the wrapped datasource instance.
	return NewDataTablesDataSource()
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
func (d *DataTablesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	//: Return early if context is cancelled.
	if ctx.Err() != nil {
		//: Return result.
		return
	}
	resp.TypeName = req.ProviderTypeName + "_data_tables"
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
func (d *DataTablesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	//: Return early if context is cancelled.
	if ctx.Err() != nil {
		//: Return result.
		return
	}
	//: Apply the schema definition when the request is the default schema request.
	if req == (datasource.SchemaRequest{}) {
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
								MarkdownDescription: "Project ID the data table belongs to",
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
							"columns": schema.ListNestedAttribute{
								MarkdownDescription: "Column definitions for the data table",
								Computed:            true,
								NestedObject:        dataTableColumnSchema(),
							},
						},
					},
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
func (d *DataTablesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *DataTablesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.DataSources

	//: Honour deferred reads when the Terraform client supports deferral.
	if req.ClientCapabilities.DeferralAllowed {
		//: Return early when the client requests deferred evaluation.
		return
	}

	//: Return early if listing data tables failed.
	if !d.executeListLogic(ctx, &data, resp) {
		//: Return after list failure.
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// executeListLogic fetches all data tables and populates the data model.
//
// Params:
//   - ctx: The request context
//   - data: The data sources model to populate
//   - resp: The read response
//
// Returns:
//   - bool: True if listing succeeded, false otherwise
func (d *DataTablesDataSource) executeListLogic(ctx context.Context, data *models.DataSources, resp *datasource.ReadResponse) (ok bool) {
	dtList, httpResp, err := d.client.APIClient.DataTableAPI.DataTablesGet(ctx).Execute()
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
			"Error listing data tables",
			fmt.Sprintf("Could not list data tables: %s\nHTTP Response: %v", err.Error(), httpResp),
		)
		//: Return false to signal listing failure.
		return false
	}

	data.DataTables = make([]models.Item, 0, constants.DefaultListCapacity)
	//: Populate the model only when the response contains data.
	if dtList != nil && dtList.Data != nil {
		//: Iterate over each data table and append it to the model.
		for _, dt := range dtList.Data {
			item := mapDataTableToItem(&dt)
			data.DataTables = append(data.DataTables, item)
		}
	}

	//: Return true to signal successful listing.
	return true
}
