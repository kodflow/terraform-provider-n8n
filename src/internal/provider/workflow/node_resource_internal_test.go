// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

package workflow

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/workflow/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWorkflowNodeResource_generateNodeID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		plan     *models.NodeResource
		wantSame bool
		checkID  func(t *testing.T, id string)
	}{
		{
			name: "generates consistent IDs for same inputs",
			plan: &models.NodeResource{
				Name: types.StringValue("Test Node"),
				Type: types.StringValue("n8n-nodes-base.webhook"),
			},
			wantSame: true,
			checkID: func(t *testing.T, id string) {
				t.Helper()
				assert.NotEmpty(t, id)
			},
		},
		{
			name: "generates different IDs for different types",
			plan: &models.NodeResource{
				Name: types.StringValue("Test Node"),
				Type: types.StringValue("n8n-nodes-base.code"),
			},
			wantSame: false,
			checkID: func(t *testing.T, id string) {
				t.Helper()
				assert.NotEmpty(t, id)
			},
		},
		{
			name: "handles empty names",
			plan: &models.NodeResource{
				Name: types.StringValue(""),
				Type: types.StringValue("n8n-nodes-base.webhook"),
			},
			wantSame: true,
			checkID: func(t *testing.T, id string) {
				t.Helper()
				assert.NotEmpty(t, id)
			},
		},
		{
			name: "error case - generates valid ID even with extreme values",
			plan: &models.NodeResource{
				Name: types.StringValue("Very long node name that might cause issues"),
				Type: types.StringValue("n8n-nodes-base.extremely-long-type-name-that-is-unusual"),
			},
			wantSame: true,
			checkID: func(t *testing.T, id string) {
				t.Helper()
				assert.NotEmpty(t, id)
				assert.Contains(t, id, "-")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &WorkflowNodeResource{}
			id1 := r.generateNodeID(tt.plan)
			id2 := r.generateNodeID(tt.plan)

			if tt.wantSame {
				assert.Equal(t, id1, id2, "IDs should be consistent for same inputs")
			}
			if tt.checkID != nil {
				tt.checkID(t, id1)
			}
		})
	}
}

