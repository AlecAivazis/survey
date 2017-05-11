package terminal

import (
	"fmt"
	"os"
)

type RuneReader struct {
	Input *os.File

	state runeReaderState
}

func NewRuneReader(input *os.File) *RuneReader {
	return &RuneReader{
		Input: input,
		state: newRuneReaderState(input),
	}
}

func (rr *RuneReader) ReadLine(mask rune) ([]rune, error) {
	line := []rune{}
	for {
		r, _, err := rr.ReadRune()
		if err != nil {
			return line, err
		}
		if r == '\r' || r == '\n' || r == KeyEndTransmission {
			Print("\r\n")
			return line, nil
		}
		if r == KeyInterrupt {
			Print("\r\n")
			return line, fmt.Errorf("interrupt")
		}
		// allow for backspace/delete editing of password
		if r == KeyBackspace || r == KeyDelete {
			if len(line) > 0 {
				line = line[:len(line)-1]
				CursorBack(1)
				EraseLine(ERASE_LINE_END)
			}
			continue
		}
		line = append(line, r)
		if mask == 0 {
			Printf("%c", r)
		} else {
			Printf("%c", mask)
		}
	}
}
