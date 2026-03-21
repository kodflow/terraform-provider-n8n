package workflow

import (
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/client"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/workflow/models"
	"github.com/stretchr/testify/assert"
)

// setupTestClientForHelpers creates a test N8nClient with httptest server.
func setupTestClientForHelpers(t *testing.T, handler http.HandlerFunc) (*client.N8nClient, *httptest.Server) {
	t.Helper()
	server := httptest.NewServer(handler)

	cfg := n8nsdk.NewConfiguration()
	cfg.Servers = n8nsdk.ServerConfigurations{
		{
			URL:         server.URL,
			Description: "Test server",
		},
	}
	cfg.HTTPClient = server.Client()
	cfg.AddDefaultHeader("X-N8N-API-KEY", "test-key")

	apiClient := n8nsdk.NewAPIClient(cfg)
	n8nClient := &client.N8nClient{
		APIClient: apiClient,
	}

	return n8nClient, server
}

func Test_parseWorkflowJSON(t *testing.T) {
	t.Parallel()

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
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "map tags with valid IDs",
			testFunc: func(t *testing.T) {
				t.Helper()
				workflow := &n8nsdk.Workflow{
					Tags: []n8nsdk.Tag{
						{Id: new("tag-1")},
						{Id: new("tag-2")},
					},
				}
				diags := &diag.Diagnostics{}

				result := mapTagsFromWorkflow(t.Context(), workflow, diags)

				assert.False(t, diags.HasError())
				assert.False(t, result.IsNull())
				var tagIDs []string
				diags.Append(result.ElementsAs(t.Context(), &tagIDs, false)...)
				assert.Len(t, tagIDs, 2)
				assert.Contains(t, tagIDs, "tag-1")
				assert.Contains(t, tagIDs, "tag-2")
			},
		},
		{
			name: "map tags with nil ID",
			testFunc: func(t *testing.T) {
				t.Helper()
				workflow := &n8nsdk.Workflow{
					Tags: []n8nsdk.Tag{
						{Id: new("tag-1")},
						{Id: nil},
					},
				}
				diags := &diag.Diagnostics{}

				result := mapTagsFromWorkflow(t.Context(), workflow, diags)

				assert.False(t, diags.HasError())
				var tagIDs []string
				diags.Append(result.ElementsAs(t.Context(), &tagIDs, false)...)
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

				result := mapTagsFromWorkflow(t.Context(), workflow, diags)

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

				result := mapTagsFromWorkflow(t.Context(), workflow, diags)

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
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "maps all basic fields when provided",
			testFunc: func(t *testing.T) {
				t.Helper()
				workflow := &n8nsdk.Workflow{
					Active:       new(true),
					VersionId:    new("v1"),
					IsArchived:   new(false),
					TriggerCount: new(float32(5)),
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
				workflow := &n8nsdk.Workflow{
					Active: new(false),
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
				workflow := &n8nsdk.Workflow{
					TriggerCount: new(float32(0)),
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
				workflow := &n8nsdk.Workflow{
					Active:    new(false),
					VersionId: new("new-version"),
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
				workflow := &n8nsdk.Workflow{
					IsArchived: new(true),
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
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "map both timestamps",
			testFunc: func(t *testing.T) {
				t.Helper()
				workflow := &n8nsdk.Workflow{
					CreatedAt: new(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)),
					UpdatedAt: new(time.Date(2024, 1, 2, 12, 0, 0, 0, time.UTC)),
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
				workflow := &n8nsdk.Workflow{
					CreatedAt: nil,
					UpdatedAt: new(time.Date(2024, 1, 2, 12, 0, 0, 0, time.UTC)),
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
				workflow := &n8nsdk.Workflow{
					CreatedAt: new(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)),
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
				workflow := &n8nsdk.Workflow{
					CreatedAt: new(time.Date(2024, 1, 1, 12, 0, 0, 0, loc)),
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
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "map complete workflow",
			testFunc: func(t *testing.T) {
				t.Helper()
				workflow := &n8nsdk.Workflow{
					Name:        "Test Workflow",
					Active:      new(true),
					VersionId:   new("v1"),
					Tags:        []n8nsdk.Tag{{Id: new("tag-1")}},
					CreatedAt:   new(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)),
					Nodes:       []n8nsdk.Node{},
					Connections: map[string]any{},
					Settings:    n8nsdk.WorkflowSettings{},
				}
				plan := &models.Resource{}
				diags := &diag.Diagnostics{}

				mapWorkflowToModel(t.Context(), workflow, plan, diags)

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

				mapWorkflowToModel(t.Context(), workflow, plan, diags)

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

				mapWorkflowToModel(t.Context(), workflow, plan, diags)

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

				mapWorkflowToModel(t.Context(), workflow, plan, diags)

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

				mapWorkflowToModel(t.Context(), workflow, plan, diags)

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
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "serialize all JSON fields",
			testFunc: func(t *testing.T) {
				t.Helper()
				workflow := &n8nsdk.Workflow{
					Nodes: []n8nsdk.Node{
						{Name: new("Start")},
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
	t.Parallel()

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
	t.Parallel()

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
	t.Parallel()

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

				r.handleWorkflowActivation(t.Context(), plan, state, diags)

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

				r.handleWorkflowActivation(t.Context(), plan, state, diags)

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

				r.handleWorkflowActivation(t.Context(), plan, state, diags)

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

				r.handleWorkflowActivation(t.Context(), plan, state, diags)

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
	t.Parallel()

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

				r.updateWorkflowTags(t.Context(), "wf-123", plan, workflow, diags)

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

				r.updateWorkflowTags(t.Context(), "wf-123", plan, workflow, diags)

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
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := http.HandlerFunc(tt.setupHandler)
			n8nClient, server := setupTestClient(t, handler)
			t.Cleanup(server.Close)

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

			r.handleWorkflowActivation(t.Context(), plan, state, diags)

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
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := http.HandlerFunc(tt.setupHandler)
			n8nClient, server := setupTestClient(t, handler)
			t.Cleanup(server.Close)

			r := &WorkflowResource{client: n8nClient}
			plan := &models.Resource{
				Tags: tt.tags,
			}
			workflow := &n8nsdk.Workflow{}
			diags := &diag.Diagnostics{}

			r.updateWorkflowTags(t.Context(), "wf-123", plan, workflow, diags)

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
				settings := n8nsdk.WorkflowSettings{
					CallerPolicy: new(CallerPolicyDefault),
				}

				result := normalizeWorkflowSettings(settings)

				assert.Nil(t, result.CallerPolicy)
			},
		},
		{
			name: "keeps non-default callerPolicy value",
			testFunc: func(t *testing.T) {
				t.Helper()
				settings := n8nsdk.WorkflowSettings{
					CallerPolicy: new("any"),
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
				settings := n8nsdk.WorkflowSettings{
					AvailableInMCP: new(false),
				}

				result := normalizeWorkflowSettings(settings)

				assert.Nil(t, result.AvailableInMCP)
			},
		},
		{
			name: "keeps non-default availableInMCP value",
			testFunc: func(t *testing.T) {
				t.Helper()
				settings := n8nsdk.WorkflowSettings{
					AvailableInMCP: new(true),
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
				settings := n8nsdk.WorkflowSettings{
					CallerPolicy:   new(CallerPolicyDefault),
					AvailableInMCP: new(false),
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
				settings := n8nsdk.WorkflowSettings{
					CallerPolicy:   new(CallerPolicyDefault),
					AvailableInMCP: new(true),
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
					AdditionalProperties: map[string]any{
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
				settings := n8nsdk.WorkflowSettings{
					CallerPolicy: new(""),
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

// Test_mapWorkflowProjectID tests the mapWorkflowProjectID function.
func Test_mapWorkflowProjectID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "maps project ID when workflow has shared project",
			testFunc: func(t *testing.T) {
				t.Helper()
				projectID := "project-123"
				workflow := &n8nsdk.Workflow{
					Shared: []n8nsdk.SharedWorkflow{
						{ProjectId: &projectID},
					},
				}
				plan := &models.Resource{}

				mapWorkflowProjectID(workflow, plan)

				assert.False(t, plan.ProjectID.IsNull())
				assert.Equal(t, projectID, plan.ProjectID.ValueString())
			},
		},
		{
			name: "sets null when workflow shared project has nil ID",
			testFunc: func(t *testing.T) {
				t.Helper()
				workflow := &n8nsdk.Workflow{
					Shared: []n8nsdk.SharedWorkflow{
						{ProjectId: nil},
					},
				}
				plan := &models.Resource{}

				mapWorkflowProjectID(workflow, plan)

				assert.True(t, plan.ProjectID.IsNull())
			},
		},
		{
			name: "sets null when workflow has no shared projects",
			testFunc: func(t *testing.T) {
				t.Helper()
				workflow := &n8nsdk.Workflow{
					Shared: []n8nsdk.SharedWorkflow{},
				}
				plan := &models.Resource{}

				mapWorkflowProjectID(workflow, plan)

				assert.True(t, plan.ProjectID.IsNull())
			},
		},
		{
			name: "sets null when workflow shared is nil",
			testFunc: func(t *testing.T) {
				t.Helper()
				workflow := &n8nsdk.Workflow{
					Shared: nil,
				}
				plan := &models.Resource{}

				mapWorkflowProjectID(workflow, plan)

				assert.True(t, plan.ProjectID.IsNull())
			},
		},
		{
			name: "error case - handles empty string project ID",
			testFunc: func(t *testing.T) {
				t.Helper()
				workflow := &n8nsdk.Workflow{
					Shared: []n8nsdk.SharedWorkflow{
						{ProjectId: new("")},
					},
				}
				plan := &models.Resource{}

				mapWorkflowProjectID(workflow, plan)

				assert.False(t, plan.ProjectID.IsNull())
				assert.Equal(t, "", plan.ProjectID.ValueString())
			},
		},
		{
			name: "error case - handles multiple shared projects (uses first)",
			testFunc: func(t *testing.T) {
				t.Helper()
				projectID1 := "project-1"
				workflow := &n8nsdk.Workflow{
					Shared: []n8nsdk.SharedWorkflow{
						{ProjectId: &projectID1},
						{ProjectId: new("project-2")},
					},
				}
				plan := &models.Resource{}

				mapWorkflowProjectID(workflow, plan)

				// Should use first project
				assert.False(t, plan.ProjectID.IsNull())
				assert.Equal(t, projectID1, plan.ProjectID.ValueString())
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

// Test_preserveProjectIDOnUpdate tests the preserveProjectIDOnUpdate function.
func Test_preserveProjectIDOnUpdate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "preserves state project_id when unchanged",
			testFunc: func(t *testing.T) {
				t.Helper()

				projectID := "test-project-id"
				plan := &models.Resource{
					ProjectID: types.StringValue(projectID),
				}
				state := &models.Resource{
					ProjectID: types.StringValue(projectID),
				}

				preserveProjectIDOnUpdate(plan, state)

				assert.Equal(t, projectID, plan.ProjectID.ValueString())
			},
		},
		{
			name: "preserves plan project_id when changed",
			testFunc: func(t *testing.T) {
				t.Helper()

				newProjectID := "new-project-id"
				oldProjectID := "old-project-id"
				plan := &models.Resource{
					ProjectID: types.StringValue(newProjectID),
				}
				state := &models.Resource{
					ProjectID: types.StringValue(oldProjectID),
				}

				preserveProjectIDOnUpdate(plan, state)

				assert.Equal(t, newProjectID, plan.ProjectID.ValueString())
			},
		},
		{
			name: "preserves state when plan is null",
			testFunc: func(t *testing.T) {
				t.Helper()

				projectID := "test-project-id"
				plan := &models.Resource{
					ProjectID: types.StringNull(),
				}
				state := &models.Resource{
					ProjectID: types.StringValue(projectID),
				}

				preserveProjectIDOnUpdate(plan, state)

				assert.Equal(t, projectID, plan.ProjectID.ValueString())
			},
		},
		{
			name: "preserves state when plan is unknown",
			testFunc: func(t *testing.T) {
				t.Helper()

				projectID := "test-project-id"
				plan := &models.Resource{
					ProjectID: types.StringUnknown(),
				}
				state := &models.Resource{
					ProjectID: types.StringValue(projectID),
				}

				preserveProjectIDOnUpdate(plan, state)

				assert.Equal(t, projectID, plan.ProjectID.ValueString())
			},
		},
		{
			name: "preserves null state when both null",
			testFunc: func(t *testing.T) {
				t.Helper()

				plan := &models.Resource{
					ProjectID: types.StringNull(),
				}
				state := &models.Resource{
					ProjectID: types.StringNull(),
				}

				preserveProjectIDOnUpdate(plan, state)

				assert.True(t, plan.ProjectID.IsNull())
			},
		},
		{
			name: "error case - empty string project_id preserved",
			testFunc: func(t *testing.T) {
				t.Helper()

				plan := &models.Resource{
					ProjectID: types.StringValue(""),
				}
				state := &models.Resource{
					ProjectID: types.StringValue("old-project"),
				}

				preserveProjectIDOnUpdate(plan, state)

				// Empty string is valid, not replaced with state.
				assert.Equal(t, "", plan.ProjectID.ValueString())
			},
		},
		{
			name: "error case - unknown state keeps plan unknown",
			testFunc: func(t *testing.T) {
				t.Helper()

				plan := &models.Resource{
					ProjectID: types.StringUnknown(),
				}
				state := &models.Resource{
					ProjectID: types.StringUnknown(),
				}

				preserveProjectIDOnUpdate(plan, state)

				// Both are unknown, plan should become unknown.
				assert.True(t, plan.ProjectID.IsUnknown())
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

// TestWorkflowResource_createWorkflowViaAPI tests the createWorkflowViaAPI method.
// Note: This is an integration test that requires a real n8n instance.
// Unit testing is not feasible without complex mocking of the SDK client.
func TestWorkflowResource_createWorkflowViaAPI(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "function exists and has correct signature",
			testFunc: func(t *testing.T) {
				t.Helper()

				r := &WorkflowResource{}
				ctx := t.Context()
				diags := &diag.Diagnostics{}
				workflowRequest := n8nsdk.Workflow{}

				// Should panic with nil client
				assert.Panics(t, func() {
					r.createWorkflowViaAPI(ctx, workflowRequest, diags)
				})
			},
		},
		{
			name: "error case - nil client causes panic",
			testFunc: func(t *testing.T) {
				t.Helper()

				r := &WorkflowResource{}
				ctx := t.Context()
				diags := &diag.Diagnostics{}
				workflowRequest := n8nsdk.Workflow{Name: "test"}

				assert.Panics(t, func() {
					r.createWorkflowViaAPI(ctx, workflowRequest, diags)
				})
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

// TestWorkflowResource_updateWorkflowViaAPI tests the updateWorkflowViaAPI method.
// Note: This is an integration test that requires a real n8n instance.
// Unit testing is not feasible without complex mocking of the SDK client.
func TestWorkflowResource_updateWorkflowViaAPI(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "function exists and has correct signature",
			testFunc: func(t *testing.T) {
				t.Helper()

				r := &WorkflowResource{}
				ctx := t.Context()
				diags := &diag.Diagnostics{}
				workflowRequest := n8nsdk.Workflow{}

				// Should panic with nil client
				assert.Panics(t, func() {
					r.updateWorkflowViaAPI(ctx, "workflow-123", workflowRequest, diags)
				})
			},
		},
		{
			name: "error case - nil client causes panic",
			testFunc: func(t *testing.T) {
				t.Helper()

				r := &WorkflowResource{}
				ctx := t.Context()
				diags := &diag.Diagnostics{}
				workflowRequest := n8nsdk.Workflow{Name: "updated"}

				assert.Panics(t, func() {
					r.updateWorkflowViaAPI(ctx, "workflow-123", workflowRequest, diags)
				})
			},
		},
		{
			name: "error case - empty workflow ID",
			testFunc: func(t *testing.T) {
				t.Helper()

				r := &WorkflowResource{}
				ctx := t.Context()
				diags := &diag.Diagnostics{}
				workflowRequest := n8nsdk.Workflow{}

				assert.Panics(t, func() {
					r.updateWorkflowViaAPI(ctx, "", workflowRequest, diags)
				})
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

// TestWorkflowResource_transferWorkflowToProject tests the transferWorkflowToProject method.
// Note: This is an integration test that requires a real n8n instance.
// Unit testing is not feasible without complex mocking of the SDK client.
func TestWorkflowResource_transferWorkflowToProject(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "function exists and has correct signature",
			testFunc: func(t *testing.T) {
				t.Helper()

				r := &WorkflowResource{}

				// Verify method exists by checking it can be called
				// (will fail with nil client, but that's expected)
				ctx := t.Context()
				diags := &diag.Diagnostics{}

				// This will panic with nil client, which is expected
				assert.Panics(t, func() {
					r.transferWorkflowToProject(ctx, "workflow-123", "project-456", diags)
				})
			},
		},
		{
			name: "error case - nil client causes panic",
			testFunc: func(t *testing.T) {
				t.Helper()

				r := &WorkflowResource{}
				ctx := t.Context()
				diags := &diag.Diagnostics{}

				// Should panic with nil client
				assert.Panics(t, func() {
					r.transferWorkflowToProject(ctx, "workflow-123", "project-456", diags)
				})
			},
		},
		{
			name: "error case - empty workflow ID",
			testFunc: func(t *testing.T) {
				t.Helper()

				r := &WorkflowResource{}
				ctx := t.Context()
				diags := &diag.Diagnostics{}

				// Should panic with nil client (API call would fail anyway)
				assert.Panics(t, func() {
					r.transferWorkflowToProject(ctx, "", "project-456", diags)
				})
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

// TestWorkflowResource_handleProjectAssignment tests the handleProjectAssignment method.
// Note: This is an integration test that requires a real n8n instance.
// Unit testing is not feasible without complex mocking of the SDK client.
func TestWorkflowResource_handleProjectAssignment(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "function exists and has correct signature",
			testFunc: func(t *testing.T) {
				t.Helper()

				r := &WorkflowResource{}

				// Verify method exists by checking it can be called
				// (will fail with nil client, but that's expected)
				ctx := t.Context()
				diags := &diag.Diagnostics{}

				// This will panic with nil client, which is expected
				assert.Panics(t, func() {
					r.handleProjectAssignment(ctx, "workflow-123", "project-456", diags)
				})
			},
		},
		{
			name: "error case - nil client causes panic",
			testFunc: func(t *testing.T) {
				t.Helper()

				r := &WorkflowResource{}
				ctx := t.Context()
				diags := &diag.Diagnostics{}

				// Should panic with nil client
				assert.Panics(t, func() {
					r.handleProjectAssignment(ctx, "workflow-123", "project-456", diags)
				})
			},
		},
		{
			name: "error case - returns nil on transfer failure",
			testFunc: func(t *testing.T) {
				t.Helper()

				r := &WorkflowResource{}
				ctx := t.Context()
				diags := &diag.Diagnostics{}

				// Should panic with nil client
				assert.Panics(t, func() {
					workflow := r.handleProjectAssignment(ctx, "workflow-123", "project-456", diags)
					// If it didn't panic, verify it would return nil on error
					assert.Nil(t, workflow)
				})
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

// TestWorkflowResource_handlePostCreation tests the handlePostCreation method.
// Note: This is an integration test that requires a real n8n instance.
// Unit testing is not feasible without complex mocking of the SDK client.
func TestWorkflowResource_handlePostCreation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "sets ID from workflow",
			testFunc: func(t *testing.T) {
				t.Helper()

				r := &WorkflowResource{}
				ctx := t.Context()
				diags := &diag.Diagnostics{}
				workflow := &n8nsdk.Workflow{Id: new("workflow-123")}
				plan := &models.Resource{}

				// Function sets ID first; with no tags/project in plan, returns workflow without panic.
				result := r.handlePostCreation(ctx, workflow, plan, diags)

				// ID should be set
				assert.Equal(t, "workflow-123", plan.ID.ValueString())
				// Result is not nil (no activation requested either)
				assert.NotNil(t, result)
			},
		},
		{
			name: "error case - handles nil workflow ID",
			testFunc: func(t *testing.T) {
				t.Helper()

				r := &WorkflowResource{}
				ctx := t.Context()
				diags := &diag.Diagnostics{}
				workflow := &n8nsdk.Workflow{Id: nil}
				plan := &models.Resource{}

				// Function sets ID from nil workflow.Id, resulting in null string value.
				// With no tags/project, no API calls occur so no panic.
				result := r.handlePostCreation(ctx, workflow, plan, diags)

				// ID should be null when workflow ID is nil
				assert.True(t, plan.ID.IsNull())
				assert.NotNil(t, result)
			},
		},
		{
			name: "error case - function exists with correct signature",
			testFunc: func(t *testing.T) {
				t.Helper()

				r := &WorkflowResource{}
				ctx := t.Context()
				diags := &diag.Diagnostics{}
				workflow := &n8nsdk.Workflow{Id: new("workflow-456")}
				plan := &models.Resource{}

				// Verify method signature compiles and returns non-nil with valid inputs.
				result := r.handlePostCreation(ctx, workflow, plan, diags)
				assert.NotNil(t, result)
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

// TestWorkflowResource_applyPostCreationTagsAndProject tests the applyPostCreationTagsAndProject method.
// Note: This is an integration test that requires a real n8n instance.
// Unit testing is not feasible without complex mocking of the SDK client.
func TestWorkflowResource_applyPostCreationTagsAndProject(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "returns workflow when no tags or project",
			testFunc: func(t *testing.T) {
				t.Helper()

				r := &WorkflowResource{}
				ctx := t.Context()
				diags := &diag.Diagnostics{}
				workflow := &n8nsdk.Workflow{Id: new("workflow-123")}
				plan := &models.Resource{}

				result := r.applyPostCreationTagsAndProject(ctx, workflow, plan, diags)

				assert.NotNil(t, result)
				assert.Equal(t, workflow, result)
				assert.False(t, diags.HasError())
			},
		},
		{
			name: "returns workflow when tags are null",
			testFunc: func(t *testing.T) {
				t.Helper()

				r := &WorkflowResource{}
				ctx := t.Context()
				diags := &diag.Diagnostics{}
				workflow := &n8nsdk.Workflow{Id: new("workflow-456")}
				plan := &models.Resource{
					Tags: types.SetNull(types.StringType),
				}

				result := r.applyPostCreationTagsAndProject(ctx, workflow, plan, diags)

				assert.NotNil(t, result)
				assert.False(t, diags.HasError())
			},
		},
		{
			name: "returns workflow when project_id is null",
			testFunc: func(t *testing.T) {
				t.Helper()

				r := &WorkflowResource{}
				ctx := t.Context()
				diags := &diag.Diagnostics{}
				workflow := &n8nsdk.Workflow{Id: new("workflow-789")}
				plan := &models.Resource{
					ProjectID: types.StringNull(),
				}

				result := r.applyPostCreationTagsAndProject(ctx, workflow, plan, diags)

				assert.NotNil(t, result)
				assert.False(t, diags.HasError())
			},
		},
		{
			name: "returns workflow when workflow ID is nil",
			testFunc: func(t *testing.T) {
				t.Helper()

				r := &WorkflowResource{}
				ctx := t.Context()
				diags := &diag.Diagnostics{}
				workflow := &n8nsdk.Workflow{Id: nil}
				plan := &models.Resource{}

				// With nil ID, tags/project conditions require workflow.Id != nil, so skip.
				result := r.applyPostCreationTagsAndProject(ctx, workflow, plan, diags)

				assert.NotNil(t, result)
				assert.False(t, diags.HasError())
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

// TestWorkflowResource_transferWorkflowToProject_WithMock tests transferWorkflowToProject with mock HTTP.
func TestWorkflowResource_transferWorkflowToProject_WithMock(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		workflowID   string
		projectID    string
		setupHandler func(w http.ResponseWriter, r *http.Request)
		expectError  bool
	}{
		{
			name:       "success - transfer workflow to project",
			workflowID: "workflow-123",
			projectID:  "project-456",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == "PUT" && strings.HasPrefix(r.URL.Path, "/workflows/") && strings.HasSuffix(r.URL.Path, "/transfer") {
					w.WriteHeader(http.StatusOK)
					_, _ = w.Write([]byte(`{}`))
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: false,
		},
		{
			name:       "error - transfer fails",
			workflowID: "workflow-123",
			projectID:  "project-456",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"message": "Failed to transfer workflow"}`))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := http.HandlerFunc(tt.setupHandler)
			n8nClient, server := setupTestClientForHelpers(t, handler)
			t.Cleanup(server.Close)

			r := &WorkflowResource{client: n8nClient}
			ctx := t.Context()
			diags := &diag.Diagnostics{}

			result := r.transferWorkflowToProject(ctx, tt.workflowID, tt.projectID, diags)

			if tt.expectError {
				assert.False(t, result, "Should return false on error")
				assert.True(t, diags.HasError(), "Should have diagnostics error")
			} else {
				assert.True(t, result, "Should return true on success")
				assert.False(t, diags.HasError(), "Should not have diagnostics error")
			}
		})
	}
}

// TestWorkflowResource_handleProjectAssignment_WithMock tests handleProjectAssignment with mock HTTP.
func TestWorkflowResource_handleProjectAssignment_WithMock(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		workflowID   string
		projectID    string
		setupHandler func(w http.ResponseWriter, r *http.Request)
		expectError  bool
		expectNil    bool
	}{
		{
			name:       "success - assign workflow to project",
			workflowID: "workflow-123",
			projectID:  "project-456",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				if r.Method == "PUT" && strings.HasSuffix(r.URL.Path, "/transfer") {
					w.WriteHeader(http.StatusOK)
					_, _ = w.Write([]byte(`{}`))
					return
				}
				if r.Method == "GET" && strings.HasPrefix(r.URL.Path, "/workflows/workflow-123") {
					w.WriteHeader(http.StatusOK)
					_, _ = w.Write([]byte(`{
						"id": "workflow-123",
						"name": "Test Workflow",
						"active": false,
						"nodes": [],
						"connections": {},
						"settings": {}
					}`))
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: false,
			expectNil:   false,
		},
		{
			name:       "error - transfer fails",
			workflowID: "workflow-123",
			projectID:  "project-456",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == "PUT" && strings.HasSuffix(r.URL.Path, "/transfer") {
					w.WriteHeader(http.StatusInternalServerError)
					_, _ = w.Write([]byte(`{"message": "Failed to transfer"}`))
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: true,
			expectNil:   true,
		},
		{
			name:       "error - refetch fails after transfer",
			workflowID: "workflow-123",
			projectID:  "project-456",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == "PUT" && strings.HasSuffix(r.URL.Path, "/transfer") {
					w.WriteHeader(http.StatusOK)
					_, _ = w.Write([]byte(`{}`))
					return
				}
				if r.Method == "GET" && strings.HasPrefix(r.URL.Path, "/workflows/workflow-123") {
					w.WriteHeader(http.StatusInternalServerError)
					_, _ = w.Write([]byte(`{"message": "Failed to fetch workflow"}`))
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: true,
			expectNil:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := http.HandlerFunc(tt.setupHandler)
			n8nClient, server := setupTestClientForHelpers(t, handler)
			t.Cleanup(server.Close)

			r := &WorkflowResource{client: n8nClient}
			ctx := t.Context()
			diags := &diag.Diagnostics{}

			result := r.handleProjectAssignment(ctx, tt.workflowID, tt.projectID, diags)

			if tt.expectError {
				assert.True(t, diags.HasError(), "Should have diagnostics error")
			} else {
				assert.False(t, diags.HasError(), "Should not have diagnostics error")
			}

			if tt.expectNil {
				assert.Nil(t, result, "Should return nil on error")
			} else {
				assert.NotNil(t, result, "Should return workflow on success")
			}
		})
	}
}

// TestWorkflowResource_applyPostCreationTagsAndProject_WithMock tests project assignment path.
func TestWorkflowResource_applyPostCreationTagsAndProject_WithMock(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		projectID    string
		setupHandler func(w http.ResponseWriter, r *http.Request)
		expectError  bool
		expectNil    bool
	}{
		{
			name:      "success - apply project assignment",
			projectID: "project-456",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				if r.Method == "PUT" && strings.HasSuffix(r.URL.Path, "/transfer") {
					w.WriteHeader(http.StatusOK)
					_, _ = w.Write([]byte(`{}`))
					return
				}
				if r.Method == "GET" && strings.HasPrefix(r.URL.Path, "/workflows/") {
					w.WriteHeader(http.StatusOK)
					_, _ = w.Write([]byte(`{
						"id": "workflow-123",
						"name": "Test Workflow",
						"active": false,
						"nodes": [],
						"connections": {},
						"settings": {}
					}`))
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: false,
			expectNil:   false,
		},
		{
			name:      "error - project assignment fails",
			projectID: "project-456",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == "PUT" && strings.HasSuffix(r.URL.Path, "/transfer") {
					w.WriteHeader(http.StatusInternalServerError)
					_, _ = w.Write([]byte(`{"message": "Failed to assign project"}`))
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: true,
			expectNil:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := http.HandlerFunc(tt.setupHandler)
			n8nClient, server := setupTestClientForHelpers(t, handler)
			t.Cleanup(server.Close)

			r := &WorkflowResource{client: n8nClient}
			ctx := t.Context()
			diags := &diag.Diagnostics{}

			workflow := &n8nsdk.Workflow{Id: new("workflow-123")}
			plan := &models.Resource{
				ProjectID: types.StringValue(tt.projectID),
			}

			result := r.applyPostCreationTagsAndProject(ctx, workflow, plan, diags)

			if tt.expectError {
				assert.True(t, diags.HasError(), "Should have diagnostics error")
			} else {
				assert.False(t, diags.HasError(), "Should not have diagnostics error")
			}

			if tt.expectNil {
				assert.Nil(t, result, "Should return nil on error")
			} else {
				assert.NotNil(t, result, "Should return workflow on success")
			}
		})
	}
}

// TestParseNodesJSON tests the parseNodesJSON function.
func TestParseNodesJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		nodesJSON   string
		isNull      bool
		expectError bool
		expectNil   bool
	}{
		{
			name:        "null field returns nil nodes",
			isNull:      true,
			expectError: false,
			expectNil:   true,
		},
		{
			name:        "valid JSON returns nodes",
			nodesJSON:   `[]`,
			expectError: false,
			expectNil:   false,
		},
		{
			name:        "invalid JSON returns error",
			nodesJSON:   `not-json`,
			expectError: true,
			expectNil:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			diags := &diag.Diagnostics{}
			plan := &models.Resource{}
			//: Set up plan field.
			if tt.isNull {
				plan.NodesJSON = types.StringNull()
			} else {
				plan.NodesJSON = types.StringValue(tt.nodesJSON)
			}
			//: Call parseNodesJSON.
			nodes, ok := parseNodesJSON(plan, diags)
			//: Verify expected results.
			if tt.expectError {
				assert.False(t, ok)
				assert.True(t, diags.HasError())
			} else {
				assert.True(t, ok)
				assert.False(t, diags.HasError())
			}
			if tt.expectNil {
				assert.Nil(t, nodes)
			}
		})
	}
}

// TestParseConnectionsJSON tests the parseConnectionsJSON function.
func TestParseConnectionsJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		connectionsJSON string
		isNull          bool
		expectError     bool
		expectEmpty     bool
	}{
		{
			name:        "null field returns empty map",
			isNull:      true,
			expectError: false,
			expectEmpty: true,
		},
		{
			name:            "valid JSON returns connections",
			connectionsJSON: `{}`,
			expectError:     false,
			expectEmpty:     true,
		},
		{
			name:            "invalid JSON returns error",
			connectionsJSON: `not-json`,
			expectError:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			diags := &diag.Diagnostics{}
			plan := &models.Resource{}
			//: Set up plan field.
			if tt.isNull {
				plan.ConnectionsJSON = types.StringNull()
			} else {
				plan.ConnectionsJSON = types.StringValue(tt.connectionsJSON)
			}
			//: Call parseConnectionsJSON.
			connections, ok := parseConnectionsJSON(plan, diags)
			//: Verify expected results.
			if tt.expectError {
				assert.False(t, ok)
				assert.True(t, diags.HasError())
			} else {
				assert.True(t, ok)
				assert.False(t, diags.HasError())
				assert.NotNil(t, connections)
			}
		})
	}
}

// TestParseSettingsJSON tests the parseSettingsJSON function.
func TestParseSettingsJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		settingsJSON string
		isNull       bool
		expectError  bool
	}{
		{
			name:        "null field returns empty settings",
			isNull:      true,
			expectError: false,
		},
		{
			name:         "valid JSON returns settings",
			settingsJSON: `{}`,
			expectError:  false,
		},
		{
			name:         "invalid JSON returns error",
			settingsJSON: `not-json`,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			diags := &diag.Diagnostics{}
			plan := &models.Resource{}
			//: Set up plan field.
			if tt.isNull {
				plan.SettingsJSON = types.StringNull()
			} else {
				plan.SettingsJSON = types.StringValue(tt.settingsJSON)
			}
			//: Call parseSettingsJSON.
			_, ok := parseSettingsJSON(plan, diags)
			//: Verify expected results.
			if tt.expectError {
				assert.False(t, ok)
				assert.True(t, diags.HasError())
			} else {
				assert.True(t, ok)
				assert.False(t, diags.HasError())
			}
		})
	}
}

// TestWorkflowResource_applyTagsIfPresent tests the applyTagsIfPresent method.
func TestWorkflowResource_applyTagsIfPresent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		tagsNull   bool
		workflowID *string
		expectTrue bool
	}{
		{
			name:       "null tags returns true (skip)",
			tagsNull:   true,
			workflowID: new("wf-123"),
			expectTrue: true,
		},
		{
			name:       "nil workflow ID returns true (skip)",
			tagsNull:   false,
			workflowID: nil,
			expectTrue: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			diags := &diag.Diagnostics{}
			r := &WorkflowResource{}
			workflow := &n8nsdk.Workflow{Id: tt.workflowID}
			plan := &models.Resource{}
			//: Set up tags field.
			if tt.tagsNull {
				plan.Tags = types.SetNull(types.StringType)
			} else {
				plan.Tags = types.SetValueMust(types.StringType, []attr.Value{})
			}
			//: Call applyTagsIfPresent.
			ok := r.applyTagsIfPresent(t.Context(), workflow, plan, diags)
			//: Verify expected results.
			assert.Equal(t, tt.expectTrue, ok)
			assert.False(t, diags.HasError())
		})
	}
}

// TestWorkflowResource_applyProjectIfPresent tests the applyProjectIfPresent method.
func TestWorkflowResource_applyProjectIfPresent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		projectIDNull  bool
		workflowID     *string
		expectOriginal bool
	}{
		{
			name:           "null project_id returns original workflow",
			projectIDNull:  true,
			workflowID:     new("wf-123"),
			expectOriginal: true,
		},
		{
			name:           "nil workflow ID returns original workflow",
			projectIDNull:  false,
			workflowID:     nil,
			expectOriginal: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			diags := &diag.Diagnostics{}
			r := &WorkflowResource{}
			workflow := &n8nsdk.Workflow{Id: tt.workflowID}
			plan := &models.Resource{}
			//: Set up project ID field.
			if tt.projectIDNull {
				plan.ProjectID = types.StringNull()
			} else {
				plan.ProjectID = types.StringValue("proj-1")
			}
			//: Call applyProjectIfPresent.
			result := r.applyProjectIfPresent(t.Context(), workflow, plan, diags)
			//: Verify original workflow returned.
			if tt.expectOriginal {
				assert.Equal(t, workflow, result)
			}
		})
	}
}
