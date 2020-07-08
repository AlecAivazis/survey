package survey

import (
	"fmt"
	"testing"

	"github.com/AlecAivazis/survey/v2/core"
)

func TestValidationError(t *testing.T) {

	err := fmt.Errorf("Football is not a valid month")

	actual, _, err := core.RunTemplate(
		ErrorTemplate,
		&ErrorTemplateData{
			Error: err,
			Icon:  defaultIcons().Error,
		},
	)
	if err != nil {
		t.Errorf("Failed to run template to format error: %s", err)
	}

	expected := fmt.Sprintf("%s Sorry, your reply was invalid: Football is not a valid month\n", defaultIcons().Error.Text)

	if actual != expected {
		t.Errorf("Formatted error was not formatted correctly. Found:\n%s\nExpected:\n%s", actual, expected)
	}
}
