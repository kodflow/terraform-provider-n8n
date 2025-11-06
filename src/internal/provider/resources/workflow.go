package resources

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
	providertypes "github.com/kodflow/n8n/src/internal/provider/types"
)

// Ensure WorkflowResource implements required interfaces.
var (
	_ resource.Resource                = &WorkflowResource{}
	_ resource.ResourceWithConfigure   = &WorkflowResource{}
	_ resource.ResourceWithImportState = &WorkflowResource{}
)

// WorkflowResource defines the resource implementation for n8n workflows.
type WorkflowResource struct {
	client *providertypes.N8nClient
}

// WorkflowResourceModel describes the resource data model.
type WorkflowResourceModel struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Active          types.Bool   `tfsdk:"active"`
	Tags            types.List   `tfsdk:"tags"`
	NodesJSON       types.String `tfsdk:"nodes_json"`
	ConnectionsJSON types.String `tfsdk:"connections_json"`
	SettingsJSON    types.String `tfsdk:"settings_json"`
	CreatedAt       types.String `tfsdk:"created_at"`
	UpdatedAt       types.String `tfsdk:"updated_at"`
	VersionId       types.String `tfsdk:"version_id"`
	IsArchived      types.Bool   `tfsdk:"is_archived"`
	TriggerCount    types.Int64  `tfsdk:"trigger_count"`
	Meta            types.Map    `tfsdk:"meta"`
	PinData         types.Map    `tfsdk:"pin_data"`
}

// NewWorkflowResource creates a new WorkflowResource instance.
func NewWorkflowResource() resource.Resource {
 // Return result.
	return &WorkflowResource{}
}

// Metadata returns the resource type name.
func (r *WorkflowResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workflow"
}

// Schema defines the schema for the resource.
func (r *WorkflowResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "n8n workflow resource using generated SDK",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Workflow identifier",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Workflow name",
				Required:            true,
			},
			"active": schema.BoolAttribute{
				MarkdownDescription: "Whether the workflow is active",
				Optional:            true,
				Computed:            true,
			},
			"tags": schema.ListAttribute{
				MarkdownDescription: "List of tag IDs associated with this workflow",
				ElementType:         types.StringType,
				Optional:            true,
			},
			"nodes_json": schema.StringAttribute{
				MarkdownDescription: "Workflow nodes as JSON string. Must be valid JSON array of node objects.",
				Optional:            true,
				Computed:            true,
			},
			"connections_json": schema.StringAttribute{
				MarkdownDescription: "Workflow connections as JSON string. Must be valid JSON object mapping node connections.",
				Optional:            true,
				Computed:            true,
			},
			"settings_json": schema.StringAttribute{
				MarkdownDescription: "Workflow settings as JSON string. Must be valid JSON object.",
				Optional:            true,
				Computed:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "Timestamp when the workflow was created",
				Computed:            true,
			},
			"updated_at": schema.StringAttribute{
				MarkdownDescription: "Timestamp when the workflow was last updated",
				Computed:            true,
			},
			"version_id": schema.StringAttribute{
				MarkdownDescription: "Version identifier of the workflow",
				Computed:            true,
			},
			"is_archived": schema.BoolAttribute{
				MarkdownDescription: "Whether the workflow is archived",
				Computed:            true,
			},
			"trigger_count": schema.Int64Attribute{
				MarkdownDescription: "Number of triggers in the workflow",
				Computed:            true,
			},
			"meta": schema.MapAttribute{
				MarkdownDescription: "Workflow metadata",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"pin_data": schema.MapAttribute{
				MarkdownDescription: "Pinned test data for the workflow",
				ElementType:         types.StringType,
				Computed:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *WorkflowResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*providertypes.N8nClient)
	// Check condition.
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *providertypes.N8nClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		// Return result.
		return
	}

	r.client = client
}

