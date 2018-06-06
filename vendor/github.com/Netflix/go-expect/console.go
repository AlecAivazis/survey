// Copyright 2018 Netflix, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package expect

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"unicode/utf8"

	"github.com/kr/pty"
)

// Console is an interface to automate input and output for interactive
// applications. Console can block until a specified output is received and send
// input back on it's tty. Console can also multiplex other sources of input
// and multiplex its output to other writers.
type Console struct {
	opts       ConsoleOpts
	ptm        *os.File
	pts        *os.File
	runeReader *bufio.Reader
	closers    []io.Closer
}

// ConsoleOpt allows setting Console options.
type ConsoleOpt func(*ConsoleOpts) error

// ConsoleOpts provides additional options on creating a Console.
type ConsoleOpts struct {
	Logger  *log.Logger
	Stdins  []io.Reader
	Stdouts []io.Writer
	Closers []io.Closer
}

// WithStdout adds writers that Console duplicates writes to, similar to the
// Unix tee(1) command.
//
// Each write is written to each listed writer, one at a time. Console is the
// last writer, writing to it's internal buffer for matching expects.
// If a listed writer returns an error, that overall write operation stops and
// returns the error; it does not continue down the list.
func WithStdout(writers ...io.Writer) ConsoleOpt {
	return func(opts *ConsoleOpts) error {
		opts.Stdouts = append(opts.Stdouts, writers...)
		return nil
	}
}

// WithStdin adds readers that bytes read are written to Console's  tty. If a
// listed reader returns an error, that reader will not be continued to read.
func WithStdin(readers ...io.Reader) ConsoleOpt {
	return func(opts *ConsoleOpts) error {
		opts.Stdins = append(opts.Stdins, readers...)
		return nil
	}
}

// WithCloser adds closers that are closed in order when Console is closed.
func WithCloser(closer ...io.Closer) ConsoleOpt {
	return func(opts *ConsoleOpts) error {
		opts.Closers = append(opts.Closers, closer...)
		return nil
	}
}

// WithLogger adds a logger for Console to log debugging information to. By
// default Console will discard logs.
func WithLogger(logger *log.Logger) ConsoleOpt {
	return func(opts *ConsoleOpts) error {
		opts.Logger = logger
		return nil
	}
}

// NewConsole returns a new Console with the given options.
func NewConsole(opts ...ConsoleOpt) (*Console, error) {
	options := ConsoleOpts{
		Logger: log.New(ioutil.Discard, "", 0),
	}

	for _, opt := range opts {
		if err := opt(&options); err != nil {
			return nil, err
		}
	}

	ptm, pts, err := pty.Open()
	if err != nil {
		return nil, err
	}
	closers := append(options.Closers, pts, ptm)

	c := &Console{
		opts:       options,
		ptm:        ptm,
		pts:        pts,
		runeReader: bufio.NewReaderSize(ptm, utf8.UTFMax),
		closers:    closers,
	}

	for _, r := range options.Stdins {
		go func(r io.Reader) {
			_, err := io.Copy(c, r)
			if err != nil {
				c.Logf("failed to copy stdin: %s", err)
			}
		}(r)
	}

	return c, nil
}

// Tty returns Console's pts (slave part of a pty). A pseudoterminal, or pty is
// a pair of psuedo-devices, one of which, the slave, emulates a real text
// terminal device.
func (c *Console) Tty() *os.File {
	return c.pts
}

// Read reads bytes b from Console's tty.
func (c *Console) Read(b []byte) (int, error) {
	return c.ptm.Read(b)
}

// Write writes bytes b to Console's tty.
func (c *Console) Write(b []byte) (int, error) {
	c.Logf("console write: %q", b)
	return c.ptm.Write(b)
}

// Close closes Console's tty. Calling Close will unblock Expect and ExpectEOF.
func (c *Console) Close() error {
	for _, fd := range c.closers {
		err := fd.Close()
		if err != nil {
			c.Logf("failed to close: %s", err)
		}
	}
	return nil
}

// Send writes string s to Console's tty.
func (c *Console) Send(s string) (int, error) {
	c.Logf("console send: %q", s)
	return c.ptm.WriteString(s)
}

// SendLine writes string s to Console's tty with a trailing newline.
func (c *Console) SendLine(s string) (int, error) {
	return c.Send(fmt.Sprintf("%s\n", s))
}

// Log prints to Console's logger.
// Arguments are handled in the manner of fmt.Print.
func (c *Console) Log(v ...interface{}) {
	c.opts.Logger.Print(v...)
}

// Logf prints to Console's logger.
// Arguments are handled in the manner of fmt.Printf.
func (c *Console) Logf(format string, v ...interface{}) {
	c.opts.Logger.Printf(format, v...)
}
