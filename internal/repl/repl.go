package repl

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/workspace/repo/internal/arithmetic"
	"github.com/workspace/repo/internal/parser"
	"github.com/workspace/repo/internal/validator"
)

// Loop holds the dependencies for the REPL core loop.
type Loop struct {
	in        io.Reader
	out       io.Writer
	parser    parser.Parser
	validator validator.Validator
	engine    arithmetic.Engine
}

// New creates a new Loop with the given dependencies.
func New(in io.Reader, out io.Writer, p parser.Parser, v validator.Validator, e arithmetic.Engine) *Loop {
	return &Loop{
		in:        in,
		out:       out,
		parser:    p,
		validator: v,
		engine:    e,
	}
}

// Run starts the REPL loop, reading from in and writing to out until exit/quit or EOF.
func (l *Loop) Run(ctx context.Context) error {
	scanner := bufio.NewScanner(l.in)
	for {
		fmt.Fprint(l.out, "> ")
		if !scanner.Scan() {
			return scanner.Err()
		}
		text := strings.TrimSpace(scanner.Text())
		if text == "exit" || text == "quit" {
			return nil
		}
		cmd, err := l.parser.Parse(text)
		if err != nil {
			continue
		}
		if err := l.validator.Validate(cmd); err != nil {
			fmt.Fprintln(l.out, err)
			continue
		}
		result, err := l.engine.Execute(cmd)
		if err != nil {
			fmt.Fprintln(l.out, err)
			continue
		}
		fmt.Fprintln(l.out, result)
	}
}
