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
// It fetches data table details from the n8n API using ID or name-based filtering.
type DataTableDataSource struct {
	// client is the N8n API client used for operations.
	client *client.N8nClient
}

// NewDataTableDataSource creates a new DataTableDataSource instance.
//
// Returns:
//   - *DataTableDataSource: A new DataTableDataSource instance
func NewDataTableDataSource() (dataTableDataSource *DataTableDataSource) {
	//: Return a new empty datasource instance.
	return &DataTableDataSource{}
}

// NewDataTableDataSourceWrapper creates a new DataTableDataSource instance for Terraform.
// This wrapper function is used by the provider to maintain compatibility with the framework.
//
// Returns:
//   - datasource.DataSource: the wrapped DataTableDataSource instance
func NewDataTableDataSourceWrapper() (dataSource datasource.DataSource) {
	//: Return the wrapped datasource instance.
	return NewDataTableDataSource()
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
func (d *DataTableDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	//: Return early if context is cancelled.
	if ctx.Err() != nil {
		//: Return result.
		return
	}
	resp.TypeName = req.ProviderTypeName + "_data_table"
}

// dataTableColumnSchema returns the nested column attribute schema for data table datasources.
//
// Returns:
//   - schema.NestedAttributeObject: The nested attribute object definition for columns
func dataTableColumnSchema() (nestedAttr schema.NestedAttributeObject) {
	//: Return the schema definition for column nested attributes.
	return schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Column identifier",
				Computed:            true,
			},
			"index": schema.Int64Attribute{
				MarkdownDescription: "Column position index",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Column name",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Column data type",
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
func (d *DataTableDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	//: Return early if context is cancelled.
	if ctx.Err() != nil {
		//: Return result.
		return
	}
	//: Apply the schema definition when the request is the default schema request.
	if req == (datasource.SchemaRequest{}) {
		resp.Schema = schema.Schema{
			MarkdownDescription: "Fetches a single n8n data table by ID or name.",
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
func (d *DataTableDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *DataTableDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	data := &models.DataSource{}

	resp.Diagnostics.Append(req.Config.Get(ctx, data)...)
	//: Return early if config parsing produced errors.
	if resp.Diagnostics.HasError() {
		//: Return with error after diagnostics are populated.
		return
	}

	//: Require at least one identifier to perform the lookup.
	if data.ID.IsNull() && data.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Missing Required Attribute",
			"Either 'id' or 'name' must be specified",
		)
		//: Return after adding the missing attribute error.
		return
	}

	//: Dispatch to ID-based or name-based read depending on which identifier is set.
	var ok bool
	//: Use direct GET by ID when the ID attribute is set.
	if !data.ID.IsNull() {
		ok = d.executeReadByID(ctx, data, resp)
		//: Fall back to listing and filtering by name when ID is not provided.
	} else {
		ok = d.executeReadByName(ctx, data, resp)
	}

	//: Return without setting state if the retrieval failed.
	if !ok {
		//: Return after read failure.
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

// executeReadByID retrieves a data table using the direct GET endpoint.
//
// Params:
//   - ctx: The request context
//   - data: The data source model
//   - resp: The read response
//
// Returns:
//   - bool: True if retrieval succeeded, false otherwise
func (d *DataTableDataSource) executeReadByID(ctx context.Context, data *models.DataSource, resp *datasource.ReadResponse) (ok bool) {
	dt, httpResp, err := d.client.APIClient.DataTableAPI.DataTablesIDGet(ctx, data.ID.ValueString()).Execute()
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
			"Error retrieving data table",
			fmt.Sprintf("Could not retrieve data table with ID %s: %s\nHTTP Response: %v", data.ID.ValueString(), err.Error(), httpResp),
		)
		//: Return false to signal retrieval failure.
		return false
	}

	//: Map the API response to the datasource model.
	mapDataTableToDataSourceModel(dt, data)

	//: Return true to signal successful retrieval.
	return true
}

// executeReadByName retrieves a data table by listing and filtering by name.
//
// Params:
//   - ctx: The request context
//   - data: The data source model
//   - resp: The read response
//
// Returns:
//   - bool: True if retrieval succeeded, false otherwise
func (d *DataTableDataSource) executeReadByName(ctx context.Context, data *models.DataSource, resp *datasource.ReadResponse) (ok bool) {
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

	//: Search the response list for the data table matching the requested name.
	return resolveNamedDataTable(dtList, data, resp)
}

// resolveNamedDataTable searches a list response for a named data table and maps it to the model.
//
// Params:
//   - dtList: The API list response containing data tables
//   - data: The data source model to populate when found
//   - resp: The read response for reporting not-found errors
//
// Returns:
//   - bool: True if the table was found and mapped, false otherwise
func resolveNamedDataTable(dtList *n8nsdk.DataTableList, data *models.DataSource, resp *datasource.ReadResponse) (ok bool) {
	//: Search only when the list response contains data.
	if dtList != nil && dtList.Data != nil {
		foundDT, found := findDataTableByName(dtList.Data, data.Name.ValueString())
		//: Map the found table to the model when a match exists.
		if found {
			mapDataTableToDataSourceModel(foundDT, data)
			//: Return true to signal successful name-based lookup.
			return true
		}
	}

	//: Report an error when no matching data table was found.
	resp.Diagnostics.AddError(
		"Data Table Not Found",
		fmt.Sprintf("Could not find data table with name: %s", data.Name.ValueString()),
	)
	//: Return false to signal the table was not found.
	return false
}
