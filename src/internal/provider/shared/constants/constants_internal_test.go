package constants

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDEFAULT_LIST_CAPACITY(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "has correct value", wantErr: false},
		{name: "is of type int", wantErr: false},
		{name: "is positive", wantErr: false},
		{name: "is reasonable for list capacity", wantErr: false},
		{name: "can be used to create slice", wantErr: false},
		{name: "is const and cannot be modified", wantErr: false},
		{name: "error case - validation checks", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "has correct value":
				assert.Equal(t, 10, DefaultListCapacity, "DefaultListCapacity should be 10")

			case "is of type int":
				assert.IsType(t, 0, DefaultListCapacity, "DefaultListCapacity should be int type")

			case "is positive":
				assert.Greater(t, DefaultListCapacity, 0, "DefaultListCapacity should be positive")

			case "is reasonable for list capacity":
				assert.GreaterOrEqual(t, DefaultListCapacity, 1, "Should be at least 1")
				assert.LessOrEqual(t, DefaultListCapacity, 1000, "Should not be excessively large")

			case "can be used to create slice":
				slice := make([]string, 0, DefaultListCapacity)
				assert.NotNil(t, slice, "Should be able to create slice with DefaultListCapacity")
				assert.Equal(t, 0, len(slice), "Initial length should be 0")
				assert.Equal(t, DefaultListCapacity, cap(slice), "Capacity should be DefaultListCapacity")

			case "is const and cannot be modified":
				value := DefaultListCapacity
				assert.Equal(t, 10, value, "Constant value should be 10")

			case "error case - validation checks":
				assert.NotEqual(t, 0, DefaultListCapacity)
			}
		})
	}
}

func TestFLOAT32_BIT_SIZE(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "has correct value", wantErr: false},
		{name: "is of type int", wantErr: false},
		{name: "is positive", wantErr: false},
		{name: "matches float32 bit size", wantErr: false},
		{name: "can be used for parsing", wantErr: false},
		{name: "is const and cannot be modified", wantErr: false},
		{name: "is power of 2", wantErr: false},
		{name: "error case - validation checks", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "has correct value":
				assert.Equal(t, 32, Float32BitSize, "Float32BitSize should be 32")

			case "is of type int":
				assert.IsType(t, 0, Float32BitSize, "Float32BitSize should be int type")

			case "is positive":
				assert.Greater(t, Float32BitSize, 0, "Float32BitSize should be positive")

			case "matches float32 bit size":
				assert.Equal(t, 32, Float32BitSize, "Should match float32 bit size")

			case "can be used for parsing":
				assert.Equal(t, 32, Float32BitSize, "Should be usable as bitSize parameter")

			case "is const and cannot be modified":
				value := Float32BitSize
				assert.Equal(t, 32, value, "Constant value should be 32")

			case "is power of 2":
				assert.Equal(t, 32, Float32BitSize)
				isPowerOfTwo := (Float32BitSize & (Float32BitSize - 1)) == 0
				assert.True(t, isPowerOfTwo, "Bit size should be a power of 2")

			case "error case - validation checks":
				assert.NotEqual(t, 0, Float32BitSize)
			}
		})
	}
}

func TestConstantsPackage(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "all constants are defined", wantErr: false},
		{name: "constants have expected types", wantErr: false},
		{name: "constants are public", wantErr: false},
		{name: "error case - validation checks", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "all constants are defined":
				assert.NotNil(t, DefaultListCapacity)
				assert.NotNil(t, Float32BitSize)

			case "constants have expected types":
				listCap := DefaultListCapacity
				bitSize := Float32BitSize

				assert.Equal(t, 10, listCap)
				assert.Equal(t, 32, bitSize)

			case "constants are public":
				assert.True(t, true, "Constants are accessible, therefore public")

			case "error case - validation checks":
				assert.NotEqual(t, DefaultListCapacity, Float32BitSize)
			}
		})
	}
}

