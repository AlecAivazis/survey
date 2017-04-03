package survey

import (
	"bufio"
	"fmt"
	"os"
	"regexp"

	tm "github.com/buger/goterm"
)

// Confirm is a regular text input that accept yes/no answers.
type Confirm struct {
	Message string
	Default bool
	Answer  *bool
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

func yesNo(t bool) string {
	if t {
		return "Yes"
	}
	return "No"
}

// Prompt prompts the user with a simple text field and expects a reply followed
// by a carriage return.
func (confirm *Confirm) Prompt() (string, error) {
	out, err := runTemplate(
		ConfirmQuestionTemplate,
		ConfirmTemplateData{Confirm: *confirm},
	)
	if err != nil {
		return "", err
	}

	// print the question we were given to kick off the prompt
	fmt.Print(out)

	// a scanner to look at the input from stdin
	scanner := bufio.NewScanner(os.Stdin)
	// wait for a response
	yesRx := regexp.MustCompile("^(?i:y(?:es)?)$")
	noRx := regexp.MustCompile("^(?i:n(?:o)?)$")
	answer := confirm.Default
	for scanner.Scan() {
		// get the availible text in the scanner
		res := scanner.Text()
		// if there is no answer
		if res == "" {
			// use the default
			break
		}
		// is answer yes?
		if yesRx.Match([]byte(res)) {
			answer = true
			break
		}

		// is answer "no"
		if noRx.Match([]byte(res)) {
			answer = false
			break
		}

		// we didnt get a valid answer, so print error and prompt again
		out, err := runTemplate(ErrorTemplate, fmt.Errorf("%q is not a valid answer, try again", res))
		if err != nil {
			return "", err
		}
		// send the message to the user
		fmt.Print(out)
		return confirm.Prompt()
	}

	// return the value
	*confirm.Answer = answer
	return yesNo(answer), nil
}

// Cleanup overwrite the line with the finalized formatted version
func (confirm *Confirm) Cleanup(val string) error {
	// get the current cursor location
	loc, err := CursorLocation()
	// if something went wrong
	if err != nil {
		// bubble
		return err
	}

	var initLoc int
	// if we are printing at the end of the console
	if loc.col == tm.Height() {
		initLoc = loc.col - 2
	} else {
		initLoc = loc.col - 1
	}

	// move to the beginning of the current line
	tm.MoveCursor(initLoc, 1)

	out, err := runTemplate(
		ConfirmQuestionTemplate,
		ConfirmTemplateData{Confirm: *confirm, Answer: val},
	)
	if err != nil {
		return err
	}

	tm.Print(out, AnsiClearLine)
	tm.Flush()

	// nothing went wrong
	return nil
}
