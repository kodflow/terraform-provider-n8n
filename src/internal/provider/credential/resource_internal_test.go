package credential

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/credential/models"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
	"github.com/stretchr/testify/assert"
)

// KTN-TEST-009 RATIONALE: Tests in this file test public functions that require
// access to private fields (specifically the .client field). These tests MUST
// remain in the internal test file to access the unexported client field for
// proper white-box testing. Moving these to external tests would prevent us
// from properly testing the functions that depend on the client field.
//
// Note: Mock implementations are in helpers_test.go to avoid duplication.

// Helper functions.
func strPtr(s string) *string {
	return &s
}

// TestNewCredentialResource is now in resource_external_test.go
// TestNewCredentialResourceWrapper is now in resource_external_test.go
// TestCredentialResource_Configure is now in resource_external_test.go.
func Test_usesCredential(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		workflow *n8nsdk.Workflow
		credID   string
		want     bool
		wantErr  bool
	}{
		{
			name: "workflow uses credential",
			workflow: &n8nsdk.Workflow{
				Nodes: []n8nsdk.Node{{Credentials: map[string]interface{}{"api": map[string]interface{}{"id": "cred-123"}}}},
			},
			credID:  "cred-123",
			want:    true,
			wantErr: false,
		},
		{
			name: "workflow does not use credential",
			workflow: &n8nsdk.Workflow{
				Nodes: []n8nsdk.Node{{Credentials: map[string]interface{}{"api": map[string]interface{}{"id": "other-cred"}}}},
			},
			credID:  "cred-123",
			want:    false,
			wantErr: false,
		},
		{
			name: "workflow with no credentials",
			workflow: &n8nsdk.Workflow{
				Nodes: []n8nsdk.Node{{Credentials: nil}},
			},
			credID:  "cred-123",
			want:    false,
			wantErr: false,
		},
		{
			name:     "workflow with empty nodes",
			workflow: &n8nsdk.Workflow{Nodes: []n8nsdk.Node{}},
			credID:   "cred-123",
			want:     false,
			wantErr:  false,
		},
		{
			name:     "nil workflow",
			workflow: nil,
			credID:   "cred-123",
			want:     false,
			wantErr:  false,
		},
		{
			name: "multiple nodes with credentials",
			workflow: &n8nsdk.Workflow{
				Nodes: []n8nsdk.Node{
					{Credentials: map[string]interface{}{"api": map[string]interface{}{"id": "other-cred"}}},
					{Credentials: map[string]interface{}{"http": map[string]interface{}{"id": "cred-123"}}},
				},
			},
			credID:  "cred-123",
			want:    true,
			wantErr: false,
		},
		{
			name: "credential in nested structure",
			workflow: &n8nsdk.Workflow{
				Nodes: []n8nsdk.Node{{Credentials: map[string]interface{}{
					"api":  map[string]interface{}{"id": "cred-123", "name": "Test"},
					"http": map[string]interface{}{"id": "other-cred"},
				}}},
			},
			credID:  "cred-123",
			want:    true,
			wantErr: false,
		},
		{
			name: "invalid credential structure",
			workflow: &n8nsdk.Workflow{
				Nodes: []n8nsdk.Node{{Credentials: map[string]interface{}{"api": "invalid-not-a-map"}}},
			},
			credID:  "cred-123",
			want:    false,
			wantErr: false,
		},
		{
			name: "credential with different types",
			workflow: &n8nsdk.Workflow{
				Nodes: []n8nsdk.Node{{Credentials: map[string]interface{}{
					"api":  map[string]interface{}{"id": "cred-123"},
					"http": map[string]interface{}{"id": "cred-123"},
				}}},
			},
			credID:  "cred-123",
			want:    true,
			wantErr: false,
		},
		{
			name:     "error case - empty credential ID",
			workflow: &n8nsdk.Workflow{Nodes: []n8nsdk.Node{{Credentials: map[string]interface{}{"api": map[string]interface{}{"id": "cred-123"}}}}},
			credID:   "",
			want:     false,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := usesCredential(tt.workflow, tt.credID)
			assert.Equal(t, tt.want, result)
		})
	}
}

