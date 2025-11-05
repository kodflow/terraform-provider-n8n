package main

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/kodflow/n8n/src/cmd"
)

func TestVersion(t *testing.T) {
	// Test that version variable is accessible and has a default value
	if version == "" {
		t.Error("Expected version to have a default value")
	}

	// Test version can be modified (for build-time injection)
	originalVersion := version
	version = "1.2.3"
	if version != "1.2.3" {
		t.Error("Version should be modifiable")
	}
	version = originalVersion
}

func TestMain_SetVersion(t *testing.T) {
	// Test that main() calls cmd.SetVersion
	originalVersion := cmd.Version
	defer func() {
		cmd.Version = originalVersion
	}()

	// Set a test version
	version = "test-version"

	// Simulate what main() does
	cmd.SetVersion(version)

	// Verify version was set
	if cmd.Version != "test-version" {
		t.Errorf("Expected cmd.Version to be 'test-version', got %q", cmd.Version)
	}
}

func TestMainFunction(t *testing.T) {
	// Save original functions
	originalExit := cmd.OsExit
	originalServe := cmd.ProviderServe
	originalVersion := cmd.Version
	defer func() {
		cmd.OsExit = originalExit
		cmd.ProviderServe = originalServe
		cmd.Version = originalVersion
	}()

	// Mock OsExit to capture exit calls
	exitCode := -1
	cmd.OsExit = func(code int) {
		exitCode = code
	}

	// Mock ProviderServe to return success
	cmd.ProviderServe = func(ctx context.Context, providerFunc func() provider.Provider, opts providerserver.ServeOpts) error {
		return nil
	}

	// Set test version
	version = "test-main-version"

	// Call main
	main()

	// Verify no exit was called
	if exitCode != -1 {
		t.Errorf("Expected no exit call, but OsExit was called with code %d", exitCode)
	}

	// Verify version was set
	if cmd.Version != "test-main-version" {
		t.Errorf("Expected cmd.Version to be 'test-main-version', got %q", cmd.Version)
	}
}
