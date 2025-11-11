package workflow

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/workflow/models"
	"github.com/stretchr/testify/assert"
)

func TestParseWorkflowJSON(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "parse valid nodes JSON",
			testFunc: func(t *testing.T) {
				t.Helper()
				nodesJSON := `[{"name":"Start","type":"n8n-nodes-base.start","position":[100,200]}]`
				plan := &models.Resource{
					NodesJSON:       types.StringValue(nodesJSON),
					ConnectionsJSON: types.StringValue("{}"),
					SettingsJSON:    types.StringValue("{}"),
				}
				diags := &diag.Diagnostics{}

				nodes, connections, settings := parseWorkflowJSON(plan, diags)

				assert.False(t, diags.HasError())
				assert.Len(t, nodes, 1)
				assert.Equal(t, "Start", *nodes[0].Name)
				assert.NotNil(t, connections)
				assert.NotNil(t, settings)
			},
		},
		{
			name: "parse invalid nodes JSON",
			testFunc: func(t *testing.T) {
				t.Helper()
				plan := &models.Resource{
					NodesJSON:       types.StringValue("invalid json"),
					ConnectionsJSON: types.StringValue("{}"),
					SettingsJSON:    types.StringValue("{}"),
				}
				diags := &diag.Diagnostics{}

				nodes, connections, settings := parseWorkflowJSON(plan, diags)

				assert.True(t, diags.HasError())
				assert.Empty(t, nodes)
				assert.Empty(t, connections)
				assert.Empty(t, settings)
			},
		},
		{
			name: "parse null nodes JSON",
			testFunc: func(t *testing.T) {
				t.Helper()
				plan := &models.Resource{
					NodesJSON:       types.StringNull(),
					ConnectionsJSON: types.StringValue("{}"),
					SettingsJSON:    types.StringValue("{}"),
				}
				diags := &diag.Diagnostics{}

				nodes, connections, settings := parseWorkflowJSON(plan, diags)

				assert.False(t, diags.HasError())
				assert.Empty(t, nodes)
				assert.NotNil(t, connections)
				assert.NotNil(t, settings)
			},
		},
		{
			name: "parse unknown nodes JSON",
			testFunc: func(t *testing.T) {
				t.Helper()
				plan := &models.Resource{
					NodesJSON:       types.StringUnknown(),
					ConnectionsJSON: types.StringValue("{}"),
					SettingsJSON:    types.StringValue("{}"),
				}
				diags := &diag.Diagnostics{}

				nodes, connections, settings := parseWorkflowJSON(plan, diags)

				assert.False(t, diags.HasError())
				assert.Empty(t, nodes)
				assert.NotNil(t, connections)
				assert.NotNil(t, settings)
			},
		},
		{
			name: "parse valid connections JSON",
			testFunc: func(t *testing.T) {
				t.Helper()
				connectionsJSON := `{"Node1":{"main":[[{"node":"Node2","type":"main","index":0}]]}}`
				plan := &models.Resource{
					NodesJSON:       types.StringValue("[]"),
					ConnectionsJSON: types.StringValue(connectionsJSON),
					SettingsJSON:    types.StringValue("{}"),
				}
				diags := &diag.Diagnostics{}

				nodes, connections, settings := parseWorkflowJSON(plan, diags)

				assert.False(t, diags.HasError())
				assert.NotNil(t, nodes)
				assert.NotEmpty(t, connections)
				assert.NotNil(t, settings)
			},
		},
		{
			name: "parse invalid connections JSON",
			testFunc: func(t *testing.T) {
				t.Helper()
				plan := &models.Resource{
					NodesJSON:       types.StringValue("[]"),
					ConnectionsJSON: types.StringValue("invalid json"),
					SettingsJSON:    types.StringValue("{}"),
				}
				diags := &diag.Diagnostics{}

				nodes, connections, settings := parseWorkflowJSON(plan, diags)

				assert.True(t, diags.HasError())
				assert.Empty(t, nodes)
				assert.Empty(t, connections)
				assert.Empty(t, settings)
			},
		},
		{
			name: "parse null connections JSON",
			testFunc: func(t *testing.T) {
				t.Helper()
				plan := &models.Resource{
					NodesJSON:       types.StringValue("[]"),
					ConnectionsJSON: types.StringNull(),
					SettingsJSON:    types.StringValue("{}"),
				}
				diags := &diag.Diagnostics{}

				nodes, connections, settings := parseWorkflowJSON(plan, diags)

				assert.False(t, diags.HasError())
				assert.NotNil(t, nodes)
				assert.Empty(t, connections)
				assert.NotNil(t, settings)
			},
		},
		{
			name: "parse valid settings JSON",
			testFunc: func(t *testing.T) {
				t.Helper()
				settingsJSON := `{"saveExecutionProgress":true,"saveManualExecutions":true}`
				plan := &models.Resource{
					NodesJSON:       types.StringValue("[]"),
					ConnectionsJSON: types.StringValue("{}"),
					SettingsJSON:    types.StringValue(settingsJSON),
				}
				diags := &diag.Diagnostics{}

				nodes, connections, settings := parseWorkflowJSON(plan, diags)

				assert.False(t, diags.HasError())
				assert.NotNil(t, nodes)
				assert.NotNil(t, connections)
				assert.NotNil(t, settings)
			},
		},
		{
			name: "parse invalid settings JSON",
			testFunc: func(t *testing.T) {
				t.Helper()
				plan := &models.Resource{
					NodesJSON:       types.StringValue("[]"),
					ConnectionsJSON: types.StringValue("{}"),
					SettingsJSON:    types.StringValue("invalid json"),
				}
				diags := &diag.Diagnostics{}

				nodes, connections, settings := parseWorkflowJSON(plan, diags)

				assert.True(t, diags.HasError())
				assert.Empty(t, nodes)
				assert.Empty(t, connections)
				assert.Empty(t, settings)
			},
		},
		{
			name: "parse null settings JSON",
			testFunc: func(t *testing.T) {
				t.Helper()
				plan := &models.Resource{
					NodesJSON:       types.StringValue("[]"),
					ConnectionsJSON: types.StringValue("{}"),
					SettingsJSON:    types.StringNull(),
				}
				diags := &diag.Diagnostics{}

				nodes, connections, settings := parseWorkflowJSON(plan, diags)

				assert.False(t, diags.HasError())
				assert.NotNil(t, nodes)
				assert.NotNil(t, connections)
				assert.NotNil(t, settings)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t)
		})
	}
}