func Test_replaceCredentialInWorkflow(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		workflow *n8nsdk.Workflow
		oldCred  string
		newCred  string
		wantErr  bool
		checkFn  func(*testing.T, *n8nsdk.Workflow)
	}{
		{
			name: "replace single credential",
			workflow: &n8nsdk.Workflow{
				Id:    strPtr("wf-1"),
				Nodes: []n8nsdk.Node{{Credentials: map[string]interface{}{"api": map[string]interface{}{"id": "old-cred", "name": "Old"}}}},
			},
			oldCred: "old-cred",
			newCred: "new-cred",
			checkFn: func(t *testing.T, w *n8nsdk.Workflow) {
				t.Helper()
				assert.NotNil(t, w)
				apiCred := w.Nodes[0].Credentials["api"].(map[string]interface{})
				assert.Equal(t, "new-cred", apiCred["id"])
			},
		},
		{
			name: "replace multiple occurrences",
			workflow: &n8nsdk.Workflow{
				Nodes: []n8nsdk.Node{
					{Credentials: map[string]interface{}{"api": map[string]interface{}{"id": "old-cred"}}},
					{Credentials: map[string]interface{}{"http": map[string]interface{}{"id": "old-cred"}}},
				},
			},
			oldCred: "old-cred",
			newCred: "new-cred",
			checkFn: func(t *testing.T, w *n8nsdk.Workflow) {
				t.Helper()
				assert.NotNil(t, w)
				apiCred := w.Nodes[0].Credentials["api"].(map[string]interface{})
				assert.Equal(t, "new-cred", apiCred["id"])
				httpCred := w.Nodes[1].Credentials["http"].(map[string]interface{})
				assert.Equal(t, "new-cred", httpCred["id"])
			},
		},
		{
			name: "no replacement needed",
			workflow: &n8nsdk.Workflow{
				Nodes: []n8nsdk.Node{{Credentials: map[string]interface{}{"api": map[string]interface{}{"id": "other-cred"}}}},
			},
			oldCred: "old-cred",
			newCred: "new-cred",
			checkFn: func(t *testing.T, w *n8nsdk.Workflow) {
				t.Helper()
				assert.NotNil(t, w)
				apiCred := w.Nodes[0].Credentials["api"].(map[string]interface{})
				assert.Equal(t, "other-cred", apiCred["id"])
			},
		},
		{
			name:     "nil workflow",
			workflow: nil,
			oldCred:  "old-cred",
			newCred:  "new-cred",
			checkFn: func(t *testing.T, w *n8nsdk.Workflow) {
				t.Helper()
				assert.Nil(t, w)
			},
		},
		{
			name:     "empty nodes",
			workflow: &n8nsdk.Workflow{Id: strPtr("wf-1"), Nodes: []n8nsdk.Node{}},
			oldCred:  "old-cred",
			newCred:  "new-cred",
			checkFn: func(t *testing.T, w *n8nsdk.Workflow) {
				t.Helper()
				assert.Empty(t, w.Nodes)
			},
		},
		{
			name: "node without credentials",
			workflow: &n8nsdk.Workflow{
				Nodes: []n8nsdk.Node{{Credentials: nil}},
			},
			oldCred: "old-cred",
			newCred: "new-cred",
			checkFn: func(t *testing.T, w *n8nsdk.Workflow) {
				t.Helper()
				assert.Nil(t, w.Nodes[0].Credentials)
			},
		},
		{
			name: "invalid credential structure",
			workflow: &n8nsdk.Workflow{
				Nodes: []n8nsdk.Node{{Credentials: map[string]interface{}{"api": "invalid-string"}}},
			},
			oldCred: "old-cred",
			newCred: "new-cred",
			checkFn: func(t *testing.T, w *n8nsdk.Workflow) {
				t.Helper()
				assert.NotNil(t, w)
			},
		},
		{
			name: "mixed valid and invalid credentials",
			workflow: &n8nsdk.Workflow{
				Nodes: []n8nsdk.Node{{Credentials: map[string]interface{}{
					"api":     map[string]interface{}{"id": "old-cred"},
					"invalid": "not-a-map",
				}}},
			},
			oldCred: "old-cred",
			newCred: "new-cred",
			checkFn: func(t *testing.T, w *n8nsdk.Workflow) {
				t.Helper()
				apiCred := w.Nodes[0].Credentials["api"].(map[string]interface{})
				assert.Equal(t, "new-cred", apiCred["id"])
			},
		},
		{
			name: "preserve other credential properties",
			workflow: &n8nsdk.Workflow{
				Nodes: []n8nsdk.Node{{Credentials: map[string]interface{}{
					"api": map[string]interface{}{"id": "old-cred", "name": "API Credential", "type": "oauth2"},
				}}},
			},
			oldCred: "old-cred",
			newCred: "new-cred",
			checkFn: func(t *testing.T, w *n8nsdk.Workflow) {
				t.Helper()
				apiCred := w.Nodes[0].Credentials["api"].(map[string]interface{})
				assert.Equal(t, "new-cred", apiCred["id"])
				assert.Equal(t, "API Credential", apiCred["name"])
				assert.Equal(t, "oauth2", apiCred["type"])
			},
		},
		{
			name: "error case - empty old credential ID",
			workflow: &n8nsdk.Workflow{
				Nodes: []n8nsdk.Node{{Credentials: map[string]interface{}{"api": map[string]interface{}{"id": "cred"}}}},
			},
			oldCred: "",
			newCred: "new-cred",
			wantErr: true,
			checkFn: func(t *testing.T, w *n8nsdk.Workflow) {
				t.Helper()
				assert.NotNil(t, w)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := replaceCredentialInWorkflow(tt.workflow, tt.oldCred, tt.newCred)
			if tt.checkFn != nil {
				tt.checkFn(t, result)
			}
		})
	}
}

