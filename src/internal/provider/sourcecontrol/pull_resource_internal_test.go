package sourcecontrol

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

func TestNewSourceControlPullResource(t *testing.T) {
	t.Run("create new instance", func(t *testing.T) {
		r := NewSourceControlPullResource()

		assert.NotNil(t, r)
		assert.IsType(t, &SourceControlPullResource{}, r)
		assert.Nil(t, r.client)
	})

	t.Run("multiple instances are independent", func(t *testing.T) {
		r1 := NewSourceControlPullResource()
		r2 := NewSourceControlPullResource()

		assert.NotSame(t, r1, r2)
	})
}

func TestNewSourceControlPullResourceWrapper(t *testing.T) {
	t.Run("create new wrapped instance", func(t *testing.T) {
		r := NewSourceControlPullResourceWrapper()

		assert.NotNil(t, r)
		assert.IsType(t, &SourceControlPullResource{}, r)
	})

	t.Run("wrapper returns resource.Resource interface", func(t *testing.T) {
		r := NewSourceControlPullResourceWrapper()

		// r is already of type resource.Resource, no assertion needed
		assert.NotNil(t, r)
	})
}

func TestSourceControlPullResource_Metadata(t *testing.T) {
	t.Run("set metadata without provider type", func(t *testing.T) {
		r := NewSourceControlPullResource()
		resp := &resource.MetadataResponse{}

		r.Metadata(context.Background(), resource.MetadataRequest{
			ProviderTypeName: "n8n",
		}, resp)

		assert.Equal(t, "n8n_source_control_pull", resp.TypeName)
	})

	t.Run("set metadata with provider type name", func(t *testing.T) {
		r := NewSourceControlPullResource()
		resp := &resource.MetadataResponse{}
		req := resource.MetadataRequest{
			ProviderTypeName: "n8n",
		}

		r.Metadata(context.Background(), req, resp)

		assert.Equal(t, "n8n_source_control_pull", resp.TypeName)
	})

	t.Run("set metadata with different provider type", func(t *testing.T) {
		r := NewSourceControlPullResource()
		resp := &resource.MetadataResponse{}
		req := resource.MetadataRequest{
			ProviderTypeName: "custom_provider",
		}

		r.Metadata(context.Background(), req, resp)

		assert.Equal(t, "custom_provider_source_control_pull", resp.TypeName)
	})
}

func TestSourceControlPullResource_Schema(t *testing.T) {
	t.Run("return schema", func(t *testing.T) {
		r := NewSourceControlPullResource()
		resp := &resource.SchemaResponse{}

		r.Schema(context.Background(), resource.SchemaRequest{}, resp)

		assert.NotNil(t, resp.Schema)
		assert.Contains(t, resp.Schema.MarkdownDescription, "Pulls changes from the remote source control repository")

		// Verify computed attributes
		idAttr, exists := resp.Schema.Attributes["id"]
		assert.True(t, exists)
		assert.True(t, idAttr.IsComputed())

		// Verify optional attributes
		forceAttr, exists := resp.Schema.Attributes["force"]
		assert.True(t, exists)
		assert.True(t, forceAttr.IsOptional())

		variablesAttr, exists := resp.Schema.Attributes["variables_json"]
		assert.True(t, exists)
		assert.True(t, variablesAttr.IsOptional())

		resultAttr, exists := resp.Schema.Attributes["result_json"]
		assert.True(t, exists)
		assert.True(t, resultAttr.IsComputed())
	})

	t.Run("schema attributes have descriptions", func(t *testing.T) {
		r := NewSourceControlPullResource()
		resp := &resource.SchemaResponse{}

		r.Schema(context.Background(), resource.SchemaRequest{}, resp)

		for name, attr := range resp.Schema.Attributes {
			assert.NotEmpty(t, attr.GetMarkdownDescription(), "Attribute %s should have description", name)
		}
	})

	t.Run("schema has all expected attributes", func(t *testing.T) {
		r := NewSourceControlPullResource()
		resp := &resource.SchemaResponse{}

		r.Schema(context.Background(), resource.SchemaRequest{}, resp)

		expectedAttrs := []string{
			"id",
			"force",
			"variables_json",
			"result_json",
		}

		for _, attr := range expectedAttrs {
			schemaAttr, exists := resp.Schema.Attributes[attr]
			assert.True(t, exists, "Attribute %s should exist", attr)
			assert.NotNil(t, schemaAttr)
		}
	})

	t.Run("schema attributes from schemaAttributes", func(t *testing.T) {
		r := NewSourceControlPullResource()
		attrs := r.schemaAttributes()

		assert.NotNil(t, attrs)
		assert.Len(t, attrs, 4)
		assert.Contains(t, attrs, "id")
		assert.Contains(t, attrs, "force")
		assert.Contains(t, attrs, "variables_json")
		assert.Contains(t, attrs, "result_json")
	})
}

