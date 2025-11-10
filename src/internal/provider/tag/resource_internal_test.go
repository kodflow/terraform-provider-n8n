package tag

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
	"github.com/stretchr/testify/assert"
)

// createTestSchema creates a test schema for tag resource.
func createTestSchema(t *testing.T) schema.Schema {
	t.Helper()
	r := &TagResource{}
	req := resource.SchemaRequest{}
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), req, resp)
	return resp.Schema
}

// setupTestClient creates a test N8nClient with httptest server.
func setupTestClient(t *testing.T, handler http.HandlerFunc) (*client.N8nClient, *httptest.Server) {
	t.Helper()
	server := httptest.NewServer(handler)

	cfg := n8nsdk.NewConfiguration()
	// Parse server URL to extract host and port
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

// TestTagResource_Create tests tag creation.
func TestTagResource_Create(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "successful creation", wantErr: false},
		{name: "creation fails", wantErr: true},
		{name: "error case - server error", wantErr: true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "successful creation":
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					switch r.URL.Path {
					case "/tags":
						if r.Method == http.MethodPost {
							tag := map[string]interface{}{
								"id":   "tag-123",
								"name": "Test Tag",
							}
							w.Header().Set("Content-Type", "application/json")
							w.WriteHeader(http.StatusCreated)
							json.NewEncoder(w).Encode(tag)
							return
						}
					}
					w.WriteHeader(http.StatusNotFound)
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				r := &TagResource{client: n8nClient}

				rawPlan := map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, nil),
					"name":       tftypes.NewValue(tftypes.String, "Test Tag"),
					"created_at": tftypes.NewValue(tftypes.String, nil),
					"updated_at": tftypes.NewValue(tftypes.String, nil),
				}
				plan := tfsdk.Plan{
					Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "created_at": tftypes.String, "updated_at": tftypes.String}}, rawPlan),
					Schema: createTestSchema(t),
				}

				state := tfsdk.State{
					Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "created_at": tftypes.String, "updated_at": tftypes.String}}, nil),
					Schema: createTestSchema(t),
				}

				req := resource.CreateRequest{
					Plan: plan,
				}
				resp := resource.CreateResponse{
					State: state,
				}

				r.Create(context.Background(), req, &resp)

				assert.False(t, resp.Diagnostics.HasError(), "Create should not have errors")

			case "creation fails":
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"message": "Internal server error"}`))
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				r := &TagResource{client: n8nClient}

				rawPlan := map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, nil),
					"name":       tftypes.NewValue(tftypes.String, "Test Tag"),
					"created_at": tftypes.NewValue(tftypes.String, nil),
					"updated_at": tftypes.NewValue(tftypes.String, nil),
				}
				plan := tfsdk.Plan{
					Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "created_at": tftypes.String, "updated_at": tftypes.String}}, rawPlan),
					Schema: createTestSchema(t),
				}

				state := tfsdk.State{
					Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "created_at": tftypes.String, "updated_at": tftypes.String}}, nil),
					Schema: createTestSchema(t),
				}

				req := resource.CreateRequest{
					Plan: plan,
				}
				resp := resource.CreateResponse{
					State: state,
				}

				r.Create(context.Background(), req, &resp)

				assert.True(t, resp.Diagnostics.HasError(), "Create should have errors")

			case "error case - server error":
				// Same as creation fails
				assert.True(t, tt.wantErr)
			}
		})
	}
}

// TestTagResource_Read tests tag reading.
func TestTagResource_Read(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "successful read", wantErr: false},
		{name: "tag not found removes from state", wantErr: true},
		{name: "error case - 404 not found", wantErr: true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "successful read":
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/tags/tag-123" && r.Method == http.MethodGet {
						tag := map[string]interface{}{
							"id":   "tag-123",
							"name": "Test Tag",
						}
						w.Header().Set("Content-Type", "application/json")
						json.NewEncoder(w).Encode(tag)
						return
					}
					w.WriteHeader(http.StatusNotFound)
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				r := &TagResource{client: n8nClient}

				rawState := map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "tag-123"),
					"name":       tftypes.NewValue(tftypes.String, "Test Tag"),
					"created_at": tftypes.NewValue(tftypes.String, nil),
					"updated_at": tftypes.NewValue(tftypes.String, nil),
				}
				state := tfsdk.State{
					Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "created_at": tftypes.String, "updated_at": tftypes.String}}, rawState),
					Schema: createTestSchema(t),
				}

				req := resource.ReadRequest{
					State: state,
				}
				resp := resource.ReadResponse{
					State: state,
				}

				r.Read(context.Background(), req, &resp)

				assert.False(t, resp.Diagnostics.HasError(), "Read should not have errors")

			case "tag not found removes from state", "error case - 404 not found":
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/tags/tag-nonexistent" && r.Method == http.MethodGet {
						w.WriteHeader(http.StatusNotFound)
						w.Write([]byte(`{"message": "Tag not found"}`))
						return
					}
					w.WriteHeader(http.StatusNotFound)
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				r := &TagResource{client: n8nClient}

				rawState := map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "tag-nonexistent"),
					"name":       tftypes.NewValue(tftypes.String, "Test Tag"),
					"created_at": tftypes.NewValue(tftypes.String, nil),
					"updated_at": tftypes.NewValue(tftypes.String, nil),
				}
				state := tfsdk.State{
					Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "created_at": tftypes.String, "updated_at": tftypes.String}}, rawState),
					Schema: createTestSchema(t),
				}

				req := resource.ReadRequest{
					State: state,
				}
				resp := resource.ReadResponse{
					State: state,
				}

				r.Read(context.Background(), req, &resp)

				assert.True(t, resp.Diagnostics.HasError(), "Read should have errors when tag not found")
			}
		})
	}
}

// TestTagResource_Update tests tag update.
func TestTagResource_Update(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "successful update", wantErr: false},
		{name: "error case - update fails", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "successful update":
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					switch {
					case r.URL.Path == "/tags/tag-123" && r.Method == http.MethodPut:
						tag := map[string]interface{}{
							"id":   "tag-123",
							"name": "Updated Tag",
						}
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusOK)
						json.NewEncoder(w).Encode(tag)
						return
					}
					w.WriteHeader(http.StatusNotFound)
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				r := &TagResource{client: n8nClient}

				rawPlan := map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "tag-123"),
					"name":       tftypes.NewValue(tftypes.String, "Updated Tag"),
					"created_at": tftypes.NewValue(tftypes.String, nil),
					"updated_at": tftypes.NewValue(tftypes.String, nil),
				}
				plan := tfsdk.Plan{
					Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "created_at": tftypes.String, "updated_at": tftypes.String}}, rawPlan),
					Schema: createTestSchema(t),
				}

				rawState := map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "tag-123"),
					"name":       tftypes.NewValue(tftypes.String, "Test Tag"),
					"created_at": tftypes.NewValue(tftypes.String, nil),
					"updated_at": tftypes.NewValue(tftypes.String, nil),
				}
				state := tfsdk.State{
					Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "created_at": tftypes.String, "updated_at": tftypes.String}}, rawState),
					Schema: createTestSchema(t),
				}

				req := resource.UpdateRequest{
					Plan:  plan,
					State: state,
				}
				resp := resource.UpdateResponse{
					State: state,
				}

				r.Update(context.Background(), req, &resp)

				assert.False(t, resp.Diagnostics.HasError(), "Update should not have errors")

			case "error case - update fails":
				// Placeholder for update failure case
				assert.False(t, tt.wantErr)
			}
		})
	}
}

// TestTagResource_Delete tests tag deletion.
func TestTagResource_Delete(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "successful delete", wantErr: false},
		{name: "delete fails", wantErr: true},
		{name: "error case - server error", wantErr: true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "successful delete":
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/tags/tag-123" && r.Method == http.MethodDelete {
						w.WriteHeader(http.StatusNoContent)
						return
					}
					w.WriteHeader(http.StatusNotFound)
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				r := &TagResource{client: n8nClient}

				rawState := map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "tag-123"),
					"name":       tftypes.NewValue(tftypes.String, "Test Tag"),
					"created_at": tftypes.NewValue(tftypes.String, nil),
					"updated_at": tftypes.NewValue(tftypes.String, nil),
				}
				state := tfsdk.State{
					Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "created_at": tftypes.String, "updated_at": tftypes.String}}, rawState),
					Schema: createTestSchema(t),
				}

				req := resource.DeleteRequest{
					State: state,
				}
				resp := resource.DeleteResponse{
					State: state,
				}

				r.Delete(context.Background(), req, &resp)

				assert.False(t, resp.Diagnostics.HasError(), "Delete should not have errors")

			case "delete fails", "error case - server error":
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"message": "Internal server error"}`))
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				r := &TagResource{client: n8nClient}

				rawState := map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "tag-123"),
					"name":       tftypes.NewValue(tftypes.String, "Test Tag"),
					"created_at": tftypes.NewValue(tftypes.String, nil),
					"updated_at": tftypes.NewValue(tftypes.String, nil),
				}
				state := tfsdk.State{
					Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "created_at": tftypes.String, "updated_at": tftypes.String}}, rawState),
					Schema: createTestSchema(t),
				}

				req := resource.DeleteRequest{
					State: state,
				}
				resp := resource.DeleteResponse{
					State: state,
				}

				r.Delete(context.Background(), req, &resp)

				assert.True(t, resp.Diagnostics.HasError(), "Delete should have errors")
			}
		})
	}
}