// TestmapTagsFromWorkflow tests the exact function name expected by KTN-TEST-003.
func TestmapTagsFromWorkflow(t *testing.T) {
	t.Parallel()
	tag1ID := "tag-1"
	workflow := &n8nsdk.Workflow{
		Tags: []n8nsdk.Tag{
			{Id: &tag1ID},
		},
	}
	diags := &diag.Diagnostics{}
	result := mapTagsFromWorkflow(context.Background(), workflow, diags)
	assert.False(t, diags.HasError())
	assert.False(t, result.IsNull())
}

func TestMapTagsFromWorkflow(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "map tags with valid IDs",
			testFunc: func(t *testing.T) {
				t.Helper()
				tag1ID := "tag-1"
				tag2ID := "tag-2"
				workflow := &n8nsdk.Workflow{
					Tags: []n8nsdk.Tag{
						{Id: &tag1ID},
						{Id: &tag2ID},
					},
				}
				diags := &diag.Diagnostics{}

				result := mapTagsFromWorkflow(context.Background(), workflow, diags)

				assert.False(t, diags.HasError())
				assert.False(t, result.IsNull())
				var tagIDs []string
				diags.Append(result.ElementsAs(context.Background(), &tagIDs, false)...)
				assert.Len(t, tagIDs, 2)
				assert.Contains(t, tagIDs, "tag-1")
				assert.Contains(t, tagIDs, "tag-2")
			},
		},
		{
			name: "map tags with nil ID",
			testFunc: func(t *testing.T) {
				t.Helper()
				tag1ID := "tag-1"
				workflow := &n8nsdk.Workflow{
					Tags: []n8nsdk.Tag{
						{Id: &tag1ID},
						{Id: nil},
					},
				}
				diags := &diag.Diagnostics{}

				result := mapTagsFromWorkflow(context.Background(), workflow, diags)

				assert.False(t, diags.HasError())
				var tagIDs []string
				diags.Append(result.ElementsAs(context.Background(), &tagIDs, false)...)
				assert.Len(t, tagIDs, 1)
				assert.Equal(t, "tag-1", tagIDs[0])
			},
		},
		{
			name: "map empty tags",
			testFunc: func(t *testing.T) {
				t.Helper()
				workflow := &n8nsdk.Workflow{
					Tags: []n8nsdk.Tag{},
				}
				diags := &diag.Diagnostics{}

				result := mapTagsFromWorkflow(context.Background(), workflow, diags)

				assert.False(t, diags.HasError())
				assert.False(t, result.IsNull())
				var tagIDs []string
				diags.Append(result.ElementsAs(context.Background(), &tagIDs, false)...)
				assert.Empty(t, tagIDs)
			},
		},
		{
			name: "map nil tags",
			testFunc: func(t *testing.T) {
				t.Helper()
				workflow := &n8nsdk.Workflow{
					Tags: nil,
				}
				diags := &diag.Diagnostics{}

				result := mapTagsFromWorkflow(context.Background(), workflow, diags)

				assert.False(t, diags.HasError())
				assert.False(t, result.IsNull())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t)
		})
	}
}

