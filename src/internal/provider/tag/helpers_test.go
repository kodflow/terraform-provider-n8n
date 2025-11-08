package tag

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/tag/models"
	"github.com/stretchr/testify/assert"
)

func TestMapTagToDataSourceModel(t *testing.T) {
	t.Run("map with all fields populated", func(t *testing.T) {
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
	})

	t.Run("map with nil id", func(t *testing.T) {
		name := "Tag without ID"
		tag := &n8nsdk.Tag{
			Id:   nil,
			Name: name,
		}

		data := &models.DataSource{}
		mapTagToDataSourceModel(tag, data)

		assert.True(t, data.ID.IsNull())
		assert.Equal(t, "Tag without ID", data.Name.ValueString())
	})

	t.Run("map with nil timestamps", func(t *testing.T) {
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
	})

	t.Run("map with empty name", func(t *testing.T) {
		id := "tag-789"
		tag := &n8nsdk.Tag{
			Id:   &id,
			Name: "",
		}

		data := &models.DataSource{}
		mapTagToDataSourceModel(tag, data)

		assert.Equal(t, "tag-789", data.ID.ValueString())
		assert.Equal(t, "", data.Name.ValueString())
	})

	t.Run("map with special characters in name", func(t *testing.T) {
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
	})

	t.Run("map with long name", func(t *testing.T) {
		id := "tag-long"
		longName := "This is a very long tag name that contains many characters and might be used to test edge cases in the system for handling long strings"
		tag := &n8nsdk.Tag{
			Id:   &id,
			Name: longName,
		}

		data := &models.DataSource{}
		mapTagToDataSourceModel(tag, data)

		assert.Equal(t, longName, data.Name.ValueString())
	})

	t.Run("map with zero time", func(t *testing.T) {
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
	})

	t.Run("map with different timestamp values", func(t *testing.T) {
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
	})

	t.Run("map preserves existing data fields", func(t *testing.T) {
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
	})

	t.Run("map multiple times", func(t *testing.T) {
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
	})
}

func TestFindTagByName(t *testing.T) {
	t.Run("find existing tag", func(t *testing.T) {
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
	})

	t.Run("tag not found", func(t *testing.T) {
		id1 := "tag-1"
		tags := []n8nsdk.Tag{
			{Id: &id1, Name: "Production"},
		}

		tag, found := findTagByName(tags, "NonExistent")

		assert.False(t, found)
		assert.Nil(t, tag)
	})

	t.Run("empty tag list", func(t *testing.T) {
		tags := []n8nsdk.Tag{}

		tag, found := findTagByName(tags, "AnyTag")

		assert.False(t, found)
		assert.Nil(t, tag)
	})

	t.Run("nil tag list", func(t *testing.T) {
		var tags []n8nsdk.Tag

		tag, found := findTagByName(tags, "AnyTag")

		assert.False(t, found)
		assert.Nil(t, tag)
	})

	t.Run("find first tag", func(t *testing.T) {
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
	})

	t.Run("find last tag", func(t *testing.T) {
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
	})

	t.Run("case sensitive search", func(t *testing.T) {
		id := "tag-1"
		tags := []n8nsdk.Tag{
			{Id: &id, Name: "Production"},
		}

		// Search with different case
		tag, found := findTagByName(tags, "production")

		assert.False(t, found)
		assert.Nil(t, tag)
	})

	t.Run("exact match required", func(t *testing.T) {
		id := "tag-1"
		tags := []n8nsdk.Tag{
			{Id: &id, Name: "Production Server"},
		}

		// Partial match should not work
		tag, found := findTagByName(tags, "Production")

		assert.False(t, found)
		assert.Nil(t, tag)
	})

	t.Run("empty name search", func(t *testing.T) {
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
	})

	t.Run("whitespace in names", func(t *testing.T) {
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
	})

	t.Run("special characters in names", func(t *testing.T) {
		id := "tag-special"
		tags := []n8nsdk.Tag{
			{Id: &id, Name: "Tag@#$%^&*()"},
		}

		tag, found := findTagByName(tags, "Tag@#$%^&*()")

		assert.True(t, found)
		assert.NotNil(t, tag)
		assert.Equal(t, "Tag@#$%^&*()", tag.Name)
	})

	t.Run("unicode in names", func(t *testing.T) {
		id := "tag-unicode"
		tags := []n8nsdk.Tag{
			{Id: &id, Name: "标签测试"},
		}

		tag, found := findTagByName(tags, "标签测试")

		assert.True(t, found)
		assert.NotNil(t, tag)
		assert.Equal(t, "标签测试", tag.Name)
	})

	t.Run("duplicate names returns first", func(t *testing.T) {
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
	})

	t.Run("tag with nil id", func(t *testing.T) {
		tags := []n8nsdk.Tag{
			{Id: nil, Name: "NilID"},
		}

		tag, found := findTagByName(tags, "NilID")

		assert.True(t, found)
		assert.NotNil(t, tag)
		assert.Nil(t, tag.Id)
		assert.Equal(t, "NilID", tag.Name)
	})

	t.Run("large tag list", func(t *testing.T) {
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
	})
}

func TestMapTagToDataSourceModelConcurrency(t *testing.T) {
	t.Run("concurrent mapping", func(t *testing.T) {
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
	})
}

func TestFindTagByNameConcurrency(t *testing.T) {
	t.Run("concurrent searches", func(t *testing.T) {
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
	})
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
