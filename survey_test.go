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
		ErrorTemplate,
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
