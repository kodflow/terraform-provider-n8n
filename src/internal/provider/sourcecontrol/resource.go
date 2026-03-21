// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// Package sourcecontrol implements source control pull resource functionality.
package sourcecontrol

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/client"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/sourcecontrol/models"
)

// Ensure SourceControlPullResource implements required interfaces.
var (
	_ resource.Resource              = &SourceControlPullResource{}
	_ SourceControlPullResourceIface = &SourceControlPullResource{}
	_ resource.ResourceWithConfigure = &SourceControlPullResource{}
)

// SourceControlPullResourceIface defines the interface for SourceControlPullResource.
type SourceControlPullResourceIface interface {
	resource.Resource
	Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse)
	Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse)
	Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse)
	Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse)
	Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse)
	Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse)
	Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse)
}

// SourceControlPullResource manages source control pull operations on an n8n instance.
// Creating or updating this resource triggers a pull from the configured git repository,
// importing workflows, credentials, variables, and tags into n8n.
type SourceControlPullResource struct {
	// client is the N8n API client used for operations.
	client *client.N8nClient
}

// NewSourceControlPullResource creates a new SourceControlPullResource instance.
//
// Returns:
//   - *SourceControlPullResource: new SourceControlPullResource instance
func NewSourceControlPullResource() (r *SourceControlPullResource) {
	//: Return result.
	return &SourceControlPullResource{}
}

// NewSourceControlPullResourceWrapper creates a new SourceControlPullResource instance for Terraform.
// This wrapper function is used by the provider to maintain compatibility with the framework.
//
// Returns:
//   - resource.Resource: the wrapped SourceControlPullResource instance
func NewSourceControlPullResourceWrapper() (r resource.Resource) {
	//: Return the wrapped resource instance.
	return NewSourceControlPullResource()
}

// Metadata returns the resource type name.
//
// Params:
//   - ctx: context
//   - req: metadata request
//   - resp: metadata response
//
// Returns:
//   - none
func (r *SourceControlPullResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	//: Guard against cancelled context.
	if ctx.Err() != nil {
		//: Return early when context is cancelled.
		return
	}
	resp.TypeName = req.ProviderTypeName + "_source_control_pull"
}

// Schema defines the schema for the resource.
//
// Params:
//   - ctx: context
//   - req: schema request
//   - resp: schema response
//
// Returns:
//   - none
func (r *SourceControlPullResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	//: Guard against cancelled context.
	if ctx.Err() != nil {
		//: Return early when context is cancelled.
		return
	}
	resp.Schema = schema.Schema{
		MarkdownDescription: "Triggers a source control pull on an n8n instance, importing workflows, credentials, variables, and tags from the configured git repository.",

		Attributes: map[string]schema.Attribute{
			"force": schema.BoolAttribute{
				MarkdownDescription: "Force the pull even when no remote changes are detected",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"workflows_imported": schema.StringAttribute{
				MarkdownDescription: "JSON summary of imported workflow changes",
				Computed:            true,
			},
			"credentials_imported": schema.StringAttribute{
				MarkdownDescription: "JSON summary of imported credential changes",
				Computed:            true,
			},
			"tags_imported": schema.StringAttribute{
				MarkdownDescription: "JSON summary of imported tag changes",
				Computed:            true,
			},
			"variables_imported": schema.StringAttribute{
				MarkdownDescription: "JSON summary of imported variable changes",
				Computed:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
//
// Params:
//   - ctx: context
//   - req: configure request
//   - resp: configure response
//
// Returns:
//   - none
func (r *SourceControlPullResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	//: Guard against cancelled context.
	if ctx.Err() != nil {
		//: Return early when context is cancelled.
		return
	}
	//: Check for nil value.
	if req.ProviderData == nil {
		//: Return result.
		return
	}

	clientData, ok := req.ProviderData.(*client.N8nClient)
	//: Check condition.
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.N8nClient, got: %T", req.ProviderData),
		)
		//: Return result.
		return
	}

	r.client = clientData
}

// Create triggers a source control pull and sets the initial Terraform state.
//
// Params:
//   - ctx: context
//   - req: create request
//   - resp: create response
//
// Returns:
//   - none
func (r *SourceControlPullResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan *models.Resource

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	//: Check condition.
	if resp.Diagnostics.HasError() {
		//: Return with error.
		return
	}

	//: Execute create logic.
	if !r.executePullLogic(ctx, plan, &resp.Diagnostics) {
		//: Return with error.
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// executePullLogic contains the main logic for triggering a source control pull.
// This helper function is separated for testability.
//
// Params:
//   - ctx: Context for the request
//   - data: The resource data model
//   - diags: Diagnostics to append errors to
//
// Returns:
//   - bool: True if the pull succeeded, false otherwise
func (r *SourceControlPullResource) executePullLogic(ctx context.Context, data *models.Resource, diags interface {
	AddError(summary, detail string)
	HasError() bool
}) (ok bool) {
	pullReq := n8nsdk.NewPull()
	//: Set force flag from resource data.
	if !data.Force.IsNull() && !data.Force.IsUnknown() {
		force := data.Force.ValueBool()
		pullReq.SetForce(force)
	}

	result, httpResp, err := r.client.APIClient.SourceControlAPI.SourceControlPullPost(ctx).
		Pull(*pullReq).
		Execute()
	//: Close the HTTP response body to avoid resource leaks.
	if httpResp != nil && httpResp.Body != nil {
		defer func() {
			//: Silently discard close error on response body.
			if closeErr := httpResp.Body.Close(); closeErr != nil {
				_ = closeErr
			}
		}()
	}

	//: Check for error.
	if err != nil {
		diags.AddError(
			"Error triggering source control pull",
			fmt.Sprintf("Could not trigger source control pull: %s\nHTTP Response: %v", err.Error(), httpResp),
		)
		//: Return failure.
		return false
	}

	//: Map the import result to the resource model.
	if result != nil {
		mapImportResultToResource(result, data)
	} else {
		//: Clear all imported fields when no result is returned.
		data.WorkflowsImported = types.StringNull()
		data.CredentialsImported = types.StringNull()
		data.TagsImported = types.StringNull()
		data.VariablesImported = types.StringNull()
	}

	//: Return success.
	return true
}

// Read refreshes the Terraform state. Source control pull results are not re-readable,
// so this method preserves the existing state.
//
// Params:
//   - ctx: context
//   - req: read request
//   - resp: read response
//
// Returns:
//   - none
func (r *SourceControlPullResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *models.Resource

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	//: Check condition.
	if resp.Diagnostics.HasError() {
		//: Return with error.
		return
	}

	//: Preserve existing state since pull results are not re-readable.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update triggers a new source control pull and updates the Terraform state.
//
// Params:
//   - ctx: context
//   - req: update request
//   - resp: update response
//
// Returns:
//   - none
func (r *SourceControlPullResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan *models.Resource

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	//: Check condition.
	if resp.Diagnostics.HasError() {
		//: Return with error.
		return
	}

	//: Execute pull logic on update.
	if !r.executePullLogic(ctx, plan, &resp.Diagnostics) {
		//: Return with error.
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete removes the resource from state. No API call is needed since pull is not reversible.
//
// Params:
//   - ctx: context
//   - req: delete request
//   - resp: delete response
//
// Returns:
//   - none
func (r *SourceControlPullResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	//: Guard against cancelled context.
	if ctx.Err() != nil {
		//: Return early when context is cancelled.
		return
	}
	//: Removing a pull resource from state requires no API operation.
}
