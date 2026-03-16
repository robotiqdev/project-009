package parser

import (
	"errors"
	"strings"
)

// ErrEmptyInput is returned when Parse is called with an empty or whitespace-only string.
var ErrEmptyInput = errors.New("empty input")

// Command holds the parsed operation and its arguments.
type Command struct {
	Operation string
	Args      []string
}

// Parser defines the interface for parsing input lines into Commands.
type Parser interface {
	Parse(line string) (Command, error)
}

// LineParser implements Parser using whitespace-field splitting.
type LineParser struct{}

// NewLineParser returns a new LineParser.
func NewLineParser() *LineParser {
	return &LineParser{}
}

// Parse splits line into fields; returns ErrEmptyInput if no fields found.
// The first field is Operation, remaining fields are Args.
func (p *LineParser) Parse(line string) (Command, error) {
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return Command{}, ErrEmptyInput
	}
	args := []string{}
	if len(fields) > 1 {
		args = fields[1:]
	}
	return Command{Operation: fields[0], Args: args}, nil
}
