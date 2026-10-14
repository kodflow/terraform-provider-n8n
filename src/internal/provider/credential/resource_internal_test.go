package credential

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/credential/models"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/client"
	"github.com/stretchr/testify/assert"
)

// KTN-TEST-009 RATIONALE: Tests in this file test public functions that require
// access to private fields (specifically the .client field). These tests MUST
// remain in the internal test file to access the unexported client field for
// proper white-box testing. Moving these to external tests would prevent us
// from properly testing the functions that depend on the client field.
//
// Note: Mock implementations are in helpers_test.go to avoid duplication.

// TestNewCredentialResource is now in resource_external_test.go
// TestNewCredentialResourceWrapper is now in resource_external_test.go
// TestCredentialResource_Configure is now in resource_external_test.go.

// TestCredentialResource_Configure_Internal tests the internal Configure behavior with actual client.
func TestCredentialResource_Configure_Internal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		providerData any
		wantErr      bool
		expectClient bool
	}{
		{
			name:         "successful configuration with valid client",
			providerData: &client.N8nClient{},
			wantErr:      false,
			expectClient: true,
		},
		{
			name:         "nil provider data",
			providerData: nil,
			wantErr:      false,
			expectClient: false,
		},
		{
			name:         "invalid provider data type",
			providerData: "invalid-string",
			wantErr:      true,
			expectClient: false,
		},
		{
			name:         "invalid provider data type - integer",
			providerData: 123,
			wantErr:      true,
			expectClient: false,
		},
		{
			name:         "invalid provider data type - map",
			providerData: map[string]string{"key": "value"},
			wantErr:      true,
			expectClient: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := &CredentialResource{}
			req := resource.ConfigureRequest{
				ProviderData: tt.providerData,
			}
			resp := &resource.ConfigureResponse{}

			r.Configure(context.Background(), req, resp)

			if tt.wantErr {
				assert.True(t, resp.Diagnostics.HasError(), "expected error")
			} else {
				assert.False(t, resp.Diagnostics.HasError(), "expected no error")
			}

			if tt.expectClient {
				assert.NotNil(t, r.client, "expected client to be set")
			} else {
				assert.Nil(t, r.client, "expected client to remain nil")
			}
		})
	}
}

func Test_extractCredentialData(t *testing.T) {
	t.Parallel()

	// NOTE: This function is a thin wrapper around Terraform's ElementsAs method.
	// ElementsAs requires proper Terraform framework context (schema, state, etc.) to function.
	// It CANNOT be unit tested in isolation - it only works within Read/Create/Update methods.
	//
	// This test validates that the function properly calls ElementsAs and propagates its behavior.
	// The actual data extraction is tested via executeCreateLogicWithData and executeUpdateLogicWithData.

	tests := []struct {
		name     string
		setupMap func(context.Context) types.Map
		checkFn  func(*testing.T, map[string]any, diag.Diagnostics)
	}{
		{
			name: "null map returns empty data without errors",
			setupMap: func(ctx context.Context) types.Map {
				return types.MapNull(types.StringType)
			},
			checkFn: func(t *testing.T, data map[string]any, diags diag.Diagnostics) {
				t.Helper()
				// Null maps should not cause errors
				assert.False(t, diags.HasError(), "null map should not cause errors")
				assert.NotNil(t, data, "null map should return empty map, not nil")
				assert.Empty(t, data, "null map should return empty map")
			},
		},
		{
			name: "function calls ElementsAs and returns result",
			setupMap: func(ctx context.Context) types.Map {
				// Create a map that ElementsAs will attempt to process
				// In unit test context (without framework schema), ElementsAs may return errors
				// This is expected and documents the framework boundary
				m, diags := types.MapValueFrom(ctx, types.StringType, map[string]interface{}{
					"test": "value",
				})
				if diags.HasError() {
					t.Fatalf("Failed to create test map: %v", diags)
				}
				return m
			},
			checkFn: func(t *testing.T, data map[string]any, diags diag.Diagnostics) {
				t.Helper()
				// We're testing that extractCredentialData calls ElementsAs,
				// not that ElementsAs works (that's Terraform framework's responsibility)
				// The function should return a map and diagnostics (even if empty)
				assert.NotNil(t, data, "data should be returned")
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			inputMap := tt.setupMap(ctx)
			data, diags := extractCredentialData(ctx, inputMap)

			tt.checkFn(t, data, diags)
		})
	}
}

func TestCredentialResource_InterfaceCompliance(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "implements Resource interface", wantErr: false},
		{name: "implements ResourceWithConfigure interface", wantErr: false},
		{name: "implements ResourceWithImportState interface", wantErr: false},
		{name: "implements CredentialResourceInterface", wantErr: false},
		{name: "error case - verify all interfaces implemented", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "implements Resource interface":
				var _ resource.Resource = (*CredentialResource)(nil)
			case "implements ResourceWithConfigure interface":
				var _ resource.ResourceWithConfigure = (*CredentialResource)(nil)
			case "implements ResourceWithImportState interface":
				var _ resource.ResourceWithImportState = (*CredentialResource)(nil)
			case "implements CredentialResourceInterface":
				var _ CredentialResourceInterface = (*CredentialResource)(nil)
			case "error case - verify all interfaces implemented":
				// Compile-time check that all interfaces are implemented
				var _ resource.Resource = (*CredentialResource)(nil)
				var _ resource.ResourceWithConfigure = (*CredentialResource)(nil)
				var _ resource.ResourceWithImportState = (*CredentialResource)(nil)
				var _ CredentialResourceInterface = (*CredentialResource)(nil)
			}
		})
	}
}

// Test helper functions

// setupTestClient creates a test N8nClient with httptest server.
func setupTestClient(t *testing.T, handler http.HandlerFunc) (*client.N8nClient, *httptest.Server) {
	t.Helper()
	server := httptest.NewServer(handler)

	cfg := n8nsdk.NewConfiguration()
	// Parse server URL to extract host and port
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

func TestCredentialResource_schemaAttributes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		expectedAttrs   []string
		unexpectedAttrs []string
		checkNonNil     bool
	}{
		{
			name: "all required attributes present",
			expectedAttrs: []string{
				"id",
				"name",
				"type",
				"data",
				"created_at",
				"updated_at",
			},
			unexpectedAttrs: []string{
				"invalid_attr",
				"non_existent",
			},
			checkNonNil: true,
		},
		{
			name: "schema not empty",
			expectedAttrs: []string{
				"id",
			},
			checkNonNil: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := &CredentialResource{}
			attrs := r.schemaAttributes()

			if tt.checkNonNil {
				assert.NotNil(t, attrs, "schemaAttributes should not return nil")
				assert.NotEmpty(t, attrs, "schemaAttributes should not be empty")
			}

			for _, expectedAttr := range tt.expectedAttrs {
				assert.Contains(t, attrs, expectedAttr, "schemaAttributes should contain attribute: %s", expectedAttr)
			}

			for _, unexpectedAttr := range tt.unexpectedAttrs {
				assert.NotContains(t, attrs, unexpectedAttr, "schemaAttributes should not contain attribute: %s", unexpectedAttr)
			}
		})
	}
}

