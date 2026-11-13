package credential

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/credential/models"
	"github.com/stretchr/testify/assert"
)

// Test_convertValueByType tests the convertValueByType function.
func Test_convertValueByType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		key        string
		value      string
		propSchema map[string]any
		expected   any
	}{
		{
			name:       "convert string to number",
			key:        "port",
			value:      "993",
			propSchema: map[string]any{"type": "number"},
			expected:   float64(993),
		},
		{
			name:       "convert string to integer",
			key:        "count",
			value:      "42",
			propSchema: map[string]any{"type": "integer"},
			expected:   float64(42),
		},
		{
			name:       "convert string to boolean true",
			key:        "secure",
			value:      "true",
			propSchema: map[string]any{"type": "boolean"},
			expected:   true,
		},
		{
			name:       "convert string to boolean false",
			key:        "disabled",
			value:      "false",
			propSchema: map[string]any{"type": "boolean"},
			expected:   false,
		},
		{
			name:       "keep string for string type",
			key:        "host",
			value:      "smtp.example.com",
			propSchema: map[string]any{"type": "string"},
			expected:   "smtp.example.com",
		},
		{
			name:       "keep string when no type in schema",
			key:        "unknown",
			value:      "test",
			propSchema: map[string]any{},
			expected:   "test",
		},
		{
			name:       "keep string on invalid number",
			key:        "port",
			value:      "not-a-number",
			propSchema: map[string]any{"type": "number"},
			expected:   "not-a-number",
		},
		{
			name:       "keep string on invalid boolean",
			key:        "secure",
			value:      "not-a-bool",
			propSchema: map[string]any{"type": "boolean"},
			expected:   "not-a-bool",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			result := convertValueByType(ctx, tt.key, tt.value, tt.propSchema)

			assert.Equal(t, tt.expected, result, "Converted value should match expected")
		})
	}
}

// TestCredentialResource_applySchemaTypeConversions tests applySchemaTypeConversions.
func TestCredentialResource_applySchemaTypeConversions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		schema   map[string]any
		data     map[string]any
		expected map[string]any
	}{
		{
			name: "convert number and boolean fields",
			schema: map[string]any{
				"properties": map[string]any{
					"port":   map[string]any{"type": "number"},
					"secure": map[string]any{"type": "boolean"},
					"host":   map[string]any{"type": "string"},
				},
			},
			data: map[string]any{
				"port":   "993",
				"secure": "true",
				"host":   "imap.example.com",
			},
			expected: map[string]any{
				"port":   float64(993),
				"secure": true,
				"host":   "imap.example.com",
			},
		},
		{
			name:   "no properties in schema",
			schema: map[string]any{},
			data: map[string]any{
				"key": "value",
			},
			expected: map[string]any{
				"key": "value",
			},
		},
		{
			name: "non-string values are kept as-is",
			schema: map[string]any{
				"properties": map[string]any{
					"count": map[string]any{"type": "number"},
				},
			},
			data: map[string]any{
				"count": float64(42),
			},
			expected: map[string]any{
				"count": float64(42),
			},
		},
		{
			name: "unknown property kept as string",
			schema: map[string]any{
				"properties": map[string]any{},
			},
			data: map[string]any{
				"unknown": "value",
			},
			expected: map[string]any{
				"unknown": "value",
			},
		},
		{
			name: "error case - invalid number string returns original",
			schema: map[string]any{
				"properties": map[string]any{
					"port": map[string]any{"type": "number"},
				},
			},
			data: map[string]any{
				"port": "not-a-number",
			},
			expected: map[string]any{
				"port": "not-a-number",
			},
		},
		{
			name: "error case - invalid boolean string returns original",
			schema: map[string]any{
				"properties": map[string]any{
					"enabled": map[string]any{"type": "boolean"},
				},
			},
			data: map[string]any{
				"enabled": "maybe",
			},
			expected: map[string]any{
				"enabled": "maybe",
			},
		},
		{
			name: "error case - nil data returns empty map",
			schema: map[string]any{
				"properties": map[string]any{
					"port": map[string]any{"type": "number"},
				},
			},
			data:     nil,
			expected: map[string]any{},
		},
		{
			name:     "error case - nil schema returns original data",
			schema:   nil,
			data:     map[string]any{"key": "value"},
			expected: map[string]any{"key": "value"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := &CredentialResource{}
			ctx := context.Background()

			result := r.applySchemaTypeConversions(ctx, tt.schema, tt.data)

			assert.Equal(t, tt.expected, result, "Converted data should match expected")
		})
	}
}

