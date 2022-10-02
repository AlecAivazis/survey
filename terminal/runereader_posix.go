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
	"syscall"
	"unsafe"
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
	// Because we are clearing canonical mode, we need to ensure VMIN & VTIME are
	// set to the values we expect. This combination puts things in standard
	// "blocking read" mode (see termios(3)).
	newState.Cc[syscall.VMIN] = 1
	newState.Cc[syscall.VTIME] = 0

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
	return readRunePOSIX(rr.state.reader, true)
}
