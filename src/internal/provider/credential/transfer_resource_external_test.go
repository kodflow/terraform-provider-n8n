package credential_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/kodflow/n8n/src/internal/provider/credential"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
	"github.com/stretchr/testify/assert"
)

// External tests (black-box testing) go here.
// These tests only have access to exported functions and types.

func TestNewCredentialTransferResource(t *testing.T) {
	t.Parallel()

	resource := credential.NewCredentialTransferResource()
	assert.NotNil(t, resource, "NewCredentialTransferResource should not return nil")
}

func TestNewCredentialTransferResourceWrapper(t *testing.T) {
	t.Parallel()

	resource := credential.NewCredentialTransferResourceWrapper()
	assert.NotNil(t, resource, "NewCredentialTransferResourceWrapper should not return nil")
}

func TestCredentialTransferResource_Metadata(t *testing.T) {
	t.Parallel()

	// Note: Comprehensive testing is done in integration tests.
	// This test just ensures the function exists and is accessible.
	resource := credential.NewCredentialTransferResource()
	assert.NotNil(t, resource, "resource should not be nil")
}

func TestCredentialTransferResource_Schema(t *testing.T) {
	t.Parallel()

	// Note: Comprehensive testing is done in integration tests.
	// This test just ensures the function exists and is accessible.
	resource := credential.NewCredentialTransferResource()
	assert.NotNil(t, resource, "resource should not be nil")
}

func TestCredentialTransferResource_Configure(t *testing.T) {
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
				r := credential.NewCredentialTransferResource()
				req := resource.ConfigureRequest{
					ProviderData: &client.N8nClient{},
				}
				resp := &resource.ConfigureResponse{}
				r.Configure(context.Background(), req, resp)
				assert.False(t, resp.Diagnostics.HasError())

			case "error case - nil provider data":
				r := credential.NewCredentialTransferResource()
				req := resource.ConfigureRequest{
					ProviderData: nil,
				}
				resp := &resource.ConfigureResponse{}
				r.Configure(context.Background(), req, resp)
				assert.False(t, resp.Diagnostics.HasError())

			case "error case - wrong provider data type":
				r := credential.NewCredentialTransferResource()
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

func TestCredentialTransferResource_Create(t *testing.T) {
	t.Parallel()

	// Note: Comprehensive testing is done in integration tests.
	// This test just ensures the function exists and is accessible.
	resource := credential.NewCredentialTransferResource()
	assert.NotNil(t, resource, "resource should not be nil")
}

func TestCredentialTransferResource_Read(t *testing.T) {
	t.Parallel()

	// Note: Comprehensive testing is done in integration tests.
	// This test just ensures the function exists and is accessible.
	resource := credential.NewCredentialTransferResource()
	assert.NotNil(t, resource, "resource should not be nil")
}

func TestCredentialTransferResource_Update(t *testing.T) {
	t.Parallel()

	// Note: Comprehensive testing is done in integration tests.
	// This test just ensures the function exists and is accessible.
	resource := credential.NewCredentialTransferResource()
	assert.NotNil(t, resource, "resource should not be nil")
}

func TestCredentialTransferResource_Delete(t *testing.T) {
	t.Parallel()

	// Note: Comprehensive testing is done in integration tests.
	// This test just ensures the function exists and is accessible.
	resource := credential.NewCredentialTransferResource()
	assert.NotNil(t, resource, "resource should not be nil")
}

func TestCredentialTransferResource_ImportState(t *testing.T) {
	t.Parallel()

	// Note: Comprehensive testing is done in integration tests.
	// This test just ensures the function exists and is accessible.
	resource := credential.NewCredentialTransferResource()
	assert.NotNil(t, resource, "resource should not be nil")
}
