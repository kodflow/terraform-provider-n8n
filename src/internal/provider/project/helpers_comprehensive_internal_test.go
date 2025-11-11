package project

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/project/models"
	"github.com/stretchr/testify/assert"
)

// TestFindProjectByIDOrName tests the findProjectByIDOrName function.
func TestFindProjectByIDOrName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		projects  []n8nsdk.Project
		id        types.String
		nameValue types.String
		wantFound bool
		wantID    string
	}{
		{
			name: "find by id",
			projects: []n8nsdk.Project{
				{
					Id:   func() *string { s := "proj-123"; return &s }(),
					Name: "Test Project",
				},
			},
			id:        types.StringValue("proj-123"),
			nameValue: types.StringNull(),
			wantFound: true,
			wantID:    "proj-123",
		},
		{
			name: "find by name",
			projects: []n8nsdk.Project{
				{
					Id:   func() *string { s := "proj-123"; return &s }(),
					Name: "Test Project",
				},
			},
			id:        types.StringNull(),
			nameValue: types.StringValue("Test Project"),
			wantFound: true,
			wantID:    "proj-123",
		},
		{
			name: "find by both id and name - id match",
			projects: []n8nsdk.Project{
				{
					Id:   func() *string { s := "proj-123"; return &s }(),
					Name: "Test Project",
				},
			},
			id:        types.StringValue("proj-123"),
			nameValue: types.StringValue("Test Project"),
			wantFound: true,
			wantID:    "proj-123",
		},
		{
			name: "not found",
			projects: []n8nsdk.Project{
				{
					Id:   func() *string { s := "proj-456"; return &s }(),
					Name: "Other Project",
				},
			},
			id:        types.StringValue("proj-123"),
			nameValue: types.StringNull(),
			wantFound: false,
		},
		{
			name: "nil project id",
			projects: []n8nsdk.Project{
				{
					Id:   nil,
					Name: "Test Project",
				},
			},
			id:        types.StringValue("proj-123"),
			nameValue: types.StringNull(),
			wantFound: false,
		},
		{
			name:      "empty projects list",
			projects:  []n8nsdk.Project{},
			id:        types.StringValue("proj-123"),
			nameValue: types.StringNull(),
			wantFound: false,
		},
		{
			name: "multiple projects - find second",
			projects: []n8nsdk.Project{
				{
					Id:   func() *string { s := "proj-123"; return &s }(),
					Name: "Project 1",
				},
				{
					Id:   func() *string { s := "proj-456"; return &s }(),
					Name: "Project 2",
				},
			},
			id:        types.StringValue("proj-456"),
			nameValue: types.StringNull(),
			wantFound: true,
			wantID:    "proj-456",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			project, found := findProjectByIDOrName(tt.projects, tt.id, tt.nameValue)

			assert.Equal(t, tt.wantFound, found, "findProjectByIDOrName returned unexpected found status")
			if tt.wantFound {
				assert.NotNil(t, project, "expected project to be non-nil")
				if project != nil && project.Id != nil {
					assert.Equal(t, tt.wantID, *project.Id, "project ID mismatch")
				}
			} else {
				assert.Nil(t, project, "expected project to be nil")
			}
		})
	}
}

// TestMapProjectToDataSourceModel tests the mapProjectToDataSourceModel function.
func TestMapProjectToDataSourceModel(t *testing.T) {
	t.Parallel()

	now := time.Now()

	tests := []struct {
		name    string
		project *n8nsdk.Project
		check   func(*testing.T, *models.DataSource)
	}{
		{
			name: "all fields populated",
			project: &n8nsdk.Project{
				Id:          func() *string { s := "proj-123"; return &s }(),
				Name:        "Test Project",
				Type:        func() *string { s := "team"; return &s }(),
				CreatedAt:   &now,
				UpdatedAt:   &now,
				Icon:        func() *string { s := "icon-test"; return &s }(),
				Description: func() *string { s := "test description"; return &s }(),
			},
			check: func(t *testing.T, data *models.DataSource) {
				t.Helper()
				assert.Equal(t, "proj-123", data.ID.ValueString())
				assert.Equal(t, "Test Project", data.Name.ValueString())
				assert.Equal(t, "team", data.Type.ValueString())
				assert.False(t, data.CreatedAt.IsNull())
				assert.False(t, data.UpdatedAt.IsNull())
				assert.Equal(t, "icon-test", data.Icon.ValueString())
				assert.Equal(t, "test description", data.Description.ValueString())
			},
		},
		{
			name: "minimal fields",
			project: &n8nsdk.Project{
				Id:   nil,
				Name: "Test Project",
			},
			check: func(t *testing.T, data *models.DataSource) {
				t.Helper()
				assert.True(t, data.ID.IsNull())
				assert.Equal(t, "Test Project", data.Name.ValueString())
				assert.True(t, data.Type.IsNull())
				assert.True(t, data.CreatedAt.IsNull())
				assert.True(t, data.UpdatedAt.IsNull())
				assert.True(t, data.Icon.IsNull())
				assert.True(t, data.Description.IsNull())
			},
		},
		{
			name: "nil icon sets null",
			project: &n8nsdk.Project{
				Id:   func() *string { s := "proj-123"; return &s }(),
				Name: "Test Project",
				Icon: nil,
			},
			check: func(t *testing.T, data *models.DataSource) {
				t.Helper()
				assert.True(t, data.Icon.IsNull(), "expected icon to be null")
			},
		},
		{
			name: "nil description sets null",
			project: &n8nsdk.Project{
				Id:          func() *string { s := "proj-123"; return &s }(),
				Name:        "Test Project",
				Description: nil,
			},
			check: func(t *testing.T, data *models.DataSource) {
				t.Helper()
				assert.True(t, data.Description.IsNull(), "expected description to be null")
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			data := &models.DataSource{}
			mapProjectToDataSourceModel(tt.project, data)
			tt.check(t, data)
		})
	}
}

