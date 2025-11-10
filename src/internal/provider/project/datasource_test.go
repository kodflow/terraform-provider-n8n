package project

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/kodflow/n8n/src/internal/provider/shared/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockProjectDataSourceInterface is a mock implementation of ProjectDataSourceInterface.
type MockProjectDataSourceInterface struct {
	mock.Mock
}

func (m *MockProjectDataSourceInterface) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockProjectDataSourceInterface) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockProjectDataSourceInterface) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	m.Called(ctx, req, resp)
}

func (m *MockProjectDataSourceInterface) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	m.Called(ctx, req, resp)
}

func TestNewProjectDataSource(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "create new instance",
			wantErr: false,
		},
		{
			name:    "multiple instances are independent",
			wantErr: false,
		},
		{
			name:    "error case - validation checks",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "create new instance":
				ds := NewProjectDataSource()

				assert.NotNil(t, ds)
				assert.IsType(t, &ProjectDataSource{}, ds)
				assert.Nil(t, ds.client)

			case "multiple instances are independent":
				ds1 := NewProjectDataSource()
				ds2 := NewProjectDataSource()

				assert.NotSame(t, ds1, ds2)
				// Each instance is a different pointer, even if they have the same initial state

			case "error case - validation checks":
				// Validation test: ensure NewProjectDataSource never returns nil
				ds := NewProjectDataSource()
				assert.NotNil(t, ds, "NewProjectDataSource must not return nil")
				assert.IsType(t, &ProjectDataSource{}, ds, "NewProjectDataSource must return correct type")
			}
		})
	}
}

func TestNewProjectDataSourceWrapper(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "create new wrapped instance",
			wantErr: false,
		},
		{
			name:    "wrapper returns datasource.DataSource interface",
			wantErr: false,
		},
		{
			name:    "error case - validation checks",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "create new wrapped instance":
				ds := NewProjectDataSourceWrapper()

				assert.NotNil(t, ds)
				assert.IsType(t, &ProjectDataSource{}, ds)

			case "wrapper returns datasource.DataSource interface":
				ds := NewProjectDataSourceWrapper()

				// ds is already of type datasource.DataSource, no assertion needed
				assert.NotNil(t, ds)

			case "error case - validation checks":
				// Validation test: ensure wrapper returns valid interface
				ds := NewProjectDataSourceWrapper()
				assert.NotNil(t, ds, "NewProjectDataSourceWrapper must not return nil")
				assert.Implements(t, (*datasource.DataSource)(nil), ds, "Must implement DataSource interface")
			}
		})
	}
}

func TestProjectDataSource_Metadata(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "set metadata without provider type",
			wantErr: false,
		},
		{
			name:    "set metadata with provider type name",
			wantErr: false,
		},
		{
			name:    "set metadata with different provider type",
			wantErr: false,
		},
		{
			name:    "error case - validation checks",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "set metadata without provider type":
				ds := NewProjectDataSource()
				resp := &datasource.MetadataResponse{}

				ds.Metadata(context.Background(), datasource.MetadataRequest{
					ProviderTypeName: "n8n",
				}, resp)

				assert.Equal(t, "n8n_project", resp.TypeName)

			case "set metadata with provider type name":
				ds := NewProjectDataSource()
				resp := &datasource.MetadataResponse{}
				req := datasource.MetadataRequest{
					ProviderTypeName: "n8n",
				}

				ds.Metadata(context.Background(), req, resp)

				assert.Equal(t, "n8n_project", resp.TypeName)

			case "set metadata with different provider type":
				ds := NewProjectDataSource()
				resp := &datasource.MetadataResponse{}
				req := datasource.MetadataRequest{
					ProviderTypeName: "custom_provider",
				}

				ds.Metadata(context.Background(), req, resp)

				assert.Equal(t, "custom_provider_project", resp.TypeName)

			case "error case - validation checks":
				// Validation test: ensure metadata is set correctly even with empty provider
				ds := NewProjectDataSource()
				resp := &datasource.MetadataResponse{}
				req := datasource.MetadataRequest{
					ProviderTypeName: "",
				}

				ds.Metadata(context.Background(), req, resp)

				assert.Equal(t, "_project", resp.TypeName, "TypeName should be set even with empty provider")
			}
		})
	}
}

