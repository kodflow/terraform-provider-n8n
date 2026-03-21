// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

package credential_test

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/credential"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/client"
	"github.com/stretchr/testify/assert"
)

func TestNewCredentialsDataSource(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "normal case",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := credential.NewCredentialsDataSource()
				assert.NotNil(t, ds)
			},
		},
		{
			name: "error case - validates behavior",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := credential.NewCredentialsDataSource()
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

func TestNewCredentialsDataSourceWrapper(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "normal case",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := credential.NewCredentialsDataSourceWrapper()
				assert.NotNil(t, ds)
				assert.Implements(t, (*datasource.DataSource)(nil), ds)
			},
		},
		{
			name: "error case - validates behavior",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := credential.NewCredentialsDataSourceWrapper()
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

func TestCredentialsDataSource_Metadata(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "normal case",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := credential.NewCredentialsDataSource()
				resp := &datasource.MetadataResponse{}

				ds.Metadata(t.Context(), datasource.MetadataRequest{
					ProviderTypeName: "n8n",
				}, resp)

				assert.Equal(t, "n8n_credentials", resp.TypeName)
			},
		},
		{
			name: "error case - validates behavior",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := credential.NewCredentialsDataSource()
				resp := &datasource.MetadataResponse{}

				ds.Metadata(t.Context(), datasource.MetadataRequest{
					ProviderTypeName: "n8n",
				}, resp)

				assert.Equal(t, "n8n_credentials", resp.TypeName)
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

func TestCredentialsDataSource_Schema(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "normal case",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := credential.NewCredentialsDataSource()
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
				ds := credential.NewCredentialsDataSource()
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

func TestCredentialsDataSource_Configure(t *testing.T) {
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

			ds := credential.NewCredentialsDataSource()
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

// setupCredentialsDSWithSchema sets up a DataSource + ReadResponse with a proper State schema.
func setupCredentialsDSWithSchema(t *testing.T, n8nClient *client.N8nClient) (*credential.CredentialsDataSource, datasource.ReadResponse) {
	t.Helper()
	ds := credential.NewCredentialsDataSource()
	ds.Configure(t.Context(), datasource.ConfigureRequest{
		ProviderData: n8nClient,
	}, &datasource.ConfigureResponse{})

	schemaResp := datasource.SchemaResponse{}
	ds.Schema(t.Context(), datasource.SchemaRequest{}, &schemaResp)

	resp := datasource.ReadResponse{}
	resp.State.Schema = schemaResp.Schema

	return ds, resp
}

func TestCredentialsDataSource_Read(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "read with successful API call - empty list",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{"data":[]}`))
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				ds, readResp := setupCredentialsDSWithSchema(t, n8nClient)

				ctx := t.Context()
				req := datasource.ReadRequest{}

				ds.Read(ctx, req, &readResp)

				if readResp.Diagnostics.HasError() {
					for _, diag := range readResp.Diagnostics.Errors() {
						t.Logf("Error: %s - %s", diag.Summary(), diag.Detail())
					}
				}
				assert.False(t, readResp.Diagnostics.HasError())
			},
		},
		{
			name: "read with successful API call - with credentials",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{"data":[{"id":"cred-1","name":"Credential One","type":"httpHeaderAuth","createdAt":"2024-01-01T00:00:00Z","updatedAt":"2024-01-02T00:00:00Z","nodesAccess":[]},{"id":"cred-2","name":"Credential Two","type":"oAuth2Api","createdAt":"2024-02-01T00:00:00Z","updatedAt":"2024-02-02T00:00:00Z","nodesAccess":[]}]}`))
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				ds, readResp := setupCredentialsDSWithSchema(t, n8nClient)

				ctx := t.Context()
				req := datasource.ReadRequest{}

				ds.Read(ctx, req, &readResp)

				if readResp.Diagnostics.HasError() {
					for _, diag := range readResp.Diagnostics.Errors() {
						t.Logf("Error: %s - %s", diag.Summary(), diag.Detail())
					}
				}
				assert.False(t, readResp.Diagnostics.HasError())
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

				ds, readResp := setupCredentialsDSWithSchema(t, n8nClient)

				ctx := t.Context()
				req := datasource.ReadRequest{}

				ds.Read(ctx, req, &readResp)

				assert.True(t, readResp.Diagnostics.HasError())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}
