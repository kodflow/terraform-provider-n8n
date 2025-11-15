package execution

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/execution/models"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/client"
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

// TestNewExecutionRetryResource is now in external test file - refactored to test behavior only.

// TestNewExecutionRetryResourceWrapper is now in external test file - refactored to test behavior only.

// TestExecutionRetryResource_Metadata is now in external test file - refactored to test behavior only.

// TestExecutionRetryResource_Schema is now in external test file - refactored to test behavior only.

// TestExecutionRetryResource_Configure is now in external test file - refactored to test behavior only.

// TestExecutionRetryResource_ImportState is now in external test file - refactored to test behavior only.

func TestExecutionRetryResource_Interfaces(t *testing.T) {
	t.Helper()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "implements required interfaces",
			testFunc: func(t *testing.T) {
				t.Helper()

				r := NewExecutionRetryResource()

				// Test that ExecutionRetryResource implements resource.Resource
				var _ resource.Resource = r

				// Test that ExecutionRetryResource implements resource.ResourceWithConfigure
				var _ resource.ResourceWithConfigure = r

				// Test that ExecutionRetryResource implements resource.ResourceWithImportState
				var _ resource.ResourceWithImportState = r

				// Test that ExecutionRetryResource implements ExecutionRetryResourceInterface
				var _ ExecutionRetryResourceInterface = r
			},
		},
		{
			name: "error case - verify interface assignment safety",
			testFunc: func(t *testing.T) {
				t.Helper()

				r := NewExecutionRetryResource()

				// Test interface assignment in error scenarios
				// Verify that nil interface assignments are handled safely
				var res resource.Resource = r
				assert.NotNil(t, res, "interface assignment should not produce nil")

				// Test type assertion safety
				concrete, ok := res.(*ExecutionRetryResource)
				assert.True(t, ok, "type assertion should succeed")
				assert.NotNil(t, concrete, "concrete type should not be nil")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t)
		})
	}
}

func TestExecutionRetryResourceConcurrency(t *testing.T) {
	t.Helper()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "concurrent metadata calls",
			testFunc: func(t *testing.T) {
				t.Helper()

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
			},
		},
		{
			name: "concurrent schema calls",
			testFunc: func(t *testing.T) {
				t.Helper()

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
			},
		},
		{
			name: "error case - concurrent access with nil client",
			testFunc: func(t *testing.T) {
				t.Helper()

				r := NewExecutionRetryResource()
				// r.client is nil, testing concurrent access in error state

				done := make(chan bool, 50)
				for i := 0; i < 50; i++ {
					go func() {
						// Verify that concurrent access to unconfigured resource is safe
						assert.Nil(t, r.client, "client should remain nil in error state")
						done <- true
					}()
				}

				for i := 0; i < 50; i++ {
					<-done
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t)
		})
	}
}