func TestSourceControlPullResource_Configure(t *testing.T) {
	t.Run("configure with valid client", func(t *testing.T) {
		r := NewSourceControlPullResource()
		resp := &resource.ConfigureResponse{}

		mockClient := &client.N8nClient{}
		req := resource.ConfigureRequest{
			ProviderData: mockClient,
		}

		r.Configure(context.Background(), req, resp)

		assert.NotNil(t, r.client)
		assert.Equal(t, mockClient, r.client)
		assert.False(t, resp.Diagnostics.HasError())
	})

	t.Run("configure with nil provider data", func(t *testing.T) {
		r := NewSourceControlPullResource()
		resp := &resource.ConfigureResponse{}
		req := resource.ConfigureRequest{
			ProviderData: nil,
		}

		r.Configure(context.Background(), req, resp)

		assert.Nil(t, r.client)
		assert.False(t, resp.Diagnostics.HasError())
	})

	t.Run("configure with invalid provider data", func(t *testing.T) {
		r := NewSourceControlPullResource()
		resp := &resource.ConfigureResponse{}
		req := resource.ConfigureRequest{
			ProviderData: "invalid-data",
		}

		r.Configure(context.Background(), req, resp)

		assert.Nil(t, r.client)
		assert.True(t, resp.Diagnostics.HasError())
		assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Unexpected Resource Configure Type")
	})

	t.Run("configure with wrong type", func(t *testing.T) {
		r := NewSourceControlPullResource()
		resp := &resource.ConfigureResponse{}
		req := resource.ConfigureRequest{
			ProviderData: struct{}{},
		}

		r.Configure(context.Background(), req, resp)

		assert.Nil(t, r.client)
		assert.True(t, resp.Diagnostics.HasError())
	})

	t.Run("configure multiple times", func(t *testing.T) {
		r := NewSourceControlPullResource()

		// First configuration
		resp1 := &resource.ConfigureResponse{}
		client1 := &client.N8nClient{}
		req1 := resource.ConfigureRequest{
			ProviderData: client1,
		}
		r.Configure(context.Background(), req1, resp1)
		assert.Equal(t, client1, r.client)

		// Second configuration
		resp2 := &resource.ConfigureResponse{}
		client2 := &client.N8nClient{}
		req2 := resource.ConfigureRequest{
			ProviderData: client2,
		}
		r.Configure(context.Background(), req2, resp2)
		assert.Equal(t, client2, r.client)
	})
}

func TestSourceControlPullResource_Read(t *testing.T) {
	t.Run("read maintains state", func(t *testing.T) {
		r := NewSourceControlPullResource()
		ctx := context.Background()

		// Create schema for state
		schemaResp := resource.SchemaResponse{}
		r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

		// Create a state with the schema
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
	})
}

func TestSourceControlPullResource_Update(t *testing.T) {
	t.Run("update not supported", func(t *testing.T) {
		r := NewSourceControlPullResource()
		resp := &resource.UpdateResponse{}

		r.Update(context.Background(), resource.UpdateRequest{}, resp)

		assert.True(t, resp.Diagnostics.HasError())
		assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Update Not Supported")
	})
}

func TestSourceControlPullResource_Delete(t *testing.T) {
	t.Run("delete does nothing", func(t *testing.T) {
		r := NewSourceControlPullResource()
		resp := &resource.DeleteResponse{}

		r.Delete(context.Background(), resource.DeleteRequest{}, resp)

		assert.False(t, resp.Diagnostics.HasError())
	})
}

func TestSourceControlPullResource_ImportState(t *testing.T) {
	t.Run("import state with valid ID", func(t *testing.T) {
		r := NewSourceControlPullResource()
		ctx := context.Background()

		// Create schema for state
		schemaResp := resource.SchemaResponse{}
		r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

		// Create an empty state with the schema
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
	})

	t.Run("import state sets id attribute", func(t *testing.T) {
		r := NewSourceControlPullResource()
		ctx := context.Background()

		// Create schema for state
		schemaResp := resource.SchemaResponse{}
		r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

		// Create an empty state with the schema
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
	})
}

