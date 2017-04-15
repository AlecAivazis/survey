package survey

import (
	"errors"
	"io/ioutil"
	"strings"

	"github.com/chzyer/readline"
	ansi "github.com/k0kubun/go-ansi"

	"github.com/alecaivazis/survey/core"
)

// Select is a prompt that presents a list of various options to the user
// for them to select using the arrow keys and enter.
type Select struct {
	Message       string
	Options       []string
	Default       string
	SelectedIndex int
}

// the data available to the templates when processing
type SelectTemplateData struct {
	Select
	Answer string
}

const (
	SelectQuestionTemplate = `
{{- color "green+hb"}}? {{color "reset"}}
{{- color "default+hb"}}{{ .Message }} {{color "reset"}}
{{- if .Answer}}{{color "cyan"}}{{.Answer}}{{color "reset"}}{{end}}`
	// the template used to show the list of Selects
	SelectChoicesTemplate = `
{{- range $ix, $choice := .Options}}
  {{- if eq $ix $.Select.SelectedIndex}}{{color "cyan+b"}}> {{else}}{{color "default+hb"}}  {{end}}
  {{- $choice}}
  {{- color "reset"}}
{{end}}`
)

// OnChange is called on every keypress.
func (s *Select) OnChange(line []rune, pos int, key rune) (newLine []rune, newPos int, ok bool) {
	// if the user pressed the enter key
	if key == core.KeyEnter {
		return []rune(s.Options[s.SelectedIndex]), 0, true
		// if the user pressed the up arrow
	} else if key == core.KeyArrowUp && s.SelectedIndex > 0 {
		// decrement the selected index
		s.SelectedIndex--
		// if the user pressed down and there is room to move
	} else if key == core.KeyArrowDown && s.SelectedIndex < len(s.Options)-1 {
		// increment the selected index
		s.SelectedIndex++
	}

	// render the options
	s.render()

	// if we are not pressing ent
	return []rune(s.Options[s.SelectedIndex]), 0, true
}

func (s *Select) render() error {
	for range s.Options {
		ansi.CursorPreviousLine(1)
		ansi.EraseInLine(1)
	}

	// the formatted response
	out, err := RunTemplate(
		SelectChoicesTemplate,
		SelectTemplateData{Select: *s},
	)
	if err != nil {
		return err
	}

	// ask the question
	ansi.Println(strings.TrimRight(out, "\n"))
	// nothing went wrong
	return nil
}

func (s *Select) Prompt(rl *readline.Instance) (string, error) {
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
	s.SelectedIndex = sel

	// render the initial question
	out, err := RunTemplate(
		SelectQuestionTemplate,
		SelectTemplateData{Select: *s},
	)
	if err != nil {
		return "", err
	}

	// hide the cursor
	ansi.CursorHide()
	// ask the question
	ansi.Println(out)
	for range s.Options {
		ansi.Println()
	}
	// start waiting for input
	val, err := rl.Readline()
	// show the cursor when we're done
	ansi.CursorShow()

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

func (s *Select) Cleanup(rl *readline.Instance, val string) error {
	ansi.CursorPreviousLine(1)
	ansi.EraseInLine(1)
	for range s.Options {
		ansi.CursorPreviousLine(1)
		ansi.EraseInLine(1)
	}

	// execute the output summary template with the answer
	output, err := RunTemplate(
		SelectQuestionTemplate,
		SelectTemplateData{Select: *s, Answer: val},
	)
	if err != nil {
		return err
	}
	// render the summary
	ansi.Println(output)

	// nothing went wrong
	return nil
}
