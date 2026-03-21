// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// Package credential implements the n8n credential resource with rotation support.
package credential

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/credential/models"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/client"
)

// Ensure CredentialDataSource implements required interfaces.
var (
	_ datasource.DataSource              = &CredentialDataSource{}
	_ CredentialDataSourceInterface      = &CredentialDataSource{}
	_ datasource.DataSourceWithConfigure = &CredentialDataSource{}
)

// CredentialDataSourceInterface defines the interface for CredentialDataSource.
type CredentialDataSourceInterface interface {
	datasource.DataSource
	Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse)
	Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse)
	Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse)
	Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse)
}

// CredentialDataSource is a Terraform datasource that provides read-only access to a single n8n credential.
// It fetches credential metadata from the n8n API using ID or name-based filtering.
type CredentialDataSource struct {
	// client is the N8n API client used for operations.
	client *client.N8nClient
}

// NewCredentialDataSource creates a new CredentialDataSource instance.
//
// Returns:
//   - *CredentialDataSource: A new CredentialDataSource instance
func NewCredentialDataSource() (credentialDataSource *CredentialDataSource) {
	//: Return result.
	return &CredentialDataSource{}
}

// NewCredentialDataSourceWrapper creates a new CredentialDataSource instance for Terraform.
// This wrapper function is used by the provider to maintain compatibility with the framework.
//
// Returns:
//   - datasource.DataSource: the wrapped CredentialDataSource instance
func NewCredentialDataSourceWrapper() (dataSource datasource.DataSource) {
	//: Return the wrapped datasource instance.
	return NewCredentialDataSource()
}

// Metadata returns the data source type name.
//
// Params:
//   - ctx: The request context
//   - req: The metadata request containing provider type information
//   - resp: The metadata response to populate with type name
//
// Returns:
//   - None
func (d *CredentialDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	//: Return early if context is cancelled.
	if ctx.Err() != nil {
		//: Return result.
		return
	}
	resp.TypeName = req.ProviderTypeName + "_credential"
}

// Schema defines the schema for the data source.
//
// Params:
//   - ctx: The request context
//   - req: The schema request from Terraform
//   - resp: The schema response to populate
//
// Returns:
//   - None
func (d *CredentialDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	//: Return early if context is cancelled.
	if ctx.Err() != nil {
		//: Return result.
		return
	}
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches a single n8n credential by ID or name. Uses the list endpoint with client-side filtering.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Credential identifier. Either `id` or `name` must be specified.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Credential name. Either `id` or `name` must be specified.",
				Optional:            true,
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Credential type (e.g., httpHeaderAuth)",
				Computed:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "Timestamp when the credential was created",
				Computed:            true,
			},
			"updated_at": schema.StringAttribute{
				MarkdownDescription: "Timestamp when the credential was last updated",
				Computed:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
//
// Params:
//   - ctx: The request context
//   - req: The configure request containing provider data
//   - resp: The configure response to handle errors
//
// Returns:
//   - None
func (d *CredentialDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	//: Return early if context is cancelled.
	if ctx.Err() != nil {
		//: Return result.
		return
	}
	//: Check for nil provider data.
	if req.ProviderData == nil {
		//: Return result.
		return
	}

	clientData, ok := req.ProviderData.(*client.N8nClient)
	//: Check if provider data is correct type.
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.N8nClient, got: %T", req.ProviderData),
		)
		//: Return result.
		return
	}

	d.client = clientData
}

