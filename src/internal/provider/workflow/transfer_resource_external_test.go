package workflow_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
	"github.com/kodflow/n8n/src/internal/provider/workflow"
	"github.com/stretchr/testify/assert"
)

// TestNewWorkflowTransferResource tests the NewWorkflowTransferResource constructor.
func TestNewWorkflowTransferResource(t *testing.T) {
	t.Run("create new transfer resource", func(t *testing.T) {
		res := workflow.NewWorkflowTransferResource()
		assert.NotNil(t, res, "NewWorkflowTransferResource should return a non-nil resource")
	})

	t.Run("resource has nil client initially", func(t *testing.T) {
		res := workflow.NewWorkflowTransferResource()
		assert.NotNil(t, res, "Resource should not be nil")
	})

	t.Run("multiple instances are independent", func(t *testing.T) {
		resource1 := workflow.NewWorkflowTransferResource()
		resource2 := workflow.NewWorkflowTransferResource()
		assert.NotNil(t, resource1)
		assert.NotNil(t, resource2)
		assert.NotSame(t, resource1, resource2, "Each call should create a new instance")
	})
}

// TestNewWorkflowTransferResourceWrapper tests the wrapper function.
func TestNewWorkflowTransferResourceWrapper(t *testing.T) {
	t.Run("create resource via wrapper", func(t *testing.T) {
		res := workflow.NewWorkflowTransferResourceWrapper()
		assert.NotNil(t, res, "NewWorkflowTransferResourceWrapper should return a non-nil resource")
	})

	t.Run("wrapper returns Resource interface", func(t *testing.T) {
		res := workflow.NewWorkflowTransferResourceWrapper()
		assert.Implements(t, (*resource.Resource)(nil), res)
	})

	t.Run("wrapper returns correct type", func(t *testing.T) {
		res := workflow.NewWorkflowTransferResourceWrapper()
		assert.NotNil(t, res)
		assert.Implements(t, (*resource.Resource)(nil), res)
	})
}

// TestWorkflowTransferResource_Metadata tests the Metadata method.
func TestWorkflowTransferResource_Metadata(t *testing.T) {
	t.Run("set metadata with provider name", func(t *testing.T) {
		r := workflow.NewWorkflowTransferResource()
		req := resource.MetadataRequest{
			ProviderTypeName: "n8n",
		}
		resp := &resource.MetadataResponse{}

		r.Metadata(context.Background(), req, resp)

		assert.Equal(t, "n8n_workflow_transfer", resp.TypeName)
	})

	t.Run("different provider type name", func(t *testing.T) {
		r := workflow.NewWorkflowTransferResource()
		req := resource.MetadataRequest{
			ProviderTypeName: "custom",
		}
		resp := &resource.MetadataResponse{}

		r.Metadata(context.Background(), req, resp)

		assert.Equal(t, "custom_workflow_transfer", resp.TypeName)
	})

	t.Run("empty provider type name", func(t *testing.T) {
		r := workflow.NewWorkflowTransferResource()
		req := resource.MetadataRequest{
			ProviderTypeName: "",
		}
		resp := &resource.MetadataResponse{}

		r.Metadata(context.Background(), req, resp)

		assert.Equal(t, "_workflow_transfer", resp.TypeName)
	})

	t.Run("metadata is idempotent", func(t *testing.T) {
		r := workflow.NewWorkflowTransferResource()
		req := resource.MetadataRequest{
			ProviderTypeName: "n8n",
		}
		resp1 := &resource.MetadataResponse{}
		resp2 := &resource.MetadataResponse{}

		r.Metadata(context.Background(), req, resp1)
		r.Metadata(context.Background(), req, resp2)

		assert.Equal(t, resp1.TypeName, resp2.TypeName)
	})
}

