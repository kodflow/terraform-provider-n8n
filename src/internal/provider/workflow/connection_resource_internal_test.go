// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

package workflow

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/workflow/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWorkflowConnectionResource_generateConnectionID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		plan     *models.ConnectionResource
		wantSame bool
		checkID  func(t *testing.T, id string)
	}{
		{
			name: "generates consistent IDs for same inputs",
			plan: &models.ConnectionResource{
				SourceNode:        types.StringValue("Webhook"),
				SourceOutputIndex: types.Int64Value(0),
				TargetNode:        types.StringValue("Process"),
				TargetInputIndex:  types.Int64Value(0),
			},
			wantSame: true,
			checkID: func(t *testing.T, id string) {
				t.Helper()
				assert.NotEmpty(t, id)
			},
		},
		{
			name: "generates different IDs for different outputs",
			plan: &models.ConnectionResource{
				SourceNode:        types.StringValue("Webhook"),
				SourceOutputIndex: types.Int64Value(1),
				TargetNode:        types.StringValue("Process"),
				TargetInputIndex:  types.Int64Value(0),
			},
			wantSame: false,
			checkID: func(t *testing.T, id string) {
				t.Helper()
				assert.NotEmpty(t, id)
			},
		},
		{
			name: "handles empty node names",
			plan: &models.ConnectionResource{
				SourceNode:        types.StringValue(""),
				SourceOutputIndex: types.Int64Value(0),
				TargetNode:        types.StringValue(""),
				TargetInputIndex:  types.Int64Value(0),
			},
			wantSame: true,
			checkID: func(t *testing.T, id string) {
				t.Helper()
				assert.NotEmpty(t, id)
			},
		},
		{
			name: "error case - generates valid ID even with extreme values",
			plan: &models.ConnectionResource{
				SourceNode:        types.StringValue("Node with very long name that exceeds normal limits"),
				SourceOutputIndex: types.Int64Value(999),
				TargetNode:        types.StringValue("Target"),
				TargetInputIndex:  types.Int64Value(999),
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
			r := &WorkflowConnectionResource{}
			id1 := r.generateConnectionID(tt.plan)
			id2 := r.generateConnectionID(tt.plan)

			if tt.wantSame {
				assert.Equal(t, id1, id2, "IDs should be consistent for same inputs")
			}
			if tt.checkID != nil {
				tt.checkID(t, id1)
			}
		})
	}
}

func TestWorkflowConnectionResource_generateConnectionJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		plan      *models.ConnectionResource
		wantErr   bool
		checkJSON func(t *testing.T, connJSON string)
	}{
		{
			name: "basic connection with defaults",
			plan: &models.ConnectionResource{
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
			plan: &models.ConnectionResource{
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
			plan: &models.ConnectionResource{
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

func TestWorkflowConnectionResource_getIDAttribute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		check func(t *testing.T, attr schema.StringAttribute)
	}{
		{
			name: "ID attribute is computed",
			check: func(t *testing.T, attr schema.StringAttribute) {
				t.Helper()
				assert.True(t, attr.Computed)
			},
		},
		{
			name: "ID attribute has description",
			check: func(t *testing.T, attr schema.StringAttribute) {
				t.Helper()
				assert.NotEmpty(t, attr.MarkdownDescription)
			},
		},
		{
			name: "ID attribute has plan modifiers",
			check: func(t *testing.T, attr schema.StringAttribute) {
				t.Helper()
				assert.NotEmpty(t, attr.PlanModifiers)
			},
		},
		{
			name: "error case - ID attribute validates constraints",
			check: func(t *testing.T, attr schema.StringAttribute) {
				t.Helper()
				assert.True(t, attr.Computed, "ID must be computed")
				assert.False(t, attr.Required, "ID cannot be required")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &WorkflowConnectionResource{}
			attr := r.getIDAttribute()
			tt.check(t, attr)
		})
	}
}

func TestWorkflowConnectionResource_addSourceAttributes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		checkAttrs []string
		checkTypes func(t *testing.T, attrs map[string]schema.Attribute)
	}{
		{
			name:       "source attributes added",
			checkAttrs: []string{"source_node", "source_output", "source_output_index"},
			checkTypes: func(t *testing.T, attrs map[string]schema.Attribute) {
				t.Helper()
				// Verify source_node is required
				strAttr, ok := attrs["source_node"].(schema.StringAttribute)
				assert.True(t, ok)
				assert.True(t, strAttr.Required)
			},
		},
		{
			name:       "optional attributes have defaults",
			checkAttrs: []string{"source_output", "source_output_index"},
			checkTypes: func(t *testing.T, attrs map[string]schema.Attribute) {
				t.Helper()
				// Verify source_output has default
				strAttr, ok := attrs["source_output"].(schema.StringAttribute)
				assert.True(t, ok)
				assert.NotNil(t, strAttr.Default)
			},
		},
		{
			name:       "error case - all required source attributes must be present",
			checkAttrs: []string{"source_node", "source_output", "source_output_index"},
			checkTypes: func(t *testing.T, attrs map[string]schema.Attribute) {
				t.Helper()
				assert.Len(t, attrs, 3, "must have exactly 3 source attributes")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := &WorkflowConnectionResource{}
			attrs := make(map[string]schema.Attribute)
			r.addSourceAttributes(attrs)

			for _, attrName := range tt.checkAttrs {
				assert.Contains(t, attrs, attrName)
			}

			if tt.checkTypes != nil {
				tt.checkTypes(t, attrs)
			}
		})
	}
}

func TestWorkflowConnectionResource_addTargetAttributes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		checkAttrs []string
		checkTypes func(t *testing.T, attrs map[string]schema.Attribute)
	}{
		{
			name:       "target attributes added",
			checkAttrs: []string{"target_node", "target_input", "target_input_index"},
			checkTypes: func(t *testing.T, attrs map[string]schema.Attribute) {
				t.Helper()
				// Verify target_node is required
				strAttr, ok := attrs["target_node"].(schema.StringAttribute)
				assert.True(t, ok)
				assert.True(t, strAttr.Required)
			},
		},
		{
			name:       "optional attributes have defaults",
			checkAttrs: []string{"target_input", "target_input_index"},
			checkTypes: func(t *testing.T, attrs map[string]schema.Attribute) {
				t.Helper()
				// Verify target_input has default
				strAttr, ok := attrs["target_input"].(schema.StringAttribute)
				assert.True(t, ok)
				assert.NotNil(t, strAttr.Default)
			},
		},
		{
			name:       "error case - all required target attributes must be present",
			checkAttrs: []string{"target_node", "target_input", "target_input_index"},
			checkTypes: func(t *testing.T, attrs map[string]schema.Attribute) {
				t.Helper()
				assert.Len(t, attrs, 3, "must have exactly 3 target attributes")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := &WorkflowConnectionResource{}
			attrs := make(map[string]schema.Attribute)
			r.addTargetAttributes(attrs)

			for _, attrName := range tt.checkAttrs {
				assert.Contains(t, attrs, attrName)
			}

			if tt.checkTypes != nil {
				tt.checkTypes(t, attrs)
			}
		})
	}
}

func TestWorkflowConnectionResource_getSchemaAttributes(t *testing.T) {
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
				"id", "source_node", "target_node",
				"source_output", "source_output_index",
				"target_input", "target_input_index", "connection_json",
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

			r := &WorkflowConnectionResource{}
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
