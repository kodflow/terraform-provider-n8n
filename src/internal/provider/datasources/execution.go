package datasources

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	providertypes "github.com/kodflow/n8n/src/internal/provider/types"
)

// Ensure ExecutionDataSource implements required interfaces.
var _ datasource.DataSource = &ExecutionDataSource{}
var _ datasource.DataSourceWithConfigure = &ExecutionDataSource{}

// ExecutionDataSource defines the data source implementation for a single execution.
type ExecutionDataSource struct {
	client *providertypes.N8nClient
}

// ExecutionDataSourceModel describes the data source data model.
type ExecutionDataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	WorkflowID   types.String `tfsdk:"workflow_id"`
	Finished     types.Bool   `tfsdk:"finished"`
	Mode         types.String `tfsdk:"mode"`
	StartedAt    types.String `tfsdk:"started_at"`
	StoppedAt    types.String `tfsdk:"stopped_at"`
	Status       types.String `tfsdk:"status"`
	IncludeData  types.Bool   `tfsdk:"include_data"`
}

// NewExecutionDataSource creates a new ExecutionDataSource instance.
func NewExecutionDataSource() datasource.DataSource {
	return &ExecutionDataSource{}
}

// Metadata returns the data source type name.
func (d *ExecutionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_execution"
}

// Schema defines the schema for the data source.
func (d *ExecutionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches a single n8n workflow execution by ID.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Execution identifier",
				Required:            true,
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
				MarkdownDescription: "Execution mode (e.g., 'manual', 'trigger', 'webhook')",
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
				MarkdownDescription: "Execution status (e.g., 'success', 'error', 'running')",
				Computed:            true,
			},
			"include_data": schema.BoolAttribute{
				MarkdownDescription: "Whether to include execution data in the response",
				Optional:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *ExecutionDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *ExecutionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ExecutionDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert ID string to float32 as required by the API
	executionID, err := strconv.ParseFloat(data.ID.ValueString(), 32)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Execution ID",
			fmt.Sprintf("Could not parse execution ID %s as number: %s", data.ID.ValueString(), err.Error()),
		)
		return
	}

	execution, httpResp, err := d.client.APIClient.ExecutionAPI.ExecutionsIdGet(ctx, float32(executionID)).Execute()
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Error retrieving execution",
			fmt.Sprintf("Could not retrieve execution with ID %s: %s\nHTTP Response: %v", data.ID.ValueString(), err.Error(), httpResp),
		)
		return
	}

	// Populate data model
	if execution.Id != nil {
		data.ID = types.StringValue(fmt.Sprintf("%v", *execution.Id))
	}
	if execution.WorkflowId != nil {
		data.WorkflowID = types.StringValue(fmt.Sprintf("%v", *execution.WorkflowId))
	}
	if execution.Finished != nil {
		data.Finished = types.BoolPointerValue(execution.Finished)
	}
	if execution.Mode != nil {
		data.Mode = types.StringPointerValue(execution.Mode)
	}
	if execution.StartedAt != nil {
		data.StartedAt = types.StringValue(execution.StartedAt.String())
	}
	if execution.StoppedAt.IsSet() {
		stoppedAt := execution.StoppedAt.Get()
		if stoppedAt != nil {
			data.StoppedAt = types.StringValue(stoppedAt.String())
		}
	}
	if execution.Status != nil {
		data.Status = types.StringPointerValue(execution.Status)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
