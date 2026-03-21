package execution

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/execution/models"
	"github.com/stretchr/testify/assert"
)

// TestExecutionsDataSource_executeListLogic tests the executeListLogic method.
func TestExecutionsDataSource_executeListLogic(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		workflowID   string
		status       string
		setupHandler func(w http.ResponseWriter, r *http.Request)
		expectError  bool
		expectCount  int
	}{
		{
			name: "successful list - no filters",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet && r.URL.Path == "/executions" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{
						"data": []any{
							map[string]any{"id": 1, "mode": "manual", "status": "success"},
							map[string]any{"id": 2, "mode": "trigger", "status": "error"},
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
			name:       "successful list - with workflow filter",
			workflowID: "wf-1",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet && r.URL.Path == "/executions" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{
						"data": []any{
							map[string]any{"id": 1, "mode": "manual", "status": "success"},
						},
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: false,
			expectCount: 1,
		},
		{
			name:   "API error",
			status: "success",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"message": "Internal server error"}`))
			},
			expectError: true,
		},
		{
			name: "empty list",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet && r.URL.Path == "/executions" {
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			n8nClient, server := setupTestClient(t, http.HandlerFunc(tt.setupHandler))
			defer server.Close()

			d := &ExecutionsDataSource{client: n8nClient}
			ctx := t.Context()

			var workflowID types.String
			var status types.String

			//: Set filters when provided by the test case.
			if tt.workflowID != "" {
				workflowID = types.StringValue(tt.workflowID)
			} else {
				workflowID = types.StringNull()
			}
			if tt.status != "" {
				status = types.StringValue(tt.status)
			} else {
				status = types.StringNull()
			}

			data := &models.DataSources{
				WorkflowID: workflowID,
				Status:     status,
			}
			resp := &datasource.ReadResponse{
				State: datasource.ReadResponse{}.State,
			}

			result := d.executeListLogic(ctx, data, resp)

			//: Verify the result matches the expected error state.
			if tt.expectError {
				assert.False(t, result, "Should return false on error")
				assert.True(t, resp.Diagnostics.HasError(), "Should have diagnostics error")
			} else {
				assert.True(t, result, "Should return true on success")
				assert.False(t, resp.Diagnostics.HasError(), "Should not have diagnostics error")
				assert.Len(t, data.Executions, tt.expectCount, "Execution count should match")
			}
		})
	}
}
