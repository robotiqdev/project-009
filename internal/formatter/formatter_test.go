package formatter

import (
	"errors"
	"fmt"
	"testing"

	"github.com/workspace/repo/internal/arithmetic"
	"github.com/workspace/repo/internal/parser"
)

func TestFormatResult(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected string
	}{
		{
			name:     "integer result 5.0 returns '5'",
			input:    5.0,
			expected: "5",
		},
		{
			name:     "one decimal 1.5 returns '1.5'",
			input:    1.5,
			expected: "1.5",
		},
		{
			name:     "two decimals 3.14 returns '3.14'",
			input:    3.14,
			expected: "3.14",
		},
		{
			name:     "division result 1.0/3.0 returns '0.33'",
			input:    1.0 / 3.0,
			expected: "0.33",
		},
		{
			name:     "negative -2.5 returns '-2.5'",
			input:    -2.5,
			expected: "-2.5",
		},
		{
			name:     "zero 0.0 returns '0'",
			input:    0.0,
			expected: "0",
		},
		{
			name:     "large integer 1000.0 returns '1000'",
			input:    1000.0,
			expected: "1000",
		},
		{
			name:     "exactly two decimals 0.12 returns '0.12'",
			input:    0.12,
			expected: "0.12",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatResult(tt.input)
			if got != tt.expected {
				t.Errorf("FormatResult(%v) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestFormatError(t *testing.T) {
	tests := []struct {
		name     string
		input    error
		expected string
	}{
		{
			name:     "division by zero returns human-readable message",
			input:    arithmetic.ErrDivisionByZero,
			expected: "Error: division by zero",
		},
		{
			name:     "unknown operation returns human-readable message",
			input:    arithmetic.ErrUnknownOperation,
			expected: "Error: unknown operation",
		},
		{
			name:     "invalid argument count returns human-readable message",
			input:    arithmetic.ErrInvalidArgCount,
			expected: "Error: invalid argument count",
		},
		{
			name:     "invalid number returns human-readable message",
			input:    parser.ErrInvalidNumber,
			expected: "Error: invalid number",
		},
		{
			name:     "unknown error returns generic message using error text",
			input:    errors.New("something unexpected"),
			expected: fmt.Sprintf("Error: %v", errors.New("something unexpected")),
		},
		{
			name:     "wrapped division by zero matches via errors.Is",
			input:    fmt.Errorf("wrapped: %w", arithmetic.ErrDivisionByZero),
			expected: "Error: division by zero",
		},
		{
			name:     "wrapped unknown operation matches via errors.Is",
			input:    fmt.Errorf("wrapped: %w", arithmetic.ErrUnknownOperation),
			expected: "Error: unknown operation",
		},
		{
			name:     "wrapped invalid arg count matches via errors.Is",
			input:    fmt.Errorf("wrapped: %w", arithmetic.ErrInvalidArgCount),
			expected: "Error: invalid argument count",
		},
		{
			name:     "wrapped invalid number matches via errors.Is",
			input:    fmt.Errorf("wrapped: %w", parser.ErrInvalidNumber),
			expected: "Error: invalid number",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatError(tt.input)
			if got != tt.expected {
				t.Errorf("FormatError(%v) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestFormatError_DoesNotPanic(t *testing.T) {
	// Unknown errors should produce a generic message, not panic
	unknownErr := errors.New("totally unknown error")
	got := FormatError(unknownErr)
	if got == "" {
		t.Error("FormatError with unknown error returned empty string, expected a generic message")
	}
	expected := fmt.Sprintf("Error: %v", unknownErr)
	if got != expected {
		t.Errorf("FormatError(%v) = %q, want %q", unknownErr, got, expected)
	}
}
