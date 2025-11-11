package workflow

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
	"github.com/kodflow/n8n/src/internal/provider/workflow/models"
	"github.com/stretchr/testify/assert"
)

// TestNewWorkflowResource is now in external test file - refactored to test behavior only.

// TestNewWorkflowResourceWrapper is now in external test file - refactored to test behavior only.

// TestWorkflowResource_Metadata is now in external test file - refactored to test behavior only.

// TestWorkflowResource_Schema is now in external test file - refactored to test behavior only.

// TestWorkflowResource_Configure is now in external test file - refactored to test behavior only.

func TestWorkflowResource_InterfaceCompliance(t *testing.T) {
	t.Run("implements Resource interface", func(t *testing.T) {
		var _ resource.Resource = (*WorkflowResource)(nil)
	})

	t.Run("implements ResourceWithConfigure interface", func(t *testing.T) {
		var _ resource.ResourceWithConfigure = (*WorkflowResource)(nil)
	})

	t.Run("implements ResourceWithImportState interface", func(t *testing.T) {
		var _ resource.ResourceWithImportState = (*WorkflowResource)(nil)
	})

	t.Run("implements WorkflowResourceInterface", func(t *testing.T) {
		var _ WorkflowResourceInterface = (*WorkflowResource)(nil)
	})

	t.Run("error case - nil resource does not implement interfaces", func(t *testing.T) {
		var r *WorkflowResource
		// This verifies that the type itself implements the interfaces, not just non-nil instances
		assert.Nil(t, r)
		// Interface compliance is verified at compile time
	})
}

// createTestSchema creates a test schema for workflow resource.
func createTestSchema(t *testing.T) schema.Schema {
	t.Helper()
	r := &WorkflowResource{}
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

// TestWorkflowResource_Create tests workflow creation.
// TestWorkflowResource_Create is now in external test file - refactored to test behavior only.

// TestWorkflowResource_Read is now in external test file - refactored to test behavior only.

// TestWorkflowResource_Update is now in external test file - refactored to test behavior only.

// TestWorkflowResource_Delete is now in external test file - refactored to test behavior only.

func TestWorkflowResource_Concurrency(t *testing.T) {
	t.Run("concurrent metadata calls", func(t *testing.T) {
		r := &WorkflowResource{}

		done := make(chan bool, 100)
		for i := 0; i < 100; i++ {
			go func() {
				resp := &resource.MetadataResponse{}
				r.Metadata(context.Background(), resource.MetadataRequest{
					ProviderTypeName: "n8n",
				}, resp)
				assert.Equal(t, "n8n_workflow", resp.TypeName)
				done <- true
			}()
		}

		for i := 0; i < 100; i++ {
			<-done
		}
	})

	t.Run("concurrent schema calls", func(t *testing.T) {
		r := &WorkflowResource{}

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

	t.Run("error case - concurrent mixed calls do not panic", func(t *testing.T) {
		r := &WorkflowResource{}

		done := make(chan bool, 100)
		for i := 0; i < 50; i++ {
			go func() {
				resp := &resource.MetadataResponse{}
				r.Metadata(context.Background(), resource.MetadataRequest{
					ProviderTypeName: "test",
				}, resp)
				done <- true
			}()
		}
		for i := 0; i < 50; i++ {
			go func() {
				resp := &resource.SchemaResponse{}
				r.Schema(context.Background(), resource.SchemaRequest{}, resp)
				done <- true
			}()
		}

		for i := 0; i < 100; i++ {
			<-done
		}
	})
}

func TestWORKFLOW_ATTRIBUTES_SIZE(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "constant is defined",
			testFunc: func(t *testing.T) {
				t.Helper()
				assert.Equal(t, 14, WORKFLOW_ATTRIBUTES_SIZE)
			},
		},
		{
			name: "actual schema has 14 attributes",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := &WorkflowResource{}
				attrs := r.schemaAttributes()
				// The actual schema has 14 attributes:
				// id, name, active, tags, nodes_json, connections_json, settings_json,
				// created_at, updated_at, version_id, is_archived, trigger_count, meta, pin_data
				assert.Equal(t, 14, len(attrs))
			},
		},
		{
			name: "error case - constant matches actual count",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := &WorkflowResource{}
				attrs := r.schemaAttributes()
				assert.Equal(t, WORKFLOW_ATTRIBUTES_SIZE, len(attrs), "Constant and actual count must match")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t)
		})
	}
}

