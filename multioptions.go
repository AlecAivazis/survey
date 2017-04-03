package survey

import "github.com/buger/goterm"

func cleanupMultiOptions(height int, output string) error {
	// get the current location of the cursor
	loc, err := CursorLocation()
	// if something went wrong
	if err != nil {
		// yell loudly
		return err
	}

	var initLoc int
	// if the options would fit cleanly
	if loc.col+height <= goterm.Height() {
		// start at the current location
		initLoc = loc.col - height - 1
		// otherwise we will be placed at the bottom of the terminal after this print
	} else {
		// the we have to start printing so that we just fit
		initLoc = loc.col - height - 2
	}

	// start where we were told
	goterm.MoveCursor(initLoc, 1)
	goterm.Print(output, AnsiClearLine)
	// for each choice
	for i := 0; i < height; i++ {
		// add an empty line
		goterm.Print(AnsiClearLine)
		// print the output
		goterm.Flush()
	}
	// add an empty line
	goterm.Print(AnsiClearLine)
	// print the output
	goterm.Flush()
	goterm.MoveCursor(initLoc, 1)
	goterm.Flush()

	// nothing went wrong
	return nil
}