// TestCredentialResource_convertDataToSchemaTypes tests convertDataToSchemaTypes.
func TestCredentialResource_convertDataToSchemaTypes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		credType     string
		data         map[string]any
		setupHandler func(w http.ResponseWriter, r *http.Request)
		expected     map[string]any
	}{
		{
			name:     "successful schema fetch and conversion",
			credType: "imap",
			data: map[string]any{
				"port":   "993",
				"secure": "true",
			},
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				// Handle schema request.
				if r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/credentials/schema/") {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{
						"properties": map[string]any{
							"port":   map[string]any{"type": "number"},
							"secure": map[string]any{"type": "boolean"},
						},
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expected: map[string]any{
				"port":   float64(993),
				"secure": true,
			},
		},
		{
			name:     "schema fetch fails - returns original data",
			credType: "unknown",
			data: map[string]any{
				"port": "993",
			},
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			expected: map[string]any{
				"port": "993",
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

			result := r.convertDataToSchemaTypes(ctx, tt.credType, tt.data)

			assert.Equal(t, tt.expected, result, "Converted data should match expected")
		})
	}
}

// TestCredentialResource_transferCredentialToProject tests transferCredentialToProject.
func TestCredentialResource_transferCredentialToProject(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		credentialID string
		projectID    string
		setupHandler func(w http.ResponseWriter, r *http.Request)
		expectOK     bool
	}{
		{
			name:         "successful transfer",
			credentialID: "cred-123",
			projectID:    "proj-456",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				// Handle transfer request.
				if r.Method == http.MethodPut && r.URL.Path == "/credentials/cred-123/transfer" {
					w.WriteHeader(http.StatusNoContent)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectOK: true,
		},
		{
			name:         "credential already in the project",
			credentialID: "cred-123",
			projectID:    "proj-456",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"message": "You can't transfer a credential into the project that's already owning it."}`))
			},
			expectOK: true,
		},
		{
			name:         "transfer fails",
			credentialID: "cred-fail",
			projectID:    "proj-456",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"message": "Transfer failed"}`))
			},
			expectOK: false,
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

			result := r.transferCredentialToProject(ctx, tt.credentialID, tt.projectID, &diags)

			assert.Equal(t, tt.expectOK, result, "Transfer result should match expected")
			// Check diagnostics based on expected outcome.
			if !tt.expectOK {
				assert.True(t, diags.HasError(), "Should have diagnostics error on failure")
			} else {
				assert.False(t, diags.HasError(), "Should not have diagnostics error on success")
			}
		})
	}
}

// Test_mapCredentialProjectID tests mapCredentialProjectID.
func Test_mapCredentialProjectID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		requestedProjectID types.String
		expectedIsNull     bool
		expectedValue      string
	}{
		{
			name:               "project_id is set",
			requestedProjectID: types.StringValue("proj-123"),
			expectedIsNull:     false,
			expectedValue:      "proj-123",
		},
		{
			name:               "project_id is null",
			requestedProjectID: types.StringNull(),
			expectedIsNull:     true,
			expectedValue:      "",
		},
		{
			name:               "project_id is unknown",
			requestedProjectID: types.StringUnknown(),
			expectedIsNull:     true,
			expectedValue:      "",
		},
		{
			name:               "error case - empty string project_id",
			requestedProjectID: types.StringValue(""),
			expectedIsNull:     false,
			expectedValue:      "",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			plan := &models.Resource{}
			mapCredentialProjectID(plan, tt.requestedProjectID)

			if tt.expectedIsNull {
				assert.True(t, plan.ProjectID.IsNull(), "ProjectID should be null")
			} else {
				assert.False(t, plan.ProjectID.IsNull(), "ProjectID should not be null")
				assert.Equal(t, tt.expectedValue, plan.ProjectID.ValueString(), "ProjectID value should match")
			}
		})
	}
}
