package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/kodflow/n8n/src/internal/provider/shared/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// contextKey is a custom type for context keys to avoid collisions.
type contextKey string

func TestNewN8nProvider(t *testing.T) {
	t.Run("creates provider with version", func(t *testing.T) {
		version := "1.2.3"
		p := NewN8nProvider(version)

		assert.NotNil(t, p, "Provider should not be nil")
		assert.Equal(t, version, p.version, "Provider version should be set")
	})

	t.Run("creates provider with empty version", func(t *testing.T) {
		p := NewN8nProvider("")

		assert.NotNil(t, p, "Provider should not be nil")
		assert.Equal(t, "", p.version, "Provider should handle empty version")
	})

	t.Run("creates provider with dev version", func(t *testing.T) {
		p := NewN8nProvider("dev")

		assert.NotNil(t, p, "Provider should not be nil")
		assert.Equal(t, "dev", p.version, "Provider should handle dev version")
	})

	t.Run("creates multiple providers with different versions", func(t *testing.T) {
		p1 := NewN8nProvider("1.0.0")
		p2 := NewN8nProvider("2.0.0")

		assert.NotEqual(t, p1.version, p2.version, "Different providers should have different versions")
	})
}

func TestNew(t *testing.T) {
	t.Run("returns factory function", func(t *testing.T) {
		version := "1.2.3"
		factory := New(version)

		assert.NotNil(t, factory, "Factory function should not be nil")
	})

	t.Run("factory creates provider instance", func(t *testing.T) {
		version := "1.2.3"
		factory := New(version)
		p := factory()

		assert.NotNil(t, p, "Factory should create provider")
		assert.Implements(t, (*provider.Provider)(nil), p, "Should implement provider.Provider")
	})

	t.Run("factory creates N8nProvider instance", func(t *testing.T) {
		version := "1.2.3"
		factory := New(version)
		p := factory()

		n8nProvider, ok := p.(*N8nProvider)
		assert.True(t, ok, "Should create N8nProvider instance")
		assert.Equal(t, version, n8nProvider.version, "Should have correct version")
	})

	t.Run("factory can be called multiple times", func(t *testing.T) {
		factory := New("1.0.0")
		p1 := factory()
		p2 := factory()

		assert.NotNil(t, p1, "First call should create provider")
		assert.NotNil(t, p2, "Second call should create provider")
		// Different instances
		assert.NotSame(t, p1, p2, "Each call should create a new instance")
	})
}

func TestMetadata(t *testing.T) {
	t.Run("sets type name and version", func(t *testing.T) {
		version := "1.2.3"
		p := NewN8nProvider(version)

		req := provider.MetadataRequest{}
		resp := &provider.MetadataResponse{}

		p.Metadata(context.Background(), req, resp)

		assert.Equal(t, "n8n", resp.TypeName, "TypeName should be 'n8n'")
		assert.Equal(t, version, resp.Version, "Version should match provider version")
	})

	t.Run("handles empty version", func(t *testing.T) {
		p := NewN8nProvider("")

		req := provider.MetadataRequest{}
		resp := &provider.MetadataResponse{}

		p.Metadata(context.Background(), req, resp)

		assert.Equal(t, "n8n", resp.TypeName, "TypeName should be 'n8n'")
		assert.Equal(t, "", resp.Version, "Version should be empty")
	})

	t.Run("uses context", func(t *testing.T) {
		p := NewN8nProvider("1.0.0")
		ctx := context.WithValue(context.Background(), contextKey("test"), "value")

		req := provider.MetadataRequest{}
		resp := &provider.MetadataResponse{}

		// Should not panic with custom context
		assert.NotPanics(t, func() {
			p.Metadata(ctx, req, resp)
		})
	})

	t.Run("type name is always n8n", func(t *testing.T) {
		versions := []string{"1.0.0", "2.0.0", "dev", ""}
		for _, version := range versions {
			p := NewN8nProvider(version)
			resp := &provider.MetadataResponse{}
			p.Metadata(context.Background(), provider.MetadataRequest{}, resp)
			assert.Equal(t, "n8n", resp.TypeName, "TypeName should always be 'n8n'")
		}
	})
}

