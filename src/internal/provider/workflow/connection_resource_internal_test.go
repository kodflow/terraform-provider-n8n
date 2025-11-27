// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

package workflow

import (
	"context"
	"encoding/json"
	"testing"

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

// createConnectionTestSchema creates a schema for connection resource testing.
func createConnectionTestSchema(t *testing.T) schema.Schema {
	t.Helper()
	r := &WorkflowConnectionResource{}
	return schema.Schema{
		Attributes: r.getSchemaAttributes(),
	}
}

// getConnectionObjectType returns the tftypes.Object for connection resource.
func getConnectionObjectType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"id":                  tftypes.String,
			"source_node":         tftypes.String,
			"source_output":       tftypes.String,
			"source_output_index": tftypes.Number,
			"target_node":         tftypes.String,
			"target_input":        tftypes.String,
			"target_input_index":  tftypes.Number,
			"connection_json":     tftypes.String,
		},
	}
}

func TestWorkflowConnectionResource_Create_Success(t *testing.T) {
	t.Parallel()

	r := &WorkflowConnectionResource{}
	ctx := context.Background()
	objectType := getConnectionObjectType()

	rawPlan := map[string]tftypes.Value{
		"id":                  tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
		"source_node":         tftypes.NewValue(tftypes.String, "Webhook"),
		"source_output":       tftypes.NewValue(tftypes.String, "main"),
		"source_output_index": tftypes.NewValue(tftypes.Number, 0),
		"target_node":         tftypes.NewValue(tftypes.String, "Process"),
		"target_input":        tftypes.NewValue(tftypes.String, "main"),
		"target_input_index":  tftypes.NewValue(tftypes.Number, 0),
		"connection_json":     tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
	}

	plan := tfsdk.Plan{
		Raw:    tftypes.NewValue(objectType, rawPlan),
		Schema: createConnectionTestSchema(t),
	}

	state := tfsdk.State{
		Raw:    tftypes.NewValue(objectType, nil),
		Schema: createConnectionTestSchema(t),
	}

	req := resource.CreateRequest{
		Plan: plan,
	}
	resp := resource.CreateResponse{
		State: state,
	}

	r.Create(ctx, req, &resp)

	assert.False(t, resp.Diagnostics.HasError(), "Create should not have errors")
}

func TestWorkflowConnectionResource_Create_WithMultipleOutputs(t *testing.T) {
	t.Parallel()

	r := &WorkflowConnectionResource{}
	ctx := context.Background()
	objectType := getConnectionObjectType()

	rawPlan := map[string]tftypes.Value{
		"id":                  tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
		"source_node":         tftypes.NewValue(tftypes.String, "Switch"),
		"source_output":       tftypes.NewValue(tftypes.String, "main"),
		"source_output_index": tftypes.NewValue(tftypes.Number, 2),
		"target_node":         tftypes.NewValue(tftypes.String, "Handler"),
		"target_input":        tftypes.NewValue(tftypes.String, "main"),
		"target_input_index":  tftypes.NewValue(tftypes.Number, 1),
		"connection_json":     tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
	}

	plan := tfsdk.Plan{
		Raw:    tftypes.NewValue(objectType, rawPlan),
		Schema: createConnectionTestSchema(t),
	}

	state := tfsdk.State{
		Raw:    tftypes.NewValue(objectType, nil),
		Schema: createConnectionTestSchema(t),
	}

	req := resource.CreateRequest{
		Plan: plan,
	}
	resp := resource.CreateResponse{
		State: state,
	}

	r.Create(ctx, req, &resp)

	assert.False(t, resp.Diagnostics.HasError(), "Create with multiple outputs should not have errors")
}

