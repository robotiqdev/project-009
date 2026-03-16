package arithmetic_test

import (
	"testing"

	"github.com/workspace/repo/internal/arithmetic"
)

func TestResultStructFields(t *testing.T) {
	r := arithmetic.Result{Value: 3.14, Formatted: "3.14"}

	if r.Value != 3.14 {
		t.Errorf("Result.Value = %v, want 3.14", r.Value)
	}
	if r.Formatted != "3.14" {
		t.Errorf("Result.Formatted = %q, want %q", r.Formatted, "3.14")
	}
}

func TestResultValueField_IsFloat64(t *testing.T) {
	var r arithmetic.Result
	r.Value = 0.0
	r.Value = 1.23456789
	if r.Value != 1.23456789 {
		t.Errorf("Result.Value = %v, want 1.23456789", r.Value)
	}
}

func TestResultFormattedField_IsString(t *testing.T) {
	var r arithmetic.Result
	r.Formatted = "0.00"
	if r.Formatted != "0.00" {
		t.Errorf("Result.Formatted = %q, want %q", r.Formatted, "0.00")
	}
}

func TestResultZeroValue(t *testing.T) {
	var r arithmetic.Result
	if r.Value != 0.0 {
		t.Errorf("zero Result.Value = %v, want 0.0", r.Value)
	}
	if r.Formatted != "" {
		t.Errorf("zero Result.Formatted = %q, want empty string", r.Formatted)
	}
}

func TestResultNegativeValue(t *testing.T) {
	r := arithmetic.Result{Value: -1.5, Formatted: "-1.5"}
	if r.Value != -1.5 {
		t.Errorf("Result.Value = %v, want -1.5", r.Value)
	}
	if r.Formatted != "-1.5" {
		t.Errorf("Result.Formatted = %q, want %q", r.Formatted, "-1.5")
	}
}

func TestEngineInterface_Exists(t *testing.T) {
	// Verify that Engine is an interface with Execute method by checking
	// that a nil pointer of a type implementing it can be assigned.
	var _ arithmetic.Engine = (*mockEngine)(nil)
}

// mockEngine is a local test-only implementation to verify the Engine interface shape.
type mockEngine struct{}

func (m *mockEngine) Execute(operation string, args []float64) (arithmetic.Result, error) {
	return arithmetic.Result{}, nil
}
