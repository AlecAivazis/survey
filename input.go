package survey

import (
	"github.com/AlecAivazis/survey/terminal"
	"github.com/chzyer/readline"
)

// Input is a regular text input that prints each character the user types on the screen
// and accepts the input with the enter key.
type Input struct {
	renderer
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
  {{- color "cyan"}}{{.Answer}}{{color "reset"}}{{"\n"}}
{{- else }}
  {{- if .Default}}{{color "white"}}({{.Default}}) {{color "reset"}}{{end}}
{{- end}}`

func (i *Input) Prompt(rl *readline.Instance) (line interface{}, err error) {
	// render the template
	err = i.render(
		InputQuestionTemplate,
		InputTemplateData{Input: *i},
	)
	if err != nil {
		return "", err
	}
	rl.SetConfig(simpleReadlineConfig)

	// get the next line
	ans, err := rl.Readline()
	// readline will echo the \n so we need to jump back up one row
	terminal.CursorUp(1)
	return ans, err
}

func (i *Input) Cleanup(rl *readline.Instance, val interface{}) error {
	return i.render(
		InputQuestionTemplate,
		InputTemplateData{Input: *i, Answer: val.(string)},
	)
}
