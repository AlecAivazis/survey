package survey

import (
	"testing"
)

func init() {
	// disable color output for all prompts to simplify testing
	DisableColor = true
}

func TestPasswordFormatQuestion(t *testing.T) {

	prompt := &Input{
		Message: "Tell me your secret:",
	}

	actual, err := RunTemplate(
		passwordQuestionTemplate,
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
