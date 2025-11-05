package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

// DataSources returns all data sources supported by the provider.
// Each data source factory function creates a new instance of the data source.
func DataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewWorkflowDataSource,
		NewWorkflowsDataSource,
		// TODO: Add more data sources as needed
		// NewCredentialDataSource,
		// NewCredentialsDataSource,
		// NewTagDataSource,
		// NewTagsDataSource,
		// NewVariableDataSource,
		// NewVariablesDataSource,
		// NewProjectDataSource,
		// NewProjectsDataSource,
	}
}