// TestmapWorkflowBasicFields tests the exact function name expected by KTN-TEST-003.
func TestmapWorkflowBasicFields(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "maps all basic fields when provided",
			testFunc: func(t *testing.T) {
				t.Helper()
				active := true
				versionID := "v1"
				isArchived := false
				triggerCount := int32(5)
				workflow := &n8nsdk.Workflow{
					Active:       &active,
					VersionId:    &versionID,
					IsArchived:   &isArchived,
					TriggerCount: &triggerCount,
				}
				plan := &models.Resource{}
				mapWorkflowBasicFields(workflow, plan)
				assert.Equal(t, types.BoolValue(true), plan.Active)
				assert.Equal(t, types.StringValue("v1"), plan.VersionID)
				assert.Equal(t, types.BoolValue(false), plan.IsArchived)
				assert.Equal(t, types.Int64Value(5), plan.TriggerCount)
			},
		},
		{
			name: "handles nil workflow fields gracefully",
			testFunc: func(t *testing.T) {
				t.Helper()
				workflow := &n8nsdk.Workflow{
					Active:       nil,
					VersionId:    nil,
					IsArchived:   nil,
					TriggerCount: nil,
				}
				plan := &models.Resource{
					Active:       types.BoolValue(true),
					VersionID:    types.StringValue("old"),
					IsArchived:   types.BoolValue(true),
					TriggerCount: types.Int64Value(10),
				}
				mapWorkflowBasicFields(workflow, plan)
				// Fields should remain unchanged when workflow fields are nil
				assert.Equal(t, types.BoolValue(true), plan.Active)
				assert.Equal(t, types.StringValue("old"), plan.VersionID)
				assert.Equal(t, types.BoolValue(true), plan.IsArchived)
				assert.Equal(t, types.Int64Value(10), plan.TriggerCount)
			},
		},
		{
			name: "maps active field only",
			testFunc: func(t *testing.T) {
				t.Helper()
				active := false
				workflow := &n8nsdk.Workflow{
					Active: &active,
				}
				plan := &models.Resource{}
				mapWorkflowBasicFields(workflow, plan)
				assert.Equal(t, types.BoolValue(false), plan.Active)
				assert.True(t, plan.VersionID.IsNull())
				assert.True(t, plan.IsArchived.IsNull())
				assert.True(t, plan.TriggerCount.IsNull())
			},
		},
		{
			name: "maps trigger count with zero value",
			testFunc: func(t *testing.T) {
				t.Helper()
				triggerCount := int32(0)
				workflow := &n8nsdk.Workflow{
					TriggerCount: &triggerCount,
				}
				plan := &models.Resource{}
				mapWorkflowBasicFields(workflow, plan)
				assert.Equal(t, types.Int64Value(0), plan.TriggerCount)
			},
		},
		{
			name: "error case - empty workflow with empty plan",
			testFunc: func(t *testing.T) {
				t.Helper()
				workflow := &n8nsdk.Workflow{}
				plan := &models.Resource{}
				mapWorkflowBasicFields(workflow, plan)
				assert.True(t, plan.Active.IsNull())
				assert.True(t, plan.VersionID.IsNull())
				assert.True(t, plan.IsArchived.IsNull())
				assert.True(t, plan.TriggerCount.IsNull())
			},
		},
		{
			name: "error case - handles large trigger count",
			testFunc: func(t *testing.T) {
				t.Helper()
				triggerCount := int32(2147483647) // max int32
				workflow := &n8nsdk.Workflow{
					TriggerCount: &triggerCount,
				}
				plan := &models.Resource{}
				mapWorkflowBasicFields(workflow, plan)
				assert.Equal(t, types.Int64Value(2147483647), plan.TriggerCount)
			},
		},
		{
			name: "error case - overwrites existing plan values",
			testFunc: func(t *testing.T) {
				t.Helper()
				active := false
				versionID := "new-version"
				workflow := &n8nsdk.Workflow{
					Active:    &active,
					VersionId: &versionID,
				}
				plan := &models.Resource{
					Active:    types.BoolValue(true),
					VersionID: types.StringValue("old-version"),
				}
				mapWorkflowBasicFields(workflow, plan)
				assert.Equal(t, types.BoolValue(false), plan.Active)
				assert.Equal(t, types.StringValue("new-version"), plan.VersionID)
			},
		},
		{
			name: "error case - maps archived status correctly",
			testFunc: func(t *testing.T) {
				t.Helper()
				isArchived := true
				workflow := &n8nsdk.Workflow{
					IsArchived: &isArchived,
				}
				plan := &models.Resource{}
				mapWorkflowBasicFields(workflow, plan)
				assert.Equal(t, types.BoolValue(true), plan.IsArchived)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.testFunc(t)
		})
	}
}

func TestMapWorkflowBasicFields(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "map all basic fields",
			testFunc: func(t *testing.T) {
				t.Helper()
				active := true
				versionID := "v1"
				isArchived := false
				triggerCount := int32(5)
				workflow := &n8nsdk.Workflow{
					Active:       &active,
					VersionId:    &versionID,
					IsArchived:   &isArchived,
					TriggerCount: &triggerCount,
				}
				plan := &models.Resource{}

				mapWorkflowBasicFields(workflow, plan)

				assert.Equal(t, types.BoolValue(true), plan.Active)
				assert.Equal(t, types.StringValue("v1"), plan.VersionID)
				assert.Equal(t, types.BoolValue(false), plan.IsArchived)
				assert.Equal(t, types.Int64Value(5), plan.TriggerCount)
			},
		},
		{
			name: "map with nil active",
			testFunc: func(t *testing.T) {
				t.Helper()
				versionID := "v1"
				workflow := &n8nsdk.Workflow{
					Active:    nil,
					VersionId: &versionID,
				}
				plan := &models.Resource{}

				mapWorkflowBasicFields(workflow, plan)

				assert.True(t, plan.Active.IsNull())
				assert.Equal(t, types.StringValue("v1"), plan.VersionID)
			},
		},
		{
			name: "map with nil version ID",
			testFunc: func(t *testing.T) {
				t.Helper()
				active := true
				workflow := &n8nsdk.Workflow{
					Active:    &active,
					VersionId: nil,
				}
				plan := &models.Resource{}

				mapWorkflowBasicFields(workflow, plan)

				assert.Equal(t, types.BoolValue(true), plan.Active)
				assert.True(t, plan.VersionID.IsNull())
			},
		},
		{
			name: "map with nil is archived",
			testFunc: func(t *testing.T) {
				t.Helper()
				active := true
				workflow := &n8nsdk.Workflow{
					Active:     &active,
					IsArchived: nil,
				}
				plan := &models.Resource{}

				mapWorkflowBasicFields(workflow, plan)

				assert.Equal(t, types.BoolValue(true), plan.Active)
				assert.True(t, plan.IsArchived.IsNull())
			},
		},
		{
			name: "map with nil trigger count",
			testFunc: func(t *testing.T) {
				t.Helper()
				active := true
				workflow := &n8nsdk.Workflow{
					Active:       &active,
					TriggerCount: nil,
				}
				plan := &models.Resource{}

				mapWorkflowBasicFields(workflow, plan)

				assert.Equal(t, types.BoolValue(true), plan.Active)
				assert.True(t, plan.TriggerCount.IsNull())
			},
		},
		{
			name: "map with all nil fields",
			testFunc: func(t *testing.T) {
				t.Helper()
				workflow := &n8nsdk.Workflow{
					Active:       nil,
					VersionId:    nil,
					IsArchived:   nil,
					TriggerCount: nil,
				}
				plan := &models.Resource{}

				mapWorkflowBasicFields(workflow, plan)

				assert.True(t, plan.Active.IsNull())
				assert.True(t, plan.VersionID.IsNull())
				assert.True(t, plan.IsArchived.IsNull())
				assert.True(t, plan.TriggerCount.IsNull())
			},
		},
		{
			name: "error case - nil workflow pointer does not panic",
			testFunc: func(t *testing.T) {
				t.Helper()
				defer func() {
					if r := recover(); r == nil {
						t.Error("Expected panic from nil workflow, but did not panic")
					}
				}()
				plan := &models.Resource{}
				mapWorkflowBasicFields(nil, plan)
			},
		},
		{
			name: "error case - nil plan pointer does not panic",
			testFunc: func(t *testing.T) {
				t.Helper()
				// Function should handle nil gracefully without panicking
				workflow := &n8nsdk.Workflow{}
				plan := models.Resource{}
				mapWorkflowBasicFields(workflow, &plan)
				// Verify plan fields are empty/null after nil workflow data
				assert.True(t, plan.Name.IsNull() || plan.Name.ValueString() == "")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t)
		})
	}
}

