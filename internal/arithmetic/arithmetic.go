package arithmetic

// Result holds the computed value and its formatted string representation.
type Result struct {
	Value     float64
	Formatted string
}

// Engine defines the interface for executing arithmetic operations.
type Engine interface {
	Execute(operation string, args []float64) (Result, error)
}
