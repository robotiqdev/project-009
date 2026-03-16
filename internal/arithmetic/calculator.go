package arithmetic

import "github.com/workspace/repo/internal/parser"

// Calculator is a concrete implementation of Engine.
type Calculator struct{}

// NewCalculator returns a new Calculator.
func NewCalculator() *Calculator {
	return &Calculator{}
}

// Execute performs the arithmetic operation described by cmd.
func (c *Calculator) Execute(cmd parser.Command) (float64, error) {
	return 0, nil
}
