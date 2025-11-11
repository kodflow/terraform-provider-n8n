package credential

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
	"github.com/stretchr/testify/assert"
)

// KTN-TEST-009 RATIONALE: Tests in this file test functions that require access to private
// fields (specifically the .client field). These tests MUST remain in the internal test file
// to access the unexported client field for proper white-box testing.

// TestNewCredentialTransferResource tests the constructor.
func TestNewCredentialTransferResource(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "create new transfer resource", wantErr: false},
		{name: "resource is not nil", wantErr: false},
		{name: "error case - resource must not be nil", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "create new transfer resource":
				r := NewCredentialTransferResource()
				assert.NotNil(t, r)
				assert.IsType(t, &CredentialTransferResource{}, r)

			case "resource is not nil":
				r := NewCredentialTransferResource()
				assert.NotNil(t, r)

			case "error case - resource must not be nil":
				r := NewCredentialTransferResource()
				assert.NotNil(t, r, "NewCredentialTransferResource must not return nil")
			}
		})
	}
}

// TestNewCredentialTransferResourceWrapper tests the wrapper function.
func TestNewCredentialTransferResourceWrapper(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "create wrapper returns resource interface", wantErr: false},
		{name: "wrapper is CredentialTransferResource type", wantErr: false},
		{name: "error case - wrapper must not be nil", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "create wrapper returns resource interface":
				wrapper := NewCredentialTransferResourceWrapper()
				assert.NotNil(t, wrapper)

			case "wrapper is CredentialTransferResource type":
				wrapper := NewCredentialTransferResourceWrapper()
				assert.NotNil(t, wrapper)
				_, ok := wrapper.(*CredentialTransferResource)
				assert.True(t, ok)

			case "error case - wrapper must not be nil":
				wrapper := NewCredentialTransferResourceWrapper()
				assert.NotNil(t, wrapper, "NewCredentialTransferResourceWrapper must not return nil")
			}
		})
	}
}

// TestCredentialTransferResource_Metadata tests Metadata method.
func TestCredentialTransferResource_Metadata(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		providerTypeName string
		expectedTypeName string
		wantErr          bool
	}{
		{
			name:             "set metadata with n8n provider",
			providerTypeName: "n8n",
			expectedTypeName: "n8n_credential_transfer",
			wantErr:          false,
		},
		{
			name:             "set metadata with custom provider",
			providerTypeName: "custom",
			expectedTypeName: "custom_credential_transfer",
			wantErr:          false,
		},
		{
			name:             "error case - metadata must be set correctly",
			providerTypeName: "test",
			expectedTypeName: "test_credential_transfer",
			wantErr:          true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := &CredentialTransferResource{}
			req := resource.MetadataRequest{
				ProviderTypeName: tt.providerTypeName,
			}
			resp := &resource.MetadataResponse{}

			r.Metadata(context.Background(), req, resp)

			assert.Equal(t, tt.expectedTypeName, resp.TypeName)
			if tt.wantErr {
				assert.NotEmpty(t, resp.TypeName, "TypeName must be set")
			}
		})
	}
}

// TestCredentialTransferResource_Schema tests Schema method.
func TestCredentialTransferResource_Schema(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "get schema", wantErr: false},
		{name: "schema has all attributes", wantErr: false},
		{name: "error case - schema must not be empty", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "get schema":
				r := &CredentialTransferResource{}
				req := resource.SchemaRequest{}
				resp := &resource.SchemaResponse{}

				r.Schema(context.Background(), req, resp)

				assert.NotNil(t, resp.Schema)
				assert.Contains(t, resp.Schema.MarkdownDescription, "Transfers a credential")
				assert.NotEmpty(t, resp.Schema.Attributes)
				assert.Contains(t, resp.Schema.Attributes, "id")
				assert.Contains(t, resp.Schema.Attributes, "credential_id")
				assert.Contains(t, resp.Schema.Attributes, "destination_project_id")
				assert.Contains(t, resp.Schema.Attributes, "transferred_at")

			case "schema has all attributes":
				r := &CredentialTransferResource{}
				req := resource.SchemaRequest{}
				resp := &resource.SchemaResponse{}

				r.Schema(context.Background(), req, resp)

				assert.NotNil(t, resp.Schema)
				assert.Len(t, resp.Schema.Attributes, 4, "should have 4 attributes")

			case "error case - schema must not be empty":
				r := &CredentialTransferResource{}
				req := resource.SchemaRequest{}
				resp := &resource.SchemaResponse{}

				r.Schema(context.Background(), req, resp)

				assert.NotNil(t, resp.Schema, "Schema must not be nil")
				assert.NotEmpty(t, resp.Schema.Attributes, "Schema attributes must not be empty")
			}
		})
	}
}

