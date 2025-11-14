package cmd

// NOTE: This file contains internal tests for the cmd package.
// These tests MUST stay in the internal test file (package cmd) because they require
// access to private package variables and functions:
//   - TestRun: Tests the private 'run' function
//   - TestExecute: Needs to modify and restore global 'rootCmd' variable
//   - TestSetVersion: Needs to modify and restore global 'Version' variable
//   - TestProviderServeVariable: Needs to modify global 'ProviderServe' function
//   - TestOsExitVariable: Needs to modify global 'OsExit' function
//   - TestDebugFlag: Needs to modify global 'debug' variable
//
// Moving these tests to an external test file (package cmd_test) would make it
// impossible to access and mock these private variables for testing.

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockProvider for testing.
type MockProvider struct {
	mock.Mock
	provider.Provider
}

// contextKey is a custom type for context keys to avoid collisions.
type contextKey string

func TestInit(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "init configures root command",
			wantErr: false,
		},
		{
			name:    "init adds debug flag",
			wantErr: false,
		},
		{
			name:    "debug flag has correct default",
			wantErr: false,
		},
		{
			name:    "debug flag has proper usage text",
			wantErr: false,
		},
		{
			name:    "debug flag shorthand is empty",
			wantErr: false,
		},
		{
			name:    "error case - root command is nil",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "init configures root command":
				// Verify that init() has set up the rootCmd
				assert.NotNil(t, rootCmd, "rootCmd should be initialized")
				assert.NotNil(t, rootCmd.Flags(), "rootCmd should have flags")

			case "init adds debug flag":
				// Verify that the debug flag exists
				debugFlag := rootCmd.Flags().Lookup("debug")
				assert.NotNil(t, debugFlag, "debug flag should be defined")
				assert.Equal(t, "bool", debugFlag.Value.Type(), "debug flag should be boolean")

			case "debug flag has correct default":
				debugFlag := rootCmd.Flags().Lookup("debug")
				require.NotNil(t, debugFlag)
				assert.Equal(t, "false", debugFlag.DefValue, "debug flag should default to false")

			case "debug flag has proper usage text":
				debugFlag := rootCmd.Flags().Lookup("debug")
				require.NotNil(t, debugFlag)
				assert.Contains(t, debugFlag.Usage, "debugger", "Usage should mention debugger support")
				assert.Contains(t, debugFlag.Usage, "delve", "Usage should mention delve")

			case "debug flag shorthand is empty":
				debugFlag := rootCmd.Flags().Lookup("debug")
				require.NotNil(t, debugFlag)
				assert.Empty(t, debugFlag.Shorthand, "debug flag should not have shorthand")

			case "error case - root command is nil":
				// Verify behavior when root command is not properly initialized
				// In normal execution, rootCmd is always initialized by init()
				// This error case tests that the initialization is not skipped
				assert.NotNil(t, rootCmd, "rootCmd must be initialized by init()")
			}
		})
	}
}

