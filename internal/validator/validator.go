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
	return ""
}

// ValidatedCommand represents a successfully validated command with typed args.
type ValidatedCommand struct {
	Operation string
	Args      []float64
}
