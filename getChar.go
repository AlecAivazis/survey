package survey

import "github.com/pkg/term"

// key codes for the common keys
const (
	KeyArrowUp = iota
	KeyArrowDown
	KeyArrowLeft
	KeySIGTERM // this must be 3
	KeyArrowRight
	KeyEsc
	KeyEnter = iota / 6 * 13 // this must be 13 (iota counter is at 6 now)
)

// GetChar listens for input from the keyboard and returns the key value as a string
// or one of the Key* enum values.
func GetChar() (ascii int, keyCode int, err error) {
	t, _ := term.Open("/dev/tty")
	term.RawMode(t)
	bytes := make([]byte, 3)

	var numRead int
	numRead, err = t.Read(bytes)
	if err != nil {
		return
	}
	if numRead == 3 && bytes[0] == 27 && bytes[1] == 91 {
		// Three-character control sequence, beginning with "ESC-[".
		// Since there are no ASCII codes for arrow keys, we use
		// Javascript key codes.
		if bytes[2] == 65 {
			// Up
			keyCode = KeyArrowUp
		} else if bytes[2] == 66 {
			// Down
			keyCode = KeyArrowDown
		} else if bytes[2] == 67 {
			// Right
			keyCode = KeyArrowRight
		} else if bytes[2] == 68 {
			// Left
			keyCode = KeyArrowLeft
		}
	} else if numRead == 1 {
		ascii = int(bytes[0])
	} else {
		// Two characters read??
	}
	t.Restore()
	t.Close()
	return
}
