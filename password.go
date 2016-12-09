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

// Prompt behaves like a normal int but hides the input.
func (prompt *Password) Prompt() (string, error) {
	// a string to track the input
	input := ""

	// print the question we were given to kick off the prompt
	fmt.Print(format.Ask(fmt.Sprintf("%v ", prompt.Message), ""))

	// until we're intterupted
	for {
		// wait for an input from the user
		ascii, _, err := getChar()
		// if there is an error
		if err != nil {
			// bubble up
			return "", err
		}

		// if the user sends SIGTERM (ascii 3) or presses esc (ascii 27)
		if ascii == 3 || ascii == 27 {
			// hard close
			return "", fmt.Errorf("\nGoodbye.")
		}

		// if the user presses enter (ascii 13)
		if ascii == 13 {
			// we're done with the rendering loop (the current value is good)
			break
		}

		// handle paste

		// add the character to the running total
		fmt.Print(hideChar)
		input += string(ascii)
	}
	fmt.Print("\n")

	return input, nil
}

// Cleanup hides the string with a fixed number of characters.
func (prompt *Password) Cleanup(val string) error {
	return nil
}
