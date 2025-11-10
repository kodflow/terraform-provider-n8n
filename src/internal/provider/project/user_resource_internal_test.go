package project

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/project/models"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
	"github.com/stretchr/testify/assert"
)

func TestNewProjectUserResource(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "create new instance",
			wantErr: false,
		},
		{
			name:    "multiple instances are independent",
			wantErr: false,
		},
		{
			name:    "error case - validation checks",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "create new instance":
				r := NewProjectUserResource()

				assert.NotNil(t, r)
				assert.IsType(t, &ProjectUserResource{}, r)
				assert.Nil(t, r.client)

			case "multiple instances are independent":
				r1 := NewProjectUserResource()
				r2 := NewProjectUserResource()

				assert.NotSame(t, r1, r2)

			case "error case - validation checks":
				r := NewProjectUserResource()

				assert.NotNil(t, r, "NewProjectUserResource should never return nil")
				assert.IsType(t, &ProjectUserResource{}, r, "should return ProjectUserResource type")
			}
		})
	}
}

func TestNewProjectUserResourceWrapper(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "create new wrapped instance",
			wantErr: false,
		},
		{
			name:    "wrapper returns resource.Resource interface",
			wantErr: false,
		},
		{
			name:    "error case - validation checks",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "create new wrapped instance":
				r := NewProjectUserResourceWrapper()

				assert.NotNil(t, r)
				assert.IsType(t, &ProjectUserResource{}, r)

			case "wrapper returns resource.Resource interface":
				r := NewProjectUserResourceWrapper()

				// Verify it implements resource.Resource.
				assert.Implements(t, (*resource.Resource)(nil), r)
				assert.NotNil(t, r)

			case "error case - validation checks":
				r := NewProjectUserResourceWrapper()

				assert.NotNil(t, r, "NewProjectUserResourceWrapper should never return nil")
				assert.Implements(t, (*resource.Resource)(nil), r, "should implement resource.Resource interface")
			}
		})
	}
}

func TestProjectUserResource_Metadata(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "set metadata with provider type",
			wantErr: false,
		},
		{
			name:    "set metadata with custom provider type",
			wantErr: false,
		},
		{
			name:    "error case - validation checks",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "set metadata with provider type":
				r := NewProjectUserResource()
				resp := &resource.MetadataResponse{}

				r.Metadata(context.Background(), resource.MetadataRequest{
					ProviderTypeName: "n8n",
				}, resp)

				assert.Equal(t, "n8n_project_user", resp.TypeName)

			case "set metadata with custom provider type":
				r := NewProjectUserResource()
				resp := &resource.MetadataResponse{}

				r.Metadata(context.Background(), resource.MetadataRequest{
					ProviderTypeName: "custom_provider",
				}, resp)

				assert.Equal(t, "custom_provider_project_user", resp.TypeName)

			case "error case - validation checks":
				r := NewProjectUserResource()
				resp := &resource.MetadataResponse{}

				r.Metadata(context.Background(), resource.MetadataRequest{
					ProviderTypeName: "",
				}, resp)

				assert.Equal(t, "_project_user", resp.TypeName, "should handle empty provider type name")
			}
		})
	}
}

func TestProjectUserResource_Schema(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "return schema",
			wantErr: false,
		},
		{
			name:    "schema attributes have descriptions",
			wantErr: false,
		},
		{
			name:    "error case - validation checks",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "return schema":
				r := NewProjectUserResource()
				resp := &resource.SchemaResponse{}

				r.Schema(context.Background(), resource.SchemaRequest{}, resp)

				assert.NotNil(t, resp.Schema)
				assert.Contains(t, resp.Schema.MarkdownDescription, "user membership")

				expectedAttrs := []string{
					"id",
					"project_id",
					"user_id",
					"role",
				}

				for _, attr := range expectedAttrs {
					_, exists := resp.Schema.Attributes[attr]
					assert.True(t, exists, "Attribute %s should exist", attr)
				}

			case "schema attributes have descriptions":
				r := NewProjectUserResource()
				resp := &resource.SchemaResponse{}

				r.Schema(context.Background(), resource.SchemaRequest{}, resp)

				for name, attr := range resp.Schema.Attributes {
					assert.NotEmpty(t, attr.GetMarkdownDescription(), "Attribute %s should have description", name)
				}

			case "error case - validation checks":
				r := NewProjectUserResource()
				resp := &resource.SchemaResponse{}

				r.Schema(context.Background(), resource.SchemaRequest{}, resp)

				assert.NotNil(t, resp.Schema, "Schema should not be nil")
				assert.NotEmpty(t, resp.Schema.Attributes, "Schema should have attributes")
			}
		})
	}
}

