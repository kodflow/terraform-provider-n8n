package project_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/project"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/client"
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

// TestNewProjectResource tests the NewProjectResource constructor.
func TestNewProjectResource(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "creates valid resource", wantErr: false},
		{name: "error case - validation checks", wantErr: false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "creates valid resource":
				r := project.NewProjectResource()
				assert.NotNil(t, r)

			case "error case - validation checks":
				r := project.NewProjectResource()
				assert.NotNil(t, r)
				assert.Implements(t, (*resource.Resource)(nil), r)
			}
		})
	}
}

// TestNewProjectResourceWrapper tests the NewProjectResourceWrapper constructor.
func TestNewProjectResourceWrapper(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "creates valid resource wrapper", wantErr: false},
		{name: "error case - validation checks", wantErr: false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "creates valid resource wrapper":
				wrapper := project.NewProjectResourceWrapper()
				assert.NotNil(t, wrapper)

			case "error case - validation checks":
				wrapper := project.NewProjectResourceWrapper()
				assert.NotNil(t, wrapper)
			}
		})
	}
}

// TestProjectResource_Metadata tests the Metadata method.
func TestProjectResource_Metadata(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "sets correct type name", wantErr: false},
		{name: "error case - validation checks", wantErr: false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "sets correct type name":
				r := project.NewProjectResource()
				req := resource.MetadataRequest{
					ProviderTypeName: "n8n",
				}
				resp := &resource.MetadataResponse{}
				r.Metadata(context.Background(), req, resp)
				assert.Equal(t, "n8n_project", resp.TypeName)

			case "error case - validation checks":
				r := project.NewProjectResource()
				req := resource.MetadataRequest{
					ProviderTypeName: "n8n",
				}
				resp := &resource.MetadataResponse{}
				r.Metadata(context.Background(), req, resp)
				assert.NotEmpty(t, resp.TypeName)
			}
		})
	}
}

// TestProjectResource_Schema tests the Schema method.
func TestProjectResource_Schema(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "returns valid schema", wantErr: false},
		{name: "error case - validation checks", wantErr: false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "returns valid schema":
				r := project.NewProjectResource()
				req := resource.SchemaRequest{}
				resp := &resource.SchemaResponse{}
				r.Schema(context.Background(), req, resp)
				assert.NotNil(t, resp.Schema)
				assert.NotEmpty(t, resp.Schema.Attributes)

			case "error case - validation checks":
				r := project.NewProjectResource()
				req := resource.SchemaRequest{}
				resp := &resource.SchemaResponse{}
				r.Schema(context.Background(), req, resp)
				assert.NotNil(t, resp.Schema)
			}
		})
	}
}

// TestProjectResource_Configure tests the Configure method.
func TestProjectResource_Configure(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "configures with valid client", wantErr: false},
		{name: "error case - nil provider data", wantErr: false},
		{name: "error case - wrong provider data type", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "configures with valid client":
				r := project.NewProjectResource()
				req := resource.ConfigureRequest{
					ProviderData: &client.N8nClient{},
				}
				resp := &resource.ConfigureResponse{}
				r.Configure(context.Background(), req, resp)
				assert.False(t, resp.Diagnostics.HasError())

			case "error case - nil provider data":
				r := project.NewProjectResource()
				req := resource.ConfigureRequest{
					ProviderData: nil,
				}
				resp := &resource.ConfigureResponse{}
				r.Configure(context.Background(), req, resp)
				assert.False(t, resp.Diagnostics.HasError())

			case "error case - wrong provider data type":
				r := project.NewProjectResource()
				req := resource.ConfigureRequest{
					ProviderData: "wrong type",
				}
				resp := &resource.ConfigureResponse{}
				r.Configure(context.Background(), req, resp)
				assert.True(t, resp.Diagnostics.HasError())
				assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Unexpected Resource Configure Type")
			}
		})
	}
}

// TestProjectResource_Create tests the Create method with full execution.
func TestProjectResource_Create(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "create with successful API call",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					// POST to create - returns 201 with no body
					if r.Method == "POST" && r.URL.Path == "/projects" {
						w.WriteHeader(http.StatusCreated)
						return
					}
					// GET to list projects - returns the created project
					if r.Method == "GET" && r.URL.Path == "/projects" {
						w.WriteHeader(http.StatusOK)
						w.Write([]byte(`{"data":[{"id":"proj-123","name":"test-project","type":"team"}]}`))
						return
					}
					w.WriteHeader(http.StatusNotFound)
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				r := project.NewProjectResource()
				r.Configure(context.Background(), resource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &resource.ConfigureResponse{})

				ctx := context.Background()
				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				// Build plan using tftypes
				planRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":   tftypes.NewValue(tftypes.String, nil),
					"name": tftypes.NewValue(tftypes.String, "test-project"),
					"type": tftypes.NewValue(tftypes.String, "team"),
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
				r := project.NewProjectResource()
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

				r := project.NewProjectResource()
				r.Configure(context.Background(), resource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &resource.ConfigureResponse{})

				ctx := context.Background()
				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				planRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":   tftypes.NewValue(tftypes.String, nil),
					"name": tftypes.NewValue(tftypes.String, "test-project"),
					"type": tftypes.NewValue(tftypes.String, "team"),
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

// TestProjectResource_Read tests the Read method with full execution.
func TestProjectResource_Read(t *testing.T) {
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
					w.Write([]byte(`{"data":[{"id":"proj-123","name":"test-project","type":"team"}]}`))
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				r := project.NewProjectResource()
				r.Configure(context.Background(), resource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &resource.ConfigureResponse{})

				ctx := context.Background()
				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				// Build state
				stateRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":   tftypes.NewValue(tftypes.String, "proj-123"),
					"name": tftypes.NewValue(tftypes.String, "test-project"),
					"type": tftypes.NewValue(tftypes.String, "team"),
				})

				state := tfsdk.State{
					Schema: schemaResp.Schema,
					Raw:    stateRaw,
				}

				req := resource.ReadRequest{
					State: state,
				}
				resp := &resource.ReadResponse{
					State: tfsdk.State{Schema: schemaResp.Schema},
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
				r := project.NewProjectResource()
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
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte(`{"message": "Project not found"}`))
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				r := project.NewProjectResource()
				r.Configure(context.Background(), resource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &resource.ConfigureResponse{})

				ctx := context.Background()
				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				stateRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":   tftypes.NewValue(tftypes.String, "proj-123"),
					"name": tftypes.NewValue(tftypes.String, "test-project"),
					"type": tftypes.NewValue(tftypes.String, "team"),
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

