//go:build !windows
// +build !windows

// The terminal mode manipulation code is derived heavily from:
// https://github.com/golang/crypto/blob/master/ssh/terminal/util.go:
// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package terminal

import (
	"bufio"
	"bytes"
	"fmt"
	"syscall"
	"unsafe"
)

const (
	normalKeypad      = '['
	applicationKeypad = 'O'
)

type runeReaderState struct {
	term   syscall.Termios
	reader *bufio.Reader
	buf    *bytes.Buffer
}

func newRuneReaderState(input FileReader) runeReaderState {
	buf := new(bytes.Buffer)
	return runeReaderState{
		reader: bufio.NewReader(&BufferedReader{
			In:     input,
			Buffer: buf,
		}),
		buf: buf,
	}
}

func (rr *RuneReader) Buffer() *bytes.Buffer {
	return rr.state.buf
}

// For reading runes we just want to disable echo.
func (rr *RuneReader) SetTermMode() error {
	if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(rr.stdio.In.Fd()), ioctlReadTermios, uintptr(unsafe.Pointer(&rr.state.term)), 0, 0, 0); err != 0 {
		return err
	}

	newState := rr.state.term
	newState.Lflag &^= syscall.ECHO | syscall.ECHONL | syscall.ICANON | syscall.ISIG

	if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(rr.stdio.In.Fd()), ioctlWriteTermios, uintptr(unsafe.Pointer(&newState)), 0, 0, 0); err != 0 {
		return err
	}

	return nil
}

func (rr *RuneReader) RestoreTermMode() error {
	if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(rr.stdio.In.Fd()), ioctlWriteTermios, uintptr(unsafe.Pointer(&rr.state.term)), 0, 0, 0); err != 0 {
		return err
	}
	return nil
}

func (rr *RuneReader) ReadRune() (rune, int, error) {
	r, size, err := rr.state.reader.ReadRune()
	if err != nil {
		return r, size, err
	}

	// parse ^[ sequences to look for arrow keys
	if r == '\033' {
		if rr.state.reader.Buffered() == 0 {
			// no more characters so must be `Esc` key
			return KeyEscape, 1, nil
		}
		r, size, err = rr.state.reader.ReadRune()
		if err != nil {
			return r, size, err
		}

		switch r {
		case normalKeypad:
			return rr.readNormalKeypad()

		case applicationKeypad:
			return rr.readApplicationKeypad()

		default:
			return r, size, fmt.Errorf("unexpected escape sequence from terminal: %q", []rune{'\033', r})
		}
	}

	return r, size, err
}

func (rr *RuneReader) readNormalKeypad() (rune, int, error) {
	r, size, err := rr.state.reader.ReadRune()
	if err != nil {
		return r, size, err
	}

	switch r {
	// ESC [ D
	// https://vt100.net/docs/vt102-ug/appendixc.html#SC.2.2.3
	// https://vt100.net/docs/vt102-ug/chapter5.html#S5.5.2.14
	case 'D':
		return KeyArrowLeft, 1, nil

	// ESC [ C
	// https://vt100.net/docs/vt102-ug/appendixc.html#SC.2.2.3
	// https://vt100.net/docs/vt102-ug/chapter5.html#S5.5.2.14
	case 'C':
		return KeyArrowRight, 1, nil

	// ESC [ A
	// https://vt100.net/docs/vt102-ug/appendixc.html#SC.2.2.3
	// https://vt100.net/docs/vt102-ug/chapter5.html#S5.5.2.14
	case 'A':
		return KeyArrowUp, 1, nil

	// ESC [ B
	// https://vt100.net/docs/vt102-ug/appendixc.html#SC.2.2.3
	// https://vt100.net/docs/vt102-ug/chapter5.html#S5.5.2.14
	case 'B':
		return KeyArrowDown, 1, nil

	// Home Key
	// Cursor Position (Home) (CUP)
	// ESC [ H
	// https://vt100.net/docs/vt102-ug/appendixc.html#SC.2.2.9
	case 'H':
		return SpecialKeyHome, 1, nil

	// End Key
	// ESC [ F
	case 'F':
		return SpecialKeyEnd, 1, nil

	// Delete Key
	// Tabulation Clear (TBC)
	// ESC [ 3
	case '3':
		// discard the following '~' key from buffer
		_, _ = rr.state.reader.Discard(1)
		return SpecialKeyDelete, 1, nil

	default:
		// discard the following '~' key from buffer
		_, _ = rr.state.reader.Discard(1)
		return IgnoreKey, 1, nil
	}
}

func (rr *RuneReader) readApplicationKeypad() (rune, int, error) {
	r, size, err := rr.state.reader.ReadRune()
	if err != nil {
		return r, size, err
	}

	switch r {
	// ESC O D
	// https://vt100.net/docs/vt102-ug/appendixc.html#SC.2.2.3
	// https://vt100.net/docs/vt102-ug/chapter5.html#S5.5.2.14
	case 'D':
		return KeyArrowLeft, 1, nil

	// ESC O C
	// https://vt100.net/docs/vt102-ug/appendixc.html#SC.2.2.3
	// https://vt100.net/docs/vt102-ug/chapter5.html#S5.5.2.14
	case 'C':
		return KeyArrowRight, 1, nil

	// ESC O A
	// https://vt100.net/docs/vt102-ug/appendixc.html#SC.2.2.3
	// https://vt100.net/docs/vt102-ug/chapter5.html#S5.5.2.14
	case 'A':
		return KeyArrowUp, 1, nil

	// ESC O B
	// https://vt100.net/docs/vt102-ug/appendixc.html#SC.2.2.3
	// https://vt100.net/docs/vt102-ug/chapter5.html#S5.5.2.14
	case 'B':
		return KeyArrowDown, 1, nil

	// Home Key
	// Cursor Position (Home) (CUP)
	// ESC O H
	// https://vt100.net/docs/vt102-ug/appendixc.html#SC.2.2.9
	case 'H':
		return SpecialKeyHome, 1, nil

	// End Key
	// ESC O F
	case 'F':
		return SpecialKeyEnd, 1, nil

	default:
		// discard the following '~' key from buffer
		_, _ = rr.state.reader.Discard(1)
		return IgnoreKey, 1, nil
	}
}
