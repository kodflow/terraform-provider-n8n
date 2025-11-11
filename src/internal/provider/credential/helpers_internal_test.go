package credential

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/credential/models"
	"github.com/stretchr/testify/assert"
)

// TestCredentialResource_createNewCredential tests the createNewCredential helper.
func TestCredentialResource_createNewCredential(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		credName     string
		credType     string
		setupHandler func(w http.ResponseWriter, r *http.Request)
		expectError  bool
		expectNil    bool
	}{
		{
			name:     "successful creation",
			credName: "New Credential",
			credType: "httpHeaderAuth",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPost && r.URL.Path == "/credentials" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusCreated)
					json.NewEncoder(w).Encode(map[string]any{
						"id":        "cred-new-123",
						"name":      "New Credential",
						"type":      "httpHeaderAuth",
						"createdAt": "2024-01-01T00:00:00Z",
						"updatedAt": "2024-01-01T00:00:00Z",
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: false,
			expectNil:   false,
		},
		{
			name:     "creation error",
			credName: "Failed Credential",
			credType: "httpBasicAuth",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"message": "Internal server error"}`))
			},
			expectError: true,
			expectNil:   true,
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
			diags := diag.Diagnostics{}

			result := r.createNewCredential(ctx, tt.credName, tt.credType, map[string]any{"key": "value"}, &diags)

			if tt.expectError {
				assert.True(t, diags.HasError(), "Should have diagnostics error")
			} else {
				assert.False(t, diags.HasError(), "Should not have diagnostics error")
			}

			if tt.expectNil {
				assert.Nil(t, result, "Result should be nil on error")
			} else {
				assert.NotNil(t, result, "Result should not be nil on success")
			}
		})
	}
}

// TestCredentialResource_scanAffectedWorkflows tests the scanAffectedWorkflows helper.
func TestCredentialResource_scanAffectedWorkflows(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		oldCredID         string
		newCredID         string
		setupHandler      func(w http.ResponseWriter, r *http.Request)
		expectError       bool
		expectSuccess     bool
		expectedWorkflows int
	}{
		{
			name:      "no affected workflows",
			oldCredID: "old-cred-1",
			newCredID: "new-cred-1",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet && r.URL.Path == "/workflows" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{"data": []any{}})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError:       false,
			expectSuccess:     true,
			expectedWorkflows: 0,
		},
		{
			name:      "workflows found with credential",
			oldCredID: "old-cred-2",
			newCredID: "new-cred-2",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet && r.URL.Path == "/workflows" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{
						"data": []any{
							map[string]any{
								"id":   "wf-1",
								"name": "Workflow 1",
								"nodes": []any{
									map[string]any{
										"credentials": map[string]any{
											"api": map[string]any{"id": "old-cred-2"},
										},
									},
								},
								"connections": map[string]any{},
								"settings":    map[string]any{},
							},
						},
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError:       false,
			expectSuccess:     true,
			expectedWorkflows: 1,
		},
		{
			name:      "error listing workflows",
			oldCredID: "old-cred-3",
			newCredID: "new-cred-3",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet && r.URL.Path == "/workflows" {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				if r.Method == http.MethodDelete {
					w.WriteHeader(http.StatusNoContent)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError:       true,
			expectSuccess:     false,
			expectedWorkflows: 0,
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
			diags := diag.Diagnostics{}

			workflows, success := r.scanAffectedWorkflows(ctx, tt.oldCredID, tt.newCredID, &diags)

			assert.Equal(t, tt.expectSuccess, success, "Success status should match expected")
			assert.Equal(t, tt.expectedWorkflows, len(workflows), "Number of workflows should match expected")

			if tt.expectError {
				assert.True(t, diags.HasError(), "Should have diagnostics error")
			}
		})
	}
}

// TestCredentialResource_updateAffectedWorkflows tests the updateAffectedWorkflows helper.
func TestCredentialResource_updateAffectedWorkflows(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		backups       []models.WorkflowBackup
		oldCredID     string
		newCredID     string
		setupHandler  func(w http.ResponseWriter, r *http.Request)
		expectSuccess bool
		expectedCount int
		expectError   bool
	}{
		{
			name:          "no workflows to update",
			backups:       []models.WorkflowBackup{},
			oldCredID:     "old-1",
			newCredID:     "new-1",
			setupHandler:  func(w http.ResponseWriter, r *http.Request) {},
			expectSuccess: true,
			expectedCount: 0,
			expectError:   false,
		},
		{
			name: "update single workflow successfully",
			backups: []models.WorkflowBackup{
				{
					ID: "wf-1",
					Original: &n8nsdk.Workflow{
						Id:          strPtr("wf-1"),
						Name:        "Test Workflow",
						Nodes:       []n8nsdk.Node{{Credentials: map[string]interface{}{"api": map[string]interface{}{"id": "old-1"}}}},
						Connections: map[string]interface{}{},
						Settings:    n8nsdk.WorkflowSettings{},
					},
				},
			},
			oldCredID: "old-1",
			newCredID: "new-1",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet && r.URL.Path == "/workflows/wf-1" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{
						"id":   "wf-1",
						"name": "Test Workflow",
						"nodes": []any{
							map[string]any{"credentials": map[string]any{"api": map[string]any{"id": "old-1"}}},
						},
						"connections": map[string]any{},
						"settings":    map[string]any{},
					})
					return
				}
				if r.Method == http.MethodPut && r.URL.Path == "/workflows/wf-1" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{
						"id":          "wf-1",
						"name":        "Test Workflow",
						"nodes":       []any{},
						"connections": map[string]any{},
						"settings":    map[string]any{},
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectSuccess: true,
			expectedCount: 1,
			expectError:   false,
		},
		{
			name: "error getting workflow",
			backups: []models.WorkflowBackup{
				{
					ID: "wf-error",
					Original: &n8nsdk.Workflow{
						Id:          strPtr("wf-error"),
						Name:        "Error Workflow",
						Nodes:       []n8nsdk.Node{},
						Connections: map[string]interface{}{},
						Settings:    n8nsdk.WorkflowSettings{},
					},
				},
			},
			oldCredID: "old-error",
			newCredID: "new-error",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet && r.URL.Path == "/workflows/wf-error" {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				if r.Method == http.MethodDelete {
					w.WriteHeader(http.StatusNoContent)
					return
				}
				if r.Method == http.MethodPut {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{
						"id":          "wf-1",
						"name":        "Test",
						"nodes":       []any{},
						"connections": map[string]any{},
						"settings":    map[string]any{},
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectSuccess: false,
			expectedCount: 0,
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
			diags := diag.Diagnostics{}

			updatedIDs, success := r.updateAffectedWorkflows(ctx, tt.backups, tt.oldCredID, tt.newCredID, &diags)

			assert.Equal(t, tt.expectSuccess, success, "Success status should match expected")
			assert.Equal(t, tt.expectedCount, len(updatedIDs), "Number of updated workflows should match expected")

			if tt.expectError {
				assert.True(t, diags.HasError(), "Should have diagnostics error")
			}
		})
	}
}

// TestCredentialResource_deleteCredentialBestEffort tests the deleteCredentialBestEffort helper.
func TestCredentialResource_deleteCredentialBestEffort(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		credID       string
		setupHandler func(w http.ResponseWriter, r *http.Request)
	}{
		{
			name:   "successful delete",
			credID: "cred-1",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodDelete && r.URL.Path == "/credentials/cred-1" {
					w.WriteHeader(http.StatusNoContent)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
		},
		{
			name:   "delete fails but doesn't panic",
			credID: "cred-2",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
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

			// Should not panic regardless of result
			assert.NotPanics(t, func() {
				r.deleteCredentialBestEffort(ctx, tt.credID)
			})
		})
	}
}

// TestCredentialResource_deleteOldCredential tests the deleteOldCredential helper.
func TestCredentialResource_deleteOldCredential(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		oldCredID    string
		newCredID    string
		setupHandler func(w http.ResponseWriter, r *http.Request)
		expectWarn   bool
	}{
		{
			name:      "successful delete",
			oldCredID: "old-cred-1",
			newCredID: "new-cred-1",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodDelete && r.URL.Path == "/credentials/old-cred-1" {
					w.WriteHeader(http.StatusNoContent)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectWarn: false,
		},
		{
			name:      "delete fails - logs warning",
			oldCredID: "old-cred-2",
			newCredID: "new-cred-2",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"message": "Internal server error"}`))
			},
			expectWarn: true,
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

			// Should not panic regardless of result
			assert.NotPanics(t, func() {
				r.deleteOldCredential(ctx, tt.oldCredID, tt.newCredID)
			})
		})
	}
}