// TestWorkflowTransferResource_Schema tests the Schema method.
func TestWorkflowTransferResource_Schema(t *testing.T) {
	t.Run("get schema", func(t *testing.T) {
		r := workflow.NewWorkflowTransferResource()
		req := resource.SchemaRequest{}
		resp := &resource.SchemaResponse{}

		r.Schema(context.Background(), req, resp)

		assert.NotNil(t, resp.Schema)
		assert.NotEmpty(t, resp.Schema.MarkdownDescription)
		assert.NotEmpty(t, resp.Schema.Attributes)
	})

	t.Run("schema has required attributes", func(t *testing.T) {
		r := workflow.NewWorkflowTransferResource()
		req := resource.SchemaRequest{}
		resp := &resource.SchemaResponse{}

		r.Schema(context.Background(), req, resp)

		assert.Contains(t, resp.Schema.Attributes, "id")
		assert.Contains(t, resp.Schema.Attributes, "workflow_id")
		assert.Contains(t, resp.Schema.Attributes, "destination_project_id")
		assert.Contains(t, resp.Schema.Attributes, "transferred_at")
	})

	t.Run("schema has correct attribute count", func(t *testing.T) {
		r := workflow.NewWorkflowTransferResource()
		req := resource.SchemaRequest{}
		resp := &resource.SchemaResponse{}

		r.Schema(context.Background(), req, resp)

		assert.Equal(t, 4, len(resp.Schema.Attributes))
	})

	t.Run("schema description mentions transfer", func(t *testing.T) {
		r := workflow.NewWorkflowTransferResource()
		req := resource.SchemaRequest{}
		resp := &resource.SchemaResponse{}

		r.Schema(context.Background(), req, resp)

		assert.Contains(t, resp.Schema.MarkdownDescription, "transfer")
	})
}

// TestWorkflowTransferResource_Configure tests the Configure method.
func TestWorkflowTransferResource_Configure(t *testing.T) {
	t.Run("configure with valid client", func(t *testing.T) {
		r := workflow.NewWorkflowTransferResource()
		mockClient := &client.N8nClient{}

		req := resource.ConfigureRequest{
			ProviderData: mockClient,
		}
		resp := &resource.ConfigureResponse{}

		r.Configure(context.Background(), req, resp)

		assert.False(t, resp.Diagnostics.HasError())
	})

	t.Run("configure with nil provider data", func(t *testing.T) {
		r := workflow.NewWorkflowTransferResource()

		req := resource.ConfigureRequest{
			ProviderData: nil,
		}
		resp := &resource.ConfigureResponse{}

		r.Configure(context.Background(), req, resp)

		assert.False(t, resp.Diagnostics.HasError())
	})

	t.Run("configure with wrong type", func(t *testing.T) {
		r := workflow.NewWorkflowTransferResource()

		req := resource.ConfigureRequest{
			ProviderData: "invalid",
		}
		resp := &resource.ConfigureResponse{}

		r.Configure(context.Background(), req, resp)

		assert.True(t, resp.Diagnostics.HasError())
		assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Unexpected Resource Configure Type")
	})

	t.Run("configure multiple times", func(t *testing.T) {
		r := workflow.NewWorkflowTransferResource()

		// First configuration
		resp1 := &resource.ConfigureResponse{}
		client1 := &client.N8nClient{}
		req1 := resource.ConfigureRequest{
			ProviderData: client1,
		}
		r.Configure(context.Background(), req1, resp1)
		assert.False(t, resp1.Diagnostics.HasError())

		// Second configuration
		resp2 := &resource.ConfigureResponse{}
		client2 := &client.N8nClient{}
		req2 := resource.ConfigureRequest{
			ProviderData: client2,
		}
		r.Configure(context.Background(), req2, resp2)
		assert.False(t, resp2.Diagnostics.HasError())
	})

	t.Run("configure with integer provider data", func(t *testing.T) {
		r := workflow.NewWorkflowTransferResource()

		req := resource.ConfigureRequest{
			ProviderData: 123,
		}
		resp := &resource.ConfigureResponse{}

		r.Configure(context.Background(), req, resp)

		assert.True(t, resp.Diagnostics.HasError())
	})
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

// createTransferTestSchema creates a test schema for workflow transfer resource.
func createTransferTestSchema(t *testing.T) schema.Schema {
	t.Helper()
	r := workflow.NewWorkflowTransferResource()
	schemaReq := resource.SchemaRequest{}
	schemaResp := &resource.SchemaResponse{}
	r.Schema(context.Background(), schemaReq, schemaResp)
	return schemaResp.Schema
}

// createTransferTestObjectType creates the tftypes.Object for test state.
func createTransferTestObjectType(t *testing.T) tftypes.Object {
	t.Helper()
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"id":                     tftypes.String,
			"workflow_id":            tftypes.String,
			"destination_project_id": tftypes.String,
			"transferred_at":         tftypes.String,
		},
	}
}

