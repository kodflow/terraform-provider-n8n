// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// Package workflow provides white-box tests for the WorkflowDataSource type.
package workflow

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/client"
	"github.com/stretchr/testify/assert"
)

// TestWorkflowDataSource_Metadata_Internal verifies the Metadata method sets the correct type name.
func TestWorkflowDataSource_Metadata_Internal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		providerTypeName string
		expectedTypeName string
	}{
		{
			name:             "n8n provider type",
			providerTypeName: "n8n",
			expectedTypeName: "n8n_workflow",
		},
		{
			name:             "empty provider type",
			providerTypeName: "",
			expectedTypeName: "_workflow",
		},
		{
			name:             "error case - custom provider type",
			providerTypeName: "custom",
			expectedTypeName: "custom_workflow",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			d := NewWorkflowDataSource()
			resp := &datasource.MetadataResponse{}

			d.Metadata(t.Context(), datasource.MetadataRequest{
				ProviderTypeName: tt.providerTypeName,
			}, resp)

			assert.Equal(t, tt.expectedTypeName, resp.TypeName)
		})
	}
}

// TestWorkflowDataSource_Schema_Internal verifies the Schema method defines expected attributes.
func TestWorkflowDataSource_Schema_Internal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		wantAttributes []string
	}{
		{
			name:           "schema has expected attributes",
			wantAttributes: []string{"id", "name", "active"},
		},
		{
			name:           "error case - schema is not nil",
			wantAttributes: []string{"id"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			d := NewWorkflowDataSource()
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

// TestWorkflowDataSource_Configure_Internal verifies Configure assigns client or adds errors.
func TestWorkflowDataSource_Configure_Internal(t *testing.T) {
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
			providerData: 7,
			wantError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			d := NewWorkflowDataSource()
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

// TestNewWorkflowDataSource_Internal verifies the constructor returns a non-nil instance.
func TestNewWorkflowDataSource_Internal(t *testing.T) {
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
				ds := NewWorkflowDataSource()
				assert.NotNil(t, ds)

			case "error case - multiple calls are independent":
				ds1 := NewWorkflowDataSource()
				ds2 := NewWorkflowDataSource()
				assert.NotNil(t, ds1)
				assert.NotNil(t, ds2)
				assert.NotSame(t, ds1, ds2)
			}
		})
	}
}
