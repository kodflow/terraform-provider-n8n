package tag

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
)

// Ensure TagDataSource implements required interfaces.
var (
	_ datasource.DataSource              = &TagDataSource{}
	_ datasource.DataSourceWithConfigure = &TagDataSource{}
)

// TagDataSource is a Terraform datasource that provides read-only access to a single n8n tag.
// It fetches tag details from the n8n API using ID or name-based filtering.
type TagDataSource struct {
	client *client.N8nClient
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
//
// Params:
//   - None
// Returns:
//   - datasource.DataSource: A new TagDataSource instance
func NewTagDataSource() datasource.DataSource {
	// Return result.
	return &TagDataSource{}
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
func (d *TagDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tag"
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
//
// Params:
//   - ctx: The request context
//   - req: The configure request containing provider data
//   - resp: The configure response to handle errors
//
// Returns:
//   - None
func (d *TagDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Check for nil provider data.
	if req.ProviderData == nil {
		return
	}

	clientData, ok := req.ProviderData.(*client.N8nClient)
	// Check if provider data is correct type.
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
//   - ctx: The request context
//   - req: The read request containing configuration
//   - resp: The read response to populate with state
//
// Returns:
//   - None
func (d *TagDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	data := &TagDataSourceModel{}

	
	resp.Diagnostics.Append(req.Config.Get(ctx, data)...)
	// If there are errors from config parsing, return early.
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate that at least one identifier is provided.
	if data.ID.IsNull() && data.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Missing Required Attribute",
			"Either 'id' or 'name' must be specified",
		)
		return
	}

	// If ID is provided, use the direct GET endpoint.
	if !data.ID.IsNull() {
		tag, httpResp, err := d.client.APIClient.TagsAPI.TagsIdGet(ctx, data.ID.ValueString()).Execute()
		// Close HTTP response body if present.
		if httpResp != nil && httpResp.Body != nil {
			defer httpResp.Body.Close()
		}
		// Check if API call returned an error.
		if err != nil {
			resp.Diagnostics.AddError(
				"Error retrieving tag",
				fmt.Sprintf("Could not retrieve tag with ID %s: %s\nHTTP Response: %v", data.ID.ValueString(), err.Error(), httpResp),
			)
			return
		}

		// Map tag to model.
		mapTagToDataSourceModel(tag, data)
	} else {
		// If only name is provided, list all tags and filter client-side.
		tagList, httpResp, err := d.client.APIClient.TagsAPI.TagsGet(ctx).Execute()
		// Close HTTP response body if present.
		if httpResp != nil && httpResp.Body != nil {
			defer httpResp.Body.Close()
		}
		// Check if API call returned an error.
		if err != nil {
			resp.Diagnostics.AddError(
				"Error listing tags",
				fmt.Sprintf("Could not list tags: %s\nHTTP Response: %v", err.Error(), httpResp),
			)
			return
		}

		// Find tag by name in the response data.
		var tag *n8nsdk.Tag
		var found bool
		// Check if tag list data is not empty.
		if tagList.Data != nil {
			tag, found = findTagByName(tagList.Data, data.Name.ValueString())
		}

		// Return error if tag was not found.
		if !found {
			resp.Diagnostics.AddError(
				"Tag Not Found",
				fmt.Sprintf("Could not find tag with name: %s", data.Name.ValueString()),
			)
			return
		}

		// Map tag to model.
		mapTagToDataSourceModel(tag, data)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}