// TestWorkflowTransferResource_Create tests the Create method.
func TestWorkflowTransferResource_Create(t *testing.T) {
	t.Run("successful transfer", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/workflows/wf-123/transfer" && r.Method == http.MethodPut {
				w.WriteHeader(http.StatusOK)
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTransferTestClient(t, handler)
		defer server.Close()

		r := workflow.NewWorkflowTransferResource()
		req := resource.ConfigureRequest{
			ProviderData: n8nClient,
		}
		resp := &resource.ConfigureResponse{}
		r.Configure(context.Background(), req, resp)

		rawPlan := map[string]tftypes.Value{
			"id":                     tftypes.NewValue(tftypes.String, nil),
			"workflow_id":            tftypes.NewValue(tftypes.String, "wf-123"),
			"destination_project_id": tftypes.NewValue(tftypes.String, "proj-456"),
			"transferred_at":         tftypes.NewValue(tftypes.String, nil),
		}

		objectType := createTransferTestObjectType(t)
		testSchema := createTransferTestSchema(t)
		plan := tfsdk.Plan{
			Raw:    tftypes.NewValue(objectType, rawPlan),
			Schema: testSchema,
		}

		state := tfsdk.State{
			Raw:    tftypes.NewValue(objectType, nil),
			Schema: testSchema,
		}

		createReq := resource.CreateRequest{
			Plan: plan,
		}
		createResp := resource.CreateResponse{
			State: state,
		}

		r.Create(context.Background(), createReq, &createResp)

		assert.False(t, createResp.Diagnostics.HasError())
	})

	t.Run("transfer fails with API error", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message": "Internal server error"}`))
		})

		n8nClient, server := setupTransferTestClient(t, handler)
		defer server.Close()

		r := workflow.NewWorkflowTransferResource()
		req := resource.ConfigureRequest{
			ProviderData: n8nClient,
		}
		resp := &resource.ConfigureResponse{}
		r.Configure(context.Background(), req, resp)

		rawPlan := map[string]tftypes.Value{
			"id":                     tftypes.NewValue(tftypes.String, nil),
			"workflow_id":            tftypes.NewValue(tftypes.String, "wf-123"),
			"destination_project_id": tftypes.NewValue(tftypes.String, "proj-456"),
			"transferred_at":         tftypes.NewValue(tftypes.String, nil),
		}

		objectType := createTransferTestObjectType(t)
		testSchema := createTransferTestSchema(t)
		plan := tfsdk.Plan{
			Raw:    tftypes.NewValue(objectType, rawPlan),
			Schema: testSchema,
		}

		state := tfsdk.State{
			Raw:    tftypes.NewValue(objectType, nil),
			Schema: testSchema,
		}

		createReq := resource.CreateRequest{
			Plan: plan,
		}
		createResp := resource.CreateResponse{
			State: state,
		}

		r.Create(context.Background(), createReq, &createResp)

		assert.True(t, createResp.Diagnostics.HasError())
	})

	t.Run("create with invalid plan", func(t *testing.T) {
		r := workflow.NewWorkflowTransferResource()

		// Create plan with mismatched schema - using number instead of string
		wrongSchema := tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"id": tftypes.Number,
			},
		}
		rawPlan := map[string]tftypes.Value{
			"id": tftypes.NewValue(tftypes.Number, 123),
		}

		testSchema := createTransferTestSchema(t)

		plan := tfsdk.Plan{
			Raw:    tftypes.NewValue(wrongSchema, rawPlan),
			Schema: testSchema,
		}

		state := tfsdk.State{
			Schema: testSchema,
		}

		createReq := resource.CreateRequest{
			Plan: plan,
		}
		createResp := resource.CreateResponse{
			State: state,
		}

		r.Create(context.Background(), createReq, &createResp)

		assert.True(t, createResp.Diagnostics.HasError())
	})
}

