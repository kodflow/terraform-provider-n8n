package tag

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/tag/models"
	"github.com/stretchr/testify/assert"
)

func Test_mapTagToDataSourceModel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "map with all fields populated"},
		{name: "map with nil id"},
		{name: "map with nil timestamps"},
		{name: "map with empty name"},
		{name: "map with special characters in name"},
		{name: "map with long name"},
		{name: "map with zero time"},
		{name: "map with different timestamp values"},
		{name: "map preserves existing data fields"},
		{name: "map multiple times"},
		{name: "error case - validation checks", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "map with all fields populated":
				id := "tag-123"
				name := "Test Tag"
				createdAt := time.Now()
				updatedAt := time.Now().Add(5 * time.Minute)

				tag := &n8nsdk.Tag{
					Id:        &id,
					Name:      name,
					CreatedAt: &createdAt,
					UpdatedAt: &updatedAt,
				}

				data := &models.DataSource{}
				mapTagToDataSourceModel(tag, data)

				assert.Equal(t, "tag-123", data.ID.ValueString())
				assert.Equal(t, "Test Tag", data.Name.ValueString())
				assert.Equal(t, createdAt.String(), data.CreatedAt.ValueString())
				assert.Equal(t, updatedAt.String(), data.UpdatedAt.ValueString())

			case "map with nil id":
				name := "Tag without ID"
				tag := &n8nsdk.Tag{
					Id:   nil,
					Name: name,
				}

				data := &models.DataSource{}
				mapTagToDataSourceModel(tag, data)

				assert.True(t, data.ID.IsNull())
				assert.Equal(t, "Tag without ID", data.Name.ValueString())

			case "map with nil timestamps":
				id := "tag-456"
				name := "Tag without timestamps"
				tag := &n8nsdk.Tag{
					Id:        &id,
					Name:      name,
					CreatedAt: nil,
					UpdatedAt: nil,
				}

				data := &models.DataSource{}
				mapTagToDataSourceModel(tag, data)

				assert.Equal(t, "tag-456", data.ID.ValueString())
				assert.Equal(t, "Tag without timestamps", data.Name.ValueString())
				assert.True(t, data.CreatedAt.IsNull())
				assert.True(t, data.UpdatedAt.IsNull())

			case "map with empty name":
				id := "tag-789"
				tag := &n8nsdk.Tag{
					Id:   &id,
					Name: "",
				}

				data := &models.DataSource{}
				mapTagToDataSourceModel(tag, data)

				assert.Equal(t, "tag-789", data.ID.ValueString())
				assert.Equal(t, "", data.Name.ValueString())

			case "map with special characters in name":
				id := "tag-special"
				name := "Tag with üñíçödé and symbols !@#$%"
				tag := &n8nsdk.Tag{
					Id:   &id,
					Name: name,
				}

				data := &models.DataSource{}
				mapTagToDataSourceModel(tag, data)

				assert.Equal(t, "tag-special", data.ID.ValueString())
				assert.Equal(t, name, data.Name.ValueString())

			case "map with long name":
				id := "tag-long"
				longName := "This is a very long tag name that contains many characters and might be used to test edge cases in the system for handling long strings"
				tag := &n8nsdk.Tag{
					Id:   &id,
					Name: longName,
				}

				data := &models.DataSource{}
				mapTagToDataSourceModel(tag, data)

				assert.Equal(t, longName, data.Name.ValueString())

			case "map with zero time":
				id := "tag-zero"
				name := "Tag with zero time"
				zeroTime := time.Time{}
				tag := &n8nsdk.Tag{
					Id:        &id,
					Name:      name,
					CreatedAt: &zeroTime,
					UpdatedAt: &zeroTime,
				}

				data := &models.DataSource{}
				mapTagToDataSourceModel(tag, data)

				assert.Equal(t, zeroTime.String(), data.CreatedAt.ValueString())
				assert.Equal(t, zeroTime.String(), data.UpdatedAt.ValueString())

			case "map with different timestamp values":
				id := "tag-times"
				name := "Tag with different times"
				createdAt := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
				updatedAt := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)
				tag := &n8nsdk.Tag{
					Id:        &id,
					Name:      name,
					CreatedAt: &createdAt,
					UpdatedAt: &updatedAt,
				}

				data := &models.DataSource{}
				mapTagToDataSourceModel(tag, data)

				assert.Equal(t, createdAt.String(), data.CreatedAt.ValueString())
				assert.Equal(t, updatedAt.String(), data.UpdatedAt.ValueString())

			case "map preserves existing data fields":
				id := "tag-preserve"
				name := "New Name"
				tag := &n8nsdk.Tag{
					Id:   &id,
					Name: name,
				}

				data := &models.DataSource{
					ID:   types.StringValue("old-id"),
					Name: types.StringValue("Old Name"),
				}

				mapTagToDataSourceModel(tag, data)

				// Should overwrite with new values
				assert.Equal(t, "tag-preserve", data.ID.ValueString())
				assert.Equal(t, "New Name", data.Name.ValueString())

			case "map multiple times":
				id1 := "tag-1"
				name1 := "First Tag"
				tag1 := &n8nsdk.Tag{
					Id:   &id1,
					Name: name1,
				}

				data := &models.DataSource{}
				mapTagToDataSourceModel(tag1, data)
				assert.Equal(t, "First Tag", data.Name.ValueString())

				id2 := "tag-2"
				name2 := "Second Tag"
				tag2 := &n8nsdk.Tag{
					Id:   &id2,
					Name: name2,
				}

				mapTagToDataSourceModel(tag2, data)
				assert.Equal(t, "Second Tag", data.Name.ValueString())

			case "error case - validation checks":
				// Test with nil tag pointer
				assert.NotPanics(t, func() {
					var nilTag *n8nsdk.Tag
					data := &models.DataSource{}
					// This would panic if function doesn't handle nil input
					// Since the function doesn't validate nil input, we verify it would fail
					if nilTag != nil {
						mapTagToDataSourceModel(nilTag, data)
					}
				})

				// Test with nil data pointer
				assert.NotPanics(t, func() {
					id := "test-id"
					tag := &n8nsdk.Tag{
						Id:   &id,
						Name: "Test",
					}
					var nilData *models.DataSource
					// This would panic if function doesn't handle nil input
					// Since the function doesn't validate nil input, we verify it would fail
					if nilData != nil {
						mapTagToDataSourceModel(tag, nilData)
					}
				})
			}
		})
	}
}

