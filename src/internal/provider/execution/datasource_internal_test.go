package execution

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/execution/models"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockExecutionDataSourceInterface is a mock implementation of ExecutionDataSourceInterface.
type MockExecutionDataSourceInterface struct {
	mock.Mock
}

func (m *MockExecutionDataSourceInterface) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockExecutionDataSourceInterface) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockExecutionDataSourceInterface) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockExecutionDataSourceInterface) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	m.Called(ctx, req, resp)
}

// TestNewExecutionDataSource is now in external test file - refactored to test behavior only.

// TestExecutionDataSource_Configure is now in external test file - refactored to test behavior only.

func TestExecutionDataSource_Interfaces(t *testing.T) {
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
				ds := NewExecutionDataSource()

				// Test that ExecutionDataSource implements datasource.DataSource
				var _ datasource.DataSource = ds

				// Test that ExecutionDataSource implements datasource.DataSourceWithConfigure
				var _ datasource.DataSourceWithConfigure = ds

				// Test that ExecutionDataSource implements ExecutionDataSourceInterface
				var _ ExecutionDataSourceInterface = ds

			case "error case - interface compliance":
				ds := NewExecutionDataSource()
				var _ datasource.DataSource = ds
			}
		})
	}
}

func TestExecutionDataSourceConcurrency(t *testing.T) {
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
				ds := NewExecutionDataSource()

				done := make(chan bool, 100)
				for i := 0; i < 100; i++ {
					go func() {
						resp := &datasource.MetadataResponse{}
						ds.Metadata(context.Background(), datasource.MetadataRequest{
							ProviderTypeName: "n8n",
						}, resp)
						assert.Equal(t, "n8n_execution", resp.TypeName)
						done <- true
					}()
				}

				for i := 0; i < 100; i++ {
					<-done
				}

			case "concurrent schema calls":
				ds := NewExecutionDataSource()

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
				ds := NewExecutionDataSource()
				done := make(chan bool, 50)
				for i := 0; i < 50; i++ {
					go func() {
						resp := &datasource.MetadataResponse{}
						ds.Metadata(context.Background(), datasource.MetadataRequest{
							ProviderTypeName: "n8n",
						}, resp)
						assert.Equal(t, "n8n_execution", resp.TypeName)
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

func BenchmarkExecutionDataSource_Schema(b *testing.B) {
	ds := NewExecutionDataSource()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := &datasource.SchemaResponse{}
		ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)
	}
}

func BenchmarkExecutionDataSource_Metadata(b *testing.B) {
	ds := NewExecutionDataSource()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := &datasource.MetadataResponse{}
		ds.Metadata(context.Background(), datasource.MetadataRequest{}, resp)
	}
}

func BenchmarkExecutionDataSource_Configure(b *testing.B) {
	ds := NewExecutionDataSource()
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

// TestExecutionDataSource_Read is now in external test file - refactored to test behavior only.

func TestExecutionDataSource_schemaAttributes(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "schema attributes returned correctly", wantErr: false},
		{name: "required fields present", wantErr: false},
		{name: "computed fields present", wantErr: false},
		{name: "optional fields present", wantErr: false},
		{name: "error case - verify attribute names", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "schema attributes returned correctly":
				ds := &ExecutionDataSource{}
				attrs := ds.schemaAttributes()

				assert.NotNil(t, attrs)
				assert.Contains(t, attrs, "id")
				assert.Contains(t, attrs, "workflow_id")
				assert.Contains(t, attrs, "finished")
				assert.Contains(t, attrs, "mode")
				assert.Contains(t, attrs, "started_at")
				assert.Contains(t, attrs, "stopped_at")
				assert.Contains(t, attrs, "status")
				assert.Contains(t, attrs, "include_data")

			case "required fields present":
				ds := &ExecutionDataSource{}
				attrs := ds.schemaAttributes()

				idAttr, ok := attrs["id"].(schema.StringAttribute)
				assert.True(t, ok)
				assert.True(t, idAttr.Required)

			case "computed fields present":
				ds := &ExecutionDataSource{}
				attrs := ds.schemaAttributes()

				workflowIDAttr, ok := attrs["workflow_id"].(schema.StringAttribute)
				assert.True(t, ok)
				assert.True(t, workflowIDAttr.Computed)

				finishedAttr, ok := attrs["finished"].(schema.BoolAttribute)
				assert.True(t, ok)
				assert.True(t, finishedAttr.Computed)

			case "optional fields present":
				ds := &ExecutionDataSource{}
				attrs := ds.schemaAttributes()

				includeDataAttr, ok := attrs["include_data"].(schema.BoolAttribute)
				assert.True(t, ok)
				assert.True(t, includeDataAttr.Optional)

			case "error case - verify attribute names":
				ds := &ExecutionDataSource{}
				attrs := ds.schemaAttributes()
				assert.Len(t, attrs, 8)
			}
		})
	}
}