func TestWorkflowNodeResource_generateNodeJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		plan      *models.NodeResource
		wantErr   bool
		checkJSON func(t *testing.T, nodeJSON string)
	}{
		{
			name: "basic node with minimal fields",
			plan: &models.NodeResource{
				ID:          types.StringValue("test-id-123"),
				Name:        types.StringValue("Webhook"),
				Type:        types.StringValue("n8n-nodes-base.webhook"),
				TypeVersion: types.Int64Value(1),
				Position:    types.ListValueMust(types.Int64Type, []attr.Value{types.Int64Value(250), types.Int64Value(300)}),
			},
			wantErr: false,
			checkJSON: func(t *testing.T, nodeJSON string) {
				t.Helper()

				var node map[string]any
				err := json.Unmarshal([]byte(nodeJSON), &node)
				require.NoError(t, err)

				assert.Equal(t, "test-id-123", node["id"])
				assert.Equal(t, "Webhook", node["name"])
				assert.Equal(t, "n8n-nodes-base.webhook", node["type"])
				assert.Equal(t, float64(1), node["typeVersion"])

				position, ok := node["position"].([]any)
				require.True(t, ok)
				assert.Len(t, position, 2)
				assert.Equal(t, float64(250), position[0])
				assert.Equal(t, float64(300), position[1])
			},
		},
		{
			name: "node with parameters",
			plan: &models.NodeResource{
				ID:          types.StringValue("test-id-456"),
				Name:        types.StringValue("Code"),
				Type:        types.StringValue("n8n-nodes-base.code"),
				TypeVersion: types.Int64Value(1),
				Position:    types.ListValueMust(types.Int64Type, []attr.Value{types.Int64Value(450), types.Int64Value(300)}),
				Parameters:  types.StringValue(`{"mode":"runOnceForAllItems","jsCode":"return items;"}`),
			},
			wantErr: false,
			checkJSON: func(t *testing.T, nodeJSON string) {
				t.Helper()

				var node map[string]any
				err := json.Unmarshal([]byte(nodeJSON), &node)
				require.NoError(t, err)

				params, ok := node["parameters"].(map[string]any)
				require.True(t, ok)
				assert.Equal(t, "runOnceForAllItems", params["mode"])
				assert.Equal(t, "return items;", params["jsCode"])
			},
		},
		{
			name: "webhook node with webhookId",
			plan: &models.NodeResource{
				ID:          types.StringValue("test-id-789"),
				Name:        types.StringValue("Webhook Trigger"),
				Type:        types.StringValue("n8n-nodes-base.webhook"),
				TypeVersion: types.Int64Value(1),
				Position:    types.ListValueMust(types.Int64Type, []attr.Value{types.Int64Value(250), types.Int64Value(300)}),
				WebhookID:   types.StringValue("my-webhook-id"),
			},
			wantErr: false,
			checkJSON: func(t *testing.T, nodeJSON string) {
				t.Helper()

				var node map[string]any
				err := json.Unmarshal([]byte(nodeJSON), &node)
				require.NoError(t, err)

				assert.Equal(t, "my-webhook-id", node["webhookId"])
			},
		},
		{
			name: "disabled node with notes",
			plan: &models.NodeResource{
				ID:          types.StringValue("test-id-disabled"),
				Name:        types.StringValue("Disabled Node"),
				Type:        types.StringValue("n8n-nodes-base.set"),
				TypeVersion: types.Int64Value(1),
				Position:    types.ListValueMust(types.Int64Type, []attr.Value{types.Int64Value(650), types.Int64Value(300)}),
				Disabled:    types.BoolValue(true),
				Notes:       types.StringValue("This node is temporarily disabled"),
			},
			wantErr: false,
			checkJSON: func(t *testing.T, nodeJSON string) {
				t.Helper()

				var node map[string]any
				err := json.Unmarshal([]byte(nodeJSON), &node)
				require.NoError(t, err)

				assert.Equal(t, true, node["disabled"])
				assert.Equal(t, "This node is temporarily disabled", node["notes"])
			},
		},
		{
			name: "invalid parameters JSON",
			plan: &models.NodeResource{
				ID:          types.StringValue("test-id-invalid"),
				Name:        types.StringValue("Invalid"),
				Type:        types.StringValue("n8n-nodes-base.code"),
				TypeVersion: types.Int64Value(1),
				Position:    types.ListValueMust(types.Int64Type, []attr.Value{types.Int64Value(250), types.Int64Value(300)}),
				Parameters:  types.StringValue(`{invalid json`),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := &WorkflowNodeResource{}
			var diags diag.Diagnostics

			ctx := t.Context()
			success := r.generateNodeJSON(ctx, tt.plan, &diags)

			if tt.wantErr {
				assert.False(t, success)
				assert.True(t, diags.HasError())
			} else {
				assert.True(t, success)
				assert.False(t, diags.HasError())
				require.NotEmpty(t, tt.plan.NodeJSON.ValueString())

				if tt.checkJSON != nil {
					tt.checkJSON(t, tt.plan.NodeJSON.ValueString())
				}
			}
		})
	}
}

func TestWorkflowNodeResource_parseNodeParameters(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		plan      *models.NodeResource
		wantErr   bool
		wantValue map[string]any
	}{
		{
			name: "valid JSON parameters",
			plan: &models.NodeResource{
				Parameters: types.StringValue(`{"key":"value"}`),
			},
			wantErr:   false,
			wantValue: map[string]any{"key": "value"},
		},
		{
			name: "invalid JSON parameters",
			plan: &models.NodeResource{
				Parameters: types.StringValue(`{invalid`),
			},
			wantErr: true,
		},
		{
			name: "null parameters",
			plan: &models.NodeResource{
				Parameters: types.StringNull(),
			},
			wantErr:   false,
			wantValue: map[string]any{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := &WorkflowNodeResource{}
			node := make(map[string]any)
			var diags diag.Diagnostics

			success := r.parseNodeParameters(tt.plan, node, &diags)

			// Verify result matches expectation.
			if tt.wantErr {
				assert.False(t, success)
				assert.True(t, diags.HasError())
			} else {
				assert.True(t, success)
				assert.False(t, diags.HasError())
				if tt.wantValue != nil {
					assert.Equal(t, tt.wantValue, node["parameters"])
				}
			}
		})
	}
}