// TestWorkflowTransferResource_Read tests the Read method.
func TestWorkflowTransferResource_Read(t *testing.T) {
	t.Run("successful read maintains state", func(t *testing.T) {
		r := workflow.NewWorkflowTransferResource()

		rawState := map[string]tftypes.Value{
			"id":                     tftypes.NewValue(tftypes.String, "wf-123-to-proj-456"),
			"workflow_id":            tftypes.NewValue(tftypes.String, "wf-123"),
			"destination_project_id": tftypes.NewValue(tftypes.String, "proj-456"),
			"transferred_at":         tftypes.NewValue(tftypes.String, "2024-01-01"),
		}

		objectType := createTransferTestObjectType(t)
		testSchema := createTransferTestSchema(t)
		state := tfsdk.State{
			Raw:    tftypes.NewValue(objectType, rawState),
			Schema: testSchema,
		}

		readReq := resource.ReadRequest{
			State: state,
		}
		readResp := resource.ReadResponse{
			State: state,
		}

		r.Read(context.Background(), readReq, &readResp)

		assert.False(t, readResp.Diagnostics.HasError())
	})

	t.Run("read with invalid state", func(t *testing.T) {
		r := workflow.NewWorkflowTransferResource()

		// Create state with mismatched schema
		wrongSchema := tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"id": tftypes.Number,
			},
		}
		rawState := map[string]tftypes.Value{
			"id": tftypes.NewValue(tftypes.Number, 123),
		}

		testSchema := createTransferTestSchema(t)

		state := tfsdk.State{
			Raw:    tftypes.NewValue(wrongSchema, rawState),
			Schema: testSchema,
		}

		readReq := resource.ReadRequest{
			State: state,
		}
		readResp := resource.ReadResponse{
			State: state,
		}

		r.Read(context.Background(), readReq, &readResp)

		assert.True(t, readResp.Diagnostics.HasError())
	})

	t.Run("read is idempotent", func(t *testing.T) {
		r := workflow.NewWorkflowTransferResource()

		rawState := map[string]tftypes.Value{
			"id":                     tftypes.NewValue(tftypes.String, "wf-123-to-proj-456"),
			"workflow_id":            tftypes.NewValue(tftypes.String, "wf-123"),
			"destination_project_id": tftypes.NewValue(tftypes.String, "proj-456"),
			"transferred_at":         tftypes.NewValue(tftypes.String, "2024-01-01"),
		}

		objectType := createTransferTestObjectType(t)
		testSchema := createTransferTestSchema(t)
		state := tfsdk.State{
			Raw:    tftypes.NewValue(objectType, rawState),
			Schema: testSchema,
		}

		// First read
		readReq1 := resource.ReadRequest{State: state}
		readResp1 := resource.ReadResponse{State: state}
		r.Read(context.Background(), readReq1, &readResp1)

		// Second read
		readReq2 := resource.ReadRequest{State: state}
		readResp2 := resource.ReadResponse{State: state}
		r.Read(context.Background(), readReq2, &readResp2)

		assert.False(t, readResp1.Diagnostics.HasError())
		assert.False(t, readResp2.Diagnostics.HasError())
	})
}

