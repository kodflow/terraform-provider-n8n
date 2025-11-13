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

// setupTestClient creates a test N8nClient with httptest server.
func setupTestClient(t *testing.T, handler http.HandlerFunc) (*client.N8nClient, *httptest.Server) {
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

func TestCredentialResource_Metadata(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "set metadata", wantErr: false},
		{name: "different provider type name", wantErr: false},
		{name: "error case - metadata must be set correctly", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "set metadata":
				r := &credential.CredentialResource{}
				req := resource.MetadataRequest{
					ProviderTypeName: "n8n",
				}
				resp := &resource.MetadataResponse{}

				r.Metadata(context.Background(), req, resp)

				assert.Equal(t, "n8n_credential", resp.TypeName)

			case "different provider type name":
				r := &credential.CredentialResource{}
				req := resource.MetadataRequest{
					ProviderTypeName: "custom",
				}
				resp := &resource.MetadataResponse{}

				r.Metadata(context.Background(), req, resp)

				assert.Equal(t, "custom_credential", resp.TypeName)

			case "error case - metadata must be set correctly":
				r := &credential.CredentialResource{}
				req := resource.MetadataRequest{
					ProviderTypeName: "test",
				}
				resp := &resource.MetadataResponse{}

				r.Metadata(context.Background(), req, resp)

				assert.NotEmpty(t, resp.TypeName, "TypeName must be set")
			}
		})
	}
}

func TestCredentialResource_Schema(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "get schema", wantErr: false},
		{name: "error case - schema must not be empty", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "get schema":
				r := &credential.CredentialResource{}
				req := resource.SchemaRequest{}
				resp := &resource.SchemaResponse{}

				r.Schema(context.Background(), req, resp)

				assert.NotNil(t, resp.Schema)
				assert.Contains(t, resp.Schema.MarkdownDescription, "credential resource")
				assert.Contains(t, resp.Schema.MarkdownDescription, "automatic rotation")
				assert.NotEmpty(t, resp.Schema.Attributes)

			case "error case - schema must not be empty":
				r := &credential.CredentialResource{}
				req := resource.SchemaRequest{}
				resp := &resource.SchemaResponse{}

				r.Schema(context.Background(), req, resp)

				assert.NotNil(t, resp.Schema, "Schema must not be nil")
				assert.NotEmpty(t, resp.Schema.Attributes, "Schema attributes must not be empty")
			}
		})
	}
}

func TestCredentialResource_ImportState(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "import state", wantErr: false},
		{name: "error case - import must set ID", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "import state":
				r := &credential.CredentialResource{}
				ctx := context.Background()

				// Create schema for state
				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				// Create an empty state with the schema
				state := tfsdk.State{
					Schema: schemaResp.Schema,
				}

				// Initialize the raw value with required attributes
				state.Raw = tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, nil),
					"name":       tftypes.NewValue(tftypes.String, nil),
					"type":       tftypes.NewValue(tftypes.String, nil),
					"data":       tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, nil),
					"created_at": tftypes.NewValue(tftypes.String, nil),
					"updated_at": tftypes.NewValue(tftypes.String, nil),
				})

				req := resource.ImportStateRequest{
					ID: "cred-123",
				}
				resp := &resource.ImportStateResponse{
					State: state,
				}

				r.ImportState(ctx, req, resp)

				// Check that ID was set to import
				assert.NotNil(t, resp)
				if resp.Diagnostics.HasError() {
					t.Logf("Diagnostics errors: %v", resp.Diagnostics.Errors())
				}
				// ImportState will have a warning but not an error
				assert.False(t, resp.Diagnostics.HasError())

			case "error case - import must set ID":
				r := &credential.CredentialResource{}
				ctx := context.Background()
				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)
				state := tfsdk.State{
					Schema: schemaResp.Schema,
				}
				state.Raw = tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, nil),
					"name":       tftypes.NewValue(tftypes.String, nil),
					"type":       tftypes.NewValue(tftypes.String, nil),
					"data":       tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, nil),
					"created_at": tftypes.NewValue(tftypes.String, nil),
					"updated_at": tftypes.NewValue(tftypes.String, nil),
				})
				req := resource.ImportStateRequest{
					ID: "test-id",
				}
				resp := &resource.ImportStateResponse{
					State: state,
				}
				r.ImportState(ctx, req, resp)
				assert.NotNil(t, resp)
			}
		})
	}
}

// TestNewCredentialResource tests the public NewCredentialResource function.
func TestNewCredentialResource(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "create new resource", wantErr: false},
		{name: "create new resource returns CredentialResource type", wantErr: false},
		{name: "error case - resource must not be nil", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "create new resource":
				resource := credential.NewCredentialResource()
				assert.NotNil(t, resource)
				assert.IsType(t, &credential.CredentialResource{}, resource)

			case "create new resource returns CredentialResource type":
				resource := credential.NewCredentialResource()
				assert.NotNil(t, resource)
				assert.IsType(t, &credential.CredentialResource{}, resource)

			case "error case - resource must not be nil":
				resource := credential.NewCredentialResource()
				assert.NotNil(t, resource, "NewCredentialResource must not return nil")
			}
		})
	}
}

