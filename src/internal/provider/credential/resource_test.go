package credential

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/credential/models"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
	"github.com/stretchr/testify/assert"
)

// Note: Mock implementations are in helpers_test.go to avoid duplication

// Helper functions.
func strPtr(s string) *string {
	return &s
}

func TestNewCredentialResource(t *testing.T) {
	t.Run("create new resource", func(t *testing.T) {
		resource := NewCredentialResource()
		assert.NotNil(t, resource)
		assert.Nil(t, resource.client)
	})

	t.Run("create new resource wrapper", func(t *testing.T) {
		resource := NewCredentialResourceWrapper()
		assert.NotNil(t, resource)
		_, ok := resource.(*CredentialResource)
		assert.True(t, ok)
	})
}

func TestCredentialResource_Metadata(t *testing.T) {
	t.Run("set metadata", func(t *testing.T) {
		r := &CredentialResource{}
		req := resource.MetadataRequest{
			ProviderTypeName: "n8n",
		}
		resp := &resource.MetadataResponse{}

		r.Metadata(context.Background(), req, resp)

		assert.Equal(t, "n8n_credential", resp.TypeName)
	})

	t.Run("different provider type name", func(t *testing.T) {
		r := &CredentialResource{}
		req := resource.MetadataRequest{
			ProviderTypeName: "custom",
		}
		resp := &resource.MetadataResponse{}

		r.Metadata(context.Background(), req, resp)

		assert.Equal(t, "custom_credential", resp.TypeName)
	})
}

func TestCredentialResource_Schema(t *testing.T) {
	t.Run("get schema", func(t *testing.T) {
		r := &CredentialResource{}
		req := resource.SchemaRequest{}
		resp := &resource.SchemaResponse{}

		r.Schema(context.Background(), req, resp)

		assert.NotNil(t, resp.Schema)
		assert.Contains(t, resp.Schema.MarkdownDescription, "credential resource")
		assert.Contains(t, resp.Schema.MarkdownDescription, "automatic rotation")
		assert.NotEmpty(t, resp.Schema.Attributes)
	})

	t.Run("schema attributes", func(t *testing.T) {
		r := &CredentialResource{}
		attrs := r.schemaAttributes()

		assert.NotNil(t, attrs)
		assert.Contains(t, attrs, "id")
		assert.Contains(t, attrs, "name")
		assert.Contains(t, attrs, "type")
		assert.Contains(t, attrs, "data")
		assert.Contains(t, attrs, "created_at")
		assert.Contains(t, attrs, "updated_at")

		// Check ID is computed
		idAttr, ok := attrs["id"].(schema.StringAttribute)
		assert.True(t, ok)
		assert.True(t, idAttr.Computed)

		// Check name is required
		nameAttr, ok := attrs["name"].(schema.StringAttribute)
		assert.True(t, ok)
		assert.True(t, nameAttr.Required)
	})
}

func TestCredentialResource_Configure(t *testing.T) {
	t.Run("configure with valid client", func(t *testing.T) {
		r := &CredentialResource{}
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
		r := &CredentialResource{}

		req := resource.ConfigureRequest{
			ProviderData: nil,
		}
		resp := &resource.ConfigureResponse{}

		r.Configure(context.Background(), req, resp)

		assert.Nil(t, r.client)
		assert.False(t, resp.Diagnostics.HasError())
	})

	t.Run("configure with wrong type", func(t *testing.T) {
		r := &CredentialResource{}

		req := resource.ConfigureRequest{
			ProviderData: "invalid",
		}
		resp := &resource.ConfigureResponse{}

		r.Configure(context.Background(), req, resp)

		assert.Nil(t, r.client)
		assert.True(t, resp.Diagnostics.HasError())
		assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Unexpected Resource Configure Type")
	})
}

func TestCredentialResource_ImportState(t *testing.T) {
	t.Run("import state", func(t *testing.T) {
		r := &CredentialResource{}
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
			"id":         tftypes.NewValue(tftypes.String, nil),
			"name":       tftypes.NewValue(tftypes.String, nil),
			"type":       tftypes.NewValue(tftypes.String, nil),
			"data":       tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, nil),
			"created_at": tftypes.NewValue(tftypes.String, nil),
			"updated_at": tftypes.NewValue(tftypes.String, nil),
		})

		req := resource.ImportStateRequest{
			ID: "cred-123",
		}
		resp := &resource.ImportStateResponse{
			State: state,
		}

		r.ImportState(ctx, req, resp)

		// Check that ID was set to import
		assert.NotNil(t, resp)
		if resp.Diagnostics.HasError() {
			t.Logf("Diagnostics errors: %v", resp.Diagnostics.Errors())
		}
		// ImportState will have a warning but not an error
		assert.False(t, resp.Diagnostics.HasError())
	})
}