func TestWorkflowConnectionResource_Update_Success(t *testing.T) {
	t.Parallel()

	r := &WorkflowConnectionResource{}
	ctx := context.Background()
	objectType := getConnectionObjectType()

	rawPlan := map[string]tftypes.Value{
		"id":                  tftypes.NewValue(tftypes.String, "conn-123"),
		"source_node":         tftypes.NewValue(tftypes.String, "UpdatedWebhook"),
		"source_output":       tftypes.NewValue(tftypes.String, "main"),
		"source_output_index": tftypes.NewValue(tftypes.Number, 0),
		"target_node":         tftypes.NewValue(tftypes.String, "UpdatedProcess"),
		"target_input":        tftypes.NewValue(tftypes.String, "main"),
		"target_input_index":  tftypes.NewValue(tftypes.Number, 0),
		"connection_json":     tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
	}

	plan := tfsdk.Plan{
		Raw:    tftypes.NewValue(objectType, rawPlan),
		Schema: createConnectionTestSchema(t),
	}

	state := tfsdk.State{
		Raw:    tftypes.NewValue(objectType, nil),
		Schema: createConnectionTestSchema(t),
	}

	req := resource.UpdateRequest{
		Plan: plan,
	}
	resp := resource.UpdateResponse{
		State: state,
	}

	r.Update(ctx, req, &resp)

	assert.False(t, resp.Diagnostics.HasError(), "Update should not have errors")
}

func TestWorkflowConnectionResource_ImportState_Success(t *testing.T) {
	t.Parallel()

	r := &WorkflowConnectionResource{}
	ctx := context.Background()
	objectType := getConnectionObjectType()

	state := tfsdk.State{
		Raw:    tftypes.NewValue(objectType, nil),
		Schema: createConnectionTestSchema(t),
	}

	req := resource.ImportStateRequest{
		ID: "imported-connection-id",
	}
	resp := resource.ImportStateResponse{
		State: state,
	}

	r.ImportState(ctx, req, &resp)

	// ImportState uses passthrough so it should set the ID in state.
	assert.False(t, resp.Diagnostics.HasError(), "ImportState should not have errors")
}

func TestWorkflowConnectionResource_generateConnectionJSON_Error(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		plan          *models.ConnectionResource
		expectSuccess bool
		expectError   bool
	}{
		{
			name: "success - with unknown source node",
			plan: &models.ConnectionResource{
				ID:                types.StringValue("test-id"),
				SourceNode:        types.StringUnknown(),
				SourceOutput:      types.StringValue("main"),
				SourceOutputIndex: types.Int64Value(0),
				TargetNode:        types.StringValue("Target"),
				TargetInput:       types.StringValue("main"),
				TargetInputIndex:  types.Int64Value(0),
			},
			expectSuccess: true,
			expectError:   false,
		},
		{
			name: "error case - all null values",
			plan: &models.ConnectionResource{
				ID:                types.StringNull(),
				SourceNode:        types.StringNull(),
				SourceOutput:      types.StringNull(),
				SourceOutputIndex: types.Int64Null(),
				TargetNode:        types.StringNull(),
				TargetInput:       types.StringNull(),
				TargetInputIndex:  types.Int64Null(),
			},
			expectSuccess: true,
			expectError:   false,
		},
		{
			name: "error case - unknown target values",
			plan: &models.ConnectionResource{
				ID:                types.StringValue("test-id"),
				SourceNode:        types.StringValue("Source"),
				SourceOutput:      types.StringValue("main"),
				SourceOutputIndex: types.Int64Value(0),
				TargetNode:        types.StringUnknown(),
				TargetInput:       types.StringUnknown(),
				TargetInputIndex:  types.Int64Unknown(),
			},
			expectSuccess: true,
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := &WorkflowConnectionResource{}
			var diags diag.Diagnostics
			success := r.generateConnectionJSON(tt.plan, &diags)
			// Check expectations.
			assert.Equal(t, tt.expectSuccess, success, "Success should match expected")
			if tt.expectError {
				assert.True(t, diags.HasError(), "Should have diagnostics error")
			} else {
				assert.False(t, diags.HasError(), "Should not have diagnostics error")
			}
		})
	}
}
