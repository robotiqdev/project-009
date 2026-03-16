package repl

import (
	"context"
	"io"

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
	return &Loop{}
}

// Run starts the REPL loop, reading from in and writing to out until exit/quit or EOF.
func (l *Loop) Run(ctx context.Context) error {
	return nil
}
