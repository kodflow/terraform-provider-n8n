package project

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestProjectsDataSource_schemaAttributes tests the schemaAttributes method.
func TestProjectsDataSource_schemaAttributes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "returns valid schema attributes",
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
			case "returns valid schema attributes":
				ds := &ProjectsDataSource{}
				attrs := ds.schemaAttributes()
				assert.NotNil(t, attrs)
				assert.Contains(t, attrs, "projects")
				assert.NotNil(t, attrs["projects"])

			case "error case - validation checks":
				ds := &ProjectsDataSource{}
				attrs := ds.schemaAttributes()
				assert.NotNil(t, attrs)
				assert.NotEmpty(t, attrs)
			}
		})
	}
}

// TestProjectsDataSource_projectAttributes tests the projectAttributes method.
func TestProjectsDataSource_projectAttributes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "returns valid project attributes",
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
			case "returns valid project attributes":
				ds := &ProjectsDataSource{}
				attrs := ds.projectAttributes()
				assert.NotNil(t, attrs)
				assert.Contains(t, attrs, "id")
				assert.Contains(t, attrs, "name")
				assert.NotNil(t, attrs["id"])
				assert.NotNil(t, attrs["name"])

			case "error case - validation checks":
				ds := &ProjectsDataSource{}
				attrs := ds.projectAttributes()
				assert.NotNil(t, attrs)
				assert.NotEmpty(t, attrs)
			}
		})
	}
}
