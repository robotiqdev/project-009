package help

import (
	"bytes"
	"strings"
	"testing"
)

func TestUsageTextContainsHelpFlag(t *testing.T) {
	if !strings.Contains(UsageText, "--help") {
		t.Errorf("UsageText does not contain '--help'")
	}
}

func TestUsageTextContainsAddOperation(t *testing.T) {
	if !strings.Contains(UsageText, "add") {
		t.Errorf("UsageText does not contain 'add'")
	}
}

func TestUsageTextContainsSubtractOperation(t *testing.T) {
	if !strings.Contains(UsageText, "subtract") {
		t.Errorf("UsageText does not contain 'subtract'")
	}
}

func TestUsageTextContainsMultiplyOperation(t *testing.T) {
	if !strings.Contains(UsageText, "multiply") {
		t.Errorf("UsageText does not contain 'multiply'")
	}
}

func TestUsageTextContainsDivideOperation(t *testing.T) {
	if !strings.Contains(UsageText, "divide") {
		t.Errorf("UsageText does not contain 'divide'")
	}
}

func TestUsageTextContainsExitInstruction(t *testing.T) {
	if !strings.Contains(UsageText, "exit") {
		t.Errorf("UsageText does not contain 'exit'")
	}
}

func TestUsageTextContainsAddExample(t *testing.T) {
	if !strings.Contains(UsageText, "add 1 2") {
		t.Errorf("UsageText does not contain example 'add 1 2'")
	}
}

func TestUsageTextIsNonEmpty(t *testing.T) {
	if len(strings.TrimSpace(UsageText)) == 0 {
		t.Errorf("UsageText is empty")
	}
}

func TestUsageTextContainsAllRequiredSubstrings(t *testing.T) {
	required := []string{
		"--help",
		"add",
		"subtract",
		"multiply",
		"divide",
		"exit",
		"add 1 2",
	}
	for _, substr := range required {
		if !strings.Contains(UsageText, substr) {
			t.Errorf("UsageText missing required substring: %q", substr)
		}
	}
}

// Compile-time check: *UsagePrinter satisfies Printer interface.
var _ Printer = (*UsagePrinter)(nil)

func TestNewUsagePrinterReturnsNonNil(t *testing.T) {
	p := NewUsagePrinter()
	if p == nil {
		t.Error("NewUsagePrinter() returned nil")
	}
}

func TestPrintWritesToWriter(t *testing.T) {
	p := NewUsagePrinter()
	var buf bytes.Buffer
	p.Print(&buf)
	if buf.Len() == 0 {
		t.Error("Print() wrote nothing to the writer; buffer is empty")
	}
}

func TestPrintBufferContainsAdd(t *testing.T) {
	p := NewUsagePrinter()
	var buf bytes.Buffer
	p.Print(&buf)
	if !strings.Contains(buf.String(), "add") {
		t.Errorf("Print() output does not contain 'add'; got: %q", buf.String())
	}
}

func TestPrintBufferContainsSubtract(t *testing.T) {
	p := NewUsagePrinter()
	var buf bytes.Buffer
	p.Print(&buf)
	if !strings.Contains(buf.String(), "subtract") {
		t.Errorf("Print() output does not contain 'subtract'; got: %q", buf.String())
	}
}

func TestPrintBufferContainsMultiply(t *testing.T) {
	p := NewUsagePrinter()
	var buf bytes.Buffer
	p.Print(&buf)
	if !strings.Contains(buf.String(), "multiply") {
		t.Errorf("Print() output does not contain 'multiply'; got: %q", buf.String())
	}
}

func TestPrintBufferContainsDivide(t *testing.T) {
	p := NewUsagePrinter()
	var buf bytes.Buffer
	p.Print(&buf)
	if !strings.Contains(buf.String(), "divide") {
		t.Errorf("Print() output does not contain 'divide'; got: %q", buf.String())
	}
}

func TestPrintBufferContainsHelpFlag(t *testing.T) {
	p := NewUsagePrinter()
	var buf bytes.Buffer
	p.Print(&buf)
	if !strings.Contains(buf.String(), "--help") {
		t.Errorf("Print() output does not contain '--help'; got: %q", buf.String())
	}
}

func TestPrintBufferNonEmpty(t *testing.T) {
	p := NewUsagePrinter()
	var buf bytes.Buffer
	p.Print(&buf)
	if buf.Len() == 0 {
		t.Error("Print() output buffer is empty")
	}
}

func TestPrintAllRequiredSubstrings(t *testing.T) {
	p := NewUsagePrinter()
	var buf bytes.Buffer
	p.Print(&buf)
	output := buf.String()
	required := []string{"add", "subtract", "multiply", "divide", "--help"}
	for _, substr := range required {
		if !strings.Contains(output, substr) {
			t.Errorf("Print() output missing required substring: %q", substr)
		}
	}
}
