package tag

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
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

// TestTagResource_Create_Internal tests the Create method with full execution.
func TestTagResource_Create_Internal(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "create with successful API call",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte(`{"id":"tag-123","name":"test-tag","createdAt":"2024-01-01T00:00:00Z","updatedAt":"2024-01-01T00:00:00Z"}`))
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				r := NewTagResource()
				r.Configure(context.Background(), resource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &resource.ConfigureResponse{})

				ctx := context.Background()
				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				// Build plan using tftypes
				planRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, nil),
					"name":       tftypes.NewValue(tftypes.String, "test-tag"),
					"created_at": tftypes.NewValue(tftypes.String, nil),
					"updated_at": tftypes.NewValue(tftypes.String, nil),
				})

				plan := tfsdk.Plan{
					Schema: schemaResp.Schema,
					Raw:    planRaw,
				}

				state := tfsdk.State{
					Schema: schemaResp.Schema,
				}

				req := resource.CreateRequest{
					Plan: plan,
				}
				resp := &resource.CreateResponse{
					State: state,
				}

				// Call Create
				r.Create(ctx, req, resp)

				// Verify success
				assert.False(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "error - create with invalid plan",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := NewTagResource()
				ctx := context.Background()

				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				// Create invalid plan that will fail Get()
				planRaw := tftypes.NewValue(tftypes.String, "invalid")

				plan := tfsdk.Plan{
					Schema: schemaResp.Schema,
					Raw:    planRaw,
				}

				req := resource.CreateRequest{
					Plan: plan,
				}
				resp := &resource.CreateResponse{}

				r.Create(ctx, req, resp)

				// Verify error
				assert.True(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "error - create with API error",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"message": "Internal server error"}`))
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				r := NewTagResource()
				r.Configure(context.Background(), resource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &resource.ConfigureResponse{})

				ctx := context.Background()
				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				planRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, nil),
					"name":       tftypes.NewValue(tftypes.String, "test-tag"),
					"created_at": tftypes.NewValue(tftypes.String, nil),
					"updated_at": tftypes.NewValue(tftypes.String, nil),
				})

				plan := tfsdk.Plan{
					Schema: schemaResp.Schema,
					Raw:    planRaw,
				}

				state := tfsdk.State{
					Schema: schemaResp.Schema,
				}

				req := resource.CreateRequest{
					Plan: plan,
				}
				resp := &resource.CreateResponse{
					State: state,
				}

				r.Create(ctx, req, resp)

				// Verify error
				assert.True(t, resp.Diagnostics.HasError())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}

// TestTagResource_Read_Internal tests the Read method with full execution.
func TestTagResource_Read_Internal(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "read with successful API call",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{"id":"tag-123","name":"test-tag","createdAt":"2024-01-01T00:00:00Z","updatedAt":"2024-01-01T00:00:00Z"}`))
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				r := NewTagResource()
				r.Configure(context.Background(), resource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &resource.ConfigureResponse{})

				ctx := context.Background()
				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				// Build state
				stateRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "tag-123"),
					"name":       tftypes.NewValue(tftypes.String, "test-tag"),
					"created_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
					"updated_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
				})

				state := tfsdk.State{
					Schema: schemaResp.Schema,
					Raw:    stateRaw,
				}

				req := resource.ReadRequest{
					State: state,
				}
				resp := &resource.ReadResponse{
					State: state,
				}

				// Call Read
				r.Read(ctx, req, resp)

				// Verify success
				assert.False(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "error - read with invalid state",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := NewTagResource()
				ctx := context.Background()

				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				// Create invalid state
				stateRaw := tftypes.NewValue(tftypes.String, "invalid")

				state := tfsdk.State{
					Schema: schemaResp.Schema,
					Raw:    stateRaw,
				}

				req := resource.ReadRequest{
					State: state,
				}
				resp := &resource.ReadResponse{}

				r.Read(ctx, req, resp)

				// Verify error
				assert.True(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "error - read with API error",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte(`{"message": "Tag not found"}`))
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				r := NewTagResource()
				r.Configure(context.Background(), resource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &resource.ConfigureResponse{})

				ctx := context.Background()
				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				stateRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "tag-123"),
					"name":       tftypes.NewValue(tftypes.String, "test-tag"),
					"created_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
					"updated_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
				})

				state := tfsdk.State{
					Schema: schemaResp.Schema,
					Raw:    stateRaw,
				}

				req := resource.ReadRequest{
					State: state,
				}
				resp := &resource.ReadResponse{
					State: state,
				}

				r.Read(ctx, req, resp)

				// Verify error
				assert.True(t, resp.Diagnostics.HasError())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}

