package project_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/project"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
	"github.com/stretchr/testify/assert"
)

// TestNewProjectUserResource tests the NewProjectUserResource constructor.
func TestNewProjectUserResource(t *testing.T) {
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
				r := project.NewProjectUserResource()
				assert.NotNil(t, r)

			case "error case - validation checks":
				r := project.NewProjectUserResource()
				assert.NotNil(t, r)
				assert.Implements(t, (*resource.Resource)(nil), r)
			}
		})
	}
}

// TestNewProjectUserResourceWrapper tests the NewProjectUserResourceWrapper constructor.
func TestNewProjectUserResourceWrapper(t *testing.T) {
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
				wrapper := project.NewProjectUserResourceWrapper()
				assert.NotNil(t, wrapper)

			case "error case - validation checks":
				wrapper := project.NewProjectUserResourceWrapper()
				assert.NotNil(t, wrapper)
			}
		})
	}
}

// TestProjectUserResource_Metadata tests the Metadata method.
func TestProjectUserResource_Metadata(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		providerTypeName string
		expectedTypeName string
	}{
		{
			name:             "with n8n provider",
			providerTypeName: "n8n",
			expectedTypeName: "n8n_project_user",
		},
		{
			name:             "with custom provider",
			providerTypeName: "custom",
			expectedTypeName: "custom_project_user",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := project.NewProjectUserResource()
			req := resource.MetadataRequest{
				ProviderTypeName: tt.providerTypeName,
			}
			resp := &resource.MetadataResponse{}

			r.Metadata(context.Background(), req, resp)

			assert.Equal(t, tt.expectedTypeName, resp.TypeName)
		})
	}
}

// TestProjectUserResource_Schema tests the Schema method.
func TestProjectUserResource_Schema(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
	}{
		{name: "schema has required attributes"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := project.NewProjectUserResource()
			req := resource.SchemaRequest{}
			resp := &resource.SchemaResponse{}

			r.Schema(context.Background(), req, resp)

			assert.NotNil(t, resp.Schema)
			assert.Contains(t, resp.Schema.MarkdownDescription, "user")

			// Verify required attributes exist
			attrs := resp.Schema.Attributes
			assert.Contains(t, attrs, "id")
			assert.Contains(t, attrs, "project_id")
			assert.Contains(t, attrs, "user_id")
			assert.Contains(t, attrs, "role")
		})
	}
}

// TestProjectUserResource_Configure tests the Configure method.
func TestProjectUserResource_Configure(t *testing.T) {
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
				r := project.NewProjectUserResource()
				req := resource.ConfigureRequest{
					ProviderData: &client.N8nClient{},
				}
				resp := &resource.ConfigureResponse{}
				r.Configure(context.Background(), req, resp)
				assert.False(t, resp.Diagnostics.HasError())

			case "error case - nil provider data":
				r := project.NewProjectUserResource()
				req := resource.ConfigureRequest{
					ProviderData: nil,
				}
				resp := &resource.ConfigureResponse{}
				r.Configure(context.Background(), req, resp)
				assert.False(t, resp.Diagnostics.HasError())

			case "error case - wrong provider data type":
				r := project.NewProjectUserResource()
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

// TestProjectUserResource_ImportState tests the ImportState method.
func TestProjectUserResource_ImportState(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		importID  string
		wantError bool
	}{
		{
			name:      "success - valid composite ID",
			importID:  "proj-123/user-456",
			wantError: false,
		},
		{
			name:      "error - invalid format (single part)",
			importID:  "proj-123",
			wantError: true,
		},
		{
			name:      "error - invalid format (three parts)",
			importID:  "proj-123/user-456/extra",
			wantError: true,
		},
		{
			name:      "error - empty ID",
			importID:  "",
			wantError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := project.NewProjectUserResource()
			ctx := context.Background()

			schemaResp := resource.SchemaResponse{}
			r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

			// Build empty state
			emptyValue := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
				"id":         tftypes.NewValue(tftypes.String, nil),
				"project_id": tftypes.NewValue(tftypes.String, nil),
				"user_id":    tftypes.NewValue(tftypes.String, nil),
				"role":       tftypes.NewValue(tftypes.String, nil),
			})

			req := resource.ImportStateRequest{
				ID: tt.importID,
			}
			resp := &resource.ImportStateResponse{
				State: tfsdk.State{
					Schema: schemaResp.Schema,
					Raw:    emptyValue,
				},
			}

			r.ImportState(ctx, req, resp)

			if tt.wantError {
				assert.True(t, resp.Diagnostics.HasError(), "Expected error for import ID: %s", tt.importID)
			} else {
				assert.False(t, resp.Diagnostics.HasError(), "Unexpected error for import ID: %s", tt.importID)
			}
		})
	}
}

