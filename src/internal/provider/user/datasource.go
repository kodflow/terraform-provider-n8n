// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// Package user implements user management resources and data sources.
package user

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/client"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/user/models"
)

// Ensure UserDataSource implements required interfaces.
var (
	_ datasource.DataSource              = &UserDataSource{}
	_ UserDataSourceInterface            = &UserDataSource{}
	_ datasource.DataSourceWithConfigure = &UserDataSource{}
)

// UserDataSourceInterface defines the interface for UserDataSource.
type UserDataSourceInterface interface {
	datasource.DataSource
	Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse)
	Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse)
	Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse)
	Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse)
}

// UserDataSource is a Terraform datasource implementation for retrieving a single user.
// It provides read-only access to n8n user details through the n8n API.
type UserDataSource struct {
	// client is the N8n API client used for operations.
	client *client.N8nClient
}

// NewUserDataSource creates a new UserDataSource instance.
//
// Returns:
//   - datasource.DataSource: new user data source instance
func NewUserDataSource() *UserDataSource {
	// Return result.
	return &UserDataSource{}
}

// NewUserDataSourceWrapper creates a new UserDataSource instance for Terraform.
// This wrapper function is used by the provider to maintain compatibility with the framework.
//
// Returns:
//   - datasource.DataSource: the wrapped UserDataSource instance
func NewUserDataSourceWrapper() datasource.DataSource {
	// Return the wrapped datasource instance.
	return NewUserDataSource()
}

// Metadata returns the data source type name.
//
// Params:
//   - ctx: context.Context for request lifecycle
//   - req: datasource.MetadataRequest containing provider type name
//   - resp: datasource.MetadataResponse to set the type name
func (d *UserDataSource) Metadata(_ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

// Schema defines the schema for the data source.
//
// Params:
//   - ctx: context.Context for request lifecycle
//   - req: datasource.SchemaRequest for schema definition
//   - resp: datasource.SchemaResponse to set the schema
func (d *UserDataSource) Schema(_ctx context.Context, _req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches a single n8n user by ID or email. The API accepts both ID and email as identifiers.",
		Attributes:          d.schemaAttributes(),
	}
}

// schemaAttributes returns the attribute definitions for the user data source schema.
//
// Returns:
//   - map[string]schema.Attribute: the data source attribute definitions
func (d *UserDataSource) schemaAttributes() map[string]schema.Attribute {
	// Return schema attributes.
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "User identifier. Either `id` or `email` must be specified. Only available for instance owners.",
			Optional:            true,
			Computed:            true,
		},
		"email": schema.StringAttribute{
			MarkdownDescription: "User email address. Either `id` or `email` must be specified.",
			Optional:            true,
			Computed:            true,
		},
		"first_name": schema.StringAttribute{
			MarkdownDescription: "User's first name",
			Computed:            true,
		},
		"last_name": schema.StringAttribute{
			MarkdownDescription: "User's last name",
			Computed:            true,
		},
		"is_pending": schema.BoolAttribute{
			MarkdownDescription: "Whether the user finished setting up their account in response to the invitation (false) or not (true)",
			Computed:            true,
		},
		"created_at": schema.StringAttribute{
			MarkdownDescription: "Timestamp when the user was created",
			Computed:            true,
		},
		"updated_at": schema.StringAttribute{
			MarkdownDescription: "Timestamp when the user was last updated",
			Computed:            true,
		},
		"role": schema.StringAttribute{
			MarkdownDescription: "User's global role (e.g., 'global:owner', 'global:admin', 'global:member')",
			Computed:            true,
		},
	}
}

// Configure adds the provider configured client to the data source.
//
// Params:
//   - ctx: context.Context for request lifecycle
//   - req: datasource.ConfigureRequest containing provider data
//   - resp: datasource.ConfigureResponse for diagnostics
func (d *UserDataSource) Configure(_ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Check for nil value.
	if req.ProviderData == nil {
		// Return with error.
		return
	}

	clientData, ok := req.ProviderData.(*client.N8nClient)
	// Check condition.
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.N8nClient, got: %T", req.ProviderData),
		)
		// Return result.
		return
	}

	d.client = clientData
}

// Read refreshes the Terraform state with the latest data.
//
// Params:
//   - ctx: context.Context for request lifecycle
//   - req: datasource.ReadRequest containing configuration
//   - resp: datasource.ReadResponse to set state
func (d *UserDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	data := &models.DataSource{}

	resp.Diagnostics.Append(req.Config.Get(ctx, data)...)
	// Check condition.
	if resp.Diagnostics.HasError() {
		// Return with error.
		return
	}

	// Validate and get identifier
	identifier := d.getIdentifier(data, resp)
	// Check condition.
	if identifier == "" {
		// Return result.
		return
	}

	// Fetch user from API
	user := d.fetchUser(ctx, identifier, resp)
	// Check for nil value.
	if user == nil {
		// Return result.
		return
	}

	// Populate data model
	d.populateUserData(user, data)

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

// getIdentifier validates and returns the user identifier.
//
// Params:
//   - data: user data source model
//   - resp: read response
//
// Returns:
//   - string: the identifier (ID or email)
func (d *UserDataSource) getIdentifier(data *models.DataSource, resp *datasource.ReadResponse) string {
	// Validate that at least one identifier is provided
	if data.ID.IsNull() && data.Email.IsNull() {
		resp.Diagnostics.AddError(
			"Missing Required Attribute",
			"Either 'id' or 'email' must be specified",
		)
		// Return result.
		return ""
	}

	// Use ID if provided, otherwise use email
	identifier := data.ID.ValueString()
	// Check condition.
	if identifier == "" {
		identifier = data.Email.ValueString()
	}
	// Return result.
	return identifier
}

// fetchUser retrieves a user from the API.
//
// Params:
//   - ctx: context for request
//   - identifier: user ID or email
//   - resp: read response
//
// Returns:
//   - *n8nsdk.User: the user or nil if error occurred
func (d *UserDataSource) fetchUser(ctx context.Context, identifier string, resp *datasource.ReadResponse) *n8nsdk.User {
	user, httpResp, err := d.client.APIClient.UserAPI.UsersIdGet(ctx, identifier).Execute()
	// Check for non-nil value.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}
	// Check for error.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error retrieving user",
			fmt.Sprintf("Could not retrieve user with identifier %s: %s\nHTTP Response: %v", identifier, err.Error(), httpResp),
		)
		// Return with error.
		return nil
	}
	// Return result.
	return user
}

// populateUserData populates the data model from the user.
//
// Params:
//   - user: source user
//   - data: target data model
func (d *UserDataSource) populateUserData(user *n8nsdk.User, data *models.DataSource) {
	// Check condition.
	if user.Id != nil {
		data.ID = types.StringValue(*user.Id)
	}
	data.Email = types.StringValue(user.Email)
	// Check condition.
	if user.FirstName != nil {
		data.FirstName = types.StringPointerValue(user.FirstName)
	}
	// Check condition.
	if user.LastName != nil {
		data.LastName = types.StringPointerValue(user.LastName)
	}
	// Check condition.
	if user.IsPending != nil {
		data.IsPending = types.BoolPointerValue(user.IsPending)
	}
	// Check condition.
	if user.CreatedAt != nil {
		data.CreatedAt = types.StringValue(user.CreatedAt.String())
	}
	// Check condition.
	if user.UpdatedAt != nil {
		data.UpdatedAt = types.StringValue(user.UpdatedAt.String())
	}
	// Check condition.
	if user.Role != nil {
		data.Role = types.StringPointerValue(user.Role)
	}
}
