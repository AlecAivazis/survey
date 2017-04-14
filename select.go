package survey

import (
	"strings"

	tm "github.com/buger/goterm"
	"github.com/chzyer/readline"

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
	selectQuestionTemplate = `
{{- color "green+hb"}}? {{color "reset"}}
{{- color "default+hb"}}{{ .Message }} {{color "reset"}}
{{- if .Answer}}{{color "cyan"}}{{.Answer}}{{color "reset"}}{{end}}`
	// the template used to show the list of Selects
	selectChoicesTemplate = `
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

func (s *Select) render() {
	// the formatted response
	out, err := core.RunTemplate(
		selectChoicesTemplate,
		SelectTemplateData{Select: *s},
	)
	if err != nil {
		panic(err)
	}

	// ask the question
	tm.Print(strings.TrimRight(out, "\n"))
	tm.Flush()

	// move up one line for each option
	for range s.Options {
		tm.Print(core.AnsiCursorUp)
	}
}

func (s *Select) Prompt(rl *readline.Instance) (string, error) {
	config := &readline.Config{
		Listener: s,
		Stdout:   &core.DevNull{},
	}
	rl.SetConfig(config)

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
	out, err := core.RunTemplate(
		selectQuestionTemplate,
		SelectTemplateData{Select: *s},
	)
	if err != nil {
		return "", err
	}

	// ask the question
	tm.Println(out)

	// start waiting for input
	return rl.Readline()
}

func (s *Select) Cleanup(rl *readline.Instance, val string) error {
	// remove the original message we left behind
	tm.Print(core.AnsiCursorUp)
	tm.Print(core.AnsiClearLine)

	// execute the output summary template with the answer
	output, err := core.RunTemplate(
		selectQuestionTemplate,
		SelectTemplateData{Select: *s, Answer: val},
	)
	if err != nil {
		return err
	}
	// render the summary
	tm.Print(output)
	tm.Flush()

	// nothing went wrong
	return nil
}
