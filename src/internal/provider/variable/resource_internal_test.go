package variable

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
	"github.com/kodflow/n8n/src/internal/provider/variable/models"
	"github.com/stretchr/testify/assert"
)

// setupTestResourceClient creates a test N8nClient with httptest server.
func setupTestResourceClient(t *testing.T, handler http.HandlerFunc) (*client.N8nClient, *httptest.Server) {
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

// TestVariableResource_Create is tested via external tests and helper method tests.
// This ensures coverage of the Create method through the private methods it calls.

// TestVariableResource_executeVariableCreate tests the executeVariableCreate method.
func TestVariableResource_executeVariableCreate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "success - create variable"},
		{name: "error - API error", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "success - create variable":
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/variables" && r.Method == http.MethodPost {
						w.WriteHeader(http.StatusCreated)
						return
					}
					w.WriteHeader(http.StatusNotFound)
				})

				n8nClient, server := setupTestResourceClient(t, handler)
				defer server.Close()

				r := &VariableResource{client: n8nClient}
				plan := &models.Resource{
					Key:   types.StringValue("test-key"),
					Value: types.StringValue("test-value"),
				}
				resp := &resource.CreateResponse{}

				result := r.executeVariableCreate(context.Background(), plan, resp)

				assert.True(t, result)
				assert.False(t, resp.Diagnostics.HasError())

			case "error - API error":
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
				})

				n8nClient, server := setupTestResourceClient(t, handler)
				defer server.Close()

				r := &VariableResource{client: n8nClient}
				plan := &models.Resource{
					Key:   types.StringValue("test-key"),
					Value: types.StringValue("test-value"),
				}
				resp := &resource.CreateResponse{}

				result := r.executeVariableCreate(context.Background(), plan, resp)

				assert.False(t, result)
				assert.True(t, resp.Diagnostics.HasError())
			}
		})
	}
}

// TestVariableResource_findCreatedVariable tests the findCreatedVariable method.
func TestVariableResource_findCreatedVariable(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "success - variable found"},
		{name: "error - API error", wantErr: true},
		{name: "error - variable not found", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "success - variable found":
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/variables" && r.Method == http.MethodGet {
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusOK)
						json.NewEncoder(w).Encode(n8nsdk.VariableList{
							Data: []n8nsdk.Variable{
								{Id: ptrString("var-123"), Key: "test-key", Value: "test-value"},
							},
						})
						return
					}
					w.WriteHeader(http.StatusNotFound)
				})

				n8nClient, server := setupTestResourceClient(t, handler)
				defer server.Close()

				r := &VariableResource{client: n8nClient}
				plan := &models.Resource{
					Key: types.StringValue("test-key"),
				}
				resp := &resource.CreateResponse{}

				variable := r.findCreatedVariable(context.Background(), plan, resp)

				assert.NotNil(t, variable)
				assert.False(t, resp.Diagnostics.HasError())

			case "error - API error":
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
				})

				n8nClient, server := setupTestResourceClient(t, handler)
				defer server.Close()

				r := &VariableResource{client: n8nClient}
				plan := &models.Resource{
					Key: types.StringValue("test-key"),
				}
				resp := &resource.CreateResponse{}

				variable := r.findCreatedVariable(context.Background(), plan, resp)

				assert.Nil(t, variable)
				assert.True(t, resp.Diagnostics.HasError())

			case "error - variable not found":
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/variables" && r.Method == http.MethodGet {
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusOK)
						json.NewEncoder(w).Encode(n8nsdk.VariableList{
							Data: []n8nsdk.Variable{},
						})
						return
					}
					w.WriteHeader(http.StatusNotFound)
				})

				n8nClient, server := setupTestResourceClient(t, handler)
				defer server.Close()

				r := &VariableResource{client: n8nClient}
				plan := &models.Resource{
					Key: types.StringValue("test-key"),
				}
				resp := &resource.CreateResponse{}

				variable := r.findCreatedVariable(context.Background(), plan, resp)

				assert.Nil(t, variable)
				assert.True(t, resp.Diagnostics.HasError())
			}
		})
	}
}

