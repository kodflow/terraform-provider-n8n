// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// Package cmd provides the CLI entry point for the n8n Terraform provider.
package cmd

import (
	"context"
	"errors"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	n8nprovider "github.com/kodflow/terraform-provider-n8n/src/internal/provider"
	"github.com/spf13/cobra"
)

var (
	// ErrUnexpectedArguments is returned when the root command receives unexpected arguments.
	ErrUnexpectedArguments = errors.New("this command takes no arguments")

	// Version is the build version, injected at compile time.
	Version string = "dev"

	// OsExit allows mocking os.Exit for testing.
	OsExit func(int) = os.Exit

	// ProviderServe allows mocking providerserver.Serve for testing.
	ProviderServe func(ctx context.Context, providerFunc func() provider.Provider, opts providerserver.ServeOpts) error = providerserver.Serve

	// debug enables support for debuggers like delve when set to true.
	debug bool

	// rootCmd represents the base command.
	rootCmd *cobra.Command = newRootCmd()
)

// newRootCmd creates and configures the root cobra command with flags.
//
// Returns:
//   - cmd: configured root command
func newRootCmd() (cmd *cobra.Command) {
	//: Build the root command with provider configuration.
	c := &cobra.Command{
		Use:   "terraform-provider-n8n",
		Short: "Terraform provider for n8n automation platform",
		Long:  `A Terraform provider that allows you to manage n8n resources.`,
		RunE:  run,
	}

	//: Register debug flag for debugger support.
	c.Flags().BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")

	//: Return configured command.
	return c
}

// run starts the Terraform provider server.
//
// Params:
//   - cmd: cobra command providing lifecycle context
//   - args: command arguments (must be empty for this command)
//
// Returns:
//   - err: Error if provider server fails to start
func run(cmd *cobra.Command, args []string) (err error) {
	//: Guard against unexpected arguments.
	if len(args) > 0 {
		//: Return sentinel error for unexpected arguments.
		return ErrUnexpectedArguments
	}

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/kodflow/n8n",
		Debug:   debug,
	}

	//: Fall back to background context when cmd is unavailable.
	ctx := context.Background()
	//: Prefer cmd context for proper signal/cancellation propagation.
	if cmd != nil && cmd.Context() != nil {
		ctx = cmd.Context()
	}
	err = ProviderServe(ctx, n8nprovider.New(Version), opts)
	//: Check for provider serve error.
	if err != nil {
		//: Return error from provider server.
		return err
	}

	//: Return nil on successful provider serve.
	return nil
}

// Execute runs the root command and returns the exit code.
//
// Returns:
//   - n: Exit code (0 for success, 1 for error)
func Execute() (n int) {
	//: Check for error on root command execution.
	if err := rootCmd.Execute(); err != nil {
		//: Return error code on failure.
		return 1
	}
	//: Return success code.
	return 0
}

// SetVersion sets the version for the provider.
//
// Params:
//   - v: Version string to set
func SetVersion(v string) {
	Version = v
}
