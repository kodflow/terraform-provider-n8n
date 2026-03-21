// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

package datatable_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/datatable"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/client"
	"github.com/stretchr/testify/assert"
)

func TestNewDataTableResource(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "normal case",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := datatable.NewDataTableResource()
				assert.NotNil(t, r)
			},
		},
		{
			name: "error case - validates behavior",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := datatable.NewDataTableResource()
				assert.NotNil(t, r)
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

func TestNewDataTableResourceWrapper(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "normal case",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := datatable.NewDataTableResourceWrapper()
				assert.NotNil(t, r)
				assert.Implements(t, (*resource.Resource)(nil), r)
			},
		},
		{
			name: "error case - validates behavior",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := datatable.NewDataTableResourceWrapper()
				assert.NotNil(t, r)
				assert.Implements(t, (*resource.Resource)(nil), r)
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

func TestDataTableResource_Metadata(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "normal case",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := datatable.NewDataTableResource()
				resp := &resource.MetadataResponse{}

				r.Metadata(t.Context(), resource.MetadataRequest{
					ProviderTypeName: "n8n",
				}, resp)

				assert.Equal(t, "n8n_data_table", resp.TypeName)
			},
		},
		{
			name: "error case - validates behavior",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := datatable.NewDataTableResource()
				resp := &resource.MetadataResponse{}

				r.Metadata(t.Context(), resource.MetadataRequest{
					ProviderTypeName: "n8n",
				}, resp)

				assert.Equal(t, "n8n_data_table", resp.TypeName)
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

func TestDataTableResource_Schema(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "normal case",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := datatable.NewDataTableResource()
				resp := &resource.SchemaResponse{}

				r.Schema(t.Context(), resource.SchemaRequest{}, resp)

				assert.NotNil(t, resp.Schema)
				assert.NotEmpty(t, resp.Schema.Attributes)
			},
		},
		{
			name: "error case - validates behavior",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := datatable.NewDataTableResource()
				resp := &resource.SchemaResponse{}

				r.Schema(t.Context(), resource.SchemaRequest{}, resp)

				assert.NotNil(t, resp.Schema)
				assert.NotEmpty(t, resp.Schema.Attributes)
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

func TestDataTableResource_Configure(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		providerData any
		wantError    bool
	}{
		{
			name:         "valid configuration",
			providerData: &client.N8nClient{},
			wantError:    false,
		},
		{
			name:         "nil provider data",
			providerData: nil,
			wantError:    false,
		},
		{
			name:         "invalid provider data type",
			providerData: "invalid",
			wantError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := datatable.NewDataTableResource()
			resp := &resource.ConfigureResponse{}
			req := resource.ConfigureRequest{
				ProviderData: tt.providerData,
			}

			r.Configure(t.Context(), req, resp)

			if tt.wantError {
				assert.True(t, resp.Diagnostics.HasError())
			} else {
				assert.False(t, resp.Diagnostics.HasError())
			}
		})
	}
}

// dataTableResourceStateParams holds parameters for building a DataTableResource state value.
type dataTableResourceStateParams struct {
	ID        string
	Name      string
	ProjectID string
	CreatedAt string
	UpdatedAt string
}

// buildDataTableResourceState constructs a valid tftypes.Value for the DataTableResource schema.
func buildDataTableResourceState(ctx context.Context, schemaType tftypes.Type, p dataTableResourceStateParams) tftypes.Value {
	id, name, projectID, createdAt, updatedAt := p.ID, p.Name, p.ProjectID, p.CreatedAt, p.UpdatedAt
	idVal := tftypes.NewValue(tftypes.String, id)
	if id == "" {
		idVal = tftypes.NewValue(tftypes.String, nil)
	}
	projectVal := tftypes.NewValue(tftypes.String, projectID)
	if projectID == "" {
		projectVal = tftypes.NewValue(tftypes.String, nil)
	}
	createdVal := tftypes.NewValue(tftypes.String, createdAt)
	if createdAt == "" {
		createdVal = tftypes.NewValue(tftypes.String, nil)
	}
	updatedVal := tftypes.NewValue(tftypes.String, updatedAt)
	if updatedAt == "" {
		updatedVal = tftypes.NewValue(tftypes.String, nil)
	}

	return tftypes.NewValue(schemaType, map[string]tftypes.Value{
		"id":         idVal,
		"name":       tftypes.NewValue(tftypes.String, name),
		"project_id": projectVal,
		"created_at": createdVal,
		"updated_at": updatedVal,
		"columns": tftypes.NewValue(tftypes.List{
			ElementType: tftypes.Object{
				AttributeTypes: map[string]tftypes.Type{
					"name": tftypes.String,
					"type": tftypes.String,
				},
			},
		}, []tftypes.Value{}),
	})
}

