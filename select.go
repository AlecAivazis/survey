package survey

import (
	"errors"
	"io/ioutil"
	"strings"

	"github.com/AlecAivazis/survey/core"
	"github.com/AlecAivazis/survey/terminal"
	"github.com/chzyer/readline"
)

// Select is a prompt that presents a list of various options to the user
// for them to select using the arrow keys and enter.
type Select struct {
	Message       string
	Options       []string
	Default       string
	selectedIndex int
}

// the data available to the templates when processing
type SelectTemplateData struct {
	Select
	SelectedIndex int
	Answer        string
}

const (
	SelectQuestionTemplate = `
{{- color "green+hb"}}? {{color "reset"}}
{{- color "default+hb"}}{{ .Message }} {{color "reset"}}
{{- if .Answer}}{{color "cyan"}}{{.Answer}}{{color "reset"}}{{end}}`
	// the template used to show the list of Selects
	SelectChoicesTemplate = `
{{- range $ix, $choice := .Options}}
  {{- if eq $ix $.SelectedIndex}}{{color "cyan+b"}}> {{else}}{{color "default+hb"}}  {{end}}
  {{- $choice}}
  {{- color "reset"}}
{{end}}`
)

// OnChange is called on every keypress.
func (s *Select) OnChange(line []rune, pos int, key rune) (newLine []rune, newPos int, ok bool) {
	// if the user pressed the enter key
	if key == terminal.KeyEnter {
		return []rune(s.Options[s.selectedIndex]), 0, true
		// if the user pressed the up arrow
	} else if key == terminal.KeyArrowUp && s.selectedIndex > 0 {
		// decrement the selected index
		s.selectedIndex--
		// if the user pressed down and there is room to move
	} else if key == terminal.KeyArrowDown && s.selectedIndex < len(s.Options)-1 {
		// increment the selected index
		s.selectedIndex++
	}

	// render the options
	s.render()

	// if we are not pressing ent
	return []rune(s.Options[s.selectedIndex]), 0, true
}

func (s *Select) render() error {
	for range s.Options {
		terminal.CursorPreviousLine(1)
		terminal.EraseLine(terminal.ERASE_LINE_ALL)
	}

	// the formatted response
	out, err := core.RunTemplate(
		SelectChoicesTemplate,
		SelectTemplateData{Select: *s, SelectedIndex: s.selectedIndex},
	)
	if err != nil {
		return err
	}

	// ask the question
	terminal.Println(strings.TrimRight(out, "\n"))
	// nothing went wrong
	return nil
}

func (s *Select) Prompt(rl *readline.Instance) (interface{}, error) {
	config := &readline.Config{
		Listener: s,
		Stdout:   ioutil.Discard,
	}
	rl.SetConfig(config)

	// if there are no options to render
	if len(s.Options) == 0 {
		// we failed
		return "", errors.New("please provide options to select from")
	}

	// start off with the first option selected
	sel := 0
	// if there is a default
	if s.Default != "" {
		// find the choice
		for i, opt := range s.Options {
			// if the option correponds to the default
			if opt == s.Default {
				// we found our initial value
				sel = i
				// stop looking
				break
			}
		}
	}
	// save the selected index
	s.selectedIndex = sel

	// render the initial question
	out, err := core.RunTemplate(
		SelectQuestionTemplate,
		SelectTemplateData{Select: *s, SelectedIndex: sel},
	)
	if err != nil {
		return "", err
	}

	// hide the cursor
	terminal.CursorHide()
	// ask the question
	terminal.Println(out)
	for range s.Options {
		terminal.Println()
	}
	// start waiting for input
	val, err := rl.Readline()
	// show the cursor when we're done
	terminal.CursorShow()

	//  if the value is empty (not sure why)
	if val == "" {
		// if there is a default value
		if s.Default != "" {
			// use the default value
			val = s.Default
		} else {
			// there is no default value so use the first
			val = s.Options[0]
		}
	}

	// return rl.Readline()
	return val, err
}

func (s *Select) Cleanup(rl *readline.Instance, val interface{}) error {
	terminal.CursorPreviousLine(1)
	terminal.EraseLine(terminal.ERASE_LINE_ALL)
	for range s.Options {
		terminal.CursorPreviousLine(1)
		terminal.EraseLine(terminal.ERASE_LINE_ALL)
	}

	// execute the output summary template with the answer
	output, err := core.RunTemplate(
		SelectQuestionTemplate,
		SelectTemplateData{Select: *s, Answer: val.(string)},
	)
	if err != nil {
		return err
	}
	// render the summary
	terminal.Println(output)

	// nothing went wrong
	return nil
}
