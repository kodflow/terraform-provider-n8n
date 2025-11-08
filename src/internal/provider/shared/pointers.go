// Package shared provides utilities shared across the provider.
// The pointer utilities are inspired by github.com/kitsunium/sdk/pkg/lib/pointer
// but implemented locally for Bazel compatibility.
package shared

// Ptr returns a pointer to the given value.
// This is a generic function that works with any type.
//
//go:inline
func Ptr[T any](v T) *T {
	return &v
}

// String returns a pointer to the given string.
//
//go:inline
func String(s string) *string {
	return &s
}

// Bool returns a pointer to the given bool.
//
//go:inline
func Bool(b bool) *bool {
	return &b
}

// Int returns a pointer to the given int.
//
//go:inline
func Int(i int) *int {
	return &i
}

// Int32 returns a pointer to the given int32.
//
//go:inline
func Int32(i int32) *int32 {
	return &i
}

// Float32 returns a pointer to the given float32.
//
//go:inline
func Float32(f float32) *float32 {
	return &f
}