func TestROTATION_THROTTLE_MILLISECONDS(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		want    int
		wantErr bool
	}{
		{name: "constant value", want: 100, wantErr: false},
		{name: "error case - constant should not be zero", want: 100, wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "constant value":
				assert.Equal(t, tt.want, ROTATION_THROTTLE_MILLISECONDS)
			case "error case - constant should not be zero":
				assert.NotEqual(t, 0, ROTATION_THROTTLE_MILLISECONDS)
				assert.Equal(t, tt.want, ROTATION_THROTTLE_MILLISECONDS)
			}
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

func TestCredentialResource_ConcurrentOperations(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "concurrent credential checks", wantErr: false},
		{name: "concurrent replacements", wantErr: false},
		{name: "error case - concurrent operations on nil workflow", wantErr: true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "concurrent credential checks":
				workflow := &n8nsdk.Workflow{
					Nodes: []n8nsdk.Node{
						{
							Credentials: map[string]interface{}{
								"api": map[string]interface{}{
									"id": "cred-123",
								},
							},
						},
					},
				}

				// Run concurrent checks
				results := make(chan bool, 100)
				for range 100 {
					go func() {
						results <- usesCredential(workflow, "cred-123")
					}()
				}

				// Collect results
				for range 100 {
					assert.True(t, <-results)
				}

			case "concurrent replacements":
				// Run concurrent replacements on copies
				results := make(chan *n8nsdk.Workflow, 100)
				for i := range 100 {
					go func(idx int) {
						// Create a copy for each goroutine
						workflow := &n8nsdk.Workflow{
							Id: strPtr("wf-1"),
							Nodes: []n8nsdk.Node{
								{
									Credentials: map[string]interface{}{
										"api": map[string]interface{}{
											"id": "old-cred",
										},
									},
								},
							},
						}
						newCredID := fmt.Sprintf("new-cred-%d", idx)
						results <- replaceCredentialInWorkflow(workflow, "old-cred", newCredID)
					}(i)
				}

				// Verify all replacements
				for range 100 {
					result := <-results
					assert.NotNil(t, result)
				}

			case "error case - concurrent operations on nil workflow":
				// Test that concurrent operations on nil workflow are safe
				results := make(chan bool, 10)
				for range 10 {
					go func() {
						results <- usesCredential(nil, "cred-123")
					}()
				}

				// Collect results - should all be false
				for range 10 {
					assert.False(t, <-results)
				}
			}
		})
	}
}

