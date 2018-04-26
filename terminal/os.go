package terminal

import "os"

var (
	Stdout = NewAnsiStdout()
	Stderr = NewAnsiStderr()
	Stdin  = os.Stdin
)
