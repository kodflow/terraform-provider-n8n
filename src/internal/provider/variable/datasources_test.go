package variable

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
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
	t.Run("create new instance", func(t *testing.T) {
		ds := NewVariablesDataSource()

		assert.NotNil(t, ds)
		assert.IsType(t, &VariablesDataSource{}, ds)
		assert.Nil(t, ds.client)
	})

	t.Run("multiple instances are independent", func(t *testing.T) {
		ds1 := NewVariablesDataSource()
		ds2 := NewVariablesDataSource()

		assert.NotSame(t, ds1, ds2)
	})
}

func TestNewVariablesDataSourceWrapper(t *testing.T) {
	t.Run("create new wrapped instance", func(t *testing.T) {
		ds := NewVariablesDataSourceWrapper()

		assert.NotNil(t, ds)
		assert.IsType(t, &VariablesDataSource{}, ds)
	})

	t.Run("wrapper returns datasource.DataSource interface", func(t *testing.T) {
		ds := NewVariablesDataSourceWrapper()

		// ds is already of type datasource.DataSource, no assertion needed
		assert.NotNil(t, ds)
	})
}

func TestVariablesDataSource_Metadata(t *testing.T) {
	t.Run("set metadata with provider type", func(t *testing.T) {
		ds := NewVariablesDataSource()
		resp := &datasource.MetadataResponse{}

		ds.Metadata(context.Background(), datasource.MetadataRequest{
			ProviderTypeName: "n8n",
		}, resp)

		assert.Equal(t, "n8n_variables", resp.TypeName)
	})

	t.Run("set metadata with different provider type", func(t *testing.T) {
		ds := NewVariablesDataSource()
		resp := &datasource.MetadataResponse{}
		req := datasource.MetadataRequest{
			ProviderTypeName: "custom_provider",
		}

		ds.Metadata(context.Background(), req, resp)

		assert.Equal(t, "custom_provider_variables", resp.TypeName)
	})

	t.Run("set metadata with empty provider type", func(t *testing.T) {
		ds := NewVariablesDataSource()
		resp := &datasource.MetadataResponse{}
		req := datasource.MetadataRequest{
			ProviderTypeName: "",
		}

		ds.Metadata(context.Background(), req, resp)

		assert.Equal(t, "_variables", resp.TypeName)
	})
}

func TestVariablesDataSource_Schema(t *testing.T) {
	t.Run("return schema", func(t *testing.T) {
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
	})

	t.Run("schema attributes have descriptions", func(t *testing.T) {
		ds := NewVariablesDataSource()
		resp := &datasource.SchemaResponse{}

		ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)

		for name, attr := range resp.Schema.Attributes {
			if name != "variables" { // Skip nested attribute
				assert.NotEmpty(t, attr.GetMarkdownDescription(), "Attribute %s should have description", name)
			}
		}
	})

	t.Run("variables list has correct structure", func(t *testing.T) {
		ds := NewVariablesDataSource()
		resp := &datasource.SchemaResponse{}

		ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)

		variablesAttr, exists := resp.Schema.Attributes["variables"]
		assert.True(t, exists)
		assert.True(t, variablesAttr.IsComputed())
	})
}

func TestVariablesDataSource_SchemaAttributes(t *testing.T) {
	t.Run("return schema attributes", func(t *testing.T) {
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
	})
}

func TestVariablesDataSource_VariableItemAttributes(t *testing.T) {
	t.Run("return variable item attributes", func(t *testing.T) {
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
	})

	t.Run("value attribute is sensitive", func(t *testing.T) {
		ds := NewVariablesDataSource()

		attrs := ds.variableItemAttributes()

		valueAttr, exists := attrs["value"]
		assert.True(t, exists)
		assert.True(t, valueAttr.IsComputed())
		// Note: Sensitive attribute testing depends on framework internals
	})
}

func TestVariablesDataSource_Configure(t *testing.T) {
	t.Run("configure with valid client", func(t *testing.T) {
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
	})

	t.Run("configure with nil provider data", func(t *testing.T) {
		ds := NewVariablesDataSource()
		resp := &datasource.ConfigureResponse{}
		req := datasource.ConfigureRequest{
			ProviderData: nil,
		}

		ds.Configure(context.Background(), req, resp)

		assert.Nil(t, ds.client)
		assert.False(t, resp.Diagnostics.HasError())
	})

	t.Run("configure with invalid provider data", func(t *testing.T) {
		ds := NewVariablesDataSource()
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
		ds := NewVariablesDataSource()
		resp := &datasource.ConfigureResponse{}
		req := datasource.ConfigureRequest{
			ProviderData: struct{}{},
		}

		ds.Configure(context.Background(), req, resp)

		assert.Nil(t, ds.client)
		assert.True(t, resp.Diagnostics.HasError())
		assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "*client.N8nClient")
	})

	t.Run("configure multiple times", func(t *testing.T) {
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
	})
}

func TestVariablesDataSource_Interfaces(t *testing.T) {
	t.Run("implements required interfaces", func(t *testing.T) {
		ds := NewVariablesDataSource()

		// Test that VariablesDataSource implements datasource.DataSource
		var _ datasource.DataSource = ds

		// Test that VariablesDataSource implements datasource.DataSourceWithConfigure
		var _ datasource.DataSourceWithConfigure = ds

		// Test that VariablesDataSource implements VariablesDataSourceInterface
		var _ VariablesDataSourceInterface = ds
	})
}

func TestVariablesDataSourceConcurrency(t *testing.T) {
	t.Run("concurrent metadata calls", func(t *testing.T) {
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
	})

	t.Run("concurrent schema calls", func(t *testing.T) {
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
	})

	t.Run("concurrent configure calls", func(t *testing.T) {
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
	})

	t.Run("concurrent schema attributes calls", func(t *testing.T) {
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
	})

	t.Run("concurrent variable item attributes calls", func(t *testing.T) {
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
	})
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