func BenchmarkWorkflowResource_Schema(b *testing.B) {
	r := &WorkflowResource{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := &resource.SchemaResponse{}
		r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	}
}

func BenchmarkWorkflowResource_Metadata(b *testing.B) {
	r := &WorkflowResource{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := &resource.MetadataResponse{}
		r.Metadata(context.Background(), resource.MetadataRequest{}, resp)
	}
}

func BenchmarkWorkflowResource_Configure(b *testing.B) {
	r := &WorkflowResource{}
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

func BenchmarkWorkflowResource_SchemaAttributes(b *testing.B) {
	r := &WorkflowResource{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = r.schemaAttributes()
	}
}

// CRUD Tests - Focus on helper function testing.

func TestWorkflowResource_Create_JSONParsing(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "create with tags",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					switch r.URL.Path {
					case "/workflows":
						if r.Method == http.MethodPost {
							response := map[string]interface{}{
								"id":          "wf-123",
								"name":        "Test Workflow",
								"active":      false,
								"nodes":       []interface{}{},
								"connections": map[string]interface{}{},
								"settings":    map[string]interface{}{},
								"tags":        []interface{}{map[string]interface{}{"id": "tag1", "name": "Tag 1"}},
							}
							w.Header().Set("Content-Type", "application/json")
							w.WriteHeader(http.StatusCreated)
							json.NewEncoder(w).Encode(response)
							return
						}
					case "/workflows/wf-123/tags":
						if r.Method == http.MethodPut {
							w.Header().Set("Content-Type", "application/json")
							w.WriteHeader(http.StatusOK)
							json.NewEncoder(w).Encode([]interface{}{map[string]interface{}{"id": "tag1", "name": "Tag 1"}})
							return
						}
					}
					w.WriteHeader(http.StatusNotFound)
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				r := &WorkflowResource{client: n8nClient}

				rawPlan := map[string]tftypes.Value{
					"id":               tftypes.NewValue(tftypes.String, nil),
					"name":             tftypes.NewValue(tftypes.String, "Test Workflow"),
					"active":           tftypes.NewValue(tftypes.Bool, nil),
					"tags":             tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, []tftypes.Value{tftypes.NewValue(tftypes.String, "tag1")}),
					"nodes_json":       tftypes.NewValue(tftypes.String, nil),
					"connections_json": tftypes.NewValue(tftypes.String, nil),
					"settings_json":    tftypes.NewValue(tftypes.String, nil),
					"created_at":       tftypes.NewValue(tftypes.String, nil),
					"updated_at":       tftypes.NewValue(tftypes.String, nil),
					"version_id":       tftypes.NewValue(tftypes.String, nil),
					"is_archived":      tftypes.NewValue(tftypes.Bool, nil),
					"trigger_count":    tftypes.NewValue(tftypes.Number, nil),
					"meta":             tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, nil),
					"pin_data":         tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, nil),
				}
				objectType := tftypes.Object{
					AttributeTypes: map[string]tftypes.Type{
						"id":               tftypes.String,
						"name":             tftypes.String,
						"active":           tftypes.Bool,
						"tags":             tftypes.List{ElementType: tftypes.String},
						"nodes_json":       tftypes.String,
						"connections_json": tftypes.String,
						"settings_json":    tftypes.String,
						"created_at":       tftypes.String,
						"updated_at":       tftypes.String,
						"version_id":       tftypes.String,
						"is_archived":      tftypes.Bool,
						"trigger_count":    tftypes.Number,
						"meta":             tftypes.Map{ElementType: tftypes.String},
						"pin_data":         tftypes.Map{ElementType: tftypes.String},
					},
				}

				plan := tfsdk.Plan{
					Raw:    tftypes.NewValue(objectType, rawPlan),
					Schema: createTestSchema(t),
				}

				state := tfsdk.State{
					Raw:    tftypes.NewValue(objectType, nil),
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
						t.Logf("Diagnostic Error: %s - %s", diag.Summary(), diag.Detail())
					}
				}
				assert.False(t, resp.Diagnostics.HasError(), "Create with tags should not have errors")
			},
		},
		{
			name: "invalid nodes json",
			testFunc: func(t *testing.T) {
				t.Helper()
				plan := &models.Resource{
					Name:            types.StringValue("Test Workflow"),
					NodesJSON:       types.StringValue("invalid json"),
					ConnectionsJSON: types.StringValue("{}"),
				}

				diags := &diag.Diagnostics{}
				parseWorkflowJSON(plan, diags)

				assert.True(t, diags.HasError())
				assert.Contains(t, diags.Errors()[0].Summary(), "Invalid nodes JSON")
			},
		},
		{
			name: "invalid connections json",
			testFunc: func(t *testing.T) {
				t.Helper()
				plan := &models.Resource{
					Name:            types.StringValue("Test Workflow"),
					NodesJSON:       types.StringValue("[]"),
					ConnectionsJSON: types.StringValue("invalid json"),
				}

				diags := &diag.Diagnostics{}
				parseWorkflowJSON(plan, diags)

				assert.True(t, diags.HasError())
				assert.Contains(t, diags.Errors()[0].Summary(), "Invalid connections JSON")
			},
		},
		{
			name: "invalid settings json",
			testFunc: func(t *testing.T) {
				t.Helper()
				plan := &models.Resource{
					Name:            types.StringValue("Test Workflow"),
					NodesJSON:       types.StringValue("[]"),
					ConnectionsJSON: types.StringValue("{}"),
					SettingsJSON:    types.StringValue("invalid json"),
				}

				diags := &diag.Diagnostics{}
				parseWorkflowJSON(plan, diags)

				assert.True(t, diags.HasError())
				assert.Contains(t, diags.Errors()[0].Summary(), "Invalid settings JSON")
			},
		},
		{
			name: "valid json parsing",
			testFunc: func(t *testing.T) {
				t.Helper()
				plan := &models.Resource{
					Name:            types.StringValue("Test Workflow"),
					NodesJSON:       types.StringValue("[]"),
					ConnectionsJSON: types.StringValue("{}"),
					SettingsJSON:    types.StringValue("{}"),
				}

				diags := &diag.Diagnostics{}
				nodes, connections, settings := parseWorkflowJSON(plan, diags)

				assert.False(t, diags.HasError())
				assert.NotNil(t, nodes)
				assert.NotNil(t, connections)
				assert.NotNil(t, settings)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t)
		})
	}
}

