// +build !windows

package core

import (
	"fmt"
)

func EraseInLine(mode int) {
	fmt.Printf("\x1b[%dK", mode)
}
