package tag

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/client"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/tag/models"
	"github.com/stretchr/testify/assert"
)

func TestTagDataSource_validateIdentifier(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		data        *models.DataSource
		expectValid bool
	}{
		{
			name: "valid with ID",
			data: &models.DataSource{
				ID: types.StringValue("tag-123"),
			},
			expectValid: true,
		},
		{
			name: "valid with name",
			data: &models.DataSource{
				Name: types.StringValue("Test Tag"),
			},
			expectValid: true,
		},
		{
			name: "valid with both ID and name",
			data: &models.DataSource{
				ID:   types.StringValue("tag-123"),
				Name: types.StringValue("Test Tag"),
			},
			expectValid: true,
		},
		{
			name: "invalid with both null",
			data: &models.DataSource{
				ID:   types.StringNull(),
				Name: types.StringNull(),
			},
			expectValid: false,
		},
		{
			name:        "invalid with uninitialized values",
			data:        &models.DataSource{},
			expectValid: false,
		},
		{
			name: "valid with empty string ID",
			data: &models.DataSource{
				ID: types.StringValue(""),
			},
			expectValid: true,
		},
		{
			name: "valid with empty string name",
			data: &models.DataSource{
				Name: types.StringValue(""),
			},
			expectValid: true,
		},
		{
			name: "valid with unknown ID",
			data: &models.DataSource{
				ID: types.StringUnknown(),
			},
			expectValid: true,
		},
		{
			name: "valid with unknown name",
			data: &models.DataSource{
				Name: types.StringUnknown(),
			},
			expectValid: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			d := &TagDataSource{}
			resp := &datasource.ReadResponse{}

			valid := d.validateIdentifier(tt.data, resp)

			if tt.expectValid {
				assert.True(t, valid)
				assert.False(t, resp.Diagnostics.HasError())
			} else {
				assert.False(t, valid)
				assert.True(t, resp.Diagnostics.HasError())
			}
		})
	}
}

func TestTagDataSource_fetchTagByID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		handler     http.HandlerFunc
		tagID       string
		expectNil   bool
		expectError bool
	}{
		{
			name:  "tag found by ID",
			tagID: "tag-123",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/tags/tag-123" && r.Method == http.MethodGet {
					id := "tag-123"
					response := map[string]interface{}{
						"id":   id,
						"name": "Test Tag",
					}
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(response)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			expectNil:   false,
			expectError: false,
		},
		{
			name:  "tag not found",
			tagID: "tag-999",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/tags/tag-999" && r.Method == http.MethodGet {
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte(`{"message": "Tag not found"}`))
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			expectNil:   true,
			expectError: true,
		},
		{
			name:  "API returns error",
			tagID: "tag-123",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/tags/tag-123" && r.Method == http.MethodGet {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"message": "Internal server error"}`))
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			expectNil:   true,
			expectError: true,
		},
		{
			name:  "bad request error",
			tagID: "invalid-id",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/tags/invalid-id" && r.Method == http.MethodGet {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(`{"message": "Bad request"}`))
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			expectNil:   true,
			expectError: true,
		},
		{
			name:  "unauthorized error",
			tagID: "tag-123",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/tags/tag-123" && r.Method == http.MethodGet {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte(`{"message": "Unauthorized"}`))
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			expectNil:   true,
			expectError: true,
		},
		{
			name:  "forbidden error",
			tagID: "tag-123",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/tags/tag-123" && r.Method == http.MethodGet {
					w.WriteHeader(http.StatusForbidden)
					w.Write([]byte(`{"message": "Forbidden"}`))
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			expectNil:   true,
			expectError: true,
		},
		{
			name:  "empty tag ID",
			tagID: "",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/tags/" && r.Method == http.MethodGet {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			expectNil:   true,
			expectError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			n8nClient, server := setupTagTestClient(t, tt.handler)
			defer server.Close()

			d := &TagDataSource{client: n8nClient}
			data := &models.DataSource{}
			data.ID = types.StringValue(tt.tagID)
			resp := &datasource.ReadResponse{}

			tag := d.fetchTagByID(context.Background(), data, resp)

			if tt.expectNil {
				assert.Nil(t, tag)
			} else {
				assert.NotNil(t, tag)
			}

			if tt.expectError {
				assert.True(t, resp.Diagnostics.HasError())
			} else {
				assert.False(t, resp.Diagnostics.HasError())
			}
		})
	}
}

