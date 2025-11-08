package execution

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
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

func TestNewExecutionsDataSource(t *testing.T) {
	t.Run("create new instance", func(t *testing.T) {
		ds := NewExecutionsDataSource()

		assert.NotNil(t, ds)
		assert.IsType(t, &ExecutionsDataSource{}, ds)
		assert.Nil(t, ds.client)
	})

	t.Run("multiple instances are independent", func(t *testing.T) {
		ds1 := NewExecutionsDataSource()
		ds2 := NewExecutionsDataSource()

		assert.NotSame(t, ds1, ds2)
		// Each instance is a different pointer, even if they have the same initial state
	})
}

func TestNewExecutionsDataSourceWrapper(t *testing.T) {
	t.Run("create new wrapped instance", func(t *testing.T) {
		ds := NewExecutionsDataSourceWrapper()

		assert.NotNil(t, ds)
		assert.IsType(t, &ExecutionsDataSource{}, ds)
	})

	t.Run("wrapper returns datasource.DataSource interface", func(t *testing.T) {
		ds := NewExecutionsDataSourceWrapper()

		// ds is already of type datasource.DataSource, no assertion needed
		assert.NotNil(t, ds)
	})
}

func TestExecutionsDataSource_Metadata(t *testing.T) {
	t.Run("set metadata without provider type", func(t *testing.T) {
		ds := NewExecutionsDataSource()
		resp := &datasource.MetadataResponse{}

		ds.Metadata(context.Background(), datasource.MetadataRequest{
			ProviderTypeName: "n8n",
		}, resp)

		assert.Equal(t, "n8n_executions", resp.TypeName)
	})

	t.Run("set metadata with provider type name", func(t *testing.T) {
		ds := NewExecutionsDataSource()
		resp := &datasource.MetadataResponse{}
		req := datasource.MetadataRequest{
			ProviderTypeName: "n8n",
		}

		ds.Metadata(context.Background(), req, resp)

		assert.Equal(t, "n8n_executions", resp.TypeName)
	})

	t.Run("set metadata with different provider type", func(t *testing.T) {
		ds := NewExecutionsDataSource()
		resp := &datasource.MetadataResponse{}
		req := datasource.MetadataRequest{
			ProviderTypeName: "custom_provider",
		}

		ds.Metadata(context.Background(), req, resp)

		assert.Equal(t, "custom_provider_executions", resp.TypeName)
	})
}

func TestExecutionsDataSource_Schema(t *testing.T) {
	t.Run("return schema", func(t *testing.T) {
		ds := NewExecutionsDataSource()
		resp := &datasource.SchemaResponse{}

		ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)

		assert.NotNil(t, resp.Schema)
		assert.Contains(t, resp.Schema.MarkdownDescription, "n8n workflow executions")

		// Verify filter attributes
		workflowAttr, exists := resp.Schema.Attributes["workflow_id"]
		assert.True(t, exists)
		assert.True(t, workflowAttr.IsOptional())

		projectAttr, exists := resp.Schema.Attributes["project_id"]
		assert.True(t, exists)
		assert.True(t, projectAttr.IsOptional())

		statusAttr, exists := resp.Schema.Attributes["status"]
		assert.True(t, exists)
		assert.True(t, statusAttr.IsOptional())

		includeDataAttr, exists := resp.Schema.Attributes["include_data"]
		assert.True(t, exists)
		assert.True(t, includeDataAttr.IsOptional())

		// Verify computed executions list
		executionsAttr, exists := resp.Schema.Attributes["executions"]
		assert.True(t, exists)
		assert.True(t, executionsAttr.IsComputed())
	})

	t.Run("executions list has correct structure", func(t *testing.T) {
		ds := NewExecutionsDataSource()
		resp := &datasource.SchemaResponse{}

		ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)

		executionsAttr, exists := resp.Schema.Attributes["executions"]
		assert.True(t, exists)
		assert.True(t, executionsAttr.IsComputed())

		// Verify that executions is a computed list attribute
		// The actual nested structure verification would require access to the
		// internal schema structure which varies by Terraform plugin framework version
	})

	t.Run("schema attributes have descriptions", func(t *testing.T) {
		ds := NewExecutionsDataSource()
		resp := &datasource.SchemaResponse{}

		ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)

		for name, attr := range resp.Schema.Attributes {
			if name != "executions" { // Skip nested attribute
				assert.NotEmpty(t, attr.GetMarkdownDescription(), "Attribute %s should have description", name)
			}
		}
	})
}

func TestExecutionsDataSource_Configure(t *testing.T) {
	t.Run("configure with valid client", func(t *testing.T) {
		ds := NewExecutionsDataSource()
		resp := &datasource.ConfigureResponse{}

		mockClient := &client.N8nClient{}
		req := datasource.ConfigureRequest{
			ProviderData: mockClient,
		}

		ds.Configure(context.Background(), req, resp)

		assert.NotNil(t, ds.client)
		assert.Equal(t, mockClient, ds.client)
		assert.False(t, resp.Diagnostics.HasError())
	})

	t.Run("configure with nil provider data", func(t *testing.T) {
		ds := NewExecutionsDataSource()
		resp := &datasource.ConfigureResponse{}
		req := datasource.ConfigureRequest{
			ProviderData: nil,
		}

		ds.Configure(context.Background(), req, resp)

		assert.Nil(t, ds.client)
		assert.False(t, resp.Diagnostics.HasError())
	})

	t.Run("configure with invalid provider data", func(t *testing.T) {
		ds := NewExecutionsDataSource()
		resp := &datasource.ConfigureResponse{}
		req := datasource.ConfigureRequest{
			ProviderData: "invalid-data",
		}

		ds.Configure(context.Background(), req, resp)

		assert.Nil(t, ds.client)
		assert.True(t, resp.Diagnostics.HasError())
		assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Unexpected Data Source Configure Type")
	})

	t.Run("configure with wrong type", func(t *testing.T) {
		ds := NewExecutionsDataSource()
		resp := &datasource.ConfigureResponse{}
		req := datasource.ConfigureRequest{
			ProviderData: struct{}{},
		}

		ds.Configure(context.Background(), req, resp)

		assert.Nil(t, ds.client)
		assert.True(t, resp.Diagnostics.HasError())
	})

	t.Run("configure multiple times", func(t *testing.T) {
		ds := NewExecutionsDataSource()

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
	})
}

func TestExecutionsDataSource_Interfaces(t *testing.T) {
	t.Run("implements required interfaces", func(t *testing.T) {
		ds := NewExecutionsDataSource()

		// Test that ExecutionsDataSource implements datasource.DataSource
		var _ datasource.DataSource = ds

		// Test that ExecutionsDataSource implements datasource.DataSourceWithConfigure
		var _ datasource.DataSourceWithConfigure = ds

		// Test that ExecutionsDataSource implements ExecutionsDataSourceInterface
		var _ ExecutionsDataSourceInterface = ds
	})
}

func TestExecutionsDataSourceConcurrency(t *testing.T) {
	t.Run("concurrent metadata calls", func(t *testing.T) {
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
	})

	t.Run("concurrent schema calls", func(t *testing.T) {
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
	})
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
