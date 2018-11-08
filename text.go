package survey

import "gopkg.in/AlecAivazis/survey.v1/core"

type level int

const (
	Info level = iota
	Warning
	Danger
)

/*
Text is a prompt that prints message without answer.

	prompt := &survey.Text{ Message: "This is a important message.",Level:survey.Danger }
	survey.AskOne(prompt, nil, nil)
*/
type Text struct {
	core.Renderer
	Message string
	Help    string
	Level   level
}

func (*Text) NeedAnswer() bool {
	return false
}

// data available to the templates when processing
type TextTemplateData struct {
	Text
	ShowHelp bool
}

// Templates with Color formatting. See Documentation: https://github.com/mgutz/ansi#style-format
var InfoTemplate = `
{{- if eq .Level 1 }}
	{{- color "yellow+hb"}}{{ TextIcon }} {{ .Message }}{{color "reset"}}{{"\n"}}
{{- else if eq .Level 2 }}
	{{- color "red+hb"}}{{ TextIcon }} {{ .Message }} {{color "reset"}}{{"\n"}}
{{- else }}
	{{- color "cyan+hb"}}{{ TextIcon }} {{ .Message }}{{color "reset"}}{{"\n"}}
{{- end }}`

func (info *Text) Prompt() (interface{}, error) {
	err := info.Render(
		InfoTemplate,
		TextTemplateData{Text: *info},
	)
	if err != nil {
		return "", err
	}
	return "", nil
}

func (info *Text) Cleanup(interface{}) error {
	return nil
}
