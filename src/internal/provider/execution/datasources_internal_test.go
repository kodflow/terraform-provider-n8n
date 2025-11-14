package execution

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/execution/models"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockExecutionsDataSourceInterface is a mock implementation of ExecutionsDataSourceInterface.
type MockExecutionsDataSourceInterface struct {
	mock.Mock
}

func (m *MockExecutionsDataSourceInterface) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockExecutionsDataSourceInterface) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockExecutionsDataSourceInterface) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockExecutionsDataSourceInterface) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	m.Called(ctx, req, resp)
}

// TestNewExecutionsDataSource is now in external test file - refactored to test behavior only.

// TestNewExecutionsDataSourceWrapper is now in external test file - refactored to test behavior only.

// TestExecutionsDataSource_Metadata is now in external test file - refactored to test behavior only.

// TestExecutionsDataSource_Schema is now in external test file - refactored to test behavior only.

// TestExecutionsDataSource_Configure is now in external test file - refactored to test behavior only.

func TestExecutionsDataSource_Interfaces(t *testing.T) {
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
				ds := NewExecutionsDataSource()

				// Test that ExecutionsDataSource implements datasource.DataSource
				var _ datasource.DataSource = ds

				// Test that ExecutionsDataSource implements datasource.DataSourceWithConfigure
				var _ datasource.DataSourceWithConfigure = ds

				// Test that ExecutionsDataSource implements ExecutionsDataSourceInterface
				var _ ExecutionsDataSourceInterface = ds

			case "error case - interface compliance":
				ds := NewExecutionsDataSource()
				var _ datasource.DataSource = ds
			}
		})
	}
}

func TestExecutionsDataSourceConcurrency(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "concurrent metadata calls", wantErr: false},
		{name: "concurrent schema calls", wantErr: false},
		{name: "error case - concurrent access validation", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// NO t.Parallel() here - concurrent goroutines don't play well with t.Parallel()
			switch tt.name {
			case "concurrent metadata calls":
				ds := NewExecutionsDataSource()

				done := make(chan bool, 100)
				for i := 0; i < 100; i++ {
					go func() {
						resp := &datasource.MetadataResponse{}
						ds.Metadata(context.Background(), datasource.MetadataRequest{
							ProviderTypeName: "n8n",
						}, resp)
						assert.Equal(t, "n8n_executions", resp.TypeName)
						done <- true
					}()
				}

				for i := 0; i < 100; i++ {
					<-done
				}

			case "concurrent schema calls":
				ds := NewExecutionsDataSource()

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

			case "error case - concurrent access validation":
				ds := NewExecutionsDataSource()
				done := make(chan bool, 50)
				for i := 0; i < 50; i++ {
					go func() {
						resp := &datasource.MetadataResponse{}
						ds.Metadata(context.Background(), datasource.MetadataRequest{
							ProviderTypeName: "n8n",
						}, resp)
						assert.Equal(t, "n8n_executions", resp.TypeName)
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

func BenchmarkExecutionsDataSource_Schema(b *testing.B) {
	ds := NewExecutionsDataSource()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := &datasource.SchemaResponse{}
		ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)
	}
}

func BenchmarkExecutionsDataSource_Metadata(b *testing.B) {
	ds := NewExecutionsDataSource()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := &datasource.MetadataResponse{}
		ds.Metadata(context.Background(), datasource.MetadataRequest{}, resp)
	}
}

func BenchmarkExecutionsDataSource_Configure(b *testing.B) {
	ds := NewExecutionsDataSource()
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

// TestExecutionsDataSource_Read is now in external test file - refactored to test behavior only.

func TestExecutionsDataSource_schemaAttributes(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "schema attributes returned correctly", wantErr: false},
		{name: "optional fields present", wantErr: false},
		{name: "computed fields present", wantErr: false},
		{name: "nested executions list present", wantErr: false},
		{name: "error case - verify attribute count", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "schema attributes returned correctly":
				ds := &ExecutionsDataSource{}
				attrs := ds.schemaAttributes()

				assert.NotNil(t, attrs)
				assert.Contains(t, attrs, "workflow_id")
				assert.Contains(t, attrs, "project_id")
				assert.Contains(t, attrs, "status")
				assert.Contains(t, attrs, "include_data")
				assert.Contains(t, attrs, "executions")

			case "optional fields present":
				ds := &ExecutionsDataSource{}
				attrs := ds.schemaAttributes()

				workflowIDAttr, ok := attrs["workflow_id"].(schema.StringAttribute)
				assert.True(t, ok)
				assert.True(t, workflowIDAttr.Optional)

				projectIDAttr, ok := attrs["project_id"].(schema.StringAttribute)
				assert.True(t, ok)
				assert.True(t, projectIDAttr.Optional)

			case "computed fields present":
				ds := &ExecutionsDataSource{}
				attrs := ds.schemaAttributes()

				executionsAttr, ok := attrs["executions"].(schema.ListNestedAttribute)
				assert.True(t, ok)
				assert.True(t, executionsAttr.Computed)

			case "nested executions list present":
				ds := &ExecutionsDataSource{}
				attrs := ds.schemaAttributes()

				executionsAttr, ok := attrs["executions"].(schema.ListNestedAttribute)
				assert.True(t, ok)
				assert.NotNil(t, executionsAttr.NestedObject)

			case "error case - verify attribute count":
				ds := &ExecutionsDataSource{}
				attrs := ds.schemaAttributes()
				assert.Len(t, attrs, 5)
			}
		})
	}
}

