package workflow

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
	"github.com/kodflow/n8n/src/internal/provider/workflow/models"
	"github.com/stretchr/testify/assert"
)

// Helper functions for tests.
func strPtr(s string) *string {
	return &s
}

func TestParseWorkflowJSON(t *testing.T) {
	t.Run("parse valid nodes JSON", func(t *testing.T) {
		plan := &models.Resource{
			NodesJSON: types.StringValue(`[{"name":"Start","type":"n8n-nodes-base.start","position":[100,200]}]`),
		}
		diags := &diag.Diagnostics{}

		nodes, connections, settings := parseWorkflowJSON(plan, diags)

		if diags.HasError() {
			for _, err := range diags.Errors() {
				t.Logf("Diagnostic Error: %s - %s", err.Summary(), err.Detail())
			}
		}
		assert.False(t, diags.HasError())
		assert.Len(t, nodes, 1)
		assert.NotNil(t, nodes[0].Name)
		assert.Equal(t, "Start", *nodes[0].Name)
		assert.NotNil(t, connections)
		assert.NotNil(t, settings)
	})

	t.Run("parse valid connections JSON", func(t *testing.T) {
		plan := &models.Resource{
			ConnectionsJSON: types.StringValue(`{"Start":{"main":[[{"node":"HTTP","type":"main","index":0}]]}}`),
		}
		diags := &diag.Diagnostics{}

		nodes, connections, settings := parseWorkflowJSON(plan, diags)

		if diags.HasError() {
			for _, err := range diags.Errors() {
				t.Logf("Diagnostic Error: %s - %s", err.Summary(), err.Detail())
			}
		}
		assert.False(t, diags.HasError())
		assert.NotNil(t, nodes)
		assert.NotNil(t, connections)
		assert.Contains(t, connections, "Start")
		assert.NotNil(t, settings)
	})

	t.Run("parse valid settings JSON", func(t *testing.T) {
		plan := &models.Resource{
			SettingsJSON: types.StringValue(`{"saveDataErrorExecution":"all","saveDataSuccessExecution":"all"}`),
		}
		diags := &diag.Diagnostics{}

		nodes, connections, settings := parseWorkflowJSON(plan, diags)

		if diags.HasError() {
			for _, err := range diags.Errors() {
				t.Logf("Diagnostic Error: %s - %s", err.Summary(), err.Detail())
			}
		}
		assert.False(t, diags.HasError())
		assert.NotNil(t, nodes)
		assert.NotNil(t, connections)
		assert.NotNil(t, settings)
	})

	t.Run("parse invalid nodes JSON", func(t *testing.T) {
		plan := &models.Resource{
			NodesJSON: types.StringValue(`invalid json`),
		}
		diags := &diag.Diagnostics{}

		nodes, connections, settings := parseWorkflowJSON(plan, diags)

		assert.True(t, diags.HasError())
		assert.Contains(t, diags.Errors()[0].Summary(), "Invalid nodes JSON")
		assert.Empty(t, nodes)
		assert.Empty(t, connections)
		assert.NotNil(t, settings)
	})

	t.Run("parse invalid connections JSON", func(t *testing.T) {
		plan := &models.Resource{
			ConnectionsJSON: types.StringValue(`{invalid`),
		}
		diags := &diag.Diagnostics{}

		_, _, _ = parseWorkflowJSON(plan, diags)

		assert.True(t, diags.HasError())
		assert.Contains(t, diags.Errors()[0].Summary(), "Invalid connections JSON")
	})

	t.Run("parse invalid settings JSON", func(t *testing.T) {
		plan := &models.Resource{
			SettingsJSON: types.StringValue(`not valid`),
		}
		diags := &diag.Diagnostics{}

		_, _, _ = parseWorkflowJSON(plan, diags)

		assert.True(t, diags.HasError())
		assert.Contains(t, diags.Errors()[0].Summary(), "Invalid settings JSON")
	})

	t.Run("parse null JSON fields", func(t *testing.T) {
		plan := &models.Resource{
			NodesJSON:       types.StringNull(),
			ConnectionsJSON: types.StringNull(),
			SettingsJSON:    types.StringNull(),
		}
		diags := &diag.Diagnostics{}

		nodes, connections, settings := parseWorkflowJSON(plan, diags)

		if diags.HasError() {
			for _, err := range diags.Errors() {
				t.Logf("Diagnostic Error: %s - %s", err.Summary(), err.Detail())
			}
		}
		assert.False(t, diags.HasError())
		assert.Empty(t, nodes)
		assert.Empty(t, connections)
		assert.NotNil(t, settings)
	})

	t.Run("parse unknown JSON fields", func(t *testing.T) {
		plan := &models.Resource{
			NodesJSON:       types.StringUnknown(),
			ConnectionsJSON: types.StringUnknown(),
			SettingsJSON:    types.StringUnknown(),
		}
		diags := &diag.Diagnostics{}

		nodes, connections, settings := parseWorkflowJSON(plan, diags)

		if diags.HasError() {
			for _, err := range diags.Errors() {
				t.Logf("Diagnostic Error: %s - %s", err.Summary(), err.Detail())
			}
		}
		assert.False(t, diags.HasError())
		assert.Empty(t, nodes)
		assert.Empty(t, connections)
		assert.NotNil(t, settings)
	})
}

