package arithmetic

import "errors"

var (
	ErrDivisionByZero   = errors.New("division by zero")
	ErrUnknownOperation = errors.New("unknown operation")
	ErrInvalidArgCount  = errors.New("invalid argument count")
)
