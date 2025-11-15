package variable

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/client"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/variable/models"
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

// TestVariableResource_executeCreateLogic tests the executeCreateLogic method with error cases.
func TestVariableResource_executeCreateLogic(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		variableKey  string
		variableVal  string
		setupHandler func(w http.ResponseWriter, r *http.Request)
		expectError  bool
		expectID     string
	}{
		{
			name:        "successful creation",
			variableKey: "TEST_VAR",
			variableVal: "test_value",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPost && r.URL.Path == "/variables" {
					w.WriteHeader(http.StatusCreated)
					return
				}
				if r.Method == http.MethodGet && r.URL.Path == "/variables" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{
						"data": []map[string]any{
							{
								"id":    "var-123",
								"key":   "TEST_VAR",
								"value": "test_value",
								"type":  "string",
							},
						},
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: false,
			expectID:    "var-123",
		},
		{
			name:        "API error on create",
			variableKey: "FAILED_VAR",
			variableVal: "value",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"message": "Internal server error"}`))
			},
			expectError: true,
		},
		{
			name:        "error finding created variable",
			variableKey: "MISSING_VAR",
			variableVal: "value",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPost && r.URL.Path == "/variables" {
					w.WriteHeader(http.StatusCreated)
					return
				}
				if r.Method == http.MethodGet && r.URL.Path == "/variables" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{
						"data": []map[string]any{},
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := http.HandlerFunc(tt.setupHandler)
			n8nClient, server := setupTestClient(t, handler)
			defer server.Close()

			r := &VariableResource{client: n8nClient}
			ctx := context.Background()
			plan := &models.Resource{
				Key:   types.StringValue(tt.variableKey),
				Value: types.StringValue(tt.variableVal),
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
				assert.Equal(t, tt.expectID, plan.ID.ValueString(), "Variable ID should match")
			}
		})
	}
}

// TestVariableResource_executeReadLogic tests the executeReadLogic method with error cases.
func TestVariableResource_executeReadLogic(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		variableID   string
		setupHandler func(w http.ResponseWriter, r *http.Request)
		expectError  bool
		expectKey    string
	}{
		{
			name:       "successful read",
			variableID: "var-123",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet && r.URL.Path == "/variables" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{
						"data": []map[string]any{
							{
								"id":    "var-123",
								"key":   "RETRIEVED_VAR",
								"value": "retrieved_value",
								"type":  "string",
							},
						},
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: false,
			expectKey:   "RETRIEVED_VAR",
		},
		{
			name:       "API error on list",
			variableID: "var-500",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"message": "Internal server error"}`))
			},
			expectError: true,
		},
		{
			name:       "API error not found",
			variableID: "var-404",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`{"message": "Variable not found"}`))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := http.HandlerFunc(tt.setupHandler)
			n8nClient, server := setupTestClient(t, handler)
			defer server.Close()

			r := &VariableResource{client: n8nClient}
			ctx := context.Background()
			state := &models.Resource{
				ID: types.StringValue(tt.variableID),
			}
			resp := &resource.ReadResponse{
				State: resource.ReadResponse{}.State,
			}

			result := r.executeReadLogic(ctx, state, resp)

			if tt.expectError {
				assert.False(t, result, "Should return false on error")
			} else {
				assert.True(t, result, "Should return true on success")
				assert.False(t, resp.Diagnostics.HasError(), "Should not have diagnostics error")
				assert.Equal(t, tt.expectKey, state.Key.ValueString(), "Variable key should match")
			}
		})
	}
}

