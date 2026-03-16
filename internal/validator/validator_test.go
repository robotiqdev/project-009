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

// dummyArgs are valid numeric string args used to avoid triggering arg-count
// or not-a-number validators when testing operation-name validation only.
var dummyArgs = []string{"1", "2"}

// TestNewCommandValidator_ReturnsNonNil verifies the constructor returns a usable value.
func TestNewCommandValidator_ReturnsNonNil(t *testing.T) {
	v := validator.NewCommandValidator()
	if v == nil {
		t.Fatal("NewCommandValidator() returned nil")
	}
}

// TestValidate_ValidOperations verifies that supported operation names produce no error.
func TestValidate_ValidOperations(t *testing.T) {
	validOps := []string{"add", "subtract", "multiply", "divide"}

	v := validator.NewCommandValidator()

	for _, op := range validOps {
		t.Run(op, func(t *testing.T) {
			_, err := v.Validate(validator.Command{Operation: op, Args: dummyArgs})
			if err != nil {
				t.Errorf("Validate(%q) returned unexpected error: %v", op, err)
			}
		})
	}
}

// TestValidate_UnknownOperation_UppercaseVariant verifies that 'ADD' (wrong case)
// returns a ValidationError with Kind == ErrUnknownOperation (operations are case-sensitive).
func TestValidate_UnknownOperation_UppercaseVariant(t *testing.T) {
	v := validator.NewCommandValidator()

	_, err := v.Validate(validator.Command{Operation: "ADD", Args: dummyArgs})
	if err == nil {
		t.Fatal("Validate(\"ADD\") expected an error, got nil")
	}

	ve, ok := err.(*validator.ValidationError)
	if !ok {
		t.Fatalf("Validate(\"ADD\") returned error of type %T, want *validator.ValidationError", err)
	}
	if ve.Kind != validator.ErrUnknownOperation {
		t.Errorf("ValidationError.Kind = %d, want ErrUnknownOperation (%d)", int(ve.Kind), int(validator.ErrUnknownOperation))
	}
}

// TestValidate_UnknownOperation_Mod verifies that 'mod' (unsupported operation)
// returns a ValidationError with Kind == ErrUnknownOperation.
func TestValidate_UnknownOperation_Mod(t *testing.T) {
	v := validator.NewCommandValidator()

	_, err := v.Validate(validator.Command{Operation: "mod", Args: dummyArgs})
	if err == nil {
		t.Fatal("Validate(\"mod\") expected an error, got nil")
	}

	ve, ok := err.(*validator.ValidationError)
	if !ok {
		t.Fatalf("Validate(\"mod\") returned error of type %T, want *validator.ValidationError", err)
	}
	if ve.Kind != validator.ErrUnknownOperation {
		t.Errorf("ValidationError.Kind = %d, want ErrUnknownOperation (%d)", int(ve.Kind), int(validator.ErrUnknownOperation))
	}
}

// TestValidate_UnknownOperation_Empty verifies that an empty operation string
// returns a ValidationError with Kind == ErrUnknownOperation.
func TestValidate_UnknownOperation_Empty(t *testing.T) {
	v := validator.NewCommandValidator()

	_, err := v.Validate(validator.Command{Operation: "", Args: dummyArgs})
	if err == nil {
		t.Fatal("Validate(\"\") expected an error, got nil")
	}

	ve, ok := err.(*validator.ValidationError)
	if !ok {
		t.Fatalf("Validate(\"\") returned error of type %T, want *validator.ValidationError", err)
	}
	if ve.Kind != validator.ErrUnknownOperation {
		t.Errorf("ValidationError.Kind = %d, want ErrUnknownOperation (%d)", int(ve.Kind), int(validator.ErrUnknownOperation))
	}
}

