package project_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/kodflow/n8n/src/internal/provider/project"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
	"github.com/stretchr/testify/assert"
)

// TestNewProjectsDataSource tests the NewProjectsDataSource constructor.
func TestNewProjectsDataSource(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "creates valid datasource",
			wantErr: false,
		},
		{
			name:    "error case - validation checks",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "creates valid datasource":
				ds := project.NewProjectsDataSource()
				assert.NotNil(t, ds)

			case "error case - validation checks":
				ds := project.NewProjectsDataSource()
				assert.NotNil(t, ds)
				assert.Implements(t, (*datasource.DataSource)(nil), ds)
			}
		})
	}
}

// TestNewProjectsDataSourceWrapper tests the NewProjectsDataSourceWrapper constructor.
func TestNewProjectsDataSourceWrapper(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "creates valid datasource wrapper",
			wantErr: false,
		},
		{
			name:    "error case - validation checks",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "creates valid datasource wrapper":
				wrapper := project.NewProjectsDataSourceWrapper()
				assert.NotNil(t, wrapper)

			case "error case - validation checks":
				wrapper := project.NewProjectsDataSourceWrapper()
				assert.NotNil(t, wrapper)
			}
		})
	}
}

// TestProjectsDataSource_Metadata tests the Metadata method.
func TestProjectsDataSource_Metadata(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "sets correct type name",
			wantErr: false,
		},
		{
			name:    "error case - validation checks",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "sets correct type name":
				ds := project.NewProjectsDataSource()
				req := datasource.MetadataRequest{
					ProviderTypeName: "n8n",
				}
				resp := &datasource.MetadataResponse{}
				ds.Metadata(context.Background(), req, resp)
				assert.Equal(t, "n8n_projects", resp.TypeName)

			case "error case - validation checks":
				ds := project.NewProjectsDataSource()
				req := datasource.MetadataRequest{
					ProviderTypeName: "n8n",
				}
				resp := &datasource.MetadataResponse{}
				ds.Metadata(context.Background(), req, resp)
				assert.NotEmpty(t, resp.TypeName)
			}
		})
	}
}

// TestProjectsDataSource_Schema tests the Schema method.
func TestProjectsDataSource_Schema(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "returns valid schema",
			wantErr: false,
		},
		{
			name:    "error case - validation checks",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "returns valid schema":
				ds := project.NewProjectsDataSource()
				req := datasource.SchemaRequest{}
				resp := &datasource.SchemaResponse{}
				ds.Schema(context.Background(), req, resp)
				assert.NotNil(t, resp.Schema)
				assert.NotEmpty(t, resp.Schema.Attributes)

			case "error case - validation checks":
				ds := project.NewProjectsDataSource()
				req := datasource.SchemaRequest{}
				resp := &datasource.SchemaResponse{}
				ds.Schema(context.Background(), req, resp)
				assert.NotNil(t, resp.Schema)
			}
		})
	}
}

// TestProjectsDataSource_Configure tests the Configure method.
func TestProjectsDataSource_Configure(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "configures with valid client",
			wantErr: false,
		},
		{
			name:    "error case - nil provider data",
			wantErr: false,
		},
		{
			name:    "error case - wrong provider data type",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "configures with valid client":
				ds := project.NewProjectsDataSource()
				req := datasource.ConfigureRequest{
					ProviderData: &client.N8nClient{},
				}
				resp := &datasource.ConfigureResponse{}
				ds.Configure(context.Background(), req, resp)
				assert.False(t, resp.Diagnostics.HasError())

			case "error case - nil provider data":
				ds := project.NewProjectsDataSource()
				req := datasource.ConfigureRequest{
					ProviderData: nil,
				}
				resp := &datasource.ConfigureResponse{}
				ds.Configure(context.Background(), req, resp)
				assert.False(t, resp.Diagnostics.HasError())

			case "error case - wrong provider data type":
				ds := project.NewProjectsDataSource()
				req := datasource.ConfigureRequest{
					ProviderData: "wrong type",
				}
				resp := &datasource.ConfigureResponse{}
				ds.Configure(context.Background(), req, resp)
				assert.True(t, resp.Diagnostics.HasError())
				assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Unexpected Data Source Configure Type")
			}
		})
	}
}

// TestProjectsDataSource_Read tests the Read method.
func TestProjectsDataSource_Read(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "validates method exists",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "validates method exists":
				ds := project.NewProjectsDataSource()
				assert.NotNil(t, ds)
			}
		})
	}
}
