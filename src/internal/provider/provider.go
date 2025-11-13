// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE.md in the project root for license information.

// Package provider implements the n8n Terraform provider.
package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/kodflow/n8n/src/internal/provider/credential"
	"github.com/kodflow/n8n/src/internal/provider/execution"
	"github.com/kodflow/n8n/src/internal/provider/project"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
	"github.com/kodflow/n8n/src/internal/provider/shared/models"
	"github.com/kodflow/n8n/src/internal/provider/sourcecontrol"
	"github.com/kodflow/n8n/src/internal/provider/tag"
	"github.com/kodflow/n8n/src/internal/provider/user"
	"github.com/kodflow/n8n/src/internal/provider/variable"
	"github.com/kodflow/n8n/src/internal/provider/workflow"
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
func (p *N8nProvider) Metadata(_ctx context.Context, _req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "n8n"
	resp.Version = p.version
}

// Schema defines the provider configuration schema.
// Requires API key and base URL for n8n instance authentication.
//
// Params:
//   - ctx: context for the operation
//   - req: schema request from Terraform
//   - resp: response object to populate with the provider schema
func (p *N8nProvider) Schema(_ctx context.Context, _req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Terraform provider for n8n automation platform",
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				MarkdownDescription: "API key for n8n instance authentication. Can also be set via N8N_API_TOKEN environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
			"base_url": schema.StringAttribute{
				MarkdownDescription: "Base URL of the n8n instance (e.g., https://n8n.example.com). Can also be set via N8N_URL environment variable.",
				Optional:            true,
			},
		},
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

	// Exit early if configuration parsing encountered errors
	if resp.Diagnostics.HasError() {
		return
	}

	// Read from environment variables if not set in configuration
	apiKey := config.APIKey.ValueString()
	// Use environment variable N8N_API_TOKEN if api_key is not provided
	if apiKey == "" {
		// Read API token from environment
		if envAPIKey := os.Getenv("N8N_API_TOKEN"); envAPIKey != "" {
			apiKey = envAPIKey
		}
	}

	baseURL := config.BaseURL.ValueString()
	// Use environment variable N8N_URL if base_url is not provided
	if baseURL == "" {
		// Read base URL from environment
		if envBaseURL := os.Getenv("N8N_URL"); envBaseURL != "" {
			baseURL = envBaseURL
		}
	}

	// Validate required configuration
	// Check that API key is provided either via config or environment
	if apiKey == "" {
		resp.Diagnostics.AddError(
			"Missing API Key",
			"The provider requires an API key. Set the api_key attribute in the provider configuration or the N8N_API_TOKEN environment variable.",
		)
	}

	// Check that base URL is provided either via config or environment
	if baseURL == "" {
		resp.Diagnostics.AddError(
			"Missing Base URL",
			"The provider requires a base URL. Set the base_url attribute in the provider configuration or the N8N_URL environment variable.",
		)
	}

	// Exit early if validation failed
	if resp.Diagnostics.HasError() {
		return
	}

	// Create n8n client using the generated SDK
	n8nClient := client.NewN8nClient(baseURL, apiKey)

	// Make client available to resources and data sources
	resp.DataSourceData = n8nClient
	resp.ResourceData = n8nClient
}

// Resources returns the list of resources supported by this provider.
// Returns factory functions for all supported resources.
//
// Params:
//   - ctx: context for the operation
//
// Returns:
//   - []func() resource.Resource: list of resource factory functions
func (p *N8nProvider) Resources(_ctx context.Context) []func() resource.Resource {
	// Return result.
	return []func() resource.Resource{
		// Workflow domain
		workflow.NewWorkflowResourceWrapper,
		workflow.NewWorkflowTransferResourceWrapper,
		// Project domain
		project.NewProjectResourceWrapper,
		project.NewProjectUserResourceWrapper,
		// Credential domain
		credential.NewCredentialResourceWrapper,
		credential.NewCredentialTransferResourceWrapper,
		// Execution domain
		execution.NewExecutionRetryResourceWrapper,
		// Tag domain
		tag.NewTagResourceWrapper,
		// Variable domain
		variable.NewVariableResourceWrapper,
		// User domain
		user.NewUserResourceWrapper,
		// Source control domain
		sourcecontrol.NewSourceControlPullResourceWrapper,
	}
}

// DataSources returns the list of data sources supported by this provider.
// Returns factory functions for all supported data sources.
//
// Params:
//   - ctx: context for the operation
//
// Returns:
//   - []func() datasource.DataSource: list of data source factory functions
func (p *N8nProvider) DataSources(_ctx context.Context) []func() datasource.DataSource {
	// Return result.
	return []func() datasource.DataSource{
		// Workflow domain
		workflow.NewWorkflowDataSourceWrapper,
		workflow.NewWorkflowsDataSourceWrapper,
		// Project domain
		project.NewProjectDataSourceWrapper,
		project.NewProjectsDataSourceWrapper,
		// Execution domain
		execution.NewExecutionDataSourceWrapper,
		execution.NewExecutionsDataSourceWrapper,
		// Tag domain
		tag.NewTagDataSourceWrapper,
		tag.NewTagsDataSourceWrapper,
		// Variable domain
		variable.NewVariableDataSourceWrapper,
		variable.NewVariablesDataSourceWrapper,
		// User domain
		user.NewUserDataSourceWrapper,
		user.NewUsersDataSourceWrapper,
	}
}

// NewN8nProvider creates and initializes a new N8nProvider instance with the specified version.
// This is the recommended constructor for creating provider instances.
//
// Params:
//   - version: provider version string
//
// Returns:
//   - *N8nProvider: initialized provider instance
func NewN8nProvider(version string) *N8nProvider {
	// Construct provider with version for Terraform metadata reporting
	return &N8nProvider{
		version: version,
	}
}

// New returns a provider factory function that creates N8nProvider instances.
// This function is required by the Terraform plugin framework for provider initialization.
//
// Params:
//   - version: version string to assign to created provider instances
//
// Returns:
//   - func() provider.Provider: factory function that creates provider instances
func New(version string) func() provider.Provider {
	// Lazy initialization pattern required by Terraform plugin framework
	return func() provider.Provider {
		// Delegate to constructor for consistent provider initialization
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
//   - TerraformProvider: the validated provider instance
func ValidateProvider(p TerraformProvider) TerraformProvider {
	// Provider validation ensures complete interface implementation
	return p
}
