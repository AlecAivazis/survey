package survey

import (
	"errors"

	"github.com/AlecAivazis/survey/core"
)

// PageSize is the default maximum number of items to show in select/multiselect prompts
var PageSize = 7

// Validator is a function passed to a Question in order to redefine
type Validator func(interface{}) error

// Converter is a function passed to a Question in order to convert value types
type Converter func(interface{}) (interface{}, error)

// Question is the core data structure for a survey questionnaire.
type Question struct {
	Name     string
	Prompt   Prompt
	Validate Validator
	Convert  Converter
}

// Prompt is the primary interface for the objects that can take user input
// and return a string value.
type Prompt interface {
	Prompt() (interface{}, error)
	Cleanup(interface{}) error
	Error(error) error
}

// AskOne asks a single question without performing validation on the answer.
func AskOne(p Prompt, t interface{}, v Validator, c Converter) error {
	err := Ask([]*Question{{Prompt: p, Validate: v, Convert: c}}, t)
	if err != nil {
		return err
	}

	return nil
}

// Ask performs the prompt loop
func Ask(qs []*Question, t interface{}) error {

	// if we weren't passed a place to record the answers
	if t == nil {
		// we can't go any further
		return errors.New("cannot call Ask() with a nil reference to record the answers")
	}

	// go over every question
	for _, q := range qs {
		// grab the user input and save it
		ans, err := q.Prompt.Prompt()
		convertedAns := ans
		// if there was a problem
		if err != nil {
			return err
		}

		// if there's a converter
		if q.Convert != nil {
			var invalid error

			// wait for a valid response
			for convertedAns, invalid = q.Convert(ans); invalid != nil; convertedAns, invalid = q.Convert(ans) {
				err := q.Prompt.Error(invalid)
				// if there was a problem
				if err != nil {
					return err
				}

				// ask for more input
				ans, err = q.Prompt.Prompt()
				// if there was a problem
				if err != nil {
					return err
				}
			}
		}

		// if there is a validate handler for this question
		if q.Validate != nil {
			// wait for a valid response
			for invalid := q.Validate(convertedAns); invalid != nil; invalid = q.Validate(convertedAns) {
				err := q.Prompt.Error(invalid)
				// if there was a problem
				if err != nil {
					return err
				}

				// ask for more input
				ans, err = q.Prompt.Prompt()
				// if there was a problem
				if err != nil {
					return err
				}
			}
		}

		// tell the prompt to cleanup with the validated value
		q.Prompt.Cleanup(ans)

		// if something went wrong
		if err != nil {
			// stop listening
			return err
		}

		// add it to the map
		err = core.WriteAnswer(t, q.Name, convertedAns)
		// if something went wrong
		if err != nil {
			return err
		}

	}
	// return the response
	return nil
}

// paginate returns a single page of choices given the page size, the total list of
// possible choices, and the current selected index in the total list.
func paginate(page int, choices []string, sel int) ([]string, int) {
	// the number of elements to show in a single page
	var pageSize int
	// if the select has a specific page size
	if page != 0 {
		// use the specified one
		pageSize = page
		// otherwise the select does not have a page size
	} else {
		// use the package default
		pageSize = PageSize
	}

	var start, end, cursor int

	if len(choices) < pageSize {
		// if we dont have enough options to fill a page
		start = 0
		end = len(choices)
		cursor = sel

	} else if sel < pageSize/2 {
		// if we are in the first half page
		start = 0
		end = pageSize
		cursor = sel

	} else if len(choices)-sel-1 < pageSize/2 {
		// if we are in the last half page
		start = len(choices) - pageSize
		end = len(choices)
		cursor = sel - start

	} else {
		// somewhere in the middle
		above := pageSize / 2
		below := pageSize - above

		cursor = pageSize / 2
		start = sel - above
		end = sel + below
	}

	// return the subset we care about and the index
	return choices[start:end], cursor
}
