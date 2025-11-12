package credential_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/credential"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
	"github.com/stretchr/testify/assert"
)

// External tests (black-box testing) go here.
// These tests only have access to exported functions and types.

func TestNewCredentialTransferResource(t *testing.T) {
	t.Parallel()

	resource := credential.NewCredentialTransferResource()
	assert.NotNil(t, resource, "NewCredentialTransferResource should not return nil")
}

func TestNewCredentialTransferResourceWrapper(t *testing.T) {
	t.Parallel()

	resource := credential.NewCredentialTransferResourceWrapper()
	assert.NotNil(t, resource, "NewCredentialTransferResourceWrapper should not return nil")
}

func TestCredentialTransferResource_Metadata(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		providerTypeName string
		expectedTypeName string
	}{
		{
			name:             "with n8n provider",
			providerTypeName: "n8n",
			expectedTypeName: "n8n_credential_transfer",
		},
		{
			name:             "with custom provider",
			providerTypeName: "custom",
			expectedTypeName: "custom_credential_transfer",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := credential.NewCredentialTransferResource()
			req := resource.MetadataRequest{
				ProviderTypeName: tt.providerTypeName,
			}
			resp := &resource.MetadataResponse{}

			r.Metadata(context.Background(), req, resp)

			assert.Equal(t, tt.expectedTypeName, resp.TypeName)
		})
	}
}

func TestCredentialTransferResource_Schema(t *testing.T) {
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

			r := credential.NewCredentialTransferResource()
			req := resource.SchemaRequest{}
			resp := &resource.SchemaResponse{}

			r.Schema(context.Background(), req, resp)

			assert.NotNil(t, resp.Schema)
			assert.Contains(t, resp.Schema.MarkdownDescription, "transfer")

			// Verify required attributes exist
			attrs := resp.Schema.Attributes
			assert.Contains(t, attrs, "id")
			assert.Contains(t, attrs, "credential_id")
			assert.Contains(t, attrs, "destination_project_id")
			assert.Contains(t, attrs, "transferred_at")
		})
	}
}

func TestCredentialTransferResource_Configure(t *testing.T) {
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
				r := credential.NewCredentialTransferResource()
				req := resource.ConfigureRequest{
					ProviderData: &client.N8nClient{},
				}
				resp := &resource.ConfigureResponse{}
				r.Configure(context.Background(), req, resp)
				assert.False(t, resp.Diagnostics.HasError())

			case "error case - nil provider data":
				r := credential.NewCredentialTransferResource()
				req := resource.ConfigureRequest{
					ProviderData: nil,
				}
				resp := &resource.ConfigureResponse{}
				r.Configure(context.Background(), req, resp)
				assert.False(t, resp.Diagnostics.HasError())

			case "error case - wrong provider data type":
				r := credential.NewCredentialTransferResource()
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

func TestCredentialTransferResource_Create(t *testing.T) {
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
				if r.Method == http.MethodPut && r.URL.Path == "/credentials/cred-123/transfer" {
					w.WriteHeader(http.StatusOK)
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

			n8nClient, server := setupTransferTestClient(t, tt.handler)
			defer server.Close()

			r := credential.NewCredentialTransferResource()
			// Configure the resource with the mock client
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
					"id":                     tftypes.NewValue(tftypes.String, nil),
					"credential_id":          tftypes.NewValue(tftypes.String, "cred-123"),
					"destination_project_id": tftypes.NewValue(tftypes.String, "proj-456"),
					"transferred_at":         tftypes.NewValue(tftypes.String, nil),
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

func setupTransferTestClient(t *testing.T, handler http.HandlerFunc) (*client.N8nClient, *httptest.Server) {
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

func TestCredentialTransferResource_Read(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "success - read maintains state",
			wantErr: false,
		},
		{
			name:    "error - invalid state",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := credential.NewCredentialTransferResource()
			ctx := context.Background()

			schemaResp := resource.SchemaResponse{}
			r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

			switch tt.name {
			case "success - read maintains state":
				// Build valid state
				stateRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":                     tftypes.NewValue(tftypes.String, "cred-123-to-proj-456"),
					"credential_id":          tftypes.NewValue(tftypes.String, "cred-123"),
					"destination_project_id": tftypes.NewValue(tftypes.String, "proj-456"),
					"transferred_at":         tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
				})

				state := tfsdk.State{
					Schema: schemaResp.Schema,
					Raw:    stateRaw,
				}

				req := resource.ReadRequest{
					State: state,
				}
				resp := &resource.ReadResponse{
					State: tfsdk.State{Schema: schemaResp.Schema},
				}

				r.Read(ctx, req, resp)

				assert.False(t, resp.Diagnostics.HasError())

			case "error - invalid state":
				// Create invalid state
				stateRaw := tftypes.NewValue(tftypes.String, "invalid")

				state := tfsdk.State{
					Schema: schemaResp.Schema,
					Raw:    stateRaw,
				}

				req := resource.ReadRequest{
					State: state,
				}
				resp := &resource.ReadResponse{}

				r.Read(ctx, req, resp)

				assert.True(t, resp.Diagnostics.HasError())
			}
		})
	}
}

func TestCredentialTransferResource_Update(t *testing.T) {
	t.Run("error - update not supported", func(t *testing.T) {
		r := credential.NewCredentialTransferResource()
		ctx := context.Background()

		req := resource.UpdateRequest{}
		resp := &resource.UpdateResponse{}

		r.Update(ctx, req, resp)

		assert.True(t, resp.Diagnostics.HasError())
		assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Update Not Supported")
	})
}

func TestCredentialTransferResource_Delete(t *testing.T) {
	t.Run("success - delete removes from state", func(t *testing.T) {
		r := credential.NewCredentialTransferResource()
		ctx := context.Background()

		req := resource.DeleteRequest{}
		resp := &resource.DeleteResponse{}

		// Should not panic and should not return errors
		assert.NotPanics(t, func() {
			r.Delete(ctx, req, resp)
		})

		assert.False(t, resp.Diagnostics.HasError())
	})
}

func TestCredentialTransferResource_ImportState(t *testing.T) {
	t.Run("import state passthrough", func(t *testing.T) {
		r := credential.NewCredentialTransferResource()
		ctx := context.Background()

		schemaResp := resource.SchemaResponse{}
		r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

		// Build empty state
		emptyValue := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
			"id":                     tftypes.NewValue(tftypes.String, nil),
			"credential_id":          tftypes.NewValue(tftypes.String, nil),
			"destination_project_id": tftypes.NewValue(tftypes.String, nil),
			"transferred_at":         tftypes.NewValue(tftypes.String, nil),
		})

		req := resource.ImportStateRequest{
			ID: "cred-123-to-proj-456",
		}
		resp := &resource.ImportStateResponse{
			State: tfsdk.State{
				Schema: schemaResp.Schema,
				Raw:    emptyValue,
			},
		}

		r.ImportState(ctx, req, resp)

		assert.False(t, resp.Diagnostics.HasError())
	})
}