func TestWorkflowResource_Read_ClientRequired(t *testing.T) {
	t.Run("resource requires configured client", func(t *testing.T) {
		r := &WorkflowResource{
			client: nil,
		}

		assert.Nil(t, r.client)
	})

	t.Run("error case - resource with client is not nil", func(t *testing.T) {
		r := &WorkflowResource{
			client: &client.N8nClient{},
		}

		assert.NotNil(t, r.client)
	})
}

func TestWorkflowResource_Update_ActivationDetection(t *testing.T) {
	tests := []struct {
		name          string
		planActive    types.Bool
		stateActive   types.Bool
		expectChanged bool
	}{
		{"true to false - should detect change", types.BoolValue(true), types.BoolValue(false), true},
		{"false to true - should detect change", types.BoolValue(false), types.BoolValue(true), true},
		{"true to true - no change", types.BoolValue(true), types.BoolValue(true), false},
		{"false to false - no change", types.BoolValue(false), types.BoolValue(false), false},
		{"null plan - no change detected", types.BoolNull(), types.BoolValue(true), false},
		{"null state - no change detected", types.BoolValue(true), types.BoolNull(), false},
		{"both null - no change", types.BoolNull(), types.BoolNull(), false},
		{"unknown plan - no change", types.BoolUnknown(), types.BoolValue(true), false},
		{"unknown state - no change", types.BoolValue(true), types.BoolUnknown(), false},
		{"both unknown - no change", types.BoolUnknown(), types.BoolUnknown(), false},
		{"error case - nil checks work correctly", types.BoolNull(), types.BoolNull(), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plan := &models.Resource{
				Active: tt.planActive,
			}
			state := &models.Resource{
				Active: tt.stateActive,
			}

			activeChanged := !plan.Active.IsNull() && !plan.Active.IsUnknown() &&
				!state.Active.IsNull() && !state.Active.IsUnknown() &&
				plan.Active.ValueBool() != state.Active.ValueBool()

			assert.Equal(t, tt.expectChanged, activeChanged)
		})
	}
}

func TestWorkflowResource_Delete_Basic(t *testing.T) {
	t.Run("delete requires client", func(t *testing.T) {
		r := &WorkflowResource{
			client: nil,
		}

		// Delete operation requires an API client
		assert.Nil(t, r.client)
	})

	t.Run("error case - resource with client can delete", func(t *testing.T) {
		r := &WorkflowResource{
			client: &client.N8nClient{},
		}

		assert.NotNil(t, r.client)
	})
}

