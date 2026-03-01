// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

package datatable

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/shared/client"
	"github.com/stretchr/testify/assert"
)

func TestDataTablesDataSource_Read(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectError bool
		expectCount int
	}{
		{
			name: "successful list with data",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/data-tables" && r.Method == http.MethodGet {
					response := map[string]interface{}{
						"data": []interface{}{
							map[string]interface{}{
								"id":        "dt-1",
								"name":      "Table One",
								"projectId": "proj-1",
								"columns":   []interface{}{},
								"createdAt": "2024-01-01T00:00:00Z",
								"updatedAt": "2024-01-01T00:00:00Z",
							},
							map[string]interface{}{
								"id":        "dt-2",
								"name":      "Table Two",
								"projectId": "proj-1",
								"columns":   []interface{}{},
								"createdAt": "2024-01-02T00:00:00Z",
								"updatedAt": "2024-01-02T00:00:00Z",
							},
						},
					}
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(response)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			expectError: false,
			expectCount: 2,
		},
		{
			name: "successful list with empty data",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/data-tables" && r.Method == http.MethodGet {
					response := map[string]interface{}{
						"data": []interface{}{},
					}
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(response)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			expectError: false,
			expectCount: 0,
		},
		{
			name: "API returns error",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/data-tables" && r.Method == http.MethodGet {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"message": "Internal server error"}`))
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			expectError: true,
		},
		{
			name: "API returns not found",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/data-tables" && r.Method == http.MethodGet {
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte(`{"message": "Not found"}`))
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			expectError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(tt.handler)
			defer server.Close()

			cfg := n8nsdk.NewConfiguration()
			cfg.Servers = n8nsdk.ServerConfigurations{
				{URL: server.URL, Description: "Test server"},
			}
			cfg.HTTPClient = server.Client()
			cfg.AddDefaultHeader("X-N8N-API-KEY", "test-key")
			n8nClient := &client.N8nClient{APIClient: n8nsdk.NewAPIClient(cfg)}

			d := &DataTablesDataSource{client: n8nClient}
			var schemaResp datasource.SchemaResponse
			d.Schema(context.Background(), datasource.SchemaRequest{}, &schemaResp)
			req := datasource.ReadRequest{}
			resp := datasource.ReadResponse{}
			resp.State = tfsdk.State{Schema: schemaResp.Schema}

			d.Read(context.Background(), req, &resp)

			if tt.expectError {
				assert.True(t, resp.Diagnostics.HasError())
				return
			}

			assert.False(t, resp.Diagnostics.HasError(), "Read should not report diagnostics for success cases")
		})
	}
}