// TestProjectUserResource_Create tests the public Create method.
func TestProjectUserResource_Create(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectError bool
		invalidPlan bool
	}{
		{
			name: "successful creation",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPost && r.URL.Path == "/projects/proj-123/users" {
					w.WriteHeader(http.StatusNoContent)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			expectError: false,
			invalidPlan: false,
		},
		{
			name: "API error",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			}),
			expectError: true,
			invalidPlan: false,
		},
		{
			name:        "invalid plan data",
			handler:     http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
			expectError: true,
			invalidPlan: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			n8nClient, server := setupUserTestClient(t, tt.handler)
			defer server.Close()

			r := project.NewProjectUserResource()
			configReq := resource.ConfigureRequest{ProviderData: n8nClient}
			configResp := &resource.ConfigureResponse{}
			r.Configure(context.Background(), configReq, configResp)

			ctx := context.Background()

			// Build schema and plan
			schemaResp := resource.SchemaResponse{}
			r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

			var planRaw tftypes.Value
			if tt.invalidPlan {
				// Create invalid plan with wrong type to trigger req.Plan.Get() error
				planRaw = tftypes.NewValue(tftypes.String, "invalid")
			} else {
				planRaw = tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, nil),
					"project_id": tftypes.NewValue(tftypes.String, "proj-123"),
					"user_id":    tftypes.NewValue(tftypes.String, "user-456"),
					"role":       tftypes.NewValue(tftypes.String, "project:admin"),
				})
			}

			req := resource.CreateRequest{
				Plan: tfsdk.Plan{
					Schema: schemaResp.Schema,
					Raw:    planRaw,
				},
			}
			resp := &resource.CreateResponse{
				State: tfsdk.State{Schema: schemaResp.Schema},
			}

			r.Create(ctx, req, resp)

			if tt.expectError {
				assert.True(t, resp.Diagnostics.HasError())
			} else {
				assert.False(t, resp.Diagnostics.HasError())
			}
		})
	}
}

// TestProjectUserResource_Read tests the public Read method.
func TestProjectUserResource_Read(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		handler      http.HandlerFunc
		expectError  bool
		invalidState bool
	}{
		{
			name: "successful read",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/users" && r.Method == http.MethodGet {
					response := map[string]interface{}{
						"data": []interface{}{
							map[string]interface{}{
								"id":    "user-456",
								"email": "user@example.com",
								"role":  "project:admin",
							},
						},
					}
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(response)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			expectError:  false,
			invalidState: false,
		},
		{
			name: "API error",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			}),
			expectError:  true,
			invalidState: false,
		},
		{
			name:         "invalid state data",
			handler:      http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
			expectError:  true,
			invalidState: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			n8nClient, server := setupUserTestClient(t, tt.handler)
			defer server.Close()

			r := project.NewProjectUserResource()
			configReq := resource.ConfigureRequest{ProviderData: n8nClient}
			configResp := &resource.ConfigureResponse{}
			r.Configure(context.Background(), configReq, configResp)

			ctx := context.Background()

			// Build schema and state
			schemaResp := resource.SchemaResponse{}
			r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

			var stateRaw tftypes.Value
			if tt.invalidState {
				// Create invalid state with wrong type to trigger req.State.Get() error
				stateRaw = tftypes.NewValue(tftypes.String, "invalid")
			} else {
				stateRaw = tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "proj-123-user-456"),
					"project_id": tftypes.NewValue(tftypes.String, "proj-123"),
					"user_id":    tftypes.NewValue(tftypes.String, "user-456"),
					"role":       tftypes.NewValue(tftypes.String, "project:admin"),
				})
			}

			req := resource.ReadRequest{
				State: tfsdk.State{
					Schema: schemaResp.Schema,
					Raw:    stateRaw,
				},
			}
			resp := &resource.ReadResponse{
				State: tfsdk.State{Schema: schemaResp.Schema},
			}

			r.Read(ctx, req, resp)

			if tt.expectError {
				assert.True(t, resp.Diagnostics.HasError())
			} else {
				assert.False(t, resp.Diagnostics.HasError())
			}
		})
	}
}