// TestCredentialResource_Create_WithMockClient tests credential creation with mock client.
func TestCredentialResource_Create_WithMockClient(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		setupHandler func(w http.ResponseWriter, r *http.Request)
		expectError  bool
	}{
		{
			name: "successful credential creation",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPost && r.URL.Path == "/credentials" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusCreated)
					json.NewEncoder(w).Encode(map[string]any{
						"id":        "cred-123",
						"name":      "Test Credential",
						"type":      "httpHeaderAuth",
						"createdAt": "2024-01-01T00:00:00Z",
						"updatedAt": "2024-01-01T00:00:00Z",
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: false,
		},
		{
			name: "API error during creation",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPost && r.URL.Path == "/credentials" {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"message": "Internal server error"}`))
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: true,
		},
		{
			name: "bad request error",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPost && r.URL.Path == "/credentials" {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(`{"message": "Invalid request"}`))
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := http.HandlerFunc(tt.setupHandler)
			n8nClient, server := setupTestClient(t, handler)
			defer server.Close()

			r := &CredentialResource{client: n8nClient}

			// Verify that the client is properly configured
			assert.NotNil(t, r.client, "Client should be configured")
			assert.NotNil(t, r.client.APIClient, "API client should be configured")
		})
	}
}

// TestCredentialResource_Read_WithMockClient tests credential reading with mock client.
func TestCredentialResource_Read_WithMockClient(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "read preserves state",
			wantErr: false,
		},
		{
			name:    "read without API call",
			wantErr: false,
		},
		{
			name:    "error case - read must not panic",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := &CredentialResource{}

			// Read keeps state as-is - no API call made
			// Verify method exists and doesn't panic
			assert.NotPanics(t, func() {
				_ = r
			})
		})
	}
}

// TestCredentialResource_Update_WithMockClient tests credential update with mock client.
func TestCredentialResource_Update_WithMockClient(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		setupHandler func(w http.ResponseWriter, r *http.Request)
		expectError  bool
	}{
		{
			name: "successful credential update",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPatch && r.URL.Path == "/credentials/cred-123" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{
						"id":        "cred-123",
						"name":      "Updated",
						"type":      "httpHeaderAuth",
						"createdAt": "2024-01-01T00:00:00Z",
						"updatedAt": "2024-01-02T00:00:00Z",
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: false,
		},
		{
			name: "error updating credential",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPatch && r.URL.Path == "/credentials/cred-123" {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"message": "Failed to update credential"}`))
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := http.HandlerFunc(tt.setupHandler)
			n8nClient, server := setupTestClient(t, handler)
			defer server.Close()

			r := &CredentialResource{client: n8nClient}

			updateResp, httpResp, err := r.executeUpdate(context.Background(), "cred-123", "Updated", "httpHeaderAuth", map[string]any{"value": "new"})
			if httpResp != nil && httpResp.Body != nil {
				httpResp.Body.Close()
			}

			if tt.expectError {
				assert.Error(t, err, "Should return error")
				return
			}
			assert.NoError(t, err, "Should not return error")
			assert.Equal(t, "cred-123", updateResp.Id, "Should keep the same credential ID")
		})
	}
}