func TestMapTagsFromWorkflow(t *testing.T) {
	t.Run("map tags with IDs", func(t *testing.T) {
		workflow := &n8nsdk.Workflow{
			Tags: []n8nsdk.Tag{
				{Id: strPtr("tag1")},
				{Id: strPtr("tag2")},
				{Id: strPtr("tag3")},
			},
		}
		diags := &diag.Diagnostics{}

		tagList := mapTagsFromWorkflow(context.Background(), workflow, diags)

		if diags.HasError() {
			for _, err := range diags.Errors() {
				t.Logf("Diagnostic Error: %s - %s", err.Summary(), err.Detail())
			}
		}
		assert.False(t, diags.HasError())
		assert.False(t, tagList.IsNull())

		var tags []string
		diags.Append(tagList.ElementsAs(context.Background(), &tags, false)...)
		assert.Len(t, tags, 3)
		assert.Contains(t, tags, "tag1")
		assert.Contains(t, tags, "tag2")
		assert.Contains(t, tags, "tag3")
	})

	t.Run("map empty tags", func(t *testing.T) {
		workflow := &n8nsdk.Workflow{
			Tags: []n8nsdk.Tag{},
		}
		diags := &diag.Diagnostics{}

		tagList := mapTagsFromWorkflow(context.Background(), workflow, diags)

		if diags.HasError() {
			for _, err := range diags.Errors() {
				t.Logf("Diagnostic Error: %s - %s", err.Summary(), err.Detail())
			}
		}
		assert.False(t, diags.HasError())
		assert.False(t, tagList.IsNull())
	})

	t.Run("map tags with nil IDs", func(t *testing.T) {
		workflow := &n8nsdk.Workflow{
			Tags: []n8nsdk.Tag{
				{Id: nil},
				{Id: strPtr("valid-tag")},
			},
		}
		diags := &diag.Diagnostics{}

		tagList := mapTagsFromWorkflow(context.Background(), workflow, diags)

		if diags.HasError() {
			for _, err := range diags.Errors() {
				t.Logf("Diagnostic Error: %s - %s", err.Summary(), err.Detail())
			}
		}
		assert.False(t, diags.HasError())
		assert.False(t, tagList.IsNull())

		var tags []string
		diags.Append(tagList.ElementsAs(context.Background(), &tags, false)...)
		assert.Len(t, tags, 1)
		assert.Contains(t, tags, "valid-tag")
	})
}

func TestMapWorkflowBasicFields(t *testing.T) {
	t.Run("map all basic fields", func(t *testing.T) {
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

		assert.True(t, plan.Active.ValueBool())
		assert.Equal(t, "v1", plan.VersionID.ValueString())
		assert.False(t, plan.IsArchived.ValueBool())
		assert.Equal(t, int64(5), plan.TriggerCount.ValueInt64())
	})

	t.Run("map nil basic fields", func(t *testing.T) {
		workflow := &n8nsdk.Workflow{}
		plan := &models.Resource{}

		mapWorkflowBasicFields(workflow, plan)

		assert.True(t, plan.Active.IsNull())
		assert.True(t, plan.VersionID.IsNull())
		assert.True(t, plan.IsArchived.IsNull())
		assert.True(t, plan.TriggerCount.IsNull())
	})

	t.Run("map zero trigger count", func(t *testing.T) {
		triggerCount := int32(0)
		workflow := &n8nsdk.Workflow{
			TriggerCount: &triggerCount,
		}
		plan := &models.Resource{}

		mapWorkflowBasicFields(workflow, plan)

		assert.Equal(t, int64(0), plan.TriggerCount.ValueInt64())
	})
}