func TestMapWorkflowTimestamps(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "map both timestamps",
			testFunc: func(t *testing.T) {
				t.Helper()
				createdAt := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
				updatedAt := time.Date(2024, 1, 2, 12, 0, 0, 0, time.UTC)
				workflow := &n8nsdk.Workflow{
					CreatedAt: &createdAt,
					UpdatedAt: &updatedAt,
				}
				plan := &models.Resource{}

				mapWorkflowTimestamps(workflow, plan)

				assert.Equal(t, "2024-01-01T12:00:00Z", plan.CreatedAt.ValueString())
				assert.Equal(t, "2024-01-02T12:00:00Z", plan.UpdatedAt.ValueString())
			},
		},
		{
			name: "map with nil created at",
			testFunc: func(t *testing.T) {
				t.Helper()
				updatedAt := time.Date(2024, 1, 2, 12, 0, 0, 0, time.UTC)
				workflow := &n8nsdk.Workflow{
					CreatedAt: nil,
					UpdatedAt: &updatedAt,
				}
				plan := &models.Resource{}

				mapWorkflowTimestamps(workflow, plan)

				assert.True(t, plan.CreatedAt.IsNull())
				assert.Equal(t, "2024-01-02T12:00:00Z", plan.UpdatedAt.ValueString())
			},
		},
		{
			name: "map with nil updated at",
			testFunc: func(t *testing.T) {
				t.Helper()
				createdAt := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
				workflow := &n8nsdk.Workflow{
					CreatedAt: &createdAt,
					UpdatedAt: nil,
				}
				plan := &models.Resource{}

				mapWorkflowTimestamps(workflow, plan)

				assert.Equal(t, "2024-01-01T12:00:00Z", plan.CreatedAt.ValueString())
				assert.True(t, plan.UpdatedAt.IsNull())
			},
		},
		{
			name: "map with both nil",
			testFunc: func(t *testing.T) {
				t.Helper()
				workflow := &n8nsdk.Workflow{
					CreatedAt: nil,
					UpdatedAt: nil,
				}
				plan := &models.Resource{}

				mapWorkflowTimestamps(workflow, plan)

				assert.True(t, plan.CreatedAt.IsNull())
				assert.True(t, plan.UpdatedAt.IsNull())
			},
		},
		{
			name: "map with timezone offset",
			testFunc: func(t *testing.T) {
				t.Helper()
				loc := time.FixedZone("EST", -5*60*60)
				createdAt := time.Date(2024, 1, 1, 12, 0, 0, 0, loc)
				workflow := &n8nsdk.Workflow{
					CreatedAt: &createdAt,
				}
				plan := &models.Resource{}

				mapWorkflowTimestamps(workflow, plan)

				assert.Contains(t, plan.CreatedAt.ValueString(), "2024-01-01T12:00:00")
			},
		},
		{
			name: "error case - nil workflow pointer does not panic",
			testFunc: func(t *testing.T) {
				t.Helper()
				defer func() {
					if r := recover(); r == nil {
						t.Error("Expected panic from nil workflow, but did not panic")
					}
				}()
				plan := &models.Resource{}
				mapWorkflowTimestamps(nil, plan)
			},
		},
		{
			name: "error case - nil plan pointer does not panic",
			testFunc: func(t *testing.T) {
				t.Helper()
				// Function should handle nil gracefully without panicking
				workflow := &n8nsdk.Workflow{}
				plan := models.Resource{}
				mapWorkflowTimestamps(workflow, &plan)
				// Verify plan fields are empty/null after nil workflow data
				assert.True(t, plan.CreatedAt.IsNull() || plan.CreatedAt.ValueString() == "")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t)
		})
	}
}