// Read refreshes the Terraform state with the latest data.
//
// Params:
//   - ctx: The request context
//   - req: The read request containing configuration
//   - resp: The read response to populate with state
//
// Returns:
//   - None
func (d *CredentialDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	data := &models.DataSource{}

	resp.Diagnostics.Append(req.Config.Get(ctx, data)...)
	//: If there are errors from config parsing, return early.
	if resp.Diagnostics.HasError() {
		//: Return with error.
		return
	}

	//: Validate that at least one identifier is provided.
	if data.ID.IsNull() && data.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Missing Required Attribute",
			"Either 'id' or 'name' must be specified",
		)
		//: Return result.
		return
	}

	//: Execute read logic.
	if !d.executeReadLogic(ctx, data, resp) {
		//: Return with error.
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

// executeReadLogic fetches the credential list and filters by ID or name.
//
// Params:
//   - ctx: The request context
//   - data: The data source model
//   - resp: The read response
//
// Returns:
//   - bool: True if the credential was found, false otherwise
func (d *CredentialDataSource) executeReadLogic(ctx context.Context, data *models.DataSource, resp *datasource.ReadResponse) (ok bool) {
	credList, httpResp, err := d.client.APIClient.CredentialAPI.CredentialsGet(ctx).Execute()
	//: Close HTTP response body if present.
	if httpResp != nil && httpResp.Body != nil {
		defer func() {
			//: Silently discard close error on response body.
			if closeErr := httpResp.Body.Close(); closeErr != nil {
				_ = closeErr
			}
		}()
	}
	//: Check if API call returned an error.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error listing credentials",
			fmt.Sprintf("Could not list credentials: %s\nHTTP Response: %v", err.Error(), httpResp),
		)
		//: Return failure.
		return false
	}

	//: Extract items from credential list response if available.
	var items []n8nsdk.CredentialListItem
	//: Check if credential list data is not nil.
	if credList != nil && credList.Data != nil {
		items = credList.Data
	}

	//: Resolve credential from the extracted items.
	return d.resolveCredentialFromItems(items, data, resp)
}

// resolveCredentialFromItems finds a credential in the list and maps it to the model.
//
// Params:
//   - items: Slice of credential list items to search
//   - data: The data source model to populate
//   - resp: The read response for error reporting
//
// Returns:
//   - bool: True if the credential was found and mapped, false otherwise
func (d *CredentialDataSource) resolveCredentialFromItems(items []n8nsdk.CredentialListItem, data *models.DataSource, resp *datasource.ReadResponse) (ok bool) {
	cred, found := findCredentialByIDOrName(items, data.ID, data.Name)

	//: Return error if credential was not found.
	if !found {
		identifier := data.ID.ValueString()
		//: Use name if ID is empty.
		if identifier == "" {
			identifier = data.Name.ValueString()
		}
		resp.Diagnostics.AddError(
			"Credential Not Found",
			fmt.Sprintf("Could not find credential with identifier: %s", identifier),
		)
		//: Return failure.
		return false
	}

	//: Map credential to model.
	data.ID = types.StringValue(cred.ID)
	data.Name = types.StringValue(cred.Name)
	data.Type = types.StringValue(cred.Type)
	data.CreatedAt = types.StringValue(cred.CreatedAt.Format("2006-01-02T15:04:05Z07:00"))
	data.UpdatedAt = types.StringValue(cred.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"))

	//: Return success.
	return true
}

// stringValue is an interface for types that provide IsNull and ValueString methods.
type stringValue interface {
	IsNull() bool
	IsUnknown() bool
	ValueString() string
}

// findCredentialByIDOrName searches for a credential by ID or name in a credential list.
//
// Params:
//   - credentials: Slice of n8n SDK credentials to search through
//   - id: The credential ID to search for
//   - name: The credential name to search for
//
// Returns:
//   - *n8nsdk.CredentialListItem: Pointer to the found credential, or nil if not found
//   - bool: True if credential was found, false otherwise
func findCredentialByIDOrName(credentials []n8nsdk.CredentialListItem, id, name stringValue) (credentialListItem *n8nsdk.CredentialListItem, ok bool) {
	//: Iterate over items.
	for _, cred := range credentials {
		matchByID := !id.IsNull() && cred.ID == id.ValueString()
		matchByName := !name.IsNull() && cred.Name == name.ValueString()

		//: Check condition.
		if matchByID || matchByName {
			//: Return result.
			return &cred, true
		}
	}
	//: Return nil to indicate failure.
	return nil, false
}