// TestTagResource_Update_Internal tests the Update method with full execution.
func TestTagResource_Update_Internal(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "update with successful API call",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{"id":"tag-123","name":"updated-tag","createdAt":"2024-01-01T00:00:00Z","updatedAt":"2024-01-02T00:00:00Z"}`))
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				r := NewTagResource()
				r.Configure(context.Background(), resource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &resource.ConfigureResponse{})

				ctx := context.Background()
				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				// Build plan
				planRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "tag-123"),
					"name":       tftypes.NewValue(tftypes.String, "updated-tag"),
					"created_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
					"updated_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
				})

				plan := tfsdk.Plan{
					Schema: schemaResp.Schema,
					Raw:    planRaw,
				}

				// Build state
				stateRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "tag-123"),
					"name":       tftypes.NewValue(tftypes.String, "old-tag"),
					"created_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
					"updated_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
				})

				state := tfsdk.State{
					Schema: schemaResp.Schema,
					Raw:    stateRaw,
				}

				req := resource.UpdateRequest{
					Plan:  plan,
					State: state,
				}
				resp := &resource.UpdateResponse{
					State: state,
				}

				// Call Update
				r.Update(ctx, req, resp)

				// Verify success
				assert.False(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "error - update with invalid plan",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := NewTagResource()
				ctx := context.Background()

				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				// Create invalid plan
				planRaw := tftypes.NewValue(tftypes.String, "invalid")

				plan := tfsdk.Plan{
					Schema: schemaResp.Schema,
					Raw:    planRaw,
				}

				req := resource.UpdateRequest{
					Plan: plan,
				}
				resp := &resource.UpdateResponse{}

				r.Update(ctx, req, resp)

				// Verify error
				assert.True(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "error - update with API error",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"message": "Internal server error"}`))
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				r := NewTagResource()
				r.Configure(context.Background(), resource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &resource.ConfigureResponse{})

				ctx := context.Background()
				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				planRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "tag-123"),
					"name":       tftypes.NewValue(tftypes.String, "updated-tag"),
					"created_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
					"updated_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
				})

				plan := tfsdk.Plan{
					Schema: schemaResp.Schema,
					Raw:    planRaw,
				}

				stateRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "tag-123"),
					"name":       tftypes.NewValue(tftypes.String, "old-tag"),
					"created_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
					"updated_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
				})

				state := tfsdk.State{
					Schema: schemaResp.Schema,
					Raw:    stateRaw,
				}

				req := resource.UpdateRequest{
					Plan:  plan,
					State: state,
				}
				resp := &resource.UpdateResponse{
					State: state,
				}

				r.Update(ctx, req, resp)

				// Verify error
				assert.True(t, resp.Diagnostics.HasError())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}

// TestTagResource_Delete_Internal tests the Delete method with full execution.
func TestTagResource_Delete_Internal(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "delete with successful API call",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNoContent)
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				r := NewTagResource()
				r.Configure(context.Background(), resource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &resource.ConfigureResponse{})

				ctx := context.Background()
				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				// Build state
				stateRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "tag-123"),
					"name":       tftypes.NewValue(tftypes.String, "test-tag"),
					"created_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
					"updated_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
				})

				state := tfsdk.State{
					Schema: schemaResp.Schema,
					Raw:    stateRaw,
				}

				req := resource.DeleteRequest{
					State: state,
				}
				resp := &resource.DeleteResponse{}

				// Call Delete
				r.Delete(ctx, req, resp)

				// Verify success
				assert.False(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "error - delete with invalid state",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := NewTagResource()
				ctx := context.Background()

				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				// Create invalid state
				stateRaw := tftypes.NewValue(tftypes.String, "invalid")

				state := tfsdk.State{
					Schema: schemaResp.Schema,
					Raw:    stateRaw,
				}

				req := resource.DeleteRequest{
					State: state,
				}
				resp := &resource.DeleteResponse{}

				r.Delete(ctx, req, resp)

				// Verify error
				assert.True(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "error - delete with API error",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"message": "Internal server error"}`))
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				r := NewTagResource()
				r.Configure(context.Background(), resource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &resource.ConfigureResponse{})

				ctx := context.Background()
				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				stateRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "tag-123"),
					"name":       tftypes.NewValue(tftypes.String, "test-tag"),
					"created_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
					"updated_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
				})

				state := tfsdk.State{
					Schema: schemaResp.Schema,
					Raw:    stateRaw,
				}

				req := resource.DeleteRequest{
					State: state,
				}
				resp := &resource.DeleteResponse{}

				r.Delete(ctx, req, resp)

				// Verify error
				assert.True(t, resp.Diagnostics.HasError())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}

// TestTagResource_ImportState_Internal tests the ImportState method.
func TestTagResource_ImportState_Internal(t *testing.T) {
	t.Run("import state passthrough", func(t *testing.T) {
		t.Helper()
		r := NewTagResource()
		ctx := context.Background()

		schemaResp := resource.SchemaResponse{}
		r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

		// Build empty state
		emptyValue := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, nil),
			"name":       tftypes.NewValue(tftypes.String, nil),
			"created_at": tftypes.NewValue(tftypes.String, nil),
			"updated_at": tftypes.NewValue(tftypes.String, nil),
		})

		req := resource.ImportStateRequest{
			ID: "tag-123",
		}
		resp := &resource.ImportStateResponse{
			State: tfsdk.State{
				Schema: schemaResp.Schema,
				Raw:    emptyValue,
			},
		}

		// Call ImportState
		r.ImportState(ctx, req, resp)

		// Verify no errors (ImportStatePassthroughID doesn't return errors for valid IDs)
		assert.False(t, resp.Diagnostics.HasError())
	})
}
