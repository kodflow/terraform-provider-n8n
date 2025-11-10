package project

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/project/models"
	"github.com/stretchr/testify/assert"
)

func TestFindProjectByIDOrName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "find by ID", wantErr: false},
		{name: "find by name", wantErr: false},
		{name: "find by ID and name (ID takes precedence)", wantErr: false},
		{name: "not found", wantErr: false},
		{name: "empty projects list", wantErr: false},
		{name: "null search parameters", wantErr: false},
		{name: "project with nil ID", wantErr: false},
		{name: "case sensitive name matching", wantErr: false},
		{name: "error case - validation checks", wantErr: false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "find by ID":
				id1 := "proj-1"
				id2 := "proj-2"
				projects := []n8nsdk.Project{
					{Id: &id1, Name: "Project One"},
					{Id: &id2, Name: "Project Two"},
				}
				found, ok := findProjectByIDOrName(projects, types.StringValue("proj-2"), types.StringNull())
				assert.True(t, ok)
				assert.NotNil(t, found)
				assert.Equal(t, "proj-2", *found.Id)
				assert.Equal(t, "Project Two", found.Name)

			case "find by name":
				id1 := "proj-1"
				id2 := "proj-2"
				projects := []n8nsdk.Project{
					{Id: &id1, Name: "Project One"},
					{Id: &id2, Name: "Project Two"},
				}
				found, ok := findProjectByIDOrName(projects, types.StringNull(), types.StringValue("Project One"))
				assert.True(t, ok)
				assert.NotNil(t, found)
				assert.Equal(t, "proj-1", *found.Id)
				assert.Equal(t, "Project One", found.Name)

			case "find by ID and name (ID takes precedence)":
				id1 := "proj-1"
				id2 := "proj-2"
				projects := []n8nsdk.Project{
					{Id: &id1, Name: "Project One"},
					{Id: &id2, Name: "Project Two"},
				}
				found, ok := findProjectByIDOrName(projects, types.StringValue("proj-1"), types.StringValue("Project Two"))
				assert.True(t, ok)
				assert.NotNil(t, found)
				assert.Equal(t, "proj-1", *found.Id)
				assert.Equal(t, "Project One", found.Name)

			case "not found":
				id1 := "proj-1"
				projects := []n8nsdk.Project{{Id: &id1, Name: "Project One"}}
				found, ok := findProjectByIDOrName(projects, types.StringValue("proj-999"), types.StringValue("Non-existent"))
				assert.False(t, ok)
				assert.Nil(t, found)

			case "empty projects list":
				projects := []n8nsdk.Project{}
				found, ok := findProjectByIDOrName(projects, types.StringValue("any-id"), types.StringValue("any-name"))
				assert.False(t, ok)
				assert.Nil(t, found)

			case "null search parameters":
				id1 := "proj-1"
				projects := []n8nsdk.Project{{Id: &id1, Name: "Project One"}}
				found, ok := findProjectByIDOrName(projects, types.StringNull(), types.StringNull())
				assert.False(t, ok)
				assert.Nil(t, found)

			case "project with nil ID":
				projects := []n8nsdk.Project{{Id: nil, Name: "Project Without ID"}}
				// Should not find by ID
				found, ok := findProjectByIDOrName(projects, types.StringValue("any-id"), types.StringNull())
				assert.False(t, ok)
				assert.Nil(t, found)
				// Should find by name
				found, ok = findProjectByIDOrName(projects, types.StringNull(), types.StringValue("Project Without ID"))
				assert.True(t, ok)
				assert.NotNil(t, found)
				assert.Equal(t, "Project Without ID", found.Name)

			case "case sensitive name matching":
				id1 := "proj-1"
				projects := []n8nsdk.Project{{Id: &id1, Name: "Project One"}}
				// Exact match should work
				found, ok := findProjectByIDOrName(projects, types.StringNull(), types.StringValue("Project One"))
				assert.True(t, ok)
				assert.NotNil(t, found)
				// Different case should not match
				found, ok = findProjectByIDOrName(projects, types.StringNull(), types.StringValue("project one"))
				assert.False(t, ok)
				assert.Nil(t, found)

			case "error case - validation checks":
				projects := []n8nsdk.Project{}
				found, ok := findProjectByIDOrName(projects, types.StringNull(), types.StringNull())
				assert.False(t, ok)
				assert.Nil(t, found)
			}
		})
	}
}

