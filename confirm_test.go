package survey

import (
	"testing"

	"github.com/AlecAivazis/survey/core"
)

func init() {
	// disable color output for all prompts to simplify testing
	core.DisableColor = true
}

func TestConfirmFormatQuestion(t *testing.T) {

	prompt := &Confirm{
		Message: "Is pizza your favorite food?",
		Default: true,
	}

	actual, err := core.RunTemplate(
		ConfirmQuestionTemplate,
		ConfirmTemplateData{Confirm: *prompt},
	)
	if err != nil {
		t.Errorf("Failed to run template to format input question: %s", err)
	}

	expected := `? Is pizza your favorite food? (Y/n) `

	if actual != expected {
		t.Errorf("Formatted input question was not formatted correctly. Found:\n%s\nExpected:\n%s", actual, expected)
	}
}

func TestConfirmFormatQuestionDefaultFalse(t *testing.T) {

	prompt := &Confirm{
		Message: "Is pizza your favorite food?",
		Default: false,
	}

	actual, err := core.RunTemplate(
		ConfirmQuestionTemplate,
		ConfirmTemplateData{Confirm: *prompt},
	)
	if err != nil {
		t.Errorf("Failed to run template to format input answer: %s", err)
	}

	expected := `? Is pizza your favorite food? (y/N) `

	if actual != expected {
		t.Errorf("Formatted input answer was not formatted correctly. Found:\n%s\nExpected:\n%s", actual, expected)
	}
}

func TestConfirmFormatAnswer(t *testing.T) {

	// default false
	prompt := &Confirm{
		Message: "Is pizza your favorite food?",
	}

	actual, err := core.RunTemplate(
		ConfirmQuestionTemplate,
		ConfirmTemplateData{Confirm: *prompt, Answer: "Yes"},
	)
	if err != nil {
		t.Errorf("Failed to run template to format input answer: %s", err)
	}

	expected := `? Is pizza your favorite food? Yes`

	if actual != expected {
		t.Errorf("Formatted input answer was not formatted correctly. Found:\n%s\nExpected:\n%s", actual, expected)
	}
}