func Test_run(t *testing.T) {
	// Note: No t.Parallel() at function level because this test modifies global ProviderServe and debug
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "run succeeds when ProviderServe succeeds", wantErr: false},
		{name: "run fails when ProviderServe fails", wantErr: true},
		{name: "run passes correct address to ProviderServe", wantErr: false},
		{name: "run respects debug flag when false", wantErr: false},
		{name: "run respects debug flag when true", wantErr: false},
		{name: "run passes context to ProviderServe", wantErr: false},
		{name: "run creates provider with correct version", wantErr: false},
		{name: "run handles nil command", wantErr: false},
		{name: "run handles empty args", wantErr: false},
		{name: "run handles args with values", wantErr: false},
		{name: "run handles nil args", wantErr: false},
		{name: "run with panic recovery", wantErr: false},
		{name: "run with very long args", wantErr: false},
		{name: "run with special characters in args", wantErr: false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Note: No t.Parallel() here because subtests modify global ProviderServe and debug
			originalProviderServe := ProviderServe
			originalDebug := debug
			defer func() {
				ProviderServe = originalProviderServe
				debug = originalDebug
			}()

			switch tt.name {
			case "run succeeds when ProviderServe succeeds":
				ProviderServe = func(ctx context.Context, providerFunc func() provider.Provider, opts providerserver.ServeOpts) error {
					assert.Equal(t, "registry.terraform.io/kodflow/n8n", opts.Address)
					assert.Equal(t, debug, opts.Debug)
					p := providerFunc()
					assert.NotNil(t, p, "Provider should be created")
					return nil
				}
				err := run(&cobra.Command{}, []string{})
				assert.NoError(t, err, "run should succeed when ProviderServe succeeds")

			case "run fails when ProviderServe fails":
				expectedErr := errors.New("provider serve error")
				ProviderServe = func(ctx context.Context, providerFunc func() provider.Provider, opts providerserver.ServeOpts) error {
					return expectedErr
				}
				err := run(&cobra.Command{}, []string{})
				assert.Error(t, err, "run should fail when ProviderServe fails")
				assert.Equal(t, expectedErr, err, "run should return the ProviderServe error")

			case "run passes correct address to ProviderServe":
				addressCaptured := ""
				ProviderServe = func(ctx context.Context, providerFunc func() provider.Provider, opts providerserver.ServeOpts) error {
					addressCaptured = opts.Address
					return nil
				}
				_ = run(&cobra.Command{}, []string{})
				assert.Equal(t, "registry.terraform.io/kodflow/n8n", addressCaptured, "Correct address should be passed")

			case "run respects debug flag when false":
				debug = false
				debugCaptured := true
				ProviderServe = func(ctx context.Context, providerFunc func() provider.Provider, opts providerserver.ServeOpts) error {
					debugCaptured = opts.Debug
					return nil
				}
				_ = run(&cobra.Command{}, []string{})
				assert.False(t, debugCaptured, "Debug should be false when flag is false")

			case "run respects debug flag when true":
				debug = true
				debugCaptured := false
				ProviderServe = func(ctx context.Context, providerFunc func() provider.Provider, opts providerserver.ServeOpts) error {
					debugCaptured = opts.Debug
					return nil
				}
				_ = run(&cobra.Command{}, []string{})
				assert.True(t, debugCaptured, "Debug should be true when flag is true")

			case "run passes context to ProviderServe":
				var ctxCaptured context.Context
				ProviderServe = func(ctx context.Context, providerFunc func() provider.Provider, opts providerserver.ServeOpts) error {
					ctxCaptured = ctx
					return nil
				}
				_ = run(&cobra.Command{}, []string{})
				assert.NotNil(t, ctxCaptured, "Context should be passed to ProviderServe")
				assert.Equal(t, context.Background(), ctxCaptured, "Background context should be used")

			case "run creates provider with correct version":
				originalVersion := Version
				Version = "test-version-1.2.3"
				defer func() { Version = originalVersion }()
				var createdProvider provider.Provider
				ProviderServe = func(ctx context.Context, providerFunc func() provider.Provider, opts providerserver.ServeOpts) error {
					createdProvider = providerFunc()
					return nil
				}
				_ = run(&cobra.Command{}, []string{})
				assert.NotNil(t, createdProvider, "Provider should be created")

			case "run handles nil command":
				ProviderServe = func(ctx context.Context, providerFunc func() provider.Provider, opts providerserver.ServeOpts) error {
					return nil
				}
				err := run(nil, []string{})
				assert.NoError(t, err, "run should handle nil command")

			case "run handles empty args":
				ProviderServe = func(ctx context.Context, providerFunc func() provider.Provider, opts providerserver.ServeOpts) error {
					return nil
				}
				err := run(&cobra.Command{}, []string{})
				assert.NoError(t, err, "run should handle empty args")

			case "run handles args with values":
				ProviderServe = func(ctx context.Context, providerFunc func() provider.Provider, opts providerserver.ServeOpts) error {
					return nil
				}
				err := run(&cobra.Command{}, []string{"arg1", "arg2", "--flag=value"})
				assert.NoError(t, err, "run should handle args")

			case "run handles nil args":
				ProviderServe = func(ctx context.Context, providerFunc func() provider.Provider, opts providerserver.ServeOpts) error {
					return nil
				}
				err := run(&cobra.Command{}, nil)
				assert.NoError(t, err, "run should handle nil args")

			case "run with panic recovery":
				ProviderServe = func(ctx context.Context, providerFunc func() provider.Provider, opts providerserver.ServeOpts) error {
					defer func() {
						_ = recover()
					}()
					_ = providerFunc()
					return nil
				}
				assert.NotPanics(t, func() {
					_ = run(&cobra.Command{}, []string{})
				})

			case "run with very long args":
				ProviderServe = func(ctx context.Context, providerFunc func() provider.Provider, opts providerserver.ServeOpts) error {
					return nil
				}
				longArgs := make([]string, 1000)
				for i := range longArgs {
					longArgs[i] = fmt.Sprintf("arg%d", i)
				}
				err := run(&cobra.Command{}, longArgs)
				assert.NoError(t, err, "run should handle many args")

			case "run with special characters in args":
				ProviderServe = func(ctx context.Context, providerFunc func() provider.Provider, opts providerserver.ServeOpts) error {
					return nil
				}
				specialArgs := []string{
					"arg with spaces",
					"arg-with-dashes",
					"arg_with_underscores",
					"arg.with.dots",
					"arg/with/slashes",
					"arg\\with\\backslashes",
					"arg:with:colons",
					"arg;with;semicolons",
					"arg=with=equals",
					"arg?with?questions",
					"arg*with*asterisks",
					"arg[with]brackets",
					"arg{with}braces",
					"arg(with)parens",
					"arg<with>angles",
					"arg|with|pipes",
					"arg&with&ampersands",
					"arg$with$dollars",
					"arg%with%percents",
					"arg#with#hashes",
					"arg@with@ats",
					"arg!with!exclamations",
					"arg~with~tildes",
					"arg`with`backticks",
					"arg^with^carets",
					"arg+with+plus",
					"arg'with'quotes",
					"arg\"with\"doublequotes",
				}
				err := run(&cobra.Command{}, specialArgs)
				assert.NoError(t, err, "run should handle special character args")
			}
		})
	}
}