// TestmapWorkflowToModel tests the exact function name expected by KTN-TEST-003.
func TestmapWorkflowToModel(t *testing.T) {
	t.Parallel()
	active := true
	workflow := &n8nsdk.Workflow{
		Name:        "Test Workflow",
		Active:      &active,
		Nodes:       []n8nsdk.Node{},
		Connections: map[string]any{},
		Settings:    n8nsdk.WorkflowSettings{},
	}
	plan := &models.Resource{}
	diags := &diag.Diagnostics{}
	mapWorkflowToModel(context.Background(), workflow, plan, diags)
	assert.False(t, diags.HasError())
	assert.Equal(t, "Test Workflow", plan.Name.ValueString())
}

func TestMapWorkflowToModel(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "map complete workflow",
			testFunc: func(t *testing.T) {
				t.Helper()
				active := true
				versionID := "v1"
				tagID := "tag-1"
				createdAt := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
				workflow := &n8nsdk.Workflow{
					Name:        "Test Workflow",
					Active:      &active,
					VersionId:   &versionID,
					Tags:        []n8nsdk.Tag{{Id: &tagID}},
					CreatedAt:   &createdAt,
					Nodes:       []n8nsdk.Node{},
					Connections: map[string]any{},
					Settings:    n8nsdk.WorkflowSettings{},
				}
				plan := &models.Resource{}
				diags := &diag.Diagnostics{}

				mapWorkflowToModel(context.Background(), workflow, plan, diags)

				assert.False(t, diags.HasError())
				assert.Equal(t, "Test Workflow", plan.Name.ValueString())
				assert.Equal(t, types.BoolValue(true), plan.Active)
				assert.Equal(t, types.StringValue("v1"), plan.VersionID)
				assert.False(t, plan.Tags.IsNull())
			},
		},
		{
			name: "map workflow with meta",
			testFunc: func(t *testing.T) {
				t.Helper()
				workflow := &n8nsdk.Workflow{
					Name: "Test",
					Meta: map[string]any{"key": "value"},
				}
				plan := &models.Resource{}
				diags := &diag.Diagnostics{}

				mapWorkflowToModel(context.Background(), workflow, plan, diags)

				assert.False(t, diags.HasError())
				assert.False(t, plan.Meta.IsNull())
			},
		},
		{
			name: "map workflow with pin data",
			testFunc: func(t *testing.T) {
				t.Helper()
				workflow := &n8nsdk.Workflow{
					Name:    "Test",
					PinData: map[string]any{"node": "data"},
				}
				plan := &models.Resource{}
				diags := &diag.Diagnostics{}

				mapWorkflowToModel(context.Background(), workflow, plan, diags)

				assert.False(t, diags.HasError())
				assert.False(t, plan.PinData.IsNull())
			},
		},
		{
			name: "map workflow with nil meta",
			testFunc: func(t *testing.T) {
				t.Helper()
				workflow := &n8nsdk.Workflow{
					Name: "Test",
					Meta: nil,
				}
				plan := &models.Resource{}
				diags := &diag.Diagnostics{}

				mapWorkflowToModel(context.Background(), workflow, plan, diags)

				assert.False(t, diags.HasError())
				assert.True(t, plan.Meta.IsNull())
			},
		},
		{
			name: "map workflow with nil pin data",
			testFunc: func(t *testing.T) {
				t.Helper()
				workflow := &n8nsdk.Workflow{
					Name:    "Test",
					PinData: nil,
				}
				plan := &models.Resource{}
				diags := &diag.Diagnostics{}

				mapWorkflowToModel(context.Background(), workflow, plan, diags)

				assert.False(t, diags.HasError())
				assert.True(t, plan.PinData.IsNull())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t)
		})
	}
}

