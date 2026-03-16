package validator

import (
	"fmt"
	"strconv"
)

// ErrorKind represents the type of validation error.
type ErrorKind int

const (
	ErrUnknownOperation ErrorKind = iota
	ErrWrongArgCount
	ErrNotANumber
	ErrDivisionByZero
)

// ValidationError represents a validation error with a kind and message.
type ValidationError struct {
	Kind    ErrorKind
	Message string
}

// Error implements the error interface.
func (e *ValidationError) Error() string {
	return e.Message
}

// ValidatedCommand represents a successfully validated command with typed args.
type ValidatedCommand struct {
	Operation string
	Args      []float64
}

// Command represents a raw (unvalidated) command input.
type Command struct {
	Operation string
	Args      []string
}

var validOps = map[string]int{
	"add":      2,
	"subtract": 2,
	"multiply": 2,
	"divide":   2,
}

// CommandValidator validates Command inputs.
type CommandValidator struct{}

// NewCommandValidator returns a new CommandValidator.
func NewCommandValidator() *CommandValidator {
	return &CommandValidator{}
}

func (v *CommandValidator) validateArgCount(cmd Command, expected int) error {
	if len(cmd.Args) != expected {
		return &ValidationError{
			Kind:    ErrWrongArgCount,
			Message: fmt.Sprintf("%s requires exactly %d argument(s), got %d", cmd.Operation, expected, len(cmd.Args)),
		}
	}
	return nil
}

func (v *CommandValidator) validateOperation(cmd Command) error {
	if _, ok := validOps[cmd.Operation]; !ok {
		return &ValidationError{
			Kind:    ErrUnknownOperation,
			Message: fmt.Sprintf("unknown operation %q; supported: add, subtract, multiply, divide", cmd.Operation),
		}
	}
	return nil
}

// Validate validates a Command and returns a ValidatedCommand or a ValidationError.
func (v *CommandValidator) Validate(cmd Command) (*ValidatedCommand, error) {
	if err := v.validateOperation(cmd); err != nil {
		return nil, err
	}

	if err := v.validateArgCount(cmd, validOps[cmd.Operation]); err != nil {
		return nil, err
	}

	args := make([]float64, len(cmd.Args))
	for i, a := range cmd.Args {
		f, err := strconv.ParseFloat(a, 64)
		if err != nil {
			return nil, &ValidationError{
				Kind:    ErrNotANumber,
				Message: fmt.Sprintf("invalid number %q: must be a valid decimal number", a),
			}
		}
		args[i] = f
	}

	return &ValidatedCommand{
		Operation: cmd.Operation,
		Args:      args,
	}, nil
}
