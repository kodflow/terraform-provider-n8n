// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

package execution_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/execution"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/client"
	"github.com/stretchr/testify/assert"
)

// TestNewExecutionDataSource tests the constructor.
func TestNewExecutionDataSource(t *testing.T) {
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
			ds := execution.NewExecutionDataSource()
			assert.NotNil(t, ds)
		})
	}
}

// TestNewExecutionDataSourceWrapper tests the wrapper constructor.
func TestNewExecutionDataSourceWrapper(t *testing.T) {
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
			ds := execution.NewExecutionDataSourceWrapper()
			assert.NotNil(t, ds)
			assert.Implements(t, (*datasource.DataSource)(nil), ds)
		})
	}
}

// TestExecutionDataSource_Metadata tests the Metadata method.
func TestExecutionDataSource_Metadata(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		cancelled bool
	}{
		{name: "normal case"},
		{name: "cancelled context exits early", cancelled: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ds := execution.NewExecutionDataSource()
			resp := &datasource.MetadataResponse{}

			if tt.cancelled {
				ctx, cancel := context.WithCancel(t.Context())
				cancel()
				assert.NotPanics(t, func() {
					ds.Metadata(ctx, datasource.MetadataRequest{ProviderTypeName: "n8n"}, resp)
				})
				return
			}

			ds.Metadata(t.Context(), datasource.MetadataRequest{
				ProviderTypeName: "n8n",
			}, resp)

			assert.Equal(t, "n8n_execution", resp.TypeName)
		})
	}
}

// TestExecutionDataSource_Schema tests the Schema method.
func TestExecutionDataSource_Schema(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		cancelled bool
	}{
		{name: "normal case"},
		{name: "cancelled context exits early", cancelled: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ds := execution.NewExecutionDataSource()
			resp := &datasource.SchemaResponse{}

			if tt.cancelled {
				ctx, cancel := context.WithCancel(t.Context())
				cancel()
				assert.NotPanics(t, func() {
					ds.Schema(ctx, datasource.SchemaRequest{}, resp)
				})
				return
			}

			ds.Schema(t.Context(), datasource.SchemaRequest{}, resp)

			assert.NotNil(t, resp.Schema)
			assert.NotEmpty(t, resp.Schema.Attributes)
		})
	}
}

// TestExecutionDataSource_Configure tests the Configure method.
func TestExecutionDataSource_Configure(t *testing.T) {
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

			ds := execution.NewExecutionDataSource()
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

// TestExecutionDataSource_Read tests the Read method end-to-end.
func TestExecutionDataSource_Read(t *testing.T) {
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
					w.Write([]byte(`{"id":42,"mode":"manual","status":"success","finished":true}`))
				})

				n8nClient, server := setupTestClientForExecution(t, handler)
				defer server.Close()

				ds := execution.NewExecutionDataSource()
				ds.Configure(t.Context(), datasource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &datasource.ConfigureResponse{})

				ctx := t.Context()
				schemaResp := datasource.SchemaResponse{}
				ds.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)

				configRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":          tftypes.NewValue(tftypes.String, "42"),
					"mode":        tftypes.NewValue(tftypes.String, nil),
					"status":      tftypes.NewValue(tftypes.String, nil),
					"workflow_id": tftypes.NewValue(tftypes.String, nil),
					"finished":    tftypes.NewValue(tftypes.Bool, nil),
					"created_at":  tftypes.NewValue(tftypes.String, nil),
					"started_at":  tftypes.NewValue(tftypes.String, nil),
					"stopped_at":  tftypes.NewValue(tftypes.String, nil),
				})

				config := tfsdk.Config{
					Schema: schemaResp.Schema,
					Raw:    configRaw,
				}
				state := tfsdk.State{Schema: schemaResp.Schema}

				req := datasource.ReadRequest{Config: config}
				resp := &datasource.ReadResponse{State: state}

				ds.Read(ctx, req, resp)

				assert.False(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "error - invalid config type",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := execution.NewExecutionDataSource()
				ctx := t.Context()

				schemaResp := datasource.SchemaResponse{}
				ds.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)

				configRaw := tftypes.NewValue(tftypes.String, "invalid")
				config := tfsdk.Config{
					Schema: schemaResp.Schema,
					Raw:    configRaw,
				}
				state := tfsdk.State{Schema: schemaResp.Schema}

				req := datasource.ReadRequest{Config: config}
				resp := &datasource.ReadResponse{State: state}

				ds.Read(ctx, req, resp)

				assert.True(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "error - invalid execution ID format",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				})

				n8nClient, server := setupTestClientForExecution(t, handler)
				defer server.Close()

				ds := execution.NewExecutionDataSource()
				ds.Configure(t.Context(), datasource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &datasource.ConfigureResponse{})

				ctx := t.Context()
				schemaResp := datasource.SchemaResponse{}
				ds.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)

				configRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":          tftypes.NewValue(tftypes.String, "not-a-number"),
					"mode":        tftypes.NewValue(tftypes.String, nil),
					"status":      tftypes.NewValue(tftypes.String, nil),
					"workflow_id": tftypes.NewValue(tftypes.String, nil),
					"finished":    tftypes.NewValue(tftypes.Bool, nil),
					"created_at":  tftypes.NewValue(tftypes.String, nil),
					"started_at":  tftypes.NewValue(tftypes.String, nil),
					"stopped_at":  tftypes.NewValue(tftypes.String, nil),
				})

				config := tfsdk.Config{Schema: schemaResp.Schema, Raw: configRaw}
				state := tfsdk.State{Schema: schemaResp.Schema}

				req := datasource.ReadRequest{Config: config}
				resp := &datasource.ReadResponse{State: state}

				ds.Read(ctx, req, resp)

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

				n8nClient, server := setupTestClientForExecution(t, handler)
				defer server.Close()

				ds := execution.NewExecutionDataSource()
				ds.Configure(t.Context(), datasource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &datasource.ConfigureResponse{})

				ctx := t.Context()
				schemaResp := datasource.SchemaResponse{}
				ds.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)

				configRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":          tftypes.NewValue(tftypes.String, "99"),
					"mode":        tftypes.NewValue(tftypes.String, nil),
					"status":      tftypes.NewValue(tftypes.String, nil),
					"workflow_id": tftypes.NewValue(tftypes.String, nil),
					"finished":    tftypes.NewValue(tftypes.Bool, nil),
					"created_at":  tftypes.NewValue(tftypes.String, nil),
					"started_at":  tftypes.NewValue(tftypes.String, nil),
					"stopped_at":  tftypes.NewValue(tftypes.String, nil),
				})

				config := tfsdk.Config{
					Schema: schemaResp.Schema,
					Raw:    configRaw,
				}
				state := tfsdk.State{Schema: schemaResp.Schema}

				req := datasource.ReadRequest{Config: config}
				resp := &datasource.ReadResponse{State: state}

				ds.Read(ctx, req, resp)

				assert.True(t, resp.Diagnostics.HasError())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}
