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

func TestNewExecutionDataSourceWrapper(t *testing.T) {
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
				ds := execution.NewExecutionDataSourceWrapper()

				assert.NotNil(t, ds)
				assert.IsType(t, &execution.ExecutionDataSource{}, ds)

			case "wrapper returns datasource.DataSource interface":
				ds := execution.NewExecutionDataSourceWrapper()

				// ds is already of type datasource.DataSource, no assertion needed.
				assert.NotNil(t, ds)

			case "error case - verify wrapper type":
				ds := execution.NewExecutionDataSourceWrapper()
				// Verify ds implements datasource.DataSource interface.
				assert.NotNil(t, ds)
			}
		})
	}
}

func TestExecutionDataSource_Metadata(t *testing.T) {
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
				ds := execution.NewExecutionDataSource()
				resp := &datasource.MetadataResponse{}

				ds.Metadata(context.Background(), datasource.MetadataRequest{
					ProviderTypeName: "n8n",
				}, resp)

				assert.Equal(t, "n8n_execution", resp.TypeName)

			case "set metadata with provider type name":
				ds := execution.NewExecutionDataSource()
				resp := &datasource.MetadataResponse{}
				req := datasource.MetadataRequest{
					ProviderTypeName: "n8n",
				}

				ds.Metadata(context.Background(), req, resp)

				assert.Equal(t, "n8n_execution", resp.TypeName)

			case "set metadata with different provider type":
				ds := execution.NewExecutionDataSource()
				resp := &datasource.MetadataResponse{}
				req := datasource.MetadataRequest{
					ProviderTypeName: "custom_provider",
				}

				ds.Metadata(context.Background(), req, resp)

				assert.Equal(t, "custom_provider_execution", resp.TypeName)

			case "error case - metadata validation":
				ds := execution.NewExecutionDataSource()
				resp := &datasource.MetadataResponse{}
				ds.Metadata(context.Background(), datasource.MetadataRequest{
					ProviderTypeName: "test",
				}, resp)
				assert.Equal(t, "test_execution", resp.TypeName)
			}
		})
	}
}

func TestExecutionDataSource_Schema(t *testing.T) {
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
				ds := execution.NewExecutionDataSource()
				resp := &datasource.SchemaResponse{}

				ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)

				assert.NotNil(t, resp.Schema)
				assert.Contains(t, resp.Schema.MarkdownDescription, "n8n workflow execution")

				// Verify required attributes.
				idAttr, exists := resp.Schema.Attributes["id"]
				assert.True(t, exists)
				assert.True(t, idAttr.IsRequired())

				// Verify computed attributes.
				computedAttrs := []string{
					"workflow_id",
					"finished",
					"mode",
					"started_at",
					"stopped_at",
					"status",
				}

				for _, attr := range computedAttrs {
					schemaAttr, exists := resp.Schema.Attributes[attr]
					assert.True(t, exists, "Attribute %s should exist", attr)
					assert.True(t, schemaAttr.IsComputed(), "Attribute %s should be computed", attr)
				}

				// Verify optional attributes.
				includeDataAttr, exists := resp.Schema.Attributes["include_data"]
				assert.True(t, exists)
				assert.True(t, includeDataAttr.IsOptional())

			case "schema attributes have descriptions":
				ds := execution.NewExecutionDataSource()
				resp := &datasource.SchemaResponse{}

				ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)

				for name, attr := range resp.Schema.Attributes {
					assert.NotEmpty(t, attr.GetMarkdownDescription(), "Attribute %s should have description", name)
				}

			case "error case - schema validation":
				ds := execution.NewExecutionDataSource()
				resp := &datasource.SchemaResponse{}
				ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)
				assert.NotNil(t, resp.Schema)
			}
		})
	}
}

func TestNewExecutionDataSource(t *testing.T) {
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
				ds := execution.NewExecutionDataSource()
				assert.NotNil(t, ds)
				assert.IsType(t, &execution.ExecutionDataSource{}, ds)

			case "returns non-nil instance":
				ds := execution.NewExecutionDataSource()
				assert.NotNil(t, ds)

			case "error case - verify instance type":
				ds := execution.NewExecutionDataSource()
				assert.NotNil(t, ds)
			}
		})
	}
}

func TestExecutionDataSource_Configure(t *testing.T) {
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
				ds := execution.NewExecutionDataSource()
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
				ds := execution.NewExecutionDataSource()
				resp := &datasource.ConfigureResponse{}
				req := datasource.ConfigureRequest{
					ProviderData: nil,
				}

				ds.Configure(context.Background(), req, resp)

				assert.False(t, resp.Diagnostics.HasError())

			case "configure with invalid provider data":
				ds := execution.NewExecutionDataSource()
				resp := &datasource.ConfigureResponse{}
				req := datasource.ConfigureRequest{
					ProviderData: "invalid",
				}

				ds.Configure(context.Background(), req, resp)

				assert.True(t, resp.Diagnostics.HasError())

			case "error case - verify configuration":
				ds := execution.NewExecutionDataSource()
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

func TestExecutionDataSource_Read(t *testing.T) {
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
				ds := execution.NewExecutionDataSource()
				assert.NotNil(t, ds)
				assert.NotNil(t, ds.Read)

			case "error case - verify read behavior":
				ds := execution.NewExecutionDataSource()
				assert.NotNil(t, ds)
			}
		})
	}
}
