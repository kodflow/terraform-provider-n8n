package provider

import (
	"testing"
)

func TestProviderNew(t *testing.T) {
	provider := New("test")()
	if provider == nil {
		t.Fatal("Provider should not be nil")
	}
}

func TestProviderMetadata(t *testing.T) {
	provider := New("test")()
	if provider == nil {
		t.Fatal("Provider should not be nil")
	}

	// Basic smoke test - provider should be instantiable
	// Full integration tests will be added once terraform-plugin-go supports Go 1.25
}
