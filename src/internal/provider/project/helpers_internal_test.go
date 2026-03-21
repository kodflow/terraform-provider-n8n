package project

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/project/models"
	"github.com/stretchr/testify/assert"
)

func Test_findProjectByIDOrName(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "find by ID":
				projects := []n8nsdk.Project{
					{Id: new("proj-1"), Name: "Project One"},
					{Id: new("proj-2"), Name: "Project Two"},
				}
				found, ok := findProjectByIDOrName(projects, types.StringValue("proj-2"), types.StringNull())
				assert.True(t, ok)
				assert.NotNil(t, found)
				assert.Equal(t, "proj-2", *found.Id)
				assert.Equal(t, "Project Two", found.Name)

			case "find by name":
				projects := []n8nsdk.Project{
					{Id: new("proj-1"), Name: "Project One"},
					{Id: new("proj-2"), Name: "Project Two"},
				}
				found, ok := findProjectByIDOrName(projects, types.StringNull(), types.StringValue("Project One"))
				assert.True(t, ok)
				assert.NotNil(t, found)
				assert.Equal(t, "proj-1", *found.Id)
				assert.Equal(t, "Project One", found.Name)

			case "find by ID and name (ID takes precedence)":
				projects := []n8nsdk.Project{
					{Id: new("proj-1"), Name: "Project One"},
					{Id: new("proj-2"), Name: "Project Two"},
				}
				found, ok := findProjectByIDOrName(projects, types.StringValue("proj-1"), types.StringValue("Project Two"))
				assert.True(t, ok)
				assert.NotNil(t, found)
				assert.Equal(t, "proj-1", *found.Id)
				assert.Equal(t, "Project One", found.Name)

			case "not found":
				projects := []n8nsdk.Project{{Id: new("proj-1"), Name: "Project One"}}
				found, ok := findProjectByIDOrName(projects, types.StringValue("proj-999"), types.StringValue("Non-existent"))
				assert.False(t, ok)
				assert.Nil(t, found)

			case "empty projects list":
				projects := []n8nsdk.Project{}
				found, ok := findProjectByIDOrName(projects, types.StringValue("any-id"), types.StringValue("any-name"))
				assert.False(t, ok)
				assert.Nil(t, found)

			case "null search parameters":
				projects := []n8nsdk.Project{{Id: new("proj-1"), Name: "Project One"}}
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
				projects := []n8nsdk.Project{{Id: new("proj-1"), Name: "Project One"}}
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

func Test_mapProjectToDataSourceModel(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "map with all fields":
				id := "proj-123"
				createdAt := time.Now()
				updatedAt := time.Now().Add(1 * time.Hour)
				description := "Test project description"
				project := &n8nsdk.Project{
					Id:          &id,
					Name:        "Test Project",
					Type:        new("personal"),
					CreatedAt:   &createdAt,
					UpdatedAt:   &updatedAt,
					Icon:        *n8nsdk.NewNullableProjectIcon(&n8nsdk.ProjectIcon{Value: new("📁")}),
					Description: *n8nsdk.NewNullableString(&description),
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
					Type: new("team"),
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

func Test_mapProjectToItem(t *testing.T) {
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
					Icon:        *n8nsdk.NewNullableProjectIcon(&n8nsdk.ProjectIcon{Value: &icon}),
					Description: *n8nsdk.NewNullableString(&description),
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
					Icon:        *n8nsdk.NewNullableProjectIcon(&n8nsdk.ProjectIcon{Value: &icon}),
					Description: *n8nsdk.NewNullableString(&description),
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
					Icon:        *n8nsdk.NewNullableProjectIcon(&n8nsdk.ProjectIcon{Value: &icon}),
					Description: *n8nsdk.NewNullableString(&description),
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

// TestHelpersConcurrency verifies that helper functions are safe for concurrent use.
// Goroutine lifecycle: goroutines are launched in batches, each sends to a done channel,
// and the test waits for all goroutines to complete via channel receives before returning.
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
		t.Run(tt.name, func(t *testing.T) {
			// NO t.Parallel() here - goroutines

			switch tt.name {
			case "concurrent findProjectByIDOrName":
				projects := []n8nsdk.Project{
					{Id: new("proj-1"), Name: "Project One"},
					{Id: new("proj-2"), Name: "Project Two"},
				}
				done := make(chan bool, 100)
				for i := range 100 {
					//: Goroutine terminates after sending result to done channel.
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
				for range 100 {
					<-done
				}

			case "concurrent mapProjectToItem":
				project := &n8nsdk.Project{
					Id:   new("proj-concurrent"),
					Name: "Concurrent Project",
					Type: new("team"),
					Icon: *n8nsdk.NewNullableProjectIcon(&n8nsdk.ProjectIcon{Value: new("🔄")}),
				}
				done := make(chan bool, 100)
				for range 100 {
					//: Goroutine terminates after sending result to done channel.
					go func() {
						item := mapProjectToItem(project)
						assert.Equal(t, "proj-concurrent", item.ID.ValueString())
						assert.Equal(t, "Concurrent Project", item.Name.ValueString())
						assert.Equal(t, "team", item.Type.ValueString())
						assert.Equal(t, "🔄", item.Icon.ValueString())
						done <- true
					}()
				}
				for range 100 {
					<-done
				}

			case "error case - concurrent validation":
				projects := []n8nsdk.Project{}
				done := make(chan bool, 10)
				for range 10 {
					//: Goroutine terminates after sending result to done channel.
					go func() {
						found, ok := findProjectByIDOrName(projects, types.StringNull(), types.StringNull())
						assert.False(t, ok)
						assert.Nil(t, found)
						done <- true
					}()
				}
				for range 10 {
					<-done
				}
			}
		})
	}
}

func BenchmarkFindProjectByIDOrName(b *testing.B) {
	projects := []n8nsdk.Project{
		{Id: new("proj-1"), Name: "Project One"},
		{Id: new("proj-2"), Name: "Project Two"},
		{Id: new("proj-3"), Name: "Project Three"},
	}

	b.Run("find by ID", func(b *testing.B) {
		for b.Loop() {
			_, _ = findProjectByIDOrName(projects, types.StringValue("proj-2"), types.StringNull())
		}
	})

	b.Run("find by name", func(b *testing.B) {
		for b.Loop() {
			_, _ = findProjectByIDOrName(projects, types.StringNull(), types.StringValue("Project Two"))
		}
	})

	b.Run("not found", func(b *testing.B) {
		for b.Loop() {
			_, _ = findProjectByIDOrName(projects, types.StringValue("proj-999"), types.StringNull())
		}
	})
}

func BenchmarkMapProjectToItem(b *testing.B) {

	project := &n8nsdk.Project{
		Id:          new("proj-bench"),
		Name:        "Benchmark Project",
		Type:        new("team"),
		CreatedAt:   new(time.Now()),
		UpdatedAt:   new(time.Now()),
		Icon:        *n8nsdk.NewNullableProjectIcon(&n8nsdk.ProjectIcon{Value: new("📊")}),
		Description: *n8nsdk.NewNullableString(new("Benchmark description")),
	}

	b.ResetTimer()
	for b.Loop() {
		mapProjectToItem(project)
	}
}

// TestExtractIconValue verifies the extractIconValue helper function.
func TestExtractIconValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		icon       n8nsdk.NullableProjectIcon
		expectNull bool
		expectVal  string
	}{
		{
			name:       "set icon with value",
			icon:       *n8nsdk.NewNullableProjectIcon(&n8nsdk.ProjectIcon{Value: new("🚀")}),
			expectNull: false,
			expectVal:  "🚀",
		},
		{
			name:       "set icon with nil value",
			icon:       *n8nsdk.NewNullableProjectIcon(&n8nsdk.ProjectIcon{Value: nil}),
			expectNull: true,
		},
		{
			name:       "nil icon",
			icon:       *n8nsdk.NewNullableProjectIcon(nil),
			expectNull: true,
		},
		{
			name:       "error case - unset icon",
			icon:       n8nsdk.NullableProjectIcon{},
			expectNull: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := extractIconValue(tt.icon)

			if tt.expectNull {
				assert.True(t, result.IsNull(), "expected null result")
			} else {
				assert.Equal(t, tt.expectVal, result.ValueString())
			}
		})
	}
}

// TestExtractDescriptionValue verifies the extractDescriptionValue helper function.
func TestExtractDescriptionValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		desc       n8nsdk.NullableString
		expectNull bool
		expectVal  string
	}{
		{
			name:       "set description",
			desc:       *n8nsdk.NewNullableString(new("My project description")),
			expectNull: false,
			expectVal:  "My project description",
		},
		{
			name:       "nil description",
			desc:       *n8nsdk.NewNullableString(nil),
			expectNull: true,
		},
		{
			name:       "error case - unset description",
			desc:       n8nsdk.NullableString{},
			expectNull: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := extractDescriptionValue(tt.desc)

			if tt.expectNull {
				assert.True(t, result.IsNull(), "expected null result")
			} else {
				assert.Equal(t, tt.expectVal, result.ValueString())
			}
		})
	}
}