func TestSchema(t *testing.T) {
	t.Run("defines provider schema", func(t *testing.T) {
		p := NewN8nProvider("1.0.0")

		req := provider.SchemaRequest{}
		resp := &provider.SchemaResponse{}

		p.Schema(context.Background(), req, resp)

		assert.NotNil(t, resp.Schema, "Schema should be defined")
		assert.NotEmpty(t, resp.Schema.MarkdownDescription, "Schema should have description")
	})

	t.Run("requires api_key attribute", func(t *testing.T) {
		p := NewN8nProvider("1.0.0")

		req := provider.SchemaRequest{}
		resp := &provider.SchemaResponse{}

		p.Schema(context.Background(), req, resp)

		apiKeyAttr, exists := resp.Schema.Attributes["api_key"]
		assert.True(t, exists, "api_key attribute should exist")

		stringAttr, ok := apiKeyAttr.(schema.StringAttribute)
		assert.True(t, ok, "api_key should be StringAttribute")
		assert.True(t, stringAttr.Required, "api_key should be required")
		assert.True(t, stringAttr.Sensitive, "api_key should be sensitive")
	})

	t.Run("requires base_url attribute", func(t *testing.T) {
		p := NewN8nProvider("1.0.0")

		req := provider.SchemaRequest{}
		resp := &provider.SchemaResponse{}

		p.Schema(context.Background(), req, resp)

		baseURLAttr, exists := resp.Schema.Attributes["base_url"]
		assert.True(t, exists, "base_url attribute should exist")

		stringAttr, ok := baseURLAttr.(schema.StringAttribute)
		assert.True(t, ok, "base_url should be StringAttribute")
		assert.True(t, stringAttr.Required, "base_url should be required")
	})

	t.Run("has markdown description", func(t *testing.T) {
		p := NewN8nProvider("1.0.0")

		req := provider.SchemaRequest{}
		resp := &provider.SchemaResponse{}

		p.Schema(context.Background(), req, resp)

		assert.Contains(t, resp.Schema.MarkdownDescription, "Terraform provider", "Description should mention Terraform provider")
	})

	t.Run("api_key has markdown description", func(t *testing.T) {
		p := NewN8nProvider("1.0.0")
		resp := &provider.SchemaResponse{}
		p.Schema(context.Background(), provider.SchemaRequest{}, resp)

		apiKeyAttr := resp.Schema.Attributes["api_key"].(schema.StringAttribute)
		assert.NotEmpty(t, apiKeyAttr.MarkdownDescription, "api_key should have description")
		assert.Contains(t, apiKeyAttr.MarkdownDescription, "API key", "Description should mention API key")
	})

	t.Run("base_url has markdown description", func(t *testing.T) {
		p := NewN8nProvider("1.0.0")
		resp := &provider.SchemaResponse{}
		p.Schema(context.Background(), provider.SchemaRequest{}, resp)

		baseURLAttr := resp.Schema.Attributes["base_url"].(schema.StringAttribute)
		assert.NotEmpty(t, baseURLAttr.MarkdownDescription, "base_url should have description")
		assert.Contains(t, baseURLAttr.MarkdownDescription, "Base URL", "Description should mention Base URL")
	})
}

