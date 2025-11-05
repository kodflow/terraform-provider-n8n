package resources

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	providertypes "github.com/kodflow/n8n/src/internal/provider/types"
)

// Ensure SourceControlPullResource implements required interfaces.
var (
	_ resource.Resource                = &SourceControlPullResource{}
	_ resource.ResourceWithConfigure   = &SourceControlPullResource{}
	_ resource.ResourceWithImportState = &SourceControlPullResource{}
)

// SourceControlPullResource defines the resource implementation for pulling from source control.
type SourceControlPullResource struct {
	client *providertypes.N8nClient
}

// SourceControlPullResourceModel describes the resource data model.
type SourceControlPullResourceModel struct {
	ID            types.String `tfsdk:"id"`
	Force         types.Bool   `tfsdk:"force"`
	VariablesJSON types.String `tfsdk:"variables_json"`
	ResultJSON    types.String `tfsdk:"result_json"`
}

// NewSourceControlPullResource creates a new SourceControlPullResource instance.
func NewSourceControlPullResource() resource.Resource {
	return &SourceControlPullResource{}
}

// Metadata returns the resource type name.
func (r *SourceControlPullResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_source_control_pull"
}

// Schema defines the schema for the resource.
func (r *SourceControlPullResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Pulls changes from the remote source control repository. Requires the Source Control feature to be licensed and connected to a repository. This resource triggers a pull operation and captures the import results.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Resource identifier (generated)",
				Computed:            true,
			},
			"force": schema.BoolAttribute{
				MarkdownDescription: "Whether to force pull, overwriting local changes",
				Optional:            true,
			},
			"variables_json": schema.StringAttribute{
				MarkdownDescription: "Variables to use during pull as JSON string (map[string]interface{})",
				Optional:            true,
			},
			"result_json": schema.StringAttribute{
				MarkdownDescription: "Import result as JSON string containing imported workflows, credentials, tags, and variables",
				Computed:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *SourceControlPullResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*providertypes.N8nClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *providertypes.N8nClient, got: %T", req.ProviderData),
		)
		return
	}

	r.client = client
}

// Create triggers the source control pull.
func (r *SourceControlPullResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SourceControlPullResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build pull request
	pullReq := n8nsdk.NewPull()

	// Set force flag if provided
	if !plan.Force.IsNull() && !plan.Force.IsUnknown() {
		force := plan.Force.ValueBool()
		pullReq.SetForce(force)
	}

	// Parse variables JSON if provided
	if !plan.VariablesJSON.IsNull() && !plan.VariablesJSON.IsUnknown() {
		var variables map[string]interface{}
		if err := json.Unmarshal([]byte(plan.VariablesJSON.ValueString()), &variables); err != nil {
			resp.Diagnostics.AddError(
				"Invalid variables JSON",
				fmt.Sprintf("Could not parse variables_json: %s", err.Error()),
			)
			return
		}
		pullReq.SetVariables(variables)
	}

	// Execute pull operation
	importResult, httpResp, err := r.client.APIClient.SourceControlAPI.SourceControlPullPost(ctx).Pull(*pullReq).Execute()
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Error pulling from source control",
			fmt.Sprintf("Could not pull from source control: %s\nHTTP Response: %v", err.Error(), httpResp),
		)
		return
	}

	// Generate a unique ID for this pull operation (using timestamp or similar)
	plan.ID = types.StringValue(fmt.Sprintf("pull-%d", httpResp.StatusCode))

	// Serialize import result to JSON
	if importResult != nil {
		resultJSON, err := json.Marshal(importResult)
		if err != nil {
			resp.Diagnostics.AddWarning(
				"Could not serialize import result",
				fmt.Sprintf("Import was successful but result could not be serialized: %s", err.Error()),
			)
		} else {
			plan.ResultJSON = types.StringValue(string(resultJSON))
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read refreshes the resource state. For pull operations, we just keep the current state.
func (r *SourceControlPullResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SourceControlPullResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Pull operations are one-time actions, so we just maintain the state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update is not supported for source control pull operations.
func (r *SourceControlPullResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update Not Supported",
		"Source control pull resources cannot be updated. To pull again, create a new resource.",
	)
}

// Delete removes the resource from state but doesn't perform any API operation.
func (r *SourceControlPullResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Pull operations cannot be undone, so we just remove from state
}

// ImportState imports the resource into Terraform state.
func (r *SourceControlPullResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