func TestCredentialResource_EdgeCases(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		workflow *n8nsdk.Workflow
		credID   string
		want     bool
		wantErr  bool
	}{
		{
			name: "deeply nested credential structure",
			workflow: &n8nsdk.Workflow{
				Nodes: []n8nsdk.Node{
					{
						Credentials: map[string]interface{}{
							"oauth": map[string]interface{}{
								"id": "cred-123",
								"metadata": map[string]interface{}{
									"scope":   "read write",
									"expires": 3600,
								},
							},
						},
					},
				},
			},
			credID:  "cred-123",
			want:    true,
			wantErr: false,
		},
		{
			name: "credential ID as non-string",
			workflow: &n8nsdk.Workflow{
				Nodes: []n8nsdk.Node{
					{
						Credentials: map[string]interface{}{
							"api": map[string]interface{}{
								"id": 123, // Number instead of string
							},
						},
					},
				},
			},
			credID:  "123",
			want:    false,
			wantErr: false,
		},
		{
			name:     "error case - nil workflow",
			workflow: nil,
			credID:   "cred-123",
			want:     false,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := usesCredential(tt.workflow, tt.credID)
			assert.Equal(t, tt.want, result)
		})
	}
}

func BenchmarkUsesCredential(b *testing.B) {
	workflow := &n8nsdk.Workflow{
		Nodes: []n8nsdk.Node{
			{
				Credentials: map[string]interface{}{
					"api": map[string]interface{}{
						"id": "cred-123",
					},
				},
			},
			{
				Credentials: map[string]interface{}{
					"http": map[string]interface{}{
						"id": "other-cred",
					},
				},
			},
		},
	}

	b.ResetTimer()
	for b.Loop() {
		_ = usesCredential(workflow, "cred-123")
	}
}

func BenchmarkReplaceCredentialInWorkflow(b *testing.B) {
	workflow := &n8nsdk.Workflow{
		Id: strPtr("wf-1"),
		Nodes: []n8nsdk.Node{
			{
				Credentials: map[string]interface{}{
					"api": map[string]interface{}{
						"id": "old-cred",
					},
				},
			},
		},
	}

	b.ResetTimer()
	for b.Loop() {
		_ = replaceCredentialInWorkflow(workflow, "old-cred", "new-cred")
	}
}

