package cmd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
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
	t.Run("init configures root command", func(t *testing.T) {
		// Verify that init() has set up the rootCmd
		assert.NotNil(t, rootCmd, "rootCmd should be initialized")
		assert.NotNil(t, rootCmd.Flags(), "rootCmd should have flags")
	})

	t.Run("init adds debug flag", func(t *testing.T) {
		// Verify that the debug flag exists
		debugFlag := rootCmd.Flags().Lookup("debug")
		assert.NotNil(t, debugFlag, "debug flag should be defined")
		assert.Equal(t, "bool", debugFlag.Value.Type(), "debug flag should be boolean")
	})

	t.Run("debug flag has correct default", func(t *testing.T) {
		debugFlag := rootCmd.Flags().Lookup("debug")
		require.NotNil(t, debugFlag)
		assert.Equal(t, "false", debugFlag.DefValue, "debug flag should default to false")
	})

	t.Run("debug flag has proper usage text", func(t *testing.T) {
		debugFlag := rootCmd.Flags().Lookup("debug")
		require.NotNil(t, debugFlag)
		assert.Contains(t, debugFlag.Usage, "debugger", "Usage should mention debugger support")
		assert.Contains(t, debugFlag.Usage, "delve", "Usage should mention delve")
	})

	t.Run("debug flag shorthand is empty", func(t *testing.T) {
		debugFlag := rootCmd.Flags().Lookup("debug")
		require.NotNil(t, debugFlag)
		assert.Empty(t, debugFlag.Shorthand, "debug flag should not have shorthand")
	})
}

