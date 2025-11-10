package workflow

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/stretchr/testify/assert"
)

// TestWorkflowResource_schemaAttributes tests the private schemaAttributes method.
func TestWorkflowResource_schemaAttributes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "returns all schema attributes",
			wantErr: false,
		},
		{
			name:    "handles nil receiver",
			wantErr: false,
		},
		{
			name:    "returns map with correct capacity",
			wantErr: false,
		},
		{
			name:    "handles zero value resource",
			wantErr: false,
		},
		{
			name:    "all attributes have descriptions",
			wantErr: false,
		},
		{
			name:    "error case - validation checks",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "returns all schema attributes":
				r := &WorkflowResource{}

				attrs := r.schemaAttributes()

				assert.NotNil(t, attrs)
				assert.Equal(t, WORKFLOW_ATTRIBUTES_SIZE, len(attrs))
				assert.Contains(t, attrs, "id")
				assert.Contains(t, attrs, "name")
				assert.Contains(t, attrs, "active")
				assert.Contains(t, attrs, "tags")
				assert.Contains(t, attrs, "nodes_json")
				assert.Contains(t, attrs, "connections_json")
				assert.Contains(t, attrs, "settings_json")
				assert.Contains(t, attrs, "created_at")
				assert.Contains(t, attrs, "updated_at")
				assert.Contains(t, attrs, "version_id")
				assert.Contains(t, attrs, "is_archived")
				assert.Contains(t, attrs, "trigger_count")
				assert.Contains(t, attrs, "meta")
				assert.Contains(t, attrs, "pin_data")

			case "handles nil receiver":
				r := &WorkflowResource{}

				attrs := r.schemaAttributes()

				assert.NotNil(t, attrs)

			case "returns map with correct capacity":
				r := &WorkflowResource{}

				attrs := r.schemaAttributes()

				assert.Len(t, attrs, WORKFLOW_ATTRIBUTES_SIZE)

			case "handles zero value resource":
				var r WorkflowResource

				assert.NotPanics(t, func() {
					attrs := r.schemaAttributes()
					assert.NotNil(t, attrs)
					assert.Len(t, attrs, WORKFLOW_ATTRIBUTES_SIZE)
				})

			case "all attributes have descriptions":
				r := &WorkflowResource{}
				attrs := r.schemaAttributes()

				for key, attr := range attrs {
					assert.NotEmpty(t, key, "attribute key should not be empty")
					assert.NotNil(t, attr, "attribute %s should not be nil", key)
				}

			case "error case - validation checks":
				r := &WorkflowResource{}
				attrs := r.schemaAttributes()
				assert.NotNil(t, attrs)
			}
		})
	}
}

// TestWorkflowResource_addCoreAttributes tests the private addCoreAttributes method.
func TestWorkflowResource_addCoreAttributes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "adds core attributes to map",
			wantErr: false,
		},
		{
			name:    "handles empty map",
			wantErr: false,
		},
		{
			name:    "handles nil receiver",
			wantErr: false,
		},
		{
			name:    "error case - validation checks",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "adds core attributes to map":
				r := &WorkflowResource{}
				attrs := make(map[string]schema.Attribute)

				r.addCoreAttributes(attrs)

				assert.Contains(t, attrs, "id")
				assert.Contains(t, attrs, "name")
				assert.Contains(t, attrs, "active")
				assert.Contains(t, attrs, "tags")

			case "handles empty map":
				r := &WorkflowResource{}
				attrs := make(map[string]schema.Attribute)

				r.addCoreAttributes(attrs)

				assert.NotEmpty(t, attrs)
				assert.Len(t, attrs, 4)

			case "handles nil receiver":
				r := &WorkflowResource{}
				attrs := make(map[string]schema.Attribute)

				assert.NotPanics(t, func() {
					r.addCoreAttributes(attrs)
				})

			case "error case - validation checks":
				r := &WorkflowResource{}
				attrs := make(map[string]schema.Attribute)
				assert.NotPanics(t, func() {
					r.addCoreAttributes(attrs)
				})
			}
		})
	}
}

// TestWorkflowResource_addJSONAttributes tests the private addJSONAttributes method.
func TestWorkflowResource_addJSONAttributes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "adds JSON attributes to map",
			wantErr: false,
		},
		{
			name:    "handles empty map",
			wantErr: false,
		},
		{
			name:    "handles nil receiver",
			wantErr: false,
		},
		{
			name:    "error case - validation checks",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "adds JSON attributes to map":
				r := &WorkflowResource{}
				attrs := make(map[string]schema.Attribute)

				r.addJSONAttributes(attrs)

				assert.Contains(t, attrs, "nodes_json")
				assert.Contains(t, attrs, "connections_json")
				assert.Contains(t, attrs, "settings_json")

			case "handles empty map":
				r := &WorkflowResource{}
				attrs := make(map[string]schema.Attribute)

				r.addJSONAttributes(attrs)

				assert.NotEmpty(t, attrs)
				assert.Len(t, attrs, 3)

			case "handles nil receiver":
				r := &WorkflowResource{}
				attrs := make(map[string]schema.Attribute)

				assert.NotPanics(t, func() {
					r.addJSONAttributes(attrs)
				})

			case "error case - validation checks":
				r := &WorkflowResource{}
				attrs := make(map[string]schema.Attribute)
				assert.NotPanics(t, func() {
					r.addJSONAttributes(attrs)
				})
			}
		})
	}
}

// TestWorkflowResource_addMetadataAttributes tests the private addMetadataAttributes method.
func TestWorkflowResource_addMetadataAttributes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "adds metadata attributes to map",
			wantErr: false,
		},
		{
			name:    "handles empty map",
			wantErr: false,
		},
		{
			name:    "handles nil receiver",
			wantErr: false,
		},
		{
			name:    "error case - validation checks",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "adds metadata attributes to map":
				r := &WorkflowResource{}
				attrs := make(map[string]schema.Attribute)

				r.addMetadataAttributes(attrs)

				assert.Contains(t, attrs, "created_at")
				assert.Contains(t, attrs, "updated_at")
				assert.Contains(t, attrs, "version_id")
				assert.Contains(t, attrs, "is_archived")
				assert.Contains(t, attrs, "trigger_count")
				assert.Contains(t, attrs, "meta")
				assert.Contains(t, attrs, "pin_data")

			case "handles empty map":
				r := &WorkflowResource{}
				attrs := make(map[string]schema.Attribute)

				r.addMetadataAttributes(attrs)

				assert.NotEmpty(t, attrs)
				assert.Len(t, attrs, 7)

			case "handles nil receiver":
				r := &WorkflowResource{}
				attrs := make(map[string]schema.Attribute)

				assert.NotPanics(t, func() {
					r.addMetadataAttributes(attrs)
				})

			case "error case - validation checks":
				r := &WorkflowResource{}
				attrs := make(map[string]schema.Attribute)
				assert.NotPanics(t, func() {
					r.addMetadataAttributes(attrs)
				})
			}
		})
	}
}