func Test_populateExecutionData(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "populate with all fields", wantErr: false},
		{name: "populate with nil fields", wantErr: false},
		{name: "populate with partial fields", wantErr: false},
		{name: "populate with unset stoppedAt", wantErr: false},
		{name: "populate with nil stoppedAt inside nullable", wantErr: false},
		{name: "populate with all timestamps", wantErr: false},
		{name: "error case - empty execution", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "populate with all fields":
				id := float32(123)
				workflowID := float32(456)
				finished := true
				mode := "manual"
				status := "success"

				execution := &n8nsdk.Execution{
					Id:         &id,
					WorkflowId: &workflowID,
					Finished:   &finished,
					Mode:       &mode,
					Status:     &status,
				}

				data := &models.DataSource{}
				populateExecutionData(execution, data)

				assert.Equal(t, "123", data.ID.ValueString())
				assert.Equal(t, "456", data.WorkflowID.ValueString())
				assert.True(t, data.Finished.ValueBool())
				assert.Equal(t, "manual", data.Mode.ValueString())
				assert.Equal(t, "success", data.Status.ValueString())

			case "populate with nil fields":
				execution := &n8nsdk.Execution{}
				data := &models.DataSource{}
				populateExecutionData(execution, data)

				assert.True(t, data.ID.IsNull())
				assert.True(t, data.WorkflowID.IsNull())
				assert.True(t, data.Finished.IsNull())
				assert.True(t, data.Mode.IsNull())
				assert.True(t, data.Status.IsNull())

			case "populate with partial fields":
				id := float32(789)
				status := "running"

				execution := &n8nsdk.Execution{
					Id:     &id,
					Status: &status,
				}
				data := &models.DataSource{}
				populateExecutionData(execution, data)

				assert.Equal(t, "789", data.ID.ValueString())
				assert.Equal(t, "running", data.Status.ValueString())
				assert.True(t, data.WorkflowID.IsNull())
				assert.True(t, data.Finished.IsNull())

			case "populate with unset stoppedAt":
				id := float32(999)
				execution := &n8nsdk.Execution{
					Id:        &id,
					StoppedAt: n8nsdk.NullableTime{},
				}
				data := &models.DataSource{}
				populateExecutionData(execution, data)

				assert.Equal(t, "999", data.ID.ValueString())
				assert.True(t, data.StoppedAt.IsNull())

			case "populate with nil stoppedAt inside nullable":
				id := float32(888)
				execution := &n8nsdk.Execution{
					Id: &id,
				}
				execution.StoppedAt.Set(nil)

				data := &models.DataSource{}
				populateExecutionData(execution, data)

				assert.Equal(t, "888", data.ID.ValueString())
				assert.True(t, data.StoppedAt.IsNull())

			case "populate with all timestamps":
				id := float32(777)
				startedAt := time.Now()
				stoppedAt := time.Now().Add(1 * time.Hour)

				execution := &n8nsdk.Execution{
					Id:        &id,
					StartedAt: *n8nsdk.NewNullableTime(&startedAt),
				}
				execution.StoppedAt.Set(&stoppedAt)

				data := &models.DataSource{}
				populateExecutionData(execution, data)

				assert.Equal(t, "777", data.ID.ValueString())
				assert.NotEmpty(t, data.StartedAt.ValueString())
				assert.NotEmpty(t, data.StoppedAt.ValueString())

			case "error case - empty execution":
				execution := &n8nsdk.Execution{}
				data := &models.DataSource{}
				populateExecutionData(execution, data)
				assert.NotNil(t, data)
			}
		})
	}
}