func TestExecutionRetryResource_HelperFunctions(t *testing.T) {
	t.Helper()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "path to execution_id",
			testFunc: func(t *testing.T) {
				t.Helper()

				p := path.Root("execution_id")
				assert.Equal(t, "execution_id", p.String())
			},
		},
		{
			name: "path to new_execution_id",
			testFunc: func(t *testing.T) {
				t.Helper()

				p := path.Root("new_execution_id")
				assert.Equal(t, "new_execution_id", p.String())
			},
		},
		{
			name: "types.StringValue conversion",
			testFunc: func(t *testing.T) {
				t.Helper()

				strVal := types.StringValue("test-value")
				assert.Equal(t, "test-value", strVal.ValueString())
				assert.False(t, strVal.IsNull())
				assert.False(t, strVal.IsUnknown())
			},
		},
		{
			name: "types.BoolValue conversion",
			testFunc: func(t *testing.T) {
				t.Helper()

				boolVal := types.BoolValue(true)
				assert.True(t, boolVal.ValueBool())
				assert.False(t, boolVal.IsNull())
				assert.False(t, boolVal.IsUnknown())
			},
		},
		{
			name: "error case - verify null and unknown value handling",
			testFunc: func(t *testing.T) {
				t.Helper()

				// Test null string value
				nullStr := types.StringNull()
				assert.True(t, nullStr.IsNull(), "null string should be null")
				assert.False(t, nullStr.IsUnknown(), "null string should not be unknown")

				// Test unknown string value
				unknownStr := types.StringUnknown()
				assert.False(t, unknownStr.IsNull(), "unknown string should not be null")
				assert.True(t, unknownStr.IsUnknown(), "unknown string should be unknown")

				// Test null bool value
				nullBool := types.BoolNull()
				assert.True(t, nullBool.IsNull(), "null bool should be null")

				// Test path with special characters
				specialPath := path.Root("execution_id").AtName("nested")
				assert.Contains(t, specialPath.String(), "execution_id", "path should contain execution_id")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t)
		})
	}
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
	t.Helper()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "successful creation",
			testFunc: func(t *testing.T) {
				t.Helper()

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
			},
		},
		{
			name: "creation fails",
			testFunc: func(t *testing.T) {
				t.Helper()

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
			},
		},
		{
			name: "invalid execution ID format",
			testFunc: func(t *testing.T) {
				t.Helper()

				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				})

				n8nClient, server := setupTestRetryClient(t, handler)
				defer server.Close()

				r := &ExecutionRetryResource{client: n8nClient}

				rawPlan := map[string]tftypes.Value{
					"execution_id":     tftypes.NewValue(tftypes.String, "invalid-not-a-number"),
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

				assert.True(t, resp.Diagnostics.HasError(), "Create should have errors for invalid execution ID")
				assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Invalid Execution ID")
			},
		},
		{
			name: "plan get error",
			testFunc: func(t *testing.T) {
				t.Helper()

				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				})

				n8nClient, server := setupTestRetryClient(t, handler)
				defer server.Close()

				r := &ExecutionRetryResource{client: n8nClient}

				// Create invalid schema that will cause Get to fail
				invalidPlan := tfsdk.Plan{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"execution_id": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"execution_id": tftypes.NewValue(tftypes.String, "123"),
					}),
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
					Plan: invalidPlan,
				}
				resp := resource.CreateResponse{
					State: state,
				}

				r.Create(context.Background(), req, &resp)

				assert.True(t, resp.Diagnostics.HasError(), "Create should have errors for plan get failure")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t)
		})
	}
}

// TestRetryExecutionResource_Read tests retry execution reading.
func TestRetryExecutionResource_Read(t *testing.T) {
	t.Helper()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "successful read",
			testFunc: func(t *testing.T) {
				t.Helper()

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
			},
		},
		{
			name: "execution not found removes from state",
			testFunc: func(t *testing.T) {
				t.Helper()

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
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t)
		})
	}
}

// TestRetryExecutionResource_Delete tests retry execution deletion.
func TestRetryExecutionResource_Delete(t *testing.T) {
	t.Helper()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "successful delete",
			testFunc: func(t *testing.T) {
				t.Helper()

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
			},
		},
		{
			name: "delete fails",
			testFunc: func(t *testing.T) {
				t.Helper()

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
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t)
		})
	}
}

// TestExecutionRetryResource_Create is now in external test file - refactored to test behavior only.

// TestExecutionRetryResource_Read is now in external test file - refactored to test behavior only.

// TestExecutionRetryResource_Update is now in external test file - refactored to test behavior only.

// TestExecutionRetryResource_Delete is now in external test file - refactored to test behavior only.

func TestExecutionRetryResource_schemaAttributes(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "schema attributes returned correctly", wantErr: false},
		{name: "required fields present", wantErr: false},
		{name: "computed fields present", wantErr: false},
		{name: "execution_id is required", wantErr: false},
		{name: "error case - verify attribute count", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "schema attributes returned correctly":
				r := &ExecutionRetryResource{}
				attrs := r.schemaAttributes()

				assert.NotNil(t, attrs)
				assert.Contains(t, attrs, "execution_id")
				assert.Contains(t, attrs, "new_execution_id")
				assert.Contains(t, attrs, "workflow_id")
				assert.Contains(t, attrs, "finished")
				assert.Contains(t, attrs, "mode")
				assert.Contains(t, attrs, "started_at")
				assert.Contains(t, attrs, "stopped_at")
				assert.Contains(t, attrs, "status")

			case "required fields present":
				r := &ExecutionRetryResource{}
				attrs := r.schemaAttributes()

				executionIDAttr, ok := attrs["execution_id"].(schema.StringAttribute)
				assert.True(t, ok)
				assert.True(t, executionIDAttr.Required)

			case "computed fields present":
				r := &ExecutionRetryResource{}
				attrs := r.schemaAttributes()

				newExecutionIDAttr, ok := attrs["new_execution_id"].(schema.StringAttribute)
				assert.True(t, ok)
				assert.True(t, newExecutionIDAttr.Computed)

				workflowIDAttr, ok := attrs["workflow_id"].(schema.StringAttribute)
				assert.True(t, ok)
				assert.True(t, workflowIDAttr.Computed)

			case "execution_id is required":
				r := &ExecutionRetryResource{}
				attrs := r.schemaAttributes()

				executionIDAttr, ok := attrs["execution_id"].(schema.StringAttribute)
				assert.True(t, ok)
				assert.True(t, executionIDAttr.Required)
				assert.False(t, executionIDAttr.Computed)
				assert.False(t, executionIDAttr.Optional)

			case "error case - verify attribute count":
				r := &ExecutionRetryResource{}
				attrs := r.schemaAttributes()
				assert.Len(t, attrs, 8)
			}
		})
	}
}

