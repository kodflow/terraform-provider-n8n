// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE.md in the project root for license information.

// Package variable implements environment variable management resources and data sources.
package variable

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
	"github.com/kodflow/n8n/src/internal/provider/shared/constants"
	"github.com/kodflow/n8n/src/internal/provider/variable/models"
)

// Ensure VariablesDataSource implements required interfaces.
var (
	_ datasource.DataSource              = &VariablesDataSource{}
	_ VariablesDataSourceInterface       = &VariablesDataSource{}
	_ datasource.DataSourceWithConfigure = &VariablesDataSource{}
)

// VariablesDataSourceInterface defines the interface for VariablesDataSource.
type VariablesDataSourceInterface interface {
	datasource.DataSource
	Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse)
	Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse)
	Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse)
	Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse)
}

// VariablesDataSource provides a Terraform datasource for read-only access to n8n variables.
// It enables users to fetch and list variables from their n8n instance through the n8n API.
type VariablesDataSource struct {
	// client is the N8n API client used for operations.
	client *client.N8nClient
}

// NewVariablesDataSource creates and returns a new VariablesDataSource instance.
//
// Returns:
//   - datasource.DataSource: a new VariablesDataSource instance
func NewVariablesDataSource() *VariablesDataSource {
	// Return result.
	return &VariablesDataSource{}
}

// NewVariablesDataSourceWrapper creates a new VariablesDataSource instance for Terraform.
// This wrapper function is used by the provider to maintain compatibility with the framework.
//
// Returns:
//   - datasource.DataSource: the wrapped VariablesDataSource instance
func NewVariablesDataSourceWrapper() datasource.DataSource {
	// Return the wrapped datasource instance.
	return NewVariablesDataSource()
}

// Metadata returns the data source type name.
//
// Params:
//   - ctx: context for the metadata request
//   - req: metadata request with provider type name
//   - resp: metadata response to populate
//
// Returns:
//   - (none)
func (d *VariablesDataSource) Metadata(_ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_variables"
}

// Schema defines the schema for the data source.
//
// Params:
//   - ctx: context for the schema request
//   - req: schema request with provider information
//   - resp: schema response to populate
//
// Returns:
//   - (none)
func (d *VariablesDataSource) Schema(_ctx context.Context, _req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches a list of n8n variables with optional filtering",
		Attributes:          d.schemaAttributes(),
	}
}

// schemaAttributes returns the attribute definitions for the variables data source schema.
//
// Returns:
//   - map[string]schema.Attribute: the data source attribute definitions
func (d *VariablesDataSource) schemaAttributes() map[string]schema.Attribute {
	// Return schema attributes.
	return map[string]schema.Attribute{
		"project_id": schema.StringAttribute{
			MarkdownDescription: "Filter variables by project ID",
			Optional:            true,
		},
		"state": schema.StringAttribute{
			MarkdownDescription: "Filter variables by state (e.g., 'empty')",
			Optional:            true,
		},
		"variables": schema.ListNestedAttribute{
			MarkdownDescription: "List of variables",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: d.variableItemAttributes()},
		},
	}
}

// variableItemAttributes returns the nested attribute definitions for individual variable items.
//
// Returns:
//   - map[string]schema.Attribute: the variable item attribute definitions
func (d *VariablesDataSource) variableItemAttributes() map[string]schema.Attribute {
	// Return schema attributes.
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "Variable identifier",
			Computed:            true,
		},
		"key": schema.StringAttribute{
			MarkdownDescription: "Variable key",
			Computed:            true,
		},
		"value": schema.StringAttribute{
			MarkdownDescription: "Variable value",
			Computed:            true,
			Sensitive:           true,
		},
		"type": schema.StringAttribute{
			MarkdownDescription: "Variable type",
			Computed:            true,
		},
		"project_id": schema.StringAttribute{
			MarkdownDescription: "Project ID the variable belongs to",
			Computed:            true,
		},
	}
}

