package user_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/client"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/user"
	"github.com/stretchr/testify/assert"
)

// setupTestClient creates a test N8nClient with httptest server.
func setupTestClient(t *testing.T, handler http.HandlerFunc) (*client.N8nClient, *httptest.Server) {
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

// TestNewUserResource tests the NewUserResource constructor.
func TestNewUserResource(t *testing.T) {
	t.Parallel()

	r := user.NewUserResource()
	assert.NotNil(t, r, "NewUserResource should return a non-nil resource")
}

// TestNewUserResourceWrapper tests the NewUserResourceWrapper constructor.
func TestNewUserResourceWrapper(t *testing.T) {
	t.Parallel()

	r := user.NewUserResourceWrapper()
	assert.NotNil(t, r, "NewUserResourceWrapper should return a non-nil resource")
}

// TestUserResource_Metadata tests the Metadata method.
func TestUserResource_Metadata(t *testing.T) {
	t.Parallel()

	r := user.NewUserResource()
	req := resource.MetadataRequest{
		ProviderTypeName: "n8n",
	}
	resp := &resource.MetadataResponse{}

	r.Metadata(context.Background(), req, resp)

	assert.Equal(t, "n8n_user", resp.TypeName, "TypeName should be set correctly")
}

// TestUserResource_Schema tests the Schema method.
func TestUserResource_Schema(t *testing.T) {
	t.Parallel()

	r := user.NewUserResource()
	req := resource.SchemaRequest{}
	resp := &resource.SchemaResponse{}

	r.Schema(context.Background(), req, resp)

	assert.NotNil(t, resp.Schema, "Schema should not be nil")
	assert.NotNil(t, resp.Schema.Attributes, "Schema attributes should not be nil")
}

// TestUserResource_Configure tests the Configure method.
func TestUserResource_Configure(t *testing.T) {
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

			r := user.NewUserResource()
			req := resource.ConfigureRequest{
				ProviderData: tt.providerData,
			}
			resp := &resource.ConfigureResponse{}

			r.Configure(context.Background(), req, resp)

			if tt.expectError {
				assert.True(t, resp.Diagnostics.HasError(), "Configure should return error for invalid provider data")
			} else {
				assert.False(t, resp.Diagnostics.HasError(), "Configure should not return error")
			}
		})
	}
}

