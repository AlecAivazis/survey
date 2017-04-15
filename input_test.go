package survey

import (
	"testing"

	"github.com/alecaivazis/survey/core"
)

func init() {
	// disable color output for all prompts to simplify testing
	core.DisableColor = true
}

func TestInputFormatQuestion(t *testing.T) {

	prompt := &Input{
		Message: "What is your favorite month:",
		Default: "April",
	}

	actual, err := core.RunTemplate(
		InputQuestionTemplate,
		InputTemplateData{Input: *prompt},
	)
	if err != nil {
		t.Errorf("Failed to run template to format input question: %s", err)
	}

	expected := `? What is your favorite month: (April) `

	if actual != expected {
		t.Errorf("Formatted input question was not formatted correctly. Found:\n%s\nExpected:\n%s", actual, expected)
	}
}

func TestInputFormatAnswer(t *testing.T) {

	prompt := &Input{
		Message: "What is your favorite month:",
		Default: "April",
	}

	actual, err := core.RunTemplate(
		InputQuestionTemplate,
		InputTemplateData{Input: *prompt, Answer: "October"},
	)
	if err != nil {
		t.Errorf("Failed to run template to format input answer: %s", err)
	}

	expected := `? What is your favorite month: October`

	if actual != expected {
		t.Errorf("Formatted input answer was not formatted correctly. Found:\n%s\nExpected:\n%s", actual, expected)
	}
}
