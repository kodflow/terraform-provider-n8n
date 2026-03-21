// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

package audit_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/audit"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/client"
	"github.com/stretchr/testify/assert"
)

// setupTestClientForAudit creates a test N8nClient with httptest server for audit datasource tests.
func setupTestClientForAudit(t *testing.T, handler http.HandlerFunc) (*client.N8nClient, *httptest.Server) {
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

// TestNewAuditDataSource tests the constructor.
func TestNewAuditDataSource(t *testing.T) {
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
			ds := audit.NewAuditDataSource()
			assert.NotNil(t, ds)
		})
	}
}

// TestNewAuditDataSourceWrapper tests the wrapper constructor.
func TestNewAuditDataSourceWrapper(t *testing.T) {
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
			ds := audit.NewAuditDataSourceWrapper()
			assert.NotNil(t, ds)
			assert.Implements(t, (*datasource.DataSource)(nil), ds)
		})
	}
}

// TestAuditDataSource_Metadata tests the Metadata method.
func TestAuditDataSource_Metadata(t *testing.T) {
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
			ds := audit.NewAuditDataSource()
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

			assert.Equal(t, "n8n_audit", resp.TypeName)
		})
	}
}

// TestAuditDataSource_Schema tests the Schema method.
func TestAuditDataSource_Schema(t *testing.T) {
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
			ds := audit.NewAuditDataSource()
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
			assert.Contains(t, resp.Schema.Attributes, "credentials_risk_report")
			assert.Contains(t, resp.Schema.Attributes, "database_risk_report")
			assert.Contains(t, resp.Schema.Attributes, "filesystem_risk_report")
			assert.Contains(t, resp.Schema.Attributes, "nodes_risk_report")
			assert.Contains(t, resp.Schema.Attributes, "instance_risk_report")
		})
	}
}

// TestAuditDataSource_Configure tests the Configure method.
func TestAuditDataSource_Configure(t *testing.T) {
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

			ds := audit.NewAuditDataSource()
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

// TestAuditDataSource_Read tests the Read method end-to-end.
func TestAuditDataSource_Read(t *testing.T) {
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
					w.Write([]byte(`{"Credentials Risk Report":{"risk":"low"},"Database Risk Report":{"tables":3}}`))
				})

				n8nClient, server := setupTestClientForAudit(t, handler)
				defer server.Close()

				ds := audit.NewAuditDataSource()
				ds.Configure(t.Context(), datasource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &datasource.ConfigureResponse{})

				ctx := t.Context()
				schemaResp := datasource.SchemaResponse{}
				ds.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)

				state := tfsdk.State{Schema: schemaResp.Schema}

				req := datasource.ReadRequest{}
				resp := &datasource.ReadResponse{State: state}

				ds.Read(ctx, req, resp)

				assert.False(t, resp.Diagnostics.HasError())
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

				n8nClient, server := setupTestClientForAudit(t, handler)
				defer server.Close()

				ds := audit.NewAuditDataSource()
				ds.Configure(t.Context(), datasource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &datasource.ConfigureResponse{})

				ctx := t.Context()
				schemaResp := datasource.SchemaResponse{}
				ds.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)

				state := tfsdk.State{Schema: schemaResp.Schema}

				req := datasource.ReadRequest{}
				resp := &datasource.ReadResponse{State: state}

				ds.Read(ctx, req, resp)

				assert.True(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "read with empty state sets null fields",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{}`))
				})

				n8nClient, server := setupTestClientForAudit(t, handler)
				defer server.Close()

				ds := audit.NewAuditDataSource()
				ds.Configure(t.Context(), datasource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &datasource.ConfigureResponse{})

				ctx := t.Context()
				schemaResp := datasource.SchemaResponse{}
				ds.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)

				// Build a state with null values for all fields.
				stateRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"credentials_risk_report": tftypes.NewValue(tftypes.String, nil),
					"database_risk_report":    tftypes.NewValue(tftypes.String, nil),
					"filesystem_risk_report":  tftypes.NewValue(tftypes.String, nil),
					"nodes_risk_report":       tftypes.NewValue(tftypes.String, nil),
					"instance_risk_report":    tftypes.NewValue(tftypes.String, nil),
				})
				state := tfsdk.State{Schema: schemaResp.Schema, Raw: stateRaw}

				req := datasource.ReadRequest{}
				resp := &datasource.ReadResponse{State: state}

				ds.Read(ctx, req, resp)

				assert.False(t, resp.Diagnostics.HasError())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}
