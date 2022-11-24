package terminal

import (
	"bytes"
	"fmt"
	"syscall"
	"unsafe"

	"github.com/AlecAivazis/survey/v2/log"
)

var COORDINATE_SYSTEM_BEGIN Short = 0

// shared variable to save the cursor location from CursorSave()
var cursorLoc Coord

type Cursor struct {
	In  FileReader
	Out FileWriter
}

func (c *Cursor) Up(n int) error {
	return c.cursorMove(0, -1*n, false)
}

func (c *Cursor) Down(n int) error {
	return c.cursorMove(0, n, false)
}

func (c *Cursor) Forward(n int) error {
	return c.cursorMove(n, 0, false)
}

func (c *Cursor) Back(n int) error {
	return c.cursorMove(-1*n, 0, false)
}

// save the cursor location
func (c *Cursor) Save() error {
	loc, err := c.Location(nil)
	if err != nil {
		return err
	}
	cursorLoc = *loc
	return nil
}

func (c *Cursor) Restore() error {
	handle := syscall.Handle(c.Out.Fd())
	// restore it to the original position
	_, _, err := procSetConsoleCursorPosition.Call(uintptr(handle), uintptr(*(*int32)(unsafe.Pointer(&cursorLoc))))
	return normalizeError(err)
}

func (cur Coord) CursorIsAtLineEnd(size *Coord) bool {
	return cur.X == size.X
}

func (cur Coord) CursorIsAtLineBegin() bool {
	return cur.X == 0
}

func (c *Cursor) cursorMove(x int, y int, xIsAbs bool) error {
	handle := syscall.Handle(c.Out.Fd())

	var csbi consoleScreenBufferInfo
	if _, _, err := procGetConsoleScreenBufferInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&csbi))); normalizeError(err) != nil {
		log.Printf("cursorMove READ ERROR: %v", err)
		return err
	}

	var cursor Coord
	if xIsAbs {
		cursor.X = Short(x)
	} else {
		cursor.X = csbi.cursorPosition.X + Short(x)
	}
	cursor.Y = csbi.cursorPosition.Y + Short(y)

	xAbsLabel := ""
	if xIsAbs {
		xAbsLabel = " (abs)"
	}
	log.Printf("cursorMove X:%d%s Y:%d => %d, %d", x, xAbsLabel, y, cursor.X, cursor.Y)

	_, _, err := procSetConsoleCursorPosition.Call(uintptr(handle), uintptr(*(*int32)(unsafe.Pointer(&cursor))))
	if normalizeError(err) != nil {
		log.Printf("cursorMove WRITE ERROR: %v", err)
		return err
	}
	return nil
}

func (c *Cursor) NextLine(n int) error {
	return c.cursorMove(0, n, true)
}

func (c *Cursor) PreviousLine(n int) error {
	return c.cursorMove(0, -1*n, true)
}

func (c *Cursor) MoveNextLine(cur *Coord, terminalSize *Coord) error {
	if err := c.cursorMove(0, 1, true); err != nil {
		// moving to the next line will fail when at the bottom of the terminal viewport
		_, err = fmt.Fprint(c.Out, "\n")
		return err
	}
	return nil
}

func (c *Cursor) HorizontalAbsolute(x int) error {
	return c.cursorMove(0, 0, true)
}

func (c *Cursor) Show() error {
	handle := syscall.Handle(c.Out.Fd())

	var cci consoleCursorInfo
	if _, _, err := procGetConsoleCursorInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&cci))); normalizeError(err) != nil {
		return err
	}
	cci.visible = 1

	_, _, err := procSetConsoleCursorInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&cci)))
	return normalizeError(err)
}

func (c *Cursor) Hide() error {
	handle := syscall.Handle(c.Out.Fd())

	var cci consoleCursorInfo
	if _, _, err := procGetConsoleCursorInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&cci))); normalizeError(err) != nil {
		return err
	}
	cci.visible = 0

	_, _, err := procSetConsoleCursorInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&cci)))
	return normalizeError(err)
}

func (c *Cursor) Location(buf *bytes.Buffer) (*Coord, error) {
	handle := syscall.Handle(c.Out.Fd())

	var csbi consoleScreenBufferInfo
	if _, _, err := procGetConsoleScreenBufferInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&csbi))); normalizeError(err) != nil {
		return nil, err
	}

	return &csbi.cursorPosition, nil
}

func (c *Cursor) Size(buf *bytes.Buffer) (*Coord, error) {
	handle := syscall.Handle(c.Out.Fd())

	var csbi consoleScreenBufferInfo
	if _, _, err := procGetConsoleScreenBufferInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&csbi))); normalizeError(err) != nil {
		return nil, err
	}
	// windows' coordinate system begins at (0, 0)
	csbi.size.X--
	csbi.size.Y--
	return &csbi.size, nil
}