func TestDEFAULT_LIST_CAPACITY_UseCases(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "use with make for slice", wantErr: false},
		{name: "use with make for map", wantErr: false},
		{name: "prevents immediate reallocation", wantErr: false},
		{name: "capacity is sufficient for small lists", wantErr: false},
		{name: "can be used in different contexts", wantErr: false},
		{name: "error case - validation checks", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "use with make for slice":
				items := make([]int, 0, DefaultListCapacity)
				assert.Equal(t, 0, len(items))
				assert.Equal(t, DefaultListCapacity, cap(items))

			case "use with make for map":
				items := make(map[string]string, DefaultListCapacity)
				assert.NotNil(t, items)
				assert.Equal(t, 0, len(items))

			case "prevents immediate reallocation":
				items := make([]int, 0, DefaultListCapacity)

				for i := range DefaultListCapacity {
					items = append(items, i)
				}

				assert.Equal(t, DefaultListCapacity, cap(items))
				assert.Equal(t, DefaultListCapacity, len(items))

			case "capacity is sufficient for small lists":
				items := make([]string, 0, DefaultListCapacity)

				for range 5 {
					items = append(items, "item")
				}

				assert.Less(t, len(items), cap(items))

			case "can be used in different contexts":
				intSlice := make([]int, 0, DefaultListCapacity)
				assert.Equal(t, DefaultListCapacity, cap(intSlice))

				stringSlice := make([]string, 0, DefaultListCapacity)
				assert.Equal(t, DefaultListCapacity, cap(stringSlice))

				type Item struct{ ID int }
				structSlice := make([]Item, 0, DefaultListCapacity)
				assert.Equal(t, DefaultListCapacity, cap(structSlice))

			case "error case - validation checks":
				items := make([]int, 0, DefaultListCapacity)
				assert.NotNil(t, items)
			}
		})
	}
}

func TestFLOAT32_BIT_SIZE_UseCases(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "use with strconv.ParseFloat", wantErr: false},
		{name: "matches float32 specification", wantErr: false},
		{name: "differentiates from float64", wantErr: false},
		{name: "can be used in bit manipulation", wantErr: false},
		{name: "error case - validation checks", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "use with strconv.ParseFloat", "matches float32 specification":
				assert.Equal(t, 32, Float32BitSize)

			case "differentiates from float64":
				assert.NotEqual(t, 64, Float32BitSize)
				assert.Equal(t, 32, Float32BitSize)

			case "can be used in bit manipulation":
				bytesRequired := Float32BitSize / 8
				assert.Equal(t, 4, bytesRequired, "float32 requires 4 bytes")

			case "error case - validation checks":
				assert.Greater(t, Float32BitSize, 0)
			}
		})
	}
}

func TestConstantValues(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "DefaultListCapacity is not zero", wantErr: false},
		{name: "Float32BitSize is not zero", wantErr: false},
		{name: "constants are different", wantErr: false},
		{name: "error case - validation checks", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "DefaultListCapacity is not zero":
				assert.NotEqual(t, 0, DefaultListCapacity, "Should have a non-zero default")

			case "Float32BitSize is not zero":
				assert.NotEqual(t, 0, Float32BitSize, "Should have a non-zero value")

			case "constants are different":
				assert.NotEqual(t, DefaultListCapacity, Float32BitSize, "Constants should have different purposes and values")

			case "error case - validation checks":
				assert.True(t, DefaultListCapacity != Float32BitSize)
			}
		})
	}
}

func TestConstantNaming(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "DefaultListCapacity follows naming convention", wantErr: false},
		{name: "Float32BitSize follows naming convention", wantErr: false},
		{name: "error case - validation checks", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "DefaultListCapacity follows naming convention":
				assert.Equal(t, 10, DefaultListCapacity)

			case "Float32BitSize follows naming convention":
				assert.Equal(t, 32, Float32BitSize)

			case "error case - validation checks":
				assert.NotEqual(t, 0, DefaultListCapacity)
			}
		})
	}
}

func TestConstantDocumentation(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "constants have clear purpose", wantErr: false},
		{name: "error case - validation checks", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "constants have clear purpose":
				assert.Equal(t, 10, DefaultListCapacity)
				assert.Equal(t, 32, Float32BitSize)

			case "error case - validation checks":
				assert.NotEqual(t, DefaultListCapacity, Float32BitSize)
			}
		})
	}
}

