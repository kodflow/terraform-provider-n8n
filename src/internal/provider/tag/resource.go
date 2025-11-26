// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// Package tag implements tag management resources and data sources.
package tag

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/client"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/tag/models"
)

// Ensure TagResource implements required interfaces.
var (
	_ resource.Resource                = &TagResource{}
	_ TagResourceInterface             = &TagResource{}
	_ resource.ResourceWithConfigure   = &TagResource{}
	_ resource.ResourceWithImportState = &TagResource{}
)

// TagResourceInterface defines the interface for TagResource.
type TagResourceInterface interface {
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

// TagResource defines the resource implementation for n8n tags.
// Terraform resource that manages CRUD operations for n8n tags via the n8n API.
type TagResource struct {
	// client is the N8n API client used for operations.
	client *client.N8nClient
}

// NewTagResource creates a new TagResource instance.
//
// Returns:
//   - resource.Resource: new TagResource instance
func NewTagResource() *TagResource {
	// Return result.
	return &TagResource{}
}

// NewTagResourceWrapper creates a new TagResource instance for Terraform.
// This wrapper function is used by the provider to maintain compatibility with the framework.
//
// Returns:
//   - resource.Resource: the wrapped TagResource instance
func NewTagResourceWrapper() resource.Resource {
	// Return the wrapped resource instance.
	return NewTagResource()
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
func (r *TagResource) Metadata(_ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tag"
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
func (r *TagResource) Schema(_ctx context.Context, _req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "n8n tag resource using generated SDK",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Tag identifier",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Tag name",
				Required:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "Timestamp when the tag was created",
				Computed:            true,
			},
			"updated_at": schema.StringAttribute{
				MarkdownDescription: "Timestamp when the tag was last updated",
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
func (r *TagResource) Configure(_ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create creates the resource and sets the initial Terraform state.
//
// Params:
//   - ctx: context
//   - req: create request
//   - resp: create response
//
// Returns:
//   - none
func (r *TagResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan *models.Resource

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	// Check condition.
	if resp.Diagnostics.HasError() {
		// Return with error.
		return
	}

	// Execute create logic
	if !r.executeCreateLogic(ctx, plan, resp) {
		// Return with error.
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// executeCreateLogic contains the main logic for creating a tag.
// This helper function is separated for testability.
//
// Params:
//   - ctx: Context for the request
//   - plan: The planned resource data
//   - resp: Create response
//
// Returns:
//   - bool: True if creation succeeded, false otherwise
func (r *TagResource) executeCreateLogic(ctx context.Context, plan *models.Resource, resp *resource.CreateResponse) bool {
	tagRequest := n8nsdk.Tag{
		Name: plan.Name.ValueString(),
	}

	tag, httpResp, err := r.client.APIClient.TagsAPI.TagsPost(ctx).
		Tag(tagRequest).
		Execute()
	// Check for non-nil value.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}

	// Check for error.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating tag",
			fmt.Sprintf("Could not create tag: %s\nHTTP Response: %v", err.Error(), httpResp),
		)
		// Return failure.
		return false
	}

	plan.ID = types.StringPointerValue(tag.Id)
	plan.Name = types.StringValue(tag.Name)
	// Check for non-nil value.
	if tag.CreatedAt != nil {
		plan.CreatedAt = types.StringValue(tag.CreatedAt.Format("2006-01-02T15:04:05Z07:00"))
	}
	// Check for non-nil value.
	if tag.UpdatedAt != nil {
		plan.UpdatedAt = types.StringValue(tag.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"))
	}

	// Return success.
	return true
}

// Read refreshes the Terraform state with the latest data.
//
// Params:
//   - ctx: context
//   - req: read request
//   - resp: read response
//
// Returns:
//   - none
func (r *TagResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *models.Resource

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	// Check condition.
	if resp.Diagnostics.HasError() {
		// Return with error.
		return
	}

	// Execute read logic
	if !r.executeReadLogic(ctx, state, resp) {
		// Return with error.
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// executeReadLogic contains the main logic for reading a tag.
// This helper function is separated for testability.
//
// Params:
//   - ctx: Context for the request
//   - state: The current resource state
//   - resp: Read response
//
// Returns:
//   - bool: True if read succeeded, false otherwise
func (r *TagResource) executeReadLogic(ctx context.Context, state *models.Resource, resp *resource.ReadResponse) bool {
	tag, httpResp, err := r.client.APIClient.TagsAPI.TagsIdGet(ctx, state.ID.ValueString()).Execute()
	// Check for non-nil value.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}

	// Check for error.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading tag",
			fmt.Sprintf("Could not read tag ID %s: %s\nHTTP Response: %v", state.ID.ValueString(), err.Error(), httpResp),
		)
		// Return failure.
		return false
	}

	state.Name = types.StringValue(tag.Name)
	// Check for non-nil value.
	if tag.CreatedAt != nil {
		state.CreatedAt = types.StringValue(tag.CreatedAt.Format("2006-01-02T15:04:05Z07:00"))
	}
	// Check for non-nil value.
	if tag.UpdatedAt != nil {
		state.UpdatedAt = types.StringValue(tag.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"))
	}

	// Return success.
	return true
}

// Update updates the resource and sets the updated Terraform state on success.
//
// Params:
//   - ctx: context
//   - req: update request
//   - resp: update response
//
// Returns:
//   - none
func (r *TagResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state *models.Resource

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	// Check condition.
	if resp.Diagnostics.HasError() {
		// Return with error.
		return
	}

	// Execute update logic
	if !r.executeUpdateLogic(ctx, plan, state, resp) {
		// Return with error.
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// executeUpdateLogic contains the main logic for updating a tag.
// This helper function is separated for testability.
//
// Params:
//   - ctx: Context for the request
//   - plan: The planned resource data
//   - state: The current resource state
//   - resp: Update response
//
// Returns:
//   - bool: True if update succeeded, false otherwise
func (r *TagResource) executeUpdateLogic(ctx context.Context, plan, state *models.Resource, resp *resource.UpdateResponse) bool {
	// Use state.ID for the tag ID since plan.ID may be Unknown
	// for Computed attributes.
	tagID := state.ID.ValueString()

	// Copy ID from state to plan for consistency.
	plan.ID = state.ID

	tagRequest := n8nsdk.Tag{
		Name: plan.Name.ValueString(),
	}

	tag, httpResp, err := r.client.APIClient.TagsAPI.TagsIdPut(ctx, tagID).
		Tag(tagRequest).
		Execute()
	// Check for non-nil value.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}

	// Check for error.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating tag",
			fmt.Sprintf("Could not update tag ID %s: %s\nHTTP Response: %v", tagID, err.Error(), httpResp),
		)
		// Return failure.
		return false
	}

	plan.Name = types.StringValue(tag.Name)
	// Copy created_at from API response or fall back to state value.
	// The PUT API may not return created_at.
	if tag.CreatedAt != nil {
		plan.CreatedAt = types.StringValue(tag.CreatedAt.Format("2006-01-02T15:04:05Z07:00"))
	} else {
		plan.CreatedAt = state.CreatedAt
	}
	// Copy updated_at from API response or fall back to state value.
	if tag.UpdatedAt != nil {
		plan.UpdatedAt = types.StringValue(tag.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"))
	} else {
		plan.UpdatedAt = state.UpdatedAt
	}

	// Return success.
	return true
}

// Delete deletes the resource and removes the Terraform state on success.
//
// Params:
//   - ctx: context
//   - req: delete request
//   - resp: delete response
//
// Returns:
//   - none
func (r *TagResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *models.Resource

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	// Check condition.
	if resp.Diagnostics.HasError() {
		// Return with error.
		return
	}

	// Execute delete logic
	r.executeDeleteLogic(ctx, state, resp)
}

// executeDeleteLogic contains the main logic for deleting a tag.
// This helper function is separated for testability.
//
// Params:
//   - ctx: Context for the request
//   - state: The current resource state
//   - resp: Delete response
//
// Returns:
//   - bool: True if delete succeeded, false otherwise
func (r *TagResource) executeDeleteLogic(ctx context.Context, state *models.Resource, resp *resource.DeleteResponse) bool {
	_, httpResp, err := r.client.APIClient.TagsAPI.TagsIdDelete(ctx, state.ID.ValueString()).Execute()
	// Check for non-nil value.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}

	// Check for error.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting tag",
			fmt.Sprintf("Could not delete tag ID %s: %s\nHTTP Response: %v", state.ID.ValueString(), err.Error(), httpResp),
		)
		// Return failure.
		return false
	}

	// Return success.
	return true
}

// ImportState imports the resource into Terraform state.
//
// Params:
//   - ctx: context
//   - req: import state request
//   - resp: import state response
//
// Returns:
//   - none
func (r *TagResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