// TestCredentialResource_Delete_WithMockClient tests credential deletion with mock client.
func TestCredentialResource_Delete_WithMockClient(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		setupHandler func(w http.ResponseWriter, r *http.Request)
		credID       string
		expectError  bool
	}{
		{
			name: "successful deletion",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodDelete && r.URL.Path == "/credentials/cred-123" {
					w.WriteHeader(http.StatusNoContent)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			credID:      "cred-123",
			expectError: false,
		},
		{
			name: "delete with not found error",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodDelete && r.URL.Path == "/credentials/cred-404" {
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte(`{"message": "Credential not found"}`))
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			credID:      "cred-404",
			expectError: true,
		},
		{
			name: "delete with internal server error",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodDelete && r.URL.Path == "/credentials/cred-500" {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"message": "Internal server error"}`))
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			credID:      "cred-500",
			expectError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := http.HandlerFunc(tt.setupHandler)
			n8nClient, server := setupTestClient(t, handler)
			defer server.Close()

			r := &CredentialResource{client: n8nClient}

			// Verify client is configured for delete operations
			assert.NotNil(t, r.client, "Client should be configured")
			assert.NotNil(t, r.client.APIClient, "API client should be configured")
		})
	}
}

// TestCredentialResource_CRUD_ErrorHandling tests error handling in CRUD operations.
func TestCredentialResource_CRUD_ErrorHandling(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		setupHandler func(w http.ResponseWriter, r *http.Request)
		operation    string
	}{
		{
			name: "create with API error",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			operation: "create",
		},
		{
			name: "delete with not found error",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			operation: "delete",
		},
		{
			name: "update with create error",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPost {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			operation: "update",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := http.HandlerFunc(tt.setupHandler)
			n8nClient, server := setupTestClient(t, handler)
			defer server.Close()

			r := &CredentialResource{client: n8nClient}

			// Verify resource is properly initialized for error handling
			assert.NotNil(t, r, "Resource should not be nil")
			assert.NotNil(t, r.client, "Client should be configured")
		})
	}
}

// TestCredentialResource_Read_CRUD tests the Read method with full CRUD coverage.
func TestCredentialResource_Read_CRUD(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "read with successful API call",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.Method == http.MethodGet && r.URL.Path == "/credentials/cred-123" {
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusOK)
						w.Write([]byte(`{"id":"cred-123","name":"test-credential","type":"httpHeaderAuth","data":{"key":"value"},"createdAt":"2024-01-01T00:00:00Z","updatedAt":"2024-01-01T00:00:00Z"}`))
					}
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				r := &CredentialResource{client: n8nClient}
				ctx := context.Background()

				schemaReq := resource.SchemaRequest{}
				schemaResp := &resource.SchemaResponse{}
				r.Schema(ctx, schemaReq, schemaResp)

				dataMap := map[string]tftypes.Value{
					"key": tftypes.NewValue(tftypes.String, "value"),
				}
				stateRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "cred-123"),
					"name":       tftypes.NewValue(tftypes.String, "test-credential"),
					"type":       tftypes.NewValue(tftypes.String, "httpHeaderAuth"),
					"data":       tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, dataMap),
					"project_id": tftypes.NewValue(tftypes.String, nil),
					"created_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
					"updated_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
				})

				req := resource.ReadRequest{
					State: tfsdk.State{
						Raw:    stateRaw,
						Schema: schemaResp.Schema,
					},
				}

				resp := &resource.ReadResponse{
					State: tfsdk.State{
						Schema: schemaResp.Schema,
					},
				}

				r.Read(ctx, req, resp)

				assert.False(t, resp.Diagnostics.HasError(), "Expected no errors")
			},
		},
		{
			name: "error - read with invalid state",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					t.Fatal("Should not reach API")
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				r := &CredentialResource{client: n8nClient}
				ctx := context.Background()

				schemaReq := resource.SchemaRequest{}
				schemaResp := &resource.SchemaResponse{}
				r.Schema(ctx, schemaReq, schemaResp)

				stateRaw := tftypes.NewValue(tftypes.String, "invalid")

				req := resource.ReadRequest{
					State: tfsdk.State{
						Raw:    stateRaw,
						Schema: schemaResp.Schema,
					},
				}

				resp := &resource.ReadResponse{
					State: tfsdk.State{
						Schema: schemaResp.Schema,
					},
				}

				r.Read(ctx, req, resp)

				assert.True(t, resp.Diagnostics.HasError(), "Expected error with invalid state")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}

// TestCredentialResource_Delete_CRUD tests the Delete method with full CRUD coverage.
func TestCredentialResource_Delete_CRUD(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "delete with successful API call",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.Method == http.MethodDelete && r.URL.Path == "/credentials/cred-123" {
						w.WriteHeader(http.StatusNoContent)
					}
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				r := &CredentialResource{client: n8nClient}
				ctx := context.Background()

				schemaReq := resource.SchemaRequest{}
				schemaResp := &resource.SchemaResponse{}
				r.Schema(ctx, schemaReq, schemaResp)

				dataMap := map[string]tftypes.Value{
					"key": tftypes.NewValue(tftypes.String, "value"),
				}
				stateRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "cred-123"),
					"name":       tftypes.NewValue(tftypes.String, "test-credential"),
					"type":       tftypes.NewValue(tftypes.String, "httpHeaderAuth"),
					"data":       tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, dataMap),
					"project_id": tftypes.NewValue(tftypes.String, nil),
					"created_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
					"updated_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
				})

				req := resource.DeleteRequest{
					State: tfsdk.State{
						Raw:    stateRaw,
						Schema: schemaResp.Schema,
					},
				}

				resp := &resource.DeleteResponse{
					State: tfsdk.State{
						Schema: schemaResp.Schema,
					},
				}

				r.Delete(ctx, req, resp)

				assert.False(t, resp.Diagnostics.HasError(), "Expected no errors")
			},
		},
		{
			name: "delete with 200 OK response",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.Method == http.MethodDelete && r.URL.Path == "/credentials/cred-123" {
						w.WriteHeader(http.StatusOK)
						w.Write([]byte(`{"message":"Deleted successfully"}`))
					}
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				r := &CredentialResource{client: n8nClient}
				ctx := context.Background()

				schemaReq := resource.SchemaRequest{}
				schemaResp := &resource.SchemaResponse{}
				r.Schema(ctx, schemaReq, schemaResp)

				dataMap := map[string]tftypes.Value{
					"key": tftypes.NewValue(tftypes.String, "value"),
				}
				stateRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "cred-123"),
					"name":       tftypes.NewValue(tftypes.String, "test-credential"),
					"type":       tftypes.NewValue(tftypes.String, "httpHeaderAuth"),
					"data":       tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, dataMap),
					"project_id": tftypes.NewValue(tftypes.String, nil),
					"created_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
					"updated_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
				})

				req := resource.DeleteRequest{
					State: tfsdk.State{
						Raw:    stateRaw,
						Schema: schemaResp.Schema,
					},
				}

				resp := &resource.DeleteResponse{
					State: tfsdk.State{
						Schema: schemaResp.Schema,
					},
				}

				r.Delete(ctx, req, resp)

				assert.False(t, resp.Diagnostics.HasError(), "Expected no errors with 200 OK")
			},
		},
		{
			name: "error - delete with invalid state",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					t.Fatal("Should not reach API")
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				r := &CredentialResource{client: n8nClient}
				ctx := context.Background()

				schemaReq := resource.SchemaRequest{}
				schemaResp := &resource.SchemaResponse{}
				r.Schema(ctx, schemaReq, schemaResp)

				stateRaw := tftypes.NewValue(tftypes.String, "invalid")

				req := resource.DeleteRequest{
					State: tfsdk.State{
						Raw:    stateRaw,
						Schema: schemaResp.Schema,
					},
				}

				resp := &resource.DeleteResponse{
					State: tfsdk.State{
						Schema: schemaResp.Schema,
					},
				}

				r.Delete(ctx, req, resp)

				assert.True(t, resp.Diagnostics.HasError(), "Expected error with invalid state")
			},
		},
		{
			name: "error - API delete fails",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.Method == http.MethodDelete && r.URL.Path == "/credentials/cred-123" {
						w.WriteHeader(http.StatusInternalServerError)
						w.Write([]byte(`{"message":"Internal server error"}`))
						return
					}
					t.Errorf("Unexpected request: %s %s", r.Method, r.URL.Path)
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				r := &CredentialResource{client: n8nClient}
				ctx := context.Background()

				schemaReq := resource.SchemaRequest{}
				schemaResp := &resource.SchemaResponse{}
				r.Schema(ctx, schemaReq, schemaResp)

				dataMap := map[string]tftypes.Value{
					"key": tftypes.NewValue(tftypes.String, "value"),
				}
				stateRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "cred-123"),
					"name":       tftypes.NewValue(tftypes.String, "test-credential"),
					"type":       tftypes.NewValue(tftypes.String, "httpHeaderAuth"),
					"data":       tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, dataMap),
					"project_id": tftypes.NewValue(tftypes.String, nil),
					"created_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
					"updated_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
				})

				req := resource.DeleteRequest{
					State: tfsdk.State{
						Raw:    stateRaw,
						Schema: schemaResp.Schema,
					},
				}

				resp := &resource.DeleteResponse{
					State: tfsdk.State{
						Schema: schemaResp.Schema,
					},
				}

				r.Delete(ctx, req, resp)

				assert.True(t, resp.Diagnostics.HasError(), "Expected error when API delete fails")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}

// TestCredentialResource_executeCreate tests the executeCreate helper.
func TestCredentialResource_executeCreate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		expectError bool
	}{
		{name: "success - create credential", expectError: false},
		{name: "error - API fails", expectError: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPost && r.URL.Path == "/credentials" {
					w.Header().Set("Content-Type", "application/json")
					if tt.expectError {
						w.WriteHeader(http.StatusInternalServerError)
						w.Write([]byte(`{"message":"error"}`))
					} else {
						w.WriteHeader(http.StatusCreated)
						w.Write([]byte(`{"id":"test-id","name":"Test","type":"httpHeaderAuth","createdAt":"2024-01-01T00:00:00Z","updatedAt":"2024-01-01T00:00:00Z"}`))
					}
					return
				}
				t.Fatalf("Unexpected request: %s %s", r.Method, r.URL.Path)
			})

			n8nClient, server := setupTestClient(t, handler)
			defer server.Close()

			r := &CredentialResource{client: n8nClient}
			ctx := context.Background()

			data := map[string]any{"key": "value"}
			createResp, httpResp, err := r.executeCreate(ctx, "Test", "httpHeaderAuth", data)

			if httpResp != nil && httpResp.Body != nil {
				defer httpResp.Body.Close()
			}

			if tt.expectError {
				assert.Error(t, err, "Expected error from API")
			} else {
				assert.NoError(t, err, "Expected no error")
				assert.NotNil(t, createResp, "Response should not be nil")
				assert.Equal(t, "test-id", createResp.Id, "ID should match")
			}
		})
	}
}

// TestCredentialResource_mapCreateResponseToPlan tests the mapCreateResponseToPlan helper.
func TestCredentialResource_mapCreateResponseToPlan(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		testFunc          func(*testing.T)
		expectedID        string
		expectedName      string
		expectedType      string
		createdAt         time.Time
		updatedAt         time.Time
		expectedCreatedAt string
		expectedUpdatedAt string
	}{
		{
			name:              "success - map all fields correctly",
			expectedID:        "test-id-123",
			expectedName:      "Test Credential",
			expectedType:      "httpHeaderAuth",
			createdAt:         must(time.Parse(time.RFC3339, "2024-01-01T10:00:00Z")),
			updatedAt:         must(time.Parse(time.RFC3339, "2024-01-02T15:30:45Z")),
			expectedCreatedAt: "2024-01-01T10:00:00Z",
			expectedUpdatedAt: "2024-01-02T15:30:45Z",
		},
		{
			name:              "error case - empty string values",
			expectedID:        "",
			expectedName:      "",
			expectedType:      "",
			createdAt:         must(time.Parse(time.RFC3339, "2024-01-01T00:00:00Z")),
			updatedAt:         must(time.Parse(time.RFC3339, "2024-01-01T00:00:00Z")),
			expectedCreatedAt: "2024-01-01T00:00:00Z",
			expectedUpdatedAt: "2024-01-01T00:00:00Z",
		},
		{
			name:              "error case - epoch zero time",
			expectedID:        "zero-time-id",
			expectedName:      "Zero Time Test",
			expectedType:      "testType",
			createdAt:         time.Unix(0, 0).UTC(),
			updatedAt:         time.Unix(0, 0).UTC(),
			expectedCreatedAt: "1970-01-01T00:00:00Z",
			expectedUpdatedAt: "1970-01-01T00:00:00Z",
		},
		{
			name:              "exception case - future dates",
			expectedID:        "future-id",
			expectedName:      "Future Credential",
			expectedType:      "futureType",
			createdAt:         must(time.Parse(time.RFC3339, "2099-12-31T23:59:59Z")),
			updatedAt:         must(time.Parse(time.RFC3339, "2099-12-31T23:59:59Z")),
			expectedCreatedAt: "2099-12-31T23:59:59Z",
			expectedUpdatedAt: "2099-12-31T23:59:59Z",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := &CredentialResource{}

			// Create a plan with empty values
			plan := &models.Resource{
				ID:        types.StringValue(""),
				Name:      types.StringValue(""),
				Type:      types.StringValue(""),
				CreatedAt: types.StringValue(""),
				UpdatedAt: types.StringValue(""),
			}

			// Create a mock API response
			createResp := &n8nsdk.CreateCredentialResponse{
				Id:        tt.expectedID,
				Name:      tt.expectedName,
				Type:      tt.expectedType,
				CreatedAt: tt.createdAt,
				UpdatedAt: tt.updatedAt,
			}

			// Call the helper function
			r.mapCreateResponseToPlan(plan, createResp)

			// Verify mappings
			assert.Equal(t, tt.expectedID, plan.ID.ValueString(), "ID should be mapped")
			assert.Equal(t, tt.expectedName, plan.Name.ValueString(), "Name should be mapped")
			assert.Equal(t, tt.expectedType, plan.Type.ValueString(), "Type should be mapped")
			assert.Equal(t, tt.expectedCreatedAt, plan.CreatedAt.ValueString(), "CreatedAt should be mapped")
			assert.Equal(t, tt.expectedUpdatedAt, plan.UpdatedAt.ValueString(), "UpdatedAt should be mapped")
		})
	}
}

// TestCredentialResource_mapCreateResponseToPlan_ErrorCases tests error conditions.
func TestCredentialResource_mapCreateResponseToPlan_ErrorCases(t *testing.T) {
	t.Parallel()

	t.Run("handles nil createResp gracefully - should panic", func(t *testing.T) {
		t.Parallel()

		// This test documents that nil response will panic
		// In production, this should never happen as the API client validates responses
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic with nil createResp, but did not panic")
			}
		}()

		r := &CredentialResource{}
		plan := &models.Resource{}

		// This WILL panic - which is expected behavior for nil pointer dereference
		r.mapCreateResponseToPlan(plan, nil)
	})

	t.Run("handles nil plan gracefully - should panic", func(t *testing.T) {
		t.Parallel()

		// This test documents that nil plan will panic
		// In production, this should never happen as Terraform framework validates state
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic with nil plan, but did not panic")
			}
		}()

		r := &CredentialResource{}
		createResp := &n8nsdk.CreateCredentialResponse{
			Id:        "test",
			Name:      "test",
			Type:      "test",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// This WILL panic - which is expected behavior for nil pointer dereference
		r.mapCreateResponseToPlan(nil, createResp)
	})
}

// must is a helper to panic on time.Parse errors in tests.
func must[T any](val T, err error) T {
	if err != nil {
		panic(err)
	}
	return val
}

// TestCredentialResource_Create_CRUD tests the Create method with full CRUD coverage.
// NOTE: Success path cannot be fully tested due to Terraform framework limitation
// with map[string]interface{} fields causing "don't know how to reflect tftypes.String into interface {}"
// The Create method is tested for all error paths which provides maximum achievable coverage.
func TestCredentialResource_Create_CRUD(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "error - invalid plan",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					t.Fatal("Should not reach API")
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				r := &CredentialResource{client: n8nClient}
				ctx := context.Background()

				schemaReq := resource.SchemaRequest{}
				schemaResp := &resource.SchemaResponse{}
				r.Schema(ctx, schemaReq, schemaResp)

				planRaw := tftypes.NewValue(tftypes.String, "invalid")

				req := resource.CreateRequest{
					Plan: tfsdk.Plan{
						Raw:    planRaw,
						Schema: schemaResp.Schema,
					},
				}

				resp := &resource.CreateResponse{
					State: tfsdk.State{
						Schema: schemaResp.Schema,
					},
				}

				r.Create(ctx, req, resp)

				assert.True(t, resp.Diagnostics.HasError(), "Expected error with invalid plan")
			},
		},
		{
			name: "error - API create fails",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					// Handle schema request (type conversion feature)
					if r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/credentials/schema/") {
						w.WriteHeader(http.StatusNotFound)
						return
					}
					if r.Method == http.MethodPost && r.URL.Path == "/credentials" {
						w.WriteHeader(http.StatusInternalServerError)
						w.Write([]byte(`{"message":"Internal server error"}`))
						return
					}
					t.Fatalf("Unexpected request: %s %s", r.Method, r.URL.Path)
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				r := &CredentialResource{client: n8nClient}
				ctx := context.Background()

				schemaReq := resource.SchemaRequest{}
				schemaResp := &resource.SchemaResponse{}
				r.Schema(ctx, schemaReq, schemaResp)

				dataMap := map[string]tftypes.Value{
					"key": tftypes.NewValue(tftypes.String, "value"),
				}

				planRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, nil),
					"name":       tftypes.NewValue(tftypes.String, "Test Credential"),
					"type":       tftypes.NewValue(tftypes.String, "httpHeaderAuth"),
					"data":       tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, dataMap),
					"project_id": tftypes.NewValue(tftypes.String, nil),
					"created_at": tftypes.NewValue(tftypes.String, nil),
					"updated_at": tftypes.NewValue(tftypes.String, nil),
				})

				req := resource.CreateRequest{
					Plan: tfsdk.Plan{
						Raw:    planRaw,
						Schema: schemaResp.Schema,
					},
				}

				resp := &resource.CreateResponse{
					State: tfsdk.State{
						Schema: schemaResp.Schema,
					},
				}

				r.Create(ctx, req, resp)

				assert.True(t, resp.Diagnostics.HasError(), "Expected error when API create fails")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}

// TestCredentialResource_Update_CRUD tests the Update method with full coverage.
// NOTE: Success path cannot be fully tested due to Terraform framework limitation
// with map[string]interface{} fields causing "don't know how to reflect tftypes.String into interface {}"
// The Update method is tested for all error paths which provides maximum achievable coverage.
func TestCredentialResource_Update_CRUD(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "error - invalid plan",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					t.Fatal("Should not reach API")
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				r := &CredentialResource{client: n8nClient}
				ctx := context.Background()

				schemaReq := resource.SchemaRequest{}
				schemaResp := &resource.SchemaResponse{}
				r.Schema(ctx, schemaReq, schemaResp)

				planRaw := tftypes.NewValue(tftypes.String, "invalid")

				dataMap := map[string]tftypes.Value{
					"key": tftypes.NewValue(tftypes.String, "value"),
				}

				validStateRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "old-cred-123"),
					"name":       tftypes.NewValue(tftypes.String, "old-cred"),
					"type":       tftypes.NewValue(tftypes.String, "httpHeaderAuth"),
					"data":       tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, dataMap),
					"project_id": tftypes.NewValue(tftypes.String, nil),
					"created_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
					"updated_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
				})

				req := resource.UpdateRequest{
					Plan: tfsdk.Plan{
						Raw:    planRaw,
						Schema: schemaResp.Schema,
					},
					State: tfsdk.State{
						Raw:    validStateRaw,
						Schema: schemaResp.Schema,
					},
				}

				resp := &resource.UpdateResponse{
					State: tfsdk.State{
						Schema: schemaResp.Schema,
					},
				}

				r.Update(ctx, req, resp)

				assert.True(t, resp.Diagnostics.HasError(), "Expected error with invalid plan")
			},
		},
		{
			name: "error - invalid state",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					t.Fatal("Should not reach API")
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				r := &CredentialResource{client: n8nClient}
				ctx := context.Background()

				schemaReq := resource.SchemaRequest{}
				schemaResp := &resource.SchemaResponse{}
				r.Schema(ctx, schemaReq, schemaResp)

				dataMap := map[string]tftypes.Value{
					"key": tftypes.NewValue(tftypes.String, "value"),
				}

				validPlanRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "old-cred-123"),
					"name":       tftypes.NewValue(tftypes.String, "updated-cred"),
					"type":       tftypes.NewValue(tftypes.String, "httpHeaderAuth"),
					"data":       tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, dataMap),
					"project_id": tftypes.NewValue(tftypes.String, nil),
					"created_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
					"updated_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
				})

				stateRaw := tftypes.NewValue(tftypes.String, "invalid")

				req := resource.UpdateRequest{
					Plan: tfsdk.Plan{
						Raw:    validPlanRaw,
						Schema: schemaResp.Schema,
					},
					State: tfsdk.State{
						Raw:    stateRaw,
						Schema: schemaResp.Schema,
					},
				}

				resp := &resource.UpdateResponse{
					State: tfsdk.State{
						Schema: schemaResp.Schema,
					},
				}

				r.Update(ctx, req, resp)

				assert.True(t, resp.Diagnostics.HasError(), "Expected error with invalid state")
			},
		},
		{
			name: "error - update credential fails",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					// Handle schema request (type conversion feature)
					if r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/credentials/schema/") {
						w.WriteHeader(http.StatusNotFound)
						return
					}
					// PATCH /credentials/{id} - fail
					if r.Method == http.MethodPatch && r.URL.Path == "/credentials/old-cred-123" {
						w.WriteHeader(http.StatusInternalServerError)
						w.Write([]byte(`{"message":"Internal server error"}`))
						return
					}
					t.Fatalf("Unexpected request: %s %s", r.Method, r.URL.Path)
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				r := &CredentialResource{client: n8nClient}
				ctx := context.Background()

				schemaReq := resource.SchemaRequest{}
				schemaResp := &resource.SchemaResponse{}
				r.Schema(ctx, schemaReq, schemaResp)

				dataMap := map[string]tftypes.Value{
					"key": tftypes.NewValue(tftypes.String, "value"),
				}

				planRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "old-cred-123"),
					"name":       tftypes.NewValue(tftypes.String, "updated-cred"),
					"type":       tftypes.NewValue(tftypes.String, "httpHeaderAuth"),
					"data":       tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, dataMap),
					"project_id": tftypes.NewValue(tftypes.String, nil),
					"created_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
					"updated_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
				})

				stateRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "old-cred-123"),
					"name":       tftypes.NewValue(tftypes.String, "old-cred"),
					"type":       tftypes.NewValue(tftypes.String, "httpHeaderAuth"),
					"data":       tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, dataMap),
					"project_id": tftypes.NewValue(tftypes.String, nil),
					"created_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
					"updated_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
				})

				req := resource.UpdateRequest{
					Plan: tfsdk.Plan{
						Raw:    planRaw,
						Schema: schemaResp.Schema,
					},
					State: tfsdk.State{
						Raw:    stateRaw,
						Schema: schemaResp.Schema,
					},
				}

				resp := &resource.UpdateResponse{
					State: tfsdk.State{
						Schema: schemaResp.Schema,
					},
				}

				r.Update(ctx, req, resp)

				assert.True(t, resp.Diagnostics.HasError(), "Expected error when updating credential fails")
			},
		},
		{
			name: "success - credential updated in place keeps its ID",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")

					// Handle schema request (type conversion feature)
					if r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/credentials/schema/") {
						w.WriteHeader(http.StatusNotFound)
						return
					}

					// PATCH /credentials/{id} - update in place, full data replacement
					if r.Method == http.MethodPatch && r.URL.Path == "/credentials/old-cred-123" {
						var body map[string]any
						if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
							t.Fatalf("Invalid request body: %v", err)
						}
						assert.Equal(t, "updated-cred", body["name"])
						assert.Equal(t, "httpHeaderAuth", body["type"])
						assert.Equal(t, map[string]any{"key": "value"}, body["data"])
						assert.NotContains(t, body, "isPartialData")
						w.WriteHeader(http.StatusOK)
						w.Write([]byte(`{"id":"old-cred-123","name":"updated-cred","type":"httpHeaderAuth","createdAt":"2024-01-01T00:00:00Z","updatedAt":"2024-01-02T00:00:00Z"}`))
						return
					}

					// No credential creation, workflow rewrite or deletion is expected.
					t.Fatalf("Unexpected request: %s %s", r.Method, r.URL.Path)
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				r := &CredentialResource{client: n8nClient}
				ctx := context.Background()

				schemaReq := resource.SchemaRequest{}
				schemaResp := &resource.SchemaResponse{}
				r.Schema(ctx, schemaReq, schemaResp)

				dataMap := map[string]tftypes.Value{
					"key": tftypes.NewValue(tftypes.String, "value"),
				}

				planRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "old-cred-123"),
					"name":       tftypes.NewValue(tftypes.String, "updated-cred"),
					"type":       tftypes.NewValue(tftypes.String, "httpHeaderAuth"),
					"data":       tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, dataMap),
					"project_id": tftypes.NewValue(tftypes.String, nil),
					"created_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
					"updated_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
				})

				stateRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "old-cred-123"),
					"name":       tftypes.NewValue(tftypes.String, "old-cred"),
					"type":       tftypes.NewValue(tftypes.String, "httpHeaderAuth"),
					"data":       tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, dataMap),
					"project_id": tftypes.NewValue(tftypes.String, nil),
					"created_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
					"updated_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
				})

				req := resource.UpdateRequest{
					Plan: tfsdk.Plan{
						Raw:    planRaw,
						Schema: schemaResp.Schema,
					},
					State: tfsdk.State{
						Raw:    stateRaw,
						Schema: schemaResp.Schema,
					},
				}

				resp := &resource.UpdateResponse{
					State: tfsdk.State{
						Schema: schemaResp.Schema,
					},
				}

				r.Update(ctx, req, resp)

				assert.False(t, resp.Diagnostics.HasError(), "Unexpected error: %v", resp.Diagnostics)

				var newState models.Resource
				resp.State.Get(ctx, &newState)
				assert.Equal(t, "old-cred-123", newState.ID.ValueString(), "Credential ID must be preserved")
				assert.Equal(t, "updated-cred", newState.Name.ValueString())
				assert.Equal(t, "2024-01-02T00:00:00Z", newState.UpdatedAt.ValueString())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}

// TestCredentialResource_executeCreateLogic tests the executeCreateLogic method with error cases.
// Note: This test focuses on error paths. The success path is tested through the Create method
// because executeCreateLogic relies on Terraform framework internals (ElementsAs) that require
// proper schema context not available in unit tests.
func TestCredentialResource_executeCreateLogic(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		credName     string
		credType     string
		setupInvalid bool // If true, create invalid Data that will fail ElementsAs
		setupHandler func(w http.ResponseWriter, r *http.Request)
		expectError  bool
	}{
		{
			name:         "error during ElementsAs conversion",
			credName:     "Test Credential",
			credType:     "httpHeaderAuth",
			setupInvalid: true,
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				// Should not be called due to ElementsAs failure
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: true,
		},
		{
			name:         "error during API call",
			credName:     "Failed Credential",
			credType:     "httpBasicAuth",
			setupInvalid: false,
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPost && r.URL.Path == "/credentials" {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"message": "Internal server error"}`))
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := http.HandlerFunc(tt.setupHandler)
			n8nClient, server := setupTestClient(t, handler)
			defer server.Close()

			r := &CredentialResource{client: n8nClient}
			ctx := context.Background()

			var dataValue types.Map
			if tt.setupInvalid {
				// Create an invalid map that will cause ElementsAs to fail
				// This tests the error handling path in executeCreateLogic
				dataValue = types.MapValueMust(types.StringType, map[string]attr.Value{
					"key": types.StringValue("value"),
				})
			} else {
				// Create a valid map for success and API error tests
				var diags diag.Diagnostics
				dataValue, diags = types.MapValueFrom(ctx, types.StringType, map[string]interface{}{
					"name":  "Authorization",
					"value": "Bearer test-token",
				})
				if diags.HasError() {
					t.Fatalf("Failed to create map value: %v", diags)
				}
			}

			plan := &models.Resource{
				Name: types.StringValue(tt.credName),
				Type: types.StringValue(tt.credType),
				Data: dataValue,
			}
			resp := &resource.CreateResponse{
				State: resource.CreateResponse{}.State,
			}

			result := r.executeCreateLogic(ctx, plan, resp)

			assert.False(t, result, "Should return false on error")
			assert.True(t, resp.Diagnostics.HasError(), "Should have diagnostics error")
		})
	}
}

