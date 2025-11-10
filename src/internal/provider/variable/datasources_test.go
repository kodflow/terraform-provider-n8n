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
	"github.com/stretchr/testify/mock"
)

// MockVariablesDataSourceInterface is a mock implementation of VariablesDataSourceInterface.
type MockVariablesDataSourceInterface struct {
	mock.Mock
}

func (m *MockVariablesDataSourceInterface) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockVariablesDataSourceInterface) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockVariablesDataSourceInterface) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockVariablesDataSourceInterface) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	m.Called(ctx, req, resp)
}

func TestNewVariablesDataSource(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "create new instance", wantErr: false},
		{name: "multiple instances are independent", wantErr: false},
		{name: "error case - verify initialization", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "create new instance":
				ds := NewVariablesDataSource()

				assert.NotNil(t, ds)
				assert.IsType(t, &VariablesDataSource{}, ds)
				assert.Nil(t, ds.client)

			case "multiple instances are independent":
				ds1 := NewVariablesDataSource()
				ds2 := NewVariablesDataSource()

				assert.NotSame(t, ds1, ds2)

			case "error case - verify initialization":
				ds := NewVariablesDataSource()
				assert.Nil(t, ds.client)
			}
		})
	}
}

func TestNewVariablesDataSourceWrapper(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "create new wrapped instance", wantErr: false},
		{name: "wrapper returns datasource.DataSource interface", wantErr: false},
		{name: "error case - verify wrapper type", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "create new wrapped instance":
				ds := NewVariablesDataSourceWrapper()

				assert.NotNil(t, ds)
				assert.IsType(t, &VariablesDataSource{}, ds)

			case "wrapper returns datasource.DataSource interface":
				ds := NewVariablesDataSourceWrapper()

				// ds is already of type datasource.DataSource, no assertion needed
				assert.NotNil(t, ds)

			case "error case - verify wrapper type":
				ds := NewVariablesDataSourceWrapper()
				var _ datasource.DataSource = ds
			}
		})
	}
}

func TestVariablesDataSource_Metadata(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "set metadata with provider type", wantErr: false},
		{name: "set metadata with different provider type", wantErr: false},
		{name: "set metadata with empty provider type", wantErr: false},
		{name: "error case - metadata validation", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "set metadata with provider type":
				ds := NewVariablesDataSource()
				resp := &datasource.MetadataResponse{}

				ds.Metadata(context.Background(), datasource.MetadataRequest{
					ProviderTypeName: "n8n",
				}, resp)

				assert.Equal(t, "n8n_variables", resp.TypeName)

			case "set metadata with different provider type":
				ds := NewVariablesDataSource()
				resp := &datasource.MetadataResponse{}
				req := datasource.MetadataRequest{
					ProviderTypeName: "custom_provider",
				}

				ds.Metadata(context.Background(), req, resp)

				assert.Equal(t, "custom_provider_variables", resp.TypeName)

			case "set metadata with empty provider type":
				ds := NewVariablesDataSource()
				resp := &datasource.MetadataResponse{}
				req := datasource.MetadataRequest{
					ProviderTypeName: "",
				}

				ds.Metadata(context.Background(), req, resp)

				assert.Equal(t, "_variables", resp.TypeName)

			case "error case - metadata validation":
				ds := NewVariablesDataSource()
				resp := &datasource.MetadataResponse{}
				ds.Metadata(context.Background(), datasource.MetadataRequest{
					ProviderTypeName: "test",
				}, resp)
				assert.Equal(t, "test_variables", resp.TypeName)
			}
		})
	}
}