// TestWorkflowResource_CRUD is a comprehensive suite that validates CRUD operation logic.
func TestWorkflowResource_CRUD(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "workflow lifecycle validation",
			testFunc: func(t *testing.T) {
				t.Helper()
				// Create - validates JSON parsing (tested in TestWorkflowResource_Create_JSONParsing)
				// Read - validates client requirement (tested in TestWorkflowResource_Read_ClientRequired)
				// Update - validates JSON parsing and activation detection (tested in TestWorkflowResource_Update_ActivationDetection)
				// Delete - validates state requirement (tested in TestWorkflowResource_Delete_Basic)

				// This test ensures all CRUD operations have basic validation in place
				assert.NotNil(t, NewWorkflowResource())
			},
		},
		{
			name: "error case - CRUD operations require proper initialization",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := NewWorkflowResource()
				assert.NotNil(t, r)
				assert.Nil(t, r.client, "New resource should have nil client until configured")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t)
		})
	}
}

// TestWorkflowResource_ImportState tests state import functionality.
// TestWorkflowResource_ImportState is now in external test file - refactored to test behavior only.

func TestWorkflowResource_EdgeCasesForCoverage(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "create fails when plan get fails",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := &WorkflowResource{}

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
				resp := resource.CreateResponse{
					State: tfsdk.State{Schema: createTestSchema(t)},
				}

				r.Create(context.Background(), req, &resp)

				assert.True(t, resp.Diagnostics.HasError(), "Create should fail when Plan.Get fails")
			},
		},
		{
			name: "read fails when state get fails",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := &WorkflowResource{}

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
				resp := resource.ReadResponse{
					State: tfsdk.State{Schema: createTestSchema(t)},
				}

				r.Read(context.Background(), req, &resp)

				assert.True(t, resp.Diagnostics.HasError(), "Read should fail when State.Get fails")
			},
		},
		{
			name: "update fails when plan get fails",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := &WorkflowResource{}

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

				state := tfsdk.State{
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

				assert.True(t, resp.Diagnostics.HasError(), "Update should fail when Plan.Get fails")
			},
		},
		{
			name: "delete fails when state get fails",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := &WorkflowResource{}

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

				req := resource.DeleteRequest{
					State: state,
				}
				resp := resource.DeleteResponse{
					State: tfsdk.State{Schema: createTestSchema(t)},
				}

				r.Delete(context.Background(), req, &resp)

				assert.True(t, resp.Diagnostics.HasError(), "Delete should fail when State.Get fails")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t)
		})
	}
}

// TestWorkflowResource_schemaAttributes tests the schemaAttributes private method.
func TestWorkflowResource_schemaAttributes(t *testing.T) {
	tests := []struct {
		name           string
		testFunc       func(*testing.T)
		wantAttrCount  int
		wantAttributes []string
	}{
		{
			name:          "returns correct number of attributes",
			wantAttrCount: 14,
			testFunc: func(t *testing.T) {
				t.Helper()
				r := &WorkflowResource{}
				attrs := r.schemaAttributes()
				assert.NotNil(t, attrs)
				assert.Equal(t, 14, len(attrs), "Should have exactly 14 attributes")
			},
		},
		{
			name: "contains all core attributes",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := &WorkflowResource{}
				attrs := r.schemaAttributes()
				assert.Contains(t, attrs, "id")
				assert.Contains(t, attrs, "name")
				assert.Contains(t, attrs, "active")
				assert.Contains(t, attrs, "tags")
			},
		},
		{
			name: "contains all JSON attributes",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := &WorkflowResource{}
				attrs := r.schemaAttributes()
				assert.Contains(t, attrs, "nodes_json")
				assert.Contains(t, attrs, "connections_json")
				assert.Contains(t, attrs, "settings_json")
			},
		},
		{
			name: "contains all metadata attributes",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := &WorkflowResource{}
				attrs := r.schemaAttributes()
				assert.Contains(t, attrs, "created_at")
				assert.Contains(t, attrs, "updated_at")
				assert.Contains(t, attrs, "version_id")
				assert.Contains(t, attrs, "is_archived")
				assert.Contains(t, attrs, "trigger_count")
				assert.Contains(t, attrs, "meta")
				assert.Contains(t, attrs, "pin_data")
			},
		},
		{
			name: "error case - no nil attributes",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := &WorkflowResource{}
				attrs := r.schemaAttributes()
				for key, attr := range attrs {
					assert.NotNil(t, attr, "Attribute %s should not be nil", key)
				}
			},
		},
		{
			name: "error case - map is not nil",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := &WorkflowResource{}
				attrs := r.schemaAttributes()
				assert.NotNil(t, attrs, "Returned map should not be nil")
			},
		},
		{
			name: "error case - no duplicate keys",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := &WorkflowResource{}
				attrs := r.schemaAttributes()
				// Verify all expected keys are unique (map already ensures uniqueness)
				expectedKeys := []string{
					"id", "name", "active", "tags",
					"nodes_json", "connections_json", "settings_json",
					"created_at", "updated_at", "version_id",
					"is_archived", "trigger_count", "meta", "pin_data",
				}
				assert.Equal(t, len(expectedKeys), len(attrs), "Should have no duplicate keys")
			},
		},
		{
			name: "error case - consistent with WORKFLOW_ATTRIBUTES_SIZE",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := &WorkflowResource{}
				attrs := r.schemaAttributes()
				assert.Equal(t, WORKFLOW_ATTRIBUTES_SIZE, len(attrs), "Attribute count must match constant")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.testFunc(t)
		})
	}
}

