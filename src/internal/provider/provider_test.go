package provider_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	n8nprovider "github.com/kodflow/n8n/src/internal/provider"
)

func TestProviderNew(t *testing.T) {
	tests := []struct {
		name        string
		version     string
		expectError bool
		checkResult func(*testing.T, provider.Provider)
	}{
		{
			name:        "creates provider with valid version",
			version:     "1.0.0",
			expectError: false,
			checkResult: func(t *testing.T, p provider.Provider) {
				t.Helper()
				if p == nil {
					t.Fatal("Provider should not be nil")
				}
			},
		},
		{
			name:        "creates provider with empty version",
			version:     "",
			expectError: false,
			checkResult: func(t *testing.T, p provider.Provider) {
				t.Helper()
				if p == nil {
					t.Fatal("Provider should not be nil even with empty version")
				}
			},
		},
		{
			name:        "creates provider with dev version",
			version:     "dev",
			expectError: false,
			checkResult: func(t *testing.T, p provider.Provider) {
				t.Helper()
				if p == nil {
					t.Fatal("Provider should not be nil with dev version")
				}
			},
		},
		{
			name:        "creates provider with semver version",
			version:     "2.3.4-rc1",
			expectError: false,
			checkResult: func(t *testing.T, p provider.Provider) {
				t.Helper()
				if p == nil {
					t.Fatal("Provider should not be nil with semver version")
				}
			},
		},
		{
			name:        "creates provider with long version string",
			version:     "1.0.0-alpha.beta.gamma.delta.epsilon",
			expectError: false,
			checkResult: func(t *testing.T, p provider.Provider) {
				t.Helper()
				if p == nil {
					t.Fatal("Provider should not be nil with long version")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory := n8nprovider.New(tt.version)
			if factory == nil {
				t.Fatal("Factory function should not be nil")
			}

			p := factory()
			if tt.expectError {
				if p != nil {
					t.Errorf("Expected error but got provider: %v", p)
				}
			} else {
				if tt.checkResult != nil {
					tt.checkResult(t, p)
				}
			}
		})
	}
}

func TestProviderMetadata(t *testing.T) {
	tests := []struct {
		name            string
		version         string
		expectError     bool
		expectedType    string
		expectedVersion string
	}{
		{
			name:            "returns correct metadata with valid version",
			version:         "1.0.0",
			expectError:     false,
			expectedType:    "n8n",
			expectedVersion: "1.0.0",
		},
		{
			name:            "returns correct metadata with empty version",
			version:         "",
			expectError:     false,
			expectedType:    "n8n",
			expectedVersion: "",
		},
		{
			name:            "returns correct metadata with dev version",
			version:         "dev",
			expectError:     false,
			expectedType:    "n8n",
			expectedVersion: "dev",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := n8nprovider.New(tt.version)()
			if p == nil {
				t.Fatal("Provider should not be nil")
			}

			var resp provider.MetadataResponse
			p.Metadata(context.Background(), provider.MetadataRequest{}, &resp)

			if tt.expectError {
				if resp.TypeName != "" {
					t.Errorf("Expected error but got TypeName: %s", resp.TypeName)
				}
			} else {
				if resp.TypeName != tt.expectedType {
					t.Errorf("Expected TypeName %q, got %q", tt.expectedType, resp.TypeName)
				}
				if resp.Version != tt.expectedVersion {
					t.Errorf("Expected Version %q, got %q", tt.expectedVersion, resp.Version)
				}
			}
		})
	}
}
