package help

import (
	"fmt"
	"io"
)

// Printer is the interface for printing usage information.
type Printer interface {
	Print(w io.Writer)
}

// UsagePrinter implements the Printer interface.
type UsagePrinter struct{}

// NewUsagePrinter returns a new UsagePrinter.
func NewUsagePrinter() *UsagePrinter {
	return &UsagePrinter{}
}

// Print writes the usage text to w.
func (u *UsagePrinter) Print(w io.Writer) {
	fmt.Fprint(w, UsageText)
}

// UsageText contains the usage documentation for the calculator program.
const UsageText = `Usage: calculator [--help]

An interactive calculator supporting basic arithmetic operations.

Operations:
  add <num1> <num2>        Add two numbers
  subtract <num1> <num2>   Subtract num2 from num1
  multiply <num1> <num2>   Multiply two numbers
  divide <num1> <num2>     Divide num1 by num2

Output is formatted to 2 decimal places.

Examples:
  > add 1 2
  3.00
  > subtract 10 4
  6.00
  > multiply 3 5
  15.00
  > divide 10 2
  5.00

To exit, type 'exit' or 'quit', or press Ctrl+D.

Flags:
  --help    Show this help message
`
