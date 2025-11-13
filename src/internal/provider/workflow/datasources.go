// Package workflow implements workflow management resources and data sources.
package workflow

import (
	"context"
	"fmt"

	"github.com/kodflow/n8n/src/internal/provider/shared/constants"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
	"github.com/kodflow/n8n/src/internal/provider/workflow/models"
)

// Ensure WorkflowsDataSource implements required interfaces.
var (
	_ datasource.DataSource              = &WorkflowsDataSource{}
	_ WorkflowsDataSourceInterface       = &WorkflowsDataSource{}
	_ datasource.DataSourceWithConfigure = &WorkflowsDataSource{}
)

// WorkflowsDataSourceInterface defines the interface for WorkflowsDataSource.
type WorkflowsDataSourceInterface interface {
	datasource.DataSource
	Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse)
	Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse)
	Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse)
	Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse)
}

// WorkflowsDataSource provides a Terraform datasource for read-only access to n8n workflows.
// It enables users to fetch and list workflows from their n8n instance through the n8n API with optional filtering.
type WorkflowsDataSource struct {
	// client is the N8n API client used for operations.
	client *client.N8nClient
}

// NewWorkflowsDataSource creates a new WorkflowsDataSource instance.
//
// Returns:
//   - datasource.DataSource: A new WorkflowsDataSource instance
func NewWorkflowsDataSource() *WorkflowsDataSource {
	// Return result.
	return &WorkflowsDataSource{}
}

// NewWorkflowsDataSourceWrapper creates a new WorkflowsDataSource instance for Terraform.
// This wrapper function is used by the provider to maintain compatibility with the framework.
//
// Returns:
//   - datasource.DataSource: the wrapped WorkflowsDataSource instance
func NewWorkflowsDataSourceWrapper() datasource.DataSource {
	// Return the wrapped datasource instance.
	return NewWorkflowsDataSource()
}

// Metadata returns the data source type name.
//
// Params:
//   - ctx: context for the operation
//   - req: metadata request containing provider type name
//   - resp: metadata response to populate with type name
func (d *WorkflowsDataSource) Metadata(_ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workflows"
}

// Schema defines the schema for the data source.
//
// Params:
//   - ctx: context for the operation
//   - req: schema request
//   - resp: schema response to populate with schema definition
func (d *WorkflowsDataSource) Schema(_ctx context.Context, _req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
func (d *WorkflowsDataSource) Configure(_ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Check for nil value.
	if req.ProviderData == nil {
		// Return with error.
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
func (d *WorkflowsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.DataSources

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	// Check condition.
	if resp.Diagnostics.HasError() {
		// Return with error.
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
	data.Workflows = make([]models.Item, 0, constants.DEFAULT_LIST_CAPACITY)
	// Check for non-nil value.
	if workflowList.Data != nil {
		// Iterate over items.
		for _, workflow := range workflowList.Data {
			workflowModel := models.Item{
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
