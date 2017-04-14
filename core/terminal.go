package core

import (
	"fmt"

	"github.com/chzyer/readline"
)

// Terminal is our own version of the terminal to disable the default bell behavior
type Terminal struct {
	*readline.Terminal
}

// NewTerminal creates a wrapper over the terminal
func NewTerminal() (*Terminal, error) {
	return nil, nil
}

// Bell is overridden to disable the sound
func (t *Terminal) Bell() {
	return
}

// SoundBell uses the default Bell() behavior in order to not lose functionality.
func (t *Terminal) SoundBell() {
	fmt.Fprintf(t, "%c", readline.CharBell)
	t.Bell()
}

// DevNull is a writer that writes no where (ie, leave the printing up to me)
type DevNull struct{}

func (dn *DevNull) Write(p []byte) (n int, err error) {
	return len(p), nil
}
