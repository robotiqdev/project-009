package arithmetic

import "github.com/workspace/repo/internal/parser"

// Engine defines the interface for executing arithmetic commands.
type Engine interface {
	Execute(cmd parser.Command) (float64, error)
}
