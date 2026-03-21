// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// Package project implements n8n project management resources and data sources.
package project

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/project/models"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/client"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/constants"
)

// Ensure ProjectMembersDataSource implements required interfaces.
var (
	_ datasource.DataSource              = &ProjectMembersDataSource{}
	_ ProjectMembersDataSourceInterface  = &ProjectMembersDataSource{}
	_ datasource.DataSourceWithConfigure = &ProjectMembersDataSource{}
)

// ProjectMembersDataSourceInterface defines the interface for ProjectMembersDataSource.
type ProjectMembersDataSourceInterface interface {
	datasource.DataSource
	Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse)
	Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse)
	Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse)
	Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse)
}

// ProjectMembersDataSource is a Terraform datasource that provides read-only access to project members.
// It fetches members of a specific n8n project from the n8n API.
type ProjectMembersDataSource struct {
	// client is the N8n API client used for operations.
	client *client.N8nClient
}

// NewProjectMembersDataSource creates a new ProjectMembersDataSource instance.
//
// Returns:
//   - *ProjectMembersDataSource: A new ProjectMembersDataSource instance
func NewProjectMembersDataSource() (projectMembersDataSource *ProjectMembersDataSource) {
	//: Return new empty ProjectMembersDataSource instance.
	return &ProjectMembersDataSource{}
}

// NewProjectMembersDataSourceWrapper creates a new ProjectMembersDataSource instance for Terraform.
// This wrapper function is used by the provider to maintain compatibility with the framework.
//
// Returns:
//   - datasource.DataSource: the wrapped ProjectMembersDataSource instance
func NewProjectMembersDataSourceWrapper() (dataSource datasource.DataSource) {
	//: Return the wrapped datasource instance.
	return NewProjectMembersDataSource()
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
func (d *ProjectMembersDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	//: Return early if context is cancelled.
	if ctx.Err() != nil {
		//: Return when context is cancelled.
		return
	}
	resp.TypeName = req.ProviderTypeName + "_project_members"
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
func (d *ProjectMembersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	//: Return early if context is cancelled.
	if ctx.Err() != nil {
		//: Return when context is cancelled.
		return
	}
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches the list of members for an n8n project.",

		Attributes: map[string]schema.Attribute{
			"project_id": schema.StringAttribute{
				MarkdownDescription: "Project identifier",
				Required:            true,
			},
			"members": schema.ListNestedAttribute{
				MarkdownDescription: "List of project members",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"user_id": schema.StringAttribute{
							MarkdownDescription: "User identifier",
							Computed:            true,
						},
						"role": schema.StringAttribute{
							MarkdownDescription: "User role in the project (e.g., project:admin, project:editor, project:viewer)",
							Computed:            true,
						},
						"email": schema.StringAttribute{
							MarkdownDescription: "User email address",
							Computed:            true,
						},
						"first_name": schema.StringAttribute{
							MarkdownDescription: "User first name",
							Computed:            true,
						},
						"last_name": schema.StringAttribute{
							MarkdownDescription: "User last name",
							Computed:            true,
						},
					},
				},
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
func (d *ProjectMembersDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	//: Return early if context is cancelled.
	if ctx.Err() != nil {
		//: Return when context is cancelled.
		return
	}
	//: Skip configuration when provider data is not yet available.
	if req.ProviderData == nil {
		//: Return early when no provider data is set.
		return
	}

	clientData, ok := req.ProviderData.(*client.N8nClient)
	//: Validate the provider data type is the expected N8nClient.
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.N8nClient, got: %T", req.ProviderData),
		)
		//: Return after reporting the type mismatch error.
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
func (d *ProjectMembersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	data := &models.DataSourceMembers{}

	resp.Diagnostics.Append(req.Config.Get(ctx, data)...)
	//: Check for errors from config parsing and return early if any.
	if resp.Diagnostics.HasError() {
		//: Return after config parsing failure.
		return
	}

	//: Execute read logic to fetch members from the API.
	if !d.executeReadLogic(ctx, data, resp) {
		//: Return after read logic failure.
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

// executeReadLogic retrieves project members from the API.
//
// Params:
//   - ctx: The request context
//   - data: The data source model
//   - resp: The read response
//
// Returns:
//   - ok: True if read succeeded, false otherwise
func (d *ProjectMembersDataSource) executeReadLogic(ctx context.Context, data *models.DataSourceMembers, resp *datasource.ReadResponse) (ok bool) {
	memberList, httpResp, err := d.client.APIClient.ProjectsAPI.ProjectsProjectIdUsersGet(ctx, data.ProjectID.ValueString()).Execute()
	//: Close HTTP response body if present to prevent resource leaks.
	if httpResp != nil && httpResp.Body != nil {
		defer func() {
			//: Silently discard close error on response body.
			if closeErr := httpResp.Body.Close(); closeErr != nil {
				resp.Diagnostics.AddWarning("Failed to close response body", closeErr.Error())
			}
		}()
	}
	//: Check if API call returned an error.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading project members",
			fmt.Sprintf("Could not read members for project %s: %s\nHTTP Response: %v",
				data.ProjectID.ValueString(), err.Error(), httpResp),
		)
		//: Return failure when API call fails.
		return false
	}

	data.Members = make([]models.MemberItem, 0, constants.DefaultListCapacity)
	//: Populate members list when response contains data.
	if memberList != nil && memberList.Data != nil {
		//: Iterate over items and convert each to MemberItem.
		for _, member := range memberList.Data {
			data.Members = append(data.Members, memberToItem(member))
		}
	}

	//: Return success when all members have been processed.
	return true
}

// memberToItem converts an SDK ProjectMember to a MemberItem model.
//
// Params:
//   - member: The SDK project member
//
// Returns:
//   - memberItem: The converted member item
func memberToItem(member n8nsdk.ProjectMember) (memberItem models.MemberItem) {
	//: Return the converted member item.
	return models.MemberItem{
		UserID:    stringPtrToTF(member.ID),
		Role:      stringPtrToTF(member.Role),
		Email:     stringPtrToTF(member.Email),
		FirstName: stringPtrToTF(member.FirstName),
		LastName:  stringPtrToTF(member.LastName),
	}
}

// stringPtrToTF converts a *string pointer to a Terraform types.String value.
//
// Params:
//   - s: Pointer to string, may be nil
//
// Returns:
//   - types.String: StringValue if non-nil, StringNull otherwise
func stringPtrToTF(s *string) (result types.String) {
	//: Return string value when pointer is non-nil.
	if s != nil {
		//: Return wrapped string value.
		return types.StringValue(*s)
	}
	//: Return null string when pointer is nil.
	return types.StringNull()
}
