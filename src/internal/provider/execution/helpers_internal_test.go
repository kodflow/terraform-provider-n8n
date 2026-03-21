package execution

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

// TestExtractTagIDs tests the extractTagIDs helper.
func TestExtractTagIDs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		tagIDs      []string
		expectCount int
	}{
		{
			name:        "multiple tags",
			tagIDs:      []string{"tag-1", "tag-2", "tag-3"},
			expectCount: 3,
		},
		{
			name:        "single tag",
			tagIDs:      []string{"tag-1"},
			expectCount: 1,
		},
		{
			name:        "empty set",
			tagIDs:      []string{},
			expectCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			tagSet, diagsCreate := types.SetValueFrom(ctx, types.StringType, tt.tagIDs)
			//: Check condition.
			if diagsCreate.HasError() {
				t.Fatalf("failed to create set: %v", diagsCreate)
			}

			var diags diag.Diagnostics
			result := extractTagIDs(ctx, tagSet, &diags)

			assert.False(t, diags.HasError(), "Should not have diagnostics error")
			assert.Len(t, result, tt.expectCount)
		})
	}
}

// TestExtractTagIDsWithErrorPath tests extractTagIDs with a type mismatch to cover error branches.
func TestExtractTagIDsWithErrorPath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "null set returns empty slice",
			testFunc: func(t *testing.T) {
				t.Helper()
				ctx := t.Context()
				nullSet := types.SetNull(types.StringType)
				var diags diag.Diagnostics
				result := extractTagIDs(ctx, nullSet, &diags)
				assert.False(t, diags.HasError())
				assert.Empty(t, result)
			},
		},
		{
			name: "unknown set returns empty slice",
			testFunc: func(t *testing.T) {
				t.Helper()
				ctx := t.Context()
				unknownSet := types.SetUnknown(types.StringType)
				var diags diag.Diagnostics
				result := extractTagIDs(ctx, unknownSet, &diags)
				// Unknown sets may produce diagnostics or return empty
				assert.Empty(t, result)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.testFunc(t)
		})
	}
}