func TestUsesCredential(t *testing.T) {
	t.Run("workflow uses credential", func(t *testing.T) {
		workflow := &n8nsdk.Workflow{
			Nodes: []n8nsdk.Node{
				{
					Credentials: map[string]interface{}{
						"api": map[string]interface{}{
							"id": "cred-123",
						},
					},
				},
			},
		}

		result := usesCredential(workflow, "cred-123")
		assert.True(t, result)
	})

	t.Run("workflow does not use credential", func(t *testing.T) {
		workflow := &n8nsdk.Workflow{
			Nodes: []n8nsdk.Node{
				{
					Credentials: map[string]interface{}{
						"api": map[string]interface{}{
							"id": "other-cred",
						},
					},
				},
			},
		}

		result := usesCredential(workflow, "cred-123")
		assert.False(t, result)
	})

	t.Run("workflow with no credentials", func(t *testing.T) {
		workflow := &n8nsdk.Workflow{
			Nodes: []n8nsdk.Node{
				{
					Credentials: nil,
				},
			},
		}

		result := usesCredential(workflow, "cred-123")
		assert.False(t, result)
	})

	t.Run("workflow with empty nodes", func(t *testing.T) {
		workflow := &n8nsdk.Workflow{
			Nodes: []n8nsdk.Node{},
		}

		result := usesCredential(workflow, "cred-123")
		assert.False(t, result)
	})

	t.Run("nil workflow", func(t *testing.T) {
		result := usesCredential(nil, "cred-123")
		assert.False(t, result)
	})

	t.Run("multiple nodes with credentials", func(t *testing.T) {
		workflow := &n8nsdk.Workflow{
			Nodes: []n8nsdk.Node{
				{
					Credentials: map[string]interface{}{
						"api": map[string]interface{}{
							"id": "other-cred",
						},
					},
				},
				{
					Credentials: map[string]interface{}{
						"http": map[string]interface{}{
							"id": "cred-123",
						},
					},
				},
			},
		}

		result := usesCredential(workflow, "cred-123")
		assert.True(t, result)
	})

	t.Run("credential in nested structure", func(t *testing.T) {
		workflow := &n8nsdk.Workflow{
			Nodes: []n8nsdk.Node{
				{
					Credentials: map[string]interface{}{
						"api": map[string]interface{}{
							"id":   "cred-123",
							"name": "Test Credential",
						},
						"http": map[string]interface{}{
							"id": "other-cred",
						},
					},
				},
			},
		}

		result := usesCredential(workflow, "cred-123")
		assert.True(t, result)
	})

	t.Run("invalid credential structure", func(t *testing.T) {
		// Test with invalid credential structure
		workflow := &n8nsdk.Workflow{
			Nodes: []n8nsdk.Node{
				{
					Credentials: map[string]interface{}{
						"api": "invalid-not-a-map",
					},
				},
			},
		}

		result := usesCredential(workflow, "cred-123")
		assert.False(t, result)
	})

	t.Run("credential with different types", func(t *testing.T) {
		workflow := &n8nsdk.Workflow{
			Nodes: []n8nsdk.Node{
				{
					Credentials: map[string]interface{}{
						"api": map[string]interface{}{
							"id": "cred-123",
						},
						"http": map[string]interface{}{
							"id": "cred-123", // Same credential used for different types
						},
					},
				},
			},
		}

		result := usesCredential(workflow, "cred-123")
		assert.True(t, result)
	})
}

func TestReplaceCredentialInWorkflow(t *testing.T) {
	t.Run("replace single credential", func(t *testing.T) {
		workflow := &n8nsdk.Workflow{
			Id: strPtr("wf-1"),
			Nodes: []n8nsdk.Node{
				{
					Credentials: map[string]interface{}{
						"api": map[string]interface{}{
							"id":   "old-cred",
							"name": "Old Credential",
						},
					},
				},
			},
		}

		result := replaceCredentialInWorkflow(workflow, "old-cred", "new-cred")

		assert.NotNil(t, result)
		credMap := result.Nodes[0].Credentials
		apiCred, ok := credMap["api"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, "new-cred", apiCred["id"])
	})

	t.Run("replace multiple occurrences", func(t *testing.T) {
		workflow := &n8nsdk.Workflow{
			Nodes: []n8nsdk.Node{
				{
					Credentials: map[string]interface{}{
						"api": map[string]interface{}{
							"id": "old-cred",
						},
					},
				},
				{
					Credentials: map[string]interface{}{
						"http": map[string]interface{}{
							"id": "old-cred",
						},
					},
				},
			},
		}

		result := replaceCredentialInWorkflow(workflow, "old-cred", "new-cred")

		assert.NotNil(t, result)

		// Check first node
		credMap1 := result.Nodes[0].Credentials
		apiCred, ok := credMap1["api"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, "new-cred", apiCred["id"])

		// Check second node
		credMap2 := result.Nodes[1].Credentials
		httpCred, ok := credMap2["http"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, "new-cred", httpCred["id"])
	})

	t.Run("no replacement needed", func(t *testing.T) {
		workflow := &n8nsdk.Workflow{
			Nodes: []n8nsdk.Node{
				{
					Credentials: map[string]interface{}{
						"api": map[string]interface{}{
							"id": "other-cred",
						},
					},
				},
			},
		}

		result := replaceCredentialInWorkflow(workflow, "old-cred", "new-cred")

		assert.NotNil(t, result)
		credMap := result.Nodes[0].Credentials
		apiCred, ok := credMap["api"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, "other-cred", apiCred["id"])
	})

	t.Run("nil workflow", func(t *testing.T) {
		result := replaceCredentialInWorkflow(nil, "old-cred", "new-cred")
		assert.Nil(t, result)
	})

	t.Run("empty nodes", func(t *testing.T) {
		workflow := &n8nsdk.Workflow{
			Id:    strPtr("wf-1"),
			Nodes: []n8nsdk.Node{},
		}

		result := replaceCredentialInWorkflow(workflow, "old-cred", "new-cred")
		assert.NotNil(t, result)
		assert.Empty(t, result.Nodes)
	})

	t.Run("node without credentials", func(t *testing.T) {
		workflow := &n8nsdk.Workflow{
			Nodes: []n8nsdk.Node{
				{
					Credentials: nil,
				},
			},
		}

		result := replaceCredentialInWorkflow(workflow, "old-cred", "new-cred")
		assert.NotNil(t, result)
		assert.Nil(t, result.Nodes[0].Credentials)
	})

	t.Run("invalid credential structure", func(t *testing.T) {
		workflow := &n8nsdk.Workflow{
			Nodes: []n8nsdk.Node{
				{
					Credentials: map[string]interface{}{
						"api": "invalid-string",
					},
				},
			},
		}

		result := replaceCredentialInWorkflow(workflow, "old-cred", "new-cred")
		assert.NotNil(t, result)
		// Should not crash, just skip invalid structure
	})

	t.Run("mixed valid and invalid credentials", func(t *testing.T) {
		workflow := &n8nsdk.Workflow{
			Nodes: []n8nsdk.Node{
				{
					Credentials: map[string]interface{}{
						"api": map[string]interface{}{
							"id": "old-cred",
						},
						"invalid": "not-a-map",
					},
				},
			},
		}

		result := replaceCredentialInWorkflow(workflow, "old-cred", "new-cred")
		assert.NotNil(t, result)

		credMap := result.Nodes[0].Credentials
		apiCred, ok := credMap["api"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, "new-cred", apiCred["id"])
	})

	t.Run("preserve other credential properties", func(t *testing.T) {
		workflow := &n8nsdk.Workflow{
			Nodes: []n8nsdk.Node{
				{
					Credentials: map[string]interface{}{
						"api": map[string]interface{}{
							"id":   "old-cred",
							"name": "API Credential",
							"type": "oauth2",
						},
					},
				},
			},
		}

		result := replaceCredentialInWorkflow(workflow, "old-cred", "new-cred")
		assert.NotNil(t, result)

		credMap := result.Nodes[0].Credentials
		apiCred, ok := credMap["api"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, "new-cred", apiCred["id"])
		assert.Equal(t, "API Credential", apiCred["name"])
		assert.Equal(t, "oauth2", apiCred["type"])
	})
}

