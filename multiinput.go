package survey

import (
	"fmt"
	"log"

	"github.com/AlecAivazis/survey/v2/core"
)

/*
MultiInput utilises the Input type to provide the user with a means of entering
multiple values for a field. Response type is a []string.

	names := make([]string, 0)
	prompt := &survey.MultiInput{ Message: "What are your friend's names?" }
	survey.AskOne(prompt, &names, nil)
*/
type MultiInput struct {
	core.Renderer
	Message string
	Default []string
	Help    string
}

// data available to the templates when processing
type MultiInputTemplateData struct {
	MultiInput
	Answer     []string
	ShowAnswer bool
	ShowHelp   bool
}

var MultiInputQuestionTemplate = `
{{- color "green+hb"}}{{ QuestionIcon }} {{color "reset"}}
{{- color "default+hb"}}{{ .Message }} {{color "reset"}}
{{- if .ShowAnswer}}
  {{- color "cyan"}}{{.Answer}}{{color "reset"}}{{"\n"}}
{{- else }}
  {{- "  "}}{{- color "cyan"}}[Enter empty line to finish] {{color "reset"}}
  {{- if .Default}}{{color "white"}}({{.Default}}) {{color "reset"}}{{end}}
{{- end}}`

func (i *MultiInput) Prompt() (interface{}, error) {
	// ctore answers
	answers := make([]string, 0)
	// current loop iteration
	counter := 0

	// render the template
	err := i.Render(
		MultiInputQuestionTemplate,
		MultiInputTemplateData{MultiInput: *i},
	)
	if err != nil {
		return "", err
	}

	// required to prevent the question from being removed
	fmt.Println("")

	// run a loop to accept n answers
	for {
		counter++
		answer := ""
		message := fmt.Sprintf("#%d:", counter)

		prompt := &Input{
			Message: message,
			Help:    i.Help,
		}

		err := AskOne(prompt, &answer, nil)
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}

		// check if answer is the escape sequence
		if answer == "" {
			break
		}

		// return answers
		answers = append(answers, answer)
	}

	// reset cursor to remove entries
	cursor := i.NewCursor()
	cursor.PreviousLine(counter + 1)

	// if the answers are empty
	if len(answers) == 0 {
		// use the default value
		return i.Default, err
	}

	return answers, nil
}

func (i *MultiInput) Cleanup(val interface{}) error {
	return i.Render(
		MultiInputQuestionTemplate,
		MultiInputTemplateData{MultiInput: *i, Answer: val.([]string), ShowAnswer: true, ShowHelp: false},
	)
}
