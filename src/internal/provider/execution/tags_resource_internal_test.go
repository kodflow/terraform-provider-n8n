package execution

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/execution/models"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/client"
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

// makeTagIDsSet creates a types.Set of tag ID strings for testing.
func makeTagIDsSet(t *testing.T, ids []string) types.Set {
	t.Helper()
	val, diags := types.SetValueFrom(t.Context(), types.StringType, ids)
	// Check condition.
	if diags.HasError() {
		t.Fatalf("failed to create tag IDs set: %v", diags)
	}
	return val
}

// TestParseExecutionID tests the parseExecutionID helper.
func TestParseExecutionID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		input       string
		expectValue float32
		expectError bool
	}{
		{
			name:        "valid integer",
			input:       "42",
			expectValue: 42.0,
			expectError: false,
		},
		{
			name:        "valid float",
			input:       "3.14",
			expectValue: 3.14,
			expectError: false,
		},
		{
			name:        "invalid string",
			input:       "not-a-number",
			expectError: true,
		},
		{
			name:        "empty string",
			input:       "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result, err := parseExecutionID(tt.input)

			if tt.expectError {
				assert.Error(t, err, "Should return error for invalid input")
			} else {
				assert.NoError(t, err, "Should not return error")
				assert.InDelta(t, tt.expectValue, result, 0.01, "Value should match")
			}
		})
	}
}

// TestExecutionTagsResource_executeCreateLogic tests the executeCreateLogic method.
func TestExecutionTagsResource_executeCreateLogic(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		executionID  string
		tagIDs       []string
		setupHandler func(w http.ResponseWriter, r *http.Request)
		expectError  bool
	}{
		{
			name:        "successful create",
			executionID: "42",
			tagIDs:      []string{"tag-1", "tag-2"},
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPut && r.URL.Path == "/executions/42/tags" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode([]any{
						map[string]any{"id": "tag-1", "name": "Tag 1"},
						map[string]any{"id": "tag-2", "name": "Tag 2"},
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: false,
		},
		{
			name:        "invalid execution ID",
			executionID: "invalid",
			tagIDs:      []string{},
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
			expectError: true,
		},
		{
			name:        "API error",
			executionID: "99",
			tagIDs:      []string{"tag-1"},
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

			r := &ExecutionTagsResource{client: n8nClient}
			ctx := t.Context()
			plan := &models.TagsResource{
				ExecutionID: types.StringValue(tt.executionID),
				TagIDs:      makeTagIDsSet(t, tt.tagIDs),
			}
			resp := &resource.CreateResponse{
				State: resource.CreateResponse{}.State,
			}

			result := r.executeCreateLogic(ctx, plan, resp)

			if tt.expectError {
				assert.False(t, result, "Should return false on error")
				assert.True(t, resp.Diagnostics.HasError(), "Should have diagnostics error")
			} else {
				assert.True(t, result, "Should return true on success")
				assert.False(t, resp.Diagnostics.HasError(), "Should not have diagnostics error")
			}
		})
	}
}

