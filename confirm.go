package survey

import (
	"fmt"
	"regexp"

	"github.com/alecaivazis/survey/core"
	"github.com/alecaivazis/survey/terminal"
	"github.com/chzyer/readline"
)

// Confirm is a regular text input that accept yes/no answers.
type Confirm struct {
	Message string
	Default bool
}

// data available to the templates when processing
type ConfirmTemplateData struct {
	Confirm
	Answer string
}

// Templates with Color formatting. See Documentation: https://github.com/mgutz/ansi#style-format
var ConfirmQuestionTemplate = `
{{- color "green+hb"}}? {{color "reset"}}
{{- color "default+hb"}}{{ .Message }} {{color "reset"}}
{{- if .Answer}}
  {{- color "cyan"}}{{.Answer}}{{color "reset"}}
{{- else }}
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

func (c *Confirm) getBool(rl *readline.Instance) (bool, error) {
	// start waiting for input
	val, err := rl.Readline()
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
	default:
		// we didnt get a valid answer, so print error and prompt again
		out, err := core.RunTemplate(
			ErrorTemplate, fmt.Errorf("%q is not a valid answer, please try again.", val),
		)
		// if something went wrong
		if err != nil {
			// use the default value and bubble up
			return c.Default, err
		}
		// send the message to the user
		terminal.Print(out)

		answer, err = c.getBool(rl)
		// if something went wrong
		if err != nil {
			// use the default value
			return c.Default, err
		}
	}

	return answer, nil
}

// Prompt prompts the user with a simple text field and expects a reply followed
// by a carriage return.
func (c *Confirm) Prompt(rl *readline.Instance) (interface{}, error) {
	// render the question template
	out, err := core.RunTemplate(
		ConfirmQuestionTemplate,
		ConfirmTemplateData{Confirm: *c},
	)
	if err != nil {
		return "", err
	}

	// use the result of the template as the prompt for the readline instance
	rl.SetPrompt(fmt.Sprintf(out))

	// start waiting for input
	answer, err := c.getBool(rl)
	// if something went wrong
	if err != nil {
		// bubble up
		return "", err
	}

	// convert the boolean into the appropriate string
	return answer, nil
}

// Cleanup overwrite the line with the finalized formatted version
func (c *Confirm) Cleanup(rl *readline.Instance, val interface{}) error {
	// go up one line
	terminal.CursorPreviousLine(1)
	// clear the line
	terminal.EraseInLine(0)

	// the string version of the answer
	ans := ""
	// if the value was previously true
	if val.(bool) {
		ans = "true"
	} else {
		ans = "false"
	}

	// render the template
	out, err := core.RunTemplate(
		ConfirmQuestionTemplate,
		ConfirmTemplateData{Confirm: *c, Answer: ans},
	)
	if err != nil {
		return err
	}

	// print the summary
	terminal.Println(out)

	// we're done
	return nil
}
