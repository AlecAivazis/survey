package main

import (
	tm "github.com/buger/goterm"
	"time"
)

func main() {

	// get the current location of the cursor
	loc, err := CursorLocation()
	// if something went wrong
	if err != nil {
		// yell loudly
		panic(err)
	}
	initialRow := loc.col + 1
	for {
		// overwrite the line after the cursore
		tm.MoveCursor(initialRow, 1)

		tm.Print("Current Time: ", time.Now().Format(time.RFC1123))

		tm.Flush() // Call it every time at the end of rendering

		time.Sleep(time.Second)
	}
}
