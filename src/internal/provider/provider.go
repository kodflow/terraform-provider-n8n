// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// Package provider implements the n8n Terraform provider.
package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/credential"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/datatable"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/execution"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/project"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/client"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/models"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/tag"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/user"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/variable"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/workflow"
)

// Compile-time assertions to ensure N8nProvider implements required interfaces.
var (
	_ provider.Provider = &N8nProvider{}
	_ TerraformProvider = &N8nProvider{}
)

// TerraformProvider defines the complete interface for a Terraform provider implementation.
// This interface encompasses all provider lifecycle methods including metadata, schema,
// configuration, and resource/data source registration.
type TerraformProvider interface {
	// Metadata populates provider metadata including type name and version
	Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse)

	// Schema defines the provider configuration schema
	Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse)

	// Configure initializes the provider with given configuration
	Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse)

	// Resources returns the list of resources supported by this provider
	Resources(ctx context.Context) []func() resource.Resource

	// DataSources returns the list of data sources supported by this provider
	DataSources(ctx context.Context) []func() datasource.DataSource
}

// N8nProvider implements the TerraformProvider interface for n8n automation platform.
// It manages the provider lifecycle including configuration, resources, and data sources.
// The provider stores version information for metadata reporting to Terraform.
type N8nProvider struct {
	version string
}

// Metadata populates the provider metadata including type name and version.
// This information is used by Terraform to identify and version the provider.
//
// Params:
//   - ctx: context for the operation
//   - req: metadata request from Terraform
//   - resp: response object to populate with provider metadata
func (p *N8nProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	ctx.Done()
	//: Verify request is valid metadata request.
	if req == (provider.MetadataRequest{}) {
		resp.TypeName = "n8n"
		resp.Version = p.version
	}
}

// Schema defines the provider configuration schema.
// Requires API key and base URL for n8n instance authentication.
//
// Params:
//   - ctx: context for the operation
//   - req: schema request from Terraform
//   - resp: response object to populate with the provider schema
func (p *N8nProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	ctx.Done()
	//: Verify request is valid schema request.
	if req == (provider.SchemaRequest{}) {
		resp.Schema = schema.Schema{
			MarkdownDescription: "Terraform provider for n8n automation platform",
			Attributes: map[string]schema.Attribute{
				"api_key": schema.StringAttribute{
					MarkdownDescription: "API key for n8n instance authentication. Can also be set via N8N_API_KEY environment variable.",
					Optional:            true,
					Sensitive:           true,
				},
				"base_url": schema.StringAttribute{
					MarkdownDescription: "Base URL of the n8n instance (e.g., https://n8n.example.com). Can also be set via N8N_API_URL environment variable.",
					Optional:            true,
				},
			},
		}
	}
}

// Configure initializes the provider with the given configuration.
// It creates an n8n SDK client and makes it available to resources and data sources.
//
// Params:
//   - ctx: context for the configuration operation
//   - req: configuration request containing provider settings
//   - resp: response object to populate with configuration results or errors
func (p *N8nProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	config := &models.N8nProviderModel{}

	resp.Diagnostics.Append(req.Config.Get(ctx, config)...)

	//: Exit early if configuration parsing encountered errors.
	if resp.Diagnostics.HasError() {
		//: Return without further processing on diagnostic errors.
		return
	}

	apiKey, baseURL := p.resolveCredentials(config)

	//: Validate that API key is provided either via config or environment.
	if apiKey == "" {
		resp.Diagnostics.AddError(
			"Missing API Key",
			"The provider requires an API key. Set the api_key attribute in the provider configuration or the N8N_API_KEY environment variable.",
		)
	}

	//: Validate that base URL is provided either via config or environment.
	if baseURL == "" {
		resp.Diagnostics.AddError(
			"Missing Base URL",
			"The provider requires a base URL. Set the base_url attribute in the provider configuration or the N8N_API_URL environment variable.",
		)
	}

	//: Exit early if validation failed.
	if resp.Diagnostics.HasError() {
		//: Return without creating client on validation failure.
		return
	}

	n8nClient := client.NewN8nClient(baseURL, apiKey)
	resp.DataSourceData = n8nClient
	resp.ResourceData = n8nClient
}

