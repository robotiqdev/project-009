package parser

import (
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    Command
		expectError bool
	}{
		{
			name:  "normal input 'add 1.5 2.3' returns Command with operation and args",
			input: "add 1.5 2.3",
			expected: Command{
				Operation: "add",
				Args:      []string{"1.5", "2.3"},
			},
			expectError: false,
		},
		{
			name:  "extra whitespace '  add   1  2  ' returns same result as normal input",
			input: "  add   1  2  ",
			expected: Command{
				Operation: "add",
				Args:      []string{"1", "2"},
			},
			expectError: false,
		},
		{
			name:        "empty string returns error",
			input:       "",
			expected:    Command{},
			expectError: true,
		},
		{
			name:  "single token 'add' returns Command with operation and empty args",
			input: "add",
			expected: Command{
				Operation: "add",
				Args:      []string{},
			},
			expectError: false,
		},
		{
			name:  "many tokens 'add 1 2 3 4' returns all extra args included",
			input: "add 1 2 3 4",
			expected: Command{
				Operation: "add",
				Args:      []string{"1", "2", "3", "4"},
			},
			expectError: false,
		},
	}

	p := NewLineParser()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := p.Parse(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("Parse(%q) expected error, got nil", tt.input)
				}
				return
			}

			if err != nil {
				t.Errorf("Parse(%q) unexpected error: %v", tt.input, err)
				return
			}

			if got.Operation != tt.expected.Operation {
				t.Errorf("Parse(%q) Operation = %q, want %q", tt.input, got.Operation, tt.expected.Operation)
			}

			if len(got.Args) != len(tt.expected.Args) {
				t.Errorf("Parse(%q) Args length = %d, want %d; got %v, want %v",
					tt.input, len(got.Args), len(tt.expected.Args), got.Args, tt.expected.Args)
				return
			}

			for i, arg := range got.Args {
				if arg != tt.expected.Args[i] {
					t.Errorf("Parse(%q) Args[%d] = %q, want %q", tt.input, i, arg, tt.expected.Args[i])
				}
			}
		})
	}
}

func TestParseEmptyInputError(t *testing.T) {
	p := NewLineParser()
	_, err := p.Parse("")
	if err == nil {
		t.Fatal("Parse(\"\") expected ErrEmptyInput, got nil")
	}
	if err != ErrEmptyInput {
		t.Errorf("Parse(\"\") error = %v, want ErrEmptyInput", err)
	}
}

func TestParseWhitespaceOnlyReturnsError(t *testing.T) {
	p := NewLineParser()
	_, err := p.Parse("   ")
	if err == nil {
		t.Fatal("Parse(\"   \") expected error for whitespace-only input, got nil")
	}
}

func TestNewLineParserReturnsParser(t *testing.T) {
	p := NewLineParser()
	if p == nil {
		t.Fatal("NewLineParser() returned nil")
	}
}

func TestParserInterface(t *testing.T) {
	// Verify that *LineParser implements Parser interface
	var _ Parser = NewLineParser()
}
