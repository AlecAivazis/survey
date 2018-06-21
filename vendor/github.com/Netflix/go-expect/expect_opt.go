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
	"bytes"
	"regexp"
	"strings"
)

// ExpectOpt allows settings Expect options.
type ExpectOpt func(*ExpectOpts) error

// ExpectOpts provides additional options on Expect.
type ExpectOpts struct {
	Matchers []Matcher
	EOF      bool
}

// Matcher provides an interface for finding a match in content read from
// Console's tty.
type Matcher interface {
	Match(buf *bytes.Buffer) bool
}

// StringMatcher fulfills the Matcher interface to match strings against a given
// bytes.Buffer.
type StringMatcher struct {
	str string
}

func (sm *StringMatcher) Match(buf *bytes.Buffer) bool {
	if strings.Contains(buf.String(), sm.str) {
		return true
	}
	return false
}

// RegexpMatcher fulfills the Matcher interface to match Regexp against a given
// bytes.Buffer.
type RegexpMatcher struct {
	re *regexp.Regexp
}

func (rm *RegexpMatcher) Match(buf *bytes.Buffer) bool {
	return rm.re.Match(buf.Bytes())
}

// String adds an Expect condition to exit if the content read from Console's
// tty contains any of the given strings.
func String(strs ...string) ExpectOpt {
	return func(opts *ExpectOpts) error {
		for _, str := range strs {
			opts.Matchers = append(opts.Matchers, &StringMatcher{
				str: str,
			})
		}
		return nil
	}
}

// Regexp adds an Expect condition to exit if the content read from Console's
// tty matches the given Regexp.
func Regexp(res ...*regexp.Regexp) ExpectOpt {
	return func(opts *ExpectOpts) error {
		for _, re := range res {
			opts.Matchers = append(opts.Matchers, &RegexpMatcher{
				re: re,
			})
		}
		return nil
	}
}

// RegexpPattern adds an Expect condition to exit if the content read from
// Console's tty matches the given Regexp patterns. Expect returns an error if
// the patterns were unsuccessful in compiling the Regexp.
func RegexpPattern(ps ...string) ExpectOpt {
	return func(opts *ExpectOpts) error {
		var res []*regexp.Regexp
		for _, p := range ps {
			re, err := regexp.Compile(p)
			if err != nil {
				return err
			}
			res = append(res, re)
		}
		return Regexp(res...)(opts)
	}
}

// EOF adds an Expect condition to exit if io.EOF is returned from reading
// Console's tty.
func EOF(opts *ExpectOpts) error {
	opts.EOF = true
	return nil
}
