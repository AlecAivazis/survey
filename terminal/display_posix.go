// +build !windows

package terminal

import (
	"fmt"
)

func EraseInLine(mode EraseLineMode) {
	fmt.Printf("\x1b[%dK", mode)
}
