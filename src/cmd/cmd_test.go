package cmd_test

import (
	"context"
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/kodflow/n8n/src/cmd"
)

func TestSetVersion(t *testing.T) {
	tests := []struct {
		name    string
		version string
	}{
		{
			name:    "sets version to 1.0.0",
			version: "1.0.0",
		},
		{
			name:    "sets version to dev",
			version: "dev",
		},
		{
			name:    "sets version to empty string",
			version: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd.SetVersion(tt.version)
			if cmd.Version != tt.version {
				t.Errorf("Expected version %q, got %q", tt.version, cmd.Version)
			}
		})
	}
}

func TestExecute_Success(t *testing.T) {
	// Save original functions
	originalExit := cmd.OsExit
	originalServe := cmd.ProviderServe
	defer func() {
		cmd.OsExit = originalExit
		cmd.ProviderServe = originalServe
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

	// Execute should not call OsExit on success
	cmd.Execute()

	// Verify OsExit was not called (exitCode remains -1)
	if exitCode != -1 {
		t.Errorf("Expected no exit call, but OsExit was called with code %d", exitCode)
	}
}

func TestExecute_WithError(t *testing.T) {
	// Save original functions
	originalExit := cmd.OsExit
	originalServe := cmd.ProviderServe
	defer func() {
		cmd.OsExit = originalExit
		cmd.ProviderServe = originalServe
	}()

	// Mock OsExit to capture exit calls
	exitCode := -1
	cmd.OsExit = func(code int) {
		exitCode = code
	}

	// Mock ProviderServe to return an error
	cmd.ProviderServe = func(ctx context.Context, providerFunc func() provider.Provider, opts providerserver.ServeOpts) error {
		return errors.New("mock error")
	}

	// Execute should call OsExit(1) on error
	cmd.Execute()

	// Verify OsExit was called with code 1
	if exitCode != 1 {
		t.Errorf("Expected exit code 1, got %d", exitCode)
	}
}

func TestRun_Success(t *testing.T) {
	// Save original ProviderServe
	originalServe := cmd.ProviderServe
	defer func() {
		cmd.ProviderServe = originalServe
	}()

	// Mock ProviderServe to return success
	called := false
	cmd.ProviderServe = func(ctx context.Context, providerFunc func() provider.Provider, opts providerserver.ServeOpts) error {
		called = true
		// Verify options
		if opts.Address != "registry.terraform.io/kodflow/n8n" {
			t.Errorf("Expected address 'registry.terraform.io/kodflow/n8n', got %q", opts.Address)
		}
		return nil
	}

	// Call run - we need to access the root command and invoke it
	// Since run is not exported, we test it through Execute
	cmd.Execute()

	// Verify ProviderServe was called
	if !called {
		t.Error("Expected ProviderServe to be called")
	}
}

func TestRun_WithError(t *testing.T) {
	// Save original functions
	originalExit := cmd.OsExit
	originalServe := cmd.ProviderServe
	defer func() {
		cmd.OsExit = originalExit
		cmd.ProviderServe = originalServe
	}()

	// Mock OsExit
	exitCode := -1
	cmd.OsExit = func(code int) {
		exitCode = code
	}

	// Mock ProviderServe to return an error
	expectedErr := errors.New("mock provider error")
	cmd.ProviderServe = func(ctx context.Context, providerFunc func() provider.Provider, opts providerserver.ServeOpts) error {
		return expectedErr
	}

	// Execute should call OsExit(1) on error
	cmd.Execute()

	// Verify error caused exit
	if exitCode != 1 {
		t.Errorf("Expected exit code 1, got %d", exitCode)
	}
}
