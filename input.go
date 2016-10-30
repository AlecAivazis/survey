package probe

import (
	"fmt"
)

type Input struct {
	Question string
}

// Inputs prompt the user with a simple text field and exepect a reply followed
// by a newline or carriage return.
func (input *Input) Prompt() (string, error) {
	// print the question we were given to kick off the prompt
	fmt.Print(fmt.Sprintf("%v ", input.Question))

	// a string to hold the user's input
	var res string
	// wait for a newline or carriage return
	fmt.Scanln(&res)
	// return the value
	return res, nil
}
