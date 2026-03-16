package repl_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/workspace/repo/internal/repl"
	"github.com/workspace/repo/internal/validator"
)

// --- Stub implementations ---

// stubParser returns a fixed Command for any input.
type stubParser struct {
	cmd validator.Command
	err error
}

func (p *stubParser) Parse(input string) (validator.Command, error) {
	return p.cmd, p.err
}

// stubValidator is configured with a sequence of (result, error) pairs.
// Each call to Validate returns the next pair in the sequence.
type stubValidator struct {
	results []*validator.ValidatedCommand
	errors  []error
	calls   int
}

func (v *stubValidator) Validate(cmd validator.Command) (*validator.ValidatedCommand, error) {
	i := v.calls
	v.calls++
	if i < len(v.errors) {
		return v.results[i], v.errors[i]
	}
	// Default: success with zero-value command
	return &validator.ValidatedCommand{Operation: "add", Args: []float64{0, 0}}, nil
}

// stubEngine always returns the configured result.
type stubEngine struct {
	result repl.Result
	err    error
}

func (e *stubEngine) Execute(operation string, args []float64) (repl.Result, error) {
	return e.result, e.err
}

// --- Helper ---

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

// newValidatedCmd creates a ValidatedCommand for use in stubs.
func newValidatedCmd(op string, args []float64) *validator.ValidatedCommand {
	return &validator.ValidatedCommand{Operation: op, Args: args}
}

// --- Tests ---

// TestRun_ValidatorReturnsErrUnknownOperation verifies that when the validator
// returns an ErrUnknownOperation error, the REPL outputs "Error: " prefix and
// does not exit the loop.
func TestRun_ValidatorReturnsErrUnknownOperation(t *testing.T) {
	validationErr := &validator.ValidationError{
		Kind:    validator.ErrUnknownOperation,
		Message: `unknown operation "foo"; supported: add, subtract, multiply, divide`,
	}

	p := &stubParser{cmd: validator.Command{Operation: "foo", Args: []string{"1", "2"}}}
	v := &stubValidator{
		results: []*validator.ValidatedCommand{nil},
		errors:  []error{validationErr},
	}
	e := &stubEngine{result: repl.Result{Formatted: "3"}}

	input := strings.NewReader("foo 1 2\n")
	var out bytes.Buffer

	r := repl.New(p, v, e)
	r.Run(input, &out)

	got := out.String()
	if !contains(got, "Error: ") {
		t.Errorf("expected output to contain 'Error: ', got: %q", got)
	}
	if !contains(got, validationErr.Message) {
		t.Errorf("expected output to contain error message %q, got: %q", validationErr.Message, got)
	}
}

// TestRun_ValidatorReturnsErrWrongArgCount verifies that when the validator
// returns an ErrWrongArgCount error, the REPL outputs "Error: " prefix and
// does not exit the loop.
func TestRun_ValidatorReturnsErrWrongArgCount(t *testing.T) {
	validationErr := &validator.ValidationError{
		Kind:    validator.ErrWrongArgCount,
		Message: "add requires exactly 2 argument(s), got 1",
	}

	p := &stubParser{cmd: validator.Command{Operation: "add", Args: []string{"1"}}}
	v := &stubValidator{
		results: []*validator.ValidatedCommand{nil},
		errors:  []error{validationErr},
	}
	e := &stubEngine{result: repl.Result{Formatted: "3"}}

	input := strings.NewReader("add 1\n")
	var out bytes.Buffer

	r := repl.New(p, v, e)
	r.Run(input, &out)

	got := out.String()
	if !contains(got, "Error: ") {
		t.Errorf("expected output to contain 'Error: ', got: %q", got)
	}
	if !contains(got, validationErr.Message) {
		t.Errorf("expected output to contain error message %q, got: %q", validationErr.Message, got)
	}
}

