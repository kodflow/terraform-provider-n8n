package datasources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	providertypes "github.com/kodflow/n8n/src/internal/provider/types"
)

// Ensure ProjectDataSource implements required interfaces.
var _ datasource.DataSource = &ProjectDataSource{}
var _ datasource.DataSourceWithConfigure = &ProjectDataSource{}

// ProjectDataSource defines the data source implementation for a single project.
type ProjectDataSource struct {
	client *providertypes.N8nClient
}

// ProjectDataSourceModel describes the data source data model.
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
func NewProjectDataSource() datasource.DataSource {
	return &ProjectDataSource{}
}

// Metadata returns the data source type name.
func (d *ProjectDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

// Schema defines the schema for the data source.
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
func (d *ProjectDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *ProjectDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ProjectDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate that at least one identifier is provided
	if data.ID.IsNull() && data.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Missing Required Attribute",
			"Either 'id' or 'name' must be specified",
		)
		return
	}

	// List all projects and filter client-side (API limitation)
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

	// Filter by ID or name
	var found bool
	if projectList.Data != nil {
		for _, project := range projectList.Data {
			matchByID := !data.ID.IsNull() && project.Id != nil && *project.Id == data.ID.ValueString()
			matchByName := !data.Name.IsNull() && project.Name == data.Name.ValueString()

			if matchByID || matchByName {
				// Populate data model
				if project.Id != nil {
					data.ID = types.StringValue(*project.Id)
				}
				data.Name = types.StringValue(project.Name)
				if project.Type != nil {
					data.Type = types.StringPointerValue(project.Type)
				}
				if project.CreatedAt != nil {
					data.CreatedAt = types.StringValue(project.CreatedAt.String())
				}
				if project.UpdatedAt != nil {
					data.UpdatedAt = types.StringValue(project.UpdatedAt.String())
				}
				if project.Icon != nil {
					data.Icon = types.StringPointerValue(project.Icon)
				}
				if project.Description != nil {
					data.Description = types.StringPointerValue(project.Description)
				}
				found = true
				break
			}
		}
	}

	if !found {
		identifier := data.ID.ValueString()
		if identifier == "" {
			identifier = data.Name.ValueString()
		}
		resp.Diagnostics.AddError(
			"Project Not Found",
			fmt.Sprintf("Could not find project with identifier: %s", identifier),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
