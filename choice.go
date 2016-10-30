package probe

import (
	tm "github.com/buger/goterm"
	"time"
)

type Choice struct {
	Question string
	Choices  []string
}

func (prompt *Choice) Prompt() (string, error) {
	// get the current location of the cursor
	loc, err := CursorLocation()
	// if something went wrong
	if err != nil {
		// yell loudly
		panic(err)
	}

	// the height of the prompt's output
	height := 2

	// the starting point of the list depends on wether or not we
	// are at the bottom of the current terminal session
	var initialLocation int
	// if we are at the bottom
	if loc.col == tm.Height() {
		// the we have to start
		initialLocation = loc.col - height
		// otherwise we are not at the bottom of the terminal
	} else {
		// start at the current location
		initialLocation = loc.col
	}

	for {

		tm.Print("Current Time: ", time.Now().Format(time.RFC1123))
		tm.Print("\nHello")

		tm.Flush() // Call it every time at the end of rendering
		// make sure we overwrite the first line
		tm.MoveCursor(initialLocation, 1)

		time.Sleep(time.Second)
	}

	return "hello", nil
}
