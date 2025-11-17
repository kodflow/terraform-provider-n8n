// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

package node

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/workflow/node/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWorkflowNodeResource_Metadata(t *testing.T) {
	t.Parallel()

	r := NewWorkflowNodeResource()
	req := resource.MetadataRequest{
		ProviderTypeName: "n8n",
	}
	resp := &resource.MetadataResponse{}

	r.Metadata(context.Background(), req, resp)

	assert.Equal(t, "n8n_workflow_node", resp.TypeName)
}

func TestWorkflowNodeResource_Schema(t *testing.T) {
	t.Parallel()

	r := NewWorkflowNodeResource()
	req := resource.SchemaRequest{}
	resp := &resource.SchemaResponse{}

	r.Schema(context.Background(), req, resp)

	require.NotNil(t, resp.Schema)
	assert.Contains(t, resp.Schema.Attributes, "id")
	assert.Contains(t, resp.Schema.Attributes, "name")
	assert.Contains(t, resp.Schema.Attributes, "type")
	assert.Contains(t, resp.Schema.Attributes, "position")
	assert.Contains(t, resp.Schema.Attributes, "parameters")
	assert.Contains(t, resp.Schema.Attributes, "node_json")
}

func TestWorkflowNodeResource_GenerateNodeID(t *testing.T) {
	t.Parallel()

	r := &WorkflowNodeResource{}

	plan := &models.Resource{
		Name: types.StringValue("Test Node"),
		Type: types.StringValue("n8n-nodes-base.webhook"),
	}

	id1 := r.generateNodeID(plan)
	id2 := r.generateNodeID(plan)

	// Should generate consistent IDs.
	assert.Equal(t, id1, id2)
	assert.NotEmpty(t, id1)
}

func TestWorkflowNodeResource_GenerateNodeJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		plan      *models.Resource
		wantErr   bool
		checkJSON func(t *testing.T, nodeJSON string)
	}{
		{
			name: "basic node with minimal fields",
			plan: &models.Resource{
				ID:          types.StringValue("test-id-123"),
				Name:        types.StringValue("Webhook"),
				Type:        types.StringValue("n8n-nodes-base.webhook"),
				TypeVersion: types.Int64Value(1),
				Position:    types.ListValueMust(types.Int64Type, []attr.Value{types.Int64Value(250), types.Int64Value(300)}),
			},
			wantErr: false,
			checkJSON: func(t *testing.T, nodeJSON string) {
				t.Helper()

				var node map[string]interface{}
				err := json.Unmarshal([]byte(nodeJSON), &node)
				require.NoError(t, err)

				assert.Equal(t, "test-id-123", node["id"])
				assert.Equal(t, "Webhook", node["name"])
				assert.Equal(t, "n8n-nodes-base.webhook", node["type"])
				assert.Equal(t, float64(1), node["typeVersion"])

				position, ok := node["position"].([]interface{})
				require.True(t, ok)
				assert.Len(t, position, 2)
				assert.Equal(t, float64(250), position[0])
				assert.Equal(t, float64(300), position[1])
			},
		},
		{
			name: "node with parameters",
			plan: &models.Resource{
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

				var node map[string]interface{}
				err := json.Unmarshal([]byte(nodeJSON), &node)
				require.NoError(t, err)

				params, ok := node["parameters"].(map[string]interface{})
				require.True(t, ok)
				assert.Equal(t, "runOnceForAllItems", params["mode"])
				assert.Equal(t, "return items;", params["jsCode"])
			},
		},
		{
			name: "webhook node with webhookId",
			plan: &models.Resource{
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

				var node map[string]interface{}
				err := json.Unmarshal([]byte(nodeJSON), &node)
				require.NoError(t, err)

				assert.Equal(t, "my-webhook-id", node["webhookId"])
			},
		},
		{
			name: "disabled node with notes",
			plan: &models.Resource{
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

				var node map[string]interface{}
				err := json.Unmarshal([]byte(nodeJSON), &node)
				require.NoError(t, err)

				assert.Equal(t, true, node["disabled"])
				assert.Equal(t, "This node is temporarily disabled", node["notes"])
			},
		},
		{
			name: "invalid parameters JSON",
			plan: &models.Resource{
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

			ctx := context.Background()
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