func TestMapWorkflowTimestamps(t *testing.T) {
	t.Run("map valid timestamps", func(t *testing.T) {
		createdAt := time.Now()
		updatedAt := time.Now().Add(1 * time.Hour)

		workflow := &n8nsdk.Workflow{
			CreatedAt: &createdAt,
			UpdatedAt: &updatedAt,
		}
		plan := &models.Resource{}

		mapWorkflowTimestamps(workflow, plan)

		assert.False(t, plan.CreatedAt.IsNull())
		assert.False(t, plan.UpdatedAt.IsNull())
		assert.Contains(t, plan.CreatedAt.ValueString(), createdAt.Format("2006-01-02"))
		assert.Contains(t, plan.UpdatedAt.ValueString(), updatedAt.Format("2006-01-02"))
	})

	t.Run("map nil timestamps", func(t *testing.T) {
		workflow := &n8nsdk.Workflow{}
		plan := &models.Resource{}

		mapWorkflowTimestamps(workflow, plan)

		assert.True(t, plan.CreatedAt.IsNull())
		assert.True(t, plan.UpdatedAt.IsNull())
	})
}

func TestMapWorkflowToModel(t *testing.T) {
	t.Run("map complete workflow", func(t *testing.T) {
		active := true
		id := "wf-123"
		createdAt := time.Now()
		updatedAt := time.Now()

		workflow := &n8nsdk.Workflow{
			Id:        &id,
			Name:      "Test Workflow",
			Active:    &active,
			CreatedAt: &createdAt,
			UpdatedAt: &updatedAt,
			Tags: []n8nsdk.Tag{
				{Id: strPtr("tag1")},
			},
		}
		plan := &models.Resource{}
		diags := &diag.Diagnostics{}

		mapWorkflowToModel(context.Background(), workflow, plan, diags)

		if diags.HasError() {
			for _, err := range diags.Errors() {
				t.Logf("Diagnostic Error: %s - %s", err.Summary(), err.Detail())
			}
		}
		assert.False(t, diags.HasError())
		assert.Equal(t, "Test Workflow", plan.Name.ValueString())
		assert.True(t, plan.Active.ValueBool())
		assert.False(t, plan.CreatedAt.IsNull())
		assert.False(t, plan.UpdatedAt.IsNull())
		assert.False(t, plan.Tags.IsNull())
	})

	t.Run("map workflow with metadata", func(t *testing.T) {
		workflow := &n8nsdk.Workflow{
			Name: "Test Workflow",
			Meta: map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
			},
		}
		plan := &models.Resource{}
		diags := &diag.Diagnostics{}

		mapWorkflowToModel(context.Background(), workflow, plan, diags)

		if diags.HasError() {
			for _, err := range diags.Errors() {
				t.Logf("Diagnostic Error: %s - %s", err.Summary(), err.Detail())
			}
		}
		assert.False(t, diags.HasError())
		assert.False(t, plan.Meta.IsNull())
	})

	t.Run("map workflow with pin data", func(t *testing.T) {
		workflow := &n8nsdk.Workflow{
			Name: "Test Workflow",
			PinData: map[string]interface{}{
				"node1": "data1",
			},
		}
		plan := &models.Resource{}
		diags := &diag.Diagnostics{}

		mapWorkflowToModel(context.Background(), workflow, plan, diags)

		if diags.HasError() {
			for _, err := range diags.Errors() {
				t.Logf("Diagnostic Error: %s - %s", err.Summary(), err.Detail())
			}
		}
		assert.False(t, diags.HasError())
		assert.False(t, plan.PinData.IsNull())
	})
}