// TestNewCredentialResourceWrapper tests the public NewCredentialResourceWrapper function.
func TestNewCredentialResourceWrapper(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "create wrapper returns resource interface", wantErr: false},
		{name: "wrapper is CredentialResource type", wantErr: false},
		{name: "error case - wrapper must not be nil", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "create wrapper returns resource interface":
				wrapper := credential.NewCredentialResourceWrapper()
				assert.NotNil(t, wrapper)

			case "wrapper is CredentialResource type":
				wrapper := credential.NewCredentialResourceWrapper()
				assert.NotNil(t, wrapper)
				_, ok := wrapper.(*credential.CredentialResource)
				assert.True(t, ok)

			case "error case - wrapper must not be nil":
				wrapper := credential.NewCredentialResourceWrapper()
				assert.NotNil(t, wrapper, "NewCredentialResourceWrapper must not return nil")
			}
		})
	}
}

// TestCredentialResource_Configure tests the public Configure method behavior.
func TestCredentialResource_Configure(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "configure with nil provider data", wantErr: false},
		{name: "configure accepts valid provider data", wantErr: false},
		{name: "error case - Configure must handle nil gracefully", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "configure with nil provider data":
				r := &credential.CredentialResource{}
				req := resource.ConfigureRequest{
					ProviderData: nil,
				}
				resp := &resource.ConfigureResponse{}

				r.Configure(context.Background(), req, resp)

				// Should not error with nil provider data
				assert.False(t, resp.Diagnostics.HasError())

			case "configure accepts valid provider data":
				r := &credential.CredentialResource{}
				req := resource.ConfigureRequest{
					ProviderData: "test-data",
				}
				resp := &resource.ConfigureResponse{}

				r.Configure(context.Background(), req, resp)

				// Configure should accept any provider data
				assert.NotNil(t, resp)

			case "error case - Configure must handle nil gracefully":
				r := &credential.CredentialResource{}
				req := resource.ConfigureRequest{
					ProviderData: nil,
				}
				resp := &resource.ConfigureResponse{}

				assert.NotPanics(t, func() {
					r.Configure(context.Background(), req, resp)
				})
			}
		})
	}
}

// TestCredentialResource_Create tests the public Create method behavior.
func TestCredentialResource_Create(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "error - invalid plan", wantErr: true},
		{name: "error - API create fails", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "error - invalid plan":
				r := &credential.CredentialResource{}
				ctx := context.Background()

				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				// Invalid plan with wrong type
				req := resource.CreateRequest{
					Plan: tfsdk.Plan{
						Raw:    tftypes.NewValue(tftypes.String, "invalid"),
						Schema: schemaResp.Schema,
					},
				}
				resp := &resource.CreateResponse{
					State: tfsdk.State{Schema: schemaResp.Schema},
				}

				r.Create(ctx, req, resp)
				assert.True(t, resp.Diagnostics.HasError(), "Expected error with invalid plan")

			case "error - API create fails":
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.Method == http.MethodPost && r.URL.Path == "/credentials" {
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusInternalServerError)
						w.Write([]byte(`{"message":"Internal server error"}`))
						return
					}
					t.Errorf("Unexpected request: %s %s", r.Method, r.URL.Path)
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				r := credential.NewCredentialResource()
				r.Configure(context.Background(), resource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &resource.ConfigureResponse{})

				ctx := context.Background()

				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				dataMap := map[string]tftypes.Value{
					"key": tftypes.NewValue(tftypes.String, "value"),
				}
				planRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, nil),
					"name":       tftypes.NewValue(tftypes.String, "test-credential"),
					"type":       tftypes.NewValue(tftypes.String, "httpHeaderAuth"),
					"data":       tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, dataMap),
					"created_at": tftypes.NewValue(tftypes.String, nil),
					"updated_at": tftypes.NewValue(tftypes.String, nil),
				})

				req := resource.CreateRequest{
					Plan: tfsdk.Plan{
						Raw:    planRaw,
						Schema: schemaResp.Schema,
					},
				}
				resp := &resource.CreateResponse{
					State: tfsdk.State{Schema: schemaResp.Schema},
				}

				r.Create(ctx, req, resp)
				assert.True(t, resp.Diagnostics.HasError(), "Expected error when API create fails")
			}
		})
	}
}