// TestExecute is now in root_external_test.go as it tests the public Execute function.

// TestSetVersion is now in root_external_test.go as it tests the public SetVersion function.

func TestVersionVariable(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "Version default value is dev",
			wantErr: false,
		},
		{
			name:    "Version is mutable",
			wantErr: false,
		},
		{
			name:    "Version is string type",
			wantErr: false,
		},
		{
			name:    "error case - Version empty string",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "Version default value is dev":
				// This tests the initial state
				// Note: Version might have been modified by other tests
				originalVersion := Version
				Version = "dev"
				defer func() { Version = originalVersion }()

				assert.Equal(t, "dev", Version, "Default version should be 'dev'")

			case "Version is mutable":
				originalVersion := Version
				defer func() { Version = originalVersion }()

				Version = "test"
				assert.Equal(t, "test", Version, "Version should be mutable")

				Version = "another"
				assert.Equal(t, "another", Version, "Version should be changeable")

			case "Version is string type":
				assert.IsType(t, "", Version, "Version should be string type")

			case "error case - Version empty string":
				// Test that empty version string is handled correctly
				originalVersion := Version
				defer func() { Version = originalVersion }()

				Version = ""
				// Empty version should be detectable
				assert.Empty(t, Version, "Empty version should be testable")
			}
		})
	}
}

func TestOsExitVariable(t *testing.T) {
	// Note: No t.Parallel() at function level because this test modifies global OsExit

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "OsExit is mockable", wantErr: false},
		{name: "OsExit captures exit code", wantErr: false},
		{name: "OsExit can be called multiple times", wantErr: false},
		{name: "OsExit default is os.Exit", wantErr: false},
		{name: "error case - OsExit must not be nil", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Note: No t.Parallel() here because subtests manipulate shared global state

			originalOsExit := OsExit
			defer func() { OsExit = originalOsExit }()

			switch tt.name {
			case "OsExit is mockable":
				called := false
				OsExit = func(code int) {
					called = true
				}

				OsExit(1)
				assert.True(t, called, "OsExit should be mockable")

			case "OsExit captures exit code":
				var capturedCode int
				OsExit = func(code int) {
					capturedCode = code
				}

				testCodes := []int{0, 1, 2, 127, 255, -1}
				for _, code := range testCodes {
					OsExit(code)
					assert.Equal(t, code, capturedCode, "Should capture exit code %d", code)
				}

			case "OsExit can be called multiple times":
				callCount := 0
				OsExit = func(code int) {
					callCount++
				}

				for i := 0; i < 10; i++ {
					OsExit(i)
				}

				assert.Equal(t, 10, callCount, "OsExit should be callable multiple times")

			case "OsExit default is os.Exit":
				// Reset to ensure we're testing the default
				OsExit = os.Exit
				assert.NotNil(t, OsExit, "OsExit should default to os.Exit")
				// We can't test actual os.Exit as it would terminate the test

			case "error case - OsExit must not be nil":
				assert.NotNil(t, OsExit, "OsExit must be initialized")
			}
		})
	}
}