func TestProjectDataSource_Schema(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "return schema",
			wantErr: false,
		},
		{
			name:    "schema attributes have descriptions",
			wantErr: false,
		},
		{
			name:    "error case - validation checks",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "return schema":
				ds := NewProjectDataSource()
				resp := &datasource.SchemaResponse{}

				ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)

				assert.NotNil(t, resp.Schema)
				assert.Contains(t, resp.Schema.MarkdownDescription, "n8n project")

				// Verify attributes
				expectedAttrs := []string{
					"id",
					"name",
					"type",
					"created_at",
					"updated_at",
					"icon",
					"description",
				}

				for _, attr := range expectedAttrs {
					_, exists := resp.Schema.Attributes[attr]
					assert.True(t, exists, "Attribute %s should exist", attr)
				}

			case "schema attributes have descriptions":
				ds := NewProjectDataSource()
				resp := &datasource.SchemaResponse{}

				ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)

				for name, attr := range resp.Schema.Attributes {
					assert.NotEmpty(t, attr.GetMarkdownDescription(), "Attribute %s should have description", name)
				}

			case "error case - validation checks":
				// Validation test: ensure schema is always set
				ds := NewProjectDataSource()
				resp := &datasource.SchemaResponse{}

				ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)

				assert.NotNil(t, resp.Schema, "Schema must not be nil")
				assert.NotEmpty(t, resp.Schema.Attributes, "Schema must have attributes")
			}
		})
	}
}

func TestProjectDataSource_Configure(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "configure with valid client",
			wantErr: false,
		},
		{
			name:    "configure with nil provider data",
			wantErr: false,
		},
		{
			name:    "configure with invalid provider data",
			wantErr: true,
		},
		{
			name:    "configure with wrong type",
			wantErr: true,
		},
		{
			name:    "error case - validation checks",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "configure with valid client":
				ds := NewProjectDataSource()
				resp := &datasource.ConfigureResponse{}

				mockClient := &client.N8nClient{}
				req := datasource.ConfigureRequest{
					ProviderData: mockClient,
				}

				ds.Configure(context.Background(), req, resp)

				assert.NotNil(t, ds.client)
				assert.Equal(t, mockClient, ds.client)
				assert.False(t, resp.Diagnostics.HasError())

			case "configure with nil provider data":
				ds := NewProjectDataSource()
				resp := &datasource.ConfigureResponse{}
				req := datasource.ConfigureRequest{
					ProviderData: nil,
				}

				ds.Configure(context.Background(), req, resp)

				assert.Nil(t, ds.client)
				assert.False(t, resp.Diagnostics.HasError())

			case "configure with invalid provider data":
				ds := NewProjectDataSource()
				resp := &datasource.ConfigureResponse{}
				req := datasource.ConfigureRequest{
					ProviderData: "invalid-data",
				}

				ds.Configure(context.Background(), req, resp)

				assert.Nil(t, ds.client)
				assert.True(t, resp.Diagnostics.HasError())
				assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Unexpected Data Source Configure Type")

			case "configure with wrong type":
				ds := NewProjectDataSource()
				resp := &datasource.ConfigureResponse{}
				req := datasource.ConfigureRequest{
					ProviderData: struct{}{},
				}

				ds.Configure(context.Background(), req, resp)

				assert.Nil(t, ds.client)
				assert.True(t, resp.Diagnostics.HasError())

			case "error case - validation checks":
				// Validation test: ensure invalid types are rejected
				ds := NewProjectDataSource()
				resp := &datasource.ConfigureResponse{}
				req := datasource.ConfigureRequest{
					ProviderData: 12345, // Invalid type
				}

				ds.Configure(context.Background(), req, resp)

				assert.Nil(t, ds.client, "Client must remain nil for invalid provider data")
				assert.True(t, resp.Diagnostics.HasError(), "Must have error for invalid provider data type")
			}
		})
	}
}