// TestVariableResource_executeUpdateLogic tests the executeUpdateLogic method with error cases.
func TestVariableResource_executeUpdateLogic(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		variableID   string
		newKey       string
		newValue     string
		setupHandler func(w http.ResponseWriter, r *http.Request)
		expectError  bool
	}{
		{
			name:       "successful update",
			variableID: "var-123",
			newKey:     "UPDATED_VAR",
			newValue:   "updated_value",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPut && r.URL.Path == "/variables/var-123" {
					w.WriteHeader(http.StatusNoContent)
					return
				}
				if r.Method == http.MethodGet && r.URL.Path == "/variables" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{
						"data": []map[string]any{
							{
								"id":    "var-123",
								"key":   "UPDATED_VAR",
								"value": "updated_value",
								"type":  "string",
							},
						},
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: false,
		},
		{
			name:       "variable not found",
			variableID: "var-404",
			newKey:     "UPDATED_VAR",
			newValue:   "value",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`{"message": "Variable not found"}`))
			},
			expectError: true,
		},
		{
			name:       "API error",
			variableID: "var-500",
			newKey:     "UPDATED_VAR",
			newValue:   "value",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"message": "Internal server error"}`))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := http.HandlerFunc(tt.setupHandler)
			n8nClient, server := setupTestClient(t, handler)
			defer server.Close()

			r := &VariableResource{client: n8nClient}
			ctx := context.Background()
			plan := &models.Resource{
				ID:    types.StringValue(tt.variableID),
				Key:   types.StringValue(tt.newKey),
				Value: types.StringValue(tt.newValue),
			}
			resp := &resource.UpdateResponse{
				State: resource.UpdateResponse{}.State,
			}

			result := r.executeUpdateLogic(ctx, plan, resp)

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

// TestVariableResource_executeDeleteLogic tests the executeDeleteLogic method with error cases.
func TestVariableResource_executeDeleteLogic(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		variableID   string
		setupHandler func(w http.ResponseWriter, r *http.Request)
		expectError  bool
	}{
		{
			name:       "successful deletion",
			variableID: "var-123",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodDelete && r.URL.Path == "/variables/var-123" {
					w.WriteHeader(http.StatusNoContent)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: false,
		},
		{
			name:       "variable not found",
			variableID: "var-404",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`{"message": "Variable not found"}`))
			},
			expectError: true,
		},
		{
			name:       "API error",
			variableID: "var-500",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"message": "Internal server error"}`))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := http.HandlerFunc(tt.setupHandler)
			n8nClient, server := setupTestClient(t, handler)
			defer server.Close()

			r := &VariableResource{client: n8nClient}
			ctx := context.Background()
			state := &models.Resource{
				ID: types.StringValue(tt.variableID),
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

// TestVariableResource_executeVariableCreate tests the executeVariableCreate helper method.
func TestVariableResource_executeVariableCreate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		variableKey  string
		variableVal  string
		setupHandler func(w http.ResponseWriter, r *http.Request)
		expectError  bool
	}{
		{
			name:        "successful create",
			variableKey: "TEST_VAR",
			variableVal: "test_value",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPost && r.URL.Path == "/variables" {
					w.WriteHeader(http.StatusCreated)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: false,
		},
		{
			name:        "API error 500",
			variableKey: "FAILED_VAR",
			variableVal: "value",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"message": "Internal server error"}`))
			},
			expectError: true,
		},
		{
			name:        "API error 400",
			variableKey: "INVALID_VAR",
			variableVal: "value",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"message": "Bad request"}`))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := http.HandlerFunc(tt.setupHandler)
			n8nClient, server := setupTestClient(t, handler)
			defer server.Close()

			r := &VariableResource{client: n8nClient}
			ctx := context.Background()
			plan := &models.Resource{
				Key:   types.StringValue(tt.variableKey),
				Value: types.StringValue(tt.variableVal),
			}
			resp := &resource.CreateResponse{
				State: resource.CreateResponse{}.State,
			}

			result := r.executeVariableCreate(ctx, plan, resp)

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