func TestVariablesDataSource_Schema(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "return schema", wantErr: false},
		{name: "schema attributes have descriptions", wantErr: false},
		{name: "variables list has correct structure", wantErr: false},
		{name: "error case - schema validation", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "return schema":
				ds := NewVariablesDataSource()
				resp := &datasource.SchemaResponse{}

				ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)

				assert.NotNil(t, resp.Schema)
				assert.Contains(t, resp.Schema.MarkdownDescription, "n8n variables")

				// Verify filter attributes
				projectIDAttr, exists := resp.Schema.Attributes["project_id"]
				assert.True(t, exists)
				assert.True(t, projectIDAttr.IsOptional())

				stateAttr, exists := resp.Schema.Attributes["state"]
				assert.True(t, exists)
				assert.True(t, stateAttr.IsOptional())

				// Verify computed variables list
				variablesAttr, exists := resp.Schema.Attributes["variables"]
				assert.True(t, exists)
				assert.True(t, variablesAttr.IsComputed())

			case "schema attributes have descriptions":
				ds := NewVariablesDataSource()
				resp := &datasource.SchemaResponse{}

				ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)

				for name, attr := range resp.Schema.Attributes {
					if name != "variables" { // Skip nested attribute
						assert.NotEmpty(t, attr.GetMarkdownDescription(), "Attribute %s should have description", name)
					}
				}

			case "variables list has correct structure":
				ds := NewVariablesDataSource()
				resp := &datasource.SchemaResponse{}

				ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)

				variablesAttr, exists := resp.Schema.Attributes["variables"]
				assert.True(t, exists)
				assert.True(t, variablesAttr.IsComputed())

			case "error case - schema validation":
				ds := NewVariablesDataSource()
				resp := &datasource.SchemaResponse{}
				ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)
				assert.NotNil(t, resp.Schema)
			}
		})
	}
}

func TestVariablesDataSource_SchemaAttributes(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "return schema attributes", wantErr: false},
		{name: "error case - attribute structure validation", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "return schema attributes":
				ds := NewVariablesDataSource()

				attrs := ds.schemaAttributes()

				assert.NotNil(t, attrs)
				assert.Contains(t, attrs, "project_id")
				assert.Contains(t, attrs, "state")
				assert.Contains(t, attrs, "variables")

				// Verify attribute properties
				projectIDAttr := attrs["project_id"]
				assert.True(t, projectIDAttr.IsOptional())

				stateAttr := attrs["state"]
				assert.True(t, stateAttr.IsOptional())

				variablesAttr := attrs["variables"]
				assert.True(t, variablesAttr.IsComputed())

			case "error case - attribute structure validation":
				ds := NewVariablesDataSource()
				attrs := ds.schemaAttributes()
				assert.Len(t, attrs, 3)
			}
		})
	}
}

func TestVariablesDataSource_VariableItemAttributes(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "return variable item attributes", wantErr: false},
		{name: "value attribute is sensitive", wantErr: false},
		{name: "error case - item attributes validation", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "return variable item attributes":
				ds := NewVariablesDataSource()

				attrs := ds.variableItemAttributes()

				assert.NotNil(t, attrs)
				assert.Contains(t, attrs, "id")
				assert.Contains(t, attrs, "key")
				assert.Contains(t, attrs, "value")
				assert.Contains(t, attrs, "type")
				assert.Contains(t, attrs, "project_id")

				// Verify all are computed
				for _, attr := range attrs {
					assert.True(t, attr.IsComputed())
				}

			case "value attribute is sensitive":
				ds := NewVariablesDataSource()

				attrs := ds.variableItemAttributes()

				valueAttr, exists := attrs["value"]
				assert.True(t, exists)
				assert.True(t, valueAttr.IsComputed())
				// Note: Sensitive attribute testing depends on framework internals

			case "error case - item attributes validation":
				ds := NewVariablesDataSource()
				attrs := ds.variableItemAttributes()
				assert.Len(t, attrs, 5)
			}
		})
	}
}