func TestSerializeWorkflowJSON(t *testing.T) {
	t.Run("serialize nodes", func(t *testing.T) {
		nodeName := "Start"
		nodeType := "n8n-nodes-base.start"
		workflow := &n8nsdk.Workflow{
			Nodes: []n8nsdk.Node{
				{Name: &nodeName, Type: &nodeType},
			},
		}
		plan := &models.Resource{}

		serializeWorkflowJSON(workflow, plan)

		assert.False(t, plan.NodesJSON.IsNull())
		assert.Contains(t, plan.NodesJSON.ValueString(), "Start")
	})

	t.Run("serialize connections", func(t *testing.T) {
		workflow := &n8nsdk.Workflow{
			Connections: map[string]interface{}{
				"Start": map[string]interface{}{
					"main": []interface{}{},
				},
			},
		}
		plan := &models.Resource{}

		serializeWorkflowJSON(workflow, plan)

		assert.False(t, plan.ConnectionsJSON.IsNull())
		assert.Contains(t, plan.ConnectionsJSON.ValueString(), "Start")
	})

	t.Run("serialize settings", func(t *testing.T) {
		workflow := &n8nsdk.Workflow{
			Settings: n8nsdk.WorkflowSettings{
				SaveDataErrorExecution:   strPtr("all"),
				SaveDataSuccessExecution: strPtr("all"),
			},
		}
		plan := &models.Resource{}

		serializeWorkflowJSON(workflow, plan)

		assert.False(t, plan.SettingsJSON.IsNull())
	})

	t.Run("serialize nil fields", func(t *testing.T) {
		workflow := &n8nsdk.Workflow{}
		plan := &models.Resource{}

		serializeWorkflowJSON(workflow, plan)

		assert.True(t, plan.NodesJSON.IsNull())
		assert.True(t, plan.ConnectionsJSON.IsNull())
		assert.False(t, plan.SettingsJSON.IsNull()) // Settings always serializes
	})
}

func TestConvertTagIDsToTagIdsInner(t *testing.T) {
	t.Run("convert tag IDs", func(t *testing.T) {
		tagIDs := []string{"tag1", "tag2", "tag3"}

		result := convertTagIDsToTagIdsInner(tagIDs)

		assert.Len(t, result, 3)
		assert.Equal(t, "tag1", result[0].Id)
		assert.Equal(t, "tag2", result[1].Id)
		assert.Equal(t, "tag3", result[2].Id)
	})

	t.Run("convert empty tag IDs", func(t *testing.T) {
		tagIDs := []string{}

		result := convertTagIDsToTagIdsInner(tagIDs)

		assert.Empty(t, result)
	})

	t.Run("convert single tag ID", func(t *testing.T) {
		tagIDs := []string{"single-tag"}

		result := convertTagIDsToTagIdsInner(tagIDs)

		assert.Len(t, result, 1)
		assert.Equal(t, "single-tag", result[0].Id)
	})
}

func BenchmarkParseWorkflowJSON(b *testing.B) {
	plan := &models.Resource{
		NodesJSON:       types.StringValue(`[{"name":"Start","type":"n8n-nodes-base.start"}]`),
		ConnectionsJSON: types.StringValue(`{"Start":{"main":[[]]}}`),
		SettingsJSON:    types.StringValue(`{"saveDataErrorExecution":"all"}`),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		diags := &diag.Diagnostics{}
		_, _, _ = parseWorkflowJSON(plan, diags)
	}
}

func BenchmarkMapWorkflowToModel(b *testing.B) {
	active := true
	id := "wf-123"
	createdAt := time.Now()

	workflow := &n8nsdk.Workflow{
		Id:        &id,
		Name:      "Benchmark Workflow",
		Active:    &active,
		CreatedAt: &createdAt,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		plan := &models.Resource{}
		diags := &diag.Diagnostics{}
		mapWorkflowToModel(context.Background(), workflow, plan, diags)
	}
}

func BenchmarkSerializeWorkflowJSON(b *testing.B) {
	nodeName := "Start"
	nodeType := "n8n-nodes-base.start"
	workflow := &n8nsdk.Workflow{
		Nodes: []n8nsdk.Node{
			{Name: &nodeName, Type: &nodeType},
		},
		Connections: map[string]interface{}{
			"Start": map[string]interface{}{},
		},
		Settings: n8nsdk.WorkflowSettings{},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		plan := &models.Resource{}
		serializeWorkflowJSON(workflow, plan)
	}
}

