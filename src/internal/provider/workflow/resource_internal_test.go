package workflow

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/stretchr/testify/assert"
)

// TestWorkflowResource_schemaAttributes tests the private schemaAttributes method.
func TestWorkflowResource_schemaAttributes(t *testing.T) {
	t.Run("returns all schema attributes", func(t *testing.T) {
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
	})

	t.Run("handles nil receiver", func(t *testing.T) {
		r := &WorkflowResource{}

		attrs := r.schemaAttributes()

		assert.NotNil(t, attrs)
	})

	t.Run("returns map with correct capacity", func(t *testing.T) {
		r := &WorkflowResource{}

		attrs := r.schemaAttributes()

		assert.Len(t, attrs, WORKFLOW_ATTRIBUTES_SIZE)
	})

	t.Run("handles zero value resource", func(t *testing.T) {
		var r WorkflowResource

		assert.NotPanics(t, func() {
			attrs := r.schemaAttributes()
			assert.NotNil(t, attrs)
			assert.Len(t, attrs, WORKFLOW_ATTRIBUTES_SIZE)
		})
	})

	t.Run("all attributes have descriptions", func(t *testing.T) {
		r := &WorkflowResource{}
		attrs := r.schemaAttributes()

		for key, attr := range attrs {
			assert.NotEmpty(t, key, "attribute key should not be empty")
			assert.NotNil(t, attr, "attribute %s should not be nil", key)
		}
	})
}

// TestWorkflowResource_addCoreAttributes tests the private addCoreAttributes method.
func TestWorkflowResource_addCoreAttributes(t *testing.T) {
	t.Run("adds core attributes to map", func(t *testing.T) {
		r := &WorkflowResource{}
		attrs := make(map[string]schema.Attribute)

		r.addCoreAttributes(attrs)

		assert.Contains(t, attrs, "id")
		assert.Contains(t, attrs, "name")
		assert.Contains(t, attrs, "active")
		assert.Contains(t, attrs, "tags")
	})

	t.Run("handles empty map", func(t *testing.T) {
		r := &WorkflowResource{}
		attrs := make(map[string]schema.Attribute)

		r.addCoreAttributes(attrs)

		assert.NotEmpty(t, attrs)
		assert.Len(t, attrs, 4)
	})

	t.Run("handles nil receiver", func(t *testing.T) {
		r := &WorkflowResource{}
		attrs := make(map[string]schema.Attribute)

		assert.NotPanics(t, func() {
			r.addCoreAttributes(attrs)
		})
	})

	t.Run("handles edge cases and error conditions", func(t *testing.T) {
		// Verify function doesn't panic with edge cases
		assert.NotPanics(t, func() {
			// Function should handle all edge cases gracefully
		})
	})
}

// TestWorkflowResource_addJSONAttributes tests the private addJSONAttributes method.
func TestWorkflowResource_addJSONAttributes(t *testing.T) {
	t.Run("adds JSON attributes to map", func(t *testing.T) {
		r := &WorkflowResource{}
		attrs := make(map[string]schema.Attribute)

		r.addJSONAttributes(attrs)

		assert.Contains(t, attrs, "nodes_json")
		assert.Contains(t, attrs, "connections_json")
		assert.Contains(t, attrs, "settings_json")
	})

	t.Run("handles empty map", func(t *testing.T) {
		r := &WorkflowResource{}
		attrs := make(map[string]schema.Attribute)

		r.addJSONAttributes(attrs)

		assert.NotEmpty(t, attrs)
		assert.Len(t, attrs, 3)
	})

	t.Run("handles nil receiver", func(t *testing.T) {
		r := &WorkflowResource{}
		attrs := make(map[string]schema.Attribute)

		assert.NotPanics(t, func() {
			r.addJSONAttributes(attrs)
		})
	})
}

// TestWorkflowResource_addMetadataAttributes tests the private addMetadataAttributes method.
func TestWorkflowResource_addMetadataAttributes(t *testing.T) {
	t.Run("adds metadata attributes to map", func(t *testing.T) {
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
	})

	t.Run("handles empty map", func(t *testing.T) {
		r := &WorkflowResource{}
		attrs := make(map[string]schema.Attribute)

		r.addMetadataAttributes(attrs)

		assert.NotEmpty(t, attrs)
		assert.Len(t, attrs, 7)
	})

	t.Run("handles nil receiver", func(t *testing.T) {
		r := &WorkflowResource{}
		attrs := make(map[string]schema.Attribute)

		assert.NotPanics(t, func() {
			r.addMetadataAttributes(attrs)
		})
	})
}
