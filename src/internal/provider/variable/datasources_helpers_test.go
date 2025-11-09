package variable

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
	"github.com/kodflow/n8n/src/internal/provider/variable/models"
	"github.com/stretchr/testify/assert"
)

// Test helper functions for datasources Read functionality.

func TestVariablesDataSource_BuildAPIRequestWithFilters(t *testing.T) {
	t.Run("build request with no filters", func(t *testing.T) {
		ds := &VariablesDataSource{}
		apiClient := n8nsdk.NewAPIClient(n8nsdk.NewConfiguration())
		ds.client = &client.N8nClient{APIClient: apiClient}

		data := &models.DataSources{}

		req := ds.buildAPIRequestWithFilters(context.Background(), data)

		assert.NotNil(t, req)
	})

	t.Run("build request with project_id filter", func(t *testing.T) {
		ds := &VariablesDataSource{}
		apiClient := n8nsdk.NewAPIClient(n8nsdk.NewConfiguration())
		ds.client = &client.N8nClient{APIClient: apiClient}

		data := &models.DataSources{
			ProjectID: types.StringValue("proj-123"),
		}

		req := ds.buildAPIRequestWithFilters(context.Background(), data)

		assert.NotNil(t, req)
	})

	t.Run("build request with state filter", func(t *testing.T) {
		ds := &VariablesDataSource{}
		apiClient := n8nsdk.NewAPIClient(n8nsdk.NewConfiguration())
		ds.client = &client.N8nClient{APIClient: apiClient}

		data := &models.DataSources{
			State: types.StringValue("active"),
		}

		req := ds.buildAPIRequestWithFilters(context.Background(), data)

		assert.NotNil(t, req)
	})

	t.Run("build request with both filters", func(t *testing.T) {
		ds := &VariablesDataSource{}
		apiClient := n8nsdk.NewAPIClient(n8nsdk.NewConfiguration())
		ds.client = &client.N8nClient{APIClient: apiClient}

		data := &models.DataSources{
			ProjectID: types.StringValue("proj-456"),
			State:     types.StringValue("inactive"),
		}

		req := ds.buildAPIRequestWithFilters(context.Background(), data)

		assert.NotNil(t, req)
	})

	t.Run("build request with null filters", func(t *testing.T) {
		ds := &VariablesDataSource{}
		apiClient := n8nsdk.NewAPIClient(n8nsdk.NewConfiguration())
		ds.client = &client.N8nClient{APIClient: apiClient}

		data := &models.DataSources{
			ProjectID: types.StringNull(),
			State:     types.StringNull(),
		}

		req := ds.buildAPIRequestWithFilters(context.Background(), data)

		assert.NotNil(t, req)
	})
}

func TestVariablesDataSource_PopulateVariables(t *testing.T) {
	t.Run("populate with valid data", func(t *testing.T) {
		ds := &VariablesDataSource{}
		data := &models.DataSources{}

		id1 := "var-1"
		id2 := "var-2"
		varType := "string"
		projID := "proj-123"

		variableList := &n8nsdk.VariableList{
			Data: []n8nsdk.Variable{
				{
					Id:    &id1,
					Key:   "API_KEY",
					Value: "secret",
					Type:  &varType,
					Project: &n8nsdk.Project{
						Id: &projID,
					},
				},
				{
					Id:    &id2,
					Key:   "DB_URL",
					Value: "postgres://",
					Type:  &varType,
				},
			},
		}

		ds.populateVariables(data, variableList)

		assert.Len(t, data.Variables, 2)
		assert.Equal(t, "var-1", data.Variables[0].ID.ValueString())
		assert.Equal(t, "API_KEY", data.Variables[0].Key.ValueString())
		assert.Equal(t, "proj-123", data.Variables[0].ProjectID.ValueString())
		assert.Equal(t, "var-2", data.Variables[1].ID.ValueString())
		assert.Equal(t, "DB_URL", data.Variables[1].Key.ValueString())
	})

	t.Run("populate with nil data", func(t *testing.T) {
		ds := &VariablesDataSource{}
		data := &models.DataSources{}

		variableList := &n8nsdk.VariableList{
			Data: nil,
		}

		ds.populateVariables(data, variableList)

		assert.Empty(t, data.Variables)
	})

	t.Run("populate with empty data", func(t *testing.T) {
		ds := &VariablesDataSource{}
		data := &models.DataSources{}

		variableList := &n8nsdk.VariableList{
			Data: []n8nsdk.Variable{},
		}

		ds.populateVariables(data, variableList)

		assert.Empty(t, data.Variables)
	})

	t.Run("populate with multiple variables", func(t *testing.T) {
		ds := &VariablesDataSource{}
		data := &models.DataSources{}

		id1 := "v1"
		id2 := "v2"
		id3 := "v3"

		variableList := &n8nsdk.VariableList{
			Data: []n8nsdk.Variable{
				{Id: &id1, Key: "KEY1", Value: "val1"},
				{Id: &id2, Key: "KEY2", Value: "val2"},
				{Id: &id3, Key: "KEY3", Value: "val3"},
			},
		}

		ds.populateVariables(data, variableList)

		assert.Len(t, data.Variables, 3)
	})
}

