package main

import "github.com/kodflow/n8n/src/cmd"

var (
	// version represents the build version of the provider, injected at compile time.
	// This value is replaced at build time using -ldflags "-X main.version=x.y.z".
	version string = "dev"
)

// main is the entry point for the n8n Terraform provider.
// It sets the provider version and executes the root command.
func main() {
	cmd.SetVersion(version)
	cmd.Execute()
}