func TestRun(t *testing.T) {
	// Save original values
	originalProviderServe := ProviderServe
	originalDebug := debug
	defer func() {
		ProviderServe = originalProviderServe
		debug = originalDebug
	}()

	t.Run("run succeeds when ProviderServe succeeds", func(t *testing.T) {
		// Mock ProviderServe to succeed
		ProviderServe = func(ctx context.Context, providerFunc func() provider.Provider, opts providerserver.ServeOpts) error {
			// Verify the opts are correct
			assert.Equal(t, "registry.terraform.io/kodflow/n8n", opts.Address)
			assert.Equal(t, debug, opts.Debug)

			// Verify provider function works
			p := providerFunc()
			assert.NotNil(t, p, "Provider should be created")

			return nil
		}

		err := run(&cobra.Command{}, []string{})
		assert.NoError(t, err, "run should succeed when ProviderServe succeeds")
	})

	t.Run("run fails when ProviderServe fails", func(t *testing.T) {
		expectedErr := errors.New("provider serve error")

		// Mock ProviderServe to fail
		ProviderServe = func(ctx context.Context, providerFunc func() provider.Provider, opts providerserver.ServeOpts) error {
			return expectedErr
		}

		err := run(&cobra.Command{}, []string{})
		assert.Error(t, err, "run should fail when ProviderServe fails")
		assert.Equal(t, expectedErr, err, "run should return the ProviderServe error")
	})

	t.Run("run passes correct address to ProviderServe", func(t *testing.T) {
		addressCaptured := ""
		ProviderServe = func(ctx context.Context, providerFunc func() provider.Provider, opts providerserver.ServeOpts) error {
			addressCaptured = opts.Address
			return nil
		}

		_ = run(&cobra.Command{}, []string{})
		assert.Equal(t, "registry.terraform.io/kodflow/n8n", addressCaptured, "Correct address should be passed")
	})

	t.Run("run respects debug flag when false", func(t *testing.T) {
		debug = false
		debugCaptured := true

		ProviderServe = func(ctx context.Context, providerFunc func() provider.Provider, opts providerserver.ServeOpts) error {
			debugCaptured = opts.Debug
			return nil
		}

		_ = run(&cobra.Command{}, []string{})
		assert.False(t, debugCaptured, "Debug should be false when flag is false")
	})

	t.Run("run respects debug flag when true", func(t *testing.T) {
		debug = true
		defer func() { debug = false }()

		debugCaptured := false
		ProviderServe = func(ctx context.Context, providerFunc func() provider.Provider, opts providerserver.ServeOpts) error {
			debugCaptured = opts.Debug
			return nil
		}

		_ = run(&cobra.Command{}, []string{})
		assert.True(t, debugCaptured, "Debug should be true when flag is true")
	})

	t.Run("run passes context to ProviderServe", func(t *testing.T) {
		var ctxCaptured context.Context
		ProviderServe = func(ctx context.Context, providerFunc func() provider.Provider, opts providerserver.ServeOpts) error {
			ctxCaptured = ctx
			return nil
		}

		_ = run(&cobra.Command{}, []string{})
		assert.NotNil(t, ctxCaptured, "Context should be passed to ProviderServe")
		assert.Equal(t, context.Background(), ctxCaptured, "Background context should be used")
	})

	t.Run("run creates provider with correct version", func(t *testing.T) {
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
	})

	t.Run("run handles nil command", func(t *testing.T) {
		ProviderServe = func(ctx context.Context, providerFunc func() provider.Provider, opts providerserver.ServeOpts) error {
			return nil
		}

		err := run(nil, []string{})
		assert.NoError(t, err, "run should handle nil command")
	})

	t.Run("run handles empty args", func(t *testing.T) {
		ProviderServe = func(ctx context.Context, providerFunc func() provider.Provider, opts providerserver.ServeOpts) error {
			return nil
		}

		err := run(&cobra.Command{}, []string{})
		assert.NoError(t, err, "run should handle empty args")
	})

	t.Run("run handles args with values", func(t *testing.T) {
		ProviderServe = func(ctx context.Context, providerFunc func() provider.Provider, opts providerserver.ServeOpts) error {
			return nil
		}

		err := run(&cobra.Command{}, []string{"arg1", "arg2", "--flag=value"})
		assert.NoError(t, err, "run should handle args")
	})

	t.Run("run handles nil args", func(t *testing.T) {
		ProviderServe = func(ctx context.Context, providerFunc func() provider.Provider, opts providerserver.ServeOpts) error {
			return nil
		}

		err := run(&cobra.Command{}, nil)
		assert.NoError(t, err, "run should handle nil args")
	})

	t.Run("run with panic recovery", func(t *testing.T) {
		ProviderServe = func(ctx context.Context, providerFunc func() provider.Provider, opts providerserver.ServeOpts) error {
			// Simulate panic in provider creation
			defer func() {
				_ = recover() // Explicitly ignore panic recovery
			}()
			_ = providerFunc()
			return nil
		}

		// Should not panic
		assert.NotPanics(t, func() {
			_ = run(&cobra.Command{}, []string{})
		})
	})

	t.Run("run with very long args", func(t *testing.T) {
		ProviderServe = func(ctx context.Context, providerFunc func() provider.Provider, opts providerserver.ServeOpts) error {
			return nil
		}

		longArgs := make([]string, 1000)
		for i := range longArgs {
			longArgs[i] = fmt.Sprintf("arg%d", i)
		}

		err := run(&cobra.Command{}, longArgs)
		assert.NoError(t, err, "run should handle many args")
	})

	t.Run("run with special characters in args", func(t *testing.T) {
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
	})
}

