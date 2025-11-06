package datasources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	providertypes "github.com/kodflow/n8n/src/internal/provider/types"
)

// Ensure WorkflowsDataSource implements required interfaces.
var _ datasource.DataSource = &WorkflowsDataSource{}
var _ datasource.DataSourceWithConfigure = &WorkflowsDataSource{}

// WorkflowsDataSource defines the data source implementation for listing workflows.
type WorkflowsDataSource struct {
	client *providertypes.N8nClient
}

// WorkflowsDataSourceModel describes the data source data model.
type WorkflowsDataSourceModel struct {
	Workflows []WorkflowModel `tfsdk:"workflows"`
	Active    types.Bool      `tfsdk:"active"`
}

// WorkflowModel represents a single workflow in the list.
type WorkflowModel struct {
	ID     types.String `tfsdk:"id"`
	Name   types.String `tfsdk:"name"`
	Active types.Bool   `tfsdk:"active"`
}

// NewWorkflowsDataSource creates a new WorkflowsDataSource instance.
func NewWorkflowsDataSource() datasource.DataSource {
 // Return result.
	return &WorkflowsDataSource{}
}

// Metadata returns the data source type name.
func (d *WorkflowsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workflows"
}

// Schema defines the schema for the data source.
func (d *WorkflowsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches a list of n8n workflows",

		Attributes: map[string]schema.Attribute{
			"active": schema.BoolAttribute{
				MarkdownDescription: "Filter by active status",
				Optional:            true,
			},
			"workflows": schema.ListNestedAttribute{
				MarkdownDescription: "List of workflows",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Workflow identifier",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Workflow name",
							Computed:            true,
						},
						"active": schema.BoolAttribute{
							MarkdownDescription: "Whether the workflow is active",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *WorkflowsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *WorkflowsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data WorkflowsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	// Check condition.
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the request
	apiReq := d.client.APIClient.WorkflowAPI.WorkflowsGet(ctx)

	// Apply filters if specified
	if !data.Active.IsNull() {
		apiReq = apiReq.Active(data.Active.ValueBool())
	}

	workflowList, httpResp, err := apiReq.Execute()
	// Check for non-nil value.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}

	// Check for error.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading workflows",
			fmt.Sprintf("Could not read workflows: %s\nHTTP Response: %v", err.Error(), httpResp),
		)
		// Return result.
		return
	}

	// Map response to state
	data.Workflows = make([]WorkflowModel, 0)
	// Check for non-nil value.
	if workflowList.Data != nil {
  // Iterate over items.
		for _, workflow := range workflowList.Data {
			workflowModel := WorkflowModel{
				ID:   types.StringPointerValue(workflow.Id),
				Name: types.StringValue(workflow.Name),
			}
			// Check for non-nil value.
			if workflow.Active != nil {
				workflowModel.Active = types.BoolPointerValue(workflow.Active)
			}
			data.Workflows = append(data.Workflows, workflowModel)
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