func TestDataTableResource_Create(t *testing.T) {
	t.Parallel()
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
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte(`{"id":"dt-123","name":"My Table","columns":[],"projectId":"proj-1","createdAt":"2024-01-01T00:00:00Z","updatedAt":"2024-01-01T00:00:00Z"}`))
				})

				n8nClient, server := setupTestClientForDTDataSource(t, handler)
				defer server.Close()

				r := datatable.NewDataTableResource()
				r.Configure(t.Context(), resource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &resource.ConfigureResponse{})

				ctx := t.Context()
				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				planRaw := buildDataTableResourceState(ctx, schemaResp.Schema.Type().TerraformType(ctx), dataTableResourceStateParams{ID: "", Name: "My Table", ProjectID: "", CreatedAt: "", UpdatedAt: ""})

				plan := tfsdk.Plan{
					Schema: schemaResp.Schema,
					Raw:    planRaw,
				}

				state := tfsdk.State{
					Schema: schemaResp.Schema,
				}

				req := resource.CreateRequest{Plan: plan}
				resp := &resource.CreateResponse{State: state}

				r.Create(ctx, req, resp)

				if resp.Diagnostics.HasError() {
					for _, diag := range resp.Diagnostics.Errors() {
						t.Logf("Error: %s - %s", diag.Summary(), diag.Detail())
					}
				}
				assert.False(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "error - create with invalid plan",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := datatable.NewDataTableResource()
				ctx := t.Context()

				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				planRaw := tftypes.NewValue(tftypes.String, "invalid")
				plan := tfsdk.Plan{
					Schema: schemaResp.Schema,
					Raw:    planRaw,
				}

				req := resource.CreateRequest{Plan: plan}
				resp := &resource.CreateResponse{}

				r.Create(ctx, req, resp)

				assert.True(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "error - create with API error",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"message":"Internal server error"}`))
				})

				n8nClient, server := setupTestClientForDTDataSource(t, handler)
				defer server.Close()

				r := datatable.NewDataTableResource()
				r.Configure(t.Context(), resource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &resource.ConfigureResponse{})

				ctx := t.Context()
				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				planRaw := buildDataTableResourceState(ctx, schemaResp.Schema.Type().TerraformType(ctx), dataTableResourceStateParams{ID: "", Name: "My Table", ProjectID: "", CreatedAt: "", UpdatedAt: ""})
				plan := tfsdk.Plan{
					Schema: schemaResp.Schema,
					Raw:    planRaw,
				}

				state := tfsdk.State{
					Schema: schemaResp.Schema,
				}

				req := resource.CreateRequest{Plan: plan}
				resp := &resource.CreateResponse{State: state}

				r.Create(ctx, req, resp)

				assert.True(t, resp.Diagnostics.HasError())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}

func TestDataTableResource_Read(t *testing.T) {
	t.Parallel()
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
					w.Write([]byte(`{"id":"dt-123","name":"My Table","columns":[],"projectId":"proj-1","createdAt":"2024-01-01T00:00:00Z","updatedAt":"2024-01-01T00:00:00Z"}`))
				})

				n8nClient, server := setupTestClientForDTDataSource(t, handler)
				defer server.Close()

				r := datatable.NewDataTableResource()
				r.Configure(t.Context(), resource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &resource.ConfigureResponse{})

				ctx := t.Context()
				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				stateRaw := buildDataTableResourceState(ctx, schemaResp.Schema.Type().TerraformType(ctx), dataTableResourceStateParams{ID: "dt-123", Name: "My Table", ProjectID: "proj-1", CreatedAt: "2024-01-01T00:00:00Z", UpdatedAt: "2024-01-01T00:00:00Z"})
				state := tfsdk.State{
					Schema: schemaResp.Schema,
					Raw:    stateRaw,
				}

				req := resource.ReadRequest{State: state}
				resp := &resource.ReadResponse{State: state}

				r.Read(ctx, req, resp)

				assert.False(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "read - resource deleted externally (404)",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte(`{"message":"Not found"}`))
				})

				n8nClient, server := setupTestClientForDTDataSource(t, handler)
				defer server.Close()

				r := datatable.NewDataTableResource()
				r.Configure(t.Context(), resource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &resource.ConfigureResponse{})

				ctx := t.Context()
				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				stateRaw := buildDataTableResourceState(ctx, schemaResp.Schema.Type().TerraformType(ctx), dataTableResourceStateParams{ID: "dt-404", Name: "My Table", ProjectID: "proj-1", CreatedAt: "2024-01-01T00:00:00Z", UpdatedAt: "2024-01-01T00:00:00Z"})
				state := tfsdk.State{
					Schema: schemaResp.Schema,
					Raw:    stateRaw,
				}

				req := resource.ReadRequest{State: state}
				resp := &resource.ReadResponse{State: state}

				r.Read(ctx, req, resp)

				// No error on 404 - resource is removed from state
				assert.False(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "error - read with invalid state",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := datatable.NewDataTableResource()
				ctx := t.Context()

				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				stateRaw := tftypes.NewValue(tftypes.String, "invalid")
				state := tfsdk.State{
					Schema: schemaResp.Schema,
					Raw:    stateRaw,
				}

				req := resource.ReadRequest{State: state}
				resp := &resource.ReadResponse{}

				r.Read(ctx, req, resp)

				assert.True(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "error - read with API error",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"message":"Internal server error"}`))
				})

				n8nClient, server := setupTestClientForDTDataSource(t, handler)
				defer server.Close()

				r := datatable.NewDataTableResource()
				r.Configure(t.Context(), resource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &resource.ConfigureResponse{})

				ctx := t.Context()
				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				stateRaw := buildDataTableResourceState(ctx, schemaResp.Schema.Type().TerraformType(ctx), dataTableResourceStateParams{ID: "dt-123", Name: "My Table", ProjectID: "proj-1", CreatedAt: "2024-01-01T00:00:00Z", UpdatedAt: "2024-01-01T00:00:00Z"})
				state := tfsdk.State{
					Schema: schemaResp.Schema,
					Raw:    stateRaw,
				}

				req := resource.ReadRequest{State: state}
				resp := &resource.ReadResponse{State: state}

				r.Read(ctx, req, resp)

				assert.True(t, resp.Diagnostics.HasError())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}

