package datasources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	providertypes "github.com/kodflow/n8n/src/internal/provider/types"
)

// DefaultListCapacity is the initial capacity for slices when listing items.
const DefaultListCapacity int = 10

// Ensure ExecutionsDataSource implements required interfaces.
var (
	_ datasource.DataSource              = &ExecutionsDataSource{}
	_ datasource.DataSourceWithConfigure = &ExecutionsDataSource{}
)

// ExecutionsDataSource is a Terraform datasource that provides read-only access to n8n executions.
// It enables querying and filtering workflow executions from the n8n API.
type ExecutionsDataSource struct {
	client *providertypes.N8nClient
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

// NewExecutionsDataSource creates a new ExecutionsDataSource instance.
func NewExecutionsDataSource() datasource.DataSource {
	// Return result.
	return &ExecutionsDataSource{}
}

// Metadata returns the data source type name.
func (d *ExecutionsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_executions"
}

// Schema defines the schema for the data source.
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
func (d *ExecutionsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *ExecutionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ExecutionsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	// Check condition.
	if resp.Diagnostics.HasError() {
		return
	}

	// Build API request with optional filters
	apiReq := d.client.APIClient.ExecutionAPI.ExecutionsGet(ctx)

	// Check condition.
	if !data.WorkflowID.IsNull() {
		apiReq = apiReq.WorkflowId(data.WorkflowID.ValueString())
	}
	// Check condition.
	if !data.ProjectID.IsNull() {
		apiReq = apiReq.ProjectId(data.ProjectID.ValueString())
	}
	// Check condition.
	if !data.Status.IsNull() {
		apiReq = apiReq.Status(data.Status.ValueString())
	}
	// Check condition.
	if !data.IncludeData.IsNull() {
		apiReq = apiReq.IncludeData(data.IncludeData.ValueBool())
	}

	executionList, httpResp, err := apiReq.Execute()
	// Check for non-nil value.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}
	// Check for error.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error listing executions",
			fmt.Sprintf("Could not list executions: %s\nHTTP Response: %v", err.Error(), httpResp),
		)
		// Return result.
		return
	}

	data.Executions = make([]ExecutionItemModel, 0, DefaultListCapacity)
	// Check for non-nil value.
	if executionList.Data != nil {
		// Iterate over items.
		for _, execution := range executionList.Data {
			item := ExecutionItemModel{}
			// Check for non-nil value.
			if execution.Id != nil {
				item.ID = types.StringValue(fmt.Sprintf("%v", *execution.Id))
			}
			// Check for non-nil value.
			if execution.WorkflowId != nil {
				item.WorkflowID = types.StringValue(fmt.Sprintf("%v", *execution.WorkflowId))
			}
			// Check for non-nil value.
			if execution.Finished != nil {
				item.Finished = types.BoolPointerValue(execution.Finished)
			}
			// Check for non-nil value.
			if execution.Mode != nil {
				item.Mode = types.StringPointerValue(execution.Mode)
			}
			// Check for non-nil value.
			if execution.StartedAt != nil {
				item.StartedAt = types.StringValue(execution.StartedAt.String())
			}
			// Check condition.
			if execution.StoppedAt.IsSet() {
				stoppedAt := execution.StoppedAt.Get()
				// Check for non-nil value.
				if stoppedAt != nil {
					item.StoppedAt = types.StringValue(stoppedAt.String())
				}
			}
			// Check for non-nil value.
			if execution.Status != nil {
				item.Status = types.StringPointerValue(execution.Status)
			}
			data.Executions = append(data.Executions, item)
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