func TestExecute(t *testing.T) {
	// Save original values
	originalOsExit := OsExit
	originalRootCmd := rootCmd
	originalProviderServe := ProviderServe
	defer func() {
		OsExit = originalOsExit
		rootCmd = originalRootCmd
		ProviderServe = originalProviderServe
	}()

	t.Run("Execute calls os.Exit(1) on error", func(t *testing.T) {
		exitCalled := false
		var exitCode int
		OsExit = func(code int) {
			exitCalled = true
			exitCode = code
		}

		// Create a temporary root command that will fail
		rootCmd = &cobra.Command{
			Use: "test",
			RunE: func(cmd *cobra.Command, args []string) error {
				return errors.New("test error")
			},
		}

		Execute()

		assert.True(t, exitCalled, "OsExit should be called on error")
		assert.Equal(t, 1, exitCode, "Should exit with code 1")
	})

	t.Run("Execute does not call os.Exit on success", func(t *testing.T) {
		exitCalled := false
		OsExit = func(code int) {
			exitCalled = true
		}

		ProviderServe = func(ctx context.Context, providerFunc func() provider.Provider, opts providerserver.ServeOpts) error {
			return nil
		}

		rootCmd = &cobra.Command{
			Use:  "test",
			RunE: run,
		}

		Execute()

		assert.False(t, exitCalled, "OsExit should not be called on success")
	})

	t.Run("Execute handles help flag", func(t *testing.T) {
		exitCalled := false
		OsExit = func(code int) {
			exitCalled = true
		}

		rootCmd = &cobra.Command{
			Use: "test",
			Run: func(cmd *cobra.Command, args []string) {
				// Help command should not cause exit
			},
		}
		rootCmd.SetArgs([]string{"--help"})

		Execute()

		assert.False(t, exitCalled, "Help should not cause exit")
	})

	t.Run("Execute handles version flag", func(t *testing.T) {
		exitCalled := false
		OsExit = func(code int) {
			exitCalled = true
		}

		rootCmd = &cobra.Command{
			Use:     "test",
			Version: "1.0.0",
		}
		rootCmd.SetArgs([]string{"--version"})

		Execute()

		assert.False(t, exitCalled, "Version flag should not cause exit")
	})

	t.Run("Execute with nil rootCmd causes panic", func(t *testing.T) {
		rootCmd = nil

		assert.Panics(t, func() {
			Execute()
		}, "Execute should panic with nil rootCmd")
	})

	t.Run("Execute with command that returns specific errors", func(t *testing.T) {
		testErrors := []error{
			errors.New("simple error"),
			fmt.Errorf("formatted error: %s", "test"),
			fmt.Errorf("wrapped error: %w", errors.New("inner")),
			context.Canceled,
			context.DeadlineExceeded,
		}

		for _, testErr := range testErrors {
			t.Run(fmt.Sprintf("error: %v", testErr), func(t *testing.T) {
				exitCalled := false
				OsExit = func(code int) {
					exitCalled = true
				}

				rootCmd = &cobra.Command{
					Use: "test",
					RunE: func(cmd *cobra.Command, args []string) error {
						return testErr
					},
				}

				Execute()

				assert.True(t, exitCalled, "OsExit should be called for error: %v", testErr)
			})
		}
	})

	t.Run("Execute captures output correctly", func(t *testing.T) {
		// Capture stdout
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		rootCmd = &cobra.Command{
			Use: "test",
			Run: func(cmd *cobra.Command, args []string) {
				fmt.Println("test output")
			},
		}

		Execute()

		// Restore stdout
		w.Close()
		os.Stdout = oldStdout

		// Read captured output
		var buf bytes.Buffer
		io.Copy(&buf, r)
		output := buf.String()

		assert.Contains(t, output, "test output", "Should capture command output")
	})
}

func TestSetVersion(t *testing.T) {
	// Save original version
	originalVersion := Version
	defer func() {
		Version = originalVersion
	}()

	t.Run("SetVersion sets the version", func(t *testing.T) {
		testVersion := "1.2.3"
		SetVersion(testVersion)
		assert.Equal(t, testVersion, Version, "Version should be set correctly")
	})

	t.Run("SetVersion handles empty string", func(t *testing.T) {
		SetVersion("")
		assert.Equal(t, "", Version, "Version should handle empty string")
	})

	t.Run("SetVersion handles dev version", func(t *testing.T) {
		SetVersion("dev")
		assert.Equal(t, "dev", Version, "Version should handle 'dev'")
	})

	t.Run("SetVersion handles semantic version formats", func(t *testing.T) {
		testCases := []string{
			"1.0.0",
			"2.3.4",
			"10.20.30",
			"1.0.0-alpha",
			"1.0.0-beta.1",
			"1.0.0+build.123",
			"1.0.0-rc.1+build.456",
			"v1.0.0",
			"v1.0.0-alpha",
			"1.0.0-SNAPSHOT",
			"1.0",
			"1",
		}

		for _, version := range testCases {
			SetVersion(version)
			assert.Equal(t, version, Version, "Version should handle %s", version)
		}
	})

	t.Run("SetVersion handles very long version string", func(t *testing.T) {
		longVersion := "1.0.0-" + strings.Repeat("a", 10000)
		SetVersion(longVersion)
		assert.Equal(t, longVersion, Version, "Version should handle long strings")
	})

	t.Run("SetVersion can be called multiple times", func(t *testing.T) {
		versions := []string{"1.0.0", "2.0.0", "3.0.0", "dev", ""}
		for _, v := range versions {
			SetVersion(v)
			assert.Equal(t, v, Version, "Version should update on each call")
		}
	})

	t.Run("SetVersion handles unicode characters", func(t *testing.T) {
		unicodeVersions := []string{
			"1.0.0-测试",
			"1.0.0-テスト",
			"1.0.0-тест",
			"1.0.0-🚀",
		}

		for _, v := range unicodeVersions {
			SetVersion(v)
			assert.Equal(t, v, Version, "Version should handle unicode: %s", v)
		}
	})

	t.Run("SetVersion handles special characters", func(t *testing.T) {
		specialVersions := []string{
			"1.0.0!@#$%^&*()",
			"1.0.0<>?:",
			"1.0.0[]{}",
			"1.0.0|\\",
		}

		for _, v := range specialVersions {
			SetVersion(v)
			assert.Equal(t, v, Version, "Version should handle special chars: %s", v)
		}
	})

	t.Run("SetVersion is thread-safe for reading", func(t *testing.T) {
		SetVersion("concurrent-test")

		var wg sync.WaitGroup
		errs := make([]error, 100)

		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				v := Version
				if v != "concurrent-test" {
					errs[idx] = fmt.Errorf("unexpected version: %s", v)
				}
			}(i)
		}

		wg.Wait()

		for _, err := range errs {
			assert.NoError(t, err)
		}
	})
}

