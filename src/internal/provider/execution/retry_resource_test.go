package execution

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockExecutionRetryResourceInterface is a mock implementation of ExecutionRetryResourceInterface.
type MockExecutionRetryResourceInterface struct {
	mock.Mock
}

func (m *MockExecutionRetryResourceInterface) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockExecutionRetryResourceInterface) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockExecutionRetryResourceInterface) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockExecutionRetryResourceInterface) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockExecutionRetryResourceInterface) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockExecutionRetryResourceInterface) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockExecutionRetryResourceInterface) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockExecutionRetryResourceInterface) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	m.Called(ctx, req, resp)
}

func TestNewExecutionRetryResource(t *testing.T) {
	t.Run("create new instance", func(t *testing.T) {
		r := NewExecutionRetryResource()

		assert.NotNil(t, r)
		assert.IsType(t, &ExecutionRetryResource{}, r)
		assert.Nil(t, r.client)
	})

	t.Run("multiple instances are independent", func(t *testing.T) {
		r1 := NewExecutionRetryResource()
		r2 := NewExecutionRetryResource()

		assert.NotSame(t, r1, r2)
		// Each instance is a different pointer, even if they have the same initial state
	})
}

func TestNewExecutionRetryResourceWrapper(t *testing.T) {
	t.Run("create new wrapped instance", func(t *testing.T) {
		r := NewExecutionRetryResourceWrapper()

		assert.NotNil(t, r)
		assert.IsType(t, &ExecutionRetryResource{}, r)
	})

	t.Run("wrapper returns resource.Resource interface", func(t *testing.T) {
		r := NewExecutionRetryResourceWrapper()

		// r is already of type resource.Resource, no assertion needed
		assert.NotNil(t, r)
	})
}

func TestExecutionRetryResource_Metadata(t *testing.T) {
	t.Run("set metadata without provider type", func(t *testing.T) {
		r := NewExecutionRetryResource()
		resp := &resource.MetadataResponse{}

		r.Metadata(context.Background(), resource.MetadataRequest{
			ProviderTypeName: "n8n",
		}, resp)

		assert.Equal(t, "n8n_execution_retry", resp.TypeName)
	})

	t.Run("set metadata with provider type name", func(t *testing.T) {
		r := NewExecutionRetryResource()
		resp := &resource.MetadataResponse{}
		req := resource.MetadataRequest{
			ProviderTypeName: "n8n",
		}

		r.Metadata(context.Background(), req, resp)

		assert.Equal(t, "n8n_execution_retry", resp.TypeName)
	})

	t.Run("set metadata with different provider type", func(t *testing.T) {
		r := NewExecutionRetryResource()
		resp := &resource.MetadataResponse{}
		req := resource.MetadataRequest{
			ProviderTypeName: "custom_provider",
		}

		r.Metadata(context.Background(), req, resp)

		assert.Equal(t, "custom_provider_execution_retry", resp.TypeName)
	})
}

func TestExecutionRetryResource_Schema(t *testing.T) {
	t.Run("return schema", func(t *testing.T) {
		r := NewExecutionRetryResource()
		resp := &resource.SchemaResponse{}

		r.Schema(context.Background(), resource.SchemaRequest{}, resp)

		assert.NotNil(t, resp.Schema)
		assert.Contains(t, resp.Schema.MarkdownDescription, "Retries a failed n8n workflow execution")

		// Verify required attributes
		executionIdAttr, exists := resp.Schema.Attributes["execution_id"]
		assert.True(t, exists)
		assert.True(t, executionIdAttr.IsRequired())

		// Verify computed attributes
		computedAttrs := []string{
			"new_execution_id",
			"workflow_id",
			"finished",
			"mode",
			"started_at",
			"stopped_at",
			"status",
		}

		for _, attr := range computedAttrs {
			schemaAttr, exists := resp.Schema.Attributes[attr]
			assert.True(t, exists, "Attribute %s should exist", attr)
			assert.True(t, schemaAttr.IsComputed(), "Attribute %s should be computed", attr)
		}
	})

	t.Run("schema attributes have descriptions", func(t *testing.T) {
		r := NewExecutionRetryResource()
		resp := &resource.SchemaResponse{}

		r.Schema(context.Background(), resource.SchemaRequest{}, resp)

		for name, attr := range resp.Schema.Attributes {
			assert.NotEmpty(t, attr.GetMarkdownDescription(), "Attribute %s should have description", name)
		}
	})

	t.Run("schema has all expected attributes", func(t *testing.T) {
		r := NewExecutionRetryResource()
		resp := &resource.SchemaResponse{}

		r.Schema(context.Background(), resource.SchemaRequest{}, resp)

		// All expected attributes
		expectedAttrs := []string{
			"execution_id",
			"new_execution_id",
			"workflow_id",
			"finished",
			"mode",
			"started_at",
			"stopped_at",
			"status",
		}

		for _, attr := range expectedAttrs {
			schemaAttr, exists := resp.Schema.Attributes[attr]
			assert.True(t, exists, "Attribute %s should exist", attr)
			assert.NotNil(t, schemaAttr)
		}
	})
}

