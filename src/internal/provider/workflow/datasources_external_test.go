package workflow_test

import (
	"testing"

	"github.com/kodflow/n8n/src/internal/provider/workflow"
	"github.com/stretchr/testify/assert"
)

// TestNewWorkflowsDataSourceExternal tests the public NewWorkflowsDataSource function (black-box).
func TestNewWorkflowsDataSourceExternal(t *testing.T) {
	t.Run("create new instance from external package", func(t *testing.T) {
		ds := workflow.NewWorkflowsDataSource()

		assert.NotNil(t, ds)
	})

	t.Run("multiple calls return independent instances", func(t *testing.T) {
		ds1 := workflow.NewWorkflowsDataSource()
		ds2 := workflow.NewWorkflowsDataSource()

		assert.NotSame(t, ds1, ds2, "each call should return a new instance")
	})

	t.Run("returned instance is not nil even on repeated calls", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			ds := workflow.NewWorkflowsDataSource()
			assert.NotNil(t, ds, "instance should never be nil")
		}
	})
}

// TestNewWorkflowsDataSourceWrapperExternal tests the public NewWorkflowsDataSourceWrapper function (black-box).
func TestNewWorkflowsDataSourceWrapperExternal(t *testing.T) {
	t.Run("create wrapped instance from external package", func(t *testing.T) {
		ds := workflow.NewWorkflowsDataSourceWrapper()

		assert.NotNil(t, ds)
	})

	t.Run("multiple calls return independent instances", func(t *testing.T) {
		ds1 := workflow.NewWorkflowsDataSourceWrapper()
		ds2 := workflow.NewWorkflowsDataSourceWrapper()

		assert.NotSame(t, ds1, ds2, "each call should return a new instance")
	})

	t.Run("returned instance is not nil even on repeated calls", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			ds := workflow.NewWorkflowsDataSourceWrapper()
			assert.NotNil(t, ds, "instance should never be nil")
		}
	})
}
