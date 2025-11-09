package workflow

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

func TestNewWorkflowTransferResource(t *testing.T) {
	t.Run("create new resource", func(t *testing.T) {
		r := NewWorkflowTransferResource()
		assert.NotNil(t, r)
		assert.Nil(t, r.client)
	})

	t.Run("create new resource wrapper", func(t *testing.T) {
		r := NewWorkflowTransferResourceWrapper()
		assert.NotNil(t, r)
		_, ok := r.(*WorkflowTransferResource)
		assert.True(t, ok)
	})
}

func TestWorkflowTransferResource_Metadata(t *testing.T) {
	t.Run("set metadata", func(t *testing.T) {
		r := &WorkflowTransferResource{}
		req := resource.MetadataRequest{
			ProviderTypeName: "n8n",
		}
		resp := &resource.MetadataResponse{}

		r.Metadata(context.Background(), req, resp)

		assert.Equal(t, "n8n_workflow_transfer", resp.TypeName)
	})

	t.Run("different provider type name", func(t *testing.T) {
		r := &WorkflowTransferResource{}
		req := resource.MetadataRequest{
			ProviderTypeName: "custom",
		}
		resp := &resource.MetadataResponse{}

		r.Metadata(context.Background(), req, resp)

		assert.Equal(t, "custom_workflow_transfer", resp.TypeName)
	})

	t.Run("empty provider type name", func(t *testing.T) {
		r := &WorkflowTransferResource{}
		req := resource.MetadataRequest{
			ProviderTypeName: "",
		}
		resp := &resource.MetadataResponse{}

		r.Metadata(context.Background(), req, resp)

		assert.Equal(t, "_workflow_transfer", resp.TypeName)
	})
}

func TestWorkflowTransferResource_Schema(t *testing.T) {
	t.Run("get schema", func(t *testing.T) {
		r := &WorkflowTransferResource{}
		req := resource.SchemaRequest{}
		resp := &resource.SchemaResponse{}

		r.Schema(context.Background(), req, resp)

		assert.NotNil(t, resp.Schema)
		assert.Contains(t, resp.Schema.MarkdownDescription, "Transfers a workflow")
		assert.NotEmpty(t, resp.Schema.Attributes)
	})

	t.Run("schema has required attributes", func(t *testing.T) {
		r := &WorkflowTransferResource{}
		resp := &resource.SchemaResponse{}

		r.Schema(context.Background(), resource.SchemaRequest{}, resp)

		// Check for id attribute (computed)
		idAttr, exists := resp.Schema.Attributes["id"]
		assert.True(t, exists)
		assert.True(t, idAttr.IsComputed())

		// Check for workflow_id attribute (required)
		workflowAttr, exists := resp.Schema.Attributes["workflow_id"]
		assert.True(t, exists)
		assert.True(t, workflowAttr.IsRequired())

		// Check for destination_project_id attribute (required)
		destAttr, exists := resp.Schema.Attributes["destination_project_id"]
		assert.True(t, exists)
		assert.True(t, destAttr.IsRequired())

		// Check for transferred_at attribute (computed)
		transferredAttr, exists := resp.Schema.Attributes["transferred_at"]
		assert.True(t, exists)
		assert.True(t, transferredAttr.IsComputed())
	})

	t.Run("schema attributes have descriptions", func(t *testing.T) {
		r := &WorkflowTransferResource{}
		resp := &resource.SchemaResponse{}

		r.Schema(context.Background(), resource.SchemaRequest{}, resp)

		for name, attr := range resp.Schema.Attributes {
			assert.NotEmpty(t, attr.GetMarkdownDescription(), "Attribute %s should have description", name)
		}
	})

	t.Run("schema has correct attribute count", func(t *testing.T) {
		r := &WorkflowTransferResource{}
		resp := &resource.SchemaResponse{}

		r.Schema(context.Background(), resource.SchemaRequest{}, resp)

		// Should have id, workflow_id, destination_project_id, transferred_at
		assert.Len(t, resp.Schema.Attributes, 4)
	})
}