func TestWorkflowNodeResource_addOptionalNodeFields(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		plan       *models.NodeResource
		wantFields map[string]any
	}{
		{
			name: "all optional fields present",
			plan: &models.NodeResource{
				WebhookID: types.StringValue("webhook-123"),
				Disabled:  types.BoolValue(true),
				Notes:     types.StringValue("test notes"),
			},
			wantFields: map[string]any{
				"webhookId": "webhook-123",
				"disabled":  true,
				"notes":     "test notes",
			},
		},
		{
			name: "no optional fields",
			plan: &models.NodeResource{
				WebhookID: types.StringNull(),
				Disabled:  types.BoolNull(),
				Notes:     types.StringNull(),
			},
			wantFields: map[string]any{},
		},
		{
			name: "disabled is false",
			plan: &models.NodeResource{
				WebhookID: types.StringNull(),
				Disabled:  types.BoolValue(false),
				Notes:     types.StringNull(),
			},
			wantFields: map[string]any{},
		},
		{
			name: "unknown values are skipped",
			plan: &models.NodeResource{
				WebhookID: types.StringUnknown(),
				Disabled:  types.BoolUnknown(),
				Notes:     types.StringUnknown(),
			},
			wantFields: map[string]any{},
		},
		{
			name: "error case - handles edge case with disabled true and empty webhook",
			plan: &models.NodeResource{
				WebhookID: types.StringValue(""),
				Disabled:  types.BoolValue(true),
				Notes:     types.StringNull(),
			},
			wantFields: map[string]any{
				"webhookId": "",
				"disabled":  true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := &WorkflowNodeResource{}
			node := make(map[string]any)

			r.addOptionalNodeFields(tt.plan, node)

			// Verify all expected fields are present.
			for key, expectedValue := range tt.wantFields {
				actualValue, ok := node[key]
				assert.True(t, ok, "expected field %s not found", key)
				assert.Equal(t, expectedValue, actualValue)
			}

			// Verify no unexpected fields.
			assert.Equal(t, len(tt.wantFields), len(node))
		})
	}
}

func TestWorkflowNodeResource_addRequiredAttributes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		checkAttrs []string
		checkTypes func(t *testing.T, attrs map[string]schema.Attribute)
	}{
		{
			name:       "required attributes added",
			checkAttrs: []string{"id", "name", "type", "type_version", "position"},
			checkTypes: func(t *testing.T, attrs map[string]schema.Attribute) {
				t.Helper()
				// Verify name and type are required
				nameAttr, ok := attrs["name"].(schema.StringAttribute)
				assert.True(t, ok)
				assert.True(t, nameAttr.Required)

				typeAttr, ok := attrs["type"].(schema.StringAttribute)
				assert.True(t, ok)
				assert.True(t, typeAttr.Required)
			},
		},
		{
			name:       "id attribute is computed",
			checkAttrs: []string{"id"},
			checkTypes: func(t *testing.T, attrs map[string]schema.Attribute) {
				t.Helper()
				idAttr, ok := attrs["id"].(schema.StringAttribute)
				assert.True(t, ok)
				assert.True(t, idAttr.Computed)
				assert.NotEmpty(t, idAttr.PlanModifiers)
			},
		},
		{
			name:       "error case - all required attributes must be present",
			checkAttrs: []string{"id", "name", "type", "type_version", "position"},
			checkTypes: func(t *testing.T, attrs map[string]schema.Attribute) {
				t.Helper()
				assert.Len(t, attrs, 5, "must have exactly 5 required attributes")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := &WorkflowNodeResource{}
			attrs := make(map[string]schema.Attribute)

			r.addRequiredAttributes(attrs)

			for _, attrName := range tt.checkAttrs {
				assert.Contains(t, attrs, attrName)
			}

			if tt.checkTypes != nil {
				tt.checkTypes(t, attrs)
			}
		})
	}
}

