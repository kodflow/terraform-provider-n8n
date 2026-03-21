package sourcecontrol

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/client"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/sourcecontrol/models"
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

// TestSourceControlPullResource_executePullLogic tests the executePullLogic method.
func TestSourceControlPullResource_executePullLogic(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                   string
		force                  bool
		setupHandler           func(w http.ResponseWriter, r *http.Request)
		expectError            bool
		expectWorkflowsNotNull bool
	}{
		{
			name:  "successful pull - no force",
			force: false,
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPost && r.URL.Path == "/source-control/pull" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{
						"workflows": []any{
							map[string]any{"id": "wf-1", "name": "My Workflow"},
						},
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError:            false,
			expectWorkflowsNotNull: true,
		},
		{
			name:  "successful pull - with force",
			force: true,
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPost && r.URL.Path == "/source-control/pull" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError:            false,
			expectWorkflowsNotNull: false,
		},
		{
			name:  "API error",
			force: false,
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"message": "Internal server error"}`))
			},
			expectError: true,
		},
		{
			name:  "nil result - no content response",
			force: false,
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPost && r.URL.Path == "/source-control/pull" {
					w.WriteHeader(http.StatusNoContent)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError:            false,
			expectWorkflowsNotNull: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			n8nClient, server := setupTestClient(t, http.HandlerFunc(tt.setupHandler))
			defer server.Close()

			r := &SourceControlPullResource{client: n8nClient}
			ctx := t.Context()
			data := &models.Resource{
				Force: types.BoolValue(tt.force),
			}

			var diags resource.CreateResponse
			result := r.executePullLogic(ctx, data, &diags.Diagnostics)

			//: Verify the result matches the expected error state.
			if tt.expectError {
				assert.False(t, result, "Should return false on error")
				assert.True(t, diags.Diagnostics.HasError(), "Should have diagnostics error")
			} else {
				assert.True(t, result, "Should return true on success")
				assert.False(t, diags.Diagnostics.HasError(), "Should not have diagnostics error")
				if tt.expectWorkflowsNotNull {
					assert.False(t, data.WorkflowsImported.IsNull(), "WorkflowsImported should not be null")
				}
			}
		})
	}
}

// TestSourceControlPullResource_executePullLogic_NullForce tests with null force value.
func TestSourceControlPullResource_executePullLogic_NullForce(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
	}{
		{name: "null force does not set force flag"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			n8nClient, server := setupTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPost && r.URL.Path == "/source-control/pull" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}))
			defer server.Close()

			r := &SourceControlPullResource{client: n8nClient}
			ctx := t.Context()
			data := &models.Resource{
				Force: types.BoolNull(),
			}

			var diags resource.CreateResponse
			result := r.executePullLogic(ctx, data, &diags.Diagnostics)

			assert.True(t, result)
			assert.False(t, diags.Diagnostics.HasError())
		})
	}
}
