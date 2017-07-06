package survey

import (
	"fmt"
	"testing"

	"github.com/AlecAivazis/survey/core"
)

func init() {
	// disable color output for all prompts to simplify testing
	core.DisableColor = true
}

func TestValidationError(t *testing.T) {

	err := fmt.Errorf("Football is not a valid month")

	actual, err := core.RunTemplate(
		core.ErrorTemplate,
		err,
	)
	if err != nil {
		t.Errorf("Failed to run template to format error: %s", err)
	}

	expected := `✘ Sorry, your reply was invalid: Football is not a valid month
`

	if actual != expected {
		t.Errorf("Formatted error was not formatted correctly. Found:\n%s\nExpected:\n%s", actual, expected)
	}
}

func TestAsk_returnsErrorIfTargetIsNil(t *testing.T) {
	// pass an empty place to leave the answers
	err := Ask([]*Question{}, nil)

	// if we didn't get an error
	if err == nil {
		// the test failed
		t.Error("Did not encounter error when asking with no where to record.")
	}
}

func TestPagination_tooFew(t *testing.T) {
	// a small list of options
	choices := []string{"choice1", "choice2", "choice3"}
	// a page bigger than the total number
	pageSize := 4
	// the current selection
	sel := 2

	// compute the page info
	choices, idx := paginate(pageSize, choices, sel)

	if choices[0] != "choice1" && choices[1] != "choice2" && choices[2] != "choice3" {
		t.Error("Did not recieve the right page for the first half")
	}

	// with the second index highlighted
	if idx != 2 {
		t.Error("Did not recieve the correct index")
	}

}

func TestPagination_firstHalf(t *testing.T) {
	// the choices for the test
	choices := []string{"choice1", "choice2", "choice3", "choice4", "choice5", "choice6"}

	// section the choices into groups of 4 so the choice is somewhere in the middle
	// to verify there is no displacement of the page
	pageSize := 4

	// test the second item
	sel := 2

	// compute the page info
	choices, idx := paginate(pageSize, choices, sel)

	// we should see the first three options
	if choices[0] != "choice1" && choices[1] != "choice2" && choices[2] != "choice3" && choices[3] != "choice3" {
		t.Error("Did not recieve the right page for the first half")
	}

	// with the second index highlighted
	if idx != 2 {
		t.Error("Did not recieve the correct index")
	}
}

func TestPagination_middle(t *testing.T) {
	// the choices for the test
	choices := []string{"choice1", "choice2", "choice3", "choice4", "choice5", "choice6"}

	// section the choices into groups of 3
	pageSize := 2

	// test the second item
	sel := 3

	// compute the page info
	choices, idx := paginate(pageSize, choices, sel)

	// we should see the first three options
	if choices[0] != "choice3" && choices[1] != "choice4" {
		t.Error("Did not recieve the right page for the middle half")
	}

	// with the second index highlighted
	if idx != 1 {
		t.Error("Did not recieve the correct index")
	}
}

func TestPagination_lastHalf(t *testing.T) {
	// the choices for the test
	choices := []string{"choice1", "choice2", "choice3", "choice4", "choice5", "choice6"}

	// section the choices into groups of 3
	pageSize := 3

	// test the second item
	sel := 5

	// compute the page info
	choices, idx := paginate(pageSize, choices, sel)

	// we should see the first three options
	if choices[0] != "choice4" && choices[1] != "choice5" && choices[2] != "choice6" {
		t.Error("Did not recieve the right page for the last half")
	}

	// with the second index highlighted
	if idx != 2 {
		t.Error("Did not recieve the correct index")
	}
}
