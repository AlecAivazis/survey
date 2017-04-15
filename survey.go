package survey

import (
	"fmt"
	"os"

	"github.com/alecaivazis/survey/terminal"
	"github.com/chzyer/readline"
)

// Validator is a function passed to a Question in order to redefine
type Validator func(string) error

// Question is the core data structure for a survey questionnaire.
type Question struct {
	Name     string
	Prompt   Prompt
	Validate Validator
}

// Prompt is the primary interface for the objects that can take user input
// and return a string value.
type Prompt interface {
	Prompt(*readline.Instance) (string, error)
	Cleanup(*readline.Instance, string) error
}

var ErrorTemplate = `{{color "red"}}✘ Sorry, your reply was invalid: {{.Error}}{{color "reset"}}
`

// AskOne asks a single question without performing validation on the answer.
func AskOne(p Prompt) (string, error) {
	answers, err := Ask([]*Question{{Name: "q1", Prompt: p}})
	if err != nil {
		return "", err
	}
	return answers["q1"], nil
}

// AskOneValidate asks a single question and validates the answer with v.
func AskOneValidate(p Prompt, v Validator) (string, error) {
	answers, err := Ask([]*Question{{Name: "q1", Prompt: p, Validate: v}})
	return answers["q1"], err
}

func handleError(err error) {
	// tell the user what happened
	fmt.Println(err.Error())
	// quit the survey
	os.Exit(1)
}

// Ask performs the prompt loop
func Ask(qs []*Question) (map[string]string, error) {
	// grab the readline instance
	rl, err := terminal.GetReadline()
	if err != nil {
		handleError(err)
	}

	// the response map
	res := make(map[string]string)
	// go over every question
	for _, q := range qs {
		// grab the user input and save it
		ans, err := q.Prompt.Prompt(rl)
		// if there was a problem
		if err != nil {
			handleError(err)
		}

		// if there is a validate handler for this question
		if q.Validate != nil {
			// wait for a valid response
			for invalid := q.Validate(ans); invalid != nil; invalid = q.Validate(ans) {
				out, err := RunTemplate(ErrorTemplate, invalid)
				if err != nil {
					return nil, err
				}
				// send the message to the user
				fmt.Print(out)
				// ask for more input
				ans, err = q.Prompt.Prompt(rl)
				// if there was a problem
				if err != nil {
					handleError(err)
				}
			}
		}

		// tell the prompt to cleanup with the validated value
		q.Prompt.Cleanup(rl, ans)

		// if something went wrong
		if err != nil {
			// stop listening
			return nil, err
		}
		// add it to the map
		res[q.Name] = ans
	}
	// return the response
	return res, nil
}
