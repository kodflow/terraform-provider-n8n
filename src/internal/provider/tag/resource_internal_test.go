package tag

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
	"github.com/stretchr/testify/assert"
)

// KTN-TEST-009 RATIONALE: Tests in this file test public functions that require
// access to private fields (specifically the .client field). These tests MUST
// remain in the internal test file to access the unexported client field for
// proper white-box testing.

func TestTagResource_CreateAPICall(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setupMock func() *httptest.Server
		tagName   string
		wantError bool
		checkFunc func(*testing.T, *n8nsdk.Tag, error)
	}{
		{
			name: "successful creation",
			setupMock: func() *httptest.Server {
				id := "tag-123"
				now := time.Now()
				tag := &n8nsdk.Tag{
					Id:        &id,
					Name:      "Production",
					CreatedAt: &now,
					UpdatedAt: &now,
				}

				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/api/v1/tags" && r.Method == http.MethodPost {
						w.WriteHeader(http.StatusCreated)
						json.NewEncoder(w).Encode(tag)
						return
					}
					w.WriteHeader(http.StatusNotFound)
				}))
			},
			tagName:   "Production",
			wantError: false,
			checkFunc: func(t *testing.T, tag *n8nsdk.Tag, err error) {
				t.Helper()
				assert.NoError(t, err)
				assert.NotNil(t, tag)
				assert.Equal(t, "tag-123", *tag.Id)
				assert.Equal(t, "Production", tag.Name)
			},
		},
		{
			name: "creation with nil timestamps",
			setupMock: func() *httptest.Server {
				id := "tag-456"
				tag := &n8nsdk.Tag{
					Id:        &id,
					Name:      "Development",
					CreatedAt: nil,
					UpdatedAt: nil,
				}

				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusCreated)
					json.NewEncoder(w).Encode(tag)
				}))
			},
			tagName:   "Development",
			wantError: false,
			checkFunc: func(t *testing.T, tag *n8nsdk.Tag, err error) {
				t.Helper()
				assert.NoError(t, err)
				assert.NotNil(t, tag)
				assert.Nil(t, tag.CreatedAt)
				assert.Nil(t, tag.UpdatedAt)
			},
		},
		{
			name: "API error during creation",
			setupMock: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"message": "Internal server error"}`))
				}))
			},
			tagName:   "FailTag",
			wantError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := tt.setupMock()
			defer server.Close()

			apiClient := n8nsdk.NewAPIClient(&n8nsdk.Configuration{
				Servers: n8nsdk.ServerConfigurations{
					{URL: server.URL + "/api/v1"},
				},
			})

			r := &TagResource{
				client: &client.N8nClient{
					APIClient: apiClient,
				},
			}

			tagRequest := n8nsdk.Tag{
				Name: tt.tagName,
			}

			tag, httpResp, err := r.client.APIClient.TagsAPI.TagsPost(context.Background()).
				Tag(tagRequest).
				Execute()
			if httpResp != nil && httpResp.Body != nil {
				defer httpResp.Body.Close()
			}

			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.checkFunc != nil {
					tt.checkFunc(t, tag, err)
				}
			}
		})
	}
}

func TestTagResource_ReadAPICall(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setupMock func() *httptest.Server
		tagID     string
		wantError bool
		checkFunc func(*testing.T, *n8nsdk.Tag, error)
	}{
		{
			name: "successful read",
			setupMock: func() *httptest.Server {
				id := "tag-123"
				now := time.Now()
				tag := &n8nsdk.Tag{
					Id:        &id,
					Name:      "Production",
					CreatedAt: &now,
					UpdatedAt: &now,
				}

				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/tags/tag-123" && r.Method == http.MethodGet {
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusOK)
						json.NewEncoder(w).Encode(tag)
						return
					}
					w.WriteHeader(http.StatusNotFound)
				}))
			},
			tagID:     "tag-123",
			wantError: false,
			checkFunc: func(t *testing.T, tag *n8nsdk.Tag, err error) {
				t.Helper()
				assert.NoError(t, err)
				assert.NotNil(t, tag)
				assert.Equal(t, "Production", tag.Name)
			},
		},
		{
			name: "read with nil timestamps",
			setupMock: func() *httptest.Server {
				id := "tag-456"
				tag := &n8nsdk.Tag{
					Id:        &id,
					Name:      "Development",
					CreatedAt: nil,
					UpdatedAt: nil,
				}

				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(tag)
				}))
			},
			tagID:     "tag-456",
			wantError: false,
			checkFunc: func(t *testing.T, tag *n8nsdk.Tag, err error) {
				t.Helper()
				assert.NoError(t, err)
				assert.NotNil(t, tag)
				assert.Nil(t, tag.CreatedAt)
				assert.Nil(t, tag.UpdatedAt)
			},
		},
		{
			name: "API error during read",
			setupMock: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte(`{"message": "Tag not found"}`))
				}))
			},
			tagID:     "nonexistent",
			wantError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := tt.setupMock()
			defer server.Close()

			apiClient := n8nsdk.NewAPIClient(&n8nsdk.Configuration{
				Servers: n8nsdk.ServerConfigurations{
					{URL: server.URL + "/api/v1"},
				},
			})

			r := &TagResource{
				client: &client.N8nClient{
					APIClient: apiClient,
				},
			}

			tag, httpResp, err := r.client.APIClient.TagsAPI.TagsIdGet(context.Background(), tt.tagID).Execute()
			if httpResp != nil && httpResp.Body != nil {
				defer httpResp.Body.Close()
			}

			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.checkFunc != nil {
					tt.checkFunc(t, tag, err)
				}
			}
		})
	}
}

func TestTagResource_UpdateAPICall(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setupMock func() *httptest.Server
		tagID     string
		newName   string
		wantError bool
		checkFunc func(*testing.T, *n8nsdk.Tag, error)
	}{
		{
			name: "successful update",
			setupMock: func() *httptest.Server {
				id := "tag-123"
				now := time.Now()
				tag := &n8nsdk.Tag{
					Id:        &id,
					Name:      "UpdatedName",
					CreatedAt: &now,
					UpdatedAt: &now,
				}

				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/api/v1/tags/tag-123" && r.Method == http.MethodPut {
						w.WriteHeader(http.StatusOK)
						json.NewEncoder(w).Encode(tag)
						return
					}
					w.WriteHeader(http.StatusNotFound)
				}))
			},
			tagID:     "tag-123",
			newName:   "UpdatedName",
			wantError: false,
			checkFunc: func(t *testing.T, tag *n8nsdk.Tag, err error) {
				t.Helper()
				assert.NoError(t, err)
				assert.NotNil(t, tag)
				assert.Equal(t, "UpdatedName", tag.Name)
			},
		},
		{
			name: "update with nil timestamps",
			setupMock: func() *httptest.Server {
				id := "tag-456"
				tag := &n8nsdk.Tag{
					Id:        &id,
					Name:      "NewName",
					CreatedAt: nil,
					UpdatedAt: nil,
				}

				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(tag)
				}))
			},
			tagID:     "tag-456",
			newName:   "NewName",
			wantError: false,
			checkFunc: func(t *testing.T, tag *n8nsdk.Tag, err error) {
				t.Helper()
				assert.NoError(t, err)
				assert.NotNil(t, tag)
				assert.Nil(t, tag.CreatedAt)
				assert.Nil(t, tag.UpdatedAt)
			},
		},
		{
			name: "API error during update",
			setupMock: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"message": "Update failed"}`))
				}))
			},
			tagID:     "tag-fail",
			newName:   "FailName",
			wantError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := tt.setupMock()
			defer server.Close()

			apiClient := n8nsdk.NewAPIClient(&n8nsdk.Configuration{
				Servers: n8nsdk.ServerConfigurations{
					{URL: server.URL + "/api/v1"},
				},
			})

			r := &TagResource{
				client: &client.N8nClient{
					APIClient: apiClient,
				},
			}

			tagRequest := n8nsdk.Tag{
				Name: tt.newName,
			}

			tag, httpResp, err := r.client.APIClient.TagsAPI.TagsIdPut(context.Background(), tt.tagID).
				Tag(tagRequest).
				Execute()
			if httpResp != nil && httpResp.Body != nil {
				defer httpResp.Body.Close()
			}

			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.checkFunc != nil {
					tt.checkFunc(t, tag, err)
				}
			}
		})
	}
}

