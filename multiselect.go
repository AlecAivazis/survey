package survey

import (
	"errors"
	"io/ioutil"
	"strings"

	"github.com/AlecAivazis/survey/core"
	"github.com/AlecAivazis/survey/terminal"
	"github.com/chzyer/readline"
)

// MultiSelect is a prompt that presents a list of various options to the user
// for them to select using the arrow keys and enter.
type MultiSelect struct {
	core.Renderer
	Message       string
	Options       []string
	Default       []string
	Help          string
	selectedIndex int
	checked       map[int]bool
	showingHelp   bool
}

// data available to the templates when processing
type MultiSelectTemplateData struct {
	MultiSelect
	Answer        string
	Checked       map[int]bool
	SelectedIndex int
	ShowHelp      bool
}

var MultiSelectQuestionTemplate = `
{{- if .ShowHelp }}{{- color "cyan"}}{{ HelpIcon }} {{ .Help }}{{color "reset"}}{{"\n"}}{{end}}
{{- color "green+hb"}}{{ QuestionIcon }} {{color "reset"}}
{{- color "default+hb"}}{{ .Message }}{{color "reset"}}
{{- if .Answer}}{{color "cyan"}} {{.Answer}}{{color "reset"}}{{"\n"}}
{{- else }}
  {{- if and .Help (not .ShowHelp)}} {{color "cyan"}}[{{ HelpInputRune }} for help]{{color "reset"}}{{end}}
  {{- "\n"}}
  {{- range $ix, $option := .Options}}
    {{- if eq $ix $.SelectedIndex}}{{color "cyan"}}{{ SelectFocusIcon }}{{color "reset"}}{{else}} {{end}}
    {{- if index $.Checked $ix}}{{color "green"}} {{ MarkedOptionIcon }} {{else}}{{color "default+hb"}} {{ UnmarkedOptionIcon }} {{end}}
    {{- color "reset"}}
    {{- " "}}{{$option}}{{"\n"}}
  {{- end}}
{{- end}}`

// OnChange is called on every keypress.
func (m *MultiSelect) OnChange(line []rune, pos int, key rune) (newLine []rune, newPos int, ok bool) {
	if key == terminal.KeyEnter {
		// just pass on the current value
		return line, 0, true
	} else if key == terminal.KeyArrowUp && m.selectedIndex > 0 {
		// decrement the selected index
		m.selectedIndex--
	} else if key == terminal.KeyArrowDown && m.selectedIndex < len(m.Options)-1 {
		// if the user pressed down and there is room to move
		// increment the selected index
		m.selectedIndex++
	} else if key == terminal.KeySpace {
		if old, ok := m.checked[m.selectedIndex]; !ok {
			// otherwise just invert the current value
			m.checked[m.selectedIndex] = true
		} else {
			// otherwise just invert the current value
			m.checked[m.selectedIndex] = !old
		}
	} else if key == core.HelpInputRune {
		m.showingHelp = true
	}

	// render the options
	m.Render(
		MultiSelectQuestionTemplate,
		MultiSelectTemplateData{
			MultiSelect:   *m,
			SelectedIndex: m.selectedIndex,
			Checked:       m.checked,
			ShowHelp:      m.showingHelp,
		},
	)

	// if we are not pressing ent
	return line, 0, true
}

func (m *MultiSelect) Prompt(rl *readline.Instance) (interface{}, error) {
	// the readline config
	config := &readline.Config{
		Listener: m,
		Stdout:   ioutil.Discard,
	}
	rl.SetConfig(config)

	// compute the default state
	m.checked = make(map[int]bool)
	// if there is a default
	if len(m.Default) > 0 {
		for _, dflt := range m.Default {
			for i, opt := range m.Options {
				// if the option correponds to the default
				if opt == dflt {
					// we found our initial value
					m.checked[i] = true
					// stop looking
					break
				}
			}
		}
	}

	// if there are no options to render
	if len(m.Options) == 0 {
		// we failed
		return "", errors.New("please provide options to select from")
	}

	// hide the cursor
	terminal.CursorHide()
	// show the cursor when we're done
	defer terminal.CursorShow()

	// ask the question
	err := m.Render(
		MultiSelectQuestionTemplate,
		MultiSelectTemplateData{
			MultiSelect:   *m,
			SelectedIndex: m.selectedIndex,
			Checked:       m.checked,
		},
	)
	if err != nil {
		return "", err
	}

	// start waiting for input
	_, err = rl.Readline()
	// if something went wrong
	if err != nil {
		return "", err
	}

	answers := []string{}
	for ix, option := range m.Options {
		if val, ok := m.checked[ix]; ok && val {
			answers = append(answers, option)
		}
	}

	return answers, nil
}

// Cleanup removes the options section, and renders the ask like a normal question.
func (m *MultiSelect) Cleanup(rl *readline.Instance, val interface{}) error {
	// execute the output summary template with the answer
	return m.Render(
		MultiSelectQuestionTemplate,
		MultiSelectTemplateData{
			MultiSelect:   *m,
			SelectedIndex: m.selectedIndex,
			Checked:       m.checked,
			Answer:        strings.Join(val.([]string), ", "),
		},
	)
}