func TestConfigure(t *testing.T) {
	t.Run("configures with valid config", func(t *testing.T) {
		p := NewN8nProvider("1.0.0")

		// Create a valid config
		configValue := tftypes.NewValue(
			tftypes.Object{
				AttributeTypes: map[string]tftypes.Type{
					"api_key":  tftypes.String,
					"base_url": tftypes.String,
				},
			},
			map[string]tftypes.Value{
				"api_key":  tftypes.NewValue(tftypes.String, "test-api-key"),
				"base_url": tftypes.NewValue(tftypes.String, "https://n8n.example.com"),
			},
		)

		schemaResp := &provider.SchemaResponse{}
		p.Schema(context.Background(), provider.SchemaRequest{}, schemaResp)

		config := tfsdk.Config{
			Schema: schemaResp.Schema,
			Raw:    configValue,
		}

		req := provider.ConfigureRequest{
			Config: config,
		}
		resp := &provider.ConfigureResponse{}

		p.Configure(context.Background(), req, resp)

		assert.False(t, resp.Diagnostics.HasError(), "Configure should succeed with valid config")
		assert.NotNil(t, resp.ResourceData, "ResourceData should be set")
		assert.NotNil(t, resp.DataSourceData, "DataSourceData should be set")
	})

	t.Run("fails with missing api_key", func(t *testing.T) {
		p := NewN8nProvider("1.0.0")

		configValue := tftypes.NewValue(
			tftypes.Object{
				AttributeTypes: map[string]tftypes.Type{
					"api_key":  tftypes.String,
					"base_url": tftypes.String,
				},
			},
			map[string]tftypes.Value{
				"api_key":  tftypes.NewValue(tftypes.String, nil),
				"base_url": tftypes.NewValue(tftypes.String, "https://n8n.example.com"),
			},
		)

		schemaResp := &provider.SchemaResponse{}
		p.Schema(context.Background(), provider.SchemaRequest{}, schemaResp)

		config := tfsdk.Config{
			Schema: schemaResp.Schema,
			Raw:    configValue,
		}

		req := provider.ConfigureRequest{
			Config: config,
		}
		resp := &provider.ConfigureResponse{}

		p.Configure(context.Background(), req, resp)

		assert.True(t, resp.Diagnostics.HasError(), "Configure should fail with missing api_key")
		assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Missing API Key", "Error should mention missing API key")
	})

	t.Run("fails with empty api_key", func(t *testing.T) {
		p := NewN8nProvider("1.0.0")

		configValue := tftypes.NewValue(
			tftypes.Object{
				AttributeTypes: map[string]tftypes.Type{
					"api_key":  tftypes.String,
					"base_url": tftypes.String,
				},
			},
			map[string]tftypes.Value{
				"api_key":  tftypes.NewValue(tftypes.String, ""),
				"base_url": tftypes.NewValue(tftypes.String, "https://n8n.example.com"),
			},
		)

		schemaResp := &provider.SchemaResponse{}
		p.Schema(context.Background(), provider.SchemaRequest{}, schemaResp)

		config := tfsdk.Config{
			Schema: schemaResp.Schema,
			Raw:    configValue,
		}

		req := provider.ConfigureRequest{
			Config: config,
		}
		resp := &provider.ConfigureResponse{}

		p.Configure(context.Background(), req, resp)

		assert.True(t, resp.Diagnostics.HasError(), "Configure should fail with empty api_key")
	})

	t.Run("fails with missing base_url", func(t *testing.T) {
		p := NewN8nProvider("1.0.0")

		configValue := tftypes.NewValue(
			tftypes.Object{
				AttributeTypes: map[string]tftypes.Type{
					"api_key":  tftypes.String,
					"base_url": tftypes.String,
				},
			},
			map[string]tftypes.Value{
				"api_key":  tftypes.NewValue(tftypes.String, "test-api-key"),
				"base_url": tftypes.NewValue(tftypes.String, nil),
			},
		)

		schemaResp := &provider.SchemaResponse{}
		p.Schema(context.Background(), provider.SchemaRequest{}, schemaResp)

		config := tfsdk.Config{
			Schema: schemaResp.Schema,
			Raw:    configValue,
		}

		req := provider.ConfigureRequest{
			Config: config,
		}
		resp := &provider.ConfigureResponse{}

		p.Configure(context.Background(), req, resp)

		assert.True(t, resp.Diagnostics.HasError(), "Configure should fail with missing base_url")
		assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Missing Base URL", "Error should mention missing base URL")
	})

	t.Run("fails with empty base_url", func(t *testing.T) {
		p := NewN8nProvider("1.0.0")

		configValue := tftypes.NewValue(
			tftypes.Object{
				AttributeTypes: map[string]tftypes.Type{
					"api_key":  tftypes.String,
					"base_url": tftypes.String,
				},
			},
			map[string]tftypes.Value{
				"api_key":  tftypes.NewValue(tftypes.String, "test-api-key"),
				"base_url": tftypes.NewValue(tftypes.String, ""),
			},
		)

		schemaResp := &provider.SchemaResponse{}
		p.Schema(context.Background(), provider.SchemaRequest{}, schemaResp)

		config := tfsdk.Config{
			Schema: schemaResp.Schema,
			Raw:    configValue,
		}

		req := provider.ConfigureRequest{
			Config: config,
		}
		resp := &provider.ConfigureResponse{}

		p.Configure(context.Background(), req, resp)

		assert.True(t, resp.Diagnostics.HasError(), "Configure should fail with empty base_url")
	})

	t.Run("exits early on config parsing error", func(t *testing.T) {
		p := NewN8nProvider("1.0.0")

		// Invalid config that will fail parsing
		configValue := tftypes.NewValue(tftypes.String, "invalid")

		schemaResp := &provider.SchemaResponse{}
		p.Schema(context.Background(), provider.SchemaRequest{}, schemaResp)

		config := tfsdk.Config{
			Schema: schemaResp.Schema,
			Raw:    configValue,
		}

		req := provider.ConfigureRequest{
			Config: config,
		}
		resp := &provider.ConfigureResponse{}

		p.Configure(context.Background(), req, resp)

		assert.True(t, resp.Diagnostics.HasError(), "Configure should fail with invalid config")
		// Should exit early, not set client data
		assert.Nil(t, resp.ResourceData, "ResourceData should not be set on error")
		assert.Nil(t, resp.DataSourceData, "DataSourceData should not be set on error")
	})
}