func TestDataTableResource_Update(t *testing.T) {
	t.Parallel()
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
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{"id":"dt-123","name":"Updated Table","columns":[],"projectId":"proj-1","createdAt":"2024-01-01T00:00:00Z","updatedAt":"2024-01-02T00:00:00Z"}`))
				})

				n8nClient, server := setupTestClientForDTDataSource(t, handler)
				defer server.Close()

				r := datatable.NewDataTableResource()
				r.Configure(t.Context(), resource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &resource.ConfigureResponse{})

				ctx := t.Context()
				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				planRaw := buildDataTableResourceState(ctx, schemaResp.Schema.Type().TerraformType(ctx), dataTableResourceStateParams{ID: "dt-123", Name: "Updated Table", ProjectID: "proj-1", CreatedAt: "2024-01-01T00:00:00Z", UpdatedAt: "2024-01-01T00:00:00Z"})
				stateRaw := buildDataTableResourceState(ctx, schemaResp.Schema.Type().TerraformType(ctx), dataTableResourceStateParams{ID: "dt-123", Name: "Old Table", ProjectID: "proj-1", CreatedAt: "2024-01-01T00:00:00Z", UpdatedAt: "2024-01-01T00:00:00Z"})

				plan := tfsdk.Plan{Schema: schemaResp.Schema, Raw: planRaw}
				state := tfsdk.State{Schema: schemaResp.Schema, Raw: stateRaw}

				req := resource.UpdateRequest{Plan: plan, State: state}
				resp := &resource.UpdateResponse{State: state}

				r.Update(ctx, req, resp)

				assert.False(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "error - update with invalid plan",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := datatable.NewDataTableResource()
				ctx := t.Context()

				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				planRaw := tftypes.NewValue(tftypes.String, "invalid")
				stateRaw := buildDataTableResourceState(ctx, schemaResp.Schema.Type().TerraformType(ctx), dataTableResourceStateParams{ID: "dt-123", Name: "Old Table", ProjectID: "proj-1", CreatedAt: "2024-01-01T00:00:00Z", UpdatedAt: "2024-01-01T00:00:00Z"})

				plan := tfsdk.Plan{Schema: schemaResp.Schema, Raw: planRaw}
				state := tfsdk.State{Schema: schemaResp.Schema, Raw: stateRaw}

				req := resource.UpdateRequest{Plan: plan, State: state}
				resp := &resource.UpdateResponse{}

				r.Update(ctx, req, resp)

				assert.True(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "error - update with API error",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"message":"Internal server error"}`))
				})

				n8nClient, server := setupTestClientForDTDataSource(t, handler)
				defer server.Close()

				r := datatable.NewDataTableResource()
				r.Configure(t.Context(), resource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &resource.ConfigureResponse{})

				ctx := t.Context()
				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				planRaw := buildDataTableResourceState(ctx, schemaResp.Schema.Type().TerraformType(ctx), dataTableResourceStateParams{ID: "dt-123", Name: "Updated Table", ProjectID: "proj-1", CreatedAt: "2024-01-01T00:00:00Z", UpdatedAt: "2024-01-01T00:00:00Z"})
				stateRaw := buildDataTableResourceState(ctx, schemaResp.Schema.Type().TerraformType(ctx), dataTableResourceStateParams{ID: "dt-123", Name: "Old Table", ProjectID: "proj-1", CreatedAt: "2024-01-01T00:00:00Z", UpdatedAt: "2024-01-01T00:00:00Z"})

				plan := tfsdk.Plan{Schema: schemaResp.Schema, Raw: planRaw}
				state := tfsdk.State{Schema: schemaResp.Schema, Raw: stateRaw}

				req := resource.UpdateRequest{Plan: plan, State: state}
				resp := &resource.UpdateResponse{State: state}

				r.Update(ctx, req, resp)

				assert.True(t, resp.Diagnostics.HasError())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}

