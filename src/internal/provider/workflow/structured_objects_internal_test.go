package workflow

import (
	"context"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/workflow/models"
)

func TestMapWorkflowToModelStructuredObjects(t *testing.T) {
	tests := []struct {
		name string
		data map[string]any
		want map[string]string
	}{
		{
			name: "pinned node items returned by the n8n API",
			data: map[string]any{"Fetch pluviometer": []any{map[string]any{"json": map[string]any{"rain": float64(12)}}}},
			want: map[string]string{"Fetch pluviometer": `[{"json":{"rain":12}}]`},
		},
		{
			name: "nested metadata and JSON scalars",
			data: map[string]any{
				"nested":  map[string]any{"enabled": true, "items": []any{float64(1), "two"}},
				"boolean": true, "number": float64(3.5), "null": nil,
			},
			want: map[string]string{
				"nested":  `{"enabled":true,"items":[1,"two"]}`,
				"boolean": "true", "number": "3.5", "null": "null",
			},
		},
		{
			name: "existing strings are preserved without JSON quoting",
			data: map[string]any{"plain": "text", "json_string": `{"already":"encoded"}`, "empty": ""},
			want: map[string]string{"plain": "text", "json_string": `{"already":"encoded"}`, "empty": ""},
		},
		{name: "empty object", data: map[string]any{}, want: map[string]string{}},
		{name: "absent object", data: nil, want: nil},
	}
	for _, field := range []string{"meta", "pin_data"} {
		for _, tc := range tests {
			t.Run(field+"/"+tc.name, func(t *testing.T) {
				ctx := context.Background()
				workflow := &n8nsdk.Workflow{Name: "pluviometer"}
				if field == "meta" {
					workflow.Meta = tc.data
				} else {
					workflow.PinData = tc.data
				}
				model := &models.Resource{}
				var diagnostics diag.Diagnostics
				mapWorkflowToModel(ctx, workflow, model, &diagnostics)
				if diagnostics.HasError() {
					t.Fatalf("mapping failed: %v", diagnostics)
				}
				var value types.Map
				if field == "meta" {
					value = model.Meta
				} else {
					value = model.PinData
				}
				if tc.data == nil {
					if !value.IsNull() {
						t.Fatal("absent API object must remain a null map")
					}
					return
				}
				var got map[string]string
				diagnostics = value.ElementsAs(ctx, &got, false)
				if diagnostics.HasError() || !reflect.DeepEqual(tc.want, got) {
					t.Fatalf("expected %#v, got %#v (diagnostics: %v)", tc.want, got, diagnostics)
				}
			})
		}
	}
}

func TestWorkflowObjectAsStringMapInvalidValue(t *testing.T) {
	value, diagnostics := workflowObjectAsStringMap(context.Background(), map[string]any{
		"invalid": make(chan int),
	})
	if !diagnostics.HasError() || !value.IsNull() {
		t.Fatalf("expected a diagnostic and a null map, got %v and %v", diagnostics, value)
	}
}