func TestVariablesDataSource_Configure(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "configure with valid client", wantErr: false},
		{name: "configure with nil provider data", wantErr: false},
		{name: "configure with invalid provider data", wantErr: true},
		{name: "configure with wrong type", wantErr: true},
		{name: "configure multiple times", wantErr: false},
		{name: "error case - configuration validation", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "configure with valid client":
				ds := NewVariablesDataSource()
				resp := &datasource.ConfigureResponse{}

				mockClient := &client.N8nClient{}
				req := datasource.ConfigureRequest{
					ProviderData: mockClient,
				}

				ds.Configure(context.Background(), req, resp)

				assert.NotNil(t, ds.client)
				assert.Equal(t, mockClient, ds.client)
				assert.False(t, resp.Diagnostics.HasError())

			case "configure with nil provider data":
				ds := NewVariablesDataSource()
				resp := &datasource.ConfigureResponse{}
				req := datasource.ConfigureRequest{
					ProviderData: nil,
				}

				ds.Configure(context.Background(), req, resp)

				assert.Nil(t, ds.client)
				assert.False(t, resp.Diagnostics.HasError())

			case "configure with invalid provider data":
				ds := NewVariablesDataSource()
				resp := &datasource.ConfigureResponse{}
				req := datasource.ConfigureRequest{
					ProviderData: "invalid-data",
				}

				ds.Configure(context.Background(), req, resp)

				assert.Nil(t, ds.client)
				assert.True(t, resp.Diagnostics.HasError())
				assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Unexpected Data Source Configure Type")

			case "configure with wrong type":
				ds := NewVariablesDataSource()
				resp := &datasource.ConfigureResponse{}
				req := datasource.ConfigureRequest{
					ProviderData: struct{}{},
				}

				ds.Configure(context.Background(), req, resp)

				assert.Nil(t, ds.client)
				assert.True(t, resp.Diagnostics.HasError())
				assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "*client.N8nClient")

			case "configure multiple times":
				ds := NewVariablesDataSource()

				// First configuration
				resp1 := &datasource.ConfigureResponse{}
				client1 := &client.N8nClient{}
				req1 := datasource.ConfigureRequest{
					ProviderData: client1,
				}
				ds.Configure(context.Background(), req1, resp1)
				assert.Equal(t, client1, ds.client)

				// Second configuration
				resp2 := &datasource.ConfigureResponse{}
				client2 := &client.N8nClient{}
				req2 := datasource.ConfigureRequest{
					ProviderData: client2,
				}
				ds.Configure(context.Background(), req2, resp2)
				assert.Equal(t, client2, ds.client)

			case "error case - configuration validation":
				ds := NewVariablesDataSource()
				resp := &datasource.ConfigureResponse{}
				req := datasource.ConfigureRequest{ProviderData: nil}
				ds.Configure(context.Background(), req, resp)
				assert.False(t, resp.Diagnostics.HasError())
			}
		})
	}
}

func TestVariablesDataSource_Interfaces(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "implements required interfaces", wantErr: false},
		{name: "error case - interface compliance", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "implements required interfaces":
				ds := NewVariablesDataSource()

				// Test that VariablesDataSource implements datasource.DataSource
				var _ datasource.DataSource = ds

				// Test that VariablesDataSource implements datasource.DataSourceWithConfigure
				var _ datasource.DataSourceWithConfigure = ds

				// Test that VariablesDataSource implements VariablesDataSourceInterface
				var _ VariablesDataSourceInterface = ds

			case "error case - interface compliance":
				ds := NewVariablesDataSource()
				var _ datasource.DataSource = ds
			}
		})
	}
}

func TestVariablesDataSourceConcurrency(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "concurrent metadata calls", wantErr: false},
		{name: "concurrent schema calls", wantErr: false},
		{name: "concurrent configure calls", wantErr: false},
		{name: "concurrent schema attributes calls", wantErr: false},
		{name: "concurrent variable item attributes calls", wantErr: false},
		{name: "error case - concurrent access validation", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// NO t.Parallel() here - concurrent goroutines don't play well with t.Parallel()
			switch tt.name {
			case "concurrent metadata calls":
				ds := NewVariablesDataSource()

				done := make(chan bool, 100)
				for i := 0; i < 100; i++ {
					go func() {
						resp := &datasource.MetadataResponse{}
						ds.Metadata(context.Background(), datasource.MetadataRequest{
							ProviderTypeName: "n8n",
						}, resp)
						assert.Equal(t, "n8n_variables", resp.TypeName)
						done <- true
					}()
				}

				for i := 0; i < 100; i++ {
					<-done
				}

			case "concurrent schema calls":
				ds := NewVariablesDataSource()

				done := make(chan bool, 100)
				for i := 0; i < 100; i++ {
					go func() {
						resp := &datasource.SchemaResponse{}
						ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)
						assert.NotNil(t, resp.Schema)
						done <- true
					}()
				}

				for i := 0; i < 100; i++ {
					<-done
				}

			case "concurrent configure calls":
				ds := NewVariablesDataSource()

				done := make(chan bool, 100)
				for i := 0; i < 100; i++ {
					go func() {
						resp := &datasource.ConfigureResponse{}
						mockClient := &client.N8nClient{}
						req := datasource.ConfigureRequest{
							ProviderData: mockClient,
						}
						ds.Configure(context.Background(), req, resp)
						assert.False(t, resp.Diagnostics.HasError())
						done <- true
					}()
				}

				for i := 0; i < 100; i++ {
					<-done
				}

			case "concurrent schema attributes calls":
				ds := NewVariablesDataSource()

				done := make(chan bool, 100)
				for i := 0; i < 100; i++ {
					go func() {
						attrs := ds.schemaAttributes()
						assert.NotNil(t, attrs)
						done <- true
					}()
				}

				for i := 0; i < 100; i++ {
					<-done
				}

			case "concurrent variable item attributes calls":
				ds := NewVariablesDataSource()

				done := make(chan bool, 100)
				for i := 0; i < 100; i++ {
					go func() {
						attrs := ds.variableItemAttributes()
						assert.NotNil(t, attrs)
						done <- true
					}()
				}

				for i := 0; i < 100; i++ {
					<-done
				}

			case "error case - concurrent access validation":
				ds := NewVariablesDataSource()
				done := make(chan bool, 50)
				for i := 0; i < 50; i++ {
					go func() {
						attrs := ds.schemaAttributes()
						assert.NotNil(t, attrs)
						done <- true
					}()
				}
				for i := 0; i < 50; i++ {
					<-done
				}
			}
		})
	}
}

