package project

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
)

// Ensure ProjectDataSource implements required interfaces.
var (
	_ datasource.DataSource              = &ProjectDataSource{}
	_ datasource.DataSourceWithConfigure = &ProjectDataSource{}
)

// ProjectDataSource is a Terraform datasource that provides read-only access to a single n8n project.
// It fetches project details from the n8n API using ID or name-based filtering.
type ProjectDataSource struct {
	client *client.N8nClient
}

// ProjectDataSourceModel maps the Terraform schema to a single project from the n8n API.
// It contains project metadata including timestamps, type, and descriptive information.
type ProjectDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Type        types.String `tfsdk:"type"`
	CreatedAt   types.String `tfsdk:"created_at"`
	UpdatedAt   types.String `tfsdk:"updated_at"`
	Icon        types.String `tfsdk:"icon"`
	Description types.String `tfsdk:"description"`
}

// NewProjectDataSource creates a new ProjectDataSource instance.
//
// Returns:
//   - datasource.DataSource: New ProjectDataSource instance
func NewProjectDataSource() datasource.DataSource {
	// Return result.
	return &ProjectDataSource{}
}

// Metadata returns the data source type name.
//
// Params:
//   - ctx: Context for the request
//   - req: Metadata request
//   - resp: Metadata response to populate
//
// Returns:
func (d *ProjectDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

// Schema defines the schema for the data source.
//
// Params:
//   - ctx: Context for the request
//   - req: Schema request
//   - resp: Schema response to populate
//
// Returns:
func (d *ProjectDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches a single n8n project by ID or name. Since the n8n API doesn't provide a GET /projects/{id} endpoint, this datasource uses the LIST endpoint with client-side filtering.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Project identifier. Either `id` or `name` must be specified.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Project name. Either `id` or `name` must be specified.",
				Optional:            true,
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Project type (e.g., 'team', 'personal')",
				Computed:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "Timestamp when the project was created",
				Computed:            true,
			},
			"updated_at": schema.StringAttribute{
				MarkdownDescription: "Timestamp when the project was last updated",
				Computed:            true,
			},
			"icon": schema.StringAttribute{
				MarkdownDescription: "Project icon",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Project description",
				Computed:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
//
// Params:
//   - ctx: Context for the request
//   - req: Configure request
//   - resp: Configure response to populate
//
// Returns:
func (d *ProjectDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Check for nil value.
	if req.ProviderData == nil {
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
//   - ctx: Context for the request
//   - req: Read request
//   - resp: Read response to populate
//
// Returns:
func (d *ProjectDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	data := &ProjectDataSourceModel{}

	resp.Diagnostics.Append(req.Config.Get(ctx, data)...)
	// Check for diagnostics errors.
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

	// List all projects and filter client-side (API limitation).
	projectList, httpResp, err := d.client.APIClient.ProjectsAPI.ProjectsGet(ctx).Execute()
	// Close response body if available.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}
	// Check for API errors.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error listing projects",
			fmt.Sprintf("Could not list projects: %s\nHTTP Response: %v", err.Error(), httpResp),
		)
		return
	}

	// Find project by ID or name.
	var project *n8nsdk.Project
	var found bool
	// Search for project in the list.
	if projectList.Data != nil {
		project, found = findProjectByIDOrName(projectList.Data, data.ID, data.Name)
	}

	// Check if project was found.
	if !found {
		identifier := data.ID.ValueString()
		// Use name if ID is empty.
		if identifier == "" {
			identifier = data.Name.ValueString()
		}
		resp.Diagnostics.AddError(
			"Project Not Found",
			fmt.Sprintf("Could not find project with identifier: %s", identifier),
		)
		return
	}

	// Map project to model.
	mapProjectToDataSourceModel(project, data)

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}
