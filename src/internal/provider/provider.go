package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Compile-time assertion to ensure N8nProvider implements the provider.Provider interface.
var _ provider.Provider = &N8nProvider{}

// N8nProvider implements the Terraform provider for n8n automation platform.
// It manages the provider lifecycle including configuration, resources, and data sources.
// The provider stores version information for metadata reporting to Terraform.
type N8nProvider struct {
	version string
}

// N8nProviderModel represents the configuration schema for the n8n Terraform provider.
// Currently, the provider requires no configuration parameters but this struct
// serves as a placeholder for future configuration options.
type N8nProviderModel struct{}

// Interface defines the contract for the n8n Terraform provider.
// It encompasses all required provider operations including metadata, schema,
// configuration, resource management, and data source management.
type Interface interface {
	provider.Provider
}

// Metadata populates the provider metadata including type name and version.
// This information is used by Terraform to identify and version the provider.
//
// Params:
//   - ctx: The context for the operation (currently unused).
//   - req: The metadata request from Terraform (currently unused).
//   - resp: The response object to populate with provider metadata.
func (p *N8nProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "n8n"
	resp.Version = p.version
}

// Schema defines the provider configuration schema.
// Currently returns an empty schema as no provider-level configuration is required.
//
// Params:
//   - ctx: The context for the operation (currently unused).
//   - req: The schema request from Terraform (currently unused).
//   - resp: The response object to populate with the provider schema.
func (p *N8nProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Terraform provider for n8n automation platform",
	}
}

// Configure initializes the provider with the given configuration.
// It parses the provider configuration and handles any configuration errors.
//
// Params:
//   - ctx: The context for the configuration operation.
//   - req: The configuration request containing provider settings.
//   - resp: The response object to populate with configuration results or errors.
func (p *N8nProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config N8nProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	// Exit early if configuration parsing encountered errors
	if resp.Diagnostics.HasError() {
		// Stop processing due to configuration validation errors
		return
	}
}

// Resources returns the list of resources supported by this provider.
// Currently returns an empty list as resources are not yet implemented.
//
// Params:
//   - ctx: The context for the operation (currently unused).
func (p *N8nProvider) Resources(_ context.Context) []func() resource.Resource {
	// Return nil to avoid unnecessary allocation for empty slice
	return nil
}

// DataSources returns the list of data sources supported by this provider.
// Currently returns an empty list as data sources are not yet implemented.
//
// Params:
//   - ctx: The context for the operation (currently unused).
func (p *N8nProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	// Return nil to avoid unnecessary allocation for empty slice
	return nil
}

// NewN8nProvider creates and initializes a new N8nProvider instance with the specified version.
// This is the recommended constructor for creating provider instances.
//
// Params:
//   - version: The version string for the provider.
func NewN8nProvider(version string) *N8nProvider {
	// Create and return a new provider instance with the given version
	return &N8nProvider{
		version: version,
	}
}

// New returns a provider factory function that creates N8nProvider instances.
// This function is required by the Terraform plugin framework for provider initialization.
//
// Params:
//   - version: The version string to assign to created provider instances.
func New(version string) func() provider.Provider {
	// Return a factory function that creates new provider instances
	return func() provider.Provider {
		// Create and return a new provider instance with the specified version
		return NewN8nProvider(version)
	}
}
