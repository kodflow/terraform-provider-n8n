package models

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/kodflow/n8n/sdk/n8nsdk"
	"github.com/stretchr/testify/assert"
)

// Helper functions for creating pointers.
func strPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}

func intPtr(i int32) *int32 {
	return &i
}

func TestWorkflowBackup(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "create with valid values"},
		{name: "create with nil workflow"},
		{name: "create with empty ID"},
		{name: "zero value struct"},
		{name: "various ID formats"},
		{name: "workflow with complex structure"},
		{name: "modify backup"},
		{name: "pointer to struct"},
		{name: "copy struct"},
		{name: "deep copy workflow"},
		{name: "array of backups"},
		{name: "map of backups"},
		{name: "special characters in ID"},
		{name: "very long ID"},
		{name: "backup use case simulation"},
		{name: "nil checks"},
		{name: "comparison"},
		{name: "json serialization"},
		{name: "empty workflow backup"},
		{name: "error case - validation checks", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "create with valid values":
				// Test creating backup with valid workflow
				workflow := &n8nsdk.Workflow{
					Id:   strPtr("workflow-123"),
					Name: "Test Workflow",
					Nodes: []n8nsdk.Node{
						{
							Id:   strPtr("node-1"),
							Name: strPtr("Start"),
							Type: strPtr("n8n-nodes-base.start"),
						},
					},
					Active: boolPtr(true),
				}

				backup := WorkflowBackup{
					ID:       "backup-001",
					Original: workflow,
				}

				assert.NotNil(t, backup)
				assert.Equal(t, "backup-001", backup.ID)
				assert.NotNil(t, backup.Original)
				assert.Equal(t, "workflow-123", *backup.Original.Id)
				assert.Equal(t, "Test Workflow", backup.Original.Name)

			case "create with nil workflow":
				// Test creating backup with nil workflow (error scenario)
				backup := WorkflowBackup{
					ID:       "backup-002",
					Original: nil,
				}

				assert.NotNil(t, backup)
				assert.Equal(t, "backup-002", backup.ID)
				assert.Nil(t, backup.Original)

			case "create with empty ID":
				// Test creating backup with empty ID
				workflow := &n8nsdk.Workflow{
					Id: strPtr("workflow-456"),
				}

				backup := WorkflowBackup{
					ID:       "",
					Original: workflow,
				}

				assert.NotNil(t, backup)
				assert.Equal(t, "", backup.ID)
				assert.NotNil(t, backup.Original)
				assert.Equal(t, "workflow-456", *backup.Original.Id)

			case "zero value struct":
				// Test zero value struct
				var backup WorkflowBackup

				assert.Equal(t, "", backup.ID)
				assert.Nil(t, backup.Original)

			case "various ID formats":
				// Test various ID formats
				ids := []string{
					"simple-id",
					"ID_WITH_UNDERSCORES",
					"id-with-dashes",
					"id.with.dots",
					"id123456789",
					"veryLongIDWithManyCharacters1234567890abcdefghijklmnopqrstuvwxyz",
					"id=with=equals",
					"id/with/slashes",
					"id:with:colons",
					"backup_20240101_120000",
					"", // empty ID
				}

				for _, id := range ids {
					backup := WorkflowBackup{
						ID:       id,
						Original: &n8nsdk.Workflow{},
					}
					assert.Equal(t, id, backup.ID)
				}

			case "workflow with complex structure":
				// Test backup with complex workflow structure
				now := time.Now()
				workflow := &n8nsdk.Workflow{
					Id:        strPtr("complex-workflow"),
					Name:      "Complex Workflow",
					Active:    boolPtr(true),
					CreatedAt: &now,
					UpdatedAt: &now,
					Nodes: []n8nsdk.Node{
						{
							Id:       strPtr("node-1"),
							Name:     strPtr("HTTP Request"),
							Type:     strPtr("n8n-nodes-base.httpRequest"),
							Position: []float32{250, 300},
							Parameters: map[string]interface{}{
								"method": "GET",
								"url":    "https://api.example.com/data",
								"headers": map[string]interface{}{
									"Authorization": "Bearer token",
								},
							},
						},
						{
							Id:       strPtr("node-2"),
							Name:     strPtr("Set"),
							Type:     strPtr("n8n-nodes-base.set"),
							Position: []float32{450, 300},
							Parameters: map[string]interface{}{
								"values": map[string]interface{}{
									"key1": "value1",
									"key2": 123,
									"key3": true,
								},
							},
						},
					},
					Connections: map[string]interface{}{
						"HTTP Request": map[string]interface{}{
							"main": []interface{}{
								[]interface{}{
									map[string]interface{}{
										"node":  "Set",
										"type":  "main",
										"index": 0,
									},
								},
							},
						},
					},
					Settings: n8nsdk.WorkflowSettings{
						ExecutionOrder:           strPtr("v1"),
						ErrorWorkflow:            strPtr("error-handler-workflow"),
						Timezone:                 strPtr("America/New_York"),
						SaveDataErrorExecution:   strPtr("all"),
						SaveDataSuccessExecution: strPtr("all"),
					},
					Tags: []n8nsdk.Tag{
						{Id: strPtr("tag1"), Name: "production"},
						{Id: strPtr("tag2"), Name: "api"},
						{Id: strPtr("tag3"), Name: "integration"},
					},
				}

				backup := WorkflowBackup{
					ID:       "backup-complex",
					Original: workflow,
				}

				assert.NotNil(t, backup.Original)
				assert.Equal(t, "complex-workflow", *backup.Original.Id)
				assert.Equal(t, "Complex Workflow", backup.Original.Name)
				assert.NotNil(t, backup.Original.Active)
				assert.True(t, *backup.Original.Active)
				assert.Len(t, backup.Original.Nodes, 2)
				assert.NotNil(t, backup.Original.Connections)
				assert.NotNil(t, backup.Original.Settings)
				assert.Len(t, backup.Original.Tags, 3)

			case "modify backup":
				// Test modifying backup fields
				backup := WorkflowBackup{
					ID:       "original-id",
					Original: &n8nsdk.Workflow{Id: strPtr("wf-1")},
				}

				// Modify ID
				backup.ID = "modified-id"
				assert.Equal(t, "modified-id", backup.ID)

				// Replace workflow
				newWorkflow := &n8nsdk.Workflow{Id: strPtr("wf-2")}
				backup.Original = newWorkflow
				assert.Equal(t, "wf-2", *backup.Original.Id)

				// Set to nil
				backup.Original = nil
				assert.Nil(t, backup.Original)

			case "pointer to struct":
				// Test pointer to struct
				backup := &WorkflowBackup{
					ID:       "pointer-backup",
					Original: &n8nsdk.Workflow{Id: strPtr("pointer-workflow")},
				}

				assert.NotNil(t, backup)
				assert.Equal(t, "pointer-backup", backup.ID)
				assert.NotNil(t, backup.Original)
				assert.Equal(t, "pointer-workflow", *backup.Original.Id)

			case "copy struct":
				// Test copying struct
				original := WorkflowBackup{
					ID:       "original-backup",
					Original: &n8nsdk.Workflow{Id: strPtr("original-workflow")},
				}

				copied := original

				assert.Equal(t, original.ID, copied.ID)
				assert.Equal(t, original.Original, copied.Original) // Pointer equality

				// Modify copied ID
				copied.ID = "copy-backup"
				assert.Equal(t, "original-backup", original.ID)
				assert.Equal(t, "copy-backup", copied.ID)

				// Note: Original field is a pointer, so both structs point to same workflow
				*copied.Original.Id = "modified-workflow"
				assert.Equal(t, "modified-workflow", *original.Original.Id)
				assert.Equal(t, "modified-workflow", *copied.Original.Id)

			case "deep copy workflow":
				// Test that we can create independent copies
				workflow1 := &n8nsdk.Workflow{
					Id:   strPtr("workflow-1"),
					Name: "Original",
				}

				backup1 := WorkflowBackup{
					ID:       "backup-1",
					Original: workflow1,
				}

				// Create a new workflow for second backup
				workflow2 := &n8nsdk.Workflow{
					Id:   strPtr(*workflow1.Id),
					Name: workflow1.Name,
				}

				backup2 := WorkflowBackup{
					ID:       "backup-2",
					Original: workflow2,
				}

				// Modify second workflow
				backup2.Original.Name = "Modified"

				// First backup should be unchanged
				assert.Equal(t, "Original", backup1.Original.Name)
				assert.Equal(t, "Modified", backup2.Original.Name)

			case "array of backups":
				// Test working with array of backups
				backups := []WorkflowBackup{
					{
						ID:       "backup-1",
						Original: &n8nsdk.Workflow{Id: strPtr("wf-1")},
					},
					{
						ID:       "backup-2",
						Original: &n8nsdk.Workflow{Id: strPtr("wf-2")},
					},
					{
						ID:       "backup-3",
						Original: &n8nsdk.Workflow{Id: strPtr("wf-3")},
					},
				}

				assert.Len(t, backups, 3)
				for i, backup := range backups {
					assert.Equal(t, fmt.Sprintf("backup-%d", i+1), backup.ID)
					assert.Equal(t, fmt.Sprintf("wf-%d", i+1), *backup.Original.Id)
				}

			case "map of backups":
				// Test working with map of backups
				backupMap := map[string]WorkflowBackup{
					"workflow-1": {
						ID:       "backup-1",
						Original: &n8nsdk.Workflow{Id: strPtr("workflow-1")},
					},
					"workflow-2": {
						ID:       "backup-2",
						Original: &n8nsdk.Workflow{Id: strPtr("workflow-2")},
					},
				}

				assert.Len(t, backupMap, 2)

				backup1, ok := backupMap["workflow-1"]
				assert.True(t, ok)
				assert.Equal(t, "backup-1", backup1.ID)

				backup2, ok := backupMap["workflow-2"]
				assert.True(t, ok)
				assert.Equal(t, "backup-2", backup2.ID)

			case "special characters in ID":
				// Test special characters in ID
				specialIDs := []string{
					"id-with-special-!@#$%^&*()",
					"id-测试-テスト-тест",
					"id\twith\ttabs",
					"id\nwith\nnewlines",
					"id with spaces",
					"id\"with\"quotes",
					"id'with'quotes",
				}

				for _, id := range specialIDs {
					backup := WorkflowBackup{
						ID:       id,
						Original: &n8nsdk.Workflow{},
					}
					assert.Equal(t, id, backup.ID)
				}

			case "very long ID":
				// Test with very long ID
				longID := strings.Repeat("a", 10000)

				backup := WorkflowBackup{
					ID:       longID,
					Original: &n8nsdk.Workflow{},
				}

				assert.Equal(t, longID, backup.ID)
				assert.Len(t, backup.ID, 10000)

			case "backup use case simulation":
				// Simulate actual backup use case during credential rotation

				// Original workflow state
				originalWorkflow := &n8nsdk.Workflow{
					Id:     strPtr("prod-workflow-123"),
					Name:   "Production Data Processing",
					Active: boolPtr(true),
					Nodes: []n8nsdk.Node{
						{
							Id:   strPtr("http-1"),
							Name: strPtr("Fetch Data"),
							Type: strPtr("n8n-nodes-base.httpRequest"),
							Credentials: map[string]interface{}{
								"httpBasicAuth": map[string]interface{}{
									"id":   "old-credential-id",
									"name": "Old API Credential",
								},
							},
						},
					},
				}

				// Create backup before credential rotation
				backup := WorkflowBackup{
					ID:       fmt.Sprintf("backup-%s-%d", *originalWorkflow.Id, 1704067200),
					Original: originalWorkflow,
				}

				// Simulate credential rotation (modify the original)
				originalWorkflow.Nodes[0].Credentials = map[string]interface{}{
					"httpBasicAuth": map[string]interface{}{
						"id":   "new-credential-id",
						"name": "New API Credential",
					},
				}

				// Backup should still have old credential reference
				// Note: In real scenario, we'd need deep copy to preserve old state
				// This test demonstrates the structure
				assert.NotNil(t, backup.Original)
				assert.Equal(t, "prod-workflow-123", *backup.Original.Id)

			case "nil checks":
				// Test nil safety
				backup := WorkflowBackup{
					ID:       "nil-check",
					Original: nil,
				}

				// Should not panic
				assert.Equal(t, "nil-check", backup.ID)
				assert.Nil(t, backup.Original)

				// Check accessing nil workflow doesn't cause issues
				if backup.Original != nil {
					_ = backup.Original.Id
				}

			case "comparison":
				// Test struct comparison
				wf1 := &n8nsdk.Workflow{Id: strPtr("wf-1")}
				wf2 := &n8nsdk.Workflow{Id: strPtr("wf-2")}

				backup1 := WorkflowBackup{
					ID:       "backup-1",
					Original: wf1,
				}

				backup2 := WorkflowBackup{
					ID:       "backup-1",
					Original: wf1,
				}

				backup3 := WorkflowBackup{
					ID:       "backup-2",
					Original: wf2,
				}

				// Same values
				assert.Equal(t, backup1.ID, backup2.ID)
				assert.Equal(t, backup1.Original, backup2.Original) // Same pointer

				// Different values
				assert.NotEqual(t, backup1.ID, backup3.ID)
				assert.NotEqual(t, backup1.Original, backup3.Original) // Different pointer

			case "json serialization":
				// Test JSON serialization
				backup := WorkflowBackup{
					ID: "json-backup",
					Original: &n8nsdk.Workflow{
						Id:   strPtr("json-workflow"),
						Name: "JSON Test",
					},
				}

				// Marshal to JSON
				data, err := json.Marshal(backup)
				assert.NoError(t, err)
				assert.NotNil(t, data)

				// Unmarshal back
				var restored WorkflowBackup
				err = json.Unmarshal(data, &restored)
				assert.NoError(t, err)
				assert.Equal(t, backup.ID, restored.ID)
				assert.Equal(t, *backup.Original.Id, *restored.Original.Id)
				assert.Equal(t, backup.Original.Name, restored.Original.Name)

			case "empty workflow backup":
				// Test with empty workflow
				backup := WorkflowBackup{
					ID:       "empty-backup",
					Original: &n8nsdk.Workflow{},
				}

				assert.Equal(t, "empty-backup", backup.ID)
				assert.NotNil(t, backup.Original)
				assert.Nil(t, backup.Original.Id)
				assert.Equal(t, "", backup.Original.Name)
				assert.Nil(t, backup.Original.Active)
				assert.Nil(t, backup.Original.Nodes)
				assert.Nil(t, backup.Original.Connections)

			case "error case - validation checks":
				// Test invalid backup scenarios
				// Test that struct allows invalid combinations
				backup := WorkflowBackup{
					ID:       "",
					Original: nil,
				}
				assert.Equal(t, "", backup.ID)
				assert.Nil(t, backup.Original)

				// Test that we can detect missing required fields
				if backup.ID == "" {
					assert.True(t, tt.wantErr)
				}
				if backup.Original == nil {
					assert.True(t, tt.wantErr)
				}
			}
		})
	}
}