// Configure adds the provider configured client to the data source.
//
// Params:
//   - ctx: context for the configure request
//   - req: configure request with provider data
//   - resp: configure response for diagnostics
//
// Returns:
//   - (none)
func (d *VariablesDataSource) Configure(_ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// If no provider data is available, exit early
	if req.ProviderData == nil {
		// Return result.
		return
	}

	clientData, ok := req.ProviderData.(*client.N8nClient)
	// If the provider data is not a valid N8nClient, report an error
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.N8nClient, got: %T", req.ProviderData),
		)
		// Return early on type assertion failure
		return
	}

	d.client = clientData
}

// Read refreshes the Terraform state with the latest data.
//
// Params:
//   - ctx: context for the read request
//   - req: read request with configuration
//   - resp: read response to populate with state
//
// Returns:
//   - (none)
func (d *VariablesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.DataSources

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	// If there are errors in loading the config, return early
	if resp.Diagnostics.HasError() {
		// Return with error.
		return
	}

	// Build API request with optional filters
	apiReq := d.buildAPIRequestWithFilters(ctx, &data)

	variableList, httpResp, err := apiReq.Execute()
	// Ensure the HTTP response body is properly closed to prevent resource leaks
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}
	// If the API request failed, report the error and return early
	if err != nil {
		resp.Diagnostics.AddError(
			"Error listing variables",
			fmt.Sprintf("Could not list variables: %s\nHTTP Response: %v", err.Error(), httpResp),
		)
		// Return with error.
		return
	}

	d.populateVariables(&data, variableList)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// buildAPIRequestWithFilters builds the API request with optional filters.
//
// Params:
//   - ctx: context for the request
//   - data: DataSources model with filter values
//
// Returns:
//   - n8nsdk.VariablesAPIVariablesGetRequest: API request with filters applied
func (d *VariablesDataSource) buildAPIRequestWithFilters(ctx context.Context, data *models.DataSources) n8nsdk.VariablesAPIVariablesGetRequest {
	apiReq := d.client.APIClient.VariablesAPI.VariablesGet(ctx)

	// If a project ID filter is specified, add it to the API request
	if !data.ProjectID.IsNull() {
		apiReq = apiReq.ProjectId(data.ProjectID.ValueString())
	}
	// If a state filter is specified, add it to the API request
	if !data.State.IsNull() {
		apiReq = apiReq.State(data.State.ValueString())
	}

	// Return result.
	return apiReq
}

// populateVariables populates the variables list from the API response.
//
// Params:
//   - data: DataSources model to populate
//   - variableList: API response with variable data
func (d *VariablesDataSource) populateVariables(data *models.DataSources, variableList *n8nsdk.VariableList) {
	data.Variables = make([]models.Item, 0, constants.DEFAULT_LIST_CAPACITY)
	// If the API returned variable data, process each variable
	if variableList.Data == nil {
		// Return result.
		return
	}

	// For each variable returned from the API, map it to the Terraform model
	for _, variable := range variableList.Data {
		item := d.mapVariableToItem(&variable)
		data.Variables = append(data.Variables, *item)
	}
}

// mapVariableToItem maps a variable from the API to a Item.
//
// Params:
//   - variable: variable from API response
//
// Returns:
//   - *models.Item: pointer to mapped variable item
func (d *VariablesDataSource) mapVariableToItem(variable *n8nsdk.Variable) *models.Item {
	item := &models.Item{}
	// If the variable has an ID, include it in the model
	if variable.Id != nil {
		item.ID = types.StringValue(*variable.Id)
	}
	item.Key = types.StringValue(variable.Key)
	item.Value = types.StringValue(variable.Value)
	// If the variable has a type, include it in the model
	if variable.Type != nil {
		item.Type = types.StringPointerValue(variable.Type)
	}
	// Extract project ID if the project object and its ID are present
	if variable.Project != nil && variable.Project.Id != nil {
		item.ProjectID = types.StringPointerValue(variable.Project.Id)
	}
	// Return pointer to avoid copying large struct.
	return item
}
