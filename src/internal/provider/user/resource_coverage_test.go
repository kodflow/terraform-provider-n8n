package user

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/n8n/src/internal/provider/user/models"
	"github.com/stretchr/testify/assert"
)

// TestUserResource_Create_FullCoverage ensures 100% coverage.
func TestUserResource_Create_FullCoverage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		setupHandler func(w http.ResponseWriter, r *http.Request)
		email        string
		role         types.String
		expectError  bool
	}{
		{
			name: "create user with nil user in response",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPost && r.URL.Path == "/users" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{
						"user": nil,
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			email:       "test@example.com",
			role:        types.StringNull(),
			expectError: true,
		},
		{
			name: "create user with nil ID in response",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPost && r.URL.Path == "/users" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{
						"user": map[string]any{
							"email": "test@example.com",
							"id":    nil,
						},
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			email:       "test@example.com",
			role:        types.StringNull(),
			expectError: true,
		},
		{
			name: "create user with get error",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPost && r.URL.Path == "/users" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{
						"user": map[string]any{
							"email": "test@example.com",
							"id":    "user-1",
						},
					})
					return
				}
				if r.Method == http.MethodGet {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			email:       "test@example.com",
			role:        types.StringNull(),
			expectError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := http.HandlerFunc(tt.setupHandler)
			n8nClient, server := setupTestClient(t, handler)
			defer server.Close()

			r := &UserResource{client: n8nClient}
			plan := &models.Resource{
				Email: types.StringValue(tt.email),
				Role:  tt.role,
			}

			req := resource.CreateRequest{
				Plan: tfsdk.Plan{Schema: createTestSchema(t)},
			}
			resp := resource.CreateResponse{
				State: tfsdk.State{Schema: createTestSchema(t)},
			}

			req.Plan.Set(context.Background(), plan)

			r.Create(context.Background(), req, &resp)

			if tt.expectError {
				assert.True(t, resp.Diagnostics.HasError())
			} else {
				assert.False(t, resp.Diagnostics.HasError())
			}
		})
	}
}

// TestUserResource_Update_FullCoverage ensures 100% coverage.
func TestUserResource_Update_FullCoverage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		setupHandler func(w http.ResponseWriter, r *http.Request)
		planEmail    string
		stateEmail   string
		planRole     types.String
		stateRole    types.String
		expectError  bool
	}{
		{
			name: "update with email change error",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			planEmail:   "new@example.com",
			stateEmail:  "old@example.com",
			planRole:    types.StringValue("global:admin"),
			stateRole:   types.StringValue("global:member"),
			expectError: true,
		},
		{
			name: "update role with no change",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet && r.URL.Path == "/users/user-1" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{
						"id":    "user-1",
						"email": "test@example.com",
						"role":  "global:admin",
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			planEmail:   "test@example.com",
			stateEmail:  "test@example.com",
			planRole:    types.StringValue("global:admin"),
			stateRole:   types.StringValue("global:admin"),
			expectError: false,
		},
		{
			name: "update role with null plan role",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet && r.URL.Path == "/users/user-1" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{
						"id":    "user-1",
						"email": "test@example.com",
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			planEmail:   "test@example.com",
			stateEmail:  "test@example.com",
			planRole:    types.StringNull(),
			stateRole:   types.StringValue("global:admin"),
			expectError: false,
		},
		{
			name: "update role with null state role",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet && r.URL.Path == "/users/user-1" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{
						"id":    "user-1",
						"email": "test@example.com",
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			planEmail:   "test@example.com",
			stateEmail:  "test@example.com",
			planRole:    types.StringValue("global:admin"),
			stateRole:   types.StringNull(),
			expectError: false,
		},
		{
			name: "update with role API error",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPatch {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			planEmail:   "test@example.com",
			stateEmail:  "test@example.com",
			planRole:    types.StringValue("global:admin"),
			stateRole:   types.StringValue("global:member"),
			expectError: true,
		},
		{
			name: "update with refresh error",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPatch && r.URL.Path == "/users/user-1/role" {
					w.WriteHeader(http.StatusOK)
					return
				}
				if r.Method == http.MethodGet {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			planEmail:   "test@example.com",
			stateEmail:  "test@example.com",
			planRole:    types.StringValue("global:admin"),
			stateRole:   types.StringValue("global:member"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := http.HandlerFunc(tt.setupHandler)
			n8nClient, server := setupTestClient(t, handler)
			defer server.Close()

			r := &UserResource{client: n8nClient}
			plan := &models.Resource{
				ID:    types.StringValue("user-1"),
				Email: types.StringValue(tt.planEmail),
				Role:  tt.planRole,
			}
			state := &models.Resource{
				ID:    types.StringValue("user-1"),
				Email: types.StringValue(tt.stateEmail),
				Role:  tt.stateRole,
			}

			req := resource.UpdateRequest{
				Plan:  tfsdk.Plan{Schema: createTestSchema(t)},
				State: tfsdk.State{Schema: createTestSchema(t)},
			}
			resp := resource.UpdateResponse{
				State: tfsdk.State{Schema: createTestSchema(t)},
			}

			req.Plan.Set(context.Background(), plan)
			req.State.Set(context.Background(), state)

			r.Update(context.Background(), req, &resp)

			if tt.expectError {
				assert.True(t, resp.Diagnostics.HasError())
			} else {
				assert.False(t, resp.Diagnostics.HasError())
			}
		})
	}
}
