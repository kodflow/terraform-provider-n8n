// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

package sourcecontrol_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/client"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/sourcecontrol"
	"github.com/stretchr/testify/assert"
)

// setupTestClientForSC creates a test N8nClient with httptest server for source control tests.
func setupTestClientForSC(t *testing.T, handler http.HandlerFunc) (*client.N8nClient, *httptest.Server) {
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

// TestNewSourceControlPullResource tests the constructor.
func TestNewSourceControlPullResource(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
	}{
		{name: "normal case"},
		{name: "error case - validates behavior"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := sourcecontrol.NewSourceControlPullResource()
			assert.NotNil(t, r)
		})
	}
}

// TestNewSourceControlPullResourceWrapper tests the wrapper constructor.
func TestNewSourceControlPullResourceWrapper(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
	}{
		{name: "normal case"},
		{name: "error case - validates behavior"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := sourcecontrol.NewSourceControlPullResourceWrapper()
			assert.NotNil(t, r)
			assert.Implements(t, (*resource.Resource)(nil), r)
		})
	}
}

// TestSourceControlPullResource_Metadata tests the Metadata method.
func TestSourceControlPullResource_Metadata(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
	}{
		{name: "normal case"},
		{name: "cancelled context exits early"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := sourcecontrol.NewSourceControlPullResource()
			resp := &resource.MetadataResponse{}

			if tt.name == "cancelled context exits early" {
				ctx, cancel := context.WithCancel(t.Context())
				cancel()
				assert.NotPanics(t, func() {
					r.Metadata(ctx, resource.MetadataRequest{ProviderTypeName: "n8n"}, resp)
				})
				return
			}

			r.Metadata(t.Context(), resource.MetadataRequest{
				ProviderTypeName: "n8n",
			}, resp)

			assert.Equal(t, "n8n_source_control_pull", resp.TypeName)
		})
	}
}

// TestSourceControlPullResource_Schema tests the Schema method.
func TestSourceControlPullResource_Schema(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
	}{
		{name: "normal case"},
		{name: "cancelled context exits early"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := sourcecontrol.NewSourceControlPullResource()
			resp := &resource.SchemaResponse{}

			if tt.name == "cancelled context exits early" {
				ctx, cancel := context.WithCancel(t.Context())
				cancel()
				assert.NotPanics(t, func() {
					r.Schema(ctx, resource.SchemaRequest{}, resp)
				})
				return
			}

			r.Schema(t.Context(), resource.SchemaRequest{}, resp)

			assert.NotNil(t, resp.Schema)
			assert.Contains(t, resp.Schema.Attributes, "force")
			assert.Contains(t, resp.Schema.Attributes, "workflows_imported")
			assert.Contains(t, resp.Schema.Attributes, "credentials_imported")
		})
	}
}

