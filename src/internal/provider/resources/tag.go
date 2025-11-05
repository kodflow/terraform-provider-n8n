package resources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	providertypes "github.com/kodflow/n8n/src/internal/provider/types"
)

// Ensure TagResource implements required interfaces
var (
	_ resource.Resource                = &TagResource{}
	_ resource.ResourceWithConfigure   = &TagResource{}
	_ resource.ResourceWithImportState = &TagResource{}
)

// TagResource defines the resource implementation for n8n tags.
type TagResource struct {
	client *providertypes.N8nClient
}

// TagResourceModel describes the resource data model.
type TagResourceModel struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	CreatedAt types.String `tfsdk:"created_at"`
	UpdatedAt types.String `tfsdk:"updated_at"`
}

// NewTagResource creates a new TagResource instance.
func NewTagResource() resource.Resource {
	return &TagResource{}
}

// Metadata returns the resource type name.
func (r *TagResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tag"
}

// Schema defines the schema for the resource.
func (r *TagResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
func (r *TagResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create creates the resource and sets the initial Terraform state.
func (r *TagResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan TagResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tagRequest := n8nsdk.Tag{
		Name: plan.Name.ValueString(),
	}

	tag, httpResp, err := r.client.APIClient.TagsAPI.TagsPost(ctx).
		Tag(tagRequest).
		Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating tag",
			fmt.Sprintf("Could not create tag: %s\nHTTP Response: %v", err.Error(), httpResp),
		)
		return
	}

	plan.ID = types.StringPointerValue(tag.Id)
	plan.Name = types.StringValue(tag.Name)
	if tag.CreatedAt != nil {
		plan.CreatedAt = types.StringValue(tag.CreatedAt.Format("2006-01-02T15:04:05Z07:00"))
	}
	if tag.UpdatedAt != nil {
		plan.UpdatedAt = types.StringValue(tag.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *TagResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state TagResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tag, httpResp, err := r.client.APIClient.TagsAPI.TagsIdGet(ctx, state.ID.ValueString()).Execute()
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading tag",
			fmt.Sprintf("Could not read tag ID %s: %s\nHTTP Response: %v", state.ID.ValueString(), err.Error(), httpResp),
		)
		return
	}

	state.Name = types.StringValue(tag.Name)
	if tag.CreatedAt != nil {
		state.CreatedAt = types.StringValue(tag.CreatedAt.Format("2006-01-02T15:04:05Z07:00"))
	}
	if tag.UpdatedAt != nil {
		state.UpdatedAt = types.StringValue(tag.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *TagResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan TagResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tagRequest := n8nsdk.Tag{
		Name: plan.Name.ValueString(),
	}

	tag, httpResp, err := r.client.APIClient.TagsAPI.TagsIdPut(ctx, plan.ID.ValueString()).
		Tag(tagRequest).
		Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating tag",
			fmt.Sprintf("Could not update tag ID %s: %s\nHTTP Response: %v", plan.ID.ValueString(), err.Error(), httpResp),
		)
		return
	}

	plan.Name = types.StringValue(tag.Name)
	if tag.CreatedAt != nil {
		plan.CreatedAt = types.StringValue(tag.CreatedAt.Format("2006-01-02T15:04:05Z07:00"))
	}
	if tag.UpdatedAt != nil {
		plan.UpdatedAt = types.StringValue(tag.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *TagResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state TagResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, httpResp, err := r.client.APIClient.TagsAPI.TagsIdDelete(ctx, state.ID.ValueString()).Execute()
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting tag",
			fmt.Sprintf("Could not delete tag ID %s: %s\nHTTP Response: %v", state.ID.ValueString(), err.Error(), httpResp),
		)
		return
	}
}

// ImportState imports the resource into Terraform state.
func (r *TagResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
