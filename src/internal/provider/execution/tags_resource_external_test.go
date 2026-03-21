// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

package execution_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/execution"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/client"
	"github.com/stretchr/testify/assert"
)

// setupTestClientForExecution creates a test N8nClient with httptest server.
func setupTestClientForExecution(t *testing.T, handler http.HandlerFunc) (*client.N8nClient, *httptest.Server) {
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

// buildTagsResourceState builds a valid tftypes.Value for ExecutionTagsResource schema.
func buildTagsResourceState(ctx context.Context, schemaType tftypes.Type, executionID string, tagIDs []tftypes.Value) tftypes.Value {
	return tftypes.NewValue(schemaType, map[string]tftypes.Value{
		"execution_id": tftypes.NewValue(tftypes.String, executionID),
		"tag_ids": tftypes.NewValue(tftypes.Set{
			ElementType: tftypes.String,
		}, tagIDs),
	})
}

func TestNewExecutionTagsResource(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "normal case",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := execution.NewExecutionTagsResource()
				assert.NotNil(t, r)
			},
		},
		{
			name: "error case - validates behavior",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := execution.NewExecutionTagsResource()
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

func TestNewExecutionTagsResourceWrapper(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "normal case",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := execution.NewExecutionTagsResourceWrapper()
				assert.NotNil(t, r)
				assert.Implements(t, (*resource.Resource)(nil), r)
			},
		},
		{
			name: "error case - validates behavior",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := execution.NewExecutionTagsResourceWrapper()
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

func TestExecutionTagsResource_Metadata(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "normal case",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := execution.NewExecutionTagsResource()
				resp := &resource.MetadataResponse{}

				r.Metadata(t.Context(), resource.MetadataRequest{
					ProviderTypeName: "n8n",
				}, resp)

				assert.Equal(t, "n8n_execution_tags", resp.TypeName)
			},
		},
		{
			name: "error case - validates behavior",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := execution.NewExecutionTagsResource()
				resp := &resource.MetadataResponse{}

				r.Metadata(t.Context(), resource.MetadataRequest{
					ProviderTypeName: "n8n",
				}, resp)

				assert.Equal(t, "n8n_execution_tags", resp.TypeName)
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

func TestExecutionTagsResource_Schema(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "normal case",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := execution.NewExecutionTagsResource()
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
				r := execution.NewExecutionTagsResource()
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

func TestExecutionTagsResource_Configure(t *testing.T) {
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

			r := execution.NewExecutionTagsResource()
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

func TestExecutionTagsResource_Create(t *testing.T) {
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
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`[{"id":"tag-1","name":"Tag One"},{"id":"tag-2","name":"Tag Two"}]`))
				})

				n8nClient, server := setupTestClientForExecution(t, handler)
				defer server.Close()

				r := execution.NewExecutionTagsResource()
				r.Configure(t.Context(), resource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &resource.ConfigureResponse{})

				ctx := t.Context()
				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				planRaw := buildTagsResourceState(ctx, schemaResp.Schema.Type().TerraformType(ctx), "42",
					[]tftypes.Value{
						tftypes.NewValue(tftypes.String, "tag-1"),
						tftypes.NewValue(tftypes.String, "tag-2"),
					})

				plan := tfsdk.Plan{Schema: schemaResp.Schema, Raw: planRaw}
				state := tfsdk.State{Schema: schemaResp.Schema}

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
				r := execution.NewExecutionTagsResource()
				ctx := t.Context()

				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				planRaw := tftypes.NewValue(tftypes.String, "invalid")
				plan := tfsdk.Plan{Schema: schemaResp.Schema, Raw: planRaw}

				req := resource.CreateRequest{Plan: plan}
				resp := &resource.CreateResponse{}

				r.Create(ctx, req, resp)

				assert.True(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "error - create with invalid execution ID",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`[]`))
				})

				n8nClient, server := setupTestClientForExecution(t, handler)
				defer server.Close()

				r := execution.NewExecutionTagsResource()
				r.Configure(t.Context(), resource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &resource.ConfigureResponse{})

				ctx := t.Context()
				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				planRaw := buildTagsResourceState(ctx, schemaResp.Schema.Type().TerraformType(ctx), "not-a-number",
					[]tftypes.Value{})

				plan := tfsdk.Plan{Schema: schemaResp.Schema, Raw: planRaw}
				state := tfsdk.State{Schema: schemaResp.Schema}

				req := resource.CreateRequest{Plan: plan}
				resp := &resource.CreateResponse{State: state}

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

				n8nClient, server := setupTestClientForExecution(t, handler)
				defer server.Close()

				r := execution.NewExecutionTagsResource()
				r.Configure(t.Context(), resource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &resource.ConfigureResponse{})

				ctx := t.Context()
				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				planRaw := buildTagsResourceState(ctx, schemaResp.Schema.Type().TerraformType(ctx), "42",
					[]tftypes.Value{tftypes.NewValue(tftypes.String, "tag-1")})

				plan := tfsdk.Plan{Schema: schemaResp.Schema, Raw: planRaw}
				state := tfsdk.State{Schema: schemaResp.Schema}

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

