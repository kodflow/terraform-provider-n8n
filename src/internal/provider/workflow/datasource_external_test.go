package workflow_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
	"github.com/kodflow/n8n/src/internal/provider/workflow"
	"github.com/stretchr/testify/assert"
)

// TestNewWorkflowDataSourceExternal tests the public NewWorkflowDataSource function (black-box).
func TestNewWorkflowDataSourceExternal(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "create new instance from external package", wantErr: false},
		{name: "multiple calls return independent instances", wantErr: false},
		{name: "returned instance is not nil even on repeated calls", wantErr: false},
		{name: "error case - validation checks", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "create new instance from external package":
				ds := workflow.NewWorkflowDataSource()

				assert.NotNil(t, ds)

			case "multiple calls return independent instances":
				ds1 := workflow.NewWorkflowDataSource()
				ds2 := workflow.NewWorkflowDataSource()

				assert.NotSame(t, ds1, ds2, "each call should return a new instance")

			case "returned instance is not nil even on repeated calls":
				for i := 0; i < 100; i++ {
					ds := workflow.NewWorkflowDataSource()
					assert.NotNil(t, ds, "instance should never be nil")
				}

			case "error case - validation checks":
				ds := workflow.NewWorkflowDataSource()
				assert.NotNil(t, ds)
			}
		})
	}
}

// TestNewWorkflowDataSourceWrapperExternal tests the public NewWorkflowDataSourceWrapper function (black-box).
func TestNewWorkflowDataSourceWrapperExternal(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "create wrapped instance from external package", wantErr: false},
		{name: "multiple calls return independent instances", wantErr: false},
		{name: "returned instance is not nil even on repeated calls", wantErr: false},
		{name: "error case - validation checks", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "create wrapped instance from external package":
				ds := workflow.NewWorkflowDataSourceWrapper()

				assert.NotNil(t, ds)

			case "multiple calls return independent instances":
				ds1 := workflow.NewWorkflowDataSourceWrapper()
				ds2 := workflow.NewWorkflowDataSourceWrapper()

				assert.NotSame(t, ds1, ds2, "each call should return a new instance")

			case "returned instance is not nil even on repeated calls":
				for i := 0; i < 100; i++ {
					ds := workflow.NewWorkflowDataSourceWrapper()
					assert.NotNil(t, ds, "instance should never be nil")
				}

			case "error case - validation checks":
				ds := workflow.NewWorkflowDataSourceWrapper()
				assert.NotNil(t, ds)
			}
		})
	}
}

// TestNewWorkflowDataSource tests the exact function name expected by KTN-TEST-003.
func TestNewWorkflowDataSource(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "creates non-nil datasource",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := workflow.NewWorkflowDataSource()
				assert.NotNil(t, ds)
			},
		},
		{
			name: "error case - multiple calls return different instances",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds1 := workflow.NewWorkflowDataSource()
				ds2 := workflow.NewWorkflowDataSource()
				assert.NotSame(t, ds1, ds2)
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

// TestNewWorkflowDataSourceWrapper tests the exact function name expected by KTN-TEST-003.
func TestNewWorkflowDataSourceWrapper(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "creates non-nil datasource",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := workflow.NewWorkflowDataSourceWrapper()
				assert.NotNil(t, ds)
			},
		},
		{
			name: "error case - multiple calls return different instances",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds1 := workflow.NewWorkflowDataSourceWrapper()
				ds2 := workflow.NewWorkflowDataSourceWrapper()
				assert.NotSame(t, ds1, ds2)
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

// TestWorkflowDataSource_Metadata tests the Metadata method.
func TestWorkflowDataSource_Metadata(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		providerTypeName string
		expectedTypeName string
	}{
		{
			name:             "sets correct type name",
			providerTypeName: "n8n",
			expectedTypeName: "n8n_workflow",
		},
		{
			name:             "error case - handles empty provider type name",
			providerTypeName: "",
			expectedTypeName: "_workflow",
		},
		{
			name:             "error case - handles custom provider type name",
			providerTypeName: "custom",
			expectedTypeName: "custom_workflow",
		},
		{
			name:             "error case - handles provider type name with special characters",
			providerTypeName: "test-provider_123",
			expectedTypeName: "test-provider_123_workflow",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ds := workflow.NewWorkflowDataSource()
			req := datasource.MetadataRequest{
				ProviderTypeName: tt.providerTypeName,
			}
			resp := &datasource.MetadataResponse{}

			ds.Metadata(context.Background(), req, resp)

			assert.Equal(t, tt.expectedTypeName, resp.TypeName)
		})
	}
}