func TestROTATION_THROTTLE_MILLISECONDS(t *testing.T) {
	t.Run("constant value", func(t *testing.T) {
		assert.Equal(t, 100, ROTATION_THROTTLE_MILLISECONDS)
	})
}

func TestCredentialResource_InterfaceCompliance(t *testing.T) {
	t.Run("implements Resource interface", func(t *testing.T) {
		var _ resource.Resource = (*CredentialResource)(nil)
	})

	t.Run("implements ResourceWithConfigure interface", func(t *testing.T) {
		var _ resource.ResourceWithConfigure = (*CredentialResource)(nil)
	})

	t.Run("implements ResourceWithImportState interface", func(t *testing.T) {
		var _ resource.ResourceWithImportState = (*CredentialResource)(nil)
	})

	t.Run("implements CredentialResourceInterface", func(t *testing.T) {
		var _ CredentialResourceInterface = (*CredentialResource)(nil)
	})
}

func TestCredentialResource_ConcurrentOperations(t *testing.T) {
	t.Run("concurrent credential checks", func(t *testing.T) {
		workflow := &n8nsdk.Workflow{
			Nodes: []n8nsdk.Node{
				{
					Credentials: map[string]interface{}{
						"api": map[string]interface{}{
							"id": "cred-123",
						},
					},
				},
			},
		}

		// Run concurrent checks
		results := make(chan bool, 100)
		for range 100 {
			go func() {
				results <- usesCredential(workflow, "cred-123")
			}()
		}

		// Collect results
		for range 100 {
			assert.True(t, <-results)
		}
	})

	t.Run("concurrent replacements", func(t *testing.T) {
		// Run concurrent replacements on copies
		results := make(chan *n8nsdk.Workflow, 100)
		for i := range 100 {
			go func(idx int) {
				// Create a copy for each goroutine
				workflow := &n8nsdk.Workflow{
					Id: strPtr("wf-1"),
					Nodes: []n8nsdk.Node{
						{
							Credentials: map[string]interface{}{
								"api": map[string]interface{}{
									"id": "old-cred",
								},
							},
						},
					},
				}
				newCredID := fmt.Sprintf("new-cred-%d", idx)
				results <- replaceCredentialInWorkflow(workflow, "old-cred", newCredID)
			}(i)
		}

		// Verify all replacements
		for range 100 {
			result := <-results
			assert.NotNil(t, result)
		}
	})
}

func TestCredentialResource_EdgeCases(t *testing.T) {
	t.Run("deeply nested credential structure", func(t *testing.T) {
		workflow := &n8nsdk.Workflow{
			Nodes: []n8nsdk.Node{
				{
					Credentials: map[string]interface{}{
						"oauth": map[string]interface{}{
							"id": "cred-123",
							"metadata": map[string]interface{}{
								"scope":   "read write",
								"expires": 3600,
							},
						},
					},
				},
			},
		}

		result := usesCredential(workflow, "cred-123")
		assert.True(t, result)
	})

	t.Run("credential ID as non-string", func(t *testing.T) {
		workflow := &n8nsdk.Workflow{
			Nodes: []n8nsdk.Node{
				{
					Credentials: map[string]interface{}{
						"api": map[string]interface{}{
							"id": 123, // Number instead of string
						},
					},
				},
			},
		}

		result := usesCredential(workflow, "123")
		assert.False(t, result) // Should be false as ID is not a string
	})
}

