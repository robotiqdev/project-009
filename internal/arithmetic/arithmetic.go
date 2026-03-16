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

// Execute is a stub that returns an error for all operations until implemented.
func (c *Calculator) Execute(operation string, args []float64) (Result, error) {
	return Result{}, fmt.Errorf("not implemented: %s", operation)
}