// TestUserResource_Create tests the Create method with full execution.
func TestUserResource_Create(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "create with successful API call",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					// Handle POST /users (create) and GET /users/:id (fetch details)
					switch r.Method {
					case "POST":
						w.WriteHeader(http.StatusCreated)
						// n8n API returns an array response
						w.Write([]byte(`[{"user":{"id":"user-123"},"error":""}]`))
					case "GET":
						w.WriteHeader(http.StatusOK)
						w.Write([]byte(`{"id":"user-123","email":"test@example.com","firstName":"John","lastName":"Doe","role":"global:member","isPending":false,"createdAt":"2024-01-01T00:00:00Z","updatedAt":"2024-01-01T00:00:00Z"}`))
					}
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				r := user.NewUserResource()
				r.Configure(context.Background(), resource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &resource.ConfigureResponse{})

				ctx := context.Background()
				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				// Build plan
				planRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, nil),
					"email":      tftypes.NewValue(tftypes.String, "test@example.com"),
					"first_name": tftypes.NewValue(tftypes.String, nil),
					"last_name":  tftypes.NewValue(tftypes.String, nil),
					"role":       tftypes.NewValue(tftypes.String, "global:member"),
					"is_pending": tftypes.NewValue(tftypes.Bool, nil),
					"created_at": tftypes.NewValue(tftypes.String, nil),
					"updated_at": tftypes.NewValue(tftypes.String, nil),
				})

				plan := tfsdk.Plan{
					Schema: schemaResp.Schema,
					Raw:    planRaw,
				}

				state := tfsdk.State{
					Schema: schemaResp.Schema,
				}

				req := resource.CreateRequest{
					Plan: plan,
				}
				resp := &resource.CreateResponse{
					State: state,
				}

				// Call Create
				r.Create(ctx, req, resp)

				// Verify success
				assert.False(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "error - create with invalid plan",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := user.NewUserResource()
				ctx := context.Background()

				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				// Create invalid plan that will fail Get()
				planRaw := tftypes.NewValue(tftypes.String, "invalid")

				plan := tfsdk.Plan{
					Schema: schemaResp.Schema,
					Raw:    planRaw,
				}

				req := resource.CreateRequest{
					Plan: plan,
				}
				resp := &resource.CreateResponse{}

				r.Create(ctx, req, resp)

				// Verify error
				assert.True(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "error - create with API error",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"message": "Internal server error"}`))
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				r := user.NewUserResource()
				r.Configure(context.Background(), resource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &resource.ConfigureResponse{})

				ctx := context.Background()
				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				planRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, nil),
					"email":      tftypes.NewValue(tftypes.String, "test@example.com"),
					"first_name": tftypes.NewValue(tftypes.String, nil),
					"last_name":  tftypes.NewValue(tftypes.String, nil),
					"role":       tftypes.NewValue(tftypes.String, "global:member"),
					"is_pending": tftypes.NewValue(tftypes.Bool, nil),
					"created_at": tftypes.NewValue(tftypes.String, nil),
					"updated_at": tftypes.NewValue(tftypes.String, nil),
				})

				plan := tfsdk.Plan{
					Schema: schemaResp.Schema,
					Raw:    planRaw,
				}

				state := tfsdk.State{
					Schema: schemaResp.Schema,
				}

				req := resource.CreateRequest{
					Plan: plan,
				}
				resp := &resource.CreateResponse{
					State: state,
				}

				r.Create(ctx, req, resp)

				// Verify error
				assert.True(t, resp.Diagnostics.HasError())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}

// TestUserResource_Read tests the Read method with full execution.
func TestUserResource_Read(t *testing.T) {
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
					w.Write([]byte(`{"id":"user-123","email":"test@example.com","firstName":"John","lastName":"Doe","role":"global:member","isPending":false,"createdAt":"2024-01-01T00:00:00Z","updatedAt":"2024-01-01T00:00:00Z"}`))
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				r := user.NewUserResource()
				r.Configure(context.Background(), resource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &resource.ConfigureResponse{})

				ctx := context.Background()
				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				// Build state
				stateRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "user-123"),
					"email":      tftypes.NewValue(tftypes.String, "test@example.com"),
					"first_name": tftypes.NewValue(tftypes.String, "John"),
					"last_name":  tftypes.NewValue(tftypes.String, "Doe"),
					"role":       tftypes.NewValue(tftypes.String, "global:member"),
					"is_pending": tftypes.NewValue(tftypes.Bool, false),
					"created_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
					"updated_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
				})

				state := tfsdk.State{
					Schema: schemaResp.Schema,
					Raw:    stateRaw,
				}

				req := resource.ReadRequest{
					State: state,
				}
				resp := &resource.ReadResponse{
					State: state,
				}

				// Call Read
				r.Read(ctx, req, resp)

				// Verify success
				assert.False(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "error - read with invalid state",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := user.NewUserResource()
				ctx := context.Background()

				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				// Create invalid state
				stateRaw := tftypes.NewValue(tftypes.String, "invalid")

				state := tfsdk.State{
					Schema: schemaResp.Schema,
					Raw:    stateRaw,
				}

				req := resource.ReadRequest{
					State: state,
				}
				resp := &resource.ReadResponse{}

				r.Read(ctx, req, resp)

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

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				r := user.NewUserResource()
				r.Configure(context.Background(), resource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &resource.ConfigureResponse{})

				ctx := context.Background()
				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				stateRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "user-123"),
					"email":      tftypes.NewValue(tftypes.String, "test@example.com"),
					"first_name": tftypes.NewValue(tftypes.String, "John"),
					"last_name":  tftypes.NewValue(tftypes.String, "Doe"),
					"role":       tftypes.NewValue(tftypes.String, "global:member"),
					"is_pending": tftypes.NewValue(tftypes.Bool, false),
					"created_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
					"updated_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
				})

				state := tfsdk.State{
					Schema: schemaResp.Schema,
					Raw:    stateRaw,
				}

				req := resource.ReadRequest{
					State: state,
				}
				resp := &resource.ReadResponse{
					State: state,
				}

				r.Read(ctx, req, resp)

				// Verify error
				assert.True(t, resp.Diagnostics.HasError())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}

// TestUserResource_Update tests the Update method with full execution.
func TestUserResource_Update(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "update with successful API call",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{"id":"user-123","email":"test@example.com","firstName":"John","lastName":"Doe","role":"global:admin","isPending":false,"createdAt":"2024-01-01T00:00:00Z","updatedAt":"2024-01-02T00:00:00Z"}`))
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				r := user.NewUserResource()
				r.Configure(context.Background(), resource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &resource.ConfigureResponse{})

				ctx := context.Background()
				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				// Build plan
				planRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "user-123"),
					"email":      tftypes.NewValue(tftypes.String, "test@example.com"),
					"first_name": tftypes.NewValue(tftypes.String, "John"),
					"last_name":  tftypes.NewValue(tftypes.String, "Doe"),
					"role":       tftypes.NewValue(tftypes.String, "global:admin"),
					"is_pending": tftypes.NewValue(tftypes.Bool, false),
					"created_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
					"updated_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
				})

				plan := tfsdk.Plan{
					Schema: schemaResp.Schema,
					Raw:    planRaw,
				}

				// Build state
				stateRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "user-123"),
					"email":      tftypes.NewValue(tftypes.String, "test@example.com"),
					"first_name": tftypes.NewValue(tftypes.String, "John"),
					"last_name":  tftypes.NewValue(tftypes.String, "Doe"),
					"role":       tftypes.NewValue(tftypes.String, "global:member"),
					"is_pending": tftypes.NewValue(tftypes.Bool, false),
					"created_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
					"updated_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
				})

				state := tfsdk.State{
					Schema: schemaResp.Schema,
					Raw:    stateRaw,
				}

				req := resource.UpdateRequest{
					Plan:  plan,
					State: state,
				}
				resp := &resource.UpdateResponse{
					State: state,
				}

				// Call Update
				r.Update(ctx, req, resp)

				// Verify success
				assert.False(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "error - update with invalid plan",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := user.NewUserResource()
				ctx := context.Background()

				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				// Create invalid plan
				planRaw := tftypes.NewValue(tftypes.String, "invalid")

				plan := tfsdk.Plan{
					Schema: schemaResp.Schema,
					Raw:    planRaw,
				}

				// Create valid state
				stateRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "user-123"),
					"email":      tftypes.NewValue(tftypes.String, "test@example.com"),
					"first_name": tftypes.NewValue(tftypes.String, "John"),
					"last_name":  tftypes.NewValue(tftypes.String, "Doe"),
					"role":       tftypes.NewValue(tftypes.String, "global:member"),
					"is_pending": tftypes.NewValue(tftypes.Bool, false),
					"created_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
					"updated_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
				})

				state := tfsdk.State{
					Schema: schemaResp.Schema,
					Raw:    stateRaw,
				}

				req := resource.UpdateRequest{
					Plan:  plan,
					State: state,
				}
				resp := &resource.UpdateResponse{
					State: state,
				}

				r.Update(ctx, req, resp)

				// Verify error
				assert.True(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "error - update with API error",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"message": "Internal server error"}`))
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				r := user.NewUserResource()
				r.Configure(context.Background(), resource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &resource.ConfigureResponse{})

				ctx := context.Background()
				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				planRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "user-123"),
					"email":      tftypes.NewValue(tftypes.String, "test@example.com"),
					"first_name": tftypes.NewValue(tftypes.String, "John"),
					"last_name":  tftypes.NewValue(tftypes.String, "Doe"),
					"role":       tftypes.NewValue(tftypes.String, "global:admin"),
					"is_pending": tftypes.NewValue(tftypes.Bool, false),
					"created_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
					"updated_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
				})

				plan := tfsdk.Plan{
					Schema: schemaResp.Schema,
					Raw:    planRaw,
				}

				stateRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "user-123"),
					"email":      tftypes.NewValue(tftypes.String, "test@example.com"),
					"first_name": tftypes.NewValue(tftypes.String, "John"),
					"last_name":  tftypes.NewValue(tftypes.String, "Doe"),
					"role":       tftypes.NewValue(tftypes.String, "global:member"),
					"is_pending": tftypes.NewValue(tftypes.Bool, false),
					"created_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
					"updated_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
				})

				state := tfsdk.State{
					Schema: schemaResp.Schema,
					Raw:    stateRaw,
				}

				req := resource.UpdateRequest{
					Plan:  plan,
					State: state,
				}
				resp := &resource.UpdateResponse{
					State: state,
				}

				r.Update(ctx, req, resp)

				// Verify error
				assert.True(t, resp.Diagnostics.HasError())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}

// TestUserResource_Delete tests the Delete method with full execution.
func TestUserResource_Delete(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "delete with successful API call",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNoContent)
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				r := user.NewUserResource()
				r.Configure(context.Background(), resource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &resource.ConfigureResponse{})

				ctx := context.Background()
				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				// Build state
				stateRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "user-123"),
					"email":      tftypes.NewValue(tftypes.String, "test@example.com"),
					"first_name": tftypes.NewValue(tftypes.String, "John"),
					"last_name":  tftypes.NewValue(tftypes.String, "Doe"),
					"role":       tftypes.NewValue(tftypes.String, "global:member"),
					"is_pending": tftypes.NewValue(tftypes.Bool, false),
					"created_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
					"updated_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
				})

				state := tfsdk.State{
					Schema: schemaResp.Schema,
					Raw:    stateRaw,
				}

				req := resource.DeleteRequest{
					State: state,
				}
				resp := &resource.DeleteResponse{}

				// Call Delete
				r.Delete(ctx, req, resp)

				// Verify success
				assert.False(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "error - delete with invalid state",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := user.NewUserResource()
				ctx := context.Background()

				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				// Create invalid state
				stateRaw := tftypes.NewValue(tftypes.String, "invalid")

				state := tfsdk.State{
					Schema: schemaResp.Schema,
					Raw:    stateRaw,
				}

				req := resource.DeleteRequest{
					State: state,
				}
				resp := &resource.DeleteResponse{}

				r.Delete(ctx, req, resp)

				// Verify error
				assert.True(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "error - delete with API error",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"message": "Internal server error"}`))
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				r := user.NewUserResource()
				r.Configure(context.Background(), resource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &resource.ConfigureResponse{})

				ctx := context.Background()
				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				stateRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "user-123"),
					"email":      tftypes.NewValue(tftypes.String, "test@example.com"),
					"first_name": tftypes.NewValue(tftypes.String, "John"),
					"last_name":  tftypes.NewValue(tftypes.String, "Doe"),
					"role":       tftypes.NewValue(tftypes.String, "global:member"),
					"is_pending": tftypes.NewValue(tftypes.Bool, false),
					"created_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
					"updated_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
				})

				state := tfsdk.State{
					Schema: schemaResp.Schema,
					Raw:    stateRaw,
				}

				req := resource.DeleteRequest{
					State: state,
				}
				resp := &resource.DeleteResponse{}

				r.Delete(ctx, req, resp)

				// Verify error
				assert.True(t, resp.Diagnostics.HasError())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}

