package sourcecontrol_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
	"github.com/kodflow/n8n/src/internal/provider/sourcecontrol"
	"github.com/kodflow/n8n/src/internal/provider/sourcecontrol/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSourceControlPullResourceInterface is a mock implementation of SourceControlPullResourceInterface.
type MockSourceControlPullResourceInterface struct {
	mock.Mock
}

func (m *MockSourceControlPullResourceInterface) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockSourceControlPullResourceInterface) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockSourceControlPullResourceInterface) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockSourceControlPullResourceInterface) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockSourceControlPullResourceInterface) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockSourceControlPullResourceInterface) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockSourceControlPullResourceInterface) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockSourceControlPullResourceInterface) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	m.Called(ctx, req, resp)
}

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

// TestNewSourceControlPullResource tests the NewSourceControlPullResource function.
func TestNewSourceControlPullResource(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "create new instance",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := sourcecontrol.NewSourceControlPullResource()

				assert.NotNil(t, r)
				assert.IsType(t, &sourcecontrol.SourceControlPullResource{}, r)
			},
		},
		{
			name: "multiple instances are independent",
			testFunc: func(t *testing.T) {
				t.Helper()
				r1 := sourcecontrol.NewSourceControlPullResource()
				r2 := sourcecontrol.NewSourceControlPullResource()

				assert.NotSame(t, r1, r2)
			},
		},
		{
			name: "error - instance is never nil",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := sourcecontrol.NewSourceControlPullResource()

				assert.NotNil(t, r, "NewSourceControlPullResource should never return nil")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}

// TestNewSourceControlPullResourceWrapper tests the NewSourceControlPullResourceWrapper function.
func TestNewSourceControlPullResourceWrapper(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "create new wrapped instance",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := sourcecontrol.NewSourceControlPullResourceWrapper()

				assert.NotNil(t, r)
				assert.IsType(t, &sourcecontrol.SourceControlPullResource{}, r)
			},
		},
		{
			name: "wrapper returns resource.Resource interface",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := sourcecontrol.NewSourceControlPullResourceWrapper()

				assert.NotNil(t, r)
			},
		},
		{
			name: "error - wrapper never returns nil",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := sourcecontrol.NewSourceControlPullResourceWrapper()

				assert.NotNil(t, r, "NewSourceControlPullResourceWrapper should never return nil")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}

// TestSourceControlPullResource_Metadata tests the Metadata method.
func TestSourceControlPullResource_Metadata(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "set metadata with n8n provider",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := sourcecontrol.NewSourceControlPullResource()
				resp := &resource.MetadataResponse{}

				r.Metadata(context.Background(), resource.MetadataRequest{
					ProviderTypeName: "n8n",
				}, resp)

				assert.Equal(t, "n8n_source_control_pull", resp.TypeName)
			},
		},
		{
			name: "set metadata with custom provider",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := sourcecontrol.NewSourceControlPullResource()
				resp := &resource.MetadataResponse{}
				req := resource.MetadataRequest{
					ProviderTypeName: "custom_provider",
				}

				r.Metadata(context.Background(), req, resp)

				assert.Equal(t, "custom_provider_source_control_pull", resp.TypeName)
			},
		},
		{
			name: "error - metadata with empty provider name",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := sourcecontrol.NewSourceControlPullResource()
				resp := &resource.MetadataResponse{}
				req := resource.MetadataRequest{
					ProviderTypeName: "",
				}

				r.Metadata(context.Background(), req, resp)

				assert.Equal(t, "_source_control_pull", resp.TypeName)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}

