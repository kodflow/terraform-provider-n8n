package cmd

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	n8nprovider "github.com/kodflow/n8n/src/internal/provider"
	"github.com/spf13/cobra"
)

var (
	// Version is the build version, injected at compile time.
	Version string = "dev"

	// OsExit allows mocking os.Exit for testing.
	OsExit func(int) = os.Exit

	// ProviderServe allows mocking providerserver.Serve for testing.
	ProviderServe func(ctx context.Context, providerFunc func() provider.Provider, opts providerserver.ServeOpts) error = providerserver.Serve

	// debug enables support for debuggers like delve when set to true.
	debug bool

	// rootCmd represents the base command.
	rootCmd *cobra.Command = &cobra.Command{
		Use:   "terraform-provider-n8n",
		Short: "Terraform provider for n8n automation platform",
		Long:  `A Terraform provider that allows you to manage n8n resources.`,
		RunE:  run,
	}
)

func init() {
	rootCmd.Flags().BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
}

// run starts the Terraform provider server.
func run(cmd *cobra.Command, args []string) error {
	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/kodflow/n8n",
		Debug:   debug,
	}

	err := ProviderServe(context.Background(), n8nprovider.New(Version), opts)
	// Check for error.
	if err != nil {
  // Return error.
		return err
	}

	// Return result.
	return nil
}

// Execute runs the root command.
func Execute() {
	// Check for error.
	if err := rootCmd.Execute(); err != nil {
		OsExit(1)
	}
}

// SetVersion sets the version for the provider.
func SetVersion(v string) {
	Version = v
}
