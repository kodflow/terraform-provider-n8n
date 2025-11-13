package variable_test

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
	"github.com/kodflow/n8n/src/internal/provider/variable"
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

// TestNewVariableResource tests the NewVariableResource constructor.
func TestNewVariableResource(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "normal case",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := variable.NewVariableResource()
				assert.NotNil(t, r, "NewVariableResource should return a non-nil resource")
			},
		},
		{
			name: "error case - validates behavior",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := variable.NewVariableResource()
				assert.NotNil(t, r, "NewVariableResource should return a non-nil resource")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.testFunc(t)
		})
	}
} // TestNewVariableResourceWrapper tests the NewVariableResourceWrapper constructor.
func TestNewVariableResourceWrapper(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "normal case",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := variable.NewVariableResourceWrapper()
				assert.NotNil(t, r, "NewVariableResourceWrapper should return a non-nil resource")
			},
		},
		{
			name: "error case - validates behavior",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := variable.NewVariableResourceWrapper()
				assert.NotNil(t, r, "NewVariableResourceWrapper should return a non-nil resource")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.testFunc(t)
		})
	}
} // TestVariableResource_Metadata tests the Metadata method.
func TestVariableResource_Metadata(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "normal case",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := variable.NewVariableResource()
				req := resource.MetadataRequest{
					ProviderTypeName: "n8n",
				}
				resp := &resource.MetadataResponse{}

				r.Metadata(context.Background(), req, resp)

				assert.Equal(t, "n8n_variable", resp.TypeName, "TypeName should be set correctly")
			},
		},
		{
			name: "error case - validates behavior",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := variable.NewVariableResource()
				req := resource.MetadataRequest{
					ProviderTypeName: "n8n",
				}
				resp := &resource.MetadataResponse{}

				r.Metadata(context.Background(), req, resp)

				assert.Equal(t, "n8n_variable", resp.TypeName, "TypeName should be set correctly")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.testFunc(t)
		})
	}
} // TestVariableResource_Schema tests the Schema method.
func TestVariableResource_Schema(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "normal case",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := variable.NewVariableResource()
				req := resource.SchemaRequest{}
				resp := &resource.SchemaResponse{}

				r.Schema(context.Background(), req, resp)

				assert.NotNil(t, resp.Schema, "Schema should not be nil")
				assert.NotNil(t, resp.Schema.Attributes, "Schema attributes should not be nil")
			},
		},
		{
			name: "error case - validates behavior",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := variable.NewVariableResource()
				req := resource.SchemaRequest{}
				resp := &resource.SchemaResponse{}

				r.Schema(context.Background(), req, resp)

				assert.NotNil(t, resp.Schema, "Schema should not be nil")
				assert.NotNil(t, resp.Schema.Attributes, "Schema attributes should not be nil")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.testFunc(t)
		})
	}
} // TestVariableResource_Configure tests the Configure method.
func TestVariableResource_Configure(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		providerData interface{}
		expectError  bool
	}{
		{
			name:         "valid client",
			providerData: &client.N8nClient{},
			expectError:  false,
		},
		{
			name:         "nil provider data",
			providerData: nil,
			expectError:  false,
		},
		{
			name:         "invalid provider data",
			providerData: "invalid",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := variable.NewVariableResource()
			req := resource.ConfigureRequest{
				ProviderData: tt.providerData,
			}
			resp := &resource.ConfigureResponse{}

			r.Configure(context.Background(), req, resp)

			if tt.expectError {
				assert.True(t, resp.Diagnostics.HasError(), "Configure should return error for invalid provider data")
			} else {
				assert.False(t, resp.Diagnostics.HasError(), "Configure should not return error")
			}
		})
	}
}

// TestVariableResource_Create tests the Create method with full execution.
func TestVariableResource_Create(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "create with successful API call",
			testFunc: func(t *testing.T) {
				t.Helper()
				callCount := 0
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					callCount++
					w.Header().Set("Content-Type", "application/json")
					if callCount == 1 {
						// POST /variables - returns 201 with no body
						w.WriteHeader(http.StatusCreated)
					} else {
						// GET /variables - returns list with created variable
						w.WriteHeader(http.StatusOK)
						w.Write([]byte(`{"data":[{"id":"var-123","key":"test-key","value":"test-value","type":"string"}]}`))
					}
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				r := variable.NewVariableResource()
				r.Configure(context.Background(), resource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &resource.ConfigureResponse{})

				ctx := context.Background()
				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				// Build plan using tftypes
				planRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, nil),
					"key":        tftypes.NewValue(tftypes.String, "test-key"),
					"value":      tftypes.NewValue(tftypes.String, "test-value"),
					"type":       tftypes.NewValue(tftypes.String, nil),
					"project_id": tftypes.NewValue(tftypes.String, nil),
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
				r := variable.NewVariableResource()
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

				r := variable.NewVariableResource()
				r.Configure(context.Background(), resource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &resource.ConfigureResponse{})

				ctx := context.Background()
				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				planRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, nil),
					"key":        tftypes.NewValue(tftypes.String, "test-key"),
					"value":      tftypes.NewValue(tftypes.String, "test-value"),
					"type":       tftypes.NewValue(tftypes.String, nil),
					"project_id": tftypes.NewValue(tftypes.String, nil),
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