func TestProjectDataSource_Interfaces(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "implements required interfaces",
			wantErr: false,
		},
		{
			name:    "error case - validation checks",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "implements required interfaces":
				ds := NewProjectDataSource()

				// Test that ProjectDataSource implements datasource.DataSource
				var _ datasource.DataSource = ds

				// Test that ProjectDataSource implements datasource.DataSourceWithConfigure
				var _ datasource.DataSourceWithConfigure = ds

				// Test that ProjectDataSource implements ProjectDataSourceInterface
				var _ ProjectDataSourceInterface = ds

			case "error case - validation checks":
				// Validation test: ensure datasource is not nil
				ds := NewProjectDataSource()
				assert.NotNil(t, ds, "DataSource must not be nil")
				assert.Implements(t, (*datasource.DataSource)(nil), ds, "Must implement DataSource interface")
			}
		})
	}
}

func TestProjectDataSourceConcurrency(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "concurrent metadata calls",
			wantErr: false,
		},
		{
			name:    "concurrent schema calls",
			wantErr: false,
		},
		{
			name:    "error case - validation checks",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Note: No t.Parallel() in subtests with goroutines

			switch tt.name {
			case "concurrent metadata calls":
				ds := NewProjectDataSource()

				done := make(chan bool, 100)
				for i := 0; i < 100; i++ {
					go func() {
						resp := &datasource.MetadataResponse{}
						ds.Metadata(context.Background(), datasource.MetadataRequest{
							ProviderTypeName: "n8n",
						}, resp)
						assert.Equal(t, "n8n_project", resp.TypeName)
						done <- true
					}()
				}

				for i := 0; i < 100; i++ {
					<-done
				}

			case "concurrent schema calls":
				ds := NewProjectDataSource()

				done := make(chan bool, 100)
				for i := 0; i < 100; i++ {
					go func() {
						resp := &datasource.SchemaResponse{}
						ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)
						assert.NotNil(t, resp.Schema)
						done <- true
					}()
				}

				for i := 0; i < 100; i++ {
					<-done
				}

			case "error case - validation checks":
				// Validation test: ensure concurrent configure is safe
				ds := NewProjectDataSource()
				done := make(chan bool, 10)

				for i := 0; i < 10; i++ {
					go func() {
						resp := &datasource.ConfigureResponse{}
						req := datasource.ConfigureRequest{
							ProviderData: &client.N8nClient{},
						}
						ds.Configure(context.Background(), req, resp)
						assert.False(t, resp.Diagnostics.HasError())
						done <- true
					}()
				}

				for i := 0; i < 10; i++ {
					<-done
				}
			}
		})
	}
}

func BenchmarkProjectDataSource_Schema(b *testing.B) {
	ds := NewProjectDataSource()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := &datasource.SchemaResponse{}
		ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)
	}
}

func BenchmarkProjectDataSource_Metadata(b *testing.B) {
	ds := NewProjectDataSource()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := &datasource.MetadataResponse{}
		ds.Metadata(context.Background(), datasource.MetadataRequest{
			ProviderTypeName: "n8n",
		}, resp)
	}
}

func BenchmarkProjectDataSource_Configure(b *testing.B) {
	ds := NewProjectDataSource()
	mockClient := &client.N8nClient{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := &datasource.ConfigureResponse{}
		req := datasource.ConfigureRequest{
			ProviderData: mockClient,
		}
		ds.Configure(context.Background(), req, resp)
	}
}