func TestVariablesDataSource_MapVariableToItem(t *testing.T) {
	t.Run("map with all fields", func(t *testing.T) {
		ds := &VariablesDataSource{}

		id := "var-123"
		varType := "string"
		projID := "proj-456"

		variable := &n8nsdk.Variable{
			Id:    &id,
			Key:   "TEST_KEY",
			Value: "test-value",
			Type:  &varType,
			Project: &n8nsdk.Project{
				Id: &projID,
			},
		}

		item := ds.mapVariableToItem(variable)

		assert.Equal(t, "var-123", item.ID.ValueString())
		assert.Equal(t, "TEST_KEY", item.Key.ValueString())
		assert.Equal(t, "test-value", item.Value.ValueString())
		assert.Equal(t, "string", item.Type.ValueString())
		assert.Equal(t, "proj-456", item.ProjectID.ValueString())
	})

	t.Run("map with nil optional fields", func(t *testing.T) {
		ds := &VariablesDataSource{}

		variable := &n8nsdk.Variable{
			Id:      nil,
			Key:     "KEY",
			Value:   "value",
			Type:    nil,
			Project: nil,
		}

		item := ds.mapVariableToItem(variable)

		assert.True(t, item.ID.IsNull())
		assert.Equal(t, "KEY", item.Key.ValueString())
		assert.Equal(t, "value", item.Value.ValueString())
		assert.True(t, item.Type.IsNull())
		assert.True(t, item.ProjectID.IsNull())
	})

	t.Run("map with project but nil project ID", func(t *testing.T) {
		ds := &VariablesDataSource{}

		id := "var-789"
		variable := &n8nsdk.Variable{
			Id:    &id,
			Key:   "KEY2",
			Value: "value2",
			Project: &n8nsdk.Project{
				Id: nil,
			},
		}

		item := ds.mapVariableToItem(variable)

		assert.Equal(t, "var-789", item.ID.ValueString())
		assert.True(t, item.ProjectID.IsNull())
	})

	t.Run("map with empty strings", func(t *testing.T) {
		ds := &VariablesDataSource{}

		id := ""
		varType := ""
		projID := ""

		variable := &n8nsdk.Variable{
			Id:    &id,
			Key:   "",
			Value: "",
			Type:  &varType,
			Project: &n8nsdk.Project{
				Id: &projID,
			},
		}

		item := ds.mapVariableToItem(variable)

		assert.Equal(t, "", item.ID.ValueString())
		assert.Equal(t, "", item.Key.ValueString())
		assert.Equal(t, "", item.Value.ValueString())
		assert.Equal(t, "", item.Type.ValueString())
		assert.Equal(t, "", item.ProjectID.ValueString())
	})

	t.Run("map preserves sensitive value", func(t *testing.T) {
		ds := &VariablesDataSource{}

		id := "var-secret"
		variable := &n8nsdk.Variable{
			Id:    &id,
			Key:   "API_SECRET",
			Value: "very-secret-value",
		}

		item := ds.mapVariableToItem(variable)

		assert.Equal(t, "var-secret", item.ID.ValueString())
		assert.Equal(t, "very-secret-value", item.Value.ValueString())
	})
}

