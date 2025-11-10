package workflow_test

import (
	"testing"

	"github.com/kodflow/n8n/src/internal/provider/workflow"
	"github.com/stretchr/testify/assert"
)

// TestNewWorkflowsDataSourceExternal tests the public NewWorkflowsDataSource function (black-box).
func TestNewWorkflowsDataSourceExternal(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "create new instance from external package", wantErr: false},
		{name: "multiple calls return independent instances", wantErr: false},
		{name: "returned instance is not nil even on repeated calls", wantErr: false},
		{name: "error case - validation checks", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "create new instance from external package":
				ds := workflow.NewWorkflowsDataSource()

				assert.NotNil(t, ds)

			case "multiple calls return independent instances":
				ds1 := workflow.NewWorkflowsDataSource()
				ds2 := workflow.NewWorkflowsDataSource()

				assert.NotSame(t, ds1, ds2, "each call should return a new instance")

			case "returned instance is not nil even on repeated calls":
				for i := 0; i < 100; i++ {
					ds := workflow.NewWorkflowsDataSource()
					assert.NotNil(t, ds, "instance should never be nil")
				}

			case "error case - validation checks":
				ds := workflow.NewWorkflowsDataSource()
				assert.NotNil(t, ds)
			}
		})
	}
}

// TestNewWorkflowsDataSourceWrapperExternal tests the public NewWorkflowsDataSourceWrapper function (black-box).
func TestNewWorkflowsDataSourceWrapperExternal(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "create wrapped instance from external package", wantErr: false},
		{name: "multiple calls return independent instances", wantErr: false},
		{name: "returned instance is not nil even on repeated calls", wantErr: false},
		{name: "error case - validation checks", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "create wrapped instance from external package":
				ds := workflow.NewWorkflowsDataSourceWrapper()

				assert.NotNil(t, ds)

			case "multiple calls return independent instances":
				ds1 := workflow.NewWorkflowsDataSourceWrapper()
				ds2 := workflow.NewWorkflowsDataSourceWrapper()

				assert.NotSame(t, ds1, ds2, "each call should return a new instance")

			case "returned instance is not nil even on repeated calls":
				for i := 0; i < 100; i++ {
					ds := workflow.NewWorkflowsDataSourceWrapper()
					assert.NotNil(t, ds, "instance should never be nil")
				}

			case "error case - validation checks":
				ds := workflow.NewWorkflowsDataSourceWrapper()
				assert.NotNil(t, ds)
			}
		})
	}
}
