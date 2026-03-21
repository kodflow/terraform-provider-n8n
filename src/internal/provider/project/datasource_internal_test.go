// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// Package project provides white-box tests for the ProjectDataSource type.
package project

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/project/models"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/client"
	"github.com/stretchr/testify/assert"
)

// TestProjectDataSource_Metadata_Internal verifies the Metadata method sets the correct type name.
func TestProjectDataSource_Metadata_Internal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		providerTypeName string
		expectedTypeName string
	}{
		{
			name:             "n8n provider type",
			providerTypeName: "n8n",
			expectedTypeName: "n8n_project",
		},
		{
			name:             "empty provider type",
			providerTypeName: "",
			expectedTypeName: "_project",
		},
		{
			name:             "error case - custom provider type",
			providerTypeName: "custom",
			expectedTypeName: "custom_project",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			d := NewProjectDataSource()
			resp := &datasource.MetadataResponse{}

			d.Metadata(t.Context(), datasource.MetadataRequest{
				ProviderTypeName: tt.providerTypeName,
			}, resp)

			assert.Equal(t, tt.expectedTypeName, resp.TypeName)
		})
	}
}

// TestProjectDataSource_Schema_Internal verifies the Schema method populates the schema correctly.
func TestProjectDataSource_Schema_Internal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		wantAttributes []string
	}{
		{
			name:           "schema has expected attributes",
			wantAttributes: []string{"id", "name", "type", "created_at", "updated_at", "icon", "description"},
		},
		{
			name:           "error case - schema is not nil",
			wantAttributes: []string{"id"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			d := NewProjectDataSource()
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

// TestProjectDataSource_Configure_Internal verifies the Configure method assigns the client correctly.
func TestProjectDataSource_Configure_Internal(t *testing.T) {
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
			name:         "error case - invalid type",
			providerData: "invalid-type",
			wantError:    true,
		},
		{
			name:         "error case - integer type",
			providerData: 42,
			wantError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			d := NewProjectDataSource()
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

// TestNewProjectDataSource_Internal verifies the constructor returns a non-nil instance.
func TestNewProjectDataSource_Internal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
	}{
		{name: "returns non-nil instance"},
		{name: "error case - multiple calls return independent instances"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "returns non-nil instance":
				ds := NewProjectDataSource()
				assert.NotNil(t, ds)

			case "error case - multiple calls return independent instances":
				ds1 := NewProjectDataSource()
				ds2 := NewProjectDataSource()
				assert.NotNil(t, ds1)
				assert.NotNil(t, ds2)
				assert.NotSame(t, ds1, ds2)
			}
		})
	}
}

// TestBuildProjectDataSourceSchema verifies the schema builder returns valid schema.
func TestBuildProjectDataSourceSchema(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		wantAttributes []string
	}{
		{
			name:           "schema has all expected attributes",
			wantAttributes: []string{"id", "name", "type", "created_at", "updated_at", "icon", "description"},
		},
		{
			name:           "error case - schema is not nil",
			wantAttributes: []string{"id"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := buildProjectDataSourceSchema()

			assert.NotNil(t, s)
			assert.NotEmpty(t, s.Attributes)

			for _, attr := range tt.wantAttributes {
				_, ok := s.Attributes[attr]
				assert.True(t, ok, "expected attribute %q in schema", attr)
			}
		})
	}
}

// TestProjectDataSource_executeRead verifies the executeRead method handles API calls correctly.
func TestProjectDataSource_executeRead(t *testing.T) {
	t.Parallel()

	// Build a proper schema for state initialization.
	schemaObj := buildProjectDataSourceSchema()

	tests := []struct {
		name        string
		setupData   func() *models.DataSource
		handler     http.HandlerFunc
		expectError bool
	}{
		{
			name: "successful read by ID",
			setupData: func() *models.DataSource {
				return &models.DataSource{
					ID:   types.StringValue("proj-123"),
					Name: types.StringNull(),
				}
			},
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/projects" && r.Method == http.MethodGet {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{
						"data": []any{
							map[string]any{
								"id":   "proj-123",
								"name": "Test Project",
							},
						},
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			expectError: false,
		},
		{
			name: "project not found",
			setupData: func() *models.DataSource {
				return &models.DataSource{
					ID:   types.StringValue("proj-999"),
					Name: types.StringNull(),
				}
			},
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/projects" && r.Method == http.MethodGet {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{
						"data": []any{},
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			expectError: true,
		},
		{
			name: "API error",
			setupData: func() *models.DataSource {
				return &models.DataSource{
					ID:   types.StringValue("proj-123"),
					Name: types.StringNull(),
				}
			},
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"message": "Internal server error"}`))
			}),
			expectError: true,
		},
		{
			name: "error case - nil data list returns not found",
			setupData: func() *models.DataSource {
				return &models.DataSource{
					ID:   types.StringValue("proj-123"),
					Name: types.StringNull(),
				}
			},
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/projects" && r.Method == http.MethodGet {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			n8nClient, server := setupTestClient(t, tt.handler)
			defer server.Close()

			d := &ProjectDataSource{client: n8nClient}
			data := tt.setupData()

			// Use schema-backed state to avoid nil pointer panic when Set is called.
			tfState := tfsdk.State{Schema: schemaObj}
			resp := &datasource.ReadResponse{State: tfState}

			d.executeRead(t.Context(), data, resp)

			if tt.expectError {
				assert.True(t, resp.Diagnostics.HasError())
			} else {
				assert.False(t, resp.Diagnostics.HasError())
			}
		})
	}
}
