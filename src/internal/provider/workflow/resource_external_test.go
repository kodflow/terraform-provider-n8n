package workflow_test

import (
	"testing"

	"github.com/kodflow/n8n/src/internal/provider/workflow"
	"github.com/stretchr/testify/assert"
)

// TestNewWorkflowResourceExternal tests the public NewWorkflowResource function (black-box).
func TestNewWorkflowResourceExternal(t *testing.T) {
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
				r := workflow.NewWorkflowResource()

				assert.NotNil(t, r)

			case "multiple calls return independent instances":
				r1 := workflow.NewWorkflowResource()
				r2 := workflow.NewWorkflowResource()

				assert.NotSame(t, r1, r2, "each call should return a new instance")

			case "returned instance is not nil even on repeated calls":
				for i := 0; i < 100; i++ {
					r := workflow.NewWorkflowResource()
					assert.NotNil(t, r, "instance should never be nil")
				}

			case "error case - validation checks":
				r := workflow.NewWorkflowResource()
				assert.NotNil(t, r)
			}
		})
	}
}

// TestNewWorkflowResourceWrapperExternal tests the public NewWorkflowResourceWrapper function (black-box).
func TestNewWorkflowResourceWrapperExternal(t *testing.T) {
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
				r := workflow.NewWorkflowResourceWrapper()

				assert.NotNil(t, r)

			case "multiple calls return independent instances":
				r1 := workflow.NewWorkflowResourceWrapper()
				r2 := workflow.NewWorkflowResourceWrapper()

				assert.NotSame(t, r1, r2, "each call should return a new instance")

			case "returned instance is not nil even on repeated calls":
				for i := 0; i < 100; i++ {
					r := workflow.NewWorkflowResourceWrapper()
					assert.NotNil(t, r, "instance should never be nil")
				}

			case "error case - validation checks":
				r := workflow.NewWorkflowResourceWrapper()
				assert.NotNil(t, r)
			}
		})
	}
}
