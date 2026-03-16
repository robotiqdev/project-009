package validator

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

// CommandValidator validates Command inputs.
type CommandValidator struct{}

// NewCommandValidator returns a new CommandValidator.
func NewCommandValidator() *CommandValidator {
	return nil
}

// Validate validates a Command and returns a ValidatedCommand or a ValidationError.
func (v *CommandValidator) Validate(cmd Command) (*ValidatedCommand, error) {
	return nil, nil
}