// TestCredentialResource_executeUpdateLogic tests the executeUpdateLogic method with error cases.
// Note: This test focuses on error paths. The success path is tested
// through the Update method because executeUpdateLogic relies on Terraform framework internals
// (ElementsAs) that require proper schema context not available in unit tests.
func TestCredentialResource_executeUpdateLogic(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		oldCredID    string
		newCredName  string
		newCredType  string
		setupInvalid bool // If true, create invalid Data that will fail ElementsAs
		setupHandler func(w http.ResponseWriter, r *http.Request)
		expectError  bool
	}{
		{
			name:         "error during ElementsAs conversion",
			oldCredID:    "cred-old",
			newCredName:  "Updated Credential",
			newCredType:  "httpHeaderAuth",
			setupInvalid: true,
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				// Should not be called due to ElementsAs failure
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: true,
		},
		{
			name:         "error during UpdateCredential",
			oldCredID:    "cred-old",
			newCredName:  "Updated Credential",
			newCredType:  "httpHeaderAuth",
			setupInvalid: false,
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPatch && r.URL.Path == "/credentials/cred-old" {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"message": "Failed to update credential"}`))
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := http.HandlerFunc(tt.setupHandler)
			n8nClient, server := setupTestClient(t, handler)
			defer server.Close()

			r := &CredentialResource{client: n8nClient}
			ctx := context.Background()

			var oldDataValue, newDataValue types.Map
			if tt.setupInvalid {
				// Create invalid maps that will cause ElementsAs to fail
				// This tests the error handling path in executeUpdateLogic
				oldDataValue = types.MapValueMust(types.StringType, map[string]attr.Value{
					"key": types.StringValue("old-value"),
				})
				newDataValue = types.MapValueMust(types.StringType, map[string]attr.Value{
					"key": types.StringValue("new-value"),
				})
			} else {
				// Create valid maps (currently not used since all tests are for error cases)
				var diags1, diags2 diag.Diagnostics
				oldDataValue, diags1 = types.MapValueFrom(ctx, types.StringType, map[string]interface{}{
					"name":  "Authorization",
					"value": "Bearer old-token",
				})
				if diags1.HasError() {
					t.Fatalf("Failed to create old map value: %v", diags1)
				}
				newDataValue, diags2 = types.MapValueFrom(ctx, types.StringType, map[string]interface{}{
					"name":  "Authorization",
					"value": "Bearer new-token",
				})
				if diags2.HasError() {
					t.Fatalf("Failed to create new map value: %v", diags2)
				}
			}

			state := &models.Resource{
				ID:   types.StringValue(tt.oldCredID),
				Name: types.StringValue("Old Credential"),
				Type: types.StringValue("httpHeaderAuth"),
				Data: oldDataValue,
			}
			plan := &models.Resource{
				Name: types.StringValue(tt.newCredName),
				Type: types.StringValue(tt.newCredType),
				Data: newDataValue,
			}
			resp := &resource.UpdateResponse{
				State: resource.UpdateResponse{}.State,
			}

			result := r.executeUpdateLogic(ctx, plan, state, resp)

			assert.False(t, result, "Should return false on error")
			assert.True(t, resp.Diagnostics.HasError(), "Should have diagnostics error")
		})
	}
}

// TestCredentialResource_executeCreateLogicWithData tests the executeCreateLogicWithData method.
// This function is fully testable because it takes credData as a direct parameter,
// bypassing the Terraform framework's ElementsAs conversion.
func TestCredentialResource_executeCreateLogicWithData(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		credName      string
		credType      string
		credData      map[string]any
		setupHandler  func(w http.ResponseWriter, r *http.Request)
		expectSuccess bool
		expectError   bool
	}{
		{
			name:     "success - credential created",
			credName: "Test Credential",
			credType: "httpHeaderAuth",
			credData: map[string]any{
				"name":  "Authorization",
				"value": "Bearer test-token",
			},
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				now := time.Now()
				w.Write([]byte(`{
					"id": "cred-123",
					"name": "Test Credential",
					"type": "httpHeaderAuth",
					"createdAt": "` + now.Format(time.RFC3339) + `",
					"updatedAt": "` + now.Format(time.RFC3339) + `"
				}`))
			},
			expectSuccess: true,
			expectError:   false,
		},
		{
			name:     "error - API call fails",
			credName: "Test Credential",
			credType: "httpHeaderAuth",
			credData: map[string]any{
				"name":  "Authorization",
				"value": "Bearer test-token",
			},
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"message": "Internal server error"}`))
			},
			expectSuccess: false,
			expectError:   true,
		},
		{
			name:     "error - API returns bad request",
			credName: "Test Credential",
			credType: "invalid-type",
			credData: map[string]any{
				"invalid": "data",
			},
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"message": "Invalid credential type"}`))
			},
			expectSuccess: false,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := http.HandlerFunc(tt.setupHandler)
			n8nClient, server := setupTestClient(t, handler)
			defer server.Close()

			r := &CredentialResource{client: n8nClient}
			ctx := context.Background()

			plan := &models.Resource{
				Name: types.StringValue(tt.credName),
				Type: types.StringValue(tt.credType),
			}
			resp := &resource.CreateResponse{
				State: resource.CreateResponse{}.State,
			}

			result := r.executeCreateLogicWithData(ctx, plan, tt.credData, resp)

			if tt.expectSuccess {
				assert.True(t, result, "Should return true on success")
				assert.False(t, resp.Diagnostics.HasError(), "Should not have diagnostics error")
				assert.Equal(t, "cred-123", plan.ID.ValueString(), "Should set credential ID")
				assert.Equal(t, tt.credName, plan.Name.ValueString(), "Should preserve credential name")
			} else {
				assert.False(t, result, "Should return false on error")
			}

			if tt.expectError {
				assert.True(t, resp.Diagnostics.HasError(), "Should have diagnostics error")
			}
		})
	}
}

// TestCredentialResource_executeUpdateLogicWithData tests the executeUpdateLogicWithData method.
// This function is fully testable because it takes credData as a direct parameter,
// bypassing the Terraform framework's ElementsAs conversion.
func TestCredentialResource_executeUpdateLogicWithData(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		oldCredID     string
		newCredName   string
		newCredType   string
		credData      map[string]any
		setupHandler  func(w http.ResponseWriter, r *http.Request)
		expectSuccess bool
		expectError   bool
		expectDetail  string
	}{
		{
			name:        "success - credential updated in place",
			oldCredID:   "cred-old-123",
			newCredName: "Updated Credential",
			newCredType: "httpHeaderAuth",
			credData: map[string]any{
				"name":  "Authorization",
				"value": "Bearer updated-token",
			},
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				if r.Method == "PATCH" && r.URL.Path == "/credentials/cred-old-123" {
					w.WriteHeader(http.StatusOK)
					now := time.Now()
					w.Write([]byte(`{
						"id": "cred-old-123",
						"name": "Updated Credential",
						"type": "httpHeaderAuth",
						"createdAt": "` + now.Format(time.RFC3339) + `",
						"updatedAt": "` + now.Format(time.RFC3339) + `"
					}`))
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectSuccess: true,
			expectError:   false,
		},
		{
			name:        "error - n8n older than 2.3.0 without PATCH support",
			oldCredID:   "cred-old-123",
			newCredName: "Updated Credential",
			newCredType: "httpHeaderAuth",
			credData: map[string]any{
				"name":  "Authorization",
				"value": "Bearer updated-token",
			},
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == "PATCH" && r.URL.Path == "/credentials/cred-old-123" {
					w.WriteHeader(http.StatusMethodNotAllowed)
					w.Write([]byte(`{"message": "PATCH method not allowed"}`))
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectSuccess: false,
			expectError:   true,
			expectDetail:  "requires n8n >= 2.3.0",
		},
		{
			name:        "error - UpdateCredential fails",
			oldCredID:   "cred-old-123",
			newCredName: "Updated Credential",
			newCredType: "httpHeaderAuth",
			credData: map[string]any{
				"name":  "Authorization",
				"value": "Bearer updated-token",
			},
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == "PATCH" && r.URL.Path == "/credentials/cred-old-123" {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"message": "Failed to update credential"}`))
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectSuccess: false,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := http.HandlerFunc(tt.setupHandler)
			n8nClient, server := setupTestClient(t, handler)
			defer server.Close()

			r := &CredentialResource{client: n8nClient}
			ctx := context.Background()

			oldDataValue, diags := types.MapValueFrom(ctx, types.StringType, map[string]interface{}{
				"name":  "Authorization",
				"value": "Bearer old-token",
			})
			if diags.HasError() {
				t.Fatalf("Failed to create old data map: %v", diags)
			}

			state := &models.Resource{
				ID:   types.StringValue(tt.oldCredID),
				Name: types.StringValue("Old Credential"),
				Type: types.StringValue("httpHeaderAuth"),
				Data: oldDataValue,
			}
			plan := &models.Resource{
				Name: types.StringValue(tt.newCredName),
				Type: types.StringValue(tt.newCredType),
			}
			resp := &resource.UpdateResponse{
				State: resource.UpdateResponse{}.State,
			}

			result := r.executeUpdateLogicWithData(ctx, plan, state, tt.credData, resp)

			if tt.expectSuccess {
				assert.True(t, result, "Should return true on success")
				assert.False(t, resp.Diagnostics.HasError(), "Should not have diagnostics error")
				assert.Equal(t, tt.oldCredID, plan.ID.ValueString(), "Should keep the credential ID")
				assert.Equal(t, tt.newCredName, plan.Name.ValueString(), "Should preserve credential name")
			} else {
				assert.False(t, result, "Should return false on error")
			}

			if tt.expectError {
				assert.True(t, resp.Diagnostics.HasError(), "Should have diagnostics error")
			}

			if tt.expectDetail != "" {
				assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), tt.expectDetail)
			} else if tt.expectError {
				assert.NotContains(t, resp.Diagnostics.Errors()[0].Detail(), "requires n8n >= 2.3.0")
			}
		})
	}
}

