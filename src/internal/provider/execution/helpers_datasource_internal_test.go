package execution

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kodflow/terraform-provider-n8n/sdk/n8nsdk"
	"github.com/kodflow/terraform-provider-n8n/src/internal/provider/execution/models"
	"github.com/stretchr/testify/assert"
)

// TestFormatExecutionID tests the formatExecutionID helper.
func TestFormatExecutionID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		id       *float32
		expected string
		isNull   bool
	}{
		{
			name:   "nil ID returns null",
			id:     nil,
			isNull: true,
		},
		{
			name:     "integer ID",
			id:       float32Ptr(42),
			expected: "42",
			isNull:   false,
		},
		{
			name:     "zero ID",
			id:       float32Ptr(0),
			expected: "0",
			isNull:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := formatExecutionID(tt.id)
			if tt.isNull {
				assert.Equal(t, types.StringNull(), result)
			} else {
				assert.Equal(t, tt.expected, result.ValueString())
			}
		})
	}
}

// TestFormatNullableTime tests the formatNullableTime helper.
func TestFormatNullableTime(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		setup  func() n8nsdk.NullableTime
		isNull bool
	}{
		{
			name: "not set returns null",
			setup: func() n8nsdk.NullableTime {
				return n8nsdk.NullableTime{}
			},
			isNull: true,
		},
		{
			name: "set time returns formatted string",
			setup: func() n8nsdk.NullableTime {
				nt := n8nsdk.NullableTime{}
				ts := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
				nt.Set(&ts)
				return nt
			},
			isNull: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := formatNullableTime(tt.setup())
			if tt.isNull {
				assert.Equal(t, types.StringNull(), result)
			} else {
				assert.False(t, result.IsNull())
				assert.NotEmpty(t, result.ValueString())
			}
		})
	}
}

// TestFormatNullableBool tests the formatNullableBool helper.
func TestFormatNullableBool(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    *bool
		isNull   bool
		expected bool
	}{
		{
			name:   "nil returns null",
			input:  nil,
			isNull: true,
		},
		{
			name:     "true value",
			input:    boolPtr(true),
			isNull:   false,
			expected: true,
		},
		{
			name:     "false value",
			input:    boolPtr(false),
			isNull:   false,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := formatNullableBool(tt.input)
			if tt.isNull {
				assert.Equal(t, types.BoolNull(), result)
			} else {
				assert.Equal(t, tt.expected, result.ValueBool())
			}
		})
	}
}

// TestFormatNullableString tests the formatNullableString helper.
func TestFormatNullableString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    *string
		isNull   bool
		expected string
	}{
		{
			name:   "nil returns null",
			input:  nil,
			isNull: true,
		},
		{
			name:     "string value",
			input:    strPtr("success"),
			isNull:   false,
			expected: "success",
		},
		{
			name:     "empty string",
			input:    strPtr(""),
			isNull:   false,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := formatNullableString(tt.input)
			if tt.isNull {
				assert.Equal(t, types.StringNull(), result)
			} else {
				assert.Equal(t, tt.expected, result.ValueString())
			}
		})
	}
}

// TestMapExecutionToDataSource tests the mapExecutionToDataSource helper.
func TestMapExecutionToDataSource(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		exec           *n8nsdk.Execution
		expectedID     string
		expectedStatus string
		expectedMode   string
	}{
		{
			name: "maps all fields correctly",
			exec: &n8nsdk.Execution{
				Id:     float32Ptr(123),
				Mode:   strPtr("manual"),
				Status: strPtr("success"),
			},
			expectedID:     "123",
			expectedStatus: "success",
			expectedMode:   "manual",
		},
		{
			name:           "nil fields produce null values",
			exec:           &n8nsdk.Execution{},
			expectedID:     "",
			expectedStatus: "",
			expectedMode:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			data := &models.DataSource{}
			mapExecutionToDataSource(tt.exec, data)
			assert.Equal(t, tt.expectedID, data.ID.ValueString())
			assert.Equal(t, tt.expectedStatus, data.Status.ValueString())
			assert.Equal(t, tt.expectedMode, data.Mode.ValueString())
		})
	}
}

// TestMapExecutionToItem tests the mapExecutionToItem helper.
func TestMapExecutionToItem(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		exec           *n8nsdk.Execution
		expectedID     string
		expectedStatus string
	}{
		{
			name: "maps all fields correctly",
			exec: &n8nsdk.Execution{
				Id:     float32Ptr(456),
				Status: strPtr("error"),
			},
			expectedID:     "456",
			expectedStatus: "error",
		},
		{
			name:           "nil fields",
			exec:           &n8nsdk.Execution{},
			expectedID:     "",
			expectedStatus: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			item := mapExecutionToItem(tt.exec)
			assert.Equal(t, tt.expectedID, item.ID.ValueString())
			assert.Equal(t, tt.expectedStatus, item.Status.ValueString())
		})
	}
}

// TestIsValueSet tests the isValueSet helper.
func TestIsValueSet(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		v        types.String
		expected bool
	}{
		{
			name:     "null returns false",
			v:        types.StringNull(),
			expected: false,
		},
		{
			name:     "unknown returns false",
			v:        types.StringUnknown(),
			expected: false,
		},
		{
			name:     "set value returns true",
			v:        types.StringValue("test"),
			expected: true,
		},
		{
			name:     "empty string returns true",
			v:        types.StringValue(""),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expected, isValueSet(tt.v))
		})
	}
}

// TestBuildExecutionItems tests the buildExecutionItems helper.
func TestBuildExecutionItems(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		execList      *n8nsdk.ExecutionList
		expectedCount int
	}{
		{
			name:          "nil execList returns empty slice",
			execList:      nil,
			expectedCount: 0,
		},
		{
			name:          "empty data returns empty slice",
			execList:      &n8nsdk.ExecutionList{Data: nil},
			expectedCount: 0,
		},
		{
			name: "populated list returns mapped items",
			execList: &n8nsdk.ExecutionList{
				Data: []n8nsdk.Execution{
					{Id: float32Ptr(1), Status: strPtr("success")},
					{Id: float32Ptr(2), Status: strPtr("error")},
				},
			},
			expectedCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			items := buildExecutionItems(tt.execList)
			assert.Len(t, items, tt.expectedCount)
		})
	}
}

// float32Ptr returns a pointer to a float32 value.
func float32Ptr(v float32) *float32 { return new(v) }

// boolPtr returns a pointer to a bool value.
func boolPtr(v bool) *bool { return new(v) }

// strPtr returns a pointer to a string value.
func strPtr(v string) *string { return new(v) }
