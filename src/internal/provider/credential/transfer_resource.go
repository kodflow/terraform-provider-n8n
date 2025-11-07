package credential

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
)

// Ensure CredentialTransferResource implements required interfaces.
var (
	_ resource.Resource                   = &CredentialTransferResource{}
	_ CredentialTransferResourceInterface = &CredentialTransferResource{}
	_ resource.ResourceWithConfigure      = &CredentialTransferResource{}
	_ resource.ResourceWithImportState    = &CredentialTransferResource{}
)

// CredentialTransferResourceInterface defines the interface for CredentialTransferResource.
type CredentialTransferResourceInterface interface {
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

// CredentialTransferResource defines the resource implementation for transferring a credential to a project.
// Terraform resource that manages CRUD operations for credential transfers to n8n projects via the n8n API.
type CredentialTransferResource struct {
	// client is the N8n API client used for credential transfer operations.
	client *client.N8nClient
}

// NewCredentialTransferResource creates a new CredentialTransferResource instance.
//
// Returns:
//   - resource.Resource: the new CredentialTransferResource instance
func NewCredentialTransferResource() *CredentialTransferResource {
	// Return result.
	return &CredentialTransferResource{}
}

// NewCredentialTransferResourceWrapper creates a new CredentialTransferResource instance for Terraform.
// This wrapper function is used by the provider to maintain compatibility with the framework.
//
// Returns:
//   - resource.Resource: the wrapped CredentialTransferResource instance
func NewCredentialTransferResourceWrapper() resource.Resource {
	// Return the wrapped resource instance.
	return NewCredentialTransferResource()
}

// Metadata returns the resource type name.
//
// Params:
//   - ctx: context for cancellation and timeouts
//   - req: metadata request from Terraform
//   - resp: metadata response to populate
func (r *CredentialTransferResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_credential_transfer"
}

// Schema defines the schema for the resource.
//
// Params:
//   - ctx: context for cancellation and timeouts
//   - req: schema request from Terraform
//   - resp: schema response to populate
func (r *CredentialTransferResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Transfers a credential to another project. This is a one-time operation resource that triggers the transfer when created.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Resource identifier (generated)",
				Computed:            true,
			},
			"credential_id": schema.StringAttribute{
				MarkdownDescription: "ID of the credential to transfer",
				Required:            true,
			},
			"destination_project_id": schema.StringAttribute{
				MarkdownDescription: "ID of the destination project",
				Required:            true,
			},
			"transferred_at": schema.StringAttribute{
				MarkdownDescription: "Timestamp when the transfer occurred",
				Computed:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
//
// Params:
//   - ctx: context for cancellation and timeouts
//   - req: configure request from Terraform with provider data
//   - resp: configure response to populate with errors if any
func (r *CredentialTransferResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Check for nil value.
	if req.ProviderData == nil {
		// Return with error.
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

// Create triggers the credential transfer.
//
// Params:
//   - ctx: context for cancellation and timeouts
//   - req: create request from Terraform with plan data
//   - resp: create response to populate with state and errors if any
func (r *CredentialTransferResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan *CredentialTransferResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	// Check condition.
	if resp.Diagnostics.HasError() {
		// Return with error.
		return
	}

	// Build transfer request
	transferReq := n8nsdk.NewCredentialsIdTransferPutRequest(plan.DestinationProjectID.ValueString())

	// Execute transfer
	httpResp, err := r.client.APIClient.CredentialAPI.CredentialsIdTransferPut(ctx, plan.CredentialID.ValueString()).
		CredentialsIdTransferPutRequest(*transferReq).Execute()
	// Check for non-nil value.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}
	// Check for error.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error transferring credential",
			fmt.Sprintf("Could not transfer credential %s to project %s: %s\nHTTP Response: %v",
				plan.CredentialID.ValueString(), plan.DestinationProjectID.ValueString(), err.Error(), httpResp),
		)
		// Return result.
		return
	}

	// Set computed fields
	plan.ID = types.StringValue(fmt.Sprintf("%s-to-%s", plan.CredentialID.ValueString(), plan.DestinationProjectID.ValueString()))
	plan.TransferredAt = types.StringValue(fmt.Sprintf("transfer-response-%d", httpResp.StatusCode))

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Read refreshes the resource state. For transfer operations, we just keep the current state.
//
// Params:
//   - ctx: context for cancellation and timeouts
//   - req: read request from Terraform with state data
//   - resp: read response to populate with state and errors if any
func (r *CredentialTransferResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *CredentialTransferResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	// Check condition.
	if resp.Diagnostics.HasError() {
		// Return with error.
		return
	}

	// Transfer operations are one-time actions, so we just maintain the state
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// Update is not supported for transfer operations.
//
// Params:
//   - ctx: context for cancellation and timeouts
//   - req: update request from Terraform with plan and state data
//   - resp: update response to populate with errors
func (r *CredentialTransferResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update Not Supported",
		"Credential transfer resources cannot be updated. To transfer again, create a new resource.",
	)
}

// Delete removes the resource from state but doesn't perform any API operation.
//
// Params:
//   - ctx: context for cancellation and timeouts
//   - req: delete request from Terraform with state data
//   - resp: delete response to populate with errors if any
func (r *CredentialTransferResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Transfer operations cannot be undone, so we just remove from state
}

// ImportState imports the resource into Terraform state.
//
// Params:
//   - ctx: context for cancellation and timeouts
//   - req: import state request from Terraform with ID data
//   - resp: import state response to populate with state and errors if any
func (r *CredentialTransferResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