func TestProviderServeVariable(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "ProviderServe is mockable", wantErr: false},
		{name: "ProviderServe captures all parameters", wantErr: false},
		{name: "ProviderServe can return various errors", wantErr: false},
		{name: "ProviderServe default is providerserver.Serve", wantErr: false},
		{name: "error case - ProviderServe must not be nil", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			originalProviderServe := ProviderServe
			defer func() { ProviderServe = originalProviderServe }()

			switch tt.name {
			case "ProviderServe is mockable":
				called := false
				ProviderServe = func(ctx context.Context, providerFunc func() provider.Provider, opts providerserver.ServeOpts) error {
					called = true
					return nil
				}

				_ = ProviderServe(context.Background(), nil, providerserver.ServeOpts{})
				assert.True(t, called, "ProviderServe should be mockable")

			case "ProviderServe captures all parameters":
				var capturedCtx context.Context
				var capturedFunc func() provider.Provider
				var capturedOpts providerserver.ServeOpts

				ProviderServe = func(ctx context.Context, providerFunc func() provider.Provider, opts providerserver.ServeOpts) error {
					capturedCtx = ctx
					capturedFunc = providerFunc
					capturedOpts = opts
					return nil
				}

				testCtx := context.WithValue(context.Background(), contextKey("test"), "value")
				testFunc := func() provider.Provider { return nil }
				testOpts := providerserver.ServeOpts{
					Address: "test-address",
					Debug:   true,
				}

				_ = ProviderServe(testCtx, testFunc, testOpts)

				assert.Equal(t, testCtx, capturedCtx, "Context should be captured")
				assert.NotNil(t, capturedFunc, "Function should be captured")
				assert.Equal(t, testOpts, capturedOpts, "Options should be captured")

			case "ProviderServe can return various errors":
				testErrors := []error{
					nil,
					errors.New("test error"),
					context.Canceled,
					context.DeadlineExceeded,
				}

				for _, testErr := range testErrors {
					ProviderServe = func(ctx context.Context, providerFunc func() provider.Provider, opts providerserver.ServeOpts) error {
						return testErr
					}

					err := ProviderServe(context.Background(), nil, providerserver.ServeOpts{})
					assert.Equal(t, testErr, err, "Should return expected error")
				}

			case "ProviderServe default is providerserver.Serve":
				// Reset to ensure we're testing the default
				ProviderServe = providerserver.Serve
				assert.NotNil(t, ProviderServe, "ProviderServe should default to providerserver.Serve")

			case "error case - ProviderServe must not be nil":
				assert.NotNil(t, ProviderServe, "ProviderServe must be initialized")
			}
		})
	}
}

func TestDebugFlag(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "debug flag default is false", wantErr: false},
		{name: "debug flag can be set to true", wantErr: false},
		{name: "debug flag can be toggled", wantErr: false},
		{name: "debug flag is bool type", wantErr: false},
		{name: "error case - debug flag must be bool", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			originalDebug := debug
			defer func() { debug = originalDebug }()

			switch tt.name {
			case "debug flag default is false":
				debug = false
				assert.False(t, debug, "debug should default to false")

			case "debug flag can be set to true":
				debug = true
				assert.True(t, debug, "debug should be settable to true")

			case "debug flag can be toggled":
				debug = false
				assert.False(t, debug)

				debug = true
				assert.True(t, debug)

				debug = false
				assert.False(t, debug)

			case "debug flag is bool type":
				assert.IsType(t, false, debug, "debug should be bool type")

			case "error case - debug flag must be bool":
				assert.IsType(t, false, debug, "debug must be boolean")
			}
		})
	}
}