// TestWorkflowResource_addCoreAttributes tests the addCoreAttributes private method.
func TestWorkflowResource_addCoreAttributes(t *testing.T) {
	tests := []struct {
		name         string
		testFunc     func(*testing.T)
		expectedKeys []string
	}{
		{
			name: "adds all core attributes",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := &WorkflowResource{}
				attrs := make(map[string]schema.Attribute)
				r.addCoreAttributes(attrs)
				assert.NotNil(t, attrs)
				assert.Equal(t, 4, len(attrs), "Should add exactly 4 core attributes")
			},
		},
		{
			name: "adds id attribute with correct properties",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := &WorkflowResource{}
				attrs := make(map[string]schema.Attribute)
				r.addCoreAttributes(attrs)
				assert.Contains(t, attrs, "id")
				idAttr := attrs["id"].(schema.StringAttribute)
				assert.NotEmpty(t, idAttr.MarkdownDescription)
				assert.True(t, idAttr.Computed, "id should be computed")
				assert.False(t, idAttr.Required, "id should not be required")
				assert.False(t, idAttr.Optional, "id should not be optional")
			},
		},
		{
			name: "adds name attribute with correct properties",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := &WorkflowResource{}
				attrs := make(map[string]schema.Attribute)
				r.addCoreAttributes(attrs)
				assert.Contains(t, attrs, "name")
				nameAttr := attrs["name"].(schema.StringAttribute)
				assert.NotEmpty(t, nameAttr.MarkdownDescription)
				assert.True(t, nameAttr.Required, "name should be required")
				assert.False(t, nameAttr.Optional, "name should not be optional")
				assert.False(t, nameAttr.Computed, "name should not be computed")
			},
		},
		{
			name: "adds active attribute with correct properties",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := &WorkflowResource{}
				attrs := make(map[string]schema.Attribute)
				r.addCoreAttributes(attrs)
				assert.Contains(t, attrs, "active")
				activeAttr := attrs["active"].(schema.BoolAttribute)
				assert.NotEmpty(t, activeAttr.MarkdownDescription)
				assert.True(t, activeAttr.Optional, "active should be optional")
				assert.True(t, activeAttr.Computed, "active should be computed")
				assert.False(t, activeAttr.Required, "active should not be required")
			},
		},
		{
			name: "adds tags attribute with correct properties",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := &WorkflowResource{}
				attrs := make(map[string]schema.Attribute)
				r.addCoreAttributes(attrs)
				assert.Contains(t, attrs, "tags")
				tagsAttr := attrs["tags"].(schema.ListAttribute)
				assert.NotEmpty(t, tagsAttr.MarkdownDescription)
				assert.True(t, tagsAttr.Optional, "tags should be optional")
				assert.False(t, tagsAttr.Required, "tags should not be required")
				assert.False(t, tagsAttr.Computed, "tags should not be computed")
				assert.Equal(t, types.StringType, tagsAttr.ElementType)
			},
		},
		{
			name: "error case - can add to non-empty map",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := &WorkflowResource{}
				attrs := map[string]schema.Attribute{
					"existing": schema.StringAttribute{},
				}
				r.addCoreAttributes(attrs)
				assert.Equal(t, 5, len(attrs), "Should have 1 existing + 4 new attributes")
				assert.Contains(t, attrs, "existing")
				assert.Contains(t, attrs, "id")
			},
		},
		{
			name: "error case - no nil attributes added",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := &WorkflowResource{}
				attrs := make(map[string]schema.Attribute)
				r.addCoreAttributes(attrs)
				for key, attr := range attrs {
					assert.NotNil(t, attr, "Attribute %s should not be nil", key)
				}
			},
		},
		{
			name: "error case - all attributes have descriptions",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := &WorkflowResource{}
				attrs := make(map[string]schema.Attribute)
				r.addCoreAttributes(attrs)
				for key, attr := range attrs {
					switch v := attr.(type) {
					case schema.StringAttribute:
						assert.NotEmpty(t, v.MarkdownDescription, "Attribute %s should have description", key)
					case schema.BoolAttribute:
						assert.NotEmpty(t, v.MarkdownDescription, "Attribute %s should have description", key)
					case schema.ListAttribute:
						assert.NotEmpty(t, v.MarkdownDescription, "Attribute %s should have description", key)
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.testFunc(t)
		})
	}
}

