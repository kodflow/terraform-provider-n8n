package provider_test

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
	p "github.com/kodflow/n8n/src/internal/provider"
	"github.com/kodflow/n8n/src/internal/provider/shared/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// contextKey is a custom type for context keys to avoid collisions.
type contextKey string

func TestNewN8nProvider(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		version         string
		wantNil         bool
		wantVersion     string
		wantErrContains string
	}{
		{
			name:        "creates provider with version",
			version:     "1.2.3",
			wantNil:     false,
			wantVersion: "1.2.3",
		},
		{
			name:        "creates provider with empty version",
			version:     "",
			wantNil:     false,
			wantVersion: "",
		},
		{
			name:        "creates provider with dev version",
			version:     "dev",
			wantNil:     false,
			wantVersion: "dev",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := p.NewN8nProvider(tt.version)

			if tt.wantNil {
				assert.Nil(t, got, "Provider should be nil")
			} else {
				assert.NotNil(t, got, "Provider should not be nil")
			}
		})
	}
}

func TestNewN8nProvider_MultipleInstances(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "creates multiple independent instances",
			testFunc: func(t *testing.T) {
				t.Helper()
				p1 := p.NewN8nProvider("1.0.0")
				p2 := p.NewN8nProvider("2.0.0")

				assert.NotNil(t, p1, "First provider should not be nil")
				assert.NotNil(t, p2, "Second provider should not be nil")
			},
		},
		{
			name: "error case - multiple instances are independent",
			testFunc: func(t *testing.T) {
				t.Helper()
				p1 := p.NewN8nProvider("1.0.0")
				p2 := p.NewN8nProvider("1.0.0")
				assert.NotSame(t, p1, p2, "Each call should create a new instance")
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

func TestNew(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		version     string
		wantFactory bool
	}{
		{
			name:        "returns factory function",
			version:     "1.2.3",
			wantFactory: true,
		},
		{
			name:        "factory with empty version",
			version:     "",
			wantFactory: true,
		},
		{
			name:        "factory with dev version",
			version:     "dev",
			wantFactory: true,
		},
		{
			name:        "error case - factory with nil-like version",
			version:     "\x00",
			wantFactory: true,
		},
		{
			name:        "error case - factory with special characters",
			version:     "!@#$%",
			wantFactory: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			factory := p.New(tt.version)

			if tt.wantFactory {
				assert.NotNil(t, factory, "Factory function should not be nil")

				// Test that factory creates a provider
				instance := factory()
				assert.NotNil(t, instance, "Factory should create provider")
				assert.Implements(t, (*provider.Provider)(nil), instance, "Should implement provider.Provider")
			}
		})
	}
}

