package shared

import (
	"testing"
)

// TestPointerFunctions tests basic pointer creation.
func TestPointerFunctions(t *testing.T) {
	t.Parallel()

	t.Run("String", func(t *testing.T) {
		t.Parallel()
		s := "test"
		ptr := String(s)
		if ptr == nil || *ptr != s {
			t.Errorf("String() failed: expected %q, got %v", s, ptr)
		}
	})

	t.Run("Bool", func(t *testing.T) {
		t.Parallel()
		b := true
		ptr := Bool(b)
		if ptr == nil || *ptr != b {
			t.Errorf("Bool() failed: expected %v, got %v", b, ptr)
		}
	})

	t.Run("Int", func(t *testing.T) {
		t.Parallel()
		i := 42
		ptr := Int(i)
		if ptr == nil || *ptr != i {
			t.Errorf("Int() failed: expected %d, got %v", i, ptr)
		}
	})

	t.Run("Int32", func(t *testing.T) {
		t.Parallel()
		var i int32 = 42
		ptr := Int32(i)
		if ptr == nil || *ptr != i {
			t.Errorf("Int32() failed: expected %d, got %v", i, ptr)
		}
	})

	t.Run("Float32", func(t *testing.T) {
		t.Parallel()
		var f float32 = 3.14
		ptr := Float32(f)
		if ptr == nil || *ptr != f {
			t.Errorf("Float32() failed: expected %f, got %v", f, ptr)
		}
	})

	t.Run("Ptr generic", func(t *testing.T) {
		t.Parallel()
		s := "generic"
		ptr := Ptr(s)
		if ptr == nil || *ptr != s {
			t.Errorf("Ptr() failed: expected %q, got %v", s, ptr)
		}
	})
}

// BenchmarkString benchmarks the String function.
// This ensures inlining is working - should be < 1 ns/op.
func BenchmarkString(b *testing.B) {
	s := "test"
	b.ResetTimer()
	for b.Loop() {
		_ = String(s)
	}
}

// BenchmarkBool benchmarks the Bool function.
func BenchmarkBool(b *testing.B) {
	v := true
	b.ResetTimer()
	for b.Loop() {
		_ = Bool(v)
	}
}

// BenchmarkInt benchmarks the Int function.
func BenchmarkInt(b *testing.B) {
	v := 42
	b.ResetTimer()
	for b.Loop() {
		_ = Int(v)
	}
}

// BenchmarkPtr benchmarks the generic Ptr function.
func BenchmarkPtr(b *testing.B) {
	s := "test"
	b.ResetTimer()
	for b.Loop() {
		_ = Ptr(s)
	}
}

// BenchmarkStringDirect benchmarks direct pointer creation.
// This is the baseline for comparison.
func BenchmarkStringDirect(b *testing.B) {
	s := "test"
	b.ResetTimer()
	for b.Loop() {
		_ = &s
	}
}