// TestVariableResource_findCreatedVariable tests the findCreatedVariable helper method.
func TestVariableResource_findCreatedVariable(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		variableKey  string
		setupHandler func(w http.ResponseWriter, r *http.Request)
		expectNil    bool
		expectID     string
	}{
		{
			name:        "successfully finds variable",
			variableKey: "TEST_VAR",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet && r.URL.Path == "/variables" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{
						"data": []map[string]any{
							{
								"id":    "var-123",
								"key":   "TEST_VAR",
								"value": "test_value",
								"type":  "string",
							},
						},
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectNil: false,
			expectID:  "var-123",
		},
		{
			name:        "API error on list",
			variableKey: "FAILED_VAR",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"message": "Internal server error"}`))
			},
			expectNil: true,
		},
		{
			name:        "variable not found in list",
			variableKey: "MISSING_VAR",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet && r.URL.Path == "/variables" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{
						"data": []map[string]any{
							{
								"id":    "var-999",
								"key":   "OTHER_VAR",
								"value": "other_value",
								"type":  "string",
							},
						},
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectNil: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := http.HandlerFunc(tt.setupHandler)
			n8nClient, server := setupTestClient(t, handler)
			defer server.Close()

			r := &VariableResource{client: n8nClient}
			ctx := context.Background()
			plan := &models.Resource{
				Key: types.StringValue(tt.variableKey),
			}
			resp := &resource.CreateResponse{
				State: resource.CreateResponse{}.State,
			}

			result := r.findCreatedVariable(ctx, plan, resp)

			if tt.expectNil {
				assert.Nil(t, result, "Should return nil")
				assert.True(t, resp.Diagnostics.HasError(), "Should have diagnostics error")
			} else {
				assert.NotNil(t, result, "Should return variable")
				assert.False(t, resp.Diagnostics.HasError(), "Should not have diagnostics error")
				assert.Equal(t, tt.expectID, *result.Id, "Variable ID should match")
			}
		})
	}
}

// TestVariableResource_updateStateFromVariable tests the updateStateFromVariable helper method.
func TestVariableResource_updateStateFromVariable(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		variable         *n8nsdk.Variable
		expectedKey      string
		expectedValue    string
		expectedType     string
		expectedTypeNull bool
		expectedProjID   string
		expectedProjNull bool
	}{
		{
			name: "updates all fields",
			variable: &n8nsdk.Variable{
				Id:    stringPtr("var-123"),
				Key:   "TEST_VAR",
				Value: "test_value",
				Type:  stringPtr("string"),
				Project: &n8nsdk.Project{
					Id: stringPtr("proj-456"),
				},
			},
			expectedKey:      "TEST_VAR",
			expectedValue:    "test_value",
			expectedType:     "string",
			expectedTypeNull: false,
			expectedProjID:   "proj-456",
			expectedProjNull: false,
		},
		{
			name: "handles nil type",
			variable: &n8nsdk.Variable{
				Id:    stringPtr("var-123"),
				Key:   "TEST_VAR",
				Value: "test_value",
				Type:  nil,
			},
			expectedKey:      "TEST_VAR",
			expectedValue:    "test_value",
			expectedTypeNull: true,
			expectedProjNull: true,
		},
		{
			name: "handles nil project",
			variable: &n8nsdk.Variable{
				Id:      stringPtr("var-123"),
				Key:     "TEST_VAR",
				Value:   "test_value",
				Type:    stringPtr("string"),
				Project: nil,
			},
			expectedKey:      "TEST_VAR",
			expectedValue:    "test_value",
			expectedType:     "string",
			expectedTypeNull: false,
			expectedProjNull: true,
		},
		{
			name: "error case - handles project with nil ID",
			variable: &n8nsdk.Variable{
				Id:    stringPtr("var-123"),
				Key:   "TEST_VAR",
				Value: "test_value",
				Type:  stringPtr("string"),
				Project: &n8nsdk.Project{
					Id: nil,
				},
			},
			expectedKey:      "TEST_VAR",
			expectedValue:    "test_value",
			expectedType:     "string",
			expectedTypeNull: false,
			expectedProjNull: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := &VariableResource{}
			state := &models.Resource{
				ID: types.StringValue("var-123"),
			}

			r.updateStateFromVariable(tt.variable, state)

			assert.Equal(t, tt.expectedKey, state.Key.ValueString(), "Key should match")
			assert.Equal(t, tt.expectedValue, state.Value.ValueString(), "Value should match")

			if tt.expectedTypeNull {
				assert.True(t, state.Type.IsNull(), "Type should be null")
			} else {
				assert.Equal(t, tt.expectedType, state.Type.ValueString(), "Type should match")
			}

			if tt.expectedProjNull {
				assert.True(t, state.ProjectID.IsNull(), "ProjectID should be null")
			} else {
				assert.Equal(t, tt.expectedProjID, state.ProjectID.ValueString(), "ProjectID should match")
			}
		})
	}
}

// TestVariableResource_executeVariableUpdate tests the executeVariableUpdate helper method.
func TestVariableResource_executeVariableUpdate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		variableID   string
		newKey       string
		newValue     string
		setupHandler func(w http.ResponseWriter, r *http.Request)
		expectError  bool
	}{
		{
			name:       "successful update",
			variableID: "var-123",
			newKey:     "UPDATED_VAR",
			newValue:   "updated_value",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPut && r.URL.Path == "/variables/var-123" {
					w.WriteHeader(http.StatusNoContent)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: false,
		},
		{
			name:       "API error 404",
			variableID: "var-404",
			newKey:     "UPDATED_VAR",
			newValue:   "value",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`{"message": "Variable not found"}`))
			},
			expectError: true,
		},
		{
			name:       "API error 500",
			variableID: "var-500",
			newKey:     "UPDATED_VAR",
			newValue:   "value",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"message": "Internal server error"}`))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := http.HandlerFunc(tt.setupHandler)
			n8nClient, server := setupTestClient(t, handler)
			defer server.Close()

			r := &VariableResource{client: n8nClient}
			ctx := context.Background()
			plan := &models.Resource{
				ID:    types.StringValue(tt.variableID),
				Key:   types.StringValue(tt.newKey),
				Value: types.StringValue(tt.newValue),
			}
			resp := &resource.UpdateResponse{
				State: resource.UpdateResponse{}.State,
			}

			result := r.executeVariableUpdate(ctx, plan, resp)

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