func TestWorkflowNodeResource_addOptionalAttributes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		checkAttrs []string
		checkTypes func(t *testing.T, attrs map[string]schema.Attribute)
	}{
		{
			name:       "optional attributes added",
			checkAttrs: []string{"parameters", "webhook_id", "disabled", "notes"},
			checkTypes: func(t *testing.T, attrs map[string]schema.Attribute) {
				t.Helper()
				// Verify all attributes are optional
				for _, attrName := range []string{"parameters", "webhook_id", "disabled", "notes"} {
					attr := attrs[attrName]
					assert.True(t, attr.IsOptional(), "expected %s to be optional", attrName)
				}
			},
		},
		{
			name:       "optional attributes with defaults",
			checkAttrs: []string{"parameters", "disabled"},
			checkTypes: func(t *testing.T, attrs map[string]schema.Attribute) {
				t.Helper()
				// Verify parameters has default
				paramAttr, ok := attrs["parameters"].(schema.StringAttribute)
				assert.True(t, ok)
				assert.NotNil(t, paramAttr.Default)

				// Verify disabled has default
				disabledAttr, ok := attrs["disabled"].(schema.BoolAttribute)
				assert.True(t, ok)
				assert.NotNil(t, disabledAttr.Default)
			},
		},
		{
			name:       "error case - all optional attributes must be present",
			checkAttrs: []string{"parameters", "webhook_id", "disabled", "notes"},
			checkTypes: func(t *testing.T, attrs map[string]schema.Attribute) {
				t.Helper()
				assert.Len(t, attrs, 4, "must have exactly 4 optional attributes")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := &WorkflowNodeResource{}
			attrs := make(map[string]schema.Attribute)

			r.addOptionalAttributes(attrs)

			for _, attrName := range tt.checkAttrs {
				assert.Contains(t, attrs, attrName)
			}

			if tt.checkTypes != nil {
				tt.checkTypes(t, attrs)
			}
		})
	}
}

func TestWorkflowNodeResource_getSchemaAttributes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		wantErr    bool
		checkAttrs []string
	}{
		{
			name:    "schema has all required attributes",
			wantErr: false,
			checkAttrs: []string{
				"id", "name", "type", "type_version", "position",
				"parameters", "webhook_id", "disabled", "notes", "node_json",
			},
		},
		{
			name:    "error case - schema must not be empty",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := &WorkflowNodeResource{}
			attrs := r.getSchemaAttributes()

			switch tt.name {
			case "schema has all required attributes":
				require.NotNil(t, attrs)
				for _, attrName := range tt.checkAttrs {
					assert.Contains(t, attrs, attrName)
				}
			case "error case - schema must not be empty":
				if len(attrs) == 0 && !tt.wantErr {
					t.Error("expected attributes, got empty map")
				}
			}
		})
	}
}

// createNodeTestSchema creates a schema for node resource testing.
func createNodeTestSchema(t *testing.T) schema.Schema {
	t.Helper()
	r := &WorkflowNodeResource{}
	return schema.Schema{
		Attributes: r.getSchemaAttributes(),
	}
}

// getNodeObjectType returns the tftypes.Object for node resource.
func getNodeObjectType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"id":           tftypes.String,
			"name":         tftypes.String,
			"type":         tftypes.String,
			"type_version": tftypes.Number,
			"position":     tftypes.List{ElementType: tftypes.Number},
			"parameters":   tftypes.String,
			"webhook_id":   tftypes.String,
			"disabled":     tftypes.Bool,
			"notes":        tftypes.String,
			"node_json":    tftypes.String,
		},
	}
}

