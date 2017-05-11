package survey

import (
	"fmt"
	"os"
	"regexp"

	"github.com/AlecAivazis/survey/core"
	"github.com/AlecAivazis/survey/terminal"
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

func (c *Confirm) getBool(showHelp bool) (bool, error) {
	rr := terminal.NewRuneReader(os.Stdin)
	rr.SetTermMode()
	defer rr.RestoreTermMode()
	// start waiting for input
	for {
		line, err := rr.ReadLine(0)
		if err != nil {
			return false, err
		}
		// move back up a line to compensate for the \n echoed from terminal
		terminal.CursorPreviousLine(1)
		val := string(line)

		// get the answer that matches the
		var answer bool
		switch {
		case yesRx.Match([]byte(val)):
			answer = true
		case noRx.Match([]byte(val)):
			answer = false
		case val == "":
			answer = c.Default
		case val == string(core.HelpInputRune) && c.Help != "":
			err := c.Render(
				ConfirmQuestionTemplate,
				ConfirmTemplateData{Confirm: *c, ShowHelp: true},
			)
			if err != nil {
				// use the default value and bubble up
				return c.Default, err
			}
			showHelp = true
			continue
		default:
			// we didnt get a valid answer, so print error and prompt again
			e := fmt.Errorf("%q is not a valid answer, please try again.", val)
			err := c.Render(
				ConfirmQuestionTemplate,
				ConfirmTemplateData{Confirm: *c, Error: &e, ShowHelp: showHelp},
			)
			if err != nil {
				// use the default value and bubble up
				return c.Default, err
			}
			continue
		}
		return answer, nil
	}
	// should not get here
	return c.Default, nil
}

// Prompt prompts the user with a simple text field and expects a reply followed
// by a carriage return.
func (c *Confirm) Prompt() (interface{}, error) {
	// render the question template
	err := c.Render(
		ConfirmQuestionTemplate,
		ConfirmTemplateData{Confirm: *c},
	)
	if err != nil {
		return "", err
	}

	// get input and return
	return c.getBool(false)
}

// Cleanup overwrite the line with the finalized formatted version
func (c *Confirm) Cleanup(val interface{}) error {
	// if the value was previously true
	ans := yesNo(val.(bool))
	// render the template
	return c.Render(
		ConfirmQuestionTemplate,
		ConfirmTemplateData{Confirm: *c, Answer: ans},
	)
}
