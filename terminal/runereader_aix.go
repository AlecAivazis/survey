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

    "golang.org/x/sys/unix"
)

const (
	normalKeypad      = '['
	applicationKeypad = 'O'
)

const (
        ioctlReadTermios  = 0x5401
        ioctlWriteTermios = 0x5402
)

type runeReaderState struct {
	term   syscall.Termios
	reader *bufio.Reader
	buf    *bytes.Buffer
}

func ioctl_aix(fd, request, argp uintptr) error {
        return unix.IoctlSetInt(int(fd), uint(request), int(argp))
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
    if err := ioctl_aix(uintptr(rr.stdio.In.Fd()), ioctlReadTermios, uintptr(unsafe.Pointer(&rr.state.term))); err != nil {
		return err
	}

	newState := rr.state.term
	newState.Lflag &^= syscall.ECHO | syscall.ECHONL | syscall.ICANON | syscall.ISIG

    if err := ioctl_aix(uintptr(rr.stdio.In.Fd()), ioctlWriteTermios, uintptr(unsafe.Pointer(&newState))); err != nil {
		return err
	}

	return nil
}

func (rr *RuneReader) RestoreTermMode() error {
    if err := ioctl_aix(uintptr(rr.stdio.In.Fd()), ioctlWriteTermios, uintptr(unsafe.Pointer(&rr.state.term))); err != nil {
		return err
	}
	return nil
}

// ReadRune Parse escape sequences such as ESC [ A for arrow keys.
// See https://vt100.net/docs/vt102-ug/appendixc.html
func (rr *RuneReader) ReadRune() (rune, int, error) {
	r, size, err := rr.state.reader.ReadRune()
	if err != nil {
		return r, size, err
	}

	if r != KeyEscape {
		return r, size, err
	}

	if rr.state.reader.Buffered() == 0 {
		// no more characters so must be `Esc` key
		return KeyEscape, 1, nil
	}

	r, size, err = rr.state.reader.ReadRune()
	if err != nil {
		return r, size, err
	}

	// ESC O ... or ESC [ ...?
	if r != normalKeypad && r != applicationKeypad {
		return r, size, fmt.Errorf("unexpected escape sequence from terminal: %q", []rune{KeyEscape, r})
	}

	keypad := r

	r, size, err = rr.state.reader.ReadRune()
	if err != nil {
		return r, size, err
	}

	switch r {
	case 'A': // ESC [ A or ESC O A
		return KeyArrowUp, 1, nil
	case 'B': // ESC [ B or ESC O B
		return KeyArrowDown, 1, nil
	case 'C': // ESC [ C or ESC O C
		return KeyArrowRight, 1, nil
	case 'D': // ESC [ D or ESC O D
		return KeyArrowLeft, 1, nil
	case 'F': // ESC [ F or ESC O F
		return SpecialKeyEnd, 1, nil
	case 'H': // ESC [ H or ESC O H
		return SpecialKeyHome, 1, nil
	case '3': // ESC [ 3
		if keypad == normalKeypad {
			// discard the following '~' key from buffer
			_, _ = rr.state.reader.Discard(1)
			return SpecialKeyDelete, 1, nil
		}
	}

	// discard the following '~' key from buffer
	_, _ = rr.state.reader.Discard(1)
	return IgnoreKey, 1, nil
}
