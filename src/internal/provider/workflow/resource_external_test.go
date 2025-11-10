package workflow_test

import (
	"testing"

	"github.com/kodflow/n8n/src/internal/provider/workflow"
	"github.com/stretchr/testify/assert"
)

// TestNewWorkflowResourceExternal tests the public NewWorkflowResource function (black-box).
func TestNewWorkflowResourceExternal(t *testing.T) {
	t.Run("create new instance from external package", func(t *testing.T) {
		r := workflow.NewWorkflowResource()

		assert.NotNil(t, r)
	})

	t.Run("multiple calls return independent instances", func(t *testing.T) {
		r1 := workflow.NewWorkflowResource()
		r2 := workflow.NewWorkflowResource()

		assert.NotSame(t, r1, r2, "each call should return a new instance")
	})

	t.Run("returned instance is not nil even on repeated calls", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			r := workflow.NewWorkflowResource()
			assert.NotNil(t, r, "instance should never be nil")
		}
	})
}

// TestNewWorkflowResourceWrapperExternal tests the public NewWorkflowResourceWrapper function (black-box).
func TestNewWorkflowResourceWrapperExternal(t *testing.T) {
	t.Run("create wrapped instance from external package", func(t *testing.T) {
		r := workflow.NewWorkflowResourceWrapper()

		assert.NotNil(t, r)
	})

	t.Run("multiple calls return independent instances", func(t *testing.T) {
		r1 := workflow.NewWorkflowResourceWrapper()
		r2 := workflow.NewWorkflowResourceWrapper()

		assert.NotSame(t, r1, r2, "each call should return a new instance")
	})

	t.Run("returned instance is not nil even on repeated calls", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			r := workflow.NewWorkflowResourceWrapper()
			assert.NotNil(t, r, "instance should never be nil")
		}
	})
}
