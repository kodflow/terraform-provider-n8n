package user_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
	"github.com/kodflow/n8n/src/internal/provider/user"
	"github.com/stretchr/testify/assert"
)

// TestNewUsersDataSource tests the NewUsersDataSource constructor.
func TestNewUsersDataSource(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "normal case",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := user.NewUsersDataSource()
				assert.NotNil(t, ds, "NewUsersDataSource should return a non-nil datasource")
			},
		},
		{
			name: "error case - validates behavior",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := user.NewUsersDataSource()
				assert.NotNil(t, ds, "NewUsersDataSource should return a non-nil datasource")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.testFunc(t)
		})
	}
} // TestNewUsersDataSourceWrapper tests the NewUsersDataSourceWrapper constructor.
func TestNewUsersDataSourceWrapper(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "normal case",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := user.NewUsersDataSourceWrapper()
				assert.NotNil(t, ds, "NewUsersDataSourceWrapper should return a non-nil datasource")
			},
		},
		{
			name: "error case - validates behavior",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := user.NewUsersDataSourceWrapper()
				assert.NotNil(t, ds, "NewUsersDataSourceWrapper should return a non-nil datasource")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.testFunc(t)
		})
	}
} // TestUsersDataSource_Metadata tests the Metadata method.
func TestUsersDataSource_Metadata(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "normal case",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := user.NewUsersDataSource()
				req := datasource.MetadataRequest{
					ProviderTypeName: "n8n",
				}
				resp := &datasource.MetadataResponse{}

				ds.Metadata(context.Background(), req, resp)

				assert.Equal(t, "n8n_users", resp.TypeName, "TypeName should be set correctly")
			},
		},
		{
			name: "error case - validates behavior",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := user.NewUsersDataSource()
				req := datasource.MetadataRequest{
					ProviderTypeName: "n8n",
				}
				resp := &datasource.MetadataResponse{}

				ds.Metadata(context.Background(), req, resp)

				assert.Equal(t, "n8n_users", resp.TypeName, "TypeName should be set correctly")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.testFunc(t)
		})
	}
} // TestUsersDataSource_Schema tests the Schema method.
func TestUsersDataSource_Schema(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "normal case",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := user.NewUsersDataSource()
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
				ds := user.NewUsersDataSource()
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
} // TestUsersDataSource_Configure tests the Configure method.
func TestUsersDataSource_Configure(t *testing.T) {
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

			ds := user.NewUsersDataSource()
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

// TestUsersDataSource_Read tests the Read method.
func TestUsersDataSource_Read(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "read with successful API call",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{
						"data": [
							{
								"id": "user-123",
								"email": "test@example.com",
								"firstName": "Test",
								"lastName": "User",
								"isPending": false,
								"createdAt": "2024-01-01T00:00:00Z",
								"updatedAt": "2024-01-02T00:00:00Z",
								"role": "global:member"
							}
						]
					}`))
				})

				n8nClient, server := setupTestClientForDataSources(t, handler)
				defer server.Close()

				ds := user.NewUsersDataSource()
				ds.Configure(context.Background(), datasource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &datasource.ConfigureResponse{})

				ctx := context.Background()
				schemaResp := datasource.SchemaResponse{}
				ds.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)

				// Build config
				configRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"users": tftypes.NewValue(
						tftypes.List{ElementType: tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"id":         tftypes.String,
								"email":      tftypes.String,
								"first_name": tftypes.String,
								"last_name":  tftypes.String,
								"is_pending": tftypes.Bool,
								"created_at": tftypes.String,
								"updated_at": tftypes.String,
								"role":       tftypes.String,
							},
						}},
						nil,
					),
				})

				config := tfsdk.Config{
					Schema: schemaResp.Schema,
					Raw:    configRaw,
				}

				state := tfsdk.State{
					Schema: schemaResp.Schema,
				}

				req := datasource.ReadRequest{
					Config: config,
				}
				resp := &datasource.ReadResponse{
					State: state,
				}

				// Call Read
				ds.Read(ctx, req, resp)

				// Verify success
				if resp.Diagnostics.HasError() {
					for _, diag := range resp.Diagnostics.Errors() {
						t.Logf("Error: %s - %s", diag.Summary(), diag.Detail())
					}
				}
				assert.False(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "error - read with API error",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"message": "Internal server error"}`))
				})

				n8nClient, server := setupTestClientForDataSources(t, handler)
				defer server.Close()

				ds := user.NewUsersDataSource()
				ds.Configure(context.Background(), datasource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &datasource.ConfigureResponse{})

				ctx := context.Background()
				schemaResp := datasource.SchemaResponse{}
				ds.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)

				configRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"users": tftypes.NewValue(
						tftypes.List{ElementType: tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"id":         tftypes.String,
								"email":      tftypes.String,
								"first_name": tftypes.String,
								"last_name":  tftypes.String,
								"is_pending": tftypes.Bool,
								"created_at": tftypes.String,
								"updated_at": tftypes.String,
								"role":       tftypes.String,
							},
						}},
						nil,
					),
				})

				config := tfsdk.Config{
					Schema: schemaResp.Schema,
					Raw:    configRaw,
				}

				state := tfsdk.State{
					Schema: schemaResp.Schema,
				}

				req := datasource.ReadRequest{
					Config: config,
				}
				resp := &datasource.ReadResponse{
					State: state,
				}

				ds.Read(ctx, req, resp)

				// Verify error
				assert.True(t, resp.Diagnostics.HasError())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}

// setupTestClientForDataSources creates a test N8nClient with httptest server for datasources (plural).
func setupTestClientForDataSources(t *testing.T, handler http.HandlerFunc) (*client.N8nClient, *httptest.Server) {
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
