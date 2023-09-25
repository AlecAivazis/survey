package survey

import (
	"strings"
)

type Multiline struct {
	Renderer
	Message string
	Default string
	Help    string
}

// data available to the templates when processing
type MultilineTemplateData struct {
	Multiline
	Answer     string
	ShowAnswer bool
	ShowHelp   bool
	Config     *PromptConfig
}

// Templates with Color formatting. See Documentation: https://github.com/mgutz/ansi#style-format
var MultilineQuestionTemplate = `
{{- if .ShowHelp }}{{- color .Config.Icons.Help.Format }}{{ .Config.Icons.Help.Text }} {{ .Help }}{{color "reset"}}{{"\n"}}{{end}}
{{- color .Config.Icons.Question.Format }}{{ .Config.Icons.Question.Text }} {{color "reset"}}
{{- color "default+hb"}}{{ .Message }} {{color "reset"}}
{{- if .ShowAnswer}}
  {{- if .Answer}}{{"\n"}}{{color "cyan"}}{{.Answer}}{{color "reset"}}{{end}}
{{- else }}
  {{- if .Default}}{{color "white"}}({{.Default}}) {{color "reset"}}{{end}}
  {{- color "cyan"}}[Enter 2 empty lines to finish]{{color "reset"}}
{{- end}}
`

func (i *Multiline) Prompt(config *PromptConfig) (interface{}, error) {
	// render the template
	err := i.Render(
		MultilineQuestionTemplate,
		MultilineTemplateData{
			Multiline: *i,
			Config:    config,
		},
	)
	if err != nil {
		return "", err
	}

	// start reading runes from the standard in
	rr := i.NewRuneReader()
	_ = rr.SetTermMode()
	defer func() {
		_ = rr.RestoreTermMode()
	}()

	cursor := i.NewCursor()

	multiline := make([]string, 0)

	emptyOnce := false
	// get the next line
	for {
		var line []rune
		line, err = rr.ReadLine(0)
		if err != nil {
			return string(line), err
		}

		if string(line) == "" {
			if emptyOnce {
				break
			}
			emptyOnce = true
		} else {
			emptyOnce = false
		}
		multiline = append(multiline, string(line))
	}

	// adjust for terminating newlines
	cursor.PreviousLine(2)

	// render the displayed value or use the default
	val := strings.Join(multiline, "\n")
	if len(val) == 0 {
		return i.Default, err
	}
	i.AppendRenderedText(val)

	// remove the extra newline from the answer
	ans := strings.TrimSuffix(val, "\n")
	return ans, err
}

func (i *Multiline) Cleanup(config *PromptConfig, val interface{}) error {
	return i.Render(
		MultilineQuestionTemplate,
		MultilineTemplateData{
			Multiline:  *i,
			Answer:     val.(string),
			ShowAnswer: true,
			Config:     config,
		},
	)
}
