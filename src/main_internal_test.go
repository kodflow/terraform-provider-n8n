package main

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/kodflow/terraform-provider-n8n/src/cmd"
)

// TestMainFunction tests the main function's behavior without calling it directly.
// We can't call main() directly because it calls os.Exit(), so we test its components.
func TestMainFunction(t *testing.T) {
	// Save original values
	originalVersion := version
	originalCmdVersion := cmd.Version
	defer func() {
		version = originalVersion
		cmd.Version = originalCmdVersion
	}()

	testCases := []struct {
		name            string
		inputVersion    string
		expectedVersion string
		setupFunc       func()
		assertFunc      func(t *testing.T)
	}{
		{
			name:            "main sets version correctly",
			inputVersion:    "main-test-1.0.0",
			expectedVersion: "main-test-1.0.0",
			setupFunc: func() {
				version = "main-test-1.0.0"
				cmd.SetVersion(version)
			},
			assertFunc: func(t *testing.T) {
				t.Helper()
				if cmd.Version != "main-test-1.0.0" {
					t.Errorf("Expected cmd.Version to be %q, got %q", "main-test-1.0.0", cmd.Version)
				}
			},
		},
		{
			name:            "main function flow with SetVersion",
			inputVersion:    "flow-test-2.0.0",
			expectedVersion: "flow-test-2.0.0",
			setupFunc: func() {
				version = "flow-test-2.0.0"
				cmd.SetVersion(version)
			},
			assertFunc: func(t *testing.T) {
				t.Helper()
				if cmd.Version != "flow-test-2.0.0" {
					t.Errorf("SetVersion should set cmd.Version to %q, got %q", "flow-test-2.0.0", cmd.Version)
				}
			},
		},
		{
			name:            "version variable has correct default",
			inputVersion:    "dev",
			expectedVersion: "dev",
			setupFunc: func() {
				version = "dev"
			},
			assertFunc: func(t *testing.T) {
				t.Helper()
				if version != "dev" {
					t.Errorf("Default version should be 'dev', got %q", version)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			if tc.setupFunc != nil {
				tc.setupFunc()
			}

			// Assert
			if tc.assertFunc != nil {
				tc.assertFunc(t)
			}
		})
	}
}

// TestExecuteReturnsExitCode tests that Execute returns proper exit codes.
func TestExecuteReturnsExitCode(t *testing.T) {
	// Save original values
	originalVersion := version
	originalCmdVersion := cmd.Version
	originalProviderServe := cmd.ProviderServe

	// Restore after test
	defer func() {
		version = originalVersion
		cmd.Version = originalCmdVersion
		cmd.ProviderServe = originalProviderServe
	}()

	t.Run("returns 0 on success", func(t *testing.T) {
		// Mock ProviderServe to succeed
		cmd.ProviderServe = func(ctx context.Context, providerFunc func() provider.Provider, opts providerserver.ServeOpts) error {
			return nil
		}

		// Set version
		version = "test-1.0.0"
		cmd.SetVersion(version)

		// Execute should return 0 (success)
		exitCode := cmd.Execute()
		if exitCode != 0 {
			t.Errorf("Expected exit code 0, got %d", exitCode)
		}
	})

	t.Run("returns 1 on error", func(t *testing.T) {
		// Mock ProviderServe to fail
		cmd.ProviderServe = func(ctx context.Context, providerFunc func() provider.Provider, opts providerserver.ServeOpts) error {
			return os.ErrInvalid
		}

		// Set version
		version = "test-error-1.0.0"
		cmd.SetVersion(version)

		// Execute should return 1 (error)
		exitCode := cmd.Execute()
		if exitCode != 1 {
			t.Errorf("Expected exit code 1, got %d", exitCode)
		}
	})
}

// TestVersionVariable tests the version variable.
func TestVersionVariable(t *testing.T) {
	t.Run("default version is dev", func(t *testing.T) {
		if version != "dev" {
			t.Errorf("Expected default version 'dev', got '%s'", version)
		}
	})
}