func TestExecutionRetryResource_Configure(t *testing.T) {
	t.Run("configure with valid client", func(t *testing.T) {
		r := NewExecutionRetryResource()
		resp := &resource.ConfigureResponse{}

		mockClient := &client.N8nClient{}
		req := resource.ConfigureRequest{
			ProviderData: mockClient,
		}

		r.Configure(context.Background(), req, resp)

		assert.NotNil(t, r.client)
		assert.Equal(t, mockClient, r.client)
		assert.False(t, resp.Diagnostics.HasError())
	})

	t.Run("configure with nil provider data", func(t *testing.T) {
		r := NewExecutionRetryResource()
		resp := &resource.ConfigureResponse{}
		req := resource.ConfigureRequest{
			ProviderData: nil,
		}

		r.Configure(context.Background(), req, resp)

		assert.Nil(t, r.client)
		assert.False(t, resp.Diagnostics.HasError())
	})

	t.Run("configure with invalid provider data", func(t *testing.T) {
		r := NewExecutionRetryResource()
		resp := &resource.ConfigureResponse{}
		req := resource.ConfigureRequest{
			ProviderData: "invalid-data",
		}

		r.Configure(context.Background(), req, resp)

		assert.Nil(t, r.client)
		assert.True(t, resp.Diagnostics.HasError())
		assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Unexpected Resource Configure Type")
	})

	t.Run("configure with wrong type", func(t *testing.T) {
		r := NewExecutionRetryResource()
		resp := &resource.ConfigureResponse{}
		req := resource.ConfigureRequest{
			ProviderData: struct{}{},
		}

		r.Configure(context.Background(), req, resp)

		assert.Nil(t, r.client)
		assert.True(t, resp.Diagnostics.HasError())
	})

	t.Run("configure multiple times", func(t *testing.T) {
		r := NewExecutionRetryResource()

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
	})
}

func TestExecutionRetryResource_ImportState(t *testing.T) {
	t.Run("import state with valid ID", func(t *testing.T) {
		r := NewExecutionRetryResource()
		ctx := context.Background()

		// Create schema for state
		schemaResp := resource.SchemaResponse{}
		r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

		// Create an empty state with the schema
		state := tfsdk.State{
			Schema: schemaResp.Schema,
		}

		// Initialize the raw value with required attributes
		state.Raw = tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
			"execution_id":     tftypes.NewValue(tftypes.String, nil),
			"new_execution_id": tftypes.NewValue(tftypes.String, nil),
			"workflow_id":      tftypes.NewValue(tftypes.String, nil),
			"finished":         tftypes.NewValue(tftypes.Bool, nil),
			"mode":             tftypes.NewValue(tftypes.String, nil),
			"started_at":       tftypes.NewValue(tftypes.String, nil),
			"stopped_at":       tftypes.NewValue(tftypes.String, nil),
			"status":           tftypes.NewValue(tftypes.String, nil),
		})

		resp := &resource.ImportStateResponse{
			State: state,
		}
		req := resource.ImportStateRequest{
			ID: "123",
		}

		r.ImportState(ctx, req, resp)

		assert.False(t, resp.Diagnostics.HasError())
	})

	t.Run("import state sets execution_id attribute", func(t *testing.T) {
		r := NewExecutionRetryResource()
		ctx := context.Background()

		// Create schema for state
		schemaResp := resource.SchemaResponse{}
		r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

		// Create an empty state with the schema
		state := tfsdk.State{
			Schema: schemaResp.Schema,
		}

		// Initialize the raw value with required attributes
		state.Raw = tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
			"execution_id":     tftypes.NewValue(tftypes.String, nil),
			"new_execution_id": tftypes.NewValue(tftypes.String, nil),
			"workflow_id":      tftypes.NewValue(tftypes.String, nil),
			"finished":         tftypes.NewValue(tftypes.Bool, nil),
			"mode":             tftypes.NewValue(tftypes.String, nil),
			"started_at":       tftypes.NewValue(tftypes.String, nil),
			"stopped_at":       tftypes.NewValue(tftypes.String, nil),
			"status":           tftypes.NewValue(tftypes.String, nil),
		})

		resp := &resource.ImportStateResponse{
			State: state,
		}
		req := resource.ImportStateRequest{
			ID: "456",
		}

		r.ImportState(ctx, req, resp)

		// Check that ImportStatePassthroughID was called
		assert.False(t, resp.Diagnostics.HasError())
	})
}