// TestMapProjectToItem tests the mapProjectToItem function.
func TestMapProjectToItem(t *testing.T) {
	t.Parallel()

	now := time.Now()

	tests := []struct {
		name    string
		project *n8nsdk.Project
		check   func(*testing.T, models.Item)
	}{
		{
			name: "all fields populated",
			project: &n8nsdk.Project{
				Id:          func() *string { s := "proj-123"; return &s }(),
				Name:        "Test Project",
				Type:        func() *string { s := "team"; return &s }(),
				CreatedAt:   &now,
				UpdatedAt:   &now,
				Icon:        func() *string { s := "icon-test"; return &s }(),
				Description: func() *string { s := "test description"; return &s }(),
			},
			check: func(t *testing.T, item models.Item) {
				t.Helper()
				assert.Equal(t, "proj-123", item.ID.ValueString())
				assert.Equal(t, "Test Project", item.Name.ValueString())
				assert.Equal(t, "team", item.Type.ValueString())
				assert.False(t, item.CreatedAt.IsNull())
				assert.False(t, item.UpdatedAt.IsNull())
				assert.Equal(t, "icon-test", item.Icon.ValueString())
				assert.Equal(t, "test description", item.Description.ValueString())
			},
		},
		{
			name: "minimal fields",
			project: &n8nsdk.Project{
				Id:   nil,
				Name: "Test Project",
			},
			check: func(t *testing.T, item models.Item) {
				t.Helper()
				assert.True(t, item.ID.IsNull())
				assert.Equal(t, "Test Project", item.Name.ValueString())
				assert.True(t, item.Type.IsNull())
				assert.True(t, item.CreatedAt.IsNull())
				assert.True(t, item.UpdatedAt.IsNull())
				assert.True(t, item.Icon.IsNull())
				assert.True(t, item.Description.IsNull())
			},
		},
		{
			name: "nil icon not set",
			project: &n8nsdk.Project{
				Id:   func() *string { s := "proj-123"; return &s }(),
				Name: "Test Project",
				Icon: nil,
			},
			check: func(t *testing.T, item models.Item) {
				t.Helper()
				assert.True(t, item.Icon.IsNull())
			},
		},
		{
			name: "nil description not set",
			project: &n8nsdk.Project{
				Id:          func() *string { s := "proj-123"; return &s }(),
				Name:        "Test Project",
				Description: nil,
			},
			check: func(t *testing.T, item models.Item) {
				t.Helper()
				assert.True(t, item.Description.IsNull())
			},
		},
		{
			name: "nil type not set",
			project: &n8nsdk.Project{
				Id:   func() *string { s := "proj-123"; return &s }(),
				Name: "Test Project",
				Type: nil,
			},
			check: func(t *testing.T, item models.Item) {
				t.Helper()
				assert.True(t, item.Type.IsNull())
			},
		},
		{
			name: "nil timestamps not set",
			project: &n8nsdk.Project{
				Id:        func() *string { s := "proj-123"; return &s }(),
				Name:      "Test Project",
				CreatedAt: nil,
				UpdatedAt: nil,
			},
			check: func(t *testing.T, item models.Item) {
				t.Helper()
				assert.True(t, item.CreatedAt.IsNull())
				assert.True(t, item.UpdatedAt.IsNull())
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			item := mapProjectToItem(tt.project)
			tt.check(t, item)
		})
	}
}
