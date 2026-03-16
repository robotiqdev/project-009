package arithmetic

import "fmt"

// Result holds the computed value and its formatted string representation.
type Result struct {
	Value     float64
	Formatted string
}

// Engine defines the interface for executing arithmetic operations.
type Engine interface {
	Execute(operation string, args []float64) (Result, error)
}

// Calculator is the stub implementation of Engine (to be implemented in TASK-4467).
type Calculator struct{}

// NewCalculator returns a new Calculator instance.
func NewCalculator() *Calculator {
	return &Calculator{}
}

func add(a, b float64) Result {
	v := a + b
	return Result{Value: v, Formatted: fmt.Sprintf("%.2f", v)}
}

func subtract(a, b float64) Result {
	v := a - b
	return Result{Value: v, Formatted: fmt.Sprintf("%.2f", v)}
}

func multiply(a, b float64) Result {
	v := a * b
	return Result{Value: v, Formatted: fmt.Sprintf("%.2f", v)}
}

// divide computes a/b. Division by zero is not passed here (validator catches it first),
// but if it were, Go returns +Inf for float64 division by zero, not a panic.
func divide(a, b float64) Result {
	v := a / b
	return Result{Value: v, Formatted: fmt.Sprintf("%.2f", v)}
}

// Execute dispatches to the appropriate arithmetic function based on operation.
// Returns an error for unknown operations as a defensive guard.
func (c *Calculator) Execute(operation string, args []float64) (Result, error) {
	switch operation {
	case "add":
		return add(args[0], args[1]), nil
	case "subtract":
		return subtract(args[0], args[1]), nil
	case "multiply":
		return multiply(args[0], args[1]), nil
	case "divide":
		return divide(args[0], args[1]), nil
	default:
		return Result{}, fmt.Errorf("unknown operation: %s", operation)
	}
}
