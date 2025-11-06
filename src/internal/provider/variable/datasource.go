package variable

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
)

// Ensure VariableDataSource implements required interfaces.
var (
	_ datasource.DataSource              = &VariableDataSource{}
	_ datasource.DataSourceWithConfigure = &VariableDataSource{}
)

// VariableDataSource is a Terraform datasource implementation for retrieving a single variable.
// It provides read-only access to n8n variable details through the n8n API.
type VariableDataSource struct {
	client *client.N8nClient
}

// VariableDataSourceModel maps Terraform schema attributes for variable data.
// It represents a single variable with all related attributes from the n8n API.
type VariableDataSourceModel struct {
	ID        types.String `tfsdk:"id"`
	Key       types.String `tfsdk:"key"`
	Value     types.String `tfsdk:"value"`
	Type      types.String `tfsdk:"type"`
	ProjectID types.String `tfsdk:"project_id"`
}

// NewVariableDataSource creates a new VariableDataSource instance.
//
// Params:
//   - (no parameters)
//
// Returns:
//   - datasource.DataSource: New VariableDataSource instance
func NewVariableDataSource() datasource.DataSource {
	// Return result.
	return &VariableDataSource{}
}

// Metadata returns the data source type name.
//
// Params:
//   - ctx: Context for the operation
//   - req: Metadata request from Terraform
//   - resp: Metadata response to send to Terraform
//
// Returns:
//   - None
func (d *VariableDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_variable"
}

// Schema defines the schema for the data source.
//
// Params:
//   - ctx: Context for the operation
//   - req: Schema request from Terraform
//   - resp: Schema response to send to Terraform
//
// Returns:
//   - None
func (d *VariableDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches a single n8n variable by ID or key. Since the n8n API doesn't provide a GET /variables/{id} endpoint, this datasource uses the LIST endpoint with client-side filtering.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Variable identifier. Either `id` or `key` must be specified.",
				Optional:            true,
				Computed:            true,
			},
			"key": schema.StringAttribute{
				MarkdownDescription: "Variable key. Either `id` or `key` must be specified.",
				Optional:            true,
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
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
//
// Params:
//   - ctx: Context for the operation
//   - req: Configure request from Terraform
//   - resp: Configure response to send to Terraform
//
// Returns:
//   - None
func (d *VariableDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Check for nil value.
	if req.ProviderData == nil {
		return
	}

	clientData, ok := req.ProviderData.(*client.N8nClient)
	// Check condition.
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.N8nClient, got: %T", req.ProviderData),
		)
		// Return result.
		return
	}

	d.client = clientData
}

// Read refreshes the Terraform state with the latest data.
//
// Params:
//   - ctx: Context for the operation
//   - req: Read request from Terraform containing the state
//   - resp: Read response to send to Terraform
//
// Returns:
//   - None
func (d *VariableDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	data := &VariableDataSourceModel{}

	// Initialize pointer to empty model
	

	resp.Diagnostics.Append(req.Config.Get(ctx, data)...)
	// Check if diagnostics have errors
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate that at least one identifier is provided
	if data.ID.IsNull() && data.Key.IsNull() {
		resp.Diagnostics.AddError(
			"Missing Required Attribute",
			"Either 'id' or 'key' must be specified",
		)
		// Return early
		return
	}

	// Build API request
	apiReq := d.client.APIClient.VariablesAPI.VariablesGet(ctx)

	// If project_id is provided, filter by it
	if !data.ProjectID.IsNull() {
		apiReq = apiReq.ProjectId(data.ProjectID.ValueString())
	}

	// List all variables and filter client-side (API limitation)
	variableList, httpResp, err := apiReq.Execute()
	// Close HTTP response body if present
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}
	// Handle API errors
	if err != nil {
		resp.Diagnostics.AddError(
			"Error listing variables",
			fmt.Sprintf("Could not list variables: %s\nHTTP Response: %v", err.Error(), httpResp),
		)
		return
	}

	// Find variable by ID or key
	var variable *n8nsdk.Variable
	var found bool
	// Check if variable data exists in response
	if variableList.Data != nil {
		variable, found = findVariableByIDOrKey(variableList.Data, data.ID, data.Key)
	}

	// Verify variable was found
	if !found {
		identifier := data.ID.ValueString()
		// Use key as fallback identifier
		if identifier == "" {
			identifier = data.Key.ValueString()
		}
		resp.Diagnostics.AddError(
			"Variable Not Found",
			fmt.Sprintf("Could not find variable with identifier: %s", identifier),
		)
		return
	}

	// Map variable to model
	mapVariableToDataSourceModel(variable, data)

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}
