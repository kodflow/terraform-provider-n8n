package workflow

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
)

// Ensure WorkflowDataSource implements required interfaces.
var (
	_ datasource.DataSource              = &WorkflowDataSource{}
	_ datasource.DataSourceWithConfigure = &WorkflowDataSource{}
)

// WorkflowDataSource provides a Terraform datasource for read-only access to individual n8n workflows.
// It enables users to fetch workflow details by ID from their n8n instance through the n8n API.
type WorkflowDataSource struct {
	client *client.N8nClient
}

// WorkflowDataSourceModel maps the Terraform schema attributes for a single workflow datasource.
// It represents workflow metadata including its identifier, name, and activation status.
type WorkflowDataSourceModel struct {
	ID     types.String `tfsdk:"id"`
	Name   types.String `tfsdk:"name"`
	Active types.Bool   `tfsdk:"active"`
}

// NewWorkflowDataSource creates and returns a new WorkflowDataSource instance.
//
// Returns:
//   - datasource.DataSource: a new WorkflowDataSource instance
func NewWorkflowDataSource() datasource.DataSource {
	// Return result.
	return &WorkflowDataSource{}
}

// Metadata returns the data source type name.
//
// Params:
//   - ctx: context for the operation
//   - req: metadata request from Terraform
//   - resp: metadata response to populate
//
// Returns:
func (d *WorkflowDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workflow"
}

// Schema defines the schema for the data source.
//
// Params:
//   - ctx: context for the operation
//   - req: schema request from Terraform
//   - resp: schema response to populate
//
// Returns:
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
//
// Params:
//   - ctx: context for the operation
//   - req: configure request from Terraform
//   - resp: configure response to populate
//
// Returns:
func (d *WorkflowDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
//   - ctx: context for the operation
//   - req: read request from Terraform
//   - resp: read response to populate
//
// Returns:
func (d *WorkflowDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data WorkflowDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	// Check condition.
	if resp.Diagnostics.HasError() {
		return
	}

	workflow, httpResp, err := d.client.APIClient.WorkflowAPI.WorkflowsIdGet(ctx, data.ID.ValueString()).Execute()
	// Check for non-nil value.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}

	// Check for error.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading workflow",
			fmt.Sprintf("Could not read workflow ID %s: %s\nHTTP Response: %v", data.ID.ValueString(), err.Error(), httpResp),
		)
		// Return result.
		return
	}

	data.Name = types.StringValue(workflow.Name)
	// Check for non-nil value.
	if workflow.Active != nil {
		data.Active = types.BoolPointerValue(workflow.Active)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
