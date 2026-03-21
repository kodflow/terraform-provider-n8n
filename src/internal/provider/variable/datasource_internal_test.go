package variable

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/client"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/variable/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockVariableDataSourceInterface is a mock implementation of VariableDataSourceInterface.
type MockVariableDataSourceInterface struct {
	mock.Mock
}

func (m *MockVariableDataSourceInterface) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockVariableDataSourceInterface) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockVariableDataSourceInterface) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockVariableDataSourceInterface) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	m.Called(ctx, req, resp)
}

// setupTestDataSourceClient creates a test N8nClient with httptest server.
func setupTestDataSourceClient(t *testing.T, handler http.HandlerFunc) (*client.N8nClient, *httptest.Server) {
	t.Helper()
	server := httptest.NewServer(handler)

	cfg := n8nsdk.NewConfiguration()
	cfg.Servers = n8nsdk.ServerConfigurations{
		{
			URL:         server.URL,
			Description: "Test server",
		},
	}
	cfg.HTTPClient = server.Client()
	cfg.AddDefaultHeader("X-N8N-API-KEY", "test-key")

	apiClient := n8nsdk.NewAPIClient(cfg)
	n8nClient := &client.N8nClient{
		APIClient: apiClient,
	}

	return n8nClient, server
}

// TestNewVariableDataSource is now in external test file - refactored to test behavior only.

// TestNewVariableDataSourceWrapper is now in external test file - refactored to test behavior only.

// TestVariableDataSource_Metadata is now in external test file - refactored to test behavior only.

// TestVariableDataSource_Schema is now in external test file - refactored to test behavior only.

// TestVariableDataSource_Configure is now in external test file - refactored to test behavior only.

// Compile-time interface satisfaction checks.
var _ datasource.DataSource = (*VariableDataSource)(nil)
var _ datasource.DataSourceWithConfigure = (*VariableDataSource)(nil)
var _ VariableDataSourceInterface = (*VariableDataSource)(nil)

func TestVariableDataSource_Interfaces(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "constructor returns non-nil instance",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := NewVariableDataSource()
				assert.NotNil(t, ds)
			},
		},
		{
			name: "error case - nil instance type compliance",
			testFunc: func(t *testing.T) {
				t.Helper()
				var ds *VariableDataSource
				assert.Nil(t, ds)
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

func TestVariableDataSource_ValidateIdentifiers(t *testing.T) {
	t.Parallel()

	t.Helper()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "valid with ID provided",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := NewVariableDataSource()
				assert.NotNil(t, ds)
			},
		},
		{
			name: "valid with key provided",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := NewVariableDataSource()
				assert.NotNil(t, ds)
			},
		},
		{
			name: "error case - validation requires either ID or key",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := NewVariableDataSource()
				assert.NotNil(t, ds)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t)
		})
	}
}

func TestVariableDataSourceConcurrency(t *testing.T) {
	t.Parallel()

	t.Helper()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "concurrent metadata calls",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := NewVariableDataSource()

				var wg sync.WaitGroup
				for range 100 {
					wg.Go(func() {
						resp := &datasource.MetadataResponse{}
						ds.Metadata(t.Context(), datasource.MetadataRequest{
							ProviderTypeName: "n8n",
						}, resp)
						assert.Equal(t, "n8n_variable", resp.TypeName)
					})
				}
				wg.Wait()
			},
		},
		{
			name: "concurrent schema calls",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := NewVariableDataSource()

				var wg sync.WaitGroup
				for range 100 {
					wg.Go(func() {
						resp := &datasource.SchemaResponse{}
						ds.Schema(t.Context(), datasource.SchemaRequest{}, resp)
						assert.NotNil(t, resp.Schema)
					})
				}
				wg.Wait()
			},
		},
		{
			name: "concurrent configure calls error handling",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := NewVariableDataSource()

				var wg sync.WaitGroup
				for range 50 {
					wg.Go(func() {
						resp := &datasource.ConfigureResponse{}
						req := datasource.ConfigureRequest{
							ProviderData: "invalid",
						}
						ds.Configure(t.Context(), req, resp)
						assert.True(t, resp.Diagnostics.HasError())
					})
				}
				wg.Wait()
			},
		},
		{
			name: "error case - concurrent mixed operations",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := NewVariableDataSource()

				var wg sync.WaitGroup
				for range 50 {
					wg.Go(func() {
						resp := &datasource.MetadataResponse{}
						ds.Metadata(t.Context(), datasource.MetadataRequest{
							ProviderTypeName: "n8n",
						}, resp)
					})
				}
				for range 50 {
					wg.Go(func() {
						resp := &datasource.SchemaResponse{}
						ds.Schema(t.Context(), datasource.SchemaRequest{}, resp)
					})
				}
				wg.Wait()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t)
		})
	}
}

