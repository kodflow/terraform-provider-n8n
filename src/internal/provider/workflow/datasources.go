package workflow

import (
	"github.com/kodflow/n8n/src/internal/provider/shared/constants"
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
)

// Ensure WorkflowsDataSource implements required interfaces.
var (
	_ datasource.DataSource              = &WorkflowsDataSource{}
	_ datasource.DataSourceWithConfigure = &WorkflowsDataSource{}
)

// WorkflowsDataSource provides a Terraform datasource for read-only access to n8n workflows.
// It enables users to fetch and list workflows from their n8n instance through the n8n API with optional filtering.
type WorkflowsDataSource struct {
	client *client.N8nClient
}

// NewWorkflowsDataSource creates a new WorkflowsDataSource instance.
//
// Returns:
//   - datasource.DataSource: A new WorkflowsDataSource instance
func NewWorkflowsDataSource() datasource.DataSource {
	// Return result.
	return &WorkflowsDataSource{}
}

// WorkflowsDataSourceModel maps the Terraform schema attributes for the workflows datasource.
// It represents the complete set of workflows data returned by the n8n API with optional active status filtering.
type WorkflowsDataSourceModel struct {
	Workflows []WorkflowItemModel `tfsdk:"workflows"`
	Active    types.Bool          `tfsdk:"active"`
}

// WorkflowItemModel maps individual workflow attributes within the Terraform schema.
// Each item represents a single workflow with its identifier, name, and activation status.
type WorkflowItemModel struct {
	ID     types.String `tfsdk:"id"`
	Name   types.String `tfsdk:"name"`
	Active types.Bool   `tfsdk:"active"`
}

// Metadata returns the data source type name.
//
// Params:
//   - ctx: context for the operation
//   - req: metadata request containing provider type name
//   - resp: metadata response to populate with type name
//
// Returns:
func (d *WorkflowsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workflows"
}

// Schema defines the schema for the data source.
//
// Params:
//   - ctx: context for the operation
//   - req: schema request
//   - resp: schema response to populate with schema definition
//
// Returns:
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
//
// Params:
//   - ctx: context for the operation
//   - req: configure request containing provider data
//   - resp: configure response for error handling
//
// Returns:
func (d *WorkflowsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
//   - req: read request containing configuration
//   - resp: read response to populate with state data
//
// Returns:
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
	data.Workflows = make([]WorkflowItemModel, 0, constants.DefaultListCapacity)
	// Check for non-nil value.
	if workflowList.Data != nil {
		// Iterate over items.
		for _, workflow := range workflowList.Data {
			workflowModel := WorkflowItemModel{
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
