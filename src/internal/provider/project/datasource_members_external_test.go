// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

package project_test

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/project"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/client"
	"github.com/stretchr/testify/assert"
)

// TestNewProjectMembersDataSource tests the NewProjectMembersDataSource constructor.
func TestNewProjectMembersDataSource(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "creates valid datasource",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := project.NewProjectMembersDataSource()
				assert.NotNil(t, ds)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.testFunc(t)
		})
	}
}

// TestNewProjectMembersDataSourceWrapper tests the NewProjectMembersDataSourceWrapper constructor.
func TestNewProjectMembersDataSourceWrapper(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "creates valid datasource",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := project.NewProjectMembersDataSourceWrapper()
				assert.NotNil(t, ds)
				assert.Implements(t, (*datasource.DataSource)(nil), ds)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.testFunc(t)
		})
	}
}

// TestProjectMembersDataSource_Metadata tests the Metadata method.
func TestProjectMembersDataSource_Metadata(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "sets correct type name",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := project.NewProjectMembersDataSource()
				resp := &datasource.MetadataResponse{}

				ds.Metadata(t.Context(), datasource.MetadataRequest{
					ProviderTypeName: "n8n",
				}, resp)

				assert.Equal(t, "n8n_project_members", resp.TypeName)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.testFunc(t)
		})
	}
}

// TestProjectMembersDataSource_Schema tests the Schema method.
func TestProjectMembersDataSource_Schema(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "returns valid schema",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := project.NewProjectMembersDataSource()
				resp := &datasource.SchemaResponse{}

				ds.Schema(t.Context(), datasource.SchemaRequest{}, resp)

				assert.NotNil(t, resp.Schema)
				assert.Contains(t, resp.Schema.Attributes, "project_id")
				assert.Contains(t, resp.Schema.Attributes, "members")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.testFunc(t)
		})
	}
}

// TestProjectMembersDataSource_Configure tests the Configure method.
func TestProjectMembersDataSource_Configure(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		providerData any
		wantError    bool
	}{
		{
			name:         "valid configuration",
			providerData: &client.N8nClient{},
			wantError:    false,
		},
		{
			name:         "nil provider data",
			providerData: nil,
			wantError:    false,
		},
		{
			name:         "invalid provider data type",
			providerData: "invalid",
			wantError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ds := project.NewProjectMembersDataSource()
			resp := &datasource.ConfigureResponse{}
			req := datasource.ConfigureRequest{
				ProviderData: tt.providerData,
			}

			ds.Configure(t.Context(), req, resp)

			if tt.wantError {
				assert.True(t, resp.Diagnostics.HasError())
			} else {
				assert.False(t, resp.Diagnostics.HasError())
			}
		})
	}
}

// TestProjectMembersDataSource_Read tests the Read method.
func TestProjectMembersDataSource_Read(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "successful read",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{"data":[{"id":"user-1","email":"alice@example.com","firstName":"Alice","lastName":"Smith","role":"project:admin"}]}`))
				})

				n8nClient, server := setupTestClientForDataSource(t, handler)
				defer server.Close()

				ds := project.NewProjectMembersDataSource()
				ds.Configure(t.Context(), datasource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &datasource.ConfigureResponse{})

				ctx := t.Context()
				schemaResp := datasource.SchemaResponse{}
				ds.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)

				configRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"project_id": tftypes.NewValue(tftypes.String, "proj-123"),
					"members": tftypes.NewValue(tftypes.List{ElementType: tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"user_id":    tftypes.String,
							"role":       tftypes.String,
							"email":      tftypes.String,
							"first_name": tftypes.String,
							"last_name":  tftypes.String,
						},
					}}, nil),
				})

				req := datasource.ReadRequest{
					Config: tfsdk.Config{
						Schema: schemaResp.Schema,
						Raw:    configRaw,
					},
				}
				resp := &datasource.ReadResponse{
					State: tfsdk.State{
						Schema: schemaResp.Schema,
					},
				}

				ds.Read(ctx, req, resp)

				assert.False(t, resp.Diagnostics.HasError(), "Expected no errors")
			},
		},
		{
			name: "error - invalid config",
			testFunc: func(t *testing.T) {
				t.Helper()
				ds := project.NewProjectMembersDataSource()
				ctx := t.Context()
				schemaResp := datasource.SchemaResponse{}
				ds.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)

				configRaw := tftypes.NewValue(tftypes.String, "invalid")

				req := datasource.ReadRequest{
					Config: tfsdk.Config{
						Schema: schemaResp.Schema,
						Raw:    configRaw,
					},
				}
				resp := &datasource.ReadResponse{
					State: tfsdk.State{
						Schema: schemaResp.Schema,
					},
				}

				ds.Read(ctx, req, resp)

				assert.True(t, resp.Diagnostics.HasError())
			},
		},
		{
			name: "error - api error",
			testFunc: func(t *testing.T) {
				t.Helper()
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"message":"internal error"}`))
				})

				n8nClient, server := setupTestClientForDataSource(t, handler)
				defer server.Close()

				ds := project.NewProjectMembersDataSource()
				ds.Configure(t.Context(), datasource.ConfigureRequest{
					ProviderData: n8nClient,
				}, &datasource.ConfigureResponse{})

				ctx := t.Context()
				schemaResp := datasource.SchemaResponse{}
				ds.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)

				configRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
					"project_id": tftypes.NewValue(tftypes.String, "proj-123"),
					"members": tftypes.NewValue(tftypes.List{ElementType: tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"user_id":    tftypes.String,
							"role":       tftypes.String,
							"email":      tftypes.String,
							"first_name": tftypes.String,
							"last_name":  tftypes.String,
						},
					}}, nil),
				})

				req := datasource.ReadRequest{
					Config: tfsdk.Config{
						Schema: schemaResp.Schema,
						Raw:    configRaw,
					},
				}
				resp := &datasource.ReadResponse{
					State: tfsdk.State{
						Schema: schemaResp.Schema,
					},
				}

				ds.Read(ctx, req, resp)

				assert.True(t, resp.Diagnostics.HasError())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t)
		})
	}
}
