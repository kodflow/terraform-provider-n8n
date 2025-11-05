package datasources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	providertypes "github.com/kodflow/n8n/src/internal/provider/types"
)

// Ensure ProjectsDataSource implements required interfaces.
var _ datasource.DataSource = &ProjectsDataSource{}
var _ datasource.DataSourceWithConfigure = &ProjectsDataSource{}

// ProjectsDataSource defines the data source implementation for listing projects.
type ProjectsDataSource struct {
	client *providertypes.N8nClient
}

// ProjectsDataSourceModel describes the data source data model.
type ProjectsDataSourceModel struct {
	Projects []ProjectItemModel `tfsdk:"projects"`
}

// ProjectItemModel represents a single project in the list.
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
func NewProjectsDataSource() datasource.DataSource {
	return &ProjectsDataSource{}
}

// Metadata returns the data source type name.
func (d *ProjectsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_projects"
}

// Schema defines the schema for the data source.
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
func (d *ProjectsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*providertypes.N8nClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *providertypes.N8nClient, got: %T", req.ProviderData),
		)
		return
	}

	d.client = client
}

// Read refreshes the Terraform state with the latest data.
func (d *ProjectsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ProjectsDataSourceModel

	projectList, httpResp, err := d.client.APIClient.ProjectsAPI.ProjectsGet(ctx).Execute()
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Error listing projects",
			fmt.Sprintf("Could not list projects: %s\nHTTP Response: %v", err.Error(), httpResp),
		)
		return
	}

	data.Projects = make([]ProjectItemModel, 0)
	if projectList.Data != nil {
		for _, project := range projectList.Data {
			item := ProjectItemModel{
				Name: types.StringValue(project.Name),
			}
			if project.Id != nil {
				item.ID = types.StringValue(*project.Id)
			}
			if project.Type != nil {
				item.Type = types.StringPointerValue(project.Type)
			}
			if project.CreatedAt != nil {
				item.CreatedAt = types.StringValue(project.CreatedAt.String())
			}
			if project.UpdatedAt != nil {
				item.UpdatedAt = types.StringValue(project.UpdatedAt.String())
			}
			if project.Icon != nil {
				item.Icon = types.StringPointerValue(project.Icon)
			}
			if project.Description != nil {
				item.Description = types.StringPointerValue(project.Description)
			}
			data.Projects = append(data.Projects, item)
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
