// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE.md in the project root for license information.

// Package user implements user management resources and data sources.
package user

import (
	"context"
	"fmt"

	"github.com/kodflow/n8n/src/internal/provider/shared/constants"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
	"github.com/kodflow/n8n/src/internal/provider/user/models"
)

// Ensure UsersDataSource implements required interfaces.
var (
	_ datasource.DataSource              = &UsersDataSource{}
	_ UsersDataSourceInterface           = &UsersDataSource{}
	_ datasource.DataSourceWithConfigure = &UsersDataSource{}
)

// UsersDataSourceInterface defines the interface for UsersDataSource.
type UsersDataSourceInterface interface {
	datasource.DataSource
	Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse)
	Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse)
	Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse)
	Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse)
}

// UsersDataSource is a Terraform datasource implementation for listing users.
// It provides read-only access to all n8n users through the n8n API.
type UsersDataSource struct {
	// client is the N8n API client used for operations.
	client *client.N8nClient
}

// NewUsersDataSource creates a new UsersDataSource instance.
//
// Returns:
//   - datasource.DataSource: the initialized UsersDataSource as a DataSource interface
func NewUsersDataSource() *UsersDataSource {
	// Return result.
	return &UsersDataSource{}
}

// NewUsersDataSourceWrapper creates a new UsersDataSource instance for Terraform.
// This wrapper function is used by the provider to maintain compatibility with the framework.
//
// Returns:
//   - datasource.DataSource: the wrapped UsersDataSource instance
func NewUsersDataSourceWrapper() datasource.DataSource {
	// Return the wrapped datasource instance.
	return NewUsersDataSource()
}

// Metadata returns the data source type name.
// Params:
//   - ctx: context for the request
//   - req: metadata request containing provider type name
//   - resp: metadata response to be populated with the datasource type name
func (d *UsersDataSource) Metadata(_ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_users"
}

// Schema defines the schema for the data source.
// Params:
//   - ctx: context for the request
//   - req: schema request for the data source
//   - resp: schema response to be populated with the datasource schema
func (d *UsersDataSource) Schema(_ctx context.Context, _req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches a list of all n8n users. Only available for instance owners.",
		Attributes:          d.schemaAttributes(),
	}
}

// schemaAttributes returns the attribute definitions for the users data source schema.
//
// Returns:
//   - map[string]schema.Attribute: the data source attribute definitions
func (d *UsersDataSource) schemaAttributes() map[string]schema.Attribute {
	// Return schema attributes.
	return map[string]schema.Attribute{
		"users": schema.ListNestedAttribute{
			MarkdownDescription: "List of users",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: d.userItemAttributes()},
		},
	}
}

// userItemAttributes returns the nested attribute definitions for individual user items.
//
// Returns:
//   - map[string]schema.Attribute: the user item attribute definitions
func (d *UsersDataSource) userItemAttributes() map[string]schema.Attribute {
	// Return schema attributes.
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "User identifier",
			Computed:            true,
		},
		"email": schema.StringAttribute{
			MarkdownDescription: "User email address",
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
			MarkdownDescription: "Whether the user finished setting up their account",
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
			MarkdownDescription: "User's global role",
			Computed:            true,
		},
	}
}

// Configure adds the provider configured client to the data source.
// Params:
//   - ctx: context for the request
//   - req: configure request containing the provider data
//   - resp: configure response for diagnostics
func (d *UsersDataSource) Configure(_ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
// Params:
//   - ctx: context for the request
//   - req: read request containing current state
//   - resp: read response to be populated with latest data
func (d *UsersDataSource) Read(ctx context.Context, _req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.DataSources

	userList, httpResp, err := d.client.APIClient.UserAPI.UsersGet(ctx).Execute()
	// Check for non-nil value.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}
	// Check for error.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error listing users",
			fmt.Sprintf("Could not list users: %s\nHTTP Response: %v", err.Error(), httpResp),
		)
		// Return result.
		return
	}

	data.Users = make([]models.Item, 0, constants.DEFAULT_LIST_CAPACITY)
	// Check if user data is present.
	if userList.Data != nil {
		// Iterate through each user in the list and map to item model.
		for _, user := range userList.Data {
			item := mapUserToItem(&user)
			data.Users = append(data.Users, item)
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
