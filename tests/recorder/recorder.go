package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/AlecAivazis/survey/core"
	"github.com/kr/pty"
)

type Data struct {
	Type  string
	Bytes []byte
}

type DataRecorder struct {
	Data    []Data
	Imports map[string]bool
}

func NewDataRecorder() *DataRecorder {
	return &DataRecorder{
		Data:    []Data{},
		Imports: map[string]bool{},
	}
}

func (w *DataRecorder) Write(p []byte) (n int, err error) {
	tmp := make([]byte, len(p))
	copy(tmp, p)
	w.Data = append(w.Data, Data{"OUTPUT", tmp})
	return os.Stdout.Write(p)
}

// ProcessedData will first merge all the sequential OUTPUT chunks together then will
// split then on newlines and the \r\b sequnce used by readline.RuneBuffer. We do this
// so that we can create autoplay code through repeated runs of the recorder.  Variance
// will typically occure to buffering within the pty output.
func (w *DataRecorder) ProcessedData() []Data {
	// first we merge all the sequential OUTPUT chunks
	merged := []Data{}
	for _, chunk := range w.Data {
		if len(merged) == 0 {
			merged = append(merged, chunk)
			continue
		}
		cursor := len(merged) - 1
		prevChunk := merged[cursor]
		if chunk.Type == "OUTPUT" && prevChunk.Type == "OUTPUT" {
			prevChunk.Bytes = append(prevChunk.Bytes, chunk.Bytes...)
			merged[cursor] = prevChunk
		} else {
			merged = append(merged, chunk)
		}
	}

	// Now we split up the OUTPUT blocks baked on
	// sequential crbs chunks and all linebreaks to
	// make the generated code more readable
	processed := []Data{}
	crbs := []byte("\r\b")
	crnl := []byte("\r\n")
	for _, data := range merged {
		if data.Type == "INPUT" {
			processed = append(processed, data)
			continue
		}
		chunks := bytes.SplitAfter(data.Bytes, crbs)
		if len(chunks) > 1 {
			tally := 0
			for _, chunk := range chunks {
				if bytes.Equal(chunk, crbs) {
					tally++
					continue
				} else {
					if tally > 0 {
						processed = append(processed, Data{"OUTCODE", []byte(fmt.Sprintf(`strings.Repeat(%q, %d)`, crbs, tally))})
						w.Imports["strings"] = true
						tally = 0
					}
					for _, line := range bytes.SplitAfter(chunk, crnl) {
						if len(line) > 0 {
							processed = append(processed, Data{"OUTPUT", line})
						}
					}
				}
			}
		} else {
			for _, line := range bytes.SplitAfter(data.Bytes, crnl) {
				if len(line) > 0 {
					processed = append(processed, Data{"OUTPUT", line})
				}
			}
		}
	}
	return processed
}

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: go run record/recorder.go -- <test>.go\n")
		os.Exit(1)
	}
	test := os.Args[2]
	fh, tty, _ := pty.Open()
	defer tty.Close()
	defer fh.Close()
	rec := NewDataRecorder()
	c := exec.Command("go", "run", test)
	c.Stdin = tty
	c.Stdout = tty
	c.Stderr = tty
	// start streaming the pty output to our recorder
	go func() {
		io.Copy(rec, fh)
	}()

	// put stdin in raw mode so the terminal will not double-echo and
	// mess up what we are testing
	state, _ := terminal.MakeRaw(syscall.Stdin)
	defer terminal.Restore(syscall.Stdin, state)

	// create a input Buffer so we can read a rune at a time
	inputBuf := bufio.NewReaderSize(os.Stdin, 1024)
	go func() {
		for {
			r, _, err := inputBuf.ReadRune()
			if err != nil {
				break
			}
			b := []byte(string(r))
			rec.Data = append(rec.Data, Data{"INPUT", b})
			fh.Write(b)
		}
	}()
	c.Run()
	tty.Close()
	fh.Close()

	// put Stdin back in normal state
	terminal.Restore(syscall.Stdin, state)

	// generate the autoplay code via template
	results, err := core.RunTemplate(
		DriverTemplate, struct {
			IO      []Data
			File    string
			Imports map[string]bool
		}{
			rec.ProcessedData(),
			test,
			rec.Imports,
		},
	)
	if err != nil {
		panic(err)
	}

	out, err := os.Create(fmt.Sprintf("autoplay/%s", test))
	if err != nil {
		panic(err)
	}
	fmt.Fprint(out, results)
	out.Close()
}

var DriverTemplate = `
////////////////////////////////////////////////////////////////////////////////
//                          DO NOT MODIFY THIS FILE!
//
//  This file was automatically generated via the command:
//
//      go run recorder/recorder.go -- {{.File}}
//
////////////////////////////////////////////////////////////////////////////////
package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"{{range $k, $v :=.Imports}}
	{{ printf "%q" $k }}
{{end}}
	"github.com/kr/pty"
)

func main() {
	fh, tty, _ := pty.Open()
	defer tty.Close()
	defer fh.Close()
	c := exec.Command("go", "run", "{{.File}}")
	c.Stdin = tty
	c.Stdout = tty
	c.Stderr = tty
	c.Start()
	buf := bufio.NewReaderSize(fh, 1024)
{{range .IO}}{{if eq .Type "INPUT" }}
	fh.Write([]byte({{printf "%q" .Bytes}}))
{{- else if eq .Type "OUTCODE" }}
	expect({{printf "%s" .Bytes}}, buf)
{{- else }}
	expect({{printf "%q" .Bytes}}, buf)
{{- end}}{{end}}

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
`