func TestMapProjectToDataSourceModel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "map with all fields", wantErr: false},
		{name: "map with nil fields", wantErr: false},
		{name: "overwrite existing data", wantErr: false},
		{name: "error case - validation checks", wantErr: false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "map with all fields":
				id := "proj-123"
				projectType := "personal"
				createdAt := time.Now()
				updatedAt := time.Now().Add(1 * time.Hour)
				icon := "📁"
				description := "Test project description"
				project := &n8nsdk.Project{
					Id:          &id,
					Name:        "Test Project",
					Type:        &projectType,
					CreatedAt:   &createdAt,
					UpdatedAt:   &updatedAt,
					Icon:        &icon,
					Description: &description,
				}
				data := &models.DataSource{}
				mapProjectToDataSourceModel(project, data)
				assert.Equal(t, "proj-123", data.ID.ValueString())
				assert.Equal(t, "Test Project", data.Name.ValueString())
				assert.Equal(t, "personal", data.Type.ValueString())
				assert.Equal(t, createdAt.String(), data.CreatedAt.ValueString())
				assert.Equal(t, updatedAt.String(), data.UpdatedAt.ValueString())
				assert.Equal(t, "📁", data.Icon.ValueString())
				assert.Equal(t, "Test project description", data.Description.ValueString())

			case "map with nil fields":
				project := &n8nsdk.Project{Name: "Minimal Project"}
				data := &models.DataSource{}
				mapProjectToDataSourceModel(project, data)
				assert.True(t, data.ID.IsNull())
				assert.Equal(t, "Minimal Project", data.Name.ValueString())
				assert.True(t, data.Type.IsNull())
				assert.True(t, data.CreatedAt.IsNull())
				assert.True(t, data.UpdatedAt.IsNull())
				assert.True(t, data.Icon.IsNull())
				assert.True(t, data.Description.IsNull())

			case "overwrite existing data":
				id := "new-id"
				projectType := "team"
				data := &models.DataSource{
					ID:          types.StringValue("old-id"),
					Name:        types.StringValue("Old Name"),
					Type:        types.StringValue("personal"),
					Icon:        types.StringValue("🔧"),
					Description: types.StringValue("Old description"),
				}
				project := &n8nsdk.Project{
					Id:   &id,
					Name: "New Project",
					Type: &projectType,
				}
				mapProjectToDataSourceModel(project, data)
				assert.Equal(t, "new-id", data.ID.ValueString())
				assert.Equal(t, "New Project", data.Name.ValueString())
				assert.Equal(t, "team", data.Type.ValueString())
				assert.True(t, data.Icon.IsNull())
				assert.True(t, data.Description.IsNull())

			case "error case - validation checks":
				project := &n8nsdk.Project{Name: "Test"}
				data := &models.DataSource{}
				mapProjectToDataSourceModel(project, data)
				assert.NotNil(t, data)
			}
		})
	}
}

func TestMapProjectToItem(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "map with all fields", wantErr: false},
		{name: "map with minimal fields", wantErr: false},
		{name: "map empty string values", wantErr: false},
		{name: "map special characters", wantErr: false},
		{name: "error case - validation checks", wantErr: false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "map with all fields":
				id := "proj-456"
				projectType := "organization"
				createdAt := time.Now()
				updatedAt := time.Now().Add(2 * time.Hour)
				icon := "🚀"
				description := "Item project description"
				project := &n8nsdk.Project{
					Id:          &id,
					Name:        "Item Project",
					Type:        &projectType,
					CreatedAt:   &createdAt,
					UpdatedAt:   &updatedAt,
					Icon:        &icon,
					Description: &description,
				}
				item := mapProjectToItem(project)
				assert.Equal(t, "proj-456", item.ID.ValueString())
				assert.Equal(t, "Item Project", item.Name.ValueString())
				assert.Equal(t, "organization", item.Type.ValueString())
				assert.Equal(t, createdAt.String(), item.CreatedAt.ValueString())
				assert.Equal(t, updatedAt.String(), item.UpdatedAt.ValueString())
				assert.Equal(t, "🚀", item.Icon.ValueString())
				assert.Equal(t, "Item project description", item.Description.ValueString())

			case "map with minimal fields":
				project := &n8nsdk.Project{Name: "Minimal Item"}
				item := mapProjectToItem(project)
				assert.True(t, item.ID.IsNull())
				assert.Equal(t, "Minimal Item", item.Name.ValueString())
				assert.True(t, item.Type.IsNull())
				assert.True(t, item.CreatedAt.IsNull())
				assert.True(t, item.UpdatedAt.IsNull())
				assert.True(t, item.Icon.IsNull())
				assert.True(t, item.Description.IsNull())

			case "map empty string values":
				id := ""
				projectType := ""
				icon := ""
				description := ""
				project := &n8nsdk.Project{
					Id:          &id,
					Name:        "",
					Type:        &projectType,
					Icon:        &icon,
					Description: &description,
				}
				item := mapProjectToItem(project)
				assert.Equal(t, "", item.ID.ValueString())
				assert.Equal(t, "", item.Name.ValueString())
				assert.Equal(t, "", item.Type.ValueString())
				assert.Equal(t, "", item.Icon.ValueString())
				assert.Equal(t, "", item.Description.ValueString())

			case "map special characters":
				id := "proj-!@#$%^&*()"
				projectType := "type-with-üñíçödé"
				icon := "🌍🌎🌏"
				description := "Description with\nnewlines\tand\ttabs"
				project := &n8nsdk.Project{
					Id:          &id,
					Name:        "Name with 特殊字符",
					Type:        &projectType,
					Icon:        &icon,
					Description: &description,
				}
				item := mapProjectToItem(project)
				assert.Equal(t, id, item.ID.ValueString())
				assert.Equal(t, "Name with 特殊字符", item.Name.ValueString())
				assert.Equal(t, projectType, item.Type.ValueString())
				assert.Equal(t, icon, item.Icon.ValueString())
				assert.Equal(t, description, item.Description.ValueString())

			case "error case - validation checks":
				project := &n8nsdk.Project{Name: "Test"}
				item := mapProjectToItem(project)
				assert.NotNil(t, item)
			}
		})
	}
}

