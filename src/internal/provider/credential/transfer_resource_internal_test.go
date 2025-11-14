package credential

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/credential/models"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
	"github.com/stretchr/testify/assert"
)

func TestCredentialTransferResource_executeCreateLogic(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		credentialID string
		projectID    string
		setupHandler func(w http.ResponseWriter, r *http.Request)
		expectError  bool
		expectID     string
	}{
		{
			name:         "successful creation",
			credentialID: "cred-123",
			projectID:    "proj-456",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPut && r.URL.Path == "/credentials/cred-123/transfer" {
					w.WriteHeader(http.StatusOK)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: false,
			expectID:    "cred-123-to-proj-456",
		},
		{
			name:         "API error",
			credentialID: "cred-123",
			projectID:    "proj-456",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"message": "Internal server error"}`))
			},
			expectError: true,
		},
		{
			name:         "credential not found",
			credentialID: "cred-999",
			projectID:    "proj-456",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`{"message": "Credential not found"}`))
			},
			expectError: true,
		},
		{
			name:         "project not found",
			credentialID: "cred-123",
			projectID:    "proj-999",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`{"message": "Project not found"}`))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := http.HandlerFunc(tt.setupHandler)
			n8nClient, server := setupTransferTestClient(t, handler)
			defer server.Close()

			r := &CredentialTransferResource{client: n8nClient}
			ctx := context.Background()
			plan := &models.TransferResource{
				CredentialID:         types.StringValue(tt.credentialID),
				DestinationProjectID: types.StringValue(tt.projectID),
			}
			resp := &resource.CreateResponse{}

			result := r.executeCreateLogic(ctx, plan, resp)

			if tt.expectError {
				assert.False(t, result, "Should return false on error")
				assert.True(t, resp.Diagnostics.HasError(), "Should have diagnostics error")
			} else {
				assert.True(t, result, "Should return true on success")
				assert.False(t, resp.Diagnostics.HasError(), "Should not have diagnostics error")
				assert.Equal(t, tt.expectID, plan.ID.ValueString(), "Composite ID should match")
				assert.Contains(t, plan.TransferredAt.ValueString(), "transfer-response-", "TransferredAt should be set")
			}
		})
	}
}

func setupTransferTestClient(t *testing.T, handler http.HandlerFunc) (*client.N8nClient, *httptest.Server) {
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
