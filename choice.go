package probe

import (
	"errors"
	"fmt"
	tm "github.com/buger/goterm"
	"time"

	"github.com/alecaivazis/probe/format"
)

// Choice is a prompt that presents a
type Choice struct {
	Question string
	Choices  []string
}

func (prompt *Choice) Prompt() (string, error) {
	// ask the question
	fmt.Println(format.FormatAsk(prompt.Question))

	// get the current location of the cursor
	loc, err := CursorLocation()
	// if something went wrong
	if err != nil {
		// TODO: don't panic but quit better
		// yell loudly
		panic(err)
	}

	// the height of the prompt's output
	height := len(prompt.Choices)

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

	// start off with the first option selected
	// sel := 0

	for {
		// wait for an input from the user
		ascii, keyCode, err := getChar()
		// if there is an error
		if err != nil {
			// bubble up
			return "", err
		}

		// if the user sends SIGTERM (3) or presses esc (27)
		if ascii == 3 || ascii == 27 {
			// hard close
			return "", errors.New("Goodbye.")
		}

		// we need to render the options

		tm.Print(ascii, keyCode, err, "\n")
		tm.Print("Current Time: ", time.Now().Format(time.RFC1123))
		tm.Print("\nHello")

		tm.Flush() // Call it every time at the end of rendering

		// make sure we overwrite the first line next time we print
		tm.MoveCursor(initialLocation, 1)
	}

	return "hello", nil
}

func renderOptions(opts []string, selected bool) string {
	return ""
}
