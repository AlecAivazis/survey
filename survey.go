package survey

import (
	"fmt"
	// tm "github.com/buger/goterm"
	"github.com/alecaivazis/survey/format"
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
	Prompt() (string, error)
	Cleanup(string) error
}

// AskOne asks a single question without performing validation on the answer.
// If an error occurs, an empty string is returned.
func AskOne(p Prompt) string {
	answers, err := Ask([]*Question{{Name: "q1", Prompt: p}})
	if err != nil {
		return ""
	}
	return answers["q1"]
}

// AskOneValidate asks a single question and validates the answer with v.
func AskOneValidate(p Prompt, v Validator) (string, error) {
	answers, err := Ask([]*Question{{Name: "q1", Prompt: p, Validate: v}})
	return answers["q1"], err
}

// Ask performs the prompt loop
func Ask(qs []*Question) (map[string]string, error) {
	// the response map
	res := make(map[string]string)
	// go over every question
	for _, q := range qs {
		// grab the user input and save it
		ans, err := q.Prompt.Prompt()

		// if there is a validate handler for this question
		if q.Validate != nil {
			// wait for a valid response
			for invalid := q.Validate(ans); invalid != nil; invalid = q.Validate(ans) {
				// the error message
				msg := "Sorry, your reply was invalid: "
				// send the message to the user
				fmt.Print(format.ErrorColor, format.Error, msg, invalid.Error(), format.ResetFormat, "\n")
				// ask for more input
				ans, err = q.Prompt.Prompt()
			}
		}

		// tell the prompt to cleanup with the validated value
		q.Prompt.Cleanup(ans)

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
