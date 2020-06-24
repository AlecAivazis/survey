package survey

import (
	"bufio"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/AlecAivazis/survey/v2/core"
	"github.com/AlecAivazis/survey/v2/terminal"
	goterm "golang.org/x/crypto/ssh/terminal"
)

type Renderer struct {
	stdio          terminal.Stdio
	lineCount      int
	errorLineCount int
}

type ErrorTemplateData struct {
	Error error
	Icon  Icon
}

var ErrorTemplate = `{{color .Icon.Format }}{{ .Icon.Text }} Sorry, your reply was invalid: {{ .Error.Error }}{{color "reset"}}
`

func (r *Renderer) WithStdio(stdio terminal.Stdio) {
	r.stdio = stdio
}

func (r *Renderer) Stdio() terminal.Stdio {
	return r.stdio
}

func (r *Renderer) NewRuneReader() *terminal.RuneReader {
	return terminal.NewRuneReader(r.stdio)
}

func (r *Renderer) NewCursor() *terminal.Cursor {
	return &terminal.Cursor{
		In:  r.stdio.In,
		Out: r.stdio.Out,
	}
}

func (r *Renderer) Error(config *PromptConfig, invalid error) error {
	// since errors are printed on top we need to reset the prompt
	// as well as any previous error print
	r.resetPrompt(r.lineCount + r.errorLineCount)

	// we just cleared the prompt lines
	r.lineCount = 0
	userOut, layoutOut, err := core.RunTemplate(ErrorTemplate, &ErrorTemplateData{
		Error: invalid,
		Icon:  config.Icons.Error,
	})
	if err != nil {
		return err
	}
	// keep track of how many lines are printed so we can clean up later
	r.errorLineCount = r.countLines(layoutOut)

	// send the message to the user
	fmt.Fprint(terminal.NewAnsiStdout(r.stdio.Out), userOut)
	return nil
}

func (r *Renderer) resetPrompt(lines int) {
	// clean out current line in case tmpl didnt end in newline
	cursor := r.NewCursor()
	cursor.HorizontalAbsolute(0)
	terminal.EraseLine(r.stdio.Out, terminal.ERASE_LINE_ALL)
	// clean up what we left behind last time
	for i := 0; i < lines; i++ {
		cursor.PreviousLine(1)
		terminal.EraseLine(r.stdio.Out, terminal.ERASE_LINE_ALL)
	}
}

func (r *Renderer) Render(tmpl string, data interface{}) error {
	r.resetPrompt(r.lineCount)
	// render the template summarizing the current state
	userOut, layoutOut, err := core.RunTemplate(tmpl, data)
	if err != nil {
		return err
	}

	// keep track of how many lines are printed so we can clean up later
	r.lineCount = r.countLines(layoutOut)

	// print the summary
	fmt.Fprint(terminal.NewAnsiStdout(r.stdio.Out), userOut)

	// nothing went wrong
	return nil
}

func (r *Renderer) termWidth() (int, error) {
	fd := int(r.stdio.Out.Fd())
	termWidth, _, err := goterm.GetSize(fd)
	return termWidth, err
}

// countLines will return the count of `\n` with the addition of any
// lines that have wrapped due to narrow terminal width
func (r *Renderer) countLines(out string) int {
	w, err := r.termWidth()
	if err != nil || w == 0 {
		// if we got an error due to terminal.GetSize not being supported
		// on current platform then just assume a very wide terminal
		w = 10000
	}

	count := 0
	s := bufio.NewScanner(strings.NewReader(out))
	for s.Scan() {
		line := s.Text()
		count += 1 + int(utf8.RuneCountInString(line)/w)
	}

	// if the prompt doesn't end on a newline, subtract off one '\n'
	if count != 0 && out[len(out)-1] != '\n' {
		count -= 1
	}

	return count
}