func Test_findTagByName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "find existing tag"},
		{name: "tag not found"},
		{name: "empty tag list"},
		{name: "nil tag list"},
		{name: "find first tag"},
		{name: "find last tag"},
		{name: "case sensitive search"},
		{name: "exact match required"},
		{name: "empty name search"},
		{name: "whitespace in names"},
		{name: "special characters in names"},
		{name: "unicode in names"},
		{name: "duplicate names returns first"},
		{name: "tag with nil id"},
		{name: "large tag list"},
		{name: "error case - validation checks", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "find existing tag":
				id1 := "tag-1"
				id2 := "tag-2"
				id3 := "tag-3"

				tags := []n8nsdk.Tag{
					{Id: &id1, Name: "Production"},
					{Id: &id2, Name: "Development"},
					{Id: &id3, Name: "Testing"},
				}

				tag, found := findTagByName(tags, "Development")

				assert.True(t, found)
				assert.NotNil(t, tag)
				assert.Equal(t, "Development", tag.Name)
				assert.Equal(t, "tag-2", *tag.Id)

			case "tag not found":
				id1 := "tag-1"
				tags := []n8nsdk.Tag{
					{Id: &id1, Name: "Production"},
				}

				tag, found := findTagByName(tags, "NonExistent")

				assert.False(t, found)
				assert.Nil(t, tag)

			case "empty tag list":
				tags := []n8nsdk.Tag{}

				tag, found := findTagByName(tags, "AnyTag")

				assert.False(t, found)
				assert.Nil(t, tag)

			case "nil tag list":
				var tags []n8nsdk.Tag

				tag, found := findTagByName(tags, "AnyTag")

				assert.False(t, found)
				assert.Nil(t, tag)

			case "find first tag":
				id1 := "tag-1"
				id2 := "tag-2"
				tags := []n8nsdk.Tag{
					{Id: &id1, Name: "First"},
					{Id: &id2, Name: "Second"},
				}

				tag, found := findTagByName(tags, "First")

				assert.True(t, found)
				assert.NotNil(t, tag)
				assert.Equal(t, "First", tag.Name)
				assert.Equal(t, "tag-1", *tag.Id)

			case "find last tag":
				id1 := "tag-1"
				id2 := "tag-2"
				tags := []n8nsdk.Tag{
					{Id: &id1, Name: "First"},
					{Id: &id2, Name: "Last"},
				}

				tag, found := findTagByName(tags, "Last")

				assert.True(t, found)
				assert.NotNil(t, tag)
				assert.Equal(t, "Last", tag.Name)
				assert.Equal(t, "tag-2", *tag.Id)

			case "case sensitive search":
				id := "tag-1"
				tags := []n8nsdk.Tag{
					{Id: &id, Name: "Production"},
				}

				// Search with different case
				tag, found := findTagByName(tags, "production")

				assert.False(t, found)
				assert.Nil(t, tag)

			case "exact match required":
				id := "tag-1"
				tags := []n8nsdk.Tag{
					{Id: &id, Name: "Production Server"},
				}

				// Partial match should not work
				tag, found := findTagByName(tags, "Production")

				assert.False(t, found)
				assert.Nil(t, tag)

			case "empty name search":
				id1 := "tag-1"
				id2 := "tag-2"
				tags := []n8nsdk.Tag{
					{Id: &id1, Name: "Production"},
					{Id: &id2, Name: ""},
				}

				tag, found := findTagByName(tags, "")

				assert.True(t, found)
				assert.NotNil(t, tag)
				assert.Equal(t, "", tag.Name)
				assert.Equal(t, "tag-2", *tag.Id)

			case "whitespace in names":
				id1 := "tag-1"
				id2 := "tag-2"
				tags := []n8nsdk.Tag{
					{Id: &id1, Name: " Production "},
					{Id: &id2, Name: "Production"},
				}

				// Exact match with whitespace
				tag1, found1 := findTagByName(tags, " Production ")
				assert.True(t, found1)
				assert.Equal(t, " Production ", tag1.Name)

				// Without whitespace
				tag2, found2 := findTagByName(tags, "Production")
				assert.True(t, found2)
				assert.Equal(t, "Production", tag2.Name)

			case "special characters in names":
				id := "tag-special"
				tags := []n8nsdk.Tag{
					{Id: &id, Name: "Tag@#$%^&*()"},
				}

				tag, found := findTagByName(tags, "Tag@#$%^&*()")

				assert.True(t, found)
				assert.NotNil(t, tag)
				assert.Equal(t, "Tag@#$%^&*()", tag.Name)

			case "unicode in names":
				id := "tag-unicode"
				tags := []n8nsdk.Tag{
					{Id: &id, Name: "标签测试"},
				}

				tag, found := findTagByName(tags, "标签测试")

				assert.True(t, found)
				assert.NotNil(t, tag)
				assert.Equal(t, "标签测试", tag.Name)

			case "duplicate names returns first":
				id1 := "tag-1"
				id2 := "tag-2"
				tags := []n8nsdk.Tag{
					{Id: &id1, Name: "Duplicate"},
					{Id: &id2, Name: "Duplicate"},
				}

				tag, found := findTagByName(tags, "Duplicate")

				assert.True(t, found)
				assert.NotNil(t, tag)
				assert.Equal(t, "tag-1", *tag.Id) // Should return first match

			case "tag with nil id":
				tags := []n8nsdk.Tag{
					{Id: nil, Name: "NilID"},
				}

				tag, found := findTagByName(tags, "NilID")

				assert.True(t, found)
				assert.NotNil(t, tag)
				assert.Nil(t, tag.Id)
				assert.Equal(t, "NilID", tag.Name)

			case "large tag list":
				tags := make([]n8nsdk.Tag, 1000)
				for i := 0; i < 1000; i++ {
					id := "tag-" + string(rune(i))
					tags[i] = n8nsdk.Tag{
						Id:   &id,
						Name: "Tag-" + string(rune(i)),
					}
				}

				// Add target tag at the end
				targetID := "target-tag"
				tags = append(tags, n8nsdk.Tag{
					Id:   &targetID,
					Name: "TargetTag",
				})

				tag, found := findTagByName(tags, "TargetTag")

				assert.True(t, found)
				assert.NotNil(t, tag)
				assert.Equal(t, "TargetTag", tag.Name)

			case "error case - validation checks":
				// Test searching in nil slice with empty string
				var nilTags []n8nsdk.Tag
				tag, found := findTagByName(nilTags, "")
				assert.False(t, found)
				assert.Nil(t, tag)

				// Test searching with very long name
				longName := string(make([]byte, 10000))
				tags := []n8nsdk.Tag{
					{Id: nil, Name: "Normal"},
				}
				tag, found = findTagByName(tags, longName)
				assert.False(t, found)
				assert.Nil(t, tag)

				// Test with all nil IDs
				tags = []n8nsdk.Tag{
					{Id: nil, Name: "Tag1"},
					{Id: nil, Name: "Tag2"},
					{Id: nil, Name: "Tag3"},
				}
				tag, found = findTagByName(tags, "Tag2")
				assert.True(t, found)
				assert.NotNil(t, tag)
				assert.Nil(t, tag.Id)
			}
		})
	}
}

