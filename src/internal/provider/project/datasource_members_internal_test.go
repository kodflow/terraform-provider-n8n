package project

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/project/models"
	"github.com/stretchr/testify/assert"
)

// TestProjectMembersDataSource_executeReadLogic tests the executeReadLogic method.
func TestProjectMembersDataSource_executeReadLogic(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		projectID    string
		setupHandler func(w http.ResponseWriter, r *http.Request)
		expectError  bool
		expectCount  int
	}{
		{
			name:      "successful read",
			projectID: "proj-123",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet && r.URL.Path == "/projects/proj-123/users" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					fn1, fn2 := "Alice", "Bob"
					ln1, ln2 := "Smith", "Jones"
					email1, email2 := "alice@example.com", "bob@example.com"
					role1, role2 := "project:admin", "project:viewer"
					id1, id2 := "user-1", "user-2"
					json.NewEncoder(w).Encode(map[string]any{
						"data": []any{
							map[string]any{
								"id":        id1,
								"email":     email1,
								"firstName": fn1,
								"lastName":  ln1,
								"role":      role1,
							},
							map[string]any{
								"id":        id2,
								"email":     email2,
								"firstName": fn2,
								"lastName":  ln2,
								"role":      role2,
							},
						},
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: false,
			expectCount: 2,
		},
		{
			name:      "empty members",
			projectID: "proj-empty",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet && r.URL.Path == "/projects/proj-empty/users" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{
						"data": []any{},
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: false,
			expectCount: 0,
		},
		{
			name:      "API error",
			projectID: "proj-500",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"message": "Internal server error"}`))
			},
			expectError: true,
		},
		{
			name:      "members with nil fields",
			projectID: "proj-nil",
			setupHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet && r.URL.Path == "/projects/proj-nil/users" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{
						"data": []any{
							map[string]any{},
						},
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: false,
			expectCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Use the existing setupTestClient helper from the project package tests.
			handler := http.HandlerFunc(tt.setupHandler)
			n8nClient, server := setupTestClient(t, handler)
			defer server.Close()

			d := &ProjectMembersDataSource{client: n8nClient}
			ctx := t.Context()
			data := &models.DataSourceMembers{
				ProjectID: types.StringValue(tt.projectID),
			}
			resp := &datasource.ReadResponse{
				State: datasource.ReadResponse{}.State,
			}

			result := d.executeReadLogic(ctx, data, resp)

			if tt.expectError {
				assert.False(t, result, "Should return false on error")
				assert.True(t, resp.Diagnostics.HasError(), "Should have diagnostics error")
			} else {
				assert.True(t, result, "Should return true on success")
				assert.False(t, resp.Diagnostics.HasError(), "Should not have diagnostics error")
				assert.Len(t, data.Members, tt.expectCount, "Members count should match")
			}
		})
	}
}

// TestMemberToItem tests the memberToItem helper function.
func TestMemberToItem(t *testing.T) {
	t.Parallel()

	userID := "user-123"
	role := "admin"

	tests := []struct {
		name            string
		member          n8nsdk.ProjectMember
		expectUserIDSet bool
		expectRoleSet   bool
	}{
		{
			name: "member with all fields",
			member: n8nsdk.ProjectMember{
				ID:        &userID,
				Role:      &role,
				Email:     new("test@example.com"),
				FirstName: new("Test"),
				LastName:  new("User"),
			},
			expectUserIDSet: true,
			expectRoleSet:   true,
		},
		{
			name:            "member with nil fields",
			member:          n8nsdk.ProjectMember{},
			expectUserIDSet: false,
			expectRoleSet:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			item := memberToItem(tt.member)
			if tt.expectUserIDSet {
				assert.Equal(t, userID, item.UserID.ValueString())
				assert.Equal(t, role, item.Role.ValueString())
			} else {
				assert.True(t, item.UserID.IsNull())
				assert.True(t, item.Role.IsNull())
			}
		})
	}
}

// TestStringPtrToTF tests the stringPtrToTF helper function.
func TestStringPtrToTF(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		input      *string
		expectNull bool
		expectVal  string
	}{
		{
			name:       "non-nil pointer returns string value",
			input:      new("test-value"),
			expectNull: false,
			expectVal:  "test-value",
		},
		{
			name:       "nil pointer returns null string",
			input:      nil,
			expectNull: true,
			expectVal:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := stringPtrToTF(tt.input)
			if tt.expectNull {
				assert.Equal(t, types.StringNull(), result)
			} else {
				assert.Equal(t, types.StringValue(tt.expectVal), result)
			}
		})
	}
}
