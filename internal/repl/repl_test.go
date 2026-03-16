package repl_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/workspace/repo/internal/parser"
	"github.com/workspace/repo/internal/repl"
)

// --- Mock implementations ---

// mockParser returns a fixed Command for any non-empty input, or ErrEmptyInput for empty.
type mockParser struct {
	parseFunc func(line string) (parser.Command, error)
}

func (m *mockParser) Parse(line string) (parser.Command, error) {
	if m.parseFunc != nil {
		return m.parseFunc(line)
	}
	if strings.TrimSpace(line) == "" {
		return parser.Command{}, parser.ErrEmptyInput
	}
	return parser.Command{Operation: line, Args: []string{}}, nil
}

// mockValidator always validates successfully unless configured otherwise.
type mockValidator struct {
	validateFunc func(cmd parser.Command) error
}

func (m *mockValidator) Validate(cmd parser.Command) error {
	if m.validateFunc != nil {
		return m.validateFunc(cmd)
	}
	return nil
}

// mockEngine returns a fixed result unless configured otherwise.
type mockEngine struct {
	executeFunc func(cmd parser.Command) (float64, error)
}

func (m *mockEngine) Execute(cmd parser.Command) (float64, error) {
	if m.executeFunc != nil {
		return m.executeFunc(cmd)
	}
	return 42.0, nil
}

// --- Helpers ---

func newTestLoop(input string) (*repl.Loop, *bytes.Buffer) {
	in := strings.NewReader(input)
	out := &bytes.Buffer{}
	p := &mockParser{}
	v := &mockValidator{}
	e := &mockEngine{}
	loop := repl.New(in, out, p, v, e)
	return loop, out
}

// --- Tests ---

// TestSingleCommandThenExit verifies that a single command followed by 'exit'
// runs without error and that the prompt '> ' appears in the output.
func TestSingleCommandThenExit(t *testing.T) {
	in := strings.NewReader("add 1 2\nexit\n")
	out := &bytes.Buffer{}
	p := &mockParser{}
	v := &mockValidator{}
	e := &mockEngine{}

	loop := repl.New(in, out, p, v, e)
	err := loop.Run(context.Background())

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "> ") {
		t.Errorf("expected output to contain prompt '> ', got: %q", output)
	}
}

// TestEOFOnEmptyReader verifies that when the reader is empty (immediate EOF),
// the loop exits cleanly with no error.
func TestEOFOnEmptyReader(t *testing.T) {
	in := strings.NewReader("")
	out := &bytes.Buffer{}
	p := &mockParser{}
	v := &mockValidator{}
	e := &mockEngine{}

	loop := repl.New(in, out, p, v, e)
	err := loop.Run(context.Background())

	if err != nil {
		t.Fatalf("expected no error on EOF, got: %v", err)
	}
}

// TestQuitCommandExitsCleanly verifies that 'quit' causes the loop to exit
// with no error.
func TestQuitCommandExitsCleanly(t *testing.T) {
	in := strings.NewReader("quit\n")
	out := &bytes.Buffer{}
	p := &mockParser{}
	v := &mockValidator{}
	e := &mockEngine{}

	loop := repl.New(in, out, p, v, e)
	err := loop.Run(context.Background())

	if err != nil {
		t.Fatalf("expected no error for 'quit', got: %v", err)
	}
}

// TestMultipleCommandsEachProducesOutput verifies that multiple valid commands
// each produce output, the loop continues after each, and exits cleanly.
func TestMultipleCommandsEachProducesOutput(t *testing.T) {
	in := strings.NewReader("add 1 2\nmul 3 4\nsub 10 5\nexit\n")
	out := &bytes.Buffer{}

	commandCount := 0
	p := &mockParser{
		parseFunc: func(line string) (parser.Command, error) {
			if strings.TrimSpace(line) == "" || line == "exit" || line == "quit" {
				return parser.Command{}, parser.ErrEmptyInput
			}
			return parser.Command{Operation: line, Args: []string{}}, nil
		},
	}
	v := &mockValidator{}
	e := &mockEngine{
		executeFunc: func(cmd parser.Command) (float64, error) {
			commandCount++
			return float64(commandCount) * 10.0, nil
		},
	}

	loop := repl.New(in, out, p, v, e)
	err := loop.Run(context.Background())

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	output := out.String()

	// Prompt should appear multiple times (once per iteration)
	promptCount := strings.Count(output, "> ")
	if promptCount < 3 {
		t.Errorf("expected at least 3 prompts for 3 commands + exit, got %d in output: %q", promptCount, output)
	}

	// Output should contain results from each command execution
	if commandCount < 3 {
		t.Errorf("expected engine to be called at least 3 times, was called %d times", commandCount)
	}
}

// TestPromptAppearsBeforeEachCommand verifies that the prompt '> ' is printed
// before each command is read, including before the exit command.
func TestPromptAppearsBeforeEachCommand(t *testing.T) {
	in := strings.NewReader("add 1 2\nexit\n")
	out := &bytes.Buffer{}
	p := &mockParser{}
	v := &mockValidator{}
	e := &mockEngine{}

	loop := repl.New(in, out, p, v, e)
	err := loop.Run(context.Background())

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	output := out.String()
	// At minimum, one prompt before "add 1 2" and one before "exit"
	promptCount := strings.Count(output, "> ")
	if promptCount < 2 {
		t.Errorf("expected at least 2 prompts, got %d in output: %q", promptCount, output)
	}
}

// TestExitCommandExitsCleanly verifies that 'exit' causes the loop to exit
// with no error.
func TestExitCommandExitsCleanly(t *testing.T) {
	in := strings.NewReader("exit\n")
	out := &bytes.Buffer{}
	p := &mockParser{}
	v := &mockValidator{}
	e := &mockEngine{}

	loop := repl.New(in, out, p, v, e)
	err := loop.Run(context.Background())

	if err != nil {
		t.Fatalf("expected no error for 'exit', got: %v", err)
	}
}

// TestEmptyInputContinuesLoop verifies that an empty line (parse error) causes
// the loop to continue rather than exit or error.
func TestEmptyInputContinuesLoop(t *testing.T) {
	in := strings.NewReader("\nadd 1 2\nexit\n")
	out := &bytes.Buffer{}

	callCount := 0
	p := &mockParser{
		parseFunc: func(line string) (parser.Command, error) {
			if strings.TrimSpace(line) == "" {
				return parser.Command{}, parser.ErrEmptyInput
			}
			return parser.Command{Operation: line, Args: []string{}}, nil
		},
	}
	v := &mockValidator{}
	e := &mockEngine{
		executeFunc: func(cmd parser.Command) (float64, error) {
			callCount++
			return 1.0, nil
		},
	}

	loop := repl.New(in, out, p, v, e)
	err := loop.Run(context.Background())

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// The valid command "add 1 2" should have been executed despite the empty line
	if callCount != 1 {
		t.Errorf("expected engine called once for valid command, got %d calls", callCount)
	}
}

// TestNewReturnsNonNilLoop verifies that New returns a non-nil *Loop.
func TestNewReturnsNonNilLoop(t *testing.T) {
	in := strings.NewReader("")
	out := &bytes.Buffer{}
	p := &mockParser{}
	v := &mockValidator{}
	e := &mockEngine{}

	loop := repl.New(in, out, p, v, e)
	if loop == nil {
		t.Fatal("expected New to return non-nil *Loop")
	}
}
