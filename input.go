package survey

import (
	"fmt"
	tm "github.com/buger/goterm"

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