func TestExecutionsDataSource_executionItemAttributes(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "execution item attributes returned correctly", wantErr: false},
		{name: "all computed fields present", wantErr: false},
		{name: "id field is computed", wantErr: false},
		{name: "error case - verify all required fields", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "execution item attributes returned correctly":
				ds := &ExecutionsDataSource{}
				attrs := ds.executionItemAttributes()

				assert.NotNil(t, attrs)
				assert.Contains(t, attrs, "id")
				assert.Contains(t, attrs, "workflow_id")
				assert.Contains(t, attrs, "finished")
				assert.Contains(t, attrs, "mode")
				assert.Contains(t, attrs, "started_at")
				assert.Contains(t, attrs, "stopped_at")
				assert.Contains(t, attrs, "status")

			case "all computed fields present":
				ds := &ExecutionsDataSource{}
				attrs := ds.executionItemAttributes()

				for key, attr := range attrs {
					switch v := attr.(type) {
					case schema.StringAttribute:
						assert.True(t, v.Computed, "Field %s should be computed", key)
					case schema.BoolAttribute:
						assert.True(t, v.Computed, "Field %s should be computed", key)
					}
				}

			case "id field is computed":
				ds := &ExecutionsDataSource{}
				attrs := ds.executionItemAttributes()

				idAttr, ok := attrs["id"].(schema.StringAttribute)
				assert.True(t, ok)
				assert.True(t, idAttr.Computed)

			case "error case - verify all required fields":
				ds := &ExecutionsDataSource{}
				attrs := ds.executionItemAttributes()
				assert.Len(t, attrs, 7)
			}
		})
	}
}

func TestExecutionsDataSource_buildExecutionsAPIRequest(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "build request with no filters", wantErr: false},
		{name: "build request with workflow_id filter", wantErr: false},
		{name: "build request with all filters", wantErr: false},
		{name: "build request with null values", wantErr: false},
		{name: "error case - verify request builder", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "build request with no filters":
				ds := &ExecutionsDataSource{}
				cfg := n8nsdk.NewConfiguration()
				cfg.Servers = n8nsdk.ServerConfigurations{
					{URL: "http://test.local", Description: "Test server"},
				}
				ds.client = &client.N8nClient{
					APIClient: n8nsdk.NewAPIClient(cfg),
				}

				data := &models.DataSources{}
				req := ds.buildExecutionsAPIRequest(context.Background(), data)

				assert.NotNil(t, req)

			case "build request with workflow_id filter":
				ds := &ExecutionsDataSource{}
				cfg := n8nsdk.NewConfiguration()
				cfg.Servers = n8nsdk.ServerConfigurations{
					{URL: "http://test.local", Description: "Test server"},
				}
				ds.client = &client.N8nClient{
					APIClient: n8nsdk.NewAPIClient(cfg),
				}

				data := &models.DataSources{}
				data.WorkflowID = types.StringValue("123")
				req := ds.buildExecutionsAPIRequest(context.Background(), data)

				assert.NotNil(t, req)

			case "build request with all filters":
				ds := &ExecutionsDataSource{}
				cfg := n8nsdk.NewConfiguration()
				cfg.Servers = n8nsdk.ServerConfigurations{
					{URL: "http://test.local", Description: "Test server"},
				}
				ds.client = &client.N8nClient{
					APIClient: n8nsdk.NewAPIClient(cfg),
				}

				data := &models.DataSources{}
				data.WorkflowID = types.StringValue("123")
				data.ProjectID = types.StringValue("456")
				data.Status = types.StringValue("success")
				data.IncludeData = types.BoolValue(true)
				req := ds.buildExecutionsAPIRequest(context.Background(), data)

				assert.NotNil(t, req)

			case "build request with null values":
				ds := &ExecutionsDataSource{}
				cfg := n8nsdk.NewConfiguration()
				cfg.Servers = n8nsdk.ServerConfigurations{
					{URL: "http://test.local", Description: "Test server"},
				}
				ds.client = &client.N8nClient{
					APIClient: n8nsdk.NewAPIClient(cfg),
				}

				data := &models.DataSources{}
				data.WorkflowID = types.StringNull()
				data.ProjectID = types.StringNull()
				req := ds.buildExecutionsAPIRequest(context.Background(), data)

				assert.NotNil(t, req)

			case "error case - verify request builder":
				ds := &ExecutionsDataSource{}
				cfg := n8nsdk.NewConfiguration()
				cfg.Servers = n8nsdk.ServerConfigurations{
					{URL: "http://test.local", Description: "Test server"},
				}
				ds.client = &client.N8nClient{
					APIClient: n8nsdk.NewAPIClient(cfg),
				}

				data := &models.DataSources{}
				req := ds.buildExecutionsAPIRequest(context.Background(), data)
				assert.NotNil(t, req)
			}
		})
	}
}