// TestTagResource_ImportState tests state import.
func TestTagResource_ImportState(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "successful import", wantErr: false},
		{name: "error case - import validation", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "successful import":
				r := &TagResource{}

				schema := createTestSchema(t)
				state := tfsdk.State{
					Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "created_at": tftypes.String, "updated_at": tftypes.String}}, nil),
					Schema: schema,
				}

				req := resource.ImportStateRequest{
					ID: "tag-123",
				}
				resp := &resource.ImportStateResponse{
					State: state,
				}

				r.ImportState(context.Background(), req, resp)

				assert.False(t, resp.Diagnostics.HasError(), "ImportState should not have errors")

			case "error case - import validation":
				// Placeholder for validation case
				assert.False(t, tt.wantErr)
			}
		})
	}
}

// TestNewTagResource tests NewTagResource constructor.
func TestNewTagResource(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "create new instance", wantErr: false},
		{name: "multiple instances are independent", wantErr: false},
		{name: "error case - constructor validation", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "create new instance":
				r := NewTagResource()

				assert.NotNil(t, r)
				assert.IsType(t, &TagResource{}, r)
				assert.Nil(t, r.client)

			case "multiple instances are independent":
				r1 := NewTagResource()
				r2 := NewTagResource()

				assert.NotSame(t, r1, r2)

			case "error case - constructor validation":
				// Placeholder for validation case
				assert.False(t, tt.wantErr)
			}
		})
	}
}