// TestRun_ValidatorReturnsErrNotANumber verifies that when the validator
// returns an ErrNotANumber error, the REPL outputs "Error: " prefix and
// does not exit the loop.
func TestRun_ValidatorReturnsErrNotANumber(t *testing.T) {
	validationErr := &validator.ValidationError{
		Kind:    validator.ErrNotANumber,
		Message: `invalid number "abc": must be a valid decimal number`,
	}

	p := &stubParser{cmd: validator.Command{Operation: "add", Args: []string{"abc", "2"}}}
	v := &stubValidator{
		results: []*validator.ValidatedCommand{nil},
		errors:  []error{validationErr},
	}
	e := &stubEngine{result: repl.Result{Formatted: "3"}}

	input := strings.NewReader("add abc 2\n")
	var out bytes.Buffer

	r := repl.New(p, v, e)
	r.Run(input, &out)

	got := out.String()
	if !contains(got, "Error: ") {
		t.Errorf("expected output to contain 'Error: ', got: %q", got)
	}
	if !contains(got, validationErr.Message) {
		t.Errorf("expected output to contain error message %q, got: %q", validationErr.Message, got)
	}
}

// TestRun_ValidatorReturnsErrDivisionByZero verifies that when the validator
// returns an ErrDivisionByZero error, the REPL outputs "Error: " prefix and
// does not exit the loop.
func TestRun_ValidatorReturnsErrDivisionByZero(t *testing.T) {
	validationErr := &validator.ValidationError{
		Kind:    validator.ErrDivisionByZero,
		Message: "division by zero is not allowed",
	}

	p := &stubParser{cmd: validator.Command{Operation: "divide", Args: []string{"5", "0"}}}
	v := &stubValidator{
		results: []*validator.ValidatedCommand{nil},
		errors:  []error{validationErr},
	}
	e := &stubEngine{result: repl.Result{Formatted: "3"}}

	input := strings.NewReader("divide 5 0\n")
	var out bytes.Buffer

	r := repl.New(p, v, e)
	r.Run(input, &out)

	got := out.String()
	if !contains(got, "Error: ") {
		t.Errorf("expected output to contain 'Error: ', got: %q", got)
	}
	if !contains(got, validationErr.Message) {
		t.Errorf("expected output to contain error message %q, got: %q", validationErr.Message, got)
	}
}

// TestRun_LoopDoesNotExitOnValidationError verifies that after a validation
// error the REPL continues looping and processes subsequent inputs.
// We send two inputs: first causes an error, second is processed successfully.
// The function must return (not hang) after EOF, proving both lines were read.
func TestRun_LoopDoesNotExitOnValidationError(t *testing.T) {
	validationErr := &validator.ValidationError{
		Kind:    validator.ErrUnknownOperation,
		Message: `unknown operation "bad"; supported: add, subtract, multiply, divide`,
	}
	successResult := newValidatedCmd("add", []float64{1, 2})

	p := &stubParser{cmd: validator.Command{Operation: "add", Args: []string{"1", "2"}}}
	v := &stubValidator{
		results: []*validator.ValidatedCommand{nil, successResult},
		errors:  []error{validationErr, nil},
	}
	e := &stubEngine{result: repl.Result{Formatted: "3"}}

	// Two input lines — error line followed by valid line
	input := strings.NewReader("bad\nadd 1 2\n")
	var out bytes.Buffer

	r := repl.New(p, v, e)
	r.Run(input, &out)

	got := out.String()

	// Error output must appear (loop didn't exit before processing first line)
	if !contains(got, "Error: ") {
		t.Errorf("expected output to contain 'Error: ' from first line, got: %q", got)
	}

	// Success output must also appear (loop continued after error)
	if !contains(got, "3") {
		t.Errorf("expected output to contain result '3' from second line, got: %q", got)
	}
}

// TestRun_ErrorThenSuccess_BothHandledCorrectly verifies a sequence where
// the first command produces a validation error and the second succeeds.
// Both the error output and the success output must be present.
func TestRun_ErrorThenSuccess_BothHandledCorrectly(t *testing.T) {
	validationErr := &validator.ValidationError{
		Kind:    validator.ErrDivisionByZero,
		Message: "division by zero is not allowed",
	}
	successResult := newValidatedCmd("add", []float64{10, 5})

	p := &stubParser{cmd: validator.Command{Operation: "add", Args: []string{"10", "5"}}}
	v := &stubValidator{
		results: []*validator.ValidatedCommand{nil, successResult},
		errors:  []error{validationErr, nil},
	}
	e := &stubEngine{result: repl.Result{Formatted: "15"}}

	input := strings.NewReader("divide 5 0\nadd 10 5\n")
	var out bytes.Buffer

	r := repl.New(p, v, e)
	r.Run(input, &out)

	got := out.String()

	// First command: error output
	if !contains(got, "Error: ") {
		t.Errorf("expected output to contain 'Error: ' for first command, got: %q", got)
	}
	if !contains(got, "division by zero is not allowed") {
		t.Errorf("expected output to contain division-by-zero message, got: %q", got)
	}

	// Second command: success output
	if !contains(got, "15") {
		t.Errorf("expected output to contain result '15' for second command, got: %q", got)
	}

	// The error message must appear before the result
	errIdx := strings.Index(got, "Error: ")
	resultIdx := strings.Index(got, "15")
	if errIdx >= resultIdx {
		t.Errorf("expected 'Error: ' to appear before '15' in output, got: %q", got)
	}
}

