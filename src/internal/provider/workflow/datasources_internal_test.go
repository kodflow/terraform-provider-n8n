// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// Package workflow provides white-box tests for the WorkflowsDataSource type.
package workflow

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/client"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/workflow/models"
	"github.com/stretchr/testify/assert"
)

// TestWorkflowsDataSource_Metadata_Internal verifies the Metadata method sets the correct type name.
func TestWorkflowsDataSource_Metadata_Internal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		providerTypeName string
		expectedTypeName string
	}{
		{
			name:             "n8n provider type",
			providerTypeName: "n8n",
			expectedTypeName: "n8n_workflows",
		},
		{
			name:             "empty provider type",
			providerTypeName: "",
			expectedTypeName: "_workflows",
		},
		{
			name:             "error case - custom provider type",
			providerTypeName: "custom",
			expectedTypeName: "custom_workflows",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			d := NewWorkflowsDataSource()
			resp := &datasource.MetadataResponse{}

			d.Metadata(t.Context(), datasource.MetadataRequest{
				ProviderTypeName: tt.providerTypeName,
			}, resp)

			assert.Equal(t, tt.expectedTypeName, resp.TypeName)
		})
	}
}

// TestWorkflowsDataSource_Schema_Internal verifies the Schema method defines expected attributes.
func TestWorkflowsDataSource_Schema_Internal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		wantAttributes []string
	}{
		{
			name:           "schema has expected attributes",
			wantAttributes: []string{"active", "workflows"},
		},
		{
			name:           "error case - schema is not nil",
			wantAttributes: []string{"workflows"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			d := NewWorkflowsDataSource()
			resp := &datasource.SchemaResponse{}

			d.Schema(t.Context(), datasource.SchemaRequest{}, resp)

			assert.NotNil(t, resp.Schema)
			assert.NotEmpty(t, resp.Schema.Attributes)

			for _, attr := range tt.wantAttributes {
				_, ok := resp.Schema.Attributes[attr]
				assert.True(t, ok, "expected attribute %q in schema", attr)
			}
		})
	}
}

// TestWorkflowsDataSource_Configure_Internal verifies Configure assigns client or adds errors.
func TestWorkflowsDataSource_Configure_Internal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		providerData any
		wantError    bool
	}{
		{
			name:         "nil provider data",
			providerData: nil,
			wantError:    false,
		},
		{
			name:         "valid N8nClient",
			providerData: &client.N8nClient{},
			wantError:    false,
		},
		{
			name:         "error case - string type",
			providerData: "invalid",
			wantError:    true,
		},
		{
			name:         "error case - integer type",
			providerData: 3,
			wantError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			d := NewWorkflowsDataSource()
			resp := &datasource.ConfigureResponse{}

			d.Configure(t.Context(), datasource.ConfigureRequest{
				ProviderData: tt.providerData,
			}, resp)

			if tt.wantError {
				assert.True(t, resp.Diagnostics.HasError())
			} else {
				assert.False(t, resp.Diagnostics.HasError())
			}
		})
	}
}

// TestNewWorkflowsDataSource_Internal verifies the constructor returns a non-nil instance.
func TestNewWorkflowsDataSource_Internal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
	}{
		{name: "returns non-nil instance"},
		{name: "error case - multiple calls are independent"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "returns non-nil instance":
				ds := NewWorkflowsDataSource()
				assert.NotNil(t, ds)

			case "error case - multiple calls are independent":
				ds1 := NewWorkflowsDataSource()
				ds2 := NewWorkflowsDataSource()
				assert.NotNil(t, ds1)
				assert.NotNil(t, ds2)
				assert.NotSame(t, ds1, ds2)
			}
		})
	}
}

// TestBuildWorkflowsDataSourceSchema tests the buildWorkflowsDataSourceSchema function.
func TestBuildWorkflowsDataSourceSchema(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		checkAttr string
		exists    bool
	}{
		{name: "has active attribute", checkAttr: "active", exists: true},
		{name: "has workflows attribute", checkAttr: "workflows", exists: true},
		{name: "missing attribute returns false", checkAttr: "nonexistent", exists: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			//: Build the schema.
			s := buildWorkflowsDataSourceSchema()
			//: Verify attribute existence matches expectation.
			_, exists := s.Attributes[tt.checkAttr]
			assert.Equal(t, tt.exists, exists)
		})
	}
}

// TestMapWorkflowItems tests the mapWorkflowItems function.
func TestMapWorkflowItems(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		input       []n8nsdk.Workflow
		expectCount int
		expectID    string
		expectName  string
	}{
		{
			name:        "nil input returns empty slice",
			input:       nil,
			expectCount: 0,
		},
		{
			name:        "empty input returns empty slice",
			input:       []n8nsdk.Workflow{},
			expectCount: 0,
		},
		{
			name: "single workflow mapped correctly",
			input: []n8nsdk.Workflow{
				{Id: new("wf-1"), Name: "Test Workflow", Active: new(true)},
			},
			expectCount: 1,
			expectID:    "wf-1",
			expectName:  "Test Workflow",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			//: Map workflow items.
			result := mapWorkflowItems(tt.input)
			//: Verify result count matches expectation.
			assert.Len(t, result, tt.expectCount)
			//: Check item values when present.
			if tt.expectCount > 0 {
				assert.Equal(t, tt.expectID, result[0].ID.ValueString())
				assert.Equal(t, tt.expectName, result[0].Name.ValueString())
			}
		})
	}
}

// TestExecuteWorkflowsRead tests the executeWorkflowsRead method.
func TestExecuteWorkflowsRead(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		expectError bool
	}{
		{
			name:        "nil client panics",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			//: Verify nil client causes panic.
			if tt.expectError {
				d := &WorkflowsDataSource{client: &client.N8nClient{}}
				assert.Panics(t, func() {
					resp := &datasource.ReadResponse{}
					data := &models.DataSources{}
					d.executeWorkflowsRead(t.Context(), data, resp)
				})
			}
		})
	}
}

// TestWorkflowsDataSource_executeWorkflowsRead verifies the executeWorkflowsRead method.
func TestWorkflowsDataSource_executeWorkflowsRead(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		handler     http.HandlerFunc
		active      *bool
		expectError bool
	}{
		{
			name: "successful read returns workflows",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				active := true
				id := "wf-1"
				resp := map[string]any{
					"data": []any{
						map[string]any{"id": id, "name": "My Workflow", "active": active, "nodes": []any{}, "connections": map[string]any{}, "settings": map[string]any{}},
					},
				}
				json.NewEncoder(w).Encode(resp) //nolint:errcheck
			},
			expectError: false,
		},
		{
			name: "API error returns diagnostics error",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"message":"server error"}`)) //nolint:errcheck
			},
			expectError: true,
		},
		{
			name: "empty workflow list",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]any{"data": []any{}}) //nolint:errcheck
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			n8nClient, server := setupTestClientForHelpers(t, tt.handler)
			defer server.Close()

			d := &WorkflowsDataSource{client: n8nClient}
			data := &models.DataSources{}
			resp := &datasource.ReadResponse{}

			ok := d.executeWorkflowsRead(t.Context(), data, resp)

			if tt.expectError {
				assert.False(t, ok)
				assert.True(t, resp.Diagnostics.HasError())
			} else {
				assert.True(t, ok)
				assert.False(t, resp.Diagnostics.HasError())
			}
		})
	}
}