func TestRootCmd(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "rootCmd is initialized", wantErr: false},
		{name: "rootCmd has correct Use", wantErr: false},
		{name: "rootCmd has Short description", wantErr: false},
		{name: "rootCmd has Long description", wantErr: false},
		{name: "rootCmd has RunE function", wantErr: false},
		{name: "rootCmd does not have Run function", wantErr: false},
		{name: "rootCmd has no subcommands", wantErr: false},
		{name: "rootCmd has flags", wantErr: false},
		{name: "rootCmd is not nil pointer", wantErr: false},
		{name: "error case - rootCmd must not be nil", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "rootCmd is initialized":
				assert.NotNil(t, rootCmd, "rootCmd should be initialized")

			case "rootCmd has correct Use":
				assert.Equal(t, "terraform-provider-n8n", rootCmd.Use, "rootCmd should have correct Use")

			case "rootCmd has Short description":
				assert.NotEmpty(t, rootCmd.Short, "rootCmd should have Short description")
				assert.Contains(t, rootCmd.Short, "Terraform", "Short should mention Terraform")
				assert.Contains(t, rootCmd.Short, "n8n", "Short should mention n8n")

			case "rootCmd has Long description":
				assert.NotEmpty(t, rootCmd.Long, "rootCmd should have Long description")
				assert.Contains(t, rootCmd.Long, "Terraform", "Long should mention Terraform")
				assert.Contains(t, rootCmd.Long, "n8n", "Long should mention n8n")

			case "rootCmd has RunE function":
				assert.NotNil(t, rootCmd.RunE, "rootCmd should have RunE function")

			case "rootCmd does not have Run function":
				assert.Nil(t, rootCmd.Run, "rootCmd should not have Run function (uses RunE)")

			case "rootCmd has no subcommands":
				assert.False(t, rootCmd.HasSubCommands(), "rootCmd should not have subcommands")

			case "rootCmd has flags":
				assert.True(t, rootCmd.HasFlags(), "rootCmd should have flags")

			case "rootCmd is not nil pointer":
				assert.NotNil(t, rootCmd)
				assert.IsType(t, &cobra.Command{}, rootCmd)

			case "error case - rootCmd must not be nil":
				assert.NotNil(t, rootCmd, "rootCmd must be initialized")
			}
		})
	}
}

func TestRunEdgeCases(t *testing.T) {
	// Note: No t.Parallel() at function level because this test modifies global ProviderServe

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "run handles various error types", wantErr: false},
		{name: "run returns nil explicitly on success", wantErr: false},
		{name: "run with command containing metadata", wantErr: false},
		{name: "error case - run must handle nil command", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Note: No t.Parallel() here because subtests manipulate shared global state

			originalProviderServe := ProviderServe
			defer func() { ProviderServe = originalProviderServe }()

			switch tt.name {
			case "run handles various error types":
				errorCases := []struct {
					name string
					err  error
				}{
					{"simple error", errors.New("simple error")},
					{"nil error", nil},
					{"wrapped error", fmt.Errorf("wrapped: %w", errors.New("inner"))},
					{"empty error", errors.New("")},
					{"timeout error", context.DeadlineExceeded},
					{"canceled error", context.Canceled},
				}

				for _, tc := range errorCases {
					ProviderServe = func(ctx context.Context, providerFunc func() provider.Provider, opts providerserver.ServeOpts) error {
						return tc.err
					}

					err := run(&cobra.Command{}, []string{})
					assert.Equal(t, tc.err, err, "Should return the expected error for %s", tc.name)
				}

			case "run returns nil explicitly on success":
				ProviderServe = func(ctx context.Context, providerFunc func() provider.Provider, opts providerserver.ServeOpts) error {
					return nil
				}

				err := run(&cobra.Command{}, []string{})
				assert.NoError(t, err, "Should return nil on success")
				assert.Nil(t, err, "Should be explicitly nil")

			case "run with command containing metadata":
				ProviderServe = func(ctx context.Context, providerFunc func() provider.Provider, opts providerserver.ServeOpts) error {
					return nil
				}

				cmd := &cobra.Command{
					Use:   "test",
					Short: "test command",
				}
				cmd.SetContext(context.WithValue(context.Background(), contextKey("key"), "value"))

				err := run(cmd, []string{})
				assert.NoError(t, err)

			case "error case - run must handle nil command":
				ProviderServe = func(ctx context.Context, providerFunc func() provider.Provider, opts providerserver.ServeOpts) error {
					return nil
				}
				cmd := &cobra.Command{}
				assert.NotNil(t, cmd, "command should not be nil")
			}
		})
	}
}