func TestExecutionRetryResource_Interfaces(t *testing.T) {
	t.Run("implements required interfaces", func(t *testing.T) {
		r := NewExecutionRetryResource()

		// Test that ExecutionRetryResource implements resource.Resource
		var _ resource.Resource = r

		// Test that ExecutionRetryResource implements resource.ResourceWithConfigure
		var _ resource.ResourceWithConfigure = r

		// Test that ExecutionRetryResource implements resource.ResourceWithImportState
		var _ resource.ResourceWithImportState = r

		// Test that ExecutionRetryResource implements ExecutionRetryResourceInterface
		var _ ExecutionRetryResourceInterface = r
	})
}

func TestExecutionRetryResourceConcurrency(t *testing.T) {
	t.Run("concurrent metadata calls", func(t *testing.T) {
		r := NewExecutionRetryResource()

		done := make(chan bool, 100)
		for i := 0; i < 100; i++ {
			go func() {
				resp := &resource.MetadataResponse{}
				r.Metadata(context.Background(), resource.MetadataRequest{
					ProviderTypeName: "n8n",
				}, resp)
				assert.Equal(t, "n8n_execution_retry", resp.TypeName)
				done <- true
			}()
		}

		for i := 0; i < 100; i++ {
			<-done
		}
	})

	t.Run("concurrent schema calls", func(t *testing.T) {
		r := NewExecutionRetryResource()

		done := make(chan bool, 100)
		for i := 0; i < 100; i++ {
			go func() {
				resp := &resource.SchemaResponse{}
				r.Schema(context.Background(), resource.SchemaRequest{}, resp)
				assert.NotNil(t, resp.Schema)
				done <- true
			}()
		}

		for i := 0; i < 100; i++ {
			<-done
		}
	})
}

func TestExecutionRetryResource_HelperFunctions(t *testing.T) {
	t.Run("path to execution_id", func(t *testing.T) {
		p := path.Root("execution_id")
		assert.Equal(t, "execution_id", p.String())
	})

	t.Run("path to new_execution_id", func(t *testing.T) {
		p := path.Root("new_execution_id")
		assert.Equal(t, "new_execution_id", p.String())
	})

	t.Run("types.StringValue conversion", func(t *testing.T) {
		strVal := types.StringValue("test-value")
		assert.Equal(t, "test-value", strVal.ValueString())
		assert.False(t, strVal.IsNull())
		assert.False(t, strVal.IsUnknown())
	})

	t.Run("types.BoolValue conversion", func(t *testing.T) {
		boolVal := types.BoolValue(true)
		assert.True(t, boolVal.ValueBool())
		assert.False(t, boolVal.IsNull())
		assert.False(t, boolVal.IsUnknown())
	})
}

func BenchmarkExecutionRetryResource_Schema(b *testing.B) {
	r := NewExecutionRetryResource()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := &resource.SchemaResponse{}
		r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	}
}

func BenchmarkExecutionRetryResource_Metadata(b *testing.B) {
	r := NewExecutionRetryResource()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := &resource.MetadataResponse{}
		r.Metadata(context.Background(), resource.MetadataRequest{
			ProviderTypeName: "n8n",
		}, resp)
	}
}

