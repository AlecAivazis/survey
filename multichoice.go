package survey

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"strings"

	"github.com/alecaivazis/survey/core"
	"github.com/alecaivazis/survey/terminal"
	"github.com/chzyer/readline"
)

// MultiChoice is a prompt that presents a list of various options to the user
// for them to select using the arrow keys and enter.
type MultiChoice struct {
	Message       string
	Options       []string
	Defaults      []string
	Answer        *[]string
	selectedIndex int
	checked       map[int]bool
}

// data available to the templates when processing
type MultiChoiceTemplateData struct {
	MultiChoice
	Answer        []string
	Checked       map[int]bool
	SelectedIndex int
}

var MultiChoiceQuestionTemplate = `
{{- color "green+hb"}}? {{color "reset"}}
{{- color "default+hb"}}{{ .Message }} {{color "reset"}}
{{- if .Answer}}{{color "cyan"}}{{.Answer | printf "%q"}}{{color "reset"}}{{end}}`

var MultiChoiceOptionsTemplate = `
{{- range $ix, $option := .Options}}
  {{- if eq $ix $.SelectedIndex}}{{color "cyan"}}❯{{color "reset"}}{{else}} {{end}}
  {{- if index $.Checked $ix}}{{color "green"}} ◉ {{else}}{{color "default+hb"}} ◯ {{end}}
  {{- color "reset"}}
  {{- " "}}{{$option}}
{{end}}`

// OnChange is called on every keypress.
func (m *MultiChoice) OnChange(line []rune, pos int, key rune) (newLine []rune, newPos int, ok bool) {
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

	}

	// render the options
	m.render()

	// if we are not pressing ent
	return line, 0, true
}

func (m *MultiChoice) render() error {
	// clean up what we left behind last time
	for range m.Options {
		terminal.CursorPreviousLine(1)
		terminal.EraseInLine(1)
	}

	// render the template summarizing the current state
	out, err := core.RunTemplate(
		MultiChoiceOptionsTemplate,
		MultiChoiceTemplateData{
			MultiChoice:   *m,
			SelectedIndex: m.selectedIndex,
			Checked:       m.checked,
		},
	)
	if err != nil {
		return err
	}

	// print the summary
	terminal.Println(strings.TrimRight(out, "\n"))

	// nothing went wrong
	return nil
}

func (m *MultiChoice) Prompt(rl *readline.Instance) (string, error) {
	// if the user didn't pass an answer reference
	if m.Answer == nil {
		// build one
		answer := []string{}
		m.Answer = &answer
	}

	// the readline config
	config := &readline.Config{
		Listener: m,
		Stdout:   ioutil.Discard,
	}
	rl.SetConfig(config)

	// compute the default state
	m.checked = make(map[int]bool)
	// if there is a default
	if len(m.Defaults) > 0 {
		for _, dflt := range m.Defaults {
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
	// generate the template for the current state of the prompt
	out, err := core.RunTemplate(
		MultiChoiceQuestionTemplate,
		MultiChoiceTemplateData{
			MultiChoice:   *m,
			SelectedIndex: m.selectedIndex,
			Checked:       m.checked,
		},
	)
	if err != nil {
		return "", err
	}
	// hide the cursor
	terminal.CursorHide()
	// ask the question
	terminal.Println(out)
	for range m.Options {
		terminal.Println()
	}

	// start waiting for input
	_, err = rl.Readline()
	// if something went wrong
	if err != nil {
		return "", err
	}
	// show the cursor when we're done
	terminal.CursorShow()

	answers := []string{}
	for ix, option := range m.Options {
		if val, ok := m.checked[ix]; ok && val {
			answers = append(answers, option)
		}
	}
	*m.Answer = answers

	// nothing went wrong
	return m.value()
}

func (m *MultiChoice) value() (string, error) {
	answers := []string{}
	for ix, option := range m.Options {
		if val, ok := m.checked[ix]; ok && val {
			answers = append(answers, option)
		}
	}
	// return the selected option
	js, err := json.Marshal(answers)
	if err != nil {
		return "", err
	}
	return string(js), nil
}

// Cleanup removes the options section, and renders the ask like a normal question.
func (m *MultiChoice) Cleanup(rl *readline.Instance, val string) error {
	terminal.CursorPreviousLine(1)
	terminal.EraseInLine(1)
	for range m.Options {
		terminal.CursorPreviousLine(1)
		terminal.EraseInLine(1)
	}

	// parse the value into a list of strings
	var value []string
	json.Unmarshal([]byte(val), &value)

	// execute the output summary template with the answer
	output, err := core.RunTemplate(
		MultiChoiceQuestionTemplate,
		MultiChoiceTemplateData{
			MultiChoice:   *m,
			SelectedIndex: m.selectedIndex,
			Checked:       m.checked,
			Answer:        value,
		},
	)
	if err != nil {
		return err
	}
	// render the summary
	terminal.Println(output)

	// nothing went wrong
	return nil
}
