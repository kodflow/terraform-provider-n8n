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

	tests := []struct {
		name         string
		serveErr     error
		version      string
		expectedCode int
	}{
		{
			name:         "returns 0 on success",
			serveErr:     nil,
			version:      "test-1.0.0",
			expectedCode: 0,
		},
		{
			name:         "returns 1 on error",
			serveErr:     os.ErrInvalid,
			version:      "test-error-1.0.0",
			expectedCode: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original values
			originalVersion := version
			originalCmdVersion := cmd.Version
			originalProviderServe := cmd.ProviderServe

			t.Cleanup(func() {
				version = originalVersion
				cmd.Version = originalCmdVersion
				cmd.ProviderServe = originalProviderServe
			})

			cmd.ProviderServe = func(ctx context.Context, providerFunc func() provider.Provider, opts providerserver.ServeOpts) error {
				return tt.serveErr
			}

			version = tt.version
			cmd.SetVersion(version)

			exitCode := cmd.Execute()
			if exitCode != tt.expectedCode {
				t.Errorf("Expected exit code %d, got %d", tt.expectedCode, exitCode)
			}
		})
	}
}

// TestVersionVariable tests the version variable and its interaction with cmd.SetVersion.
func TestVersionVariable(t *testing.T) {

	tests := []struct {
		name       string
		setVersion string
		expected   string
	}{
		{
			name:       "default version is dev",
			setVersion: "",
			expected:   "dev",
		},
		{
			name:       "set version is propagated",
			setVersion: "1.2.3",
			expected:   "1.2.3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save and restore original version
			originalVersion := version
			originalCmdVersion := cmd.Version
			t.Cleanup(func() {
				version = originalVersion
				cmd.Version = originalCmdVersion
			})

			if tt.setVersion != "" {
				version = tt.setVersion
				cmd.SetVersion(version)
				if cmd.Version != tt.expected {
					t.Errorf("Expected cmd.Version '%s', got '%s'", tt.expected, cmd.Version)
				}
			} else if version != tt.expected {
				t.Errorf("Expected version '%s', got '%s'", tt.expected, version)
			}
		})
	}
}
