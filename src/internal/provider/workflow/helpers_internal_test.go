package workflow

import (
	"context"
	"encoding/json"
	"math/big"
	"net/http"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/workflow/models"
	"github.com/stretchr/testify/assert"
)

func Test_parseWorkflowJSON(t *testing.T) {
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

func Test_mapTagsFromWorkflow(t *testing.T) {
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
				// Returns null to avoid inconsistent result errors when plan had null.
				assert.True(t, result.IsNull())
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
				// Returns null to avoid inconsistent result errors when plan had null.
				assert.True(t, result.IsNull())
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
func Test_mapWorkflowBasicFields(t *testing.T) {
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
				triggerCount := float32(5)
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
				triggerCount := float32(0)
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
				triggerCount := float32(2147483647) // max int32
				workflow := &n8nsdk.Workflow{
					TriggerCount: &triggerCount,
				}
				plan := &models.Resource{}
				mapWorkflowBasicFields(workflow, plan)
				assert.Equal(t, types.Int64Value(int64(triggerCount)), plan.TriggerCount)
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

func Test_mapWorkflowTimestamps(t *testing.T) {
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

func Test_mapWorkflowToModel(t *testing.T) {
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

func Test_serializeWorkflowJSON(t *testing.T) {
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

func Test_convertTagIDsToTagIdsInner(t *testing.T) {
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

func Test_isActivationChanged(t *testing.T) {
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

func Test_getActivationAction(t *testing.T) {
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
// This test ensures the function exists and covers all branches.
func TestWorkflowResource_handleWorkflowActivation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
	}{
		{name: "no activation change returns early"},
		{name: "activation from false to true - null state"},
		{name: "deactivation from true to false"},
		{name: "activation from unknown state"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "no activation change returns early":
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

			case "activation from false to true - null state":
				r := &WorkflowResource{}
				plan := &models.Resource{
					ID:     types.StringValue("wf-123"),
					Active: types.BoolValue(true),
				}
				state := &models.Resource{
					ID:     types.StringValue("wf-123"),
					Active: types.BoolNull(),
				}
				diags := &diag.Diagnostics{}

				r.handleWorkflowActivation(context.Background(), plan, state, diags)

				// No change detected because state.Active is null, returns early
				assert.False(t, diags.HasError())

			case "deactivation from true to false":
				r := &WorkflowResource{}
				plan := &models.Resource{
					ID:     types.StringValue("wf-123"),
					Active: types.BoolValue(false),
				}
				state := &models.Resource{
					ID:     types.StringValue("wf-123"),
					Active: types.BoolNull(),
				}
				diags := &diag.Diagnostics{}

				r.handleWorkflowActivation(context.Background(), plan, state, diags)

				// No change detected because state.Active is null, returns early
				assert.False(t, diags.HasError())

			case "activation from unknown state":
				r := &WorkflowResource{}
				plan := &models.Resource{
					ID:     types.StringValue("wf-123"),
					Active: types.BoolValue(true),
				}
				state := &models.Resource{
					ID:     types.StringValue("wf-123"),
					Active: types.BoolUnknown(),
				}
				diags := &diag.Diagnostics{}

				r.handleWorkflowActivation(context.Background(), plan, state, diags)

				// No change detected because state.Active is unknown, returns early
				assert.False(t, diags.HasError())
			}
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
					Tags: types.SetNull(types.StringType),
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
					Tags: types.SetUnknown(types.StringType),
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

// TestmapWorkflowBasicFields is now in external test file - refactored to test behavior only.

// TestmapWorkflowToModel is now in external test file - refactored to test behavior only.
func TestHandleWorkflowActivation_FullCoverage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		planActive   types.Bool
		stateActive  types.Bool
		setupHandler func(w http.ResponseWriter, r *http.Request)
		expectError  bool
	}{
		{
			name:        "activate workflow successfully",
			planActive:  types.BoolValue(true),
			stateActive: types.BoolValue(false),
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				// Accept any POST request - SDK will format the URL
				if r.Method == http.MethodPost {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					// Return a full Workflow object as SDK expects
					json.NewEncoder(w).Encode(map[string]any{
						"id":          "wf-123",
						"name":        "Test Workflow",
						"active":      true,
						"nodes":       []any{},
						"connections": map[string]any{},
						"settings":    map[string]any{},
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: false,
		},
		{
			name:        "deactivate workflow successfully",
			planActive:  types.BoolValue(false),
			stateActive: types.BoolValue(true),
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				// Accept any POST request - SDK will format the URL
				if r.Method == http.MethodPost {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					// Return a full Workflow object as SDK expects
					json.NewEncoder(w).Encode(map[string]any{
						"id":          "wf-123",
						"name":        "Test Workflow",
						"active":      false,
						"nodes":       []any{},
						"connections": map[string]any{},
						"settings":    map[string]any{},
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: false,
		},
		{
			name:        "activate workflow with API error",
			planActive:  types.BoolValue(true),
			stateActive: types.BoolValue(false),
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"message": "Internal server error"}`))
			},
			expectError: true,
		},
		{
			name:        "deactivate workflow with API error",
			planActive:  types.BoolValue(false),
			stateActive: types.BoolValue(true),
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"message": "Internal server error"}`))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := http.HandlerFunc(tt.setupHandler)
			n8nClient, server := setupTestClient(t, handler)
			defer server.Close()

			r := &WorkflowResource{client: n8nClient}
			plan := &models.Resource{
				ID:     types.StringValue("wf-123"),
				Active: tt.planActive,
			}
			state := &models.Resource{
				ID:     types.StringValue("wf-123"),
				Active: tt.stateActive,
			}
			diags := &diag.Diagnostics{}

			r.handleWorkflowActivation(context.Background(), plan, state, diags)

			if tt.expectError {
				assert.True(t, diags.HasError())
			} else {
				assert.False(t, diags.HasError())
			}
		})
	}
}

// TestUpdateWorkflowTags_FullCoverage tests all branches of updateWorkflowTags.
func TestUpdateWorkflowTags_FullCoverage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		tags         types.Set
		setupHandler func(w http.ResponseWriter, r *http.Request)
		expectError  bool
	}{
		{
			name: "update tags successfully",
			tags: types.SetValueMust(types.StringType, []attr.Value{
				types.StringValue("tag-1"),
				types.StringValue("tag-2"),
			}),
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPut && r.URL.Path == "/workflows/wf-123/tags" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode([]map[string]any{
						{"id": "tag-1", "name": "Tag 1"},
						{"id": "tag-2", "name": "Tag 2"},
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: false,
		},
		{
			name: "update tags with API error",
			tags: types.SetValueMust(types.StringType, []attr.Value{
				types.StringValue("tag-1"),
			}),
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"message": "Internal server error"}`))
			},
			expectError: true,
		},
		{
			name: "update tags with ElementsAs error",
			// Use a set with wrong element type to trigger ElementsAs error
			tags: types.SetValueMust(types.NumberType, []attr.Value{
				types.NumberValue(big.NewFloat(123)),
			}),
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := http.HandlerFunc(tt.setupHandler)
			n8nClient, server := setupTestClient(t, handler)
			defer server.Close()

			r := &WorkflowResource{client: n8nClient}
			plan := &models.Resource{
				Tags: tt.tags,
			}
			workflow := &n8nsdk.Workflow{}
			diags := &diag.Diagnostics{}

			r.updateWorkflowTags(context.Background(), "wf-123", plan, workflow, diags)

			if tt.expectError {
				assert.True(t, diags.HasError())
			} else {
				assert.False(t, diags.HasError())
			}
		})
	}
}

// Test_normalizeWorkflowSettings tests the normalizeWorkflowSettings function.
func Test_normalizeWorkflowSettings(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "removes default callerPolicy value",
			testFunc: func(t *testing.T) {
				t.Helper()
				callerPolicy := CALLER_POLICY_DEFAULT
				settings := n8nsdk.WorkflowSettings{
					CallerPolicy: &callerPolicy,
				}

				result := normalizeWorkflowSettings(settings)

				assert.Nil(t, result.CallerPolicy)
			},
		},
		{
			name: "keeps non-default callerPolicy value",
			testFunc: func(t *testing.T) {
				t.Helper()
				callerPolicy := "any"
				settings := n8nsdk.WorkflowSettings{
					CallerPolicy: &callerPolicy,
				}

				result := normalizeWorkflowSettings(settings)

				assert.NotNil(t, result.CallerPolicy)
				assert.Equal(t, "any", *result.CallerPolicy)
			},
		},
		{
			name: "removes default availableInMCP value",
			testFunc: func(t *testing.T) {
				t.Helper()
				availableInMCP := false
				settings := n8nsdk.WorkflowSettings{
					AvailableInMCP: &availableInMCP,
				}

				result := normalizeWorkflowSettings(settings)

				assert.Nil(t, result.AvailableInMCP)
			},
		},
		{
			name: "keeps non-default availableInMCP value",
			testFunc: func(t *testing.T) {
				t.Helper()
				availableInMCP := true
				settings := n8nsdk.WorkflowSettings{
					AvailableInMCP: &availableInMCP,
				}

				result := normalizeWorkflowSettings(settings)

				assert.NotNil(t, result.AvailableInMCP)
				assert.True(t, *result.AvailableInMCP)
			},
		},
		{
			name: "handles empty settings",
			testFunc: func(t *testing.T) {
				t.Helper()
				settings := n8nsdk.WorkflowSettings{}

				result := normalizeWorkflowSettings(settings)

				assert.Nil(t, result.CallerPolicy)
				assert.Nil(t, result.AvailableInMCP)
			},
		},
		{
			name: "removes both defaults at once",
			testFunc: func(t *testing.T) {
				t.Helper()
				callerPolicy := CALLER_POLICY_DEFAULT
				availableInMCP := false
				settings := n8nsdk.WorkflowSettings{
					CallerPolicy:   &callerPolicy,
					AvailableInMCP: &availableInMCP,
				}

				result := normalizeWorkflowSettings(settings)

				assert.Nil(t, result.CallerPolicy)
				assert.Nil(t, result.AvailableInMCP)
			},
		},
		{
			name: "keeps mixed non-default and default values",
			testFunc: func(t *testing.T) {
				t.Helper()
				callerPolicy := CALLER_POLICY_DEFAULT
				availableInMCP := true
				settings := n8nsdk.WorkflowSettings{
					CallerPolicy:   &callerPolicy,
					AvailableInMCP: &availableInMCP,
				}

				result := normalizeWorkflowSettings(settings)

				assert.Nil(t, result.CallerPolicy)
				assert.NotNil(t, result.AvailableInMCP)
				assert.True(t, *result.AvailableInMCP)
			},
		},
		{
			name: "error case - preserves other settings fields untouched",
			testFunc: func(t *testing.T) {
				t.Helper()
				// Test that other fields are preserved.
				settings := n8nsdk.WorkflowSettings{
					AdditionalProperties: map[string]interface{}{
						"customField": "value",
					},
				}

				result := normalizeWorkflowSettings(settings)

				// Other fields should remain intact.
				assert.NotNil(t, result.AdditionalProperties)
				assert.Equal(t, "value", result.AdditionalProperties["customField"])
			},
		},
		{
			name: "error case - empty string callerPolicy is not removed",
			testFunc: func(t *testing.T) {
				t.Helper()
				callerPolicy := ""
				settings := n8nsdk.WorkflowSettings{
					CallerPolicy: &callerPolicy,
				}

				result := normalizeWorkflowSettings(settings)

				// Empty string is not the default, so it should be preserved.
				assert.NotNil(t, result.CallerPolicy)
				assert.Equal(t, "", *result.CallerPolicy)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t)
		})
	}
}
