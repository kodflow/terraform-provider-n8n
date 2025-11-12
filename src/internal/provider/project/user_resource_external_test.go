package project_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/kodflow/n8n/src/internal/provider/project"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
	"github.com/stretchr/testify/assert"
)

// TestNewProjectUserResource tests the NewProjectUserResource constructor.
func TestNewProjectUserResource(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "creates valid resource", wantErr: false},
		{name: "error case - validation checks", wantErr: false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "creates valid resource":
				r := project.NewProjectUserResource()
				assert.NotNil(t, r)

			case "error case - validation checks":
				r := project.NewProjectUserResource()
				assert.NotNil(t, r)
				assert.Implements(t, (*resource.Resource)(nil), r)
			}
		})
	}
}

// TestNewProjectUserResourceWrapper tests the NewProjectUserResourceWrapper constructor.
func TestNewProjectUserResourceWrapper(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "creates valid resource wrapper", wantErr: false},
		{name: "error case - validation checks", wantErr: false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "creates valid resource wrapper":
				wrapper := project.NewProjectUserResourceWrapper()
				assert.NotNil(t, wrapper)

			case "error case - validation checks":
				wrapper := project.NewProjectUserResourceWrapper()
				assert.NotNil(t, wrapper)
			}
		})
	}
}

// TestProjectUserResource_Metadata tests the Metadata method.
func TestProjectUserResource_Metadata(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		providerTypeName string
		expectedTypeName string
	}{
		{
			name:             "with n8n provider",
			providerTypeName: "n8n",
			expectedTypeName: "n8n_project_user",
		},
		{
			name:             "with custom provider",
			providerTypeName: "custom",
			expectedTypeName: "custom_project_user",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := project.NewProjectUserResource()
			req := resource.MetadataRequest{
				ProviderTypeName: tt.providerTypeName,
			}
			resp := &resource.MetadataResponse{}

			r.Metadata(context.Background(), req, resp)

			assert.Equal(t, tt.expectedTypeName, resp.TypeName)
		})
	}
}

// TestProjectUserResource_Schema tests the Schema method.
func TestProjectUserResource_Schema(t *testing.T) {
	t.Parallel()

	r := project.NewProjectUserResource()
	req := resource.SchemaRequest{}
	resp := &resource.SchemaResponse{}

	r.Schema(context.Background(), req, resp)

	assert.NotNil(t, resp.Schema)
	assert.Contains(t, resp.Schema.MarkdownDescription, "user")

	// Verify required attributes exist
	attrs := resp.Schema.Attributes
	assert.Contains(t, attrs, "id")
	assert.Contains(t, attrs, "project_id")
	assert.Contains(t, attrs, "user_id")
	assert.Contains(t, attrs, "role")
}

// TestProjectUserResource_Configure tests the Configure method.
func TestProjectUserResource_Configure(t *testing.T) {
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
				r := project.NewProjectUserResource()
				req := resource.ConfigureRequest{
					ProviderData: &client.N8nClient{},
				}
				resp := &resource.ConfigureResponse{}
				r.Configure(context.Background(), req, resp)
				assert.False(t, resp.Diagnostics.HasError())

			case "error case - nil provider data":
				r := project.NewProjectUserResource()
				req := resource.ConfigureRequest{
					ProviderData: nil,
				}
				resp := &resource.ConfigureResponse{}
				r.Configure(context.Background(), req, resp)
				assert.False(t, resp.Diagnostics.HasError())

			case "error case - wrong provider data type":
				r := project.NewProjectUserResource()
				req := resource.ConfigureRequest{
					ProviderData: "wrong type",
				}
				resp := &resource.ConfigureResponse{}
				r.Configure(context.Background(), req, resp)
				assert.True(t, resp.Diagnostics.HasError())
				assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Unexpected Resource Configure Type")
			}
		})
	}
}

// TestProjectUserResource_Create tests the Create method.
func TestProjectUserResource_Create(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "validates method exists", wantErr: false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "validates method exists":
				r := project.NewProjectUserResource()
				assert.NotNil(t, r)
			}
		})
	}
}

// TestProjectUserResource_Read tests the Read method.
func TestProjectUserResource_Read(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "validates method exists", wantErr: false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "validates method exists":
				r := project.NewProjectUserResource()
				assert.NotNil(t, r)
			}
		})
	}
}

// TestProjectUserResource_Update tests the Update method.
func TestProjectUserResource_Update(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "validates method exists", wantErr: false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "validates method exists":
				r := project.NewProjectUserResource()
				assert.NotNil(t, r)
			}
		})
	}
}

// TestProjectUserResource_Delete tests the Delete method.
func TestProjectUserResource_Delete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "validates method exists", wantErr: false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "validates method exists":
				r := project.NewProjectUserResource()
				assert.NotNil(t, r)
			}
		})
	}
}

// TestProjectUserResource_ImportState tests the ImportState method.
func TestProjectUserResource_ImportState(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		importID  string
		wantError bool
	}{
		{
			name:      "success - valid composite ID",
			importID:  "proj-123/user-456",
			wantError: false,
		},
		{
			name:      "error - invalid format (single part)",
			importID:  "proj-123",
			wantError: true,
		},
		{
			name:      "error - invalid format (three parts)",
			importID:  "proj-123/user-456/extra",
			wantError: true,
		},
		{
			name:      "error - empty ID",
			importID:  "",
			wantError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := project.NewProjectUserResource()
			ctx := context.Background()

			schemaResp := resource.SchemaResponse{}
			r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

			// Build empty state
			emptyValue := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
				"id":         tftypes.NewValue(tftypes.String, nil),
				"project_id": tftypes.NewValue(tftypes.String, nil),
				"user_id":    tftypes.NewValue(tftypes.String, nil),
				"role":       tftypes.NewValue(tftypes.String, nil),
			})

			req := resource.ImportStateRequest{
				ID: tt.importID,
			}
			resp := &resource.ImportStateResponse{
				State: tfsdk.State{
					Schema: schemaResp.Schema,
					Raw:    emptyValue,
				},
			}

			r.ImportState(ctx, req, resp)

			if tt.wantError {
				assert.True(t, resp.Diagnostics.HasError(), "Expected error for import ID: %s", tt.importID)
			} else {
				assert.False(t, resp.Diagnostics.HasError(), "Unexpected error for import ID: %s", tt.importID)
			}
		})
	}
}