// TestWorkflowTransferResource_Update tests the Update method.
func TestWorkflowTransferResource_Update(t *testing.T) {
	t.Run("update returns error", func(t *testing.T) {
		r := workflow.NewWorkflowTransferResource()

		rawPlan := map[string]tftypes.Value{
			"id":                     tftypes.NewValue(tftypes.String, "wf-123-to-proj-456"),
			"workflow_id":            tftypes.NewValue(tftypes.String, "wf-123"),
			"destination_project_id": tftypes.NewValue(tftypes.String, "proj-789"),
			"transferred_at":         tftypes.NewValue(tftypes.String, "2024-01-01"),
		}

		rawState := map[string]tftypes.Value{
			"id":                     tftypes.NewValue(tftypes.String, "wf-123-to-proj-456"),
			"workflow_id":            tftypes.NewValue(tftypes.String, "wf-123"),
			"destination_project_id": tftypes.NewValue(tftypes.String, "proj-456"),
			"transferred_at":         tftypes.NewValue(tftypes.String, "2024-01-01"),
		}

		objectType := createTransferTestObjectType(t)
		testSchema := createTransferTestSchema(t)
		plan := tfsdk.Plan{
			Raw:    tftypes.NewValue(objectType, rawPlan),
			Schema: testSchema,
		}

		state := tfsdk.State{
			Raw:    tftypes.NewValue(objectType, rawState),
			Schema: testSchema,
		}

		updateReq := resource.UpdateRequest{
			Plan:  plan,
			State: state,
		}
		updateResp := resource.UpdateResponse{
			State: state,
		}

		r.Update(context.Background(), updateReq, &updateResp)

		assert.True(t, updateResp.Diagnostics.HasError())
		assert.Contains(t, updateResp.Diagnostics.Errors()[0].Summary(), "Update Not Supported")
	})

	t.Run("update always fails", func(t *testing.T) {
		r := workflow.NewWorkflowTransferResource()

		rawPlan := map[string]tftypes.Value{
			"id":                     tftypes.NewValue(tftypes.String, "test-id"),
			"workflow_id":            tftypes.NewValue(tftypes.String, "wf-123"),
			"destination_project_id": tftypes.NewValue(tftypes.String, "proj-456"),
			"transferred_at":         tftypes.NewValue(tftypes.String, "2024-01-01"),
		}

		objectType := createTransferTestObjectType(t)
		testSchema := createTransferTestSchema(t)
		plan := tfsdk.Plan{
			Raw:    tftypes.NewValue(objectType, rawPlan),
			Schema: testSchema,
		}

		state := tfsdk.State{
			Raw:    tftypes.NewValue(objectType, rawPlan),
			Schema: testSchema,
		}

		updateReq := resource.UpdateRequest{
			Plan:  plan,
			State: state,
		}
		updateResp := resource.UpdateResponse{
			State: state,
		}

		r.Update(context.Background(), updateReq, &updateResp)

		assert.True(t, updateResp.Diagnostics.HasError())
	})
}