// TestCredentialResource_updateAffectedWorkflows_ErrorCases tests error cases.
func TestCredentialResource_updateAffectedWorkflows_ErrorCases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		backups       []models.WorkflowBackup
		oldCredID     string
		newCredID     string
		setupHandler  func(w http.ResponseWriter, r *http.Request)
		expectSuccess bool
		expectedCount int
		expectError   bool
	}{
		{
			name: "error updating workflow",
			backups: []models.WorkflowBackup{
				{
					ID: "wf-update-error",
					Original: &n8nsdk.Workflow{
						Id:          strPtr("wf-update-error"),
						Name:        "Error Workflow",
						Nodes:       []n8nsdk.Node{{Credentials: map[string]interface{}{"api": map[string]interface{}{"id": "old-update"}}}},
						Connections: map[string]interface{}{},
						Settings:    n8nsdk.WorkflowSettings{},
					},
				},
			},
			oldCredID: "old-update",
			newCredID: "new-update",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet && r.URL.Path == "/workflows/wf-update-error" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{
						"id":   "wf-update-error",
						"name": "Error Workflow",
						"nodes": []any{
							map[string]any{"credentials": map[string]any{"api": map[string]any{"id": "old-update"}}},
						},
						"connections": map[string]any{},
						"settings":    map[string]any{},
					})
					return
				}
				if r.Method == http.MethodPut && r.URL.Path == "/workflows/wf-update-error" {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				if r.Method == http.MethodDelete {
					w.WriteHeader(http.StatusNoContent)
					return
				}
				if r.Method == http.MethodPut {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{
						"id":          "wf-1",
						"name":        "Test",
						"nodes":       []any{},
						"connections": map[string]any{},
						"settings":    map[string]any{},
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectSuccess: false,
			expectedCount: 0,
			expectError:   true,
		},
		{
			name: "update multiple workflows successfully with throttling",
			backups: []models.WorkflowBackup{
				{
					ID: "wf-1",
					Original: &n8nsdk.Workflow{
						Id:          strPtr("wf-1"),
						Name:        "Workflow 1",
						Nodes:       []n8nsdk.Node{{Credentials: map[string]interface{}{"api": map[string]interface{}{"id": "old-multi"}}}},
						Connections: map[string]interface{}{},
						Settings:    n8nsdk.WorkflowSettings{},
					},
				},
				{
					ID: "wf-2",
					Original: &n8nsdk.Workflow{
						Id:          strPtr("wf-2"),
						Name:        "Workflow 2",
						Nodes:       []n8nsdk.Node{{Credentials: map[string]interface{}{"api": map[string]interface{}{"id": "old-multi"}}}},
						Connections: map[string]interface{}{},
						Settings:    n8nsdk.WorkflowSettings{},
					},
				},
			},
			oldCredID: "old-multi",
			newCredID: "new-multi",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					if r.URL.Path == "/workflows/wf-1" {
						json.NewEncoder(w).Encode(map[string]any{
							"id":   "wf-1",
							"name": "Workflow 1",
							"nodes": []any{
								map[string]any{"credentials": map[string]any{"api": map[string]any{"id": "old-multi"}}},
							},
							"connections": map[string]any{},
							"settings":    map[string]any{},
						})
						return
					}
					if r.URL.Path == "/workflows/wf-2" {
						json.NewEncoder(w).Encode(map[string]any{
							"id":   "wf-2",
							"name": "Workflow 2",
							"nodes": []any{
								map[string]any{"credentials": map[string]any{"api": map[string]any{"id": "old-multi"}}},
							},
							"connections": map[string]any{},
							"settings":    map[string]any{},
						})
						return
					}
				}
				if r.Method == http.MethodPut {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{
						"id":          r.URL.Path[len("/workflows/"):],
						"name":        "Updated",
						"nodes":       []any{},
						"connections": map[string]any{},
						"settings":    map[string]any{},
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectSuccess: true,
			expectedCount: 2,
			expectError:   false,
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
			diags := diag.Diagnostics{}

			updatedIDs, success := r.updateAffectedWorkflows(ctx, tt.backups, tt.oldCredID, tt.newCredID, &diags)

			assert.Equal(t, tt.expectSuccess, success, "Success status should match expected")
			assert.Equal(t, tt.expectedCount, len(updatedIDs), "Number of updated workflows should match expected")

			if tt.expectError {
				assert.True(t, diags.HasError(), "Should have diagnostics error")
			}
		})
	}
}