func TestWorkflowTransferResource_Configure(t *testing.T) {
	t.Run("configure with valid client", func(t *testing.T) {
		r := &WorkflowTransferResource{}
		mockClient := &client.N8nClient{}

		req := resource.ConfigureRequest{
			ProviderData: mockClient,
		}
		resp := &resource.ConfigureResponse{}

		r.Configure(context.Background(), req, resp)

		assert.Equal(t, mockClient, r.client)
		assert.False(t, resp.Diagnostics.HasError())
	})

	t.Run("configure with nil provider data", func(t *testing.T) {
		r := &WorkflowTransferResource{}

		req := resource.ConfigureRequest{
			ProviderData: nil,
		}
		resp := &resource.ConfigureResponse{}

		r.Configure(context.Background(), req, resp)

		assert.Nil(t, r.client)
		assert.False(t, resp.Diagnostics.HasError())
	})

	t.Run("configure with wrong type", func(t *testing.T) {
		r := &WorkflowTransferResource{}

		req := resource.ConfigureRequest{
			ProviderData: "invalid",
		}
		resp := &resource.ConfigureResponse{}

		r.Configure(context.Background(), req, resp)

		assert.Nil(t, r.client)
		assert.True(t, resp.Diagnostics.HasError())
		assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Unexpected Resource Configure Type")
	})

	t.Run("configure with struct type", func(t *testing.T) {
		r := &WorkflowTransferResource{}

		req := resource.ConfigureRequest{
			ProviderData: struct{}{},
		}
		resp := &resource.ConfigureResponse{}

		r.Configure(context.Background(), req, resp)

		assert.Nil(t, r.client)
		assert.True(t, resp.Diagnostics.HasError())
	})
}

// createTransferTestSchema creates a test schema for workflow transfer resource.
func createTransferTestSchema(t *testing.T) schema.Schema {
	t.Helper()
	r := &WorkflowTransferResource{}
	req := resource.SchemaRequest{}
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), req, resp)
	return resp.Schema
}

// setupTransferTestClient creates a test N8nClient with httptest server.
func setupTransferTestClient(t *testing.T, handler http.HandlerFunc) (*client.N8nClient, *httptest.Server) {
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

func TestWorkflowTransferResource_Create(t *testing.T) {
	t.Run("successful transfer", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/workflows/wf-123/transfer" && r.Method == http.MethodPut {
				// Verify request body
				var reqBody map[string]interface{}
				err := json.NewDecoder(r.Body).Decode(&reqBody)
				assert.NoError(t, err)
				assert.Equal(t, "proj-456", reqBody["destinationProjectId"])

				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]interface{}{
					"success": true,
				})
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTransferTestClient(t, handler)
		defer server.Close()

		r := &WorkflowTransferResource{client: n8nClient}

		rawPlan := map[string]tftypes.Value{
			"id":                     tftypes.NewValue(tftypes.String, nil),
			"workflow_id":            tftypes.NewValue(tftypes.String, "wf-123"),
			"destination_project_id": tftypes.NewValue(tftypes.String, "proj-456"),
			"transferred_at":         tftypes.NewValue(tftypes.String, nil),
		}
		plan := tfsdk.Plan{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "workflow_id": tftypes.String, "destination_project_id": tftypes.String, "transferred_at": tftypes.String}}, rawPlan),
			Schema: createTransferTestSchema(t),
		}

		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "workflow_id": tftypes.String, "destination_project_id": tftypes.String, "transferred_at": tftypes.String}}, nil),
			Schema: createTransferTestSchema(t),
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

	t.Run("transfer fails", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/workflows/wf-123/transfer" && r.Method == http.MethodPut {
				w.WriteHeader(http.StatusBadRequest)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]interface{}{
					"message": "Invalid destination project",
				})
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTransferTestClient(t, handler)
		defer server.Close()

		r := &WorkflowTransferResource{client: n8nClient}

		rawPlan := map[string]tftypes.Value{
			"id":                     tftypes.NewValue(tftypes.String, nil),
			"workflow_id":            tftypes.NewValue(tftypes.String, "wf-123"),
			"destination_project_id": tftypes.NewValue(tftypes.String, "invalid-proj"),
			"transferred_at":         tftypes.NewValue(tftypes.String, nil),
		}
		plan := tfsdk.Plan{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "workflow_id": tftypes.String, "destination_project_id": tftypes.String, "transferred_at": tftypes.String}}, rawPlan),
			Schema: createTransferTestSchema(t),
		}

		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "workflow_id": tftypes.String, "destination_project_id": tftypes.String, "transferred_at": tftypes.String}}, nil),
			Schema: createTransferTestSchema(t),
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

func TestWorkflowTransferResource_Read(t *testing.T) {
	t.Run("successful read", func(t *testing.T) {
		r := &WorkflowTransferResource{}

		rawState := map[string]tftypes.Value{
			"id":                     tftypes.NewValue(tftypes.String, "wf-123-to-proj-456"),
			"workflow_id":            tftypes.NewValue(tftypes.String, "wf-123"),
			"destination_project_id": tftypes.NewValue(tftypes.String, "proj-456"),
			"transferred_at":         tftypes.NewValue(tftypes.String, "transfer-response-200"),
		}
		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "workflow_id": tftypes.String, "destination_project_id": tftypes.String, "transferred_at": tftypes.String}}, rawState),
			Schema: createTransferTestSchema(t),
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
}

