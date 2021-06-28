package survey

import (
	"github.com/AlecAivazis/survey/v2/core"
	"github.com/AlecAivazis/survey/v2/terminal"
)

/*
Input is a regular text input that prints each character the user types on the screen
and accepts the input with the enter key. Response type is a string.

	name := ""
	prompt := &survey.Input{ Message: "What is your name?" }
	survey.AskOne(prompt, &name)
*/
type Input struct {
	Renderer
	Message       string
	Default       string
	Help          string
	Suggest       func(toComplete string) []string
	typedAnswer   string
	answer        string
	options       []core.OptionAnswer
	selectedIndex int
	showingHelp   bool
}

// data available to the templates when processing
type InputTemplateData struct {
	Input
	ShowAnswer    bool
	ShowHelp      bool
	Answer        string
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
  {{- .Answer}} [Use arrows to move, enter to select, type to continue]
  {{- "\n"}}
  {{- range $ix, $choice := .PageEntries}}
    {{- if eq $ix $.SelectedIndex }}{{color $.Config.Icons.SelectFocus.Format }}{{ $.Config.Icons.SelectFocus.Text }} {{else}}{{color "default"}}  {{end}}
    {{- $choice.Value}}
    {{- color "reset"}}{{"\n"}}
  {{- end}}
{{- else }}
  {{- if or (and .Help (not .ShowHelp)) .Suggest }}{{color "cyan"}}[
    {{- if and .Help (not .ShowHelp)}}{{ print .Config.HelpInput }} for help {{- if and .Suggest}}, {{end}}{{end -}}
    {{- if and .Suggest }}{{color "cyan"}}{{ print .Config.SuggestInput }} for suggestions{{end -}}
  ]{{color "reset"}} {{end}}
  {{- if .Default}}{{color "white"}}({{.Default}}) {{color "reset"}}{{end}}
  {{- .Answer -}}
{{- end}}`

func (i *Input) OnSpecialKey(config *PromptConfig) func(rune) bool {
	return func(key rune) bool {
		if key == terminal.KeyArrowUp && len(i.options) > 0 {
			if i.selectedIndex == 0 {
				i.selectedIndex = len(i.options) - 1
			} else {
				i.selectedIndex--
			}
			i.answer = i.options[i.selectedIndex].Value
			return true
		} else if (key == terminal.KeyArrowDown || key == terminal.KeyTab) && len(i.options) > 0 {
			if i.selectedIndex == len(i.options)-1 {
				i.selectedIndex = 0
			} else {
				i.selectedIndex++
			}
			i.answer = i.options[i.selectedIndex].Value
			return true
		} else if key == terminal.KeyTab && i.Suggest != nil {
			options := i.Suggest(i.answer)
			i.selectedIndex = 0
			i.typedAnswer = i.answer
			if len(options) > 0 {
				i.answer = options[0]
				if len(options) == 1 {
					i.options = nil
				} else {
					i.options = core.OptionAnswerList(options)
				}
			}
			return true
		}

		return false
	}
}

func (i *Input) Prompt(config *PromptConfig) (interface{}, error) {
	// render the template
	err := i.Render(
		InputQuestionTemplate,
		InputTemplateData{
			Input:    *i,
			Config:   config,
			ShowHelp: i.showingHelp,
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
	if !config.ShowCursor {
		cursor.Hide()       // hide the cursor
		defer cursor.Show() // show the cursor when we're done
	}

	line, err := rr.ReadLine(0, i.OnSpecialKey(config))
	if err != nil {
		return "", err
	}
	i.answer = string(line)
	// readline print an empty line, go up before we render the follow up
	cursor.Up(1)

	// if we ran into the help string
	if i.answer == config.HelpInput && i.Help != "" {
		// show the help and prompt again
		i.showingHelp = true
		return i.Prompt(config)
	}

	pageSize := config.PageSize
	opts, idx := paginate(pageSize, i.options, i.selectedIndex)
	err = i.Render(
		InputQuestionTemplate,
		InputTemplateData{
			Input:         *i,
			Answer:        i.answer,
			ShowHelp:      i.showingHelp,
			SelectedIndex: idx,
			PageEntries:   opts,
			Config:        config,
		},
	)
	if err != nil {
		return "", err
	}

	// if the line is empty
	if len(i.answer) == 0 {
		// use the default value
		return i.Default, err
	}

	lineStr := i.answer

	i.AppendRenderedText(lineStr)

	// we're done
	return lineStr, err
}

func (i *Input) Cleanup(config *PromptConfig, val interface{}) error {
	// use the default answer when cleaning up the prompt if necessary
	ans := i.answer
	if ans == "" && i.Default != "" {
		ans = i.Default
	}

	// render the cleanup
	return i.Render(
		InputQuestionTemplate,
		InputTemplateData{
			Input:      *i,
			ShowAnswer: true,
			Config:     config,
			Answer:     ans,
		},
	)
}
