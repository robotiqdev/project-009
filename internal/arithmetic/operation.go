package arithmetic

// Operation represents an arithmetic operation.
type Operation string

const (
	OpAdd      Operation = "add"
	OpSubtract Operation = "subtract"
	OpMultiply Operation = "multiply"
	OpDivide   Operation = "divide"
)

// String returns the string representation of the Operation.
func (o Operation) String() string {
	return string(o)
}