func TestVariablesDataSource_Read(t *testing.T) {
	t.Run("successful read with variables", func(t *testing.T) {
		// Create a mock HTTP server
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/variables" && r.Method == http.MethodGet {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]interface{}{
					"data": []map[string]interface{}{
						{
							"id":    "var-1",
							"key":   "API_KEY",
							"value": "secret",
							"type":  "string",
						},
					},
				})
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestVariablesDataSourceClient(t, handler)
		defer server.Close()

		ds := &VariablesDataSource{client: n8nClient}

		req := datasource.ReadRequest{
			Config: createTestVariablesDataSourceConfigHelper(t),
		}
		resp := &datasource.ReadResponse{
			State: createTestVariablesDataSourceStateHelper(t),
		}

		ds.Read(context.Background(), req, resp)

		assert.False(t, resp.Diagnostics.HasError())
	})

	t.Run("read fails on API error", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/variables" && r.Method == http.MethodGet {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"message": "error"}`))
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestVariablesDataSourceClient(t, handler)
		defer server.Close()

		ds := &VariablesDataSource{client: n8nClient}

		req := datasource.ReadRequest{
			Config: createTestVariablesDataSourceConfigHelper(t),
		}
		resp := &datasource.ReadResponse{
			State: createTestVariablesDataSourceStateHelper(t),
		}

		ds.Read(context.Background(), req, resp)

		assert.True(t, resp.Diagnostics.HasError())
		assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Error listing variables")
	})

	t.Run("read with empty response", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/variables" && r.Method == http.MethodGet {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]interface{}{
					"data": []interface{}{},
				})
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestVariablesDataSourceClient(t, handler)
		defer server.Close()

		ds := &VariablesDataSource{client: n8nClient}

		req := datasource.ReadRequest{
			Config: createTestVariablesDataSourceConfigHelper(t),
		}
		resp := &datasource.ReadResponse{
			State: createTestVariablesDataSourceStateHelper(t),
		}

		ds.Read(context.Background(), req, resp)

		assert.False(t, resp.Diagnostics.HasError())
	})

	t.Run("read with null data", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/variables" && r.Method == http.MethodGet {
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{}`))
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupTestVariablesDataSourceClient(t, handler)
		defer server.Close()

		ds := &VariablesDataSource{client: n8nClient}

		req := datasource.ReadRequest{
			Config: createTestVariablesDataSourceConfigHelper(t),
		}
		resp := &datasource.ReadResponse{
			State: createTestVariablesDataSourceStateHelper(t),
		}

		ds.Read(context.Background(), req, resp)

		assert.False(t, resp.Diagnostics.HasError())
	})

	t.Run("read with invalid config", func(t *testing.T) {
		n8nClient, server := setupTestVariablesDataSourceClient(t, nil)
		defer server.Close()

		ds := &VariablesDataSource{client: n8nClient}

		// Create config with wrong schema
		wrongSchema := datasource.SchemaResponse{}
		req := datasource.ReadRequest{
			Config: tfsdk.Config{
				Schema: wrongSchema.Schema,
				Raw:    tftypes.NewValue(tftypes.Object{}, nil),
			},
		}
		resp := &datasource.ReadResponse{
			State: createTestVariablesDataSourceStateHelper(t),
		}

		ds.Read(context.Background(), req, resp)

		assert.True(t, resp.Diagnostics.HasError())
	})
}

// Helper functions for Read tests.
func setupTestVariablesDataSourceClient(t *testing.T, handler http.HandlerFunc) (*client.N8nClient, *httptest.Server) {
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
		BaseURL:   server.URL,
	}

	return n8nClient, server
}

func createTestVariablesDataSourceStateHelper(t *testing.T) tfsdk.State {
	t.Helper()

	ds := NewVariablesDataSource()
	schemaResp := &datasource.SchemaResponse{}
	ds.Schema(context.Background(), datasource.SchemaRequest{}, schemaResp)

	return tfsdk.State{
		Schema: schemaResp.Schema,
		Raw:    tftypes.NewValue(schemaResp.Schema.Type().TerraformType(context.Background()), nil),
	}
}

func createTestVariablesDataSourceConfigHelper(t *testing.T) tfsdk.Config {
	t.Helper()

	ds := NewVariablesDataSource()
	schemaResp := &datasource.SchemaResponse{}
	ds.Schema(context.Background(), datasource.SchemaRequest{}, schemaResp)

	// Create a proper config with null values for optional fields
	tfType := schemaResp.Schema.Type().TerraformType(context.Background())
	configValue := tftypes.NewValue(tfType, map[string]tftypes.Value{
		"project_id": tftypes.NewValue(tftypes.String, nil),
		"state":      tftypes.NewValue(tftypes.String, nil),
		"variables":  tftypes.NewValue(tftypes.List{ElementType: tftypes.Object{AttributeTypes: map[string]tftypes.Type{
			"id":         tftypes.String,
			"key":        tftypes.String,
			"value":      tftypes.String,
			"type":       tftypes.String,
			"project_id": tftypes.String,
		}}}, nil),
	})

	return tfsdk.Config{
		Schema: schemaResp.Schema,
		Raw:    configValue,
	}
}
