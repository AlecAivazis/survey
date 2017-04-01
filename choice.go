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
  {{- if eq $ix $.Selected}}{{color "cyan"}} ➤ {{else}}{{color "default+hb"}}   {{end}}
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

	// get the current location of the cursor
	loc, err := CursorLocation()
	// if something went wrong
	if err != nil {
		// bubble up
		return "", err
	}

	// the height of the prompt's output
	height := len(prompt.Choices)

	// the starting point of the list depends on wether or not we
	// are at the bottom of the current terminal session
	var initialLocation int
	// if the options would fit cleanly
	if loc.col+height < tm.Height() {
		// start at the current location
		initialLocation = loc.col
		// otherwise we will be placed at the bottom of the terminal after this print
	} else {
		// the we have to start printing so that we just fit
		initialLocation = tm.Height() - height
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
	prompt.refreshOptions(sel, initialLocation)

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

		prompt.refreshOptions(sel, initialLocation)
	}

	// return the selected choice
	return prompt.Choices[sel], nil
}

// Cleanup removes the choices section, and renders the ask like a normal question.
func (prompt *Choice) Cleanup(val string) error {

	// the height of the prompt's output
	height := len(prompt.Choices)

	// get the current location of the cursor
	loc, err := CursorLocation()
	// if something went wrong
	if err != nil {
		// yell loudly
		return err
	}

	var initLoc int
	// if the options would fit cleanly
	if loc.col+height <= tm.Height() {
		// start at the current location
		initLoc = loc.col - height - 1
		// otherwise we will be placed at the bottom of the terminal after this print
	} else {

		// the we have to start printing so that we just fit
		initLoc = loc.col - height - 2
	}

	// find the index of the selected choice

	// start where we were told
	tm.MoveCursor(initLoc, 1)
	out, err := runTemplate(
		ChoiceQuestionTemplate,
		ChoiceTemplateData{Choice: *prompt, Answer: val},
	)
	if err != nil {
		return err
	}
	// ask the question
	tm.Print(out, AnsiClearLine)
	// for each choice
	for range prompt.Choices {
		// add an empty line
		tm.Print(AnsiClearLine)
		// print the output
		tm.Flush()
	}
	// add an empty line
	tm.Print(AnsiClearLine)
	// print the output
	tm.Flush()
	tm.MoveCursor(initLoc, 1)
	tm.Flush()

	// nothing went wrong
	return nil
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
