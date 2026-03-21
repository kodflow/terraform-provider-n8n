// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// Package execution implements execution management resources.
package execution

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/execution/models"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/client"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/constants"
)

// Ensure ExecutionTagsResource implements required interfaces.
var (
	_ resource.Resource                = &ExecutionTagsResource{}
	_ ExecutionTagsResourceInterface   = &ExecutionTagsResource{}
	_ resource.ResourceWithConfigure   = &ExecutionTagsResource{}
	_ resource.ResourceWithImportState = &ExecutionTagsResource{}
)

// ExecutionTagsResourceInterface defines the interface for ExecutionTagsResource.
type ExecutionTagsResourceInterface interface {
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

// ExecutionTagsResource defines the resource implementation for execution annotation tags.
// Terraform resource that manages CRUD operations for execution tags via the n8n API.
type ExecutionTagsResource struct {
	// client is the N8n API client used for operations.
	client *client.N8nClient
}

// NewExecutionTagsResource creates a new ExecutionTagsResource instance.
//
// Returns:
//   - *ExecutionTagsResource: new ExecutionTagsResource instance
func NewExecutionTagsResource() (executionTagsResource *ExecutionTagsResource) {
	//: Return result.
	return &ExecutionTagsResource{}
}

// NewExecutionTagsResourceWrapper creates a new ExecutionTagsResource instance for Terraform.
// This wrapper function is used by the provider to maintain compatibility with the framework.
//
// Returns:
//   - resource.Resource: the wrapped ExecutionTagsResource instance
func NewExecutionTagsResourceWrapper() (r resource.Resource) {
	//: Return the wrapped resource instance.
	return NewExecutionTagsResource()
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
func (r *ExecutionTagsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	//: Guard against cancelled context.
	if ctx.Err() != nil {
		//: Return early when context is cancelled.
		return
	}
	resp.TypeName = req.ProviderTypeName + "_execution_tags"
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
func (r *ExecutionTagsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	//: Guard against cancelled context.
	if ctx.Err() != nil {
		//: Return early when context is cancelled.
		return
	}
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages annotation tags for an n8n execution.",

		Attributes: map[string]schema.Attribute{
			"execution_id": schema.StringAttribute{
				MarkdownDescription: "Execution identifier",
				Required:            true,
			},
			"tag_ids": schema.SetAttribute{
				MarkdownDescription: "Set of tag IDs to associate with the execution",
				Required:            true,
				ElementType:         types.StringType,
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
func (r *ExecutionTagsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// parseExecutionID converts a string execution ID to float32 for the API.
//
// Params:
//   - executionID: The execution ID string to parse
//
// Returns:
//   - float32: The parsed execution ID
//   - error: Error if parsing fails
func parseExecutionID(executionID string) (id float32, err error) {
	val, err := strconv.ParseFloat(executionID, constants.Float32BitSize)
	//: Check for error.
	if err != nil {
		//: Return error.
		return 0, fmt.Errorf("invalid execution ID %q: %w", executionID, err)
	}
	//: Return result.
	return float32(val), nil
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
func (r *ExecutionTagsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan *models.TagsResource

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	//: Check condition.
	if resp.Diagnostics.HasError() {
		//: Return with error.
		return
	}

	//: Execute create logic.
	if !r.executeCreateLogic(ctx, plan, resp) {
		//: Return with error.
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// executeCreateLogic contains the main logic for creating execution tags.
// This helper function is separated for testability.
//
// Params:
//   - ctx: Context for the request
//   - plan: The planned resource data
//   - resp: Create response
//
// Returns:
//   - bool: True if creation succeeded, false otherwise
func (r *ExecutionTagsResource) executeCreateLogic(ctx context.Context, plan *models.TagsResource, resp *resource.CreateResponse) (ok bool) {
	execID, err := parseExecutionID(plan.ExecutionID.ValueString())
	//: Check for error.
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Execution ID",
			fmt.Sprintf("Could not parse execution ID: %s", err.Error()),
		)
		//: Return failure.
		return false
	}

	tagIDs := extractTagIDs(ctx, plan.TagIDs, &resp.Diagnostics)
	//: Check condition.
	if resp.Diagnostics.HasError() {
		//: Return failure.
		return false
	}

	_, httpResp, err := r.client.APIClient.ExecutionAPI.ExecutionsIdTagsPut(ctx, execID).
		TagIds(tagIDs).
		Execute()
	//: Check for non-nil value.
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
		resp.Diagnostics.AddError(
			"Error setting execution tags",
			fmt.Sprintf("Could not set tags for execution %s: %s\nHTTP Response: %v",
				plan.ExecutionID.ValueString(), err.Error(), httpResp),
		)
		//: Return failure.
		return false
	}

	//: Return success.
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
func (r *ExecutionTagsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *models.TagsResource

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	//: Check condition.
	if resp.Diagnostics.HasError() {
		//: Return with error.
		return
	}

	//: Execute read logic.
	if !r.executeReadLogic(ctx, state, resp) {
		//: Return with error.
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// collectTagIDStrings builds a slice of tag ID strings from API response tags.
//
// Params:
//   - tags: The tags returned from the API
//
// Returns:
//   - []string: The tag ID strings
func collectTagIDStrings(tags []n8nsdk.Tag) (ids []string) {
	ids = make([]string, 0, len(tags))
	//: Iterate over items.
	for _, tag := range tags {
		//: Check for non-nil value.
		if tag.Id != nil {
			ids = append(ids, *tag.Id)
		}
	}
	//: Return result.
	return ids
}

// executeReadLogic contains the main logic for reading execution tags.
// This helper function is separated for testability.
//
// Params:
//   - ctx: Context for the request
//   - state: The current resource state
//   - resp: Read response
//
// Returns:
//   - bool: True if read succeeded, false otherwise
func (r *ExecutionTagsResource) executeReadLogic(ctx context.Context, state *models.TagsResource, resp *resource.ReadResponse) (ok bool) {
	execID, err := parseExecutionID(state.ExecutionID.ValueString())
	//: Check for error.
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Execution ID",
			fmt.Sprintf("Could not parse execution ID: %s", err.Error()),
		)
		//: Return failure.
		return false
	}

	tags, httpResp, err := r.client.APIClient.ExecutionAPI.ExecutionsIdTagsGet(ctx, execID).Execute()
	//: Check for non-nil value.
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
		resp.Diagnostics.AddError(
			"Error reading execution tags",
			fmt.Sprintf("Could not read tags for execution %s: %s\nHTTP Response: %v",
				state.ExecutionID.ValueString(), err.Error(), httpResp),
		)
		//: Return failure.
		return false
	}

	tagIDsValue, diags := types.SetValueFrom(ctx, types.StringType, collectTagIDStrings(tags))
	resp.Diagnostics.Append(diags...)
	//: Check condition.
	if resp.Diagnostics.HasError() {
		//: Return failure.
		return false
	}

	state.TagIDs = tagIDsValue

	//: Return success.
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
func (r *ExecutionTagsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state *models.TagsResource

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	//: Check condition.
	if resp.Diagnostics.HasError() {
		//: Return with error.
		return
	}

	//: Execute update logic.
	if !r.executeUpdateLogic(ctx, plan, state, resp) {
		//: Return with error.
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// executeUpdateLogic contains the main logic for updating execution tags.
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
func (r *ExecutionTagsResource) executeUpdateLogic(ctx context.Context, plan, state *models.TagsResource, resp *resource.UpdateResponse) (ok bool) {
	//: Copy execution_id from state since it cannot change.
	plan.ExecutionID = state.ExecutionID

	execID, err := parseExecutionID(state.ExecutionID.ValueString())
	//: Check for error.
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Execution ID",
			fmt.Sprintf("Could not parse execution ID: %s", err.Error()),
		)
		//: Return failure.
		return false
	}

	tagIDs := extractTagIDs(ctx, plan.TagIDs, &resp.Diagnostics)
	//: Check condition.
	if resp.Diagnostics.HasError() {
		//: Return failure.
		return false
	}

	_, httpResp, err := r.client.APIClient.ExecutionAPI.ExecutionsIdTagsPut(ctx, execID).
		TagIds(tagIDs).
		Execute()
	//: Check for non-nil value.
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
		resp.Diagnostics.AddError(
			"Error updating execution tags",
			fmt.Sprintf("Could not update tags for execution %s: %s\nHTTP Response: %v",
				state.ExecutionID.ValueString(), err.Error(), httpResp),
		)
		//: Return failure.
		return false
	}

	//: Return success.
	return true
}

// Delete deletes the resource by clearing all tags from the execution.
//
// Params:
//   - ctx: context
//   - req: delete request
//   - resp: delete response
//
// Returns:
//   - none
func (r *ExecutionTagsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *models.TagsResource

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	//: Check condition.
	if resp.Diagnostics.HasError() {
		//: Return with error.
		return
	}

	//: Execute delete logic.
	r.executeDeleteLogic(ctx, state, resp)
}

// executeDeleteLogic contains the main logic for deleting execution tags.
// This helper function is separated for testability.
//
// Params:
//   - ctx: Context for the request
//   - state: The current resource state
//   - resp: Delete response
//
// Returns:
//   - bool: True if delete succeeded, false otherwise
func (r *ExecutionTagsResource) executeDeleteLogic(ctx context.Context, state *models.TagsResource, resp *resource.DeleteResponse) (ok bool) {
	execID, err := parseExecutionID(state.ExecutionID.ValueString())
	//: Check for error.
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Execution ID",
			fmt.Sprintf("Could not parse execution ID: %s", err.Error()),
		)
		//: Return failure.
		return false
	}

	//: Clear all tags by setting an empty array.
	_, httpResp, err := r.client.APIClient.ExecutionAPI.ExecutionsIdTagsPut(ctx, execID).
		TagIds(nil).
		Execute()
	//: Check for non-nil value.
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
		resp.Diagnostics.AddError(
			"Error deleting execution tags",
			fmt.Sprintf("Could not clear tags for execution %s: %s\nHTTP Response: %v",
				state.ExecutionID.ValueString(), err.Error(), httpResp),
		)
		//: Return failure.
		return false
	}

	//: Return success.
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
func (r *ExecutionTagsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("execution_id"), req, resp)
}
