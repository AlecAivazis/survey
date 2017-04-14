package survey

// import (
// 	"encoding/json"
// 	"fmt"
// 	"strings"

// 	"github.com/alecaivazis/survey/core"
// 	tm "github.com/buger/goterm"
// )

// // MultiChoice is a prompt that presents a list of various options to the user
// // for them to select using the arrow keys and enter.
// type MultiChoice struct {
// 	Message  string
// 	Options  []string
// 	Defaults []string
// 	Answer   *[]string
// }

// // data available to the templates when processing
// type MultiChoiceTemplateData struct {
// 	MultiChoice
// 	Answer   []string
// 	Checked  map[int]bool
// 	Selected int
// }

// var MultiChoiceQuestionTemplate = `
// {{- color "green+hb"}}? {{color "reset"}}
// {{- color "default+hb"}}{{ .Message }} {{color "reset"}}
// {{- if .Answer}}{{color "cyan"}}{{.Answer | printf "%q"}}{{color "reset"}}{{end}}`

// var MultiChoiceOptionsTemplate = `
// {{- range $ix, $option := .Options}}
//   {{- if eq $ix $.Selected}}{{color "cyan"}}❯{{color "reset"}}{{else}} {{end}}
//   {{- if index $.Checked $ix}}{{color "green"}} ◉ {{else}}{{color "default+hb"}} ◯ {{end}}
//   {{- color "reset"}}
//   {{- " "}}{{$option}}
// {{end}}`

// // Prompt shows the list, and listens for input from the user using /dev/tty.
// func (prompt *MultiChoice) Prompt() (string, error) {
// 	if prompt.Answer == nil {
// 		answer := []string{}
// 		prompt.Answer = &answer
// 	}
// 	out, err := core.RunTemplate(
// 		MultiChoiceQuestionTemplate,
// 		MultiChoiceTemplateData{MultiChoice: *prompt},
// 	)
// 	if err != nil {
// 		return "", err
// 	}
// 	// ask the question
// 	fmt.Println(out)

// 	initialLocation, err := InitialLocation(len(prompt.Options))
// 	if err != nil {
// 		return "", err
// 	}

// 	sel := 0
// 	checked := map[int]bool{}
// 	// if there is a default
// 	if len(prompt.Defaults) > 0 {
// 		for _, dflt := range prompt.Defaults {
// 			for i, opt := range prompt.Options {
// 				// if the option correponds to the default
// 				if opt == dflt {
// 					// we found our initial value
// 					checked[i] = true
// 					// stop looking
// 					break
// 				}
// 			}
// 		}
// 	}

// 	// print the options to start
// 	err = prompt.refreshOptions(checked, sel, initialLocation)
// 	if err != nil {
// 		return "", err
// 	}

// 	for {
// 		// wait for an input from the user
// 		_, keycode, err := GetInput(nil)
// 		// if there is an error
// 		if err != nil {
// 			// bubble up
// 			return "", err
// 		}

// 		// if the user pressed the up arrow and we can decrement sel
// 		if keycode == KeyArrowUp && sel > 0 {
// 			// decrement the selected index
// 			sel--
// 		}
// 		// if the user pressed the down arrow and we can decrement sel
// 		if keycode == KeyArrowDown && sel < len(prompt.Options)-1 {
// 			// decrement the selected index
// 			sel++
// 		}

// 		if keycode == KeySpace {
// 			if val, ok := checked[sel]; !ok {
// 				checked[sel] = true
// 			} else {
// 				checked[sel] = !val
// 			}
// 		}

// 		// // if the user presses enter
// 		if keycode == KeyEnter {
// 			// we're done with the rendering loop (the current value is good)
// 			break
// 		}

// 		err = prompt.refreshOptions(checked, sel, initialLocation)
// 		if err != nil {
// 			return "", err
// 		}
// 	}

// 	answers := []string{}
// 	for ix, option := range prompt.Options {
// 		if val, ok := checked[ix]; ok && val {
// 			answers = append(answers, option)
// 		}
// 	}
// 	*prompt.Answer = answers
// 	// return the selected option
// 	js, err := json.Marshal(answers)
// 	if err != nil {
// 		return "", err
// 	}
// 	return string(js), nil
// }

// // Cleanup removes the options section, and renders the ask like a normal question.
// func (prompt *MultiChoice) Cleanup(val string) error {
// 	output, err := core.RunTemplate(
// 		MultiChoiceQuestionTemplate,
// 		MultiChoiceTemplateData{MultiChoice: *prompt, Answer: *prompt.Answer},
// 	)
// 	if err != nil {
// 		return err
// 	}
// 	return cleanupMultiOptions(len(prompt.Options), output)
// }

// func (prompt *MultiChoice) refreshOptions(checked map[int]bool, sel int, initLoc int) error {
// 	out, err := core.RunTemplate(
// 		MultiChoiceOptionsTemplate,
// 		MultiChoiceTemplateData{MultiChoice: *prompt, Selected: sel, Checked: checked},
// 	)
// 	if err != nil {
// 		return err
// 	}
// 	// ask the question
// 	tm.Print(strings.TrimRight(out, "\n"))
// 	tm.Flush()
// 	// make sure we overwrite the first line next time we print
// 	tm.MoveCursor(initLoc, 1)
// 	return nil
// }
