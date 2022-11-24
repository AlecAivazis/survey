package terminal

import (
	"io"
	"os"

	"github.com/mattn/go-colorable"
)

// NewAnsiStdout returns a writer connected to standard output that interprets ANSI escape codes to in a platform-agnostic way.
//
// Deprecated: use the mattn/go-colorable module instead of this method.
func NewAnsiStdout(out FileWriter) io.Writer {
	if f, ok := out.(*os.File); ok {
		return colorable.NewColorable(f)
	}
	return out
}

// NewAnsiStdout returns a writer connected to standard error that interprets ANSI escape codes to in a platform-agnostic way.
//
// Deprecated: use the mattn/go-colorable module instead of this method.
func NewAnsiStderr(out FileWriter) io.Writer {
	if f, ok := out.(*os.File); ok {
		return colorable.NewColorable(f)
	}
	return out
}
