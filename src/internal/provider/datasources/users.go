package datasources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	providertypes "github.com/kodflow/n8n/src/internal/provider/types"
)

// Ensure UsersDataSource implements required interfaces.
var (
	_ datasource.DataSource              = &UsersDataSource{}
	_ datasource.DataSourceWithConfigure = &UsersDataSource{}
)

// UsersDataSource is a Terraform datasource implementation for listing users.
// It provides read-only access to all n8n users through the n8n API.
type UsersDataSource struct {
	client *providertypes.N8nClient
}

// UsersDataSourceModel maps Terraform schema attributes for user list data.
// It represents the complete data structure returned from the n8n users API.
type UsersDataSourceModel struct {
	Users []UserItemModel `tfsdk:"users"`
}

// UserItemModel represents a single user in the list.
// It maps individual user attributes from the n8n API to Terraform schema.
type UserItemModel struct {
	ID        types.String `tfsdk:"id"`
	Email     types.String `tfsdk:"email"`
	FirstName types.String `tfsdk:"first_name"`
	LastName  types.String `tfsdk:"last_name"`
	IsPending types.Bool   `tfsdk:"is_pending"`
	CreatedAt types.String `tfsdk:"created_at"`
	UpdatedAt types.String `tfsdk:"updated_at"`
	Role      types.String `tfsdk:"role"`
}

// NewUsersDataSource creates a new UsersDataSource instance.
func NewUsersDataSource() datasource.DataSource {
	// Return result.
	return &UsersDataSource{}
}

// Metadata returns the data source type name.
func (d *UsersDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_users"
}

// Schema defines the schema for the data source.
func (d *UsersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches a list of all n8n users. Only available for instance owners.",

		Attributes: map[string]schema.Attribute{
			"users": schema.ListNestedAttribute{
				MarkdownDescription: "List of users",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "User identifier",
							Computed:            true,
						},
						"email": schema.StringAttribute{
							MarkdownDescription: "User email address",
							Computed:            true,
						},
						"first_name": schema.StringAttribute{
							MarkdownDescription: "User's first name",
							Computed:            true,
						},
						"last_name": schema.StringAttribute{
							MarkdownDescription: "User's last name",
							Computed:            true,
						},
						"is_pending": schema.BoolAttribute{
							MarkdownDescription: "Whether the user finished setting up their account",
							Computed:            true,
						},
						"created_at": schema.StringAttribute{
							MarkdownDescription: "Timestamp when the user was created",
							Computed:            true,
						},
						"updated_at": schema.StringAttribute{
							MarkdownDescription: "Timestamp when the user was last updated",
							Computed:            true,
						},
						"role": schema.StringAttribute{
							MarkdownDescription: "User's global role",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *UsersDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *UsersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data UsersDataSourceModel

	userList, httpResp, err := d.client.APIClient.UserAPI.UsersGet(ctx).Execute()
	// Check for non-nil value.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}
	// Check for error.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error listing users",
			fmt.Sprintf("Could not list users: %s\nHTTP Response: %v", err.Error(), httpResp),
		)
		// Return result.
		return
	}

	data.Users = make([]UserItemModel, 0, DefaultListCapacity)
	// Check for non-nil value.
	if userList.Data != nil {
		// Iterate over items.
		for _, user := range userList.Data {
			item := UserItemModel{
				Email: types.StringValue(user.Email),
			}
			// Check for non-nil value.
			if user.Id != nil {
				item.ID = types.StringValue(*user.Id)
			}
			// Check for non-nil value.
			if user.FirstName != nil {
				item.FirstName = types.StringPointerValue(user.FirstName)
			}
			// Check for non-nil value.
			if user.LastName != nil {
				item.LastName = types.StringPointerValue(user.LastName)
			}
			// Check for non-nil value.
			if user.IsPending != nil {
				item.IsPending = types.BoolPointerValue(user.IsPending)
			}
			// Check for non-nil value.
			if user.CreatedAt != nil {
				item.CreatedAt = types.StringValue(user.CreatedAt.String())
			}
			// Check for non-nil value.
			if user.UpdatedAt != nil {
				item.UpdatedAt = types.StringValue(user.UpdatedAt.String())
			}
			// Check for non-nil value.
			if user.Role != nil {
				item.Role = types.StringPointerValue(user.Role)
			}
			data.Users = append(data.Users, item)
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