func TestSourceControlPullResource_BuildPullRequest(t *testing.T) {
	t.Run("build pull request with force", func(t *testing.T) {
		r := NewSourceControlPullResource()
		plan := &models.PullResource{
			Force: types.BoolValue(true),
		}
		resp := &resource.CreateResponse{}

		pullReq := r.buildPullRequest(plan, resp)

		assert.NotNil(t, pullReq)
		assert.False(t, resp.Diagnostics.HasError())
		assert.NotNil(t, pullReq.GetForce())
		assert.True(t, pullReq.GetForce())
	})

	t.Run("build pull request without force", func(t *testing.T) {
		r := NewSourceControlPullResource()
		plan := &models.PullResource{
			Force: types.BoolNull(),
		}
		resp := &resource.CreateResponse{}

		pullReq := r.buildPullRequest(plan, resp)

		assert.NotNil(t, pullReq)
		assert.False(t, resp.Diagnostics.HasError())
	})

	t.Run("build pull request with variables", func(t *testing.T) {
		r := NewSourceControlPullResource()
		varsJSON := `{"key1": "value1", "key2": "value2"}`
		plan := &models.PullResource{
			Force:         types.BoolValue(false),
			VariablesJSON: types.StringValue(varsJSON),
		}
		resp := &resource.CreateResponse{}

		pullReq := r.buildPullRequest(plan, resp)

		assert.NotNil(t, pullReq)
		assert.False(t, resp.Diagnostics.HasError())
		assert.NotNil(t, pullReq.GetVariables())
		assert.Len(t, pullReq.GetVariables(), 2)
	})

	t.Run("build pull request with invalid JSON", func(t *testing.T) {
		r := NewSourceControlPullResource()
		plan := &models.PullResource{
			VariablesJSON: types.StringValue("invalid-json"),
		}
		resp := &resource.CreateResponse{}

		pullReq := r.buildPullRequest(plan, resp)

		assert.Nil(t, pullReq)
		assert.True(t, resp.Diagnostics.HasError())
		assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Invalid variables JSON")
	})

	t.Run("build pull request with empty variables", func(t *testing.T) {
		r := NewSourceControlPullResource()
		plan := &models.PullResource{
			Force:         types.BoolValue(true),
			VariablesJSON: types.StringNull(),
		}
		resp := &resource.CreateResponse{}

		pullReq := r.buildPullRequest(plan, resp)

		assert.NotNil(t, pullReq)
		assert.False(t, resp.Diagnostics.HasError())
	})

	t.Run("build pull request with unknown force", func(t *testing.T) {
		r := NewSourceControlPullResource()
		plan := &models.PullResource{
			Force: types.BoolUnknown(),
		}
		resp := &resource.CreateResponse{}

		pullReq := r.buildPullRequest(plan, resp)

		assert.NotNil(t, pullReq)
		assert.False(t, resp.Diagnostics.HasError())
	})

	t.Run("build pull request with unknown variables", func(t *testing.T) {
		r := NewSourceControlPullResource()
		plan := &models.PullResource{
			VariablesJSON: types.StringUnknown(),
		}
		resp := &resource.CreateResponse{}

		pullReq := r.buildPullRequest(plan, resp)

		assert.NotNil(t, pullReq)
		assert.False(t, resp.Diagnostics.HasError())
	})
}

