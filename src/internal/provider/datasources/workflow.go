package datasources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	providertypes "github.com/kodflow/n8n/src/internal/provider/types"
)

// Ensure WorkflowDataSource implements required interfaces
var _ datasource.DataSource = &WorkflowDataSource{}
var _ datasource.DataSourceWithConfigure = &WorkflowDataSource{}

// WorkflowDataSource defines the data source implementation for a single workflow.
type WorkflowDataSource struct {
	client *providertypes.N8nClient
}

// WorkflowDataSourceModel describes the data source data model.
type WorkflowDataSourceModel struct {
	ID     types.String `tfsdk:"id"`
	Name   types.String `tfsdk:"name"`
	Active types.Bool   `tfsdk:"active"`
}

// NewWorkflowDataSource creates a new WorkflowDataSource instance.
func NewWorkflowDataSource() datasource.DataSource {
	return &WorkflowDataSource{}
}

// Metadata returns the data source type name.
func (d *WorkflowDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workflow"
}

// Schema defines the schema for the data source.
func (d *WorkflowDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches a single n8n workflow by ID",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Workflow identifier",
				Required:            true,
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
	}
}

// Configure adds the provider configured client to the data source.
func (d *WorkflowDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *WorkflowDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data WorkflowDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	workflow, httpResp, err := d.client.APIClient.WorkflowAPI.WorkflowsIdGet(ctx, data.ID.ValueString()).Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading workflow",
			fmt.Sprintf("Could not read workflow ID %s: %s\nHTTP Response: %v", data.ID.ValueString(), err.Error(), httpResp),
		)
		return
	}

	data.Name = types.StringValue(workflow.Name)
	if workflow.Active != nil {
		data.Active = types.BoolPointerValue(workflow.Active)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
