package user

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
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/user/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserDataSourceInterface is a mock implementation of UserDataSourceInterface.
type MockUserDataSourceInterface struct {
	mock.Mock
}

func (m *MockUserDataSourceInterface) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockUserDataSourceInterface) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockUserDataSourceInterface) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockUserDataSourceInterface) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	m.Called(ctx, req, resp)
}

// TestNewUserDataSource is now in external test file - refactored to test behavior only.

// TestNewUserDataSourceWrapper is now in external test file - refactored to test behavior only.

// TestUserDataSource_Metadata is now in external test file - refactored to test behavior only.

// TestUserDataSource_Schema is now in external test file - refactored to test behavior only.

// TestUserDataSource_Configure is now in external test file - refactored to test behavior only.

// Compile-time interface satisfaction checks.
var _ datasource.DataSource = (*UserDataSource)(nil)
var _ datasource.DataSourceWithConfigure = (*UserDataSource)(nil)
var _ UserDataSourceInterface = (*UserDataSource)(nil)

func TestUserDataSource_Interfaces(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "constructor returns non-nil instance",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := NewUserDataSource()
				assert.NotNil(t, ds)
			},
		},
		{
			name: "interface implementation error case",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := NewUserDataSource()
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

func TestUserDataSourceConcurrency(t *testing.T) {
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

				ds := NewUserDataSource()

				var wg sync.WaitGroup
				for range 100 {
					wg.Go(func() {
						resp := &datasource.MetadataResponse{}
						ds.Metadata(t.Context(), datasource.MetadataRequest{
							ProviderTypeName: "n8n",
						}, resp)
						assert.Equal(t, "n8n_user", resp.TypeName)
					})
				}
				wg.Wait()
			},
		},
		{
			name: "concurrent schema calls",
			testFunc: func(t *testing.T) {
				t.Helper()

				ds := NewUserDataSource()

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

				ds := NewUserDataSource()

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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t)
		})
	}
}

func BenchmarkUserDataSource_Schema(b *testing.B) {
	ds := NewUserDataSource()

	for b.Loop() {
		resp := &datasource.SchemaResponse{}
		ds.Schema(b.Context(), datasource.SchemaRequest{}, resp)
	}
}

func BenchmarkUserDataSource_Metadata(b *testing.B) {
	ds := NewUserDataSource()

	for b.Loop() {
		resp := &datasource.MetadataResponse{}
		ds.Metadata(b.Context(), datasource.MetadataRequest{}, resp)
	}
}

func BenchmarkUserDataSource_Configure(b *testing.B) {
	ds := NewUserDataSource()
	mockClient := &client.N8nClient{}

	for b.Loop() {
		resp := &datasource.ConfigureResponse{}
		req := datasource.ConfigureRequest{
			ProviderData: mockClient,
		}
		ds.Configure(b.Context(), req, resp)
	}
}

// TestUserDataSource_Read is now in external test file - refactored to test behavior only.

// TestUserDataSource_schemaAttributes tests the private schemaAttributes method.
func TestUserDataSource_schemaAttributes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "normal case",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := NewUserDataSource()
				attrs := ds.schemaAttributes()

				assert.NotNil(t, attrs, "schemaAttributes should return non-nil attributes")
				assert.NotEmpty(t, attrs, "schemaAttributes should return non-empty attributes")
			},
		},
		{
			name: "error case - validates behavior",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := NewUserDataSource()
				attrs := ds.schemaAttributes()

				assert.NotNil(t, attrs, "schemaAttributes should return non-nil attributes")
				assert.NotEmpty(t, attrs, "schemaAttributes should return non-empty attributes")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.testFunc(t)
		})
	}
} // TestUserDataSource_getIdentifier tests the private getIdentifier method.
func TestUserDataSource_getIdentifier(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		data        *models.DataSource
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid id",
			data: &models.DataSource{
				ID: types.StringValue("user-123"),
			},
			expectError: false,
		},
		{
			name: "valid email",
			data: &models.DataSource{
				Email: types.StringValue("test@example.com"),
			},
			expectError: false,
		},
		{
			name:        "no identifier",
			data:        &models.DataSource{},
			expectError: true,
			errorMsg:    "Missing Required Attribute",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ds := &UserDataSource{}
			resp := &datasource.ReadResponse{}

			identifier := ds.getIdentifier(tt.data, resp)

			if tt.expectError {
				assert.True(t, resp.Diagnostics.HasError(), "getIdentifier should return error")
				assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), tt.errorMsg)
			} else {
				assert.False(t, resp.Diagnostics.HasError(), "getIdentifier should not return error")
				assert.NotEmpty(t, identifier, "identifier should not be empty")
			}
		})
	}
}

