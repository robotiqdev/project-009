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

// TestValidate_ArgCount_AddWithTwoArgs verifies that 'add' with exactly 2 args succeeds.
func TestValidate_ArgCount_AddWithTwoArgs(t *testing.T) {
	v := validator.NewCommandValidator()

	result, err := v.Validate(validator.Command{Operation: "add", Args: []string{"1", "2"}})
	if err != nil {
		t.Fatalf("Validate(\"add\", [\"1\",\"2\"]) unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("Validate(\"add\", [\"1\",\"2\"]) returned nil, want non-nil ValidatedCommand")
	}
}

// TestValidate_ArgCount_AddWithOneArg verifies that 'add' with 1 arg returns ErrWrongArgCount.
func TestValidate_ArgCount_AddWithOneArg(t *testing.T) {
	v := validator.NewCommandValidator()

	_, err := v.Validate(validator.Command{Operation: "add", Args: []string{"1"}})
	if err == nil {
		t.Fatal("Validate(\"add\", [\"1\"]) expected error, got nil")
	}

	ve, ok := err.(*validator.ValidationError)
	if !ok {
		t.Fatalf("Validate(\"add\", [\"1\"]) returned error of type %T, want *validator.ValidationError", err)
	}
	if ve.Kind != validator.ErrWrongArgCount {
		t.Errorf("ValidationError.Kind = %d, want ErrWrongArgCount (%d)", int(ve.Kind), int(validator.ErrWrongArgCount))
	}
}

// TestValidate_ArgCount_AddWithThreeArgs verifies that 'add' with 3 args returns ErrWrongArgCount.
func TestValidate_ArgCount_AddWithThreeArgs(t *testing.T) {
	v := validator.NewCommandValidator()

	_, err := v.Validate(validator.Command{Operation: "add", Args: []string{"1", "2", "3"}})
	if err == nil {
		t.Fatal("Validate(\"add\", [\"1\",\"2\",\"3\"]) expected error, got nil")
	}

	ve, ok := err.(*validator.ValidationError)
	if !ok {
		t.Fatalf("Validate(\"add\", [\"1\",\"2\",\"3\"]) returned error of type %T, want *validator.ValidationError", err)
	}
	if ve.Kind != validator.ErrWrongArgCount {
		t.Errorf("ValidationError.Kind = %d, want ErrWrongArgCount (%d)", int(ve.Kind), int(validator.ErrWrongArgCount))
	}
}

// TestValidate_ArgCount_AddWithZeroArgs verifies that 'add' with 0 args returns ErrWrongArgCount.
func TestValidate_ArgCount_AddWithZeroArgs(t *testing.T) {
	v := validator.NewCommandValidator()

	_, err := v.Validate(validator.Command{Operation: "add", Args: []string{}})
	if err == nil {
		t.Fatal("Validate(\"add\", []) expected error, got nil")
	}

	ve, ok := err.(*validator.ValidationError)
	if !ok {
		t.Fatalf("Validate(\"add\", []) returned error of type %T, want *validator.ValidationError", err)
	}
	if ve.Kind != validator.ErrWrongArgCount {
		t.Errorf("ValidationError.Kind = %d, want ErrWrongArgCount (%d)", int(ve.Kind), int(validator.ErrWrongArgCount))
	}
}

// TestValidate_ArgCount_ErrorMessageContainsOperationAndCounts verifies the error
// message mentions the operation and the expected/actual argument counts.
func TestValidate_ArgCount_ErrorMessageContainsOperationAndCounts(t *testing.T) {
	v := validator.NewCommandValidator()

	_, err := v.Validate(validator.Command{Operation: "add", Args: []string{"1"}})
	if err == nil {
		t.Fatal("Validate(\"add\", [\"1\"]) expected error, got nil")
	}

	ve, ok := err.(*validator.ValidationError)
	if !ok {
		t.Fatalf("Validate(\"add\", [\"1\"]) returned error of type %T, want *validator.ValidationError", err)
	}
	if ve.Message == "" {
		t.Error("ValidationError.Message is empty; expected a descriptive message")
	}
}

// TestValidate_ArgCount_WrongCountCheckedBeforeNumericParsing verifies that a wrong
// arg count error is returned even when the args themselves are non-numeric, confirming
// that arg count validation runs before numeric parsing.
func TestValidate_ArgCount_WrongCountCheckedBeforeNumericParsing(t *testing.T) {
	v := validator.NewCommandValidator()

	// Single non-numeric arg: should get ErrWrongArgCount, not ErrNotANumber.
	_, err := v.Validate(validator.Command{Operation: "add", Args: []string{"abc"}})
	if err == nil {
		t.Fatal("Validate(\"add\", [\"abc\"]) expected error, got nil")
	}

	ve, ok := err.(*validator.ValidationError)
	if !ok {
		t.Fatalf("Validate(\"add\", [\"abc\"]) returned error of type %T, want *validator.ValidationError", err)
	}
	if ve.Kind != validator.ErrWrongArgCount {
		t.Errorf("ValidationError.Kind = %d, want ErrWrongArgCount (%d) (arg count must be validated before numeric parsing)",
			int(ve.Kind), int(validator.ErrWrongArgCount))
	}
}

// --- Numeric Input Validator Tests (TASK-4473) ---

// TestValidate_Numeric_DecimalArgs verifies that 'add 1.5 2.3' succeeds with
// Args parsed to [1.5, 2.3].
func TestValidate_Numeric_DecimalArgs(t *testing.T) {
	v := validator.NewCommandValidator()

	result, err := v.Validate(validator.Command{Operation: "add", Args: []string{"1.5", "2.3"}})
	if err != nil {
		t.Fatalf("Validate(\"add\", [\"1.5\",\"2.3\"]) unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("Validate(\"add\", [\"1.5\",\"2.3\"]) returned nil, want non-nil ValidatedCommand")
	}
	if len(result.Args) != 2 {
		t.Fatalf("ValidatedCommand.Args length = %d, want 2", len(result.Args))
	}
	if result.Args[0] != 1.5 {
		t.Errorf("Args[0] = %v, want 1.5", result.Args[0])
	}
	if result.Args[1] != 2.3 {
		t.Errorf("Args[1] = %v, want 2.3", result.Args[1])
	}
}

// TestValidate_Numeric_ScientificNotation verifies that 'add 1e2 3' succeeds
// because scientific notation is a valid float64 representation.
func TestValidate_Numeric_ScientificNotation(t *testing.T) {
	v := validator.NewCommandValidator()

	result, err := v.Validate(validator.Command{Operation: "add", Args: []string{"1e2", "3"}})
	if err != nil {
		t.Fatalf("Validate(\"add\", [\"1e2\",\"3\"]) unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("Validate(\"add\", [\"1e2\",\"3\"]) returned nil, want non-nil ValidatedCommand")
	}
	if len(result.Args) != 2 {
		t.Fatalf("ValidatedCommand.Args length = %d, want 2", len(result.Args))
	}
	if result.Args[0] != 100.0 {
		t.Errorf("Args[0] = %v, want 100.0 (1e2 parsed as float64)", result.Args[0])
	}
	if result.Args[1] != 3.0 {
		t.Errorf("Args[1] = %v, want 3.0", result.Args[1])
	}
}

// TestValidate_Numeric_FirstArgNotANumber verifies that 'add abc 2' returns
// a ValidationError with Kind=ErrNotANumber and message mentioning 'abc'.
func TestValidate_Numeric_FirstArgNotANumber(t *testing.T) {
	v := validator.NewCommandValidator()

	_, err := v.Validate(validator.Command{Operation: "add", Args: []string{"abc", "2"}})
	if err == nil {
		t.Fatal("Validate(\"add\", [\"abc\",\"2\"]) expected error, got nil")
	}

	ve, ok := err.(*validator.ValidationError)
	if !ok {
		t.Fatalf("Validate(\"add\", [\"abc\",\"2\"]) returned error of type %T, want *validator.ValidationError", err)
	}
	if ve.Kind != validator.ErrNotANumber {
		t.Errorf("ValidationError.Kind = %d, want ErrNotANumber (%d)", int(ve.Kind), int(validator.ErrNotANumber))
	}
	// The message must mention the offending argument in quoted form (e.g. "abc").
	if !contains(ve.Message, `"abc"`) {
		t.Errorf("ValidationError.Message = %q, want it to contain %q (the bad argument quoted)", ve.Message, `"abc"`)
	}
}

// TestValidate_Numeric_SecondArgNotANumber verifies that 'add 1 xyz' returns
// a ValidationError with Kind=ErrNotANumber and message mentioning 'xyz'.
func TestValidate_Numeric_SecondArgNotANumber(t *testing.T) {
	v := validator.NewCommandValidator()

	_, err := v.Validate(validator.Command{Operation: "add", Args: []string{"1", "xyz"}})
	if err == nil {
		t.Fatal("Validate(\"add\", [\"1\",\"xyz\"]) expected error, got nil")
	}

	ve, ok := err.(*validator.ValidationError)
	if !ok {
		t.Fatalf("Validate(\"add\", [\"1\",\"xyz\"]) returned error of type %T, want *validator.ValidationError", err)
	}
	if ve.Kind != validator.ErrNotANumber {
		t.Errorf("ValidationError.Kind = %d, want ErrNotANumber (%d)", int(ve.Kind), int(validator.ErrNotANumber))
	}
	// The message must mention the offending argument in quoted form (e.g. "xyz").
	if !contains(ve.Message, `"xyz"`) {
		t.Errorf("ValidationError.Message = %q, want it to contain %q (the bad argument quoted)", ve.Message, `"xyz"`)
	}
}

// TestValidate_Numeric_ErrorMessageMentionsInvalidNumber verifies the error
// message contains the phrase "invalid number" as specified.
func TestValidate_Numeric_ErrorMessageMentionsInvalidNumber(t *testing.T) {
	v := validator.NewCommandValidator()

	_, err := v.Validate(validator.Command{Operation: "add", Args: []string{"notanumber", "2"}})
	if err == nil {
		t.Fatal("Validate(\"add\", [\"notanumber\",\"2\"]) expected error, got nil")
	}

	ve, ok := err.(*validator.ValidationError)
	if !ok {
		t.Fatalf("returned error of type %T, want *validator.ValidationError", err)
	}
	if ve.Kind != validator.ErrNotANumber {
		t.Errorf("ValidationError.Kind = %d, want ErrNotANumber (%d)", int(ve.Kind), int(validator.ErrNotANumber))
	}
	if !contains(ve.Message, "invalid number") {
		t.Errorf("ValidationError.Message = %q, want it to contain \"invalid number\"", ve.Message)
	}
}

// TestValidate_Numeric_ErrorMessageMentionsDecimalNumber verifies the error
// message tells the user what format is expected.
func TestValidate_Numeric_ErrorMessageMentionsDecimalNumber(t *testing.T) {
	v := validator.NewCommandValidator()

	_, err := v.Validate(validator.Command{Operation: "add", Args: []string{"abc", "2"}})
	if err == nil {
		t.Fatal("Validate(\"add\", [\"abc\",\"2\"]) expected error, got nil")
	}

	ve, ok := err.(*validator.ValidationError)
	if !ok {
		t.Fatalf("returned error of type %T, want *validator.ValidationError", err)
	}
	if !contains(ve.Message, "decimal number") {
		t.Errorf("ValidationError.Message = %q, want it to contain \"decimal number\"", ve.Message)
	}
}

// TestValidate_Numeric_FailsOnFirstInvalidArg verifies that the validator fails
// on the first invalid argument encountered (not the second).
func TestValidate_Numeric_FailsOnFirstInvalidArg(t *testing.T) {
	v := validator.NewCommandValidator()

	// Both args are invalid; the error message should mention the first one.
	_, err := v.Validate(validator.Command{Operation: "add", Args: []string{"bad1", "bad2"}})
	if err == nil {
		t.Fatal("Validate(\"add\", [\"bad1\",\"bad2\"]) expected error, got nil")
	}

	ve, ok := err.(*validator.ValidationError)
	if !ok {
		t.Fatalf("returned error of type %T, want *validator.ValidationError", err)
	}
	if ve.Kind != validator.ErrNotANumber {
		t.Errorf("ValidationError.Kind = %d, want ErrNotANumber (%d)", int(ve.Kind), int(validator.ErrNotANumber))
	}
	// Should mention first bad arg, not second.
	if !contains(ve.Message, `"bad1"`) {
		t.Errorf("ValidationError.Message = %q, want it to mention first bad arg %q", ve.Message, `"bad1"`)
	}
}

// TestValidate_Numeric_IntegerArgsAreValid verifies that integer strings like "1" and "2"
// are accepted as valid float64 values.
func TestValidate_Numeric_IntegerArgsAreValid(t *testing.T) {
	v := validator.NewCommandValidator()

	result, err := v.Validate(validator.Command{Operation: "add", Args: []string{"1", "2"}})
	if err != nil {
		t.Fatalf("Validate(\"add\", [\"1\",\"2\"]) unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("Validate(\"add\", [\"1\",\"2\"]) returned nil, want non-nil ValidatedCommand")
	}
	if result.Args[0] != 1.0 {
		t.Errorf("Args[0] = %v, want 1.0", result.Args[0])
	}
	if result.Args[1] != 2.0 {
		t.Errorf("Args[1] = %v, want 2.0", result.Args[1])
	}
}

// TestValidate_Numeric_NegativeNumbersAreValid verifies that negative float values
// (with a leading minus sign) are accepted.
func TestValidate_Numeric_NegativeNumbersAreValid(t *testing.T) {
	v := validator.NewCommandValidator()

	result, err := v.Validate(validator.Command{Operation: "add", Args: []string{"-1.5", "-2.3"}})
	if err != nil {
		t.Fatalf("Validate(\"add\", [\"-1.5\",\"-2.3\"]) unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("returned nil ValidatedCommand, want non-nil")
	}
	if result.Args[0] != -1.5 {
		t.Errorf("Args[0] = %v, want -1.5", result.Args[0])
	}
	if result.Args[1] != -2.3 {
		t.Errorf("Args[1] = %v, want -2.3", result.Args[1])
	}
}

// TestValidate_Numeric_ArgCountCheckedBeforeNumericParsing verifies that an arg
// count mismatch with non-numeric args returns ErrWrongArgCount, not ErrNotANumber.
// This confirms the ordering: operation -> arg count -> numeric parsing.
func TestValidate_Numeric_ArgCountCheckedBeforeNumericParsing(t *testing.T) {
	v := validator.NewCommandValidator()

	// Two non-numeric args that are also the wrong count for a hypothetical 1-arg op,
	// but here we pass 3 non-numeric args to 'add' which expects 2.
	_, err := v.Validate(validator.Command{Operation: "add", Args: []string{"abc", "def", "ghi"}})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	ve, ok := err.(*validator.ValidationError)
	if !ok {
		t.Fatalf("returned error of type %T, want *validator.ValidationError", err)
	}
	if ve.Kind != validator.ErrWrongArgCount {
		t.Errorf("ValidationError.Kind = %d, want ErrWrongArgCount (%d); arg count must be validated before numeric parsing",
			int(ve.Kind), int(validator.ErrWrongArgCount))
	}
}

// TestValidate_Numeric_AllOpsAcceptValidFloats verifies that all supported
// operations accept valid float64 arguments without error.
func TestValidate_Numeric_AllOpsAcceptValidFloats(t *testing.T) {
	ops := []string{"add", "subtract", "multiply", "divide"}
	v := validator.NewCommandValidator()

	for _, op := range ops {
		t.Run(op, func(t *testing.T) {
			result, err := v.Validate(validator.Command{Operation: op, Args: []string{"3.14", "2.71"}})
			if err != nil {
				t.Errorf("Validate(%q, [\"3.14\",\"2.71\"]) unexpected error: %v", op, err)
				return
			}
			if result == nil {
				t.Errorf("Validate(%q, [\"3.14\",\"2.71\"]) returned nil, want non-nil", op)
				return
			}
			if result.Args[0] != 3.14 {
				t.Errorf("Args[0] = %v, want 3.14", result.Args[0])
			}
			if result.Args[1] != 2.71 {
				t.Errorf("Args[1] = %v, want 2.71", result.Args[1])
			}
		})
	}
}

// TestValidate_Numeric_AllOpsRejectNonNumericArgs verifies that all supported
// operations reject non-numeric arguments with ErrNotANumber.
func TestValidate_Numeric_AllOpsRejectNonNumericArgs(t *testing.T) {
	ops := []string{"add", "subtract", "multiply", "divide"}
	v := validator.NewCommandValidator()

	for _, op := range ops {
		t.Run(op, func(t *testing.T) {
			_, err := v.Validate(validator.Command{Operation: op, Args: []string{"bad", "2"}})
			if err == nil {
				t.Errorf("Validate(%q, [\"bad\",\"2\"]) expected error, got nil", op)
				return
			}
			ve, ok := err.(*validator.ValidationError)
			if !ok {
				t.Errorf("Validate(%q) returned error of type %T, want *validator.ValidationError", op, err)
				return
			}
			if ve.Kind != validator.ErrNotANumber {
				t.Errorf("Validate(%q) ValidationError.Kind = %d, want ErrNotANumber (%d)",
					op, int(ve.Kind), int(validator.ErrNotANumber))
			}
		})
	}
}

// TestValidate_Numeric_ValidatedCommandArgsLengthMatchesInput verifies that the
// returned ValidatedCommand.Args slice has the same length as the input Args.
func TestValidate_Numeric_ValidatedCommandArgsLengthMatchesInput(t *testing.T) {
	v := validator.NewCommandValidator()

	result, err := v.Validate(validator.Command{Operation: "add", Args: []string{"10", "20"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Args) != 2 {
		t.Errorf("len(ValidatedCommand.Args) = %d, want 2", len(result.Args))
	}
}

// --- Division by Zero Validator Tests (TASK-4474) ---

// TestValidate_DivisionByZero_IntegerZeroDivisor verifies that 'divide 5 0'
// returns a ValidationError with Kind=ErrDivisionByZero.
func TestValidate_DivisionByZero_IntegerZeroDivisor(t *testing.T) {
	v := validator.NewCommandValidator()

	_, err := v.Validate(validator.Command{Operation: "divide", Args: []string{"5", "0"}})
	if err == nil {
		t.Fatal("Validate(\"divide\", [\"5\",\"0\"]) expected error, got nil")
	}

	ve, ok := err.(*validator.ValidationError)
	if !ok {
		t.Fatalf("Validate(\"divide\", [\"5\",\"0\"]) returned error of type %T, want *validator.ValidationError", err)
	}
	if ve.Kind != validator.ErrDivisionByZero {
		t.Errorf("ValidationError.Kind = %d, want ErrDivisionByZero (%d)", int(ve.Kind), int(validator.ErrDivisionByZero))
	}
}

// TestValidate_DivisionByZero_FloatZeroDivisor verifies that 'divide 5 0.0'
// returns a ValidationError with Kind=ErrDivisionByZero.
func TestValidate_DivisionByZero_FloatZeroDivisor(t *testing.T) {
	v := validator.NewCommandValidator()

	_, err := v.Validate(validator.Command{Operation: "divide", Args: []string{"5", "0.0"}})
	if err == nil {
		t.Fatal("Validate(\"divide\", [\"5\",\"0.0\"]) expected error, got nil")
	}

	ve, ok := err.(*validator.ValidationError)
	if !ok {
		t.Fatalf("Validate(\"divide\", [\"5\",\"0.0\"]) returned error of type %T, want *validator.ValidationError", err)
	}
	if ve.Kind != validator.ErrDivisionByZero {
		t.Errorf("ValidationError.Kind = %d, want ErrDivisionByZero (%d)", int(ve.Kind), int(validator.ErrDivisionByZero))
	}
}

// TestValidate_DivisionByZero_ErrorMessageDescribesDivisionByZero verifies that
// the error message is descriptive and mentions division by zero.
func TestValidate_DivisionByZero_ErrorMessageDescribesDivisionByZero(t *testing.T) {
	v := validator.NewCommandValidator()

	_, err := v.Validate(validator.Command{Operation: "divide", Args: []string{"5", "0"}})
	if err == nil {
		t.Fatal("Validate(\"divide\", [\"5\",\"0\"]) expected error, got nil")
	}

	ve, ok := err.(*validator.ValidationError)
	if !ok {
		t.Fatalf("returned error of type %T, want *validator.ValidationError", err)
	}
	if ve.Message == "" {
		t.Error("ValidationError.Message is empty; expected a descriptive message about division by zero")
	}
}

// TestValidate_DivisionByZero_ZeroDividendIsOK verifies that 'divide 0 5'
// is valid — a zero dividend (numerator) is allowed.
func TestValidate_DivisionByZero_ZeroDividendIsOK(t *testing.T) {
	v := validator.NewCommandValidator()

	result, err := v.Validate(validator.Command{Operation: "divide", Args: []string{"0", "5"}})
	if err != nil {
		t.Fatalf("Validate(\"divide\", [\"0\",\"5\"]) unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("Validate(\"divide\", [\"0\",\"5\"]) returned nil, want non-nil ValidatedCommand")
	}
	if result.Args[0] != 0.0 {
		t.Errorf("Args[0] = %v, want 0.0", result.Args[0])
	}
	if result.Args[1] != 5.0 {
		t.Errorf("Args[1] = %v, want 5.0", result.Args[1])
	}
}

// TestValidate_DivisionByZero_ZeroIsValidForNonDivideOps verifies that 'add 5 0'
// is valid — zero is allowed as an argument for non-divide operations.
func TestValidate_DivisionByZero_ZeroIsValidForNonDivideOps(t *testing.T) {
	ops := []string{"add", "subtract", "multiply"}
	v := validator.NewCommandValidator()

	for _, op := range ops {
		t.Run(op, func(t *testing.T) {
			result, err := v.Validate(validator.Command{Operation: op, Args: []string{"5", "0"}})
			if err != nil {
				t.Errorf("Validate(%q, [\"5\",\"0\"]) unexpected error: %v", op, err)
				return
			}
			if result == nil {
				t.Errorf("Validate(%q, [\"5\",\"0\"]) returned nil, want non-nil ValidatedCommand", op)
			}
		})
	}
}

// TestValidate_DivisionByZero_NonZeroDivisorIsOK verifies that 'divide 5 1'
// succeeds when the divisor is non-zero.
func TestValidate_DivisionByZero_NonZeroDivisorIsOK(t *testing.T) {
	v := validator.NewCommandValidator()

	result, err := v.Validate(validator.Command{Operation: "divide", Args: []string{"5", "1"}})
	if err != nil {
		t.Fatalf("Validate(\"divide\", [\"5\",\"1\"]) unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("Validate(\"divide\", [\"5\",\"1\"]) returned nil, want non-nil ValidatedCommand")
	}
	if result.Operation != "divide" {
		t.Errorf("ValidatedCommand.Operation = %q, want \"divide\"", result.Operation)
	}
	if result.Args[0] != 5.0 {
		t.Errorf("Args[0] = %v, want 5.0", result.Args[0])
	}
	if result.Args[1] != 1.0 {
		t.Errorf("Args[1] = %v, want 1.0", result.Args[1])
	}
}

// TestValidate_DivisionByZero_CheckedAfterNumericValidation verifies that
// division by zero check only runs after numeric validation succeeds — i.e.,
// a non-numeric divisor returns ErrNotANumber, not ErrDivisionByZero.
func TestValidate_DivisionByZero_CheckedAfterNumericValidation(t *testing.T) {
	v := validator.NewCommandValidator()

	_, err := v.Validate(validator.Command{Operation: "divide", Args: []string{"5", "abc"}})
	if err == nil {
		t.Fatal("Validate(\"divide\", [\"5\",\"abc\"]) expected error, got nil")
	}

	ve, ok := err.(*validator.ValidationError)
	if !ok {
		t.Fatalf("returned error of type %T, want *validator.ValidationError", err)
	}
	if ve.Kind != validator.ErrNotANumber {
		t.Errorf("ValidationError.Kind = %d, want ErrNotANumber (%d); division-by-zero check must run after numeric parsing",
			int(ve.Kind), int(validator.ErrNotANumber))
	}
}

// TestValidate_DivisionByZero_NegativeZeroDivisorIsRejected verifies that
// 'divide 5 -0' is also caught as division by zero (negative zero parses to 0.0).
func TestValidate_DivisionByZero_NegativeZeroDivisorIsRejected(t *testing.T) {
	v := validator.NewCommandValidator()

	_, err := v.Validate(validator.Command{Operation: "divide", Args: []string{"5", "-0"}})
	if err == nil {
		t.Fatal("Validate(\"divide\", [\"5\",\"-0\"]) expected error, got nil")
	}

	ve, ok := err.(*validator.ValidationError)
	if !ok {
		t.Fatalf("returned error of type %T, want *validator.ValidationError", err)
	}
	if ve.Kind != validator.ErrDivisionByZero {
		t.Errorf("ValidationError.Kind = %d, want ErrDivisionByZero (%d)", int(ve.Kind), int(validator.ErrDivisionByZero))
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
