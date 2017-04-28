
////////////////////////////////////////////////////////////////////////////////
//                          DO NOT MODIFY THIS FILE!
//
//  This file was automatically generated via the command:
//
//      go run recorder/recorder.go -- selectThanInput.go
//
////////////////////////////////////////////////////////////////////////////////
package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"github.com/kr/pty"
)

func main() {
	fh, tty, _ := pty.Open()
	defer tty.Close()
	defer fh.Close()
	c := exec.Command("go", "run", "selectThanInput.go")
	c.Stdin = tty
	c.Stdout = tty
	c.Stderr = tty
	c.Start()
	buf := bufio.NewReaderSize(fh, 1024)

	expect("\x1b[2K\x1b[1;92m? \x1b[0m\x1b[1;99mChoose a color:\x1b[0m\r\n", buf)
	expect("\x1b[1;36m❯ red\x1b[0m\r\n", buf)
	expect("\x1b[1;99m  blue\x1b[0m\r\n", buf)
	expect("\x1b[1;99m  green\x1b[0m\r\n", buf)
	expect("\x1b[?25l\x1b[2K\x1b[1F\x1b[2K\x1b[1F\x1b[2K\x1b[1F\x1b[2K\x1b[1F\x1b[2K\x1b[1;92m? \x1b[0m\x1b[1;99mChoose a color:\x1b[0m\r\n", buf)
	expect("\x1b[1;36m❯ red\x1b[0m\r\n", buf)
	expect("\x1b[1;99m  blue\x1b[0m\r\n", buf)
	expect("\x1b[1;99m  green\x1b[0m\r\n", buf)
	fh.Write([]byte("\x1b"))
	fh.Write([]byte("["))
	fh.Write([]byte("B"))
	expect("\x1b[2K\x1b[1F\x1b[2K\x1b[1F\x1b[2K\x1b[1F\x1b[2K\x1b[1F\x1b[2K\x1b[1;92m? \x1b[0m\x1b[1;99mChoose a color:\x1b[0m\r\n", buf)
	expect("\x1b[1;99m  red\x1b[0m\r\n", buf)
	expect("\x1b[1;36m❯ blue\x1b[0m\r\n", buf)
	expect("\x1b[1;99m  green\x1b[0m\r\n", buf)
	fh.Write([]byte("\r"))
	expect("\x1b[?25h\x1b[2K\x1b[1F\x1b[2K\x1b[1F\x1b[2K\x1b[1F\x1b[2K\x1b[1F\x1b[2K\x1b[1;92m? \x1b[0m\x1b[1;99mChoose a color:\x1b[0m\x1b[36m blue\x1b[0m\r\n", buf)
	expect("\x1b[2K\x1b[1;92m? \x1b[0m\x1b[1;99mWhat is your name? \x1b[0m", buf)
	fh.Write([]byte("L"))
	expect("L", buf)
	fh.Write([]byte("a"))
	expect("a", buf)
	fh.Write([]byte("r"))
	expect("r", buf)
	fh.Write([]byte("r"))
	expect("r", buf)
	fh.Write([]byte("y"))
	expect("y", buf)
	fh.Write([]byte(" "))
	expect(" ", buf)
	fh.Write([]byte("W"))
	expect("W", buf)
	fh.Write([]byte("a"))
	expect("a", buf)
	fh.Write([]byte("l"))
	expect("l", buf)
	fh.Write([]byte("l"))
	expect("l", buf)
	fh.Write([]byte("\r"))
	expect("\r\n", buf)
	expect("\x1b[1A\x1b[2K\x1b[1;92m? \x1b[0m\x1b[1;99mWhat is your name? \x1b[0m\x1b[36mLarry Wall\x1b[0m\r\n", buf)
	expect("Larry Wall chose blue.\r\n", buf)

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