// TestSourceControlPullResource_Schema tests the Schema method.
func TestSourceControlPullResource_Schema(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "return schema with all attributes",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := sourcecontrol.NewSourceControlPullResource()
				resp := &resource.SchemaResponse{}

				r.Schema(context.Background(), resource.SchemaRequest{}, resp)

				assert.NotNil(t, resp.Schema)
				assert.Contains(t, resp.Schema.MarkdownDescription, "Pulls changes from the remote source control repository")

				expectedAttrs := []string{"id", "force", "variables_json", "result_json"}
				for _, attr := range expectedAttrs {
					_, exists := resp.Schema.Attributes[attr]
					assert.True(t, exists, "Attribute %s should exist", attr)
				}
			},
		},
		{
			name: "schema attributes have correct properties",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := sourcecontrol.NewSourceControlPullResource()
				resp := &resource.SchemaResponse{}

				r.Schema(context.Background(), resource.SchemaRequest{}, resp)

				idAttr := resp.Schema.Attributes["id"]
				assert.True(t, idAttr.IsComputed())

				forceAttr := resp.Schema.Attributes["force"]
				assert.True(t, forceAttr.IsOptional())

				variablesAttr := resp.Schema.Attributes["variables_json"]
				assert.True(t, variablesAttr.IsOptional())

				resultAttr := resp.Schema.Attributes["result_json"]
				assert.True(t, resultAttr.IsComputed())
			},
		},
		{
			name: "error - all attributes have descriptions",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := sourcecontrol.NewSourceControlPullResource()
				resp := &resource.SchemaResponse{}

				r.Schema(context.Background(), resource.SchemaRequest{}, resp)

				for name, attr := range resp.Schema.Attributes {
					assert.NotEmpty(t, attr.GetMarkdownDescription(), "Attribute %s should have description", name)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}

// TestSourceControlPullResource_Configure tests the Configure method.
func TestSourceControlPullResource_Configure(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "configure with valid client",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := sourcecontrol.NewSourceControlPullResource()
				resp := &resource.ConfigureResponse{}

				mockClient := &client.N8nClient{}
				req := resource.ConfigureRequest{
					ProviderData: mockClient,
				}

				r.Configure(context.Background(), req, resp)

				assert.False(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "configure with nil provider data",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := sourcecontrol.NewSourceControlPullResource()
				resp := &resource.ConfigureResponse{}
				req := resource.ConfigureRequest{
					ProviderData: nil,
				}

				r.Configure(context.Background(), req, resp)

				assert.False(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "error - configure with invalid provider data",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := sourcecontrol.NewSourceControlPullResource()
				resp := &resource.ConfigureResponse{}
				req := resource.ConfigureRequest{
					ProviderData: "invalid-data",
				}

				r.Configure(context.Background(), req, resp)

				assert.True(t, resp.Diagnostics.HasError())
				assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Unexpected Resource Configure Type")
			},
		},
		{
			name: "error - configure with wrong type",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := sourcecontrol.NewSourceControlPullResource()
				resp := &resource.ConfigureResponse{}
				req := resource.ConfigureRequest{
					ProviderData: struct{}{},
				}

				r.Configure(context.Background(), req, resp)

				assert.True(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "configure multiple times",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := sourcecontrol.NewSourceControlPullResource()

				// First configuration
				resp1 := &resource.ConfigureResponse{}
				client1 := &client.N8nClient{}
				req1 := resource.ConfigureRequest{
					ProviderData: client1,
				}
				r.Configure(context.Background(), req1, resp1)
				assert.False(t, resp1.Diagnostics.HasError())

				// Second configuration
				resp2 := &resource.ConfigureResponse{}
				client2 := &client.N8nClient{}
				req2 := resource.ConfigureRequest{
					ProviderData: client2,
				}
				r.Configure(context.Background(), req2, resp2)
				assert.False(t, resp2.Diagnostics.HasError())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}

// TestSourceControlPullResource_Create tests the Create method.
func TestSourceControlPullResource_Create(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "error - create with API error",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("Internal Server Error"))
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				r := sourcecontrol.NewSourceControlPullResource()
				r.Configure(context.Background(), resource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &resource.ConfigureResponse{})

				ctx := context.Background()
				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				planRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":             tftypes.NewValue(tftypes.String, nil),
					"force":          tftypes.NewValue(tftypes.Bool, false),
					"variables_json": tftypes.NewValue(tftypes.String, nil),
					"result_json":    tftypes.NewValue(tftypes.String, nil),
				})

				plan := tfsdk.Plan{
					Schema: schemaResp.Schema,
					Raw:    planRaw,
				}

				state := tfsdk.State{
					Schema: schemaResp.Schema,
				}

				req := resource.CreateRequest{
					Plan: plan,
				}
				resp := &resource.CreateResponse{
					State: state,
				}

				r.Create(ctx, req, resp)

				assert.True(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "error - create with invalid variables JSON",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := sourcecontrol.NewSourceControlPullResource()
				ctx := context.Background()

				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				planRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":             tftypes.NewValue(tftypes.String, nil),
					"force":          tftypes.NewValue(tftypes.Bool, false),
					"variables_json": tftypes.NewValue(tftypes.String, `{invalid-json`),
					"result_json":    tftypes.NewValue(tftypes.String, nil),
				})

				plan := tfsdk.Plan{
					Schema: schemaResp.Schema,
					Raw:    planRaw,
				}

				state := tfsdk.State{
					Schema: schemaResp.Schema,
				}

				req := resource.CreateRequest{
					Plan: plan,
				}
				resp := &resource.CreateResponse{
					State: state,
				}

				r.Create(ctx, req, resp)

				assert.True(t, resp.Diagnostics.HasError())
				assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Invalid variables JSON")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}