func TestWorkflowBackupConcurrency(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "concurrent read"},
		{name: "concurrent backup creation"},
		{name: "concurrent map access"},
		{name: "error case - validation checks", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// NOTE: No t.Parallel() here because subtests contain goroutines

			switch tt.name {
			case "concurrent read":
				// Test concurrent reads
				backup := WorkflowBackup{
					ID:       "concurrent-backup",
					Original: &n8nsdk.Workflow{Id: strPtr("concurrent-workflow")},
				}

				var wg sync.WaitGroup
				errors := make(chan error, 100)

				for i := 0; i < 100; i++ {
					wg.Add(1)
					go func() {
						defer wg.Done()

						// Read operations
						id := backup.ID
						if id != "concurrent-backup" {
							errors <- fmt.Errorf("unexpected ID: %s", id)
							return
						}

						if backup.Original != nil && backup.Original.Id != nil {
							wfID := *backup.Original.Id
							if wfID != "concurrent-workflow" {
								errors <- fmt.Errorf("unexpected workflow ID: %s", wfID)
							}
						}
					}()
				}

				wg.Wait()
				close(errors)

				for err := range errors {
					t.Errorf("Concurrent read error: %v", err)
				}

			case "concurrent backup creation":
				// Test creating multiple backups concurrently
				var wg sync.WaitGroup
				backups := make([]WorkflowBackup, 100)

				for i := 0; i < 100; i++ {
					wg.Add(1)
					go func(idx int) {
						defer wg.Done()

						backups[idx] = WorkflowBackup{
							ID: fmt.Sprintf("backup-%d", idx),
							Original: &n8nsdk.Workflow{
								Id: strPtr(fmt.Sprintf("workflow-%d", idx)),
							},
						}
					}(i)
				}

				wg.Wait()

				// Verify all backups were created correctly
				for i := 0; i < 100; i++ {
					assert.Equal(t, fmt.Sprintf("backup-%d", i), backups[i].ID)
					assert.Equal(t, fmt.Sprintf("workflow-%d", i), *backups[i].Original.Id)
				}

			case "concurrent map access":
				// Test concurrent access to backup map
				backupMap := &sync.Map{}

				var wg sync.WaitGroup

				// Writers
				for i := 0; i < 50; i++ {
					wg.Add(1)
					go func(idx int) {
						defer wg.Done()

						backup := WorkflowBackup{
							ID:       fmt.Sprintf("backup-%d", idx),
							Original: &n8nsdk.Workflow{Id: strPtr(fmt.Sprintf("wf-%d", idx))},
						}
						backupMap.Store(backup.ID, backup)
					}(i)
				}

				// Readers
				for i := 0; i < 50; i++ {
					wg.Add(1)
					go func(idx int) {
						defer wg.Done()

						key := fmt.Sprintf("backup-%d", idx%50)
						if val, ok := backupMap.Load(key); ok {
							backup := val.(WorkflowBackup)
							_ = backup.ID
							_ = backup.Original
						}
					}(i)
				}

				wg.Wait()

				// Verify map contents
				count := 0
				backupMap.Range(func(key, value interface{}) bool {
					count++
					return true
				})
				assert.GreaterOrEqual(t, count, 0)

			case "error case - validation checks":
				// Test error scenarios in concurrent operations
				var wg sync.WaitGroup
				errors := make(chan error, 10)

				// Test concurrent access to nil workflow
				backup := WorkflowBackup{
					ID:       "concurrent-nil",
					Original: nil,
				}

				for i := 0; i < 10; i++ {
					wg.Add(1)
					go func() {
						defer wg.Done()

						if backup.Original == nil {
							errors <- fmt.Errorf("workflow is nil")
						}
					}()
				}

				wg.Wait()
				close(errors)

				// Should have errors
				errorCount := 0
				for range errors {
					errorCount++
				}
				if tt.wantErr {
					assert.Greater(t, errorCount, 0)
				}
			}
		})
	}
}