func TestHelpersConcurrency(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "concurrent findProjectByIDOrName", wantErr: false},
		{name: "concurrent mapProjectToItem", wantErr: false},
		{name: "error case - concurrent validation", wantErr: false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// NO t.Parallel() here - goroutines

			switch tt.name {
			case "concurrent findProjectByIDOrName":
				id1 := "proj-1"
				id2 := "proj-2"
				projects := []n8nsdk.Project{
					{Id: &id1, Name: "Project One"},
					{Id: &id2, Name: "Project Two"},
				}
				done := make(chan bool, 100)
				for i := 0; i < 100; i++ {
					go func(i int) {
						if i%2 == 0 {
							found, ok := findProjectByIDOrName(projects, types.StringValue("proj-1"), types.StringNull())
							assert.True(t, ok)
							assert.NotNil(t, found)
						} else {
							found, ok := findProjectByIDOrName(projects, types.StringNull(), types.StringValue("Project Two"))
							assert.True(t, ok)
							assert.NotNil(t, found)
						}
						done <- true
					}(i)
				}
				for i := 0; i < 100; i++ {
					<-done
				}

			case "concurrent mapProjectToItem":
				id := "proj-concurrent"
				projectType := "team"
				icon := "🔄"
				project := &n8nsdk.Project{
					Id:   &id,
					Name: "Concurrent Project",
					Type: &projectType,
					Icon: &icon,
				}
				done := make(chan bool, 100)
				for i := 0; i < 100; i++ {
					go func() {
						item := mapProjectToItem(project)
						assert.Equal(t, "proj-concurrent", item.ID.ValueString())
						assert.Equal(t, "Concurrent Project", item.Name.ValueString())
						assert.Equal(t, "team", item.Type.ValueString())
						assert.Equal(t, "🔄", item.Icon.ValueString())
						done <- true
					}()
				}
				for i := 0; i < 100; i++ {
					<-done
				}

			case "error case - concurrent validation":
				projects := []n8nsdk.Project{}
				done := make(chan bool, 10)
				for i := 0; i < 10; i++ {
					go func() {
						found, ok := findProjectByIDOrName(projects, types.StringNull(), types.StringNull())
						assert.False(t, ok)
						assert.Nil(t, found)
						done <- true
					}()
				}
				for i := 0; i < 10; i++ {
					<-done
				}
			}
		})
	}
}

func BenchmarkFindProjectByIDOrName(b *testing.B) {
	id1 := "proj-1"
	id2 := "proj-2"
	id3 := "proj-3"
	projects := []n8nsdk.Project{
		{Id: &id1, Name: "Project One"},
		{Id: &id2, Name: "Project Two"},
		{Id: &id3, Name: "Project Three"},
	}

	b.Run("find by ID", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = findProjectByIDOrName(projects, types.StringValue("proj-2"), types.StringNull())
		}
	})

	b.Run("find by name", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = findProjectByIDOrName(projects, types.StringNull(), types.StringValue("Project Two"))
		}
	})

	b.Run("not found", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = findProjectByIDOrName(projects, types.StringValue("proj-999"), types.StringNull())
		}
	})
}

func BenchmarkMapProjectToItem(b *testing.B) {
	id := "proj-bench"
	projectType := "team"
	createdAt := time.Now()
	updatedAt := time.Now()
	icon := "📊"
	description := "Benchmark description"

	project := &n8nsdk.Project{
		Id:          &id,
		Name:        "Benchmark Project",
		Type:        &projectType,
		CreatedAt:   &createdAt,
		UpdatedAt:   &updatedAt,
		Icon:        &icon,
		Description: &description,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mapProjectToItem(project)
	}
}