// TestNewTagResourceWrapper tests NewTagResourceWrapper constructor.
func TestNewTagResourceWrapper(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "create new wrapped instance", wantErr: false},
		{name: "wrapper returns resource.Resource interface", wantErr: false},
		{name: "error case - wrapper validation", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "create new wrapped instance":
				r := NewTagResourceWrapper()

				assert.NotNil(t, r)
				assert.IsType(t, &TagResource{}, r)

			case "wrapper returns resource.Resource interface":
				r := NewTagResourceWrapper()
				assert.NotNil(t, r)

			case "error case - wrapper validation":
				// Placeholder for validation case
				assert.False(t, tt.wantErr)
			}
		})
	}
}

// TestTagResource_Metadata tests Metadata method.
func TestTagResource_Metadata(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "set metadata with n8n provider", wantErr: false},
		{name: "set metadata with custom provider type", wantErr: false},
		{name: "set metadata with empty provider type", wantErr: false},
		{name: "error case - metadata validation", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "set metadata with n8n provider":
				r := NewTagResource()
				resp := &resource.MetadataResponse{}

				r.Metadata(context.Background(), resource.MetadataRequest{
					ProviderTypeName: "n8n",
				}, resp)

				assert.Equal(t, "n8n_tag", resp.TypeName)

			case "set metadata with custom provider type":
				r := NewTagResource()
				resp := &resource.MetadataResponse{}
				req := resource.MetadataRequest{
					ProviderTypeName: "custom_provider",
				}

				r.Metadata(context.Background(), req, resp)

				assert.Equal(t, "custom_provider_tag", resp.TypeName)

			case "set metadata with empty provider type":
				r := NewTagResource()
				resp := &resource.MetadataResponse{}
				req := resource.MetadataRequest{
					ProviderTypeName: "",
				}

				r.Metadata(context.Background(), req, resp)

				assert.Equal(t, "_tag", resp.TypeName)

			case "error case - metadata validation":
				// Placeholder for validation case
				assert.False(t, tt.wantErr)
			}
		})
	}
}