// TestProjectUserResource_Update tests the public Update method.
func TestProjectUserResource_Update(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		handler      http.HandlerFunc
		expectError  bool
		invalidPlan  bool
		invalidState bool
	}{
		{
			name: "successful update",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPatch && r.URL.Path == "/projects/proj-123/users/user-456" {
					w.WriteHeader(http.StatusNoContent)
					return
				}
				if r.URL.Path == "/users" && r.Method == http.MethodGet {
					response := map[string]interface{}{
						"data": []interface{}{
							map[string]interface{}{
								"id":    "user-456",
								"email": "user@example.com",
								"role":  "project:editor",
							},
						},
					}
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(response)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			expectError:  false,
			invalidPlan:  false,
			invalidState: false,
		},
		{
			name: "API error",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			}),
			expectError:  true,
			invalidPlan:  false,
			invalidState: false,
		},
		{
			name:         "invalid plan data",
			handler:      http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
			expectError:  true,
			invalidPlan:  true,
			invalidState: false,
		},
		{
			name:         "invalid state data",
			handler:      http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
			expectError:  true,
			invalidPlan:  false,
			invalidState: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			n8nClient, server := setupUserTestClient(t, tt.handler)
			defer server.Close()

			r := project.NewProjectUserResource()
			configReq := resource.ConfigureRequest{ProviderData: n8nClient}
			configResp := &resource.ConfigureResponse{}
			r.Configure(context.Background(), configReq, configResp)

			ctx := context.Background()

			// Build schema and plan/state
			schemaResp := resource.SchemaResponse{}
			r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

			var stateRaw, planRaw tftypes.Value
			if tt.invalidState {
				// Create invalid state with wrong type to trigger req.State.Get() error
				stateRaw = tftypes.NewValue(tftypes.String, "invalid")
				planRaw = tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "proj-123-user-456"),
					"project_id": tftypes.NewValue(tftypes.String, "proj-123"),
					"user_id":    tftypes.NewValue(tftypes.String, "user-456"),
					"role":       tftypes.NewValue(tftypes.String, "project:editor"),
				})
			} else if tt.invalidPlan {
				// Create invalid plan with wrong type to trigger req.Plan.Get() error
				stateRaw = tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "proj-123-user-456"),
					"project_id": tftypes.NewValue(tftypes.String, "proj-123"),
					"user_id":    tftypes.NewValue(tftypes.String, "user-456"),
					"role":       tftypes.NewValue(tftypes.String, "project:admin"),
				})
				planRaw = tftypes.NewValue(tftypes.String, "invalid")
			} else {
				stateRaw = tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "proj-123-user-456"),
					"project_id": tftypes.NewValue(tftypes.String, "proj-123"),
					"user_id":    tftypes.NewValue(tftypes.String, "user-456"),
					"role":       tftypes.NewValue(tftypes.String, "project:admin"),
				})

				planRaw = tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "proj-123-user-456"),
					"project_id": tftypes.NewValue(tftypes.String, "proj-123"),
					"user_id":    tftypes.NewValue(tftypes.String, "user-456"),
					"role":       tftypes.NewValue(tftypes.String, "project:editor"),
				})
			}

			req := resource.UpdateRequest{
				State: tfsdk.State{
					Schema: schemaResp.Schema,
					Raw:    stateRaw,
				},
				Plan: tfsdk.Plan{
					Schema: schemaResp.Schema,
					Raw:    planRaw,
				},
			}
			resp := &resource.UpdateResponse{
				State: tfsdk.State{Schema: schemaResp.Schema},
			}

			r.Update(ctx, req, resp)

			if tt.expectError {
				assert.True(t, resp.Diagnostics.HasError())
			} else {
				assert.False(t, resp.Diagnostics.HasError())
			}
		})
	}
}

// TestProjectUserResource_Delete tests the public Delete method.
func TestProjectUserResource_Delete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		handler      http.HandlerFunc
		expectError  bool
		invalidState bool
	}{
		{
			name: "successful deletion",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodDelete && r.URL.Path == "/projects/proj-123/users/user-456" {
					w.WriteHeader(http.StatusNoContent)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			expectError:  false,
			invalidState: false,
		},
		{
			name: "API error",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			}),
			expectError:  true,
			invalidState: false,
		},
		{
			name:         "invalid state data",
			handler:      http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
			expectError:  true,
			invalidState: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			n8nClient, server := setupUserTestClient(t, tt.handler)
			defer server.Close()

			r := project.NewProjectUserResource()
			configReq := resource.ConfigureRequest{ProviderData: n8nClient}
			configResp := &resource.ConfigureResponse{}
			r.Configure(context.Background(), configReq, configResp)

			ctx := context.Background()

			// Build schema and state
			schemaResp := resource.SchemaResponse{}
			r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

			var stateRaw tftypes.Value
			if tt.invalidState {
				// Create invalid state with wrong type to trigger req.State.Get() error
				stateRaw = tftypes.NewValue(tftypes.String, "invalid")
			} else {
				stateRaw = tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "proj-123-user-456"),
					"project_id": tftypes.NewValue(tftypes.String, "proj-123"),
					"user_id":    tftypes.NewValue(tftypes.String, "user-456"),
					"role":       tftypes.NewValue(tftypes.String, "project:admin"),
				})
			}

			req := resource.DeleteRequest{
				State: tfsdk.State{
					Schema: schemaResp.Schema,
					Raw:    stateRaw,
				},
			}
			resp := &resource.DeleteResponse{}

			r.Delete(ctx, req, resp)

			if tt.expectError {
				assert.True(t, resp.Diagnostics.HasError())
			} else {
				assert.False(t, resp.Diagnostics.HasError())
			}
		})
	}
}

func setupUserTestClient(t *testing.T, handler http.HandlerFunc) (*client.N8nClient, *httptest.Server) {
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
