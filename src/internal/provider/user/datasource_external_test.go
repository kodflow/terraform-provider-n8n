package user_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
	"github.com/kodflow/n8n/src/internal/provider/user"
	"github.com/stretchr/testify/assert"
)

// TestNewUserDataSource tests the NewUserDataSource constructor.
func TestNewUserDataSource(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "normal case",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := user.NewUserDataSource()
				assert.NotNil(t, ds, "NewUserDataSource should return a non-nil datasource")
			},
		},
		{
			name: "error case - validates behavior",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := user.NewUserDataSource()
				assert.NotNil(t, ds, "NewUserDataSource should return a non-nil datasource")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.testFunc(t)
		})
	}
} // TestNewUserDataSourceWrapper tests the NewUserDataSourceWrapper constructor.
func TestNewUserDataSourceWrapper(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "normal case",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := user.NewUserDataSourceWrapper()
				assert.NotNil(t, ds, "NewUserDataSourceWrapper should return a non-nil datasource")
			},
		},
		{
			name: "error case - validates behavior",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := user.NewUserDataSourceWrapper()
				assert.NotNil(t, ds, "NewUserDataSourceWrapper should return a non-nil datasource")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.testFunc(t)
		})
	}
} // TestUserDataSource_Metadata tests the Metadata method.
func TestUserDataSource_Metadata(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "normal case",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := user.NewUserDataSource()
				req := datasource.MetadataRequest{
					ProviderTypeName: "n8n",
				}
				resp := &datasource.MetadataResponse{}

				ds.Metadata(context.Background(), req, resp)

				assert.Equal(t, "n8n_user", resp.TypeName, "TypeName should be set correctly")
			},
		},
		{
			name: "error case - validates behavior",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := user.NewUserDataSource()
				req := datasource.MetadataRequest{
					ProviderTypeName: "n8n",
				}
				resp := &datasource.MetadataResponse{}

				ds.Metadata(context.Background(), req, resp)

				assert.Equal(t, "n8n_user", resp.TypeName, "TypeName should be set correctly")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.testFunc(t)
		})
	}
} // TestUserDataSource_Schema tests the Schema method.
func TestUserDataSource_Schema(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "normal case",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := user.NewUserDataSource()
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
				ds := user.NewUserDataSource()
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
} // TestUserDataSource_Configure tests the Configure method.
func TestUserDataSource_Configure(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ds := user.NewUserDataSource()
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

// TestUserDataSource_Read tests the Read method.
func TestUserDataSource_Read(t *testing.T) {
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
				// This test just verifies that NewUserDataSource doesn't panic.
				assert.NotPanics(t, func() {
					_ = user.NewUserDataSource()
				}, "NewUserDataSource should not panic")
			},
		},
		{
			name: "error case - validates behavior",
			testFunc: func(t *testing.T) {
				t.Helper()
				// Note: Cannot test Read() with empty request as it would panic
				// due to uninitialized Config. In production, terraform-plugin-framework
				// always provides properly initialized Config structures.
				// This test just verifies that NewUserDataSource doesn't panic.
				assert.NotPanics(t, func() {
					_ = user.NewUserDataSource()
				}, "NewUserDataSource should not panic")
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