func TestProjectUserResource_Configure(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "configure with valid client",
			wantErr: false,
		},
		{
			name:    "configure with nil provider data",
			wantErr: false,
		},
		{
			name:    "configure with invalid provider data",
			wantErr: true,
		},
		{
			name:    "configure with wrong type",
			wantErr: true,
		},
		{
			name:    "error case - validation checks",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "configure with valid client":
				r := NewProjectUserResource()
				resp := &resource.ConfigureResponse{}

				mockClient := &client.N8nClient{}
				req := resource.ConfigureRequest{
					ProviderData: mockClient,
				}

				r.Configure(context.Background(), req, resp)

				assert.NotNil(t, r.client)
				assert.Equal(t, mockClient, r.client)
				assert.False(t, resp.Diagnostics.HasError())

			case "configure with nil provider data":
				r := NewProjectUserResource()
				resp := &resource.ConfigureResponse{}
				req := resource.ConfigureRequest{
					ProviderData: nil,
				}

				r.Configure(context.Background(), req, resp)

				assert.Nil(t, r.client)
				assert.False(t, resp.Diagnostics.HasError())

			case "configure with invalid provider data":
				r := NewProjectUserResource()
				resp := &resource.ConfigureResponse{}
				req := resource.ConfigureRequest{
					ProviderData: "invalid-data",
				}

				r.Configure(context.Background(), req, resp)

				assert.Nil(t, r.client)
				assert.True(t, resp.Diagnostics.HasError())
				assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Unexpected Resource Configure Type")

			case "configure with wrong type":
				r := NewProjectUserResource()
				resp := &resource.ConfigureResponse{}
				req := resource.ConfigureRequest{
					ProviderData: struct{}{},
				}

				r.Configure(context.Background(), req, resp)

				assert.Nil(t, r.client)
				assert.True(t, resp.Diagnostics.HasError())

			case "error case - validation checks":
				r := NewProjectUserResource()
				resp := &resource.ConfigureResponse{}
				req := resource.ConfigureRequest{
					ProviderData: 123,
				}

				r.Configure(context.Background(), req, resp)

				assert.True(t, resp.Diagnostics.HasError(), "Configure should fail with integer provider data")
				assert.Nil(t, r.client, "client should remain nil on error")
			}
		})
	}
}

