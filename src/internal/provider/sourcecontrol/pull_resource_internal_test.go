package sourcecontrol

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/client"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/sourcecontrol/models"
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

// TestSourceControlPullResource_schemaAttributes tests the schemaAttributes private function.
func TestSourceControlPullResource_schemaAttributes(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "returns all required attributes",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := NewSourceControlPullResource()

				attrs := r.schemaAttributes()

				assert.NotNil(t, attrs)
				assert.Len(t, attrs, 4)
				assert.Contains(t, attrs, "id")
				assert.Contains(t, attrs, "force")
				assert.Contains(t, attrs, "variables_json")
				assert.Contains(t, attrs, "result_json")
			},
		},
		{
			name: "id attribute is computed",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := NewSourceControlPullResource()

				attrs := r.schemaAttributes()

				idAttr := attrs["id"]
				assert.True(t, idAttr.IsComputed())
				assert.False(t, idAttr.IsOptional())
				assert.False(t, idAttr.IsRequired())
			},
		},
		{
			name: "force attribute is optional",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := NewSourceControlPullResource()

				attrs := r.schemaAttributes()

				forceAttr := attrs["force"]
				assert.True(t, forceAttr.IsOptional())
				assert.False(t, forceAttr.IsComputed())
				assert.False(t, forceAttr.IsRequired())
			},
		},
		{
			name: "variables_json attribute is optional",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := NewSourceControlPullResource()

				attrs := r.schemaAttributes()

				varsAttr := attrs["variables_json"]
				assert.True(t, varsAttr.IsOptional())
				assert.False(t, varsAttr.IsComputed())
				assert.False(t, varsAttr.IsRequired())
			},
		},
		{
			name: "result_json attribute is computed",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := NewSourceControlPullResource()

				attrs := r.schemaAttributes()

				resultAttr := attrs["result_json"]
				assert.True(t, resultAttr.IsComputed())
				assert.False(t, resultAttr.IsOptional())
				assert.False(t, resultAttr.IsRequired())
			},
		},
		{
			name: "all attributes have descriptions",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := NewSourceControlPullResource()

				attrs := r.schemaAttributes()

				for name, attr := range attrs {
					assert.NotEmpty(t, attr.GetMarkdownDescription(), "Attribute %s should have description", name)
				}
			},
		},
		{
			name: "error - verify returned map is not nil",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := NewSourceControlPullResource()

				attrs := r.schemaAttributes()

				assert.NotNil(t, attrs)
				assert.Greater(t, len(attrs), 0)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}

