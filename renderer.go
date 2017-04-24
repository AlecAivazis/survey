package survey

import (
	"io/ioutil"
	"strings"

	"github.com/AlecAivazis/survey/core"
	"github.com/AlecAivazis/survey/terminal"
	"github.com/chzyer/readline"
)

// basic readline config that does not put terminal in Raw mode
// so that basic editing works (backspace, arrow keys etc).  This
// is used by all text-input prompt types
var simpleReadlineConfig = &readline.Config{
	Stdout: ioutil.Discard,
	FuncMakeRaw: func() error {
		return nil
	},
	FuncExitRaw: func() error {
		return nil
	},
	HistoryLimit: -1,
	EOFPrompt:    "\n",
}

type renderer struct {
	lineCount int
	mask      rune
}

func (r *renderer) render(tmpl string, data interface{}) error {
	// clean out current line in case tmpl didnt end in newline
	terminal.EraseLine(terminal.ERASE_LINE_ALL)
	// clean up what we left behind last time
	for i := 0; i < r.lineCount; i++ {
		terminal.CursorPreviousLine(1)
		terminal.EraseLine(terminal.ERASE_LINE_ALL)
	}

	// render the template summarizing the current state
	out, err := core.RunTemplate(tmpl, data)
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