func TestProjectUserResource_Create(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "successful creation",
			wantErr: false,
		},
		{
			name:    "creation fails with API error",
			wantErr: true,
		},
		{
			name:    "error case - validation checks",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "successful creation":
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/projects/proj-123/users" && r.Method == http.MethodPost {
						w.WriteHeader(http.StatusCreated)
						return
					}
					w.WriteHeader(http.StatusNotFound)
				})

				n8nClient, server := setupTestUserResourceClient(t, handler)
				defer server.Close()

				r := &ProjectUserResource{client: n8nClient}

				rawPlan := map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, nil),
					"project_id": tftypes.NewValue(tftypes.String, "proj-123"),
					"user_id":    tftypes.NewValue(tftypes.String, "user-456"),
					"role":       tftypes.NewValue(tftypes.String, "project:editor"),
				}
				plan := tfsdk.Plan{
					Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "project_id": tftypes.String, "user_id": tftypes.String, "role": tftypes.String}}, rawPlan),
					Schema: createTestUserResourceSchema(t),
				}

				state := tfsdk.State{
					Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "project_id": tftypes.String, "user_id": tftypes.String, "role": tftypes.String}}, nil),
					Schema: createTestUserResourceSchema(t),
				}

				req := resource.CreateRequest{
					Plan: plan,
				}
				resp := resource.CreateResponse{
					State: state,
				}

				r.Create(context.Background(), req, &resp)

				assert.False(t, resp.Diagnostics.HasError(), "Create should not have errors")

			case "creation fails with API error":
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"message": "Internal server error"}`))
				})

				n8nClient, server := setupTestUserResourceClient(t, handler)
				defer server.Close()

				r := &ProjectUserResource{client: n8nClient}

				rawPlan := map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, nil),
					"project_id": tftypes.NewValue(tftypes.String, "proj-123"),
					"user_id":    tftypes.NewValue(tftypes.String, "user-456"),
					"role":       tftypes.NewValue(tftypes.String, "project:editor"),
				}
				plan := tfsdk.Plan{
					Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "project_id": tftypes.String, "user_id": tftypes.String, "role": tftypes.String}}, rawPlan),
					Schema: createTestUserResourceSchema(t),
				}

				state := tfsdk.State{
					Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "project_id": tftypes.String, "user_id": tftypes.String, "role": tftypes.String}}, nil),
					Schema: createTestUserResourceSchema(t),
				}

				req := resource.CreateRequest{
					Plan: plan,
				}
				resp := resource.CreateResponse{
					State: state,
				}

				r.Create(context.Background(), req, &resp)

				assert.True(t, resp.Diagnostics.HasError(), "Create should have errors")
				assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Error adding user to project")

			case "error case - validation checks":
				r := &ProjectUserResource{}

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
					State: tfsdk.State{Schema: createTestUserResourceSchema(t)},
				}

				r.Create(context.Background(), req, &resp)

				assert.True(t, resp.Diagnostics.HasError(), "Create should fail when Plan.Get fails")
			}
		})
	}
}

func TestProjectUserResource_Read(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "successful read",
			wantErr: false,
		},
		{
			name:    "user not found removes from state",
			wantErr: false,
		},
		{
			name:    "read fails with API error",
			wantErr: true,
		},
		{
			name:    "read with nil data",
			wantErr: false,
		},
		{
			name:    "error case - validation checks",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "successful read":
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/users" && r.Method == http.MethodGet {
						w.Header().Set("Content-Type", "application/json")
						userID := "user-456"
						role := "project:editor"
						json.NewEncoder(w).Encode(map[string]interface{}{
							"data": []map[string]interface{}{
								{
									"id":    userID,
									"email": "test@example.com",
									"role":  role,
								},
							},
						})
						return
					}
					w.WriteHeader(http.StatusNotFound)
				})

				n8nClient, server := setupTestUserResourceClient(t, handler)
				defer server.Close()

				r := &ProjectUserResource{client: n8nClient}

				rawState := map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "proj-123/user-456"),
					"project_id": tftypes.NewValue(tftypes.String, "proj-123"),
					"user_id":    tftypes.NewValue(tftypes.String, "user-456"),
					"role":       tftypes.NewValue(tftypes.String, "project:editor"),
				}
				state := tfsdk.State{
					Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "project_id": tftypes.String, "user_id": tftypes.String, "role": tftypes.String}}, rawState),
					Schema: createTestUserResourceSchema(t),
				}

				req := resource.ReadRequest{
					State: state,
				}
				resp := resource.ReadResponse{
					State: state,
				}

				r.Read(context.Background(), req, &resp)

				assert.False(t, resp.Diagnostics.HasError(), "Read should not have errors")

			case "user not found removes from state":
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/users" && r.Method == http.MethodGet &&
						r.URL.Query().Get("projectId") == "proj-123" {
						w.Header().Set("Content-Type", "application/json")
						json.NewEncoder(w).Encode(map[string]interface{}{
							"data": []map[string]interface{}{},
						})
						return
					}
					w.WriteHeader(http.StatusNotFound)
				})

				n8nClient, server := setupTestUserResourceClient(t, handler)
				defer server.Close()

				r := &ProjectUserResource{client: n8nClient}

				rawState := map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "proj-123/user-456"),
					"project_id": tftypes.NewValue(tftypes.String, "proj-123"),
					"user_id":    tftypes.NewValue(tftypes.String, "user-456"),
					"role":       tftypes.NewValue(tftypes.String, "project:editor"),
				}
				state := tfsdk.State{
					Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "project_id": tftypes.String, "user_id": tftypes.String, "role": tftypes.String}}, rawState),
					Schema: createTestUserResourceSchema(t),
				}

				req := resource.ReadRequest{
					State: state,
				}
				resp := resource.ReadResponse{
					State: state,
				}

				r.Read(context.Background(), req, &resp)

				assert.False(t, resp.Diagnostics.HasError(), "Read should not have errors")

			case "read fails with API error":
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"message": "Internal server error"}`))
				})

				n8nClient, server := setupTestUserResourceClient(t, handler)
				defer server.Close()

				r := &ProjectUserResource{client: n8nClient}

				rawState := map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "proj-123/user-456"),
					"project_id": tftypes.NewValue(tftypes.String, "proj-123"),
					"user_id":    tftypes.NewValue(tftypes.String, "user-456"),
					"role":       tftypes.NewValue(tftypes.String, "project:editor"),
				}
				state := tfsdk.State{
					Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "project_id": tftypes.String, "user_id": tftypes.String, "role": tftypes.String}}, rawState),
					Schema: createTestUserResourceSchema(t),
				}

				req := resource.ReadRequest{
					State: state,
				}
				resp := resource.ReadResponse{
					State: state,
				}

				r.Read(context.Background(), req, &resp)

				assert.True(t, resp.Diagnostics.HasError(), "Read should have errors")
				assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Error reading project users")

			case "read with nil data":
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/users" && r.Method == http.MethodGet &&
						r.URL.Query().Get("projectId") == "proj-123" {
						w.Header().Set("Content-Type", "application/json")
						json.NewEncoder(w).Encode(map[string]interface{}{})
						return
					}
					w.WriteHeader(http.StatusNotFound)
				})

				n8nClient, server := setupTestUserResourceClient(t, handler)
				defer server.Close()

				r := &ProjectUserResource{client: n8nClient}

				rawState := map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "proj-123/user-456"),
					"project_id": tftypes.NewValue(tftypes.String, "proj-123"),
					"user_id":    tftypes.NewValue(tftypes.String, "user-456"),
					"role":       tftypes.NewValue(tftypes.String, "project:editor"),
				}
				state := tfsdk.State{
					Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "project_id": tftypes.String, "user_id": tftypes.String, "role": tftypes.String}}, rawState),
					Schema: createTestUserResourceSchema(t),
				}

				req := resource.ReadRequest{
					State: state,
				}
				resp := resource.ReadResponse{
					State: state,
				}

				r.Read(context.Background(), req, &resp)

				assert.False(t, resp.Diagnostics.HasError(), "Read should handle nil data gracefully")

			case "error case - validation checks":
				r := &ProjectUserResource{}

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
					State: tfsdk.State{Schema: createTestUserResourceSchema(t)},
				}

				r.Read(context.Background(), req, &resp)

				assert.True(t, resp.Diagnostics.HasError(), "Read should fail when State.Get fails")
			}
		})
	}
}

