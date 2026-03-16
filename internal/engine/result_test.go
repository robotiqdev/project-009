package engine_test

import (
	"errors"
	"testing"

	"calculator/internal/engine"
)

// TestResultStruct verifies that Result is a struct with a float64 Value field.
func TestResultStruct(t *testing.T) {
	r := engine.Result{Value: 42.0}
	if r.Value != 42.0 {
		t.Errorf("expected Result.Value to be 42.0, got %f", r.Value)
	}
}

// TestResultZeroValue verifies that the zero value of Result has Value == 0.
func TestResultZeroValue(t *testing.T) {
	var r engine.Result
	if r.Value != 0 {
		t.Errorf("expected zero-value Result.Value to be 0, got %f", r.Value)
	}
}

// TestResultNegativeValue verifies that Result can hold negative float64 values.
func TestResultNegativeValue(t *testing.T) {
	r := engine.Result{Value: -3.14}
	if r.Value != -3.14 {
		t.Errorf("expected Result.Value to be -3.14, got %f", r.Value)
	}
}

// TestErrDivisionByZeroIsNotNil verifies that the sentinel error is declared and non-nil.
func TestErrDivisionByZeroIsNotNil(t *testing.T) {
	if engine.ErrDivisionByZero == nil {
		t.Error("expected ErrDivisionByZero to be non-nil")
	}
}

// TestErrUnknownOperationIsNotNil verifies that the sentinel error is declared and non-nil.
func TestErrUnknownOperationIsNotNil(t *testing.T) {
	if engine.ErrUnknownOperation == nil {
		t.Error("expected ErrUnknownOperation to be non-nil")
	}
}

// TestErrDivisionByZeroMessage verifies the error message text.
func TestErrDivisionByZeroMessage(t *testing.T) {
	want := "division by zero"
	if engine.ErrDivisionByZero.Error() != want {
		t.Errorf("expected ErrDivisionByZero message %q, got %q", want, engine.ErrDivisionByZero.Error())
	}
}

// TestErrUnknownOperationMessage verifies the error message text.
func TestErrUnknownOperationMessage(t *testing.T) {
	want := "unknown operation"
	if engine.ErrUnknownOperation.Error() != want {
		t.Errorf("expected ErrUnknownOperation message %q, got %q", want, engine.ErrUnknownOperation.Error())
	}
}

// TestErrDivisionByZeroErrorsIs verifies that errors.Is works with the sentinel error.
func TestErrDivisionByZeroErrorsIs(t *testing.T) {
	err := engine.ErrDivisionByZero
	if !errors.Is(err, engine.ErrDivisionByZero) {
		t.Error("expected errors.Is to return true for ErrDivisionByZero")
	}
}

// TestErrUnknownOperationErrorsIs verifies that errors.Is works with the sentinel error.
func TestErrUnknownOperationErrorsIs(t *testing.T) {
	err := engine.ErrUnknownOperation
	if !errors.Is(err, engine.ErrUnknownOperation) {
		t.Error("expected errors.Is to return true for ErrUnknownOperation")
	}
}

// TestSentinelErrorsAreDistinct verifies the two sentinel errors are not equal.
func TestSentinelErrorsAreDistinct(t *testing.T) {
	if errors.Is(engine.ErrDivisionByZero, engine.ErrUnknownOperation) {
		t.Error("expected ErrDivisionByZero and ErrUnknownOperation to be distinct errors")
	}
	if errors.Is(engine.ErrUnknownOperation, engine.ErrDivisionByZero) {
		t.Error("expected ErrUnknownOperation and ErrDivisionByZero to be distinct errors")
	}
}