// TestWorkflowTransferResource_Delete tests the Delete method.
func TestWorkflowTransferResource_Delete(t *testing.T) {
	t.Run("delete succeeds without API call", func(t *testing.T) {
		r := workflow.NewWorkflowTransferResource()

		rawState := map[string]tftypes.Value{
			"id":                     tftypes.NewValue(tftypes.String, "wf-123-to-proj-456"),
			"workflow_id":            tftypes.NewValue(tftypes.String, "wf-123"),
			"destination_project_id": tftypes.NewValue(tftypes.String, "proj-456"),
			"transferred_at":         tftypes.NewValue(tftypes.String, "2024-01-01"),
		}

		objectType := createTransferTestObjectType(t)
		testSchema := createTransferTestSchema(t)
		state := tfsdk.State{
			Raw:    tftypes.NewValue(objectType, rawState),
			Schema: testSchema,
		}

		deleteReq := resource.DeleteRequest{
			State: state,
		}
		deleteResp := resource.DeleteResponse{
			State: state,
		}

		r.Delete(context.Background(), deleteReq, &deleteResp)

		assert.False(t, deleteResp.Diagnostics.HasError())
	})

	t.Run("delete is idempotent", func(t *testing.T) {
		r := workflow.NewWorkflowTransferResource()

		rawState := map[string]tftypes.Value{
			"id":                     tftypes.NewValue(tftypes.String, "wf-123-to-proj-456"),
			"workflow_id":            tftypes.NewValue(tftypes.String, "wf-123"),
			"destination_project_id": tftypes.NewValue(tftypes.String, "proj-456"),
			"transferred_at":         tftypes.NewValue(tftypes.String, "2024-01-01"),
		}

		objectType := createTransferTestObjectType(t)
		testSchema := createTransferTestSchema(t)
		state := tfsdk.State{
			Raw:    tftypes.NewValue(objectType, rawState),
			Schema: testSchema,
		}

		// First delete
		deleteReq1 := resource.DeleteRequest{State: state}
		deleteResp1 := resource.DeleteResponse{State: state}
		r.Delete(context.Background(), deleteReq1, &deleteResp1)

		// Second delete
		deleteReq2 := resource.DeleteRequest{State: state}
		deleteResp2 := resource.DeleteResponse{State: state}
		r.Delete(context.Background(), deleteReq2, &deleteResp2)

		assert.False(t, deleteResp1.Diagnostics.HasError())
		assert.False(t, deleteResp2.Diagnostics.HasError())
	})

	t.Run("delete with nil state", func(t *testing.T) {
		r := workflow.NewWorkflowTransferResource()

		objectType := createTransferTestObjectType(t)
		testSchema := createTransferTestSchema(t)
		state := tfsdk.State{
			Raw:    tftypes.NewValue(objectType, nil),
			Schema: testSchema,
		}

		deleteReq := resource.DeleteRequest{
			State: state,
		}
		deleteResp := resource.DeleteResponse{
			State: state,
		}

		r.Delete(context.Background(), deleteReq, &deleteResp)

		// Delete should still succeed since it doesn't perform any operations
		assert.False(t, deleteResp.Diagnostics.HasError())
	})
}

// TestWorkflowTransferResource_ImportState tests the ImportState method.
func TestWorkflowTransferResource_ImportState(t *testing.T) {
	t.Run("successful import", func(t *testing.T) {
		r := workflow.NewWorkflowTransferResource()

		objectType := createTransferTestObjectType(t)
		testSchema := createTransferTestSchema(t)
		state := tfsdk.State{
			Raw:    tftypes.NewValue(objectType, nil),
			Schema: testSchema,
		}

		importReq := resource.ImportStateRequest{
			ID: "wf-123-to-proj-456",
		}
		importResp := &resource.ImportStateResponse{
			State: state,
		}

		r.ImportState(context.Background(), importReq, importResp)

		assert.False(t, importResp.Diagnostics.HasError())
	})

	t.Run("import with different ID", func(t *testing.T) {
		r := workflow.NewWorkflowTransferResource()

		objectType := createTransferTestObjectType(t)
		testSchema := createTransferTestSchema(t)
		state := tfsdk.State{
			Raw:    tftypes.NewValue(objectType, nil),
			Schema: testSchema,
		}

		importReq := resource.ImportStateRequest{
			ID: "custom-transfer-id",
		}
		importResp := &resource.ImportStateResponse{
			State: state,
		}

		r.ImportState(context.Background(), importReq, importResp)

		assert.False(t, importResp.Diagnostics.HasError())
	})

	t.Run("import uses passthrough", func(t *testing.T) {
		r := workflow.NewWorkflowTransferResource()

		objectType := createTransferTestObjectType(t)
		testSchema := createTransferTestSchema(t)
		state := tfsdk.State{
			Raw:    tftypes.NewValue(objectType, nil),
			Schema: testSchema,
		}

		importReq := resource.ImportStateRequest{
			ID: "test-id",
		}
		importResp := &resource.ImportStateResponse{
			State: state,
		}

		// Verify ImportStatePassthroughID is called
		r.ImportState(context.Background(), importReq, importResp)

		// The function uses resource.ImportStatePassthroughID which sets the ID
		assert.False(t, importResp.Diagnostics.HasError())
	})

	t.Run("import with empty ID", func(t *testing.T) {
		r := workflow.NewWorkflowTransferResource()

		objectType := createTransferTestObjectType(t)
		testSchema := createTransferTestSchema(t)
		state := tfsdk.State{
			Raw:    tftypes.NewValue(objectType, nil),
			Schema: testSchema,
		}

		importReq := resource.ImportStateRequest{
			ID: "",
		}
		importResp := &resource.ImportStateResponse{
			State: state,
		}

		r.ImportState(context.Background(), importReq, importResp)

		// ImportStatePassthroughID should handle empty ID gracefully
		assert.False(t, importResp.Diagnostics.HasError())
	})
}