// TestVariableResource_Read is tested via external tests and helper method tests.
// This ensures coverage of the Read method through the private methods it calls.

// TestVariableResource_findVariableByID tests the findVariableByID method.
func TestVariableResource_findVariableByID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "success - variable found",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := &VariableResource{}
				variableList := &n8nsdk.VariableList{
					Data: []n8nsdk.Variable{
						{Id: ptrString("var-123"), Key: "key1", Value: "value1"},
						{Id: ptrString("var-456"), Key: "key2", Value: "value2"},
					},
				}

				variable := r.findVariableByID(variableList, "var-123")

				assert.NotNil(t, variable)
				assert.Equal(t, "var-123", *variable.Id)
			},
		},
		{
			name: "variable not found",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := &VariableResource{}
				variableList := &n8nsdk.VariableList{
					Data: []n8nsdk.Variable{
						{Id: ptrString("var-456"), Key: "key2", Value: "value2"},
					},
				}

				variable := r.findVariableByID(variableList, "var-123")

				assert.Nil(t, variable)
			},
		},
		{
			name: "nil data",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := &VariableResource{}
				variableList := &n8nsdk.VariableList{
					Data: nil,
				}

				variable := r.findVariableByID(variableList, "var-123")

				assert.Nil(t, variable)
			},
		},
		{
			name: "error case - empty list",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := &VariableResource{}
				variableList := &n8nsdk.VariableList{
					Data: []n8nsdk.Variable{},
				}

				variable := r.findVariableByID(variableList, "var-123")

				assert.Nil(t, variable)
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

// TestVariableResource_updateStateFromVariable tests the updateStateFromVariable method.
func TestVariableResource_updateStateFromVariable(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "complete variable data",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := &VariableResource{}
				typeStr := "string"
				projectID := "proj-123"
				variable := &n8nsdk.Variable{
					Id:    ptrString("var-123"),
					Key:   "test-key",
					Value: "test-value",
					Type:  &typeStr,
					Project: &n8nsdk.Project{
						Id: &projectID,
					},
				}
				state := &models.Resource{}

				r.updateStateFromVariable(variable, state)

				assert.Equal(t, "test-key", state.Key.ValueString())
				assert.Equal(t, "test-value", state.Value.ValueString())
				assert.Equal(t, "string", state.Type.ValueString())
				assert.Equal(t, "proj-123", state.ProjectID.ValueString())
			},
		},
		{
			name: "minimal variable data",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := &VariableResource{}
				variable := &n8nsdk.Variable{
					Key:   "test-key",
					Value: "test-value",
				}
				state := &models.Resource{}

				r.updateStateFromVariable(variable, state)

				assert.Equal(t, "test-key", state.Key.ValueString())
				assert.Equal(t, "test-value", state.Value.ValueString())
			},
		},
		{
			name: "error case - nil project",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := &VariableResource{}
				variable := &n8nsdk.Variable{
					Key:     "test-key",
					Value:   "test-value",
					Project: nil,
				}
				state := &models.Resource{}

				r.updateStateFromVariable(variable, state)

				assert.Equal(t, "test-key", state.Key.ValueString())
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

// TestVariableResource_Update is tested via external tests and helper method tests.
// This ensures coverage of the Update method through the private methods it calls.

// TestVariableResource_executeVariableUpdate tests the executeVariableUpdate method.
func TestVariableResource_executeVariableUpdate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "success - update variable"},
		{name: "error - API error", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "success - update variable":
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/variables/var-123" && r.Method == http.MethodPut {
						w.WriteHeader(http.StatusNoContent)
						return
					}
					w.WriteHeader(http.StatusNotFound)
				})

				n8nClient, server := setupTestResourceClient(t, handler)
				defer server.Close()

				r := &VariableResource{client: n8nClient}
				plan := &models.Resource{
					ID:    types.StringValue("var-123"),
					Key:   types.StringValue("test-key"),
					Value: types.StringValue("test-value"),
				}
				resp := &resource.UpdateResponse{}

				result := r.executeVariableUpdate(context.Background(), plan, resp)

				assert.True(t, result)
				assert.False(t, resp.Diagnostics.HasError())

			case "error - API error":
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
				})

				n8nClient, server := setupTestResourceClient(t, handler)
				defer server.Close()

				r := &VariableResource{client: n8nClient}
				plan := &models.Resource{
					ID:    types.StringValue("var-123"),
					Key:   types.StringValue("test-key"),
					Value: types.StringValue("test-value"),
				}
				resp := &resource.UpdateResponse{}

				result := r.executeVariableUpdate(context.Background(), plan, resp)

				assert.False(t, result)
				assert.True(t, resp.Diagnostics.HasError())
			}
		})
	}
}

