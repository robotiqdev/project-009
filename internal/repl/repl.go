package repl

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/workspace/repo/internal/validator"
)

// Result holds the output of an engine execution.
type Result struct {
	Formatted string
}

// Parser parses raw input into a Command.
type Parser interface {
	Parse(input string) (validator.Command, error)
}

// Validator validates a Command and returns a ValidatedCommand or error.
type Validator interface {
	Validate(cmd validator.Command) (*validator.ValidatedCommand, error)
}

// Engine executes a validated command and returns a Result.
type Engine interface {
	Execute(operation string, args []float64) (Result, error)
}

// Repl is a read-eval-print loop for the calculator.
type Repl struct {
	parser    Parser
	validator Validator
	engine    Engine
}

// New creates a new Repl with the given dependencies.
func New(p Parser, v Validator, e Engine) *Repl {
	return &Repl{parser: p, validator: v, engine: e}
}

// Run starts the REPL, reading from in and writing to out.
// It loops until EOF or an exit/quit command is received.
func (r *Repl) Run(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "exit" || line == "quit" {
			return
		}

		cmd, err := r.parser.Parse(line)
		if err != nil {
			fmt.Fprintf(out, "Error: %s\n", err.Error())
			continue
		}

		vcmd, err := r.validator.Validate(cmd)
		if err != nil {
			fmt.Fprintf(out, "Error: %s\n", err.Error())
			continue
		}

		result, err := r.engine.Execute(vcmd.Operation, vcmd.Args)
		if err != nil {
			fmt.Fprintf(out, "Error: %s\n", err.Error())
			continue
		}

		fmt.Fprintf(out, "%s\n", result.Formatted)
	}
}