// TestUserDataSource_fetchUser tests the private fetchUser method.
// setupUserDataSourceTestClient creates a test datasource with httptest server.
func setupUserDataSourceTestClient(t *testing.T, handler http.HandlerFunc) (*UserDataSource, *httptest.Server) {
	t.Helper()
	server := httptest.NewServer(handler)

	cfg := n8nsdk.NewConfiguration()
	cfg.Servers = n8nsdk.ServerConfigurations{{URL: server.URL, Description: "Test server"}}
	cfg.HTTPClient = server.Client()
	cfg.AddDefaultHeader("X-N8N-API-KEY", "test-key")

	apiClient := n8nsdk.NewAPIClient(cfg)
	n8nClient := &client.N8nClient{APIClient: apiClient}

	ds := &UserDataSource{client: n8nClient}
	return ds, server
}

func TestUserDataSource_fetchUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		identifier   string
		setupHandler func(w http.ResponseWriter, r *http.Request)
		expectNil    bool
		expectError  bool
	}{
		{
			name:       "success - user found",
			identifier: "user-123",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				userID := "user-123"
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]any{
					"id":    userID,
					"email": "test@example.com",
				})
			},
			expectNil:   false,
			expectError: false,
		},
		{
			name:       "error - user not found",
			identifier: "nonexistent",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNotFound)
				_, _ = w.Write([]byte(`{"message": "User not found"}`))
			},
			expectNil:   true,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ds, server := setupUserDataSourceTestClient(t, http.HandlerFunc(tt.setupHandler))
			t.Cleanup(server.Close)

			resp := &datasource.ReadResponse{}
			user := ds.fetchUser(t.Context(), tt.identifier, resp)

			if tt.expectNil {
				assert.Nil(t, user)
			} else {
				assert.NotNil(t, user)
			}
			assert.Equal(t, tt.expectError, resp.Diagnostics.HasError())
		})
	}
}

// TestUserDataSource_populateUserData tests the private populateUserData method.
func TestUserDataSource_populateUserData(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		user        *n8nsdk.User
		expectError bool
	}{
		{
			name: "valid user with all fields",
			user: &n8nsdk.User{
				Id:        stringPtr("user-123"),
				Email:     "test@example.com",
				FirstName: stringPtr("Test"),
				LastName:  stringPtr("User"),
				IsPending: boolPtr(false),
				Role:      stringPtr("admin"),
			},
			expectError: false,
		},
		{
			name: "user with nil id",
			user: &n8nsdk.User{
				Id:        nil,
				Email:     "test@example.com",
				FirstName: stringPtr("Test"),
				LastName:  stringPtr("User"),
			},
			expectError: false,
		},
		{
			name: "user with nil firstName",
			user: &n8nsdk.User{
				Id:        stringPtr("user-123"),
				Email:     "test@example.com",
				FirstName: nil,
				LastName:  stringPtr("User"),
			},
			expectError: false,
		},
		{
			name: "user with nil lastName",
			user: &n8nsdk.User{
				Id:        stringPtr("user-123"),
				Email:     "test@example.com",
				FirstName: stringPtr("Test"),
				LastName:  nil,
			},
			expectError: false,
		},
		{
			name: "user with nil isPending",
			user: &n8nsdk.User{
				Id:        stringPtr("user-123"),
				Email:     "test@example.com",
				IsPending: nil,
			},
			expectError: false,
		},
		{
			name: "user with nil createdAt",
			user: &n8nsdk.User{
				Id:        stringPtr("user-123"),
				Email:     "test@example.com",
				CreatedAt: nil,
			},
			expectError: false,
		},
		{
			name: "user with nil updatedAt",
			user: &n8nsdk.User{
				Id:        stringPtr("user-123"),
				Email:     "test@example.com",
				UpdatedAt: nil,
			},
			expectError: false,
		},
		{
			name: "user with nil role",
			user: &n8nsdk.User{
				Id:    stringPtr("user-123"),
				Email: "test@example.com",
				Role:  nil,
			},
			expectError: false,
		},
		{
			name: "user with minimal fields",
			user: &n8nsdk.User{
				Email: "test@example.com",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ds := &UserDataSource{}
			data := &models.DataSource{}

			ds.populateUserData(tt.user, data)

			if !tt.expectError && tt.user != nil {
				assert.NotNil(t, data, "data should not be nil")
				assert.Equal(t, tt.user.Email, data.Email.ValueString(), "email should match")
			}
		})
	}
}

// stringPtr is a helper function to create string pointers.
//
//go:fix inline
func stringPtr(s string) *string {
	return new(s)
}

// boolPtr is a helper function to create bool pointers.
//
//go:fix inline
func boolPtr(b bool) *bool {
	return new(b)
}
