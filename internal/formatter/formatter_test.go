package formatter

import "testing"

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