func BenchmarkConvertTagIDsToTagIdsInner(b *testing.B) {
	tagIDs := []string{"tag1", "tag2", "tag3", "tag4", "tag5"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = convertTagIDsToTagIdsInner(tagIDs)
	}
}

func TestHandleWorkflowActivation(t *testing.T) {
	t.Run("activate workflow when active changes to true", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/workflows/wf-123/activate" && r.Method == http.MethodPost {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"id":          "wf-123",
					"name":        "Test Workflow",
					"active":      true,
					"nodes":       []interface{}{},
					"connections": map[string]interface{}{},
					"settings":    map[string]interface{}{},
				})
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupHelperTestClient(t, handler)
		defer server.Close()

		r := &WorkflowResource{client: n8nClient}

		plan := &models.Resource{
			ID:     types.StringValue("wf-123"),
			Active: types.BoolValue(true),
		}
		state := &models.Resource{
			ID:     types.StringValue("wf-123"),
			Active: types.BoolValue(false),
		}
		diags := &diag.Diagnostics{}

		r.handleWorkflowActivation(context.Background(), plan, state, diags)

		if diags.HasError() {
			for _, err := range diags.Errors() {
				t.Logf("Diagnostic Error: %s - %s", err.Summary(), err.Detail())
			}
		}
		assert.False(t, diags.HasError(), "Activation should not have errors")
	})

	t.Run("deactivate workflow when active changes to false", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/workflows/wf-123/deactivate" && r.Method == http.MethodPost {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"id":          "wf-123",
					"name":        "Test Workflow",
					"active":      false,
					"nodes":       []interface{}{},
					"connections": map[string]interface{}{},
					"settings":    map[string]interface{}{},
				})
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupHelperTestClient(t, handler)
		defer server.Close()

		r := &WorkflowResource{client: n8nClient}

		plan := &models.Resource{
			ID:     types.StringValue("wf-123"),
			Active: types.BoolValue(false),
		}
		state := &models.Resource{
			ID:     types.StringValue("wf-123"),
			Active: types.BoolValue(true),
		}
		diags := &diag.Diagnostics{}

		r.handleWorkflowActivation(context.Background(), plan, state, diags)

		if diags.HasError() {
			for _, err := range diags.Errors() {
				t.Logf("Diagnostic Error: %s - %s", err.Summary(), err.Detail())
			}
		}
		assert.False(t, diags.HasError(), "Deactivation should not have errors")
	})

	t.Run("no change when active status is same", func(t *testing.T) {
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

		if diags.HasError() {
			for _, err := range diags.Errors() {
				t.Logf("Diagnostic Error: %s - %s", err.Summary(), err.Detail())
			}
		}
		assert.False(t, diags.HasError(), "No change should not have errors")
	})

	t.Run("error when activation fails", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message": "Activation failed"}`))
		})

		n8nClient, server := setupHelperTestClient(t, handler)
		defer server.Close()

		r := &WorkflowResource{client: n8nClient}

		plan := &models.Resource{
			ID:     types.StringValue("wf-123"),
			Active: types.BoolValue(true),
		}
		state := &models.Resource{
			ID:     types.StringValue("wf-123"),
			Active: types.BoolValue(false),
		}
		diags := &diag.Diagnostics{}

		r.handleWorkflowActivation(context.Background(), plan, state, diags)

		assert.True(t, diags.HasError(), "Activation should have errors")
		assert.Contains(t, diags.Errors()[0].Summary(), "Error changing workflow activation status")
	})

	t.Run("no change when active is null", func(t *testing.T) {
		r := &WorkflowResource{}

		plan := &models.Resource{
			ID:     types.StringValue("wf-123"),
			Active: types.BoolNull(),
		}
		state := &models.Resource{
			ID:     types.StringValue("wf-123"),
			Active: types.BoolValue(false),
		}
		diags := &diag.Diagnostics{}

		r.handleWorkflowActivation(context.Background(), plan, state, diags)

		if diags.HasError() {
			for _, err := range diags.Errors() {
				t.Logf("Diagnostic Error: %s - %s", err.Summary(), err.Detail())
			}
		}
		assert.False(t, diags.HasError(), "Null active should not cause errors")
	})
}

