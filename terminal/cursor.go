// +build !windows

package terminal

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Move the cursor n cells to up.
func CursorUp(n int) {
	fmt.Printf("\x1b[%dA", n)
}

// Move the cursor n cells to down.
func CursorDown(n int) {
	fmt.Printf("\x1b[%dB", n)
}

// Move the cursor n cells to right.
func CursorForward(n int) {
	fmt.Printf("\x1b[%dC", n)
}

// Move the cursor n cells to left.
func CursorBack(n int) {
	fmt.Printf("\x1b[%dD", n)
}

// Move cursor to beginning of the line n lines down.
func CursorNextLine(n int) {
	fmt.Printf("\x1b[%dE", n)
}

// Move cursor to beginning of the line n lines up.
func CursorPreviousLine(n int) {
	fmt.Printf("\x1b[%dF", n)
}

// Move cursor horizontally to x.
func CursorHorizontalAbsolute(x int) {
	fmt.Printf("\x1b[%dG", x)
}

// Show the cursor.
func CursorShow() {
	fmt.Print("\x1b[?25h")
}

// Hide the cursor.
func CursorHide() {
	fmt.Print("\x1b[?25l")
}

// CursorLocation returns the current location of the cursor in the terminal
func CursorLocation() (*Coord, error) {
	// print the escape sequence to recieve the position in our stdin
	fmt.Print("\x1b[6n")

	// read from stdin to get the response
	reader := bufio.NewReader(os.Stdin)
	// spec says we read 'til R, so do that
	text, err := reader.ReadSlice('R')
	if err != nil {
		return nil, err
	}

	// spec also says they're split by ;, so do that too
	if strings.Contains(string(text), ";") {
		// a regex to parse the output of the ansi code
		re := regexp.MustCompile(`\d+;\d+`)
		line := re.FindString(string(text))

		// find the column and rows embedded in the string
		coords := strings.Split(line, ";")

		// try to cast the col number to an int
		col, err := strconv.Atoi(coords[1])
		if err != nil {
			return nil, err
		}

		// try to cast the row number to an int
		row, err := strconv.Atoi(coords[0])
		if err != nil {
			return nil, err
		}

		// return the coordinate object with the col and row we calculated
		return &Coord{Short(col), Short(row)}, nil
	}

	// it didn't work so return an error
	return nil, fmt.Errorf("could not compute the cursor position using ascii escape sequences")
}
