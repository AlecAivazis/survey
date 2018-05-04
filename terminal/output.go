// +build !windows

package terminal

import (
	"io"
	"os"
)

var (
	Stdout io.Writer = NewAnsiStdout()
	Stderr io.Writer = NewAnsiStderr()
)

// Returns special stdout, which converts escape sequences to Windows API calls
// on Windows environment.
func NewAnsiStdout() io.Writer {
	return os.Stdout
}

// Returns special stderr, which converts escape sequences to Windows API calls
// on Windows environment.
func NewAnsiStderr() io.Writer {
	return os.Stderr
}
