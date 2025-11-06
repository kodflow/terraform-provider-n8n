// Package credential contains resources and datasources for credential management
package credential

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// CredentialTransferResourceModel describes the resource data model.
// Maps n8n credential transfer attributes to Terraform schema, storing transfer metadata and timestamps.
type CredentialTransferResourceModel struct {
	ID                   types.String `tfsdk:"id"`
	CredentialID         types.String `tfsdk:"credential_id"`
	DestinationProjectID types.String `tfsdk:"destination_project_id"`
	TransferredAt        types.String `tfsdk:"transferred_at"`
}