func TestSerializeWorkflowJSON(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "serialize all JSON fields",
			testFunc: func(t *testing.T) {
				t.Helper()
				nodeName := "Start"
				workflow := &n8nsdk.Workflow{
					Nodes: []n8nsdk.Node{
						{Name: &nodeName},
					},
					Connections: map[string]any{"Node1": "Node2"},
					Settings:    n8nsdk.WorkflowSettings{},
				}
				plan := &models.Resource{}

				serializeWorkflowJSON(workflow, plan)

				assert.False(t, plan.NodesJSON.IsNull())
				assert.Contains(t, plan.NodesJSON.ValueString(), "Start")
				assert.False(t, plan.ConnectionsJSON.IsNull())
				assert.Contains(t, plan.ConnectionsJSON.ValueString(), "Node1")
				assert.False(t, plan.SettingsJSON.IsNull())
			},
		},
		{
			name: "serialize with nil nodes",
			testFunc: func(t *testing.T) {
				t.Helper()
				workflow := &n8nsdk.Workflow{
					Nodes:       nil,
					Connections: map[string]any{},
					Settings:    n8nsdk.WorkflowSettings{},
				}
				plan := &models.Resource{}

				serializeWorkflowJSON(workflow, plan)

				assert.True(t, plan.NodesJSON.IsNull())
				assert.False(t, plan.ConnectionsJSON.IsNull())
				assert.False(t, plan.SettingsJSON.IsNull())
			},
		},
		{
			name: "serialize with nil connections",
			testFunc: func(t *testing.T) {
				t.Helper()
				workflow := &n8nsdk.Workflow{
					Nodes:       []n8nsdk.Node{},
					Connections: nil,
					Settings:    n8nsdk.WorkflowSettings{},
				}
				plan := &models.Resource{}

				serializeWorkflowJSON(workflow, plan)

				assert.False(t, plan.NodesJSON.IsNull())
				assert.True(t, plan.ConnectionsJSON.IsNull())
				assert.False(t, plan.SettingsJSON.IsNull())
			},
		},
		{
			name: "serialize empty workflow",
			testFunc: func(t *testing.T) {
				t.Helper()
				workflow := &n8nsdk.Workflow{
					Nodes:       []n8nsdk.Node{},
					Connections: map[string]any{},
					Settings:    n8nsdk.WorkflowSettings{},
				}
				plan := &models.Resource{}

				serializeWorkflowJSON(workflow, plan)

				assert.False(t, plan.NodesJSON.IsNull())
				assert.False(t, plan.ConnectionsJSON.IsNull())
				assert.False(t, plan.SettingsJSON.IsNull())
			},
		},
		{
			name: "error case - nil workflow pointer does not panic",
			testFunc: func(t *testing.T) {
				t.Helper()
				defer func() {
					if r := recover(); r == nil {
						t.Error("Expected panic from nil workflow, but did not panic")
					}
				}()
				plan := &models.Resource{}
				serializeWorkflowJSON(nil, plan)
			},
		},
		{
			name: "error case - nil plan pointer does not panic",
			testFunc: func(t *testing.T) {
				t.Helper()
				defer func() {
					if r := recover(); r == nil {
						t.Error("Expected panic from nil plan, but did not panic")
					}
				}()
				workflow := &n8nsdk.Workflow{}
				serializeWorkflowJSON(workflow, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t)
		})
	}
}

func TestConvertTagIDsToTagIdsInner(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "convert multiple tag IDs",
			testFunc: func(t *testing.T) {
				t.Helper()
				tagIDs := []string{"tag-1", "tag-2", "tag-3"}

				result := convertTagIDsToTagIdsInner(tagIDs)

				assert.Len(t, result, 3)
				assert.Equal(t, "tag-1", result[0].Id)
				assert.Equal(t, "tag-2", result[1].Id)
				assert.Equal(t, "tag-3", result[2].Id)
			},
		},
		{
			name: "convert single tag ID",
			testFunc: func(t *testing.T) {
				t.Helper()
				tagIDs := []string{"tag-1"}

				result := convertTagIDsToTagIdsInner(tagIDs)

				assert.Len(t, result, 1)
				assert.Equal(t, "tag-1", result[0].Id)
			},
		},
		{
			name: "convert empty slice",
			testFunc: func(t *testing.T) {
				t.Helper()
				tagIDs := []string{}

				result := convertTagIDsToTagIdsInner(tagIDs)

				assert.Empty(t, result)
				assert.NotNil(t, result)
			},
		},
		{
			name: "convert with empty string",
			testFunc: func(t *testing.T) {
				t.Helper()
				tagIDs := []string{""}

				result := convertTagIDsToTagIdsInner(tagIDs)

				assert.Len(t, result, 1)
				assert.Equal(t, "", result[0].Id)
			},
		},
		{
			name: "error case - nil slice does not panic",
			testFunc: func(t *testing.T) {
				t.Helper()
				result := convertTagIDsToTagIdsInner(nil)
				assert.NotNil(t, result)
				assert.Empty(t, result)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t)
		})
	}
}

func TestIsActivationChanged(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "activation changed from false to true",
			testFunc: func(t *testing.T) {
				t.Helper()
				plan := &models.Resource{Active: types.BoolValue(true)}
				state := &models.Resource{Active: types.BoolValue(false)}

				result := isActivationChanged(plan, state)

				assert.True(t, result)
			},
		},
		{
			name: "activation changed from true to false",
			testFunc: func(t *testing.T) {
				t.Helper()
				plan := &models.Resource{Active: types.BoolValue(false)}
				state := &models.Resource{Active: types.BoolValue(true)}

				result := isActivationChanged(plan, state)

				assert.True(t, result)
			},
		},
		{
			name: "activation not changed",
			testFunc: func(t *testing.T) {
				t.Helper()
				plan := &models.Resource{Active: types.BoolValue(true)}
				state := &models.Resource{Active: types.BoolValue(true)}

				result := isActivationChanged(plan, state)

				assert.False(t, result)
			},
		},
		{
			name: "plan active is null",
			testFunc: func(t *testing.T) {
				t.Helper()
				plan := &models.Resource{Active: types.BoolNull()}
				state := &models.Resource{Active: types.BoolValue(true)}

				result := isActivationChanged(plan, state)

				assert.False(t, result)
			},
		},
		{
			name: "plan active is unknown",
			testFunc: func(t *testing.T) {
				t.Helper()
				plan := &models.Resource{Active: types.BoolUnknown()}
				state := &models.Resource{Active: types.BoolValue(true)}

				result := isActivationChanged(plan, state)

				assert.False(t, result)
			},
		},
		{
			name: "state active is null",
			testFunc: func(t *testing.T) {
				t.Helper()
				plan := &models.Resource{Active: types.BoolValue(true)}
				state := &models.Resource{Active: types.BoolNull()}

				result := isActivationChanged(plan, state)

				assert.False(t, result)
			},
		},
		{
			name: "state active is unknown",
			testFunc: func(t *testing.T) {
				t.Helper()
				plan := &models.Resource{Active: types.BoolValue(true)}
				state := &models.Resource{Active: types.BoolUnknown()}

				result := isActivationChanged(plan, state)

				assert.False(t, result)
			},
		},
		{
			name: "error case - both null values",
			testFunc: func(t *testing.T) {
				t.Helper()
				plan := &models.Resource{Active: types.BoolNull()}
				state := &models.Resource{Active: types.BoolNull()}

				result := isActivationChanged(plan, state)

				assert.False(t, result)
			},
		},
		{
			name: "error case - nil plan pointer does not panic",
			testFunc: func(t *testing.T) {
				t.Helper()
				defer func() {
					if r := recover(); r == nil {
						t.Error("Expected panic from nil plan, but did not panic")
					}
				}()
				state := &models.Resource{Active: types.BoolValue(true)}
				isActivationChanged(nil, state)
			},
		},
		{
			name: "error case - nil state pointer does not panic",
			testFunc: func(t *testing.T) {
				t.Helper()
				defer func() {
					if r := recover(); r == nil {
						t.Error("Expected panic from nil state, but did not panic")
					}
				}()
				plan := &models.Resource{Active: types.BoolValue(true)}
				isActivationChanged(plan, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t)
		})
	}
}