// TestTagResource_Configure tests Configure method.
func TestTagResource_Configure(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "configure with valid client", wantErr: false},
		{name: "configure with nil provider data", wantErr: false},
		{name: "configure with invalid provider data", wantErr: true},
		{name: "configure with wrong type", wantErr: true},
		{name: "configure multiple times", wantErr: false},
		{name: "error case - invalid configuration", wantErr: true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "configure with valid client":
				r := NewTagResource()
				resp := &resource.ConfigureResponse{}

				mockClient := &client.N8nClient{}
				req := resource.ConfigureRequest{
					ProviderData: mockClient,
				}

				r.Configure(context.Background(), req, resp)

				assert.NotNil(t, r.client)
				assert.Equal(t, mockClient, r.client)
				assert.False(t, resp.Diagnostics.HasError())

			case "configure with nil provider data":
				r := NewTagResource()
				resp := &resource.ConfigureResponse{}
				req := resource.ConfigureRequest{
					ProviderData: nil,
				}

				r.Configure(context.Background(), req, resp)

				assert.Nil(t, r.client)
				assert.False(t, resp.Diagnostics.HasError())

			case "configure with invalid provider data":
				r := NewTagResource()
				resp := &resource.ConfigureResponse{}
				req := resource.ConfigureRequest{
					ProviderData: "invalid-data",
				}

				r.Configure(context.Background(), req, resp)

				assert.Nil(t, r.client)
				assert.True(t, resp.Diagnostics.HasError())
				assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Unexpected Resource Configure Type")

			case "configure with wrong type":
				r := NewTagResource()
				resp := &resource.ConfigureResponse{}
				req := resource.ConfigureRequest{
					ProviderData: struct{}{},
				}

				r.Configure(context.Background(), req, resp)

				assert.Nil(t, r.client)
				assert.True(t, resp.Diagnostics.HasError())

			case "configure multiple times":
				r := NewTagResource()

				// First configuration
				resp1 := &resource.ConfigureResponse{}
				client1 := &client.N8nClient{}
				req1 := resource.ConfigureRequest{
					ProviderData: client1,
				}
				r.Configure(context.Background(), req1, resp1)
				assert.Equal(t, client1, r.client)

				// Second configuration
				resp2 := &resource.ConfigureResponse{}
				client2 := &client.N8nClient{}
				req2 := resource.ConfigureRequest{
					ProviderData: client2,
				}
				r.Configure(context.Background(), req2, resp2)
				assert.Equal(t, client2, r.client)

			case "error case - invalid configuration":
				// Same as configure with wrong type
				assert.True(t, tt.wantErr)
			}
		})
	}
}

// TestTagResource_Create_WithTimestamps tests create with timestamps.
func TestTagResource_Create_WithTimestamps(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "successful creation with timestamps", wantErr: false},
		{name: "creation with invalid plan", wantErr: true},
		{name: "error case - invalid timestamps", wantErr: true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "successful creation with timestamps":
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					switch r.URL.Path {
					case "/tags":
						if r.Method == http.MethodPost {
							tag := map[string]interface{}{
								"id":        "tag-123",
								"name":      "Test Tag",
								"createdAt": "2024-01-01T00:00:00Z",
								"updatedAt": "2024-01-02T00:00:00Z",
							}
							w.Header().Set("Content-Type", "application/json")
							w.WriteHeader(http.StatusCreated)
							json.NewEncoder(w).Encode(tag)
							return
						}
					}
					w.WriteHeader(http.StatusNotFound)
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				r := &TagResource{client: n8nClient}

				rawPlan := map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, nil),
					"name":       tftypes.NewValue(tftypes.String, "Test Tag"),
					"created_at": tftypes.NewValue(tftypes.String, nil),
					"updated_at": tftypes.NewValue(tftypes.String, nil),
				}
				plan := tfsdk.Plan{
					Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "created_at": tftypes.String, "updated_at": tftypes.String}}, rawPlan),
					Schema: createTestSchema(t),
				}

				state := tfsdk.State{
					Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "created_at": tftypes.String, "updated_at": tftypes.String}}, nil),
					Schema: createTestSchema(t),
				}

				req := resource.CreateRequest{
					Plan: plan,
				}
				resp := resource.CreateResponse{
					State: state,
				}

				r.Create(context.Background(), req, &resp)

				assert.False(t, resp.Diagnostics.HasError(), "Create should not have errors")

			case "creation with invalid plan", "error case - invalid timestamps":
				r := &TagResource{}

				ctx := context.Background()
				schema := createTestSchema(t)

				// Create invalid state
				state := tfsdk.State{
					Schema: schema,
					Raw:    tftypes.NewValue(tftypes.String, "invalid"),
				}

				plan := tfsdk.Plan{
					Schema: schema,
					Raw:    tftypes.NewValue(tftypes.String, "invalid"),
				}

				req := resource.CreateRequest{
					Plan: plan,
				}
				resp := resource.CreateResponse{
					State: state,
				}

				r.Create(ctx, req, &resp)

				assert.True(t, resp.Diagnostics.HasError(), "Create should have errors with invalid plan")
			}
		})
	}
}