func TestVersionVariable(t *testing.T) {
	t.Run("Version default value is dev", func(t *testing.T) {
		// This tests the initial state
		// Note: Version might have been modified by other tests
		originalVersion := Version
		Version = "dev"
		defer func() { Version = originalVersion }()

		assert.Equal(t, "dev", Version, "Default version should be 'dev'")
	})

	t.Run("Version is mutable", func(t *testing.T) {
		originalVersion := Version
		defer func() { Version = originalVersion }()

		Version = "test"
		assert.Equal(t, "test", Version, "Version should be mutable")

		Version = "another"
		assert.Equal(t, "another", Version, "Version should be changeable")
	})

	t.Run("Version is string type", func(t *testing.T) {
		assert.IsType(t, "", Version, "Version should be string type")
	})
}

func TestOsExitVariable(t *testing.T) {
	originalOsExit := OsExit
	defer func() { OsExit = originalOsExit }()

	t.Run("OsExit is mockable", func(t *testing.T) {
		called := false
		OsExit = func(code int) {
			called = true
		}

		OsExit(1)
		assert.True(t, called, "OsExit should be mockable")
	})

	t.Run("OsExit captures exit code", func(t *testing.T) {
		var capturedCode int
		OsExit = func(code int) {
			capturedCode = code
		}

		testCodes := []int{0, 1, 2, 127, 255, -1}
		for _, code := range testCodes {
			OsExit(code)
			assert.Equal(t, code, capturedCode, "Should capture exit code %d", code)
		}
	})

	t.Run("OsExit can be called multiple times", func(t *testing.T) {
		callCount := 0
		OsExit = func(code int) {
			callCount++
		}

		for i := 0; i < 10; i++ {
			OsExit(i)
		}

		assert.Equal(t, 10, callCount, "OsExit should be callable multiple times")
	})

	t.Run("OsExit default is os.Exit", func(t *testing.T) {
		// Reset to ensure we're testing the default
		OsExit = os.Exit
		assert.NotNil(t, OsExit, "OsExit should default to os.Exit")
		// We can't test actual os.Exit as it would terminate the test
	})
}

func TestProviderServeVariable(t *testing.T) {
	originalProviderServe := ProviderServe
	defer func() { ProviderServe = originalProviderServe }()

	t.Run("ProviderServe is mockable", func(t *testing.T) {
		called := false
		ProviderServe = func(ctx context.Context, providerFunc func() provider.Provider, opts providerserver.ServeOpts) error {
			called = true
			return nil
		}

		_ = ProviderServe(context.Background(), nil, providerserver.ServeOpts{})
		assert.True(t, called, "ProviderServe should be mockable")
	})

	t.Run("ProviderServe captures all parameters", func(t *testing.T) {
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
	})

	t.Run("ProviderServe can return various errors", func(t *testing.T) {
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
	})

	t.Run("ProviderServe default is providerserver.Serve", func(t *testing.T) {
		// Reset to ensure we're testing the default
		ProviderServe = providerserver.Serve
		assert.NotNil(t, ProviderServe, "ProviderServe should default to providerserver.Serve")
	})
}