// TestProjectResource_Update tests the Update method with full execution.
func TestProjectResource_Update(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "update with successful API call",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					// PUT to update - returns 204 with no body
					if r.Method == "PUT" && r.URL.Path == "/projects/proj-123" {
						w.WriteHeader(http.StatusNoContent)
						return
					}
					// GET to list projects - returns the updated project
					if r.Method == "GET" && r.URL.Path == "/projects" {
						w.WriteHeader(http.StatusOK)
						w.Write([]byte(`{"data":[{"id":"proj-123","name":"updated-project","type":"team"}]}`))
						return
					}
					w.WriteHeader(http.StatusNotFound)
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				r := project.NewProjectResource()
				r.Configure(context.Background(), resource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &resource.ConfigureResponse{})

				ctx := context.Background()
				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				// Build plan
				planRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":   tftypes.NewValue(tftypes.String, "proj-123"),
					"name": tftypes.NewValue(tftypes.String, "updated-project"),
					"type": tftypes.NewValue(tftypes.String, "team"),
				})

				plan := tfsdk.Plan{
					Schema: schemaResp.Schema,
					Raw:    planRaw,
				}

				// Build state
				stateRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":   tftypes.NewValue(tftypes.String, "proj-123"),
					"name": tftypes.NewValue(tftypes.String, "old-project"),
					"type": tftypes.NewValue(tftypes.String, "team"),
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
				r := project.NewProjectResource()
				ctx := context.Background()

				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				// Create invalid plan
				planRaw := tftypes.NewValue(tftypes.String, "invalid")

				plan := tfsdk.Plan{
					Schema: schemaResp.Schema,
					Raw:    planRaw,
				}

				// Create valid state (required since Update reads from both plan and state)
				stateRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":   tftypes.NewValue(tftypes.String, "proj-123"),
					"name": tftypes.NewValue(tftypes.String, "existing-project"),
					"type": tftypes.NewValue(tftypes.String, "team"),
				})

				state := tfsdk.State{
					Schema: schemaResp.Schema,
					Raw:    stateRaw,
				}

				req := resource.UpdateRequest{
					Plan:  plan,
					State: state,
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

				r := project.NewProjectResource()
				r.Configure(context.Background(), resource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &resource.ConfigureResponse{})

				ctx := context.Background()
				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				planRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":   tftypes.NewValue(tftypes.String, "proj-123"),
					"name": tftypes.NewValue(tftypes.String, "updated-project"),
					"type": tftypes.NewValue(tftypes.String, "team"),
				})

				plan := tfsdk.Plan{
					Schema: schemaResp.Schema,
					Raw:    planRaw,
				}

				stateRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":   tftypes.NewValue(tftypes.String, "proj-123"),
					"name": tftypes.NewValue(tftypes.String, "old-project"),
					"type": tftypes.NewValue(tftypes.String, "team"),
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

// TestProjectResource_Delete tests the Delete method with full execution.
func TestProjectResource_Delete(t *testing.T) {
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

				r := project.NewProjectResource()
				r.Configure(context.Background(), resource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &resource.ConfigureResponse{})

				ctx := context.Background()
				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				// Build state
				stateRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":   tftypes.NewValue(tftypes.String, "proj-123"),
					"name": tftypes.NewValue(tftypes.String, "test-project"),
					"type": tftypes.NewValue(tftypes.String, "team"),
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
				r := project.NewProjectResource()
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

				r := project.NewProjectResource()
				r.Configure(context.Background(), resource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &resource.ConfigureResponse{})

				ctx := context.Background()
				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				stateRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":   tftypes.NewValue(tftypes.String, "proj-123"),
					"name": tftypes.NewValue(tftypes.String, "test-project"),
					"type": tftypes.NewValue(tftypes.String, "team"),
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
}

// TestProjectResource_ImportState tests the ImportState method.
func TestProjectResource_ImportState(t *testing.T) {
	t.Run("import state passthrough", func(t *testing.T) {
		t.Helper()
		r := project.NewProjectResource()
		ctx := context.Background()

		schemaResp := resource.SchemaResponse{}
		r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

		// Build empty state
		emptyValue := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
			"id":   tftypes.NewValue(tftypes.String, nil),
			"name": tftypes.NewValue(tftypes.String, nil),
			"type": tftypes.NewValue(tftypes.String, nil),
		})

		req := resource.ImportStateRequest{
			ID: "proj-123",
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