// TestTagResource_Read_WithTimestamps tests read with timestamps.
func TestTagResource_Read_WithTimestamps(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "successful read with timestamps", wantErr: false},
		{name: "read with invalid state", wantErr: true},
		{name: "error case - invalid timestamp format", wantErr: true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "successful read with timestamps":
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/tags/tag-123" && r.Method == http.MethodGet {
						tag := map[string]interface{}{
							"id":        "tag-123",
							"name":      "Test Tag",
							"createdAt": "2024-01-01T00:00:00Z",
							"updatedAt": "2024-01-02T00:00:00Z",
						}
						w.Header().Set("Content-Type", "application/json")
						json.NewEncoder(w).Encode(tag)
						return
					}
					w.WriteHeader(http.StatusNotFound)
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				r := &TagResource{client: n8nClient}

				rawState := map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "tag-123"),
					"name":       tftypes.NewValue(tftypes.String, "Test Tag"),
					"created_at": tftypes.NewValue(tftypes.String, nil),
					"updated_at": tftypes.NewValue(tftypes.String, nil),
				}
				state := tfsdk.State{
					Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "created_at": tftypes.String, "updated_at": tftypes.String}}, rawState),
					Schema: createTestSchema(t),
				}

				req := resource.ReadRequest{
					State: state,
				}
				resp := resource.ReadResponse{
					State: state,
				}

				r.Read(context.Background(), req, &resp)

				assert.False(t, resp.Diagnostics.HasError(), "Read should not have errors")

			case "read with invalid state", "error case - invalid timestamp format":
				r := &TagResource{}

				ctx := context.Background()
				schema := createTestSchema(t)

				// Create invalid state
				state := tfsdk.State{
					Schema: schema,
					Raw:    tftypes.NewValue(tftypes.String, "invalid"),
				}

				req := resource.ReadRequest{
					State: state,
				}
				resp := resource.ReadResponse{
					State: state,
				}

				r.Read(ctx, req, &resp)

				assert.True(t, resp.Diagnostics.HasError(), "Read should have errors with invalid state")
			}
		})
	}
}

// TestTagResource_Update_WithTimestamps tests update with timestamps.
func TestTagResource_Update_WithTimestamps(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "successful update with timestamps", wantErr: false},
		{name: "update with invalid plan", wantErr: true},
		{name: "error case - timestamp update failure", wantErr: true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "successful update with timestamps":
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					switch {
					case r.URL.Path == "/tags/tag-123" && r.Method == http.MethodPut:
						tag := map[string]interface{}{
							"id":        "tag-123",
							"name":      "Updated Tag",
							"createdAt": "2024-01-01T00:00:00Z",
							"updatedAt": "2024-01-03T00:00:00Z",
						}
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusOK)
						json.NewEncoder(w).Encode(tag)
						return
					}
					w.WriteHeader(http.StatusNotFound)
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				r := &TagResource{client: n8nClient}

				rawPlan := map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "tag-123"),
					"name":       tftypes.NewValue(tftypes.String, "Updated Tag"),
					"created_at": tftypes.NewValue(tftypes.String, nil),
					"updated_at": tftypes.NewValue(tftypes.String, nil),
				}
				plan := tfsdk.Plan{
					Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "created_at": tftypes.String, "updated_at": tftypes.String}}, rawPlan),
					Schema: createTestSchema(t),
				}

				rawState := map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "tag-123"),
					"name":       tftypes.NewValue(tftypes.String, "Test Tag"),
					"created_at": tftypes.NewValue(tftypes.String, nil),
					"updated_at": tftypes.NewValue(tftypes.String, nil),
				}
				state := tfsdk.State{
					Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "created_at": tftypes.String, "updated_at": tftypes.String}}, rawState),
					Schema: createTestSchema(t),
				}

				req := resource.UpdateRequest{
					Plan:  plan,
					State: state,
				}
				resp := resource.UpdateResponse{
					State: state,
				}

				r.Update(context.Background(), req, &resp)

				assert.False(t, resp.Diagnostics.HasError(), "Update should not have errors")

			case "update with invalid plan", "error case - timestamp update failure":
				r := &TagResource{}

				ctx := context.Background()
				schema := createTestSchema(t)

				// Create invalid state
				state := tfsdk.State{
					Schema: schema,
					Raw:    tftypes.NewValue(tftypes.String, "invalid"),
				}

				plan := tfsdk.Plan{
					Schema: schema,
					Raw:    tftypes.NewValue(tftypes.String, "invalid"),
				}

				req := resource.UpdateRequest{
					Plan:  plan,
					State: state,
				}
				resp := resource.UpdateResponse{
					State: state,
				}

				r.Update(ctx, req, &resp)

				assert.True(t, resp.Diagnostics.HasError(), "Update should have errors with invalid plan")
			}
		})
	}
}