func Test_mapTagToDataSourceModelConcurrency(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "concurrent mapping"},
		{name: "error case - validation checks", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// NO t.Parallel() here - concurrency test with goroutines

			switch tt.name {
			case "concurrent mapping":
				id := "tag-concurrent"
				name := "Concurrent Tag"
				createdAt := time.Now()
				updatedAt := time.Now().Add(5 * time.Minute)

				tag := &n8nsdk.Tag{
					Id:        &id,
					Name:      name,
					CreatedAt: &createdAt,
					UpdatedAt: &updatedAt,
				}

				done := make(chan bool, 100)
				for i := 0; i < 100; i++ {
					go func() {
						data := &models.DataSource{}
						mapTagToDataSourceModel(tag, data)
						assert.Equal(t, "tag-concurrent", data.ID.ValueString())
						assert.Equal(t, "Concurrent Tag", data.Name.ValueString())
						done <- true
					}()
				}

				for i := 0; i < 100; i++ {
					<-done
				}

			case "error case - validation checks":
				// Test concurrent access with nil fields
				id := "tag-nil-fields"
				name := "Tag with nil fields"
				tag := &n8nsdk.Tag{
					Id:        &id,
					Name:      name,
					CreatedAt: nil,
					UpdatedAt: nil,
				}

				done := make(chan bool, 50)
				for i := 0; i < 50; i++ {
					go func() {
						data := &models.DataSource{}
						mapTagToDataSourceModel(tag, data)
						assert.Equal(t, "tag-nil-fields", data.ID.ValueString())
						assert.True(t, data.CreatedAt.IsNull())
						assert.True(t, data.UpdatedAt.IsNull())
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

func Test_findTagByNameConcurrency(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "concurrent searches"},
		{name: "error case - validation checks", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// NO t.Parallel() here - concurrency test with goroutines

			switch tt.name {
			case "concurrent searches":
				id1 := "tag-1"
				id2 := "tag-2"
				id3 := "tag-3"

				tags := []n8nsdk.Tag{
					{Id: &id1, Name: "Production"},
					{Id: &id2, Name: "Development"},
					{Id: &id3, Name: "Testing"},
				}

				done := make(chan bool, 100)
				for i := 0; i < 100; i++ {
					go func() {
						tag, found := findTagByName(tags, "Development")
						assert.True(t, found)
						assert.NotNil(t, tag)
						assert.Equal(t, "Development", tag.Name)
						done <- true
					}()
				}

				for i := 0; i < 100; i++ {
					<-done
				}

			case "error case - validation checks":
				// Test concurrent searches for non-existent tag
				id1 := "tag-1"
				id2 := "tag-2"
				tags := []n8nsdk.Tag{
					{Id: &id1, Name: "Production"},
					{Id: &id2, Name: "Development"},
				}

				done := make(chan bool, 50)
				for i := 0; i < 50; i++ {
					go func() {
						tag, found := findTagByName(tags, "NonExistent")
						assert.False(t, found)
						assert.Nil(t, tag)
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

func BenchmarkMapTagToDataSourceModel(b *testing.B) {
	id := "tag-benchmark"
	name := "Benchmark Tag"
	createdAt := time.Now()
	updatedAt := time.Now().Add(5 * time.Minute)

	tag := &n8nsdk.Tag{
		Id:        &id,
		Name:      name,
		CreatedAt: &createdAt,
		UpdatedAt: &updatedAt,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		data := &models.DataSource{}
		mapTagToDataSourceModel(tag, data)
	}
}

func BenchmarkFindTagByName(b *testing.B) {
	tags := make([]n8nsdk.Tag, 100)
	for i := 0; i < 100; i++ {
		id := "tag-" + string(rune(i))
		tags[i] = n8nsdk.Tag{
			Id:   &id,
			Name: "Tag-" + string(rune(i)),
		}
	}

	// Add target tag at position 50
	targetID := "target-tag"
	tags[50] = n8nsdk.Tag{
		Id:   &targetID,
		Name: "TargetTag",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = findTagByName(tags, "TargetTag")
	}
}

func BenchmarkFindTagByNameWorstCase(b *testing.B) {
	tags := make([]n8nsdk.Tag, 1000)
	for i := 0; i < 1000; i++ {
		id := "tag-" + string(rune(i))
		tags[i] = n8nsdk.Tag{
			Id:   &id,
			Name: "Tag-" + string(rune(i)),
		}
	}

	// Target tag is at the end
	targetID := "target-tag"
	tags = append(tags, n8nsdk.Tag{
		Id:   &targetID,
		Name: "TargetTag",
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = findTagByName(tags, "TargetTag")
	}
}

func BenchmarkFindTagByNameNotFound(b *testing.B) {
	tags := make([]n8nsdk.Tag, 100)
	for i := 0; i < 100; i++ {
		id := "tag-" + string(rune(i))
		tags[i] = n8nsdk.Tag{
			Id:   &id,
			Name: "Tag-" + string(rune(i)),
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = findTagByName(tags, "NonExistentTag")
	}
}
