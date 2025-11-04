package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/kodflow/n8n/src/internal/provider"
)

var (
	// version represents the build version of the provider, injected at compile time.
	version = "dev"
)

// main initializes and starts the Terraform provider server with the configured options.
// It supports debug mode for development and handles provider registration with Terraform.
func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/kodflow/n8n",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), provider.New(version), opts)

	// Handle fatal errors during provider server startup
	if err != nil {
		log.Fatal(err.Error())
	}
}
