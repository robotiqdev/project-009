package help

import (
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
