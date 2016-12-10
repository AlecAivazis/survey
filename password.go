package survey

import (
	"fmt"

	"github.com/alecaivazis/survey/format"
)

// Password is like a normal Input but the text shows up as *'s and
// there is no default.
type Password struct {
	Message string
}

// the character to use to hide the input
var hideChar = "*"

// this function will be passed to the input handler to hide the input
func hideInput(input string) string {
	out := ""
	// fmt.Print(input, "h")
	for i := 0; i < len(input); i++ {
		out += hideChar
	}

	return out
}

// Prompt behaves like a normal int but hides the input.
func (prompt *Password) Prompt() (string, error) {
	// print the question we were given to kick off the prompt
	fmt.Print(format.Ask(fmt.Sprintf("%v ", prompt.Message), ""))

	// a running total
	value := ""

	// listen for input (this will ignore crazy characters like arrow keys)
	for val, keyCode, err := GetInput(hideInput); true; value, keyCode, err = GetInput(hideInput) {
		// if there is an error
		if err != nil {
			// bubble up
			return "", err
		}

		// add val to the running total
		value += val

		if keyCode == KeyEnter {
			fmt.Print("\n")
			return value, nil
		}
	}

	return value, nil
}

// Cleanup hides the string with a fixed number of characters.
func (prompt *Password) Cleanup(val string) error {
	return nil
}
