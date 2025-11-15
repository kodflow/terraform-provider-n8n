// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE.md in the project root for license information.

// Package main provides the entry point for the n8n Terraform provider.
package main

import (
	"os"

	"github.com/kodflow/terraform-provider-n8n/src/cmd"
)

var (
	// version represents the build version of the provider, injected at compile time.
	// This value is replaced at build time using -ldflags "-X main.version=x.y.z".
	version string = "dev"
)

// main is the entry point for the n8n Terraform provider.
// It sets the provider version and executes the root command.
func main() {
	cmd.SetVersion(version)
	os.Exit(cmd.Execute())
}