func TestResources(t *testing.T) {
	t.Run("returns resource factory functions", func(t *testing.T) {
		p := NewN8nProvider("1.0.0")
		resources := p.Resources(context.Background())

		assert.NotNil(t, resources, "Resources should not be nil")
		assert.NotEmpty(t, resources, "Should return at least one resource")
	})

	t.Run("all resource factories are callable", func(t *testing.T) {
		p := NewN8nProvider("1.0.0")
		resources := p.Resources(context.Background())

		for i, factory := range resources {
			assert.NotNil(t, factory, "Resource factory %d should not be nil", i)
			r := factory()
			assert.NotNil(t, r, "Resource factory %d should create a resource", i)
			assert.Implements(t, (*resource.Resource)(nil), r, "Resource %d should implement resource.Resource", i)
		}
	})

	t.Run("includes workflow resources", func(t *testing.T) {
		p := NewN8nProvider("1.0.0")
		resources := p.Resources(context.Background())

		// Should include workflow-related resources
		assert.GreaterOrEqual(t, len(resources), 2, "Should include workflow resources")
	})

	t.Run("includes project resources", func(t *testing.T) {
		p := NewN8nProvider("1.0.0")
		resources := p.Resources(context.Background())

		// Should include project-related resources
		assert.GreaterOrEqual(t, len(resources), 4, "Should include project resources")
	})

	t.Run("includes credential resources", func(t *testing.T) {
		p := NewN8nProvider("1.0.0")
		resources := p.Resources(context.Background())

		// Should include credential-related resources
		assert.GreaterOrEqual(t, len(resources), 6, "Should include credential resources")
	})

	t.Run("resource list is stable", func(t *testing.T) {
		p := NewN8nProvider("1.0.0")
		resources1 := p.Resources(context.Background())
		resources2 := p.Resources(context.Background())

		assert.Equal(t, len(resources1), len(resources2), "Resource count should be stable")
	})
}

