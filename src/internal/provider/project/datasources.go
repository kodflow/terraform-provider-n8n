package project

import (
	"github.com/kodflow/n8n/src/internal/provider/shared/constants"
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
)

// Ensure ProjectsDataSource implements required interfaces.
var (
	_ datasource.DataSource              = &ProjectsDataSource{}
	_ datasource.DataSourceWithConfigure = &ProjectsDataSource{}
)

// ProjectsDataSource is a Terraform datasource that provides read-only access to all n8n projects.
// It enables querying and iterating through all available projects from the n8n API.
type ProjectsDataSource struct {
	client *client.N8nClient
}

// ProjectsDataSourceModel maps the Terraform schema to the datasource response.
// It represents a list of projects retrieved from the n8n API with all project details.
type ProjectsDataSourceModel struct {
	Projects []ProjectItemModel `tfsdk:"projects"`
}

// ProjectItemModel represents a single project in the list returned from the n8n API.
// It contains project metadata including name, type, timestamps, and descriptive information.
type ProjectItemModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Type        types.String `tfsdk:"type"`
	CreatedAt   types.String `tfsdk:"created_at"`
	UpdatedAt   types.String `tfsdk:"updated_at"`
	Icon        types.String `tfsdk:"icon"`
	Description types.String `tfsdk:"description"`
}

// NewProjectsDataSource creates a new ProjectsDataSource instance.
//
// Returns:
//   - datasource.DataSource: a new ProjectsDataSource instance configured for accessing n8n projects
func NewProjectsDataSource() datasource.DataSource {
	// Return result.
	return &ProjectsDataSource{}
}

// Metadata returns the data source type name.
//
// Params:
//   - ctx: context for the request
//   - req: metadata request containing provider type name
//   - resp: metadata response to populate
//
// Returns:
func (d *ProjectsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_projects"
}

// Schema defines the schema for the data source.
//
// Params:
//   - ctx: context for the request
//   - req: schema request from the framework
//   - resp: schema response to populate with schema definition
//
// Returns:
func (d *ProjectsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches a list of all n8n projects",

		Attributes: map[string]schema.Attribute{
			"projects": schema.ListNestedAttribute{
				MarkdownDescription: "List of projects",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Project identifier",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Project name",
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
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
//
// Params:
//   - ctx: context for the request
//   - req: configure request containing provider data
//   - resp: configure response to report errors
//
// Returns:
func (d *ProjectsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
//   - ctx: context for the request
//   - req: read request from Terraform
//   - resp: read response to populate with data
//
// Returns:
func (d *ProjectsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ProjectsDataSourceModel

	projectList, httpResp, err := d.client.APIClient.ProjectsAPI.ProjectsGet(ctx).Execute()
	// Close the HTTP response body if it is not nil to prevent resource leaks.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}
	// Check for error.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error listing projects",
			fmt.Sprintf("Could not list projects: %s\nHTTP Response: %v", err.Error(), httpResp),
		)
		// Return result.
		return
	}

	data.Projects = make([]ProjectItemModel, 0, constants.DefaultListCapacity)
	// Iterate through all projects and convert them to the model format.
	if projectList.Data != nil {
		// Convert each project from the API response to the ProjectItemModel format.
		for _, project := range projectList.Data {
			item := mapProjectToProjectItemModel(&project)
			data.Projects = append(data.Projects, item)
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
