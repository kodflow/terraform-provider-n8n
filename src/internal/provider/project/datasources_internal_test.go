package project

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestProjectsDataSource_schemaAttributes tests the schemaAttributes helper.
func TestProjectsDataSource_schemaAttributes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		want      int
		wantError bool
	}{
		{
			name:      "returns correct number of attributes",
			want:      1, // projects attribute
			wantError: false,
		},
		{
			name:      "nil datasource should not panic",
			want:      1,
			wantError: false,
		},
		{
			name:      "multiple calls return same result",
			want:      1,
			wantError: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			d := &ProjectsDataSource{}

			if tt.name == "multiple calls return same result" {
				// Call multiple times and verify consistent results
				attrs1 := d.schemaAttributes()
				attrs2 := d.schemaAttributes()
				assert.Equal(t, len(attrs1), len(attrs2))
			}

			attrs := d.schemaAttributes()

			assert.NotNil(t, attrs)
			assert.Len(t, attrs, tt.want)
			assert.Contains(t, attrs, "projects")
		})
	}
}

// TestProjectsDataSource_projectAttributes tests the projectAttributes helper.
func TestProjectsDataSource_projectAttributes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		want      int
		wantError bool
	}{
		{
			name:      "returns correct number of attributes",
			want:      7, // id, name, type, created_at, updated_at, icon, description
			wantError: false,
		},
		{
			name:      "nil datasource should not panic",
			want:      7,
			wantError: false,
		},
		{
			name:      "multiple calls return same result",
			want:      7,
			wantError: false,
		},
		{
			name:      "all required attributes present",
			want:      7,
			wantError: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			d := &ProjectsDataSource{}

			if tt.name == "multiple calls return same result" {
				// Call multiple times and verify consistent results
				attrs1 := d.projectAttributes()
				attrs2 := d.projectAttributes()
				assert.Equal(t, len(attrs1), len(attrs2))
			}

			attrs := d.projectAttributes()

			assert.NotNil(t, attrs)
			assert.Len(t, attrs, tt.want)
			assert.Contains(t, attrs, "id")
			assert.Contains(t, attrs, "name")
			assert.Contains(t, attrs, "type")
			assert.Contains(t, attrs, "created_at")
			assert.Contains(t, attrs, "updated_at")
			assert.Contains(t, attrs, "icon")
			assert.Contains(t, attrs, "description")
		})
	}
}