// TestValidate_UnknownOperation_ErrorMessageContainsOperation verifies the error
// message mentions the unrecognized operation name.
func TestValidate_UnknownOperation_ErrorMessageContainsOperation(t *testing.T) {
	v := validator.NewCommandValidator()

	_, err := v.Validate(validator.Command{Operation: "modulo", Args: dummyArgs})
	if err == nil {
		t.Fatal("Validate(\"modulo\") expected an error, got nil")
	}

	ve, ok := err.(*validator.ValidationError)
	if !ok {
		t.Fatalf("Validate(\"modulo\") returned error of type %T, want *validator.ValidationError", err)
	}
	if ve.Kind != validator.ErrUnknownOperation {
		t.Errorf("ValidationError.Kind = %d, want ErrUnknownOperation (%d)", int(ve.Kind), int(validator.ErrUnknownOperation))
	}
	if ve.Message == "" {
		t.Error("ValidationError.Message is empty; expected a descriptive message")
	}
}

// TestValidate_UnknownOperation_ErrorMessageListsSupportedOps verifies the error
// message lists the supported operations.
func TestValidate_UnknownOperation_ErrorMessageListsSupportedOps(t *testing.T) {
	v := validator.NewCommandValidator()

	_, err := v.Validate(validator.Command{Operation: "pow", Args: dummyArgs})
	if err == nil {
		t.Fatal("Validate(\"pow\") expected an error, got nil")
	}

	ve, ok := err.(*validator.ValidationError)
	if !ok {
		t.Fatalf("Validate(\"pow\") returned error of type %T, want *validator.ValidationError", err)
	}

	msg := ve.Message
	for _, supported := range []string{"add", "subtract", "multiply", "divide"} {
		if !contains(msg, supported) {
			t.Errorf("ValidationError.Message %q does not mention supported operation %q", msg, supported)
		}
	}
}

// TestValidate_ValidOperation_ReturnsValidatedCommand verifies that a valid operation
// returns a non-nil ValidatedCommand.
func TestValidate_ValidOperation_ReturnsValidatedCommand(t *testing.T) {
	v := validator.NewCommandValidator()

	result, err := v.Validate(validator.Command{Operation: "add", Args: dummyArgs})
	if err != nil {
		t.Fatalf("Validate(\"add\") unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("Validate(\"add\") returned nil ValidatedCommand, want non-nil")
	}
}

// TestValidate_ValidOperation_OperationPreserved verifies the operation name is
// preserved in the returned ValidatedCommand.
func TestValidate_ValidOperation_OperationPreserved(t *testing.T) {
	v := validator.NewCommandValidator()

	result, err := v.Validate(validator.Command{Operation: "multiply", Args: dummyArgs})
	if err != nil {
		t.Fatalf("Validate(\"multiply\") unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("Validate(\"multiply\") returned nil ValidatedCommand, want non-nil")
	}
	if result.Operation != "multiply" {
		t.Errorf("ValidatedCommand.Operation = %q, want %q", result.Operation, "multiply")
	}
}

// TestValidate_CaseSensitivity verifies that mixed-case variants are not accepted.
func TestValidate_CaseSensitivity(t *testing.T) {
	caseVariants := []string{"Add", "ADD", "aDd", "Subtract", "SUBTRACT", "Multiply", "MULTIPLY", "Divide", "DIVIDE"}

	v := validator.NewCommandValidator()

	for _, op := range caseVariants {
		t.Run(op, func(t *testing.T) {
			_, err := v.Validate(validator.Command{Operation: op, Args: dummyArgs})
			if err == nil {
				t.Errorf("Validate(%q) expected error for case-variant, got nil", op)
				return
			}
			ve, ok := err.(*validator.ValidationError)
			if !ok {
				t.Errorf("Validate(%q) returned error of type %T, want *validator.ValidationError", op, err)
				return
			}
			if ve.Kind != validator.ErrUnknownOperation {
				t.Errorf("Validate(%q) ValidationError.Kind = %d, want ErrUnknownOperation (%d)",
					op, int(ve.Kind), int(validator.ErrUnknownOperation))
			}
		})
	}
}

// contains is a helper to check substring presence.
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		func() bool {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
			return false
		}())
}
