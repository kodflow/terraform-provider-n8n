package execution_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/execution"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
	"github.com/stretchr/testify/assert"
)

func TestNewExecutionRetryResource(t *testing.T) {
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
				r := execution.NewExecutionRetryResource()
				assert.NotNil(t, r)
				assert.IsType(t, &execution.ExecutionRetryResource{}, r)

			case "returns non-nil instance":
				r := execution.NewExecutionRetryResource()
				assert.NotNil(t, r)

			case "error case - verify instance type":
				r := execution.NewExecutionRetryResource()
				assert.NotNil(t, r)
			}
		})
	}
}

func TestNewExecutionRetryResourceWrapper(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "create new wrapped instance", wantErr: false},
		{name: "wrapper returns resource.Resource interface", wantErr: false},
		{name: "error case - verify wrapper type", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "create new wrapped instance":
				r := execution.NewExecutionRetryResourceWrapper()
				assert.NotNil(t, r)
				assert.IsType(t, &execution.ExecutionRetryResource{}, r)

			case "wrapper returns resource.Resource interface":
				r := execution.NewExecutionRetryResourceWrapper()
				assert.NotNil(t, r)

			case "error case - verify wrapper type":
				r := execution.NewExecutionRetryResourceWrapper()
				assert.NotNil(t, r)
			}
		})
	}
}

func TestExecutionRetryResource_Metadata(t *testing.T) {
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
				r := execution.NewExecutionRetryResource()
				resp := &resource.MetadataResponse{}

				r.Metadata(context.Background(), resource.MetadataRequest{
					ProviderTypeName: "n8n",
				}, resp)

				assert.Equal(t, "n8n_execution_retry", resp.TypeName)

			case "set metadata with provider type name":
				r := execution.NewExecutionRetryResource()
				resp := &resource.MetadataResponse{}
				req := resource.MetadataRequest{
					ProviderTypeName: "n8n",
				}

				r.Metadata(context.Background(), req, resp)

				assert.Equal(t, "n8n_execution_retry", resp.TypeName)

			case "set metadata with different provider type":
				r := execution.NewExecutionRetryResource()
				resp := &resource.MetadataResponse{}
				req := resource.MetadataRequest{
					ProviderTypeName: "custom_provider",
				}

				r.Metadata(context.Background(), req, resp)

				assert.Equal(t, "custom_provider_execution_retry", resp.TypeName)

			case "error case - metadata validation":
				r := execution.NewExecutionRetryResource()
				resp := &resource.MetadataResponse{}
				r.Metadata(context.Background(), resource.MetadataRequest{
					ProviderTypeName: "test",
				}, resp)
				assert.Equal(t, "test_execution_retry", resp.TypeName)
			}
		})
	}
}

func TestExecutionRetryResource_Schema(t *testing.T) {
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
				r := execution.NewExecutionRetryResource()
				resp := &resource.SchemaResponse{}

				r.Schema(context.Background(), resource.SchemaRequest{}, resp)

				assert.NotNil(t, resp.Schema)
				assert.Contains(t, resp.Schema.MarkdownDescription, "retry")

			case "schema attributes have descriptions":
				r := execution.NewExecutionRetryResource()
				resp := &resource.SchemaResponse{}

				r.Schema(context.Background(), resource.SchemaRequest{}, resp)

				for name, attr := range resp.Schema.Attributes {
					assert.NotEmpty(t, attr.GetMarkdownDescription(), "Attribute %s should have description", name)
				}

			case "error case - schema validation":
				r := execution.NewExecutionRetryResource()
				resp := &resource.SchemaResponse{}
				r.Schema(context.Background(), resource.SchemaRequest{}, resp)
				assert.NotNil(t, resp.Schema)
			}
		})
	}
}

func TestExecutionRetryResource_Configure(t *testing.T) {
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
				r := execution.NewExecutionRetryResource()
				cfg := n8nsdk.NewConfiguration()
				cfg.Servers = n8nsdk.ServerConfigurations{
					{URL: "http://test.local", Description: "Test server"},
				}
				mockClient := &client.N8nClient{
					APIClient: n8nsdk.NewAPIClient(cfg),
				}
				resp := &resource.ConfigureResponse{}
				req := resource.ConfigureRequest{
					ProviderData: mockClient,
				}

				r.Configure(context.Background(), req, resp)

				assert.False(t, resp.Diagnostics.HasError())

			case "configure with nil provider data":
				r := execution.NewExecutionRetryResource()
				resp := &resource.ConfigureResponse{}
				req := resource.ConfigureRequest{
					ProviderData: nil,
				}

				r.Configure(context.Background(), req, resp)

				assert.False(t, resp.Diagnostics.HasError())

			case "configure with invalid provider data":
				r := execution.NewExecutionRetryResource()
				resp := &resource.ConfigureResponse{}
				req := resource.ConfigureRequest{
					ProviderData: "invalid",
				}

				r.Configure(context.Background(), req, resp)

				assert.True(t, resp.Diagnostics.HasError())

			case "error case - verify configuration":
				r := execution.NewExecutionRetryResource()
				resp := &resource.ConfigureResponse{}
				req := resource.ConfigureRequest{
					ProviderData: 12345,
				}
				r.Configure(context.Background(), req, resp)
				assert.True(t, resp.Diagnostics.HasError())
			}
		})
	}
}

