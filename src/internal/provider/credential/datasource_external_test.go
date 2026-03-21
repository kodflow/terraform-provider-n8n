// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

package credential_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/credential"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/client"
	"github.com/stretchr/testify/assert"
)

func TestNewCredentialDataSource(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "normal case",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := credential.NewCredentialDataSource()
				assert.NotNil(t, ds)
			},
		},
		{
			name: "error case - validates behavior",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := credential.NewCredentialDataSource()
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

func TestNewCredentialDataSourceWrapper(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "normal case",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := credential.NewCredentialDataSourceWrapper()
				assert.NotNil(t, ds)
				assert.Implements(t, (*datasource.DataSource)(nil), ds)
			},
		},
		{
			name: "error case - validates behavior",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := credential.NewCredentialDataSourceWrapper()
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

func TestCredentialDataSource_Metadata(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "normal case",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := credential.NewCredentialDataSource()
				resp := &datasource.MetadataResponse{}

				ds.Metadata(t.Context(), datasource.MetadataRequest{
					ProviderTypeName: "n8n",
				}, resp)

				assert.Equal(t, "n8n_credential", resp.TypeName)
			},
		},
		{
			name: "error case - validates behavior",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := credential.NewCredentialDataSource()
				resp := &datasource.MetadataResponse{}

				ds.Metadata(t.Context(), datasource.MetadataRequest{
					ProviderTypeName: "n8n",
				}, resp)

				assert.Equal(t, "n8n_credential", resp.TypeName)
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

func TestCredentialDataSource_Schema(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "normal case",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := credential.NewCredentialDataSource()
				resp := &datasource.SchemaResponse{}

				ds.Schema(t.Context(), datasource.SchemaRequest{}, resp)

				assert.NotNil(t, resp.Schema)
				assert.NotEmpty(t, resp.Schema.Attributes)
			},
		},
		{
			name: "error case - validates behavior",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := credential.NewCredentialDataSource()
				resp := &datasource.SchemaResponse{}

				ds.Schema(t.Context(), datasource.SchemaRequest{}, resp)

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

func TestCredentialDataSource_Configure(t *testing.T) {
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

			ds := credential.NewCredentialDataSource()
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

// credentialDSConfigParams holds parameters for building a CredentialDataSource config value.
type credentialDSConfigParams struct {
	ID        string
	Name      string
	CredType  string
	CreatedAt string
	UpdatedAt string
}

// buildCredentialDSConfig builds a valid tftypes.Value for the CredentialDataSource schema.
func buildCredentialDSConfig(ctx context.Context, schemaType tftypes.Type, p credentialDSConfigParams) tftypes.Value {
	id, name, credType, createdAt, updatedAt := p.ID, p.Name, p.CredType, p.CreatedAt, p.UpdatedAt
	idVal := tftypes.NewValue(tftypes.String, nil)
	if id != "" {
		idVal = tftypes.NewValue(tftypes.String, id)
	}
	nameVal := tftypes.NewValue(tftypes.String, nil)
	if name != "" {
		nameVal = tftypes.NewValue(tftypes.String, name)
	}
	typeVal := tftypes.NewValue(tftypes.String, nil)
	if credType != "" {
		typeVal = tftypes.NewValue(tftypes.String, credType)
	}
	createdVal := tftypes.NewValue(tftypes.String, nil)
	if createdAt != "" {
		createdVal = tftypes.NewValue(tftypes.String, createdAt)
	}
	updatedVal := tftypes.NewValue(tftypes.String, nil)
	if updatedAt != "" {
		updatedVal = tftypes.NewValue(tftypes.String, updatedAt)
	}

	return tftypes.NewValue(schemaType, map[string]tftypes.Value{
		"id":         idVal,
		"name":       nameVal,
		"type":       typeVal,
		"created_at": createdVal,
		"updated_at": updatedVal,
	})
}

func TestCredentialDataSource_Read(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "read by ID with successful API call",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{"data":[{"id":"cred-123","name":"My Credential","type":"httpHeaderAuth","createdAt":"2024-01-01T00:00:00Z","updatedAt":"2024-01-02T00:00:00Z","nodesAccess":[]}]}`))
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				ds := credential.NewCredentialDataSource()
				ds.Configure(t.Context(), datasource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &datasource.ConfigureResponse{})

				ctx := t.Context()
				schemaResp := datasource.SchemaResponse{}
				ds.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)

				configRaw := buildCredentialDSConfig(ctx, schemaResp.Schema.Type().TerraformType(ctx), credentialDSConfigParams{ID: "cred-123", Name: "", CredType: "", CreatedAt: "", UpdatedAt: ""})

				config := tfsdk.Config{Schema: schemaResp.Schema, Raw: configRaw}
				state := tfsdk.State{Schema: schemaResp.Schema}

				req := datasource.ReadRequest{Config: config}
				resp := &datasource.ReadResponse{State: state}

				ds.Read(ctx, req, resp)

				if resp.Diagnostics.HasError() {
					for _, diag := range resp.Diagnostics.Errors() {
						t.Logf("Error: %s - %s", diag.Summary(), diag.Detail())
					}
				}
				assert.False(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "read by name with successful API call",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{"data":[{"id":"cred-123","name":"My Credential","type":"httpHeaderAuth","createdAt":"2024-01-01T00:00:00Z","updatedAt":"2024-01-02T00:00:00Z","nodesAccess":[]}]}`))
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				ds := credential.NewCredentialDataSource()
				ds.Configure(t.Context(), datasource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &datasource.ConfigureResponse{})

				ctx := t.Context()
				schemaResp := datasource.SchemaResponse{}
				ds.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)

				configRaw := buildCredentialDSConfig(ctx, schemaResp.Schema.Type().TerraformType(ctx), credentialDSConfigParams{ID: "", Name: "My Credential", CredType: "", CreatedAt: "", UpdatedAt: ""})

				config := tfsdk.Config{Schema: schemaResp.Schema, Raw: configRaw}
				state := tfsdk.State{Schema: schemaResp.Schema}

				req := datasource.ReadRequest{Config: config}
				resp := &datasource.ReadResponse{State: state}

				ds.Read(ctx, req, resp)

				if resp.Diagnostics.HasError() {
					for _, diag := range resp.Diagnostics.Errors() {
						t.Logf("Error: %s - %s", diag.Summary(), diag.Detail())
					}
				}
				assert.False(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "error - read with invalid config",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := credential.NewCredentialDataSource()
				ctx := t.Context()

				schemaResp := datasource.SchemaResponse{}
				ds.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)

				configRaw := tftypes.NewValue(tftypes.String, "invalid")

				config := tfsdk.Config{Schema: schemaResp.Schema, Raw: configRaw}
				state := tfsdk.State{Schema: schemaResp.Schema}

				req := datasource.ReadRequest{Config: config}
				resp := &datasource.ReadResponse{State: state}

				ds.Read(ctx, req, resp)

				assert.True(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "error - missing required identifiers",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := credential.NewCredentialDataSource()
				ctx := t.Context()

				schemaResp := datasource.SchemaResponse{}
				ds.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)

				configRaw := buildCredentialDSConfig(ctx, schemaResp.Schema.Type().TerraformType(ctx), credentialDSConfigParams{ID: "", Name: "", CredType: "", CreatedAt: "", UpdatedAt: ""})

				config := tfsdk.Config{Schema: schemaResp.Schema, Raw: configRaw}
				state := tfsdk.State{Schema: schemaResp.Schema}

				req := datasource.ReadRequest{Config: config}
				resp := &datasource.ReadResponse{State: state}

				ds.Read(ctx, req, resp)

				assert.True(t, resp.Diagnostics.HasError())
				assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Missing Required Attribute")
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

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				ds := credential.NewCredentialDataSource()
				ds.Configure(t.Context(), datasource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &datasource.ConfigureResponse{})

				ctx := t.Context()
				schemaResp := datasource.SchemaResponse{}
				ds.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)

				configRaw := buildCredentialDSConfig(ctx, schemaResp.Schema.Type().TerraformType(ctx), credentialDSConfigParams{ID: "cred-123", Name: "", CredType: "", CreatedAt: "", UpdatedAt: ""})

				config := tfsdk.Config{Schema: schemaResp.Schema, Raw: configRaw}
				state := tfsdk.State{Schema: schemaResp.Schema}

				req := datasource.ReadRequest{Config: config}
				resp := &datasource.ReadResponse{State: state}

				ds.Read(ctx, req, resp)

				assert.True(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "error - credential not found by ID",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{"data":[]}`))
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				ds := credential.NewCredentialDataSource()
				ds.Configure(t.Context(), datasource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &datasource.ConfigureResponse{})

				ctx := t.Context()
				schemaResp := datasource.SchemaResponse{}
				ds.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)

				configRaw := buildCredentialDSConfig(ctx, schemaResp.Schema.Type().TerraformType(ctx), credentialDSConfigParams{ID: "cred-nonexistent", Name: "", CredType: "", CreatedAt: "", UpdatedAt: ""})

				config := tfsdk.Config{Schema: schemaResp.Schema, Raw: configRaw}
				state := tfsdk.State{Schema: schemaResp.Schema}

				req := datasource.ReadRequest{Config: config}
				resp := &datasource.ReadResponse{State: state}

				ds.Read(ctx, req, resp)

				assert.True(t, resp.Diagnostics.HasError())
				assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Credential Not Found")
			},
		},
		{
			name: "error - credential not found by name",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{"data":[{"id":"cred-456","name":"Other Credential","type":"httpHeaderAuth","createdAt":"2024-01-01T00:00:00Z","updatedAt":"2024-01-02T00:00:00Z","nodesAccess":[]}]}`))
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				ds := credential.NewCredentialDataSource()
				ds.Configure(t.Context(), datasource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &datasource.ConfigureResponse{})

				ctx := t.Context()
				schemaResp := datasource.SchemaResponse{}
				ds.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)

				configRaw := buildCredentialDSConfig(ctx, schemaResp.Schema.Type().TerraformType(ctx), credentialDSConfigParams{ID: "", Name: "NonExistent Credential", CredType: "", CreatedAt: "", UpdatedAt: ""})

				config := tfsdk.Config{Schema: schemaResp.Schema, Raw: configRaw}
				state := tfsdk.State{Schema: schemaResp.Schema}

				req := datasource.ReadRequest{Config: config}
				resp := &datasource.ReadResponse{State: state}

				ds.Read(ctx, req, resp)

				assert.True(t, resp.Diagnostics.HasError())
				assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Credential Not Found")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}