func TestWorkflowTransferResource_Update(t *testing.T) {
	t.Run("update returns error", func(t *testing.T) {
		r := &WorkflowTransferResource{}
		resp := &resource.UpdateResponse{}

		r.Update(context.Background(), resource.UpdateRequest{}, resp)

		assert.True(t, resp.Diagnostics.HasError())
		assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Update Not Supported")
	})
}

func TestWorkflowTransferResource_Delete(t *testing.T) {
	t.Run("successful delete", func(t *testing.T) {
		r := &WorkflowTransferResource{}
		resp := &resource.DeleteResponse{}

		// Delete should not cause any errors
		r.Delete(context.Background(), resource.DeleteRequest{}, resp)

		assert.False(t, resp.Diagnostics.HasError())
	})

	t.Run("delete with state", func(t *testing.T) {
		r := &WorkflowTransferResource{}

		rawState := map[string]tftypes.Value{
			"id":                     tftypes.NewValue(tftypes.String, "wf-123-to-proj-456"),
			"workflow_id":            tftypes.NewValue(tftypes.String, "wf-123"),
			"destination_project_id": tftypes.NewValue(tftypes.String, "proj-456"),
			"transferred_at":         tftypes.NewValue(tftypes.String, "2024-01-01"),
		}
		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "workflow_id": tftypes.String, "destination_project_id": tftypes.String, "transferred_at": tftypes.String}}, rawState),
			Schema: createTransferTestSchema(t),
		}

		req := resource.DeleteRequest{
			State: state,
		}
		resp := &resource.DeleteResponse{}

		r.Delete(context.Background(), req, resp)

		assert.False(t, resp.Diagnostics.HasError(), "Delete with state should not have errors")
	})

	t.Run("delete multiple times", func(t *testing.T) {
		r := &WorkflowTransferResource{}

		// First delete
		resp1 := &resource.DeleteResponse{}
		r.Delete(context.Background(), resource.DeleteRequest{}, resp1)
		assert.False(t, resp1.Diagnostics.HasError())

		// Second delete
		resp2 := &resource.DeleteResponse{}
		r.Delete(context.Background(), resource.DeleteRequest{}, resp2)
		assert.False(t, resp2.Diagnostics.HasError())
	})
}

func TestWorkflowTransferResource_ImportState(t *testing.T) {
	t.Run("import state", func(t *testing.T) {
		r := &WorkflowTransferResource{}
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
			"id":                     tftypes.NewValue(tftypes.String, nil),
			"workflow_id":            tftypes.NewValue(tftypes.String, nil),
			"destination_project_id": tftypes.NewValue(tftypes.String, nil),
			"transferred_at":         tftypes.NewValue(tftypes.String, nil),
		})

		req := resource.ImportStateRequest{
			ID: "transfer-123",
		}
		resp := &resource.ImportStateResponse{
			State: state,
		}

		r.ImportState(ctx, req, resp)

		assert.NotNil(t, resp)
		assert.False(t, resp.Diagnostics.HasError())
	})
}

func TestWorkflowTransferResource_InterfaceCompliance(t *testing.T) {
	t.Run("implements Resource interface", func(t *testing.T) {
		var _ resource.Resource = (*WorkflowTransferResource)(nil)
	})

	t.Run("implements ResourceWithConfigure interface", func(t *testing.T) {
		var _ resource.ResourceWithConfigure = (*WorkflowTransferResource)(nil)
	})

	t.Run("implements ResourceWithImportState interface", func(t *testing.T) {
		var _ resource.ResourceWithImportState = (*WorkflowTransferResource)(nil)
	})

	t.Run("implements WorkflowTransferResourceInterface", func(t *testing.T) {
		var _ WorkflowTransferResourceInterface = (*WorkflowTransferResource)(nil)
	})
}

