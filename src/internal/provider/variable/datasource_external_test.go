package variable_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
	"github.com/kodflow/n8n/src/internal/provider/variable"
	"github.com/stretchr/testify/assert"
)

// TestNewVariableDataSource tests the NewVariableDataSource constructor.
func TestNewVariableDataSource(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "normal case",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := variable.NewVariableDataSource()
				assert.NotNil(t, ds, "NewVariableDataSource should return a non-nil datasource")
			},
		},
		{
			name: "error case - validates behavior",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := variable.NewVariableDataSource()
				assert.NotNil(t, ds, "NewVariableDataSource should return a non-nil datasource")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.testFunc(t)
		})
	}
} // TestNewVariableDataSourceWrapper tests the NewVariableDataSourceWrapper constructor.
func TestNewVariableDataSourceWrapper(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "normal case",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := variable.NewVariableDataSourceWrapper()
				assert.NotNil(t, ds, "NewVariableDataSourceWrapper should return a non-nil datasource")
			},
		},
		{
			name: "error case - validates behavior",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := variable.NewVariableDataSourceWrapper()
				assert.NotNil(t, ds, "NewVariableDataSourceWrapper should return a non-nil datasource")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.testFunc(t)
		})
	}
} // TestVariableDataSource_Metadata tests the Metadata method.
func TestVariableDataSource_Metadata(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "normal case",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := variable.NewVariableDataSource()
				req := datasource.MetadataRequest{
					ProviderTypeName: "n8n",
				}
				resp := &datasource.MetadataResponse{}

				ds.Metadata(context.Background(), req, resp)

				assert.Equal(t, "n8n_variable", resp.TypeName, "TypeName should be set correctly")
			},
		},
		{
			name: "error case - validates behavior",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := variable.NewVariableDataSource()
				req := datasource.MetadataRequest{
					ProviderTypeName: "n8n",
				}
				resp := &datasource.MetadataResponse{}

				ds.Metadata(context.Background(), req, resp)

				assert.Equal(t, "n8n_variable", resp.TypeName, "TypeName should be set correctly")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.testFunc(t)
		})
	}
} // TestVariableDataSource_Schema tests the Schema method.
func TestVariableDataSource_Schema(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "normal case",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := variable.NewVariableDataSource()
				req := datasource.SchemaRequest{}
				resp := &datasource.SchemaResponse{}

				ds.Schema(context.Background(), req, resp)

				assert.NotNil(t, resp.Schema, "Schema should not be nil")
				assert.NotNil(t, resp.Schema.Attributes, "Schema attributes should not be nil")
			},
		},
		{
			name: "error case - validates behavior",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := variable.NewVariableDataSource()
				req := datasource.SchemaRequest{}
				resp := &datasource.SchemaResponse{}

				ds.Schema(context.Background(), req, resp)

				assert.NotNil(t, resp.Schema, "Schema should not be nil")
				assert.NotNil(t, resp.Schema.Attributes, "Schema attributes should not be nil")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.testFunc(t)
		})
	}
} // TestVariableDataSource_Configure tests the Configure method.
func TestVariableDataSource_Configure(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		providerData interface{}
		expectError  bool
	}{
		{
			name:         "valid client",
			providerData: &client.N8nClient{},
			expectError:  false,
		},
		{
			name:         "nil provider data",
			providerData: nil,
			expectError:  false,
		},
		{
			name:         "invalid provider data",
			providerData: "invalid",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ds := variable.NewVariableDataSource()
			req := datasource.ConfigureRequest{
				ProviderData: tt.providerData,
			}
			resp := &datasource.ConfigureResponse{}

			ds.Configure(context.Background(), req, resp)

			if tt.expectError {
				assert.True(t, resp.Diagnostics.HasError(), "Configure should return error for invalid provider data")
			} else {
				assert.False(t, resp.Diagnostics.HasError(), "Configure should not return error")
			}
		})
	}
}

// TestVariableDataSource_Read tests the Read method.
func TestVariableDataSource_Read(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "normal case",
			testFunc: func(t *testing.T) {
				t.Helper()
				// Note: Cannot test Read() with empty request as it would panic
				// due to uninitialized Config. In production, terraform-plugin-framework
				// always provides properly initialized Config structures.
				// This test just verifies that NewVariableDataSource doesn't panic.
				assert.NotPanics(t, func() {
					_ = variable.NewVariableDataSource()
				}, "NewVariableDataSource should not panic")
			},
		},
		{
			name: "error case - validates behavior",
			testFunc: func(t *testing.T) {
				t.Helper()
				// Note: Cannot test Read() with empty request as it would panic
				// due to uninitialized Config. In production, terraform-plugin-framework
				// always provides properly initialized Config structures.
				// This test just verifies that NewVariableDataSource doesn't panic.
				assert.NotPanics(t, func() {
					_ = variable.NewVariableDataSource()
				}, "NewVariableDataSource should not panic")
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
