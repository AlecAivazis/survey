package survey

import (
	"errors"
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/core"
	"github.com/AlecAivazis/survey/terminal"
)

// Select is a prompt that presents a list of various options to the user
// for them to select using the arrow keys and enter.
type Select struct {
	core.Renderer
	Message       string
	Options       []string
	Default       string
	Help          string
	selectedIndex int
	useDefault    bool
	showingHelp   bool
}

// the data available to the templates when processing
type SelectTemplateData struct {
	Select
	SelectedIndex int
	Answer        string
	ShowAnswer    bool
	ShowHelp      bool
}

var SelectQuestionTemplate = `
{{- if .ShowHelp }}{{- color "cyan"}}{{ HelpIcon }} {{ .Help }}{{color "reset"}}{{"\n"}}{{end}}
{{- color "green+hb"}}{{ QuestionIcon }} {{color "reset"}}
{{- color "default+hb"}}{{ .Message }}{{color "reset"}}
{{- if .ShowAnswer}}{{color "cyan"}} {{.Answer}}{{color "reset"}}{{"\n"}}
{{- else}}
  {{- if and .Help (not .ShowHelp)}} {{color "cyan"}}[{{ HelpInputRune }} for help]{{color "reset"}}{{end}}
  {{- "\n"}}
  {{- range $ix, $choice := .Options}}
    {{- if eq $ix $.SelectedIndex}}{{color "cyan+b"}}{{ SelectFocusIcon }} {{else}}{{color "default+hb"}}  {{end}}
    {{- $choice}}
    {{- color "reset"}}{{"\n"}}
  {{- end}}
{{- end}}`

// OnChange is called on every keypress.
func (s *Select) OnChange(line []rune, pos int, key rune) (newLine []rune, newPos int, ok bool) {
	// if the user pressed the enter key
	if key == terminal.KeyEnter {
		return []rune(s.Options[s.selectedIndex]), 0, true
		// if the user pressed the up arrow
	} else if key == terminal.KeyArrowUp && s.selectedIndex > 0 {
		s.useDefault = false
		// decrement the selected index
		s.selectedIndex--
		// if the user pressed down and there is room to move
	} else if key == terminal.KeyArrowDown && s.selectedIndex < len(s.Options)-1 {
		s.useDefault = false
		// increment the selected index
		s.selectedIndex++
		// only show the help message if we have one
	} else if key == core.HelpInputRune && s.Help != "" {
		s.showingHelp = true
	}

	// render the options
	s.Render(
		SelectQuestionTemplate,
		SelectTemplateData{Select: *s, SelectedIndex: s.selectedIndex, ShowHelp: s.showingHelp},
	)

	// if we are not pressing ent
	return []rune(s.Options[s.selectedIndex]), 0, true
}

func (s *Select) Prompt() (interface{}, error) {
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

	// ask the question
	err := s.Render(
		SelectQuestionTemplate,
		SelectTemplateData{Select: *s, SelectedIndex: sel},
	)
	if err != nil {
		return "", err
	}

	// hide the cursor
	terminal.CursorHide()
	// show the cursor when we're done
	defer terminal.CursorShow()

	// by default, use the default value
	s.useDefault = true

	rr := terminal.NewRuneReader(os.Stdin)
	rr.SetTermMode()
	defer rr.RestoreTermMode()
	// start waiting for input
	for {
		r, _, err := rr.ReadRune()
		if err != nil {
			return "", err
		}
		if r == '\r' || r == '\n' {
			break
		}
		if r == terminal.KeyInterrupt {
			return "", fmt.Errorf("interrupt")
		}
		if r == terminal.KeyEndTransmission {
			break
		}
		s.OnChange(nil, 0, r)
	}

	var val string
	// if we are supposed to use the default value
	if s.useDefault {
		// if there is a default value
		if s.Default != "" {
			// use the default value
			val = s.Default
		} else {
			// there is no default value so use the first
			val = s.Options[0]
		}
		// otherwise the selected index points to the value
	} else {
		// the
		val = s.Options[s.selectedIndex]
	}

	return val, err
}

func (s *Select) Cleanup(val interface{}) error {
	return s.Render(
		SelectQuestionTemplate,
		SelectTemplateData{Select: *s, Answer: val.(string), ShowAnswer: true},
	)
}