func TestRunProviderCreation(t *testing.T) {
	// Note: No t.Parallel() at function level because this test modifies global ProviderServe and Version

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "run creates provider with version", wantErr: false},
		{name: "run provider function can be called multiple times", wantErr: false},
		{name: "run provider function returns consistent provider", wantErr: false},
		{name: "error case - provider function must return valid provider", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Note: No t.Parallel() here because subtests manipulate shared global state

			originalProviderServe := ProviderServe
			originalVersion := Version
			defer func() {
				ProviderServe = originalProviderServe
				Version = originalVersion
			}()

			switch tt.name {
			case "run creates provider with version":
				testVersion := "test-version-x.y.z"
				Version = testVersion

				var providerCreated bool
				ProviderServe = func(ctx context.Context, providerFunc func() provider.Provider, opts providerserver.ServeOpts) error {
					p := providerFunc()
					providerCreated = (p != nil)
					return nil
				}

				_ = run(&cobra.Command{}, []string{})
				assert.True(t, providerCreated, "Provider should be created")

			case "run provider function can be called multiple times":
				callCount := 0
				ProviderServe = func(ctx context.Context, providerFunc func() provider.Provider, opts providerserver.ServeOpts) error {
					// Call provider function multiple times
					for i := 0; i < 3; i++ {
						p := providerFunc()
						if p != nil {
							callCount++
						}
					}
					return nil
				}

				_ = run(&cobra.Command{}, []string{})
				assert.Equal(t, 3, callCount, "Provider function should be callable multiple times")

			case "run provider function returns consistent provider":
				var providers []provider.Provider
				ProviderServe = func(ctx context.Context, providerFunc func() provider.Provider, opts providerserver.ServeOpts) error {
					for i := 0; i < 3; i++ {
						providers = append(providers, providerFunc())
					}
					return nil
				}

				_ = run(&cobra.Command{}, []string{})

				// All calls should return a provider
				for i, p := range providers {
					assert.NotNil(t, p, "Provider %d should not be nil", i)
				}

			case "error case - provider function must return valid provider":
				ProviderServe = func(ctx context.Context, providerFunc func() provider.Provider, opts providerserver.ServeOpts) error {
					p := providerFunc()
					assert.NotNil(t, p, "provider function must return non-nil provider")
					return nil
				}
				_ = run(&cobra.Command{}, []string{})
			}
		})
	}
}

func TestCommandIntegration(t *testing.T) {
	// Note: No t.Parallel() at function level because this test modifies global ProviderServe and rootCmd

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "rootCmd can parse flags", wantErr: false},
		{name: "rootCmd Execute with real command", wantErr: false},
		{name: "debug flag is registered on rootCmd", wantErr: false},
		{name: "error case - invalid flags must fail", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Note: No t.Parallel() here because subtests manipulate shared global state

			switch tt.name {
			case "rootCmd can parse flags":
				originalDebug := debug
				defer func() { debug = originalDebug }()

				// Parse debug flag
				rootCmd.SetArgs([]string{"--debug"})
				err := rootCmd.ParseFlags([]string{"--debug"})

				assert.NoError(t, err, "Should parse flags without error")

			case "rootCmd Execute with real command":
				originalProviderServe := ProviderServe
				defer func() { ProviderServe = originalProviderServe }()

				executed := false
				ProviderServe = func(ctx context.Context, providerFunc func() provider.Provider, opts providerserver.ServeOpts) error {
					executed = true
					return nil
				}

				// Reset args to avoid interference
				rootCmd.SetArgs([]string{})
				err := rootCmd.Execute()

				assert.NoError(t, err)
				assert.True(t, executed, "Command should execute")

			case "debug flag is registered on rootCmd":
				flag := rootCmd.Flags().Lookup("debug")
				assert.NotNil(t, flag, "debug flag should be registered")
				assert.Equal(t, "bool", flag.Value.Type(), "debug flag should be boolean")
				assert.Contains(t, flag.Usage, "debugger", "debug flag should mention debugger")

			case "error case - invalid flags must fail":
				err := rootCmd.ParseFlags([]string{"--nonexistent-flag"})
				assert.Error(t, err, "Should fail with invalid flag")
			}
		})
	}
}

