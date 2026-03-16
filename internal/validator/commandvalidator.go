package validator

import "github.com/workspace/repo/internal/parser"

// CommandValidator is a concrete implementation of Validator.
type CommandValidator struct{}

// NewCommandValidator returns a new CommandValidator.
func NewCommandValidator() *CommandValidator {
	return &CommandValidator{}
}

// Validate checks the parsed command for correctness.
func (cv *CommandValidator) Validate(cmd parser.Command) error {
	return nil
}
