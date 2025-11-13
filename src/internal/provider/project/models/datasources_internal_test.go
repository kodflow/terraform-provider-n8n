// Package models defines data structures for project resources.
package models

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestDataSources(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "create with projects list", wantErr: false},
		{name: "create with empty projects", wantErr: false},
		{name: "create with nil projects", wantErr: false},
		{name: "copy struct", wantErr: false},
		{name: "zero value struct", wantErr: false},
		{name: "multiple projects with different types", wantErr: false},
		{name: "append projects", wantErr: false},
		{name: "projects with partial data", wantErr: false},
		{name: "error case - validation checks", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "create with projects list":
				datasources := DataSources{
					Projects: []Item{
						{
							ID:          types.StringValue("proj-1"),
							Name:        types.StringValue("Project 1"),
							Type:        types.StringValue("personal"),
							CreatedAt:   types.StringValue("2024-01-01T00:00:00Z"),
							UpdatedAt:   types.StringValue("2024-01-01T00:00:00Z"),
							Icon:        types.StringValue("📁"),
							Description: types.StringValue("First project"),
						},
						{
							ID:          types.StringValue("proj-2"),
							Name:        types.StringValue("Project 2"),
							Type:        types.StringValue("team"),
							CreatedAt:   types.StringValue("2024-01-02T00:00:00Z"),
							UpdatedAt:   types.StringValue("2024-01-02T00:00:00Z"),
							Icon:        types.StringValue("🚀"),
							Description: types.StringValue("Second project"),
						},
					},
				}

				assert.Len(t, datasources.Projects, 2)
				assert.Equal(t, "proj-1", datasources.Projects[0].ID.ValueString())
				assert.Equal(t, "Project 1", datasources.Projects[0].Name.ValueString())
				assert.Equal(t, "proj-2", datasources.Projects[1].ID.ValueString())
				assert.Equal(t, "Project 2", datasources.Projects[1].Name.ValueString())

			case "create with empty projects":
				datasources := DataSources{
					Projects: []Item{},
				}

				assert.Len(t, datasources.Projects, 0)
				assert.NotNil(t, datasources.Projects)

			case "create with nil projects":
				datasources := DataSources{
					Projects: nil,
				}

				assert.Nil(t, datasources.Projects)
				assert.Len(t, datasources.Projects, 0)

			case "copy struct":
				original := DataSources{
					Projects: []Item{
						{
							ID:   types.StringValue("original-proj"),
							Name: types.StringValue("Original"),
						},
					},
				}

				copied := original

				assert.Equal(t, original.Projects[0].ID.ValueString(), copied.Projects[0].ID.ValueString())
				assert.Equal(t, original.Projects[0].Name.ValueString(), copied.Projects[0].Name.ValueString())

				// Note: slice is shared between original and copied
				copied.Projects[0].ID = types.StringValue("modified-proj")
				assert.Equal(t, "modified-proj", original.Projects[0].ID.ValueString())

			case "zero value struct":
				var datasources DataSources
				assert.Nil(t, datasources.Projects)
				assert.Len(t, datasources.Projects, 0)

			case "multiple projects with different types":
				datasources := DataSources{
					Projects: []Item{
						{
							ID:   types.StringValue("proj-personal"),
							Name: types.StringValue("Personal Project"),
							Type: types.StringValue("personal"),
						},
						{
							ID:   types.StringValue("proj-team"),
							Name: types.StringValue("Team Project"),
							Type: types.StringValue("team"),
						},
						{
							ID:   types.StringValue("proj-org"),
							Name: types.StringValue("Org Project"),
							Type: types.StringValue("organization"),
						},
					},
				}

				assert.Len(t, datasources.Projects, 3)
				assert.Equal(t, "personal", datasources.Projects[0].Type.ValueString())
				assert.Equal(t, "team", datasources.Projects[1].Type.ValueString())
				assert.Equal(t, "organization", datasources.Projects[2].Type.ValueString())

			case "append projects":
				datasources := DataSources{
					Projects: []Item{
						{
							ID:   types.StringValue("proj-1"),
							Name: types.StringValue("First"),
						},
					},
				}

				assert.Len(t, datasources.Projects, 1)

				datasources.Projects = append(datasources.Projects, Item{
					ID:   types.StringValue("proj-2"),
					Name: types.StringValue("Second"),
				})

				assert.Len(t, datasources.Projects, 2)
				assert.Equal(t, "proj-2", datasources.Projects[1].ID.ValueString())
				assert.Equal(t, "Second", datasources.Projects[1].Name.ValueString())

			case "projects with partial data":
				datasources := DataSources{
					Projects: []Item{
						{
							ID:   types.StringValue("proj-full"),
							Name: types.StringValue("Full Data"),
							Type: types.StringValue("personal"),
							Icon: types.StringValue("📁"),
						},
						{
							ID:   types.StringValue("proj-minimal"),
							Name: types.StringValue("Minimal Data"),
							// Other fields null
						},
					},
				}

				assert.Len(t, datasources.Projects, 2)
				assert.False(t, datasources.Projects[0].Type.IsNull())
				assert.False(t, datasources.Projects[0].Icon.IsNull())
				assert.True(t, datasources.Projects[1].Type.IsNull())
				assert.True(t, datasources.Projects[1].Icon.IsNull())

			case "error case - validation checks":
				datasources := DataSources{
					Projects: []Item{
						{
							ID:   types.StringValue(""),
							Name: types.StringValue(""),
							Type: types.StringValue("invalid-type"),
						},
					},
				}
				assert.Len(t, datasources.Projects, 1)
				assert.Equal(t, "", datasources.Projects[0].ID.ValueString())
			}
		})
	}
}