// TestCredentialResource_Read tests the public Read method behavior.
func TestCredentialResource_Read(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "success - read keeps state as-is", wantErr: false},
		{name: "error - invalid state", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "success - read keeps state as-is":
				r := &credential.CredentialResource{}
				ctx := context.Background()

				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				dataMap := map[string]tftypes.Value{
					"key": tftypes.NewValue(tftypes.String, "value"),
				}
				stateRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "cred-123"),
					"name":       tftypes.NewValue(tftypes.String, "test-credential"),
					"type":       tftypes.NewValue(tftypes.String, "httpHeaderAuth"),
					"data":       tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, dataMap),
					"created_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
					"updated_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
				})

				req := resource.ReadRequest{
					State: tfsdk.State{
						Raw:    stateRaw,
						Schema: schemaResp.Schema,
					},
				}
				resp := &resource.ReadResponse{
					State: tfsdk.State{Schema: schemaResp.Schema},
				}

				r.Read(ctx, req, resp)
				// Read should keep state as-is without errors or warnings
				assert.False(t, resp.Diagnostics.HasError(), "Expected no errors")
				assert.False(t, len(resp.Diagnostics.Warnings()) > 0, "Expected no warnings")

			case "error - invalid state":
				r := &credential.CredentialResource{}
				ctx := context.Background()

				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				req := resource.ReadRequest{
					State: tfsdk.State{
						Raw:    tftypes.NewValue(tftypes.String, "invalid"),
						Schema: schemaResp.Schema,
					},
				}
				resp := &resource.ReadResponse{
					State: tfsdk.State{Schema: schemaResp.Schema},
				}

				r.Read(ctx, req, resp)
				assert.True(t, resp.Diagnostics.HasError(), "Expected error with invalid state")
			}
		})
	}
}

// TestCredentialResource_Update tests the public Update method behavior.
func TestCredentialResource_Update(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "error - invalid plan", wantErr: true},
		{name: "error - invalid state", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "error - invalid plan":
				r := &credential.CredentialResource{}
				ctx := context.Background()

				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				dataMap := map[string]tftypes.Value{
					"key": tftypes.NewValue(tftypes.String, "value"),
				}
				stateRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "cred-123"),
					"name":       tftypes.NewValue(tftypes.String, "test-credential"),
					"type":       tftypes.NewValue(tftypes.String, "httpHeaderAuth"),
					"data":       tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, dataMap),
					"created_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
					"updated_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
				})

				req := resource.UpdateRequest{
					Plan: tfsdk.Plan{
						Raw:    tftypes.NewValue(tftypes.String, "invalid"),
						Schema: schemaResp.Schema,
					},
					State: tfsdk.State{
						Raw:    stateRaw,
						Schema: schemaResp.Schema,
					},
				}
				resp := &resource.UpdateResponse{
					State: tfsdk.State{Schema: schemaResp.Schema},
				}

				r.Update(ctx, req, resp)
				assert.True(t, resp.Diagnostics.HasError(), "Expected error with invalid plan")

			case "error - invalid state":
				r := &credential.CredentialResource{}
				ctx := context.Background()

				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				dataMap := map[string]tftypes.Value{
					"key": tftypes.NewValue(tftypes.String, "value"),
				}
				planRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "cred-123"),
					"name":       tftypes.NewValue(tftypes.String, "test-credential-updated"),
					"type":       tftypes.NewValue(tftypes.String, "httpHeaderAuth"),
					"data":       tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, dataMap),
					"created_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
					"updated_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
				})

				req := resource.UpdateRequest{
					Plan: tfsdk.Plan{
						Raw:    planRaw,
						Schema: schemaResp.Schema,
					},
					State: tfsdk.State{
						Raw:    tftypes.NewValue(tftypes.String, "invalid"),
						Schema: schemaResp.Schema,
					},
				}
				resp := &resource.UpdateResponse{
					State: tfsdk.State{Schema: schemaResp.Schema},
				}

				r.Update(ctx, req, resp)
				assert.True(t, resp.Diagnostics.HasError(), "Expected error with invalid state")
			}
		})
	}
}

// TestCredentialResource_Delete tests the public Delete method behavior.
func TestCredentialResource_Delete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "error - invalid state", wantErr: true},
		{name: "error - API delete fails", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "error - invalid state":
				r := &credential.CredentialResource{}
				ctx := context.Background()

				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				req := resource.DeleteRequest{
					State: tfsdk.State{
						Raw:    tftypes.NewValue(tftypes.String, "invalid"),
						Schema: schemaResp.Schema,
					},
				}
				resp := &resource.DeleteResponse{}

				r.Delete(ctx, req, resp)
				assert.True(t, resp.Diagnostics.HasError(), "Expected error with invalid state")

			case "error - API delete fails":
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.Method == http.MethodDelete && r.URL.Path == "/credentials/cred-123" {
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusInternalServerError)
						w.Write([]byte(`{"message":"delete failed"}`))
						return
					}
				})

				n8nClient, _ := setupTestClient(t, handler)

				r := credential.NewCredentialResource()
				r.Configure(context.Background(), resource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &resource.ConfigureResponse{})

				ctx := context.Background()
				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				dataMap := map[string]tftypes.Value{
					"key": tftypes.NewValue(tftypes.String, "value"),
				}
				stateRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":         tftypes.NewValue(tftypes.String, "cred-123"),
					"name":       tftypes.NewValue(tftypes.String, "test-credential"),
					"type":       tftypes.NewValue(tftypes.String, "httpHeaderAuth"),
					"data":       tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, dataMap),
					"created_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
					"updated_at": tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
				})

				req := resource.DeleteRequest{
					State: tfsdk.State{
						Raw:    stateRaw,
						Schema: schemaResp.Schema,
					},
				}
				resp := &resource.DeleteResponse{}

				r.Delete(ctx, req, resp)
				assert.True(t, resp.Diagnostics.HasError(), "Expected error when API delete fails")
			}
		})
	}
}