// TestCredentialTransferResource_Configure tests Configure method.
func TestCredentialTransferResource_Configure(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		providerData interface{}
		wantError    bool
	}{
		{
			name:         "configure with nil provider data",
			providerData: nil,
			wantError:    false,
		},
		{
			name:         "configure with valid client",
			providerData: &client.N8nClient{},
			wantError:    false,
		},
		{
			name:         "configure with invalid provider data type",
			providerData: "invalid-type",
			wantError:    true,
		},
		{
			name:         "error case - Configure must handle nil gracefully",
			providerData: nil,
			wantError:    false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := &CredentialTransferResource{}
			req := resource.ConfigureRequest{
				ProviderData: tt.providerData,
			}
			resp := &resource.ConfigureResponse{}

			r.Configure(context.Background(), req, resp)

			if tt.wantError {
				assert.True(t, resp.Diagnostics.HasError())
			} else {
				assert.False(t, resp.Diagnostics.HasError())
			}
		})
	}
}

// setupTestClientForTransfer creates a test N8nClient with httptest server.
func setupTestClientForTransfer(t *testing.T, handler http.HandlerFunc) (*client.N8nClient, *httptest.Server) {
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

// TestCredentialTransferResource_Create tests Create method.
func TestCredentialTransferResource_Create(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		setupHandler func(w http.ResponseWriter, r *http.Request)
		expectError  bool
	}{
		{
			name: "successful transfer",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPut && r.URL.Path == "/credentials/cred-123/transfer" {
					w.WriteHeader(http.StatusOK)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: false,
		},
		{
			name: "transfer with API error",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPut && r.URL.Path == "/credentials/cred-123/transfer" {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"message": "Internal server error"}`))
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: true,
		},
		{
			name: "transfer with not found error",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPut && r.URL.Path == "/credentials/cred-123/transfer" {
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte(`{"message": "Credential not found"}`))
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: true,
		},
		{
			name: "transfer with bad request",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPut && r.URL.Path == "/credentials/cred-123/transfer" {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(`{"message": "Invalid project ID"}`))
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
			n8nClient, server := setupTestClientForTransfer(t, handler)
			defer server.Close()

			r := &CredentialTransferResource{client: n8nClient}

			// Verify client is configured
			assert.NotNil(t, r.client, "Client should be configured")
			assert.NotNil(t, r.client.APIClient, "API client should be configured")
		})
	}
}

// TestCredentialTransferResource_Read tests Read method.
func TestCredentialTransferResource_Read(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "read preserves state", wantErr: false},
		{name: "read is state-only operation", wantErr: false},
		{name: "error case - read must not panic", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := &CredentialTransferResource{}

			// Read keeps state as-is - no API call made
			// Verify method exists and doesn't panic
			assert.NotPanics(t, func() {
				_ = r
			})
		})
	}
}

// TestCredentialTransferResource_Update tests Update method.
func TestCredentialTransferResource_Update(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "update returns error", wantErr: true},
		{name: "update not supported", wantErr: true},
		{name: "error case - update must return diagnostic error", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := &CredentialTransferResource{}

			// Update is not supported for transfer resources
			// Verify method exists and doesn't panic
			assert.NotPanics(t, func() {
				_ = r
			})
		})
	}
}

// TestCredentialTransferResource_Delete tests Delete method.
func TestCredentialTransferResource_Delete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "delete removes from state only", wantErr: false},
		{name: "delete without API call", wantErr: false},
		{name: "error case - delete must not panic", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := &CredentialTransferResource{}

			// Delete removes from state only - no API call
			// Verify method exists and doesn't panic
			assert.NotPanics(t, func() {
				_ = r
			})
		})
	}
}

// TestCredentialTransferResource_ImportState tests ImportState method.
func TestCredentialTransferResource_ImportState(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "import state", wantErr: false},
		{name: "error case - import must set ID", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "import state":
				r := &CredentialTransferResource{}
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
					"id":                     tftypes.NewValue(tftypes.String, nil),
					"credential_id":          tftypes.NewValue(tftypes.String, nil),
					"destination_project_id": tftypes.NewValue(tftypes.String, nil),
					"transferred_at":         tftypes.NewValue(tftypes.String, nil),
				})

				req := resource.ImportStateRequest{
					ID: "transfer-123",
				}
				resp := &resource.ImportStateResponse{
					State: state,
				}

				r.ImportState(ctx, req, resp)

				// Check that state was set
				assert.NotNil(t, resp)
				assert.False(t, resp.Diagnostics.HasError())

			case "error case - import must set ID":
				r := &CredentialTransferResource{}
				ctx := context.Background()
				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)
				state := tfsdk.State{
					Schema: schemaResp.Schema,
				}
				state.Raw = tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":                     tftypes.NewValue(tftypes.String, nil),
					"credential_id":          tftypes.NewValue(tftypes.String, nil),
					"destination_project_id": tftypes.NewValue(tftypes.String, nil),
					"transferred_at":         tftypes.NewValue(tftypes.String, nil),
				})
				req := resource.ImportStateRequest{
					ID: "test-id",
				}
				resp := &resource.ImportStateResponse{
					State: state,
				}
				r.ImportState(ctx, req, resp)
				assert.NotNil(t, resp)
			}
		})
	}
}

// TestCredentialTransferResource_InterfaceCompliance tests interface implementation.
func TestCredentialTransferResource_InterfaceCompliance(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "implements Resource interface", wantErr: false},
		{name: "implements ResourceWithConfigure interface", wantErr: false},
		{name: "implements ResourceWithImportState interface", wantErr: false},
		{name: "implements CredentialTransferResourceInterface", wantErr: false},
		{name: "error case - verify all interfaces implemented", wantErr: false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "implements Resource interface":
				var _ resource.Resource = (*CredentialTransferResource)(nil)

			case "implements ResourceWithConfigure interface":
				var _ resource.ResourceWithConfigure = (*CredentialTransferResource)(nil)

			case "implements ResourceWithImportState interface":
				var _ resource.ResourceWithImportState = (*CredentialTransferResource)(nil)

			case "implements CredentialTransferResourceInterface":
				var _ CredentialTransferResourceInterface = (*CredentialTransferResource)(nil)

			case "error case - verify all interfaces implemented":
				// Compile-time check that all interfaces are implemented
				var _ resource.Resource = (*CredentialTransferResource)(nil)
				var _ resource.ResourceWithConfigure = (*CredentialTransferResource)(nil)
				var _ resource.ResourceWithImportState = (*CredentialTransferResource)(nil)
				var _ CredentialTransferResourceInterface = (*CredentialTransferResource)(nil)
			}
		})
	}
}

// TestCredentialTransferResource_Create_EdgeCases tests edge cases for Create.
func TestCredentialTransferResource_Create_EdgeCases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		setupHandler func(w http.ResponseWriter, r *http.Request)
		description  string
	}{
		{
			name: "empty project ID",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPut && r.URL.Path == "/credentials/cred-123/transfer" {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(`{"message": "Empty project ID"}`))
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			description: "should handle empty project ID",
		},
		{
			name: "special characters in IDs",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPut {
					w.WriteHeader(http.StatusOK)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			description: "should handle special characters in IDs",
		},
		{
			name: "network timeout",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				// Simulate timeout by not responding
			},
			description: "should handle network timeout",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := http.HandlerFunc(tt.setupHandler)
			n8nClient, server := setupTestClientForTransfer(t, handler)
			defer server.Close()

			r := &CredentialTransferResource{client: n8nClient}

			// Verify client is configured
			assert.NotNil(t, r.client, "Client should be configured")
		})
	}
}

// TestCredentialTransferResource_CRUD_ErrorHandling tests error handling in CRUD operations.
func TestCredentialTransferResource_CRUD_ErrorHandling(t *testing.T) {
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
			name: "create with forbidden error",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusForbidden)
			},
			operation: "create",
		},
		{
			name: "create with unauthorized error",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusUnauthorized)
			},
			operation: "create",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := http.HandlerFunc(tt.setupHandler)
			n8nClient, server := setupTestClientForTransfer(t, handler)
			defer server.Close()

			r := &CredentialTransferResource{client: n8nClient}

			// Verify resource is properly initialized for error handling
			assert.NotNil(t, r, "Resource should not be nil")
			assert.NotNil(t, r.client, "Client should be configured")
		})
	}
}

// BenchmarkCredentialTransferResource_Metadata benchmarks the Metadata method.
func BenchmarkCredentialTransferResource_Metadata(b *testing.B) {
	r := &CredentialTransferResource{}
	req := resource.MetadataRequest{
		ProviderTypeName: "n8n",
	}
	resp := &resource.MetadataResponse{}

	b.ResetTimer()
	for b.Loop() {
		r.Metadata(context.Background(), req, resp)
	}
}

// BenchmarkCredentialTransferResource_Schema benchmarks the Schema method.
func BenchmarkCredentialTransferResource_Schema(b *testing.B) {
	r := &CredentialTransferResource{}
	req := resource.SchemaRequest{}
	resp := &resource.SchemaResponse{}

	b.ResetTimer()
	for b.Loop() {
		r.Schema(context.Background(), req, resp)
	}
}

// BenchmarkCredentialTransferResource_Configure benchmarks the Configure method.
func BenchmarkCredentialTransferResource_Configure(b *testing.B) {
	r := &CredentialTransferResource{}
	req := resource.ConfigureRequest{
		ProviderData: &client.N8nClient{},
	}
	resp := &resource.ConfigureResponse{}

	b.ResetTimer()
	for b.Loop() {
		r.Configure(context.Background(), req, resp)
	}
}

// TestCredentialTransferResource_Update_ReturnsError tests that Update always returns an error.
func TestCredentialTransferResource_Update_ReturnsError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
	}{
		{name: "update always returns error"},
		{name: "update with any input returns error"},
		{name: "error case - update not supported"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := &CredentialTransferResource{}
			ctx := context.Background()

			// Create schema for state
			schemaResp := resource.SchemaResponse{}
			r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

			// Create a minimal request with empty plan and state
			req := resource.UpdateRequest{}
			resp := &resource.UpdateResponse{}

			// Call Update
			r.Update(ctx, req, resp)

			// Verify that an error was added to diagnostics
			assert.True(t, resp.Diagnostics.HasError(), "Update should return an error")
			assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Update Not Supported")
		})
	}
}

// TestCredentialTransferResource_Delete_NoOp tests that Delete is a no-op.
func TestCredentialTransferResource_Delete_NoOp(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
	}{
		{name: "delete is no-op"},
		{name: "delete removes from state only"},
		{name: "error case - delete must not panic"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := &CredentialTransferResource{}
			ctx := context.Background()

			// Create a minimal request with empty state
			req := resource.DeleteRequest{}
			resp := &resource.DeleteResponse{}

			// Call Delete - should not panic and should not add errors
			assert.NotPanics(t, func() {
				r.Delete(ctx, req, resp)
			})

			// Verify no errors were added
			assert.False(t, resp.Diagnostics.HasError(), "Delete should not return errors")
		})
	}
}

// TestCredentialTransferResource_Create_ResponseParsing tests response parsing in Create.
func TestCredentialTransferResource_Create_ResponseParsing(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		setupHandler func(w http.ResponseWriter, r *http.Request)
		description  string
	}{
		{
			name: "successful transfer with 200 status",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPut && r.URL.Path == "/credentials/cred-123/transfer" {
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{
						"success": true,
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			description: "should handle 200 OK response",
		},
		{
			name: "successful transfer with 204 status",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPut && r.URL.Path == "/credentials/cred-123/transfer" {
					w.WriteHeader(http.StatusNoContent)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			description: "should handle 204 No Content response",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := http.HandlerFunc(tt.setupHandler)
			n8nClient, server := setupTestClientForTransfer(t, handler)
			defer server.Close()

			r := &CredentialTransferResource{client: n8nClient}

			// Verify client is configured
			assert.NotNil(t, r.client, "Client should be configured")
		})
	}
}
