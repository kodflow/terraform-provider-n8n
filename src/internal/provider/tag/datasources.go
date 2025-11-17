// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// Package tag implements tag management resources and data sources.
package tag

import (
	"context"
	"fmt"

	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/constants"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/client"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/tag/models"
)

// Ensure TagsDataSource implements required interfaces.
var (
	_ datasource.DataSource              = &TagsDataSource{}
	_ TagsDataSourceInterface            = &TagsDataSource{}
	_ datasource.DataSourceWithConfigure = &TagsDataSource{}
)

// TagsDataSourceInterface defines the interface for TagsDataSource.
type TagsDataSourceInterface interface {
	datasource.DataSource
	Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse)
	Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse)
	Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse)
	Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse)
}

// TagsDataSource is a Terraform datasource implementation for listing tags.
// It provides read-only access to all n8n tags through the n8n API.
type TagsDataSource struct {
	// client is the N8n API client used for operations.
	client *client.N8nClient
}

// NewTagsDataSource creates a new TagsDataSource instance.
//
// Returns:
//   - datasource.DataSource: une nouvelle instance de TagsDataSource
func NewTagsDataSource() *TagsDataSource {
	// Return result.
	return &TagsDataSource{}
}

// NewTagsDataSourceWrapper creates a new TagsDataSource instance for Terraform.
// This wrapper function is used by the provider to maintain compatibility with the framework.
//
// Returns:
//   - datasource.DataSource: the wrapped TagsDataSource instance
func NewTagsDataSourceWrapper() datasource.DataSource {
	// Return the wrapped datasource instance.
	return NewTagsDataSource()
}

// Metadata returns the data source type name.
//
// Params:
//   - ctx: contexte de la requête
//   - req: requête de métadonnées
//   - resp: réponse de métadonnées
func (d *TagsDataSource) Metadata(_ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tags"
}

// Schema defines the schema for the data source.
//
// Params:
//   - ctx: contexte de la requête
//   - req: requête de schéma
//   - resp: réponse de schéma
func (d *TagsDataSource) Schema(_ctx context.Context, _req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches a list of all n8n tags",

		Attributes: map[string]schema.Attribute{
			"tags": schema.ListNestedAttribute{
				MarkdownDescription: "List of tags",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Tag identifier",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Tag name",
							Computed:            true,
						},
						"created_at": schema.StringAttribute{
							MarkdownDescription: "Timestamp when the tag was created",
							Computed:            true,
						},
						"updated_at": schema.StringAttribute{
							MarkdownDescription: "Timestamp when the tag was last updated",
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
//   - ctx: contexte de la requête
//   - req: requête de configuration
//   - resp: réponse de configuration
func (d *TagsDataSource) Configure(_ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Check for nil value.
	if req.ProviderData == nil {
		// Return result.
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
//   - ctx: contexte de la requête
//   - req: requête de lecture
//   - resp: réponse de lecture
func (d *TagsDataSource) Read(ctx context.Context, _req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.DataSources

	tagList, httpResp, err := d.client.APIClient.TagsAPI.TagsGet(ctx).Execute()
	// Check for non-nil value.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}
	// Check for error.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error listing tags",
			fmt.Sprintf("Could not list tags: %s\nHTTP Response: %v", err.Error(), httpResp),
		)
		// Return result.
		return
	}

	data.Tags = make([]models.Item, 0, constants.DEFAULT_LIST_CAPACITY)
	// Check for non-nil value.
	if tagList.Data != nil {
		// Iterate over items.
		for _, tag := range tagList.Data {
			item := models.Item{
				Name: types.StringValue(tag.Name),
			}
			// Check for non-nil value.
			if tag.Id != nil {
				item.ID = types.StringValue(*tag.Id)
			}
			// Check for non-nil value.
			if tag.CreatedAt != nil {
				item.CreatedAt = types.StringValue(tag.CreatedAt.String())
			}
			// Check for non-nil value.
			if tag.UpdatedAt != nil {
				item.UpdatedAt = types.StringValue(tag.UpdatedAt.String())
			}
			data.Tags = append(data.Tags, item)
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