// TestUserResource_ImportState tests the ImportState method.
func TestUserResource_ImportState(t *testing.T) {
	t.Parallel()

	r := &user.UserResource{}
	ctx := context.Background()

	// Create schema for state
	schemaResp := resource.SchemaResponse{}
	r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

	// Create an empty state with the schema
	state := tfsdk.State{
		Schema: schemaResp.Schema,
	}

	// Initialize the raw value with required attributes
	state.Raw = tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
		"id":         tftypes.NewValue(tftypes.String, nil),
		"email":      tftypes.NewValue(tftypes.String, nil),
		"first_name": tftypes.NewValue(tftypes.String, nil),
		"last_name":  tftypes.NewValue(tftypes.String, nil),
		"role":       tftypes.NewValue(tftypes.String, nil),
		"is_pending": tftypes.NewValue(tftypes.Bool, nil),
		"created_at": tftypes.NewValue(tftypes.String, nil),
		"updated_at": tftypes.NewValue(tftypes.String, nil),
	})

	req := resource.ImportStateRequest{
		ID: "user-123",
	}
	resp := &resource.ImportStateResponse{
		State: state,
	}

	r.ImportState(ctx, req, resp)

	// Check that import succeeded
	assert.NotNil(t, resp)
	assert.False(t, resp.Diagnostics.HasError(), "ImportState should not have errors")
}