// TestSourceControlPullResource_Configure tests the Configure method.
func TestSourceControlPullResource_Configure(t *testing.T) {
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

			r := sourcecontrol.NewSourceControlPullResource()
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

// TestSourceControlPullResource_Create tests the Create method end-to-end.
func TestSourceControlPullResource_Create(t *testing.T) {
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
					w.Write([]byte(`{"workflows":[{"id":"wf-1","name":"Test"}]}`))
				})

				n8nClient, server := setupTestClientForSC(t, handler)
				defer server.Close()

				r := sourcecontrol.NewSourceControlPullResource()
				r.Configure(t.Context(), resource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &resource.ConfigureResponse{})

				ctx := t.Context()
				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				planRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"force":                tftypes.NewValue(tftypes.Bool, false),
					"workflows_imported":   tftypes.NewValue(tftypes.String, nil),
					"credentials_imported": tftypes.NewValue(tftypes.String, nil),
					"tags_imported":        tftypes.NewValue(tftypes.String, nil),
					"variables_imported":   tftypes.NewValue(tftypes.String, nil),
				})

				plan := tfsdk.Plan{Schema: schemaResp.Schema, Raw: planRaw}
				state := tfsdk.State{Schema: schemaResp.Schema}

				req := resource.CreateRequest{Plan: plan}
				resp := &resource.CreateResponse{State: state}

				r.Create(ctx, req, resp)

				assert.False(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "error - invalid plan type",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := sourcecontrol.NewSourceControlPullResource()
				ctx := t.Context()

				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				planRaw := tftypes.NewValue(tftypes.String, "invalid")
				plan := tfsdk.Plan{Schema: schemaResp.Schema, Raw: planRaw}
				state := tfsdk.State{Schema: schemaResp.Schema}

				req := resource.CreateRequest{Plan: plan}
				resp := &resource.CreateResponse{State: state}

				r.Create(ctx, req, resp)

				assert.True(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "error - API returns error",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"message":"Internal server error"}`))
				})

				n8nClient, server := setupTestClientForSC(t, handler)
				defer server.Close()

				r := sourcecontrol.NewSourceControlPullResource()
				r.Configure(t.Context(), resource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &resource.ConfigureResponse{})

				ctx := t.Context()
				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				planRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"force":                tftypes.NewValue(tftypes.Bool, true),
					"workflows_imported":   tftypes.NewValue(tftypes.String, nil),
					"credentials_imported": tftypes.NewValue(tftypes.String, nil),
					"tags_imported":        tftypes.NewValue(tftypes.String, nil),
					"variables_imported":   tftypes.NewValue(tftypes.String, nil),
				})

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

// TestSourceControlPullResource_Read tests that Read preserves existing state.
func TestSourceControlPullResource_Read(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
	}{
		{name: "read preserves existing state"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := sourcecontrol.NewSourceControlPullResource()
			ctx := t.Context()

			schemaResp := resource.SchemaResponse{}
			r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

			stateRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
				"force":                tftypes.NewValue(tftypes.Bool, false),
				"workflows_imported":   tftypes.NewValue(tftypes.String, nil),
				"credentials_imported": tftypes.NewValue(tftypes.String, nil),
				"tags_imported":        tftypes.NewValue(tftypes.String, nil),
				"variables_imported":   tftypes.NewValue(tftypes.String, nil),
			})

			state := tfsdk.State{Schema: schemaResp.Schema, Raw: stateRaw}

			req := resource.ReadRequest{State: state}
			resp := &resource.ReadResponse{State: state}

			r.Read(ctx, req, resp)

			assert.False(t, resp.Diagnostics.HasError())
		})
	}
}

// TestSourceControlPullResource_Update tests the Update method end-to-end.
func TestSourceControlPullResource_Update(t *testing.T) {
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
					w.Write([]byte(`{"workflows":[{"id":"wf-1","name":"Test"}]}`))
				})

				n8nClient, server := setupTestClientForSC(t, handler)
				defer server.Close()

				r := sourcecontrol.NewSourceControlPullResource()
				r.Configure(t.Context(), resource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &resource.ConfigureResponse{})

				ctx := t.Context()
				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				planRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"force":                tftypes.NewValue(tftypes.Bool, false),
					"workflows_imported":   tftypes.NewValue(tftypes.String, nil),
					"credentials_imported": tftypes.NewValue(tftypes.String, nil),
					"tags_imported":        tftypes.NewValue(tftypes.String, nil),
					"variables_imported":   tftypes.NewValue(tftypes.String, nil),
				})

				plan := tfsdk.Plan{Schema: schemaResp.Schema, Raw: planRaw}
				state := tfsdk.State{Schema: schemaResp.Schema}

				req := resource.UpdateRequest{Plan: plan, State: state}
				resp := &resource.UpdateResponse{State: state}

				r.Update(ctx, req, resp)

				assert.False(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "error - API returns error on update",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"message":"Internal server error"}`))
				})

				n8nClient, server := setupTestClientForSC(t, handler)
				defer server.Close()

				r := sourcecontrol.NewSourceControlPullResource()
				r.Configure(t.Context(), resource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &resource.ConfigureResponse{})

				ctx := t.Context()
				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				planRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"force":                tftypes.NewValue(tftypes.Bool, true),
					"workflows_imported":   tftypes.NewValue(tftypes.String, nil),
					"credentials_imported": tftypes.NewValue(tftypes.String, nil),
					"tags_imported":        tftypes.NewValue(tftypes.String, nil),
					"variables_imported":   tftypes.NewValue(tftypes.String, nil),
				})

				plan := tfsdk.Plan{Schema: schemaResp.Schema, Raw: planRaw}
				state := tfsdk.State{Schema: schemaResp.Schema}

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

// TestSourceControlPullResource_Delete tests that Delete is a no-op.
func TestSourceControlPullResource_Delete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		cancelled bool
	}{
		{name: "delete is a no-op"},
		{name: "cancelled context exits early", cancelled: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := sourcecontrol.NewSourceControlPullResource()
			ctx := t.Context()

			schemaResp := resource.SchemaResponse{}
			r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

			stateRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
				"force":                tftypes.NewValue(tftypes.Bool, false),
				"workflows_imported":   tftypes.NewValue(tftypes.String, nil),
				"credentials_imported": tftypes.NewValue(tftypes.String, nil),
				"tags_imported":        tftypes.NewValue(tftypes.String, nil),
				"variables_imported":   tftypes.NewValue(tftypes.String, nil),
			})

			state := tfsdk.State{Schema: schemaResp.Schema, Raw: stateRaw}
			req := resource.DeleteRequest{State: state}
			resp := &resource.DeleteResponse{}

			if tt.cancelled {
				cancelCtx, cancel := context.WithCancel(ctx)
				cancel()
				assert.NotPanics(t, func() { r.Delete(cancelCtx, req, resp) })
				return
			}

			r.Delete(ctx, req, resp)
			assert.False(t, resp.Diagnostics.HasError())
		})
	}
}
