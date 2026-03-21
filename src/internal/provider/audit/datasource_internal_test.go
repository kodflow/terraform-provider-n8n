package audit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/audit/models"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/client"
	"github.com/stretchr/testify/assert"
)

// setupTestClient creates a test N8nClient with httptest server for audit tests.
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

// TestAuditDataSource_executeAuditLogic tests the executeAuditLogic method.
func TestAuditDataSource_executeAuditLogic(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                     string
		setupHandler             func(w http.ResponseWriter, r *http.Request)
		expectError              bool
		expectCredentialsNotNull bool
	}{
		{
			name: "successful audit",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPost && r.URL.Path == "/audit" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{
						"Credentials Risk Report": map[string]any{"risk": "low"},
						"Database Risk Report":    map[string]any{"tables": 3},
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError:              false,
			expectCredentialsNotNull: true,
		},
		{
			name: "API error",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"message": "Internal server error"}`))
			},
			expectError: true,
		},
		{
			name: "empty audit response",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPost && r.URL.Path == "/audit" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError:              false,
			expectCredentialsNotNull: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			n8nClient, server := setupTestClient(t, http.HandlerFunc(tt.setupHandler))
			defer server.Close()

			d := &AuditDataSource{client: n8nClient}
			ctx := t.Context()
			data := &models.DataSource{}
			resp := &datasource.ReadResponse{
				State: datasource.ReadResponse{}.State,
			}

			result := d.executeAuditLogic(ctx, data, resp)

			//: Verify the result matches the expected error state.
			if tt.expectError {
				assert.False(t, result, "Should return false on error")
				assert.True(t, resp.Diagnostics.HasError(), "Should have diagnostics error")
			} else {
				assert.True(t, result, "Should return true on success")
				assert.False(t, resp.Diagnostics.HasError(), "Should not have diagnostics error")
				if tt.expectCredentialsNotNull {
					assert.False(t, data.CredentialsRiskReport.IsNull(), "CredentialsRiskReport should not be null")
				} else {
					assert.True(t, data.CredentialsRiskReport.IsNull(), "CredentialsRiskReport should be null")
				}
			}
		})
	}
}