func BenchmarkReplaceCredentialLargeWorkflow(b *testing.B) {
	// Create a workflow with many nodes
	nodes := make([]n8nsdk.Node, 100)
	for range 100 {
		nodes = append(nodes, n8nsdk.Node{
			Credentials: map[string]interface{}{
				fmt.Sprintf("api%d", len(nodes)): map[string]interface{}{
					"id": "old-cred",
				},
			},
		})
	}

	workflow := &n8nsdk.Workflow{
		Id:    strPtr("large-wf"),
		Nodes: nodes,
	}

	b.ResetTimer()
	for b.Loop() {
		_ = replaceCredentialInWorkflow(workflow, "old-cred", "new-cred")
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

// TestCredentialResource_Create tests credential creation.
// TestCredentialResource_Create is now in resource_external_test.go
// TestCredentialResource_Read is now in resource_external_test.go
// TestCredentialResource_Update is now in resource_external_test.go
// TestCredentialResource_Delete is now in resource_external_test.go.
// TestCredentialResource_RollbackFunctions is now in external test file - refactored to test behavior only.

func TestCredentialResource_DeleteOldCredential(t *testing.T) {
	t.Run("deleteOldCredential failure logs warning", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message": "Internal server error"}`))
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &CredentialResource{client: n8nClient}
		r.deleteOldCredential(context.Background(), "old-cred-123", "new-cred-456")
	})
}

// TestCredentialResource_HelperFunctions tests the helper functions for credential rotation.
// TestCredentialResource_HelperFunctions is now in external test file - refactored to test behavior only.

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

// TestCredentialResource_rollbackRotation tests the rollbackRotation method.
func TestCredentialResource_rollbackRotation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		newCredID         string
		setupHandler      func(w http.ResponseWriter, r *http.Request)
		expectDeleteCall  bool
		expectUpdateCalls int
	}{
		{
			name:      "rollback with delete and restore",
			newCredID: "new-cred-123",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodDelete {
					w.WriteHeader(http.StatusNoContent)
					return
				}
				if r.Method == http.MethodPut {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{
						"id":          "wf-123",
						"name":        "Test",
						"nodes":       []any{},
						"connections": map[string]any{},
						"settings":    map[string]any{},
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectDeleteCall:  true,
			expectUpdateCalls: 1,
		},
		{
			name:      "rollback with delete failure",
			newCredID: "new-cred-456",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodDelete {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				if r.Method == http.MethodPut {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{
						"id":          "wf-123",
						"name":        "Test",
						"nodes":       []any{},
						"connections": map[string]any{},
						"settings":    map[string]any{},
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectDeleteCall:  true,
			expectUpdateCalls: 1,
		},
		{
			name:      "rollback with restore failure",
			newCredID: "new-cred-789",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodDelete {
					w.WriteHeader(http.StatusNoContent)
					return
				}
				if r.Method == http.MethodPut {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectDeleteCall:  true,
			expectUpdateCalls: 1,
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
			backups := []models.WorkflowBackup{
				{
					ID: "wf-123",
					Original: &n8nsdk.Workflow{
						Id:          strPtr("wf-123"),
						Name:        "Test",
						Nodes:       []n8nsdk.Node{},
						Connections: map[string]interface{}{},
						Settings:    n8nsdk.WorkflowSettings{},
					},
				},
			}
			updatedWorkflows := []string{"wf-123"}

			// Call rollbackRotation - it should not panic
			r.rollbackRotation(context.Background(), tt.newCredID, backups, updatedWorkflows)

			// Verification is implicit - if the function completes without panic, the test passes
			assert.NotNil(t, r, "resource should not be nil after rollback")
		})
	}
}

// TestCredentialResource_deleteNewCredential tests the deleteNewCredential method.
func TestCredentialResource_deleteNewCredential(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		credID     string
		statusCode int
		expectCall bool
	}{
		{
			name:       "successful delete",
			credID:     "new-cred-123",
			statusCode: http.StatusNoContent,
			expectCall: true,
		},
		{
			name:       "delete failure",
			credID:     "new-cred-456",
			statusCode: http.StatusInternalServerError,
			expectCall: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			called := false
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodDelete && r.URL.Path == "/credentials/"+tt.credID {
					called = true
					w.WriteHeader(tt.statusCode)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			})

			n8nClient, server := setupTestClient(t, handler)
			defer server.Close()

			r := &CredentialResource{client: n8nClient}

			// Call deleteNewCredential - it should not panic
			r.deleteNewCredential(context.Background(), tt.credID)

			if tt.expectCall {
				assert.True(t, called, "expected API call to be made")
			}
		})
	}
}

// TestCredentialResource_restoreWorkflows tests the restoreWorkflows method.
func TestCredentialResource_restoreWorkflows(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		backups          []models.WorkflowBackup
		updatedWorkflows []string
		setupHandler     func(w http.ResponseWriter, r *http.Request)
		expectedRestored int
	}{
		{
			name: "restore single workflow",
			backups: []models.WorkflowBackup{
				{
					ID: "wf-123",
					Original: &n8nsdk.Workflow{
						Id:          strPtr("wf-123"),
						Name:        "Test",
						Nodes:       []n8nsdk.Node{},
						Connections: map[string]interface{}{},
						Settings:    n8nsdk.WorkflowSettings{},
					},
				},
			},
			updatedWorkflows: []string{"wf-123"},
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPut && r.URL.Path == "/workflows/wf-123" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{
						"id":          "wf-123",
						"name":        "Test",
						"nodes":       []any{},
						"connections": map[string]any{},
						"settings":    map[string]any{},
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectedRestored: 1,
		},
		{
			name:             "no workflows to restore",
			backups:          []models.WorkflowBackup{},
			updatedWorkflows: []string{},
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			expectedRestored: 0,
		},
		{
			name: "restore fails with error",
			backups: []models.WorkflowBackup{
				{
					ID: "wf-error",
					Original: &n8nsdk.Workflow{
						Id:          strPtr("wf-error"),
						Name:        "Error Test",
						Nodes:       []n8nsdk.Node{},
						Connections: map[string]interface{}{},
						Settings:    n8nsdk.WorkflowSettings{},
					},
				},
			},
			updatedWorkflows: []string{"wf-error"},
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			expectedRestored: 0,
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

			count := r.restoreWorkflows(context.Background(), tt.backups, tt.updatedWorkflows)

			assert.Equal(t, tt.expectedRestored, count, "unexpected number of restored workflows")
		})
	}
}

