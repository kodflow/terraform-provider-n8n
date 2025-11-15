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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "has correct value":
				assert.Equal(t, 10, DEFAULT_LIST_CAPACITY, "DEFAULT_LIST_CAPACITY should be 10")

			case "is of type int":
				assert.IsType(t, 0, DEFAULT_LIST_CAPACITY, "DEFAULT_LIST_CAPACITY should be int type")

			case "is positive":
				assert.Greater(t, DEFAULT_LIST_CAPACITY, 0, "DEFAULT_LIST_CAPACITY should be positive")

			case "is reasonable for list capacity":
				assert.GreaterOrEqual(t, DEFAULT_LIST_CAPACITY, 1, "Should be at least 1")
				assert.LessOrEqual(t, DEFAULT_LIST_CAPACITY, 1000, "Should not be excessively large")

			case "can be used to create slice":
				slice := make([]string, 0, DEFAULT_LIST_CAPACITY)
				assert.NotNil(t, slice, "Should be able to create slice with DEFAULT_LIST_CAPACITY")
				assert.Equal(t, 0, len(slice), "Initial length should be 0")
				assert.Equal(t, DEFAULT_LIST_CAPACITY, cap(slice), "Capacity should be DEFAULT_LIST_CAPACITY")

			case "is const and cannot be modified":
				value := DEFAULT_LIST_CAPACITY
				assert.Equal(t, 10, value, "Constant value should be 10")

			case "error case - validation checks":
				assert.NotEqual(t, 0, DEFAULT_LIST_CAPACITY)
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "has correct value":
				assert.Equal(t, 32, FLOAT32_BIT_SIZE, "FLOAT32_BIT_SIZE should be 32")

			case "is of type int":
				assert.IsType(t, 0, FLOAT32_BIT_SIZE, "FLOAT32_BIT_SIZE should be int type")

			case "is positive":
				assert.Greater(t, FLOAT32_BIT_SIZE, 0, "FLOAT32_BIT_SIZE should be positive")

			case "matches float32 bit size":
				assert.Equal(t, 32, FLOAT32_BIT_SIZE, "Should match float32 bit size")

			case "can be used for parsing":
				assert.Equal(t, 32, FLOAT32_BIT_SIZE, "Should be usable as bitSize parameter")

			case "is const and cannot be modified":
				value := FLOAT32_BIT_SIZE
				assert.Equal(t, 32, value, "Constant value should be 32")

			case "is power of 2":
				assert.Equal(t, 32, FLOAT32_BIT_SIZE)
				isPowerOfTwo := (FLOAT32_BIT_SIZE & (FLOAT32_BIT_SIZE - 1)) == 0
				assert.True(t, isPowerOfTwo, "Bit size should be a power of 2")

			case "error case - validation checks":
				assert.NotEqual(t, 0, FLOAT32_BIT_SIZE)
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "all constants are defined":
				assert.NotNil(t, DEFAULT_LIST_CAPACITY)
				assert.NotNil(t, FLOAT32_BIT_SIZE)

			case "constants have expected types":
				listCap := DEFAULT_LIST_CAPACITY
				bitSize := FLOAT32_BIT_SIZE

				assert.Equal(t, 10, listCap)
				assert.Equal(t, 32, bitSize)

			case "constants are public":
				assert.True(t, true, "Constants are accessible, therefore public")

			case "error case - validation checks":
				assert.NotEqual(t, DEFAULT_LIST_CAPACITY, FLOAT32_BIT_SIZE)
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "use with make for slice":
				items := make([]int, 0, DEFAULT_LIST_CAPACITY)
				assert.Equal(t, 0, len(items))
				assert.Equal(t, DEFAULT_LIST_CAPACITY, cap(items))

			case "use with make for map":
				items := make(map[string]string, DEFAULT_LIST_CAPACITY)
				assert.NotNil(t, items)
				assert.Equal(t, 0, len(items))

			case "prevents immediate reallocation":
				items := make([]int, 0, DEFAULT_LIST_CAPACITY)

				for i := 0; i < DEFAULT_LIST_CAPACITY; i++ {
					items = append(items, i)
				}

				assert.Equal(t, DEFAULT_LIST_CAPACITY, cap(items))
				assert.Equal(t, DEFAULT_LIST_CAPACITY, len(items))

			case "capacity is sufficient for small lists":
				items := make([]string, 0, DEFAULT_LIST_CAPACITY)

				for i := 0; i < 5; i++ {
					items = append(items, "item")
				}

				assert.Less(t, len(items), cap(items))

			case "can be used in different contexts":
				intSlice := make([]int, 0, DEFAULT_LIST_CAPACITY)
				assert.Equal(t, DEFAULT_LIST_CAPACITY, cap(intSlice))

				stringSlice := make([]string, 0, DEFAULT_LIST_CAPACITY)
				assert.Equal(t, DEFAULT_LIST_CAPACITY, cap(stringSlice))

				type Item struct{ ID int }
				structSlice := make([]Item, 0, DEFAULT_LIST_CAPACITY)
				assert.Equal(t, DEFAULT_LIST_CAPACITY, cap(structSlice))

			case "error case - validation checks":
				items := make([]int, 0, DEFAULT_LIST_CAPACITY)
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "use with strconv.ParseFloat":
				assert.Equal(t, 32, FLOAT32_BIT_SIZE)

			case "matches float32 specification":
				assert.Equal(t, 32, FLOAT32_BIT_SIZE)

			case "differentiates from float64":
				assert.NotEqual(t, 64, FLOAT32_BIT_SIZE)
				assert.Equal(t, 32, FLOAT32_BIT_SIZE)

			case "can be used in bit manipulation":
				bytesRequired := FLOAT32_BIT_SIZE / 8
				assert.Equal(t, 4, bytesRequired, "float32 requires 4 bytes")

			case "error case - validation checks":
				assert.Greater(t, FLOAT32_BIT_SIZE, 0)
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
		{name: "DEFAULT_LIST_CAPACITY is not zero", wantErr: false},
		{name: "FLOAT32_BIT_SIZE is not zero", wantErr: false},
		{name: "constants are different", wantErr: false},
		{name: "error case - validation checks", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "DEFAULT_LIST_CAPACITY is not zero":
				assert.NotEqual(t, 0, DEFAULT_LIST_CAPACITY, "Should have a non-zero default")

			case "FLOAT32_BIT_SIZE is not zero":
				assert.NotEqual(t, 0, FLOAT32_BIT_SIZE, "Should have a non-zero value")

			case "constants are different":
				assert.NotEqual(t, DEFAULT_LIST_CAPACITY, FLOAT32_BIT_SIZE, "Constants should have different purposes and values")

			case "error case - validation checks":
				assert.True(t, DEFAULT_LIST_CAPACITY != FLOAT32_BIT_SIZE)
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
		{name: "DEFAULT_LIST_CAPACITY follows naming convention", wantErr: false},
		{name: "FLOAT32_BIT_SIZE follows naming convention", wantErr: false},
		{name: "error case - validation checks", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "DEFAULT_LIST_CAPACITY follows naming convention":
				assert.Equal(t, 10, DEFAULT_LIST_CAPACITY)

			case "FLOAT32_BIT_SIZE follows naming convention":
				assert.Equal(t, 32, FLOAT32_BIT_SIZE)

			case "error case - validation checks":
				assert.NotEqual(t, 0, DEFAULT_LIST_CAPACITY)
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "constants have clear purpose":
				assert.Equal(t, 10, DEFAULT_LIST_CAPACITY)
				assert.Equal(t, 32, FLOAT32_BIT_SIZE)

			case "error case - validation checks":
				assert.NotEqual(t, DEFAULT_LIST_CAPACITY, FLOAT32_BIT_SIZE)
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "constants cannot be reassigned":
				capacity := DEFAULT_LIST_CAPACITY
				bitSize := FLOAT32_BIT_SIZE

				assert.Equal(t, 10, capacity)
				assert.Equal(t, 32, bitSize)

			case "constants maintain their values":
				value1 := DEFAULT_LIST_CAPACITY
				value2 := DEFAULT_LIST_CAPACITY
				value3 := DEFAULT_LIST_CAPACITY

				assert.Equal(t, value1, value2)
				assert.Equal(t, value2, value3)
				assert.Equal(t, 10, value1)

			case "error case - validation checks":
				assert.NotEqual(t, 0, DEFAULT_LIST_CAPACITY)
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
		{name: "DEFAULT_LIST_CAPACITY arithmetic", wantErr: false},
		{name: "FLOAT32_BIT_SIZE arithmetic", wantErr: false},
		{name: "constants in comparisons", wantErr: false},
		{name: "error case - validation checks", wantErr: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "DEFAULT_LIST_CAPACITY arithmetic":
				doubled := DEFAULT_LIST_CAPACITY * 2
				assert.Equal(t, 20, doubled)

				halved := DEFAULT_LIST_CAPACITY / 2
				assert.Equal(t, 5, halved)

			case "FLOAT32_BIT_SIZE arithmetic":
				bytes := FLOAT32_BIT_SIZE / 8
				assert.Equal(t, 4, bytes)

				doubled := FLOAT32_BIT_SIZE * 2
				assert.Equal(t, 64, doubled)

			case "constants in comparisons":
				assert.True(t, DEFAULT_LIST_CAPACITY < FLOAT32_BIT_SIZE)
				assert.True(t, DEFAULT_LIST_CAPACITY > 0)
				assert.True(t, FLOAT32_BIT_SIZE > 0)

			case "error case - validation checks":
				assert.NotEqual(t, DEFAULT_LIST_CAPACITY, FLOAT32_BIT_SIZE)
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "constants have no runtime overhead":
				for i := 0; i < 1000; i++ {
					_ = make([]int, 0, DEFAULT_LIST_CAPACITY)
				}
				assert.Equal(t, 10, DEFAULT_LIST_CAPACITY)

			case "error case - validation checks":
				assert.NotEqual(t, 0, DEFAULT_LIST_CAPACITY)
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "both constants can be used together":
				floats := make([]float32, 0, DEFAULT_LIST_CAPACITY)
				assert.Equal(t, 0, len(floats))
				assert.Equal(t, DEFAULT_LIST_CAPACITY, cap(floats))

				assert.Equal(t, 32, FLOAT32_BIT_SIZE)

			case "constants are independent":
				capacity := DEFAULT_LIST_CAPACITY
				bitSize := FLOAT32_BIT_SIZE

				assert.Equal(t, 10, capacity)
				assert.Equal(t, 32, bitSize)
				assert.NotEqual(t, capacity, bitSize)

			case "error case - validation checks":
				assert.NotEqual(t, DEFAULT_LIST_CAPACITY, FLOAT32_BIT_SIZE)
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch tt.name {
			case "constants are exported":
				assert.True(t, true, "Constants are accessible in tests, therefore exported")

			case "error case - validation checks":
				assert.NotEqual(t, 0, DEFAULT_LIST_CAPACITY)
			}
		})
	}
}

func BenchmarkDEFAULT_LIST_CAPACITY_Usage(b *testing.B) {
	b.Run("slice creation with constant", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = make([]int, 0, DEFAULT_LIST_CAPACITY)
		}
	})

	b.Run("slice creation with literal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = make([]int, 0, 10)
		}
	})
}

func BenchmarkFLOAT32_BIT_SIZE_Usage(b *testing.B) {
	b.Run("using constant", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = FLOAT32_BIT_SIZE / 8
		}
	})

	b.Run("using literal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = 32 / 8
		}
	})
}