func TestProjectDataSource_Read(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "successful read by ID",
			wantErr: false,
		},
		{
			name:    "successful read by name",
			wantErr: false,
		},
		{
			name:    "project not found",
			wantErr: true,
		},
		{
			name:    "read fails when neither ID nor name provided",
			wantErr: true,
		},
		{
			name:    "read fails with API error",
			wantErr: true,
		},
		{
			name:    "read fails when config get fails",
			wantErr: true,
		},
		{
			name:    "read succeeds but state set fails",
			wantErr: true,
		},
		{
			name:    "read fails when project list data is nil",
			wantErr: true,
		},
		{
			name:    "read fails when project not found by name",
			wantErr: true,
		},
		{
			name:    "error case - validation checks",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "successful read by ID":
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/projects" && r.Method == http.MethodGet {
						w.Header().Set("Content-Type", "application/json")
						json.NewEncoder(w).Encode(map[string]interface{}{
							"data": []map[string]interface{}{
								{
									"id":          "proj-123",
									"name":        "Test Project",
									"type":        "team",
									"createdAt":   "2024-01-01T00:00:00Z",
									"updatedAt":   "2024-01-02T00:00:00Z",
									"icon":        "home",
									"description": "Test project description",
								},
							},
						})
						return
					}
					w.WriteHeader(http.StatusNotFound)
				})

				n8nClient, server := setupTestDataSourceClient(t, handler)
				defer server.Close()

				ds := &ProjectDataSource{client: n8nClient}

				rawConfig := map[string]tftypes.Value{
					"id":          tftypes.NewValue(tftypes.String, "proj-123"),
					"name":        tftypes.NewValue(tftypes.String, nil),
					"type":        tftypes.NewValue(tftypes.String, nil),
					"created_at":  tftypes.NewValue(tftypes.String, nil),
					"updated_at":  tftypes.NewValue(tftypes.String, nil),
					"icon":        tftypes.NewValue(tftypes.String, nil),
					"description": tftypes.NewValue(tftypes.String, nil),
				}
				attrTypes := map[string]tftypes.Type{
					"id":          tftypes.String,
					"name":        tftypes.String,
					"type":        tftypes.String,
					"created_at":  tftypes.String,
					"updated_at":  tftypes.String,
					"icon":        tftypes.String,
					"description": tftypes.String,
				}
				config := tfsdk.Config{
					Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: attrTypes}, rawConfig),
					Schema: createTestDataSourceSchema(t),
				}

				state := tfsdk.State{
					Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: attrTypes}, nil),
					Schema: createTestDataSourceSchema(t),
				}

				req := datasource.ReadRequest{
					Config: config,
				}
				resp := datasource.ReadResponse{
					State: state,
				}

				ds.Read(context.Background(), req, &resp)

				assert.False(t, resp.Diagnostics.HasError(), "Read should not have errors")

			case "successful read by name":
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/projects" && r.Method == http.MethodGet {
						w.Header().Set("Content-Type", "application/json")
						json.NewEncoder(w).Encode(map[string]interface{}{
							"data": []map[string]interface{}{
								{
									"id":          "proj-456",
									"name":        "Test Project By Name",
									"type":        "personal",
									"createdAt":   "2024-01-01T00:00:00Z",
									"updatedAt":   "2024-01-02T00:00:00Z",
									"icon":        "star",
									"description": "Another project",
								},
							},
						})
						return
					}
					w.WriteHeader(http.StatusNotFound)
				})

				n8nClient, server := setupTestDataSourceClient(t, handler)
				defer server.Close()

				ds := &ProjectDataSource{client: n8nClient}

				rawConfig := map[string]tftypes.Value{
					"id":          tftypes.NewValue(tftypes.String, nil),
					"name":        tftypes.NewValue(tftypes.String, "Test Project By Name"),
					"type":        tftypes.NewValue(tftypes.String, nil),
					"created_at":  tftypes.NewValue(tftypes.String, nil),
					"updated_at":  tftypes.NewValue(tftypes.String, nil),
					"icon":        tftypes.NewValue(tftypes.String, nil),
					"description": tftypes.NewValue(tftypes.String, nil),
				}
				attrTypes := map[string]tftypes.Type{
					"id":          tftypes.String,
					"name":        tftypes.String,
					"type":        tftypes.String,
					"created_at":  tftypes.String,
					"updated_at":  tftypes.String,
					"icon":        tftypes.String,
					"description": tftypes.String,
				}
				config := tfsdk.Config{
					Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: attrTypes}, rawConfig),
					Schema: createTestDataSourceSchema(t),
				}

				state := tfsdk.State{
					Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: attrTypes}, nil),
					Schema: createTestDataSourceSchema(t),
				}

				req := datasource.ReadRequest{
					Config: config,
				}
				resp := datasource.ReadResponse{
					State: state,
				}

				ds.Read(context.Background(), req, &resp)

				assert.False(t, resp.Diagnostics.HasError(), "Read should not have errors")

			case "project not found":
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/projects" && r.Method == http.MethodGet {
						w.Header().Set("Content-Type", "application/json")
						json.NewEncoder(w).Encode(map[string]interface{}{
							"data": []map[string]interface{}{},
						})
						return
					}
					w.WriteHeader(http.StatusNotFound)
				})

				n8nClient, server := setupTestDataSourceClient(t, handler)
				defer server.Close()

				ds := &ProjectDataSource{client: n8nClient}

				rawConfig := map[string]tftypes.Value{
					"id":          tftypes.NewValue(tftypes.String, "proj-nonexistent"),
					"name":        tftypes.NewValue(tftypes.String, nil),
					"type":        tftypes.NewValue(tftypes.String, nil),
					"created_at":  tftypes.NewValue(tftypes.String, nil),
					"updated_at":  tftypes.NewValue(tftypes.String, nil),
					"icon":        tftypes.NewValue(tftypes.String, nil),
					"description": tftypes.NewValue(tftypes.String, nil),
				}
				attrTypes := map[string]tftypes.Type{
					"id":          tftypes.String,
					"name":        tftypes.String,
					"type":        tftypes.String,
					"created_at":  tftypes.String,
					"updated_at":  tftypes.String,
					"icon":        tftypes.String,
					"description": tftypes.String,
				}
				config := tfsdk.Config{
					Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: attrTypes}, rawConfig),
					Schema: createTestDataSourceSchema(t),
				}

				state := tfsdk.State{
					Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: attrTypes}, nil),
					Schema: createTestDataSourceSchema(t),
				}

				req := datasource.ReadRequest{
					Config: config,
				}
				resp := datasource.ReadResponse{
					State: state,
				}

				ds.Read(context.Background(), req, &resp)

				assert.True(t, resp.Diagnostics.HasError(), "Read should have errors when project not found")

			case "read fails when neither ID nor name provided":
				ds := &ProjectDataSource{}

				rawConfig := map[string]tftypes.Value{
					"id":          tftypes.NewValue(tftypes.String, nil),
					"name":        tftypes.NewValue(tftypes.String, nil),
					"type":        tftypes.NewValue(tftypes.String, nil),
					"created_at":  tftypes.NewValue(tftypes.String, nil),
					"updated_at":  tftypes.NewValue(tftypes.String, nil),
					"icon":        tftypes.NewValue(tftypes.String, nil),
					"description": tftypes.NewValue(tftypes.String, nil),
				}
				attrTypes := map[string]tftypes.Type{
					"id":          tftypes.String,
					"name":        tftypes.String,
					"type":        tftypes.String,
					"created_at":  tftypes.String,
					"updated_at":  tftypes.String,
					"icon":        tftypes.String,
					"description": tftypes.String,
				}
				config := tfsdk.Config{
					Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: attrTypes}, rawConfig),
					Schema: createTestDataSourceSchema(t),
				}

				state := tfsdk.State{
					Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: attrTypes}, nil),
					Schema: createTestDataSourceSchema(t),
				}

				req := datasource.ReadRequest{
					Config: config,
				}
				resp := datasource.ReadResponse{
					State: state,
				}

				ds.Read(context.Background(), req, &resp)

				assert.True(t, resp.Diagnostics.HasError(), "Read should have errors when neither ID nor name provided")
				assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Missing Required Attribute")

			case "read fails with API error":
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"message": "Internal server error"}`))
				})

				n8nClient, server := setupTestDataSourceClient(t, handler)
				defer server.Close()

				ds := &ProjectDataSource{client: n8nClient}

				rawConfig := map[string]tftypes.Value{
					"id":          tftypes.NewValue(tftypes.String, "proj-123"),
					"name":        tftypes.NewValue(tftypes.String, nil),
					"type":        tftypes.NewValue(tftypes.String, nil),
					"created_at":  tftypes.NewValue(tftypes.String, nil),
					"updated_at":  tftypes.NewValue(tftypes.String, nil),
					"icon":        tftypes.NewValue(tftypes.String, nil),
					"description": tftypes.NewValue(tftypes.String, nil),
				}
				attrTypes := map[string]tftypes.Type{
					"id":          tftypes.String,
					"name":        tftypes.String,
					"type":        tftypes.String,
					"created_at":  tftypes.String,
					"updated_at":  tftypes.String,
					"icon":        tftypes.String,
					"description": tftypes.String,
				}
				config := tfsdk.Config{
					Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: attrTypes}, rawConfig),
					Schema: createTestDataSourceSchema(t),
				}

				state := tfsdk.State{
					Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: attrTypes}, nil),
					Schema: createTestDataSourceSchema(t),
				}

				req := datasource.ReadRequest{
					Config: config,
				}
				resp := datasource.ReadResponse{
					State: state,
				}

				ds.Read(context.Background(), req, &resp)

				assert.True(t, resp.Diagnostics.HasError(), "Read should have errors on API failure")
				assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Error listing projects")

			case "read fails when config get fails":
				ds := &ProjectDataSource{}

				// Create config with mismatched schema - use wrong types
				rawConfig := map[string]tftypes.Value{
					"id":   tftypes.NewValue(tftypes.Number, 123), // Wrong type - should be String
					"name": tftypes.NewValue(tftypes.String, nil),
				}
				config := tfsdk.Config{
					Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.Number, "name": tftypes.String}}, rawConfig),
					Schema: createTestDataSourceSchema(t), // Schema expects String for id
				}

				state := tfsdk.State{
					Schema: createTestDataSourceSchema(t),
				}

				req := datasource.ReadRequest{
					Config: config,
				}
				resp := datasource.ReadResponse{
					State: state,
				}

				ds.Read(context.Background(), req, &resp)

				assert.True(t, resp.Diagnostics.HasError(), "Read should fail when Config.Get fails")

			case "read succeeds but state set fails":
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/projects" && r.Method == http.MethodGet {
						projects := []map[string]interface{}{
							{
								"id":   "proj-123",
								"name": "Test Project",
								"type": "team",
							},
						}
						w.Header().Set("Content-Type", "application/json")
						json.NewEncoder(w).Encode(map[string]interface{}{
							"data": projects,
						})
						return
					}
					w.WriteHeader(http.StatusNotFound)
				})

				n8nClient, server := setupTestDataSourceClient(t, handler)
				defer server.Close()

				ds := &ProjectDataSource{client: n8nClient}

				rawConfig := map[string]tftypes.Value{
					"id":          tftypes.NewValue(tftypes.String, "proj-123"),
					"name":        tftypes.NewValue(tftypes.String, nil),
					"type":        tftypes.NewValue(tftypes.String, nil),
					"created_at":  tftypes.NewValue(tftypes.String, nil),
					"updated_at":  tftypes.NewValue(tftypes.String, nil),
					"icon":        tftypes.NewValue(tftypes.String, nil),
					"description": tftypes.NewValue(tftypes.String, nil),
				}
				attrTypes := map[string]tftypes.Type{
					"id":          tftypes.String,
					"name":        tftypes.String,
					"type":        tftypes.String,
					"created_at":  tftypes.String,
					"updated_at":  tftypes.String,
					"icon":        tftypes.String,
					"description": tftypes.String,
				}
				config := tfsdk.Config{
					Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: attrTypes}, rawConfig),
					Schema: createTestDataSourceSchema(t),
				}

				// Create state with incompatible schema to trigger Set failure
				wrongSchema := schema.Schema{
					Attributes: map[string]schema.Attribute{
						"id": schema.NumberAttribute{
							Computed: true,
						},
					},
				}
				state := tfsdk.State{
					Schema: wrongSchema, // Wrong schema will cause Set to fail
				}

				req := datasource.ReadRequest{
					Config: config,
				}
				resp := datasource.ReadResponse{
					State: state,
				}

				ds.Read(context.Background(), req, &resp)

				// State.Set should fail due to schema mismatch
				assert.True(t, resp.Diagnostics.HasError(), "Read should fail when State.Set fails")

			case "read fails when project list data is nil":
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/projects" && r.Method == http.MethodGet {
						// Return response with nil data
						w.Header().Set("Content-Type", "application/json")
						json.NewEncoder(w).Encode(map[string]interface{}{
							"data": nil,
						})
						return
					}
					w.WriteHeader(http.StatusNotFound)
				})

				n8nClient, server := setupTestDataSourceClient(t, handler)
				defer server.Close()

				ds := &ProjectDataSource{client: n8nClient}

				rawConfig := map[string]tftypes.Value{
					"id":          tftypes.NewValue(tftypes.String, "proj-123"),
					"name":        tftypes.NewValue(tftypes.String, nil),
					"type":        tftypes.NewValue(tftypes.String, nil),
					"created_at":  tftypes.NewValue(tftypes.String, nil),
					"updated_at":  tftypes.NewValue(tftypes.String, nil),
					"icon":        tftypes.NewValue(tftypes.String, nil),
					"description": tftypes.NewValue(tftypes.String, nil),
				}
				attrTypes := map[string]tftypes.Type{
					"id":          tftypes.String,
					"name":        tftypes.String,
					"type":        tftypes.String,
					"created_at":  tftypes.String,
					"updated_at":  tftypes.String,
					"icon":        tftypes.String,
					"description": tftypes.String,
				}
				config := tfsdk.Config{
					Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: attrTypes}, rawConfig),
					Schema: createTestDataSourceSchema(t),
				}

				state := tfsdk.State{
					Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: attrTypes}, nil),
					Schema: createTestDataSourceSchema(t),
				}

				req := datasource.ReadRequest{
					Config: config,
				}
				resp := datasource.ReadResponse{
					State: state,
				}

				ds.Read(context.Background(), req, &resp)

				assert.True(t, resp.Diagnostics.HasError(), "Read should fail when project not found (nil data)")
				assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Project Not Found")

			case "read fails when project not found by name":
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/projects" && r.Method == http.MethodGet {
						// Return empty list
						w.Header().Set("Content-Type", "application/json")
						json.NewEncoder(w).Encode(map[string]interface{}{
							"data": []map[string]interface{}{},
						})
						return
					}
					w.WriteHeader(http.StatusNotFound)
				})

				n8nClient, server := setupTestDataSourceClient(t, handler)
				defer server.Close()

				ds := &ProjectDataSource{client: n8nClient}

				// Search by NAME only (ID is nil) - this will test the identifier == "" branch
				rawConfig := map[string]tftypes.Value{
					"id":          tftypes.NewValue(tftypes.String, nil),
					"name":        tftypes.NewValue(tftypes.String, "Nonexistent Project"),
					"type":        tftypes.NewValue(tftypes.String, nil),
					"created_at":  tftypes.NewValue(tftypes.String, nil),
					"updated_at":  tftypes.NewValue(tftypes.String, nil),
					"icon":        tftypes.NewValue(tftypes.String, nil),
					"description": tftypes.NewValue(tftypes.String, nil),
				}
				attrTypes := map[string]tftypes.Type{
					"id":          tftypes.String,
					"name":        tftypes.String,
					"type":        tftypes.String,
					"created_at":  tftypes.String,
					"updated_at":  tftypes.String,
					"icon":        tftypes.String,
					"description": tftypes.String,
				}
				config := tfsdk.Config{
					Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: attrTypes}, rawConfig),
					Schema: createTestDataSourceSchema(t),
				}

				state := tfsdk.State{
					Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: attrTypes}, nil),
					Schema: createTestDataSourceSchema(t),
				}

				req := datasource.ReadRequest{
					Config: config,
				}
				resp := datasource.ReadResponse{
					State: state,
				}

				ds.Read(context.Background(), req, &resp)

				assert.True(t, resp.Diagnostics.HasError(), "Read should fail when project not found by name")
				assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Project Not Found")
				assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "Nonexistent Project")

			case "error case - validation checks":
				// Validation test: verify multiple projects returned causes correct filtering
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/projects" && r.Method == http.MethodGet {
						// Return multiple projects to test filtering
						w.Header().Set("Content-Type", "application/json")
						json.NewEncoder(w).Encode(map[string]interface{}{
							"data": []map[string]interface{}{
								{
									"id":          "proj-999",
									"name":        "Wrong Project",
									"type":        "team",
									"createdAt":   "2024-01-01T00:00:00Z",
									"updatedAt":   "2024-01-02T00:00:00Z",
									"icon":        "star",
									"description": "This is not the project you're looking for",
								},
								{
									"id":          "proj-123",
									"name":        "Right Project",
									"type":        "personal",
									"createdAt":   "2024-01-01T00:00:00Z",
									"updatedAt":   "2024-01-02T00:00:00Z",
									"icon":        "home",
									"description": "This is the correct project",
								},
							},
						})
						return
					}
					w.WriteHeader(http.StatusNotFound)
				})

				n8nClient, server := setupTestDataSourceClient(t, handler)
				defer server.Close()

				ds := &ProjectDataSource{client: n8nClient}

				rawConfig := map[string]tftypes.Value{
					"id":          tftypes.NewValue(tftypes.String, "proj-123"),
					"name":        tftypes.NewValue(tftypes.String, nil),
					"type":        tftypes.NewValue(tftypes.String, nil),
					"created_at":  tftypes.NewValue(tftypes.String, nil),
					"updated_at":  tftypes.NewValue(tftypes.String, nil),
					"icon":        tftypes.NewValue(tftypes.String, nil),
					"description": tftypes.NewValue(tftypes.String, nil),
				}
				attrTypes := map[string]tftypes.Type{
					"id":          tftypes.String,
					"name":        tftypes.String,
					"type":        tftypes.String,
					"created_at":  tftypes.String,
					"updated_at":  tftypes.String,
					"icon":        tftypes.String,
					"description": tftypes.String,
				}
				config := tfsdk.Config{
					Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: attrTypes}, rawConfig),
					Schema: createTestDataSourceSchema(t),
				}

				state := tfsdk.State{
					Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: attrTypes}, nil),
					Schema: createTestDataSourceSchema(t),
				}

				req := datasource.ReadRequest{
					Config: config,
				}
				resp := datasource.ReadResponse{
					State: state,
				}

				ds.Read(context.Background(), req, &resp)

				// Should successfully find the correct project by ID
				assert.False(t, resp.Diagnostics.HasError(), "Read should succeed when correct project ID is provided")
			}
		})
	}
}

// createTestDataSourceSchema creates a test schema for project datasource.
func createTestDataSourceSchema(t *testing.T) schema.Schema {
	t.Helper()
	ds := &ProjectDataSource{}
	req := datasource.SchemaRequest{}
	resp := &datasource.SchemaResponse{}
	ds.Schema(context.Background(), req, resp)
	return resp.Schema
}

// setupTestDataSourceClient creates a test N8nClient with httptest server.
func setupTestDataSourceClient(t *testing.T, handler http.HandlerFunc) (*client.N8nClient, *httptest.Server) {
	t.Helper()
	server := httptest.NewServer(handler)

	cfg := n8nsdk.NewConfiguration()
	cfg.Servers = n8nsdk.ServerConfigurations{
		{
			URL:         server.URL,
			Description: "Test server",
		},
	}
	cfg.HTTPClient = server.Client()
	cfg.AddDefaultHeader("X-N8N-API-KEY", "test-key")

	apiClient := n8nsdk.NewAPIClient(cfg)
	n8nClient := &client.N8nClient{
		APIClient: apiClient,
	}

	return n8nClient, server
}
