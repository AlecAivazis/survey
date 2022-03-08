//go:build !windows
// +build !windows

package survey

import (
	"bytes"
	"strings"
	"testing"

	"github.com/AlecAivazis/survey/v2/terminal"
	pseudotty "github.com/creack/pty"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRenderer_countLines(t *testing.T) {
	t.Parallel()

	termWidth := 72
	pty, tty, err := pseudotty.Open()
	require.Nil(t, err)
	defer pty.Close()
	defer tty.Close()

	err = pseudotty.Setsize(tty, &pseudotty.Winsize{
		Rows: 30,
		Cols: uint16(termWidth),
	})
	require.Nil(t, err)

	r := Renderer{
		stdio: terminal.Stdio{
			In:  tty,
			Out: tty,
			Err: tty,
		},
	}

	tests := []struct {
		name  string
		buf   *bytes.Buffer
		wants int
	}{
		{
			name:  "empty",
			buf:   new(bytes.Buffer),
			wants: 0,
		},
		{
			name:  "no newline",
			buf:   bytes.NewBufferString("hello"),
			wants: 0,
		},
		{
			name:  "short line",
			buf:   bytes.NewBufferString("hello\n"),
			wants: 1,
		},
		{
			name:  "three short lines",
			buf:   bytes.NewBufferString("hello\nbeautiful\nworld\n"),
			wants: 3,
		},
		{
			name:  "full line",
			buf:   bytes.NewBufferString(strings.Repeat("A", termWidth) + "\n"),
			wants: 1,
		},
		{
			name:  "overflow",
			buf:   bytes.NewBufferString(strings.Repeat("A", termWidth+1) + "\n"),
			wants: 2,
		},
		{
			name:  "overflow fills 2nd line",
			buf:   bytes.NewBufferString(strings.Repeat("A", termWidth*2) + "\n"),
			wants: 2,
		},
		{
			name:  "overflow spills to 3rd line",
			buf:   bytes.NewBufferString(strings.Repeat("A", termWidth*2+1) + "\n"),
			wants: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := r.countLines(*tt.buf)
			assert.Equal(t, tt.wants, n)
		})
	}
}