func TestUpdateWorkflowTags(t *testing.T) {
	t.Run("update tags successfully", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/workflows/wf-123/tags" && r.Method == http.MethodPut {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode([]interface{}{
					map[string]interface{}{"id": "tag1", "name": "Tag 1"},
					map[string]interface{}{"id": "tag2", "name": "Tag 2"},
				})
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupHelperTestClient(t, handler)
		defer server.Close()

		r := &WorkflowResource{client: n8nClient}

		tags, _ := types.ListValueFrom(context.Background(), types.StringType, []string{"tag1", "tag2"})
		plan := &models.Resource{
			Tags: tags,
		}
		workflow := &n8nsdk.Workflow{}
		diags := &diag.Diagnostics{}

		r.updateWorkflowTags(context.Background(), "wf-123", plan, workflow, diags)

		if diags.HasError() {
			for _, err := range diags.Errors() {
				t.Logf("Diagnostic Error: %s - %s", err.Summary(), err.Detail())
			}
		}
		assert.False(t, diags.HasError(), "Update tags should not have errors")
	})

	t.Run("skip when tags is null", func(t *testing.T) {
		r := &WorkflowResource{}

		plan := &models.Resource{
			Tags: types.ListNull(types.StringType),
		}
		workflow := &n8nsdk.Workflow{}
		diags := &diag.Diagnostics{}

		r.updateWorkflowTags(context.Background(), "wf-123", plan, workflow, diags)

		if diags.HasError() {
			for _, err := range diags.Errors() {
				t.Logf("Diagnostic Error: %s - %s", err.Summary(), err.Detail())
			}
		}
		assert.False(t, diags.HasError(), "Null tags should not cause errors")
	})

	t.Run("skip when tags is unknown", func(t *testing.T) {
		r := &WorkflowResource{}

		plan := &models.Resource{
			Tags: types.ListUnknown(types.StringType),
		}
		workflow := &n8nsdk.Workflow{}
		diags := &diag.Diagnostics{}

		r.updateWorkflowTags(context.Background(), "wf-123", plan, workflow, diags)

		if diags.HasError() {
			for _, err := range diags.Errors() {
				t.Logf("Diagnostic Error: %s - %s", err.Summary(), err.Detail())
			}
		}
		assert.False(t, diags.HasError(), "Unknown tags should not cause errors")
	})

	t.Run("error when tag update fails", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"message": "Invalid tags"}`))
		})

		n8nClient, server := setupHelperTestClient(t, handler)
		defer server.Close()

		r := &WorkflowResource{client: n8nClient}

		tags, _ := types.ListValueFrom(context.Background(), types.StringType, []string{"tag1"})
		plan := &models.Resource{
			Tags: tags,
		}
		workflow := &n8nsdk.Workflow{}
		diags := &diag.Diagnostics{}

		r.updateWorkflowTags(context.Background(), "wf-123", plan, workflow, diags)

		assert.True(t, diags.HasError(), "Failed tag update should have errors")
	})

	t.Run("update with empty tag list", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/workflows/wf-123/tags" && r.Method == http.MethodPut {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode([]interface{}{})
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		n8nClient, server := setupHelperTestClient(t, handler)
		defer server.Close()

		r := &WorkflowResource{client: n8nClient}

		tags, _ := types.ListValueFrom(context.Background(), types.StringType, []string{})
		plan := &models.Resource{
			Tags: tags,
		}
		workflow := &n8nsdk.Workflow{}
		diags := &diag.Diagnostics{}

		r.updateWorkflowTags(context.Background(), "wf-123", plan, workflow, diags)

		if diags.HasError() {
			for _, err := range diags.Errors() {
				t.Logf("Diagnostic Error: %s - %s", err.Summary(), err.Detail())
			}
		}
		assert.False(t, diags.HasError(), "Empty tags should not cause errors")
	})
}

// setupHelperTestClient creates a test N8nClient with httptest server for helper tests.
func setupHelperTestClient(t *testing.T, handler http.HandlerFunc) (*client.N8nClient, *httptest.Server) {
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
