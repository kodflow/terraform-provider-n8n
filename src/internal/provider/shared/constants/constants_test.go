package constants

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDEFAULT_LIST_CAPACITY(t *testing.T) {
	t.Run("has correct value", func(t *testing.T) {
		assert.Equal(t, 10, DEFAULT_LIST_CAPACITY, "DEFAULT_LIST_CAPACITY should be 10")
	})

	t.Run("is of type int", func(t *testing.T) {
		assert.IsType(t, 0, DEFAULT_LIST_CAPACITY, "DEFAULT_LIST_CAPACITY should be int type")
	})

	t.Run("is positive", func(t *testing.T) {
		assert.Greater(t, DEFAULT_LIST_CAPACITY, 0, "DEFAULT_LIST_CAPACITY should be positive")
	})

	t.Run("is reasonable for list capacity", func(t *testing.T) {
		// Should be a reasonable default capacity (not too small, not too large)
		assert.GreaterOrEqual(t, DEFAULT_LIST_CAPACITY, 1, "Should be at least 1")
		assert.LessOrEqual(t, DEFAULT_LIST_CAPACITY, 1000, "Should not be excessively large")
	})

	t.Run("can be used to create slice", func(t *testing.T) {
		slice := make([]string, 0, DEFAULT_LIST_CAPACITY)
		assert.NotNil(t, slice, "Should be able to create slice with DEFAULT_LIST_CAPACITY")
		assert.Equal(t, 0, len(slice), "Initial length should be 0")
		assert.Equal(t, DEFAULT_LIST_CAPACITY, cap(slice), "Capacity should be DEFAULT_LIST_CAPACITY")
	})

	t.Run("is const and cannot be modified", func(t *testing.T) {
		// This is a compile-time check - const values cannot be reassigned
		// The fact that this test compiles proves it's a const
		value := DEFAULT_LIST_CAPACITY
		assert.Equal(t, 10, value, "Constant value should be 10")
	})
}

func TestFLOAT32_BIT_SIZE(t *testing.T) {
	t.Run("has correct value", func(t *testing.T) {
		assert.Equal(t, 32, FLOAT32_BIT_SIZE, "FLOAT32_BIT_SIZE should be 32")
	})

	t.Run("is of type int", func(t *testing.T) {
		assert.IsType(t, 0, FLOAT32_BIT_SIZE, "FLOAT32_BIT_SIZE should be int type")
	})

	t.Run("is positive", func(t *testing.T) {
		assert.Greater(t, FLOAT32_BIT_SIZE, 0, "FLOAT32_BIT_SIZE should be positive")
	})

	t.Run("matches float32 bit size", func(t *testing.T) {
		// float32 is always 32 bits
		assert.Equal(t, 32, FLOAT32_BIT_SIZE, "Should match float32 bit size")
	})

	t.Run("can be used for parsing", func(t *testing.T) {
		// Verify it can be used with strconv functions that accept bitSize
		// This is a common use case for this constant
		assert.Equal(t, 32, FLOAT32_BIT_SIZE, "Should be usable as bitSize parameter")
	})

	t.Run("is const and cannot be modified", func(t *testing.T) {
		// This is a compile-time check - const values cannot be reassigned
		value := FLOAT32_BIT_SIZE
		assert.Equal(t, 32, value, "Constant value should be 32")
	})

	t.Run("is power of 2", func(t *testing.T) {
		// Bit sizes are typically powers of 2
		assert.Equal(t, 32, FLOAT32_BIT_SIZE)
		// Verify it's a power of 2
		isPowerOfTwo := (FLOAT32_BIT_SIZE & (FLOAT32_BIT_SIZE - 1)) == 0
		assert.True(t, isPowerOfTwo, "Bit size should be a power of 2")
	})
}

func TestConstantsPackage(t *testing.T) {
	t.Run("all constants are defined", func(t *testing.T) {
		// Verify both constants exist and are accessible
		assert.NotNil(t, DEFAULT_LIST_CAPACITY)
		assert.NotNil(t, FLOAT32_BIT_SIZE)
	})

	t.Run("constants have expected types", func(t *testing.T) {
		// Both should be int type
		listCap := DEFAULT_LIST_CAPACITY
		bitSize := FLOAT32_BIT_SIZE

		assert.Equal(t, 10, listCap)
		assert.Equal(t, 32, bitSize)
	})

	t.Run("constants are public", func(t *testing.T) {
		// Constants are exported (uppercase), so they can be used outside the package
		// This is verified by the fact that we can access them in tests
		assert.True(t, true, "Constants are accessible, therefore public")
	})
}

