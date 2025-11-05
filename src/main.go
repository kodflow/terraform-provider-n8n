package main

import "github.com/kodflow/n8n/src/cmd"

var (
	// version represents the build version of the provider, injected at compile time.
	// This value is replaced at build time using -ldflags "-X main.version=x.y.z".
	version string = "dev"
)

func main() {
	cmd.SetVersion(version)
	cmd.Execute()
}
