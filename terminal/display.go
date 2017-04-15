// +build !windows

package terminal

import (
	"fmt"
)

func EraseInLine(mode int) {
	fmt.Printf("\x1b[%dK", mode)
}