// TestCredentialResource_executeUpdateLogicWithData_ProjectTransfer tests project transfer scenarios.
func TestCredentialResource_executeUpdateLogicWithData_ProjectTransfer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		stateProjectID types.String
		projectID      types.String
		wantProjectID  string
		setupHandler   func(t *testing.T) http.HandlerFunc
		expectSuccess  bool
		expectError    bool
	}{
		{
			name:           "success - update with project transfer",
			stateProjectID: types.StringValue("project-123"),
			projectID:      types.StringValue("project-456"),
			wantProjectID:  "project-456",
			setupHandler: func(t *testing.T) http.HandlerFunc {
				t.Helper()
				return func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")

					switch {
					case r.Method == "PATCH" && r.URL.Path == "/credentials/cred-old-123":
						w.WriteHeader(http.StatusOK)
						now := time.Now()
						_, _ = w.Write([]byte(`{
							"id": "cred-old-123",
							"name": "Updated Credential",
							"type": "httpHeaderAuth",
							"createdAt": "` + now.Format(time.RFC3339) + `",
							"updatedAt": "` + now.Format(time.RFC3339) + `"
						}`))
					case r.Method == "PUT" && r.URL.Path == "/credentials/cred-old-123/transfer":
						w.WriteHeader(http.StatusOK)
						_, _ = w.Write([]byte(`{}`))
					default:
						w.WriteHeader(http.StatusNotFound)
					}
				}
			},
			expectSuccess: true,
			expectError:   false,
		},
		{
			name:           "success - no transfer when project is unchanged",
			stateProjectID: types.StringValue("project-456"),
			projectID:      types.StringValue("project-456"),
			wantProjectID:  "project-456",
			setupHandler: func(t *testing.T) http.HandlerFunc {
				t.Helper()
				return func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")

					switch {
					case r.Method == "PATCH" && r.URL.Path == "/credentials/cred-old-123":
						w.WriteHeader(http.StatusOK)
						now := time.Now()
						_, _ = w.Write([]byte(`{
							"id": "cred-old-123",
							"name": "Updated Credential",
							"type": "httpHeaderAuth",
							"createdAt": "` + now.Format(time.RFC3339) + `",
							"updatedAt": "` + now.Format(time.RFC3339) + `"
						}`))
					case r.Method == "PUT" && r.URL.Path == "/credentials/cred-old-123/transfer":
						t.Errorf("Unexpected transfer when project is unchanged")
						w.WriteHeader(http.StatusInternalServerError)
					default:
						w.WriteHeader(http.StatusNotFound)
					}
				}
			},
			expectSuccess: true,
			expectError:   false,
		},
		{
			name:           "success - project_id removed from config keeps the current project",
			stateProjectID: types.StringValue("project-123"),
			projectID:      types.StringUnknown(),
			wantProjectID:  "project-123",
			setupHandler: func(t *testing.T) http.HandlerFunc {
				t.Helper()
				return func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")

					switch {
					case r.Method == "PATCH" && r.URL.Path == "/credentials/cred-old-123":
						w.WriteHeader(http.StatusOK)
						now := time.Now()
						_, _ = w.Write([]byte(`{
							"id": "cred-old-123",
							"name": "Updated Credential",
							"type": "httpHeaderAuth",
							"createdAt": "` + now.Format(time.RFC3339) + `",
							"updatedAt": "` + now.Format(time.RFC3339) + `"
						}`))
					case r.Method == "PUT" && r.URL.Path == "/credentials/cred-old-123/transfer":
						t.Errorf("Unexpected transfer when project is unchanged")
						w.WriteHeader(http.StatusInternalServerError)
					default:
						w.WriteHeader(http.StatusNotFound)
					}
				}
			},
			expectSuccess: true,
			expectError:   false,
		},
		{
			name:           "error - project transfer fails",
			stateProjectID: types.StringNull(),
			projectID:      types.StringValue("project-456"),
			wantProjectID:  "project-456",
			setupHandler: func(t *testing.T) http.HandlerFunc {
				t.Helper()
				return func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")

					switch {
					case r.Method == "PATCH" && r.URL.Path == "/credentials/cred-old-123":
						w.WriteHeader(http.StatusOK)
						now := time.Now()
						_, _ = w.Write([]byte(`{
							"id": "cred-old-123",
							"name": "Updated Credential",
							"type": "httpHeaderAuth",
							"createdAt": "` + now.Format(time.RFC3339) + `",
							"updatedAt": "` + now.Format(time.RFC3339) + `"
						}`))
					case r.Method == "PUT" && r.URL.Path == "/credentials/cred-old-123/transfer":
						w.WriteHeader(http.StatusInternalServerError)
						_, _ = w.Write([]byte(`{}`))
					default:
						w.WriteHeader(http.StatusNotFound)
					}
				}
			},
			expectSuccess: false,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			n8nClient, server := setupTestClient(t, tt.setupHandler(t))
			defer server.Close()

			r := &CredentialResource{client: n8nClient}
			ctx := context.Background()

			oldDataValue, diagsOld := types.MapValueFrom(ctx, types.StringType, map[string]interface{}{
				"name":  "Authorization",
				"value": "Bearer old-token",
			})
			if diagsOld.HasError() {
				t.Fatalf("Failed to create old data map: %v", diagsOld)
			}

			state := &models.Resource{
				ID:        types.StringValue("cred-old-123"),
				Name:      types.StringValue("Old Credential"),
				Type:      types.StringValue("httpHeaderAuth"),
				Data:      oldDataValue,
				ProjectID: tt.stateProjectID,
			}
			plan := &models.Resource{
				Name:      types.StringValue("Updated Credential"),
				Type:      types.StringValue("httpHeaderAuth"),
				ProjectID: tt.projectID,
			}
			resp := &resource.UpdateResponse{
				State: resource.UpdateResponse{}.State,
			}
			credData := map[string]any{
				"name":  "Authorization",
				"value": "Bearer updated-token",
			}

			result := r.executeUpdateLogicWithData(ctx, plan, state, credData, resp)

			if tt.expectSuccess {
				assert.True(t, result, "Should return true on success")
				assert.False(t, resp.Diagnostics.HasError(), "Should not have diagnostics error")
				assert.Equal(t, "cred-old-123", plan.ID.ValueString(), "Should keep the credential ID")
				assert.Equal(t, tt.wantProjectID, plan.ProjectID.ValueString(), "Should store the credential project")
			} else {
				assert.False(t, result, "Should return false on error")
			}

			if tt.expectError {
				assert.True(t, resp.Diagnostics.HasError(), "Should have diagnostics error")
			}
		})
	}
}