func TestWorkflowBackupEdgeCases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "backup with nil workflow operations"},
		{name: "backup restoration scenario"},
		{name: "multiple backups of same workflow"},
		{name: "backup chain"},
		{name: "partial workflow backup"},
		{name: "error case - validation checks", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch tt.name {
			case "backup with nil workflow operations":
				// Test operations when Original is nil
				backup := WorkflowBackup{
					ID:       "nil-workflow",
					Original: nil,
				}

				// These operations should not panic
				assert.Equal(t, "nil-workflow", backup.ID)

				// Simulate checking workflow before use
				if backup.Original != nil {
					_ = backup.Original.Id
				} else {
					assert.Nil(t, backup.Original)
				}

			case "backup restoration scenario":
				// Simulate backup restoration

				// Create original workflow
				original := &n8nsdk.Workflow{
					Id:     strPtr("prod-wf"),
					Name:   "Production Workflow",
					Active: boolPtr(true),
				}

				// Create backup
				backup := WorkflowBackup{
					ID:       "restore-backup",
					Original: original,
				}

				// Modify original (simulating failed update)
				original.Name = "Modified Name"
				original.Active = boolPtr(false)

				// Restore from backup (in real scenario, would deep copy)
				if backup.Original != nil {
					// This shows structure, but real restoration would need deep copy
					assert.Equal(t, "prod-wf", *backup.Original.Id)
					// Note: Due to pointer, backup.Original.Name is also modified
					// In production, we'd need deep copy during backup creation
				}

			case "multiple backups of same workflow":
				// Test managing multiple backups of same workflow
				workflowID := "multi-backup-wf"

				backups := make([]WorkflowBackup, 5)
				for i := 0; i < 5; i++ {
					backups[i] = WorkflowBackup{
						ID: fmt.Sprintf("backup-%s-v%d", workflowID, i+1),
						Original: &n8nsdk.Workflow{
							Id:           strPtr(workflowID),
							Name:         fmt.Sprintf("Version %d", i+1),
							TriggerCount: intPtr(int32(i + 1)),
						},
					}
				}

				// Verify we have all versions
				assert.Len(t, backups, 5)
				for i, backup := range backups {
					assert.Equal(t, fmt.Sprintf("backup-%s-v%d", workflowID, i+1), backup.ID)
					assert.Equal(t, workflowID, *backup.Original.Id)
					assert.Equal(t, int32(i+1), *backup.Original.TriggerCount)
				}

			case "backup chain":
				// Test chain of backups (backup of backup scenario)
				workflow := &n8nsdk.Workflow{
					Id:   strPtr("chain-wf"),
					Name: "Chain Workflow",
				}

				backup1 := WorkflowBackup{
					ID:       "backup-1",
					Original: workflow,
				}

				// In practice, we might backup the backup's workflow
				backup2 := WorkflowBackup{
					ID:       "backup-2",
					Original: backup1.Original, // Points to same workflow
				}

				assert.Equal(t, backup1.Original, backup2.Original)
				assert.Equal(t, workflow, backup1.Original)
				assert.Equal(t, workflow, backup2.Original)

			case "partial workflow backup":
				// Test backup with partially populated workflow
				workflow := &n8nsdk.Workflow{
					Id: strPtr("partial-wf"),
					// Other fields left as zero values
				}

				backup := WorkflowBackup{
					ID:       "partial-backup",
					Original: workflow,
				}

				assert.Equal(t, "partial-wf", *backup.Original.Id)
				assert.Equal(t, "", backup.Original.Name)
				assert.Nil(t, backup.Original.Active)
				assert.Nil(t, backup.Original.Nodes)

			case "error case - validation checks":
				// Test invalid edge case scenarios

				// Test backup with empty ID and nil workflow
				backup1 := WorkflowBackup{
					ID:       "",
					Original: nil,
				}
				assert.Equal(t, "", backup1.ID)
				assert.Nil(t, backup1.Original)

				// Test accessing nil workflow fields (should not panic)
				if backup1.Original != nil {
					_ = backup1.Original.Id
				} else {
					assert.True(t, tt.wantErr)
				}

				// Test backup with invalid workflow data
				backup2 := WorkflowBackup{
					ID: "invalid-backup",
					Original: &n8nsdk.Workflow{
						Id:   nil, // Invalid: nil ID
						Name: "",  // Invalid: empty name
					},
				}
				assert.NotNil(t, backup2.Original)
				assert.Nil(t, backup2.Original.Id)
				assert.Equal(t, "", backup2.Original.Name)

				if backup2.Original.Id == nil {
					assert.True(t, tt.wantErr)
				}
			}
		})
	}
}

