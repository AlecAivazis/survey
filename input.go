package survey

import (
	"bufio"
	"os"

	"github.com/AlecAivazis/survey/core"
	"github.com/AlecAivazis/survey/terminal"
)

// Input is a regular text input that prints each character the user types on the screen
// and accepts the input with the enter key.
type Input struct {
	core.Renderer
	Message string
	Default string
	Help    string
}

// data available to the templates when processing
type InputTemplateData struct {
	Input
	Answer     string
	ShowAnswer bool
	ShowHelp   bool
}

// Templates with Color formatting. See Documentation: https://github.com/mgutz/ansi#style-format
var InputQuestionTemplate = `
{{- if .ShowHelp }}{{- color "cyan"}}{{ HelpIcon }} {{ .Help }}{{color "reset"}}{{"\n"}}{{end}}
{{- color "green+hb"}}{{ QuestionIcon }} {{color "reset"}}
{{- color "default+hb"}}{{ .Message }} {{color "reset"}}
{{- if .ShowAnswer}}
  {{- color "cyan"}}{{.Answer}}{{color "reset"}}{{"\n"}}
{{- else }}
  {{- if and .Help (not .ShowHelp)}}{{color "cyan"}}[{{ HelpInputRune }} for help]{{color "reset"}} {{end}}
  {{- if .Default}}{{color "white"}}({{.Default}}) {{color "reset"}}{{end}}
{{- end}}`

func (i *Input) Prompt() (line interface{}, err error) {
	// render the template
	err = i.Render(
		InputQuestionTemplate,
		InputTemplateData{Input: *i},
	)
	if err != nil {
		return "", err
	}

	scanner := bufio.NewScanner(os.Stdin)
	// get the next line
	for scanner.Scan() {
		// terminal will echo the \n so we need to jump back up one row
		terminal.CursorPreviousLine(1)
		line = scanner.Text()

		if line == string(core.HelpInputRune) {
			err = i.Render(
				InputQuestionTemplate,
				InputTemplateData{Input: *i, ShowHelp: true},
			)
			if err != nil {
				return "", err
			}
			continue
		}
		break
	}

	if err := scanner.Err(); err != nil {
		return i.Default, err
	}
	// if the line is empty
	if line == "" || line == nil {
		// use the default value
		line = i.Default
	}

	// we're done
	return line, err
}

func (i *Input) Cleanup(val interface{}) error {
	return i.Render(
		InputQuestionTemplate,
		InputTemplateData{Input: *i, Answer: val.(string), ShowAnswer: true},
	)
}
