package validator_test

import (
	"testing"

	"github.com/workspace/repo/internal/validator"
)

// Compile-time assertion that *ValidationError implements the error interface.
var _ error = (*validator.ValidationError)(nil)

func TestValidationError_ImplementsErrorInterface(t *testing.T) {
	// Ensure *ValidationError can be used as an error value at runtime.
	var err error = &validator.ValidationError{
		Kind:    validator.ErrUnknownOperation,
		Message: "test error",
	}
	if err == nil {
		t.Fatal("expected non-nil error")
	}
}

func TestValidationError_ErrorReturnsMessage(t *testing.T) {
	tests := []struct {
		name    string
		message string
	}{
		{"simple message", "unknown operation: foo"},
		{"empty message", ""},
		{"wrong arg count message", "wrong number of arguments: expected 2, got 1"},
		{"not a number message", "argument is not a number: abc"},
		{"division by zero message", "division by zero"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := &validator.ValidationError{
				Kind:    validator.ErrUnknownOperation,
				Message: tc.message,
			}
			if got := err.Error(); got != tc.message {
				t.Errorf("Error() = %q, want %q", got, tc.message)
			}
		})
	}
}

func TestValidationError_ErrorReturnsMessageField(t *testing.T) {
	// Explicitly verify that Error() returns exactly the Message field, not some
	// other formatted representation.
	const msg = "division by zero"
	err := &validator.ValidationError{
		Kind:    validator.ErrDivisionByZero,
		Message: msg,
	}
	if got := err.Error(); got != err.Message {
		t.Errorf("Error() = %q, want Message field %q", got, err.Message)
	}
}

func TestErrorKind_ConstantsAreDistinct(t *testing.T) {
	constants := []struct {
		name  string
		value validator.ErrorKind
	}{
		{"ErrUnknownOperation", validator.ErrUnknownOperation},
		{"ErrWrongArgCount", validator.ErrWrongArgCount},
		{"ErrNotANumber", validator.ErrNotANumber},
		{"ErrDivisionByZero", validator.ErrDivisionByZero},
	}

	seen := make(map[validator.ErrorKind]string)
	for _, c := range constants {
		if prev, exists := seen[c.value]; exists {
			t.Errorf("ErrorKind %s has the same integer value as %s (%d); all constants must be distinct",
				c.name, prev, int(c.value))
		}
		seen[c.value] = c.name
	}
}

func TestErrorKind_AllConstantsDefined(t *testing.T) {
	// Verify all four expected error kind constants are accessible and have
	// sensible (zero-based iota) integer values.
	if validator.ErrUnknownOperation != 0 {
		t.Errorf("ErrUnknownOperation = %d, want 0 (first iota)", int(validator.ErrUnknownOperation))
	}
	if validator.ErrWrongArgCount != 1 {
		t.Errorf("ErrWrongArgCount = %d, want 1", int(validator.ErrWrongArgCount))
	}
	if validator.ErrNotANumber != 2 {
		t.Errorf("ErrNotANumber = %d, want 2", int(validator.ErrNotANumber))
	}
	if validator.ErrDivisionByZero != 3 {
		t.Errorf("ErrDivisionByZero = %d, want 3", int(validator.ErrDivisionByZero))
	}
}

func TestValidatedCommand_Fields(t *testing.T) {
	// Verify ValidatedCommand has the expected public fields.
	cmd := validator.ValidatedCommand{
		Operation: "add",
		Args:      []float64{1.0, 2.0},
	}
	if cmd.Operation != "add" {
		t.Errorf("Operation = %q, want %q", cmd.Operation, "add")
	}
	if len(cmd.Args) != 2 {
		t.Fatalf("len(Args) = %d, want 2", len(cmd.Args))
	}
	if cmd.Args[0] != 1.0 || cmd.Args[1] != 2.0 {
		t.Errorf("Args = %v, want [1.0 2.0]", cmd.Args)
	}
}

func TestValidationError_KindFieldIsAccessible(t *testing.T) {
	// Verify the Kind field can be set and read for each error kind.
	kinds := []struct {
		name string
		kind validator.ErrorKind
	}{
		{"ErrUnknownOperation", validator.ErrUnknownOperation},
		{"ErrWrongArgCount", validator.ErrWrongArgCount},
		{"ErrNotANumber", validator.ErrNotANumber},
		{"ErrDivisionByZero", validator.ErrDivisionByZero},
	}

	for _, tc := range kinds {
		t.Run(tc.name, func(t *testing.T) {
			err := &validator.ValidationError{Kind: tc.kind, Message: "msg"}
			if err.Kind != tc.kind {
				t.Errorf("Kind = %d, want %d", int(err.Kind), int(tc.kind))
			}
		})
	}
}