func TestDataSourcesConcurrency(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "concurrent read", wantErr: false},
		{name: "error case - concurrent access validation", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// NO t.Parallel() here - concurrent goroutines don't work well with t.Parallel()
			switch tt.name {
			case "concurrent read":
				datasources := DataSources{
					Projects: []Item{
						{
							ID:   types.StringValue("proj-1"),
							Name: types.StringValue("Project 1"),
						},
						{
							ID:   types.StringValue("proj-2"),
							Name: types.StringValue("Project 2"),
						},
					},
				}

				done := make(chan bool, 100)
				for i := 0; i < 100; i++ {
					go func() {
						_ = len(datasources.Projects)
						if len(datasources.Projects) > 0 {
							_ = datasources.Projects[0].ID.ValueString()
							_ = datasources.Projects[0].Name.ValueString()
						}
						done <- true
					}()
				}

				for i := 0; i < 100; i++ {
					<-done
				}

			case "error case - concurrent access validation":
				datasources := DataSources{
					Projects: []Item{
						{
							ID:   types.StringValue("val-id"),
							Name: types.StringValue("Val Project"),
						},
					},
				}

				done := make(chan bool, 50)
				for i := 0; i < 50; i++ {
					go func() {
						_ = len(datasources.Projects)
						if len(datasources.Projects) > 0 {
							_ = datasources.Projects[0].ID.ValueString()
						}
						done <- true
					}()
				}

				for i := 0; i < 50; i++ {
					<-done
				}
			}
		})
	}
}

func BenchmarkDataSources(b *testing.B) {
	b.Run("create", func(b *testing.B) {
		now := time.Now().Format(time.RFC3339)
		for i := 0; i < b.N; i++ {
			_ = DataSources{
				Projects: []Item{
					{
						ID:          types.StringValue("proj-1"),
						Name:        types.StringValue("Project 1"),
						Type:        types.StringValue("personal"),
						CreatedAt:   types.StringValue(now),
						UpdatedAt:   types.StringValue(now),
						Icon:        types.StringValue("📁"),
						Description: types.StringValue("Test project"),
					},
				},
			}
		}
	})

	b.Run("access fields", func(b *testing.B) {
		datasources := DataSources{
			Projects: []Item{
				{
					ID:   types.StringValue("proj-1"),
					Name: types.StringValue("Project 1"),
					Type: types.StringValue("personal"),
				},
			},
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = len(datasources.Projects)
			if len(datasources.Projects) > 0 {
				_ = datasources.Projects[0].ID.ValueString()
				_ = datasources.Projects[0].Name.ValueString()
				_ = datasources.Projects[0].Type.ValueString()
			}
		}
	})
}