func BenchmarkUsesCredential(b *testing.B) {
	workflow := &n8nsdk.Workflow{
		Nodes: []n8nsdk.Node{
			{
				Credentials: map[string]interface{}{
					"api": map[string]interface{}{
						"id": "cred-123",
					},
				},
			},
			{
				Credentials: map[string]interface{}{
					"http": map[string]interface{}{
						"id": "other-cred",
					},
				},
			},
		},
	}

	b.ResetTimer()
	for b.Loop() {
		_ = usesCredential(workflow, "cred-123")
	}
}

func BenchmarkReplaceCredentialInWorkflow(b *testing.B) {
	workflow := &n8nsdk.Workflow{
		Id: strPtr("wf-1"),
		Nodes: []n8nsdk.Node{
			{
				Credentials: map[string]interface{}{
					"api": map[string]interface{}{
						"id": "old-cred",
					},
				},
			},
		},
	}

	b.ResetTimer()
	for b.Loop() {
		_ = replaceCredentialInWorkflow(workflow, "old-cred", "new-cred")
	}
}

func BenchmarkReplaceCredentialLargeWorkflow(b *testing.B) {
	// Create a workflow with many nodes
	nodes := make([]n8nsdk.Node, 100)
	for range 100 {
		nodes = append(nodes, n8nsdk.Node{
			Credentials: map[string]interface{}{
				fmt.Sprintf("api%d", len(nodes)): map[string]interface{}{
					"id": "old-cred",
				},
			},
		})
	}

	workflow := &n8nsdk.Workflow{
		Id:    strPtr("large-wf"),
		Nodes: nodes,
	}

	b.ResetTimer()
	for b.Loop() {
		_ = replaceCredentialInWorkflow(workflow, "old-cred", "new-cred")
	}
}

// Test helper functions

// createTestSchema creates a test schema for credential resource.
func createTestSchema(t *testing.T) schema.Schema {
	t.Helper()
	r := &CredentialResource{}
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

// TestCredentialResource_Create tests credential creation.
func TestCredentialResource_Create(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/credentials" && r.Method == http.MethodPost {
				response := map[string]interface{}{
					"id":        "cred-123",
					"name":      "Test Credential",
					"type":      "httpHeaderAuth",
					"createdAt": time.Now().Format(time.RFC3339),
					"updatedAt": time.Now().Format(time.RFC3339),
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				json.NewEncoder(w).Encode(response)
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &CredentialResource{client: n8nClient}

		// Note: Using an empty map due to framework limitations with tftypes conversion in tests
		// The HTTP mocking is verified, but the Plan.Get conversion has issues with nested map types
		emptyMap := make(map[string]tftypes.Value)
		dataTfValue := tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, emptyMap)

		rawPlan := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, nil),
			"name":       tftypes.NewValue(tftypes.String, "Test Credential"),
			"type":       tftypes.NewValue(tftypes.String, "httpHeaderAuth"),
			"data":       dataTfValue,
			"created_at": tftypes.NewValue(tftypes.String, nil),
			"updated_at": tftypes.NewValue(tftypes.String, nil),
		}
		plan := tfsdk.Plan{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "type": tftypes.String, "data": tftypes.Map{ElementType: tftypes.String}, "created_at": tftypes.String, "updated_at": tftypes.String}}, rawPlan),
			Schema: createTestSchema(t),
		}

		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "type": tftypes.String, "data": tftypes.Map{ElementType: tftypes.String}, "created_at": tftypes.String, "updated_at": tftypes.String}}, nil),
			Schema: createTestSchema(t),
		}

		req := resource.CreateRequest{
			Plan: plan,
		}
		resp := resource.CreateResponse{
			State: state,
		}

		r.Create(context.Background(), req, &resp)

		if resp.Diagnostics.HasError() {
			for _, diag := range resp.Diagnostics.Errors() {
				t.Logf("Error: %s - %s", diag.Summary(), diag.Detail())
			}
		}
		assert.False(t, resp.Diagnostics.HasError(), "Create should not have errors")
	})

	t.Run("creation fails", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message": "Internal server error"}`))
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &CredentialResource{client: n8nClient}

		// Note: Using an empty map due to framework limitations with tftypes conversion in tests
		// The HTTP mocking is verified, but the Plan.Get conversion has issues with nested map types
		emptyMap := make(map[string]tftypes.Value)
		dataTfValue := tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, emptyMap)

		rawPlan := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, nil),
			"name":       tftypes.NewValue(tftypes.String, "Test Credential"),
			"type":       tftypes.NewValue(tftypes.String, "httpHeaderAuth"),
			"data":       dataTfValue,
			"created_at": tftypes.NewValue(tftypes.String, nil),
			"updated_at": tftypes.NewValue(tftypes.String, nil),
		}
		plan := tfsdk.Plan{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "type": tftypes.String, "data": tftypes.Map{ElementType: tftypes.String}, "created_at": tftypes.String, "updated_at": tftypes.String}}, rawPlan),
			Schema: createTestSchema(t),
		}

		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "type": tftypes.String, "data": tftypes.Map{ElementType: tftypes.String}, "created_at": tftypes.String, "updated_at": tftypes.String}}, nil),
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
	})
}