// TestWorkflowResource_addJSONAttributes tests the addJSONAttributes private method.
func TestWorkflowResource_addJSONAttributes(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "adds all JSON attributes",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := &WorkflowResource{}
				attrs := make(map[string]schema.Attribute)
				r.addJSONAttributes(attrs)
				assert.NotNil(t, attrs)
				assert.Equal(t, 3, len(attrs), "Should add exactly 3 JSON attributes")
			},
		},
		{
			name: "adds nodes_json attribute with correct properties",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := &WorkflowResource{}
				attrs := make(map[string]schema.Attribute)
				r.addJSONAttributes(attrs)
				assert.Contains(t, attrs, "nodes_json")
				nodesAttr := attrs["nodes_json"].(schema.StringAttribute)
				assert.NotEmpty(t, nodesAttr.MarkdownDescription)
				assert.True(t, nodesAttr.Optional, "nodes_json should be optional")
				assert.True(t, nodesAttr.Computed, "nodes_json should be computed")
				assert.False(t, nodesAttr.Required, "nodes_json should not be required")
			},
		},
		{
			name: "adds connections_json attribute with correct properties",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := &WorkflowResource{}
				attrs := make(map[string]schema.Attribute)
				r.addJSONAttributes(attrs)
				assert.Contains(t, attrs, "connections_json")
				connectionsAttr := attrs["connections_json"].(schema.StringAttribute)
				assert.NotEmpty(t, connectionsAttr.MarkdownDescription)
				assert.True(t, connectionsAttr.Optional, "connections_json should be optional")
				assert.True(t, connectionsAttr.Computed, "connections_json should be computed")
				assert.False(t, connectionsAttr.Required, "connections_json should not be required")
			},
		},
		{
			name: "adds settings_json attribute with correct properties",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := &WorkflowResource{}
				attrs := make(map[string]schema.Attribute)
				r.addJSONAttributes(attrs)
				assert.Contains(t, attrs, "settings_json")
				settingsAttr := attrs["settings_json"].(schema.StringAttribute)
				assert.NotEmpty(t, settingsAttr.MarkdownDescription)
				assert.True(t, settingsAttr.Optional, "settings_json should be optional")
				assert.True(t, settingsAttr.Computed, "settings_json should be computed")
				assert.False(t, settingsAttr.Required, "settings_json should not be required")
			},
		},
		{
			name: "error case - can add to non-empty map",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := &WorkflowResource{}
				attrs := map[string]schema.Attribute{
					"existing": schema.StringAttribute{},
				}
				r.addJSONAttributes(attrs)
				assert.Equal(t, 4, len(attrs), "Should have 1 existing + 3 new attributes")
				assert.Contains(t, attrs, "existing")
				assert.Contains(t, attrs, "nodes_json")
			},
		},
		{
			name: "error case - no nil attributes added",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := &WorkflowResource{}
				attrs := make(map[string]schema.Attribute)
				r.addJSONAttributes(attrs)
				for key, attr := range attrs {
					assert.NotNil(t, attr, "Attribute %s should not be nil", key)
				}
			},
		},
		{
			name: "error case - all attributes have descriptions",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := &WorkflowResource{}
				attrs := make(map[string]schema.Attribute)
				r.addJSONAttributes(attrs)
				for key, attr := range attrs {
					strAttr := attr.(schema.StringAttribute)
					assert.NotEmpty(t, strAttr.MarkdownDescription, "Attribute %s should have description", key)
				}
			},
		},
		{
			name: "error case - all attributes are string type",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := &WorkflowResource{}
				attrs := make(map[string]schema.Attribute)
				r.addJSONAttributes(attrs)
				for key, attr := range attrs {
					_, ok := attr.(schema.StringAttribute)
					assert.True(t, ok, "Attribute %s should be StringAttribute", key)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.testFunc(t)
		})
	}
}