func BenchmarkVariableDataSource_Schema(b *testing.B) {
	ds := NewVariableDataSource()

	for b.Loop() {
		resp := &datasource.SchemaResponse{}
		ds.Schema(b.Context(), datasource.SchemaRequest{}, resp)
	}
}

func BenchmarkVariableDataSource_Metadata(b *testing.B) {
	ds := NewVariableDataSource()

	for b.Loop() {
		resp := &datasource.MetadataResponse{}
		ds.Metadata(b.Context(), datasource.MetadataRequest{
			ProviderTypeName: "n8n",
		}, resp)
	}
}

func BenchmarkVariableDataSource_Configure(b *testing.B) {
	ds := NewVariableDataSource()
	mockClient := &client.N8nClient{}

	for b.Loop() {
		resp := &datasource.ConfigureResponse{}
		req := datasource.ConfigureRequest{
			ProviderData: mockClient,
		}
		ds.Configure(b.Context(), req, resp)
	}
}

// TestVariableDataSource_Read is tested via external tests and helper method tests.
// This ensures coverage of the Read method through the private methods it calls.

func TestVariableDataSource_validateIdentifiers(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "valid - id provided"},
		{name: "valid - key provided"},
		{name: "valid - both provided"},
		{name: "error - neither provided", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ds := &VariableDataSource{}
			resp := &datasource.ReadResponse{}

			switch tt.name {
			case "valid - id provided":
				data := &models.DataSource{
					ID:  types.StringValue("var-123"),
					Key: types.StringNull(),
				}

				result := ds.validateIdentifiers(data, resp)

				assert.True(t, result)
				assert.False(t, resp.Diagnostics.HasError())

			case "valid - key provided":
				data := &models.DataSource{
					ID:  types.StringNull(),
					Key: types.StringValue("test-key"),
				}

				result := ds.validateIdentifiers(data, resp)

				assert.True(t, result)
				assert.False(t, resp.Diagnostics.HasError())

			case "valid - both provided":
				data := &models.DataSource{
					ID:  types.StringValue("var-123"),
					Key: types.StringValue("test-key"),
				}

				result := ds.validateIdentifiers(data, resp)

				assert.True(t, result)
				assert.False(t, resp.Diagnostics.HasError())

			case "error - neither provided":
				data := &models.DataSource{
					ID:  types.StringNull(),
					Key: types.StringNull(),
				}

				result := ds.validateIdentifiers(data, resp)

				assert.False(t, result)
				assert.True(t, resp.Diagnostics.HasError())
				assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Missing Required Attribute")
			}
		})
	}
}