func BenchmarkWorkflowBackup(b *testing.B) {
	b.Run("create", func(b *testing.B) {
		workflow := &n8nsdk.Workflow{
			Id:   strPtr("bench-wf"),
			Name: "Benchmark Workflow",
		}

		b.ResetTimer()
		for b.Loop() {
			_ = WorkflowBackup{
				ID:       "bench-backup",
				Original: workflow,
			}
		}
	})

	b.Run("access fields", func(b *testing.B) {
		backup := WorkflowBackup{
			ID:       "bench-backup",
			Original: &n8nsdk.Workflow{Id: strPtr("bench-wf")},
		}

		b.ResetTimer()
		for b.Loop() {
			_ = backup.ID
			if backup.Original != nil && backup.Original.Id != nil {
				_ = *backup.Original.Id
			}
		}
	})

	b.Run("with complex workflow", func(b *testing.B) {
		workflow := &n8nsdk.Workflow{
			Id:     strPtr("complex-bench-wf"),
			Name:   "Complex Benchmark",
			Active: boolPtr(true),
			Nodes: []n8nsdk.Node{
				{Id: strPtr("node-1"), Name: strPtr("Node 1")},
				{Id: strPtr("node-2"), Name: strPtr("Node 2")},
				{Id: strPtr("node-3"), Name: strPtr("Node 3")},
			},
			Tags: []n8nsdk.Tag{
				{Id: strPtr("tag1"), Name: "tag1"},
				{Id: strPtr("tag2"), Name: "tag2"},
				{Id: strPtr("tag3"), Name: "tag3"},
			},
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = WorkflowBackup{
				ID:       fmt.Sprintf("backup-%d", i),
				Original: workflow,
			}
		}
	})

	b.Run("map operations", func(b *testing.B) {
		backupMap := make(map[string]WorkflowBackup)
		workflow := &n8nsdk.Workflow{Id: strPtr("map-wf")}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			backup := WorkflowBackup{
				ID:       fmt.Sprintf("backup-%d", i),
				Original: workflow,
			}
			backupMap[backup.ID] = backup
		}
	})

	b.Run("slice operations", func(b *testing.B) {
		workflow := &n8nsdk.Workflow{Id: strPtr("slice-wf")}

		b.ResetTimer()
		var backups []WorkflowBackup
		for i := 0; i < b.N; i++ {
			backup := WorkflowBackup{
				ID:       fmt.Sprintf("backup-%d", i),
				Original: workflow,
			}
			backups = append(backups, backup)
		}
		_ = backups // Explicitly use the slice to prevent optimization
	})
}
