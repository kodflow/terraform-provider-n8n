package user

// UsersDataSourceModel maps Terraform schema attributes for user list data.
// It represents the complete data structure returned from the n8n users API.
type UsersDataSourceModel struct {
	Users []UserItemModel `tfsdk:"users"`
}