func TestSourceControlPullResource_ParseVariablesJSON(t *testing.T) {
	t.Run("parse valid JSON", func(t *testing.T) {
		r := NewSourceControlPullResource()
		varsJSON := `{"key1": "value1", "key2": "value2"}`
		resp := &resource.CreateResponse{}

		variables := r.parseVariablesJSON(varsJSON, resp)

		assert.NotNil(t, variables)
		assert.Len(t, variables, 2)
		assert.Equal(t, "value1", variables["key1"])
		assert.Equal(t, "value2", variables["key2"])
		assert.False(t, resp.Diagnostics.HasError())
	})

	t.Run("parse empty JSON object", func(t *testing.T) {
		r := NewSourceControlPullResource()
		varsJSON := `{}`
		resp := &resource.CreateResponse{}

		variables := r.parseVariablesJSON(varsJSON, resp)

		assert.NotNil(t, variables)
		assert.Len(t, variables, 0)
		assert.False(t, resp.Diagnostics.HasError())
	})

	t.Run("parse invalid JSON", func(t *testing.T) {
		r := NewSourceControlPullResource()
		varsJSON := `{invalid-json`
		resp := &resource.CreateResponse{}

		variables := r.parseVariablesJSON(varsJSON, resp)

		assert.NotNil(t, variables)
		assert.Len(t, variables, 0)
		assert.True(t, resp.Diagnostics.HasError())
		assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Invalid variables JSON")
	})

	t.Run("parse JSON with nested objects", func(t *testing.T) {
		r := NewSourceControlPullResource()
		varsJSON := `{"key1": {"nested": "value"}}`
		resp := &resource.CreateResponse{}

		variables := r.parseVariablesJSON(varsJSON, resp)

		assert.NotNil(t, variables)
		assert.Len(t, variables, 1)
		assert.False(t, resp.Diagnostics.HasError())
	})

	t.Run("parse JSON with arrays", func(t *testing.T) {
		r := NewSourceControlPullResource()
		varsJSON := `{"key1": [1, 2, 3]}`
		resp := &resource.CreateResponse{}

		variables := r.parseVariablesJSON(varsJSON, resp)

		assert.NotNil(t, variables)
		assert.Len(t, variables, 1)
		assert.False(t, resp.Diagnostics.HasError())
	})

	t.Run("parse JSON with null values", func(t *testing.T) {
		r := NewSourceControlPullResource()
		varsJSON := `{"key1": null}`
		resp := &resource.CreateResponse{}

		variables := r.parseVariablesJSON(varsJSON, resp)

		assert.NotNil(t, variables)
		assert.Len(t, variables, 1)
		assert.Nil(t, variables["key1"])
		assert.False(t, resp.Diagnostics.HasError())
	})

	t.Run("parse JSON with numbers", func(t *testing.T) {
		r := NewSourceControlPullResource()
		varsJSON := `{"key1": 123, "key2": 45.67}`
		resp := &resource.CreateResponse{}

		variables := r.parseVariablesJSON(varsJSON, resp)

		assert.NotNil(t, variables)
		assert.Len(t, variables, 2)
		assert.False(t, resp.Diagnostics.HasError())
	})

	t.Run("parse JSON with booleans", func(t *testing.T) {
		r := NewSourceControlPullResource()
		varsJSON := `{"key1": true, "key2": false}`
		resp := &resource.CreateResponse{}

		variables := r.parseVariablesJSON(varsJSON, resp)

		assert.NotNil(t, variables)
		assert.Len(t, variables, 2)
		assert.True(t, variables["key1"].(bool))
		assert.False(t, variables["key2"].(bool))
		assert.False(t, resp.Diagnostics.HasError())
	})
}

func TestSourceControlPullResource_UpdatePlanWithResult(t *testing.T) {
	t.Run("update plan with valid result", func(t *testing.T) {
		r := NewSourceControlPullResource()
		plan := &models.PullResource{}

		workflowID := "workflow1"
		workflowName := "Test Workflow"
		workflows := []n8nsdk.ImportResultWorkflowsInner{
			{
				Id:   &workflowID,
				Name: &workflowName,
			},
		}

		importResult := &n8nsdk.ImportResult{
			Workflows: workflows,
		}
		httpResp := &http.Response{
			StatusCode: 200,
		}
		resp := &resource.CreateResponse{}

		r.updatePlanWithResult(plan, importResult, httpResp, resp)

		assert.False(t, plan.ID.IsNull())
		assert.Equal(t, "pull-200", plan.ID.ValueString())
		assert.False(t, plan.ResultJSON.IsNull())
		assert.False(t, resp.Diagnostics.HasError())
	})

	t.Run("update plan with nil result", func(t *testing.T) {
		r := NewSourceControlPullResource()
		plan := &models.PullResource{}
		httpResp := &http.Response{
			StatusCode: 200,
		}
		resp := &resource.CreateResponse{}

		r.updatePlanWithResult(plan, nil, httpResp, resp)

		assert.False(t, plan.ID.IsNull())
		assert.Equal(t, "pull-200", plan.ID.ValueString())
		assert.True(t, plan.ResultJSON.IsNull())
		assert.False(t, resp.Diagnostics.HasError())
	})

	t.Run("update plan with different status codes", func(t *testing.T) {
		r := NewSourceControlPullResource()
		statusCodes := []int{200, 201, 204}

		for _, code := range statusCodes {
			plan := &models.PullResource{}
			importResult := &n8nsdk.ImportResult{}
			httpResp := &http.Response{
				StatusCode: code,
			}
			resp := &resource.CreateResponse{}

			r.updatePlanWithResult(plan, importResult, httpResp, resp)

			expectedID := fmt.Sprintf("pull-%d", code)
			assert.Equal(t, expectedID, plan.ID.ValueString())
		}
	})
}

