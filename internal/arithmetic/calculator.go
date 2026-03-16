package arithmetic

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/workspace/repo/internal/parser"
)

// Calculator is a concrete implementation of Engine.
type Calculator struct{}

// NewCalculator returns a new Calculator.
func NewCalculator() *Calculator {
	return &Calculator{}
}

// Execute performs the arithmetic operation described by cmd.
func (c *Calculator) Execute(cmd parser.Command) (float64, error) {
	if len(cmd.Args) != 2 {
		return 0, fmt.Errorf("expected 2 arguments, got %d", len(cmd.Args))
	}
	a, err := strconv.ParseFloat(cmd.Args[0], 64)
	if err != nil {
		return 0, fmt.Errorf("invalid argument: %s", cmd.Args[0])
	}
	b, err := strconv.ParseFloat(cmd.Args[1], 64)
	if err != nil {
		return 0, fmt.Errorf("invalid argument: %s", cmd.Args[1])
	}
	switch cmd.Operation {
	case "add":
		return a + b, nil
	case "subtract", "sub":
		return a - b, nil
	case "multiply", "mul":
		return a * b, nil
	case "divide":
		if b == 0 {
			return 0, errors.New("division by zero")
		}
		return a / b, nil
	default:
		return 0, fmt.Errorf("unknown operation: %s", cmd.Operation)
	}
}
