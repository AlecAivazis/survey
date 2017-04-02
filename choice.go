package survey

import (
	"fmt"
	"strings"

	tm "github.com/buger/goterm"
)

// Choice is a prompt that presents a list of various options to the user
// for them to select using the arrow keys and enter.
type Choice struct {
	Message string
	Choices []string
	Default string
}

// data available to the templates when processing
type ChoiceTemplateData struct {
	Choice
	Answer   string
	Selected int
}

var ChoiceQuestionTemplate = `
{{- color "green+hb"}}? {{color "reset"}}
{{- color "default+hb"}}{{ .Message }} {{color "reset"}}
{{- if .Answer}}{{color "cyan"}}{{.Answer}}{{color "reset"}}{{end}}`

var ChoiceChoicesTemplate = `
{{- range $ix, $choice := .Choices}}
  {{- if eq $ix $.Selected}}{{color "cyan"}}> {{else}}{{color "default+hb"}}  {{end}}
  {{- $choice}}
  {{- color "reset"}}
{{end}}`

// Prompt shows the list, and listens for input from the user using /dev/tty.
func (prompt *Choice) Prompt() (string, error) {
	out, err := runTemplate(
		ChoiceQuestionTemplate,
		ChoiceTemplateData{Choice: *prompt},
	)
	if err != nil {
		return "", err
	}
	// ask the question
	fmt.Println(out)

	initialLocation, err := InitialLocation(len(prompt.Choices))
	if err != nil {
		return "", err
	}

	// start off with the first option selected
	sel := 0
	// if there is a default
	if prompt.Default != "" {
		// find the choice
		for i, opt := range prompt.Choices {
			// if the option correponds to the default
			if opt == prompt.Default {
				// we found our initial value
				sel = i
				// stop looking
				break
			}
		}
	}

	// print the options to start
	err = prompt.refreshOptions(sel, initialLocation)
	if err != nil {
		return "", err
	}

	for {
		// wait for an input from the user
		_, keycode, err := GetInput(nil)
		// if there is an error
		if err != nil {
			// bubble up
			return "", err
		}

		// if the user pressed the up arrow and we can decrement sel
		if keycode == KeyArrowUp && sel > 0 {
			// decrement the selected index
			sel--
		}
		// if the user pressed the down arrow and we can decrement sel
		if keycode == KeyArrowDown && sel < len(prompt.Choices)-1 {
			// decrement the selected index
			sel++
		}

		// // if the user presses enter
		if keycode == KeyEnter {
			// we're done with the rendering loop (the current value is good)
			break
		}

		err = prompt.refreshOptions(sel, initialLocation)
		if err != nil {
			return "", err
		}
	}

	// return the selected choice
	return prompt.Choices[sel], nil
}

// Cleanup removes the choices section, and renders the ask like a normal question.
func (prompt *Choice) Cleanup(val string) error {
	output, err := runTemplate(
		ChoiceQuestionTemplate,
		ChoiceTemplateData{Choice: *prompt, Answer: val},
	)
	if err != nil {
		return err
	}
	return cleanupMultiOptions(len(prompt.Choices), output)
}

func (prompt *Choice) refreshOptions(sel int, initLoc int) error {
	out, err := runTemplate(
		ChoiceChoicesTemplate,
		ChoiceTemplateData{Choice: *prompt, Selected: sel},
	)
	if err != nil {
		return err
	}
	// ask the question
	tm.Print(strings.TrimRight(out, "\n"))
	tm.Flush()
	// make sure we overwrite the first line next time we print
	tm.MoveCursor(initLoc, 1)
	return nil
}