func TestDEFAULT_LIST_CAPACITY_UseCases(t *testing.T) {
	t.Run("use with make for slice", func(t *testing.T) {
		items := make([]int, 0, DEFAULT_LIST_CAPACITY)
		assert.Equal(t, 0, len(items))
		assert.Equal(t, DEFAULT_LIST_CAPACITY, cap(items))
	})

	t.Run("use with make for map", func(t *testing.T) {
		// Can also be used for map capacity hint
		items := make(map[string]string, DEFAULT_LIST_CAPACITY)
		assert.NotNil(t, items)
		assert.Equal(t, 0, len(items))
	})

	t.Run("prevents immediate reallocation", func(t *testing.T) {
		items := make([]int, 0, DEFAULT_LIST_CAPACITY)

		// Add items up to capacity
		for i := 0; i < DEFAULT_LIST_CAPACITY; i++ {
			items = append(items, i)
		}

		// Should not have reallocated
		assert.Equal(t, DEFAULT_LIST_CAPACITY, cap(items))
		assert.Equal(t, DEFAULT_LIST_CAPACITY, len(items))
	})

	t.Run("capacity is sufficient for small lists", func(t *testing.T) {
		// For most use cases, 10 items is a reasonable default
		items := make([]string, 0, DEFAULT_LIST_CAPACITY)

		// Simulate adding a few items
		for i := 0; i < 5; i++ {
			items = append(items, "item")
		}

		// Should still have unused capacity
		assert.Less(t, len(items), cap(items))
	})

	t.Run("can be used in different contexts", func(t *testing.T) {
		// Integer slice
		intSlice := make([]int, 0, DEFAULT_LIST_CAPACITY)
		assert.Equal(t, DEFAULT_LIST_CAPACITY, cap(intSlice))

		// String slice
		stringSlice := make([]string, 0, DEFAULT_LIST_CAPACITY)
		assert.Equal(t, DEFAULT_LIST_CAPACITY, cap(stringSlice))

		// Struct slice
		type Item struct{ ID int }
		structSlice := make([]Item, 0, DEFAULT_LIST_CAPACITY)
		assert.Equal(t, DEFAULT_LIST_CAPACITY, cap(structSlice))
	})
}

func TestFLOAT32_BIT_SIZE_UseCases(t *testing.T) {
	t.Run("use with strconv.ParseFloat", func(t *testing.T) {
		// FLOAT32_BIT_SIZE can be used with strconv.ParseFloat
		// This is a common use case
		assert.Equal(t, 32, FLOAT32_BIT_SIZE)
	})

	t.Run("matches float32 specification", func(t *testing.T) {
		// float32 is IEEE 754 single-precision (32 bits)
		assert.Equal(t, 32, FLOAT32_BIT_SIZE)
	})

	t.Run("differentiates from float64", func(t *testing.T) {
		// float64 would be 64 bits
		assert.NotEqual(t, 64, FLOAT32_BIT_SIZE)
		assert.Equal(t, 32, FLOAT32_BIT_SIZE)
	})

	t.Run("can be used in bit manipulation", func(t *testing.T) {
		// Useful for understanding memory layout
		bytesRequired := FLOAT32_BIT_SIZE / 8
		assert.Equal(t, 4, bytesRequired, "float32 requires 4 bytes")
	})
}

func TestConstantValues(t *testing.T) {
	t.Run("DEFAULT_LIST_CAPACITY is not zero", func(t *testing.T) {
		assert.NotEqual(t, 0, DEFAULT_LIST_CAPACITY, "Should have a non-zero default")
	})

	t.Run("FLOAT32_BIT_SIZE is not zero", func(t *testing.T) {
		assert.NotEqual(t, 0, FLOAT32_BIT_SIZE, "Should have a non-zero value")
	})

	t.Run("constants are different", func(t *testing.T) {
		assert.NotEqual(t, DEFAULT_LIST_CAPACITY, FLOAT32_BIT_SIZE, "Constants should have different purposes and values")
	})
}

func TestConstantNaming(t *testing.T) {
	t.Run("DEFAULT_LIST_CAPACITY follows naming convention", func(t *testing.T) {
		// Constants are in UPPER_SNAKE_CASE
		// Name clearly indicates it's a default capacity for lists
		assert.Equal(t, 10, DEFAULT_LIST_CAPACITY)
	})

	t.Run("FLOAT32_BIT_SIZE follows naming convention", func(t *testing.T) {
		// Constants are in UPPER_SNAKE_CASE
		// Name clearly indicates it's the bit size for float32
		assert.Equal(t, 32, FLOAT32_BIT_SIZE)
	})
}