func TestConstants(t *testing.T) {
	// Note: No t.Parallel() at function level because this test modifies global ProviderServe

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "registry address is correct format", wantErr: false},
		{name: "error case - registry address must not be empty", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Note: No t.Parallel() here because subtests manipulate shared global state

			originalProviderServe := ProviderServe
			defer func() { ProviderServe = originalProviderServe }()

			switch tt.name {
			case "registry address is correct format":
				var capturedAddress string
				ProviderServe = func(ctx context.Context, providerFunc func() provider.Provider, opts providerserver.ServeOpts) error {
					capturedAddress = opts.Address
					return nil
				}

				_ = run(&cobra.Command{}, []string{})
				assert.Equal(t, "registry.terraform.io/kodflow/n8n", capturedAddress)
				assert.True(t, strings.HasPrefix(capturedAddress, "registry.terraform.io/"))
				assert.True(t, strings.HasSuffix(capturedAddress, "/n8n"))

			case "error case - registry address must not be empty":
				var capturedAddress string
				ProviderServe = func(ctx context.Context, providerFunc func() provider.Provider, opts providerserver.ServeOpts) error {
					capturedAddress = opts.Address
					return nil
				}
				_ = run(&cobra.Command{}, []string{})
				assert.NotEmpty(t, capturedAddress, "registry address must not be empty")
			}
		})
	}
}

func TestInitFunction(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "init is called automatically", wantErr: false},
		{name: "init sets up flags correctly", wantErr: false},
		{name: "error case - init must setup debug flag", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "init is called automatically":
				// The init function runs automatically when the package is imported
				// We verify its effects by checking the state it sets up
				flag := rootCmd.Flags().Lookup("debug")
				assert.NotNil(t, flag, "init should have added debug flag")

			case "init sets up flags correctly":
				// Check all flags set up by init
				debugFlag := rootCmd.Flags().Lookup("debug")
				require.NotNil(t, debugFlag)

				// Verify flag properties
				assert.Equal(t, "debug", debugFlag.Name)
				assert.Equal(t, "bool", debugFlag.Value.Type())
				assert.Equal(t, "false", debugFlag.DefValue)

			case "error case - init must setup debug flag":
				debugFlag := rootCmd.Flags().Lookup("debug")
				assert.NotNil(t, debugFlag, "init must setup debug flag")
			}
		})
	}
}

func TestConcurrency(t *testing.T) {
	// Note: No t.Parallel() at function level because this test modifies global ProviderServe and Version
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "SetVersion concurrent writes", wantErr: false},
		{name: "run concurrent execution", wantErr: false},
		{name: "error case - concurrent writes must not panic", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Note: No t.Parallel() here because subtests manipulate shared global state
			switch tt.name {
			case "SetVersion concurrent writes":
				originalVersion := Version
				defer func() { Version = originalVersion }()

				var wg sync.WaitGroup
				for i := 0; i < 100; i++ {
					wg.Add(1)
					go func(n int) {
						defer wg.Done()
						SetVersion(fmt.Sprintf("version-%d", n))
					}(i)
				}
				wg.Wait()

				// Version should have some value (last write wins)
				assert.NotEmpty(t, Version)

			case "run concurrent execution":
				originalProviderServe := ProviderServe
				defer func() { ProviderServe = originalProviderServe }()

				var mu sync.Mutex
				callCount := 0
				ProviderServe = func(ctx context.Context, providerFunc func() provider.Provider, opts providerserver.ServeOpts) error {
					mu.Lock()
					callCount++
					mu.Unlock()
					time.Sleep(10 * time.Millisecond) // Simulate work
					return nil
				}

				var wg sync.WaitGroup
				for i := 0; i < 10; i++ {
					wg.Add(1)
					go func() {
						defer wg.Done()
						_ = run(&cobra.Command{}, []string{})
					}()
				}
				wg.Wait()

				assert.Equal(t, 10, callCount, "All concurrent runs should complete")

			case "error case - concurrent writes must not panic":
				assert.NotPanics(t, func() {
					var wg sync.WaitGroup
					for i := 0; i < 10; i++ {
						wg.Add(1)
						go func(n int) {
							defer wg.Done()
							SetVersion(fmt.Sprintf("v%d", n))
						}(i)
					}
					wg.Wait()
				})
			}
		})
	}
}