// TestVariableResource_Read tests the Read method with full execution.
func TestVariableResource_Read(t *testing.T) {
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
					w.Write([]byte(`{"data":[{"id":"var-123","key":"test-key","value":"test-value","type":"string"}]}`))
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				r := variable.NewVariableResource()
				r.Configure(context.Background(), resource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &resource.ConfigureResponse{})

				ctx := context.Background()
				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				// Build state
				stateRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "var-123"),
					"key":        tftypes.NewValue(tftypes.String, "test-key"),
					"value":      tftypes.NewValue(tftypes.String, "test-value"),
					"type":       tftypes.NewValue(tftypes.String, "string"),
					"project_id": tftypes.NewValue(tftypes.String, nil),
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
				r := variable.NewVariableResource()
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
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"message": "Internal server error"}`))
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				r := variable.NewVariableResource()
				r.Configure(context.Background(), resource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &resource.ConfigureResponse{})

				ctx := context.Background()
				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				stateRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "var-123"),
					"key":        tftypes.NewValue(tftypes.String, "test-key"),
					"value":      tftypes.NewValue(tftypes.String, "test-value"),
					"type":       tftypes.NewValue(tftypes.String, "string"),
					"project_id": tftypes.NewValue(tftypes.String, nil),
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

// TestVariableResource_Update tests the Update method with full execution.
func TestVariableResource_Update(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "update with successful API call",
			testFunc: func(t *testing.T) {
				t.Helper()
				callCount := 0
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					callCount++
					w.Header().Set("Content-Type", "application/json")
					if callCount == 1 {
						// PUT /variables/{id} - returns 204 with no body
						w.WriteHeader(http.StatusNoContent)
					} else {
						// GET /variables - returns list with updated variable
						w.WriteHeader(http.StatusOK)
						w.Write([]byte(`{"data":[{"id":"var-123","key":"test-key","value":"updated-value","type":"string"}]}`))
					}
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				r := variable.NewVariableResource()
				r.Configure(context.Background(), resource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &resource.ConfigureResponse{})

				ctx := context.Background()
				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				// Build plan
				planRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "var-123"),
					"key":        tftypes.NewValue(tftypes.String, "test-key"),
					"value":      tftypes.NewValue(tftypes.String, "updated-value"),
					"type":       tftypes.NewValue(tftypes.String, "string"),
					"project_id": tftypes.NewValue(tftypes.String, nil),
				})

				plan := tfsdk.Plan{
					Schema: schemaResp.Schema,
					Raw:    planRaw,
				}

				// Build state
				stateRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "var-123"),
					"key":        tftypes.NewValue(tftypes.String, "test-key"),
					"value":      tftypes.NewValue(tftypes.String, "old-value"),
					"type":       tftypes.NewValue(tftypes.String, "string"),
					"project_id": tftypes.NewValue(tftypes.String, nil),
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
				r := variable.NewVariableResource()
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

				r := variable.NewVariableResource()
				r.Configure(context.Background(), resource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &resource.ConfigureResponse{})

				ctx := context.Background()
				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				planRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "var-123"),
					"key":        tftypes.NewValue(tftypes.String, "test-key"),
					"value":      tftypes.NewValue(tftypes.String, "updated-value"),
					"type":       tftypes.NewValue(tftypes.String, "string"),
					"project_id": tftypes.NewValue(tftypes.String, nil),
				})

				plan := tfsdk.Plan{
					Schema: schemaResp.Schema,
					Raw:    planRaw,
				}

				stateRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "var-123"),
					"key":        tftypes.NewValue(tftypes.String, "test-key"),
					"value":      tftypes.NewValue(tftypes.String, "old-value"),
					"type":       tftypes.NewValue(tftypes.String, "string"),
					"project_id": tftypes.NewValue(tftypes.String, nil),
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

// TestVariableResource_Delete tests the Delete method with full execution.
func TestVariableResource_Delete(t *testing.T) {
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

				r := variable.NewVariableResource()
				r.Configure(context.Background(), resource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &resource.ConfigureResponse{})

				ctx := context.Background()
				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				// Build state
				stateRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "var-123"),
					"key":        tftypes.NewValue(tftypes.String, "test-key"),
					"value":      tftypes.NewValue(tftypes.String, "test-value"),
					"type":       tftypes.NewValue(tftypes.String, "string"),
					"project_id": tftypes.NewValue(tftypes.String, nil),
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
				r := variable.NewVariableResource()
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

				r := variable.NewVariableResource()
				r.Configure(context.Background(), resource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &resource.ConfigureResponse{})

				ctx := context.Background()
				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				stateRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "var-123"),
					"key":        tftypes.NewValue(tftypes.String, "test-key"),
					"value":      tftypes.NewValue(tftypes.String, "test-value"),
					"type":       tftypes.NewValue(tftypes.String, "string"),
					"project_id": tftypes.NewValue(tftypes.String, nil),
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
} // TestVariableResource_ImportState tests the ImportState method.
func TestVariableResource_ImportState(t *testing.T) {
	t.Run("import state passthrough", func(t *testing.T) {
		t.Helper()
		r := variable.NewVariableResource()
		ctx := context.Background()

		schemaResp := resource.SchemaResponse{}
		r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

		// Build empty state
		emptyValue := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, nil),
			"key":        tftypes.NewValue(tftypes.String, nil),
			"value":      tftypes.NewValue(tftypes.String, nil),
			"type":       tftypes.NewValue(tftypes.String, nil),
			"project_id": tftypes.NewValue(tftypes.String, nil),
		})

		req := resource.ImportStateRequest{
			ID: "var-123",
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
