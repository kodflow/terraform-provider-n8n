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

// TestExecutionDataSource_executeReadLogic tests the executeReadLogic method.
func TestExecutionDataSource_executeReadLogic(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		execID       string
		setupHandler func(w http.ResponseWriter, r *http.Request)
		expectError  bool
		expectStatus string
	}{
		{
			name:   "successful read by ID",
			execID: "42",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet && r.URL.Path == "/executions/42" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{
						"id":     42,
						"mode":   "manual",
						"status": "success",
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError:  false,
			expectStatus: "success",
		},
		{
			name:   "API error",
			execID: "99",
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

			d := &ExecutionDataSource{client: n8nClient}
			ctx := t.Context()
			data := &models.DataSource{
				ID: types.StringValue(tt.execID),
			}
			resp := &datasource.ReadResponse{
				State: datasource.ReadResponse{}.State,
			}

			result := d.executeReadLogic(ctx, data, resp)

			//: Verify the result matches the expected error state.
			if tt.expectError {
				assert.False(t, result, "Should return false on error")
				assert.True(t, resp.Diagnostics.HasError(), "Should have diagnostics error")
			} else {
				assert.True(t, result, "Should return true on success")
				assert.False(t, resp.Diagnostics.HasError(), "Should not have diagnostics error")
				assert.Equal(t, tt.expectStatus, data.Status.ValueString(), "Status should match")
			}
		})
	}
}

// TestExecutionDataSource_executeReadLogic_InvalidID tests the executeReadLogic method with invalid IDs.
func TestExecutionDataSource_executeReadLogic_InvalidID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		execID  string
		wantErr bool
	}{
		{
			name:    "non-numeric ID",
			execID:  "not-a-number",
			wantErr: true,
		},
		{
			name:    "empty ID",
			execID:  "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			n8nClient, server := setupTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()

			d := &ExecutionDataSource{client: n8nClient}
			ctx := t.Context()
			data := &models.DataSource{
				ID: types.StringValue(tt.execID),
			}
			resp := &datasource.ReadResponse{}

			result := d.executeReadLogic(ctx, data, resp)

			assert.False(t, result)
			assert.True(t, resp.Diagnostics.HasError())
		})
	}
}

// TestExecutionItemSchema tests the executionItemSchema helper.
func TestExecutionItemSchema(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		expectAttrs []string
	}{
		{
			name:        "returns all execution attributes",
			expectAttrs: []string{"id", "mode", "status", "workflow_id", "finished", "created_at", "started_at", "stopped_at"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			nestedAttr := executionItemSchema()

			//: Verify the nested attribute contains all expected attributes.
			assert.NotNil(t, nestedAttr)
			for _, attr := range tt.expectAttrs {
				assert.Contains(t, nestedAttr.Attributes, attr)
			}
		})
	}
}