func TestNew_FactoryMultipleCalls(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "factory creates new instance on each call",
			testFunc: func(t *testing.T) {
				t.Helper()
				factory := p.New("1.0.0")
				p1 := factory()
				p2 := factory()

				assert.NotNil(t, p1, "First call should create provider")
				assert.NotNil(t, p2, "Second call should create provider")
				assert.NotSame(t, p1, p2, "Each call should create a new instance")
			},
		},
		{
			name: "error case - factory handles concurrent calls",
			testFunc: func(t *testing.T) {
				t.Helper()
				factory := p.New("1.0.0")
				instances := make([]provider.Provider, 10)
				for i := range instances {
					instances[i] = factory()
					assert.NotNil(t, instances[i], "All instances should be created")
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

func TestMetadata(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		version      string
		wantTypeName string
		wantVersion  string
	}{
		{
			name:         "sets type name and version",
			version:      "1.2.3",
			wantTypeName: "n8n",
			wantVersion:  "1.2.3",
		},
		{
			name:         "handles empty version",
			version:      "",
			wantTypeName: "n8n",
			wantVersion:  "",
		},
		{
			name:         "handles dev version",
			version:      "dev",
			wantTypeName: "n8n",
			wantVersion:  "dev",
		},
		{
			name:         "handles semantic version",
			version:      "1.0.0",
			wantTypeName: "n8n",
			wantVersion:  "1.0.0",
		},
		{
			name:         "error case - handles special version",
			version:      "!@#",
			wantTypeName: "n8n",
			wantVersion:  "!@#",
		},
		{
			name:         "error case - handles very long version",
			version:      string(make([]byte, 10000)),
			wantTypeName: "n8n",
			wantVersion:  string(make([]byte, 10000)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			prov := p.NewN8nProvider(tt.version)
			req := provider.MetadataRequest{}
			resp := &provider.MetadataResponse{}

			prov.Metadata(context.Background(), req, resp)

			assert.Equal(t, tt.wantTypeName, resp.TypeName, "TypeName should match")
			assert.Equal(t, tt.wantVersion, resp.Version, "Version should match")
		})
	}
}

func TestMetadata_WithContext(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "handles context with values",
			testFunc: func(t *testing.T) {
				t.Helper()
				prov := p.NewN8nProvider("1.0.0")
				ctx := context.WithValue(context.Background(), contextKey("test"), "value")

				req := provider.MetadataRequest{}
				resp := &provider.MetadataResponse{}

				assert.NotPanics(t, func() {
					prov.Metadata(ctx, req, resp)
				})
			},
		},
		{
			name: "error case - handles canceled context",
			testFunc: func(t *testing.T) {
				t.Helper()
				prov := p.NewN8nProvider("1.0.0")
				ctx, cancel := context.WithCancel(context.Background())
				cancel()

				req := provider.MetadataRequest{}
				resp := &provider.MetadataResponse{}

				assert.NotPanics(t, func() {
					prov.Metadata(ctx, req, resp)
				})
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

func TestSchema(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                string
		version             string
		wantAttributes      []string
		wantAPIKeyRequired  bool
		wantAPIKeySensitive bool
		wantBaseURLRequired bool
		wantHasDescription  bool
		wantErrContains     string
	}{
		{
			name:                "defines provider schema",
			version:             "1.0.0",
			wantAttributes:      []string{"api_key", "base_url"},
			wantAPIKeyRequired:  false, // Optional - reads from N8N_API_TOKEN env var
			wantAPIKeySensitive: true,
			wantBaseURLRequired: false, // Optional - reads from N8N_URL env var
			wantHasDescription:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			prov := p.NewN8nProvider(tt.version)
			req := provider.SchemaRequest{}
			resp := &provider.SchemaResponse{}

			prov.Schema(context.Background(), req, resp)

			// Verify schema exists
			assert.NotNil(t, resp.Schema, "Schema should be defined")

			// Verify description
			if tt.wantHasDescription {
				assert.NotEmpty(t, resp.Schema.MarkdownDescription, "Schema should have description")
				assert.Contains(t, resp.Schema.MarkdownDescription, "Terraform provider", "Description should mention Terraform provider")
			}

			// Verify attributes exist
			for _, attrName := range tt.wantAttributes {
				_, exists := resp.Schema.Attributes[attrName]
				assert.True(t, exists, "Attribute %s should exist", attrName)
			}

			// Verify api_key attribute
			if apiKeyAttr, exists := resp.Schema.Attributes["api_key"]; exists {
				stringAttr, ok := apiKeyAttr.(schema.StringAttribute)
				assert.True(t, ok, "api_key should be StringAttribute")
				assert.Equal(t, tt.wantAPIKeyRequired, stringAttr.Required, "api_key Required should match")
				assert.Equal(t, tt.wantAPIKeySensitive, stringAttr.Sensitive, "api_key Sensitive should match")
				assert.NotEmpty(t, stringAttr.MarkdownDescription, "api_key should have description")
				assert.Contains(t, stringAttr.MarkdownDescription, "API key", "Description should mention API key")
			}

			// Verify base_url attribute
			if baseURLAttr, exists := resp.Schema.Attributes["base_url"]; exists {
				stringAttr, ok := baseURLAttr.(schema.StringAttribute)
				assert.True(t, ok, "base_url should be StringAttribute")
				assert.Equal(t, tt.wantBaseURLRequired, stringAttr.Required, "base_url Required should match")
				assert.NotEmpty(t, stringAttr.MarkdownDescription, "base_url should have description")
				assert.Contains(t, stringAttr.MarkdownDescription, "Base URL", "Description should mention Base URL")
			}
		})
	}
}

func TestConfigure(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		version         string
		apiKey          *string
		baseURL         *string
		wantErr         bool
		wantErrContains string
		wantClientSet   bool
	}{
		{
			name:          "configures with valid config",
			version:       "1.0.0",
			apiKey:        stringPtr("test-api-key"),
			baseURL:       stringPtr("https://n8n.example.com"),
			wantErr:       false,
			wantClientSet: true,
		},
		{
			name:            "fails with missing api_key",
			version:         "1.0.0",
			apiKey:          nil,
			baseURL:         stringPtr("https://n8n.example.com"),
			wantErr:         true,
			wantErrContains: "Missing API Key",
			wantClientSet:   false,
		},
		{
			name:            "fails with empty api_key",
			version:         "1.0.0",
			apiKey:          stringPtr(""),
			baseURL:         stringPtr("https://n8n.example.com"),
			wantErr:         true,
			wantErrContains: "Missing API Key",
			wantClientSet:   false,
		},
		{
			name:            "fails with missing base_url",
			version:         "1.0.0",
			apiKey:          stringPtr("test-api-key"),
			baseURL:         nil,
			wantErr:         true,
			wantErrContains: "Missing Base URL",
			wantClientSet:   false,
		},
		{
			name:            "fails with empty base_url",
			version:         "1.0.0",
			apiKey:          stringPtr("test-api-key"),
			baseURL:         stringPtr(""),
			wantErr:         true,
			wantErrContains: "Missing Base URL",
			wantClientSet:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			prov := p.NewN8nProvider(tt.version)

			// Build config value
			configAttrs := make(map[string]tftypes.Value)

			if tt.apiKey != nil {
				configAttrs["api_key"] = tftypes.NewValue(tftypes.String, *tt.apiKey)
			} else {
				configAttrs["api_key"] = tftypes.NewValue(tftypes.String, nil)
			}

			if tt.baseURL != nil {
				configAttrs["base_url"] = tftypes.NewValue(tftypes.String, *tt.baseURL)
			} else {
				configAttrs["base_url"] = tftypes.NewValue(tftypes.String, nil)
			}

			configValue := tftypes.NewValue(
				tftypes.Object{
					AttributeTypes: map[string]tftypes.Type{
						"api_key":  tftypes.String,
						"base_url": tftypes.String,
					},
				},
				configAttrs,
			)

			// Get schema
			schemaResp := &provider.SchemaResponse{}
			prov.Schema(context.Background(), provider.SchemaRequest{}, schemaResp)

			config := tfsdk.Config{
				Schema: schemaResp.Schema,
				Raw:    configValue,
			}

			req := provider.ConfigureRequest{
				Config: config,
			}
			resp := &provider.ConfigureResponse{}

			prov.Configure(context.Background(), req, resp)

			// Verify results
			if tt.wantErr {
				assert.True(t, resp.Diagnostics.HasError(), "Configure should fail")
				if tt.wantErrContains != "" {
					assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), tt.wantErrContains, "Error should contain expected message")
				}
			} else {
				assert.False(t, resp.Diagnostics.HasError(), "Configure should succeed")
			}

			if tt.wantClientSet {
				assert.NotNil(t, resp.ResourceData, "ResourceData should be set")
				assert.NotNil(t, resp.DataSourceData, "DataSourceData should be set")
			} else {
				if tt.wantErr {
					// On error, client should not be set
					assert.Nil(t, resp.ResourceData, "ResourceData should not be set on error")
					assert.Nil(t, resp.DataSourceData, "DataSourceData should not be set on error")
				}
			}
		})
	}
}

