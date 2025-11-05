package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Resources returns all resources supported by the provider.
// Each resource factory function creates a new instance of the resource.
func Resources() []func() resource.Resource {
	return []func() resource.Resource{
 		NewWorkflowResource,
		NewTagResource,
		NewVariableResource,
		NewProjectResource,
		NewCredentialResource,
		NewExecutionRetryResource,
		NewSourceControlPullResource,
		NewUserResource,
		NewProjectUserResource,
		NewWorkflowTransferResource,
		NewCredentialTransferResource,
	}
}