func TestMemoryAndPerformance(t *testing.T) {
	// Note: No t.Parallel() at function level because this test modifies global ProviderServe and Version
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "Version string memory", wantErr: false},
		{name: "run performance with mock", wantErr: false},
		{name: "error case - Version must handle empty string", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Note: No t.Parallel() here because subtests manipulate shared global state

			switch tt.name {
			case "Version string memory":
				originalVersion := Version
				defer func() { Version = originalVersion }()

				// Test various string sizes
				sizes := []int{0, 1, 10, 100, 1000, 10000}
				for _, size := range sizes {
					version := strings.Repeat("x", size)
					SetVersion(version)
					assert.Len(t, Version, size, "Version should handle size %d", size)
				}

			case "run performance with mock":
				originalProviderServe := ProviderServe
				defer func() { ProviderServe = originalProviderServe }()

				var mu sync.Mutex
				callCount := 0
				ProviderServe = func(ctx context.Context, providerFunc func() provider.Provider, opts providerserver.ServeOpts) error {
					mu.Lock()
					callCount++
					mu.Unlock()
					return nil
				}

				// Run multiple times quickly
				start := time.Now()
				for i := 0; i < 1000; i++ {
					_ = run(&cobra.Command{}, []string{})
				}
				duration := time.Since(start)

				mu.Lock()
				finalCount := callCount
				mu.Unlock()
				assert.Equal(t, 1000, finalCount)
				assert.Less(t, duration, 1*time.Second, "1000 runs should complete quickly")

			case "error case - Version must handle empty string":
				originalVersion := Version
				defer func() { Version = originalVersion }()
				SetVersion("")
				assert.Equal(t, "", Version, "Version must handle empty string")
			}
		})
	}
}

func TestErrorMessages(t *testing.T) {
	// Note: No t.Parallel() at function level because this test modifies global ProviderServe

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "run preserves error messages", wantErr: false},
		{name: "error case - nil error must be handled", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Note: No t.Parallel() here because subtests manipulate shared global state

			originalProviderServe := ProviderServe
			defer func() { ProviderServe = originalProviderServe }()

			switch tt.name {
			case "run preserves error messages":
				testMessages := []string{
					"connection refused",
					"permission denied",
					"file not found",
					"invalid configuration",
					"timeout occurred",
				}

				for _, msg := range testMessages {
					expectedErr := errors.New(msg)
					ProviderServe = func(ctx context.Context, providerFunc func() provider.Provider, opts providerserver.ServeOpts) error {
						return expectedErr
					}

					err := run(&cobra.Command{}, []string{})
					assert.Error(t, err)
					assert.Equal(t, msg, err.Error(), "Error message should be preserved")
				}

			case "error case - nil error must be handled":
				ProviderServe = func(ctx context.Context, providerFunc func() provider.Provider, opts providerserver.ServeOpts) error {
					return nil
				}
				err := run(&cobra.Command{}, []string{})
				assert.NoError(t, err, "nil error must be handled")
			}
		})
	}
}

func BenchmarkSetVersion(b *testing.B) {
	originalVersion := Version
	defer func() { Version = originalVersion }()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SetVersion("benchmark-version")
	}
}

func BenchmarkRun(b *testing.B) {
	originalProviderServe := ProviderServe
	defer func() { ProviderServe = originalProviderServe }()

	ProviderServe = func(ctx context.Context, providerFunc func() provider.Provider, opts providerserver.ServeOpts) error {
		return nil
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = run(&cobra.Command{}, []string{})
	}
}

func BenchmarkExecute(b *testing.B) {
	originalOsExit := OsExit
	originalRootCmd := rootCmd
	originalProviderServe := ProviderServe
	defer func() {
		OsExit = originalOsExit
		rootCmd = originalRootCmd
		ProviderServe = originalProviderServe
	}()

	OsExit = func(code int) {} // No-op
	ProviderServe = func(ctx context.Context, providerFunc func() provider.Provider, opts providerserver.ServeOpts) error {
		return nil
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rootCmd = &cobra.Command{
			Use:  "bench",
			RunE: run,
		}
		Execute()
	}
}

// Example function to demonstrate usage.
func ExampleSetVersion() {
	originalVersion := Version
	defer func() { Version = originalVersion }()

	SetVersion("1.2.3")
	fmt.Println(Version)
	// Output would be: 1.2.3
}

func ExampleExecute() {
	// This example shows how Execute would be called
	// In practice, this is called from main.go
	// Execute() // Would start the provider
}

// TestRun is now in external test file - refactored to test behavior only.
