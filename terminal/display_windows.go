package terminal

import (
	"syscall"
	"unsafe"

	"github.com/AlecAivazis/survey/v2/log"
)

func EraseLine(out FileWriter, mode EraseLineMode) error {
	handle := syscall.Handle(out.Fd())

	var csbi consoleScreenBufferInfo
	if _, _, err := procGetConsoleScreenBufferInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&csbi))); normalizeError(err) != nil {
		return err
	}

	var w uint32
	var x Short
	cursor := csbi.cursorPosition
	switch mode {
	case ERASE_LINE_END:
		x = csbi.size.X
		log.Printf("EraseLine() END: %d", x)
	case ERASE_LINE_START:
		x = 0
		log.Printf("EraseLine() START: %d", x)
	case ERASE_LINE_ALL:
		cursor.X = 0
		x = csbi.size.X
		log.Printf("EraseLine() ALL: %d-%d", 0, x)
	}

	_, _, err := procFillConsoleOutputCharacter.Call(uintptr(handle), uintptr(' '), uintptr(x), uintptr(*(*int32)(unsafe.Pointer(&cursor))), uintptr(unsafe.Pointer(&w)))
	return normalizeError(err)
}