func TestWorkflowTransferResource_Concurrency(t *testing.T) {
	t.Run("concurrent metadata calls", func(t *testing.T) {
		r := &WorkflowTransferResource{}

		done := make(chan bool, 100)
		for i := 0; i < 100; i++ {
			go func() {
				resp := &resource.MetadataResponse{}
				r.Metadata(context.Background(), resource.MetadataRequest{
					ProviderTypeName: "n8n",
				}, resp)
				assert.Equal(t, "n8n_workflow_transfer", resp.TypeName)
				done <- true
			}()
		}

		for i := 0; i < 100; i++ {
			<-done
		}
	})

	t.Run("concurrent schema calls", func(t *testing.T) {
		r := &WorkflowTransferResource{}

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

func TestWorkflowTransferResource_SchemaStructure(t *testing.T) {
	t.Run("verify schema structure", func(t *testing.T) {
		r := &WorkflowTransferResource{}
		resp := &resource.SchemaResponse{}

		r.Schema(context.Background(), resource.SchemaRequest{}, resp)

		// Verify all attributes are StringAttribute
		for name, attr := range resp.Schema.Attributes {
			_, ok := attr.(schema.StringAttribute)
			assert.True(t, ok, "Attribute %s should be StringAttribute", name)
		}
	})
}

func BenchmarkWorkflowTransferResource_Schema(b *testing.B) {
	r := &WorkflowTransferResource{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := &resource.SchemaResponse{}
		r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	}
}

func BenchmarkWorkflowTransferResource_Metadata(b *testing.B) {
	r := &WorkflowTransferResource{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := &resource.MetadataResponse{}
		r.Metadata(context.Background(), resource.MetadataRequest{}, resp)
	}
}

// TestWorkflowTransferResource_EdgeCases tests edge cases for 100% coverage.
func TestWorkflowTransferResource_EdgeCases(t *testing.T) {
	t.Run("create fails when plan get fails", func(t *testing.T) {
		r := &WorkflowTransferResource{}

		// Create plan with mismatched schema
		wrongSchema := schema.Schema{
			Attributes: map[string]schema.Attribute{
				"id": schema.NumberAttribute{
					Required: true,
				},
			},
		}
		rawPlan := map[string]tftypes.Value{
			"id": tftypes.NewValue(tftypes.Number, 123),
		}
		plan := tfsdk.Plan{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.Number}}, rawPlan),
			Schema: wrongSchema,
		}

		req := resource.CreateRequest{
			Plan: plan,
		}

		schemaResp := resource.SchemaResponse{}
		r.Schema(context.Background(), resource.SchemaRequest{}, &schemaResp)

		resp := resource.CreateResponse{
			State: tfsdk.State{Schema: schemaResp.Schema},
		}

		r.Create(context.Background(), req, &resp)

		assert.True(t, resp.Diagnostics.HasError(), "Create should fail when Plan.Get fails")
	})

	t.Run("read fails when state get fails", func(t *testing.T) {
		r := &WorkflowTransferResource{}

		// Create state with mismatched schema
		wrongSchema := schema.Schema{
			Attributes: map[string]schema.Attribute{
				"id": schema.NumberAttribute{
					Required: true,
				},
			},
		}
		rawState := map[string]tftypes.Value{
			"id": tftypes.NewValue(tftypes.Number, 123),
		}
		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.Number}}, rawState),
			Schema: wrongSchema,
		}

		req := resource.ReadRequest{
			State: state,
		}

		schemaResp := resource.SchemaResponse{}
		r.Schema(context.Background(), resource.SchemaRequest{}, &schemaResp)

		resp := resource.ReadResponse{
			State: tfsdk.State{Schema: schemaResp.Schema},
		}

		r.Read(context.Background(), req, &resp)

		assert.True(t, resp.Diagnostics.HasError(), "Read should fail when State.Get fails")
	})

	t.Run("read succeeds but state set fails", func(t *testing.T) {
		r := &WorkflowTransferResource{}

		schemaResp := resource.SchemaResponse{}
		r.Schema(context.Background(), resource.SchemaRequest{}, &schemaResp)

		rawState := map[string]tftypes.Value{
			"id":                     tftypes.NewValue(tftypes.String, "wf-123-to-proj-456"),
			"workflow_id":            tftypes.NewValue(tftypes.String, "wf-123"),
			"destination_project_id": tftypes.NewValue(tftypes.String, "proj-456"),
			"transferred_at":         tftypes.NewValue(tftypes.String, "2024-01-01"),
		}

		attrTypes := make(map[string]tftypes.Type)
		for key := range schemaResp.Schema.Attributes {
			attrTypes[key] = tftypes.String
		}

		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: attrTypes}, rawState),
			Schema: schemaResp.Schema,
		}

		// Create response state with incompatible schema
		wrongSchema := schema.Schema{
			Attributes: map[string]schema.Attribute{
				"id": schema.NumberAttribute{
					Computed: true,
				},
			},
		}
		respState := tfsdk.State{
			Schema: wrongSchema,
		}

		req := resource.ReadRequest{
			State: state,
		}
		resp := resource.ReadResponse{
			State: respState,
		}

		r.Read(context.Background(), req, &resp)

		assert.True(t, resp.Diagnostics.HasError(), "Read should fail when State.Set fails")
	})
}

func BenchmarkWorkflowTransferResource_Configure(b *testing.B) {
	r := &WorkflowTransferResource{}
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