// TestVariableResource_findUpdatedVariable tests the findUpdatedVariable helper method.
func TestVariableResource_findUpdatedVariable(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		variableID   string
		setupHandler func(w http.ResponseWriter, r *http.Request)
		expectNil    bool
		expectedKey  string
	}{
		{
			name:       "successfully finds updated variable",
			variableID: "var-123",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet && r.URL.Path == "/variables" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{
						"data": []map[string]any{
							{
								"id":    "var-123",
								"key":   "UPDATED_VAR",
								"value": "updated_value",
								"type":  "string",
							},
						},
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectNil:   false,
			expectedKey: "UPDATED_VAR",
		},
		{
			name:       "API error on list",
			variableID: "var-500",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"message": "Internal server error"}`))
			},
			expectNil: true,
		},
		{
			name:       "variable not found in list",
			variableID: "var-404",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet && r.URL.Path == "/variables" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{
						"data": []map[string]any{
							{
								"id":    "var-999",
								"key":   "OTHER_VAR",
								"value": "other_value",
								"type":  "string",
							},
						},
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectNil: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := http.HandlerFunc(tt.setupHandler)
			n8nClient, server := setupTestClient(t, handler)
			defer server.Close()

			r := &VariableResource{client: n8nClient}
			ctx := context.Background()
			plan := &models.Resource{
				ID: types.StringValue(tt.variableID),
			}
			resp := &resource.UpdateResponse{
				State: resource.UpdateResponse{}.State,
			}

			result := r.findUpdatedVariable(ctx, plan, resp)

			if tt.expectNil {
				assert.Nil(t, result, "Should return nil")
				assert.True(t, resp.Diagnostics.HasError(), "Should have diagnostics error")
			} else {
				assert.NotNil(t, result, "Should return variable")
				assert.False(t, resp.Diagnostics.HasError(), "Should not have diagnostics error")
				assert.Equal(t, tt.expectedKey, result.Key, "Variable key should match")
			}
		})
	}
}

// stringPtr is a helper function to create string pointers for test cases.
func stringPtr(s string) *string {
	return &s
}