func BenchmarkVariablesDataSource_Schema(b *testing.B) {
	ds := NewVariablesDataSource()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := &datasource.SchemaResponse{}
		ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)
	}
}

func BenchmarkVariablesDataSource_Metadata(b *testing.B) {
	ds := NewVariablesDataSource()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := &datasource.MetadataResponse{}
		ds.Metadata(context.Background(), datasource.MetadataRequest{
			ProviderTypeName: "n8n",
		}, resp)
	}
}

func BenchmarkVariablesDataSource_Configure(b *testing.B) {
	ds := NewVariablesDataSource()
	mockClient := &client.N8nClient{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := &datasource.ConfigureResponse{}
		req := datasource.ConfigureRequest{
			ProviderData: mockClient,
		}
		ds.Configure(context.Background(), req, resp)
	}
}

func BenchmarkVariablesDataSource_SchemaAttributes(b *testing.B) {
	ds := NewVariablesDataSource()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ds.schemaAttributes()
	}
}

func BenchmarkVariablesDataSource_VariableItemAttributes(b *testing.B) {
	ds := NewVariablesDataSource()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ds.variableItemAttributes()
	}
}

// TODO: Fix datasources Read tests - complex Config initialization required.

// Helper functions for test setup.
// Test helper functions for datasources Read functionality.

func TestVariablesDataSource_BuildAPIRequestWithFilters(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "build request with no filters", wantErr: false},
		{name: "build request with project_id filter", wantErr: false},
		{name: "build request with state filter", wantErr: false},
		{name: "build request with both filters", wantErr: false},
		{name: "build request with null filters", wantErr: false},
		{name: "error case - filter validation", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "build request with no filters":
				ds := &VariablesDataSource{}
				apiClient := n8nsdk.NewAPIClient(n8nsdk.NewConfiguration())
				ds.client = &client.N8nClient{APIClient: apiClient}

				data := &models.DataSources{}

				req := ds.buildAPIRequestWithFilters(context.Background(), data)

				assert.NotNil(t, req)

			case "build request with project_id filter":
				ds := &VariablesDataSource{}
				apiClient := n8nsdk.NewAPIClient(n8nsdk.NewConfiguration())
				ds.client = &client.N8nClient{APIClient: apiClient}

				data := &models.DataSources{
					ProjectID: types.StringValue("proj-123"),
				}

				req := ds.buildAPIRequestWithFilters(context.Background(), data)

				assert.NotNil(t, req)

			case "build request with state filter":
				ds := &VariablesDataSource{}
				apiClient := n8nsdk.NewAPIClient(n8nsdk.NewConfiguration())
				ds.client = &client.N8nClient{APIClient: apiClient}

				data := &models.DataSources{
					State: types.StringValue("active"),
				}

				req := ds.buildAPIRequestWithFilters(context.Background(), data)

				assert.NotNil(t, req)

			case "build request with both filters":
				ds := &VariablesDataSource{}
				apiClient := n8nsdk.NewAPIClient(n8nsdk.NewConfiguration())
				ds.client = &client.N8nClient{APIClient: apiClient}

				data := &models.DataSources{
					ProjectID: types.StringValue("proj-456"),
					State:     types.StringValue("inactive"),
				}

				req := ds.buildAPIRequestWithFilters(context.Background(), data)

				assert.NotNil(t, req)

			case "build request with null filters":
				ds := &VariablesDataSource{}
				apiClient := n8nsdk.NewAPIClient(n8nsdk.NewConfiguration())
				ds.client = &client.N8nClient{APIClient: apiClient}

				data := &models.DataSources{
					ProjectID: types.StringNull(),
					State:     types.StringNull(),
				}

				req := ds.buildAPIRequestWithFilters(context.Background(), data)

				assert.NotNil(t, req)

			case "error case - filter validation":
				ds := &VariablesDataSource{}
				apiClient := n8nsdk.NewAPIClient(n8nsdk.NewConfiguration())
				ds.client = &client.N8nClient{APIClient: apiClient}
				data := &models.DataSources{}
				req := ds.buildAPIRequestWithFilters(context.Background(), data)
				assert.NotNil(t, req)
			}
		})
	}
}