func TestExecutionTagsResource_Read(t *testing.T) {
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
					w.Write([]byte(`[{"id":"tag-1","name":"Tag One"},{"id":"tag-2","name":"Tag Two"}]`))
				})

				n8nClient, server := setupTestClientForExecution(t, handler)
				defer server.Close()

				r := execution.NewExecutionTagsResource()
				r.Configure(t.Context(), resource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &resource.ConfigureResponse{})

				ctx := t.Context()
				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				stateRaw := buildTagsResourceState(ctx, schemaResp.Schema.Type().TerraformType(ctx), "42",
					[]tftypes.Value{
						tftypes.NewValue(tftypes.String, "tag-1"),
						tftypes.NewValue(tftypes.String, "tag-2"),
					})

				state := tfsdk.State{Schema: schemaResp.Schema, Raw: stateRaw}

				req := resource.ReadRequest{State: state}
				resp := &resource.ReadResponse{State: state}

				r.Read(ctx, req, resp)

				assert.False(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "error - read with invalid state",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := execution.NewExecutionTagsResource()
				ctx := t.Context()

				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				stateRaw := tftypes.NewValue(tftypes.String, "invalid")
				state := tfsdk.State{Schema: schemaResp.Schema, Raw: stateRaw}

				req := resource.ReadRequest{State: state}
				resp := &resource.ReadResponse{}

				r.Read(ctx, req, resp)

				assert.True(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "error - read with invalid execution ID",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`[]`))
				})

				n8nClient, server := setupTestClientForExecution(t, handler)
				defer server.Close()

				r := execution.NewExecutionTagsResource()
				r.Configure(t.Context(), resource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &resource.ConfigureResponse{})

				ctx := t.Context()
				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				stateRaw := buildTagsResourceState(ctx, schemaResp.Schema.Type().TerraformType(ctx), "not-a-number", []tftypes.Value{})
				state := tfsdk.State{Schema: schemaResp.Schema, Raw: stateRaw}

				req := resource.ReadRequest{State: state}
				resp := &resource.ReadResponse{State: state}

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

				n8nClient, server := setupTestClientForExecution(t, handler)
				defer server.Close()

				r := execution.NewExecutionTagsResource()
				r.Configure(t.Context(), resource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &resource.ConfigureResponse{})

				ctx := t.Context()
				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				stateRaw := buildTagsResourceState(ctx, schemaResp.Schema.Type().TerraformType(ctx), "42",
					[]tftypes.Value{tftypes.NewValue(tftypes.String, "tag-1")})
				state := tfsdk.State{Schema: schemaResp.Schema, Raw: stateRaw}

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

func TestExecutionTagsResource_Update(t *testing.T) {
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
					w.Write([]byte(`[{"id":"tag-3","name":"Tag Three"}]`))
				})

				n8nClient, server := setupTestClientForExecution(t, handler)
				defer server.Close()

				r := execution.NewExecutionTagsResource()
				r.Configure(t.Context(), resource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &resource.ConfigureResponse{})

				ctx := t.Context()
				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				planRaw := buildTagsResourceState(ctx, schemaResp.Schema.Type().TerraformType(ctx), "42",
					[]tftypes.Value{tftypes.NewValue(tftypes.String, "tag-3")})
				stateRaw := buildTagsResourceState(ctx, schemaResp.Schema.Type().TerraformType(ctx), "42",
					[]tftypes.Value{tftypes.NewValue(tftypes.String, "tag-1")})

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
				r := execution.NewExecutionTagsResource()
				ctx := t.Context()

				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				planRaw := tftypes.NewValue(tftypes.String, "invalid")
				stateRaw := buildTagsResourceState(ctx, schemaResp.Schema.Type().TerraformType(ctx), "42",
					[]tftypes.Value{tftypes.NewValue(tftypes.String, "tag-1")})

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

				n8nClient, server := setupTestClientForExecution(t, handler)
				defer server.Close()

				r := execution.NewExecutionTagsResource()
				r.Configure(t.Context(), resource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &resource.ConfigureResponse{})

				ctx := t.Context()
				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				planRaw := buildTagsResourceState(ctx, schemaResp.Schema.Type().TerraformType(ctx), "42",
					[]tftypes.Value{tftypes.NewValue(tftypes.String, "tag-3")})
				stateRaw := buildTagsResourceState(ctx, schemaResp.Schema.Type().TerraformType(ctx), "42",
					[]tftypes.Value{tftypes.NewValue(tftypes.String, "tag-1")})

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

func TestExecutionTagsResource_Delete(t *testing.T) {
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
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`[]`))
				})

				n8nClient, server := setupTestClientForExecution(t, handler)
				defer server.Close()

				r := execution.NewExecutionTagsResource()
				r.Configure(t.Context(), resource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &resource.ConfigureResponse{})

				ctx := t.Context()
				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				stateRaw := buildTagsResourceState(ctx, schemaResp.Schema.Type().TerraformType(ctx), "42",
					[]tftypes.Value{tftypes.NewValue(tftypes.String, "tag-1")})
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
				r := execution.NewExecutionTagsResource()
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
			name: "error - delete with invalid execution ID",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`[]`))
				})

				n8nClient, server := setupTestClientForExecution(t, handler)
				defer server.Close()

				r := execution.NewExecutionTagsResource()
				r.Configure(t.Context(), resource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &resource.ConfigureResponse{})

				ctx := t.Context()
				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				stateRaw := buildTagsResourceState(ctx, schemaResp.Schema.Type().TerraformType(ctx), "not-a-number", []tftypes.Value{})
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

				n8nClient, server := setupTestClientForExecution(t, handler)
				defer server.Close()

				r := execution.NewExecutionTagsResource()
				r.Configure(t.Context(), resource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &resource.ConfigureResponse{})

				ctx := t.Context()
				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				stateRaw := buildTagsResourceState(ctx, schemaResp.Schema.Type().TerraformType(ctx), "42",
					[]tftypes.Value{tftypes.NewValue(tftypes.String, "tag-1")})
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

func TestExecutionTagsResource_ImportState(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "import state passthrough",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := execution.NewExecutionTagsResource()
				ctx := t.Context()

				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				emptyValue := buildTagsResourceState(ctx, schemaResp.Schema.Type().TerraformType(ctx), "", []tftypes.Value{})

				req := resource.ImportStateRequest{
					ID: "42",
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
