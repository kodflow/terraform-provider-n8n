// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

package workflow_test

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/client"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/workflow"
	"github.com/stretchr/testify/assert"
)

// TestNewWorkflowVersionDataSource tests the NewWorkflowVersionDataSource constructor.
func TestNewWorkflowVersionDataSource(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "creates valid datasource",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := workflow.NewWorkflowVersionDataSource()
				assert.NotNil(t, ds)
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

// TestNewWorkflowVersionDataSourceWrapper tests the NewWorkflowVersionDataSourceWrapper constructor.
func TestNewWorkflowVersionDataSourceWrapper(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "creates valid datasource wrapper",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := workflow.NewWorkflowVersionDataSourceWrapper()
				assert.NotNil(t, ds)
				assert.Implements(t, (*datasource.DataSource)(nil), ds)
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

// TestWorkflowVersionDataSource_Metadata tests the Metadata method.
func TestWorkflowVersionDataSource_Metadata(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "sets correct type name",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := workflow.NewWorkflowVersionDataSource()
				resp := &datasource.MetadataResponse{}

				ds.Metadata(t.Context(), datasource.MetadataRequest{
					ProviderTypeName: "n8n",
				}, resp)

				assert.Equal(t, "n8n_workflow_version", resp.TypeName)
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

// TestWorkflowVersionDataSource_Schema tests the Schema method.
func TestWorkflowVersionDataSource_Schema(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "returns valid schema with required attributes",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := workflow.NewWorkflowVersionDataSource()
				resp := &datasource.SchemaResponse{}

				ds.Schema(t.Context(), datasource.SchemaRequest{}, resp)

				assert.NotNil(t, resp.Schema)
				assert.Contains(t, resp.Schema.Attributes, "workflow_id")
				assert.Contains(t, resp.Schema.Attributes, "version_id")
				assert.Contains(t, resp.Schema.Attributes, "name")
				assert.Contains(t, resp.Schema.Attributes, "nodes_json")
				assert.Contains(t, resp.Schema.Attributes, "connections_json")
				assert.Contains(t, resp.Schema.Attributes, "created_at")
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

// TestWorkflowVersionDataSource_Configure tests the Configure method.
func TestWorkflowVersionDataSource_Configure(t *testing.T) {
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

			ds := workflow.NewWorkflowVersionDataSource()
			resp := &datasource.ConfigureResponse{}
			req := datasource.ConfigureRequest{
				ProviderData: tt.providerData,
			}

			ds.Configure(t.Context(), req, resp)

			if tt.wantError {
				assert.True(t, resp.Diagnostics.HasError())
			} else {
				assert.False(t, resp.Diagnostics.HasError())
			}
		})
	}
}

// TestWorkflowVersionDataSource_Read tests the Read method.
func TestWorkflowVersionDataSource_Read(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "successful read",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{"versionId":"v1","workflowId":"wf-1","name":"My Workflow","authors":"Alice","nodes":[],"connections":{},"createdAt":"2024-01-01T00:00:00Z"}`))
				})

				n8nClient, server := setupTestClientForDataSource(t, handler)
				defer server.Close()

				ds := workflow.NewWorkflowVersionDataSource()
				ds.Configure(t.Context(), datasource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &datasource.ConfigureResponse{})

				ctx := t.Context()
				schemaResp := datasource.SchemaResponse{}
				ds.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)

				configRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"workflow_id":      tftypes.NewValue(tftypes.String, "wf-1"),
					"version_id":       tftypes.NewValue(tftypes.String, "v1"),
					"name":             tftypes.NewValue(tftypes.String, nil),
					"authors":          tftypes.NewValue(tftypes.String, nil),
					"nodes_json":       tftypes.NewValue(tftypes.String, nil),
					"connections_json": tftypes.NewValue(tftypes.String, nil),
					"created_at":       tftypes.NewValue(tftypes.String, nil),
				})

				req := datasource.ReadRequest{
					Config: tfsdk.Config{
						Schema: schemaResp.Schema,
						Raw:    configRaw,
					},
				}
				resp := &datasource.ReadResponse{
					State: tfsdk.State{
						Schema: schemaResp.Schema,
					},
				}

				ds.Read(ctx, req, resp)

				assert.False(t, resp.Diagnostics.HasError(), "Expected no errors")
			},
		},
		{
			name: "error - invalid config",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := workflow.NewWorkflowVersionDataSource()
				ctx := t.Context()
				schemaResp := datasource.SchemaResponse{}
				ds.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)

				configRaw := tftypes.NewValue(tftypes.String, "invalid")

				req := datasource.ReadRequest{
					Config: tfsdk.Config{
						Schema: schemaResp.Schema,
						Raw:    configRaw,
					},
				}
				resp := &datasource.ReadResponse{
					State: tfsdk.State{
						Schema: schemaResp.Schema,
					},
				}

				ds.Read(ctx, req, resp)

				assert.True(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "error - api error",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte(`{"message":"not found"}`))
				})

				n8nClient, server := setupTestClientForDataSource(t, handler)
				defer server.Close()

				ds := workflow.NewWorkflowVersionDataSource()
				ds.Configure(t.Context(), datasource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &datasource.ConfigureResponse{})

				ctx := t.Context()
				schemaResp := datasource.SchemaResponse{}
				ds.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)

				configRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"workflow_id":      tftypes.NewValue(tftypes.String, "wf-1"),
					"version_id":       tftypes.NewValue(tftypes.String, "v1"),
					"name":             tftypes.NewValue(tftypes.String, nil),
					"authors":          tftypes.NewValue(tftypes.String, nil),
					"nodes_json":       tftypes.NewValue(tftypes.String, nil),
					"connections_json": tftypes.NewValue(tftypes.String, nil),
					"created_at":       tftypes.NewValue(tftypes.String, nil),
				})

				req := datasource.ReadRequest{
					Config: tfsdk.Config{
						Schema: schemaResp.Schema,
						Raw:    configRaw,
					},
				}
				resp := &datasource.ReadResponse{
					State: tfsdk.State{
						Schema: schemaResp.Schema,
					},
				}

				ds.Read(ctx, req, resp)

				assert.True(t, resp.Diagnostics.HasError())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t)
		})
	}
}
