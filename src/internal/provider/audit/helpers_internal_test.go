package audit

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/audit/models"
	"github.com/stretchr/testify/assert"
)

// TestMarshalRiskReport tests the marshalRiskReport helper.
func TestMarshalRiskReport(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		report   map[string]any
		isNull   bool
		contains string
	}{
		{
			name:   "nil report returns null",
			report: nil,
			isNull: true,
		},
		{
			name:   "empty report returns null",
			report: map[string]any{},
			isNull: true,
		},
		{
			name:     "report with data returns JSON string",
			report:   map[string]any{"risk": "high", "count": 3},
			isNull:   false,
			contains: "risk",
		},
		{
			name:   "unmarshalable value returns null",
			report: map[string]any{"ch": make(chan int)},
			isNull: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := marshalRiskReport(tt.report)
			if tt.isNull {
				assert.Equal(t, types.StringNull(), result)
			} else {
				assert.False(t, result.IsNull())
				assert.Contains(t, result.ValueString(), tt.contains)
			}
		})
	}
}

// TestMapAuditToDataSource tests the mapAuditToDataSource helper.
func TestMapAuditToDataSource(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                     string
		audit                    *n8nsdk.Audit
		expectCredentialsNotNull bool
		expectDatabaseNotNull    bool
	}{
		{
			name: "maps populated audit correctly",
			audit: &n8nsdk.Audit{
				CredentialsRiskReport: map[string]any{"risk": "low"},
				DatabaseRiskReport:    map[string]any{"tables": 5},
			},
			expectCredentialsNotNull: true,
			expectDatabaseNotNull:    true,
		},
		{
			name:                     "empty audit produces null fields",
			audit:                    &n8nsdk.Audit{},
			expectCredentialsNotNull: false,
			expectDatabaseNotNull:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			data := &models.DataSource{}
			mapAuditToDataSource(tt.audit, data)
			if tt.expectCredentialsNotNull {
				assert.False(t, data.CredentialsRiskReport.IsNull())
			} else {
				assert.True(t, data.CredentialsRiskReport.IsNull())
			}
			if tt.expectDatabaseNotNull {
				assert.False(t, data.DatabaseRiskReport.IsNull())
			} else {
				assert.True(t, data.DatabaseRiskReport.IsNull())
			}
		})
	}
}
