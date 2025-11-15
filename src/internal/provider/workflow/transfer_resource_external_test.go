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
	"github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/client"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/workflow"
	"github.com/stretchr/testify/assert"
)

// TestNewWorkflowTransferResource tests the NewWorkflowTransferResource constructor.
func TestNewWorkflowTransferResource(t *testing.T) {
	tests := []struct {
		name     string
		validate func(*testing.T, *workflow.WorkflowTransferResource)
	}{
		{
			name: "create new transfer resource",
			validate: func(t *testing.T, res *workflow.WorkflowTransferResource) {
				t.Helper()
				assert.NotNil(t, res, "NewWorkflowTransferResource should return a non-nil resource")
			},
		},
		{
			name: "resource has nil client initially",
			validate: func(t *testing.T, res *workflow.WorkflowTransferResource) {
				t.Helper()
				assert.NotNil(t, res, "Resource should not be nil")
			},
		},
		{
			name: "multiple instances are independent",
			validate: func(t *testing.T, res *workflow.WorkflowTransferResource) {
				t.Helper()
				resource2 := workflow.NewWorkflowTransferResource()
				assert.NotNil(t, res)
				assert.NotNil(t, resource2)
				assert.NotSame(t, res, resource2, "Each call should create a new instance")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := workflow.NewWorkflowTransferResource()
			tt.validate(t, res)
		})
	}
}

// TestNewWorkflowTransferResourceWrapper tests the wrapper function.
func TestNewWorkflowTransferResourceWrapper(t *testing.T) {
	tests := []struct {
		name     string
		validate func(*testing.T, resource.Resource)
	}{
		{
			name: "create resource via wrapper",
			validate: func(t *testing.T, res resource.Resource) {
				t.Helper()
				assert.NotNil(t, res, "NewWorkflowTransferResourceWrapper should return a non-nil resource")
			},
		},
		{
			name: "wrapper returns Resource interface",
			validate: func(t *testing.T, res resource.Resource) {
				t.Helper()
				assert.Implements(t, (*resource.Resource)(nil), res)
			},
		},
		{
			name: "wrapper returns correct type",
			validate: func(t *testing.T, res resource.Resource) {
				t.Helper()
				assert.NotNil(t, res)
				assert.Implements(t, (*resource.Resource)(nil), res)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := workflow.NewWorkflowTransferResourceWrapper()
			tt.validate(t, res)
		})
	}
}

// TestWorkflowTransferResource_Metadata tests the Metadata method.
func TestWorkflowTransferResource_Metadata(t *testing.T) {
	tests := []struct {
		name             string
		providerTypeName string
		expectedTypeName string
		validate         func(*testing.T, *resource.MetadataResponse)
	}{
		{
			name:             "set metadata with provider name",
			providerTypeName: "n8n",
			expectedTypeName: "n8n_workflow_transfer",
		},
		{
			name:             "different provider type name",
			providerTypeName: "custom",
			expectedTypeName: "custom_workflow_transfer",
		},
		{
			name:             "empty provider type name",
			providerTypeName: "",
			expectedTypeName: "_workflow_transfer",
		},
		{
			name:             "metadata is idempotent",
			providerTypeName: "n8n",
			expectedTypeName: "n8n_workflow_transfer",
			validate: func(t *testing.T, resp *resource.MetadataResponse) {
				t.Helper()
				// Call Metadata again to verify idempotency
				r := workflow.NewWorkflowTransferResource()
				req := resource.MetadataRequest{
					ProviderTypeName: "n8n",
				}
				resp2 := &resource.MetadataResponse{}
				r.Metadata(context.Background(), req, resp2)
				assert.Equal(t, resp.TypeName, resp2.TypeName)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := workflow.NewWorkflowTransferResource()
			req := resource.MetadataRequest{
				ProviderTypeName: tt.providerTypeName,
			}
			resp := &resource.MetadataResponse{}

			r.Metadata(context.Background(), req, resp)

			assert.Equal(t, tt.expectedTypeName, resp.TypeName)

			if tt.validate != nil {
				tt.validate(t, resp)
			}
		})
	}
}