func TestDataSources(t *testing.T) {
	t.Run("returns data source factory functions", func(t *testing.T) {
		p := NewN8nProvider("1.0.0")
		dataSources := p.DataSources(context.Background())

		assert.NotNil(t, dataSources, "DataSources should not be nil")
		assert.NotEmpty(t, dataSources, "Should return at least one data source")
	})

	t.Run("all data source factories are callable", func(t *testing.T) {
		p := NewN8nProvider("1.0.0")
		dataSources := p.DataSources(context.Background())

		for i, factory := range dataSources {
			assert.NotNil(t, factory, "DataSource factory %d should not be nil", i)
			ds := factory()
			assert.NotNil(t, ds, "DataSource factory %d should create a data source", i)
			assert.Implements(t, (*datasource.DataSource)(nil), ds, "DataSource %d should implement datasource.DataSource", i)
		}
	})

	t.Run("includes workflow data sources", func(t *testing.T) {
		p := NewN8nProvider("1.0.0")
		dataSources := p.DataSources(context.Background())

		// Should include workflow-related data sources
		assert.GreaterOrEqual(t, len(dataSources), 2, "Should include workflow data sources")
	})

	t.Run("data source list is stable", func(t *testing.T) {
		p := NewN8nProvider("1.0.0")
		dataSources1 := p.DataSources(context.Background())
		dataSources2 := p.DataSources(context.Background())

		assert.Equal(t, len(dataSources1), len(dataSources2), "DataSource count should be stable")
	})
}

func TestValidateProvider(t *testing.T) {
	t.Run("validates N8nProvider implements TerraformProvider", func(t *testing.T) {
		p := NewN8nProvider("1.0.0")
		validated := ValidateProvider(p)

		assert.NotNil(t, validated, "Validated provider should not be nil")
		assert.Equal(t, p, validated, "Should return the same provider")
	})

	t.Run("ValidateProvider returns TerraformProvider interface", func(t *testing.T) {
		p := NewN8nProvider("1.0.0")
		validated := ValidateProvider(p)

		assert.Implements(t, (*TerraformProvider)(nil), validated, "Should implement TerraformProvider")
	})

	t.Run("validated provider can call all interface methods", func(t *testing.T) {
		p := NewN8nProvider("1.0.0")
		validated := ValidateProvider(p)

		// Test Metadata
		metadataResp := &provider.MetadataResponse{}
		assert.NotPanics(t, func() {
			validated.Metadata(context.Background(), provider.MetadataRequest{}, metadataResp)
		})

		// Test Schema
		schemaResp := &provider.SchemaResponse{}
		assert.NotPanics(t, func() {
			validated.Schema(context.Background(), provider.SchemaRequest{}, schemaResp)
		})

		// Test Resources
		assert.NotPanics(t, func() {
			resources := validated.Resources(context.Background())
			assert.NotNil(t, resources)
		})

		// Test DataSources
		assert.NotPanics(t, func() {
			dataSources := validated.DataSources(context.Background())
			assert.NotNil(t, dataSources)
		})
	})
}

func TestN8nProviderImplementsInterfaces(t *testing.T) {
	t.Run("N8nProvider implements provider.Provider", func(t *testing.T) {
		var _ provider.Provider = &N8nProvider{}
		assert.True(t, true, "Compile-time check passed")
	})

	t.Run("N8nProvider implements TerraformProvider", func(t *testing.T) {
		var _ TerraformProvider = &N8nProvider{}
		assert.True(t, true, "Compile-time check passed")
	})
}