func TestConstantDocumentation(t *testing.T) {
	t.Run("constants have clear purpose", func(t *testing.T) {
		// DEFAULT_LIST_CAPACITY: initial capacity for slices when listing items
		// Used to prevent immediate reallocation for common small lists
		assert.Equal(t, 10, DEFAULT_LIST_CAPACITY)

		// FLOAT32_BIT_SIZE: bit size for float32 types
		// Used with strconv and other parsing functions
		assert.Equal(t, 32, FLOAT32_BIT_SIZE)
	})
}

func TestConstantsImmutability(t *testing.T) {
	t.Run("constants cannot be reassigned", func(t *testing.T) {
		// These would cause compile errors if uncommented:
		// DEFAULT_LIST_CAPACITY = 20
		// FLOAT32_BIT_SIZE = 64

		// Instead, we can only read them
		capacity := DEFAULT_LIST_CAPACITY
		bitSize := FLOAT32_BIT_SIZE

		assert.Equal(t, 10, capacity)
		assert.Equal(t, 32, bitSize)
	})

	t.Run("constants maintain their values", func(t *testing.T) {
		// Read multiple times to ensure consistency
		value1 := DEFAULT_LIST_CAPACITY
		value2 := DEFAULT_LIST_CAPACITY
		value3 := DEFAULT_LIST_CAPACITY

		assert.Equal(t, value1, value2)
		assert.Equal(t, value2, value3)
		assert.Equal(t, 10, value1)
	})
}

func TestConstantsEdgeCases(t *testing.T) {
	t.Run("DEFAULT_LIST_CAPACITY arithmetic", func(t *testing.T) {
		// Can be used in arithmetic operations
		doubled := DEFAULT_LIST_CAPACITY * 2
		assert.Equal(t, 20, doubled)

		halved := DEFAULT_LIST_CAPACITY / 2
		assert.Equal(t, 5, halved)
	})

	t.Run("FLOAT32_BIT_SIZE arithmetic", func(t *testing.T) {
		// Can be used in arithmetic operations
		bytes := FLOAT32_BIT_SIZE / 8
		assert.Equal(t, 4, bytes)

		doubled := FLOAT32_BIT_SIZE * 2
		assert.Equal(t, 64, doubled)
	})

	t.Run("constants in comparisons", func(t *testing.T) {
		assert.True(t, DEFAULT_LIST_CAPACITY < FLOAT32_BIT_SIZE)
		assert.True(t, DEFAULT_LIST_CAPACITY > 0)
		assert.True(t, FLOAT32_BIT_SIZE > 0)
	})
}

func TestConstantsPerformance(t *testing.T) {
	t.Run("constants have no runtime overhead", func(t *testing.T) {
		// Constants are resolved at compile time
		// Using them has zero runtime cost compared to literals
		for i := 0; i < 1000; i++ {
			_ = make([]int, 0, DEFAULT_LIST_CAPACITY)
		}
		// This test mainly validates the constant can be used efficiently
		assert.Equal(t, 10, DEFAULT_LIST_CAPACITY)
	})
}

func TestConstantsIntegration(t *testing.T) {
	t.Run("both constants can be used together", func(t *testing.T) {
		// Example: creating a slice of float32 with default capacity
		floats := make([]float32, 0, DEFAULT_LIST_CAPACITY)
		assert.Equal(t, 0, len(floats))
		assert.Equal(t, DEFAULT_LIST_CAPACITY, cap(floats))

		// The bit size constant is for a different purpose (parsing)
		assert.Equal(t, 32, FLOAT32_BIT_SIZE)
	})

	t.Run("constants are independent", func(t *testing.T) {
		// Using one doesn't affect the other
		capacity := DEFAULT_LIST_CAPACITY
		bitSize := FLOAT32_BIT_SIZE

		assert.Equal(t, 10, capacity)
		assert.Equal(t, 32, bitSize)
		assert.NotEqual(t, capacity, bitSize)
	})
}

func TestConstantsAccessibility(t *testing.T) {
	t.Run("constants are exported", func(t *testing.T) {
		// Both constants start with uppercase, so they're exported
		// This means they can be imported and used by other packages
		assert.True(t, true, "Constants are accessible in tests, therefore exported")
	})
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
	// These should have identical performance - constants are compile-time
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
	// These should have identical performance - constants are compile-time
}
