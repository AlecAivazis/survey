package survey

import (
	"fmt"

	"github.com/alecaivazis/survey/format"
)

type Input struct {
	Message string
	Default string
}

// Inputs prompt the user with a simple text field and exepect a reply followed
// by a newline or carriage return.
func (input *Input) Prompt() (string, error) {
	// print the question we were given to kick off the prompt
	fmt.Print(format.Ask(fmt.Sprintf("%v ", input.Message), input.Default))

	// a string to hold the user's input
	var res string
	// wait for a newline or carriage return
	fmt.Scanln(&res)

	// if there is no answer
	if res == "" {
		// use the default
		res = input.Default
	}

	// return the value
	return res, nil
}
