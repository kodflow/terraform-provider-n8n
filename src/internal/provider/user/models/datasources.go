package models

// DataSources maps Terraform schema attributes for user list data.
// It represents the complete data structure returned from the n8n users API.
type DataSources struct {
	Users []Item `tfsdk:"users"`
}