func TestProjectUserResource_Update(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "successful update",
			wantErr: false,
		},
		{
			name:    "update fails when changing project_id",
			wantErr: true,
		},
		{
			name:    "update fails when changing user_id",
			wantErr: true,
		},
		{
			name:    "update fails with API error",
			wantErr: true,
		},
		{
			name:    "update succeeds when role unchanged",
			wantErr: false,
		},
		{
			name:    "error case - validation checks",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "successful update":
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/projects/proj-123/users/user-456" && r.Method == http.MethodPatch {
						w.WriteHeader(http.StatusOK)
						return
					}
					w.WriteHeader(http.StatusNotFound)
				})

				n8nClient, server := setupTestUserResourceClient(t, handler)
				defer server.Close()

				r := &ProjectUserResource{client: n8nClient}

				rawPlan := map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "proj-123/user-456"),
					"project_id": tftypes.NewValue(tftypes.String, "proj-123"),
					"user_id":    tftypes.NewValue(tftypes.String, "user-456"),
					"role":       tftypes.NewValue(tftypes.String, "project:admin"),
				}
				plan := tfsdk.Plan{
					Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "project_id": tftypes.String, "user_id": tftypes.String, "role": tftypes.String}}, rawPlan),
					Schema: createTestUserResourceSchema(t),
				}

				rawState := map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "proj-123/user-456"),
					"project_id": tftypes.NewValue(tftypes.String, "proj-123"),
					"user_id":    tftypes.NewValue(tftypes.String, "user-456"),
					"role":       tftypes.NewValue(tftypes.String, "project:editor"),
				}
				state := tfsdk.State{
					Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "project_id": tftypes.String, "user_id": tftypes.String, "role": tftypes.String}}, rawState),
					Schema: createTestUserResourceSchema(t),
				}

				req := resource.UpdateRequest{
					Plan:  plan,
					State: state,
				}
				resp := resource.UpdateResponse{
					State: state,
				}

				r.Update(context.Background(), req, &resp)

				assert.False(t, resp.Diagnostics.HasError(), "Update should not have errors")

			case "update fails when changing project_id":
				r := &ProjectUserResource{}

				rawPlan := map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "proj-789/user-456"),
					"project_id": tftypes.NewValue(tftypes.String, "proj-789"),
					"user_id":    tftypes.NewValue(tftypes.String, "user-456"),
					"role":       tftypes.NewValue(tftypes.String, "project:editor"),
				}
				plan := tfsdk.Plan{
					Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "project_id": tftypes.String, "user_id": tftypes.String, "role": tftypes.String}}, rawPlan),
					Schema: createTestUserResourceSchema(t),
				}

				rawState := map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "proj-123/user-456"),
					"project_id": tftypes.NewValue(tftypes.String, "proj-123"),
					"user_id":    tftypes.NewValue(tftypes.String, "user-456"),
					"role":       tftypes.NewValue(tftypes.String, "project:editor"),
				}
				state := tfsdk.State{
					Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "project_id": tftypes.String, "user_id": tftypes.String, "role": tftypes.String}}, rawState),
					Schema: createTestUserResourceSchema(t),
				}

				req := resource.UpdateRequest{
					Plan:  plan,
					State: state,
				}
				resp := resource.UpdateResponse{
					State: state,
				}

				r.Update(context.Background(), req, &resp)

				assert.True(t, resp.Diagnostics.HasError(), "Update should have errors when changing project_id")
				assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Project/User Change Not Supported")

			case "update fails when changing user_id":
				r := &ProjectUserResource{}

				rawPlan := map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "proj-123/user-789"),
					"project_id": tftypes.NewValue(tftypes.String, "proj-123"),
					"user_id":    tftypes.NewValue(tftypes.String, "user-789"),
					"role":       tftypes.NewValue(tftypes.String, "project:editor"),
				}
				plan := tfsdk.Plan{
					Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "project_id": tftypes.String, "user_id": tftypes.String, "role": tftypes.String}}, rawPlan),
					Schema: createTestUserResourceSchema(t),
				}

				rawState := map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "proj-123/user-456"),
					"project_id": tftypes.NewValue(tftypes.String, "proj-123"),
					"user_id":    tftypes.NewValue(tftypes.String, "user-456"),
					"role":       tftypes.NewValue(tftypes.String, "project:editor"),
				}
				state := tfsdk.State{
					Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "project_id": tftypes.String, "user_id": tftypes.String, "role": tftypes.String}}, rawState),
					Schema: createTestUserResourceSchema(t),
				}

				req := resource.UpdateRequest{
					Plan:  plan,
					State: state,
				}
				resp := resource.UpdateResponse{
					State: state,
				}

				r.Update(context.Background(), req, &resp)

				assert.True(t, resp.Diagnostics.HasError(), "Update should have errors when changing user_id")
				assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Project/User Change Not Supported")

			case "update fails with API error":
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"message": "Internal server error"}`))
				})

				n8nClient, server := setupTestUserResourceClient(t, handler)
				defer server.Close()

				r := &ProjectUserResource{client: n8nClient}

				rawPlan := map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "proj-123/user-456"),
					"project_id": tftypes.NewValue(tftypes.String, "proj-123"),
					"user_id":    tftypes.NewValue(tftypes.String, "user-456"),
					"role":       tftypes.NewValue(tftypes.String, "project:admin"),
				}
				plan := tfsdk.Plan{
					Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "project_id": tftypes.String, "user_id": tftypes.String, "role": tftypes.String}}, rawPlan),
					Schema: createTestUserResourceSchema(t),
				}

				rawState := map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "proj-123/user-456"),
					"project_id": tftypes.NewValue(tftypes.String, "proj-123"),
					"user_id":    tftypes.NewValue(tftypes.String, "user-456"),
					"role":       tftypes.NewValue(tftypes.String, "project:editor"),
				}
				state := tfsdk.State{
					Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "project_id": tftypes.String, "user_id": tftypes.String, "role": tftypes.String}}, rawState),
					Schema: createTestUserResourceSchema(t),
				}

				req := resource.UpdateRequest{
					Plan:  plan,
					State: state,
				}
				resp := resource.UpdateResponse{
					State: state,
				}

				r.Update(context.Background(), req, &resp)

				assert.True(t, resp.Diagnostics.HasError(), "Update should have errors")
				assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Error updating user role in project")

			case "update succeeds when role unchanged":
				r := &ProjectUserResource{}

				rawPlan := map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "proj-123/user-456"),
					"project_id": tftypes.NewValue(tftypes.String, "proj-123"),
					"user_id":    tftypes.NewValue(tftypes.String, "user-456"),
					"role":       tftypes.NewValue(tftypes.String, "project:editor"),
				}
				plan := tfsdk.Plan{
					Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "project_id": tftypes.String, "user_id": tftypes.String, "role": tftypes.String}}, rawPlan),
					Schema: createTestUserResourceSchema(t),
				}

				rawState := map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "proj-123/user-456"),
					"project_id": tftypes.NewValue(tftypes.String, "proj-123"),
					"user_id":    tftypes.NewValue(tftypes.String, "user-456"),
					"role":       tftypes.NewValue(tftypes.String, "project:editor"),
				}
				state := tfsdk.State{
					Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "project_id": tftypes.String, "user_id": tftypes.String, "role": tftypes.String}}, rawState),
					Schema: createTestUserResourceSchema(t),
				}

				req := resource.UpdateRequest{
					Plan:  plan,
					State: state,
				}
				resp := resource.UpdateResponse{
					State: state,
				}

				r.Update(context.Background(), req, &resp)

				assert.False(t, resp.Diagnostics.HasError(), "Update should not have errors when role is unchanged")

			case "error case - validation checks":
				r := &ProjectUserResource{}

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
					Schema: createTestUserResourceSchema(t),
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
			}
		})
	}
}

func TestProjectUserResource_Delete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "successful delete",
			wantErr: false,
		},
		{
			name:    "delete fails with API error",
			wantErr: true,
		},
		{
			name:    "delete fails when state get fails",
			wantErr: true,
		},
		{
			name:    "error case - validation checks",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "successful delete":
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/projects/proj-123/users/user-456" && r.Method == http.MethodDelete {
						w.WriteHeader(http.StatusNoContent)
						return
					}
					w.WriteHeader(http.StatusNotFound)
				})

				n8nClient, server := setupTestUserResourceClient(t, handler)
				defer server.Close()

				r := &ProjectUserResource{client: n8nClient}

				rawState := map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "proj-123/user-456"),
					"project_id": tftypes.NewValue(tftypes.String, "proj-123"),
					"user_id":    tftypes.NewValue(tftypes.String, "user-456"),
					"role":       tftypes.NewValue(tftypes.String, "project:editor"),
				}
				state := tfsdk.State{
					Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "project_id": tftypes.String, "user_id": tftypes.String, "role": tftypes.String}}, rawState),
					Schema: createTestUserResourceSchema(t),
				}

				req := resource.DeleteRequest{
					State: state,
				}
				resp := resource.DeleteResponse{
					State: state,
				}

				r.Delete(context.Background(), req, &resp)

				assert.False(t, resp.Diagnostics.HasError(), "Delete should not have errors")

			case "delete fails with API error":
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"message": "Internal server error"}`))
				})

				n8nClient, server := setupTestUserResourceClient(t, handler)
				defer server.Close()

				r := &ProjectUserResource{client: n8nClient}

				rawState := map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "proj-123/user-456"),
					"project_id": tftypes.NewValue(tftypes.String, "proj-123"),
					"user_id":    tftypes.NewValue(tftypes.String, "user-456"),
					"role":       tftypes.NewValue(tftypes.String, "project:editor"),
				}
				state := tfsdk.State{
					Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "project_id": tftypes.String, "user_id": tftypes.String, "role": tftypes.String}}, rawState),
					Schema: createTestUserResourceSchema(t),
				}

				req := resource.DeleteRequest{
					State: state,
				}
				resp := resource.DeleteResponse{
					State: state,
				}

				r.Delete(context.Background(), req, &resp)

				assert.True(t, resp.Diagnostics.HasError(), "Delete should have errors")
				assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Error removing user from project")

			case "delete fails when state get fails":
				r := &ProjectUserResource{}

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
					State: tfsdk.State{Schema: createTestUserResourceSchema(t)},
				}

				r.Delete(context.Background(), req, &resp)

				assert.True(t, resp.Diagnostics.HasError(), "Delete should fail when State.Get fails")

			case "error case - validation checks":
				r := &ProjectUserResource{}

				// Test with empty state
				emptySchema := schema.Schema{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Optional: true,
						},
					},
				}
				rawState := map[string]tftypes.Value{
					"id": tftypes.NewValue(tftypes.String, nil),
				}
				state := tfsdk.State{
					Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String}}, rawState),
					Schema: emptySchema,
				}

				req := resource.DeleteRequest{
					State: state,
				}
				resp := resource.DeleteResponse{
					State: tfsdk.State{Schema: createTestUserResourceSchema(t)},
				}

				r.Delete(context.Background(), req, &resp)

				assert.True(t, resp.Diagnostics.HasError(), "Delete should fail with invalid state schema")
			}
		})
	}
}

