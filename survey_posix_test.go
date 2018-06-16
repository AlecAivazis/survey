// +build !windows

package survey

import (
	"bytes"
	"testing"

	expect "github.com/Netflix/go-expect"
	"github.com/hinshun/vt10x"
	"github.com/kr/pty"
	"github.com/stretchr/testify/require"
	"gopkg.in/AlecAivazis/survey.v1/terminal"
)

func RunTest(t *testing.T, procedure func(*expect.Console), test func(terminal.Stdio) error) {
	t.Parallel()

	// Create a psuedoterminal for the virtual terminal's tty.
	ptm, pts, err := pty.Open()
	require.Nil(t, err)

	// Multiplex stdin/stdout to a virtual terminal to respond to ANSI escape
	// sequences (i.e. cursor position report).
	var state vt10x.State
	term, err := vt10x.Create(&state, pts)
	require.Nil(t, err)

	// Multiplex output to a buffer as well for the raw bytes.
	buf := new(bytes.Buffer)

	c, err := expect.NewConsole(expect.WithStdin(ptm), expect.WithStdout(term, buf), expect.WithCloser(pts, ptm, term))
	require.Nil(t, err)
	defer c.Close()

	donec := make(chan struct{})
	go func() {
		defer close(donec)
		procedure(c)
	}()

	err = test(Stdio(c))
	require.Nil(t, err)

	// Close the slave end of the pty, and read the remaining bytes from the master end.
	c.Tty().Close()
	<-donec

	t.Logf("Raw output: %q", buf.String())

	// Dump the terminal's screen.
	t.Logf("\n%s", expect.StripTrailingEmptyLines(state.String()))
}