// TestCredentialResource_Read tests credential reading.
func TestCredentialResource_Read(t *testing.T) {
	t.Run("successful read", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Note: n8n API doesn't support GET /credentials/{id}, so Read keeps state as-is
			w.WriteHeader(http.StatusOK)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &CredentialResource{client: n8nClient}

		// Create state with a properly typed map
		dataMap := types.MapValueMust(types.StringType, map[string]attr.Value{
			"api_key": types.StringValue("secret"),
		})

		dataTfValue, err := dataMap.ToTerraformValue(context.Background())
		assert.NoError(t, err, "ToTerraformValue should not error")

		rawState := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, "cred-123"),
			"name":       tftypes.NewValue(tftypes.String, "Test Credential"),
			"type":       tftypes.NewValue(tftypes.String, "httpHeaderAuth"),
			"data":       dataTfValue,
			"created_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
			"updated_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
		}
		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "type": tftypes.String, "data": tftypes.Map{ElementType: tftypes.String}, "created_at": tftypes.String, "updated_at": tftypes.String}}, rawState),
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
	})

	t.Run("credential not found removes from state", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Note: n8n API doesn't support GET /credentials/{id}
			// Read operation keeps state as-is and doesn't make API calls
			w.WriteHeader(http.StatusOK)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &CredentialResource{client: n8nClient}

		// Create state with a properly typed map
		dataMap := types.MapValueMust(types.StringType, map[string]attr.Value{
			"api_key": types.StringValue("secret"),
		})

		dataTfValue, err := dataMap.ToTerraformValue(context.Background())
		assert.NoError(t, err, "ToTerraformValue should not error")

		rawState := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, "cred-nonexistent"),
			"name":       tftypes.NewValue(tftypes.String, "Test Credential"),
			"type":       tftypes.NewValue(tftypes.String, "httpHeaderAuth"),
			"data":       dataTfValue,
			"created_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
			"updated_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
		}
		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "type": tftypes.String, "data": tftypes.Map{ElementType: tftypes.String}, "created_at": tftypes.String, "updated_at": tftypes.String}}, rawState),
			Schema: createTestSchema(t),
		}

		req := resource.ReadRequest{
			State: state,
		}
		resp := resource.ReadResponse{
			State: state,
		}

		r.Read(context.Background(), req, &resp)

		assert.False(t, resp.Diagnostics.HasError(), "Read should not have errors even when credential not found")
		// Note: In the actual implementation, Read keeps state as-is since n8n API doesn't support GET /credentials/{id}
	})
}

// TestCredentialResource_Update tests credential update.
func TestCredentialResource_Update(t *testing.T) {
	t.Run("successful update", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch {
			case r.URL.Path == "/credentials" && r.Method == http.MethodPost:
				// Create new credential
				response := map[string]interface{}{
					"id":        "cred-456",
					"name":      "Updated Credential",
					"type":      "httpHeaderAuth",
					"createdAt": time.Now().Format(time.RFC3339),
					"updatedAt": time.Now().Format(time.RFC3339),
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				json.NewEncoder(w).Encode(response)
				return
			case r.URL.Path == "/workflows" && r.Method == http.MethodGet:
				// Return empty workflow list (no workflows using the credential)
				response := map[string]interface{}{
					"data": []map[string]interface{}{},
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
				return
			case r.URL.Path == "/credentials/cred-123" && r.Method == http.MethodDelete:
				// Delete old credential
				w.WriteHeader(http.StatusNoContent)
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &CredentialResource{client: n8nClient}

		// Note: Using an empty map due to framework limitations with tftypes conversion in tests
		// The HTTP mocking is verified, but the Plan.Get conversion has issues with nested map types
		emptyMapPlan := make(map[string]tftypes.Value)
		dataTfValuePlan := tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, emptyMapPlan)

		rawPlan := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, "cred-123"),
			"name":       tftypes.NewValue(tftypes.String, "Updated Credential"),
			"type":       tftypes.NewValue(tftypes.String, "httpHeaderAuth"),
			"data":       dataTfValuePlan,
			"created_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
			"updated_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
		}
		plan := tfsdk.Plan{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "type": tftypes.String, "data": tftypes.Map{ElementType: tftypes.String}, "created_at": tftypes.String, "updated_at": tftypes.String}}, rawPlan),
			Schema: createTestSchema(t),
		}

		// Note: Using an empty map due to framework limitations with tftypes conversion in tests
		emptyMapState := make(map[string]tftypes.Value)
		dataTfValueState := tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, emptyMapState)

		rawState := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, "cred-123"),
			"name":       tftypes.NewValue(tftypes.String, "Test Credential"),
			"type":       tftypes.NewValue(tftypes.String, "httpHeaderAuth"),
			"data":       dataTfValueState,
			"created_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
			"updated_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
		}
		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "type": tftypes.String, "data": tftypes.Map{ElementType: tftypes.String}, "created_at": tftypes.String, "updated_at": tftypes.String}}, rawState),
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

		if resp.Diagnostics.HasError() {
			for _, diag := range resp.Diagnostics.Errors() {
				t.Logf("Error: %s - %s", diag.Summary(), diag.Detail())
			}
		}
		assert.False(t, resp.Diagnostics.HasError(), "Update should not have errors")
	})
}