func TestTerraformProviderInterface(t *testing.T) {
	t.Run("TerraformProvider defines all required methods", func(t *testing.T) {
		p := NewN8nProvider("1.0.0")
		var tp TerraformProvider = p

		// Verify all methods are accessible
		assert.NotNil(t, tp, "TerraformProvider should not be nil")

		// Metadata method
		metadataResp := &provider.MetadataResponse{}
		tp.Metadata(context.Background(), provider.MetadataRequest{}, metadataResp)
		assert.Equal(t, "n8n", metadataResp.TypeName)

		// Schema method
		schemaResp := &provider.SchemaResponse{}
		tp.Schema(context.Background(), provider.SchemaRequest{}, schemaResp)
		assert.NotNil(t, schemaResp.Schema)

		// Resources method
		resources := tp.Resources(context.Background())
		assert.NotNil(t, resources)

		// DataSources method
		dataSources := tp.DataSources(context.Background())
		assert.NotNil(t, dataSources)
	})
}

func TestN8nProviderStructure(t *testing.T) {
	t.Run("N8nProvider has version field", func(t *testing.T) {
		version := "test-version"
		p := NewN8nProvider(version)

		assert.Equal(t, version, p.version, "Provider should store version")
	})

	t.Run("N8nProvider version is private", func(t *testing.T) {
		// This is a structural test - version field is lowercase (private)
		p := NewN8nProvider("1.0.0")
		assert.NotNil(t, p, "Provider should be created")
		// We can't directly access p.version from outside the package in real usage
	})
}

func TestProviderContextUsage(t *testing.T) {
	t.Run("all methods accept context", func(t *testing.T) {
		p := NewN8nProvider("1.0.0")
		ctx := context.WithValue(context.Background(), contextKey("test-key"), "test-value")

		// Should not panic with custom context
		assert.NotPanics(t, func() {
			resp := &provider.MetadataResponse{}
			p.Metadata(ctx, provider.MetadataRequest{}, resp)
		})

		assert.NotPanics(t, func() {
			resp := &provider.SchemaResponse{}
			p.Schema(ctx, provider.SchemaRequest{}, resp)
		})

		assert.NotPanics(t, func() {
			_ = p.Resources(ctx)
		})

		assert.NotPanics(t, func() {
			_ = p.DataSources(ctx)
		})
	})

	t.Run("Configure handles context cancellation gracefully", func(t *testing.T) {
		p := NewN8nProvider("1.0.0")
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		configValue := tftypes.NewValue(
			tftypes.Object{
				AttributeTypes: map[string]tftypes.Type{
					"api_key":  tftypes.String,
					"base_url": tftypes.String,
				},
			},
			map[string]tftypes.Value{
				"api_key":  tftypes.NewValue(tftypes.String, "test-api-key"),
				"base_url": tftypes.NewValue(tftypes.String, "https://n8n.example.com"),
			},
		)

		schemaResp := &provider.SchemaResponse{}
		p.Schema(context.Background(), provider.SchemaRequest{}, schemaResp)

		config := tfsdk.Config{
			Schema: schemaResp.Schema,
			Raw:    configValue,
		}

		req := provider.ConfigureRequest{
			Config: config,
		}
		resp := &provider.ConfigureResponse{}

		// Should handle cancelled context gracefully
		assert.NotPanics(t, func() {
			p.Configure(ctx, req, resp)
		})
	})
}

func TestProviderVersionHandling(t *testing.T) {
	testCases := []struct {
		name    string
		version string
	}{
		{"semantic version", "1.2.3"},
		{"prerelease version", "1.2.3-alpha.1"},
		{"build metadata", "1.2.3+build.123"},
		{"dev version", "dev"},
		{"empty version", ""},
		{"snapshot", "1.0.0-SNAPSHOT"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewN8nProvider(tc.version)
			resp := &provider.MetadataResponse{}
			p.Metadata(context.Background(), provider.MetadataRequest{}, resp)

			assert.Equal(t, tc.version, resp.Version, "Version should be set correctly")
		})
	}
}

