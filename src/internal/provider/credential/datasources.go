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
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/credential/models"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/client"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/constants"
)

// Ensure CredentialsDataSource implements required interfaces.
var (
	_ datasource.DataSource              = &CredentialsDataSource{}
	_ CredentialsDataSourceInterface     = &CredentialsDataSource{}
	_ datasource.DataSourceWithConfigure = &CredentialsDataSource{}
)

// CredentialsDataSourceInterface defines the interface for CredentialsDataSource.
type CredentialsDataSourceInterface interface {
	datasource.DataSource
	Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse)
	Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse)
	Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse)
	Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse)
}

// CredentialsDataSource is a Terraform datasource implementation for listing credentials.
// It provides read-only access to all n8n credentials through the n8n API.
type CredentialsDataSource struct {
	// client is the N8n API client used for operations.
	client *client.N8nClient
}

// NewCredentialsDataSource creates a new CredentialsDataSource instance.
//
// Returns:
//   - *CredentialsDataSource: a new CredentialsDataSource instance
func NewCredentialsDataSource() (credentialsDataSource *CredentialsDataSource) {
	//: Return result.
	return &CredentialsDataSource{}
}

// NewCredentialsDataSourceWrapper creates a new CredentialsDataSource instance for Terraform.
// This wrapper function is used by the provider to maintain compatibility with the framework.
//
// Returns:
//   - datasource.DataSource: the wrapped CredentialsDataSource instance
func NewCredentialsDataSourceWrapper() (dataSource datasource.DataSource) {
	//: Return the wrapped datasource instance.
	return NewCredentialsDataSource()
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
func (d *CredentialsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	//: Guard against cancelled context.
	if ctx.Err() != nil {
		//: Return early when context is cancelled.
		return
	}
	resp.TypeName = req.ProviderTypeName + "_credentials"
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
func (d *CredentialsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	//: Guard against cancelled context.
	if ctx.Err() != nil {
		//: Return early when context is cancelled.
		return
	}
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches a list of all n8n credentials",

		Attributes: map[string]schema.Attribute{
			"credentials": schema.ListNestedAttribute{
				MarkdownDescription: "List of credentials",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Credential identifier",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Credential name",
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
func (d *CredentialsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	//: Guard against cancelled context.
	if ctx.Err() != nil {
		//: Return early when context is cancelled.
		return
	}
	//: Check for nil value.
	if req.ProviderData == nil {
		//: Return result.
		return
	}

	clientData, ok := req.ProviderData.(*client.N8nClient)
	//: Check condition.
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
//   - req: The read request from Terraform
//   - resp: The read response to populate with data
//
// Returns:
//   - None
func (d *CredentialsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.DataSources

	//: Execute list logic.
	if !d.executeListLogic(ctx, &data, resp) {
		//: Return result.
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// executeListLogic fetches all credentials and populates the data model.
//
// Params:
//   - ctx: The request context
//   - data: The data sources model to populate
//   - resp: The read response
//
// Returns:
//   - bool: True if listing succeeded, false otherwise
func (d *CredentialsDataSource) executeListLogic(ctx context.Context, data *models.DataSources, resp *datasource.ReadResponse) (ok bool) {
	credList, httpResp, err := d.client.APIClient.CredentialAPI.CredentialsGet(ctx).Execute()
	//: Check for non-nil value.
	if httpResp != nil && httpResp.Body != nil {
		defer func() {
			//: Silently discard close error on response body.
			if closeErr := httpResp.Body.Close(); closeErr != nil {
				_ = closeErr
			}
		}()
	}
	//: Check for error.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error listing credentials",
			fmt.Sprintf("Could not list credentials: %s\nHTTP Response: %v", err.Error(), httpResp),
		)
		//: Return failure.
		return false
	}

	data.Credentials = make([]models.CredentialItem, 0, constants.DefaultListCapacity)
	//: Check for non-nil value.
	if credList != nil && credList.Data != nil {
		//: Iterate over items.
		for _, cred := range credList.Data {
			item := models.CredentialItem{
				ID:        types.StringValue(cred.ID),
				Name:      types.StringValue(cred.Name),
				Type:      types.StringValue(cred.Type),
				CreatedAt: types.StringValue(cred.CreatedAt.Format("2006-01-02T15:04:05Z07:00")),
				UpdatedAt: types.StringValue(cred.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")),
			}
			data.Credentials = append(data.Credentials, item)
		}
	}

	//: Return success.
	return true
}