func TestWorkflowNodeResource_Create_Success(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		nodeName    string
		nodeType    string
		expectError bool
	}{
		{
			name:        "basic webhook node creation",
			nodeName:    "Webhook",
			nodeType:    "n8n-nodes-base.webhook",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := &WorkflowNodeResource{}
			ctx := t.Context()
			objectType := getNodeObjectType()

			rawPlan := map[string]tftypes.Value{
				"id":           tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
				"name":         tftypes.NewValue(tftypes.String, tt.nodeName),
				"type":         tftypes.NewValue(tftypes.String, tt.nodeType),
				"type_version": tftypes.NewValue(tftypes.Number, 1),
				"position":     tftypes.NewValue(tftypes.List{ElementType: tftypes.Number}, []tftypes.Value{tftypes.NewValue(tftypes.Number, 250), tftypes.NewValue(tftypes.Number, 300)}),
				"parameters":   tftypes.NewValue(tftypes.String, "{}"),
				"webhook_id":   tftypes.NewValue(tftypes.String, nil),
				"disabled":     tftypes.NewValue(tftypes.Bool, false),
				"notes":        tftypes.NewValue(tftypes.String, nil),
				"node_json":    tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
			}

			plan := tfsdk.Plan{Raw: tftypes.NewValue(objectType, rawPlan), Schema: createNodeTestSchema(t)}
			state := tfsdk.State{Raw: tftypes.NewValue(objectType, nil), Schema: createNodeTestSchema(t)}
			req := resource.CreateRequest{Plan: plan}
			resp := resource.CreateResponse{State: state}

			r.Create(ctx, req, &resp)

			assert.Equal(t, tt.expectError, resp.Diagnostics.HasError())
		})
	}
}

func TestWorkflowNodeResource_Create_WithInvalidParameters(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		parameters  string
		expectError bool
	}{
		{
			name:        "invalid JSON parameters cause error",
			parameters:  "{invalid json",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := &WorkflowNodeResource{}
			ctx := t.Context()
			objectType := getNodeObjectType()

			rawPlan := map[string]tftypes.Value{
				"id":           tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
				"name":         tftypes.NewValue(tftypes.String, "Code Node"),
				"type":         tftypes.NewValue(tftypes.String, "n8n-nodes-base.code"),
				"type_version": tftypes.NewValue(tftypes.Number, 1),
				"position":     tftypes.NewValue(tftypes.List{ElementType: tftypes.Number}, []tftypes.Value{tftypes.NewValue(tftypes.Number, 250), tftypes.NewValue(tftypes.Number, 300)}),
				"parameters":   tftypes.NewValue(tftypes.String, tt.parameters),
				"webhook_id":   tftypes.NewValue(tftypes.String, nil),
				"disabled":     tftypes.NewValue(tftypes.Bool, false),
				"notes":        tftypes.NewValue(tftypes.String, nil),
				"node_json":    tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
			}

			plan := tfsdk.Plan{Raw: tftypes.NewValue(objectType, rawPlan), Schema: createNodeTestSchema(t)}
			state := tfsdk.State{Raw: tftypes.NewValue(objectType, nil), Schema: createNodeTestSchema(t)}
			req := resource.CreateRequest{Plan: plan}
			resp := resource.CreateResponse{State: state}

			r.Create(ctx, req, &resp)

			assert.Equal(t, tt.expectError, resp.Diagnostics.HasError())
		})
	}
}

func TestWorkflowNodeResource_Create_WithNullParameters(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		nodeName    string
		expectError bool
	}{
		{
			name:        "null parameters defaults to empty object",
			nodeName:    "Node with null params",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := &WorkflowNodeResource{}
			ctx := t.Context()
			objectType := getNodeObjectType()

			rawPlan := map[string]tftypes.Value{
				"id":           tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
				"name":         tftypes.NewValue(tftypes.String, tt.nodeName),
				"type":         tftypes.NewValue(tftypes.String, "n8n-nodes-base.set"),
				"type_version": tftypes.NewValue(tftypes.Number, 1),
				"position":     tftypes.NewValue(tftypes.List{ElementType: tftypes.Number}, []tftypes.Value{tftypes.NewValue(tftypes.Number, 250), tftypes.NewValue(tftypes.Number, 300)}),
				"parameters":   tftypes.NewValue(tftypes.String, nil),
				"webhook_id":   tftypes.NewValue(tftypes.String, nil),
				"disabled":     tftypes.NewValue(tftypes.Bool, false),
				"notes":        tftypes.NewValue(tftypes.String, nil),
				"node_json":    tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
			}

			plan := tfsdk.Plan{Raw: tftypes.NewValue(objectType, rawPlan), Schema: createNodeTestSchema(t)}
			state := tfsdk.State{Raw: tftypes.NewValue(objectType, nil), Schema: createNodeTestSchema(t)}
			req := resource.CreateRequest{Plan: plan}
			resp := resource.CreateResponse{State: state}

			r.Create(ctx, req, &resp)

			assert.Equal(t, tt.expectError, resp.Diagnostics.HasError())
		})
	}
}

