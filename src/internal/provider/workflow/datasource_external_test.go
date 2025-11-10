package workflow_test

import (
	"testing"

	"github.com/kodflow/n8n/src/internal/provider/workflow"
	"github.com/stretchr/testify/assert"
)

// TestNewWorkflowDataSourceExternal tests the public NewWorkflowDataSource function (black-box).
func TestNewWorkflowDataSourceExternal(t *testing.T) {
	t.Run("create new instance from external package", func(t *testing.T) {
		ds := workflow.NewWorkflowDataSource()

		assert.NotNil(t, ds)
	})

	t.Run("multiple calls return independent instances", func(t *testing.T) {
		ds1 := workflow.NewWorkflowDataSource()
		ds2 := workflow.NewWorkflowDataSource()

		assert.NotSame(t, ds1, ds2, "each call should return a new instance")
	})

	t.Run("returned instance is not nil even on repeated calls", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			ds := workflow.NewWorkflowDataSource()
			assert.NotNil(t, ds, "instance should never be nil")
		}
	})
}

// TestNewWorkflowDataSourceWrapperExternal tests the public NewWorkflowDataSourceWrapper function (black-box).
func TestNewWorkflowDataSourceWrapperExternal(t *testing.T) {
	t.Run("create wrapped instance from external package", func(t *testing.T) {
		ds := workflow.NewWorkflowDataSourceWrapper()

		assert.NotNil(t, ds)
	})

	t.Run("multiple calls return independent instances", func(t *testing.T) {
		ds1 := workflow.NewWorkflowDataSourceWrapper()
		ds2 := workflow.NewWorkflowDataSourceWrapper()

		assert.NotSame(t, ds1, ds2, "each call should return a new instance")
	})

	t.Run("returned instance is not nil even on repeated calls", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			ds := workflow.NewWorkflowDataSourceWrapper()
			assert.NotNil(t, ds, "instance should never be nil")
		}
	})
}
