package survey

import (
	"bufio"
	"errors"
	"fmt"
	tm "github.com/buger/goterm"
	"os"

	"github.com/alecaivazis/survey/format"
)

// Input is a regular text input that prints each character the user types on the screen
// and accepts the input with the enter key.
type Input struct {
	Message string
	Default string
}

// Prompt prompts the user with a simple text field and expects a reply followed
// by a carriage return.
func (input *Input) Prompt() (string, error) {
	// print the question we were given to kick off the prompt
	fmt.Print(format.Ask(fmt.Sprintf("%v ", input.Message), input.Default))

	// a scanner to look at the input from stdin
	scanner := bufio.NewScanner(os.Stdin)
	// wait for a response
	for scanner.Scan() {
		// get the availible text in the scanner
		res := scanner.Text()
		// if there is no answer
		if res == "" {
			// use the default
			res = input.Default
		}

		// return the value
		return res, nil
	}

	return "", errors.New("Did not get input.")
}

// Cleanup overwrite the line with the finalized formatted version
func (input *Input) Cleanup(val string) error {
	// get the current cursor location
	loc, err := CursorLocation()
	// if something went wrong
	if err != nil {
		// bubble
		return err
	}

	var initLoc int
	// if we are printing at the end of the console
	if loc.col == tm.Height() {
		initLoc = loc.col - 2
	} else {
		initLoc = loc.col - 1
	}

	// move to the beginning of the current line
	tm.MoveCursor(initLoc, 1)

	tm.Print(format.Response(input.Message, val), "\x1b[0K")
	tm.Flush()

	// nothing went wrong
	return nil
}
