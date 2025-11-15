package user

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUsersDataSourceInterface is a mock implementation of UsersDataSourceInterface.
type MockUsersDataSourceInterface struct {
	mock.Mock
}

func (m *MockUsersDataSourceInterface) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockUsersDataSourceInterface) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockUsersDataSourceInterface) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockUsersDataSourceInterface) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	m.Called(ctx, req, resp)
}

// TestNewUsersDataSource is now in external test file - refactored to test behavior only.

// TestNewUsersDataSourceWrapper is now in external test file - refactored to test behavior only.

// TestUsersDataSource_Metadata is now in external test file - refactored to test behavior only.

// TestUsersDataSource_Schema is now in external test file - refactored to test behavior only.

// TestUsersDataSource_Configure is now in external test file - refactored to test behavior only.

func TestUsersDataSource_Interfaces(t *testing.T) {
	t.Helper()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "implements required interfaces",
			testFunc: func(t *testing.T) {
				t.Helper()

				ds := NewUsersDataSource()

				// Test that UsersDataSource implements datasource.DataSource
				var _ datasource.DataSource = ds

				// Test that UsersDataSource implements datasource.DataSourceWithConfigure
				var _ datasource.DataSourceWithConfigure = ds

				// Test that UsersDataSource implements UsersDataSourceInterface
				var _ UsersDataSourceInterface = ds
			},
		},
		{
			name: "interface implementation error case",
			testFunc: func(t *testing.T) {
				t.Helper()

				ds := NewUsersDataSource()

				// This test ensures the type assertions don't panic
				assert.NotNil(t, ds)
				var _ datasource.DataSource = ds
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t)
		})
	}
}

func TestUsersDataSourceConcurrency(t *testing.T) {
	t.Helper()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "concurrent metadata calls",
			testFunc: func(t *testing.T) {
				t.Helper()

				ds := NewUsersDataSource()

				done := make(chan bool, 100)
				for i := 0; i < 100; i++ {
					go func() {
						resp := &datasource.MetadataResponse{}
						ds.Metadata(context.Background(), datasource.MetadataRequest{
							ProviderTypeName: "n8n",
						}, resp)
						assert.Equal(t, "n8n_users", resp.TypeName)
						done <- true
					}()
				}

				for i := 0; i < 100; i++ {
					<-done
				}
			},
		},
		{
			name: "concurrent schema calls",
			testFunc: func(t *testing.T) {
				t.Helper()

				ds := NewUsersDataSource()

				done := make(chan bool, 100)
				for i := 0; i < 100; i++ {
					go func() {
						resp := &datasource.SchemaResponse{}
						ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)
						assert.NotNil(t, resp.Schema)
						done <- true
					}()
				}

				for i := 0; i < 100; i++ {
					<-done
				}
			},
		},
		{
			name: "concurrent configure calls error handling",
			testFunc: func(t *testing.T) {
				t.Helper()

				ds := NewUsersDataSource()

				done := make(chan bool, 50)
				for i := 0; i < 50; i++ {
					go func() {
						resp := &datasource.ConfigureResponse{}
						req := datasource.ConfigureRequest{
							ProviderData: "invalid",
						}
						ds.Configure(context.Background(), req, resp)
						assert.True(t, resp.Diagnostics.HasError())
						done <- true
					}()
				}

				for i := 0; i < 50; i++ {
					<-done
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t)
		})
	}
}

func BenchmarkUsersDataSource_Schema(b *testing.B) {
	ds := NewUsersDataSource()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := &datasource.SchemaResponse{}
		ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)
	}
}

func BenchmarkUsersDataSource_Metadata(b *testing.B) {
	ds := NewUsersDataSource()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := &datasource.MetadataResponse{}
		ds.Metadata(context.Background(), datasource.MetadataRequest{}, resp)
	}
}

func BenchmarkUsersDataSource_Configure(b *testing.B) {
	ds := NewUsersDataSource()
	mockClient := &client.N8nClient{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := &datasource.ConfigureResponse{}
		req := datasource.ConfigureRequest{
			ProviderData: mockClient,
		}
		ds.Configure(context.Background(), req, resp)
	}
}

// TestUsersDataSource_Read tests the read functionality.
// TestUsersDataSource_Read is now in external test file - refactored to test behavior only.

// TestUsersDataSource_schemaAttributes tests the private schemaAttributes method.
func TestUsersDataSource_schemaAttributes(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "returns non-nil and non-empty attributes",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := NewUsersDataSource()
				attrs := ds.schemaAttributes()
				assert.NotNil(t, attrs, "schemaAttributes should return non-nil attributes")
				assert.NotEmpty(t, attrs, "schemaAttributes should return non-empty attributes")
			},
		},
		{
			name: "error case - multiple calls return consistent results",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := NewUsersDataSource()
				attrs1 := ds.schemaAttributes()
				attrs2 := ds.schemaAttributes()
				assert.Equal(t, len(attrs1), len(attrs2), "Multiple calls should return same number of attributes")
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

// TestUsersDataSource_userItemAttributes tests the private userItemAttributes method.
func TestUsersDataSource_userItemAttributes(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "returns non-nil and non-empty attributes",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := NewUsersDataSource()
				attrs := ds.userItemAttributes()
				assert.NotNil(t, attrs, "userItemAttributes should return non-nil attributes")
				assert.NotEmpty(t, attrs, "userItemAttributes should return non-empty attributes")
			},
		},
		{
			name: "error case - multiple calls return consistent results",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := NewUsersDataSource()
				attrs1 := ds.userItemAttributes()
				attrs2 := ds.userItemAttributes()
				assert.Equal(t, len(attrs1), len(attrs2), "Multiple calls should return same number of attributes")
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
