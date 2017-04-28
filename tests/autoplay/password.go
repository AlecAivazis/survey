
////////////////////////////////////////////////////////////////////////////////
//                          DO NOT MODIFY THIS FILE!
//
//  This file was automatically generated via the command:
//
//      go run recorder/recorder.go -- password.go
//
////////////////////////////////////////////////////////////////////////////////
package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/kr/pty"
)

func main() {
	fh, tty, _ := pty.Open()
	defer tty.Close()
	defer fh.Close()
	c := exec.Command("go", "run", "password.go")
	c.Stdin = tty
	c.Stdout = tty
	c.Stderr = tty
	c.Start()
	buf := bufio.NewReaderSize(fh, 1024)

	expect("standard\r\n", buf)
	expect("\r\b", buf)
	expect(strings.Repeat("\r\b", 29), buf)
	expect("\x1b[J\x1b[1;92m? \x1b[0m\x1b[1;99mPlease type your password: \x1b[0m  \b", buf)
	fh.Write([]byte("f"))
	expect(strings.Repeat("\r\b", 30), buf)
	expect("\x1b[J\x1b[1;92m? \x1b[0m\x1b[1;99mPlease type your password: \x1b[0m *", buf)
	fh.Write([]byte("o"))
	expect(strings.Repeat("\r\b", 31), buf)
	expect("\x1b[J\x1b[1;92m? \x1b[0m\x1b[1;99mPlease type your password: \x1b[0m **", buf)
	fh.Write([]byte("o"))
	expect(strings.Repeat("\r\b", 32), buf)
	expect("\x1b[J\x1b[1;92m? \x1b[0m\x1b[1;99mPlease type your password: \x1b[0m ***", buf)
	fh.Write([]byte("b"))
	expect(strings.Repeat("\r\b", 33), buf)
	expect("\x1b[J\x1b[1;92m? \x1b[0m\x1b[1;99mPlease type your password: \x1b[0m ****", buf)
	fh.Write([]byte("a"))
	expect(strings.Repeat("\r\b", 34), buf)
	expect("\x1b[J\x1b[1;92m? \x1b[0m\x1b[1;99mPlease type your password: \x1b[0m *****", buf)
	fh.Write([]byte("r"))
	expect(strings.Repeat("\r\b", 35), buf)
	expect("\x1b[J\x1b[1;92m? \x1b[0m\x1b[1;99mPlease type your password: \x1b[0m ******", buf)
	fh.Write([]byte("\r"))
	expect(strings.Repeat("\r\b", 36), buf)
	expect("\x1b[J\x1b[1;92m? \x1b[0m\x1b[1;99mPlease type your password: \x1b[0m ******\r\b", buf)
	expect(strings.Repeat("\r\b", 35), buf)
	expect("\x1b[J\x1b[1;92m? \x1b[0m\x1b[1;99mPlease type your password: \x1b[0m ******\r\n", buf)
	expect("\x1b[JAnswered foobar.\r\n", buf)
	expect("---------------------\r\n", buf)
	expect("please make sure paste works\r\n", buf)
	expect("\r\b", buf)
	expect(strings.Repeat("\r\b", 30), buf)
	expect("\x1b[J\x1b[1;92m? \x1b[0m\x1b[1;99mPlease paste your password: \x1b[0m  \b", buf)
	fh.Write([]byte("b"))
	fh.Write([]byte("o"))
	fh.Write([]byte("o"))
	fh.Write([]byte("f"))
	fh.Write([]byte("a"))
	fh.Write([]byte("r"))
	expect(strings.Repeat("\r\b", 31), buf)
	expect("\x1b[J\x1b[1;92m? \x1b[0m\x1b[1;99mPlease paste your password: \x1b[0m *\r\b", buf)
	expect(strings.Repeat("\r\b", 31), buf)
	expect("\x1b[J\x1b[1;92m? \x1b[0m\x1b[1;99mPlease paste your password: \x1b[0m **\r\b", buf)
	expect(strings.Repeat("\r\b", 32), buf)
	expect("\x1b[J\x1b[1;92m? \x1b[0m\x1b[1;99mPlease paste your password: \x1b[0m ***\r\b", buf)
	expect(strings.Repeat("\r\b", 33), buf)
	expect("\x1b[J\x1b[1;92m? \x1b[0m\x1b[1;99mPlease paste your password: \x1b[0m ****\r\b", buf)
	expect(strings.Repeat("\r\b", 34), buf)
	expect("\x1b[J\x1b[1;92m? \x1b[0m\x1b[1;99mPlease paste your password: \x1b[0m *****\r\b", buf)
	expect(strings.Repeat("\r\b", 35), buf)
	expect("\x1b[J\x1b[1;92m? \x1b[0m\x1b[1;99mPlease paste your password: \x1b[0m ******", buf)
	fh.Write([]byte("\r"))
	expect(strings.Repeat("\r\b", 37), buf)
	expect("\x1b[J\x1b[1;92m? \x1b[0m\x1b[1;99mPlease paste your password: \x1b[0m ******\r\b", buf)
	expect(strings.Repeat("\r\b", 36), buf)
	expect("\x1b[J\x1b[1;92m? \x1b[0m\x1b[1;99mPlease paste your password: \x1b[0m ******\r\n", buf)
	expect("\x1b[JAnswered boofar.\r\n", buf)
	expect("---------------------\r\n", buf)

	c.Wait()
	tty.Close()
	fh.Close()
}

func expect(expected string, buf *bufio.Reader) {
	sofar := []rune{}
	for _, r := range expected {
		got, _, _ := buf.ReadRune()
		sofar = append(sofar, got)
		if got != r {
			fmt.Fprintln(os.Stderr)
			fmt.Fprintf(os.Stderr, "Expected: %q\n", expected[:len(sofar)])
			fmt.Fprintf(os.Stderr, "Got:      %q\n", string(sofar))
			panic(fmt.Errorf("Unexpected Rune %q, Expected %q\n", got, r))
		} else {
			fmt.Printf("%c", r)
		}
	}
}
