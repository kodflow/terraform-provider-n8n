package datasources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	providertypes "github.com/kodflow/n8n/src/internal/provider/types"
)

// Ensure TagDataSource implements required interfaces.
var (
	_ datasource.DataSource              = &TagDataSource{}
	_ datasource.DataSourceWithConfigure = &TagDataSource{}
)

// TagDataSource is a Terraform datasource that provides read-only access to a single n8n tag.
// It fetches tag details from the n8n API using ID or name-based filtering.
type TagDataSource struct {
	client *providertypes.N8nClient
}

// TagDataSourceModel maps the Terraform schema to a single tag from the n8n API.
// It contains tag metadata including name and creation/update timestamps.
type TagDataSourceModel struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	CreatedAt types.String `tfsdk:"created_at"`
	UpdatedAt types.String `tfsdk:"updated_at"`
}

// NewTagDataSource creates a new TagDataSource instance.
func NewTagDataSource() datasource.DataSource {
	// Return result.
	return &TagDataSource{}
}

// Metadata returns the data source type name.
func (d *TagDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tag"
}

// Schema defines the schema for the data source.
func (d *TagDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches a single n8n tag by ID or name. When using ID, the API's GET /tags/{id} endpoint is used directly. When using name, the LIST endpoint is used with client-side filtering.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Tag identifier. Either `id` or `name` must be specified.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Tag name. Either `id` or `name` must be specified.",
				Optional:            true,
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
	}
}

// Configure adds the provider configured client to the data source.
func (d *TagDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *TagDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data TagDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	// Check condition.
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate that at least one identifier is provided
	if data.ID.IsNull() && data.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Missing Required Attribute",
			"Either 'id' or 'name' must be specified",
		)
		// Return result.
		return
	}

	// If ID is provided, use the direct GET endpoint
	if !data.ID.IsNull() {
		tag, httpResp, err := d.client.APIClient.TagsAPI.TagsIdGet(ctx, data.ID.ValueString()).Execute()
		// Check for non-nil value.
		if httpResp != nil && httpResp.Body != nil {
			defer httpResp.Body.Close()
		}
		// Check for error.
		if err != nil {
			resp.Diagnostics.AddError(
				"Error retrieving tag",
				fmt.Sprintf("Could not retrieve tag with ID %s: %s\nHTTP Response: %v", data.ID.ValueString(), err.Error(), httpResp),
			)
			// Return result.
			return
		}

		// Populate data model
		if tag.Id != nil {
			data.ID = types.StringValue(*tag.Id)
		}
		data.Name = types.StringValue(tag.Name)
		// Check for non-nil value.
		if tag.CreatedAt != nil {
			data.CreatedAt = types.StringValue(tag.CreatedAt.String())
		}
		// Check for non-nil value.
		if tag.UpdatedAt != nil {
			data.UpdatedAt = types.StringValue(tag.UpdatedAt.String())
		}
		// Handle alternative case.
	} else {
		// If only name is provided, list all tags and filter client-side
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

		// Filter by name
		var found bool
		// Check for non-nil value.
		if tagList.Data != nil {
			// Iterate over items.
			for _, tag := range tagList.Data {
				// Check condition.
				if tag.Name == data.Name.ValueString() {
					// Populate data model
					if tag.Id != nil {
						data.ID = types.StringValue(*tag.Id)
					}
					data.Name = types.StringValue(tag.Name)
					// Check for non-nil value.
					if tag.CreatedAt != nil {
						data.CreatedAt = types.StringValue(tag.CreatedAt.String())
					}
					// Check for non-nil value.
					if tag.UpdatedAt != nil {
						data.UpdatedAt = types.StringValue(tag.UpdatedAt.String())
					}
					found = true
					break
				}
			}
		}

		// Check condition.
		if !found {
			resp.Diagnostics.AddError(
				"Tag Not Found",
				fmt.Sprintf("Could not find tag with name: %s", data.Name.ValueString()),
			)
			// Return result.
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