func TestVariableDataSource_fetchVariable(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "success - fetch by ID"},
		{name: "success - fetch by key"},
		{name: "success - with project filter"},
		{name: "error - API error", wantErr: true},
		{name: "error - variable not found", wantErr: true},
		{name: "error - variable not found with key identifier", wantErr: true},
		{name: "error - nil variable data", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "success - fetch by ID":
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/variables" && r.Method == http.MethodGet {
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusOK)
						json.NewEncoder(w).Encode(n8nsdk.VariableList{
							Data: []n8nsdk.Variable{
								{Id: ptrString("var-123"), Key: "test-key", Value: "test-value"},
							},
						})
						return
					}
					w.WriteHeader(http.StatusNotFound)
				})

				n8nClient, server := setupTestDataSourceClient(t, handler)
				defer server.Close()

				ds := &VariableDataSource{client: n8nClient}
				data := &models.DataSource{
					ID: types.StringValue("var-123"),
				}
				resp := &datasource.ReadResponse{}

				variable := ds.fetchVariable(t.Context(), data, resp)

				assert.NotNil(t, variable)
				assert.Equal(t, "var-123", *variable.Id)
				assert.False(t, resp.Diagnostics.HasError())

			case "success - fetch by key":
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/variables" && r.Method == http.MethodGet {
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusOK)
						json.NewEncoder(w).Encode(n8nsdk.VariableList{
							Data: []n8nsdk.Variable{
								{Id: ptrString("var-123"), Key: "test-key", Value: "test-value"},
							},
						})
						return
					}
					w.WriteHeader(http.StatusNotFound)
				})

				n8nClient, server := setupTestDataSourceClient(t, handler)
				defer server.Close()

				ds := &VariableDataSource{client: n8nClient}
				data := &models.DataSource{
					Key: types.StringValue("test-key"),
				}
				resp := &datasource.ReadResponse{}

				variable := ds.fetchVariable(t.Context(), data, resp)

				assert.NotNil(t, variable)
				assert.Equal(t, "test-key", variable.Key)
				assert.False(t, resp.Diagnostics.HasError())

			case "success - with project filter":
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/variables" && r.Method == http.MethodGet {
						assert.Contains(t, r.URL.RawQuery, "projectId=proj-123")
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusOK)
						json.NewEncoder(w).Encode(n8nsdk.VariableList{
							Data: []n8nsdk.Variable{
								{Id: ptrString("var-123"), Key: "test-key", Value: "test-value"},
							},
						})
						return
					}
					w.WriteHeader(http.StatusNotFound)
				})

				n8nClient, server := setupTestDataSourceClient(t, handler)
				defer server.Close()

				ds := &VariableDataSource{client: n8nClient}
				data := &models.DataSource{
					ID:        types.StringValue("var-123"),
					ProjectID: types.StringValue("proj-123"),
				}
				resp := &datasource.ReadResponse{}

				variable := ds.fetchVariable(t.Context(), data, resp)

				assert.NotNil(t, variable)
				assert.False(t, resp.Diagnostics.HasError())

			case "error - API error":
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
				})

				n8nClient, server := setupTestDataSourceClient(t, handler)
				defer server.Close()

				ds := &VariableDataSource{client: n8nClient}
				data := &models.DataSource{
					ID: types.StringValue("var-123"),
				}
				resp := &datasource.ReadResponse{}

				variable := ds.fetchVariable(t.Context(), data, resp)

				assert.Nil(t, variable)
				assert.True(t, resp.Diagnostics.HasError())

			case "error - variable not found":
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(n8nsdk.VariableList{
						Data: []n8nsdk.Variable{},
					})
				})

				n8nClient, server := setupTestDataSourceClient(t, handler)
				defer server.Close()

				ds := &VariableDataSource{client: n8nClient}
				data := &models.DataSource{
					ID: types.StringValue("nonexistent"),
				}
				resp := &datasource.ReadResponse{}

				variable := ds.fetchVariable(t.Context(), data, resp)

				assert.Nil(t, variable)
				assert.True(t, resp.Diagnostics.HasError())

			case "error - variable not found with key identifier":
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(n8nsdk.VariableList{
						Data: []n8nsdk.Variable{},
					})
				})

				n8nClient, server := setupTestDataSourceClient(t, handler)
				defer server.Close()

				ds := &VariableDataSource{client: n8nClient}
				data := &models.DataSource{
					ID:  types.StringNull(),
					Key: types.StringValue("nonexistent-key"),
				}
				resp := &datasource.ReadResponse{}

				variable := ds.fetchVariable(t.Context(), data, resp)

				assert.Nil(t, variable)
				assert.True(t, resp.Diagnostics.HasError())
				assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Variable Not Found")
				assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "nonexistent-key")

			case "error - nil variable data":
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(n8nsdk.VariableList{
						Data: nil,
					})
				})

				n8nClient, server := setupTestDataSourceClient(t, handler)
				defer server.Close()

				ds := &VariableDataSource{client: n8nClient}
				data := &models.DataSource{
					ID: types.StringValue("var-123"),
				}
				resp := &datasource.ReadResponse{}

				variable := ds.fetchVariable(t.Context(), data, resp)

				assert.Nil(t, variable)
				assert.True(t, resp.Diagnostics.HasError())
			}
		})
	}
}

// ptrString returns a pointer to the given string.
//
//go:fix inline
func ptrString(s string) *string {
	return new(s)
}

// TestVariableDataSource_findVariableInList tests the private findVariableInList method.
func TestVariableDataSource_findVariableInList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		variables []n8nsdk.Variable
		id        string
		key       string
		wantNil   bool
		wantError bool
	}{
		{
			name: "found by ID",
			variables: []n8nsdk.Variable{
				{Id: ptrString("var-1"), Key: "key-1", Value: "val-1"},
				{Id: ptrString("var-2"), Key: "key-2", Value: "val-2"},
			},
			id:        "var-1",
			wantNil:   false,
			wantError: false,
		},
		{
			name: "found by key",
			variables: []n8nsdk.Variable{
				{Id: ptrString("var-1"), Key: "key-1", Value: "val-1"},
			},
			key:       "key-1",
			wantNil:   false,
			wantError: false,
		},
		{
			name:      "nil variables slice",
			variables: nil,
			id:        "var-1",
			wantNil:   true,
			wantError: true,
		},
		{
			name:      "empty variables slice",
			variables: []n8nsdk.Variable{},
			id:        "var-1",
			wantNil:   true,
			wantError: true,
		},
		{
			name: "not found - uses key as fallback identifier",
			variables: []n8nsdk.Variable{
				{Id: ptrString("var-1"), Key: "key-1", Value: "val-1"},
			},
			key:       "nonexistent",
			wantNil:   true,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ds := &VariableDataSource{}
			resp := &datasource.ReadResponse{}

			var data models.DataSource
			if tt.id != "" {
				data.ID = types.StringValue(tt.id)
			} else {
				data.ID = types.StringNull()
			}
			if tt.key != "" {
				data.Key = types.StringValue(tt.key)
			} else {
				data.Key = types.StringNull()
			}

			result := ds.findVariableInList(tt.variables, &data, resp)

			if tt.wantNil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
			}
			assert.Equal(t, tt.wantError, resp.Diagnostics.HasError())
		})
	}
}
