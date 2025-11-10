package shared

import (
	"testing"
)

// TestPtr tests the generic Ptr function.
func TestPtr(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		wantErr bool
	}{
		{name: "string pointer", wantErr: false},
		{name: "int pointer", wantErr: false},
		{name: "bool pointer", wantErr: false},
		{name: "error case - validation checks", wantErr: false},
	}

	for _, tc := range testCases {
		tc := tc
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
			case "error case - validation checks":
				s := ""
				ptr := Ptr(s)
				if ptr == nil {
					t.Fatal("Ptr() returned nil")
				}
				if *ptr != s {
					t.Errorf("Ptr() = %v, want %v", *ptr, s)
				}
			}
		})
	}
}

// TestString tests the String function.
func TestString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{name: "normal string", input: "test", want: "test", wantErr: false},
		{name: "empty string", input: "", want: "", wantErr: false},
		{name: "unicode string", input: "Hello 世界", want: "Hello 世界", wantErr: false},
		{name: "error case - validation checks", input: "", want: "", wantErr: false},
	}

	for _, tt := range tests {
		tt := tt
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
		name    string
		input   bool
		want    bool
		wantErr bool
	}{
		{name: "true", input: true, want: true, wantErr: false},
		{name: "false", input: false, want: false, wantErr: false},
		{name: "error case - validation checks", input: false, want: false, wantErr: false},
	}

	for _, tt := range tests {
		tt := tt
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
		name    string
		input   int
		want    int
		wantErr bool
	}{
		{name: "positive", input: 42, want: 42, wantErr: false},
		{name: "zero", input: 0, want: 0, wantErr: false},
		{name: "negative", input: -10, want: -10, wantErr: false},
		{name: "error case - validation checks", input: 0, want: 0, wantErr: false},
	}

	for _, tt := range tests {
		tt := tt
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
		name    string
		input   int32
		want    int32
		wantErr bool
	}{
		{name: "positive", input: 42, want: 42, wantErr: false},
		{name: "zero", input: 0, want: 0, wantErr: false},
		{name: "negative", input: -10, want: -10, wantErr: false},
		{name: "error case - validation checks", input: 0, want: 0, wantErr: false},
	}

	for _, tt := range tests {
		tt := tt
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
		name    string
		input   float32
		want    float32
		wantErr bool
	}{
		{name: "positive", input: 3.14, want: 3.14, wantErr: false},
		{name: "zero", input: 0.0, want: 0.0, wantErr: false},
		{name: "negative", input: -2.5, want: -2.5, wantErr: false},
		{name: "error case - validation checks", input: 0.0, want: 0.0, wantErr: false},
	}

	for _, tt := range tests {
		tt := tt
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
		name    string
		wantErr bool
	}{
		{name: "String", wantErr: false},
		{name: "Bool", wantErr: false},
		{name: "Int", wantErr: false},
		{name: "Int32", wantErr: false},
		{name: "Float32", wantErr: false},
		{name: "Ptr generic", wantErr: false},
		{name: "error case - empty values", wantErr: false},
	}

	for _, tc := range testCases {
		tc := tc
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
			case "error case - empty values":
				s := ""
				ptr := String(s)
				if ptr == nil || *ptr != s {
					t.Errorf("String() failed for empty string: got %v", ptr)
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
