package sourcecontrol

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
	"github.com/kodflow/n8n/src/internal/provider/sourcecontrol/models"
)

// Ensure SourceControlPullResource implements required interfaces.
var (
	_ resource.Resource                  = &SourceControlPullResource{}
	_ SourceControlPullResourceInterface = &SourceControlPullResource{}
	_ resource.ResourceWithConfigure     = &SourceControlPullResource{}
	_ resource.ResourceWithImportState   = &SourceControlPullResource{}
)

// SourceControlPullResourceInterface defines the interface for SourceControlPullResource.
type SourceControlPullResourceInterface interface {
	resource.Resource
	Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse)
	Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse)
	Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse)
	Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse)
	Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse)
	Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse)
	Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse)
	ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse)
}

// SourceControlPullResource defines the resource implementation for pulling from source control.
// Terraform resource that manages CRUD operations for source control pull operations with the n8n API.
type SourceControlPullResource struct {
	// client is the N8n API client used for operations.
	client *client.N8nClient
}

// NewSourceControlPullResource creates a new SourceControlPullResource instance.
//
// Returns:
//   - resource.Resource: the created SourceControlPullResource instance
func NewSourceControlPullResource() *SourceControlPullResource {
	// Return result.
	return &SourceControlPullResource{}
}

// NewSourceControlPullResourceWrapper creates a new SourceControlPullResource instance for Terraform.
// This wrapper function is used by the provider to maintain compatibility with the framework.
//
// Returns:
//   - resource.Resource: the wrapped SourceControlPullResource instance
func NewSourceControlPullResourceWrapper() resource.Resource {
	// Return the wrapped resource instance.
	return NewSourceControlPullResource()
}

// Metadata returns the resource type name.
// Params:
//   - ctx: context.Context
//   - req: resource.MetadataRequest
//   - resp: *resource.MetadataResponse
func (r *SourceControlPullResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_source_control_pull"
}

// Schema defines the schema for the resource.
// Params:
//   - ctx: context.Context
//   - req: resource.SchemaRequest
//   - resp: *resource.SchemaResponse
func (r *SourceControlPullResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Pulls changes from the remote source control repository. Requires the Source Control feature to be licensed and connected to a repository. This resource triggers a pull operation and captures the import results.",
		Attributes:          r.schemaAttributes(),
	}
}

// schemaAttributes returns the attribute definitions for the source control pull resource schema.
//
// Returns:
//   - map[string]schema.Attribute: the resource attribute definitions
func (r *SourceControlPullResource) schemaAttributes() map[string]schema.Attribute {
	// Return schema attributes.
	return map[string]schema.Attribute{
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
	}
}

// Configure adds the provider configured client to the resource.
// Params:
//   - ctx: context.Context
//   - req: resource.ConfigureRequest
//   - resp: *resource.ConfigureResponse
func (r *SourceControlPullResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Check for nil value.
	if req.ProviderData == nil {
		// Return result.
		return
	}

	clientData, ok := req.ProviderData.(*client.N8nClient)
	// Check condition.
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.N8nClient, got: %T", req.ProviderData),
		)
		// Return result.
		return
	}

	r.client = clientData
}

// Create triggers the source control pull.
// Params:
//   - ctx: context.Context
//   - req: resource.CreateRequest
//   - resp: *resource.CreateResponse
func (r *SourceControlPullResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan *models.PullResource

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	// Check condition.
	if resp.Diagnostics.HasError() {
		// Return with error.
		return
	}

	// Build pull request with plan data
	pullReq := r.buildPullRequest(plan, resp)
	// Check for errors in diagnostics.
	if resp.Diagnostics.HasError() {
		// Return with error.
		return
	}

	// Execute pull operation and get result
	importResult, httpResp := r.executePullOperation(ctx, pullReq, resp)
	// Ensure response body is closed
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}
	// Check for errors in diagnostics.
	if resp.Diagnostics.HasError() {
		// Return with error.
		return
	}

	// Update plan with result data
	r.updatePlanWithResult(plan, importResult, httpResp, resp)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// buildPullRequest builds the pull request with plan data.
//
// Params:
//   - plan: models.PullResource with request data
//   - resp: response for diagnostics
//
// Returns:
//   - *n8nsdk.Pull: built pull request
func (r *SourceControlPullResource) buildPullRequest(plan *models.PullResource, resp *resource.CreateResponse) *n8nsdk.Pull {
	pullReq := n8nsdk.NewPull()

	// Set force flag if provided
	if !plan.Force.IsNull() && !plan.Force.IsUnknown() {
		force := plan.Force.ValueBool()
		pullReq.SetForce(force)
	}

	// Parse variables JSON if provided
	if !plan.VariablesJSON.IsNull() && !plan.VariablesJSON.IsUnknown() {
		variables := r.parseVariablesJSON(plan.VariablesJSON.ValueString(), resp)
		// Check for errors in diagnostics.
		if resp.Diagnostics.HasError() {
			// Return with error.
			return nil
		}
		pullReq.SetVariables(variables)
	}

	// Return result.
	return pullReq
}