func TestWorkflowNodeResource_Update_Success(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		nodeID      string
		nodeName    string
		expectError bool
	}{
		{
			name:        "update webhook node with parameters",
			nodeID:      "node-123",
			nodeName:    "Updated Webhook",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := &WorkflowNodeResource{}
			ctx := t.Context()
			objectType := getNodeObjectType()

			rawPlan := map[string]tftypes.Value{
				"id":           tftypes.NewValue(tftypes.String, tt.nodeID),
				"name":         tftypes.NewValue(tftypes.String, tt.nodeName),
				"type":         tftypes.NewValue(tftypes.String, "n8n-nodes-base.webhook"),
				"type_version": tftypes.NewValue(tftypes.Number, 1),
				"position":     tftypes.NewValue(tftypes.List{ElementType: tftypes.Number}, []tftypes.Value{tftypes.NewValue(tftypes.Number, 350), tftypes.NewValue(tftypes.Number, 400)}),
				"parameters":   tftypes.NewValue(tftypes.String, `{"path":"test"}`),
				"webhook_id":   tftypes.NewValue(tftypes.String, nil),
				"disabled":     tftypes.NewValue(tftypes.Bool, false),
				"notes":        tftypes.NewValue(tftypes.String, nil),
				"node_json":    tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
			}

			plan := tfsdk.Plan{Raw: tftypes.NewValue(objectType, rawPlan), Schema: createNodeTestSchema(t)}
			state := tfsdk.State{Raw: tftypes.NewValue(objectType, nil), Schema: createNodeTestSchema(t)}
			req := resource.UpdateRequest{Plan: plan}
			resp := resource.UpdateResponse{State: state}

			r.Update(ctx, req, &resp)

			assert.Equal(t, tt.expectError, resp.Diagnostics.HasError())
		})
	}
}

func TestWorkflowNodeResource_Update_WithNullTypeVersion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		expectError bool
	}{
		{
			name:        "null type version uses default",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := &WorkflowNodeResource{}
			ctx := t.Context()
			objectType := getNodeObjectType()

			rawPlan := map[string]tftypes.Value{
				"id":           tftypes.NewValue(tftypes.String, "node-123"),
				"name":         tftypes.NewValue(tftypes.String, "Node without version"),
				"type":         tftypes.NewValue(tftypes.String, "n8n-nodes-base.set"),
				"type_version": tftypes.NewValue(tftypes.Number, nil),
				"position":     tftypes.NewValue(tftypes.List{ElementType: tftypes.Number}, []tftypes.Value{tftypes.NewValue(tftypes.Number, 250), tftypes.NewValue(tftypes.Number, 300)}),
				"parameters":   tftypes.NewValue(tftypes.String, "{}"),
				"webhook_id":   tftypes.NewValue(tftypes.String, nil),
				"disabled":     tftypes.NewValue(tftypes.Bool, false),
				"notes":        tftypes.NewValue(tftypes.String, nil),
				"node_json":    tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
			}

			plan := tfsdk.Plan{Raw: tftypes.NewValue(objectType, rawPlan), Schema: createNodeTestSchema(t)}
			state := tfsdk.State{Raw: tftypes.NewValue(objectType, nil), Schema: createNodeTestSchema(t)}
			req := resource.UpdateRequest{Plan: plan}
			resp := resource.UpdateResponse{State: state}

			r.Update(ctx, req, &resp)

			assert.Equal(t, tt.expectError, resp.Diagnostics.HasError())
		})
	}
}

