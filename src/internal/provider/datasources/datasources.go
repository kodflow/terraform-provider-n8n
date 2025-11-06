package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

// DataSources returns all data sources supported by the provider.
// Each data source factory function creates a new instance of the data source.
func DataSources() []func() datasource.DataSource {
	// Return result.
	return []func() datasource.DataSource{
		NewWorkflowDataSource,
		NewWorkflowsDataSource,
		NewProjectDataSource,
		NewProjectsDataSource,
		NewVariableDataSource,
		NewVariablesDataSource,
		NewTagDataSource,
		NewTagsDataSource,
		NewUserDataSource,
		NewUsersDataSource,
		NewExecutionDataSource,
		NewExecutionsDataSource,
		// TODO: Add more data sources as needed
		// NewCredentialDataSource,
		// NewCredentialsDataSource,
	}
}
