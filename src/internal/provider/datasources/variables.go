package datasources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	providertypes "github.com/kodflow/n8n/src/internal/provider/types"
)

// Ensure VariablesDataSource implements required interfaces.
var _ datasource.DataSource = &VariablesDataSource{}
var _ datasource.DataSourceWithConfigure = &VariablesDataSource{}

// VariablesDataSource defines the data source implementation for listing variables.
type VariablesDataSource struct {
	client *providertypes.N8nClient
}

// VariablesDataSourceModel describes the data source data model.
type VariablesDataSourceModel struct {
	ProjectID types.String          `tfsdk:"project_id"`
	State     types.String          `tfsdk:"state"`
	Variables []VariableItemModel   `tfsdk:"variables"`
}

// VariableItemModel represents a single variable in the list.
type VariableItemModel struct {
	ID        types.String `tfsdk:"id"`
	Key       types.String `tfsdk:"key"`
	Value     types.String `tfsdk:"value"`
	Type      types.String `tfsdk:"type"`
	ProjectID types.String `tfsdk:"project_id"`
}

// NewVariablesDataSource creates a new VariablesDataSource instance.
func NewVariablesDataSource() datasource.DataSource {
 // Return result.
	return &VariablesDataSource{}
}

// Metadata returns the data source type name.
func (d *VariablesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_variables"
}

// Schema defines the schema for the data source.
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
func (d *VariablesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Check for nil value.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*providertypes.N8nClient)
	// Check condition.
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *providertypes.N8nClient, got: %T", req.ProviderData),
		)
		// Return result.
		return
	}

	d.client = client
}

// Read refreshes the Terraform state with the latest data.
func (d *VariablesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data VariablesDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	// Check condition.
	if resp.Diagnostics.HasError() {
		return
	}

	// Build API request with optional filters
	apiReq := d.client.APIClient.VariablesAPI.VariablesGet(ctx)

	// Check condition.
	if !data.ProjectID.IsNull() {
		apiReq = apiReq.ProjectId(data.ProjectID.ValueString())
	}
	// Check condition.
	if !data.State.IsNull() {
		apiReq = apiReq.State(data.State.ValueString())
	}

	variableList, httpResp, err := apiReq.Execute()
	// Check for non-nil value.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}
	// Check for error.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error listing variables",
			fmt.Sprintf("Could not list variables: %s\nHTTP Response: %v", err.Error(), httpResp),
		)
		// Return result.
		return
	}

	data.Variables = make([]VariableItemModel, 0)
	// Check for non-nil value.
	if variableList.Data != nil {
  // Iterate over items.
		for _, variable := range variableList.Data {
			item := VariableItemModel{}
			// Check for non-nil value.
			if variable.Id != nil {
				item.ID = types.StringValue(*variable.Id)
			}
			item.Key = types.StringValue(variable.Key)
			item.Value = types.StringValue(variable.Value)
			// Check for non-nil value.
			if variable.Type != nil {
				item.Type = types.StringPointerValue(variable.Type)
			}
			// Project is a nested object, extract ID if present
			if variable.Project != nil && variable.Project.Id != nil {
				item.ProjectID = types.StringPointerValue(variable.Project.Id)
			}
			data.Variables = append(data.Variables, item)
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