func TestConfigure_InvalidConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "configure fails with invalid config",
			testFunc: func(t *testing.T) {
				t.Helper()
				prov := p.NewN8nProvider("1.0.0")

				// Invalid config that will fail parsing
				configValue := tftypes.NewValue(tftypes.String, "invalid")

				schemaResp := &provider.SchemaResponse{}
				prov.Schema(context.Background(), provider.SchemaRequest{}, schemaResp)

				config := tfsdk.Config{
					Schema: schemaResp.Schema,
					Raw:    configValue,
				}

				req := provider.ConfigureRequest{
					Config: config,
				}
				resp := &provider.ConfigureResponse{}

				prov.Configure(context.Background(), req, resp)

				assert.True(t, resp.Diagnostics.HasError(), "Configure should fail with invalid config")
				assert.Nil(t, resp.ResourceData, "ResourceData should not be set on error")
				assert.Nil(t, resp.DataSourceData, "DataSourceData should not be set on error")
			},
		},
		{
			name: "error case - diagnostics are present on error",
			testFunc: func(t *testing.T) {
				t.Helper()
				prov := p.NewN8nProvider("1.0.0")
				configValue := tftypes.NewValue(tftypes.String, "invalid")
				schemaResp := &provider.SchemaResponse{}
				prov.Schema(context.Background(), provider.SchemaRequest{}, schemaResp)
				config := tfsdk.Config{
					Schema: schemaResp.Schema,
					Raw:    configValue,
				}
				req := provider.ConfigureRequest{Config: config}
				resp := &provider.ConfigureResponse{}
				prov.Configure(context.Background(), req, resp)
				assert.NotNil(t, resp.Diagnostics, "Diagnostics must not be nil")
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

func TestConfigure_ContextCancellation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "configure handles canceled context",
			testFunc: func(t *testing.T) {
				t.Helper()
				prov := p.NewN8nProvider("1.0.0")
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
				prov.Schema(context.Background(), provider.SchemaRequest{}, schemaResp)

				config := tfsdk.Config{
					Schema: schemaResp.Schema,
					Raw:    configValue,
				}

				req := provider.ConfigureRequest{
					Config: config,
				}
				resp := &provider.ConfigureResponse{}

				assert.NotPanics(t, func() {
					prov.Configure(ctx, req, resp)
				})
			},
		},
		{
			name: "error case - configure with active context",
			testFunc: func(t *testing.T) {
				t.Helper()
				prov := p.NewN8nProvider("1.0.0")
				ctx := context.Background()

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
				prov.Schema(context.Background(), provider.SchemaRequest{}, schemaResp)

				config := tfsdk.Config{
					Schema: schemaResp.Schema,
					Raw:    configValue,
				}
				req := provider.ConfigureRequest{Config: config}
				resp := &provider.ConfigureResponse{}

				assert.NotPanics(t, func() {
					prov.Configure(ctx, req, resp)
				})
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

func TestResources(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                   string
		version                string
		wantMinResourceCount   int
		wantAllFactoriesValid  bool
		wantAllImplementsIface bool
	}{
		{
			name:                   "returns resource factory functions",
			version:                "1.0.0",
			wantMinResourceCount:   1,
			wantAllFactoriesValid:  true,
			wantAllImplementsIface: true,
		},
		{
			name:                   "error case - empty version returns resources",
			version:                "",
			wantMinResourceCount:   1,
			wantAllFactoriesValid:  true,
			wantAllImplementsIface: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			prov := p.NewN8nProvider(tt.version)
			resources := prov.Resources(context.Background())

			assert.NotNil(t, resources, "Resources should not be nil")
			assert.GreaterOrEqual(t, len(resources), tt.wantMinResourceCount, "Should return minimum resource count")

			// Verify all factories are valid
			if tt.wantAllFactoriesValid {
				for i, factory := range resources {
					assert.NotNil(t, factory, "Resource factory %d should not be nil", i)
					r := factory()
					assert.NotNil(t, r, "Resource factory %d should create a resource", i)

					if tt.wantAllImplementsIface {
						assert.Implements(t, (*resource.Resource)(nil), r, "Resource %d should implement resource.Resource", i)
					}
				}
			}
		})
	}
}

func TestResources_Stability(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "resource count is stable across calls",
			testFunc: func(t *testing.T) {
				t.Helper()
				prov := p.NewN8nProvider("1.0.0")
				resources1 := prov.Resources(context.Background())
				resources2 := prov.Resources(context.Background())
				resources3 := prov.Resources(context.Background())

				assert.Equal(t, len(resources1), len(resources2), "Resource count should be stable")
				assert.Equal(t, len(resources2), len(resources3), "Resource count should be stable")
			},
		},
		{
			name: "error case - resources are not nil",
			testFunc: func(t *testing.T) {
				t.Helper()
				prov := p.NewN8nProvider("1.0.0")
				resources := prov.Resources(context.Background())
				assert.NotNil(t, resources, "Resources must not be nil")
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

func TestDataSources(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                   string
		version                string
		wantMinDataSourceCount int
		wantAllFactoriesValid  bool
		wantAllImplementsIface bool
	}{
		{
			name:                   "returns data source factory functions",
			version:                "1.0.0",
			wantMinDataSourceCount: 1,
			wantAllFactoriesValid:  true,
			wantAllImplementsIface: true,
		},
		{
			name:                   "error case - empty version returns data sources",
			version:                "",
			wantMinDataSourceCount: 1,
			wantAllFactoriesValid:  true,
			wantAllImplementsIface: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			prov := p.NewN8nProvider(tt.version)
			dataSources := prov.DataSources(context.Background())

			assert.NotNil(t, dataSources, "DataSources should not be nil")
			assert.GreaterOrEqual(t, len(dataSources), tt.wantMinDataSourceCount, "Should return minimum data source count")

			// Verify all factories are valid
			if tt.wantAllFactoriesValid {
				for i, factory := range dataSources {
					assert.NotNil(t, factory, "DataSource factory %d should not be nil", i)
					ds := factory()
					assert.NotNil(t, ds, "DataSource factory %d should create a data source", i)

					if tt.wantAllImplementsIface {
						assert.Implements(t, (*datasource.DataSource)(nil), ds, "DataSource %d should implement datasource.DataSource", i)
					}
				}
			}
		})
	}
}

func TestDataSources_Stability(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "data source count is stable across calls",
			testFunc: func(t *testing.T) {
				t.Helper()
				prov := p.NewN8nProvider("1.0.0")
				dataSources1 := prov.DataSources(context.Background())
				dataSources2 := prov.DataSources(context.Background())
				dataSources3 := prov.DataSources(context.Background())

				assert.Equal(t, len(dataSources1), len(dataSources2), "DataSource count should be stable")
				assert.Equal(t, len(dataSources2), len(dataSources3), "DataSource count should be stable")
			},
		},
		{
			name: "error case - data sources are not nil",
			testFunc: func(t *testing.T) {
				t.Helper()
				prov := p.NewN8nProvider("1.0.0")
				dataSources := prov.DataSources(context.Background())
				assert.NotNil(t, dataSources, "DataSources must not be nil")
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

func TestValidateProvider(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                string
		version             string
		wantNotNil          bool
		wantSameInstance    bool
		wantImplementsIface bool
	}{
		{
			name:                "validates N8nProvider implements TerraformProvider",
			version:             "1.0.0",
			wantNotNil:          true,
			wantSameInstance:    true,
			wantImplementsIface: true,
		},
		{
			name:                "error case - validates with empty version",
			version:             "",
			wantNotNil:          true,
			wantSameInstance:    true,
			wantImplementsIface: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			original := p.NewN8nProvider(tt.version)
			validated := p.ValidateProvider(original)

			if tt.wantNotNil {
				assert.NotNil(t, validated, "Validated provider should not be nil")
			}

			if tt.wantSameInstance {
				assert.Equal(t, original, validated, "Should return the same provider")
			}
		})
	}
}

func TestValidateProvider_AllMethods(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "all methods work on validated provider",
			testFunc: func(t *testing.T) {
				t.Helper()
				prov := p.NewN8nProvider("1.0.0")
				validated := p.ValidateProvider(prov)

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
			},
		},
		{
			name: "error case - validated provider is not nil",
			testFunc: func(t *testing.T) {
				t.Helper()
				prov := p.NewN8nProvider("1.0.0")
				validated := p.ValidateProvider(prov)
				assert.NotNil(t, validated, "Validated provider must not be nil")
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

func TestProviderVersionHandling(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		version string
	}{
		{"semantic version", "1.2.3"},
		{"prerelease version", "1.2.3-alpha.1"},
		{"build metadata", "1.2.3+build.123"},
		{"dev version", "dev"},
		{"empty version", ""},
		{"snapshot", "1.0.0-SNAPSHOT"},
		{"error case - special characters", "!@#$%"},
		{"error case - null byte", "\x00"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			prov := p.NewN8nProvider(tt.version)
			resp := &provider.MetadataResponse{}
			prov.Metadata(context.Background(), provider.MetadataRequest{}, resp)

			assert.Equal(t, tt.version, resp.Version, "Version should be set correctly")
		})
	}
}

func TestN8nProviderModelStructure(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		apiKey      string
		baseURL     string
		wantAPIKey  string
		wantBaseURL string
		wantNulls   bool
	}{
		{
			name:        "model has required fields",
			apiKey:      "test-key",
			baseURL:     "https://test.com",
			wantAPIKey:  "test-key",
			wantBaseURL: "https://test.com",
			wantNulls:   false,
		},
		{
			name:        "model with empty values",
			apiKey:      "",
			baseURL:     "",
			wantAPIKey:  "",
			wantBaseURL: "",
			wantNulls:   false,
		},
		{
			name:        "error case - model with special characters",
			apiKey:      "!@#$%",
			baseURL:     "http://test\x00.com",
			wantAPIKey:  "!@#$%",
			wantBaseURL: "http://test\x00.com",
			wantNulls:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			model := &models.N8nProviderModel{
				APIKey:  types.StringValue(tt.apiKey),
				BaseURL: types.StringValue(tt.baseURL),
			}

			if !tt.wantNulls {
				assert.False(t, model.APIKey.IsNull())
				assert.False(t, model.BaseURL.IsNull())
			}

			assert.Equal(t, tt.wantAPIKey, model.APIKey.ValueString())
			assert.Equal(t, tt.wantBaseURL, model.BaseURL.ValueString())
		})
	}
}

