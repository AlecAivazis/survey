package survey

import (
	"github.com/AlecAivazis/survey/core"
	"github.com/AlecAivazis/survey/terminal"
	"github.com/chzyer/readline"
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
	Answer   string
	ShowHelp bool
}

// Templates with Color formatting. See Documentation: https://github.com/mgutz/ansi#style-format
var InputQuestionTemplate = `
{{- if .ShowHelp }}{{- color "cyan"}}{{ .Help }}{{color "reset"}}{{"\n"}}{{end}}
{{- color "green+hb"}}? {{color "reset"}}
{{- color "default+hb"}}{{ .Message }} {{color "reset"}}
{{- if .Answer}}
  {{- color "cyan"}}{{.Answer}}{{color "reset"}}{{"\n"}}
{{- else }}
  {{- if and .Help (not .ShowHelp)}}{{color "cyan"}}[? for help] {{color "reset"}}{{end}}
  {{- if .Default}}{{color "white"}}({{.Default}}) {{color "reset"}}{{end}}
{{- end}}`

func (i *Input) Prompt(rl *readline.Instance) (line interface{}, err error) {
	// render the template
	err = i.Render(
		InputQuestionTemplate,
		InputTemplateData{Input: *i},
	)
	if err != nil {
		return "", err
	}
	rl.SetConfig(core.SimpleReadlineConfig)

	// get the next line
	line, err = rl.Readline()
	// readline will echo the \n so we need to jump back up one row
	terminal.CursorUp(1)

	if err == nil && line == "?" {
		err = i.Render(
			InputQuestionTemplate,
			InputTemplateData{Input: *i, ShowHelp: true},
		)
		if err != nil {
			return "", err
		}
		// get the next line
		line, err = rl.Readline()
		// readline will echo the \n so we need to jump back up one row
		terminal.CursorUp(1)
	}

	// if the line is empty
	if line == "" {
		// use the default value
		line = i.Default
	}

	// we're done
	return line, err
}

func (i *Input) Cleanup(rl *readline.Instance, val interface{}) error {
	return i.Render(
		InputQuestionTemplate,
		InputTemplateData{Input: *i, Answer: val.(string)},
	)
}