func TestTagDataSource_fetchTagByName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		handler     http.HandlerFunc
		tagName     string
		expectNil   bool
		expectError bool
	}{
		{
			name:    "tag found by name",
			tagName: "Test Tag",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/tags" && r.Method == http.MethodGet {
					id := "tag-123"
					response := map[string]interface{}{
						"data": []interface{}{
							map[string]interface{}{
								"id":   id,
								"name": "Test Tag",
							},
						},
					}
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(response)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			expectNil:   false,
			expectError: false,
		},
		{
			name:    "tag not found in list",
			tagName: "Missing Tag",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/tags" && r.Method == http.MethodGet {
					response := map[string]interface{}{
						"data": []interface{}{
							map[string]interface{}{
								"id":   "tag-456",
								"name": "Other Tag",
							},
						},
					}
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(response)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			expectNil:   true,
			expectError: true,
		},
		{
			name:    "API returns error",
			tagName: "Test Tag",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/tags" && r.Method == http.MethodGet {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"message": "Internal server error"}`))
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			expectNil:   true,
			expectError: true,
		},
		{
			name:    "empty tag list",
			tagName: "Test Tag",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/tags" && r.Method == http.MethodGet {
					response := map[string]interface{}{
						"data": []interface{}{},
					}
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(response)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			expectNil:   true,
			expectError: true,
		},
		{
			name:    "nil data in response",
			tagName: "Test Tag",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/tags" && r.Method == http.MethodGet {
					response := map[string]interface{}{}
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(response)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			expectNil:   true,
			expectError: true,
		},
		{
			name:    "case sensitive name matching",
			tagName: "Test Tag",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/tags" && r.Method == http.MethodGet {
					response := map[string]interface{}{
						"data": []interface{}{
							map[string]interface{}{
								"id":   "tag-123",
								"name": "test tag",
							},
						},
					}
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(response)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			expectNil:   true,
			expectError: true,
		},
		{
			name:    "multiple tags in list",
			tagName: "Second Tag",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/tags" && r.Method == http.MethodGet {
					response := map[string]interface{}{
						"data": []interface{}{
							map[string]interface{}{
								"id":   "tag-123",
								"name": "First Tag",
							},
							map[string]interface{}{
								"id":   "tag-456",
								"name": "Second Tag",
							},
							map[string]interface{}{
								"id":   "tag-789",
								"name": "Third Tag",
							},
						},
					}
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(response)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			expectNil:   false,
			expectError: false,
		},
		{
			name:    "empty string name",
			tagName: "",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/tags" && r.Method == http.MethodGet {
					response := map[string]interface{}{
						"data": []interface{}{
							map[string]interface{}{
								"id":   "tag-123",
								"name": "",
							},
						},
					}
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(response)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			expectNil:   false,
			expectError: false,
		},
		{
			name:    "special characters in name",
			tagName: "Tag-!@#$%",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/tags" && r.Method == http.MethodGet {
					response := map[string]interface{}{
						"data": []interface{}{
							map[string]interface{}{
								"id":   "tag-123",
								"name": "Tag-!@#$%",
							},
						},
					}
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(response)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			expectNil:   false,
			expectError: false,
		},
		{
			name:    "unicode characters in name",
			tagName: "标签测试",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/tags" && r.Method == http.MethodGet {
					response := map[string]interface{}{
						"data": []interface{}{
							map[string]interface{}{
								"id":   "tag-123",
								"name": "标签测试",
							},
						},
					}
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(response)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			expectNil:   false,
			expectError: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			n8nClient, server := setupTagTestClient(t, tt.handler)
			defer server.Close()

			d := &TagDataSource{client: n8nClient}
			data := &models.DataSource{}
			data.Name = types.StringValue(tt.tagName)
			resp := &datasource.ReadResponse{}

			tag := d.fetchTagByName(context.Background(), data, resp)

			if tt.expectNil {
				assert.Nil(t, tag)
			} else {
				assert.NotNil(t, tag)
			}

			if tt.expectError {
				assert.True(t, resp.Diagnostics.HasError())
			} else {
				assert.False(t, resp.Diagnostics.HasError())
			}
		})
	}
}

func setupTagTestClient(t *testing.T, handler http.HandlerFunc) (*client.N8nClient, *httptest.Server) {
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
