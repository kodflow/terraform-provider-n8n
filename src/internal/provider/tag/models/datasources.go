package models

// DataSources maps Terraform schema attributes for tag list data.
// It represents the complete data structure returned from the n8n tags API.
type DataSources struct {
	Tags []Item `tfsdk:"tags"`
}