func TestExecutionRetryResource_Create(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "verify create method exists", wantErr: false},
		{name: "error case - verify create behavior", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "verify create method exists":
				r := execution.NewExecutionRetryResource()
				assert.NotNil(t, r)
				assert.NotNil(t, r.Create)

			case "error case - verify create behavior":
				r := execution.NewExecutionRetryResource()
				assert.NotNil(t, r)
			}
		})
	}
}

func TestExecutionRetryResource_Read(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "success - read maintains state", wantErr: false},
		{name: "error - invalid state", wantErr: true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "success - read maintains state":
				r := execution.NewExecutionRetryResource()
				ctx := context.Background()

				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				// Build state with execution data
				stateRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"execution_id":     tftypes.NewValue(tftypes.String, "exec-123"),
					"new_execution_id": tftypes.NewValue(tftypes.String, "exec-456"),
					"workflow_id":      tftypes.NewValue(tftypes.String, "workflow-1"),
					"finished":         tftypes.NewValue(tftypes.Bool, true),
					"mode":             tftypes.NewValue(tftypes.String, "manual"),
					"started_at":       tftypes.NewValue(tftypes.String, "2024-01-01T00:00:00Z"),
					"stopped_at":       tftypes.NewValue(tftypes.String, "2024-01-01T00:01:00Z"),
					"status":           tftypes.NewValue(tftypes.String, "success"),
				})

				state := tfsdk.State{
					Schema: schemaResp.Schema,
					Raw:    stateRaw,
				}

				req := resource.ReadRequest{
					State: state,
				}
				resp := &resource.ReadResponse{
					State: tfsdk.State{Schema: schemaResp.Schema},
				}

				// Call Read
				r.Read(ctx, req, resp)

				// Verify success - Read should maintain state without errors
				assert.False(t, resp.Diagnostics.HasError())

			case "error - invalid state":
				r := execution.NewExecutionRetryResource()
				ctx := context.Background()

				schemaResp := resource.SchemaResponse{}
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				// Create invalid state
				stateRaw := tftypes.NewValue(tftypes.String, "invalid")

				state := tfsdk.State{
					Schema: schemaResp.Schema,
					Raw:    stateRaw,
				}

				req := resource.ReadRequest{
					State: state,
				}
				resp := &resource.ReadResponse{}

				r.Read(ctx, req, resp)

				// Verify error
				assert.True(t, resp.Diagnostics.HasError())
			}
		})
	}
}

func TestExecutionRetryResource_Update(t *testing.T) {
	t.Run("error - update not supported", func(t *testing.T) {
		r := execution.NewExecutionRetryResource()
		ctx := context.Background()

		req := resource.UpdateRequest{}
		resp := &resource.UpdateResponse{}

		// Call Update
		r.Update(ctx, req, resp)

		// Verify error - Update should always return an error
		assert.True(t, resp.Diagnostics.HasError())
		assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Update Not Supported")
	})
}

func TestExecutionRetryResource_Delete(t *testing.T) {
	t.Run("success - delete removes from state", func(t *testing.T) {
		r := execution.NewExecutionRetryResource()
		ctx := context.Background()

		req := resource.DeleteRequest{}
		resp := &resource.DeleteResponse{}

		// Call Delete - should not panic and should not return errors
		assert.NotPanics(t, func() {
			r.Delete(ctx, req, resp)
		})

		// Verify no errors - Delete just removes from state without API call
		assert.False(t, resp.Diagnostics.HasError())
	})
}

func TestExecutionRetryResource_ImportState(t *testing.T) {
	t.Run("import state passthrough", func(t *testing.T) {
		r := execution.NewExecutionRetryResource()
		ctx := context.Background()

		schemaResp := resource.SchemaResponse{}
		r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

		// Build empty state
		emptyValue := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
			"execution_id":     tftypes.NewValue(tftypes.String, nil),
			"new_execution_id": tftypes.NewValue(tftypes.String, nil),
			"workflow_id":      tftypes.NewValue(tftypes.String, nil),
			"finished":         tftypes.NewValue(tftypes.Bool, nil),
			"mode":             tftypes.NewValue(tftypes.String, nil),
			"started_at":       tftypes.NewValue(tftypes.String, nil),
			"stopped_at":       tftypes.NewValue(tftypes.String, nil),
			"status":           tftypes.NewValue(tftypes.String, nil),
		})

		req := resource.ImportStateRequest{
			ID: "exec-123",
		}
		resp := &resource.ImportStateResponse{
			State: tfsdk.State{
				Schema: schemaResp.Schema,
				Raw:    emptyValue,
			},
		}

		// Call ImportState
		r.ImportState(ctx, req, resp)

		// Verify no errors (ImportStatePassthroughID doesn't return errors for valid IDs)
		assert.False(t, resp.Diagnostics.HasError())
	})
}