// Create creates the resource and sets the initial Terraform state.
func (r *WorkflowResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan WorkflowResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	// Check condition.
	if resp.Diagnostics.HasError() {
		return
	}

	// Create workflow using SDK
	// Note: 'active' and 'tags' fields are read-only during creation.
	// Note: 'nodes', 'connections', and 'settings' are required by the API.

	// Parse nodes from JSON if provided, otherwise use empty array
	var nodes []n8nsdk.Node
	// Check condition.
	if !plan.NodesJSON.IsNull() && !plan.NodesJSON.IsUnknown() {
		// Check for error.
		if err := json.Unmarshal([]byte(plan.NodesJSON.ValueString()), &nodes); err != nil {
			resp.Diagnostics.AddError(
				"Invalid nodes JSON",
				fmt.Sprintf("Could not parse nodes_json: %s", err.Error()),
			)
			// Return result.
			return
		}
 // Handle alternative case.
	} else {
		nodes = []n8nsdk.Node{}
	}

	// Parse connections from JSON if provided, otherwise use empty object
	var connections map[string]interface{}
	// Check condition.
	if !plan.ConnectionsJSON.IsNull() && !plan.ConnectionsJSON.IsUnknown() {
		// Check for error.
		if err := json.Unmarshal([]byte(plan.ConnectionsJSON.ValueString()), &connections); err != nil {
			resp.Diagnostics.AddError(
				"Invalid connections JSON",
				fmt.Sprintf("Could not parse connections_json: %s", err.Error()),
			)
			// Return result.
			return
		}
 // Handle alternative case.
	} else {
		connections = map[string]interface{}{}
	}

	// Parse settings from JSON if provided, otherwise use empty object
	var settings n8nsdk.WorkflowSettings
	// Check condition.
	if !plan.SettingsJSON.IsNull() && !plan.SettingsJSON.IsUnknown() {
		// Check for error.
		if err := json.Unmarshal([]byte(plan.SettingsJSON.ValueString()), &settings); err != nil {
			resp.Diagnostics.AddError(
				"Invalid settings JSON",
				fmt.Sprintf("Could not parse settings_json: %s", err.Error()),
			)
			// Return result.
			return
		}
	}

	workflowRequest := n8nsdk.Workflow{
		Name:        plan.Name.ValueString(),
		Nodes:       nodes,
		Connections: connections,
		Settings:    settings,
	}

	// Call SDK API to create workflow
	workflow, httpResp, err := r.client.APIClient.WorkflowAPI.WorkflowsPost(ctx).
		Workflow(workflowRequest).
		Execute()
		// Check for non-nil value.
		if httpResp != nil && httpResp.Body != nil {
			defer httpResp.Body.Close()
		}

	// Check for error.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating workflow",
			fmt.Sprintf("Could not create workflow, unexpected error: %s\nHTTP Response: %v", err.Error(), httpResp),
		)
		// Return result.
		return
	}

	// If tags were provided, update the workflow with tags using the dedicated tags endpoint
	// (tags field is read-only in POST/PUT /workflows but can be set via PUT /workflows/{id}/tags)
	if !plan.Tags.IsNull() && !plan.Tags.IsUnknown() && workflow.Id != nil {
		var tagIDs []string
		resp.Diagnostics.Append(plan.Tags.ElementsAs(ctx, &tagIDs, false)...)
		// Check condition.
		if resp.Diagnostics.HasError() {
			return
		}

		tagIdsInner := make([]n8nsdk.TagIdsInner, len(tagIDs))
  // Iterate over items.
		for i, tagID := range tagIDs {
			tagIdsInner[i] = n8nsdk.TagIdsInner{Id: tagID}
		}

		tags, httpResp, err := r.client.APIClient.WorkflowAPI.WorkflowsIdTagsPut(ctx, *workflow.Id).
			TagIdsInner(tagIdsInner).
			Execute()
			// Check for non-nil value.
			if httpResp != nil && httpResp.Body != nil {
				defer httpResp.Body.Close()
			}

		// Check for error.
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating workflow with tags",
				fmt.Sprintf("Workflow created but could not set tags: %s\nHTTP Response: %v", err.Error(), httpResp),
			)
			// Return result.
			return
		}

		// Update workflow tags from response
		workflow.Tags = tags
	}

	// Map response to state
	plan.ID = types.StringPointerValue(workflow.Id)
	plan.Name = types.StringValue(workflow.Name)
	// Check for non-nil value.
	if workflow.Active != nil {
		plan.Active = types.BoolPointerValue(workflow.Active)
	}

	// Map tags from response
	if len(workflow.Tags) > 0 {
		tagIDs := make([]types.String, len(workflow.Tags))
  // Iterate over items.
		for i, tag := range workflow.Tags {
			// Check for non-nil value.
			if tag.Id != nil {
				tagIDs[i] = types.StringValue(*tag.Id)
			}
		}
		tagList, diags := types.ListValueFrom(ctx, types.StringType, tagIDs)
		resp.Diagnostics.Append(diags...)
		// Check condition.
		if !resp.Diagnostics.HasError() {
			plan.Tags = tagList
		}
	}

	// Map computed fields from response
	if workflow.CreatedAt != nil {
		plan.CreatedAt = types.StringValue(workflow.CreatedAt.Format("2006-01-02T15:04:05Z07:00"))
	}
	// Check for non-nil value.
	if workflow.UpdatedAt != nil {
		plan.UpdatedAt = types.StringValue(workflow.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"))
	}
	// Check for non-nil value.
	if workflow.VersionId != nil {
		plan.VersionId = types.StringPointerValue(workflow.VersionId)
	}
	// Check for non-nil value.
	if workflow.IsArchived != nil {
		plan.IsArchived = types.BoolPointerValue(workflow.IsArchived)
	}
	// Check for non-nil value.
	if workflow.TriggerCount != nil {
		plan.TriggerCount = types.Int64Value(int64(*workflow.TriggerCount))
	}
	// Check for non-nil value.
	if workflow.Meta != nil {
		metaMap, diags := types.MapValueFrom(ctx, types.StringType, workflow.Meta)
		resp.Diagnostics.Append(diags...)
		// Check condition.
		if !resp.Diagnostics.HasError() {
			plan.Meta = metaMap
		}
	}
	// Check for non-nil value.
	if workflow.PinData != nil {
		pinDataMap, diags := types.MapValueFrom(ctx, types.StringType, workflow.PinData)
		resp.Diagnostics.Append(diags...)
		// Check condition.
		if !resp.Diagnostics.HasError() {
			plan.PinData = pinDataMap
		}
	}

	// Serialize nodes, connections, and settings back to JSON
	if workflow.Nodes != nil {
		nodesJSON, err := json.Marshal(workflow.Nodes)
		// Check for nil value.
		if err == nil {
			plan.NodesJSON = types.StringValue(string(nodesJSON))
		}
	}
	// Check for non-nil value.
	if workflow.Connections != nil {
		connectionsJSON, err := json.Marshal(workflow.Connections)
		// Check for nil value.
		if err == nil {
			plan.ConnectionsJSON = types.StringValue(string(connectionsJSON))
		}
	}
	settingsJSON, err := json.Marshal(workflow.Settings)
	// Check for nil value.
	if err == nil {
		plan.SettingsJSON = types.StringValue(string(settingsJSON))
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *WorkflowResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state WorkflowResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	// Check condition.
	if resp.Diagnostics.HasError() {
		return
	}

	// Get workflow from SDK
	workflow, httpResp, err := r.client.APIClient.WorkflowAPI.WorkflowsIdGet(ctx, state.ID.ValueString()).Execute()
	// Check for non-nil value.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}

	// Check for error.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading workflow",
			fmt.Sprintf("Could not read workflow ID %s: %s\nHTTP Response: %v", state.ID.ValueString(), err.Error(), httpResp),
		)
		// Return result.
		return
	}

	// Update state
	state.Name = types.StringValue(workflow.Name)
	// Check for non-nil value.
	if workflow.Active != nil {
		state.Active = types.BoolPointerValue(workflow.Active)
	}

	// Map tags from response
	if len(workflow.Tags) > 0 {
		tagIDs := make([]types.String, len(workflow.Tags))
  // Iterate over items.
		for i, tag := range workflow.Tags {
			// Check for non-nil value.
			if tag.Id != nil {
				tagIDs[i] = types.StringValue(*tag.Id)
			}
		}
		tagList, diags := types.ListValueFrom(ctx, types.StringType, tagIDs)
		resp.Diagnostics.Append(diags...)
		// Check condition.
		if !resp.Diagnostics.HasError() {
			state.Tags = tagList
		}
 // Handle alternative case.
	} else {
		// Set to empty list if no tags
		emptyList, diags := types.ListValueFrom(ctx, types.StringType, []types.String{})
		resp.Diagnostics.Append(diags...)
		// Check condition.
		if !resp.Diagnostics.HasError() {
			state.Tags = emptyList
		}
	}

	// Map computed fields from response
	if workflow.CreatedAt != nil {
		state.CreatedAt = types.StringValue(workflow.CreatedAt.Format("2006-01-02T15:04:05Z07:00"))
	}
	// Check for non-nil value.
	if workflow.UpdatedAt != nil {
		state.UpdatedAt = types.StringValue(workflow.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"))
	}
	// Check for non-nil value.
	if workflow.VersionId != nil {
		state.VersionId = types.StringPointerValue(workflow.VersionId)
	}
	// Check for non-nil value.
	if workflow.IsArchived != nil {
		state.IsArchived = types.BoolPointerValue(workflow.IsArchived)
	}
	// Check for non-nil value.
	if workflow.TriggerCount != nil {
		state.TriggerCount = types.Int64Value(int64(*workflow.TriggerCount))
	}
	// Check for non-nil value.
	if workflow.Meta != nil {
		metaMap, diags := types.MapValueFrom(ctx, types.StringType, workflow.Meta)
		resp.Diagnostics.Append(diags...)
		// Check condition.
		if !resp.Diagnostics.HasError() {
			state.Meta = metaMap
		}
	}
	// Check for non-nil value.
	if workflow.PinData != nil {
		pinDataMap, diags := types.MapValueFrom(ctx, types.StringType, workflow.PinData)
		resp.Diagnostics.Append(diags...)
		// Check condition.
		if !resp.Diagnostics.HasError() {
			state.PinData = pinDataMap
		}
	}

	// Serialize nodes, connections, and settings back to JSON
	if workflow.Nodes != nil {
		nodesJSON, err := json.Marshal(workflow.Nodes)
		// Check for nil value.
		if err == nil {
			state.NodesJSON = types.StringValue(string(nodesJSON))
		}
	}
	// Check for non-nil value.
	if workflow.Connections != nil {
		connectionsJSON, err := json.Marshal(workflow.Connections)
		// Check for nil value.
		if err == nil {
			state.ConnectionsJSON = types.StringValue(string(connectionsJSON))
		}
	}
	settingsJSON, err := json.Marshal(workflow.Settings)
	// Check for nil value.
	if err == nil {
		state.SettingsJSON = types.StringValue(string(settingsJSON))
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *WorkflowResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan WorkflowResourceModel
	var state WorkflowResourceModel

	// Read Terraform plan and state data into the models
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	// Check condition.
	if resp.Diagnostics.HasError() {
		return
	}

	// Update workflow using SDK

	// Parse nodes from JSON if provided
	var nodes []n8nsdk.Node
	// Check condition.
	if !plan.NodesJSON.IsNull() && !plan.NodesJSON.IsUnknown() {
		// Check for error.
		if err := json.Unmarshal([]byte(plan.NodesJSON.ValueString()), &nodes); err != nil {
			resp.Diagnostics.AddError(
				"Invalid nodes JSON",
				fmt.Sprintf("Could not parse nodes_json: %s", err.Error()),
			)
			// Return result.
			return
		}
 // Handle alternative case.
	} else {
		nodes = []n8nsdk.Node{}
	}

	// Parse connections from JSON if provided
	var connections map[string]interface{}
	// Check condition.
	if !plan.ConnectionsJSON.IsNull() && !plan.ConnectionsJSON.IsUnknown() {
		// Check for error.
		if err := json.Unmarshal([]byte(plan.ConnectionsJSON.ValueString()), &connections); err != nil {
			resp.Diagnostics.AddError(
				"Invalid connections JSON",
				fmt.Sprintf("Could not parse connections_json: %s", err.Error()),
			)
			// Return result.
			return
		}
 // Handle alternative case.
	} else {
		connections = map[string]interface{}{}
	}

	// Parse settings from JSON if provided
	var settings n8nsdk.WorkflowSettings
	// Check condition.
	if !plan.SettingsJSON.IsNull() && !plan.SettingsJSON.IsUnknown() {
		// Check for error.
		if err := json.Unmarshal([]byte(plan.SettingsJSON.ValueString()), &settings); err != nil {
			resp.Diagnostics.AddError(
				"Invalid settings JSON",
				fmt.Sprintf("Could not parse settings_json: %s", err.Error()),
			)
			// Return result.
			return
		}
	}

	// Check if activation status changed - use dedicated endpoints if so
	var workflow *n8nsdk.Workflow
	var httpResp *http.Response
	var err error

	activeChanged := !plan.Active.IsNull() && !state.Active.IsNull() &&
		plan.Active.ValueBool() != state.Active.ValueBool()

	// Check condition.
	if activeChanged {
		// Use dedicated activate/deactivate endpoints
		if plan.Active.ValueBool() {
			// Activate workflow
			_, httpResp, err = r.client.APIClient.WorkflowAPI.WorkflowsIdActivatePost(ctx, plan.ID.ValueString()).Execute()
			// Check for non-nil value.
			if httpResp != nil && httpResp.Body != nil {
				defer httpResp.Body.Close()
			}
  // Handle alternative case.
		} else {
			// Deactivate workflow
			_, httpResp, err = r.client.APIClient.WorkflowAPI.WorkflowsIdDeactivatePost(ctx, plan.ID.ValueString()).Execute()
			// Check for non-nil value.
			if httpResp != nil && httpResp.Body != nil {
				defer httpResp.Body.Close()
			}
		}

		// Check for error.
		if err != nil {
			action := "activate"
			// Check condition.
			if !plan.Active.ValueBool() {
				action = "deactivate"
			}
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error changing workflow activation status to %s", action),
				fmt.Sprintf("Could not %s workflow ID %s: %s\nHTTP Response: %v", action, plan.ID.ValueString(), err.Error(), httpResp),
			)
			// Return result.
			return
		}
	}

	// Update workflow content (name, nodes, connections, settings)
	workflowRequest := n8nsdk.Workflow{
		Name:        plan.Name.ValueString(),
		Nodes:       nodes,
		Connections: connections,
		Settings:    settings,
	}

	// Call SDK API to update workflow content
	workflow, httpResp, err = r.client.APIClient.WorkflowAPI.WorkflowsIdPut(ctx, plan.ID.ValueString()).
		Workflow(workflowRequest).
		Execute()
		// Check for non-nil value.
		if httpResp != nil && httpResp.Body != nil {
			defer httpResp.Body.Close()
		}

	// Check for error.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating workflow",
			fmt.Sprintf("Could not update workflow ID %s: %s\nHTTP Response: %v", plan.ID.ValueString(), err.Error(), httpResp),
		)
		// Return result.
		return
	}

	// Update tags separately using the dedicated tags endpoint
	if !plan.Tags.IsNull() && !plan.Tags.IsUnknown() {
		var tagIDs []string
		resp.Diagnostics.Append(plan.Tags.ElementsAs(ctx, &tagIDs, false)...)
		// Check condition.
		if resp.Diagnostics.HasError() {
			return
		}

		tagIdsInner := make([]n8nsdk.TagIdsInner, len(tagIDs))
  // Iterate over items.
		for i, tagID := range tagIDs {
			tagIdsInner[i] = n8nsdk.TagIdsInner{Id: tagID}
		}

		tags, httpResp, err := r.client.APIClient.WorkflowAPI.WorkflowsIdTagsPut(ctx, plan.ID.ValueString()).
			TagIdsInner(tagIdsInner).
			Execute()
			// Check for non-nil value.
			if httpResp != nil && httpResp.Body != nil {
				defer httpResp.Body.Close()
			}

		// Check for error.
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating workflow tags",
				fmt.Sprintf("Could not update tags for workflow ID %s: %s\nHTTP Response: %v", plan.ID.ValueString(), err.Error(), httpResp),
			)
			// Return result.
			return
		}

		// Update workflow tags from response
		workflow.Tags = tags
	}

	// Update state with response
	plan.Name = types.StringValue(workflow.Name)
	// Check for non-nil value.
	if workflow.Active != nil {
		plan.Active = types.BoolPointerValue(workflow.Active)
	}

	// Map tags from response
	if len(workflow.Tags) > 0 {
		tagIDs := make([]types.String, len(workflow.Tags))
  // Iterate over items.
		for i, tag := range workflow.Tags {
			// Check for non-nil value.
			if tag.Id != nil {
				tagIDs[i] = types.StringValue(*tag.Id)
			}
		}
		tagList, diags := types.ListValueFrom(ctx, types.StringType, tagIDs)
		resp.Diagnostics.Append(diags...)
		// Check condition.
		if !resp.Diagnostics.HasError() {
			plan.Tags = tagList
		}
 // Handle alternative case.
	} else {
		// Set to empty list if no tags
		emptyList, diags := types.ListValueFrom(ctx, types.StringType, []types.String{})
		resp.Diagnostics.Append(diags...)
		// Check condition.
		if !resp.Diagnostics.HasError() {
			plan.Tags = emptyList
		}
	}

	// Map computed fields from response
	if workflow.CreatedAt != nil {
		plan.CreatedAt = types.StringValue(workflow.CreatedAt.Format("2006-01-02T15:04:05Z07:00"))
	}
	// Check for non-nil value.
	if workflow.UpdatedAt != nil {
		plan.UpdatedAt = types.StringValue(workflow.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"))
	}
	// Check for non-nil value.
	if workflow.VersionId != nil {
		plan.VersionId = types.StringPointerValue(workflow.VersionId)
	}
	// Check for non-nil value.
	if workflow.IsArchived != nil {
		plan.IsArchived = types.BoolPointerValue(workflow.IsArchived)
	}
	// Check for non-nil value.
	if workflow.TriggerCount != nil {
		plan.TriggerCount = types.Int64Value(int64(*workflow.TriggerCount))
	}
	// Check for non-nil value.
	if workflow.Meta != nil {
		metaMap, diags := types.MapValueFrom(ctx, types.StringType, workflow.Meta)
		resp.Diagnostics.Append(diags...)
		// Check condition.
		if !resp.Diagnostics.HasError() {
			plan.Meta = metaMap
		}
	}
	// Check for non-nil value.
	if workflow.PinData != nil {
		pinDataMap, diags := types.MapValueFrom(ctx, types.StringType, workflow.PinData)
		resp.Diagnostics.Append(diags...)
		// Check condition.
		if !resp.Diagnostics.HasError() {
			plan.PinData = pinDataMap
		}
	}

	// Serialize nodes, connections, and settings back to JSON
	if workflow.Nodes != nil {
		nodesJSON, err := json.Marshal(workflow.Nodes)
		// Check for nil value.
		if err == nil {
			plan.NodesJSON = types.StringValue(string(nodesJSON))
		}
	}
	// Check for non-nil value.
	if workflow.Connections != nil {
		connectionsJSON, err := json.Marshal(workflow.Connections)
		// Check for nil value.
		if err == nil {
			plan.ConnectionsJSON = types.StringValue(string(connectionsJSON))
		}
	}
	settingsJSON, err := json.Marshal(workflow.Settings)
	// Check for nil value.
	if err == nil {
		plan.SettingsJSON = types.StringValue(string(settingsJSON))
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *WorkflowResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state WorkflowResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	// Check condition.
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete workflow using SDK
	_, httpResp, err := r.client.APIClient.WorkflowAPI.WorkflowsIdDelete(ctx, state.ID.ValueString()).Execute()
	// Check for non-nil value.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}

	// Check for error.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting workflow",
			fmt.Sprintf("Could not delete workflow ID %s: %s\nHTTP Response: %v", state.ID.ValueString(), err.Error(), httpResp),
		)
		// Return result.
		return
	}
}

// ImportState imports the resource into Terraform state.
func (r *WorkflowResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