// TestExecutionTagsResource_executeReadLogic tests the executeReadLogic method.
func TestExecutionTagsResource_executeReadLogic(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		executionID  string
		setupHandler func(w http.ResponseWriter, r *http.Request)
		expectError  bool
		expectCount  int
	}{
		{
			name:        "successful read",
			executionID: "42",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet && r.URL.Path == "/executions/42/tags" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					id1, id2 := "tag-1", "tag-2"
					json.NewEncoder(w).Encode([]any{
						map[string]any{"id": id1, "name": "Tag 1"},
						map[string]any{"id": id2, "name": "Tag 2"},
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: false,
			expectCount: 2,
		},
		{
			name:        "invalid execution ID",
			executionID: "not-a-number",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
			expectError: true,
		},
		{
			name:        "API error",
			executionID: "99",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"message": "Internal server error"}`))
			},
			expectError: true,
		},
		{
			name:        "empty tags",
			executionID: "55",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet && r.URL.Path == "/executions/55/tags" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode([]any{})
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

			r := &ExecutionTagsResource{client: n8nClient}
			ctx := t.Context()
			state := &models.TagsResource{
				ExecutionID: types.StringValue(tt.executionID),
				TagIDs:      makeTagIDsSet(t, []string{}),
			}
			resp := &resource.ReadResponse{
				State: resource.ReadResponse{}.State,
			}

			result := r.executeReadLogic(ctx, state, resp)

			if tt.expectError {
				assert.False(t, result, "Should return false on error")
				assert.True(t, resp.Diagnostics.HasError(), "Should have diagnostics error")
			} else {
				assert.True(t, result, "Should return true on success")
				assert.False(t, resp.Diagnostics.HasError(), "Should not have diagnostics error")
			}
		})
	}
}

// TestExecutionTagsResource_executeUpdateLogic tests the executeUpdateLogic method.
func TestExecutionTagsResource_executeUpdateLogic(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		executionID  string
		newTagIDs    []string
		setupHandler func(w http.ResponseWriter, r *http.Request)
		expectError  bool
	}{
		{
			name:        "successful update",
			executionID: "42",
			newTagIDs:   []string{"tag-3"},
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPut && r.URL.Path == "/executions/42/tags" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode([]any{
						map[string]any{"id": "tag-3", "name": "Tag 3"},
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: false,
		},
		{
			name:        "invalid execution ID",
			executionID: "bad",
			newTagIDs:   []string{},
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
			expectError: true,
		},
		{
			name:        "API error",
			executionID: "99",
			newTagIDs:   []string{"tag-1"},
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

			r := &ExecutionTagsResource{client: n8nClient}
			ctx := t.Context()
			plan := &models.TagsResource{
				ExecutionID: types.StringValue(tt.executionID),
				TagIDs:      makeTagIDsSet(t, tt.newTagIDs),
			}
			state := &models.TagsResource{
				ExecutionID: types.StringValue(tt.executionID),
				TagIDs:      makeTagIDsSet(t, []string{"tag-old"}),
			}
			resp := &resource.UpdateResponse{
				State: resource.UpdateResponse{}.State,
			}

			result := r.executeUpdateLogic(ctx, plan, state, resp)

			if tt.expectError {
				assert.False(t, result, "Should return false on error")
				assert.True(t, resp.Diagnostics.HasError(), "Should have diagnostics error")
			} else {
				assert.True(t, result, "Should return true on success")
				assert.False(t, resp.Diagnostics.HasError(), "Should not have diagnostics error")
			}
		})
	}
}

// TestExecutionTagsResource_executeDeleteLogic tests the executeDeleteLogic method.
func TestExecutionTagsResource_executeDeleteLogic(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		executionID  string
		setupHandler func(w http.ResponseWriter, r *http.Request)
		expectError  bool
	}{
		{
			name:        "successful delete",
			executionID: "42",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPut && r.URL.Path == "/executions/42/tags" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode([]any{})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: false,
		},
		{
			name:        "invalid execution ID",
			executionID: "invalid",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
			expectError: true,
		},
		{
			name:        "API error",
			executionID: "99",
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

			r := &ExecutionTagsResource{client: n8nClient}
			ctx := t.Context()
			state := &models.TagsResource{
				ExecutionID: types.StringValue(tt.executionID),
				TagIDs:      makeTagIDsSet(t, []string{"tag-1"}),
			}
			resp := &resource.DeleteResponse{
				State: resource.DeleteResponse{}.State,
			}

			result := r.executeDeleteLogic(ctx, state, resp)

			if tt.expectError {
				assert.False(t, result, "Should return false on error")
				assert.True(t, resp.Diagnostics.HasError(), "Should have diagnostics error")
			} else {
				assert.True(t, result, "Should return true on success")
				assert.False(t, resp.Diagnostics.HasError(), "Should not have diagnostics error")
			}
		})
	}
}

// TestCollectTagIDStrings tests the collectTagIDStrings helper.
func TestCollectTagIDStrings(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		tags     []n8nsdk.Tag
		expected []string
	}{
		{
			name:     "empty tags",
			tags:     []n8nsdk.Tag{},
			expected: []string{},
		},
		{
			name: "tags with IDs",
			tags: []n8nsdk.Tag{
				{Id: new("tag-1")},
				{Id: new("tag-2")},
			},
			expected: []string{"tag-1", "tag-2"},
		},
		{
			name: "tags with nil ID",
			tags: []n8nsdk.Tag{
				{Id: new("tag-1")},
				{Id: nil},
				{Id: new("tag-3")},
			},
			expected: []string{"tag-1", "tag-3"},
		},
		{
			name:     "nil tags slice",
			tags:     nil,
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := collectTagIDStrings(tt.tags)

			assert.Equal(t, tt.expected, result)
		})
	}
}