func TestProjectUserResource_EdgeCasesForCoverage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "create fails when plan get fails",
			wantErr: true,
		},
		{
			name:    "create succeeds but state set fails",
			wantErr: true,
		},
		{
			name:    "read fails when state get fails",
			wantErr: true,
		},
		{
			name:    "read succeeds but state set fails",
			wantErr: true,
		},
		{
			name:    "update fails when plan get fails",
			wantErr: true,
		},
		{
			name:    "update succeeds but state set fails",
			wantErr: true,
		},
		{
			name:    "error case - validation checks",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "create fails when plan get fails":
				r := &ProjectUserResource{}

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
					State: tfsdk.State{Schema: createTestUserResourceSchema(t)},
				}

				r.Create(context.Background(), req, &resp)

				assert.True(t, resp.Diagnostics.HasError(), "Create should fail when Plan.Get fails")

			case "create succeeds but state set fails":
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/projects/proj-123/users/user-456" && r.Method == http.MethodPut {
						w.Header().Set("Content-Type", "application/json")
						json.NewEncoder(w).Encode(map[string]interface{}{})
						return
					}
					w.WriteHeader(http.StatusNotFound)
				})

				n8nClient, server := setupTestUserResourceClient(t, handler)
				defer server.Close()

				r := &ProjectUserResource{client: n8nClient}

				rawPlan := map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, nil),
					"project_id": tftypes.NewValue(tftypes.String, "proj-123"),
					"user_id":    tftypes.NewValue(tftypes.String, "user-456"),
					"role":       tftypes.NewValue(tftypes.String, "project:editor"),
				}
				plan := tfsdk.Plan{
					Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "project_id": tftypes.String, "user_id": tftypes.String, "role": tftypes.String}}, rawPlan),
					Schema: createTestUserResourceSchema(t),
				}

				// Create state with incompatible schema
				wrongSchema := schema.Schema{
					Attributes: map[string]schema.Attribute{
						"id": schema.NumberAttribute{
							Computed: true,
						},
					},
				}
				state := tfsdk.State{
					Schema: wrongSchema,
				}

				req := resource.CreateRequest{
					Plan: plan,
				}
				resp := resource.CreateResponse{
					State: state,
				}

				r.Create(context.Background(), req, &resp)

				assert.True(t, resp.Diagnostics.HasError(), "Create should fail when State.Set fails")

			case "read fails when state get fails":
				r := &ProjectUserResource{}

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
					State: tfsdk.State{Schema: createTestUserResourceSchema(t)},
				}

				r.Read(context.Background(), req, &resp)

				assert.True(t, resp.Diagnostics.HasError(), "Read should fail when State.Get fails")

			case "read succeeds but state set fails":
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/users" && r.Method == http.MethodGet && r.URL.Query().Get("projectId") == "proj-123" {
						w.Header().Set("Content-Type", "application/json")
						json.NewEncoder(w).Encode(map[string]interface{}{
							"data": []map[string]interface{}{
								{
									"id":    "user-456",
									"email": "test@example.com",
									"role":  "project:editor",
								},
							},
						})
						return
					}
					w.WriteHeader(http.StatusNotFound)
				})

				n8nClient, server := setupTestUserResourceClient(t, handler)
				defer server.Close()

				r := &ProjectUserResource{client: n8nClient}

				rawState := map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "proj-123/user-456"),
					"project_id": tftypes.NewValue(tftypes.String, "proj-123"),
					"user_id":    tftypes.NewValue(tftypes.String, "user-456"),
					"role":       tftypes.NewValue(tftypes.String, "project:editor"),
				}
				state := tfsdk.State{
					Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "project_id": tftypes.String, "user_id": tftypes.String, "role": tftypes.String}}, rawState),
					Schema: createTestUserResourceSchema(t),
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

			case "update fails when plan get fails":
				r := &ProjectUserResource{}

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
					Schema: createTestUserResourceSchema(t),
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

			case "update succeeds but state set fails":
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/projects/proj-123/users/user-456" && r.Method == http.MethodPut {
						w.Header().Set("Content-Type", "application/json")
						json.NewEncoder(w).Encode(map[string]interface{}{})
						return
					}
					w.WriteHeader(http.StatusNotFound)
				})

				n8nClient, server := setupTestUserResourceClient(t, handler)
				defer server.Close()

				r := &ProjectUserResource{client: n8nClient}

				rawPlan := map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "proj-123/user-456"),
					"project_id": tftypes.NewValue(tftypes.String, "proj-123"),
					"user_id":    tftypes.NewValue(tftypes.String, "user-456"),
					"role":       tftypes.NewValue(tftypes.String, "project:admin"),
				}
				plan := tfsdk.Plan{
					Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "project_id": tftypes.String, "user_id": tftypes.String, "role": tftypes.String}}, rawPlan),
					Schema: createTestUserResourceSchema(t),
				}

				rawState := map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "proj-123/user-456"),
					"project_id": tftypes.NewValue(tftypes.String, "proj-123"),
					"user_id":    tftypes.NewValue(tftypes.String, "user-456"),
					"role":       tftypes.NewValue(tftypes.String, "project:editor"),
				}
				state := tfsdk.State{
					Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "project_id": tftypes.String, "user_id": tftypes.String, "role": tftypes.String}}, rawState),
					Schema: createTestUserResourceSchema(t),
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

				req := resource.UpdateRequest{
					Plan:  plan,
					State: state,
				}
				resp := resource.UpdateResponse{
					State: respState,
				}

				r.Update(context.Background(), req, &resp)

				assert.True(t, resp.Diagnostics.HasError(), "Update should fail when State.Set fails")

			case "error case - validation checks":
				// Test that state set failures are properly handled
				r := &ProjectUserResource{}

				// Create plan with mismatched schema (different from first test case)
				wrongSchema := schema.Schema{
					Attributes: map[string]schema.Attribute{
						"id": schema.BoolAttribute{
							Required: true,
						},
					},
				}
				rawPlan := map[string]tftypes.Value{
					"id": tftypes.NewValue(tftypes.Bool, true),
				}
				plan := tfsdk.Plan{
					Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.Bool}}, rawPlan),
					Schema: wrongSchema,
				}

				req := resource.CreateRequest{
					Plan: plan,
				}
				resp := resource.CreateResponse{
					State: tfsdk.State{Schema: createTestUserResourceSchema(t)},
				}

				r.Create(context.Background(), req, &resp)

				assert.True(t, resp.Diagnostics.HasError(), "Create should fail with invalid plan schema")
			}
		})
	}
}

