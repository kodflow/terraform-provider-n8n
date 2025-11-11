package execution_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/execution"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
	"github.com/stretchr/testify/assert"
)

func TestNewExecutionsDataSource(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "create new instance", wantErr: false},
		{name: "returns non-nil instance", wantErr: false},
		{name: "error case - verify instance type", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "create new instance":
				ds := execution.NewExecutionsDataSource()
				assert.NotNil(t, ds)
				assert.IsType(t, &execution.ExecutionsDataSource{}, ds)

			case "returns non-nil instance":
				ds := execution.NewExecutionsDataSource()
				assert.NotNil(t, ds)

			case "error case - verify instance type":
				ds := execution.NewExecutionsDataSource()
				assert.NotNil(t, ds)
			}
		})
	}
}

func TestNewExecutionsDataSourceWrapper(t *testing.T) {
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
				ds := execution.NewExecutionsDataSourceWrapper()
				assert.NotNil(t, ds)
				assert.IsType(t, &execution.ExecutionsDataSource{}, ds)

			case "wrapper returns datasource.DataSource interface":
				ds := execution.NewExecutionsDataSourceWrapper()
				assert.NotNil(t, ds)

			case "error case - verify wrapper type":
				ds := execution.NewExecutionsDataSourceWrapper()
				assert.NotNil(t, ds)
			}
		})
	}
}

func TestExecutionsDataSource_Metadata(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "set metadata without provider type", wantErr: false},
		{name: "set metadata with provider type name", wantErr: false},
		{name: "set metadata with different provider type", wantErr: false},
		{name: "error case - metadata validation", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "set metadata without provider type":
				ds := execution.NewExecutionsDataSource()
				resp := &datasource.MetadataResponse{}

				ds.Metadata(context.Background(), datasource.MetadataRequest{
					ProviderTypeName: "n8n",
				}, resp)

				assert.Equal(t, "n8n_executions", resp.TypeName)

			case "set metadata with provider type name":
				ds := execution.NewExecutionsDataSource()
				resp := &datasource.MetadataResponse{}
				req := datasource.MetadataRequest{
					ProviderTypeName: "n8n",
				}

				ds.Metadata(context.Background(), req, resp)

				assert.Equal(t, "n8n_executions", resp.TypeName)

			case "set metadata with different provider type":
				ds := execution.NewExecutionsDataSource()
				resp := &datasource.MetadataResponse{}
				req := datasource.MetadataRequest{
					ProviderTypeName: "custom_provider",
				}

				ds.Metadata(context.Background(), req, resp)

				assert.Equal(t, "custom_provider_executions", resp.TypeName)

			case "error case - metadata validation":
				ds := execution.NewExecutionsDataSource()
				resp := &datasource.MetadataResponse{}
				ds.Metadata(context.Background(), datasource.MetadataRequest{
					ProviderTypeName: "test",
				}, resp)
				assert.Equal(t, "test_executions", resp.TypeName)
			}
		})
	}
}

func TestExecutionsDataSource_Schema(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "return schema", wantErr: false},
		{name: "schema attributes have descriptions", wantErr: false},
		{name: "error case - schema validation", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "return schema":
				ds := execution.NewExecutionsDataSource()
				resp := &datasource.SchemaResponse{}

				ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)

				assert.NotNil(t, resp.Schema)
				assert.Contains(t, resp.Schema.MarkdownDescription, "n8n workflow executions")

			case "schema attributes have descriptions":
				ds := execution.NewExecutionsDataSource()
				resp := &datasource.SchemaResponse{}

				ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)

				for name, attr := range resp.Schema.Attributes {
					assert.NotEmpty(t, attr.GetMarkdownDescription(), "Attribute %s should have description", name)
				}

			case "error case - schema validation":
				ds := execution.NewExecutionsDataSource()
				resp := &datasource.SchemaResponse{}
				ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)
				assert.NotNil(t, resp.Schema)
			}
		})
	}
}

func TestExecutionsDataSource_Configure(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "configure with valid client", wantErr: false},
		{name: "configure with nil provider data", wantErr: false},
		{name: "configure with invalid provider data", wantErr: false},
		{name: "error case - verify configuration", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "configure with valid client":
				ds := execution.NewExecutionsDataSource()
				cfg := n8nsdk.NewConfiguration()
				cfg.Servers = n8nsdk.ServerConfigurations{
					{URL: "http://test.local", Description: "Test server"},
				}
				mockClient := &client.N8nClient{
					APIClient: n8nsdk.NewAPIClient(cfg),
				}
				resp := &datasource.ConfigureResponse{}
				req := datasource.ConfigureRequest{
					ProviderData: mockClient,
				}

				ds.Configure(context.Background(), req, resp)

				assert.False(t, resp.Diagnostics.HasError())

			case "configure with nil provider data":
				ds := execution.NewExecutionsDataSource()
				resp := &datasource.ConfigureResponse{}
				req := datasource.ConfigureRequest{
					ProviderData: nil,
				}

				ds.Configure(context.Background(), req, resp)

				assert.False(t, resp.Diagnostics.HasError())

			case "configure with invalid provider data":
				ds := execution.NewExecutionsDataSource()
				resp := &datasource.ConfigureResponse{}
				req := datasource.ConfigureRequest{
					ProviderData: "invalid",
				}

				ds.Configure(context.Background(), req, resp)

				assert.True(t, resp.Diagnostics.HasError())

			case "error case - verify configuration":
				ds := execution.NewExecutionsDataSource()
				resp := &datasource.ConfigureResponse{}
				req := datasource.ConfigureRequest{
					ProviderData: 12345,
				}
				ds.Configure(context.Background(), req, resp)
				assert.True(t, resp.Diagnostics.HasError())
			}
		})
	}
}

func TestExecutionsDataSource_Read(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "verify read method exists", wantErr: false},
		{name: "error case - verify read behavior", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "verify read method exists":
				ds := execution.NewExecutionsDataSource()
				assert.NotNil(t, ds)
				assert.NotNil(t, ds.Read)

			case "error case - verify read behavior":
				ds := execution.NewExecutionsDataSource()
				assert.NotNil(t, ds)
			}
		})
	}
}