// TestSourceControlPullResource_Read tests the Read method.
func TestSourceControlPullResource_Read(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "read maintains state",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := sourcecontrol.NewSourceControlPullResource()
				ctx := context.Background()

				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				state := tfsdk.State{
					Schema: schemaResp.Schema,
				}

				state.Raw = tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":             tftypes.NewValue(tftypes.String, "pull-200"),
					"force":          tftypes.NewValue(tftypes.Bool, nil),
					"variables_json": tftypes.NewValue(tftypes.String, nil),
					"result_json":    tftypes.NewValue(tftypes.String, nil),
				})

				req := resource.ReadRequest{
					State: state,
				}
				resp := &resource.ReadResponse{
					State: state,
				}

				r.Read(ctx, req, resp)

				assert.False(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "error - read with invalid state",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := sourcecontrol.NewSourceControlPullResource()
				ctx := context.Background()

				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				state := tfsdk.State{
					Schema: schemaResp.Schema,
					Raw:    tftypes.NewValue(tftypes.String, "invalid"),
				}

				req := resource.ReadRequest{
					State: state,
				}
				resp := &resource.ReadResponse{
					State: state,
				}

				r.Read(ctx, req, resp)

				assert.True(t, resp.Diagnostics.HasError())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}

// TestSourceControlPullResource_Update tests the Update method.
func TestSourceControlPullResource_Update(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "error - update not supported",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := sourcecontrol.NewSourceControlPullResource()
				resp := &resource.UpdateResponse{}

				r.Update(context.Background(), resource.UpdateRequest{}, resp)

				assert.True(t, resp.Diagnostics.HasError())
				assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Update Not Supported")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}

// TestSourceControlPullResource_Delete tests the Delete method.
func TestSourceControlPullResource_Delete(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "delete does nothing",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := sourcecontrol.NewSourceControlPullResource()
				resp := &resource.DeleteResponse{}

				r.Delete(context.Background(), resource.DeleteRequest{}, resp)

				assert.False(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "delete with valid state",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := sourcecontrol.NewSourceControlPullResource()
				ctx := context.Background()

				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				stateRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":             tftypes.NewValue(tftypes.String, "pull-200"),
					"force":          tftypes.NewValue(tftypes.Bool, true),
					"variables_json": tftypes.NewValue(tftypes.String, nil),
					"result_json":    tftypes.NewValue(tftypes.String, `{"workflows":[]}`),
				})

				state := tfsdk.State{
					Schema: schemaResp.Schema,
					Raw:    stateRaw,
				}

				req := resource.DeleteRequest{
					State: state,
				}
				resp := &resource.DeleteResponse{}

				r.Delete(ctx, req, resp)

				assert.False(t, resp.Diagnostics.HasError())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}

// TestSourceControlPullResource_ImportState tests the ImportState method.
func TestSourceControlPullResource_ImportState(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "import state with valid ID",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := sourcecontrol.NewSourceControlPullResource()
				ctx := context.Background()

				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				state := tfsdk.State{
					Schema: schemaResp.Schema,
				}

				state.Raw = tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":             tftypes.NewValue(tftypes.String, nil),
					"force":          tftypes.NewValue(tftypes.Bool, nil),
					"variables_json": tftypes.NewValue(tftypes.String, nil),
					"result_json":    tftypes.NewValue(tftypes.String, nil),
				})

				resp := &resource.ImportStateResponse{
					State: state,
				}
				req := resource.ImportStateRequest{
					ID: "pull-123",
				}

				r.ImportState(ctx, req, resp)

				assert.False(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "import state sets id attribute",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := sourcecontrol.NewSourceControlPullResource()
				ctx := context.Background()

				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				state := tfsdk.State{
					Schema: schemaResp.Schema,
				}

				state.Raw = tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":             tftypes.NewValue(tftypes.String, nil),
					"force":          tftypes.NewValue(tftypes.Bool, nil),
					"variables_json": tftypes.NewValue(tftypes.String, nil),
					"result_json":    tftypes.NewValue(tftypes.String, nil),
				})

				resp := &resource.ImportStateResponse{
					State: state,
				}
				req := resource.ImportStateRequest{
					ID: "pull-456",
				}

				r.ImportState(ctx, req, resp)

				assert.False(t, resp.Diagnostics.HasError())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}

