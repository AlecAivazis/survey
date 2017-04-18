package survey

import (
	"testing"

	"github.com/AlecAivazis/survey/core"
)

func init() {
	// disable color output for all prompts to simplify testing
	core.DisableColor = true
}

func TestPasswordFormatQuestion(t *testing.T) {

	prompt := &Input{
		Message: "Tell me your secret:",
	}

	actual, err := core.RunTemplate(
		PasswordQuestionTemplate,
		*prompt,
	)
	if err != nil {
		t.Errorf("Failed to run template to format password question: %s", err)
	}

	expected := `? Tell me your secret: `

	if actual != expected {
		t.Errorf("Formatted input question was not formatted correctly. Found:\n%s\nExpected:\n%s", actual, expected)
	}
}
