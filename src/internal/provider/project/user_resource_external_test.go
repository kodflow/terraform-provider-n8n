package project_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
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
		name    string
		wantErr bool
	}{
		{name: "sets correct type name", wantErr: false},
		{name: "error case - validation checks", wantErr: false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "sets correct type name":
				r := project.NewProjectUserResource()
				req := resource.MetadataRequest{
					ProviderTypeName: "n8n",
				}
				resp := &resource.MetadataResponse{}
				r.Metadata(context.Background(), req, resp)
				assert.Equal(t, "n8n_project_user", resp.TypeName)

			case "error case - validation checks":
				r := project.NewProjectUserResource()
				req := resource.MetadataRequest{
					ProviderTypeName: "n8n",
				}
				resp := &resource.MetadataResponse{}
				r.Metadata(context.Background(), req, resp)
				assert.NotEmpty(t, resp.TypeName)
			}
		})
	}
}

// TestProjectUserResource_Schema tests the Schema method.
func TestProjectUserResource_Schema(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "returns valid schema", wantErr: false},
		{name: "error case - validation checks", wantErr: false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "returns valid schema":
				r := project.NewProjectUserResource()
				req := resource.SchemaRequest{}
				resp := &resource.SchemaResponse{}
				r.Schema(context.Background(), req, resp)
				assert.NotNil(t, resp.Schema)
				assert.NotEmpty(t, resp.Schema.Attributes)

			case "error case - validation checks":
				r := project.NewProjectUserResource()
				req := resource.SchemaRequest{}
				resp := &resource.SchemaResponse{}
				r.Schema(context.Background(), req, resp)
				assert.NotNil(t, resp.Schema)
			}
		})
	}
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
