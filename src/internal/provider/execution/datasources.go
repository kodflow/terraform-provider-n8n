package execution

import (
	"github.com/kodflow/n8n/src/internal/provider/shared/constants"
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
)


// Ensure ExecutionsDataSource implements required interfaces.
var (
	_ datasource.DataSource              = &ExecutionsDataSource{}
	_ datasource.DataSourceWithConfigure = &ExecutionsDataSource{}
)

// ExecutionsDataSource is a Terraform datasource that provides read-only access to n8n executions.
// It enables querying and filtering workflow executions from the n8n API.
type ExecutionsDataSource struct {
	client *client.N8nClient
}

// NewExecutionsDataSource creates a new ExecutionsDataSource instance.
//
// Params:
//   - none
//
// Returns:
//   - datasource.DataSource: new ExecutionsDataSource instance
func NewExecutionsDataSource() datasource.DataSource {
	// Return result.
	return &ExecutionsDataSource{}
}

// ExecutionsDataSourceModel maps the Terraform schema to the datasource response.
// It represents the filtered execution list with workflow and execution details from the n8n API.
type ExecutionsDataSourceModel struct {
	WorkflowID  types.String         `tfsdk:"workflow_id"`
	ProjectID   types.String         `tfsdk:"project_id"`
	Status      types.String         `tfsdk:"status"`
	IncludeData types.Bool           `tfsdk:"include_data"`
	Executions  []ExecutionItemModel `tfsdk:"executions"`
}

// ExecutionItemModel represents a single execution in the list returned from the n8n API.
// It contains the execution metadata including timestamps, status, and workflow reference.
type ExecutionItemModel struct {
	ID         types.String `tfsdk:"id"`
	WorkflowID types.String `tfsdk:"workflow_id"`
	Finished   types.Bool   `tfsdk:"finished"`
	Mode       types.String `tfsdk:"mode"`
	StartedAt  types.String `tfsdk:"started_at"`
	StoppedAt  types.String `tfsdk:"stopped_at"`
	Status     types.String `tfsdk:"status"`
}

// Metadata returns the data source type name.
//
// Params:
//   - ctx: context.Context for cancellation and timeout control
//   - req: datasource.MetadataRequest containing provider type name
//   - resp: datasource.MetadataResponse to populate with metadata
//
// Returns:
func (d *ExecutionsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_executions"
}

// Schema defines the schema for the data source.
//
// Params:
//   - ctx: context.Context for cancellation and timeout control
//   - req: datasource.SchemaRequest for schema definition
//   - resp: datasource.SchemaResponse to populate with schema
//
// Returns:
func (d *ExecutionsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches a list of n8n workflow executions with optional filtering",

		Attributes: map[string]schema.Attribute{
			"workflow_id": schema.StringAttribute{
				MarkdownDescription: "Filter executions by workflow ID",
				Optional:            true,
			},
			"project_id": schema.StringAttribute{
				MarkdownDescription: "Filter executions by project ID",
				Optional:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Filter executions by status (e.g., 'success', 'error', 'running')",
				Optional:            true,
			},
			"include_data": schema.BoolAttribute{
				MarkdownDescription: "Whether to include execution data in the response",
				Optional:            true,
			},
			"executions": schema.ListNestedAttribute{
				MarkdownDescription: "List of executions",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Execution identifier",
							Computed:            true,
						},
						"workflow_id": schema.StringAttribute{
							MarkdownDescription: "ID of the workflow that was executed",
							Computed:            true,
						},
						"finished": schema.BoolAttribute{
							MarkdownDescription: "Whether the execution finished",
							Computed:            true,
						},
						"mode": schema.StringAttribute{
							MarkdownDescription: "Execution mode",
							Computed:            true,
						},
						"started_at": schema.StringAttribute{
							MarkdownDescription: "Timestamp when the execution started",
							Computed:            true,
						},
						"stopped_at": schema.StringAttribute{
							MarkdownDescription: "Timestamp when the execution stopped",
							Computed:            true,
						},
						"status": schema.StringAttribute{
							MarkdownDescription: "Execution status",
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
//   - ctx: context.Context for cancellation and timeout control
//   - req: datasource.ConfigureRequest containing provider data
//   - resp: datasource.ConfigureResponse to populate with diagnostics
//
// Returns:
func (d *ExecutionsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
//   - ctx: context.Context for cancellation and timeout control
//   - req: datasource.ReadRequest containing configuration data
//   - resp: datasource.ReadResponse to populate with state and diagnostics
//
// Returns:
func (d *ExecutionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Initialize data model
	data := &ExecutionsDataSourceModel{}

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	// Check for errors in diagnostics
	if resp.Diagnostics.HasError() {
		return
	}

	// Build API request with optional filters
	apiReq := d.client.APIClient.ExecutionAPI.ExecutionsGet(ctx)
	// Filter by workflow ID if provided
	if !data.WorkflowID.IsNull() {
		apiReq = apiReq.WorkflowId(data.WorkflowID.ValueString())
	}
	// Filter by project ID if provided
	if !data.ProjectID.IsNull() {
		apiReq = apiReq.ProjectId(data.ProjectID.ValueString())
	}
	// Filter by status if provided
	if !data.Status.IsNull() {
		apiReq = apiReq.Status(data.Status.ValueString())
	}
	// Set include data flag if provided
	if !data.IncludeData.IsNull() {
		apiReq = apiReq.IncludeData(data.IncludeData.ValueBool())
	}

	executionList, httpResp, err := apiReq.Execute()
	// Close HTTP response body if present
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}
	// Check for API execution errors
	if err != nil {
		resp.Diagnostics.AddError(
			"Error listing executions",
			fmt.Sprintf("Could not list executions: %s\nHTTP Response: %v", err.Error(), httpResp),
		)
		return
	}

	data.Executions = make([]ExecutionItemModel, 0, constants.DefaultListCapacity)
	// Check if execution data is available and populate executions
	if executionList.Data != nil {
		// Process each execution and add to the executions list
		for _, execution := range executionList.Data {
			item := mapExecutionToItem(&execution)
			data.Executions = append(data.Executions, item)
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
