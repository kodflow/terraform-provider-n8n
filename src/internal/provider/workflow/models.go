package workflow

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// WorkflowResourceModel describes the resource data model.
// Maps n8n workflow attributes to Terraform schema, storing workflow configuration, nodes, connections, and metadata.
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
	VersionID       types.String `tfsdk:"version_id"`
	IsArchived      types.Bool   `tfsdk:"is_archived"`
	TriggerCount    types.Int64  `tfsdk:"trigger_count"`
	Meta            types.Map    `tfsdk:"meta"`
	PinData         types.Map    `tfsdk:"pin_data"`
}
