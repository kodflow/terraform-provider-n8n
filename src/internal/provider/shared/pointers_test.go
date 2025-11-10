package shared

import (
	"testing"
)

// TestPtr tests the generic Ptr function.
func TestPtr(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
	}{
		{"string pointer"},
		{"int pointer"},
		{"bool pointer"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			switch tc.name {
			case "string pointer":
				s := "test"
				ptr := Ptr(s)
				if ptr == nil {
					t.Fatal("Ptr() returned nil")
				}
				if *ptr != s {
					t.Errorf("Ptr() = %v, want %v", *ptr, s)
				}
			case "int pointer":
				i := 42
				ptr := Ptr(i)
				if ptr == nil {
					t.Fatal("Ptr() returned nil")
				}
				if *ptr != i {
					t.Errorf("Ptr() = %v, want %v", *ptr, i)
				}
			case "bool pointer":
				b := true
				ptr := Ptr(b)
				if ptr == nil {
					t.Fatal("Ptr() returned nil")
				}
				if *ptr != b {
					t.Errorf("Ptr() = %v, want %v", *ptr, b)
				}
			}
		})
	}
}

// TestString tests the String function.
func TestString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"normal string", "test", "test"},
		{"empty string", "", ""},
		{"unicode string", "Hello 世界", "Hello 世界"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ptr := String(tt.input)
			if ptr == nil {
				t.Fatal("String() returned nil")
			}
			if *ptr != tt.want {
				t.Errorf("String() = %v, want %v", *ptr, tt.want)
			}
		})
	}
}

// TestBool tests the Bool function.
func TestBool(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input bool
		want  bool
	}{
		{"true", true, true},
		{"false", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ptr := Bool(tt.input)
			if ptr == nil {
				t.Fatal("Bool() returned nil")
			}
			if *ptr != tt.want {
				t.Errorf("Bool() = %v, want %v", *ptr, tt.want)
			}
		})
	}
}

// TestInt tests the Int function.
func TestInt(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input int
		want  int
	}{
		{"positive", 42, 42},
		{"zero", 0, 0},
		{"negative", -10, -10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ptr := Int(tt.input)
			if ptr == nil {
				t.Fatal("Int() returned nil")
			}
			if *ptr != tt.want {
				t.Errorf("Int() = %v, want %v", *ptr, tt.want)
			}
		})
	}
}

// TestInt32 tests the Int32 function.
func TestInt32(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input int32
		want  int32
	}{
		{"positive", 42, 42},
		{"zero", 0, 0},
		{"negative", -10, -10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ptr := Int32(tt.input)
			if ptr == nil {
				t.Fatal("Int32() returned nil")
			}
			if *ptr != tt.want {
				t.Errorf("Int32() = %v, want %v", *ptr, tt.want)
			}
		})
	}
}

// TestFloat32 tests the Float32 function.
func TestFloat32(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input float32
		want  float32
	}{
		{"positive", 3.14, 3.14},
		{"zero", 0.0, 0.0},
		{"negative", -2.5, -2.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ptr := Float32(tt.input)
			if ptr == nil {
				t.Fatal("Float32() returned nil")
			}
			if *ptr != tt.want {
				t.Errorf("Float32() = %v, want %v", *ptr, tt.want)
			}
		})
	}
}

// TestPointerFunctions tests basic pointer creation.
func TestPointerFunctions(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
	}{
		{"String"},
		{"Bool"},
		{"Int"},
		{"Int32"},
		{"Float32"},
		{"Ptr generic"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			switch tc.name {
			case "String":
				s := "test"
				ptr := String(s)
				if ptr == nil || *ptr != s {
					t.Errorf("String() failed: expected %q, got %v", s, ptr)
				}
			case "Bool":
				b := true
				ptr := Bool(b)
				if ptr == nil || *ptr != b {
					t.Errorf("Bool() failed: expected %v, got %v", b, ptr)
				}
			case "Int":
				i := 42
				ptr := Int(i)
				if ptr == nil || *ptr != i {
					t.Errorf("Int() failed: expected %d, got %v", i, ptr)
				}
			case "Int32":
				var i int32 = 42
				ptr := Int32(i)
				if ptr == nil || *ptr != i {
					t.Errorf("Int32() failed: expected %d, got %v", i, ptr)
				}
			case "Float32":
				var f float32 = 3.14
				ptr := Float32(f)
				if ptr == nil || *ptr != f {
					t.Errorf("Float32() failed: expected %f, got %v", f, ptr)
				}
			case "Ptr generic":
				s := "generic"
				ptr := Ptr(s)
				if ptr == nil || *ptr != s {
					t.Errorf("Ptr() failed: expected %q, got %v", s, ptr)
				}
			}
		})
	}
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