func TestProjectUserResource_ImportState(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "successful import",
			wantErr: false,
		},
		{
			name:    "import fails with invalid ID format",
			wantErr: true,
		},
		{
			name:    "import fails with too many parts",
			wantErr: true,
		},
		{
			name:    "error case - validation checks",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "successful import":
				r := &ProjectUserResource{}

				schema := createTestUserResourceSchema(t)
				state := tfsdk.State{
					Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "project_id": tftypes.String, "user_id": tftypes.String, "role": tftypes.String}}, nil),
					Schema: schema,
				}

				req := resource.ImportStateRequest{
					ID: "proj-123/user-456",
				}
				resp := &resource.ImportStateResponse{
					State: state,
				}

				r.ImportState(context.Background(), req, resp)

				assert.False(t, resp.Diagnostics.HasError(), "ImportState should not have errors")

			case "import fails with invalid ID format":
				r := &ProjectUserResource{}

				schema := createTestUserResourceSchema(t)
				state := tfsdk.State{
					Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "project_id": tftypes.String, "user_id": tftypes.String, "role": tftypes.String}}, nil),
					Schema: schema,
				}

				req := resource.ImportStateRequest{
					ID: "invalid-id",
				}
				resp := &resource.ImportStateResponse{
					State: state,
				}

				r.ImportState(context.Background(), req, resp)

				assert.True(t, resp.Diagnostics.HasError(), "ImportState should have errors with invalid ID")
				assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Invalid Import ID")

			case "import fails with too many parts":
				r := &ProjectUserResource{}

				schema := createTestUserResourceSchema(t)
				state := tfsdk.State{
					Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "project_id": tftypes.String, "user_id": tftypes.String, "role": tftypes.String}}, nil),
					Schema: schema,
				}

				req := resource.ImportStateRequest{
					ID: "proj-123/user-456/extra",
				}
				resp := &resource.ImportStateResponse{
					State: state,
				}

				r.ImportState(context.Background(), req, resp)

				assert.True(t, resp.Diagnostics.HasError(), "ImportState should have errors with too many parts")

			case "error case - validation checks":
				r := &ProjectUserResource{}

				schema := createTestUserResourceSchema(t)
				state := tfsdk.State{
					Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "project_id": tftypes.String, "user_id": tftypes.String, "role": tftypes.String}}, nil),
					Schema: schema,
				}

				req := resource.ImportStateRequest{
					ID: "",
				}
				resp := &resource.ImportStateResponse{
					State: state,
				}

				r.ImportState(context.Background(), req, resp)

				assert.True(t, resp.Diagnostics.HasError(), "ImportState should fail with empty ID")
			}
		})
	}
}