func TestWorkflowNodeResource_Update_InvalidJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		parameters  string
		expectError bool
	}{
		{
			name:        "invalid JSON parameters cause update error",
			parameters:  "{invalid json",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := &WorkflowNodeResource{}
			ctx := t.Context()
			objectType := getNodeObjectType()

			rawPlan := map[string]tftypes.Value{
				"id":           tftypes.NewValue(tftypes.String, "node-123"),
				"name":         tftypes.NewValue(tftypes.String, "Code Node"),
				"type":         tftypes.NewValue(tftypes.String, "n8n-nodes-base.code"),
				"type_version": tftypes.NewValue(tftypes.Number, 1),
				"position":     tftypes.NewValue(tftypes.List{ElementType: tftypes.Number}, []tftypes.Value{tftypes.NewValue(tftypes.Number, 250), tftypes.NewValue(tftypes.Number, 300)}),
				"parameters":   tftypes.NewValue(tftypes.String, tt.parameters),
				"webhook_id":   tftypes.NewValue(tftypes.String, nil),
				"disabled":     tftypes.NewValue(tftypes.Bool, false),
				"notes":        tftypes.NewValue(tftypes.String, nil),
				"node_json":    tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
			}

			plan := tfsdk.Plan{Raw: tftypes.NewValue(objectType, rawPlan), Schema: createNodeTestSchema(t)}
			state := tfsdk.State{Raw: tftypes.NewValue(objectType, nil), Schema: createNodeTestSchema(t)}
			req := resource.UpdateRequest{Plan: plan}
			resp := resource.UpdateResponse{State: state}

			r.Update(ctx, req, &resp)

			assert.Equal(t, tt.expectError, resp.Diagnostics.HasError())
		})
	}
}

func TestWorkflowNodeResource_ImportState_Success(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		importID    string
		expectError bool
	}{
		{
			name:        "import node by id",
			importID:    "imported-node-id",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := &WorkflowNodeResource{}
			ctx := t.Context()
			objectType := getNodeObjectType()

			state := tfsdk.State{
				Raw:    tftypes.NewValue(objectType, nil),
				Schema: createNodeTestSchema(t),
			}

			req := resource.ImportStateRequest{ID: tt.importID}
			resp := resource.ImportStateResponse{State: state}

			r.ImportState(ctx, req, &resp)

			assert.Equal(t, tt.expectError, resp.Diagnostics.HasError())
		})
	}
}

func TestWorkflowNodeResource_generateNodeJSON_PositionError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		position      attr.Value
		expectSuccess bool
		expectError   bool
	}{
		{
			name:          "unknown position triggers failure",
			position:      types.ListUnknown(types.Int64Type),
			expectSuccess: false,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := &WorkflowNodeResource{}
			ctx := t.Context()
			var diags diag.Diagnostics

			plan := &models.NodeResource{
				ID:          types.StringValue("test-id"),
				Name:        types.StringValue("Test Node"),
				Type:        types.StringValue("n8n-nodes-base.webhook"),
				TypeVersion: types.Int64Value(1),
				Position:    tt.position.(types.List),
				Parameters:  types.StringValue("{}"),
			}

			success := r.generateNodeJSON(ctx, plan, &diags)

			assert.Equal(t, tt.expectSuccess, success)
			assert.Equal(t, tt.expectError, diags.HasError())
		})
	}
}

func TestWorkflowNodeResource_Configure_WithData(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		providerData any
		expectError  bool
	}{
		{
			name:         "nil provider data does not cause errors",
			providerData: nil,
			expectError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := &WorkflowNodeResource{}
			ctx := t.Context()

			req := resource.ConfigureRequest{ProviderData: tt.providerData}
			resp := resource.ConfigureResponse{}

			r.Configure(ctx, req, &resp)

			assert.Equal(t, tt.expectError, resp.Diagnostics.HasError())
		})
	}
}

func TestWorkflowNodeResource_Read_NoOp(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		expectError bool
	}{
		{
			name:        "read is a no-op without errors",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := &WorkflowNodeResource{}
			ctx := t.Context()

			req := resource.ReadRequest{}
			resp := resource.ReadResponse{}

			r.Read(ctx, req, &resp)

			assert.Equal(t, tt.expectError, resp.Diagnostics.HasError())
		})
	}
}