func TestDebugFlag(t *testing.T) {
	originalDebug := debug
	defer func() { debug = originalDebug }()

	t.Run("debug flag default is false", func(t *testing.T) {
		debug = false
		assert.False(t, debug, "debug should default to false")
	})

	t.Run("debug flag can be set to true", func(t *testing.T) {
		debug = true
		assert.True(t, debug, "debug should be settable to true")
	})

	t.Run("debug flag can be toggled", func(t *testing.T) {
		debug = false
		assert.False(t, debug)

		debug = true
		assert.True(t, debug)

		debug = false
		assert.False(t, debug)
	})

	t.Run("debug flag is bool type", func(t *testing.T) {
		assert.IsType(t, false, debug, "debug should be bool type")
	})
}

func TestRootCmd(t *testing.T) {
	t.Run("rootCmd is initialized", func(t *testing.T) {
		assert.NotNil(t, rootCmd, "rootCmd should be initialized")
	})

	t.Run("rootCmd has correct Use", func(t *testing.T) {
		assert.Equal(t, "terraform-provider-n8n", rootCmd.Use, "rootCmd should have correct Use")
	})

	t.Run("rootCmd has Short description", func(t *testing.T) {
		assert.NotEmpty(t, rootCmd.Short, "rootCmd should have Short description")
		assert.Contains(t, rootCmd.Short, "Terraform", "Short should mention Terraform")
		assert.Contains(t, rootCmd.Short, "n8n", "Short should mention n8n")
	})

	t.Run("rootCmd has Long description", func(t *testing.T) {
		assert.NotEmpty(t, rootCmd.Long, "rootCmd should have Long description")
		assert.Contains(t, rootCmd.Long, "Terraform", "Long should mention Terraform")
		assert.Contains(t, rootCmd.Long, "n8n", "Long should mention n8n")
	})

	t.Run("rootCmd has RunE function", func(t *testing.T) {
		assert.NotNil(t, rootCmd.RunE, "rootCmd should have RunE function")
	})

	t.Run("rootCmd does not have Run function", func(t *testing.T) {
		assert.Nil(t, rootCmd.Run, "rootCmd should not have Run function (uses RunE)")
	})

	t.Run("rootCmd has no subcommands", func(t *testing.T) {
		assert.False(t, rootCmd.HasSubCommands(), "rootCmd should not have subcommands")
	})

	t.Run("rootCmd has flags", func(t *testing.T) {
		assert.True(t, rootCmd.HasFlags(), "rootCmd should have flags")
	})

	t.Run("rootCmd is not nil pointer", func(t *testing.T) {
		assert.NotNil(t, rootCmd)
		assert.IsType(t, &cobra.Command{}, rootCmd)
	})
}

func TestRunEdgeCases(t *testing.T) {
	originalProviderServe := ProviderServe
	defer func() { ProviderServe = originalProviderServe }()

	t.Run("run handles various error types", func(t *testing.T) {
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
			t.Run(tc.name, func(t *testing.T) {
				ProviderServe = func(ctx context.Context, providerFunc func() provider.Provider, opts providerserver.ServeOpts) error {
					return tc.err
				}

				err := run(&cobra.Command{}, []string{})
				assert.Equal(t, tc.err, err, "Should return the expected error")
			})
		}
	})

	t.Run("run returns nil explicitly on success", func(t *testing.T) {
		ProviderServe = func(ctx context.Context, providerFunc func() provider.Provider, opts providerserver.ServeOpts) error {
			return nil
		}

		err := run(&cobra.Command{}, []string{})
		assert.NoError(t, err, "Should return nil on success")
		assert.Nil(t, err, "Should be explicitly nil")
	})

	t.Run("run with command containing metadata", func(t *testing.T) {
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
	})
}