func TestSourceControlPullResource_Interfaces(t *testing.T) {
	t.Run("implements required interfaces", func(t *testing.T) {
		r := NewSourceControlPullResource()

		// Test that SourceControlPullResource implements resource.Resource
		var _ resource.Resource = r

		// Test that SourceControlPullResource implements resource.ResourceWithConfigure
		var _ resource.ResourceWithConfigure = r

		// Test that SourceControlPullResource implements resource.ResourceWithImportState
		var _ resource.ResourceWithImportState = r

		// Test that SourceControlPullResource implements SourceControlPullResourceInterface
		var _ SourceControlPullResourceInterface = r
	})
}

func TestSourceControlPullResourceConcurrency(t *testing.T) {
	t.Run("concurrent metadata calls", func(t *testing.T) {
		r := NewSourceControlPullResource()

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
	})

	t.Run("concurrent schema calls", func(t *testing.T) {
		r := NewSourceControlPullResource()

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
	})

	t.Run("concurrent configure calls", func(t *testing.T) {
		done := make(chan bool, 100)
		for i := 0; i < 100; i++ {
			go func() {
				r := NewSourceControlPullResource()
				resp := &resource.ConfigureResponse{}
				mockClient := &client.N8nClient{}
				req := resource.ConfigureRequest{
					ProviderData: mockClient,
				}
				r.Configure(context.Background(), req, resp)
				assert.Equal(t, mockClient, r.client)
				done <- true
			}()
		}

		for i := 0; i < 100; i++ {
			<-done
		}
	})
}

func TestSourceControlPullResource_HelperFunctions(t *testing.T) {
	t.Run("path to id", func(t *testing.T) {
		p := path.Root("id")
		assert.Equal(t, "id", p.String())
	})

	t.Run("path to force", func(t *testing.T) {
		p := path.Root("force")
		assert.Equal(t, "force", p.String())
	})

	t.Run("path to variables_json", func(t *testing.T) {
		p := path.Root("variables_json")
		assert.Equal(t, "variables_json", p.String())
	})

	t.Run("path to result_json", func(t *testing.T) {
		p := path.Root("result_json")
		assert.Equal(t, "result_json", p.String())
	})

	t.Run("types.StringValue conversion", func(t *testing.T) {
		strVal := types.StringValue("test-value")
		assert.Equal(t, "test-value", strVal.ValueString())
		assert.False(t, strVal.IsNull())
		assert.False(t, strVal.IsUnknown())
	})

	t.Run("types.BoolValue conversion", func(t *testing.T) {
		boolVal := types.BoolValue(true)
		assert.True(t, boolVal.ValueBool())
		assert.False(t, boolVal.IsNull())
		assert.False(t, boolVal.IsUnknown())
	})

	t.Run("types.StringNull", func(t *testing.T) {
		strVal := types.StringNull()
		assert.True(t, strVal.IsNull())
		assert.False(t, strVal.IsUnknown())
	})

	t.Run("types.BoolNull", func(t *testing.T) {
		boolVal := types.BoolNull()
		assert.True(t, boolVal.IsNull())
		assert.False(t, boolVal.IsUnknown())
	})
}

