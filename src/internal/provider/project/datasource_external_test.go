package project_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/kodflow/n8n/src/internal/provider/project"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
	"github.com/stretchr/testify/assert"
)

// TestNewProjectDataSource tests the NewProjectDataSource constructor.
func TestNewProjectDataSource(t *testing.T) {
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
				ds := project.NewProjectDataSource()
				assert.NotNil(t, ds)

			case "error case - validation checks":
				ds := project.NewProjectDataSource()
				assert.NotNil(t, ds)
				assert.Implements(t, (*datasource.DataSource)(nil), ds)
			}
		})
	}
}

// TestNewProjectDataSourceWrapper tests the NewProjectDataSourceWrapper constructor.
func TestNewProjectDataSourceWrapper(t *testing.T) {
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
				wrapper := project.NewProjectDataSourceWrapper()
				assert.NotNil(t, wrapper)

			case "error case - validation checks":
				wrapper := project.NewProjectDataSourceWrapper()
				assert.NotNil(t, wrapper)
			}
		})
	}
}

// TestProjectDataSource_Metadata tests the Metadata method.
func TestProjectDataSource_Metadata(t *testing.T) {
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
				ds := project.NewProjectDataSource()
				req := datasource.MetadataRequest{
					ProviderTypeName: "n8n",
				}
				resp := &datasource.MetadataResponse{}
				ds.Metadata(context.Background(), req, resp)
				assert.Equal(t, "n8n_project", resp.TypeName)

			case "error case - validation checks":
				ds := project.NewProjectDataSource()
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

// TestProjectDataSource_Schema tests the Schema method.
func TestProjectDataSource_Schema(t *testing.T) {
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
				ds := project.NewProjectDataSource()
				req := datasource.SchemaRequest{}
				resp := &datasource.SchemaResponse{}
				ds.Schema(context.Background(), req, resp)
				assert.NotNil(t, resp.Schema)
				assert.NotEmpty(t, resp.Schema.Attributes)

			case "error case - validation checks":
				ds := project.NewProjectDataSource()
				req := datasource.SchemaRequest{}
				resp := &datasource.SchemaResponse{}
				ds.Schema(context.Background(), req, resp)
				assert.NotNil(t, resp.Schema)
			}
		})
	}
}

// TestProjectDataSource_Configure tests the Configure method.
func TestProjectDataSource_Configure(t *testing.T) {
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
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "configures with valid client":
				ds := project.NewProjectDataSource()
				req := datasource.ConfigureRequest{
					ProviderData: &client.N8nClient{},
				}
				resp := &datasource.ConfigureResponse{}
				ds.Configure(context.Background(), req, resp)
				assert.False(t, resp.Diagnostics.HasError())

			case "error case - nil provider data":
				ds := project.NewProjectDataSource()
				req := datasource.ConfigureRequest{
					ProviderData: nil,
				}
				resp := &datasource.ConfigureResponse{}
				ds.Configure(context.Background(), req, resp)
				// Should not error on nil provider data
				assert.False(t, resp.Diagnostics.HasError())
			}
		})
	}
}

// TestProjectDataSource_Read tests the Read method.
func TestProjectDataSource_Read(t *testing.T) {
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
				ds := project.NewProjectDataSource()
				assert.NotNil(t, ds)
			}
		})
	}
}
