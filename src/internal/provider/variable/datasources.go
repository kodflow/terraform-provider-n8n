package variable

import (
	"context"
	"fmt"
	"github.com/kodflow/n8n/src/internal/provider/shared/constants"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
)

// Ensure VariablesDataSource implements required interfaces.
var (
	_ datasource.DataSource              = &VariablesDataSource{}
	_ datasource.DataSourceWithConfigure = &VariablesDataSource{}
)

// VariablesDataSource provides a Terraform datasource for read-only access to n8n variables.
// It enables users to fetch and list variables from their n8n instance through the n8n API.
type VariablesDataSource struct {
	client *client.N8nClient
}

// NewVariablesDataSource creates and returns a new VariablesDataSource instance.
//
// Returns:
//   - datasource.DataSource: a new VariablesDataSource instance
func NewVariablesDataSource() datasource.DataSource {
	// Return result.
	return &VariablesDataSource{}
}

// VariablesDataSourceModel maps the Terraform schema attributes for the variables datasource.
// It represents the complete set of variables data returned by the n8n API with optional filtering.
type VariablesDataSourceModel struct {
	ProjectID types.String        `tfsdk:"project_id"`
	State     types.String        `tfsdk:"state"`
	Variables []VariableItemModel `tfsdk:"variables"`
}

// VariableItemModel maps individual variable attributes within the Terraform schema.
// Each item represents a single variable with its ID, key, value, type, and associated project.
type VariableItemModel struct {
	ID        types.String `tfsdk:"id"`
	Key       types.String `tfsdk:"key"`
	Value     types.String `tfsdk:"value"`
	Type      types.String `tfsdk:"type"`
	ProjectID types.String `tfsdk:"project_id"`
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
func (d *VariablesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
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
func (d *VariablesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches a list of n8n variables with optional filtering",

		Attributes: map[string]schema.Attribute{
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
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
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
					},
				},
			},
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
func (d *VariablesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// If no provider data is available, exit early
	if req.ProviderData == nil {
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
	var data VariablesDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	// If there are errors in loading the config, return early
	if resp.Diagnostics.HasError() {
		return
	}

	// Build API request with optional filters
	apiReq := d.client.APIClient.VariablesAPI.VariablesGet(ctx)

	// If a project ID filter is specified, add it to the API request
	if !data.ProjectID.IsNull() {
		apiReq = apiReq.ProjectId(data.ProjectID.ValueString())
	}
	// If a state filter is specified, add it to the API request
	if !data.State.IsNull() {
		apiReq = apiReq.State(data.State.ValueString())
	}

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
		// Return early on API error
		return
	}

	data.Variables = make([]VariableItemModel, 0, constants.DefaultListCapacity)
	// If the API returned variable data, process each variable
	if variableList.Data != nil {
		// For each variable returned from the API, map it to the Terraform model
		for _, variable := range variableList.Data {
			item := VariableItemModel{}
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
			data.Variables = append(data.Variables, item)
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
