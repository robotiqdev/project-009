package arithmetic

import "testing"

func TestOperationConstants_Exist(t *testing.T) {
	tests := []struct {
		name     string
		op       Operation
		expected string
	}{
		{
			name:     "OpAdd has value 'add'",
			op:       OpAdd,
			expected: "add",
		},
		{
			name:     "OpSubtract has value 'subtract'",
			op:       OpSubtract,
			expected: "subtract",
		},
		{
			name:     "OpMultiply has value 'multiply'",
			op:       OpMultiply,
			expected: "multiply",
		},
		{
			name:     "OpDivide has value 'divide'",
			op:       OpDivide,
			expected: "divide",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.op) != tt.expected {
				t.Errorf("Operation constant %q = %q, want %q", tt.name, string(tt.op), tt.expected)
			}
		})
	}
}

func TestOperationConstants_AreDistinct(t *testing.T) {
	ops := []Operation{OpAdd, OpSubtract, OpMultiply, OpDivide}
	seen := make(map[Operation]bool)

	for _, op := range ops {
		if seen[op] {
			t.Errorf("Operation constant %q is not distinct — duplicate value found", op)
		}
		seen[op] = true
	}

	if len(seen) != 4 {
		t.Errorf("expected 4 distinct Operation constants, got %d", len(seen))
	}
}

func TestOperationString_Method(t *testing.T) {
	tests := []struct {
		name     string
		op       Operation
		expected string
	}{
		{
			name:     "OpAdd String() returns 'add'",
			op:       OpAdd,
			expected: "add",
		},
		{
			name:     "OpSubtract String() returns 'subtract'",
			op:       OpSubtract,
			expected: "subtract",
		},
		{
			name:     "OpMultiply String() returns 'multiply'",
			op:       OpMultiply,
			expected: "multiply",
		},
		{
			name:     "OpDivide String() returns 'divide'",
			op:       OpDivide,
			expected: "divide",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.op.String() != tt.expected {
				t.Errorf("Operation.String() = %q, want %q", tt.op.String(), tt.expected)
			}
		})
	}
}

func TestOperation_DirectComparison(t *testing.T) {
	// Verify that Operation values can be compared directly against string tokens
	// without a lookup map, as required by the architecture
	tests := []struct {
		name  string
		token string
		op    Operation
		match bool
	}{
		{
			name:  "token 'add' matches OpAdd",
			token: "add",
			op:    OpAdd,
			match: true,
		},
		{
			name:  "token 'subtract' matches OpSubtract",
			token: "subtract",
			op:    OpSubtract,
			match: true,
		},
		{
			name:  "token 'multiply' matches OpMultiply",
			token: "multiply",
			op:    OpMultiply,
			match: true,
		},
		{
			name:  "token 'divide' matches OpDivide",
			token: "divide",
			op:    OpDivide,
			match: true,
		},
		{
			name:  "token 'add' does not match OpSubtract",
			token: "add",
			op:    OpSubtract,
			match: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Operation(tt.token) == tt.op
			if got != tt.match {
				t.Errorf("Operation(%q) == %q: got %v, want %v", tt.token, tt.op, got, tt.match)
			}
		})
	}
}
