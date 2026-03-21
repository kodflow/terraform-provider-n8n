package credential

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/credential/models"
	"github.com/stretchr/testify/assert"
)

// TestCredentialsDataSource_executeListLogic tests the executeListLogic method.
func TestCredentialsDataSource_executeListLogic(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		setupHandler func(w http.ResponseWriter, r *http.Request)
		expectError  bool
		expectCount  int
	}{
		{
			name: "successful list",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet && r.URL.Path == "/credentials" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{
						"data": []any{
							map[string]any{
								"id":        "cred-1",
								"name":      "Cred One",
								"type":      "httpHeaderAuth",
								"createdAt": "2024-01-01T00:00:00Z",
								"updatedAt": "2024-01-01T00:00:00Z",
							},
							map[string]any{
								"id":        "cred-2",
								"name":      "Cred Two",
								"type":      "httpBasicAuth",
								"createdAt": "2024-01-02T00:00:00Z",
								"updatedAt": "2024-01-02T00:00:00Z",
							},
						},
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: false,
			expectCount: 2,
		},
		{
			name: "empty list",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet && r.URL.Path == "/credentials" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{
						"data": []any{},
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: false,
			expectCount: 0,
		},
		{
			name: "API error",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"message": "Internal server error"}`))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			n8nClient, server := setupTestClient(t, http.HandlerFunc(tt.setupHandler))
			defer server.Close()

			d := &CredentialsDataSource{client: n8nClient}
			ctx := t.Context()
			var data models.DataSources
			resp := &datasource.ReadResponse{
				State: datasource.ReadResponse{}.State,
			}

			result := d.executeListLogic(ctx, &data, resp)

			if tt.expectError {
				assert.False(t, result, "Should return false on error")
				assert.True(t, resp.Diagnostics.HasError(), "Should have diagnostics error")
			} else {
				assert.True(t, result, "Should return true on success")
				assert.False(t, resp.Diagnostics.HasError(), "Should not have diagnostics error")
				assert.Len(t, data.Credentials, tt.expectCount, "Credentials count should match")
			}
		})
	}
}
