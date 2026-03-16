package arithmetic_test

import (
	"testing"

	"github.com/workspace/repo/internal/arithmetic"
)

// TestNewCalculator verifies that NewCalculator returns a non-nil *Calculator
// that implements the Engine interface.
func TestNewCalculator_ReturnsEngine(t *testing.T) {
	calc := arithmetic.NewCalculator()
	if calc == nil {
		t.Fatal("NewCalculator() returned nil")
	}
	var _ arithmetic.Engine = calc
}

// TestCalculator_Execute_Add tests the "add" operation via Execute.
func TestCalculator_Execute_Add(t *testing.T) {
	tests := []struct {
		name          string
		args          []float64
		wantValue     float64
		wantFormatted string
	}{
		{
			name:          "add(1.5, 2.5) = 4.0",
			args:          []float64{1.5, 2.5},
			wantValue:     4.0,
			wantFormatted: "4.00",
		},
		{
			name:          "add(0.0, 0.0) = 0.0",
			args:          []float64{0.0, 0.0},
			wantValue:     0.0,
			wantFormatted: "0.00",
		},
		{
			name:          "add(-1.0, 1.0) = 0.0",
			args:          []float64{-1.0, 1.0},
			wantValue:     0.0,
			wantFormatted: "0.00",
		},
	}
	calc := arithmetic.NewCalculator()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := calc.Execute("add", tt.args)
			if err != nil {
				t.Fatalf("Execute(\"add\", %v) unexpected error: %v", tt.args, err)
			}
			if result.Value != tt.wantValue {
				t.Errorf("Execute(\"add\", %v).Value = %v, want %v", tt.args, result.Value, tt.wantValue)
			}
			if result.Formatted != tt.wantFormatted {
				t.Errorf("Execute(\"add\", %v).Formatted = %q, want %q", tt.args, result.Formatted, tt.wantFormatted)
			}
		})
	}
}

// TestCalculator_Execute_Subtract tests the "subtract" operation via Execute.
func TestCalculator_Execute_Subtract(t *testing.T) {
	tests := []struct {
		name          string
		args          []float64
		wantValue     float64
		wantFormatted string
	}{
		{
			name:          "subtract(5.0, 3.0) = 2.0",
			args:          []float64{5.0, 3.0},
			wantValue:     2.0,
			wantFormatted: "2.00",
		},
		{
			name:          "subtract(1.0, 5.0) = -4.0 (negative result)",
			args:          []float64{1.0, 5.0},
			wantValue:     -4.0,
			wantFormatted: "-4.00",
		},
		{
			name:          "subtract(0.0, 0.0) = 0.0",
			args:          []float64{0.0, 0.0},
			wantValue:     0.0,
			wantFormatted: "0.00",
		},
	}
	calc := arithmetic.NewCalculator()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := calc.Execute("subtract", tt.args)
			if err != nil {
				t.Fatalf("Execute(\"subtract\", %v) unexpected error: %v", tt.args, err)
			}
			if result.Value != tt.wantValue {
				t.Errorf("Execute(\"subtract\", %v).Value = %v, want %v", tt.args, result.Value, tt.wantValue)
			}
			if result.Formatted != tt.wantFormatted {
				t.Errorf("Execute(\"subtract\", %v).Formatted = %q, want %q", tt.args, result.Formatted, tt.wantFormatted)
			}
		})
	}
}

// TestCalculator_Execute_Multiply tests the "multiply" operation via Execute.
func TestCalculator_Execute_Multiply(t *testing.T) {
	tests := []struct {
		name          string
		args          []float64
		wantValue     float64
		wantFormatted string
	}{
		{
			name:          "multiply(2.0, 3.5) = 7.0",
			args:          []float64{2.0, 3.5},
			wantValue:     7.0,
			wantFormatted: "7.00",
		},
		{
			name:          "multiply by zero: multiply(5.0, 0.0) = 0.0",
			args:          []float64{5.0, 0.0},
			wantValue:     0.0,
			wantFormatted: "0.00",
		},
		{
			name:          "multiply by zero: multiply(0.0, 3.0) = 0.0",
			args:          []float64{0.0, 3.0},
			wantValue:     0.0,
			wantFormatted: "0.00",
		},
		{
			name:          "multiply negatives: multiply(-2.0, 3.0) = -6.0",
			args:          []float64{-2.0, 3.0},
			wantValue:     -6.0,
			wantFormatted: "-6.00",
		},
	}
	calc := arithmetic.NewCalculator()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := calc.Execute("multiply", tt.args)
			if err != nil {
				t.Fatalf("Execute(\"multiply\", %v) unexpected error: %v", tt.args, err)
			}
			if result.Value != tt.wantValue {
				t.Errorf("Execute(\"multiply\", %v).Value = %v, want %v", tt.args, result.Value, tt.wantValue)
			}
			if result.Formatted != tt.wantFormatted {
				t.Errorf("Execute(\"multiply\", %v).Formatted = %q, want %q", tt.args, result.Formatted, tt.wantFormatted)
			}
		})
	}
}

// TestCalculator_Execute_Divide tests the "divide" operation via Execute.
func TestCalculator_Execute_Divide(t *testing.T) {
	tests := []struct {
		name          string
		args          []float64
		wantValue     float64
		wantFormatted string
	}{
		{
			name:          "divide(10.0, 4.0) = 2.5",
			args:          []float64{10.0, 4.0},
			wantValue:     2.5,
			wantFormatted: "2.50",
		},
		{
			name:          "divide(1.0, 3.0) = 0.33 (formatted)",
			args:          []float64{1.0, 3.0},
			wantFormatted: "0.33",
		},
		{
			name:          "divide(6.0, 2.0) = 3.0",
			args:          []float64{6.0, 2.0},
			wantValue:     3.0,
			wantFormatted: "3.00",
		},
	}
	calc := arithmetic.NewCalculator()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := calc.Execute("divide", tt.args)
			if err != nil {
				t.Fatalf("Execute(\"divide\", %v) unexpected error: %v", tt.args, err)
			}
			// Only check Value when it's explicitly set (non-zero or explicitly 0)
			if tt.wantValue != 0 && result.Value != tt.wantValue {
				t.Errorf("Execute(\"divide\", %v).Value = %v, want %v", tt.args, result.Value, tt.wantValue)
			}
			if result.Formatted != tt.wantFormatted {
				t.Errorf("Execute(\"divide\", %v).Formatted = %q, want %q", tt.args, result.Formatted, tt.wantFormatted)
			}
		})
	}
}

// TestCalculator_Execute_UnknownOperation verifies that Execute returns an error
// for an unrecognized operation string.
func TestCalculator_Execute_UnknownOperation(t *testing.T) {
	tests := []struct {
		name      string
		operation string
		args      []float64
	}{
		{
			name:      "unknown operation 'modulo'",
			operation: "modulo",
			args:      []float64{5.0, 2.0},
		},
		{
			name:      "empty operation string",
			operation: "",
			args:      []float64{1.0, 2.0},
		},
		{
			name:      "unknown operation 'power'",
			operation: "power",
			args:      []float64{2.0, 3.0},
		},
	}
	calc := arithmetic.NewCalculator()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := calc.Execute(tt.operation, tt.args)
			if err == nil {
				t.Errorf("Execute(%q, %v) expected error for unknown operation, got nil", tt.operation, tt.args)
			}
		})
	}
}

// TestCalculator_ImplementsEngine verifies Calculator satisfies the Engine interface.
func TestCalculator_ImplementsEngine(t *testing.T) {
	var _ arithmetic.Engine = arithmetic.NewCalculator()
}

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
