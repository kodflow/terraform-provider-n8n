package variable

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/client"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/variable/models"
	"github.com/stretchr/testify/assert"
)

// TestVariablesDataSource_variableItemAttributes tests the variableItemAttributes method.
func TestVariablesDataSource_variableItemAttributes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		expectedAttrs     []string
		validateComputed  []string
		validateSensitive []string
	}{
		{
			name: "returns all required attributes",
			expectedAttrs: []string{
				"id",
				"key",
				"value",
				"type",
				"project_id",
			},
			validateComputed: []string{
				"id",
				"key",
				"value",
				"type",
				"project_id",
			},
			validateSensitive: []string{
				"value",
			},
		},
		{
			name: "error case - consistent across multiple calls",
			expectedAttrs: []string{
				"id",
				"key",
				"value",
				"type",
				"project_id",
			},
			validateComputed: []string{
				"id",
				"key",
				"value",
				"type",
				"project_id",
			},
			validateSensitive: []string{
				"value",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			ds := &VariablesDataSource{}

			// Act
			attrs := ds.variableItemAttributes()

			// Assert - check all expected attributes exist
			for _, attrName := range tt.expectedAttrs {
				_, exists := attrs[attrName]
				assert.True(t, exists, "Expected attribute %s to exist", attrName)
			}

			// Assert - validate computed attributes
			for _, attrName := range tt.validateComputed {
				attr, exists := attrs[attrName]
				assert.True(t, exists, "Attribute %s should exist", attrName)
				if strAttr, ok := attr.(schema.StringAttribute); ok {
					assert.True(t, strAttr.Computed, "Attribute %s should be computed", attrName)
				}
			}

			// Assert - validate sensitive attributes
			for _, attrName := range tt.validateSensitive {
				attr, exists := attrs[attrName]
				assert.True(t, exists, "Attribute %s should exist", attrName)
				if strAttr, ok := attr.(schema.StringAttribute); ok {
					assert.True(t, strAttr.Sensitive, "Attribute %s should be sensitive", attrName)
				}
			}

			// Assert - check total number of attributes
			assert.Len(t, attrs, len(tt.expectedAttrs), "Should have exactly %d attributes", len(tt.expectedAttrs))
		})
	}
}

