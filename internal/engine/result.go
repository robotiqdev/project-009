package engine

import "errors"

// Result wraps a computed float64 value.
type Result struct {
	Value float64
}

// Sentinel errors returned by engine operations.
var (
	ErrDivisionByZero  = errors.New("division by zero")
	ErrUnknownOperation = errors.New("unknown operation")
)