// TestSourceControlPullResource_buildPullRequest tests the buildPullRequest private function.
func TestSourceControlPullResource_buildPullRequest(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "build with force true",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := NewSourceControlPullResource()
				plan := &models.PullResource{
					Force: types.BoolValue(true),
				}
				resp := &resource.CreateResponse{}

				pullReq := r.buildPullRequest(plan, resp)

				assert.NotNil(t, pullReq)
				assert.False(t, resp.Diagnostics.HasError())
				assert.True(t, pullReq.GetForce())
			},
		},
		{
			name: "build with force false",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := NewSourceControlPullResource()
				plan := &models.PullResource{
					Force: types.BoolValue(false),
				}
				resp := &resource.CreateResponse{}

				pullReq := r.buildPullRequest(plan, resp)

				assert.NotNil(t, pullReq)
				assert.False(t, resp.Diagnostics.HasError())
				assert.False(t, pullReq.GetForce())
			},
		},
		{
			name: "build with null force",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := NewSourceControlPullResource()
				plan := &models.PullResource{
					Force: types.BoolNull(),
				}
				resp := &resource.CreateResponse{}

				pullReq := r.buildPullRequest(plan, resp)

				assert.NotNil(t, pullReq)
				assert.False(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "build with unknown force",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := NewSourceControlPullResource()
				plan := &models.PullResource{
					Force: types.BoolUnknown(),
				}
				resp := &resource.CreateResponse{}

				pullReq := r.buildPullRequest(plan, resp)

				assert.NotNil(t, pullReq)
				assert.False(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "build with valid variables",
			testFunc: func(t *testing.T) {
				t.Helper()
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
			},
		},
		{
			name: "build with null variables",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := NewSourceControlPullResource()
				plan := &models.PullResource{
					Force:         types.BoolValue(true),
					VariablesJSON: types.StringNull(),
				}
				resp := &resource.CreateResponse{}

				pullReq := r.buildPullRequest(plan, resp)

				assert.NotNil(t, pullReq)
				assert.False(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "build with unknown variables",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := NewSourceControlPullResource()
				plan := &models.PullResource{
					VariablesJSON: types.StringUnknown(),
				}
				resp := &resource.CreateResponse{}

				pullReq := r.buildPullRequest(plan, resp)

				assert.NotNil(t, pullReq)
				assert.False(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "error - invalid JSON in variables",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := NewSourceControlPullResource()
				plan := &models.PullResource{
					VariablesJSON: types.StringValue("invalid-json"),
				}
				resp := &resource.CreateResponse{}

				pullReq := r.buildPullRequest(plan, resp)

				assert.Nil(t, pullReq)
				assert.True(t, resp.Diagnostics.HasError())
				assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Invalid variables JSON")
			},
		},
		{
			name: "build with force and variables",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := NewSourceControlPullResource()
				varsJSON := `{"env": "production", "debug": "false"}`
				plan := &models.PullResource{
					Force:         types.BoolValue(true),
					VariablesJSON: types.StringValue(varsJSON),
				}
				resp := &resource.CreateResponse{}

				pullReq := r.buildPullRequest(plan, resp)

				assert.NotNil(t, pullReq)
				assert.False(t, resp.Diagnostics.HasError())
				assert.True(t, pullReq.GetForce())
				assert.Len(t, pullReq.GetVariables(), 2)
			},
		},
		{
			name: "build with empty variables object",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := NewSourceControlPullResource()
				plan := &models.PullResource{
					VariablesJSON: types.StringValue("{}"),
				}
				resp := &resource.CreateResponse{}

				pullReq := r.buildPullRequest(plan, resp)

				assert.NotNil(t, pullReq)
				assert.False(t, resp.Diagnostics.HasError())
				assert.NotNil(t, pullReq.GetVariables())
				assert.Len(t, pullReq.GetVariables(), 0)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}

// TestSourceControlPullResource_parseVariablesJSON tests the parseVariablesJSON private function.
func TestSourceControlPullResource_parseVariablesJSON(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "parse valid JSON with strings",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := NewSourceControlPullResource()
				varsJSON := `{"key1": "value1", "key2": "value2"}`
				resp := &resource.CreateResponse{}

				variables := r.parseVariablesJSON(varsJSON, resp)

				assert.NotNil(t, variables)
				assert.Len(t, variables, 2)
				assert.Equal(t, "value1", variables["key1"])
				assert.Equal(t, "value2", variables["key2"])
				assert.False(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "parse empty JSON object",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := NewSourceControlPullResource()
				varsJSON := `{}`
				resp := &resource.CreateResponse{}

				variables := r.parseVariablesJSON(varsJSON, resp)

				assert.NotNil(t, variables)
				assert.Len(t, variables, 0)
				assert.False(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "parse JSON with nested objects",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := NewSourceControlPullResource()
				varsJSON := `{"key1": {"nested": "value"}}`
				resp := &resource.CreateResponse{}

				variables := r.parseVariablesJSON(varsJSON, resp)

				assert.NotNil(t, variables)
				assert.Len(t, variables, 1)
				assert.False(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "parse JSON with arrays",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := NewSourceControlPullResource()
				varsJSON := `{"key1": [1, 2, 3]}`
				resp := &resource.CreateResponse{}

				variables := r.parseVariablesJSON(varsJSON, resp)

				assert.NotNil(t, variables)
				assert.Len(t, variables, 1)
				assert.False(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "parse JSON with null values",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := NewSourceControlPullResource()
				varsJSON := `{"key1": null}`
				resp := &resource.CreateResponse{}

				variables := r.parseVariablesJSON(varsJSON, resp)

				assert.NotNil(t, variables)
				assert.Len(t, variables, 1)
				assert.Nil(t, variables["key1"])
				assert.False(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "parse JSON with numbers",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := NewSourceControlPullResource()
				varsJSON := `{"key1": 123, "key2": 45.67}`
				resp := &resource.CreateResponse{}

				variables := r.parseVariablesJSON(varsJSON, resp)

				assert.NotNil(t, variables)
				assert.Len(t, variables, 2)
				assert.False(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "parse JSON with booleans",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := NewSourceControlPullResource()
				varsJSON := `{"key1": true, "key2": false}`
				resp := &resource.CreateResponse{}

				variables := r.parseVariablesJSON(varsJSON, resp)

				assert.NotNil(t, variables)
				assert.Len(t, variables, 2)
				assert.True(t, variables["key1"].(bool))
				assert.False(t, variables["key2"].(bool))
				assert.False(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "parse JSON with special characters",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := NewSourceControlPullResource()
				varsJSON := `{"key-with-dash": "value", "key_with_underscore": "value", "key.with.dot": "value"}`
				resp := &resource.CreateResponse{}

				variables := r.parseVariablesJSON(varsJSON, resp)

				assert.NotNil(t, variables)
				assert.Len(t, variables, 3)
				assert.False(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "parse JSON with unicode characters",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := NewSourceControlPullResource()
				varsJSON := `{"键": "值", "clé": "valeur", "Schlüssel": "Wert"}`
				resp := &resource.CreateResponse{}

				variables := r.parseVariablesJSON(varsJSON, resp)

				assert.NotNil(t, variables)
				assert.Len(t, variables, 3)
				assert.False(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "error - invalid JSON syntax",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := NewSourceControlPullResource()
				varsJSON := `{invalid-json`
				resp := &resource.CreateResponse{}

				variables := r.parseVariablesJSON(varsJSON, resp)

				assert.NotNil(t, variables)
				assert.Len(t, variables, 0)
				assert.True(t, resp.Diagnostics.HasError())
				assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Invalid variables JSON")
			},
		},
		{
			name: "error - malformed JSON",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := NewSourceControlPullResource()
				varsJSON := `{"key": "value",}`
				resp := &resource.CreateResponse{}

				variables := r.parseVariablesJSON(varsJSON, resp)

				assert.NotNil(t, variables)
				assert.Len(t, variables, 0)
				assert.True(t, resp.Diagnostics.HasError())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}

// TestSourceControlPullResource_executePullOperation tests the executePullOperation private function.
func TestSourceControlPullResource_executePullOperation(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "error - 403 forbidden",
			testFunc: func(t *testing.T) {
				t.Helper()
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
			},
		},
		{
			name: "error - 500 internal server error",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("Internal Server Error"))
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
			},
		},
		{
			name: "error - 401 unauthorized",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("Unauthorized"))
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
			},
		},
		{
			name: "error - 404 not found",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte("Not Found"))
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
			},
		},
		{
			name: "error - 400 bad request",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte("Bad Request"))
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
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}

// TestSourceControlPullResource_updatePlanWithResult tests the updatePlanWithResult private function.
func TestSourceControlPullResource_updatePlanWithResult(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "update with valid result",
			testFunc: func(t *testing.T) {
				t.Helper()
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
			},
		},
		{
			name: "update with nil result",
			testFunc: func(t *testing.T) {
				t.Helper()
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
			},
		},
		{
			name: "update with 201 status code",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := NewSourceControlPullResource()
				plan := &models.PullResource{}
				importResult := &n8nsdk.ImportResult{}
				httpResp := &http.Response{
					StatusCode: 201,
				}
				resp := &resource.CreateResponse{}

				r.updatePlanWithResult(plan, importResult, httpResp, resp)

				assert.Equal(t, "pull-201", plan.ID.ValueString())
			},
		},
		{
			name: "update with 204 status code",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := NewSourceControlPullResource()
				plan := &models.PullResource{}
				importResult := &n8nsdk.ImportResult{}
				httpResp := &http.Response{
					StatusCode: 204,
				}
				resp := &resource.CreateResponse{}

				r.updatePlanWithResult(plan, importResult, httpResp, resp)

				assert.Equal(t, "pull-204", plan.ID.ValueString())
			},
		},
		{
			name: "update with empty workflows",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := NewSourceControlPullResource()
				plan := &models.PullResource{}

				importResult := &n8nsdk.ImportResult{
					Workflows: []n8nsdk.ImportResultWorkflowsInner{},
				}
				httpResp := &http.Response{
					StatusCode: 200,
				}
				resp := &resource.CreateResponse{}

				r.updatePlanWithResult(plan, importResult, httpResp, resp)

				assert.False(t, plan.ID.IsNull())
				assert.False(t, plan.ResultJSON.IsNull())
				assert.False(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "update with multiple workflows",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := NewSourceControlPullResource()
				plan := &models.PullResource{}

				wf1ID := "wf1"
				wf1Name := "Workflow 1"
				wf2ID := "wf2"
				wf2Name := "Workflow 2"
				workflows := []n8nsdk.ImportResultWorkflowsInner{
					{Id: &wf1ID, Name: &wf1Name},
					{Id: &wf2ID, Name: &wf2Name},
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
				assert.False(t, plan.ResultJSON.IsNull())
				assert.False(t, resp.Diagnostics.HasError())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}

// TestSourceControlPullResource_Delete_Internal tests the Delete method (internal test for coverage).
func TestSourceControlPullResource_Delete_Internal(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "delete does nothing - no error",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := NewSourceControlPullResource()
				ctx := context.Background()
				req := resource.DeleteRequest{}
				resp := &resource.DeleteResponse{}

				r.Delete(ctx, req, resp)

				assert.False(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "error case - delete with context",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := NewSourceControlPullResource()
				ctx := context.Background()
				req := resource.DeleteRequest{}
				resp := &resource.DeleteResponse{}

				// Delete should do nothing regardless of state
				r.Delete(ctx, req, resp)

				assert.False(t, resp.Diagnostics.HasError())
				assert.Equal(t, 0, resp.Diagnostics.ErrorsCount())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}

// TestSourceControlPullResource_Create_Internal tests the Create method (internal test for coverage).
func TestSourceControlPullResource_Create_Internal(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "error - create with invalid plan (Plan.Get fails)",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := NewSourceControlPullResource()
				ctx := context.Background()

				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				// Create an invalid plan value that will cause Plan.Get to fail
				// Using an invalid type that doesn't match the schema
				planRaw := tftypes.NewValue(tftypes.String, "invalid-not-an-object")

				plan := tfsdk.Plan{
					Schema: schemaResp.Schema,
					Raw:    planRaw,
				}

				req := resource.CreateRequest{
					Plan: plan,
				}
				resp := &resource.CreateResponse{}

				// Call Create - this should fail at req.Plan.Get and hit lines 148-151
				r.Create(ctx, req, resp)

				// Verify it has error
				assert.True(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "create success - full method with state.Set",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{"workflows": [{"id": "1", "name": "test"}], "credentials": []}`))
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				r := NewSourceControlPullResource()
				r.Configure(context.Background(), resource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &resource.ConfigureResponse{})

				ctx := context.Background()
				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				// Build plan using tftypes
				planRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"id":             tftypes.NewValue(tftypes.String, nil),
					"force":          tftypes.NewValue(tftypes.Bool, true),
					"variables_json": tftypes.NewValue(tftypes.String, `{"key": "value"}`),
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

				// Call the full Create method - this hits lines 174-175
				r.Create(ctx, req, resp)

				// Print diagnostics if there are errors for debugging
				if resp.Diagnostics.HasError() {
					for _, diag := range resp.Diagnostics.Errors() {
						t.Logf("Error: %s - %s", diag.Summary(), diag.Detail())
					}
				}

				// Verify success
				assert.False(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "create with variables JSON and body close",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					// Create response with actual body that needs closing
					w.Write([]byte(`{"workflows": [{"id": "1", "name": "test"}]}`))
				})

				n8nClient, server := setupTestClient(t, handler)
				defer server.Close()

				r := NewSourceControlPullResource()
				r.Configure(context.Background(), resource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &resource.ConfigureResponse{})

				ctx := context.Background()
				plan := &models.PullResource{
					ID:            types.StringNull(),
					Force:         types.BoolValue(true),
					VariablesJSON: types.StringValue(`{"key": "value"}`),
					ResultJSON:    types.StringNull(),
				}

				// Build request
				pullReq := r.buildPullRequest(plan, &resource.CreateResponse{})
				assert.NotNil(t, pullReq)

				// Execute pull
				resp := &resource.CreateResponse{}
				importResult, httpResp := r.executePullOperation(ctx, pullReq, resp)

				// This tests the httpResp.Body != nil path
				if httpResp != nil && httpResp.Body != nil {
					defer httpResp.Body.Close()
					assert.NotNil(t, httpResp.Body)
				}

				// Update plan
				r.updatePlanWithResult(plan, importResult, httpResp, resp)
				assert.False(t, resp.Diagnostics.HasError())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}