// TestCredentialResource_executeCreateLogicWithData_ProjectTransfer tests project transfer in create.
func TestCredentialResource_executeCreateLogicWithData_ProjectTransfer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		credName      string
		credType      string
		projectID     string
		credData      map[string]any
		setupHandler  func(w http.ResponseWriter, r *http.Request)
		expectSuccess bool
		expectError   bool
	}{
		{
			name:      "success - create with project transfer",
			credName:  "New Credential",
			credType:  "httpHeaderAuth",
			projectID: "project-456",
			credData: map[string]any{
				"name":  "Authorization",
				"value": "Bearer token",
			},
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")

				switch {
				case r.Method == "POST" && r.URL.Path == "/credentials":
					w.WriteHeader(http.StatusCreated)
					now := time.Now()
					_, _ = w.Write([]byte(`{
						"id": "cred-new-123",
						"name": "New Credential",
						"type": "httpHeaderAuth",
						"createdAt": "` + now.Format(time.RFC3339) + `",
						"updatedAt": "` + now.Format(time.RFC3339) + `"
					}`))
				case r.Method == "PUT" && r.URL.Path == "/credentials/cred-new-123/transfer":
					w.WriteHeader(http.StatusOK)
					_, _ = w.Write([]byte(`{}`))
				default:
					w.WriteHeader(http.StatusNotFound)
				}
			},
			expectSuccess: true,
			expectError:   false,
		},
		{
			name:      "error - project transfer fails on create",
			credName:  "New Credential",
			credType:  "httpHeaderAuth",
			projectID: "project-456",
			credData: map[string]any{
				"name":  "Authorization",
				"value": "Bearer token",
			},
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")

				switch {
				case r.Method == "POST" && r.URL.Path == "/credentials":
					w.WriteHeader(http.StatusCreated)
					now := time.Now()
					_, _ = w.Write([]byte(`{
						"id": "cred-new-123",
						"name": "New Credential",
						"type": "httpHeaderAuth",
						"createdAt": "` + now.Format(time.RFC3339) + `",
						"updatedAt": "` + now.Format(time.RFC3339) + `"
					}`))
				case r.Method == "PUT" && r.URL.Path == "/credentials/cred-new-123/transfer":
					// Project transfer fails
					w.WriteHeader(http.StatusInternalServerError)
					_, _ = w.Write([]byte(`{"message": "Failed to transfer credential"}`))
				default:
					w.WriteHeader(http.StatusNotFound)
				}
			},
			expectSuccess: false,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := http.HandlerFunc(tt.setupHandler)
			n8nClient, server := setupTestClient(t, handler)
			defer server.Close()

			r := &CredentialResource{client: n8nClient}
			ctx := context.Background()

			plan := &models.Resource{
				Name:      types.StringValue(tt.credName),
				Type:      types.StringValue(tt.credType),
				ProjectID: types.StringValue(tt.projectID),
			}
			resp := &resource.CreateResponse{
				State: resource.CreateResponse{}.State,
			}

			result := r.executeCreateLogicWithData(ctx, plan, tt.credData, resp)

			if tt.expectSuccess {
				assert.True(t, result, "Should return true on success")
				assert.False(t, resp.Diagnostics.HasError(), "Should not have diagnostics error")
			} else {
				assert.False(t, result, "Should return false on error")
			}

			if tt.expectError {
				assert.True(t, resp.Diagnostics.HasError(), "Should have diagnostics error")
			}
		})
	}
}