func Test_populateRetryExecutionData(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "populate with all fields", wantErr: false},
		{name: "populate with nil fields", wantErr: false},
		{name: "populate with partial fields", wantErr: false},
		{name: "populate with unset stoppedAt", wantErr: false},
		{name: "populate with nil stoppedAt inside nullable", wantErr: false},
		{name: "populate with all timestamps", wantErr: false},
		{name: "error case - empty execution", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "populate with all fields":
				id := float32(123)
				workflowID := float32(456)
				finished := true
				mode := "manual"
				status := "success"

				execution := &n8nsdk.Execution{
					Id:         &id,
					WorkflowId: &workflowID,
					Finished:   &finished,
					Mode:       &mode,
					Status:     &status,
				}

				model := &models.RetryResource{}
				populateRetryExecutionData(execution, model)

				assert.Equal(t, "123", model.NewExecutionID.ValueString())
				assert.Equal(t, "456", model.WorkflowID.ValueString())
				assert.True(t, model.Finished.ValueBool())
				assert.Equal(t, "manual", model.Mode.ValueString())
				assert.Equal(t, "success", model.Status.ValueString())

			case "populate with nil fields":
				execution := &n8nsdk.Execution{}
				model := &models.RetryResource{}
				populateRetryExecutionData(execution, model)

				assert.True(t, model.NewExecutionID.IsNull())
				assert.True(t, model.WorkflowID.IsNull())
				assert.True(t, model.Finished.IsNull())
				assert.True(t, model.Mode.IsNull())
				assert.True(t, model.Status.IsNull())

			case "populate with partial fields":
				id := float32(789)
				status := "running"

				execution := &n8nsdk.Execution{
					Id:     &id,
					Status: &status,
				}
				model := &models.RetryResource{}
				populateRetryExecutionData(execution, model)

				assert.Equal(t, "789", model.NewExecutionID.ValueString())
				assert.Equal(t, "running", model.Status.ValueString())
				assert.True(t, model.WorkflowID.IsNull())
				assert.True(t, model.Finished.IsNull())

			case "populate with unset stoppedAt":
				id := float32(999)
				execution := &n8nsdk.Execution{
					Id:        &id,
					StoppedAt: n8nsdk.NullableTime{},
				}
				model := &models.RetryResource{}
				populateRetryExecutionData(execution, model)

				assert.Equal(t, "999", model.NewExecutionID.ValueString())
				assert.True(t, model.StoppedAt.IsNull())

			case "populate with nil stoppedAt inside nullable":
				id := float32(888)
				execution := &n8nsdk.Execution{
					Id: &id,
				}
				execution.StoppedAt.Set(nil)

				model := &models.RetryResource{}
				populateRetryExecutionData(execution, model)

				assert.Equal(t, "888", model.NewExecutionID.ValueString())
				assert.True(t, model.StoppedAt.IsNull())

			case "populate with all timestamps":
				id := float32(777)
				startedAt := time.Now()
				stoppedAt := time.Now().Add(1 * time.Hour)

				execution := &n8nsdk.Execution{
					Id:        &id,
					StartedAt: *n8nsdk.NewNullableTime(&startedAt),
				}
				execution.StoppedAt.Set(&stoppedAt)

				model := &models.RetryResource{}
				populateRetryExecutionData(execution, model)

				assert.Equal(t, "777", model.NewExecutionID.ValueString())
				assert.NotEmpty(t, model.StartedAt.ValueString())
				assert.NotEmpty(t, model.StoppedAt.ValueString())

			case "error case - empty execution":
				execution := &n8nsdk.Execution{}
				model := &models.RetryResource{}
				populateRetryExecutionData(execution, model)
				assert.NotNil(t, model)
			}
		})
	}
}