// TestRun_ErrorOutput_HasCorrectFormat verifies that the error output uses
// exactly the format "Error: <message>\n" (with a newline).
func TestRun_ErrorOutput_HasCorrectFormat(t *testing.T) {
	validationErr := &validator.ValidationError{
		Kind:    validator.ErrNotANumber,
		Message: `invalid number "xyz": must be a valid decimal number`,
	}

	p := &stubParser{cmd: validator.Command{Operation: "add", Args: []string{"xyz", "2"}}}
	v := &stubValidator{
		results: []*validator.ValidatedCommand{nil},
		errors:  []error{validationErr},
	}
	e := &stubEngine{}

	input := strings.NewReader("add xyz 2\n")
	var out bytes.Buffer

	r := repl.New(p, v, e)
	r.Run(input, &out)

	got := out.String()
	expected := "Error: " + validationErr.Message + "\n"
	if !contains(got, expected) {
		t.Errorf("expected output to contain %q, got: %q", expected, got)
	}
}

// TestRun_SuccessOutput_PrintsFormattedResult verifies that on a successful
// command, the REPL prints the formatted result followed by a newline.
func TestRun_SuccessOutput_PrintsFormattedResult(t *testing.T) {
	successResult := newValidatedCmd("multiply", []float64{3, 4})

	p := &stubParser{cmd: validator.Command{Operation: "multiply", Args: []string{"3", "4"}}}
	v := &stubValidator{
		results: []*validator.ValidatedCommand{successResult},
		errors:  []error{nil},
	}
	e := &stubEngine{result: repl.Result{Formatted: "12"}}

	input := strings.NewReader("multiply 3 4\n")
	var out bytes.Buffer

	r := repl.New(p, v, e)
	r.Run(input, &out)

	got := out.String()
	if !contains(got, "12") {
		t.Errorf("expected output to contain '12', got: %q", got)
	}
	if !contains(got, "12\n") {
		t.Errorf("expected result to be followed by newline, got: %q", got)
	}
}

// TestRun_MultipleValidationErrors_AllReported verifies that multiple
// validation errors across multiple inputs are all reported and the loop
// continues processing each line.
func TestRun_MultipleValidationErrors_AllReported(t *testing.T) {
	err1 := &validator.ValidationError{
		Kind:    validator.ErrUnknownOperation,
		Message: `unknown operation "foo"; supported: add, subtract, multiply, divide`,
	}
	err2 := &validator.ValidationError{
		Kind:    validator.ErrWrongArgCount,
		Message: "add requires exactly 2 argument(s), got 3",
	}
	err3 := &validator.ValidationError{
		Kind:    validator.ErrNotANumber,
		Message: `invalid number "x": must be a valid decimal number`,
	}

	p := &stubParser{cmd: validator.Command{Operation: "add", Args: []string{"1", "2"}}}
	v := &stubValidator{
		results: []*validator.ValidatedCommand{nil, nil, nil},
		errors:  []error{err1, err2, err3},
	}
	e := &stubEngine{}

	input := strings.NewReader("foo\nadd 1 2 3\nadd x 2\n")
	var out bytes.Buffer

	r := repl.New(p, v, e)
	r.Run(input, &out)

	got := out.String()

	for _, errMsg := range []string{err1.Message, err2.Message, err3.Message} {
		if !contains(got, "Error: "+errMsg) {
			t.Errorf("expected output to contain 'Error: %s', got: %q", errMsg, got)
		}
	}
}
