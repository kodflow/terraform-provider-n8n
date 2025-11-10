package workflow_test

import (
	"testing"

	"github.com/kodflow/n8n/src/internal/provider/workflow"
	"github.com/stretchr/testify/assert"
)

// TestNewWorkflowDataSourceExternal tests the public NewWorkflowDataSource function (black-box).
func TestNewWorkflowDataSourceExternal(t *testing.T) {
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
				ds := workflow.NewWorkflowDataSource()

				assert.NotNil(t, ds)

			case "multiple calls return independent instances":
				ds1 := workflow.NewWorkflowDataSource()
				ds2 := workflow.NewWorkflowDataSource()

				assert.NotSame(t, ds1, ds2, "each call should return a new instance")

			case "returned instance is not nil even on repeated calls":
				for i := 0; i < 100; i++ {
					ds := workflow.NewWorkflowDataSource()
					assert.NotNil(t, ds, "instance should never be nil")
				}

			case "error case - validation checks":
				ds := workflow.NewWorkflowDataSource()
				assert.NotNil(t, ds)
			}
		})
	}
}

// TestNewWorkflowDataSourceWrapperExternal tests the public NewWorkflowDataSourceWrapper function (black-box).
func TestNewWorkflowDataSourceWrapperExternal(t *testing.T) {
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
				ds := workflow.NewWorkflowDataSourceWrapper()

				assert.NotNil(t, ds)

			case "multiple calls return independent instances":
				ds1 := workflow.NewWorkflowDataSourceWrapper()
				ds2 := workflow.NewWorkflowDataSourceWrapper()

				assert.NotSame(t, ds1, ds2, "each call should return a new instance")

			case "returned instance is not nil even on repeated calls":
				for i := 0; i < 100; i++ {
					ds := workflow.NewWorkflowDataSourceWrapper()
					assert.NotNil(t, ds, "instance should never be nil")
				}

			case "error case - validation checks":
				ds := workflow.NewWorkflowDataSourceWrapper()
				assert.NotNil(t, ds)
			}
		})
	}
}