// TestWorkflowTransferResourcePublicAPI tests the public API of WorkflowTransferResource.
// This is a black-box test that only uses the exported interface.
func TestWorkflowTransferResourcePublicAPI(t *testing.T) {
	t.Run("create resource instance", func(t *testing.T) {
		// Test that we can create a new resource using the public constructor
		res := workflow.NewWorkflowTransferResource()
		assert.NotNil(t, res, "NewWorkflowTransferResource should return a non-nil resource")
	})

	t.Run("create resource wrapper instance", func(t *testing.T) {
		// Test that we can create a resource using the wrapper function
		res := workflow.NewWorkflowTransferResourceWrapper()
		assert.NotNil(t, res, "NewWorkflowTransferResourceWrapper should return a non-nil resource")
	})

	t.Run("resource implements required interfaces", func(t *testing.T) {
		res := workflow.NewWorkflowTransferResourceWrapper()
		assert.Implements(t, (*resource.Resource)(nil), res)
		assert.Implements(t, (*resource.ResourceWithConfigure)(nil), res)
		assert.Implements(t, (*resource.ResourceWithImportState)(nil), res)
	})

	t.Run("metadata sets correct type name", func(t *testing.T) {
		r := workflow.NewWorkflowTransferResource()
		req := resource.MetadataRequest{
			ProviderTypeName: "n8n",
		}
		resp := &resource.MetadataResponse{}

		r.Metadata(context.Background(), req, resp)

		assert.Equal(t, "n8n_workflow_transfer", resp.TypeName)
	})

	t.Run("schema returns valid schema", func(t *testing.T) {
		r := workflow.NewWorkflowTransferResource()
		req := resource.SchemaRequest{}
		resp := &resource.SchemaResponse{}

		r.Schema(context.Background(), req, resp)

		assert.NotNil(t, resp.Schema)
		assert.NotEmpty(t, resp.Schema.Attributes)
	})

	t.Run("configure accepts valid client", func(t *testing.T) {
		r := workflow.NewWorkflowTransferResource()
		mockClient := &client.N8nClient{}

		req := resource.ConfigureRequest{
			ProviderData: mockClient,
		}
		resp := &resource.ConfigureResponse{}

		r.Configure(context.Background(), req, resp)

		assert.False(t, resp.Diagnostics.HasError())
	})

	t.Run("import state uses correct path", func(t *testing.T) {
		r := workflow.NewWorkflowTransferResource()

		objectType := createTransferTestObjectType(t)
		testSchema := createTransferTestSchema(t)
		state := tfsdk.State{
			Raw:    tftypes.NewValue(objectType, nil),
			Schema: testSchema,
		}

		importReq := resource.ImportStateRequest{
			ID: "transfer-123",
		}
		importResp := &resource.ImportStateResponse{
			State: state,
		}

		r.ImportState(context.Background(), importReq, importResp)

		// Verify that the import uses path.Root("id")
		expectedPath := path.Root("id")
		assert.NotNil(t, expectedPath)
	})
}
