package probe

import (
	"fmt"
	// tm "github.com/buger/goterm"
	"github.com/alecaivazis/probe/format"
)

// Validator is a function passed to a Question in order to redefine
type Validator func(string) error

// Question is the core data structure for a probe questionnaire.
type Question struct {
	Name     string
	Prompt   Prompt
	Validate Validator
}

// Prompt is the primary interface for the objects that can take user input
// and return a string value.
type Prompt interface {
	Prompt() (string, error)
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
				fmt.Println(format.Error, " ", msg, invalid.Error(), format.ResetFormat)
				// ask for more input
				ans, err = q.Prompt.Prompt()
			}
		}

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