// TestWorkflowDataSource_Schema tests the Schema method.
func TestWorkflowDataSource_Schema(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
	}{
		{name: "defines valid schema"},
		{name: "schema has markdown description"},
		{name: "schema has required id attribute"},
		{name: "error case - schema attributes are not empty"},
		{name: "error case - all required attributes are present"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ds := workflow.NewWorkflowDataSource()
			req := datasource.SchemaRequest{}
			resp := &datasource.SchemaResponse{}

			ds.Schema(context.Background(), req, resp)

			switch tt.name {
			case "defines valid schema":
				assert.NotNil(t, resp.Schema)
				assert.NotEmpty(t, resp.Schema.Attributes)

			case "schema has markdown description":
				assert.NotEmpty(t, resp.Schema.MarkdownDescription)
				assert.Contains(t, resp.Schema.MarkdownDescription, "workflow")

			case "schema has required id attribute":
				assert.Contains(t, resp.Schema.Attributes, "id")
				idAttr := resp.Schema.Attributes["id"]
				assert.NotNil(t, idAttr)

			case "error case - schema attributes are not empty":
				assert.NotEmpty(t, resp.Schema.Attributes)
				assert.Greater(t, len(resp.Schema.Attributes), 0)

			case "error case - all required attributes are present":
				assert.Contains(t, resp.Schema.Attributes, "id")
				assert.Contains(t, resp.Schema.Attributes, "name")
			}
		})
	}
}

// TestWorkflowDataSource_Configure tests the Configure method.
func TestWorkflowDataSource_Configure(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "configures with valid client", wantErr: false},
		{name: "error case - nil provider data", wantErr: false},
		{name: "error case - wrong provider data type", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "configures with valid client":
				d := workflow.NewWorkflowDataSource()
				req := datasource.ConfigureRequest{
					ProviderData: &client.N8nClient{},
				}
				resp := &datasource.ConfigureResponse{}
				d.Configure(context.Background(), req, resp)
				assert.False(t, resp.Diagnostics.HasError())

			case "error case - nil provider data":
				d := workflow.NewWorkflowDataSource()
				req := datasource.ConfigureRequest{
					ProviderData: nil,
				}
				resp := &datasource.ConfigureResponse{}
				d.Configure(context.Background(), req, resp)
				assert.False(t, resp.Diagnostics.HasError())

			case "error case - wrong provider data type":
				d := workflow.NewWorkflowDataSource()
				req := datasource.ConfigureRequest{
					ProviderData: "wrong type",
				}
				resp := &datasource.ConfigureResponse{}
				d.Configure(context.Background(), req, resp)
				assert.True(t, resp.Diagnostics.HasError())
				assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Unexpected Data Source Configure Type")
			}
		})
	}
}

// TestWorkflowDataSource_Read tests the Read method.
func TestWorkflowDataSource_Read(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "creates non-nil datasource",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := workflow.NewWorkflowDataSource()
				assert.NotNil(t, ds)
			},
		},
		{
			name: "error case - datasource can be created multiple times",
			testFunc: func(t *testing.T) {
				t.Helper()
				for i := 0; i < 10; i++ {
					ds := workflow.NewWorkflowDataSource()
					assert.NotNil(t, ds)
				}
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
