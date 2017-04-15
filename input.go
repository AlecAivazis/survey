package survey

import (
	"fmt"

	"github.com/chzyer/readline"

	"github.com/alecaivazis/survey/core"
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

func (i *Input) Prompt(rl *readline.Instance) (line string, err error) {
	// render the template
	out, err := RunTemplate(
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
	// we're done
	return line, err
}

func (i *Input) Cleanup(rl *readline.Instance, val string) error {
	// go up one line
	core.CursorPreviousLine(1)
	// clear the line
	core.EraseInLine(1)

	// render the template
	out, err := RunTemplate(
		InputQuestionTemplate,
		InputTemplateData{Input: *i, Answer: val},
	)
	if err != nil {
		return err
	}

	// print the summary
	core.Println(out)

	// we're done
	return err
}
