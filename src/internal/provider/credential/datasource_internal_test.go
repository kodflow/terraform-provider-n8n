package credential

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/credential/models"
	"github.com/stretchr/testify/assert"
)

// TestCredentialDataSource_executeReadLogic tests the executeReadLogic method.
func TestCredentialDataSource_executeReadLogic(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		id           string
		credName     string
		setupHandler func(w http.ResponseWriter, r *http.Request)
		expectError  bool
		expectType   string
	}{
		{
			name:     "found by ID",
			id:       "cred-123",
			credName: "",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet && r.URL.Path == "/credentials" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{
						"data": []any{
							map[string]any{
								"id":        "cred-123",
								"name":      "My Cred",
								"type":      "httpHeaderAuth",
								"createdAt": "2024-01-01T00:00:00Z",
								"updatedAt": "2024-01-01T00:00:00Z",
							},
						},
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: false,
			expectType:  "httpHeaderAuth",
		},
		{
			name:     "found by name",
			id:       "",
			credName: "My Cred",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet && r.URL.Path == "/credentials" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{
						"data": []any{
							map[string]any{
								"id":        "cred-456",
								"name":      "My Cred",
								"type":      "httpBasicAuth",
								"createdAt": "2024-01-01T00:00:00Z",
								"updatedAt": "2024-01-01T00:00:00Z",
							},
						},
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: false,
			expectType:  "httpBasicAuth",
		},
		{
			name:     "not found",
			id:       "cred-missing",
			credName: "",
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
			expectError: true,
		},
		{
			name:     "API error",
			id:       "cred-500",
			credName: "",
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

			d := &CredentialDataSource{client: n8nClient}
			ctx := t.Context()

			var idVal, nameVal types.String
			// Check condition.
			if tt.id != "" {
				idVal = types.StringValue(tt.id)
			} else {
				idVal = types.StringNull()
			}
			// Check condition.
			if tt.credName != "" {
				nameVal = types.StringValue(tt.credName)
			} else {
				nameVal = types.StringNull()
			}

			data := &models.DataSource{
				ID:   idVal,
				Name: nameVal,
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
				assert.Equal(t, tt.expectType, data.Type.ValueString(), "Type should match")
			}
		})
	}
}

// TestFindCredentialByIDOrName tests the findCredentialByIDOrName helper.
func TestFindCredentialByIDOrName(t *testing.T) {
	t.Parallel()

	credentials := []n8nsdk.CredentialListItem{
		{ID: "c1", Name: "Alpha", Type: "typeA"},
		{ID: "c2", Name: "Beta", Type: "typeB"},
	}

	tests := []struct {
		name      string
		id        types.String
		credName  types.String
		expectID  string
		expectHit bool
	}{
		{
			name:      "found by ID",
			id:        types.StringValue("c1"),
			credName:  types.StringNull(),
			expectID:  "c1",
			expectHit: true,
		},
		{
			name:      "found by name",
			id:        types.StringNull(),
			credName:  types.StringValue("Beta"),
			expectID:  "c2",
			expectHit: true,
		},
		{
			name:      "not found",
			id:        types.StringValue("c-missing"),
			credName:  types.StringNull(),
			expectHit: false,
		},
		{
			name:      "both null",
			id:        types.StringNull(),
			credName:  types.StringNull(),
			expectHit: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result, found := findCredentialByIDOrName(credentials, tt.id, tt.credName)

			assert.Equal(t, tt.expectHit, found)
			// Check condition.
			if tt.expectHit {
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectID, result.ID)
			} else {
				assert.Nil(t, result)
			}
		})
	}
}

// TestCredentialDataSource_resolveCredentialFromItems tests the resolveCredentialFromItems method.
func TestCredentialDataSource_resolveCredentialFromItems(t *testing.T) {
	t.Parallel()

	baseTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	items := []n8nsdk.CredentialListItem{
		{ID: "c1", Name: "Alpha", Type: "typeA", CreatedAt: baseTime, UpdatedAt: baseTime},
		{ID: "c2", Name: "Beta", Type: "typeB", CreatedAt: baseTime, UpdatedAt: baseTime},
	}

	tests := []struct {
		name     string
		id       types.String
		credName types.String
		useItems bool
		expectOK bool
		expectID string
	}{
		{
			name:     "found by ID",
			id:       types.StringValue("c1"),
			credName: types.StringNull(),
			useItems: true,
			expectOK: true,
			expectID: "c1",
		},
		{
			name:     "found by name",
			id:       types.StringNull(),
			credName: types.StringValue("Beta"),
			useItems: true,
			expectOK: true,
			expectID: "c2",
		},
		{
			name:     "not found returns error",
			id:       types.StringValue("missing"),
			credName: types.StringNull(),
			useItems: true,
			expectOK: false,
		},
		{
			name:     "error case - empty items list",
			id:       types.StringValue("c1"),
			credName: types.StringNull(),
			useItems: false,
			expectOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			d := &CredentialDataSource{}
			data := &models.DataSource{
				ID:   tt.id,
				Name: tt.credName,
			}
			resp := &datasource.ReadResponse{}

			var searchItems []n8nsdk.CredentialListItem
			//: Use items only for non-empty-items test cases.
			if tt.useItems {
				searchItems = items
			}

			ok := d.resolveCredentialFromItems(searchItems, data, resp)
			assert.Equal(t, tt.expectOK, ok)
			//: Verify diagnostics are set on failure.
			if !tt.expectOK {
				assert.True(t, resp.Diagnostics.HasError())
			}
			//: Verify ID was populated on success.
			if tt.expectOK {
				assert.Equal(t, tt.expectID, data.ID.ValueString())
			}
		})
	}
}