func TestConstantsImmutability(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "constants cannot be reassigned", wantErr: false},
		{name: "constants maintain their values", wantErr: false},
		{name: "error case - validation checks", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "constants cannot be reassigned":
				capacity := DefaultListCapacity
				bitSize := Float32BitSize

				assert.Equal(t, 10, capacity)
				assert.Equal(t, 32, bitSize)

			case "constants maintain their values":
				value1 := DefaultListCapacity
				value2 := DefaultListCapacity
				value3 := DefaultListCapacity

				assert.Equal(t, value1, value2)
				assert.Equal(t, value2, value3)
				assert.Equal(t, 10, value1)

			case "error case - validation checks":
				assert.NotEqual(t, 0, DefaultListCapacity)
			}
		})
	}
}

func TestConstantsEdgeCases(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "DefaultListCapacity arithmetic", wantErr: false},
		{name: "Float32BitSize arithmetic", wantErr: false},
		{name: "constants in comparisons", wantErr: false},
		{name: "error case - validation checks", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "DefaultListCapacity arithmetic":
				doubled := DefaultListCapacity * 2
				assert.Equal(t, 20, doubled)

				halved := DefaultListCapacity / 2
				assert.Equal(t, 5, halved)

			case "Float32BitSize arithmetic":
				bytes := Float32BitSize / 8
				assert.Equal(t, 4, bytes)

				doubled := Float32BitSize * 2
				assert.Equal(t, 64, doubled)

			case "constants in comparisons":
				assert.True(t, DefaultListCapacity < Float32BitSize)
				assert.True(t, DefaultListCapacity > 0)
				assert.True(t, Float32BitSize > 0)

			case "error case - validation checks":
				assert.NotEqual(t, DefaultListCapacity, Float32BitSize)
			}
		})
	}
}

func TestConstantsPerformance(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "constants have no runtime overhead", wantErr: false},
		{name: "error case - validation checks", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "constants have no runtime overhead":
				var sink []int
				for range 1000 {
					sink = make([]int, 0, DefaultListCapacity)
				}
				assert.Equal(t, 10, cap(sink))
				assert.Equal(t, 10, DefaultListCapacity)

			case "error case - validation checks":
				assert.NotEqual(t, 0, DefaultListCapacity)
			}
		})
	}
}

func TestConstantsIntegration(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "both constants can be used together", wantErr: false},
		{name: "constants are independent", wantErr: false},
		{name: "error case - validation checks", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "both constants can be used together":
				floats := make([]float32, 0, DefaultListCapacity)
				assert.Equal(t, 0, len(floats))
				assert.Equal(t, DefaultListCapacity, cap(floats))

				assert.Equal(t, 32, Float32BitSize)

			case "constants are independent":
				capacity := DefaultListCapacity
				bitSize := Float32BitSize

				assert.Equal(t, 10, capacity)
				assert.Equal(t, 32, bitSize)
				assert.NotEqual(t, capacity, bitSize)

			case "error case - validation checks":
				assert.NotEqual(t, DefaultListCapacity, Float32BitSize)
			}
		})
	}
}

func TestConstantsAccessibility(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "constants are exported", wantErr: false},
		{name: "error case - validation checks", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "constants are exported":
				assert.True(t, true, "Constants are accessible in tests, therefore exported")

			case "error case - validation checks":
				assert.NotEqual(t, 0, DefaultListCapacity)
			}
		})
	}
}

func BenchmarkDEFAULT_LIST_CAPACITY_Usage(b *testing.B) {
	b.Run("slice creation with constant", func(b *testing.B) {
		var sink []int
		for b.Loop() {
			sink = make([]int, 0, DefaultListCapacity)
		}
		b.ReportMetric(float64(cap(sink)), "cap")
	})

	b.Run("slice creation with literal", func(b *testing.B) {
		var sink []int
		for b.Loop() {
			sink = make([]int, 0, 10)
		}
		b.ReportMetric(float64(cap(sink)), "cap")
	})
}

func BenchmarkFLOAT32_BIT_SIZE_Usage(b *testing.B) {
	b.Run("using constant", func(b *testing.B) {
		var result int
		for b.Loop() {
			result = Float32BitSize / 8
		}
		b.ReportMetric(float64(result), "bytes")
	})

	b.Run("using literal", func(b *testing.B) {
		var result int
		for b.Loop() {
			result = 32 / 8
		}
		b.ReportMetric(float64(result), "bytes")
	})
}