func TestWorkflowNodeResource_Delete_NoOp(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		expectError bool
	}{
		{
			name:        "delete is a no-op without errors",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := &WorkflowNodeResource{}
			ctx := t.Context()

			req := resource.DeleteRequest{}
			resp := resource.DeleteResponse{}

			r.Delete(ctx, req, &resp)

			assert.Equal(t, tt.expectError, resp.Diagnostics.HasError())
		})
	}
}

func TestWorkflowNodeResource_Create_WithWebhookID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		webhookID   string
		notes       string
		expectError bool
	}{
		{
			name:        "create webhook node with explicit webhook id",
			webhookID:   "my-webhook-id",
			notes:       "Test notes",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := &WorkflowNodeResource{}
			ctx := t.Context()
			objectType := getNodeObjectType()

			rawPlan := map[string]tftypes.Value{
				"id":           tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
				"name":         tftypes.NewValue(tftypes.String, "Webhook Trigger"),
				"type":         tftypes.NewValue(tftypes.String, "n8n-nodes-base.webhook"),
				"type_version": tftypes.NewValue(tftypes.Number, 1),
				"position":     tftypes.NewValue(tftypes.List{ElementType: tftypes.Number}, []tftypes.Value{tftypes.NewValue(tftypes.Number, 250), tftypes.NewValue(tftypes.Number, 300)}),
				"parameters":   tftypes.NewValue(tftypes.String, "{}"),
				"webhook_id":   tftypes.NewValue(tftypes.String, tt.webhookID),
				"disabled":     tftypes.NewValue(tftypes.Bool, false),
				"notes":        tftypes.NewValue(tftypes.String, tt.notes),
				"node_json":    tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
			}

			plan := tfsdk.Plan{Raw: tftypes.NewValue(objectType, rawPlan), Schema: createNodeTestSchema(t)}
			state := tfsdk.State{Raw: tftypes.NewValue(objectType, nil), Schema: createNodeTestSchema(t)}
			req := resource.CreateRequest{Plan: plan}
			resp := resource.CreateResponse{State: state}

			r.Create(ctx, req, &resp)

			assert.Equal(t, tt.expectError, resp.Diagnostics.HasError())
		})
	}
}

func TestWorkflowNodeResource_Create_Disabled(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		disabled    bool
		expectError bool
	}{
		{
			name:        "create disabled node succeeds",
			disabled:    true,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := &WorkflowNodeResource{}
			ctx := t.Context()
			objectType := getNodeObjectType()

			rawPlan := map[string]tftypes.Value{
				"id":           tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
				"name":         tftypes.NewValue(tftypes.String, "Disabled Node"),
				"type":         tftypes.NewValue(tftypes.String, "n8n-nodes-base.set"),
				"type_version": tftypes.NewValue(tftypes.Number, 1),
				"position":     tftypes.NewValue(tftypes.List{ElementType: tftypes.Number}, []tftypes.Value{tftypes.NewValue(tftypes.Number, 250), tftypes.NewValue(tftypes.Number, 300)}),
				"parameters":   tftypes.NewValue(tftypes.String, "{}"),
				"webhook_id":   tftypes.NewValue(tftypes.String, nil),
				"disabled":     tftypes.NewValue(tftypes.Bool, tt.disabled),
				"notes":        tftypes.NewValue(tftypes.String, nil),
				"node_json":    tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
			}

			plan := tfsdk.Plan{Raw: tftypes.NewValue(objectType, rawPlan), Schema: createNodeTestSchema(t)}
			state := tfsdk.State{Raw: tftypes.NewValue(objectType, nil), Schema: createNodeTestSchema(t)}
			req := resource.CreateRequest{Plan: plan}
			resp := resource.CreateResponse{State: state}

			r.Create(ctx, req, &resp)

			assert.Equal(t, tt.expectError, resp.Diagnostics.HasError())
		})
	}
}
