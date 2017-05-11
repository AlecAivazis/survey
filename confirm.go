package survey

import (
	"fmt"
	"regexp"

	"github.com/AlecAivazis/survey/core"
	"github.com/AlecAivazis/survey/terminal"
	"github.com/chzyer/readline"
)

// Confirm is a regular text input that accept yes/no answers.
type Confirm struct {
	core.Renderer
	Message string
	Default bool
	Help    string
}

// data available to the templates when processing
type ConfirmTemplateData struct {
	Confirm
	Answer   string
	Error    *error
	ShowHelp bool
}

// Templates with Color formatting. See Documentation: https://github.com/mgutz/ansi#style-format
var ConfirmQuestionTemplate = `
{{- if .ShowHelp }}{{- color "cyan"}}{{ HelpIcon }} {{ .Help }}{{color "reset"}}{{"\n"}}{{end}}
{{- if .Error }}` + ErrorTemplate + `{{end}}
{{- color "green+hb"}}{{ QuestionIcon }} {{color "reset"}}
{{- color "default+hb"}}{{ .Message }} {{color "reset"}}
{{- if .Answer}}
  {{- color "cyan"}}{{.Answer}}{{color "reset"}}{{"\n"}}
{{- else }}
  {{- if and .Help (not .ShowHelp)}}{{color "cyan"}}[{{ HelpInputRune }} for help]{{color "reset"}} {{end}}
  {{- color "white"}}{{if .Default}}(Y/n) {{else}}(y/N) {{end}}{{color "reset"}}
{{- end}}`

// the regex for answers
var (
	yesRx = regexp.MustCompile("^(?i:y(?:es)?)$")
	noRx  = regexp.MustCompile("^(?i:n(?:o)?)$")
)

func yesNo(t bool) string {
	if t {
		return "Yes"
	}
	return "No"
}

func (c *Confirm) getBool(rl *readline.Instance, showHelp bool) (bool, error) {
	// start waiting for input
	val, err := rl.Readline()
	// move back up a line to compensate for the \n echoed from Readline
	terminal.CursorUp(1)
	// if something went wrong
	if err != nil {
		// use the default value and bubble up
		return c.Default, err
	}

	// get the answer that matches the
	var answer bool
	switch {
	case yesRx.Match([]byte(val)):
		answer = true
	case noRx.Match([]byte(val)):
		answer = false
	case val == "":
		answer = c.Default
	// only show the help message if we have one
	case val == string(core.HelpInputRune) && c.Help != "":
		err = c.Render(
			ConfirmQuestionTemplate,
			ConfirmTemplateData{Confirm: *c, ShowHelp: true},
		)
		if err != nil {
			// use the default value and bubble up
			return c.Default, err
		}
		return c.getBool(rl, true)
	default:
		// we didnt get a valid answer, so print error and prompt again
		e := fmt.Errorf("%q is not a valid answer, please try again.", val)
		err = c.Render(
			ConfirmQuestionTemplate,
			ConfirmTemplateData{Confirm: *c, Error: &e, ShowHelp: showHelp},
		)
		if err != nil {
			// use the default value and bubble up
			return c.Default, err
		}
		return c.getBool(rl, showHelp)
	}

	return answer, nil
}

// Prompt prompts the user with a simple text field and expects a reply followed
// by a carriage return.
func (c *Confirm) Prompt(rl *readline.Instance) (interface{}, error) {
	// render the question template
	err := c.Render(
		ConfirmQuestionTemplate,
		ConfirmTemplateData{Confirm: *c},
	)
	if err != nil {
		return "", err
	}

	rl.SetConfig(core.SimpleReadlineConfig)

	// get input and return
	return c.getBool(rl, false)
}

// Cleanup overwrite the line with the finalized formatted version
func (c *Confirm) Cleanup(rl *readline.Instance, val interface{}) error {
	// if the value was previously true
	ans := yesNo(val.(bool))
	// render the template
	return c.Render(
		ConfirmQuestionTemplate,
		ConfirmTemplateData{Confirm: *c, Answer: ans},
	)
}