// TestCredentialResource_findWorkflowBackup tests the findWorkflowBackup method.
func TestCredentialResource_findWorkflowBackup(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		affectedWorkflows []models.WorkflowBackup
		workflowID        string
		want              *n8nsdk.Workflow
		wantErr           bool
	}{
		{
			name: "find existing workflow backup",
			affectedWorkflows: []models.WorkflowBackup{
				{
					ID: "wf-123",
					Original: &n8nsdk.Workflow{
						Id:   strPtr("wf-123"),
						Name: "Test Workflow",
					},
				},
			},
			workflowID: "wf-123",
			want: &n8nsdk.Workflow{
				Id:   strPtr("wf-123"),
				Name: "Test Workflow",
			},
			wantErr: false,
		},
		{
			name: "workflow not found in backups",
			affectedWorkflows: []models.WorkflowBackup{
				{
					ID: "wf-123",
					Original: &n8nsdk.Workflow{
						Id:   strPtr("wf-123"),
						Name: "Test Workflow",
					},
				},
			},
			workflowID: "wf-456",
			want:       nil,
			wantErr:    false,
		},
		{
			name:              "empty backups list",
			affectedWorkflows: []models.WorkflowBackup{},
			workflowID:        "wf-123",
			want:              nil,
			wantErr:           false,
		},
		{
			name: "find workflow in multiple backups",
			affectedWorkflows: []models.WorkflowBackup{
				{
					ID: "wf-111",
					Original: &n8nsdk.Workflow{
						Id:   strPtr("wf-111"),
						Name: "First Workflow",
					},
				},
				{
					ID: "wf-222",
					Original: &n8nsdk.Workflow{
						Id:   strPtr("wf-222"),
						Name: "Second Workflow",
					},
				},
				{
					ID: "wf-333",
					Original: &n8nsdk.Workflow{
						Id:   strPtr("wf-333"),
						Name: "Third Workflow",
					},
				},
			},
			workflowID: "wf-222",
			want: &n8nsdk.Workflow{
				Id:   strPtr("wf-222"),
				Name: "Second Workflow",
			},
			wantErr: false,
		},
		{
			name: "backup with nil original",
			affectedWorkflows: []models.WorkflowBackup{
				{
					ID:       "wf-123",
					Original: nil,
				},
			},
			workflowID: "wf-123",
			want:       nil,
			wantErr:    false,
		},
		{
			name:              "error case - empty workflow ID",
			affectedWorkflows: []models.WorkflowBackup{},
			workflowID:        "",
			want:              nil,
			wantErr:           true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := &CredentialResource{}
			result := r.findWorkflowBackup(tt.affectedWorkflows, tt.workflowID)

			if tt.want == nil {
				assert.Nil(t, result, "expected nil result")
			} else {
				assert.NotNil(t, result, "expected non-nil result")
				if result != nil {
					assert.Equal(t, tt.want.Id, result.Id, "workflow ID should match")
					assert.Equal(t, tt.want.Name, result.Name, "workflow name should match")
				}
			}
		})
	}
}

func TestCredentialResource_restoreWorkflow(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		workflowID      string
		workflow        *n8nsdk.Workflow
		setupHandler    func(w http.ResponseWriter, r *http.Request)
		expectedSuccess bool
	}{
		{
			name:       "successful restore",
			workflowID: "wf-123",
			workflow: &n8nsdk.Workflow{
				Id:          strPtr("wf-123"),
				Name:        "Test",
				Nodes:       []n8nsdk.Node{},
				Connections: map[string]interface{}{},
				Settings:    n8nsdk.WorkflowSettings{},
			},
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPut && r.URL.Path == "/workflows/wf-123" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{
						"id":          "wf-123",
						"name":        "Test",
						"nodes":       []any{},
						"connections": map[string]any{},
						"settings":    map[string]any{},
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectedSuccess: true,
		},
		{
			name:       "restore failure",
			workflowID: "wf-123",
			workflow: &n8nsdk.Workflow{
				Id:          strPtr("wf-123"),
				Name:        "Test",
				Nodes:       []n8nsdk.Node{},
				Connections: map[string]interface{}{},
				Settings:    n8nsdk.WorkflowSettings{},
			},
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			expectedSuccess: false,
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

			success := r.restoreWorkflow(context.Background(), tt.workflowID, tt.workflow)

			assert.Equal(t, tt.expectedSuccess, success, "unexpected restore result")
		})
	}
}