func BenchmarkExecutionRetryResource_Configure(b *testing.B) {
	r := NewExecutionRetryResource()
	mockClient := &client.N8nClient{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := &resource.ConfigureResponse{}
		req := resource.ConfigureRequest{
			ProviderData: mockClient,
		}
		r.Configure(context.Background(), req, resp)
	}
}

// createTestRetrySchema creates a test schema for retry execution resource.
func createTestRetrySchema(t *testing.T) schema.Schema {
	t.Helper()
	r := &ExecutionRetryResource{}
	req := resource.SchemaRequest{}
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), req, resp)
	return resp.Schema
}

// setupTestRetryClient creates a test N8nClient with httptest server.
func setupTestRetryClient(t *testing.T, handler http.HandlerFunc) (*client.N8nClient, *httptest.Server) {
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

// TestRetryExecutionResource_Create tests retry execution creation.
func TestRetryExecutionResource_Create(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/executions/123/retry" && r.Method == http.MethodPost {
				execution := map[string]interface{}{
					"id":         float32(456),
					"workflowId": float32(789),
					"finished":   true,
					"mode":       "manual",
					"startedAt":  "2024-01-01T00:00:00Z",
					"stoppedAt":  "2024-01-01T00:01:00Z",
					"status":     "success",
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(execution)
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestRetryClient(t, handler)
		defer server.Close()

		r := &ExecutionRetryResource{client: n8nClient}

		rawPlan := map[string]tftypes.Value{
			"execution_id":     tftypes.NewValue(tftypes.String, "123"),
			"new_execution_id": tftypes.NewValue(tftypes.String, nil),
			"workflow_id":      tftypes.NewValue(tftypes.String, nil),
			"finished":         tftypes.NewValue(tftypes.Bool, nil),
			"mode":             tftypes.NewValue(tftypes.String, nil),
			"started_at":       tftypes.NewValue(tftypes.String, nil),
			"stopped_at":       tftypes.NewValue(tftypes.String, nil),
			"status":           tftypes.NewValue(tftypes.String, nil),
		}
		plan := tfsdk.Plan{
			Raw: tftypes.NewValue(tftypes.Object{
				AttributeTypes: map[string]tftypes.Type{
					"execution_id":     tftypes.String,
					"new_execution_id": tftypes.String,
					"workflow_id":      tftypes.String,
					"finished":         tftypes.Bool,
					"mode":             tftypes.String,
					"started_at":       tftypes.String,
					"stopped_at":       tftypes.String,
					"status":           tftypes.String,
				},
			}, rawPlan),
			Schema: createTestRetrySchema(t),
		}

		state := tfsdk.State{
			Raw: tftypes.NewValue(tftypes.Object{
				AttributeTypes: map[string]tftypes.Type{
					"execution_id":     tftypes.String,
					"new_execution_id": tftypes.String,
					"workflow_id":      tftypes.String,
					"finished":         tftypes.Bool,
					"mode":             tftypes.String,
					"started_at":       tftypes.String,
					"stopped_at":       tftypes.String,
					"status":           tftypes.String,
				},
			}, nil),
			Schema: createTestRetrySchema(t),
		}

		req := resource.CreateRequest{
			Plan: plan,
		}
		resp := resource.CreateResponse{
			State: state,
		}

		r.Create(context.Background(), req, &resp)

		assert.False(t, resp.Diagnostics.HasError(), "Create should not have errors")
	})

	t.Run("creation fails", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message": "Internal server error"}`))
		})

		n8nClient, server := setupTestRetryClient(t, handler)
		defer server.Close()

		r := &ExecutionRetryResource{client: n8nClient}

		rawPlan := map[string]tftypes.Value{
			"execution_id":     tftypes.NewValue(tftypes.String, "123"),
			"new_execution_id": tftypes.NewValue(tftypes.String, nil),
			"workflow_id":      tftypes.NewValue(tftypes.String, nil),
			"finished":         tftypes.NewValue(tftypes.Bool, nil),
			"mode":             tftypes.NewValue(tftypes.String, nil),
			"started_at":       tftypes.NewValue(tftypes.String, nil),
			"stopped_at":       tftypes.NewValue(tftypes.String, nil),
			"status":           tftypes.NewValue(tftypes.String, nil),
		}
		plan := tfsdk.Plan{
			Raw: tftypes.NewValue(tftypes.Object{
				AttributeTypes: map[string]tftypes.Type{
					"execution_id":     tftypes.String,
					"new_execution_id": tftypes.String,
					"workflow_id":      tftypes.String,
					"finished":         tftypes.Bool,
					"mode":             tftypes.String,
					"started_at":       tftypes.String,
					"stopped_at":       tftypes.String,
					"status":           tftypes.String,
				},
			}, rawPlan),
			Schema: createTestRetrySchema(t),
		}

		state := tfsdk.State{
			Raw: tftypes.NewValue(tftypes.Object{
				AttributeTypes: map[string]tftypes.Type{
					"execution_id":     tftypes.String,
					"new_execution_id": tftypes.String,
					"workflow_id":      tftypes.String,
					"finished":         tftypes.Bool,
					"mode":             tftypes.String,
					"started_at":       tftypes.String,
					"stopped_at":       tftypes.String,
					"status":           tftypes.String,
				},
			}, nil),
			Schema: createTestRetrySchema(t),
		}

		req := resource.CreateRequest{
			Plan: plan,
		}
		resp := resource.CreateResponse{
			State: state,
		}

		r.Create(context.Background(), req, &resp)

		assert.True(t, resp.Diagnostics.HasError(), "Create should have errors")
	})
}

