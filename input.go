package survey

import "github.com/AlecAivazis/survey/v2/core"

/*
Input is a regular text input that prints each character the user types on the screen
and accepts the input with the enter key. Response type is a string.

	name := ""
	prompt := &survey.Input{ Message: "What is your name?" }
	survey.AskOne(prompt, &name)
*/
type Input struct {
	Renderer
	Message string
	Default string
	Help    string
	Suggest func(toComplete string) []string
}

// data available to the templates when processing
type InputTemplateData struct {
	Input
	Answer        string
	ShowAnswer    bool
	ShowHelp      bool
	PageEntries   []core.OptionAnswer
	SelectedIndex int
	Config        *PromptConfig
}

// Templates with Color formatting. See Documentation: https://github.com/mgutz/ansi#style-format
var InputQuestionTemplate = `
{{- if .ShowHelp }}{{- color .Config.Icons.Help.Format }}{{ .Config.Icons.Help.Text }} {{ .Help }}{{color "reset"}}{{"\n"}}{{end}}
{{- color .Config.Icons.Question.Format }}{{ .Config.Icons.Question.Text }} {{color "reset"}}
{{- color "default+hb"}}{{ .Message }} {{color "reset"}}
{{- if .ShowAnswer}}
  {{- color "cyan"}}{{.Answer}}{{color "reset"}}{{"\n"}}
{{- else if .PageEntries -}}
  {{- .Answer}} [Use arrows to navegate, enter to select, type to complement answer]
  {{- "\n"}}
  {{- range $ix, $choice := .PageEntries}}
    {{- if eq $ix $.SelectedIndex }}{{color $.Config.Icons.SelectFocus.Format }}{{ $.Config.Icons.SelectFocus.Text }} {{else}}{{color "default"}}  {{end}}
    {{- $choice.Value}}
    {{- color "reset"}}{{"\n"}}
  {{- end}}
{{- else }}
  {{- if and .Help (not .ShowHelp)}}{{color "cyan"}}[{{ print .Config.HelpInput }} for help]{{color "reset"}} {{end}}
  {{- if and .Suggest }}{{color "cyan"}}[{{ print .Config.SuggestInput }} for suggestions]{{color "reset"}} {{end}}
  {{- if .Default}}{{color "white"}}({{.Default}}) {{color "reset"}}{{end}}
  {{- .Answer -}}
{{- end}}`

func (i *Input) Prompt(config *PromptConfig) (interface{}, error) {
	// render the template
	err := i.Render(
		InputQuestionTemplate,
		InputTemplateData{
			Input:  *i,
			Config: config,
		},
	)
	if err != nil {
		return "", err
	}

	// start reading runes from the standard in
	rr := i.NewRuneReader()
	rr.SetTermMode()
	defer rr.RestoreTermMode()

	cursor := i.NewCursor()

	line := []rune{}
	// get the next line
	for {
		line, err = rr.ReadLine(0)
		if err != nil {
			return string(line), err
		}
		// terminal will echo the \n so we need to jump back up one row
		cursor.Up(1)

		if string(line) == config.HelpInput && i.Help != "" {
			err = i.Render(
				InputQuestionTemplate,
				InputTemplateData{
					Input:    *i,
					ShowHelp: true,
					Config:   config,
				},
			)
			if err != nil {
				return "", err
			}
			continue
		}
		break
	}

	// if the line is empty
	if line == nil || len(line) == 0 {
		// use the default value
		return i.Default, err
	}

	lineStr := string(line)

	i.AppendRenderedText(lineStr)

	// we're done
	return lineStr, err
}

func (i *Input) Cleanup(config *PromptConfig, val interface{}) error {
	return i.Render(
		InputQuestionTemplate,
		InputTemplateData{
			Input:      *i,
			Answer:     val.(string),
			ShowAnswer: true,
			Config:     config,
		},
	)
}
