package probe

import (
	"fmt"

	"github.com/alecaivazis/probe/format"
)

type Input struct {
	Message string
}

// Inputs prompt the user with a simple text field and exepect a reply followed
// by a newline or carriage return.
func (input *Input) Prompt() (string, error) {
	// print the question we were given to kick off the prompt
	fmt.Print(format.Ask(fmt.Sprintf("%v ", input.Message)))

	// a string to hold the user's input
	var res string
	// wait for a newline or carriage return
	fmt.Scanln(&res)
	// return the value
	return res, nil
}