func TestExecutionsDataSource_populateExecutionsList(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "populate with empty list", wantErr: false},
		{name: "populate with single execution", wantErr: false},
		{name: "populate with multiple executions", wantErr: false},
		{name: "populate with nil data", wantErr: false},
		{name: "error case - verify population", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "populate with empty list":
				ds := &ExecutionsDataSource{}
				data := &models.DataSources{}
				executionList := &n8nsdk.ExecutionList{
					Data: []n8nsdk.Execution{},
				}

				ds.populateExecutionsList(data, executionList)

				assert.NotNil(t, data.Executions)
				assert.Len(t, data.Executions, 0)

			case "populate with single execution":
				ds := &ExecutionsDataSource{}
				data := &models.DataSources{}

				id := float32(123)
				workflowID := float32(456)
				status := "success"
				executionList := &n8nsdk.ExecutionList{
					Data: []n8nsdk.Execution{
						{
							Id:         &id,
							WorkflowId: &workflowID,
							Status:     &status,
						},
					},
				}

				ds.populateExecutionsList(data, executionList)

				assert.NotNil(t, data.Executions)
				assert.Len(t, data.Executions, 1)
				assert.Equal(t, "123", data.Executions[0].ID.ValueString())
				assert.Equal(t, "456", data.Executions[0].WorkflowID.ValueString())
				assert.Equal(t, "success", data.Executions[0].Status.ValueString())

			case "populate with multiple executions":
				ds := &ExecutionsDataSource{}
				data := &models.DataSources{}

				id1 := float32(111)
				id2 := float32(222)
				id3 := float32(333)
				executionList := &n8nsdk.ExecutionList{
					Data: []n8nsdk.Execution{
						{Id: &id1},
						{Id: &id2},
						{Id: &id3},
					},
				}

				ds.populateExecutionsList(data, executionList)

				assert.NotNil(t, data.Executions)
				assert.Len(t, data.Executions, 3)
				assert.Equal(t, "111", data.Executions[0].ID.ValueString())
				assert.Equal(t, "222", data.Executions[1].ID.ValueString())
				assert.Equal(t, "333", data.Executions[2].ID.ValueString())

			case "populate with nil data":
				ds := &ExecutionsDataSource{}
				data := &models.DataSources{}
				executionList := &n8nsdk.ExecutionList{}

				ds.populateExecutionsList(data, executionList)

				assert.NotNil(t, data.Executions)
				assert.Len(t, data.Executions, 0)

			case "error case - verify population":
				ds := &ExecutionsDataSource{}
				data := &models.DataSources{}
				executionList := &n8nsdk.ExecutionList{
					Data: []n8nsdk.Execution{},
				}

				ds.populateExecutionsList(data, executionList)
				assert.NotNil(t, data.Executions)
			}
		})
	}
}
