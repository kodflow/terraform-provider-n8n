package provider_test

import (
	"testing"

	"github.com/kodflow/n8n/src/internal/provider"
)

func TestProviderNew(t *testing.T) {
	t.Run("creates provider with valid version", func(t *testing.T) {
		p := provider.New("test")()
		if p == nil {
			t.Fatal("Provider should not be nil")
		}
	})

	t.Run("creates provider with empty version", func(t *testing.T) {
		p := provider.New("")()
		if p == nil {
			t.Fatal("Provider should not be nil even with empty version")
		}
	})

	t.Run("creates provider with dev version", func(t *testing.T) {
		p := provider.New("dev")()
		if p == nil {
			t.Fatal("Provider should not be nil with dev version")
		}
	})
}

func TestProviderMetadata(t *testing.T) {
	t.Run("provider is instantiable", func(t *testing.T) {
		p := provider.New("test")()
		if p == nil {
			t.Fatal("Provider should not be nil")
		}
	})

	// Basic smoke test - provider should be instantiable
	// Full integration tests will be added once terraform-plugin-go supports Go 1.25
}