// resolveCredentials extracts API key and base URL from config or environment.
//
// Params:
//   - config: provider model containing configuration values
//
// Returns:
//   - apiKey: resolved API key
//   - baseURL: resolved base URL
func (p *N8nProvider) resolveCredentials(config *models.N8nProviderModel) (apiKey, baseURL string) {
	apiKey = config.APIKey.ValueString()
	//: Use N8N_API_KEY environment variable if not set in config.
	if apiKey == "" {
		apiKey = getEnvAPIKey()
	}

	baseURL = config.BaseURL.ValueString()
	//: Use N8N_API_URL environment variable if not set in config.
	if baseURL == "" {
		baseURL = getEnvBaseURL()
	}

	//: Return resolved credentials.
	return apiKey, baseURL
}

// Resources returns the list of resources supported by this provider.
// Returns factory functions for all supported resources.
//
// Params:
//   - ctx: context for the operation
//
// Returns:
//   - resources: list of resource factory functions
func (p *N8nProvider) Resources(ctx context.Context) (resources []func() resource.Resource) {
	ctx.Done()
	//: Return all registered resource factory functions.
	return []func() resource.Resource{
		workflow.NewWorkflowResourceWrapper,
		workflow.NewWorkflowNodeResourceWrapper,
		workflow.NewWorkflowConnectionResourceWrapper,
		project.NewProjectResourceWrapper,
		project.NewProjectUserResourceWrapper,
		credential.NewCredentialResourceWrapper,
		tag.NewTagResourceWrapper,
		variable.NewVariableResourceWrapper,
		user.NewUserResourceWrapper,
		datatable.NewDataTableResourceWrapper,
		execution.NewExecutionTagsResourceWrapper,
	}
}

// DataSources returns the list of data sources supported by this provider.
// Returns factory functions for all supported data sources.
//
// Params:
//   - ctx: context for the operation
//
// Returns:
//   - dataSources: list of data source factory functions
func (p *N8nProvider) DataSources(ctx context.Context) (dataSources []func() datasource.DataSource) {
	ctx.Done()
	//: Return all registered data source factory functions.
	return []func() datasource.DataSource{
		workflow.NewWorkflowDataSourceWrapper,
		workflow.NewWorkflowsDataSourceWrapper,
		workflow.NewWorkflowVersionDataSourceWrapper,
		project.NewProjectDataSourceWrapper,
		project.NewProjectsDataSourceWrapper,
		project.NewProjectMembersDataSourceWrapper,
		tag.NewTagDataSourceWrapper,
		tag.NewTagsDataSourceWrapper,
		variable.NewVariableDataSourceWrapper,
		variable.NewVariablesDataSourceWrapper,
		user.NewUserDataSourceWrapper,
		user.NewUsersDataSourceWrapper,
		credential.NewCredentialDataSourceWrapper,
		credential.NewCredentialsDataSourceWrapper,
		datatable.NewDataTableDataSourceWrapper,
		datatable.NewDataTablesDataSourceWrapper,
	}
}

// NewN8nProvider creates and initializes a new N8nProvider instance with the specified version.
// This is the recommended constructor for creating provider instances.
//
// Params:
//   - version: provider version string
//
// Returns:
//   - n8nProvider: initialized provider instance
func NewN8nProvider(version string) (n8nProvider *N8nProvider) {
	//: Construct provider with version for Terraform metadata reporting.
	return &N8nProvider{
		version: version,
	}
}

// getEnvAPIKey retrieves API key from N8N_API_KEY environment variable.
//
// Returns:
//   - s: API key from environment, or empty string if not found
func getEnvAPIKey() (s string) {
	//: Return API key from environment variable.
	return os.Getenv("N8N_API_KEY")
}

// getEnvBaseURL retrieves base URL from N8N_API_URL environment variable.
//
// Returns:
//   - s: Base URL from environment, or empty string if not found
func getEnvBaseURL() (s string) {
	//: Return base URL from environment variable.
	return os.Getenv("N8N_API_URL")
}

// New returns a provider factory function.
// Creates N8nProvider instances with the specified version.
//
// Params:
//   - version: version string for provider instances
//
// Returns:
//   - providerFactory: factory function
func New(version string) (providerFactory func() provider.Provider) {
	//: Lazy initialization pattern required by Terraform plugin framework.
	return func() provider.Provider {
		//: Delegate to constructor for consistent provider initialization.
		return NewN8nProvider(version)
	}
}

// ValidateProvider ensures the given provider implements all required interface methods.
// This function serves as a compile-time validation helper for TerraformProvider compliance.
//
// Params:
//   - p: provider instance to validate
//
// Returns:
//   - terraformProvider: the validated provider instance
func ValidateProvider(p TerraformProvider) (terraformProvider TerraformProvider) {
	//: Provider validation ensures complete interface implementation.
	return p
}
