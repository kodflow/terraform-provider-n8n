// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

package connection

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/workflow/connection/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWorkflowConnectionResource_Metadata(t *testing.T) {
	t.Parallel()

	r := NewWorkflowConnectionResource()
	req := resource.MetadataRequest{
		ProviderTypeName: "n8n",
	}
	resp := &resource.MetadataResponse{}

	r.Metadata(context.Background(), req, resp)

	assert.Equal(t, "n8n_workflow_connection", resp.TypeName)
}

func TestWorkflowConnectionResource_Schema(t *testing.T) {
	t.Parallel()

	r := NewWorkflowConnectionResource()
	req := resource.SchemaRequest{}
	resp := &resource.SchemaResponse{}

	r.Schema(context.Background(), req, resp)

	require.NotNil(t, resp.Schema)
	assert.Contains(t, resp.Schema.Attributes, "id")
	assert.Contains(t, resp.Schema.Attributes, "source_node")
	assert.Contains(t, resp.Schema.Attributes, "target_node")
	assert.Contains(t, resp.Schema.Attributes, "source_output")
	assert.Contains(t, resp.Schema.Attributes, "source_output_index")
	assert.Contains(t, resp.Schema.Attributes, "target_input")
	assert.Contains(t, resp.Schema.Attributes, "target_input_index")
	assert.Contains(t, resp.Schema.Attributes, "connection_json")
}

func TestWorkflowConnectionResource_GenerateConnectionID(t *testing.T) {
	t.Parallel()

	r := &WorkflowConnectionResource{}

	plan := &models.Resource{
		SourceNode:        types.StringValue("Webhook"),
		SourceOutputIndex: types.Int64Value(0),
		TargetNode:        types.StringValue("Process"),
		TargetInputIndex:  types.Int64Value(0),
	}

	id1 := r.generateConnectionID(plan)
	id2 := r.generateConnectionID(plan)

	// Should generate consistent IDs.
	assert.Equal(t, id1, id2)
	assert.NotEmpty(t, id1)
}

func TestWorkflowConnectionResource_GenerateConnectionJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		plan      *models.Resource
		wantErr   bool
		checkJSON func(t *testing.T, connJSON string)
	}{
		{
			name: "basic connection with defaults",
			plan: &models.Resource{
				SourceNode:        types.StringValue("Webhook"),
				SourceOutput:      types.StringValue("main"),
				SourceOutputIndex: types.Int64Value(0),
				TargetNode:        types.StringValue("Process Data"),
				TargetInput:       types.StringValue("main"),
				TargetInputIndex:  types.Int64Value(0),
			},
			wantErr: false,
			checkJSON: func(t *testing.T, connJSON string) {
				t.Helper()

				var conn map[string]interface{}
				err := json.Unmarshal([]byte(connJSON), &conn)
				require.NoError(t, err)

				assert.Equal(t, "Webhook", conn["source_node"])
				assert.Equal(t, "main", conn["source_output"])
				assert.Equal(t, float64(0), conn["source_output_index"])
				assert.Equal(t, "Process Data", conn["target_node"])
				assert.Equal(t, "main", conn["target_input"])
				assert.Equal(t, float64(0), conn["target_input_index"])
			},
		},
		{
			name: "switch connection with output index 1",
			plan: &models.Resource{
				SourceNode:        types.StringValue("Switch"),
				SourceOutput:      types.StringValue("main"),
				SourceOutputIndex: types.Int64Value(1),
				TargetNode:        types.StringValue("Branch 2"),
				TargetInput:       types.StringValue("main"),
				TargetInputIndex:  types.Int64Value(0),
			},
			wantErr: false,
			checkJSON: func(t *testing.T, connJSON string) {
				t.Helper()

				var conn map[string]interface{}
				err := json.Unmarshal([]byte(connJSON), &conn)
				require.NoError(t, err)

				assert.Equal(t, "Switch", conn["source_node"])
				assert.Equal(t, float64(1), conn["source_output_index"])
				assert.Equal(t, "Branch 2", conn["target_node"])
			},
		},
		{
			name: "connection with output index 5",
			plan: &models.Resource{
				SourceNode:        types.StringValue("Switch Multiple"),
				SourceOutput:      types.StringValue("main"),
				SourceOutputIndex: types.Int64Value(5),
				TargetNode:        types.StringValue("Branch 6"),
				TargetInput:       types.StringValue("main"),
				TargetInputIndex:  types.Int64Value(0),
			},
			wantErr: false,
			checkJSON: func(t *testing.T, connJSON string) {
				t.Helper()

				var conn map[string]interface{}
				err := json.Unmarshal([]byte(connJSON), &conn)
				require.NoError(t, err)

				assert.Equal(t, float64(5), conn["source_output_index"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := &WorkflowConnectionResource{}
			var diags diag.Diagnostics

			success := r.generateConnectionJSON(tt.plan, &diags)

			if tt.wantErr {
				assert.False(t, success)
				assert.True(t, diags.HasError())
			} else {
				assert.True(t, success)
				assert.False(t, diags.HasError())
				require.NotEmpty(t, tt.plan.ConnectionJSON.ValueString())

				if tt.checkJSON != nil {
					tt.checkJSON(t, tt.plan.ConnectionJSON.ValueString())
				}
			}
		})
	}
}
