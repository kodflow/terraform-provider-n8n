// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// Package tag provides white-box tests for the TagsDataSource type.
package tag

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/client"
	"github.com/stretchr/testify/assert"
)

// TestTagsDataSource_Metadata_Internal verifies the Metadata method sets the correct type name.
func TestTagsDataSource_Metadata_Internal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		providerTypeName string
		expectedTypeName string
	}{
		{
			name:             "n8n provider type",
			providerTypeName: "n8n",
			expectedTypeName: "n8n_tags",
		},
		{
			name:             "empty provider type",
			providerTypeName: "",
			expectedTypeName: "_tags",
		},
		{
			name:             "error case - custom provider type",
			providerTypeName: "custom",
			expectedTypeName: "custom_tags",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			d := NewTagsDataSource()
			resp := &datasource.MetadataResponse{}

			d.Metadata(t.Context(), datasource.MetadataRequest{
				ProviderTypeName: tt.providerTypeName,
			}, resp)

			assert.Equal(t, tt.expectedTypeName, resp.TypeName)
		})
	}
}

// TestTagsDataSource_Schema_Internal verifies the Schema method populates the schema correctly.
func TestTagsDataSource_Schema_Internal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
	}{
		{name: "schema has tags list attribute"},
		{name: "error case - schema is not nil"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			d := NewTagsDataSource()
			resp := &datasource.SchemaResponse{}

			d.Schema(t.Context(), datasource.SchemaRequest{}, resp)

			assert.NotNil(t, resp.Schema)
			assert.NotEmpty(t, resp.Schema.Attributes)
			_, ok := resp.Schema.Attributes["tags"]
			assert.True(t, ok, "expected 'tags' attribute in schema")
		})
	}
}

// TestTagsDataSource_Configure_Internal verifies Configure assigns client or adds errors.
func TestTagsDataSource_Configure_Internal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		providerData any
		wantError    bool
	}{
		{
			name:         "nil provider data",
			providerData: nil,
			wantError:    false,
		},
		{
			name:         "valid N8nClient",
			providerData: &client.N8nClient{},
			wantError:    false,
		},
		{
			name:         "error case - string type",
			providerData: "invalid-type",
			wantError:    true,
		},
		{
			name:         "error case - integer type",
			providerData: 99,
			wantError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			d := NewTagsDataSource()
			resp := &datasource.ConfigureResponse{}

			d.Configure(t.Context(), datasource.ConfigureRequest{
				ProviderData: tt.providerData,
			}, resp)

			if tt.wantError {
				assert.True(t, resp.Diagnostics.HasError())
			} else {
				assert.False(t, resp.Diagnostics.HasError())
			}
		})
	}
}

// TestNewTagsDataSource_Internal verifies the constructor returns a non-nil instance.
func TestNewTagsDataSource_Internal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
	}{
		{name: "returns non-nil instance"},
		{name: "error case - multiple calls are independent"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "returns non-nil instance":
				ds := NewTagsDataSource()
				assert.NotNil(t, ds)

			case "error case - multiple calls are independent":
				ds1 := NewTagsDataSource()
				ds2 := NewTagsDataSource()
				assert.NotNil(t, ds1)
				assert.NotNil(t, ds2)
				assert.NotSame(t, ds1, ds2)
			}
		})
	}
}

// TestTagsDataSource_mapTagsToItems verifies the mapTagsToItems method correctly maps SDK tags to items.
func TestTagsDataSource_mapTagsToItems(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantLen int
	}{
		{name: "empty tags slice", wantLen: 0},
		{name: "single tag with all fields", wantLen: 1},
		{name: "multiple tags", wantLen: 3},
		{name: "tag with nil id", wantLen: 1},
		{name: "tag with nil timestamps", wantLen: 1},
		{name: "error case - nil items slice", wantLen: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			d := &TagsDataSource{}

			switch tt.name {
			case "empty tags slice":
				result := d.mapTagsToItems(nil, nil)
				assert.Len(t, result, 0)

			case "single tag with all fields":
				tagID := "tag-1"
				tags := []n8nsdk.Tag{
					{Id: &tagID, Name: "Tag One"},
				}
				result := d.mapTagsToItems(tags, nil)
				assert.Len(t, result, 1)
				assert.Equal(t, "tag-1", result[0].ID.ValueString())
				assert.Equal(t, "Tag One", result[0].Name.ValueString())

			case "multiple tags":
				id1, id2, id3 := "tag-1", "tag-2", "tag-3"
				tags := []n8nsdk.Tag{
					{Id: &id1, Name: "Tag One"},
					{Id: &id2, Name: "Tag Two"},
					{Id: &id3, Name: "Tag Three"},
				}
				result := d.mapTagsToItems(tags, nil)
				assert.Len(t, result, 3)

			case "tag with nil id":
				tags := []n8nsdk.Tag{
					{Id: nil, Name: "Tag No ID"},
				}
				result := d.mapTagsToItems(tags, nil)
				assert.Len(t, result, 1)
				assert.True(t, result[0].ID.IsNull())
				assert.Equal(t, "Tag No ID", result[0].Name.ValueString())

			case "tag with nil timestamps":
				tagID := "tag-1"
				tags := []n8nsdk.Tag{
					{Id: &tagID, Name: "Tag One", CreatedAt: nil, UpdatedAt: nil},
				}
				result := d.mapTagsToItems(tags, nil)
				assert.Len(t, result, 1)
				assert.True(t, result[0].CreatedAt.IsNull())
				assert.True(t, result[0].UpdatedAt.IsNull())

			case "error case - nil items slice":
				tagID := "tag-1"
				tags := []n8nsdk.Tag{
					{Id: &tagID, Name: "Tag One"},
				}
				result := d.mapTagsToItems(tags, nil)
				assert.Len(t, result, 1)
			}
		})
	}
}
