package models

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestDataSources(t *testing.T) {
	t.Run("create with projects list", func(t *testing.T) {
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
	})

	t.Run("create with empty projects", func(t *testing.T) {
		datasources := DataSources{
			Projects: []Item{},
		}

		assert.Len(t, datasources.Projects, 0)
		assert.NotNil(t, datasources.Projects)
	})

	t.Run("create with nil projects", func(t *testing.T) {
		datasources := DataSources{
			Projects: nil,
		}

		assert.Nil(t, datasources.Projects)
		assert.Len(t, datasources.Projects, 0)
	})

	t.Run("copy struct", func(t *testing.T) {
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
	})

	t.Run("zero value struct", func(t *testing.T) {
		var datasources DataSources
		assert.Nil(t, datasources.Projects)
		assert.Len(t, datasources.Projects, 0)
	})

	t.Run("multiple projects with different types", func(t *testing.T) {
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
	})

	t.Run("append projects", func(t *testing.T) {
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
	})

	t.Run("projects with partial data", func(t *testing.T) {
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
	})
}

func TestDataSourcesConcurrency(t *testing.T) {
	t.Run("concurrent read", func(t *testing.T) {
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
	})
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