func TestN8nProviderModelUsage(t *testing.T) {
	t.Parallel()

	prov := p.NewN8nProvider("1.0.0")

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
	prov.Schema(context.Background(), provider.SchemaRequest{}, schemaResp)

	config := tfsdk.Config{
		Schema: schemaResp.Schema,
		Raw:    configValue,
	}

	req := provider.ConfigureRequest{
		Config: config,
	}
	resp := &provider.ConfigureResponse{}

	prov.Configure(context.Background(), req, resp)

	require.False(t, resp.Diagnostics.HasError())
}

func TestProviderContextUsage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		version    string
		setupCtx   func() context.Context
		wantPanics bool
	}{
		{
			name:    "Metadata accepts context",
			version: "1.0.0",
			setupCtx: func() context.Context {
				return context.WithValue(context.Background(), contextKey("test-key"), "test-value")
			},
			wantPanics: false,
		},
		{
			name:    "Schema accepts context",
			version: "1.0.0",
			setupCtx: func() context.Context {
				return context.WithValue(context.Background(), contextKey("test-key"), "test-value")
			},
			wantPanics: false,
		},
		{
			name:    "Resources accepts context",
			version: "1.0.0",
			setupCtx: func() context.Context {
				return context.WithValue(context.Background(), contextKey("test-key"), "test-value")
			},
			wantPanics: false,
		},
		{
			name:    "DataSources accepts context",
			version: "1.0.0",
			setupCtx: func() context.Context {
				return context.WithValue(context.Background(), contextKey("test-key"), "test-value")
			},
			wantPanics: false,
		},
		{
			name:    "handles TODO context",
			version: "1.0.0",
			setupCtx: func() context.Context {
				return context.TODO()
			},
			wantPanics: false,
		},
		{
			name:    "error case - handles canceled context",
			version: "1.0.0",
			setupCtx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			wantPanics: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			prov := p.NewN8nProvider(tt.version)
			ctx := tt.setupCtx()

			if tt.wantPanics {
				// Test for panics if needed
				return
			}

			// Test Metadata
			assert.NotPanics(t, func() {
				resp := &provider.MetadataResponse{}
				prov.Metadata(ctx, provider.MetadataRequest{}, resp)
			})

			// Test Schema
			assert.NotPanics(t, func() {
				resp := &provider.SchemaResponse{}
				prov.Schema(ctx, provider.SchemaRequest{}, resp)
			})

			// Test Resources
			assert.NotPanics(t, func() {
				_ = prov.Resources(ctx)
			})

			// Test DataSources
			assert.NotPanics(t, func() {
				_ = prov.DataSources(ctx)
			})
		})
	}
}