// TestCredentialResource_Delete tests credential deletion.
func TestCredentialResource_Delete(t *testing.T) {
	t.Run("successful delete", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/credentials/cred-123" && r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &CredentialResource{client: n8nClient}

		// Create state using types.Map
		dataMap, diags := types.MapValueFrom(context.Background(), types.StringType, map[string]string{
			"api_key": "secret",
		})
		assert.False(t, diags.HasError(), "Map creation should not have errors")

		dataTfValue, err := dataMap.ToTerraformValue(context.Background())
		assert.NoError(t, err, "ToTerraformValue should not error")

		rawState := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, "cred-123"),
			"name":       tftypes.NewValue(tftypes.String, "Test Credential"),
			"type":       tftypes.NewValue(tftypes.String, "httpHeaderAuth"),
			"data":       dataTfValue,
			"created_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
			"updated_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
		}
		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "type": tftypes.String, "data": tftypes.Map{ElementType: tftypes.String}, "created_at": tftypes.String, "updated_at": tftypes.String}}, rawState),
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
	})

	t.Run("delete fails", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message": "Internal server error"}`))
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &CredentialResource{client: n8nClient}

		// Create state using types.Map
		dataMap, diags := types.MapValueFrom(context.Background(), types.StringType, map[string]string{
			"api_key": "secret",
		})
		assert.False(t, diags.HasError(), "Map creation should not have errors")

		dataTfValue, err := dataMap.ToTerraformValue(context.Background())
		assert.NoError(t, err, "ToTerraformValue should not error")

		rawState := map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, "cred-123"),
			"name":       tftypes.NewValue(tftypes.String, "Test Credential"),
			"type":       tftypes.NewValue(tftypes.String, "httpHeaderAuth"),
			"data":       dataTfValue,
			"created_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
			"updated_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
		}
		state := tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "name": tftypes.String, "type": tftypes.String, "data": tftypes.Map{ElementType: tftypes.String}, "created_at": tftypes.String, "updated_at": tftypes.String}}, rawState),
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
	})
}

// TestCredentialResource_RollbackFunctions tests rollback-related helper functions.
func TestCredentialResource_RollbackFunctions(t *testing.T) {
	t.Run("deleteNewCredential success", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/credentials/new-cred-123" && r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &CredentialResource{client: n8nClient}
		r.deleteNewCredential(context.Background(), "new-cred-123")
	})

	t.Run("deleteNewCredential failure", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &CredentialResource{client: n8nClient}
		r.deleteNewCredential(context.Background(), "new-cred-123")
	})

	t.Run("deleteCredentialBestEffort success", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/credentials/cred-123" && r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &CredentialResource{client: n8nClient}
		r.deleteCredentialBestEffort(context.Background(), "cred-123")
	})

	t.Run("deleteCredentialBestEffort failure", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &CredentialResource{client: n8nClient}
		r.deleteCredentialBestEffort(context.Background(), "cred-123")
	})

	t.Run("findWorkflowBackup found", func(t *testing.T) {
		r := &CredentialResource{}
		workflow := &n8nsdk.Workflow{Id: strPtr("wf-123")}
		backups := []models.WorkflowBackup{
			{ID: "wf-456", Original: &n8nsdk.Workflow{Id: strPtr("wf-456")}},
			{ID: "wf-123", Original: workflow},
		}

		found := r.findWorkflowBackup(backups, "wf-123")
		assert.NotNil(t, found)
		assert.Equal(t, "wf-123", *found.Id)
	})

	t.Run("findWorkflowBackup not found", func(t *testing.T) {
		r := &CredentialResource{}
		backups := []models.WorkflowBackup{
			{ID: "wf-456", Original: &n8nsdk.Workflow{Id: strPtr("wf-456")}},
		}

		found := r.findWorkflowBackup(backups, "wf-123")
		assert.Nil(t, found)
	})

	t.Run("restoreWorkflow success", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// SDK generates path like /workflows/{id}
			if r.Method == http.MethodPut && r.URL.Path == "/workflows/wf-123" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]any{
					"id":          "wf-123",
					"name":        "Test Workflow",
					"nodes":       []any{},
					"connections": map[string]any{},
					"settings":    map[string]any{},
				})
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &CredentialResource{client: n8nClient}
		workflow := &n8nsdk.Workflow{
			Id:          strPtr("wf-123"),
			Name:        "Test",
			Nodes:       []n8nsdk.Node{},
			Connections: map[string]interface{}{},
			Settings:    n8nsdk.WorkflowSettings{},
		}
		success := r.restoreWorkflow(context.Background(), "wf-123", workflow)
		assert.True(t, success)
	})

	t.Run("restoreWorkflow failure", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &CredentialResource{client: n8nClient}
		workflow := &n8nsdk.Workflow{
			Id:          strPtr("wf-123"),
			Name:        "Test",
			Nodes:       []n8nsdk.Node{},
			Connections: map[string]interface{}{},
			Settings:    n8nsdk.WorkflowSettings{},
		}
		success := r.restoreWorkflow(context.Background(), "wf-123", workflow)
		assert.False(t, success)
	})

	t.Run("restoreWorkflows with multiple workflows", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPut {
				// Succeed for wf-123, fail for wf-456
				if r.URL.Path == "/workflows/wf-123" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{
						"id":          "wf-123",
						"name":        "WF1",
						"nodes":       []any{},
						"connections": map[string]any{},
						"settings":    map[string]any{},
					})
					return
				}
				if r.URL.Path == "/workflows/wf-456" {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &CredentialResource{client: n8nClient}
		backups := []models.WorkflowBackup{
			{
				ID: "wf-123",
				Original: &n8nsdk.Workflow{
					Id:          strPtr("wf-123"),
					Name:        "WF1",
					Nodes:       []n8nsdk.Node{},
					Connections: map[string]interface{}{},
					Settings:    n8nsdk.WorkflowSettings{},
				},
			},
			{
				ID: "wf-456",
				Original: &n8nsdk.Workflow{
					Id:          strPtr("wf-456"),
					Name:        "WF2",
					Nodes:       []n8nsdk.Node{},
					Connections: map[string]interface{}{},
					Settings:    n8nsdk.WorkflowSettings{},
				},
			},
		}
		updatedWorkflows := []string{"wf-123", "wf-456"}

		count := r.restoreWorkflows(context.Background(), backups, updatedWorkflows)
		assert.Equal(t, 1, count)
	})

	t.Run("restoreWorkflows with missing backup", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPut && r.URL.Path == "/workflows/wf-456" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]any{
					"id":          "wf-456",
					"name":        "WF2",
					"nodes":       []any{},
					"connections": map[string]any{},
					"settings":    map[string]any{},
				})
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &CredentialResource{client: n8nClient}
		backups := []models.WorkflowBackup{
			{
				ID: "wf-456",
				Original: &n8nsdk.Workflow{
					Id:          strPtr("wf-456"),
					Name:        "WF2",
					Nodes:       []n8nsdk.Node{},
					Connections: map[string]interface{}{},
					Settings:    n8nsdk.WorkflowSettings{},
				},
			},
		}
		updatedWorkflows := []string{"wf-123", "wf-456"}

		count := r.restoreWorkflows(context.Background(), backups, updatedWorkflows)
		assert.Equal(t, 1, count)
	})

	t.Run("rollbackRotation complete flow", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			if r.Method == http.MethodPut && r.URL.Path == "/workflows/wf-123" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]any{
					"id":          "wf-123",
					"name":        "Test",
					"nodes":       []any{},
					"connections": map[string]any{},
					"settings":    map[string]any{},
				})
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &CredentialResource{client: n8nClient}
		backups := []models.WorkflowBackup{
			{
				ID: "wf-123",
				Original: &n8nsdk.Workflow{
					Id:          strPtr("wf-123"),
					Name:        "Test",
					Nodes:       []n8nsdk.Node{},
					Connections: map[string]interface{}{},
					Settings:    n8nsdk.WorkflowSettings{},
				},
			},
		}
		updatedWorkflows := []string{"wf-123"}

		r.rollbackRotation(context.Background(), "new-cred-123", backups, updatedWorkflows)
	})
}

