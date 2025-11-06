package datasources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	providertypes "github.com/kodflow/n8n/src/internal/provider/types"
)

// Ensure UserDataSource implements required interfaces.
var (
	_ datasource.DataSource              = &UserDataSource{}
	_ datasource.DataSourceWithConfigure = &UserDataSource{}
)

// UserDataSource is a Terraform datasource implementation for retrieving a single user.
// It provides read-only access to n8n user details through the n8n API.
type UserDataSource struct {
	client *providertypes.N8nClient
}

// UserDataSourceModel maps Terraform schema attributes for user data.
// It represents a single user with all related attributes from the n8n API.
type UserDataSourceModel struct {
	ID        types.String `tfsdk:"id"`
	Email     types.String `tfsdk:"email"`
	FirstName types.String `tfsdk:"first_name"`
	LastName  types.String `tfsdk:"last_name"`
	IsPending types.Bool   `tfsdk:"is_pending"`
	CreatedAt types.String `tfsdk:"created_at"`
	UpdatedAt types.String `tfsdk:"updated_at"`
	Role      types.String `tfsdk:"role"`
}

// NewUserDataSource creates a new UserDataSource instance.
func NewUserDataSource() datasource.DataSource {
	// Return result.
	return &UserDataSource{}
}

// Metadata returns the data source type name.
func (d *UserDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

// Schema defines the schema for the data source.
func (d *UserDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches a single n8n user by ID or email. The API accepts both ID and email as identifiers.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "User identifier. Either `id` or `email` must be specified. Only available for instance owners.",
				Optional:            true,
				Computed:            true,
			},
			"email": schema.StringAttribute{
				MarkdownDescription: "User email address. Either `id` or `email` must be specified.",
				Optional:            true,
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
				MarkdownDescription: "Whether the user finished setting up their account in response to the invitation (false) or not (true)",
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
				MarkdownDescription: "User's global role (e.g., 'global:owner', 'global:admin', 'global:member')",
				Computed:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *UserDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *UserDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data UserDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	// Check condition.
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate that at least one identifier is provided
	if data.ID.IsNull() && data.Email.IsNull() {
		resp.Diagnostics.AddError(
			"Missing Required Attribute",
			"Either 'id' or 'email' must be specified",
		)
		// Return result.
		return
	}

	// Use ID if provided, otherwise use email
	identifier := data.ID.ValueString()
	// Check condition.
	if identifier == "" {
		identifier = data.Email.ValueString()
	}

	user, httpResp, err := d.client.APIClient.UserAPI.UsersIdGet(ctx, identifier).Execute()
	// Check for non-nil value.
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}
	// Check for error.
	if err != nil {
		resp.Diagnostics.AddError(
			"Error retrieving user",
			fmt.Sprintf("Could not retrieve user with identifier %s: %s\nHTTP Response: %v", identifier, err.Error(), httpResp),
		)
		// Return result.
		return
	}

	// Populate data model
	if user.Id != nil {
		data.ID = types.StringValue(*user.Id)
	}
	data.Email = types.StringValue(user.Email)
	// Check for non-nil value.
	if user.FirstName != nil {
		data.FirstName = types.StringPointerValue(user.FirstName)
	}
	// Check for non-nil value.
	if user.LastName != nil {
		data.LastName = types.StringPointerValue(user.LastName)
	}
	// Check for non-nil value.
	if user.IsPending != nil {
		data.IsPending = types.BoolPointerValue(user.IsPending)
	}
	// Check for non-nil value.
	if user.CreatedAt != nil {
		data.CreatedAt = types.StringValue(user.CreatedAt.String())
	}
	// Check for non-nil value.
	if user.UpdatedAt != nil {
		data.UpdatedAt = types.StringValue(user.UpdatedAt.String())
	}
	// Check for non-nil value.
	if user.Role != nil {
		data.Role = types.StringPointerValue(user.Role)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