func TestTagResource_DeleteAPICall(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setupMock func() *httptest.Server
		tagID     string
		wantError bool
	}{
		{
			name: "successful deletion",
			setupMock: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/api/v1/tags/tag-123" && r.Method == http.MethodDelete {
						w.WriteHeader(http.StatusNoContent)
						return
					}
					w.WriteHeader(http.StatusNotFound)
				}))
			},
			tagID:     "tag-123",
			wantError: false,
		},
		{
			name: "API error during deletion",
			setupMock: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"message": "Delete failed"}`))
				}))
			},
			tagID:     "tag-fail",
			wantError: true,
		},
		{
			name: "delete non-existent tag",
			setupMock: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte(`{"message": "Tag not found"}`))
				}))
			},
			tagID:     "nonexistent",
			wantError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := tt.setupMock()
			defer server.Close()

			apiClient := n8nsdk.NewAPIClient(&n8nsdk.Configuration{
				Servers: n8nsdk.ServerConfigurations{
					{URL: server.URL + "/api/v1"},
				},
			})

			r := &TagResource{
				client: &client.N8nClient{
					APIClient: apiClient,
				},
			}

			_, httpResp, err := r.client.APIClient.TagsAPI.TagsIdDelete(context.Background(), tt.tagID).Execute()
			if httpResp != nil && httpResp.Body != nil {
				defer httpResp.Body.Close()
			}

			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Test to verify the resource client field is accessible.
func TestTagResource_ClientField(t *testing.T) {
	t.Parallel()

	r := &TagResource{
		client: &client.N8nClient{},
	}

	assert.NotNil(t, r.client)
}
