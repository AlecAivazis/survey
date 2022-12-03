package terminal

import (
	"fmt"
)

type EraseLineMode int
type EraseScreenMode int

const (
	ERASE_LINE_END EraseLineMode = iota
	ERASE_LINE_START
	ERASE_LINE_ALL
)

const (
	ERASE_SCREEN_END EraseScreenMode = iota
	ERASE_SCREEN_START
	ERASE_SCREEN_ALL
)

func EraseScreen(out FileWriter, mode EraseScreenMode) error {
	_, err := fmt.Fprintf(out, "\x1b[%dJ", mode)
	return err
}
