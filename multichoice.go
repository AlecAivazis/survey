package survey

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/alecaivazis/survey/core"
	"github.com/chzyer/readline"
	ansi "github.com/k0kubun/go-ansi"
)

// MultiChoice is a prompt that presents a list of various options to the user
// for them to select using the arrow keys and enter.
type MultiChoice struct {
	Message       string
	Options       []string
	Defaults      []string
	Answer        *[]string
	SelectedIndex int
	Checked       map[int]bool
}

// data available to the templates when processing
type multiChoiceTemplateData struct {
	MultiChoice
	Answer        []string
	Checked       map[int]bool
	SelectedIndex int
}

var multiChoiceQuestionTemplate = `
{{- color "green+hb"}}? {{color "reset"}}
{{- color "default+hb"}}{{ .Message }} {{color "reset"}}
{{- if .Answer}}{{color "cyan"}}{{.Answer | printf "%q"}}{{color "reset"}}{{end}}`

var multiChoiceOptionsTemplate = `
{{- range $ix, $option := .Options}}
  {{- if eq $ix $.MultiChoice.SelectedIndex}}{{color "cyan"}}❯{{color "reset"}}{{else}} {{end}}
  {{- if index $.MultiChoice.Checked $ix}}{{color "green"}} ◉ {{else}}{{color "default+hb"}} ◯ {{end}}
  {{- color "reset"}}
  {{- " "}}{{$option}}
{{end}}`

// OnChange is called on every keypress.
func (m *MultiChoice) OnChange(line []rune, pos int, key rune) (newLine []rune, newPos int, ok bool) {
	if key == core.KeyEnter {
		// just pass on the current value
		return line, 0, true
	} else if key == core.KeyArrowUp && m.SelectedIndex > 0 {
		// decrement the selected index
		m.SelectedIndex--
	} else if key == core.KeyArrowDown && m.SelectedIndex < len(m.Options)-1 {
		// if the user pressed down and there is room to move
		// increment the selected index
		m.SelectedIndex++
	} else if key == core.KeySpace {
		// otherwise just invert the current value
		m.Checked[m.SelectedIndex] = true
	}

	// print the template summarizing the current state of the selection

	// render the options
	// ansi.Print(m.SelectedIndex, m.Checked[m.SelectedIndex])
	m.render()

	// if we are not pressing ent
	return line, 0, true
}

func (m *MultiChoice) render() error {
	// clean up what we left behind last time
	for range m.Options {
		ansi.CursorPreviousLine(1)
		ansi.EraseInLine(1)
	}

	// render the template summarizing the current state
	out, err := core.RunTemplate(
		multiChoiceOptionsTemplate,
		multiChoiceTemplateData{MultiChoice: *m},
	)
	if err != nil {
		return err
	}

	// print the summary
	ansi.Println(strings.TrimRight(out, "\n"))

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
		Stdout:   &core.DevNull{},
	}
	rl.SetConfig(config)

	// compute the default state
	m.Checked = make(map[int]bool)
	// if there is a default
	if len(m.Defaults) > 0 {
		for _, dflt := range m.Defaults {
			for i, opt := range m.Options {
				// if the option correponds to the default
				if opt == dflt {
					// we found our initial value
					m.Checked[i] = true
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
		multiChoiceQuestionTemplate,
		multiChoiceTemplateData{MultiChoice: *m},
	)
	if err != nil {
		return "", err
	}
	// hide the cursor
	ansi.CursorHide()
	// ask the question
	ansi.Println(out)
	for range m.Options {
		ansi.Println()
	}

	// start waiting for input
	_, err = rl.Readline()
	// if something went wrong
	if err != nil {
		return "", err
	}
	// show the cursor when we're done
	ansi.CursorShow()

	answers := []string{}
	for ix, option := range m.Options {
		if val, ok := m.Checked[ix]; ok && val {
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
		if val, ok := m.Checked[ix]; ok && val {
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
	ansi.CursorPreviousLine(1)
	ansi.EraseInLine(1)
	for range m.Options {
		ansi.CursorPreviousLine(1)
		ansi.EraseInLine(1)
	}

	// parse the value into a list of strings
	var value []string
	json.Unmarshal([]byte(val), &value)

	// execute the output summary template with the answer
	output, err := core.RunTemplate(
		multiChoiceQuestionTemplate,
		multiChoiceTemplateData{MultiChoice: *m, Answer: value},
	)
	if err != nil {
		return err
	}
	// render the summary
	ansi.Println(output)

	// nothing went wrong
	return nil
}
