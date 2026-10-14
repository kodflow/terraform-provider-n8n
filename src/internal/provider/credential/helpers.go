// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// Package credential provides helper functions for credential operations.
package credential

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/credential/models"
)

// FLOAT64_BIT_SIZE is the bit size for float64 parsing.
const FLOAT64_BIT_SIZE int = 64

// ALREADY_IN_PROJECT_MESSAGE is returned by n8n when transferring a credential into the project that owns it.
const ALREADY_IN_PROJECT_MESSAGE string = "already owning it"

// convertDataToSchemaTypes converts string values in credential data to their proper types
// based on the credential schema from n8n API.
// This enables support for credential types that require number or boolean fields.
//
// Params:
//   - ctx: Context for the API call
//   - credType: The credential type name (e.g., "imap", "mongoDb")
//   - data: The credential data with string values
//
// Returns:
//   - map[string]any: The credential data with converted types
func (r *CredentialResource) convertDataToSchemaTypes(ctx context.Context, credType string, data map[string]any) map[string]any {
	// Fetch credential schema from n8n API.
	schema, httpResp, err := r.client.APIClient.CredentialAPI.
		CredentialsSchemaCredentialTypeNameGet(ctx, credType).
		Execute()
	// Close response body if present.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}

	// If schema fetch fails, return original data with warning.
	if err != nil {
		tflog.Warn(ctx, fmt.Sprintf(
			"Could not fetch credential schema for type '%s': %s. Using string values as-is.",
			credType, err.Error(),
		))
		// Return original data without type conversion.
		return data
	}

	// Convert data based on schema properties.
	return r.applySchemaTypeConversions(ctx, schema, data)
}

// applySchemaTypeConversions applies type conversions based on the schema.
//
// Params:
//   - ctx: Context for logging
//   - schema: The credential schema from n8n API
//   - data: The credential data with string values
//
// Returns:
//   - map[string]any: The credential data with converted types
func (r *CredentialResource) applySchemaTypeConversions(ctx context.Context, schema map[string]any, data map[string]any) map[string]any {
	// Extract properties from schema.
	properties, ok := schema["properties"].(map[string]any)
	// If no properties found, return original data.
	if !ok {
		tflog.Debug(ctx, "No properties found in credential schema, using original data")
		// Return original data.
		return data
	}

	// Create result map with converted values.
	result := make(map[string]any, len(data))
	// Iterate over data keys to convert each value.
	for key, value := range data {
		// Only convert string values.
		strValue, isString := value.(string)
		// If not a string, keep original value.
		if !isString {
			result[key] = value
			continue
		}

		// Get property schema for this key.
		propSchema, propExists := properties[key].(map[string]any)
		// If no property schema, keep original string value.
		if !propExists {
			result[key] = value
			continue
		}

		// Convert based on property type.
		result[key] = convertValueByType(ctx, key, strValue, propSchema)
	}

	// Return converted data.
	return result
}

// convertValueByType converts a string value to the appropriate type based on schema.
//
// Params:
//   - ctx: Context for logging
//   - key: The property key (for logging)
//   - value: The string value to convert
//   - propSchema: The property schema containing type information
//
// Returns:
//   - any: The converted value (or original string if conversion fails)
func convertValueByType(ctx context.Context, key, value string, propSchema map[string]any) any {
	// Get the type from property schema.
	propType, hasType := propSchema["type"].(string)
	// If no type specified, return original value.
	if !hasType {
		// Return original string value.
		return value
	}

	// Convert based on type.
	switch propType {
	// Handle numeric types (number, integer) by converting string to float64.
	case "number", "integer":
		// Try to parse as float64 (handles both int and float).
		if f, err := strconv.ParseFloat(value, FLOAT64_BIT_SIZE); err == nil {
			tflog.Debug(ctx, fmt.Sprintf("Converted '%s' from string to number: %v", key, f))
			// Return converted number.
			return f
		}
		tflog.Debug(ctx, fmt.Sprintf("Could not convert '%s' value '%s' to number, keeping as string", key, value))
	// Handle boolean type by converting string to bool.
	case "boolean":
		// Try to parse as boolean.
		if b, err := strconv.ParseBool(value); err == nil {
			tflog.Debug(ctx, fmt.Sprintf("Converted '%s' from string to boolean: %v", key, b))
			// Return converted boolean.
			return b
		}
		tflog.Debug(ctx, fmt.Sprintf("Could not convert '%s' value '%s' to boolean, keeping as string", key, value))
	// Handle other types (string, array, object, etc.) - keep as string.
	default:
	}

	// Return original string value for other types or conversion failures.
	return value
}

// transferCredentialToProject transfers a credential to a specified project.
//
// Params:
//   - ctx: Context for the API call
//   - credentialID: The credential ID to transfer
//   - projectID: The destination project ID
//   - diags: Diagnostics collector for errors
//
// Returns:
//   - bool: True if transfer succeeded, false otherwise
func (r *CredentialResource) transferCredentialToProject(ctx context.Context, credentialID, projectID string, diags *diag.Diagnostics) bool {
	transferRequest := n8nsdk.CredentialsIdTransferPutRequest{
		DestinationProjectId: projectID,
	}

	httpResp, err := r.client.APIClient.CredentialAPI.
		CredentialsIdTransferPut(ctx, credentialID).
		CredentialsIdTransferPutRequest(transferRequest).
		Execute()

	// Close response body if present.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}

	// Check for error.
	if err != nil {
		// Transferring into the project that already owns the credential is a no-op
		// (project_id is not in the state after an import).
		var apiErr *n8nsdk.GenericOpenAPIError
		if errors.As(err, &apiErr) && bytes.Contains(apiErr.Body(), []byte(ALREADY_IN_PROJECT_MESSAGE)) {
			tflog.Info(ctx, fmt.Sprintf("Credential %s is already in project %s", credentialID, projectID))
			// Return success.
			return true
		}
		diags.AddError(
			"Error transferring credential to project",
			fmt.Sprintf("Could not transfer credential ID %s to project %s: %s\nHTTP Response: %v", credentialID, projectID, err.Error(), httpResp),
		)
		// Return failure.
		return false
	}

	tflog.Info(ctx, fmt.Sprintf("Transferred credential %s to project %s", credentialID, projectID))
	// Return success.
	return true
}

// mapCredentialProjectID maps the project_id to the model.
// Since n8n API doesn't return project info in credential responses,
// we keep the value from the plan.
//
// Params:
//   - plan: The resource model to update
//   - requestedProjectID: The project ID that was requested
func mapCredentialProjectID(plan *models.Resource, requestedProjectID types.String) {
	// If project_id was set in plan, keep it.
	if !requestedProjectID.IsNull() && !requestedProjectID.IsUnknown() {
		plan.ProjectID = requestedProjectID
	} else {
		// Set to null if not specified.
		plan.ProjectID = types.StringNull()
	}
}
