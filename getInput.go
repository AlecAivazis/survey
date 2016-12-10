package survey

import (
	"errors"
	"fmt"

	"github.com/pkg/term"
)

// TerminalKey is a type used to refer to keys of interest
type TerminalKey int

// key codes for the common keys
const (
	KeyArrowUp TerminalKey = iota
	KeyArrowDown
	KeyArrowLeft
	KeyArrowRight
	KeyEsc
	KeyEnter
	KeyNull
)

// GetChar listens for input from the keyboard and returns the key value as a string
// or one of the Key* enum values.
func GetInput(format func(input string) string) (val string, keyCode TerminalKey, err error) {
	t, _ := term.Open("/dev/tty")
	term.RawMode(t)

	val = ""

InputLoop:
	for {
		bytes := make([]byte, 3)

		var numRead int
		numRead, err = t.Read(bytes)
		if err != nil {
			break InputLoop
		}

		// handle arrow-keys (three-character control sequence, beginning with "ESC-[")
		if numRead == 3 && bytes[0] == 27 && bytes[1] == 91 {
			switch bytes[2] {
			case 65:
				// Up
				keyCode = KeyArrowUp
				break InputLoop
			case 66:
				// Down
				keyCode = KeyArrowDown
				break InputLoop
			case 67:
				// Right
				keyCode = KeyArrowRight
				break InputLoop
			case 68:
				// Left
				keyCode = KeyArrowLeft
				break InputLoop
			}
		} else if numRead == 1 {
			ascii := int(bytes[0])

			// if the user sends SIGTERM (ascii 3) or presses esc (ascii 27)
			if ascii == 3 || ascii == 27 {
				// hard close
				err = errors.New("Goodbye.")
				break InputLoop
			}

			// handle the enter key
			if ascii == 13 {
				keyCode = KeyEnter
				break InputLoop
			}
		}

		// turn the ascii chars into a character
		value := string(bytes[:numRead])
		// add it to the running total
		val += value

		// if there is a formatter
		if format != nil {
			// print the
			fmt.Print(format(value))
		}
	}
	// clean up the terminal connection
	t.Restore()
	t.Close()

	// we're done here
	return val, keyCode, err
}