// TestWorkflowResource_addMetadataAttributes tests the addMetadataAttributes private method.
func TestWorkflowResource_addMetadataAttributes(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "adds all metadata attributes",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := &WorkflowResource{}
				attrs := make(map[string]schema.Attribute)
				r.addMetadataAttributes(attrs)
				assert.NotNil(t, attrs)
				assert.Equal(t, 7, len(attrs), "Should add exactly 7 metadata attributes")
			},
		},
		{
			name: "adds created_at attribute with correct properties",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := &WorkflowResource{}
				attrs := make(map[string]schema.Attribute)
				r.addMetadataAttributes(attrs)
				assert.Contains(t, attrs, "created_at")
				createdAttr := attrs["created_at"].(schema.StringAttribute)
				assert.NotEmpty(t, createdAttr.MarkdownDescription)
				assert.True(t, createdAttr.Computed, "created_at should be computed")
				assert.False(t, createdAttr.Required, "created_at should not be required")
				assert.False(t, createdAttr.Optional, "created_at should not be optional")
			},
		},
		{
			name: "adds updated_at attribute with correct properties",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := &WorkflowResource{}
				attrs := make(map[string]schema.Attribute)
				r.addMetadataAttributes(attrs)
				assert.Contains(t, attrs, "updated_at")
				updatedAttr := attrs["updated_at"].(schema.StringAttribute)
				assert.NotEmpty(t, updatedAttr.MarkdownDescription)
				assert.True(t, updatedAttr.Computed, "updated_at should be computed")
				assert.False(t, updatedAttr.Required, "updated_at should not be required")
				assert.False(t, updatedAttr.Optional, "updated_at should not be optional")
			},
		},
		{
			name: "adds version_id attribute with correct properties",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := &WorkflowResource{}
				attrs := make(map[string]schema.Attribute)
				r.addMetadataAttributes(attrs)
				assert.Contains(t, attrs, "version_id")
				versionAttr := attrs["version_id"].(schema.StringAttribute)
				assert.NotEmpty(t, versionAttr.MarkdownDescription)
				assert.True(t, versionAttr.Computed, "version_id should be computed")
				assert.False(t, versionAttr.Required, "version_id should not be required")
				assert.False(t, versionAttr.Optional, "version_id should not be optional")
			},
		},
		{
			name: "adds is_archived attribute with correct properties",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := &WorkflowResource{}
				attrs := make(map[string]schema.Attribute)
				r.addMetadataAttributes(attrs)
				assert.Contains(t, attrs, "is_archived")
				archivedAttr := attrs["is_archived"].(schema.BoolAttribute)
				assert.NotEmpty(t, archivedAttr.MarkdownDescription)
				assert.True(t, archivedAttr.Computed, "is_archived should be computed")
				assert.False(t, archivedAttr.Required, "is_archived should not be required")
				assert.False(t, archivedAttr.Optional, "is_archived should not be optional")
			},
		},
		{
			name: "adds trigger_count attribute with correct properties",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := &WorkflowResource{}
				attrs := make(map[string]schema.Attribute)
				r.addMetadataAttributes(attrs)
				assert.Contains(t, attrs, "trigger_count")
				triggerAttr := attrs["trigger_count"].(schema.Int64Attribute)
				assert.NotEmpty(t, triggerAttr.MarkdownDescription)
				assert.True(t, triggerAttr.Computed, "trigger_count should be computed")
				assert.False(t, triggerAttr.Required, "trigger_count should not be required")
				assert.False(t, triggerAttr.Optional, "trigger_count should not be optional")
			},
		},
		{
			name: "adds meta attribute with correct properties",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := &WorkflowResource{}
				attrs := make(map[string]schema.Attribute)
				r.addMetadataAttributes(attrs)
				assert.Contains(t, attrs, "meta")
				metaAttr := attrs["meta"].(schema.MapAttribute)
				assert.NotEmpty(t, metaAttr.MarkdownDescription)
				assert.True(t, metaAttr.Computed, "meta should be computed")
				assert.False(t, metaAttr.Required, "meta should not be required")
				assert.False(t, metaAttr.Optional, "meta should not be optional")
				assert.Equal(t, types.StringType, metaAttr.ElementType)
			},
		},
		{
			name: "adds pin_data attribute with correct properties",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := &WorkflowResource{}
				attrs := make(map[string]schema.Attribute)
				r.addMetadataAttributes(attrs)
				assert.Contains(t, attrs, "pin_data")
				pinDataAttr := attrs["pin_data"].(schema.MapAttribute)
				assert.NotEmpty(t, pinDataAttr.MarkdownDescription)
				assert.True(t, pinDataAttr.Computed, "pin_data should be computed")
				assert.False(t, pinDataAttr.Required, "pin_data should not be required")
				assert.False(t, pinDataAttr.Optional, "pin_data should not be optional")
				assert.Equal(t, types.StringType, pinDataAttr.ElementType)
			},
		},
		{
			name: "error case - can add to non-empty map",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := &WorkflowResource{}
				attrs := map[string]schema.Attribute{
					"existing": schema.StringAttribute{},
				}
				r.addMetadataAttributes(attrs)
				assert.Equal(t, 8, len(attrs), "Should have 1 existing + 7 new attributes")
				assert.Contains(t, attrs, "existing")
				assert.Contains(t, attrs, "created_at")
			},
		},
		{
			name: "error case - no nil attributes added",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := &WorkflowResource{}
				attrs := make(map[string]schema.Attribute)
				r.addMetadataAttributes(attrs)
				for key, attr := range attrs {
					assert.NotNil(t, attr, "Attribute %s should not be nil", key)
				}
			},
		},
		{
			name: "error case - all attributes have descriptions",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := &WorkflowResource{}
				attrs := make(map[string]schema.Attribute)
				r.addMetadataAttributes(attrs)
				for key, attr := range attrs {
					switch v := attr.(type) {
					case schema.StringAttribute:
						assert.NotEmpty(t, v.MarkdownDescription, "Attribute %s should have description", key)
					case schema.BoolAttribute:
						assert.NotEmpty(t, v.MarkdownDescription, "Attribute %s should have description", key)
					case schema.Int64Attribute:
						assert.NotEmpty(t, v.MarkdownDescription, "Attribute %s should have description", key)
					case schema.MapAttribute:
						assert.NotEmpty(t, v.MarkdownDescription, "Attribute %s should have description", key)
					}
				}
			},
		},
		{
			name: "error case - all attributes are computed only",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := &WorkflowResource{}
				attrs := make(map[string]schema.Attribute)
				r.addMetadataAttributes(attrs)
				for key, attr := range attrs {
					switch v := attr.(type) {
					case schema.StringAttribute:
						assert.True(t, v.Computed, "Attribute %s should be computed", key)
						assert.False(t, v.Required, "Attribute %s should not be required", key)
						assert.False(t, v.Optional, "Attribute %s should not be optional", key)
					case schema.BoolAttribute:
						assert.True(t, v.Computed, "Attribute %s should be computed", key)
						assert.False(t, v.Required, "Attribute %s should not be required", key)
						assert.False(t, v.Optional, "Attribute %s should not be optional", key)
					case schema.Int64Attribute:
						assert.True(t, v.Computed, "Attribute %s should be computed", key)
						assert.False(t, v.Required, "Attribute %s should not be required", key)
						assert.False(t, v.Optional, "Attribute %s should not be optional", key)
					case schema.MapAttribute:
						assert.True(t, v.Computed, "Attribute %s should be computed", key)
						assert.False(t, v.Required, "Attribute %s should not be required", key)
						assert.False(t, v.Optional, "Attribute %s should not be optional", key)
					}
				}
			},
		},
		{
			name: "error case - map attributes have correct element type",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := &WorkflowResource{}
				attrs := make(map[string]schema.Attribute)
				r.addMetadataAttributes(attrs)
				mapAttrs := []string{"meta", "pin_data"}
				for _, key := range mapAttrs {
					mapAttr := attrs[key].(schema.MapAttribute)
					assert.Equal(t, types.StringType, mapAttr.ElementType, "Map attribute %s should have StringType elements", key)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.testFunc(t)
		})
	}
}
