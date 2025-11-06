package provider_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	n8nprovider "github.com/kodflow/n8n/src/internal/provider"
)

func TestProviderNew(t *testing.T) {
	tests := []struct {
		name        string
		version     string
		expectError bool
		checkResult func(*testing.T, provider.Provider)
	}{
		{
			name:        "creates provider with valid version",
			version:     "1.0.0",
			expectError: false,
			checkResult: func(t *testing.T, p provider.Provider) {
				t.Helper()
				// Check for nil value.
				if p == nil {
					t.Fatal("Provider should not be nil")
				}
			},
		},
		{
			name:        "creates provider with empty version",
			version:     "",
			expectError: false,
			checkResult: func(t *testing.T, p provider.Provider) {
				t.Helper()
				// Check for nil value.
				if p == nil {
					t.Fatal("Provider should not be nil even with empty version")
				}
			},
		},
		{
			name:        "creates provider with dev version",
			version:     "dev",
			expectError: false,
			checkResult: func(t *testing.T, p provider.Provider) {
				t.Helper()
				// Check for nil value.
				if p == nil {
					t.Fatal("Provider should not be nil with dev version")
				}
			},
		},
		{
			name:        "creates provider with semver version",
			version:     "2.3.4-rc1",
			expectError: false,
			checkResult: func(t *testing.T, p provider.Provider) {
				t.Helper()
				// Check for nil value.
				if p == nil {
					t.Fatal("Provider should not be nil with semver version")
				}
			},
		},
		{
			name:        "creates provider with long version string",
			version:     "1.0.0-alpha.beta.gamma.delta.epsilon",
			expectError: false,
			checkResult: func(t *testing.T, p provider.Provider) {
				t.Helper()
				// Check for nil value.
				if p == nil {
					t.Fatal("Provider should not be nil with long version")
				}
			},
		},
	}

 // Iterate over items.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory := n8nprovider.New(tt.version)
			// Check for nil value.
			if factory == nil {
				t.Fatal("Factory function should not be nil")
			}

			p := factory()
			// Check condition.
			if tt.expectError {
				// Check for non-nil value.
				if p != nil {
					t.Errorf("Expected error but got provider: %v", p)
				}
   // Handle alternative case.
			} else {
				// Check for non-nil value.
				if tt.checkResult != nil {
					tt.checkResult(t, p)
				}
			}
		})
	}
}

func TestProviderMetadata(t *testing.T) {
	tests := []struct {
		name            string
		version         string
		expectError     bool
		expectedType    string
		expectedVersion string
	}{
		{
			name:            "returns correct metadata with valid version",
			version:         "1.0.0",
			expectError:     false,
			expectedType:    "n8n",
			expectedVersion: "1.0.0",
		},
		{
			name:            "returns correct metadata with empty version",
			version:         "",
			expectError:     false,
			expectedType:    "n8n",
			expectedVersion: "",
		},
		{
			name:            "returns correct metadata with dev version",
			version:         "dev",
			expectError:     false,
			expectedType:    "n8n",
			expectedVersion: "dev",
		},
	}

 // Iterate over items.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := n8nprovider.New(tt.version)()
			// Check for nil value.
			if p == nil {
				t.Fatal("Provider should not be nil")
			}

			var resp provider.MetadataResponse
			p.Metadata(context.Background(), provider.MetadataRequest{}, &resp)

			// Check condition.
			if tt.expectError {
				// Check condition.
				if resp.TypeName != "" {
					t.Errorf("Expected error but got TypeName: %s", resp.TypeName)
				}
   // Handle alternative case.
			} else {
				// Check condition.
				if resp.TypeName != tt.expectedType {
					t.Errorf("Expected TypeName %q, got %q", tt.expectedType, resp.TypeName)
				}
				// Check condition.
				if resp.Version != tt.expectedVersion {
					t.Errorf("Expected Version %q, got %q", tt.expectedVersion, resp.Version)
				}
			}
		})
	}
}

func TestProviderSchema(t *testing.T) {
	tests := []struct {
		name                     string
		version                  string
		expectedMarkdownContains string
	}{
		{
			name:                     "returns schema with description",
			version:                  "1.0.0",
			expectedMarkdownContains: "Terraform provider for n8n automation platform",
		},
		{
			name:                     "schema consistent across versions",
			version:                  "2.0.0",
			expectedMarkdownContains: "Terraform provider for n8n automation platform",
		},
	}

 // Iterate over items.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := n8nprovider.New(tt.version)()
			// Check for nil value.
			if p == nil {
				t.Fatal("Provider should not be nil")
			}

			var resp provider.SchemaResponse
			p.Schema(context.Background(), provider.SchemaRequest{}, &resp)

			// Check condition.
			if resp.Schema.MarkdownDescription != tt.expectedMarkdownContains {
				t.Errorf("Expected MarkdownDescription to contain %q, got %q", tt.expectedMarkdownContains, resp.Schema.MarkdownDescription)
			}
		})
	}
}