// TestVariableResource_findUpdatedVariable tests the findUpdatedVariable method.
func TestVariableResource_findUpdatedVariable(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "success - variable found"},
		{name: "error - API error", wantErr: true},
		{name: "error - variable not found", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "success - variable found":
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/variables" && r.Method == http.MethodGet {
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusOK)
						json.NewEncoder(w).Encode(n8nsdk.VariableList{
							Data: []n8nsdk.Variable{
								{Id: ptrString("var-123"), Key: "test-key", Value: "test-value"},
							},
						})
						return
					}
					w.WriteHeader(http.StatusNotFound)
				})

				n8nClient, server := setupTestResourceClient(t, handler)
				defer server.Close()

				r := &VariableResource{client: n8nClient}
				plan := &models.Resource{
					ID: types.StringValue("var-123"),
				}
				resp := &resource.UpdateResponse{}

				variable := r.findUpdatedVariable(context.Background(), plan, resp)

				assert.NotNil(t, variable)
				assert.False(t, resp.Diagnostics.HasError())

			case "error - API error":
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
				})

				n8nClient, server := setupTestResourceClient(t, handler)
				defer server.Close()

				r := &VariableResource{client: n8nClient}
				plan := &models.Resource{
					ID: types.StringValue("var-123"),
				}
				resp := &resource.UpdateResponse{}

				variable := r.findUpdatedVariable(context.Background(), plan, resp)

				assert.Nil(t, variable)
				assert.True(t, resp.Diagnostics.HasError())

			case "error - variable not found":
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/variables" && r.Method == http.MethodGet {
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusOK)
						json.NewEncoder(w).Encode(n8nsdk.VariableList{
							Data: []n8nsdk.Variable{},
						})
						return
					}
					w.WriteHeader(http.StatusNotFound)
				})

				n8nClient, server := setupTestResourceClient(t, handler)
				defer server.Close()

				r := &VariableResource{client: n8nClient}
				plan := &models.Resource{
					ID: types.StringValue("var-123"),
				}
				resp := &resource.UpdateResponse{}

				variable := r.findUpdatedVariable(context.Background(), plan, resp)

				assert.Nil(t, variable)
				assert.True(t, resp.Diagnostics.HasError())
			}
		})
	}
}

// TestVariableResource_Delete_CRUD tests the full Delete method execution.
func TestVariableResource_Delete_CRUD(t *testing.T) {
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

				n8nClient, server := setupTestResourceClient(t, handler)
				defer server.Close()

				r := NewVariableResource()
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
				assert.False(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "error - delete with invalid state",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := NewVariableResource()
				ctx := context.Background()

				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

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
				assert.True(t, resp.Diagnostics.HasError())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}

// TestVariableResource_ImportState_CRUD tests the full ImportState method execution.
func TestVariableResource_ImportState_CRUD(t *testing.T) {
	t.Run("import state passthrough", func(t *testing.T) {
		t.Helper()
		r := NewVariableResource()
		ctx := context.Background()

		schemaResp := resource.SchemaResponse{}
		r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

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

		r.ImportState(ctx, req, resp)
		assert.False(t, resp.Diagnostics.HasError())
	})
}