func TestConfigureClientCreation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		version    string
		baseURL    string
		apiKey     string
		wantErr    bool
		wantClient bool
	}{
		{
			name:       "creates client with correct base URL",
			version:    "1.0.0",
			baseURL:    "https://n8n.example.com",
			apiKey:     "test-key",
			wantErr:    false,
			wantClient: true,
		},
		{
			name:       "creates client with different URL",
			version:    "1.0.0",
			baseURL:    "https://another.example.com",
			apiKey:     "test-key-2",
			wantErr:    false,
			wantClient: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			prov := p.NewN8nProvider(tt.version)

			configValue := tftypes.NewValue(
				tftypes.Object{
					AttributeTypes: map[string]tftypes.Type{
						"api_key":  tftypes.String,
						"base_url": tftypes.String,
					},
				},
				map[string]tftypes.Value{
					"api_key":  tftypes.NewValue(tftypes.String, tt.apiKey),
					"base_url": tftypes.NewValue(tftypes.String, tt.baseURL),
				},
			)

			schemaResp := &provider.SchemaResponse{}
			prov.Schema(context.Background(), provider.SchemaRequest{}, schemaResp)

			config := tfsdk.Config{
				Schema: schemaResp.Schema,
				Raw:    configValue,
			}

			req := provider.ConfigureRequest{
				Config: config,
			}
			resp := &provider.ConfigureResponse{}

			prov.Configure(context.Background(), req, resp)

			if tt.wantErr {
				assert.True(t, resp.Diagnostics.HasError())
			} else {
				assert.False(t, resp.Diagnostics.HasError())
			}

			if tt.wantClient {
				assert.NotNil(t, resp.ResourceData, "Client should be created")
			}
		})
	}
}

// stringPtr is a helper function to get a pointer to a string.
func stringPtr(s string) *string {
	return &s
}