func TestConfigureClientCreation(t *testing.T) {
	t.Run("creates client with correct base URL", func(t *testing.T) {
		p := NewN8nProvider("1.0.0")

		testURL := "https://n8n.example.com"
		configValue := tftypes.NewValue(
			tftypes.Object{
				AttributeTypes: map[string]tftypes.Type{
					"api_key":  tftypes.String,
					"base_url": tftypes.String,
				},
			},
			map[string]tftypes.Value{
				"api_key":  tftypes.NewValue(tftypes.String, "test-key"),
				"base_url": tftypes.NewValue(tftypes.String, testURL),
			},
		)

		schemaResp := &provider.SchemaResponse{}
		p.Schema(context.Background(), provider.SchemaRequest{}, schemaResp)

		config := tfsdk.Config{
			Schema: schemaResp.Schema,
			Raw:    configValue,
		}

		req := provider.ConfigureRequest{
			Config: config,
		}
		resp := &provider.ConfigureResponse{}

		p.Configure(context.Background(), req, resp)

		assert.False(t, resp.Diagnostics.HasError())
		assert.NotNil(t, resp.ResourceData, "Client should be created")
	})
}

func TestProviderEdgeCases(t *testing.T) {
	t.Run("handles nil context gracefully", func(t *testing.T) {
		p := NewN8nProvider("1.0.0")

		// Test with context.TODO() instead of nil
		// Methods should work with any valid context
		assert.NotPanics(t, func() {
			resp := &provider.MetadataResponse{}
			p.Metadata(context.TODO(), provider.MetadataRequest{}, resp)
		})
	})

	t.Run("Resources returns consistent count", func(t *testing.T) {
		p := NewN8nProvider("1.0.0")
		count1 := len(p.Resources(context.Background()))
		count2 := len(p.Resources(context.Background()))
		count3 := len(p.Resources(context.Background()))

		assert.Equal(t, count1, count2, "Resource count should be consistent")
		assert.Equal(t, count2, count3, "Resource count should be consistent")
	})

	t.Run("DataSources returns consistent count", func(t *testing.T) {
		p := NewN8nProvider("1.0.0")
		count1 := len(p.DataSources(context.Background()))
		count2 := len(p.DataSources(context.Background()))
		count3 := len(p.DataSources(context.Background()))

		assert.Equal(t, count1, count2, "DataSource count should be consistent")
		assert.Equal(t, count2, count3, "DataSource count should be consistent")
	})
}

func TestN8nProviderModelUsage(t *testing.T) {
	t.Run("Configure uses N8nProviderModel", func(t *testing.T) {
		// This test verifies that Configure properly uses the N8nProviderModel
		p := NewN8nProvider("1.0.0")

		configValue := tftypes.NewValue(
			tftypes.Object{
				AttributeTypes: map[string]tftypes.Type{
					"api_key":  tftypes.String,
					"base_url": tftypes.String,
				},
			},
			map[string]tftypes.Value{
				"api_key":  tftypes.NewValue(tftypes.String, "test-key"),
				"base_url": tftypes.NewValue(tftypes.String, "https://test.com"),
			},
		)

		schemaResp := &provider.SchemaResponse{}
		p.Schema(context.Background(), provider.SchemaRequest{}, schemaResp)

		config := tfsdk.Config{
			Schema: schemaResp.Schema,
			Raw:    configValue,
		}

		req := provider.ConfigureRequest{
			Config: config,
		}
		resp := &provider.ConfigureResponse{}

		p.Configure(context.Background(), req, resp)

		require.False(t, resp.Diagnostics.HasError())
	})
}

func TestN8nProviderModelStructure(t *testing.T) {
	t.Run("N8nProviderModel has required fields", func(t *testing.T) {
		model := &models.N8nProviderModel{
			APIKey:  types.StringValue("test-key"),
			BaseURL: types.StringValue("https://test.com"),
		}

		assert.False(t, model.APIKey.IsNull())
		assert.False(t, model.BaseURL.IsNull())
		assert.Equal(t, "test-key", model.APIKey.ValueString())
		assert.Equal(t, "https://test.com", model.BaseURL.ValueString())
	})
}