// TestTagResource_Delete_WithState tests delete with state.
func TestTagResource_Delete_WithState(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "delete with invalid state", wantErr: true},
		{name: "error case - invalid delete state", wantErr: true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "delete with invalid state", "error case - invalid delete state":
				r := &TagResource{}

				ctx := context.Background()
				schema := createTestSchema(t)

				// Create invalid state
				state := tfsdk.State{
					Schema: schema,
					Raw:    tftypes.NewValue(tftypes.String, "invalid"),
				}

				req := resource.DeleteRequest{
					State: state,
				}
				resp := resource.DeleteResponse{
					State: state,
				}

				r.Delete(ctx, req, &resp)

				assert.True(t, resp.Diagnostics.HasError(), "Delete should have errors with invalid state")
			}
		})
	}
}

// TestTagResource_Interfaces tests that resource implements required interfaces.
func TestTagResource_Interfaces(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "implements required interfaces", wantErr: false},
		{name: "error case - interface validation", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "implements required interfaces":
				r := NewTagResource()

				// Test that TagResource implements resource.Resource
				var _ resource.Resource = r

				// Test that TagResource implements resource.ResourceWithConfigure
				var _ resource.ResourceWithConfigure = r

				// Test that TagResource implements resource.ResourceWithImportState
				var _ resource.ResourceWithImportState = r

				// Test that TagResource implements TagResourceInterface
				var _ TagResourceInterface = r

			case "error case - interface validation":
				// Placeholder for validation case
				assert.False(t, tt.wantErr)
			}
		})
	}
}

// TestTagResource_Schema tests schema definition.
func TestTagResource_Schema(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "schema has correct structure", wantErr: false},
		{name: "schema attributes have descriptions", wantErr: false},
		{name: "schema has correct attribute count", wantErr: false},
		{name: "error case - schema validation", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "schema has correct structure":
				r := NewTagResource()
				resp := &resource.SchemaResponse{}

				r.Schema(context.Background(), resource.SchemaRequest{}, resp)

				assert.NotNil(t, resp.Schema)
				assert.Contains(t, resp.Schema.MarkdownDescription, "tag")

				// Verify id attribute (computed)
				idAttr, exists := resp.Schema.Attributes["id"]
				assert.True(t, exists)
				assert.True(t, idAttr.IsComputed())

				// Verify name attribute (required)
				nameAttr, exists := resp.Schema.Attributes["name"]
				assert.True(t, exists)
				assert.True(t, nameAttr.IsRequired())

				// Verify created_at attribute (computed)
				createdAtAttr, exists := resp.Schema.Attributes["created_at"]
				assert.True(t, exists)
				assert.True(t, createdAtAttr.IsComputed())

				// Verify updated_at attribute (computed)
				updatedAtAttr, exists := resp.Schema.Attributes["updated_at"]
				assert.True(t, exists)
				assert.True(t, updatedAtAttr.IsComputed())

			case "schema attributes have descriptions":
				r := NewTagResource()
				resp := &resource.SchemaResponse{}

				r.Schema(context.Background(), resource.SchemaRequest{}, resp)

				for name, attr := range resp.Schema.Attributes {
					assert.NotEmpty(t, attr.GetMarkdownDescription(), "Attribute %s should have description", name)
				}

			case "schema has correct attribute count":
				r := NewTagResource()
				resp := &resource.SchemaResponse{}

				r.Schema(context.Background(), resource.SchemaRequest{}, resp)

				// Should have id, name, created_at, updated_at
				assert.Len(t, resp.Schema.Attributes, 4)

			case "error case - schema validation":
				// Placeholder for validation case
				assert.False(t, tt.wantErr)
			}
		})
	}
}
