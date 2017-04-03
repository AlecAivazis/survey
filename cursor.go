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

	"github.com/buger/goterm"
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

	// look for the cursor position in a sub process
	cmd := exec.Command("echo", AnsiPosition)
	randomBytes := &bytes.Buffer{}
	cmd.Stdout = randomBytes

	// Start command asynchronously
	cmd.Start()

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
		fmt.Print(AnsiCursorUp)

		// return the internal data structure with the location
		return &cursorCoordinate{x, y}, nil
	} else {
		return nil, errors.New("Could not find current location.")
	}
}

func InitialLocation(height int) (int, error) {
	// get the current location of the cursor
	loc, err := CursorLocation()
	// if something went wrong
	if err != nil {
		// bubble up
		return 0, err
	}

	// the starting point of the list depends on wether or not we
	// are at the bottom of the current terminal session
	var initialLocation int
	// if the options would fit cleanly
	if loc.col+height < goterm.Height() {
		// start at the current location
		initialLocation = loc.col
		// otherwise we will be placed at the bottom of the terminal after this print
	} else {
		// the we have to start printing so that we just fit
		initialLocation = goterm.Height() - height
	}

	return initialLocation, nil
}
