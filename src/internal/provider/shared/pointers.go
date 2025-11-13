// Package shared provides utilities shared across the provider.
// The pointer utilities are inspired by github.com/kitsunium/sdk/pkg/lib/pointer
// but implemented locally for Bazel compatibility.
package shared

// Ptr returns a pointer to the given value.
// This is a generic function that works with any type.
//
// Params:
//   - v: value of any type to convert to pointer
//
// Returns:
//   - *T: pointer to the value
//
//go:inline
func Ptr[T any](v T) *T {
	// Return pointer to value.
	return &v
}

// String returns a pointer to the given string.
//
// Params:
//   - s: string value to convert to pointer
//
// Returns:
//   - *string: pointer to the string
//
//go:inline
func String(s string) *string {
	// Return pointer to string.
	return &s
}

// Bool returns a pointer to the given bool.
//
// Params:
//   - b: bool value to convert to pointer
//
// Returns:
//   - *bool: pointer to the bool
//
//go:inline
func Bool(b bool) *bool {
	// Return pointer to bool.
	return &b
}

// Int returns a pointer to the given int.
//
// Params:
//   - i: int value to convert to pointer
//
// Returns:
//   - *int: pointer to the int
//
//go:inline
func Int(i int) *int {
	// Return pointer to int.
	return &i
}

// Int32 returns a pointer to the given int32.
//
// Params:
//   - i: int32 value to convert to pointer
//
// Returns:
//   - *int32: pointer to the int32
//
//go:inline
func Int32(i int32) *int32 {
	// Return pointer to int32.
	return &i
}

// Float32 returns a pointer to the given float32.
//
// Params:
//   - f: float32 value to convert to pointer
//
// Returns:
//   - *float32: pointer to the float32
//
//go:inline
func Float32(f float32) *float32 {
	// Return pointer to float32.
	return &f
}
