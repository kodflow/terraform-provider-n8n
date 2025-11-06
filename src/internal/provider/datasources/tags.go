package datasources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	providertypes "github.com/kodflow/n8n/src/internal/provider/types"
)

// Ensure TagsDataSource implements required interfaces.
var _ datasource.DataSource = &TagsDataSource{}
var _ datasource.DataSourceWithConfigure = &TagsDataSource{}

// TagsDataSource defines the data source implementation for listing tags.
type TagsDataSource struct {
	client *providertypes.N8nClient
}

// TagsDataSourceModel describes the data source data model.
type TagsDataSourceModel struct {
	Tags []TagItemModel `tfsdk:"tags"`
}

// TagItemModel represents a single tag in the list.
type TagItemModel struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	CreatedAt types.String `tfsdk:"created_at"`
	UpdatedAt types.String `tfsdk:"updated_at"`
}

// NewTagsDataSource creates a new TagsDataSource instance.
func NewTagsDataSource() datasource.DataSource {
 // Return result.
	return &TagsDataSource{}
}

// Metadata returns the data source type name.
func (d *TagsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tags"
}

// Schema defines the schema for the data source.
func (d *TagsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
func (d *TagsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Check for nil value.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*providertypes.N8nClient)
	// Check condition.
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *providertypes.N8nClient, got: %T", req.ProviderData),
		)
		// Return result.
		return
	}

	d.client = client
}

// Read refreshes the Terraform state with the latest data.
func (d *TagsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data TagsDataSourceModel

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

	data.Tags = make([]TagItemModel, 0)
	// Check for non-nil value.
	if tagList.Data != nil {
  // Iterate over items.
		for _, tag := range tagList.Data {
			item := TagItemModel{
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