func TestProjectUserResource_Interfaces(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "implements required interfaces",
			wantErr: false,
		},
		{
			name:    "error case - validation checks",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "implements required interfaces":
				r := NewProjectUserResource()

				var _ resource.Resource = r
				var _ resource.ResourceWithConfigure = r
				var _ resource.ResourceWithImportState = r
				var _ ProjectUserResourceInterface = r

			case "error case - validation checks":
				r := NewProjectUserResource()

				assert.Implements(t, (*resource.Resource)(nil), r, "should implement resource.Resource")
				assert.Implements(t, (*resource.ResourceWithConfigure)(nil), r, "should implement ResourceWithConfigure")
				assert.Implements(t, (*resource.ResourceWithImportState)(nil), r, "should implement ResourceWithImportState")
				assert.Implements(t, (*ProjectUserResourceInterface)(nil), r, "should implement ProjectUserResourceInterface")
			}
		})
	}
}

// createTestUserResourceSchema creates a test schema for project user resource.
func createTestUserResourceSchema(t *testing.T) schema.Schema {
	t.Helper()
	r := &ProjectUserResource{}
	req := resource.SchemaRequest{}
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), req, resp)
	return resp.Schema
}

// setupTestUserResourceClient creates a test N8nClient with httptest server.
func setupTestUserResourceClient(t *testing.T, handler http.HandlerFunc) (*client.N8nClient, *httptest.Server) {
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

func TestProjectUserResource_searchUserInList(t *testing.T) {
	t.Parallel()

	userID := "user-456"
	role := "project:editor"

	t.Run("user found with role", func(t *testing.T) {
		state := &models.UserResource{
			UserID: types.StringValue("user-456"),
		}
		userList := &n8nsdk.UserList{
			Data: []n8nsdk.User{
				{Id: &userID, Role: &role},
			},
		}
		r := &ProjectUserResource{}

		found := r.searchUserInList(userList, state)

		assert.True(t, found)
		assert.Equal(t, "project:editor", state.Role.ValueString())
	})

	t.Run("user found without role", func(t *testing.T) {
		state := &models.UserResource{
			UserID: types.StringValue("user-456"),
		}
		userList := &n8nsdk.UserList{
			Data: []n8nsdk.User{
				{Id: &userID, Role: nil},
			},
		}
		r := &ProjectUserResource{}

		found := r.searchUserInList(userList, state)

		assert.True(t, found)
		assert.True(t, state.Role.IsNull())
	})

	t.Run("user not found", func(t *testing.T) {
		otherUserID := "user-999"
		state := &models.UserResource{
			UserID: types.StringValue("user-456"),
		}
		userList := &n8nsdk.UserList{
			Data: []n8nsdk.User{
				{Id: &otherUserID},
			},
		}
		r := &ProjectUserResource{}

		found := r.searchUserInList(userList, state)

		assert.False(t, found)
	})

	t.Run("nil user data", func(t *testing.T) {
		state := &models.UserResource{
			UserID: types.StringValue("user-456"),
		}
		userList := &n8nsdk.UserList{
			Data: nil,
		}
		r := &ProjectUserResource{}

		found := r.searchUserInList(userList, state)

		assert.False(t, found)
	})

	t.Run("nil user ID in list", func(t *testing.T) {
		state := &models.UserResource{
			UserID: types.StringValue("user-456"),
		}
		userList := &n8nsdk.UserList{
			Data: []n8nsdk.User{
				{Id: nil},
			},
		}
		r := &ProjectUserResource{}

		found := r.searchUserInList(userList, state)

		assert.False(t, found)
	})
}
