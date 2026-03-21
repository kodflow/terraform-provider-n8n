package workflow

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/client"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/workflow/models"
	"github.com/stretchr/testify/assert"
)

// setupVersionTestClient creates a test N8nClient with httptest server for workflow version tests.
func setupVersionTestClient(t *testing.T, handler http.HandlerFunc) (*client.N8nClient, *httptest.Server) {
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

// TestWorkflowVersionDataSource_executeReadLogic tests the executeReadLogic method.
func TestWorkflowVersionDataSource_executeReadLogic(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		workflowID   string
		versionID    string
		setupHandler func(w http.ResponseWriter, r *http.Request)
		expectError  bool
		expectName   string
	}{
		{
			name:       "successful read",
			workflowID: "wf-123",
			versionID:  "ver-1",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet && r.URL.Path == "/workflows/wf-123/ver-1" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					wfName := "My Workflow"
					json.NewEncoder(w).Encode(map[string]any{
						"versionId":   "ver-1",
						"workflowId":  "wf-123",
						"name":        wfName,
						"authors":     "user1",
						"nodes":       []any{},
						"connections": map[string]any{},
						"createdAt":   "2024-01-01T00:00:00Z",
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: false,
			expectName:  "My Workflow",
		},
		{
			name:       "version not found",
			workflowID: "wf-123",
			versionID:  "ver-404",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`{"message": "Version not found"}`))
			},
			expectError: true,
		},
		{
			name:       "API error",
			workflowID: "wf-123",
			versionID:  "ver-500",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"message": "Internal server error"}`))
			},
			expectError: true,
		},
		{
			name:       "nil name and authors",
			workflowID: "wf-456",
			versionID:  "ver-2",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet && r.URL.Path == "/workflows/wf-456/ver-2" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{
						"versionId":   "ver-2",
						"workflowId":  "wf-456",
						"nodes":       []any{},
						"connections": map[string]any{},
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: false,
			expectName:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			n8nClient, server := setupVersionTestClient(t, http.HandlerFunc(tt.setupHandler))
			defer server.Close()

			d := &WorkflowVersionDataSource{client: n8nClient}
			ctx := t.Context()
			data := &models.DataSourceVersion{
				WorkflowID: types.StringValue(tt.workflowID),
				VersionID:  types.StringValue(tt.versionID),
			}
			resp := &datasource.ReadResponse{
				State: datasource.ReadResponse{}.State,
			}

			result := d.executeReadLogic(ctx, data, resp)

			if tt.expectError {
				assert.False(t, result, "Should return false on error")
				assert.True(t, resp.Diagnostics.HasError(), "Should have diagnostics error")
			} else {
				assert.True(t, result, "Should return true on success")
				assert.False(t, resp.Diagnostics.HasError(), "Should not have diagnostics error")
				// Check condition.
				if tt.expectName != "" {
					assert.Equal(t, tt.expectName, data.Name.ValueString(), "Name should match")
				}
			}
		})
	}
}

// TestSerializeVersionNodes tests the serializeVersionNodes function.
func TestSerializeVersionNodes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		nodes       []n8nsdk.Node
		expectError bool
		expectJSON  string
	}{
		{
			name:        "nil nodes serialized as empty array",
			nodes:       nil,
			expectError: false,
			expectJSON:  "null",
		},
		{
			name:        "empty nodes",
			nodes:       []n8nsdk.Node{},
			expectError: false,
			expectJSON:  "[]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			data := &models.DataSourceVersion{}
			resp := &datasource.ReadResponse{}

			result := serializeVersionNodes(tt.nodes, data, resp)

			if tt.expectError {
				assert.False(t, result)
				assert.True(t, resp.Diagnostics.HasError())
			} else {
				assert.True(t, result)
				assert.False(t, resp.Diagnostics.HasError())
				assert.Equal(t, tt.expectJSON, data.NodesJSON.ValueString())
			}
		})
	}
}

// TestMapVersionOptionalFields tests the mapVersionOptionalFields function.
func TestMapVersionOptionalFields(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		nameArg    *string
		authorsArg *string
		expectName string
		expectNil  bool
	}{
		{
			name:       "all fields set",
			nameArg:    new("Test Workflow"),
			authorsArg: new("Alice"),
			expectName: "Test Workflow",
			expectNil:  false,
		},
		{
			name:       "nil name and authors",
			nameArg:    nil,
			authorsArg: nil,
			expectNil:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			data := &models.DataSourceVersion{}
			//: Call mapVersionOptionalFields with test arguments.
			mapVersionOptionalFields(tt.nameArg, tt.authorsArg, nil, data)
			//: Verify expected state matches.
			if tt.expectNil {
				assert.True(t, data.Name.IsNull())
				assert.True(t, data.Authors.IsNull())
			} else {
				assert.Equal(t, tt.expectName, data.Name.ValueString())
				assert.Equal(t, "Alice", data.Authors.ValueString())
			}
			assert.True(t, data.CreatedAt.IsNull())
		})
	}
}

// TestSerializeVersionFields tests the serializeVersionFields function.
func TestSerializeVersionFields(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		version     *n8nsdk.WorkflowVersion
		expectError bool
		expectNodes string
	}{
		{
			name: "successful serialization",
			version: &n8nsdk.WorkflowVersion{
				Nodes:       []n8nsdk.Node{},
				Connections: map[string]any{},
			},
			expectError: false,
			expectNodes: "[]",
		},
		{
			name: "nil nodes serialized as null",
			version: &n8nsdk.WorkflowVersion{
				Nodes:       nil,
				Connections: map[string]any{},
			},
			expectError: false,
			expectNodes: "null",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			data := &models.DataSourceVersion{}
			resp := &datasource.ReadResponse{}

			//: Call serializeVersionFields with test version.
			result := serializeVersionFields(tt.version, data, resp)

			//: Verify expected result.
			if tt.expectError {
				assert.False(t, result)
				assert.True(t, resp.Diagnostics.HasError())
			} else {
				assert.True(t, result)
				assert.False(t, resp.Diagnostics.HasError())
				assert.Equal(t, tt.expectNodes, data.NodesJSON.ValueString())
			}
		})
	}
}