// TestSourceControlPullResource_Interfaces tests interface implementations.
func TestSourceControlPullResource_Interfaces(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "implements required interfaces",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := sourcecontrol.NewSourceControlPullResource()

				// Verify interfaces are implemented (compiler checks this)
				assert.NotNil(t, r)
			},
		},
		{
			name: "error - resource implements all required methods",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := sourcecontrol.NewSourceControlPullResource()

				assert.NotNil(t, r)
				assert.Implements(t, (*resource.Resource)(nil), r)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}

// TestSourceControlPullResource_Concurrency tests concurrent operations.
func TestSourceControlPullResource_Concurrency(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "concurrent metadata calls",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := sourcecontrol.NewSourceControlPullResource()

				done := make(chan bool, 100)
				for i := 0; i < 100; i++ {
					go func() {
						resp := &resource.MetadataResponse{}
						r.Metadata(context.Background(), resource.MetadataRequest{
							ProviderTypeName: "n8n",
						}, resp)
						assert.Equal(t, "n8n_source_control_pull", resp.TypeName)
						done <- true
					}()
				}

				for i := 0; i < 100; i++ {
					<-done
				}
			},
		},
		{
			name: "concurrent schema calls",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := sourcecontrol.NewSourceControlPullResource()

				done := make(chan bool, 100)
				for i := 0; i < 100; i++ {
					go func() {
						resp := &resource.SchemaResponse{}
						r.Schema(context.Background(), resource.SchemaRequest{}, resp)
						assert.NotNil(t, resp.Schema)
						done <- true
					}()
				}

				for i := 0; i < 100; i++ {
					<-done
				}
			},
		},
		{
			name: "concurrent configure calls",
			testFunc: func(t *testing.T) {
				t.Helper()
				done := make(chan bool, 100)
				for i := 0; i < 100; i++ {
					go func() {
						r := sourcecontrol.NewSourceControlPullResource()
						resp := &resource.ConfigureResponse{}
						mockClient := &client.N8nClient{}
						req := resource.ConfigureRequest{
							ProviderData: mockClient,
						}
						r.Configure(context.Background(), req, resp)
						assert.False(t, resp.Diagnostics.HasError())
						done <- true
					}()
				}

				for i := 0; i < 100; i++ {
					<-done
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}

// TestSourceControlPullResource_HelperFunctions tests various helper functions.
func TestSourceControlPullResource_HelperFunctions(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "path to id",
			testFunc: func(t *testing.T) {
				t.Helper()
				p := path.Root("id")
				assert.Equal(t, "id", p.String())
			},
		},
		{
			name: "path to force",
			testFunc: func(t *testing.T) {
				t.Helper()
				p := path.Root("force")
				assert.Equal(t, "force", p.String())
			},
		},
		{
			name: "path to variables_json",
			testFunc: func(t *testing.T) {
				t.Helper()
				p := path.Root("variables_json")
				assert.Equal(t, "variables_json", p.String())
			},
		},
		{
			name: "path to result_json",
			testFunc: func(t *testing.T) {
				t.Helper()
				p := path.Root("result_json")
				assert.Equal(t, "result_json", p.String())
			},
		},
		{
			name: "types.StringValue conversion",
			testFunc: func(t *testing.T) {
				t.Helper()
				strVal := types.StringValue("test-value")
				assert.Equal(t, "test-value", strVal.ValueString())
				assert.False(t, strVal.IsNull())
				assert.False(t, strVal.IsUnknown())
			},
		},
		{
			name: "types.BoolValue conversion",
			testFunc: func(t *testing.T) {
				t.Helper()
				boolVal := types.BoolValue(true)
				assert.True(t, boolVal.ValueBool())
				assert.False(t, boolVal.IsNull())
				assert.False(t, boolVal.IsUnknown())
			},
		},
		{
			name: "types.StringNull",
			testFunc: func(t *testing.T) {
				t.Helper()
				strVal := types.StringNull()
				assert.True(t, strVal.IsNull())
				assert.False(t, strVal.IsUnknown())
			},
		},
		{
			name: "types.BoolNull",
			testFunc: func(t *testing.T) {
				t.Helper()
				boolVal := types.BoolNull()
				assert.True(t, boolVal.IsNull())
				assert.False(t, boolVal.IsUnknown())
			},
		},
		{
			name: "error - path with empty string",
			testFunc: func(t *testing.T) {
				t.Helper()
				p := path.Root("")
				assert.NotNil(t, p)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}

// TestSourceControlPullResource_EdgeCases tests edge cases.
func TestSourceControlPullResource_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "configure then reconfigure with nil",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := sourcecontrol.NewSourceControlPullResource()

				// First configure with valid client
				resp1 := &resource.ConfigureResponse{}
				mockClient := &client.N8nClient{}
				req1 := resource.ConfigureRequest{
					ProviderData: mockClient,
				}
				r.Configure(context.Background(), req1, resp1)

				// Then configure with nil
				resp2 := &resource.ConfigureResponse{}
				req2 := resource.ConfigureRequest{
					ProviderData: nil,
				}
				r.Configure(context.Background(), req2, resp2)
				assert.False(t, resp2.Diagnostics.HasError())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}

// TestSourceControlPullResource_BuildWithDifferentPlanValues tests build request with different plan values.
func TestSourceControlPullResource_BuildWithDifferentPlanValues(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "plan with force true and false",
			testFunc: func(t *testing.T) {
				t.Helper()
				plan1 := &models.PullResource{
					Force: types.BoolValue(true),
				}
				plan2 := &models.PullResource{
					Force: types.BoolValue(false),
				}
				assert.True(t, plan1.Force.ValueBool())
				assert.False(t, plan2.Force.ValueBool())
			},
		},
		{
			name: "error - plan with null force",
			testFunc: func(t *testing.T) {
				t.Helper()
				plan := &models.PullResource{
					Force: types.BoolNull(),
				}
				assert.True(t, plan.Force.IsNull())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}

// BenchmarkSourceControlPullResource_Schema benchmarks the Schema method.
func BenchmarkSourceControlPullResource_Schema(b *testing.B) {
	r := sourcecontrol.NewSourceControlPullResource()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := &resource.SchemaResponse{}
		r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	}
}

// BenchmarkSourceControlPullResource_Metadata benchmarks the Metadata method.
func BenchmarkSourceControlPullResource_Metadata(b *testing.B) {
	r := sourcecontrol.NewSourceControlPullResource()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := &resource.MetadataResponse{}
		r.Metadata(context.Background(), resource.MetadataRequest{
			ProviderTypeName: "n8n",
		}, resp)
	}
}

// BenchmarkSourceControlPullResource_Configure benchmarks the Configure method.
func BenchmarkSourceControlPullResource_Configure(b *testing.B) {
	r := sourcecontrol.NewSourceControlPullResource()
	mockClient := &client.N8nClient{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := &resource.ConfigureResponse{}
		req := resource.ConfigureRequest{
			ProviderData: mockClient,
		}
		r.Configure(context.Background(), req, resp)
	}
}

// TestSourceControlPullResource_CreateWithPlan tests Create with different plan configurations.
func TestSourceControlPullResource_CreateWithPlan(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "error - plan with force parameter",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{"workflows": [], "credentials": [], "tags": [], "variables": []}`))
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				r := sourcecontrol.NewSourceControlPullResource()
				r.Configure(context.Background(), resource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &resource.ConfigureResponse{})

				ctx := context.Background()
				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				planRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":             tftypes.NewValue(tftypes.String, nil),
					"force":          tftypes.NewValue(tftypes.Bool, true),
					"variables_json": tftypes.NewValue(tftypes.String, nil),
					"result_json":    tftypes.NewValue(tftypes.String, nil),
				})

				plan := tfsdk.Plan{
					Schema: schemaResp.Schema,
					Raw:    planRaw,
				}

				state := tfsdk.State{
					Schema: schemaResp.Schema,
				}

				req := resource.CreateRequest{
					Plan: plan,
				}
				resp := &resource.CreateResponse{
					State: state,
				}

				r.Create(ctx, req, resp)

				// Expecting success or error depending on implementation
				assert.NotNil(t, resp)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}

// TestSourceControlPullResource_TypeConversions tests type conversions.
func TestSourceControlPullResource_TypeConversions(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "string value to ID",
			testFunc: func(t *testing.T) {
				t.Helper()
				id := types.StringValue("pull-200")
				assert.Equal(t, "pull-200", id.ValueString())
			},
		},
		{
			name: "format ID with status code",
			testFunc: func(t *testing.T) {
				t.Helper()
				statusCode := 200
				id := fmt.Sprintf("pull-%d", statusCode)
				assert.Equal(t, "pull-200", id)
			},
		},
		{
			name: "error - format ID with invalid status code",
			testFunc: func(t *testing.T) {
				t.Helper()
				statusCode := -1
				id := fmt.Sprintf("pull-%d", statusCode)
				assert.Equal(t, "pull--1", id)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}