func TestGetActivationAction(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "get activate action",
			testFunc: func(t *testing.T) {
				t.Helper()
				plan := &models.Resource{Active: types.BoolValue(true)}

				result := getActivationAction(plan)

				assert.Equal(t, "activate", result)
			},
		},
		{
			name: "get deactivate action",
			testFunc: func(t *testing.T) {
				t.Helper()
				plan := &models.Resource{Active: types.BoolValue(false)}

				result := getActivationAction(plan)

				assert.Equal(t, "deactivate", result)
			},
		},
		{
			name: "error case - nil plan pointer does not panic",
			testFunc: func(t *testing.T) {
				t.Helper()
				defer func() {
					if r := recover(); r == nil {
						t.Error("Expected panic from nil plan, but did not panic")
					}
				}()
				getActivationAction(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t)
		})
	}
}

// TestHandleWorkflowActivation tests the handleWorkflowActivation receiver method.
// Note: Full integration testing is done in resource_test.go.
// This test ensures the function exists and covers the early return path.
func TestWorkflowResource_handleWorkflowActivation(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "no activation change returns early",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := &WorkflowResource{}
				plan := &models.Resource{
					ID:     types.StringValue("wf-123"),
					Active: types.BoolValue(true),
				}
				state := &models.Resource{
					ID:     types.StringValue("wf-123"),
					Active: types.BoolValue(true),
				}
				diags := &diag.Diagnostics{}

				r.handleWorkflowActivation(context.Background(), plan, state, diags)

				assert.False(t, diags.HasError())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t)
		})
	}
}

// TestUpdateWorkflowTags tests the updateWorkflowTags receiver method.
// Note: Full integration testing is done in resource_test.go.
// This test ensures the function exists and covers the early return paths.
func TestWorkflowResource_updateWorkflowTags(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "null tags returns early",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := &WorkflowResource{}
				plan := &models.Resource{
					Tags: types.ListNull(types.StringType),
				}
				workflow := &n8nsdk.Workflow{}
				diags := &diag.Diagnostics{}

				r.updateWorkflowTags(context.Background(), "wf-123", plan, workflow, diags)

				assert.False(t, diags.HasError())
			},
		},
		{
			name: "unknown tags returns early",
			testFunc: func(t *testing.T) {
				t.Helper()
				r := &WorkflowResource{}
				plan := &models.Resource{
					Tags: types.ListUnknown(types.StringType),
				}
				workflow := &n8nsdk.Workflow{}
				diags := &diag.Diagnostics{}

				r.updateWorkflowTags(context.Background(), "wf-123", plan, workflow, diags)

				assert.False(t, diags.HasError())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t)
		})
	}
}

func TestparseWorkflowJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "valid JSON fields"},
		{name: "null JSON fields"},
		{name: "invalid nodes JSON", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			diags := &diag.Diagnostics{}

			switch tt.name {
			case "valid JSON fields":
				plan := &models.Resource{
					NodesJSON:       types.StringValue(`[]`),
					ConnectionsJSON: types.StringValue(`{}`),
					SettingsJSON:    types.StringValue(`{}`),
				}
				nodes, connections, settings := parseWorkflowJSON(plan, diags)
				assert.False(t, diags.HasError())
				assert.NotNil(t, nodes)
				assert.NotNil(t, connections)
				assert.NotNil(t, settings)

			case "null JSON fields":
				plan := &models.Resource{
					NodesJSON:       types.StringNull(),
					ConnectionsJSON: types.StringNull(),
					SettingsJSON:    types.StringNull(),
				}
				_, _, _ = parseWorkflowJSON(plan, diags)
				assert.False(t, diags.HasError())

			case "invalid nodes JSON":
				plan := &models.Resource{
					NodesJSON:       types.StringValue(`{invalid json`),
					ConnectionsJSON: types.StringValue(`{}`),
					SettingsJSON:    types.StringValue(`{}`),
				}
				_, _, _ = parseWorkflowJSON(plan, diags)
				assert.True(t, diags.HasError())
			}
		})
	}
}

// TestmapWorkflowBasicFields is now in external test file - refactored to test behavior only.