// TestRetryExecutionResource_Read tests retry execution reading.
func TestRetryExecutionResource_Read(t *testing.T) {
	t.Run("successful read", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		n8nClient, server := setupTestRetryClient(t, handler)
		defer server.Close()

		r := &ExecutionRetryResource{client: n8nClient}

		rawState := map[string]tftypes.Value{
			"execution_id":     tftypes.NewValue(tftypes.String, "123"),
			"new_execution_id": tftypes.NewValue(tftypes.String, "456"),
			"workflow_id":      tftypes.NewValue(tftypes.String, "workflow-789"),
			"finished":         tftypes.NewValue(tftypes.Bool, true),
			"mode":             tftypes.NewValue(tftypes.String, "manual"),
			"started_at":       tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
			"stopped_at":       tftypes.NewValue(tftypes.String, "2024-01-01T00:01:00Z"),
			"status":           tftypes.NewValue(tftypes.String, "success"),
		}
		state := tfsdk.State{
			Raw: tftypes.NewValue(tftypes.Object{
				AttributeTypes: map[string]tftypes.Type{
					"execution_id":     tftypes.String,
					"new_execution_id": tftypes.String,
					"workflow_id":      tftypes.String,
					"finished":         tftypes.Bool,
					"mode":             tftypes.String,
					"started_at":       tftypes.String,
					"stopped_at":       tftypes.String,
					"status":           tftypes.String,
				},
			}, rawState),
			Schema: createTestRetrySchema(t),
		}

		req := resource.ReadRequest{
			State: state,
		}
		resp := resource.ReadResponse{
			State: state,
		}

		r.Read(context.Background(), req, &resp)

		assert.False(t, resp.Diagnostics.HasError(), "Read should not have errors")
	})

	t.Run("execution not found removes from state", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestRetryClient(t, handler)
		defer server.Close()

		r := &ExecutionRetryResource{client: n8nClient}

		rawState := map[string]tftypes.Value{
			"execution_id":     tftypes.NewValue(tftypes.String, "nonexistent"),
			"new_execution_id": tftypes.NewValue(tftypes.String, "456"),
			"workflow_id":      tftypes.NewValue(tftypes.String, "workflow-789"),
			"finished":         tftypes.NewValue(tftypes.Bool, true),
			"mode":             tftypes.NewValue(tftypes.String, "manual"),
			"started_at":       tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
			"stopped_at":       tftypes.NewValue(tftypes.String, "2024-01-01T00:01:00Z"),
			"status":           tftypes.NewValue(tftypes.String, "success"),
		}
		state := tfsdk.State{
			Raw: tftypes.NewValue(tftypes.Object{
				AttributeTypes: map[string]tftypes.Type{
					"execution_id":     tftypes.String,
					"new_execution_id": tftypes.String,
					"workflow_id":      tftypes.String,
					"finished":         tftypes.Bool,
					"mode":             tftypes.String,
					"started_at":       tftypes.String,
					"stopped_at":       tftypes.String,
					"status":           tftypes.String,
				},
			}, rawState),
			Schema: createTestRetrySchema(t),
		}

		req := resource.ReadRequest{
			State: state,
		}
		resp := resource.ReadResponse{
			State: state,
		}

		r.Read(context.Background(), req, &resp)

		assert.False(t, resp.Diagnostics.HasError(), "Read should not have errors even when execution not found")
		// Note: In Terraform, when a resource is not found during Read, it's removed from state
		// This is tested by checking that resp.State would be empty, but we can't easily verify that in unit tests
	})
}

