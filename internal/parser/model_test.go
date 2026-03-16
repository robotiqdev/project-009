package parser_test

import (
	"testing"

	"github.com/example/feat1123/internal/parser"
)

// TestParsedInputStructExists verifies that the ParsedInput struct exists and is importable.
func TestParsedInputStructExists(t *testing.T) {
	var p parser.ParsedInput
	_ = p
}

// TestParsedInputOperationField verifies that ParsedInput has an Operation field of type string.
func TestParsedInputOperationField(t *testing.T) {
	p := parser.ParsedInput{Operation: "add"}
	if p.Operation != "add" {
		t.Errorf("expected Operation to be 'add', got %q", p.Operation)
	}
}

// TestParsedInputOperand1Field verifies that ParsedInput has an Operand1 field of type float64.
func TestParsedInputOperand1Field(t *testing.T) {
	p := parser.ParsedInput{Operand1: 3.14}
	if p.Operand1 != 3.14 {
		t.Errorf("expected Operand1 to be 3.14, got %f", p.Operand1)
	}
}

// TestParsedInputOperand2Field verifies that ParsedInput has an Operand2 field of type float64.
func TestParsedInputOperand2Field(t *testing.T) {
	p := parser.ParsedInput{Operand2: 2.71}
	if p.Operand2 != 2.71 {
		t.Errorf("expected Operand2 to be 2.71, got %f", p.Operand2)
	}
}

// TestParsedInputZeroValue verifies that a zero-value ParsedInput is valid (no constructor needed).
func TestParsedInputZeroValue(t *testing.T) {
	var p parser.ParsedInput
	if p.Operation != "" {
		t.Errorf("expected zero-value Operation to be empty string, got %q", p.Operation)
	}
	if p.Operand1 != 0 {
		t.Errorf("expected zero-value Operand1 to be 0, got %f", p.Operand1)
	}
	if p.Operand2 != 0 {
		t.Errorf("expected zero-value Operand2 to be 0, got %f", p.Operand2)
	}
}

// TestParsedInputAllFieldsSet verifies that all three fields can be set together.
func TestParsedInputAllFieldsSet(t *testing.T) {
	p := parser.ParsedInput{
		Operation: "multiply",
		Operand1:  6.0,
		Operand2:  7.0,
	}
	if p.Operation != "multiply" {
		t.Errorf("expected Operation 'multiply', got %q", p.Operation)
	}
	if p.Operand1 != 6.0 {
		t.Errorf("expected Operand1 6.0, got %f", p.Operand1)
	}
	if p.Operand2 != 7.0 {
		t.Errorf("expected Operand2 7.0, got %f", p.Operand2)
	}
}

// TestParsedInputFieldTypes verifies that the field types are assignable from the expected types.
// This is a compile-time check via assignment.
func TestParsedInputFieldTypes(t *testing.T) {
	var op string = "subtract"
	var op1 float64 = 10.5
	var op2 float64 = 4.5

	p := parser.ParsedInput{
		Operation: op,
		Operand1:  op1,
		Operand2:  op2,
	}

	if p.Operation != op {
		t.Errorf("Operation field type mismatch: expected %q, got %q", op, p.Operation)
	}
	if p.Operand1 != op1 {
		t.Errorf("Operand1 field type mismatch: expected %f, got %f", op1, p.Operand1)
	}
	if p.Operand2 != op2 {
		t.Errorf("Operand2 field type mismatch: expected %f, got %f", op2, p.Operand2)
	}
}