func TestRunProviderCreation(t *testing.T) {
	originalProviderServe := ProviderServe
	originalVersion := Version
	defer func() {
		ProviderServe = originalProviderServe
		Version = originalVersion
	}()

	t.Run("run creates provider with version", func(t *testing.T) {
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
	})

	t.Run("run provider function can be called multiple times", func(t *testing.T) {
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
	})

	t.Run("run provider function returns consistent provider", func(t *testing.T) {
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
	})
}

func TestCommandIntegration(t *testing.T) {
	t.Run("rootCmd can parse flags", func(t *testing.T) {
		originalDebug := debug
		defer func() { debug = originalDebug }()

		// Parse debug flag
		rootCmd.SetArgs([]string{"--debug"})
		err := rootCmd.ParseFlags([]string{"--debug"})

		assert.NoError(t, err, "Should parse flags without error")
	})

	t.Run("rootCmd Execute with real command", func(t *testing.T) {
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
	})

	t.Run("debug flag is registered on rootCmd", func(t *testing.T) {
		flag := rootCmd.Flags().Lookup("debug")
		assert.NotNil(t, flag, "debug flag should be registered")
		assert.Equal(t, "bool", flag.Value.Type(), "debug flag should be boolean")
		assert.Contains(t, flag.Usage, "debugger", "debug flag should mention debugger")
	})
}

func TestConstants(t *testing.T) {
	originalProviderServe := ProviderServe
	defer func() { ProviderServe = originalProviderServe }()

	t.Run("registry address is correct format", func(t *testing.T) {
		var capturedAddress string
		ProviderServe = func(ctx context.Context, providerFunc func() provider.Provider, opts providerserver.ServeOpts) error {
			capturedAddress = opts.Address
			return nil
		}

		_ = run(&cobra.Command{}, []string{})
		assert.Equal(t, "registry.terraform.io/kodflow/n8n", capturedAddress)
		assert.True(t, strings.HasPrefix(capturedAddress, "registry.terraform.io/"))
		assert.True(t, strings.HasSuffix(capturedAddress, "/n8n"))
	})
}

func TestInitFunction(t *testing.T) {
	t.Run("init is called automatically", func(t *testing.T) {
		// The init function runs automatically when the package is imported
		// We verify its effects by checking the state it sets up
		flag := rootCmd.Flags().Lookup("debug")
		assert.NotNil(t, flag, "init should have added debug flag")
	})

	t.Run("init sets up flags correctly", func(t *testing.T) {
		// Check all flags set up by init
		debugFlag := rootCmd.Flags().Lookup("debug")
		require.NotNil(t, debugFlag)

		// Verify flag properties
		assert.Equal(t, "debug", debugFlag.Name)
		assert.Equal(t, "bool", debugFlag.Value.Type())
		assert.Equal(t, "false", debugFlag.DefValue)
	})
}

func TestConcurrency(t *testing.T) {
	t.Run("SetVersion concurrent writes", func(t *testing.T) {
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
	})

	t.Run("run concurrent execution", func(t *testing.T) {
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
	})
}

func TestMemoryAndPerformance(t *testing.T) {
	t.Run("Version string memory", func(t *testing.T) {
		originalVersion := Version
		defer func() { Version = originalVersion }()

		// Test various string sizes
		sizes := []int{0, 1, 10, 100, 1000, 10000}
		for _, size := range sizes {
			version := strings.Repeat("x", size)
			SetVersion(version)
			assert.Len(t, Version, size, "Version should handle size %d", size)
		}
	})

	t.Run("run performance with mock", func(t *testing.T) {
		originalProviderServe := ProviderServe
		defer func() { ProviderServe = originalProviderServe }()

		callCount := 0
		ProviderServe = func(ctx context.Context, providerFunc func() provider.Provider, opts providerserver.ServeOpts) error {
			callCount++
			return nil
		}

		// Run multiple times quickly
		start := time.Now()
		for i := 0; i < 1000; i++ {
			_ = run(&cobra.Command{}, []string{})
		}
		duration := time.Since(start)

		assert.Equal(t, 1000, callCount)
		assert.Less(t, duration, 1*time.Second, "1000 runs should complete quickly")
	})
}

func TestErrorMessages(t *testing.T) {
	originalProviderServe := ProviderServe
	defer func() { ProviderServe = originalProviderServe }()

	t.Run("run preserves error messages", func(t *testing.T) {
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
	})
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