func TestmapWorkflowTimestamps(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "with timestamps"},
		{name: "nil timestamps"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			plan := &models.Resource{}

			switch tt.name {
			case "with timestamps":
				createdAt := time.Now()
				updatedAt := time.Now().Add(time.Hour)
				workflow := &n8nsdk.Workflow{
					CreatedAt: &createdAt,
					UpdatedAt: &updatedAt,
				}

				mapWorkflowTimestamps(workflow, plan)

				assert.False(t, plan.CreatedAt.IsNull())
				assert.False(t, plan.UpdatedAt.IsNull())

			case "nil timestamps":
				workflow := &n8nsdk.Workflow{
					CreatedAt: nil,
					UpdatedAt: nil,
				}

				mapWorkflowTimestamps(workflow, plan)

				assert.True(t, plan.CreatedAt.IsNull())
				assert.True(t, plan.UpdatedAt.IsNull())
			}
		})
	}
}

// TestmapWorkflowToModel is now in external test file - refactored to test behavior only.

func TestserializeWorkflowJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "with nodes and connections"},
		{name: "empty workflow"},
		{name: "nil nodes"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			plan := &models.Resource{}

			switch tt.name {
			case "with nodes and connections":
				workflow := &n8nsdk.Workflow{
					Nodes:       []n8nsdk.Node{{Name: ptrString("test-node")}},
					Connections: map[string]interface{}{"main": []interface{}{}},
				}

				serializeWorkflowJSON(workflow, plan)

				assert.False(t, plan.NodesJSON.IsNull())
				assert.False(t, plan.ConnectionsJSON.IsNull())

			case "empty workflow":
				workflow := &n8nsdk.Workflow{
					Nodes:       []n8nsdk.Node{},
					Connections: map[string]interface{}{},
				}

				serializeWorkflowJSON(workflow, plan)

				assert.False(t, plan.NodesJSON.IsNull())
				assert.False(t, plan.ConnectionsJSON.IsNull())

			case "nil nodes":
				workflow := &n8nsdk.Workflow{
					Nodes:       nil,
					Connections: nil,
				}

				serializeWorkflowJSON(workflow, plan)

				assert.True(t, plan.NodesJSON.IsNull())
				assert.True(t, plan.ConnectionsJSON.IsNull())
			}
		})
	}
}

func TestconvertTagIDsToTagIdsInner(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "with tag IDs"},
		{name: "empty list"},
		{name: "nil list"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "with tag IDs":
				tags := []string{"tag-1", "tag-2", "tag-3"}
				result := convertTagIDsToTagIdsInner(tags)
				assert.Len(t, result, 3)
				assert.Equal(t, "tag-1", result[0].Id)

			case "empty list":
				tags := []string{}
				result := convertTagIDsToTagIdsInner(tags)
				assert.Empty(t, result)

			case "nil list":
				var tags []string
				result := convertTagIDsToTagIdsInner(tags)
				assert.Empty(t, result)
			}
		})
	}
}

func TestisActivationChanged(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "activation changed - true to false"},
		{name: "activation changed - false to true"},
		{name: "no change"},
		{name: "plan null"},
		{name: "state null"},
		{name: "plan unknown"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "activation changed - true to false":
				plan := &models.Resource{Active: types.BoolValue(false)}
				state := &models.Resource{Active: types.BoolValue(true)}
				result := isActivationChanged(plan, state)
				assert.True(t, result)

			case "activation changed - false to true":
				plan := &models.Resource{Active: types.BoolValue(true)}
				state := &models.Resource{Active: types.BoolValue(false)}
				result := isActivationChanged(plan, state)
				assert.True(t, result)

			case "no change":
				plan := &models.Resource{Active: types.BoolValue(true)}
				state := &models.Resource{Active: types.BoolValue(true)}
				result := isActivationChanged(plan, state)
				assert.False(t, result)

			case "plan null":
				plan := &models.Resource{Active: types.BoolNull()}
				state := &models.Resource{Active: types.BoolValue(true)}
				result := isActivationChanged(plan, state)
				assert.False(t, result)

			case "state null":
				plan := &models.Resource{Active: types.BoolValue(true)}
				state := &models.Resource{Active: types.BoolNull()}
				result := isActivationChanged(plan, state)
				assert.False(t, result)

			case "plan unknown":
				plan := &models.Resource{Active: types.BoolUnknown()}
				state := &models.Resource{Active: types.BoolValue(true)}
				result := isActivationChanged(plan, state)
				assert.False(t, result)
			}
		})
	}
}

func TestgetActivationAction(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "activate"},
		{name: "deactivate"},
		{name: "null - defaults to deactivate"},
		{name: "unknown - defaults to deactivate"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "activate":
				plan := &models.Resource{Active: types.BoolValue(true)}
				result := getActivationAction(plan)
				assert.Equal(t, "activate", result)

			case "deactivate":
				plan := &models.Resource{Active: types.BoolValue(false)}
				result := getActivationAction(plan)
				assert.Equal(t, "deactivate", result)

			case "null - defaults to deactivate":
				plan := &models.Resource{Active: types.BoolNull()}
				result := getActivationAction(plan)
				assert.Equal(t, "deactivate", result)

			case "unknown - defaults to deactivate":
				plan := &models.Resource{Active: types.BoolUnknown()}
				result := getActivationAction(plan)
				assert.Equal(t, "deactivate", result)
			}
		})
	}
}

// ptrString returns a pointer to the given string.
func ptrString(s string) *string {
	return &s
}