// parseVariablesJSON parses the variables JSON string.
//
// Params:
//   - variablesJSON: JSON string to parse
//   - resp: response for diagnostics
//
// Returns:
//   - map[string]any: parsed variables, empty map on error
func (r *SourceControlPullResource) parseVariablesJSON(variablesJSON string, resp *resource.CreateResponse) map[string]any {
	var variables map[string]any
	// Check for error.
	if err := json.Unmarshal([]byte(variablesJSON), &variables); err != nil {
		resp.Diagnostics.AddError(
			"Invalid variables JSON",
			fmt.Sprintf("Could not parse variables_json: %s", err.Error()),
		)
		// Return empty map on error instead of nil.
		return make(map[string]any, 0)
	}
	// Return parsed variables.
	return variables
}

// executePullOperation executes the pull operation.
//
// Params:
//   - ctx: context for request cancellation
//   - pullReq: pull request to execute
//   - resp: response for diagnostics
//
// Returns:
//   - *n8nsdk.ImportResult: import result if successful
//   - *http.Response: HTTP response
func (r *SourceControlPullResource) executePullOperation(
	ctx context.Context,
	pullReq *n8nsdk.Pull,
	resp *resource.CreateResponse,
) (*n8nsdk.ImportResult, *http.Response) {
	importResult, httpResp, err := r.client.APIClient.SourceControlAPI.SourceControlPullPost(ctx).Pull(*pullReq).Execute()
	// Check for non-nil value.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}
	// Check for error.
	if err != nil {
		// Close response body before returning on error
		if httpResp != nil && httpResp.Body != nil {
			httpResp.Body.Close()
		}
		resp.Diagnostics.AddError(
			"Error pulling from source control",
			fmt.Sprintf("Could not pull from source control: %s\nHTTP Response: %v", err.Error(), httpResp),
		)
		// Return result.
		return nil, httpResp
	}
	// Return result.
	return importResult, httpResp
}

// updatePlanWithResult updates the plan with result data.
//
// Params:
//   - plan: models.PullResource to update
//   - importResult: import result from API
//   - httpResp: HTTP response
//   - resp: response for diagnostics
func (r *SourceControlPullResource) updatePlanWithResult(
	plan *models.PullResource,
	importResult *n8nsdk.ImportResult,
	httpResp *http.Response,
	resp *resource.CreateResponse,
) {
	// Generate a unique ID for this pull operation (using timestamp or similar)
	plan.ID = types.StringValue(fmt.Sprintf("pull-%d", httpResp.StatusCode))

	// Serialize import result to JSON
	if importResult == nil {
		// Return result.
		return
	}

	resultJSON, err := json.Marshal(importResult)
	// Check for error.
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Could not serialize import result",
			fmt.Sprintf("Import was successful but result could not be serialized: %s", err.Error()),
		)
		// Return with error.
		return
	}
	plan.ResultJSON = types.StringValue(string(resultJSON))
}

// Read refreshes the resource state. For pull operations, we just keep the current state.
// Params:
//   - ctx: context.Context
//   - req: resource.ReadRequest
//   - resp: *resource.ReadResponse
func (r *SourceControlPullResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *models.PullResource

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	// Check condition.
	if resp.Diagnostics.HasError() {
		// Return with error.
		return
	}

	// Pull operations are one-time actions, so we just maintain the state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update is not supported for source control pull operations.
// Params:
//   - ctx: context.Context
//   - req: resource.UpdateRequest
//   - resp: *resource.UpdateResponse
func (r *SourceControlPullResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update Not Supported",
		"Source control pull resources cannot be updated. To pull again, create a new resource.",
	)
}

// Delete removes the resource from state but doesn't perform any API operation.
// Params:
//   - ctx: context.Context
//   - req: resource.DeleteRequest
//   - resp: *resource.DeleteResponse
func (r *SourceControlPullResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Pull operations cannot be undone, so we just remove from state
	// No-op: intentionally empty
	_ = ctx
}

// ImportState imports the resource into Terraform state.
// Params:
//   - ctx: context.Context
//   - req: resource.ImportStateRequest
//   - resp: *resource.ImportStateResponse
func (r *SourceControlPullResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
