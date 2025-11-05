package datasources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	providertypes "github.com/kodflow/n8n/src/internal/provider/types"
)

// Ensure VariableDataSource implements required interfaces.
var _ datasource.DataSource = &VariableDataSource{}
var _ datasource.DataSourceWithConfigure = &VariableDataSource{}

// VariableDataSource defines the data source implementation for a single variable.
type VariableDataSource struct {
	client *providertypes.N8nClient
}

// VariableDataSourceModel describes the data source data model.
type VariableDataSourceModel struct {
	ID        types.String `tfsdk:"id"`
	Key       types.String `tfsdk:"key"`
	Value     types.String `tfsdk:"value"`
	Type      types.String `tfsdk:"type"`
	ProjectID types.String `tfsdk:"project_id"`
}

// NewVariableDataSource creates a new VariableDataSource instance.
func NewVariableDataSource() datasource.DataSource {
	return &VariableDataSource{}
}

// Metadata returns the data source type name.
func (d *VariableDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_variable"
}

// Schema defines the schema for the data source.
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
func (d *VariableDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*providertypes.N8nClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *providertypes.N8nClient, got: %T", req.ProviderData),
		)
		return
	}

	d.client = client
}

// Read refreshes the Terraform state with the latest data.
func (d *VariableDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data VariableDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate that at least one identifier is provided
	if data.ID.IsNull() && data.Key.IsNull() {
		resp.Diagnostics.AddError(
			"Missing Required Attribute",
			"Either 'id' or 'key' must be specified",
		)
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
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Error listing variables",
			fmt.Sprintf("Could not list variables: %s\nHTTP Response: %v", err.Error(), httpResp),
		)
		return
	}

	// Filter by ID or key
	var found bool
	if variableList.Data != nil {
		for _, variable := range variableList.Data {
			matchByID := !data.ID.IsNull() && variable.Id != nil && *variable.Id == data.ID.ValueString()
			matchByKey := !data.Key.IsNull() && variable.Key == data.Key.ValueString()

			if matchByID || matchByKey {
				// Populate data model
				if variable.Id != nil {
					data.ID = types.StringValue(*variable.Id)
				}
				data.Key = types.StringValue(variable.Key)
				data.Value = types.StringValue(variable.Value)
				if variable.Type != nil {
					data.Type = types.StringPointerValue(variable.Type)
				}
				// Project is a nested object, extract ID if present
				if variable.Project != nil && variable.Project.Id != nil {
					data.ProjectID = types.StringPointerValue(variable.Project.Id)
				}
				found = true
				break
			}
		}
	}

	if !found {
		identifier := data.ID.ValueString()
		if identifier == "" {
			identifier = data.Key.ValueString()
		}
		resp.Diagnostics.AddError(
			"Variable Not Found",
			fmt.Sprintf("Could not find variable with identifier: %s", identifier),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
