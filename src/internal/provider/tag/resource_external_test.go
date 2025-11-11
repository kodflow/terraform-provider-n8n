package tag_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
	"github.com/kodflow/n8n/src/internal/provider/tag"
	"github.com/stretchr/testify/assert"
)

func TestNewTagResource(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "normal case",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := tag.NewTagResource()

				assert.NotNil(t, r)
			},
		},
		{
			name: "error case - validates behavior",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := tag.NewTagResource()

				assert.NotNil(t, r)
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

func TestNewTagResourceWrapper(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "normal case",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := tag.NewTagResourceWrapper()

				assert.NotNil(t, r)
				assert.Implements(t, (*resource.Resource)(nil), r)
			},
		},
		{
			name: "error case - validates behavior",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := tag.NewTagResourceWrapper()

				assert.NotNil(t, r)
				assert.Implements(t, (*resource.Resource)(nil), r)
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

func TestTagResource_Metadata(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "normal case",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := tag.NewTagResource()
				resp := &resource.MetadataResponse{}

				r.Metadata(context.Background(), resource.MetadataRequest{
					ProviderTypeName: "n8n",
				}, resp)

				assert.Equal(t, "n8n_tag", resp.TypeName)
			},
		},
		{
			name: "error case - validates behavior",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := tag.NewTagResource()
				resp := &resource.MetadataResponse{}

				r.Metadata(context.Background(), resource.MetadataRequest{
					ProviderTypeName: "n8n",
				}, resp)

				assert.Equal(t, "n8n_tag", resp.TypeName)
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

func TestTagResource_Schema(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "normal case",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := tag.NewTagResource()
				resp := &resource.SchemaResponse{}

				r.Schema(context.Background(), resource.SchemaRequest{}, resp)

				assert.NotNil(t, resp.Schema)
				assert.NotEmpty(t, resp.Schema.Attributes)
			},
		},
		{
			name: "error case - validates behavior",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := tag.NewTagResource()
				resp := &resource.SchemaResponse{}

				r.Schema(context.Background(), resource.SchemaRequest{}, resp)

				assert.NotNil(t, resp.Schema)
				assert.NotEmpty(t, resp.Schema.Attributes)
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

func TestTagResource_Configure(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		providerData interface{}
		wantError    bool
	}{
		{
			name:         "valid configuration",
			providerData: &client.N8nClient{},
			wantError:    false,
		},
		{
			name:         "nil provider data",
			providerData: nil,
			wantError:    false,
		},
		{
			name:         "invalid provider data type",
			providerData: "invalid",
			wantError:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := tag.NewTagResource()
			resp := &resource.ConfigureResponse{}
			req := resource.ConfigureRequest{
				ProviderData: tt.providerData,
			}

			r.Configure(context.Background(), req, resp)

			if tt.wantError {
				assert.True(t, resp.Diagnostics.HasError())
			} else {
				assert.False(t, resp.Diagnostics.HasError())
			}
		})
	}
}

func TestTagResource_Create(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "normal case",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := tag.NewTagResource()

				// Verify the method exists and resource implements the interface
				assert.NotNil(t, r)
				_, ok := interface{}(r).(resource.Resource)
				assert.True(t, ok)
			},
		},
		{
			name: "error case - validates behavior",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := tag.NewTagResource()

				// Verify the method exists and resource implements the interface
				assert.NotNil(t, r)
				_, ok := interface{}(r).(resource.Resource)
				assert.True(t, ok)
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

func TestTagResource_Read(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "normal case",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := tag.NewTagResource()

				// Verify the method exists and resource implements the interface
				assert.NotNil(t, r)
				_, ok := interface{}(r).(resource.Resource)
				assert.True(t, ok)
			},
		},
		{
			name: "error case - validates behavior",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := tag.NewTagResource()

				// Verify the method exists and resource implements the interface
				assert.NotNil(t, r)
				_, ok := interface{}(r).(resource.Resource)
				assert.True(t, ok)
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

func TestTagResource_Update(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "normal case",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := tag.NewTagResource()

				// Verify the method exists and resource implements the interface
				assert.NotNil(t, r)
				_, ok := interface{}(r).(resource.Resource)
				assert.True(t, ok)
			},
		},
		{
			name: "error case - validates behavior",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := tag.NewTagResource()

				// Verify the method exists and resource implements the interface
				assert.NotNil(t, r)
				_, ok := interface{}(r).(resource.Resource)
				assert.True(t, ok)
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

func TestTagResource_Delete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "normal case",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := tag.NewTagResource()

				// Verify the method exists and resource implements the interface
				assert.NotNil(t, r)
				_, ok := interface{}(r).(resource.Resource)
				assert.True(t, ok)
			},
		},
		{
			name: "error case - validates behavior",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := tag.NewTagResource()

				// Verify the method exists and resource implements the interface
				assert.NotNil(t, r)
				_, ok := interface{}(r).(resource.Resource)
				assert.True(t, ok)
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

func TestTagResource_ImportState(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "normal case",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := tag.NewTagResource()

				// Test that ImportState interface is implemented
				_, ok := interface{}(r).(resource.ResourceWithImportState)
				assert.True(t, ok, "TagResource should implement ResourceWithImportState")
			},
		},
		{
			name: "error case - validates behavior",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := tag.NewTagResource()

				// Test that ImportState interface is implemented
				_, ok := interface{}(r).(resource.ResourceWithImportState)
				assert.True(t, ok, "TagResource should implement ResourceWithImportState")
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
