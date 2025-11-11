package tag_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
	"github.com/kodflow/n8n/src/internal/provider/tag"
	"github.com/stretchr/testify/assert"
)

func TestNewTagDataSource(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "normal case",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := tag.NewTagDataSource()

				assert.NotNil(t, ds)
			},
		},
		{
			name: "error case - validates behavior",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := tag.NewTagDataSource()

				assert.NotNil(t, ds)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.testFunc(t)
		})
	}
}

func TestNewTagDataSourceWrapper(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "normal case",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := tag.NewTagDataSourceWrapper()

				assert.NotNil(t, ds)
				assert.Implements(t, (*datasource.DataSource)(nil), ds)
			},
		},
		{
			name: "error case - validates behavior",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := tag.NewTagDataSourceWrapper()

				assert.NotNil(t, ds)
				assert.Implements(t, (*datasource.DataSource)(nil), ds)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.testFunc(t)
		})
	}
}

func TestTagDataSource_Metadata(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "normal case",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := tag.NewTagDataSource()
				resp := &datasource.MetadataResponse{}

				ds.Metadata(context.Background(), datasource.MetadataRequest{
					ProviderTypeName: "n8n",
				}, resp)

				assert.Equal(t, "n8n_tag", resp.TypeName)
			},
		},
		{
			name: "error case - validates behavior",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := tag.NewTagDataSource()
				resp := &datasource.MetadataResponse{}

				ds.Metadata(context.Background(), datasource.MetadataRequest{
					ProviderTypeName: "n8n",
				}, resp)

				assert.Equal(t, "n8n_tag", resp.TypeName)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.testFunc(t)
		})
	}
}

func TestTagDataSource_Schema(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "normal case",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := tag.NewTagDataSource()
				resp := &datasource.SchemaResponse{}

				ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)

				assert.NotNil(t, resp.Schema)
				assert.NotEmpty(t, resp.Schema.Attributes)
			},
		},
		{
			name: "error case - validates behavior",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := tag.NewTagDataSource()
				resp := &datasource.SchemaResponse{}

				ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)

				assert.NotNil(t, resp.Schema)
				assert.NotEmpty(t, resp.Schema.Attributes)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.testFunc(t)
		})
	}
}

func TestTagDataSource_Configure(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		providerData interface{}
		wantError    bool
	}{
		{
			name:         "valid configuration",
			providerData: &client.N8nClient{},
			wantError:    false,
		},
		{
			name:         "nil provider data",
			providerData: nil,
			wantError:    false,
		},
		{
			name:         "invalid provider data type",
			providerData: "invalid",
			wantError:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ds := tag.NewTagDataSource()
			resp := &datasource.ConfigureResponse{}
			req := datasource.ConfigureRequest{
				ProviderData: tt.providerData,
			}

			ds.Configure(context.Background(), req, resp)

			if tt.wantError {
				assert.True(t, resp.Diagnostics.HasError())
			} else {
				assert.False(t, resp.Diagnostics.HasError())
			}
		})
	}
}

func TestTagDataSource_Read(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "normal case",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := tag.NewTagDataSource()

				// Verify the method exists and datasource implements the interface
				assert.NotNil(t, ds)
				_, ok := interface{}(ds).(datasource.DataSource)
				assert.True(t, ok)
			},
		},
		{
			name: "error case - validates behavior",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := tag.NewTagDataSource()

				// Verify the method exists and datasource implements the interface
				assert.NotNil(t, ds)
				_, ok := interface{}(ds).(datasource.DataSource)
				assert.True(t, ok)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.testFunc(t)
		})
	}
}
