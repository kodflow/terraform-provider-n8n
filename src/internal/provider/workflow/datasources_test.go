package workflow

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockWorkflowsDataSourceInterface is a mock implementation of WorkflowsDataSourceInterface
type MockWorkflowsDataSourceInterface struct {
	mock.Mock
}

func (m *MockWorkflowsDataSourceInterface) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockWorkflowsDataSourceInterface) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockWorkflowsDataSourceInterface) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockWorkflowsDataSourceInterface) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	m.Called(ctx, req, resp)
}

func TestNewWorkflowsDataSource(t *testing.T) {
	t.Run("create new instance", func(t *testing.T) {
		ds := NewWorkflowsDataSource()

		assert.NotNil(t, ds)
		assert.IsType(t, &WorkflowsDataSource{}, ds)
		assert.Nil(t, ds.client)
	})

	t.Run("multiple instances are independent", func(t *testing.T) {
		ds1 := NewWorkflowsDataSource()
		ds2 := NewWorkflowsDataSource()

		assert.NotSame(t, ds1, ds2)
	})
}

func TestNewWorkflowsDataSourceWrapper(t *testing.T) {
	t.Run("create new wrapped instance", func(t *testing.T) {
		ds := NewWorkflowsDataSourceWrapper()

		assert.NotNil(t, ds)
		assert.IsType(t, &WorkflowsDataSource{}, ds)
	})

	t.Run("wrapper returns datasource.DataSource interface", func(t *testing.T) {
		ds := NewWorkflowsDataSourceWrapper()

		_, ok := ds.(datasource.DataSource)
		assert.True(t, ok)
	})
}

func TestWorkflowsDataSource_Metadata(t *testing.T) {
	t.Run("set metadata with provider type", func(t *testing.T) {
		ds := NewWorkflowsDataSource()
		resp := &datasource.MetadataResponse{}

		ds.Metadata(context.Background(), datasource.MetadataRequest{
			ProviderTypeName: "n8n",
		}, resp)

		assert.Equal(t, "n8n_workflows", resp.TypeName)
	})

	t.Run("set metadata with different provider type", func(t *testing.T) {
		ds := NewWorkflowsDataSource()
		resp := &datasource.MetadataResponse{}
		req := datasource.MetadataRequest{
			ProviderTypeName: "custom_provider",
		}

		ds.Metadata(context.Background(), req, resp)

		assert.Equal(t, "custom_provider_workflows", resp.TypeName)
	})

	t.Run("set metadata with empty provider type", func(t *testing.T) {
		ds := NewWorkflowsDataSource()
		resp := &datasource.MetadataResponse{}
		req := datasource.MetadataRequest{
			ProviderTypeName: "",
		}

		ds.Metadata(context.Background(), req, resp)

		assert.Equal(t, "_workflows", resp.TypeName)
	})
}

func TestWorkflowsDataSource_Schema(t *testing.T) {
	t.Run("return schema", func(t *testing.T) {
		ds := NewWorkflowsDataSource()
		resp := &datasource.SchemaResponse{}

		ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)

		assert.NotNil(t, resp.Schema)
		assert.Contains(t, resp.Schema.MarkdownDescription, "n8n workflows")

		// Verify filter attributes
		activeAttr, exists := resp.Schema.Attributes["active"]
		assert.True(t, exists)
		assert.True(t, activeAttr.IsOptional())

		// Verify computed workflows list
		workflowsAttr, exists := resp.Schema.Attributes["workflows"]
		assert.True(t, exists)
		assert.True(t, workflowsAttr.IsComputed())
	})

	t.Run("workflows list has correct structure", func(t *testing.T) {
		ds := NewWorkflowsDataSource()
		resp := &datasource.SchemaResponse{}

		ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)

		workflowsAttr, exists := resp.Schema.Attributes["workflows"]
		assert.True(t, exists)
		assert.True(t, workflowsAttr.IsComputed())
	})

	t.Run("schema attributes have descriptions", func(t *testing.T) {
		ds := NewWorkflowsDataSource()
		resp := &datasource.SchemaResponse{}

		ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)

		for name, attr := range resp.Schema.Attributes {
			assert.NotEmpty(t, attr.GetMarkdownDescription(), "Attribute %s should have description", name)
		}
	})

	t.Run("schema has correct attribute count", func(t *testing.T) {
		ds := NewWorkflowsDataSource()
		resp := &datasource.SchemaResponse{}

		ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)

		// Should have active and workflows
		assert.Len(t, resp.Schema.Attributes, 2)
	})
}

func TestWorkflowsDataSource_Configure(t *testing.T) {
	t.Run("configure with valid client", func(t *testing.T) {
		ds := NewWorkflowsDataSource()
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
		ds := NewWorkflowsDataSource()
		resp := &datasource.ConfigureResponse{}
		req := datasource.ConfigureRequest{
			ProviderData: nil,
		}

		ds.Configure(context.Background(), req, resp)

		assert.Nil(t, ds.client)
		assert.False(t, resp.Diagnostics.HasError())
	})

	t.Run("configure with invalid provider data", func(t *testing.T) {
		ds := NewWorkflowsDataSource()
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
		ds := NewWorkflowsDataSource()
		resp := &datasource.ConfigureResponse{}
		req := datasource.ConfigureRequest{
			ProviderData: struct{}{},
		}

		ds.Configure(context.Background(), req, resp)

		assert.Nil(t, ds.client)
		assert.True(t, resp.Diagnostics.HasError())
	})

	t.Run("configure multiple times", func(t *testing.T) {
		ds := NewWorkflowsDataSource()

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

func TestWorkflowsDataSource_Interfaces(t *testing.T) {
	t.Run("implements required interfaces", func(t *testing.T) {
		ds := NewWorkflowsDataSource()

		// Test that WorkflowsDataSource implements datasource.DataSource
		var _ datasource.DataSource = ds

		// Test that WorkflowsDataSource implements datasource.DataSourceWithConfigure
		var _ datasource.DataSourceWithConfigure = ds

		// Test that WorkflowsDataSource implements WorkflowsDataSourceInterface
		var _ WorkflowsDataSourceInterface = ds
	})
}

func TestWorkflowsDataSourceConcurrency(t *testing.T) {
	t.Run("concurrent metadata calls", func(t *testing.T) {
		ds := NewWorkflowsDataSource()

		done := make(chan bool, 100)
		for i := 0; i < 100; i++ {
			go func() {
				resp := &datasource.MetadataResponse{}
				ds.Metadata(context.Background(), datasource.MetadataRequest{
					ProviderTypeName: "n8n",
				}, resp)
				assert.Equal(t, "n8n_workflows", resp.TypeName)
				done <- true
			}()
		}

		for i := 0; i < 100; i++ {
			<-done
		}
	})

	t.Run("concurrent schema calls", func(t *testing.T) {
		ds := NewWorkflowsDataSource()

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

func BenchmarkWorkflowsDataSource_Schema(b *testing.B) {
	ds := NewWorkflowsDataSource()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := &datasource.SchemaResponse{}
		ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)
	}
}

func BenchmarkWorkflowsDataSource_Metadata(b *testing.B) {
	ds := NewWorkflowsDataSource()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := &datasource.MetadataResponse{}
		ds.Metadata(context.Background(), datasource.MetadataRequest{}, resp)
	}
}

func BenchmarkWorkflowsDataSource_Configure(b *testing.B) {
	ds := NewWorkflowsDataSource()
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
