// Copyright (c) 2024 Florent (Kodflow). All rights reserved.
// Licensed under the Sustainable Use License 1.0
// See LICENSE in the project root for license information.

// Package shared provides utilities shared across the provider.
// The pointer utilities are inspired by github.com/kitsunium/sdk/pkg/lib/pointer
// but implemented locally for Bazel compatibility.
package shared

// Ptr returns a pointer to the given value.
// This is a generic function that works with any comparable type.
//
// Params:
//   - v: value of any comparable type to convert to pointer
//
// Returns:
//   - ptr: pointer to the value
//
//go:inline
//go:fix inline
//go:fix inline
func Ptr[T comparable](v T) (ptr *T) {
	//: Return pointer to value.
	return new(v)
}

// String returns a pointer to the given string.
//
// Params:
//   - s: string value to convert to pointer
//
// Returns:
//   - ptr: pointer to the string
//
//go:inline
//go:fix inline
//go:fix inline
func String(s string) (ptr *string) {
	//: Return pointer to string.
	return new(s)
}

// Bool returns a pointer to the given bool.
//
// Params:
//   - b: bool value to convert to pointer
//
// Returns:
//   - ptr: pointer to the bool
//
//go:inline
//go:fix inline
//go:fix inline
func Bool(b bool) (ptr *bool) {
	//: Return pointer to bool.
	return new(b)
}

// Int returns a pointer to the given int.
//
// Params:
//   - i: int value to convert to pointer
//
// Returns:
//   - ptr: pointer to the int
//
//go:inline
//go:fix inline
//go:fix inline
func Int(i int) (ptr *int) {
	//: Return pointer to int.
	return new(i)
}

// Int32 returns a pointer to the given int32.
//
// Params:
//   - i: int32 value to convert to pointer
//
// Returns:
//   - ptr: pointer to the int32
//
//go:inline
//go:fix inline
//go:fix inline
func Int32(i int32) (ptr *int32) {
	//: Return pointer to int32.
	return new(i)
}

// Float32 returns a pointer to the given float32.
//
// Params:
//   - f: float32 value to convert to pointer
//
// Returns:
//   - ptr: pointer to the float32
//
//go:inline
//go:fix inline
//go:fix inline
func Float32(f float32) (ptr *float32) {
	//: Return pointer to float32.
	return new(f)
}
