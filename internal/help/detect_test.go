package help

import "testing"

// TestDetectHelpFlag_LongFlag verifies that --help is recognized as a help flag.
func TestDetectHelpFlag_LongFlag(t *testing.T) {
	args := []string{"--help"}
	if !DetectHelpFlag(args) {
		t.Errorf("DetectHelpFlag(%v) = false; want true", args)
	}
}

// TestDetectHelpFlag_ShortFlag verifies that -h is recognized as a help flag.
func TestDetectHelpFlag_ShortFlag(t *testing.T) {
	args := []string{"-h"}
	if !DetectHelpFlag(args) {
		t.Errorf("DetectHelpFlag(%v) = false; want true", args)
	}
}

// TestDetectHelpFlag_ArithmeticArgs verifies that normal calculator args are not help flags.
func TestDetectHelpFlag_ArithmeticArgs(t *testing.T) {
	args := []string{"add", "1", "2"}
	if DetectHelpFlag(args) {
		t.Errorf("DetectHelpFlag(%v) = true; want false", args)
	}
}

// TestDetectHelpFlag_EmptyArgs verifies that an empty arg slice is not a help flag.
func TestDetectHelpFlag_EmptyArgs(t *testing.T) {
	args := []string{}
	if DetectHelpFlag(args) {
		t.Errorf("DetectHelpFlag(%v) = true; want false", args)
	}
}

// TestDetectHelpFlag_NilArgs verifies that a nil arg slice is not a help flag.
func TestDetectHelpFlag_NilArgs(t *testing.T) {
	if DetectHelpFlag(nil) {
		t.Errorf("DetectHelpFlag(nil) = true; want false")
	}
}

// TestDetectHelpFlag_HelpAmongOtherArgs verifies that --help is detected when mixed with other args.
func TestDetectHelpFlag_HelpAmongOtherArgs(t *testing.T) {
	args := []string{"add", "--help", "1"}
	if !DetectHelpFlag(args) {
		t.Errorf("DetectHelpFlag(%v) = false; want true", args)
	}
}

// TestDetectHelpFlag_TableDriven runs multiple scenarios using table-driven style.
func TestDetectHelpFlag_TableDriven(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantHelp bool
	}{
		{
			name:     "long help flag only",
			args:     []string{"--help"},
			wantHelp: true,
		},
		{
			name:     "short help flag only",
			args:     []string{"-h"},
			wantHelp: true,
		},
		{
			name:     "help flag first among others",
			args:     []string{"--help", "add", "1", "2"},
			wantHelp: true,
		},
		{
			name:     "help flag last among others",
			args:     []string{"add", "1", "2", "--help"},
			wantHelp: true,
		},
		{
			name:     "arithmetic args only",
			args:     []string{"add", "1", "2"},
			wantHelp: false,
		},
		{
			name:     "subtract command",
			args:     []string{"subtract", "10", "4"},
			wantHelp: false,
		},
		{
			name:     "multiply command",
			args:     []string{"multiply", "3", "5"},
			wantHelp: false,
		},
		{
			name:     "divide command",
			args:     []string{"divide", "10", "2"},
			wantHelp: false,
		},
		{
			name:     "empty args",
			args:     []string{},
			wantHelp: false,
		},
		{
			name:     "nil args",
			args:     nil,
			wantHelp: false,
		},
		{
			name:     "word containing help but not flag",
			args:     []string{"helpful"},
			wantHelp: false,
		},
		{
			name:     "double dash alone",
			args:     []string{"--"},
			wantHelp: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DetectHelpFlag(tt.args)
			if got != tt.wantHelp {
				t.Errorf("DetectHelpFlag(%v) = %v; want %v", tt.args, got, tt.wantHelp)
			}
		})
	}
}