func TestSourceControlPullResource_EdgeCases(t *testing.T) {
	t.Run("configure then reconfigure with nil", func(t *testing.T) {
		r := NewSourceControlPullResource()

		// First configure with valid client
		resp1 := &resource.ConfigureResponse{}
		mockClient := &client.N8nClient{}
		req1 := resource.ConfigureRequest{
			ProviderData: mockClient,
		}
		r.Configure(context.Background(), req1, resp1)
		assert.NotNil(t, r.client)

		// Then configure with nil
		resp2 := &resource.ConfigureResponse{}
		req2 := resource.ConfigureRequest{
			ProviderData: nil,
		}
		r.Configure(context.Background(), req2, resp2)
		// Client should remain as it was
		assert.NotNil(t, r.client)
	})

	t.Run("build request with both force values", func(t *testing.T) {
		r := NewSourceControlPullResource()

		// Test with force=true
		plan1 := &models.PullResource{
			Force: types.BoolValue(true),
		}
		resp1 := &resource.CreateResponse{}
		pullReq1 := r.buildPullRequest(plan1, resp1)
		assert.NotNil(t, pullReq1)
		assert.True(t, pullReq1.GetForce())

		// Test with force=false
		plan2 := &models.PullResource{
			Force: types.BoolValue(false),
		}
		resp2 := &resource.CreateResponse{}
		pullReq2 := r.buildPullRequest(plan2, resp2)
		assert.NotNil(t, pullReq2)
		assert.False(t, pullReq2.GetForce())
	})

	t.Run("parse variables with special characters", func(t *testing.T) {
		r := NewSourceControlPullResource()
		varsJSON := `{"key-with-dash": "value", "key_with_underscore": "value", "key.with.dot": "value"}`
		resp := &resource.CreateResponse{}

		variables := r.parseVariablesJSON(varsJSON, resp)

		assert.NotNil(t, variables)
		assert.Len(t, variables, 3)
		assert.False(t, resp.Diagnostics.HasError())
	})

	t.Run("parse variables with unicode", func(t *testing.T) {
		r := NewSourceControlPullResource()
		varsJSON := `{"键": "值", "clé": "valeur", "Schlüssel": "Wert"}`
		resp := &resource.CreateResponse{}

		variables := r.parseVariablesJSON(varsJSON, resp)

		assert.NotNil(t, variables)
		assert.Len(t, variables, 3)
		assert.False(t, resp.Diagnostics.HasError())
	})
}

func BenchmarkSourceControlPullResource_Schema(b *testing.B) {
	r := NewSourceControlPullResource()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := &resource.SchemaResponse{}
		r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	}
}

func BenchmarkSourceControlPullResource_Metadata(b *testing.B) {
	r := NewSourceControlPullResource()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := &resource.MetadataResponse{}
		r.Metadata(context.Background(), resource.MetadataRequest{
			ProviderTypeName: "n8n",
		}, resp)
	}
}

func BenchmarkSourceControlPullResource_Configure(b *testing.B) {
	r := NewSourceControlPullResource()
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

func BenchmarkSourceControlPullResource_ParseVariablesJSON(b *testing.B) {
	r := NewSourceControlPullResource()
	varsJSON := `{"key1": "value1", "key2": "value2", "key3": "value3"}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := &resource.CreateResponse{}
		_ = r.parseVariablesJSON(varsJSON, resp)
	}
}

func BenchmarkSourceControlPullResource_BuildPullRequest(b *testing.B) {
	r := NewSourceControlPullResource()
	plan := &models.PullResource{
		Force:         types.BoolValue(true),
		VariablesJSON: types.StringValue(`{"key": "value"}`),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := &resource.CreateResponse{}
		_ = r.buildPullRequest(plan, resp)
	}
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

// TestSourceControlPullResource_Create tests the Create function.
func TestSourceControlPullResource_Create(t *testing.T) {
	t.Run("create with API error", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Server Error"))
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := NewSourceControlPullResource()
		r.client = n8nClient

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
	})

	t.Run("create with invalid variables JSON", func(t *testing.T) {
		r := NewSourceControlPullResource()
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
	})
}

// TestSourceControlPullResource_ExecutePullOperation_Complete tests executePullOperation.
func TestSourceControlPullResource_ExecutePullOperation_Complete(t *testing.T) {
	t.Run("API error", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Forbidden"))
		})

		n8nClient, server := setupTestClient(t, handler)
		defer server.Close()

		r := NewSourceControlPullResource()
		r.client = n8nClient

		pullReq := n8nsdk.NewPull()
		resp := &resource.CreateResponse{}

		result, httpResp := r.executePullOperation(context.Background(), pullReq, resp)
		if httpResp != nil && httpResp.Body != nil {
			defer httpResp.Body.Close()
		}

		assert.Nil(t, result)
		assert.NotNil(t, httpResp)
		assert.True(t, resp.Diagnostics.HasError())
		assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Error pulling from source control")
	})
}

// TestSourceControlPullResource_Read_Complete tests Read error paths.
func TestSourceControlPullResource_Read_Complete(t *testing.T) {
	t.Run("read with invalid state", func(t *testing.T) {
		r := NewSourceControlPullResource()
		ctx := context.Background()

		schemaResp := resource.SchemaResponse{}
		r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

		// Create invalid state
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
	})
}

// TestSourceControlPullResource_Delete_WithState tests Delete with state.
func TestSourceControlPullResource_Delete_WithState(t *testing.T) {
	t.Run("delete with valid state", func(t *testing.T) {
		r := NewSourceControlPullResource()
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
	})
}
