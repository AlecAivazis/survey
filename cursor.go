package survey

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type cursorCoordinate struct {
	col int
	row int
}

// CursorLocation returns the location (col, row) of the cursor in the current terminal
// session.
func CursorLocation() (*cursorCoordinate, error) {

	// Set the terminal to raw mode (to be undone with `-raw`)
	rawMode := exec.Command("/bin/stty", "raw")
	rawMode.Stdin = os.Stdin
	_ = rawMode.Run()
	rawMode.Wait()

	//
	// same as $ echo -e "\033[6n"
	cmd := exec.Command("echo", fmt.Sprintf("%c[6n", 27))
	randomBytes := &bytes.Buffer{}
	cmd.Stdout = randomBytes

	// Start command asynchronously
	_ = cmd.Start()

	// capture keyboard output from echo command
	reader := bufio.NewReader(os.Stdin)
	cmd.Wait()

	// by printing the command output, we are triggering input
	fmt.Print(randomBytes)
	// capture the triggered stdin from the print
	text, _ := reader.ReadSlice('R')

	// Set the terminal back from raw mode to 'cooked'
	rawModeOff := exec.Command("/bin/stty", "-raw")
	rawModeOff.Stdin = os.Stdin
	_ = rawModeOff.Run()
	rawModeOff.Wait()

	// check for the desired output
	if strings.Contains(string(text), ";") {
		re := regexp.MustCompile(`\d+;\d+`)
		line := strings.Split(re.FindString(string(text)), ";")
		// turn the coordinates into integers
		x, err := strconv.Atoi(line[0])
		if err != nil {
			return nil, err
		}
		y, err := strconv.Atoi(line[1])
		if err != nil {
			return nil, err
		}
		// go up one line to cover our tracks
		fmt.Print("\033[F")

		// return the internal data structure with the location
		return &cursorCoordinate{x, y}, nil
	} else {
		return nil, errors.New("Could not find current location.")
	}
}