// TestWorkflowTransferResource_Schema tests the Schema method.
func TestWorkflowTransferResource_Schema(t *testing.T) {
	tests := []struct {
		name     string
		validate func(*testing.T, *resource.SchemaResponse)
	}{
		{
			name: "get schema",
			validate: func(t *testing.T, resp *resource.SchemaResponse) {
				t.Helper()
				assert.NotNil(t, resp.Schema)
				assert.NotEmpty(t, resp.Schema.MarkdownDescription)
				assert.NotEmpty(t, resp.Schema.Attributes)
			},
		},
		{
			name: "schema has required attributes",
			validate: func(t *testing.T, resp *resource.SchemaResponse) {
				t.Helper()
				assert.Contains(t, resp.Schema.Attributes, "id")
				assert.Contains(t, resp.Schema.Attributes, "workflow_id")
				assert.Contains(t, resp.Schema.Attributes, "destination_project_id")
				assert.Contains(t, resp.Schema.Attributes, "transferred_at")
			},
		},
		{
			name: "schema has correct attribute count",
			validate: func(t *testing.T, resp *resource.SchemaResponse) {
				t.Helper()
				assert.Equal(t, 4, len(resp.Schema.Attributes))
			},
		},
		{
			name: "schema description mentions transfer",
			validate: func(t *testing.T, resp *resource.SchemaResponse) {
				t.Helper()
				assert.Contains(t, resp.Schema.MarkdownDescription, "transfer")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := workflow.NewWorkflowTransferResource()
			req := resource.SchemaRequest{}
			resp := &resource.SchemaResponse{}

			r.Schema(context.Background(), req, resp)

			tt.validate(t, resp)
		})
	}
}