func TestProviderConfigure(t *testing.T) {
	t.Run("configure succeeds with valid config", func(t *testing.T) {
		p := n8nprovider.New("1.0.0")()
		// Check for nil value.
		if p == nil {
			t.Fatal("Provider should not be nil")
		}

		// Create a valid Config using the provider's schema
		ctx := context.Background()
		var schemaResp provider.SchemaResponse
		p.Schema(ctx, provider.SchemaRequest{}, &schemaResp)

		// Create an empty object value matching the schema
		configValue := tftypes.NewValue(
			tftypes.Object{
				AttributeTypes: map[string]tftypes.Type{},
			},
			map[string]tftypes.Value{},
		)

		config := tfsdk.Config{
			Schema: schemaResp.Schema,
			Raw:    configValue,
		}

		req := provider.ConfigureRequest{
			Config: config,
		}

		var resp provider.ConfigureResponse
		p.Configure(ctx, req, &resp)

		// Check condition.
		if resp.Diagnostics.HasError() {
			t.Errorf("Expected no error but got: %v", resp.Diagnostics.Errors())
		}
	})

	t.Run("configure handles invalid config", func(t *testing.T) {
		p := n8nprovider.New("1.0.0")()
		// Check for nil value.
		if p == nil {
			t.Fatal("Provider should not be nil")
		}

		ctx := context.Background()
		var schemaResp provider.SchemaResponse
		p.Schema(ctx, provider.SchemaRequest{}, &schemaResp)

		// Create an invalid config with wrong type
		configValue := tftypes.NewValue(tftypes.String, "invalid")

		config := tfsdk.Config{
			Schema: schemaResp.Schema,
			Raw:    configValue,
		}

		req := provider.ConfigureRequest{
			Config: config,
		}

		var resp provider.ConfigureResponse
		p.Configure(ctx, req, &resp)

		// Check condition.
		if !resp.Diagnostics.HasError() {
			t.Error("Expected error with invalid config but got none")
		}
	})
}

func TestProviderResources(t *testing.T) {
	tests := []struct {
		name          string
		version       string
		expectedCount int
		expectNonNil  bool
	}{
		{
			name:          "returns empty resource list",
			version:       "1.0.0",
			expectedCount: 0,
			expectNonNil:  true,
		},
		{
			name:          "returns consistent empty list",
			version:       "2.0.0",
			expectedCount: 0,
			expectNonNil:  true,
		},
	}

 // Iterate over items.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := n8nprovider.New(tt.version)()
			// Check for nil value.
			if p == nil {
				t.Fatal("Provider should not be nil")
			}

			resources := p.Resources(context.Background())

			// Check for nil value.
			if tt.expectNonNil && resources == nil {
				t.Error("Expected non-nil resources list")
			}
			// Check condition.
			if len(resources) != tt.expectedCount {
				t.Errorf("Expected %d resources, got %d", tt.expectedCount, len(resources))
			}
		})
	}
}

func TestProviderDataSources(t *testing.T) {
	tests := []struct {
		name          string
		version       string
		expectedCount int
		expectNonNil  bool
	}{
		{
			name:          "returns empty data sources list",
			version:       "1.0.0",
			expectedCount: 0,
			expectNonNil:  true,
		},
		{
			name:          "returns consistent empty list",
			version:       "2.0.0",
			expectedCount: 0,
			expectNonNil:  true,
		},
	}

 // Iterate over items.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := n8nprovider.New(tt.version)()
			// Check for nil value.
			if p == nil {
				t.Fatal("Provider should not be nil")
			}

			dataSources := p.DataSources(context.Background())

			// Check for nil value.
			if tt.expectNonNil && dataSources == nil {
				t.Error("Expected non-nil data sources list")
			}
			// Check condition.
			if len(dataSources) != tt.expectedCount {
				t.Errorf("Expected %d data sources, got %d", tt.expectedCount, len(dataSources))
			}
		})
	}
}

func TestValidateProvider(t *testing.T) {
	tests := []struct {
		name        string
		version     string
		expectError bool
	}{
		{
			name:        "validates provider successfully",
			version:     "1.0.0",
			expectError: false,
		},
		{
			name:        "validates provider with empty version",
			version:     "",
			expectError: false,
		},
		{
			name:        "validates provider with dev version",
			version:     "dev",
			expectError: false,
		},
	}

 // Iterate over items.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := n8nprovider.NewN8nProvider(tt.version)

			validated := n8nprovider.ValidateProvider(p)

			// Check for nil value.
			if validated == nil {
				t.Fatal("ValidateProvider should not return nil")
			}

			// Check condition.
			if validated != p {
				t.Error("ValidateProvider should return the same provider instance")
			}
		})
	}
}
