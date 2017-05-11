package core

import (
	"strings"

	"github.com/AlecAivazis/survey/terminal"
)

type Renderer struct {
	lineCount int
	mask      rune
}

func (r *Renderer) Render(tmpl string, data interface{}) error {
	// clean out current line in case tmpl didnt end in newline
	terminal.CursorHorizontalAbsolute(0)
	terminal.EraseLine(terminal.ERASE_LINE_ALL)
	// clean up what we left behind last time
	for i := 0; i < r.lineCount; i++ {
		terminal.CursorPreviousLine(1)
		terminal.EraseLine(terminal.ERASE_LINE_ALL)
	}

	// render the template summarizing the current state
	out, err := RunTemplate(tmpl, data)
	if err != nil {
		return err
	}

	// keep track of how many lines are printed so we can clean up later
	r.lineCount = strings.Count(out, "\n")

	// print the summary
	terminal.Print(out)

	// nothing went wrong
	return nil
}