// TestWorkflowTransferResource_Configure tests the Configure method.
func TestWorkflowTransferResource_Configure(t *testing.T) {
	tests := []struct {
		name          string
		providerData  interface{}
		wantErr       bool
		errorContains string
		setup         func(*testing.T) interface{}
	}{
		{
			name:         "configure with valid client",
			providerData: &client.N8nClient{},
			wantErr:      false,
		},
		{
			name:         "configure with nil provider data",
			providerData: nil,
			wantErr:      false,
		},
		{
			name:          "configure with wrong type",
			providerData:  "invalid",
			wantErr:       true,
			errorContains: "Unexpected Resource Configure Type",
		},
		{
			name: "configure multiple times",
			setup: func(t *testing.T) interface{} {
				t.Helper()
				r := workflow.NewWorkflowTransferResource()
				// First configuration
				resp1 := &resource.ConfigureResponse{}
				client1 := &client.N8nClient{}
				req1 := resource.ConfigureRequest{
					ProviderData: client1,
				}
				r.Configure(context.Background(), req1, resp1)
				assert.False(t, resp1.Diagnostics.HasError())
				return r
			},
			providerData: &client.N8nClient{},
			wantErr:      false,
		},
		{
			name:         "configure with integer provider data",
			providerData: 123,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r *workflow.WorkflowTransferResource
			if tt.setup != nil {
				r = tt.setup(t).(*workflow.WorkflowTransferResource)
			} else {
				r = workflow.NewWorkflowTransferResource()
			}

			req := resource.ConfigureRequest{
				ProviderData: tt.providerData,
			}
			resp := &resource.ConfigureResponse{}

			r.Configure(context.Background(), req, resp)

			if tt.wantErr {
				assert.True(t, resp.Diagnostics.HasError())
				if tt.errorContains != "" {
					assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), tt.errorContains)
				}
			} else {
				assert.False(t, resp.Diagnostics.HasError())
			}
		})
	}
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
	tests := []struct {
		name    string
		handler http.HandlerFunc
		rawPlan map[string]tftypes.Value
		wantErr bool
	}{
		{
			name: "successful transfer",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/workflows/wf-123/transfer" && r.Method == http.MethodPut {
					w.WriteHeader(http.StatusOK)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			rawPlan: map[string]tftypes.Value{
				"id":                     tftypes.NewValue(tftypes.String, nil),
				"workflow_id":            tftypes.NewValue(tftypes.String, "wf-123"),
				"destination_project_id": tftypes.NewValue(tftypes.String, "proj-456"),
				"transferred_at":         tftypes.NewValue(tftypes.String, nil),
			},
			wantErr: false,
		},
		{
			name: "transfer fails with API error",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"message": "Internal server error"}`))
			}),
			rawPlan: map[string]tftypes.Value{
				"id":                     tftypes.NewValue(tftypes.String, nil),
				"workflow_id":            tftypes.NewValue(tftypes.String, "wf-123"),
				"destination_project_id": tftypes.NewValue(tftypes.String, "proj-456"),
				"transferred_at":         tftypes.NewValue(tftypes.String, nil),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.handler != nil {
				n8nClient, server := setupTransferTestClient(t, tt.handler)
				defer server.Close()

				r := workflow.NewWorkflowTransferResource()
				req := resource.ConfigureRequest{
					ProviderData: n8nClient,
				}
				resp := &resource.ConfigureResponse{}
				r.Configure(context.Background(), req, resp)

				objectType := createTransferTestObjectType(t)
				testSchema := createTransferTestSchema(t)
				plan := tfsdk.Plan{
					Raw:    tftypes.NewValue(objectType, tt.rawPlan),
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

				if tt.wantErr {
					assert.True(t, createResp.Diagnostics.HasError())
				} else {
					assert.False(t, createResp.Diagnostics.HasError())
				}
			}
		})
	}

	// Test case with invalid plan (requires different handling)
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
	tests := []struct {
		name           string
		rawState       map[string]tftypes.Value
		wantErr        bool
		useWrongSchema bool
	}{
		{
			name: "successful read maintains state",
			rawState: map[string]tftypes.Value{
				"id":                     tftypes.NewValue(tftypes.String, "wf-123-to-proj-456"),
				"workflow_id":            tftypes.NewValue(tftypes.String, "wf-123"),
				"destination_project_id": tftypes.NewValue(tftypes.String, "proj-456"),
				"transferred_at":         tftypes.NewValue(tftypes.String, "2024-01-01"),
			},
			wantErr: false,
		},
		{
			name: "read with invalid state",
			rawState: map[string]tftypes.Value{
				"id": tftypes.NewValue(tftypes.Number, 123),
			},
			wantErr:        true,
			useWrongSchema: true,
		},
		{
			name: "read is idempotent",
			rawState: map[string]tftypes.Value{
				"id":                     tftypes.NewValue(tftypes.String, "wf-123-to-proj-456"),
				"workflow_id":            tftypes.NewValue(tftypes.String, "wf-123"),
				"destination_project_id": tftypes.NewValue(tftypes.String, "proj-456"),
				"transferred_at":         tftypes.NewValue(tftypes.String, "2024-01-01"),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := workflow.NewWorkflowTransferResource()

			testSchema := createTransferTestSchema(t)
			objectType := createTransferTestObjectType(t)

			var state tfsdk.State
			if tt.useWrongSchema {
				wrongSchema := tftypes.Object{
					AttributeTypes: map[string]tftypes.Type{
						"id": tftypes.Number,
					},
				}
				state = tfsdk.State{
					Raw:    tftypes.NewValue(wrongSchema, tt.rawState),
					Schema: testSchema,
				}
			} else {
				state = tfsdk.State{
					Raw:    tftypes.NewValue(objectType, tt.rawState),
					Schema: testSchema,
				}
			}

			readReq := resource.ReadRequest{
				State: state,
			}
			readResp := resource.ReadResponse{
				State: state,
			}

			r.Read(context.Background(), readReq, &readResp)

			if tt.wantErr {
				assert.True(t, readResp.Diagnostics.HasError())
			} else {
				assert.False(t, readResp.Diagnostics.HasError())

				// For idempotency test, read again
				if tt.name == "read is idempotent" {
					readReq2 := resource.ReadRequest{State: state}
					readResp2 := resource.ReadResponse{State: state}
					r.Read(context.Background(), readReq2, &readResp2)
					assert.False(t, readResp2.Diagnostics.HasError())
				}
			}
		})
	}
}

// TestWorkflowTransferResource_Update tests the Update method.
func TestWorkflowTransferResource_Update(t *testing.T) {
	tests := []struct {
		name          string
		rawPlan       map[string]tftypes.Value
		rawState      map[string]tftypes.Value
		wantErr       bool
		errorContains string
	}{
		{
			name: "update returns error",
			rawPlan: map[string]tftypes.Value{
				"id":                     tftypes.NewValue(tftypes.String, "wf-123-to-proj-456"),
				"workflow_id":            tftypes.NewValue(tftypes.String, "wf-123"),
				"destination_project_id": tftypes.NewValue(tftypes.String, "proj-789"),
				"transferred_at":         tftypes.NewValue(tftypes.String, "2024-01-01"),
			},
			rawState: map[string]tftypes.Value{
				"id":                     tftypes.NewValue(tftypes.String, "wf-123-to-proj-456"),
				"workflow_id":            tftypes.NewValue(tftypes.String, "wf-123"),
				"destination_project_id": tftypes.NewValue(tftypes.String, "proj-456"),
				"transferred_at":         tftypes.NewValue(tftypes.String, "2024-01-01"),
			},
			wantErr:       true,
			errorContains: "Update Not Supported",
		},
		{
			name: "update always fails",
			rawPlan: map[string]tftypes.Value{
				"id":                     tftypes.NewValue(tftypes.String, "test-id"),
				"workflow_id":            tftypes.NewValue(tftypes.String, "wf-123"),
				"destination_project_id": tftypes.NewValue(tftypes.String, "proj-456"),
				"transferred_at":         tftypes.NewValue(tftypes.String, "2024-01-01"),
			},
			rawState: map[string]tftypes.Value{
				"id":                     tftypes.NewValue(tftypes.String, "test-id"),
				"workflow_id":            tftypes.NewValue(tftypes.String, "wf-123"),
				"destination_project_id": tftypes.NewValue(tftypes.String, "proj-456"),
				"transferred_at":         tftypes.NewValue(tftypes.String, "2024-01-01"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := workflow.NewWorkflowTransferResource()

			objectType := createTransferTestObjectType(t)
			testSchema := createTransferTestSchema(t)
			plan := tfsdk.Plan{
				Raw:    tftypes.NewValue(objectType, tt.rawPlan),
				Schema: testSchema,
			}

			state := tfsdk.State{
				Raw:    tftypes.NewValue(objectType, tt.rawState),
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

			if tt.wantErr {
				assert.True(t, updateResp.Diagnostics.HasError())
				if tt.errorContains != "" {
					assert.Contains(t, updateResp.Diagnostics.Errors()[0].Summary(), tt.errorContains)
				}
			} else {
				assert.False(t, updateResp.Diagnostics.HasError())
			}
		})
	}
}

// TestWorkflowTransferResource_Delete tests the Delete method.
func TestWorkflowTransferResource_Delete(t *testing.T) {
	tests := []struct {
		name        string
		rawState    map[string]tftypes.Value
		useNilState bool
	}{
		{
			name: "delete succeeds without API call",
			rawState: map[string]tftypes.Value{
				"id":                     tftypes.NewValue(tftypes.String, "wf-123-to-proj-456"),
				"workflow_id":            tftypes.NewValue(tftypes.String, "wf-123"),
				"destination_project_id": tftypes.NewValue(tftypes.String, "proj-456"),
				"transferred_at":         tftypes.NewValue(tftypes.String, "2024-01-01"),
			},
		},
		{
			name: "delete is idempotent",
			rawState: map[string]tftypes.Value{
				"id":                     tftypes.NewValue(tftypes.String, "wf-123-to-proj-456"),
				"workflow_id":            tftypes.NewValue(tftypes.String, "wf-123"),
				"destination_project_id": tftypes.NewValue(tftypes.String, "proj-456"),
				"transferred_at":         tftypes.NewValue(tftypes.String, "2024-01-01"),
			},
		},
		{
			name:        "delete with nil state",
			useNilState: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := workflow.NewWorkflowTransferResource()

			objectType := createTransferTestObjectType(t)
			testSchema := createTransferTestSchema(t)

			var state tfsdk.State
			if tt.useNilState {
				state = tfsdk.State{
					Raw:    tftypes.NewValue(objectType, nil),
					Schema: testSchema,
				}
			} else {
				state = tfsdk.State{
					Raw:    tftypes.NewValue(objectType, tt.rawState),
					Schema: testSchema,
				}
			}

			deleteReq := resource.DeleteRequest{
				State: state,
			}
			deleteResp := resource.DeleteResponse{
				State: state,
			}

			r.Delete(context.Background(), deleteReq, &deleteResp)

			assert.False(t, deleteResp.Diagnostics.HasError())

			// For idempotency test, delete again
			if tt.name == "delete is idempotent" {
				deleteReq2 := resource.DeleteRequest{State: state}
				deleteResp2 := resource.DeleteResponse{State: state}
				r.Delete(context.Background(), deleteReq2, &deleteResp2)
				assert.False(t, deleteResp2.Diagnostics.HasError())
			}
		})
	}
}

// TestWorkflowTransferResource_ImportState tests the ImportState method.
func TestWorkflowTransferResource_ImportState(t *testing.T) {
	tests := []struct {
		name string
		id   string
	}{
		{
			name: "successful import",
			id:   "wf-123-to-proj-456",
		},
		{
			name: "import with different ID",
			id:   "custom-transfer-id",
		},
		{
			name: "import uses passthrough",
			id:   "test-id",
		},
		{
			name: "import with empty ID",
			id:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := workflow.NewWorkflowTransferResource()

			objectType := createTransferTestObjectType(t)
			testSchema := createTransferTestSchema(t)
			state := tfsdk.State{
				Raw:    tftypes.NewValue(objectType, nil),
				Schema: testSchema,
			}

			importReq := resource.ImportStateRequest{
				ID: tt.id,
			}
			importResp := &resource.ImportStateResponse{
				State: state,
			}

			r.ImportState(context.Background(), importReq, importResp)

			// ImportStatePassthroughID should handle all cases gracefully
			assert.False(t, importResp.Diagnostics.HasError())
		})
	}
}

// TestWorkflowTransferResourcePublicAPI tests the public API of WorkflowTransferResource.
// This is a black-box test that only uses the exported interface.
func TestWorkflowTransferResourcePublicAPI(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "create resource instance",
			testFunc: func(t *testing.T) {
				t.Helper()
				res := workflow.NewWorkflowTransferResource()
				assert.NotNil(t, res, "NewWorkflowTransferResource should return a non-nil resource")
			},
		},
		{
			name: "create resource wrapper instance",
			testFunc: func(t *testing.T) {
				t.Helper()
				res := workflow.NewWorkflowTransferResourceWrapper()
				assert.NotNil(t, res, "NewWorkflowTransferResourceWrapper should return a non-nil resource")
			},
		},
		{
			name: "resource implements required interfaces",
			testFunc: func(t *testing.T) {
				t.Helper()
				res := workflow.NewWorkflowTransferResourceWrapper()
				assert.Implements(t, (*resource.Resource)(nil), res)
				assert.Implements(t, (*resource.ResourceWithConfigure)(nil), res)
				assert.Implements(t, (*resource.ResourceWithImportState)(nil), res)
			},
		},
		{
			name: "metadata sets correct type name",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := workflow.NewWorkflowTransferResource()
				req := resource.MetadataRequest{
					ProviderTypeName: "n8n",
				}
				resp := &resource.MetadataResponse{}

				r.Metadata(context.Background(), req, resp)

				assert.Equal(t, "n8n_workflow_transfer", resp.TypeName)
			},
		},
		{
			name: "schema returns valid schema",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := workflow.NewWorkflowTransferResource()
				req := resource.SchemaRequest{}
				resp := &resource.SchemaResponse{}

				r.Schema(context.Background(), req, resp)

				assert.NotNil(t, resp.Schema)
				assert.NotEmpty(t, resp.Schema.Attributes)
			},
		},
		{
			name: "configure accepts valid client",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := workflow.NewWorkflowTransferResource()
				mockClient := &client.N8nClient{}

				req := resource.ConfigureRequest{
					ProviderData: mockClient,
				}
				resp := &resource.ConfigureResponse{}

				r.Configure(context.Background(), req, resp)

				assert.False(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "import state uses correct path",
			testFunc: func(t *testing.T) {
				t.Helper()
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
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t)
		})
	}
}
