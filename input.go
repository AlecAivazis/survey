package survey

import (
	"fmt"

	"github.com/AlecAivazis/survey/core"
	"github.com/AlecAivazis/survey/terminal"
	"github.com/chzyer/readline"
)

// Input is a regular text input that prints each character the user types on the screen
// and accepts the input with the enter key.
type Input struct {
	Message string
	Default string
}

// data available to the templates when processing
type InputTemplateData struct {
	Input
	Answer string
}

// Templates with Color formatting. See Documentation: https://github.com/mgutz/ansi#style-format
var InputQuestionTemplate = `
{{- color "green+hb"}}? {{color "reset"}}
{{- color "default+hb"}}{{ .Message }} {{color "reset"}}
{{- if .Answer}}
  {{- color "cyan"}}{{.Answer}}{{color "reset"}}
{{- else }}
  {{- if .Default}}{{color "white"}}({{.Default}}) {{color "reset"}}{{end}}
{{- end}}`

func (i *Input) Prompt(rl *readline.Instance) (line interface{}, err error) {
	// render the template
	out, err := core.RunTemplate(
		InputQuestionTemplate,
		InputTemplateData{Input: *i},
	)
	if err != nil {
		return "", err
	}
	// make sure the prompt matches the expectation
	rl.SetPrompt(fmt.Sprintf(out))
	// get the next line
	line, err = rl.Readline()

	// if the line is empty
	if line == "" {
		// use the default value
		line = i.Default
	}

	// we're done
	return line, err
}

func (i *Input) Cleanup(rl *readline.Instance, val interface{}) error {
	// go up one line
	terminal.CursorPreviousLine(1)
	// clear the line
	terminal.EraseLine(terminal.ERASE_LINE_ALL)

	// render the template
	out, err := core.RunTemplate(
		InputQuestionTemplate,
		InputTemplateData{Input: *i, Answer: val.(string)},
	)
	if err != nil {
		return err
	}

	// print the summary
	terminal.Println(out)

	// we're done
	return err
}