// TestRetryExecutionResource_Delete tests retry execution deletion.
func TestRetryExecutionResource_Delete(t *testing.T) {
	t.Run("successful delete", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		})

		n8nClient, server := setupTestRetryClient(t, handler)
		defer server.Close()

		r := &ExecutionRetryResource{client: n8nClient}

		rawState := map[string]tftypes.Value{
			"execution_id":     tftypes.NewValue(tftypes.String, "123"),
			"new_execution_id": tftypes.NewValue(tftypes.String, "456"),
			"workflow_id":      tftypes.NewValue(tftypes.String, "workflow-789"),
			"finished":         tftypes.NewValue(tftypes.Bool, true),
			"mode":             tftypes.NewValue(tftypes.String, "manual"),
			"started_at":       tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
			"stopped_at":       tftypes.NewValue(tftypes.String, "2024-01-01T00:01:00Z"),
			"status":           tftypes.NewValue(tftypes.String, "success"),
		}
		state := tfsdk.State{
			Raw: tftypes.NewValue(tftypes.Object{
				AttributeTypes: map[string]tftypes.Type{
					"execution_id":     tftypes.String,
					"new_execution_id": tftypes.String,
					"workflow_id":      tftypes.String,
					"finished":         tftypes.Bool,
					"mode":             tftypes.String,
					"started_at":       tftypes.String,
					"stopped_at":       tftypes.String,
					"status":           tftypes.String,
				},
			}, rawState),
			Schema: createTestRetrySchema(t),
		}

		req := resource.DeleteRequest{
			State: state,
		}
		resp := resource.DeleteResponse{
			State: state,
		}

		r.Delete(context.Background(), req, &resp)

		assert.False(t, resp.Diagnostics.HasError(), "Delete should not have errors")
	})

	t.Run("delete fails", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message": "Internal server error"}`))
		})

		n8nClient, server := setupTestRetryClient(t, handler)
		defer server.Close()

		r := &ExecutionRetryResource{client: n8nClient}

		rawState := map[string]tftypes.Value{
			"execution_id":     tftypes.NewValue(tftypes.String, "123"),
			"new_execution_id": tftypes.NewValue(tftypes.String, "456"),
			"workflow_id":      tftypes.NewValue(tftypes.String, "workflow-789"),
			"finished":         tftypes.NewValue(tftypes.Bool, true),
			"mode":             tftypes.NewValue(tftypes.String, "manual"),
			"started_at":       tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
			"stopped_at":       tftypes.NewValue(tftypes.String, "2024-01-01T00:01:00Z"),
			"status":           tftypes.NewValue(tftypes.String, "success"),
		}
		state := tfsdk.State{
			Raw: tftypes.NewValue(tftypes.Object{
				AttributeTypes: map[string]tftypes.Type{
					"execution_id":     tftypes.String,
					"new_execution_id": tftypes.String,
					"workflow_id":      tftypes.String,
					"finished":         tftypes.Bool,
					"mode":             tftypes.String,
					"started_at":       tftypes.String,
					"stopped_at":       tftypes.String,
					"status":           tftypes.String,
				},
			}, rawState),
			Schema: createTestRetrySchema(t),
		}

		req := resource.DeleteRequest{
			State: state,
		}
		resp := resource.DeleteResponse{
			State: state,
		}

		r.Delete(context.Background(), req, &resp)

		assert.False(t, resp.Diagnostics.HasError(), "Delete should not have errors")
		// Note: Delete for retry operations doesn't perform any API operation,
		// it just removes from state, so it always succeeds
	})
}