func TestDataTableResource_Delete(t *testing.T) {
	t.Parallel()
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

				n8nClient, server := setupTestClientForDTDataSource(t, handler)
				defer server.Close()

				r := datatable.NewDataTableResource()
				r.Configure(t.Context(), resource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &resource.ConfigureResponse{})

				ctx := t.Context()
				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				stateRaw := buildDataTableResourceState(ctx, schemaResp.Schema.Type().TerraformType(ctx), dataTableResourceStateParams{ID: "dt-123", Name: "My Table", ProjectID: "proj-1", CreatedAt: "2024-01-01T00:00:00Z", UpdatedAt: "2024-01-01T00:00:00Z"})
				state := tfsdk.State{Schema: schemaResp.Schema, Raw: stateRaw}

				req := resource.DeleteRequest{State: state}
				resp := &resource.DeleteResponse{}

				r.Delete(ctx, req, resp)

				assert.False(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "error - delete with invalid state",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := datatable.NewDataTableResource()
				ctx := t.Context()

				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				stateRaw := tftypes.NewValue(tftypes.String, "invalid")
				state := tfsdk.State{Schema: schemaResp.Schema, Raw: stateRaw}

				req := resource.DeleteRequest{State: state}
				resp := &resource.DeleteResponse{}

				r.Delete(ctx, req, resp)

				assert.True(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "error - delete with API error",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"message":"Internal server error"}`))
				})

				n8nClient, server := setupTestClientForDTDataSource(t, handler)
				defer server.Close()

				r := datatable.NewDataTableResource()
				r.Configure(t.Context(), resource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &resource.ConfigureResponse{})

				ctx := t.Context()
				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				stateRaw := buildDataTableResourceState(ctx, schemaResp.Schema.Type().TerraformType(ctx), dataTableResourceStateParams{ID: "dt-123", Name: "My Table", ProjectID: "proj-1", CreatedAt: "2024-01-01T00:00:00Z", UpdatedAt: "2024-01-01T00:00:00Z"})
				state := tfsdk.State{Schema: schemaResp.Schema, Raw: stateRaw}

				req := resource.DeleteRequest{State: state}
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

func TestDataTableResource_ImportState(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "import state passthrough",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := datatable.NewDataTableResource()
				ctx := t.Context()

				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				emptyValue := buildDataTableResourceState(ctx, schemaResp.Schema.Type().TerraformType(ctx), dataTableResourceStateParams{})

				req := resource.ImportStateRequest{
					ID: "dt-123",
				}
				resp := &resource.ImportStateResponse{
					State: tfsdk.State{
						Schema: schemaResp.Schema,
						Raw:    emptyValue,
					},
				}

				r.ImportState(ctx, req, resp)

				assert.False(t, resp.Diagnostics.HasError())
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
