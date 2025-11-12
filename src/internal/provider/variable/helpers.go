// Package variable contains helper functions for variable operations.
package variable

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/shared"
	"github.com/kodflow/n8n/src/internal/provider/variable/models"
)

// findVariableByIDOrKey searches for a variable by ID or key in a variable list.
// Returns the found variable and a boolean indicating if it was found.
//
// Params:
//   - variables: Slice of variables to search through
//   - id: Variable ID to search for
//   - key: Variable key to search for
//
// Returns:
//   - *n8nsdk.Variable: Pointer to found variable or nil
//   - bool: True if variable was found, false otherwise
func findVariableByIDOrKey(variables []n8nsdk.Variable, id, key types.String) (*n8nsdk.Variable, bool) {
	// Iterate over items.
	for _, variable := range variables {
		// Check for null value.
		matchByID := !id.IsNull() && variable.Id != nil && *variable.Id == id.ValueString()
		// Check for null value.
		matchByKey := !key.IsNull() && variable.Key == key.ValueString()

		// Check condition.
		if matchByID || matchByKey {
			// Return result.
			return &variable, true
		}
	}
	// Return result.
	return nil, false
}

// findVariableByID searches for a variable by ID in a variable list.
// Returns the found variable and a boolean indicating if it was found.
//
// Params:
//   - variables: Slice of variables to search through
//   - id: Variable ID to search for
//
// Returns:
//   - *n8nsdk.Variable: Pointer to found variable or nil
//   - bool: True if variable was found, false otherwise
func findVariableByID(variables []n8nsdk.Variable, id string) (*n8nsdk.Variable, bool) {
	// Iterate over items.
	for _, variable := range variables {
		// Check for non-nil value.
		if variable.Id != nil && *variable.Id == id {
			// Return result.
			return &variable, true
		}
	}
	// Return result.
	return nil, false
}

// findVariableByKey searches for a variable by key in a variable list.
// Returns the found variable and a boolean indicating if it was found.
//
// Params:
//   - variables: Slice of variables to search through
//   - key: Variable key to search for
//
// Returns:
//   - *n8nsdk.Variable: Pointer to found variable or nil
//   - bool: True if variable was found, false otherwise
func findVariableByKey(variables []n8nsdk.Variable, key string) (*n8nsdk.Variable, bool) {
	// Iterate over items.
	for _, variable := range variables {
		// Check condition.
		if variable.Key == key {
			// Return result.
			return &variable, true
		}
	}
	// Return result.
	return nil, false
}

// mapVariableToDataSourceModel maps an SDK variable to the datasource model.
//
// Params:
//   - variable: SDK variable object to map
//   - data: Target datasource model to populate
//
// Returns:
//   - No return value, modifies data parameter in place
func mapVariableToDataSourceModel(variable *n8nsdk.Variable, data *models.DataSource) {
	// Check for non-nil value.
	if variable.Id != nil {
		data.ID = types.StringValue(*variable.Id)
	}
	data.Key = types.StringValue(variable.Key)
	data.Value = types.StringValue(variable.Value)
	// Check for non-nil value.
	if variable.Type != nil {
		data.Type = types.StringPointerValue(variable.Type)
	}
	// Project is a nested object, extract ID if present
	// Check for non-nil value.
	if variable.Project != nil && variable.Project.Id != nil {
		data.ProjectID = types.StringPointerValue(variable.Project.Id)
	}
}

// mapVariableToResourceModel maps an SDK variable to the resource model.
//
// Params:
//   - variable: SDK variable object to map
//   - data: Target resource model to populate
//
// Returns:
//   - No return value, modifies data parameter in place
func mapVariableToResourceModel(variable *n8nsdk.Variable, data *models.Resource) {
	// Check for non-nil value.
	if variable.Id != nil {
		data.ID = types.StringPointerValue(variable.Id)
	}
	data.Key = types.StringValue(variable.Key)
	data.Value = types.StringValue(variable.Value)
	// Check for non-nil value.
	if variable.Type != nil {
		data.Type = types.StringPointerValue(variable.Type)
	}
	// Project is a nested object, extract ID if present
	// Check for non-nil value.
	if variable.Project != nil && variable.Project.Id != nil {
		data.ProjectID = types.StringPointerValue(variable.Project.Id)
	}
}

// buildVariableRequest builds a variable request from a resource model.
//
// Params:
//   - plan: Resource model containing variable configuration
//
// Returns:
//   - n8nsdk.VariableCreate: API request structure for variable creation
func buildVariableRequest(plan *models.Resource) n8nsdk.VariableCreate {
	variableRequest := n8nsdk.VariableCreate{
		Key:   plan.Key.ValueString(),
		Value: plan.Value.ValueString(),
	}

	// Add optional fields
	// Check condition.
	if !plan.Type.IsNull() && !plan.Type.IsUnknown() {
		variableRequest.Type = shared.String(plan.Type.ValueString())
	}
	// Check condition.
	if !plan.ProjectID.IsNull() && !plan.ProjectID.IsUnknown() {
		variableRequest.ProjectId = *n8nsdk.NewNullableString(shared.String(plan.ProjectID.ValueString()))
	}

	// Return result.
	return variableRequest
}
