package core

// this file defines the interface between readline and the rest of the package

import "github.com/chzyer/readline"

// GetReadline returns the readline instance with the correct configuration
func GetReadline() (*readline.Instance, error) {
	// create an instance
	rl, err := readline.NewEx(&readline.Config{
		InterruptPrompt: "^C",
		EOFPrompt:       "Goodbye",
		HistoryLimit:    -1,
	})
	if err != nil {
		return nil, err
	}

	return rl, nil
}