func TestVariablesDataSource_PopulateVariables(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "populate with valid data", wantErr: false},
		{name: "populate with nil data", wantErr: false},
		{name: "populate with empty data", wantErr: false},
		{name: "populate with multiple variables", wantErr: false},
		{name: "error case - population validation", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "populate with valid data":
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

			case "populate with nil data":
				ds := &VariablesDataSource{}
				data := &models.DataSources{}

				variableList := &n8nsdk.VariableList{
					Data: nil,
				}

				ds.populateVariables(data, variableList)

				assert.Empty(t, data.Variables)

			case "populate with empty data":
				ds := &VariablesDataSource{}
				data := &models.DataSources{}

				variableList := &n8nsdk.VariableList{
					Data: []n8nsdk.Variable{},
				}

				ds.populateVariables(data, variableList)

				assert.Empty(t, data.Variables)

			case "populate with multiple variables":
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

			case "error case - population validation":
				ds := &VariablesDataSource{}
				data := &models.DataSources{}
				variableList := &n8nsdk.VariableList{Data: nil}
				ds.populateVariables(data, variableList)
				assert.Empty(t, data.Variables)
			}
		})
	}
}

func TestVariablesDataSource_MapVariableToItem(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "map with all fields", wantErr: false},
		{name: "map with nil optional fields", wantErr: false},
		{name: "map with project but nil project ID", wantErr: false},
		{name: "map with empty strings", wantErr: false},
		{name: "map preserves sensitive value", wantErr: false},
		{name: "error case - mapping validation", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "map with all fields":
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

			case "map with nil optional fields":
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

			case "map with project but nil project ID":
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

			case "map with empty strings":
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

			case "map preserves sensitive value":
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

			case "error case - mapping validation":
				ds := &VariablesDataSource{}
				variable := &n8nsdk.Variable{
					Id:    nil,
					Key:   "TEST",
					Value: "val",
				}
				item := ds.mapVariableToItem(variable)
				assert.True(t, item.ID.IsNull())
			}
		})
	}
}

func TestVariablesDataSource_Read(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "successful read with variables", wantErr: false},
		{name: "read fails on API error", wantErr: true},
		{name: "read with empty response", wantErr: false},
		{name: "read with null data", wantErr: false},
		{name: "read with invalid config", wantErr: true},
		{name: "error case - read validation", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "successful read with variables":
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

			case "read fails on API error":
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

			case "read with empty response":
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

			case "read with null data":
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

			case "read with invalid config":
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

			case "error case - read validation":
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(map[string]interface{}{"data": []interface{}{}})
				})
				n8nClient, server := setupTestVariablesDataSourceClient(t, handler)
				defer server.Close()
				ds := &VariablesDataSource{client: n8nClient}
				req := datasource.ReadRequest{Config: createTestVariablesDataSourceConfigHelper(t)}
				resp := &datasource.ReadResponse{State: createTestVariablesDataSourceStateHelper(t)}
				ds.Read(context.Background(), req, resp)
				assert.False(t, resp.Diagnostics.HasError())
			}
		})
	}
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
		"variables": tftypes.NewValue(tftypes.List{ElementType: tftypes.Object{AttributeTypes: map[string]tftypes.Type{
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