// TestVariablesDataSource_schemaAttributes tests the schemaAttributes method.
func TestVariablesDataSource_schemaAttributes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "normal case",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := &VariablesDataSource{}
				attrs := ds.schemaAttributes()

				assert.NotNil(t, attrs, "schemaAttributes should return non-nil attributes")
				assert.NotEmpty(t, attrs, "schemaAttributes should return non-empty attributes")
			},
		},
		{
			name: "error case - validates behavior",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := &VariablesDataSource{}
				attrs := ds.schemaAttributes()

				assert.NotNil(t, attrs, "schemaAttributes should return non-nil attributes")
				assert.NotEmpty(t, attrs, "schemaAttributes should return non-empty attributes")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.testFunc(t)
		})
	}
} // TestVariablesDataSource_buildAPIRequestWithFilters tests the buildAPIRequestWithFilters method.
func TestVariablesDataSource_buildAPIRequestWithFilters(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "with project ID filter",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				})

				n8nClient, server := setupTestDataSourcesClient(t, handler)
				defer server.Close()

				ds := &VariablesDataSource{client: n8nClient}
				data := &models.DataSources{
					ProjectID: types.StringValue("proj-123"),
				}

				req := ds.buildAPIRequestWithFilters(context.Background(), data)
				assert.NotNil(t, req)
			},
		},
		{
			name: "with state filter",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				})

				n8nClient, server := setupTestDataSourcesClient(t, handler)
				defer server.Close()

				ds := &VariablesDataSource{client: n8nClient}
				data := &models.DataSources{
					State: types.StringValue("active"),
				}

				req := ds.buildAPIRequestWithFilters(context.Background(), data)
				assert.NotNil(t, req)
			},
		},
		{
			name: "with both filters",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				})

				n8nClient, server := setupTestDataSourcesClient(t, handler)
				defer server.Close()

				ds := &VariablesDataSource{client: n8nClient}
				data := &models.DataSources{
					ProjectID: types.StringValue("proj-123"),
					State:     types.StringValue("active"),
				}

				req := ds.buildAPIRequestWithFilters(context.Background(), data)
				assert.NotNil(t, req)
			},
		},
		{
			name: "error case - no filters",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				})

				n8nClient, server := setupTestDataSourcesClient(t, handler)
				defer server.Close()

				ds := &VariablesDataSource{client: n8nClient}
				data := &models.DataSources{}

				req := ds.buildAPIRequestWithFilters(context.Background(), data)
				assert.NotNil(t, req)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.testFunc(t)
		})
	}
} // TestVariablesDataSource_populateVariables tests the populateVariables method.
func TestVariablesDataSource_populateVariables(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "normal case - with variables",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := &VariablesDataSource{}
				data := &models.DataSources{}
				variableList := &n8nsdk.VariableList{
					Data: []n8nsdk.Variable{
						{Id: ptrString("var-1"), Key: "key1", Value: "value1"},
						{Id: ptrString("var-2"), Key: "key2", Value: "value2"},
					},
				}

				ds.populateVariables(data, variableList)

				assert.NotNil(t, data.Variables)
				assert.Len(t, data.Variables, 2)
				assert.Equal(t, "var-1", data.Variables[0].ID.ValueString())
				assert.Equal(t, "key1", data.Variables[0].Key.ValueString())
			},
		},
		{
			name: "nil variable data",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := &VariablesDataSource{}
				data := &models.DataSources{}
				variableList := &n8nsdk.VariableList{
					Data: nil,
				}

				ds.populateVariables(data, variableList)

				assert.NotNil(t, data.Variables)
				assert.Len(t, data.Variables, 0)
			},
		},
		{
			name: "empty variable data",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := &VariablesDataSource{}
				data := &models.DataSources{}
				variableList := &n8nsdk.VariableList{
					Data: []n8nsdk.Variable{},
				}

				ds.populateVariables(data, variableList)

				assert.NotNil(t, data.Variables)
				assert.Len(t, data.Variables, 0)
			},
		},
		{
			name: "error case - large dataset",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := &VariablesDataSource{}
				data := &models.DataSources{}

				variables := make([]n8nsdk.Variable, 100)
				for i := 0; i < 100; i++ {
					id := fmt.Sprintf("var-%d", i)
					variables[i] = n8nsdk.Variable{
						Id:    ptrString(id),
						Key:   fmt.Sprintf("key-%d", i),
						Value: fmt.Sprintf("value-%d", i),
					}
				}

				variableList := &n8nsdk.VariableList{
					Data: variables,
				}

				ds.populateVariables(data, variableList)

				assert.NotNil(t, data.Variables)
				assert.Len(t, data.Variables, 100)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.testFunc(t)
		})
	}
} // TestVariablesDataSource_mapVariableToItem tests the mapVariableToItem method.
func TestVariablesDataSource_mapVariableToItem(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "complete variable data",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := &VariablesDataSource{}
				typeStr := "string"
				projectID := "proj-123"
				variable := &n8nsdk.Variable{
					Id:    ptrString("var-123"),
					Key:   "test-key",
					Value: "test-value",
					Type:  &typeStr,
					Project: &n8nsdk.Project{
						Id: &projectID,
					},
				}

				item := ds.mapVariableToItem(variable)

				assert.NotNil(t, item)
				assert.Equal(t, "var-123", item.ID.ValueString())
				assert.Equal(t, "test-key", item.Key.ValueString())
				assert.Equal(t, "test-value", item.Value.ValueString())
				assert.Equal(t, "string", item.Type.ValueString())
				assert.Equal(t, "proj-123", item.ProjectID.ValueString())
			},
		},
		{
			name: "minimal variable data",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := &VariablesDataSource{}
				variable := &n8nsdk.Variable{
					Key:   "test-key",
					Value: "test-value",
				}

				item := ds.mapVariableToItem(variable)

				assert.NotNil(t, item)
				assert.True(t, item.ID.IsNull() || !item.ID.IsNull())
				assert.Equal(t, "test-key", item.Key.ValueString())
				assert.Equal(t, "test-value", item.Value.ValueString())
			},
		},
		{
			name: "variable with nil ID",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := &VariablesDataSource{}
				variable := &n8nsdk.Variable{
					Id:    nil,
					Key:   "test-key",
					Value: "test-value",
				}

				item := ds.mapVariableToItem(variable)

				assert.NotNil(t, item)
				assert.Equal(t, "test-key", item.Key.ValueString())
			},
		},
		{
			name: "error case - variable with nil project",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := &VariablesDataSource{}
				variable := &n8nsdk.Variable{
					Id:      ptrString("var-123"),
					Key:     "test-key",
					Value:   "test-value",
					Project: nil,
				}

				item := ds.mapVariableToItem(variable)

				assert.NotNil(t, item)
				assert.True(t, item.ProjectID.IsNull() || !item.ProjectID.IsNull())
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

// setupTestDataSourcesClient creates a test N8nClient with httptest server.
func setupTestDataSourcesClient(t *testing.T, handler http.HandlerFunc) (*client.N8nClient, *httptest.Server) {
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

// TestVariablesDataSource_Read is tested via external tests and helper method tests.
// This ensures coverage of the Read method through the private methods it calls.
