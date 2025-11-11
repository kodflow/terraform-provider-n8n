package tag

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
	"github.com/kodflow/n8n/src/internal/provider/tag/models"
	"github.com/stretchr/testify/assert"
)

// KTN-TEST-009 RATIONALE: Tests in this file test public functions that require
// access to private fields (specifically the .client field). These tests MUST
// remain in the internal test file to access the unexported client field for
// proper white-box testing.

func TestTagsDataSource_ReadSuccess(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setupMock func() *httptest.Server
		checkFunc func(*testing.T, *models.DataSources, error)
	}{
		{
			name: "successful read with multiple tags",
			setupMock: func() *httptest.Server {
				id1 := "tag-1"
				id2 := "tag-2"
				id3 := "tag-3"
				now := time.Now()
				tags := []n8nsdk.Tag{
					{
						Id:        &id1,
						Name:      "Production",
						CreatedAt: &now,
						UpdatedAt: &now,
					},
					{
						Id:        &id2,
						Name:      "Development",
						CreatedAt: &now,
						UpdatedAt: &now,
					},
					{
						Id:        &id3,
						Name:      "Testing",
						CreatedAt: &now,
						UpdatedAt: &now,
					},
				}

				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/api/v1/tags" && r.Method == http.MethodGet {
						response := n8nsdk.TagList{
							Data: tags,
						}
						w.WriteHeader(http.StatusOK)
						json.NewEncoder(w).Encode(response)
						return
					}
					w.WriteHeader(http.StatusNotFound)
				}))
			},
			checkFunc: func(t *testing.T, data *models.DataSources, err error) {
				t.Helper()
				assert.NoError(t, err)
				assert.NotNil(t, data)
				assert.Len(t, data.Tags, 3)
				assert.Equal(t, "tag-1", data.Tags[0].ID.ValueString())
				assert.Equal(t, "Production", data.Tags[0].Name.ValueString())
			},
		},
		{
			name: "successful read with empty tag list",
			setupMock: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					response := n8nsdk.TagList{
						Data: []n8nsdk.Tag{},
					}
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(response)
				}))
			},
			checkFunc: func(t *testing.T, data *models.DataSources, err error) {
				t.Helper()
				assert.NoError(t, err)
				assert.NotNil(t, data)
				assert.Len(t, data.Tags, 0)
			},
		},
		{
			name: "successful read with nil data",
			setupMock: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					response := n8nsdk.TagList{
						Data: nil,
					}
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(response)
				}))
			},
			checkFunc: func(t *testing.T, data *models.DataSources, err error) {
				t.Helper()
				assert.NoError(t, err)
				assert.NotNil(t, data)
				assert.Len(t, data.Tags, 0)
			},
		},
		{
			name: "tags with nil id",
			setupMock: func() *httptest.Server {
				now := time.Now()
				tags := []n8nsdk.Tag{
					{
						Id:        nil,
						Name:      "TagWithoutID",
						CreatedAt: &now,
						UpdatedAt: &now,
					},
				}

				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					response := n8nsdk.TagList{
						Data: tags,
					}
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(response)
				}))
			},
			checkFunc: func(t *testing.T, data *models.DataSources, err error) {
				t.Helper()
				assert.NoError(t, err)
				assert.Len(t, data.Tags, 1)
				assert.Equal(t, "TagWithoutID", data.Tags[0].Name.ValueString())
			},
		},
		{
			name: "tags with nil timestamps",
			setupMock: func() *httptest.Server {
				id := "tag-1"
				tags := []n8nsdk.Tag{
					{
						Id:        &id,
						Name:      "TagWithoutTimestamps",
						CreatedAt: nil,
						UpdatedAt: nil,
					},
				}

				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					response := n8nsdk.TagList{
						Data: tags,
					}
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(response)
				}))
			},
			checkFunc: func(t *testing.T, data *models.DataSources, err error) {
				t.Helper()
				assert.NoError(t, err)
				assert.Len(t, data.Tags, 1)
				assert.Equal(t, "", data.Tags[0].CreatedAt.ValueString())
				assert.Equal(t, "", data.Tags[0].UpdatedAt.ValueString())
			},
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

			d := &TagsDataSource{
				client: &client.N8nClient{
					APIClient: apiClient,
				},
			}

			// Simulate the Read operation by calling the API directly
			tagList, httpResp, err := d.client.APIClient.TagsAPI.TagsGet(context.Background()).Execute()
			if httpResp != nil && httpResp.Body != nil {
				defer httpResp.Body.Close()
			}

			data := &models.DataSources{
				Tags: make([]models.Item, 0),
			}

			if err == nil && tagList != nil && tagList.Data != nil {
				for _, tag := range tagList.Data {
					item := models.Item{}
					if tag.Id != nil {
						item.ID = types.StringValue(*tag.Id)
					} else {
						item.ID = types.StringValue("")
					}
					item.Name = types.StringValue(tag.Name)
					if tag.CreatedAt != nil {
						item.CreatedAt = types.StringValue(tag.CreatedAt.String())
					} else {
						item.CreatedAt = types.StringValue("")
					}
					if tag.UpdatedAt != nil {
						item.UpdatedAt = types.StringValue(tag.UpdatedAt.String())
					} else {
						item.UpdatedAt = types.StringValue("")
					}
					data.Tags = append(data.Tags, item)
				}
			}

			if tt.checkFunc != nil {
				tt.checkFunc(t, data, err)
			}
		})
	}
}

func TestTagsDataSource_ReadError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setupMock func() *httptest.Server
		wantError bool
	}{
		{
			name: "API error during read",
			setupMock: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"message": "Internal server error"}`))
				}))
			},
			wantError: true,
		},
		{
			name: "unauthorized error",
			setupMock: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte(`{"message": "Unauthorized"}`))
				}))
			},
			wantError: true,
		},
		{
			name: "forbidden error",
			setupMock: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusForbidden)
					w.Write([]byte(`{"message": "Forbidden"}`))
				}))
			},
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

			d := &TagsDataSource{
				client: &client.N8nClient{
					APIClient: apiClient,
				},
			}

			_, httpResp, err := d.client.APIClient.TagsAPI.TagsGet(context.Background()).Execute()
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

func TestTagsDataSource_ReadHTTPResponseHandling(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setupMock func() *httptest.Server
	}{
		{
			name: "nil httpResp body",
			setupMock: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					response := n8nsdk.TagList{
						Data: []n8nsdk.Tag{},
					}
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(response)
				}))
			},
		},
		{
			name: "large number of tags",
			setupMock: func() *httptest.Server {
				tags := make([]n8nsdk.Tag, 100)
				now := time.Now()
				for i := 0; i < 100; i++ {
					id := string(rune('a' + i))
					tags[i] = n8nsdk.Tag{
						Id:        &id,
						Name:      "Tag-" + id,
						CreatedAt: &now,
						UpdatedAt: &now,
					}
				}

				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					response := n8nsdk.TagList{
						Data: tags,
					}
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(response)
				}))
			},
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

			d := &TagsDataSource{
				client: &client.N8nClient{
					APIClient: apiClient,
				},
			}

			_, httpResp, err := d.client.APIClient.TagsAPI.TagsGet(context.Background()).Execute()
			if httpResp != nil && httpResp.Body != nil {
				defer httpResp.Body.Close()
			}

			assert.NoError(t, err)
		})
	}
}

func TestTagsDataSource_ReadInterface(t *testing.T) {
	t.Parallel()

	// Test that the datasource implements the correct interface.
	var _ datasource.DataSource = &TagsDataSource{}
	var _ TagsDataSourceInterface = &TagsDataSource{}
}
