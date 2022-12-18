package terminal

import (
	"syscall"
	"unsafe"
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
	case ERASE_LINE_START:
		x = 0
	case ERASE_LINE_ALL:
		cursor.X = 0
		x = csbi.size.X
	}

	_, _, err := procFillConsoleOutputCharacter.Call(uintptr(handle), uintptr(' '), uintptr(x), uintptr(*(*int32)(unsafe.Pointer(&cursor))), uintptr(unsafe.Pointer(&w)))
	return normalizeError(err)
}

func EraseScreen(out FileWriter, mode EraseScreenMode) error {
	handle := syscall.Handle(out.Fd())

	var csbi consoleScreenBufferInfo
	if _, _, err := procGetConsoleScreenBufferInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&csbi))); normalizeError(err) != nil {
		return err
	}

	var w uint32
	termSize := uint32(csbi.size.X) * uint32(csbi.size.Y)
	cursor := csbi.cursorPosition

	_, _, err := procFillConsoleOutputCharacter.Call(uintptr(handle), uintptr(' '), uintptr(termSize), uintptr(*(*int32)(unsafe.Pointer(&cursor))), uintptr(unsafe.Pointer(&w)))
	return err
}
