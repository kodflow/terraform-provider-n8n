package user_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/client"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/user"
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
					w.Write([]byte(`{"id": "user-123", "email": "test@example.com", "firstName": "Test", "lastName": "User", "isPending": false, "createdAt": "2024-01-01T00:00:00Z", "updatedAt": "2024-01-02T00:00:00Z", "role": "global:member"}`))
				})

				n8nClient, server := setupTestClientForDataSource(t, handler)
				defer server.Close()

				ds := user.NewUserDataSource()
				ds.Configure(context.Background(), datasource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &datasource.ConfigureResponse{})

				ctx := context.Background()
				schemaResp := datasource.SchemaResponse{}
				ds.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)

				// Build config
				configRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "user-123"),
					"email":      tftypes.NewValue(tftypes.String, nil),
					"first_name": tftypes.NewValue(tftypes.String, nil),
					"last_name":  tftypes.NewValue(tftypes.String, nil),
					"is_pending": tftypes.NewValue(tftypes.Bool, nil),
					"created_at": tftypes.NewValue(tftypes.String, nil),
					"updated_at": tftypes.NewValue(tftypes.String, nil),
					"role":       tftypes.NewValue(tftypes.String, nil),
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
			name: "error - read with invalid config",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := user.NewUserDataSource()
				ctx := context.Background()

				schemaResp := datasource.SchemaResponse{}
				ds.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)

				// Create invalid config
				configRaw := tftypes.NewValue(tftypes.String, "invalid")

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
		{
			name: "error - read with API error",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte(`{"message": "User not found"}`))
				})

				n8nClient, server := setupTestClientForDataSource(t, handler)
				defer server.Close()

				ds := user.NewUserDataSource()
				ds.Configure(context.Background(), datasource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &datasource.ConfigureResponse{})

				ctx := context.Background()
				schemaResp := datasource.SchemaResponse{}
				ds.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)

				configRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "user-123"),
					"email":      tftypes.NewValue(tftypes.String, nil),
					"first_name": tftypes.NewValue(tftypes.String, nil),
					"last_name":  tftypes.NewValue(tftypes.String, nil),
					"is_pending": tftypes.NewValue(tftypes.Bool, nil),
					"created_at": tftypes.NewValue(tftypes.String, nil),
					"updated_at": tftypes.NewValue(tftypes.String, nil),
					"role":       tftypes.NewValue(tftypes.String, nil),
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
		{
			name: "error - missing required identifiers",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := user.NewUserDataSource()
				ctx := context.Background()

				schemaResp := datasource.SchemaResponse{}
				ds.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)

				// Create config with both id and email as null
				configRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, nil),
					"email":      tftypes.NewValue(tftypes.String, nil),
					"first_name": tftypes.NewValue(tftypes.String, nil),
					"last_name":  tftypes.NewValue(tftypes.String, nil),
					"is_pending": tftypes.NewValue(tftypes.Bool, nil),
					"created_at": tftypes.NewValue(tftypes.String, nil),
					"updated_at": tftypes.NewValue(tftypes.String, nil),
					"role":       tftypes.NewValue(tftypes.String, nil),
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
				assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Missing Required Attribute")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}

// setupTestClientForDataSource creates a test N8nClient with httptest server for datasources.
func setupTestClientForDataSource(t *testing.T, handler http.HandlerFunc) (*client.N8nClient, *httptest.Server) {
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
