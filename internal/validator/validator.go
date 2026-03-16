package validator

import "github.com/workspace/repo/internal/parser"

// Validator defines the interface for validating parsed commands.
type Validator interface {
	Validate(cmd parser.Command) error
}
