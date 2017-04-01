package survey

import (
	"fmt"
	"testing"
)

func init() {
	// disable color output for all prompts to simplify testing
	DisableColor = true
}

func TestValidationError(t *testing.T) {

	err := fmt.Errorf("Football is not a valid month")

	actual, err := runTemplate(
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
