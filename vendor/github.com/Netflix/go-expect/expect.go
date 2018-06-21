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
	"bytes"
	"io"
	"unicode/utf8"
)

// ExpectString reads from Console's tty until the provided string is read or
// an error occurs, and returns the buffer read by Console.
func (c *Console) ExpectString(s string) (string, error) {
	return c.Expect(String(s))
}

// ExpectEOF reads from Console's tty until EOF or an error occurs, and returns
// the buffer read by Console.
func (c *Console) ExpectEOF() (string, error) {
	return c.Expect(EOF)
}

// Expect reads from Console's tty until a condition specified from opts is
// encountered or an error occurs, and returns the buffer read by console.
// No extra bytes are read once a condition is met, so if a program isn't
// expecting input yet, it will be blocked. Sends are queued up in tty's
// internal buffer so that the next Expect will read the remaining bytes (i.e.
// rest of prompt) as well as its conditions.
func (c *Console) Expect(opts ...ExpectOpt) (string, error) {
	var options ExpectOpts
	for _, opt := range opts {
		if err := opt(&options); err != nil {
			return "", err
		}
	}

	buf := new(bytes.Buffer)
	writer := io.MultiWriter(append(c.opts.Stdouts, buf)...)
	runeWriter := bufio.NewWriterSize(writer, utf8.UTFMax)

	for {
		r, _, err := c.runeReader.ReadRune()
		if err != nil {
			if options.EOF && err == io.EOF {
				break
			}
			return buf.String(), err
		}

		c.Logf("expect read: %q", string(r))
		_, err = runeWriter.WriteRune(r)
		if err != nil {
			return buf.String(), err
		}

		// Immediately flush rune to the underlying writers.
		err = runeWriter.Flush()
		if err != nil {
			return buf.String(), err
		}

		matchFound := false
		for _, matcher := range options.Matchers {
			if matcher.Match(buf) {
				matchFound = true
				break
			}
		}

		if matchFound {
			break
		}
	}

	return buf.String(), nil
}
