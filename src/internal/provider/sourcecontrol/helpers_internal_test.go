package sourcecontrol

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/sourcecontrol/models"
	"github.com/stretchr/testify/assert"
)

// TestMarshalImportedSlice tests the marshalImportedSlice helper.
func TestMarshalImportedSlice(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		items    any
		isNull   bool
		contains string
	}{
		{
			name:   "nil slice returns null",
			items:  nil,
			isNull: true,
		},
		{
			name:   "empty slice returns null",
			items:  []any{},
			isNull: true,
		},
		{
			name:     "non-empty slice returns JSON",
			items:    []any{map[string]any{"id": "wf-1"}},
			isNull:   false,
			contains: "wf-1",
		},
		{
			name:   "unmarshalable value returns null",
			items:  make(chan int),
			isNull: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := marshalImportedSlice(tt.items)
			if tt.isNull {
				assert.True(t, result.IsNull() || result.ValueString() == "null" || result.ValueString() == "[]")
			} else {
				assert.False(t, result.IsNull())
				assert.Contains(t, result.ValueString(), tt.contains)
			}
		})
	}
}

// TestMapImportResultToResource tests the mapImportResultToResource helper.
func TestMapImportResultToResource(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                     string
		result                   *n8nsdk.ImportResult
		expectWorkflowsNotNull   bool
		expectCredentialsNotNull bool
	}{
		{
			name: "maps populated result correctly",
			result: &n8nsdk.ImportResult{
				Workflows: []n8nsdk.ImportResultWorkflowsInner{
					{Id: strPtr("wf-1")},
				},
				Credentials: []n8nsdk.ImportResultCredentialsInner{
					{Id: strPtr("cred-1")},
				},
				Tags: &n8nsdk.ImportResultTags{
					Tags: []n8nsdk.ImportResultWorkflowsInner{{Id: strPtr("tag-1")}},
				},
				Variables: &n8nsdk.ImportResultVariables{
					Added: []string{"var-1"},
				},
			},
			expectWorkflowsNotNull:   true,
			expectCredentialsNotNull: true,
		},
		{
			name:                     "empty result produces null fields",
			result:                   &n8nsdk.ImportResult{},
			expectWorkflowsNotNull:   false,
			expectCredentialsNotNull: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			data := &models.Resource{}
			mapImportResultToResource(tt.result, data)
			if tt.expectWorkflowsNotNull {
				assert.False(t, data.WorkflowsImported.IsNull())
			} else {
				assert.True(t, data.WorkflowsImported.IsNull())
			}
			if tt.expectCredentialsNotNull {
				assert.False(t, data.CredentialsImported.IsNull())
			} else {
				assert.True(t, data.CredentialsImported.IsNull())
			}
		})
	}
}

// TestMapImportResultToResource_NilResult tests nil result handling.
func TestMapImportResultToResource_NilResult(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
	}{
		{name: "nil workflows stays null"},
		{name: "nil credentials stays null"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			data := &models.Resource{}
			result := &n8nsdk.ImportResult{
				Workflows:   nil,
				Credentials: nil,
				Tags:        nil,
				Variables:   nil,
			}
			mapImportResultToResource(result, data)
			assert.Equal(t, types.StringNull(), data.WorkflowsImported)
			assert.Equal(t, types.StringNull(), data.CredentialsImported)
			assert.Equal(t, types.StringNull(), data.TagsImported)
			assert.Equal(t, types.StringNull(), data.VariablesImported)
		})
	}
}

// strPtr returns a pointer to a string value for test use.
func strPtr(v string) *string { return new(v) }