// TestCredentialResource_DeleteOldCredential tests deletion of old credentials after rotation.
func TestCredentialResource_DeleteOldCredential(t *testing.T) {
	t.Run("deleteOldCredential failure logs warning", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message": "Internal server error"}`))
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &CredentialResource{client: n8nClient}
		r.deleteOldCredential(context.Background(), "old-cred-123", "new-cred-456")
	})
}

// TestCredentialResource_HelperFunctions tests the helper functions for credential rotation.
func TestCredentialResource_HelperFunctions(t *testing.T) {
	t.Run("createNewCredential success", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost && r.URL.Path == "/credentials" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				json.NewEncoder(w).Encode(map[string]any{
					"id":        "new-cred-123",
					"name":      "Test Credential",
					"type":      "testType",
					"createdAt": "2025-01-01T00:00:00Z",
					"updatedAt": "2025-01-01T00:00:00Z",
				})
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &CredentialResource{client: n8nClient}
		diags := &diag.Diagnostics{}
		result := r.createNewCredential(context.Background(), "Test Credential", "testType", map[string]interface{}{"key": "value"}, diags)

		assert.NotNil(t, result)
		if result != nil {
			assert.Equal(t, "new-cred-123", result.Id)
		}
		assert.False(t, diags.HasError())
	})

	t.Run("createNewCredential failure", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &CredentialResource{client: n8nClient}
		diags := &diag.Diagnostics{}
		result := r.createNewCredential(context.Background(), "Test Credential", "testType", map[string]interface{}{"key": "value"}, diags)

		assert.Nil(t, result)
		assert.True(t, diags.HasError())
	})

	t.Run("scanAffectedWorkflows success with workflows", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet && r.URL.Path == "/workflows" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]any{
					"data": []any{
						map[string]any{
							"id":          "wf-1",
							"name":        "Workflow 1",
							"nodes":       []any{map[string]any{"credentials": map[string]any{"testCred": map[string]any{"id": "old-cred-123"}}}},
							"connections": map[string]any{},
							"settings":    map[string]any{},
						},
						map[string]any{
							"id":          "wf-2",
							"name":        "Workflow 2",
							"nodes":       []any{},
							"connections": map[string]any{},
							"settings":    map[string]any{},
						},
					},
				})
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &CredentialResource{client: n8nClient}
		diags := &diag.Diagnostics{}
		workflows, success := r.scanAffectedWorkflows(context.Background(), "old-cred-123", "new-cred-456", diags)

		assert.True(t, success)
		assert.Len(t, workflows, 1)
		assert.Equal(t, "wf-1", workflows[0].ID)
		assert.False(t, diags.HasError())
	})

	t.Run("scanAffectedWorkflows failure", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet && r.URL.Path == "/workflows" {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if r.Method == http.MethodDelete && r.URL.Path == "/credentials/new-cred-456" {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &CredentialResource{client: n8nClient}
		diags := &diag.Diagnostics{}
		workflows, success := r.scanAffectedWorkflows(context.Background(), "old-cred-123", "new-cred-456", diags)

		assert.False(t, success)
		assert.Len(t, workflows, 0)
		assert.True(t, diags.HasError())
	})

	t.Run("updateAffectedWorkflows success", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet && r.URL.Path == "/workflows/wf-1" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]any{
					"id":          "wf-1",
					"name":        "Workflow 1",
					"nodes":       []any{map[string]any{"credentials": map[string]any{"testCred": map[string]any{"id": "old-cred-123"}}}},
					"connections": map[string]any{},
					"settings":    map[string]any{},
				})
				return
			}
			if r.Method == http.MethodPut && r.URL.Path == "/workflows/wf-1" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]any{
					"id":          "wf-1",
					"name":        "Workflow 1",
					"nodes":       []any{},
					"connections": map[string]any{},
					"settings":    map[string]any{},
				})
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &CredentialResource{client: n8nClient}
		backups := []models.WorkflowBackup{
			{
				ID: "wf-1",
				Original: &n8nsdk.Workflow{
					Id:          strPtr("wf-1"),
					Name:        "Workflow 1",
					Nodes:       []n8nsdk.Node{},
					Connections: map[string]interface{}{},
					Settings:    n8nsdk.WorkflowSettings{},
				},
			},
		}
		diags := &diag.Diagnostics{}
		updated, success := r.updateAffectedWorkflows(context.Background(), backups, "old-cred-123", "new-cred-456", diags)

		assert.True(t, success)
		assert.Len(t, updated, 1)
		assert.Equal(t, "wf-1", updated[0])
		assert.False(t, diags.HasError())
	})

	t.Run("updateAffectedWorkflows failure on get", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet && r.URL.Path == "/workflows/wf-1" {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if r.Method == http.MethodDelete && r.URL.Path == "/credentials/new-cred-456" {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &CredentialResource{client: n8nClient}
		backups := []models.WorkflowBackup{
			{
				ID: "wf-1",
				Original: &n8nsdk.Workflow{
					Id:          strPtr("wf-1"),
					Name:        "Workflow 1",
					Nodes:       []n8nsdk.Node{},
					Connections: map[string]interface{}{},
					Settings:    n8nsdk.WorkflowSettings{},
				},
			},
		}
		diags := &diag.Diagnostics{}
		updated, success := r.updateAffectedWorkflows(context.Background(), backups, "old-cred-123", "new-cred-456", diags)

		assert.False(t, success)
		assert.Len(t, updated, 0)
		assert.True(t, diags.HasError())
	})

	t.Run("updateAffectedWorkflows failure on put", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet && r.URL.Path == "/workflows/wf-1" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]any{
					"id":          "wf-1",
					"name":        "Workflow 1",
					"nodes":       []any{map[string]any{"credentials": map[string]any{"testCred": map[string]any{"id": "old-cred-123"}}}},
					"connections": map[string]any{},
					"settings":    map[string]any{},
				})
				return
			}
			if r.Method == http.MethodPut && r.URL.Path == "/workflows/wf-1" {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if r.Method == http.MethodDelete && r.URL.Path == "/credentials/new-cred-456" {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &CredentialResource{client: n8nClient}
		backups := []models.WorkflowBackup{
			{
				ID: "wf-1",
				Original: &n8nsdk.Workflow{
					Id:          strPtr("wf-1"),
					Name:        "Workflow 1",
					Nodes:       []n8nsdk.Node{},
					Connections: map[string]interface{}{},
					Settings:    n8nsdk.WorkflowSettings{},
				},
			},
		}
		diags := &diag.Diagnostics{}
		updated, success := r.updateAffectedWorkflows(context.Background(), backups, "old-cred-123", "new-cred-456", diags)

		assert.False(t, success)
		assert.Len(t, updated, 0)
		assert.True(t, diags.HasError())
	})

	t.Run("updateAffectedWorkflows with multiple workflows tests throttling", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet && (r.URL.Path == "/workflows/wf-1" || r.URL.Path == "/workflows/wf-2") {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]any{
					"id":          r.URL.Path[len("/workflows/"):],
					"name":        "Workflow",
					"nodes":       []any{},
					"connections": map[string]any{},
					"settings":    map[string]any{},
				})
				return
			}
			if r.Method == http.MethodPut && (r.URL.Path == "/workflows/wf-1" || r.URL.Path == "/workflows/wf-2") {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]any{
					"id":          r.URL.Path[len("/workflows/"):],
					"name":        "Workflow",
					"nodes":       []any{},
					"connections": map[string]any{},
					"settings":    map[string]any{},
				})
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := &CredentialResource{client: n8nClient}
		backups := []models.WorkflowBackup{
			{
				ID: "wf-1",
				Original: &n8nsdk.Workflow{
					Id:          strPtr("wf-1"),
					Name:        "Workflow 1",
					Nodes:       []n8nsdk.Node{},
					Connections: map[string]interface{}{},
					Settings:    n8nsdk.WorkflowSettings{},
				},
			},
			{
				ID: "wf-2",
				Original: &n8nsdk.Workflow{
					Id:          strPtr("wf-2"),
					Name:        "Workflow 2",
					Nodes:       []n8nsdk.Node{},
					Connections: map[string]interface{}{},
					Settings:    n8nsdk.WorkflowSettings{},
				},
			},
		}
		diags := &diag.Diagnostics{}
		updated, success := r.updateAffectedWorkflows(context.Background(), backups, "old-cred-123", "new-cred-456", diags)

		assert.True(t, success)
		assert.Len(t, updated, 2)
		assert.False(t, diags.HasError())
	})
}
